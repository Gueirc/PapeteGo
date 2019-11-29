package main

import "log"
import "strings"

type carta string

const (
	nada    carta = ""
	papel   carta = "/papel"
	pedra   carta = "/pedra"
	tesoura carta = "/tesoura"
)

func salaJogos(jgdEntrou chan *Jogador) {

	var jgdrEsperando *Jogador
	for {
		jgdrEsperando = <-jgdEntrou
		pronto := false
		for {
			if pronto {
				break
			}
			select {
			case nJgdr := <-jgdEntrou:
				go oJogo(nJgdr, jgdrEsperando)
				pronto = true
			case msg := <-jgdrEsperando.chatJogo:
				if msg == "/quitei" {
					pronto = true
				}
			}
		}
	}
}

func sairJogo(jgd1, jgd2 *Jogador) {
	jgd1.emJogo = false
	jgd2.emJogo = false
}
func ctos(c carta) string {
	if c == papel {
		return "/papel"
	} else if c == pedra {
		return "/pedra"
	} else if c == tesoura {
		return "/tesoura"
	} else {
		return ""
	}
}
func stoc(s string) carta {
	if s == "/papel" {
		return papel
	} else if s == "/pedra" {
		return pedra
	} else if s == "/tesoura" {
		return tesoura
	} else {
		return nada
	}
}

func oJogo(jgd1, jgd2 *Jogador) {
	defer sairJogo(jgd1, jgd2)
	jgd1.emJogo = true
	jgd2.emJogo = true
	go sendChat(jgd1, "/serverJogo", []string{"/oponente", jgd2.nome})
	go sendChat(jgd2, "/serverJogo", []string{"/oponente", jgd1.nome})
	var jgd1Escolha carta
	var jgd2Escolha carta
	var jgd1Pronto bool
	var jgd2Pronto bool
	for {
		if jgd1Pronto && jgd2Pronto {
			cGanhou, err := jgoResult(jgd1Escolha, jgd2Escolha)
			log.Println(ctos(cGanhou))
			log.Println(ctos(jgd1Escolha))
			log.Println(ctos(jgd2Escolha))
			if err != 0 {
				log.Panic("Carta invalida") // pega isso devolta la no main
			}
			if cGanhou == jgd1Escolha { // fazer um sendResult pra isso
				msge1 := []string{"/resultado", ctos(jgd1Escolha), ctos(jgd2Escolha), "/ganhou"}
				go sendChat(jgd1, "/serverJogo", msge1)
				msge2 := []string{"/resultado", ctos(jgd2Escolha), ctos(jgd1Escolha), "/perdeu"}
				go sendChat(jgd2, "/serverJogo", msge2)
				// salvar no arquivo, estrelas
			} else if cGanhou == jgd2Escolha {
				msge1 := []string{"/resultado", ctos(jgd1Escolha), ctos(jgd2Escolha), "/perdeu"}
				go sendChat(jgd1, "/serverJogo", msge1)
				msge2 := []string{"/resultado", ctos(jgd2Escolha), ctos(jgd1Escolha), "/ganhou"}
				go sendChat(jgd2, "/serverJogo", msge2)

			} else {
				msge1 := []string{"/resultado", ctos(jgd1Escolha), ctos(jgd2Escolha), "/empatou"}
				go sendChat(jgd1, "/serverJogo", msge1)
				msge2 := []string{"/resultado", ctos(jgd2Escolha), ctos(jgd1Escolha), "/empatou"}
				go sendChat(jgd2, "/serverJogo", msge2)

			}

			return
		} else {
			select {
			case msg := <-jgd1.chatJogo:
				texto := strings.Split(msg, " ")
				if texto[0] == "/jogar" {
					jgd1Escolha = stoc(texto[1])
					if jgd1Escolha == nada {
						log.Panic("Suposto receber uma carta")
					}
					jgd1Pronto = true
					sendChat(jgd1, "/serverJogo", []string{"/escolhido"})
					sendChat(jgd2, "/serverJogo", []string{"/oescolhido"})
				} else if texto[0] == "/falar" {
					sendChat(jgd1, jgd1.nome, texto[1:len(texto)])
					sendChat(jgd2, jgd1.nome, texto[1:len(texto)])
				} else if texto[0] == "/quitei" {
					go sendChat(jgd2, "/serverJogo", []string{"/quitou", jgd1.nome})
					return
				}
			case msg := <-jgd2.chatJogo:
				texto := strings.Split(msg, " ")
				if texto[0] == "/jogar" {
					log.Println("opa")
					jgd2Escolha = stoc(texto[1])
					if jgd2Escolha == nada {
						log.Panic("Suposto receber uma carta")
					}
					jgd2Pronto = true
					sendChat(jgd2, "/serverJogo", []string{"/escolhido"})
					sendChat(jgd1, "/serverJogo", []string{"/oescolhido"})
				} else if texto[0] == "/falar" {
					go sendChat(jgd2, jgd2.nome, texto[1:len(texto)])
					go sendChat(jgd1, jgd2.nome, texto[1:len(texto)])
				} else if texto[0] == "/quitei" {
					go sendChat(jgd1, "/serverJogo", []string{"/quitou" + jgd2.nome})
					return
				}

			}
		}
	}
}
func jgoResult(prim, seg carta) (carta, int) {
	if prim == papel {
		if seg == papel {
			return pedra, 0
		}
		if seg == pedra {
			return papel, 0
		}
		if seg == tesoura {
			return tesoura, 0
		} else {
			return papel, -1
		}

	} else if prim == pedra {
		if seg == papel {
			return papel, 0
		}
		if seg == pedra {
			return papel, 0
		}
		if seg == tesoura {
			return pedra, 0
		} else {
			return papel, -1
		}
	} else if prim == tesoura {
		if seg == papel {
			return tesoura, 0
		}
		if seg == pedra {
			return pedra, 0
		}
		if seg == tesoura {
			return papel, 0
		} else {
			return papel, -1
		}
	} else {
		return papel, -1
	}
}
