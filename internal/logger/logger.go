package logger

import (
	"log/slog"
	"net/http"
	"os"
)

var Logger *slog.Logger

func init() {
	Logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
}

func logRequest(r *http.Request, level slog.Level, msg string, args ...any) {
	requestLogArgs := []any{
		"method", r.Method,
		"path", r.URL.Path,
		"remote_addr", r.RemoteAddr,
	}

	requestLogArgs = append(requestLogArgs, args...)

	Logger.Log(r.Context(), level, msg, requestLogArgs...)
}

func LogRequestInfo(r *http.Request, msg string, args ...any) {
	logRequest(r, slog.LevelInfo, msg, args...)
}
func LogRequestDebug(r *http.Request, msg string, args ...any) {
	logRequest(r, slog.LevelDebug, msg, args...)
}

func LogRequestError(r *http.Request, msg string, err error, args ...any) {
	args = append(args, "err", err)
	logRequest(r, slog.LevelError, msg, args...)
}

func LogRequestWarn(r *http.Request, msg string, args ...any) {
	logRequest(r, slog.LevelWarn, msg, args...)
}
