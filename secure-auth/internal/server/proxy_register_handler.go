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

func (s *Server) ProxyRegister(w http.ResponseWriter, r *http.Request) {
	internal := authsdk.NewClient(config.ProxyConfig.BASE_URL)

	var register_route = config.ProxyConfig.AUTH_REGISTER_ROUTE
	fmt.Println("route:", register_route)

	var creds models.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	body, _ := json.Marshal(creds)

	req, err := http.NewRequest("POST", register_route, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("proxy register route:", register_route)
		http.Error(w, "bad gateway 1", http.StatusBadGateway)
		return
	}
	req.Header = r.Header.Clone()

	res, err := internal.HTTPClient.Do(req)
	if err != nil {
		http.Error(w, "bad gateway 2", http.StatusBadGateway)
		return
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		http.Error(w, "failed to read upstream response 1", http.StatusInternalServerError)
		return
	}
	fmt.Println("Raw upstream response:", string(resBody))
	fmt.Println("Upstream status code:", res.StatusCode)
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
		http.Error(w, "invalid upstream response 2", http.StatusInternalServerError)
		return
	}
	type UpstreamRegisterResponse struct {
		Message string `json:"message"`
		// only if upstream returns it
	}
	var upstreamResp UpstreamRegisterResponse

	if err := json.Unmarshal(resBody, &upstreamResp); err != nil {
		http.Error(w, "invalid upstream response 3", http.StatusInternalServerError)
		return
	}

	resp := models.ResigterResponse{
		// inject username from request
		Message: upstreamResp.Message, // real token from upstream

	}

	// Write status code
	w.WriteHeader(res.StatusCode)

	// Write JSON body
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		fmt.Println("Failed encoding repackaged response:", err)
	}
	defer res.Body.Close()
}
