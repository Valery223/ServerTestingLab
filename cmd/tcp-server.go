package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/Valery223/ServerTestingLab/internal/server/tcp/server"
)

// for testing you can use netcat:
// nc localhost 8080

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	serv, _ := server.NewTCPServer("tcp", "localhost:8080", logger)
	go serv.Run()

	var inputCommand string
	for inputCommand != "stop" {
		fmt.Scanln(&inputCommand)
	}

	serv.Stop()
}
