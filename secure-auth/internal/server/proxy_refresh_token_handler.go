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

func (s *Server) ProxyRefreshToken(w http.ResponseWriter, r *http.Request) {

	//new client
	internal := authsdk.NewClient(config.ProxyConfig.BASE_URL)

	//pass route
	var refreshtoken_route = config.ProxyConfig.AUTH_REFRESH_TOKEN_ROUTE
	fmt.Println("refresh token route: ", refreshtoken_route)

	//decode body json
	var refreshReq models.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&refreshReq); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	//pass json
	body, _ := json.Marshal(refreshReq)

	req, err := http.NewRequest("POST", refreshtoken_route, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("proxy refresh route:", refreshReq)
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
	fmt.Println("Raw upstream response body:", string(resBody))
	var upstreamMap map[string]interface{}
	if err := json.Unmarshal(resBody, &upstreamMap); err != nil {
		http.Error(w, "invalid upstream response 1", http.StatusInternalServerError)
		fmt.Println("map: ", upstreamMap)
		return
	}

	type UpstreamRefreshResponse struct {
		// upstream calls it "token"

		AccessToken  string    `json:"access_token"`
		RefreshToken string    `json:"refresh_token"`
		DeviceID     string    `json:"device_id"`
		Expires      time.Time `json:"expires"` // only if upstream returns it
	}

	var upstreamResp UpstreamRefreshResponse

	if err := json.Unmarshal(resBody, &upstreamResp); err != nil {
		http.Error(w, "invalid upstream response 2", http.StatusInternalServerError)
		return
	}

	resp := models.RefreshResponse{
		// real token from upstream

		AccessToken:  upstreamResp.AccessToken,
		RefreshToken: upstreamResp.RefreshToken,
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
