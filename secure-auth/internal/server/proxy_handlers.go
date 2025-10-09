package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dvg1130/Portfolio/secure-auth/authsdk"
	"github.com/dvg1130/Portfolio/secure-auth/config"
	"github.com/dvg1130/Portfolio/secure-auth/models"
)

func (s *Server) ProxyLogin(w http.ResponseWriter, r *http.Request) {
	internal := authsdk.NewClient(config.ProxyConfig.BASE_URL)

	var auth_route = config.ProxyConfig.AUTH_LOGIN_ROUTE

	var creds models.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	fmt.Println("ProxyLogin ", creds.Username)

	body, _ := json.Marshal(creds)

	req, err := http.NewRequest("POST", config.ProxyConfig.BASE_URL+auth_route, bytes.NewBuffer(body))
	if err != nil {
		http.Error(w, "bad gateway", http.StatusBadGateway)
		return
	}
	req.Header = r.Header.Clone()

	res, err := internal.HTTPClient.Do(req)
	if err != nil {
		http.Error(w, "bad gateway", http.StatusBadGateway)
		return
	}

	// Copy headers
	for k, v := range res.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}

	// Set status code
	w.WriteHeader(res.StatusCode)

	// Copy raw response body directly
	if _, err := io.Copy(w, res.Body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()
}
