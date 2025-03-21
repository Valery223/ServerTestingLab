package httpserver

import "net/http"

type UserHandler struct{}

func (uh *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Write([]byte("Goood"))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
