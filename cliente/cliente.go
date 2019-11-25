package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
)

func defFunc(w *window.Window) {
	// crarr função dump para imprimir dados
	w.DefineFunction("dump", func(args ...*sciter.Value) *sciter.Value {
		for _, v := range args {
			fmt.Print(v.String() + " ")
		}
		fmt.Println()
		return sciter.NullValue()
	})
	//Função reg para registar, lógica disso
	w.DefineFunction("entrar", func(args ...*sciter.Value) *sciter.Value {
		var entrada string
		for _, v := range args {
			//fmt.Print(v.String() + " ")
			entrada = v.String()
		}
		if strings.Contains(entrada, "/") != true {
			dados := strings.Split(entrada, ",")      // separa entre nick e ip
			nick := strings.Split(dados[0], ":")[1]   // pega o nick
			nick = strings.ReplaceAll(nick, "\"", "") // tira o "
			ip := strings.Split(dados[1], ":")[1]
			ip = strings.TrimSuffix(ip, "}")
			ip = strings.ReplaceAll(ip, "\"", "")
			go conectServer(ip, nick)
		}

		return sciter.NullValue()
	})
}

func main() {
	// Janela
	w, err := window.New(sciter.DefaultWindowCreateFlag, &sciter.Rect{
		0, 0, 300, 120})
	if err != nil {
		log.Fatal(err)
	}
	w.LoadFile("ui.html")
	w.SetTitle("Papete Login")
	defFunc(w)
	w.Show()
	w.Run()
}

func conectServer(ip, nick string) {
	endr := ip + ":4243"
	conn, err := net.Dial("tcp", endr)
	if err != nil {
		//falar pro usuario que nao deu
		log.Panic("Erro na conexão")
	}
	defer conn.Close()
	msg := "/nick " + nick + "\n"
	conn.Write([]byte(msg))
	w, err := window.New(
		sciter.DefaultWindowCreateFlag, sciter.DefaultRect)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Vamo la")
	w.LoadFile("lobby.html")
	w.SetTitle("Papete")
	w.Show()
	w.Run()
	for {
	}

}
