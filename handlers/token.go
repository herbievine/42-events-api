package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/herbievine/42-events-api/api"
	"github.com/herbievine/42-events-api/auth"
	"github.com/herbievine/42-events-api/db"
)

const (
	baseApiUrl string = "https://api.intra.42.fr"
)

type tokenResponse struct {
	JWT string `json:"jwt"`
}

func GetToken(w http.ResponseWriter, r *http.Request, client *db.Client) {
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

	token := api.TokenResponse{}

	err = json.NewDecoder(resp.Body).Decode(&token)
	if err != nil {
		http.Error(w, "Failed to parse response body", http.StatusInternalServerError)
		return
	}

	me, err := api.Me(token.AccessToken)
	if err != nil {
		log.Fatalln(err)
		http.Error(w, "Failed to get current user", http.StatusInternalServerError)
		return
	}

	if _, err := client.Users().GetOneByID(me.ID); err != nil {
		for _, campus := range me.Campus {
			if _, err := client.Campus().GetOneByID(campus.ID); err != nil {
				_, err := client.Campus().InsertOne(db.Campus{
					CampusID:  campus.ID,
					Name:      campus.Name,
					UserCount: campus.UsersCount,
					City:      campus.City,
					Country:   campus.Country,
				})
				if err != nil {
					http.Error(w, "Failed to save campus", http.StatusInternalServerError)
					return
				}
			}
		}

		var campusIds []int
		var primaryCampusID int
		for _, campus := range me.CampusUsers {
			campusIds = append(campusIds, campus.CampusID)

			if campus.IsPrimary {
				primaryCampusID = campus.CampusID
				break
			}
		}

		_, err = client.Users().InsertOne(db.User{
			UserID:          me.ID,
			Login:           me.Login,
			ImageURL:        me.Image.Link,
			CampusIDs:       campusIds,
			PrimaryCampusID: primaryCampusID,
			LastSeen:        time.Now(),
			CreatedAt:       time.Now(),
		})
		if err != nil {
			http.Error(w, "Failed to save user", http.StatusInternalServerError)
			return
		}
	}

	jwtClaims := auth.UserClaims{
		UserID:      me.ID,
		AccessToken: token.AccessToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
		},
	}

	jwt, err := auth.Issue(jwtClaims)
	if err != nil {
		log.Println("[ERROR] Failed to create JWT:", err)

		http.Error(w, "Failed to create JWT", http.StatusInternalServerError)
		return
	}

	jwtResp := tokenResponse{
		JWT: jwt,
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(jwtResp)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
