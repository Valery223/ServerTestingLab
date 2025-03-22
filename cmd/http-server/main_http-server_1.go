package main

import (
	"log"

	httpserver "github.com/Valery223/ServerTestingLab/internal/server/http/http_server"
)

func main() {
	server := &httpserver.HTTPServer{}
	server.Init()
	log.Fatal(server.Run())
}
