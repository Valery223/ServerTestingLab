package httpserver

import (
	"context"
	"net/http"
	"time"

	"github.com/Valery223/ServerTestingLab/internal/logger"
)

type HTTPServer struct {
	server http.Server
}

func (hs *HTTPServer) Init() {
	mu := http.NewServeMux()

	userHandler := UserHandler{}
	//UserHandlers
	handleGetUser := http.HandlerFunc(userHandler.GetUser)
	handleGetUser = logginMiddleWare(handleGetUser)
	mu.HandleFunc("GET /user", handleGetUser)
	mu.HandleFunc("POST /user", userHandler.PostUser)

	//StaticHandlers
	fs := http.FileServer(http.Dir("./static"))
	handleFs := logginMiddleWare(http.StripPrefix("/static/", fs).ServeHTTP)
	mu.Handle("/static/", handleFs)

	hs.server = http.Server{Addr: ":4443",
		Handler:           mu,
		ReadHeaderTimeout: 3 * time.Second, // Random)
		IdleTimeout:       5 * time.Minute, // Again random:)
	}
	logger.Logger.Info("Server announced", "addr", hs.server.Addr)
}

func (hs *HTTPServer) Run() error {
	logger.Logger.Info("Server running", "addr", hs.server.Addr)
	return hs.server.ListenAndServeTLS("localhost.crt", "localhost.key")

}

func (hs *HTTPServer) Shutdown(ctx context.Context) error {
	logger.Logger.Info("Server started gracefull shutdawn")
	return hs.server.Shutdown(ctx)
}

func (hs *HTTPServer) Close() error {
	logger.Logger.Info("Server started forced shutdawn")
	return hs.server.Close()
}
