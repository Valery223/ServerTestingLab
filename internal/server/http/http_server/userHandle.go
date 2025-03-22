package httpserver

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Valery223/ServerTestingLab/internal/logger"
)

type UserHandler struct{}

// Echo ответ, получая имя из header "Name"
func (uh *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	defer func() {
		// Дренирование
		_, _ = io.Copy(io.Discard, r.Body)
		r.Body.Close()
	}()

	name := r.Header.Get("Name")
	if name != "" {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		n, err := w.Write([]byte(fmt.Sprintf("Received header Name: %s", name)))
		if err != nil {
			logger.Logger.Error("request progress",
				"method", r.Method,
				"path", r.URL.Path,
				"remote_addr", r.RemoteAddr,
				"error", err)
		} else {
			logger.Logger.Debug("request progress",
				"method", r.Method,
				"path", r.URL.Path,
				"remote_addr", r.RemoteAddr,
				"bytes_sent", n)
		}
	} else {
		logger.Logger.Warn("request progress",
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
			"error", "name not find")

		// В ручную сначала хочу отправлять ошибку, для обучения
		// http.Error(w, "Name header is required", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Name header is required"))
	}
}

// Логирование body, полученного из запроса
func (uh *UserHandler) PostUser(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_, _ = io.Copy(io.Discard, r.Body)
		r.Body.Close()
	}()

	switch r.Header.Get("Content-Type") {
	case "text/plain":
		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Logger.Error("get error", "error", err)
		} else {
			logger.Logger.Info("read from body",
				"body", string(body))
		}
	default:
		logger.Logger.Warn("request progress",
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
			"error", "Unsupported media type")
		http.Error(w, "Unsupported media type", http.StatusUnsupportedMediaType)

	}
}
