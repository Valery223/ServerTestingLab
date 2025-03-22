package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Valery223/ServerTestingLab/internal/logger"
	httpserver "github.com/Valery223/ServerTestingLab/internal/server/http/http_server"
)

func main() {
	serverErrors := make(chan error, 1)
	server := &httpserver.HTTPServer{}
	server.Init()

	// запуск сервера
	go func() {
		if err := server.Run(); err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	//ожидание завершения
	shutDown := make(chan os.Signal, 1)
	signal.Notify(shutDown, os.Interrupt)

	select {
	case sig := <-shutDown:
		logger.Logger.Info("Recived signal to shutdown", "signal", sig.String())
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			logger.Logger.Error("shutdown with err", "err", err)
			if err == context.DeadlineExceeded {
				logger.Logger.Error("Forcing server close")
				if err := server.Close(); err != nil {
					logger.Logger.Error("Server close error", "err", err)
				}
			}
		} else {
			logger.Logger.Info("Server shutdowned")
		}
	case err := <-serverErrors:
		logger.Logger.Error("Server error", "err", err)
	}

}
