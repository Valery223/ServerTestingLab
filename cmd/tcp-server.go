package main

import (
	"fmt"

	"github.com/Valery223/ServerTestingLab/internal/server/tcp/server"
)

// for testing you can use netcat:
// nc localhost 8080

func main() {
	serv, _ := server.NewTCPServer("tcp", "localhost:8080")
	go serv.Run()

	var inputCommand string
	for inputCommand != "stop" {
		fmt.Scanln(&inputCommand)
	}

	serv.Stop()
}
