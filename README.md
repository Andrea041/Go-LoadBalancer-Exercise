Esercizio in Go per il corso di SDCC

Per un corretto utilizzo del programma è necessario aver installato nel proprio PC:
- Go (in particolare l'applicazione è stata sviluppata nella verisone di Go 1.21.5)

L'applicazione sviluppata offre un servizio per la conversione della valuta da lira a euro (LiraToEuro) e da euro a lira (EuroToLira). In particolare per testare l'applicazione occorre eseguire i seguenti programmi:
- client.go
- loadBalancer.go
- server.go
Da terminale digitare:
- go run client.go
- go run loadBalancer.go
- go run server.go x (dove x sta ad indicare un numero da 0 a 2, poichè nel file di configurazione sono state create tre repliche del server)
E' necessario mettere in esecuzione tutti e tre i server per poter vedere una corretta esecuzione del programma.
