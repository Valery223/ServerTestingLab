package httpserver

import (
	"net/http"

	"github.com/Valery223/ServerTestingLab/internal/logger"
)

func logRequestDebug(r *http.Request, msg string, args ...any) {
	logger.Logger.Debug(msg,
		"method", r.Method,
		"path", r.URL.Path,
		"remote_addr", r.RemoteAddr,
	)
}
