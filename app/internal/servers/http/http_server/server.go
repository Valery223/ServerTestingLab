package httpserver

import (
	"context"
	"net/http"

	"github.com/Valery223/ServerTestingLab/internal/config"
	"github.com/Valery223/ServerTestingLab/internal/logger"
)

type Server struct {
	server http.Server
}

func (hs *Server) Init(cfg config.HTTPServer) {
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

	hs.server = http.Server{Addr: cfg.Address,
		Handler:           mu,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		IdleTimeout:       cfg.IdleTimeout,
	}
	logger.Logger.Info("Server announced", "addr", hs.server.Addr)
}

func (hs *Server) Run() error {
	logger.Logger.Info("Server running", "addr", hs.server.Addr)
	return hs.server.ListenAndServe()

}

func (hs *Server) Shutdown(ctx context.Context) error {
	logger.Logger.Info("Server started gracefull shutdawn")
	return hs.server.Shutdown(ctx)
}

func (hs *Server) Close() error {
	logger.Logger.Info("Server started forced shutdawn")
	return hs.server.Close()
}
