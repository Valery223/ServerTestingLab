package httpserver

import (
	"io"
	"net/http"
)

type UserHandler struct{}

func (uh *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	defer func() {
		// Дренирование
		_, _ = io.Copy(io.Discard, r.Body)
		r.Body.Close()
	}()

	_, _ = w.Write([]byte("Hello World"))
}
