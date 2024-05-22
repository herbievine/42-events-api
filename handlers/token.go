package handlers

import (
	"encoding/json"
	"net/http"
	"os"
)

const (
	baseApiUrl string = "https://api.intra.42.fr"
)

type TokenResponse struct {
	AccessToken      string `json:"access_token"`
	TokenType        string `json:"token_type"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshToken     string `json:"refresh_token"`
	Scope            string `json:"scope"`
	CreatedAt        int    `json:"created_at"`
	SecretValidUntil int    `json:"secret"`
}

func GetToken(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if code == "" || state == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	redirectUri := os.Getenv("FRONTEND_URL") + "/auth/callback"

	url, err := r.URL.Parse(baseApiUrl + "/oauth/token")
	if err != nil {
		http.Error(w, "Failed to parse URL", http.StatusInternalServerError)
		return
	}

	query := url.Query()

	query.Set("grant_type", "authorization_code")
	query.Set("client_id", os.Getenv("FORTY_TWO_API_CLIENT"))
	query.Set("client_secret", os.Getenv("FORTY_TWO_API_SECRET"))
	query.Set("code", code)
	query.Set("redirect_uri", redirectUri)
	query.Set("state", state)

	url.RawQuery = query.Encode()

	resp, err := http.Post(url.String(), "application/json", nil)
	if err != nil {
		http.Error(w, "Failed to get access token", http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to get access token", resp.StatusCode)
		return
	}

	body := TokenResponse{}

	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Failed to parse response body", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(body)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
