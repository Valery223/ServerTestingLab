package httpserver

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Valery223/ServerTestingLab/app/internal/logger"
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
	if len(name) == 0 {
		logger.LogRequestWarn(r, "missing name header")
		http.Error(w, "Name header is required", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	n, err := w.Write([]byte(fmt.Sprintf("Received header Name: %s", name)))
	if err != nil {
		logger.LogRequestWarn(r, "Error write response in body", "err", err)
	} else {
		logger.LogRequestDebug(r, "Writed in body response", "write_size", n)
	}

}

// Логирование body, полученного из запроса
func (uh *UserHandler) PostUser(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_, _ = io.Copy(io.Discard, r.Body)
		r.Body.Close()
	}()

	switch contentType := r.Header.Get("Content-Type"); contentType {
	case "text/plain":
		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.LogRequestError(r, "Read body error", err)

			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		} else {
			// Временно, для проверки
			logger.Logger.Info("read from body",
				"body", string(body))
		}
	default:
		logger.LogRequestWarn(r, "Unsupported media type",
			"content_type", contentType)
		http.Error(w, "Unsupported media type", http.StatusUnsupportedMediaType)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))

	logger.LogRequestDebug(r, "Hendle OK")
}
