package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type CampusResponse []struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	TimeZone string `json:"time_zone"`
	Language struct {
		ID         int    `json:"id"`
		Name       string `json:"name"`
		Identifier string `json:"identifier"`
	} `json:"language"`
	UsersCount         int    `json:"users_count"`
	VogsphereID        int    `json:"vogsphere_id"`
	Country            string `json:"country"`
	Address            string `json:"address"`
	Zip                string `json:"zip"`
	City               string `json:"city"`
	Website            string `json:"website"`
	Facebook           string `json:"facebook"`
	Twitter            string `json:"twitter"`
	Active             bool   `json:"active"`
	Public             bool   `json:"public"`
	EmailExtension     string `json:"email_extension"`
	DefaultHiddenPhone bool   `json:"default_hidden_phone"`
	Endpoint           struct {
		ID          int       `json:"id"`
		URL         string    `json:"url"`
		Description string    `json:"description"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	} `json:"endpoint"`
}

func GetCampusByID(token string, campusID int) (CampusResponse, error) {
	url, err := url.Parse(baseApiUrl + "/v2/campus/" + strconv.Itoa(campusID))
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	data := CampusResponse{}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
