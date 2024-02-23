Exercise in **Go** for the **SDCC course**

For correct usage of the program, it is necessary to have installed on your PC:
- Go (specifically, the application has been developed using Go version 1.21.5)

The developed application offers a service for currency conversion from lire to euros (LiraToEuro) and from euros to lire (EuroToLira). In particular, to test the application, you need to execute the following programs:
- client.go
- loadBalancer.go
- server.go
  
From the terminal, type:
- go run client.go
- go run loadBalancer.go
- go run server.go x (where x indicates a number from 0 to 2, as three replicas of the server have been created in the configuration file)

It is necessary to execute all three servers in order to see correct execution of the program.
