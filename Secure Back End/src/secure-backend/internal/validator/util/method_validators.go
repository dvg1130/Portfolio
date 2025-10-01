package validator

import "net/http"

func Method(method string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
			return
		}
		next(w, r)
	}
}
