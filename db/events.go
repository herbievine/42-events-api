package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Event struct {
	EventID      int       `json:"event_id" bson:"event_id"`
	Name         string    `json:"name" bson:"name"`
	Description  string    `json:"description" bson:"description"`
	Location     string    `json:"location" bson:"location,omitempty"`
	Type         string    `json:"type" bson:"type,omitempty"`
	Attendees    int       `json:"attendees" bson:"attendees"`
	MaxAttendees int       `json:"max_attendees" bson:"max_attendees"`
	BeginAt      time.Time `json:"begin_at" bson:"begin_at,omitempty"`
	EndAt        time.Time `json:"end_at" bson:"end_at,omitempty"`
	CampusIDs    []int     `json:"campus_ids" bson:"campus_ids"`
	CursusIDs    []int     `json:"cursus_ids" bson:"cursus_ids"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at" bson:"updated_at,omitempty"`
}

func (coll *EventCollection) GetMany(filter *bson.D) ([]Event, error) {
	var events []Event

	cursor, err := coll.collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var event Event
		err := cursor.Decode(&event)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

func (coll *EventCollection) GetManyByCampusID(campusID int) ([]Event, error) {
	filter := bson.D{
		{Key: "campus_ids", Value: campusID},
		{Key: "begin_at", Value: bson.D{
			{Key: "$gte", Value: time.Now()},
		}},
	}

	var events []Event

	cursor, err := coll.collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var event Event
		err := cursor.Decode(&event)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

func (coll *EventCollection) GetOneByID(eventID int) (*Event, error) {
	filter := bson.D{{Key: "event_id", Value: eventID}}

	var event Event
	err := coll.collection.FindOne(context.TODO(), filter).Decode(&event)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (coll *EventCollection) InsertOne(e Event) (*mongo.InsertOneResult, error) {
	return coll.collection.InsertOne(context.TODO(), e)
}

func (coll *EventCollection) InsertMany(events []Event) (*mongo.InsertManyResult, error) {
	var docs []interface{}
	for _, event := range events {
		docs = append(docs, event)
	}

	return coll.collection.InsertMany(context.TODO(), docs)
}
