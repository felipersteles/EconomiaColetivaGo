package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

//Classes utilizadas
type Doadores struct {
	Nome          string
	Senha         int
	Idade         int
	Dia_nasc      int
	Mes_nasc      int
	Disponilidade bool
	Doando        []string
}
type Pendentes struct {
	Nome      string
	Interesse []string
}

//--------------Operacoes com arquivos--------------
func EncodeToBytes(p []Doadores) []byte {

	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println("uncompressed size (bytes): ", len(buf.Bytes()))
	return buf.Bytes()
}
func PendentesToBytes(p []Pendentes) []byte {

	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println("uncompressed size (bytes): ", len(buf.Bytes()))
	return buf.Bytes()
}
func Compress(s []byte) []byte {

	zipbuf := bytes.Buffer{}
	zipped := gzip.NewWriter(&zipbuf)
	zipped.Write(s)
	zipped.Close()
	//fmt.Println("compressed size (bytes): ", len(zipbuf.Bytes()))
	return zipbuf.Bytes()
}
func Decompress(s []byte) []byte {

	rdr, _ := gzip.NewReader(bytes.NewReader(s))
	data, err := ioutil.ReadAll(rdr)
	if err != nil {
		log.Fatal(err)
	}
	rdr.Close()
	//fmt.Println("uncompressed size (bytes): ", len(data))
	return data
}
func DecodeToDoadores(s []byte) []Doadores {

	var p []Doadores
	dec := gob.NewDecoder(bytes.NewReader(s))
	err := dec.Decode(&p)
	if err != nil {
		log.Fatal(err)
	}
	return p
}
func DecodeToPendentes(s []byte) []Pendentes {

	var p []Pendentes
	dec := gob.NewDecoder(bytes.NewReader(s))
	err := dec.Decode(&p)
	if err != nil {
		log.Fatal(err)
	}
	return p
}
func WriteToFile(s []byte, file string) {

	f, err := os.Create(file)
	if err != nil {
		log.Fatal(err)
	}

	f.Write(s)
}
func ReadFromFile(path string) []byte {

	f, err := os.Open(path)
	if err != nil {
		return nil
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil
	}

	defer f.Close()
	return data
}
func remove(s []string, i int) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

//--------------Operacoes com dados--------------
func verificaPendentes(path_d, path_p string) bool {
	var fileContent []byte
	var pendentes []Pendentes

	items := 0 //doacoes requisitadas

	//var aux string
	//scanner := bufio.NewScanner(os.Stdin)

	fileContent = ReadFromFile(path_p)
	if fileContent != nil {
		fileContent = Decompress(fileContent)
		pendentes = DecodeToPendentes(fileContent)

		indice := verificarDoador(path_d)
		doador := getDoador(path_d, indice)
		for i := 0; i < len(pendentes); i++ {
			for j := 0; j < len(pendentes[i].Interesse); j++ {
				for k := 0; k < len(doador.Doando); k++ {
					if strings.ToLower(doador.Doando[k]) == strings.ToLower(pendentes[i].Interesse[j]) {
						fmt.Println(pendentes[i].Nome, "tem interesse em", doador.Doando[k]) //Imprimir interessados e interesses
						items++
					}
				}
			}
		}
		if items > 0 {
			if finalizaDoacao(path_d, path_p, indice) {
				return true
			}
		} else {
			fmt.Println("Sem items")
		}
	}

	return false
}
func finalizaDoacao(path_d, path_p string, indice int) bool {
	var fileContent []byte
	scanner := bufio.NewScanner(os.Stdin)

	doador := getDoador(path_d, indice)

	fmt.Print("\nDeseja realizar alguma doação?(S/N)")
	scanner.Scan()
	aux := scanner.Text()
	if strings.ToLower(aux) == "s" {
		fmt.Print("\nDigite o item: ")
		scanner.Scan()
		item := scanner.Text()

		fmt.Print("\nDigite o remetente: ")
		scanner.Scan()
		remetente := scanner.Text()

		fileContent = ReadFromFile(path_p)
		if fileContent != nil {
			fileContent = Decompress(fileContent)
			pendentes := DecodeToPendentes(fileContent)
			for i := 0; i < len(pendentes); i++ {
				if strings.ToLower(pendentes[i].Nome) == strings.ToLower(remetente) {
					for j := 0; j < len(pendentes[i].Interesse); j++ {
						for k := 0; k < len(doador.Doando); k++ {
							if strings.ToLower(doador.Doando[k]) == strings.ToLower(item) {
								apagaItemDoador(path_d, indice, k)
								apagaItemPendente(path_p, i, j)
								return true
							}
						}
					}
				}
			}
		}

	}

	return false
}
func apagaItemDoador(path string, indice_doador, indice_item int) {

	var fileContent []byte

	fileContent = ReadFromFile(path)
	fileContent = Decompress(fileContent)
	doadores := DecodeToDoadores(fileContent)

	doadores[indice_doador].Doando = remove(doadores[indice_doador].Doando, indice_item)

	dataOut := EncodeToBytes(doadores)
	dataOut = Compress(dataOut)
	WriteToFile(dataOut, path)
}
func apagaItemPendente(path string, indice_pendente, indice_item int) {

	var fileContent []byte

	fileContent = ReadFromFile(path)
	fileContent = Decompress(fileContent)
	pendentes := DecodeToPendentes(fileContent)

	pendentes[indice_pendente].Interesse = remove(pendentes[indice_pendente].Interesse, indice_item)

	dataOut := PendentesToBytes(pendentes)
	dataOut = Compress(dataOut)
	WriteToFile(dataOut, path)
}
func addPendente(path string) bool {
	var aux string
	var fileContent []byte
	scanner := bufio.NewScanner(os.Stdin)

	var newPendentes Pendentes
	var pendentes []Pendentes

	fmt.Print("\nDigite seu nome: ")
	scanner.Scan()
	aux = scanner.Text()
	newPendentes.Nome = aux

	fmt.Print("\nDigite seu interesse: ")
	scanner.Scan()
	item := scanner.Text()
	newPendentes.Interesse = append(newPendentes.Interesse, item)

	fileContent = ReadFromFile(path)
	if fileContent != nil {
		fileContent = Decompress(fileContent)

		pendentes = DecodeToPendentes(fileContent)
		for i := 0; i < len(pendentes); i++ {
			if strings.ToLower(pendentes[i].Nome) == strings.ToLower(newPendentes.Nome) {
				for j := 0; j < len(pendentes[i].Interesse); j++ {
					if strings.ToLower(pendentes[i].Interesse[j]) == strings.ToLower(aux) {
						fmt.Println("\nInteresse já existente!\n")
						return false
					}
				}
				pendentes[i].Interesse = append(pendentes[i].Interesse, item)
				dataOut := PendentesToBytes(pendentes)
				dataOut = Compress(dataOut)
				WriteToFile(dataOut, path)

				fmt.Println("\nInteresse registrado!\n")
				return true
			}

		}
	}

	pendentes = append(pendentes, newPendentes)

	dataOut := PendentesToBytes(pendentes)
	dataOut = Compress(dataOut)
	WriteToFile(dataOut, path)

	fmt.Println("\nInteresse registrado!\n")
	return true
}
func listaDoadores(path string) {
	var fileContent []byte

	fileContent = ReadFromFile(path)
	fileContent = Decompress(fileContent)

	doadores := DecodeToDoadores(fileContent)
	for i := 0; i < len(doadores); i++ {
		for j := 0; j < len(doadores[i].Doando); j++ {
			fmt.Println(doadores[i].Nome, "esta doando", doadores[i].Doando[j])
		}
	}
}
func listaPendetes(path string) {
	var fileContent []byte

	fileContent = ReadFromFile(path)
	fileContent = Decompress(fileContent)

	pendentes := DecodeToPendentes(fileContent)
	for i := 0; i < len(pendentes); i++ {
		for j := 0; j < len(pendentes[i].Interesse); j++ {
			fmt.Println(pendentes[i].Nome, "tem interesse em", pendentes[i].Interesse[j])
		}
	}
}
func fazerDoacao(path string, indice int) bool {

	var aux string
	var fileContent []byte

	fileContent = ReadFromFile(path)
	fileContent = Decompress(fileContent)
	doadores := DecodeToDoadores(fileContent)

	fmt.Print("Olá ", doadores[indice].Nome, " o que pretendes doar?(Ou digite sair)\nR: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	aux = scanner.Text()
	if aux != "sair" {
		doadores[indice].Doando = append(doadores[indice].Doando, aux)

		dataOut := EncodeToBytes(doadores)
		dataOut = Compress(dataOut)
		WriteToFile(dataOut, path)

		return true
	}

	return false
}
func verificarDoador(path string) int {

	var aux string
	var fileContent []byte
	var senha int

	fmt.Print("\nDigite seu nome: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	aux = scanner.Text()

	fileContent = ReadFromFile(path)
	if fileContent != nil { //verifica se o arquivo existe
		fileContent = Decompress(fileContent)
		doadores := DecodeToDoadores(fileContent)
		if doadores != nil {
			for i := 0; i < len(doadores); i++ {
				if strings.ToLower(aux) == strings.ToLower(doadores[i].Nome) {
					//Doador existe
					fmt.Println("\nDigite a senha: ")
					fmt.Scanln(&senha)
					if senha == doadores[i].Senha {
						fmt.Println("\nSenha aceita!")
						return i //Acertou a senha
					}
				}
			} //caso nao exista cria o arquivo
			fmt.Println("Doador nao encontrado. Realizando cadastro.")
			fazerCadastro(path, aux)
			return len(doadores)
		}
	} else {
		fmt.Println("Doador nao encontrado. Realizando cadastro.")
		fazerCadastro(path, aux)
		return 0
	}
	return -1
}
func fazerCadastro(path, nome string) {

	var newDoador Doadores
	var aux int
	var doadores []Doadores

	fileContent := ReadFromFile(path)
	if fileContent != nil {
		fileContent = Decompress(fileContent)
		doadores = DecodeToDoadores(fileContent)

		newDoador.Nome = nome

		fmt.Print("\nDigite sua idade: ")
		fmt.Scanln(&aux)

		newDoador.Idade = aux

		fmt.Print("\nDigite sua senha: ")
		fmt.Scanln(&aux)

		newDoador.Senha = aux

		doadores = append(doadores, newDoador)
		dataOut := EncodeToBytes(doadores)
		dataOut = Compress(dataOut)
		WriteToFile(dataOut, path)
	} else {
		newDoador.Nome = nome

		fmt.Print("\nDigite sua idade: ")
		fmt.Scanln(&aux)

		newDoador.Idade = aux

		fmt.Print("\nDigite sua senha: ")
		fmt.Scanln(&aux)

		newDoador.Senha = aux

		doadores = append(doadores, newDoador)
		dataOut := EncodeToBytes(doadores)
		dataOut = Compress(dataOut)
		WriteToFile(dataOut, path)
	}
}
func getDoador(path string, indice int) Doadores {
	var doador Doadores
	var fileContent []byte

	fileContent = ReadFromFile(path)
	fileContent = Decompress(fileContent)
	doadores := DecodeToDoadores(fileContent)
	if doadores != nil {
		doador = doadores[indice]
	}

	return doador
}
func verificarDoacao(path string) {
	var aux string
	var fileContent []byte

	fmt.Print("\nDigite seu nome: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	aux = scanner.Text()

	fileContent = ReadFromFile(path)
	fileContent = Decompress(fileContent)
	doadores := DecodeToDoadores(fileContent)
	if doadores != nil {
		for i := 0; i < len(doadores); i++ {
			//fmt.Println("comparando ", aux, " com ", doadores[i].Nome, " ")
			if aux == strings.ToLower(doadores[i].Nome) {
				for j := 0; j < len(doadores[i].Doando); j++ {
					fmt.Println(doadores[i].Doando[j])
				}
			}
		}
	}
}
