package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
)

var chJogo = make(chan string)
var fncBack = make(chan *sciter.Value)
var conn net.Conn

func defFunc(w *window.Window) {
	// crarr função dump para imprimir dados
	//w.DefineFunction("dump", func(args ...*sciter.Value) *sciter.Value {
	//	for _, v := range args {
	//		fmt.Print(v.String() + " ")
	//	}
	//	fmt.Println()
	//	return sciter.NullValue()
	//})
	//Função pra pegar coisas
	w.DefineFunction("toGo", func(args ...*sciter.Value) *sciter.Value {
		//fncLobby := args[3]
		//fncLogin := args[4]
		log.Println(args[0].String())
		switch args[0].String() {
		case "login":
			entrada := args[1].String() // formulário
			fncErr := args[2]           // função callback erro
			fncCon := args[3]           // função callback fomos conectados
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("Recuperado: ", r) // se a conexão der errado
					fncErr.Invoke(sciter.NullValue(), "[Native Script]",
						sciter.NewValue("Não foi possível conectar-se ao servidor, verifique o ip do servidor."))
				}
			}()
			if strings.Contains(entrada, "/") != true {
				dados := strings.Split(entrada, ",")      // separa entre nick e ip
				nick := strings.Split(dados[0], ":")[1]   // pega o nick
				nick = strings.ReplaceAll(nick, "\"", "") // tira o "
				ip := strings.Split(dados[1], ":")[1]
				ip = strings.TrimSuffix(ip, "}")
				ip = strings.ReplaceAll(ip, "\"", "")

				ligarNoServer(ip, nick) // panico lá em cima
				fncCon.Invoke(sciter.NullValue(), "[Native Script]")

			} else {
				fncErr.Invoke(sciter.NullValue(), "[Native Script]",
					sciter.NewValue("Digite um nome de usuário válido"))
			}

		case "aoJogo":
			chJogo <- "aoJogo"
		case "papel":
			sendPapel()
		case "pedra":
			sendPedra()
		case "tesoura":
			sendTesoura()
		case "fechar":
			sendSair()
		}
		return sciter.NullValue()
	})

}

func main() {
	// Janela
	// mudar para que não seja resizeable
	w, err := window.New(sciter.DefaultWindowCreateFlag, &sciter.Rect{
		0, 0, 0, 0})
	if err != nil {
		log.Fatal("Criar janela erro: ", err)
	}
	w.LoadFile("/home/cire/PapeteGo/src/github.com/Gueirc/PapeteGo/cliente/main.htm")
	w.SetTitle("Login")
	defFunc(w)
	go handlerCon(chJogo, w)
	w.Show()
	w.Run()
}

func ligarNoServer(ip string, nick string) {
	endr := ip + ":4243"
	var err error
	conn, err = net.Dial("tcp", endr)
	if err != nil {
		//falar pro usuario que nao deu
		log.Panic("Erro na conexão: ", err) //vai pro recover lá

	}
	msg := "/nick " + nick + "\n"
	conn.Write([]byte(msg))

	msgr, _ := bufio.NewReader(conn).ReadString('\n')
	log.Println(msgr)
}
func handlerCon(canal chan string, w *window.Window) {
	root, err := w.GetRootElement()
	if err != nil {
		log.Fatal("get root element failed: ", err.Error())
	}
	for {
		voltou := false
		pedido := <-canal
		if pedido == "aoJogo" {
			log.Println("Clico em jogar")
			sendJogar()
			for {
				if voltou == true {
					break
				}
				log.Println("Esperando")
				msg, _ := bufio.NewReader(conn).ReadString('\n')
				textoDiv := strings.Split(msg, " ")
				for i := 0; i < len(textoDiv); i++ {
					textoDiv[i] = strings.TrimSpace(textoDiv[i])
				}
				log.Println(msg)
				switch textoDiv[1] {
				case "/server":
					//Ignora o /chat e o enviador
					texto := strings.Join(textoDiv[2:], " ")
					uiMens(texto, w)
				case "/serverJogo":
					switch textoDiv[2] {
					case "/oponente":
						oponente := textoDiv[3] // nome do oponente
						log.Println("	" + oponente)
						msg := "Seu oponente é: " + oponente
						retorno, err := root.CallMethod("startJogo",
							sciter.NewValue(msg))
						if err != nil {
							log.Println("method call startJogo  failed,", err)
						} else {
							log.Println("method call startJogo successfulyy ", retorno)
						}
					case "/quitou":
						retorno, err := root.CallMethod("voltarLobby",
							sciter.NewValue("Seu oponente saiu."))
						if err != nil {
							log.Panic("method call voltarLobby failed,", err)
						} else {
							log.Println("method call voltarLobby successfulyy ", retorno)
							voltou = true
						}
					case "/escolhido":
						log.Println("passou pra escolhido")
						uiMens("Você escolheu...", w)
					case "/oescolhido":
						log.Println("passou pra oescolhido")
						uiMens("Seu oponente escolheu...", w)
					case "/resultado":
						log.Println(textoDiv)
						log.Println(textoDiv[3])
						log.Println(textoDiv[4])
						msg := "O jogo acabou. Você jogou " + ctos(textoDiv[3]) + " e seu oponente jogou " + ctos(textoDiv[4]) + "."
						log.Println(textoDiv[5])
						if textoDiv[5] == "/ganhou" {
							msg = msg + " Você ganhou"
						} else if textoDiv[5] == "/perdeu" {
							msg = msg + " Você perdeu"
						} else if textoDiv[5] == "/empatou" {
							msg = msg + " Você empatou"
						}
						retorno, err := root.CallMethod("voltarLobby",
							sciter.NewValue(msg))
						if err != nil {
							log.Panic("method call voltarLobby failed,", err)
						} else {
							log.Println("method call voltarLobby successfulyy ", retorno)
							voltou = true
						}
					}

				}
			}
		}

	}

}

func ctos(c string) string {
	if c == "/papel" {
		return "papel"
	} else if c == "/pedra" {
		return "pedra"
	} else if c == "/tesoura" {
		return "tesoura"
	} else {
		return "NAO"
	}

}

func uiMens(mens string, w *window.Window) {
	root, err := w.GetRootElement()
	if err != nil {
		log.Fatal("get root element failed: ", err.Error())
	}
	for {
		//time.Sleep(1 * time.Millisecond)
		retorno, err := root.CallMethod("mensagem",
			sciter.NewValue(mens))
		if err != nil {
			log.Println("method call startJogo  failed,", err)
		} else {
			log.Println("method call startJogo successfulyy ", retorno)
			break
		}
	}

}

func sendJogar() {
	conn.Write([]byte("/jogar X \n"))
}
func sendSair() {
	conn.Close()
}
func sendPapel() {
	conn.Write([]byte("/jogar /papel X \n"))
}
func sendPedra() {
	conn.Write([]byte("/jogar /pedra X \n"))
}
func sendTesoura() {
	conn.Write([]byte("/jogar /tesoura X \n"))
}
