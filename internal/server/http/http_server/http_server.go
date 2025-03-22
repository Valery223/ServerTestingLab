package httpserver

import (
	"net/http"
)

type HTTPServer struct {
	server http.Server
}

func (hs *HTTPServer) Init() {
	mu := http.NewServeMux()
	userHandler := UserHandler{}
	mu.HandleFunc("GET /user", userHandler.GetUser)
	mu.HandleFunc("POST /user", userHandler.PostUser)

	hs.server = http.Server{Addr: ":8081", Handler: mu}
}

func (hs *HTTPServer) Run() error {
	err := hs.server.ListenAndServe()
	return err
}
