package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"example/rpc_service"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"time"
)

type dataToPass struct {
	Args        rpc_service.Args
	ServiceToDo string
	Reply       float64
}

type Configuration struct {
	Server1 ServerK `json:"server0"`
	Server2 ServerK `json:"server1"`
	Server3 ServerK `json:"server2"`
}

type ServerK struct {
	Address string `json:"address"`
	Port    string `json:"port"`
}

func main() {
	list, err := net.Listen("tcp", ":888")
	if err != nil {
		log.Fatal("Errore nell'instaurazione della connessione: ", err)
	}

	// Lettura del file di configurazione
	config, err := readConfig("configuration.json")
	if err != nil {
		log.Fatal("Errore durante la lettura del file di configurazione:", err)
	}

	// Questa slice contiene gli indirizzi dei server replicati
	var server_address = []ServerK{config.Server1, config.Server2, config.Server3}

	log.Printf("Il load balancer si trova in ascolto sulla porta %d", 888)
	clientArrived, err := list.Accept()
	if err != nil {
		log.Fatal("Errore di connessione: ", err)
	}

	index := 0
	for {
		var buf [4096]byte
		n, err := clientArrived.Read(buf[:])
		if err != nil {
			log.Fatal("Errore in lettura: ", err)
		}
		procedureData := string(buf[:n])

		buff := decodeStruct([]byte(procedureData))

		// Round Robin per la selezione dei server replicati
		serv_choice := server_address[index]
		index = (index + 1) % len(server_address)

		structChan := make(chan dataToPass)
		go connectToSelectedDestination(serv_choice, buff, structChan)

		select {
		case bufferStruct := <-structChan:
			buffResponse := encodeStruct(bufferStruct)
			_, err = clientArrived.Write(buffResponse)
			if err != nil {
				log.Fatal("Errore nella scrittura: ", err)
			}
		case <-time.After(100 * time.Second):
			fmt.Println("Risposta persa")
		}
	}
}

func connectToSelectedDestination(server ServerK, data dataToPass, temp chan dataToPass) {
	serv := server.Address + server.Port
	fmt.Printf("Connessione al server con indirizzo: %s\n", server)
	loadClient, err := rpc.Dial("tcp", serv)
	if err != nil {
		log.Fatal("Errore di connessione: ", err)
	}

	number := rpc_service.Args{data.Args.toConvert}
	var rep float64

	if data.ServiceToDo == "LireToEuro" {
		call := loadClient.Go("Calcolo.LireToEuro", number, &rep, nil)
		call = <-call.Done
		if call.Error != nil {
			log.Fatal("Errore nella chiamata asincrona: ", call.Error)
		}
	} else if data.ServiceToDo == "EuroToLire" {
		call := loadClient.Go("Calcolo.EuroToLire", number, &rep, nil)
		call = <-call.Done
		if call.Error != nil {
			log.Fatal("Errore nella chiamata asincrona: ", call.Error)
		}
	}

	data.Reply = rep
	temp <- data
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

func readConfig(file string) (Configuration, error) {
	var config Configuration

	tmp, err := os.ReadFile(file)
	if err != nil {
		return Configuration{}, err
	}

	err = json.Unmarshal(tmp, &config)
	if err != nil {
		return Configuration{}, err
	}
	return config, nil
}
