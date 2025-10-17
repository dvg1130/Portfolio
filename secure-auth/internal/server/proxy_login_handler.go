package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

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

	// Repackage response: extract actual fields from upstream JSON
	var upstreamMap map[string]interface{}
	if err := json.Unmarshal(resBody, &upstreamMap); err != nil {
		http.Error(w, "invalid upstream response 1", http.StatusInternalServerError)
		return
	}
	type UpstreamLoginResponse struct {
		DeviceID     string    `json:"device_id"`
		Message      string    `json:"message"`
		AccessToken  string    `json:"token"` // upstream calls it "token"
		RefreshToken string    `json:"refresh_token"`
		Expires      time.Time `json:"expires"` // only if upstream returns it
	}
	var upstreamResp UpstreamLoginResponse
	if err := json.Unmarshal(resBody, &upstreamResp); err != nil {
		http.Error(w, "invalid upstream response 2", http.StatusInternalServerError)
		return
	}

	resp := models.LoginResponse{
		Username:     creds.Username,            // inject username from request
		AccessToken:  upstreamResp.AccessToken,  // real token from upstream
		RefreshToken: upstreamResp.RefreshToken, // real refresh token from upstream
		DeviceID:     upstreamResp.DeviceID,
		Expires:      upstreamResp.Expires,
	}

	// Write status code
	w.WriteHeader(res.StatusCode)

	// Write JSON body
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		fmt.Println("Failed encoding repackaged response:", err)
	}
	defer res.Body.Close()
}
