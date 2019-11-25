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
	var emEspera int = 0
	for {
		nJgdr := <-jgdEntrou

		if emEspera > 0 {
			go oJogo(nJgdr, jgdrEsperando)
			emEspera--
		} else {
			jgdrEsperando = nJgdr
			emEspera++
		}
	}
}

func sairJogo(jgd1, jgd2 *Jogador) {
	jgd1.emJogo = false
	jgd2.emJogo = false
}
func ctos(c carta) string {
	if c == "/papel" {
		return "/papel"
	} else if c == "/pedra" {
		return "/papel"
	} else if c == "/tesoura" {
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
	go sendChat(jgd1, "/serverJogo", []string{"/oponente ", jgd2.nome})
	go sendChat(jgd2, "/serverJogo", []string{"/oponente ", jgd1.nome})
	var jgd1Escolha carta
	var jgd2Escolha carta
	var jgd1Pronto bool
	var jgd2Pronto bool
	for {
		if jgd1Pronto && jgd2Pronto {
			cGanhou, err := jgoResult(jgd1Escolha, jgd2Escolha)
			if err != 0 {
				log.Panic("Carta invalida") // pega isso devolta la no main
			}
			if cGanhou == jgd1Escolha {
				// salvar no arquivo, estrelas
			} else if cGanhou == jgd2Escolha {

			} else {

			}

			msge := []string{"/resultado ", ctos(jgd2Escolha) + " ", jgd2.nome}
			go sendChat(jgd1, "/serverJogo", msge)
			msge = []string{"/resultado ", ctos(jgd1Escolha) + " ", jgd1.nome}
			go sendChat(jgd2, "/serverJogo", msge)
			break
		} else {
			select {
			case msg := <-jgd1.chatJogo:
				texto := strings.Split(msg, " ")
				if texto[0] != "/falar" {
					jgd1Escolha = stoc(texto[0])
					if jgd1Escolha == nada {
						log.Panic("Suposto receber uma carta")
					}
					jgd1Pronto = true
					go sendChat(jgd1, "/serverJogo", []string{"/escolhido"})
				} else {
					go sendChat(jgd1, jgd1.nome, texto[1:len(texto)])
					go sendChat(jgd2, jgd1.nome, texto[1:len(texto)])
				}
			case msg := <-jgd2.chatJogo:
				texto := strings.Split(msg, " ")
				if texto[0] != "/falar" {
					jgd2Escolha = stoc(texto[0])
					if jgd2Escolha == nada {
						log.Panic("Suposto receber uma carta")
					}
					jgd2Pronto = true
					go sendChat(jgd2, "/serverJogo", []string{"/escolhido"})
				} else {
					go sendChat(jgd1, jgd2.nome, texto[1:len(texto)])
					go sendChat(jgd2, jgd2.nome, texto[1:len(texto)])

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
