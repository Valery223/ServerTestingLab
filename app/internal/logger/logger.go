package logger

import (
	"log/slog"
	"net/http"
	"os"
)

type Env string

const (
	EnvLocal Env = "local"
	EnvDev   Env = "dev"
	EnvProd  Env = "prod"
)

var Logger *slog.Logger

func Init(env Env) {

	switch env {
	case EnvLocal:
		Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case EnvDev:
		Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case EnvProd:
		Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
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
