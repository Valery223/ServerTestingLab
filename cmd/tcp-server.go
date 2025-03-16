package main

import "github.com/Valery223/ServerTestingLab/internal/server/tcp/server"

// for testing you can use netcat:
// nc -l -p 8080

func main() {
	serv, _ := server.NewTCPServer("tcp", "localhost:8080")
	serv.Run()
}
