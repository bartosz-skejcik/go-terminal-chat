package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetOAuthToken(clientID, clientSecret string) (string, error) {
	url := "https://id.twitch.tv/oauth2/token"

	// Prepare the request body
	data := fmt.Sprintf(
		"client_id=%s&client_secret=%s&grant_type=client_credentials",
		clientID, clientSecret,
	)
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(data))
	if err != nil {
		return "", err
	}

	// Set the correct headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Make the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Parse the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	// Extract the access token
	if token, ok := result["access_token"].(string); ok {
		return token, nil
	}
	return "", fmt.Errorf("failed to get access token: %s", string(body))
}
