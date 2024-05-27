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
	"go.mongodb.org/mongo-driver/bson"
)

type NotificationWithEvent struct {
	db.Event
	HasRead bool `json:"has_read"`
}

func GetNotifications(w http.ResponseWriter, r *http.Request, client *db.Client) {
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

	filter := bson.D{{Key: "user_id", Value: me.ID}}

	notifications, err := client.Notifications().GetMany(filter)
	if err != nil {
		http.Error(w, "Failed to get notifications", http.StatusInternalServerError)
		return
	}

	eventIDs := make([]int, 0, len(notifications))
	for _, notification := range notifications {
		eventIDs = append(eventIDs, notification.EventID)
	}

	filter = bson.D{
		{Key: "event_id", Value: bson.D{
			{Key: "$in", Value: eventIDs},
		}},
		{Key: "begin_at", Value: bson.D{
			{Key: "$gte", Value: time.Now()},
		}},
	}

	events, err := client.Events().GetMany(&filter)
	if err != nil {
		http.Error(w, "Failed to get events", http.StatusInternalServerError)
		return
	}

	eventMap := make(map[int]db.Event)
	for _, event := range events {
		eventMap[event.EventID] = event
	}

	notificationsWithEvents := make([]NotificationWithEvent, 0, len(notifications))
	for _, notification := range notifications {
		if _, ok := eventMap[notification.EventID]; ok {
			notificationsWithEvents = append(notificationsWithEvents, NotificationWithEvent{
				Event:   eventMap[notification.EventID],
				HasRead: notification.HasRead,
			})
		}
	}

	query := r.URL.Query()

	if query.Get("sort") == "created_at" {
		sort.Slice(notificationsWithEvents, func(i, j int) bool {
			return notificationsWithEvents[i].Event.CreatedAt.After(notificationsWithEvents[j].Event.CreatedAt)
		})
	} else if query.Get("sort") == "begin_at" {
		sort.Slice(notificationsWithEvents, func(i, j int) bool {
			return notificationsWithEvents[i].Event.BeginAt.Before(notificationsWithEvents[j].Event.BeginAt)
		})
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(notificationsWithEvents)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func GetOldNotifications(w http.ResponseWriter, r *http.Request, client *db.Client) {
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

	filter := bson.D{{Key: "user_id", Value: me.ID}, {Key: "has_read", Value: true}}

	notifications, err := client.Notifications().GetMany(filter)
	if err != nil {
		http.Error(w, "Failed to get notifications", http.StatusInternalServerError)
		return
	}

	eventIDs := make([]int, 0, len(notifications))
	for _, notification := range notifications {
		eventIDs = append(eventIDs, notification.EventID)
	}

	filter = bson.D{
		{Key: "event_id", Value: bson.D{
			{Key: "$in", Value: eventIDs},
		}},
		{Key: "begin_at", Value: bson.D{
			{Key: "$gte", Value: time.Now()},
		}},
	}

	events, err := client.Events().GetMany(&filter)
	if err != nil {
		http.Error(w, "Failed to get events", http.StatusInternalServerError)
		return
	}

	eventMap := make(map[int]db.Event)
	for _, event := range events {
		eventMap[event.EventID] = event
	}

	notificationsWithEvents := make([]NotificationWithEvent, 0, len(notifications))
	for _, notification := range notifications {
		if _, ok := eventMap[notification.EventID]; ok {
			notificationsWithEvents = append(notificationsWithEvents, NotificationWithEvent{
				Event:   eventMap[notification.EventID],
				HasRead: notification.HasRead,
			})
		}
	}

	query := r.URL.Query()

	if query.Get("sort") == "created_at" {
		sort.Slice(notificationsWithEvents, func(i, j int) bool {
			return notificationsWithEvents[i].Event.CreatedAt.After(notificationsWithEvents[j].Event.CreatedAt)
		})
	} else if query.Get("sort") == "begin_at" {
		sort.Slice(notificationsWithEvents, func(i, j int) bool {
			return notificationsWithEvents[i].Event.BeginAt.Before(notificationsWithEvents[j].Event.BeginAt)
		})
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(notificationsWithEvents)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func ReadNotification(w http.ResponseWriter, r *http.Request, client *db.Client) {
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

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	filter := bson.D{
		{Key: "event_id", Value: id},
		{Key: "user_id", Value: me.ID},
	}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "has_read", Value: true}}}}

	_, err = client.Notifications().UpdateOneByFilter(filter, update)
	if err != nil {
		log.Println("[WARN] Failed to update notification", err)

		http.Error(w, "Failed to update notification", http.StatusInternalServerError)
		return
	}

	io.WriteString(w, "ok")
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusNoContent)
}
