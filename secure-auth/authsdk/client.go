package authsdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dvg1130/Portfolio/secure-auth/models"
)

type SDKClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(baseURL string) *SDKClient {

	return &SDKClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

var creds models.Credentials

func (c *SDKClient) SDKLogin(creds models.Credentials) (*models.LoginResponse, error) {
	// marshal creds into JSON
	bodyBytes, err := json.Marshal(creds)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal credentials: %w", err)
	}

	// build the request
	url := "http://127.0.0.1:8080/api/sdk/login" // make sure port matches ProxyLogin server
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

	// read and decode response as JSON
	var authResp models.LoginResponse
	if err := json.NewDecoder(res.Body).Decode(&authResp); err != nil {
		// Log raw body for debugging if JSON decode fails
		rawBody, _ := io.ReadAll(res.Body)
		fmt.Println("Failed to decode response, raw body:", string(rawBody))
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	//  optional debug log
	fmt.Println("Passed SDKLogin proxy", authResp.Username)

	return &authResp, nil
}
