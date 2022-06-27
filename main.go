package main

import (
	"fmt"
)

var PATH_DOADORES = "data/doadores.txt"
var PATH_PENDENTES = "data/pendentes.txt"

func main() {

	var menu int

	menu = -1
	for {
		fmt.Print("\n1:DOAR\n2:SOLICITAR DOACAO\n3:VERIFICAR PENDENTES\n0:SAIR\nESCOLHA UMA OPÇÃO: ")
		fmt.Scanln(&menu)
		switch menu {
		case 1:
			if fazerDoacao(PATH_DOADORES, verificarDoador(PATH_DOADORES)) {
				fmt.Println("Agradecemos a sua doação")
				break
			} else {
				fmt.Println("Ocorreu um erro e não foi possivel realizar a doação")
				break
			}
		case 2:
			listaDoadores(PATH_DOADORES)
			addPendente(PATH_PENDENTES)
		case 3:
			if !verificaPendentes(PATH_DOADORES, PATH_PENDENTES) {
				fmt.Println("\nOps! Ocorreu um erro")
			}
		case 0:
			menu = 0
		}
		if menu == 0 {
			break
		}
	}

	fmt.Println("\n\n\nAtendimento finalizado.\nCódigo desenvolvido por Felipe Teles(GitHub: felipersteles).\nPerdão pela ausencia de comentários :/\n\n\n")
}
