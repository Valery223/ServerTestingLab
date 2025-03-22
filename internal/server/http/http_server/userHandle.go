package httpserver

import (
	"fmt"
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

func (uh *UserHandler) PostUser(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_, _ = io.Copy(io.Discard, r.Body)
		r.Body.Close()
	}()

	name := r.Header.Get("Name")
	// Если значение не пустое, добавляем его в заголовки ответа под корректным названием
	if name != "" {
		w.Header().Add("Echo", name)
		// Также можно вернуть значение в теле ответа
		_, _ = w.Write([]byte(fmt.Sprintf("Received header Name: %s", name)))
	} else {
		http.Error(w, "Header 'Name' is required", http.StatusBadRequest)
	}
}
