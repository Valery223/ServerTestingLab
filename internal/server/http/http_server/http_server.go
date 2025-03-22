package httpserver

import (
	"net/http"

	"github.com/Valery223/ServerTestingLab/internal/logger"
)

type HTTPServer struct {
	server http.Server
}

func (hs *HTTPServer) Init() {
	mu := http.NewServeMux()
	userHandler := UserHandler{}
	handleGetUser := http.HandlerFunc(userHandler.GetUser)
	handleGetUser = logginMiddleWare(handleGetUser)
	mu.HandleFunc("GET /user", handleGetUser)
	mu.HandleFunc("POST /user", userHandler.PostUser)

	hs.server = http.Server{Addr: ":8081", Handler: mu}
	logger.Logger.Info("Server announced", "addr", hs.server.Addr)
}

func (hs *HTTPServer) Run() error {
	logger.Logger.Info("Server running", "addr", hs.server.Addr)
	err := hs.server.ListenAndServe()
	return err
}
