package authsdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dvg1130/Portfolio/secure-auth/models"
)

func (c *SDKClient) SDKRefreshToken(refresh models.RefreshRequest) (*models.RefreshResponse, error) {
	// marshal creds into JSON
	bodyBytes, err := json.Marshal(refresh)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal credentials: %w", err)
	}

	// build the request

	url := "http://127.0.0.1:8080/api/sdk/token/refresh"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Authorization", "Bearer "+refresh.AccessToken)
	fmt.Println("passed token: ", refresh.AccessToken)
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

	var refreshResp models.RefreshResponse
	if err := json.NewDecoder(res.Body).Decode(&refreshResp); err != nil {
		// Log raw body for debugging if JSON decode fails
		rawBody, _ := io.ReadAll(res.Body)
		fmt.Println("Failed to decode response, raw body:", string(rawBody))
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &refreshResp, nil

}
