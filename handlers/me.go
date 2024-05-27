package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/herbievine/42-events-api/api"
	"github.com/herbievine/42-events-api/db"
)

func GetMe(w http.ResponseWriter, r *http.Request, client *db.Client) {
	bearer := r.Header.Get("Authorization")
	if bearer == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	me, err := api.Me(bearer)
	if err != nil {
		http.Error(w, "Failed to get current user", http.StatusUnauthorized)
		return
	}

	user, _ := client.Users().GetOneByID(me.ID)

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
