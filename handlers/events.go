package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/herbievine/42-events-api/api"
	"github.com/herbievine/42-events-api/db"
	"go.mongodb.org/mongo-driver/bson"
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

type NewEventsResponse struct {
	EventsAdded   int `json:"events_added"`
	EventsUpdated int `json:"events_updated"`
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
	var response NewEventsResponse

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
			if event.EndAt.Before(time.Now()) {
				continue
			}

			var eventUsers api.EventUsersResponse
			page := api.Pagination{
				PageNumber: 1,
				PageSize:   100,
			}

			for {
				resp, err := api.GetEventUsersByEventId(token.AccessToken, event.ID, &page)
				if err != nil {
					log.Println("[WARN] Failed to get event users", event.ID, page, err)
					continue
				}

				if len(resp) == 0 {
					break
				}

				eventUsers = append(eventUsers, resp...)
				page.PageNumber++

				time.Sleep(1 * time.Second)
			}

			eventInDB, err := client.Events().GetOneByID(event.ID)
			if err != nil {
				events = append(events, db.Event{
					EventID:      event.ID,
					Name:         event.Name,
					Description:  event.Description,
					Location:     event.Location,
					Type:         event.Kind,
					Attendees:    len(eventUsers),
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
				}

				if _, err = client.Notifications().InsertMany(notifications); err != nil {
					http.Error(w, "Failed to save notifications", http.StatusInternalServerError)
					return
				}
			} else if eventInDB != nil {
				filter := bson.D{
					{Key: "event_id", Value: eventInDB.EventID},
				}
				update := bson.D{
					{Key: "$set", Value: bson.D{
						{Key: "name", Value: event.Name},
						{Key: "description", Value: event.Description},
						{Key: "location", Value: event.Location},
						{Key: "type", Value: event.Kind},
						{Key: "attendees", Value: len(eventUsers)},
						{Key: "max_attendees", Value: event.MaxPeople},
						{Key: "begin_at", Value: event.BeginAt},
						{Key: "end_at", Value: event.EndAt},
						{Key: "campus_ids", Value: event.CampusIds},
						{Key: "cursus_ids", Value: event.CursusIds},
						{Key: "created_at", Value: event.CreatedAt},
						{Key: "updated_at", Value: event.UpdatedAt},
					}},
				}

				res, err := client.Events().UpdateOneByFilter(filter, update)
				if err != nil {
					log.Println("[WARN] Failed to update event", err)
					continue
				}

				if res.ModifiedCount == 1 {
					log.Println("updated event", eventInDB.EventID)
					response.EventsUpdated++
				}
			}
		}
	}

	if len(events) == 0 {
		json.NewEncoder(w).Encode(response)
		w.Header().Set("Content-Type", "application/json")
		return
	}

	if _, err = client.Events().InsertMany(events); err != nil {
		http.Error(w, "Failed to save events", http.StatusInternalServerError)
		return
	}

	response.EventsAdded = len(events)

	json.NewEncoder(w).Encode(response)
	w.Header().Set("Content-Type", "application/json")
}
