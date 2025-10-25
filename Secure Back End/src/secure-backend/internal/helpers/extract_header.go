package helpers

import (
	"net/http"
	"strings"
)

func ExtractHeader_Token(w http.ResponseWriter, r *http.Request) (tokenString string) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
		return
	}

	tokenString = parts[1]
	return tokenString
}

func ExtractCookieValue(cookie_name string, w http.ResponseWriter, r *http.Request) (cookieString string) {

	cookie, err := r.Cookie(cookie_name)
	if err != nil {
		// handle missing cookie
		http.Error(w, "missinfg or invalid cookie", http.StatusBadRequest)
		return
	}
	cookieString = cookie.Value

	return cookieString
}
