package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/herbievine/42-events-api/api"
	"github.com/herbievine/42-events-api/db"
)

func GetEvents(w http.ResponseWriter, r *http.Request, client *db.Client) {
	bearer := r.Header.Get("Authorization")
	if bearer == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	me, err := api.Me(bearer)
	if err != nil {
		http.Error(w, "Failed to get current user", http.StatusInternalServerError)
		return
	}

	events := make([]db.Event, 0)
	for _, campus := range me.Campus {
		campusEvents, err := client.Events().GetManyByCampusID(campus.ID)
		if err != nil {
			log.Println(err)
			http.Error(w, "Failed to get events", http.StatusInternalServerError)
			return
		}

		events = append(events, campusEvents...)
	}

	campusEvents, err := client.Events().GetManyByCampusID(29)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to get events", http.StatusInternalServerError)
		return
	}

	events = append(events, campusEvents...)

	sort.Slice(events, func(i, j int) bool {
		return events[i].CreatedAt.After(events[j].CreatedAt)
	})

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(events)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func NewEvents(w http.ResponseWriter, r *http.Request, client *db.Client) {
	token, err := api.GetServerToken()
	if err != nil {
		http.Error(w, "Failed to get access token", http.StatusInternalServerError)
		return
	}

	campuses, err := client.Campus().GetMany()
	if err != nil {
		http.Error(w, "Failed to get campuses", http.StatusInternalServerError)
		return
	}

	var events []db.Event
	for _, campus := range campuses {
		campusEvents, err := api.GetEventsByCampusID(token.AccessToken, campus.CampusID)
		if err != nil {
			log.Println(err)
			http.Error(w, "Failed to get events", http.StatusInternalServerError)
			return
		}

		campusUsers, err := client.Users().GetManyByCampusID(campus.CampusID)
		if err != nil {
			http.Error(w, "Failed to get users for campus", http.StatusInternalServerError)
			return
		}

		for _, event := range campusEvents {
			if e, _ := client.Events().GetOneByID(event.ID); e == nil {
				events = append(events, db.Event{
					EventID:      event.ID,
					Name:         event.Name,
					Description:  event.Description,
					Location:     event.Location,
					Type:         event.Kind,
					Attendees:    event.NbrSubscribers,
					MaxAttendees: event.MaxPeople,
					BeginAt:      event.BeginAt,
					EndAt:        event.EndAt,
					CampusIDs:    event.CampusIds,
					CursusIDs:    event.CursusIds,
					CreatedAt:    event.CreatedAt,
					UpdatedAt:    event.UpdatedAt,
				})

				var notifications []db.Notification
				for _, user := range campusUsers {
					notifications = append(notifications, db.Notification{
						UserID:    user.UserID,
						EventID:   event.ID,
						HasRead:   false,
						CreatedAt: time.Now(),
						DeletedAt: time.Time{},
					})

					log.Println("new event", event.Name, "for", user.Login, "has_read", notifications[len(notifications)-1].HasRead)
				}

				if _, err = client.Notifications().InsertMany(notifications); err != nil {
					http.Error(w, "Failed to save notifications", http.StatusInternalServerError)
					return
				}
			}
		}
	}

	if len(events) == 0 {
		io.WriteString(w, "no new events")
		w.Header().Set("Content-Type", "text/plain")
		return
	}

	if _, err = client.Events().InsertMany(events); err != nil {
		http.Error(w, "Failed to save events", http.StatusInternalServerError)
		return
	}

	io.WriteString(w, "saved "+strconv.Itoa(len(events))+" events")
	w.Header().Set("Content-Type", "text/plain")
}
