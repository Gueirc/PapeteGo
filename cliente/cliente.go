package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
)

var chJogo = make(chan string)
var conn net.Conn

func defFunc(w *window.Window) {
	// crarr função dump para imprimir dados
	w.DefineFunction("dump", func(args ...*sciter.Value) *sciter.Value {
		for _, v := range args {
			fmt.Print(v.String() + " ")
		}
		fmt.Println()
		return sciter.NullValue()
	})
	//Função pra pegar coisas
	w.DefineFunction("toGo", func(args ...*sciter.Value) *sciter.Value {
		switch args[0].String() {
		case "login":
			defer func() {
			}()
			entrada := args[1].String() // formulário
			fncErr := args[2]           // função callback erro
			fncCon := args[3]           // função callback fomos conectados
			//fncLobby := args[3]
			//fncLogin := args[4]
			if r := recover(); r != nil {
				fmt.Println("Recuperado: ", r) // se a conexão der errado
				fncErr.Invoke(sciter.NullValue(), "[Native Script]",
					sciter.NewValue("Não foi possível conectar-se ao servidor."))
			}
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
		}
		return sciter.NullValue()
	})

}

func main() {
	// Janela
	w, err := window.New(sciter.DefaultWindowCreateFlag, &sciter.Rect{
		0, 0, 0, 0})
	if err != nil {
		log.Fatal("Criar janela erro: ", err)
	}
	go aoJogo(chJogo)
	w.LoadFile("/home/cire/PapeteGo/src/PapeteGo/cliente/main.htm")
	w.SetTitle("Login")
	defFunc(w)
	w.Show()
	w.Run()
}

func ligarNoServer(ip string, nick string) {
	endr := ip + ":4243"
	conn, err := net.Dial("tcp", endr)
	if err != nil {
		//falar pro usuario que nao deu
		log.Panic("Erro na conexão: ", err) //vai pro recover lá

	}
	msg := "/nick " + nick + "\n"
	conn.Write([]byte(msg))

	log.Println("Conectado")
}
func aoJogo(canal chan string) {
	for {
		pedido := <-canal
		if pedido == "aoJogo" {
			log.Println("Clico em jogar")
			sendJogar()
		}

	}

}

func sendJogar() {

}
