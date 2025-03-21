package httpserver

import (
	"net/http"
)

type HTTP_server struct {
	server http.Server
}

func (hs *HTTP_server) Init() {
	mu := http.NewServeMux()
	userHandler := UserHandler{}
	mu.Handle("/user", &userHandler)

	hs.server = http.Server{Addr: ":8081", Handler: mu}
}

func (hs *HTTP_server) Run() error {
	err := hs.server.ListenAndServe()
	return err
}
