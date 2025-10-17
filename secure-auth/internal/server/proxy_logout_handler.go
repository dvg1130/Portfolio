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

func (s *Server) ProxyLogout(w http.ResponseWriter, r *http.Request) {

	internal := authsdk.NewClient(config.ProxyConfig.BASE_URL)

	var logout_route = config.ProxyConfig.AUTH_LOGOUT_ROUTE
	fmt.Println("route:", logout_route)

	var logoutReq models.LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&logoutReq); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	body, _ := json.Marshal(logoutReq)

	req, err := http.NewRequest("POST", logout_route, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("proxy login route:", logout_route)
		http.Error(w, "bad gateway1", http.StatusBadGateway)
		return
	}
	req.Header = r.Header.Clone()

	res, err := internal.HTTPClient.Do(req)
	if err != nil {
		http.Error(w, "bad gateway2", http.StatusBadGateway)
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
	type UpstreamLogoutResponse struct {
		// upstream calls it "token"
		Message string `json:"message"` // only if upstream returns it
	}
	var upstreamResp UpstreamLogoutResponse
	if err := json.Unmarshal(resBody, &upstreamResp); err != nil {
		http.Error(w, "invalid upstream response 2", http.StatusInternalServerError)
		return
	}

	resp := models.LogoutResponse{
		// real token from upstream
		Message: upstreamResp.Message,
	}

	// Write status code
	w.WriteHeader(res.StatusCode)

	// Write JSON body
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		fmt.Println("Failed encoding repackaged response:", err)
	}
	defer res.Body.Close()
}
