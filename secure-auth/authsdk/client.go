package authsdk

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func (c *SDKClient) SDKLogin(models.Credentials) (*models.LoginResponse, error) {
	body, _ := json.Marshal(creds)

	req, err := http.NewRequest("POST", "http:127.0.0.1/api/sdk/login", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("SDKLogin:", c.BaseURL)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var authResp models.LoginResponse
	if err := json.NewDecoder(res.Body).Decode(&authResp); err != nil {
		return nil, err
	}
	//dev and testing only
	fmt.Println("Passed SDKLogin proxy", authResp.Username)

	return &authResp, nil
}
