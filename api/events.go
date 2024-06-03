package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Event struct {
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

type EventsResponse []Event

func GetEventsByCampusID(token string, campusID int) (EventsResponse, error) {
	url, err := url.Parse(baseApiUrl + "/v2/campus/" + strconv.Itoa(campusID) + "/events")
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

	data := EventsResponse{}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

type EventUsersResponse []struct {
	ID      int   `json:"id"`
	EventID int   `json:"event_id"`
	UserID  int   `json:"user_id"`
	User    User  `json:"user"`
	Event   Event `json:"event"`
}

type Pagination struct {
	PageNumber int
	PageSize   int
}

func buildPagination(page *Pagination) string {
	str := ""

	if page.PageNumber <= 0 {
		str = str + "?page[number]=1"
	} else {
		str = str + "?page[number]=" + strconv.Itoa(page.PageNumber)
	}

	if page.PageSize <= 0 {
		str = str + "?page[number]=30"
	} else if page.PageSize > 100 {
		str = str + "?page[number]=100"
	} else {
		str = str + "?page[number]=" + strconv.Itoa(page.PageSize)
	}

	return str
}

func GetEventUsersByEventId(token string, eventID int, page *Pagination) (EventUsersResponse, error) {
	url, err := url.Parse(baseApiUrl + "/v2/events/" + strconv.Itoa(eventID) + "/events_users" + buildPagination(page))
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

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Server returned " + resp.Status)
	}

	data := EventUsersResponse{}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
