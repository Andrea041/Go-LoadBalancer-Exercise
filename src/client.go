package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"example/rpc_service"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

type dataToPass struct {
	Args        rpc_service.Args
	ServiceToDo string
	Reply       float64
}

func main() {
	loadBalancer_address := "localhost:" + "888"

	client, err := net.Dial("tcp", loadBalancer_address) // connessione al load balancer
	if err != nil {
		log.Fatal("Errore nella connessione con il load balancer: ", err)
	}

	for {
		fmt.Print("Inserire il numero che si intende convertire: ")
		firstNum := keyboardInput()

		n1, err := strconv.ParseFloat(firstNum, 64)
		if err != nil {
			log.Fatal("Conversion error: ", err)
		}

		// passaggio valore della chiamata a procedura
		argsToPass := rpc_service.Args{n1}
		var result float64

		fmt.Print("Scegliere uno dei due metodi del servizio (LireToEuro or EuroToLire): ")
		input := keyboardInput()

		data := dataToPass{
			Args:        argsToPass,
			ServiceToDo: input,
			Reply:       result,
		}

		buff := encodeStruct(data)
		_, err = client.Write(buff)
		if err != nil {
			log.Fatal("Errore nella scrittura: ", err)
		}

		var responseChan = make(chan float64)
		go func() {
			var bufReply [4096]byte
			n, err := client.Read(bufReply[:])
			if err != nil {
				log.Fatal("Errore in lettura: ", err)
			}
			procedureData := string(bufReply[:n])

			bufferReply := decodeStruct([]byte(procedureData))

			responseChan <- bufferReply.Reply
		}()

		// La lettura dei dati provenienti dal load balancer viene effettuata in modo asincrono, quindi possono essere effettuate eventuali altre operazioni

		// Viene effettuato il controllo per il recupero della risposta
		select {
		case response := <-responseChan:
			fmt.Printf("La risposta al servizio richiesto è: %f\n", response)
		case <-time.After(100 * time.Second):
			fmt.Println("Risposta persa")
		}
	}
}

func keyboardInput() string {
	var scanner *bufio.Scanner

	scanner = bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		log.Fatal("Errore nell'acquisizione dell'input: ", err)
	}
	return scanner.Text()
}

func encodeStruct(data dataToPass) []byte {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(data)
	if err != nil {
		fmt.Println("Errore durante la codifica della struct: ", err)
		return nil
	}
	return buffer.Bytes()
}

func decodeStruct(buffer []byte) dataToPass {
	var data dataToPass
	decoder := gob.NewDecoder(bytes.NewReader(buffer))
	err := decoder.Decode(&data)
	if err != nil {
		fmt.Println("Errore durante la decodifica della struct: ", err)
	}
	return data
}
