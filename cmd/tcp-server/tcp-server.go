package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/Valery223/ServerTestingLab/internal/servers/tcp/server"
)

// for testing you can use netcat:
// nc localhost 8080

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	serv, _ := server.NewTCPServer("tcp", "localhost:8080", logger)
	ctx, cancel := context.WithCancel(context.Background())
	go serv.Run(ctx)

	var inputCommand string
	for inputCommand != "stop" {
		fmt.Scanln(&inputCommand)
	}

	cancel()
	time.Sleep(time.Second * 5)
}
