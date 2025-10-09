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

	body, _ := json.Marshal(creds)
	fmt.Println("proxy login route:", auth_route)

	req, err := http.NewRequest("POST", auth_route, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("proxy login route:", auth_route)
		http.Error(w, "bad gateway", http.StatusBadGateway)
		return
	}
	req.Header = r.Header.Clone()

	res, err := internal.HTTPClient.Do(req)
	if err != nil {
		http.Error(w, "bad gateway", http.StatusBadGateway)
		return
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		http.Error(w, "failed to read upstream response", http.StatusInternalServerError)
		return
	}

	// Copy headers except Content-Length
	for k, v := range res.Header {
		if k == "Content-Length" || k == "Transfer-Encoding" {
			continue
		}
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}

	// Ensure Content-Type is JSON
	w.Header().Set("Content-Type", "application/json")

	// ✅ Repackage response: extract actual fields from upstream JSON
	var upstreamMap map[string]interface{}
	if err := json.Unmarshal(resBody, &upstreamMap); err != nil {
		http.Error(w, "invalid upstream response", http.StatusInternalServerError)
		return
	}

	resp := models.LoginResponse{
		Username:     creds.Username,
		AccessToken:  fmt.Sprintf("%v", upstreamMap["token"]), // map "token" → AccessToken
		RefreshToken: fmt.Sprintf("%v", upstreamMap["refresh_token"]),
		DeviceID:     fmt.Sprintf("%v", upstreamMap["device_id"]),
	}

	// Write status code
	w.WriteHeader(res.StatusCode)

	// Write JSON body
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		fmt.Println("Failed encoding repackaged response:", err)
	}
	defer res.Body.Close()
}
