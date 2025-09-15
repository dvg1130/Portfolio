package server

import "net/http"

func ServeMuxWrapper(mux *http.ServeMux, middlewares ...func(http.Handler) http.Handler) *http.ServeMux {
	wrapped := http.Handler(mux)

	//  reverse irder so first argument is the outer wrapper
	for i := len(middlewares) - 1; i >= 0; i-- {
		wrapped = middlewares[i](wrapped)
	}

	newMux := http.NewServeMux()
	newMux.Handle("/", wrapped)
	return newMux
}
