package main

import (
	"log"

	httpserver "github.com/Valery223/ServerTestingLab/internal/server/http/http_server"
)

func main() {
	server := &httpserver.HTTP_server{}
	server.Init()
	log.Fatal(server.Run())
}
