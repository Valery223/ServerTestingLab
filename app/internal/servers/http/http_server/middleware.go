package httpserver

import (
	"net/http"
	"time"

	"github.com/Valery223/ServerTestingLab/app/internal/logger"
)

func logginMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logger.Logger.Info("request started",
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr)

		defer func() {
			duraction := time.Since(start)

			logger.Logger.Info("request completed",
				"method", r.Method,
				"path", r.URL.Path,
				"duraction", duraction.String())
		}()

		next.ServeHTTP(w, r)

	})
}
