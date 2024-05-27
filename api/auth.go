package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const baseApiUrl string = "https://api.intra.42.fr"

type MeResponse struct {
	ID             int         `json:"id"`
	Email          string      `json:"email"`
	Login          string      `json:"login"`
	FirstName      string      `json:"first_name"`
	LastName       string      `json:"last_name"`
	UsualFullName  string      `json:"usual_full_name"`
	UsualFirstName interface{} `json:"usual_first_name"`
	URL            string      `json:"url"`
	Phone          string      `json:"phone"`
	Displayname    string      `json:"displayname"`
	Kind           string      `json:"kind"`
	Image          struct {
		Link     string `json:"link"`
		Versions struct {
			Large  string `json:"large"`
			Medium string `json:"medium"`
			Small  string `json:"small"`
			Micro  string `json:"micro"`
		} `json:"versions"`
	} `json:"image"`
	Staff           bool          `json:"staff?"`
	CorrectionPoint int           `json:"correction_point"`
	PoolMonth       string        `json:"pool_month"`
	PoolYear        string        `json:"pool_year"`
	Location        interface{}   `json:"location"`
	Wallet          int           `json:"wallet"`
	AnonymizeDate   time.Time     `json:"anonymize_date"`
	DataErasureDate time.Time     `json:"data_erasure_date"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
	AlumnizedAt     interface{}   `json:"alumnized_at"`
	Alumni          bool          `json:"alumni?"`
	Active          bool          `json:"active?"`
	Groups          []interface{} `json:"groups"`
	CursusUsers     []struct {
		Grade  interface{} `json:"grade"`
		Level  float64     `json:"level"`
		Skills []struct {
			ID    int     `json:"id"`
			Name  string  `json:"name"`
			Level float64 `json:"level"`
		} `json:"skills"`
		BlackholedAt interface{} `json:"blackholed_at"`
		ID           int         `json:"id"`
		BeginAt      time.Time   `json:"begin_at"`
		EndAt        time.Time   `json:"end_at"`
		CursusID     int         `json:"cursus_id"`
		HasCoalition bool        `json:"has_coalition"`
		CreatedAt    time.Time   `json:"created_at"`
		UpdatedAt    time.Time   `json:"updated_at"`
		User         struct {
			ID             int         `json:"id"`
			Email          string      `json:"email"`
			Login          string      `json:"login"`
			FirstName      string      `json:"first_name"`
			LastName       string      `json:"last_name"`
			UsualFullName  string      `json:"usual_full_name"`
			UsualFirstName interface{} `json:"usual_first_name"`
			URL            string      `json:"url"`
			Phone          string      `json:"phone"`
			Displayname    string      `json:"displayname"`
			Kind           string      `json:"kind"`
			Image          struct {
				Link     string `json:"link"`
				Versions struct {
					Large  string `json:"large"`
					Medium string `json:"medium"`
					Small  string `json:"small"`
					Micro  string `json:"micro"`
				} `json:"versions"`
			} `json:"image"`
			Staff           bool        `json:"staff?"`
			CorrectionPoint int         `json:"correction_point"`
			PoolMonth       string      `json:"pool_month"`
			PoolYear        string      `json:"pool_year"`
			Location        interface{} `json:"location"`
			Wallet          int         `json:"wallet"`
			AnonymizeDate   time.Time   `json:"anonymize_date"`
			DataErasureDate time.Time   `json:"data_erasure_date"`
			CreatedAt       time.Time   `json:"created_at"`
			UpdatedAt       time.Time   `json:"updated_at"`
			AlumnizedAt     interface{} `json:"alumnized_at"`
			Alumni          bool        `json:"alumni?"`
			Active          bool        `json:"active?"`
		} `json:"user"`
		Cursus struct {
			ID        int       `json:"id"`
			CreatedAt time.Time `json:"created_at"`
			Name      string    `json:"name"`
			Slug      string    `json:"slug"`
			Kind      string    `json:"kind"`
		} `json:"cursus"`
	} `json:"cursus_users"`
	ProjectsUsers []struct {
		ID            int    `json:"id"`
		Occurrence    int    `json:"occurrence"`
		FinalMark     int    `json:"final_mark"`
		Status        string `json:"status"`
		Validated     bool   `json:"validated?"`
		CurrentTeamID int    `json:"current_team_id"`
		Project       struct {
			ID       int         `json:"id"`
			Name     string      `json:"name"`
			Slug     string      `json:"slug"`
			ParentID interface{} `json:"parent_id"`
		} `json:"project"`
		CursusIds   []int     `json:"cursus_ids"`
		MarkedAt    time.Time `json:"marked_at"`
		Marked      bool      `json:"marked"`
		RetriableAt time.Time `json:"retriable_at"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	} `json:"projects_users"`
	LanguagesUsers []struct {
		ID         int       `json:"id"`
		LanguageID int       `json:"language_id"`
		UserID     int       `json:"user_id"`
		Position   int       `json:"position"`
		CreatedAt  time.Time `json:"created_at"`
	} `json:"languages_users"`
	Achievements []struct {
		ID           int         `json:"id"`
		Name         string      `json:"name"`
		Description  string      `json:"description"`
		Tier         string      `json:"tier"`
		Kind         string      `json:"kind"`
		Visible      bool        `json:"visible"`
		Image        string      `json:"image"`
		NbrOfSuccess interface{} `json:"nbr_of_success"`
		UsersURL     string      `json:"users_url"`
	} `json:"achievements"`
	Titles          []interface{} `json:"titles"`
	TitlesUsers     []interface{} `json:"titles_users"`
	Partnerships    []interface{} `json:"partnerships"`
	Patroned        []interface{} `json:"patroned"`
	Patroning       []interface{} `json:"patroning"`
	ExpertisesUsers []interface{} `json:"expertises_users"`
	Roles           []interface{} `json:"roles"`
	Campus          []struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		TimeZone string `json:"time_zone"`
		Language struct {
			ID         int       `json:"id"`
			Name       string    `json:"name"`
			Identifier string    `json:"identifier"`
			CreatedAt  time.Time `json:"created_at"`
			UpdatedAt  time.Time `json:"updated_at"`
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
	} `json:"campus"`
	CampusUsers []struct {
		ID        int       `json:"id"`
		UserID    int       `json:"user_id"`
		CampusID  int       `json:"campus_id"`
		IsPrimary bool      `json:"is_primary"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	} `json:"campus_users"`
}

type TokenResponse struct {
	AccessToken      string `json:"access_token"`
	TokenType        string `json:"token_type"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshToken     string `json:"refresh_token"`
	Scope            string `json:"scope"`
	CreatedAt        int    `json:"created_at"`
	SecretValidUntil int    `json:"secret"`
}

func Me(bearer string) (*MeResponse, error) {
	url := baseApiUrl + "/v2/me"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if strings.Split(bearer, " ")[0] != "Bearer" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	} else {
		req.Header.Set("Authorization", bearer)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Server returned " + resp.Status)
	}

	body := MeResponse{}

	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return nil, err
	}

	return &body, nil
}

func GetServerToken() (*TokenResponse, error) {
	url, err := url.Parse(baseApiUrl + "/oauth/token")
	if err != nil {
		return nil, err
	}

	query := url.Query()

	query.Set("grant_type", "client_credentials")
	query.Set("client_id", os.Getenv("FORTY_TWO_API_CLIENT"))
	query.Set("client_secret", os.Getenv("FORTY_TWO_API_SECRET"))

	url.RawQuery = query.Encode()

	resp, err := http.Post(url.String(), "application/json", nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Server returned " + resp.Status)
	}

	data := TokenResponse{}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
