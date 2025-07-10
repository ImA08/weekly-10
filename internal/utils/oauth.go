package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// OAuthUserInfo contains user information received from OAuth provider
type OAuthUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

// GetGoogleUserInfo fetches user info from Google API using the access token
func GetGoogleUserInfo(accessToken string) (*OAuthUserInfo, error) {
	if accessToken == "" {
		return nil, errors.New("access token is empty")
	}

	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("Authorization", "Bearer "+accessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user info: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error response from Google: %s (%d)", string(body), resp.StatusCode)
	}

	var userInfo OAuthUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %v", err)
	}

	return &userInfo, nil
}

// ValidateGoogleIDToken validates a Google ID token
func ValidateGoogleIDToken(idToken string) (*OAuthUserInfo, error) {
	if idToken == "" {
		return nil, errors.New("ID token is empty")
	}

	// Google's token info endpoint
	url := "https://oauth2.googleapis.com/tokeninfo?id_token=" + idToken
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("invalid token: %s", string(body))
	}

	var tokenInfo OAuthUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&tokenInfo); err != nil {
		return nil, fmt.Errorf("failed to decode token info: %v", err)
	}

	// Verify that token is issued for our client
	aud := tokenInfo.ID

	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	if clientID != "" && aud != clientID {
		return nil, errors.New("token was not issued for this application")
	}

	return &tokenInfo, nil
}

// ExtractProfileFromGothUser extracts profile information from Goth user data
func ExtractProfileFromGothUser(userData map[string]interface{}) (*OAuthUserInfo, error) {
	var userInfo OAuthUserInfo

	// Extract values safely
	if id, ok := userData["sub"].(string); ok {
		userInfo.ID = id
	} else if id, ok := userData["id"].(string); ok {
		userInfo.ID = id
	}

	if email, ok := userData["email"].(string); ok {
		userInfo.Email = email
	}

	if verified, ok := userData["email_verified"].(bool); ok {
		userInfo.VerifiedEmail = verified
	}

	if name, ok := userData["name"].(string); ok {
		userInfo.Name = name
	}

	if givenName, ok := userData["given_name"].(string); ok {
		userInfo.GivenName = givenName
	}

	if familyName, ok := userData["family_name"].(string); ok {
		userInfo.FamilyName = familyName
	}

	if picture, ok := userData["picture"].(string); ok {
		userInfo.Picture = picture
	}

	if locale, ok := userData["locale"].(string); ok {
		userInfo.Locale = locale
	}

	// Validate we have at least email
	if userInfo.Email == "" {
		return nil, errors.New("no email found in user data")
	}

	return &userInfo, nil
}

// GenerateUsername generates a username from email
func GenerateUsername(email string) string {
	parts := strings.Split(email, "@")
	username := parts[0]

	// Remove special characters
	username = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			return r
		}
		return '_'
	}, username)

	return username
}
