// oauth.go
package discord

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	apiEndpoint  = "https://discord.com/api/v10"
	clientID     = "1247616920807669884"              // Replace with your actual client ID
	clientSecret = "j2Qetp1Q0HiJ1MRrQN1eSfGAXmzYnn5G" // Replace with your actual client secret
	redirectURI  = "http://localhost:8080/status"     // Replace with your actual redirect URI
)

// OAuthTokenResponse holds the structure for an OAuth token response
type OAuthTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

// ExchangeCode exchanges an OAuth2 authorization code for an access token
func ExchangeCode(code string) (*OAuthTokenResponse, error) {
	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/oauth2/token", apiEndpoint), strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error exchanging code for token, status code: %d", resp.StatusCode)
	}

	var tokenResponse OAuthTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return nil, err
	}

	return &tokenResponse, nil
}
