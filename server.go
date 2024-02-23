package main

import (
	"encoding/json"
	"example/rpc_service"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"strconv"
)

type Configuration struct {
	Server1 ServerConfig `json:"server0"`
	Server2 ServerConfig `json:"server1"`
	Server3 ServerConfig `json:"server2"`
}

type ServerConfig struct {
	Address string `json:"address"`
	Port    string `json:"port"`
}

func main() {
	remoteservice := new(rpc_service.Calcolo)

	// Leggi il file di configurazione
	config, err := readConfig("configuration.json")
	if err != nil {
		log.Fatal("Errore durante in lettura del file di configurazione:", err)
	}

	// Questa slice contiene gli indirizzi dei server replicati
	var server_address = []string{config.Server1.Port, config.Server2.Port, config.Server3.Port}

	server := rpc.NewServer()
	err = server.Register(remoteservice)
	if err != nil {
		log.Fatal("Il formato del servizio di Calcolo è errato: ", err)
	}

	if len(os.Args) < 2 {
		log.Fatal("Input non valido")
	}
	addr, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal("Errore nella conversione: ", err)
	}

	list, err := net.Listen("tcp", server_address[addr])
	if err != nil {
		log.Fatal("Errore nell'instaurazione della connessione: ", err)
	}

	log.Printf("RPC server in ascolto sulla porta %s", server_address[addr])
	server.Accept(list)
	fmt.Printf("Connection accepted\n")
}

func readConfig(file string) (Configuration, error) {
	var config Configuration

	fileContent, err := os.ReadFile(file)
	if err != nil {
		return Configuration{}, err
	}

	err = json.Unmarshal(fileContent, &config)
	if err != nil {
		return Configuration{}, err
	}
	return config, nil
}
