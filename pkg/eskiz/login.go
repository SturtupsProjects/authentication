package eskiz

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type LoginResponse struct {
	Message string `json:"message"`
	Data    struct {
		Token string `json:"token"`
	} `json:"data"`
	TokenType string `json:"token_type"`
}

func LoginToEskiz(email, password string) (string, error) {
	url := "https://notify.eskiz.uz/api/auth/login"

	// Create request body
	payload := map[string]string{
		"email":    email,
		"password": password,
	}
	jsonPayload, _ := json.Marshal(payload)

	// Create HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Decode response
	var loginResp LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return "", err
	}

	// Check if token is returned
	if loginResp.Data.Token == "" {
		return "", errors.New("failed to retrieve token")
	}

	return loginResp.Data.Token, nil
}
