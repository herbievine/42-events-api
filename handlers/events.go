package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type EventsResponse []struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Location       string    `json:"location"`
	Kind           string    `json:"kind"`
	MaxPeople      int       `json:"max_people"`
	NbrSubscribers int       `json:"nbr_subscribers"`
	BeginAt        time.Time `json:"begin_at"`
	EndAt          time.Time `json:"end_at"`
	CampusIds      []int     `json:"campus_ids"`
	CursusIds      []int     `json:"cursus_ids"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func GetEvents(w http.ResponseWriter, r *http.Request) {
	bearer := r.Header.Get("Authorization")
	if bearer == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	me, err := me(bearer)
	if err != nil {
		http.Error(w, "Failed to get current user", http.StatusInternalServerError)
		return
	}

	campus := me.Campus[0].ID

	url := baseApiUrl + "/v2/campus/" + strconv.Itoa(campus) + "/events"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Authorization", bearer)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "Failed to get events", http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to get events", resp.StatusCode)
		return
	}

	body := EventsResponse{}

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
