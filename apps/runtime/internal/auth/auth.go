package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type DeviceCodeResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type ErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description,omitempty"`
}

func StartDeviceFlow(apiUrl string) (*DeviceCodeResponse, error) {
	resp, err := http.Post(fmt.Sprintf("%s/auth/device/code", apiUrl), "application/json", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err == nil {
			return nil, fmt.Errorf("cloud error: %s", errResp.Error)
		}
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var res DeviceCodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return &res, nil
}

func PollForToken(apiUrl string, deviceCode string, interval int) (*TokenResponse, error) {
	if interval <= 0 {
		interval = 5
	}

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			token, err := requestToken(apiUrl, deviceCode)
			if err != nil {
				// Continue polling if it's "authorization_pending"
				if err.Error() == "cloud error: authorization_pending" {
					continue
				}
				return nil, err
			}
			return token, nil
		}
	}
}

func requestToken(apiUrl string, deviceCode string) (*TokenResponse, error) {
	payload, _ := json.Marshal(map[string]string{
		"device_code": deviceCode,
		"grant_type":  "urn:ietf:params:oauth:grant-type:device_code",
	})

	resp, err := http.Post(fmt.Sprintf("%s/auth/device/token", apiUrl), "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err == nil {
			return nil, fmt.Errorf("cloud error: %s", errResp.Error)
		}
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var res TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return &res, nil
}
