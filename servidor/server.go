package main

import (
	"github.com/firstrow/tcp_server"
	"log"
	"strings"
)

type Jogador struct {
	nome     string
	temNome  bool
	emJogo   bool
	chatJogo chan string
	tcp_conn *tcp_server.Client
}

//enviar para o chat de jgdr a mensagem tex de enviador
func sendChat(jgdr *Jogador, enviador string, tex []string) {
	enviar := ""
	for i := 0; i < len(tex); i++ {
		enviar = enviar + " " + tex[i]
	}
	jgdr.tcp_conn.Send("/chat " + enviador + enviar + "\n")
}

func main() {

	jgdsConectados := make(map[*tcp_server.Client]*Jogador)

	paraSalaJogos := make(chan *Jogador)
	go salaJogos(paraSalaJogos)

	server := tcp_server.New("localhost:4243")

	server.OnNewClient(func(c *tcp_server.Client) {
		// Cliente conectado
		c.Send("Conectado\n")
		log.Println("Novo cliente conectado")
		nJgd := &Jogador{"", false, false, make(chan string, 1), c}
		jgdsConectados[c] = nJgd

	})

	server.OnNewMessage(func(c *tcp_server.Client, message string) {
		log.Println("Recebido mensagem: " + message)
		texto := strings.Split(message, " ")
		jgdrAtual := jgdsConectados[c] // Jogador que falou
		//recebe nick, a primeira coisa, na tela de login
		if jgdrAtual.temNome != true {
			if texto[0] == "/nick" {
				// tira o /n do final do nome
				texto[1] = strings.TrimSuffix(texto[1], "\n")
				jgdrAtual.nome = texto[1]
				jgdrAtual.temNome = true
			} else {
				log.Println("Assumido que receba nick, fechando conexão")
				c.Close()
			}
		} else if jgdrAtual.emJogo {
			jgdrAtual.chatJogo <- message // manda pro chat do jogo
		} else if texto[0] == "/falar" { // se tiver falado pra todo mundo
			//tratar fala indevida?
			for client, jgdr := range jgdsConectados {
				if jgdr.emJogo != true { // clientes nao podem conversar com
					// clientes dentro do jogo
					if client == c {
						// mandar coisas de depois do /falar até o final
						// jgdr == jgdrAtual
						go sendChat(jgdr, jgdrAtual.nome, texto[1:len(texto)])
					} else {
						// mandar pras outras pessoas
						go sendChat(jgdr, jgdrAtual.nome, texto[1:len(texto)])
					}
				}
			}
		} else if texto[0] == "/jogar" {
			log.Println("Recebido pedido de jogo")
			msg := []string{"Esperando oponente, aguarde."}
			sendChat(jgdrAtual, "/server", msg)
			paraSalaJogos <- jgdrAtual
		}
	})
	server.OnClientConnectionClosed(func(c *tcp_server.Client, err error) {
		log.Println("Cliente Desconectado", err)
		jgdr := jgdsConectados[c]
		jgdr.chatJogo <- "/quitei"
		delete(jgdsConectados, c)

	})
	server.Listen()

}
