package authsdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dvg1130/Portfolio/secure-auth/models"
)

func (c *SDKClient) SDKLogout(logout models.LogoutRequest) (*models.LogoutResponse, error) {
	// marshal creds into JSON
	bodyBytes, err := json.Marshal(logout)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal credentials: %w", err)
	}

	// build the request
	url := "http://127.0.0.1:8080/api/sdk/logout"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request to %s: %w", url, err)
	}

	// sSet Content-Type to JSON
	req.Header.Set("Content-Type", "application/json")

	// send the request
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	var logoutResp models.LogoutResponse
	if err := json.NewDecoder(res.Body).Decode(&logoutResp); err != nil {
		// Log raw body for debugging if JSON decode fails
		rawBody, _ := io.ReadAll(res.Body)
		fmt.Println("Failed to decode response, raw body:", string(rawBody))
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &logoutResp, nil

}
