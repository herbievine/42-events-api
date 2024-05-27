package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Notification struct {
	UserID    int       `json:"user_id" bson:"user_id"`
	EventID   int       `json:"event_id" bson:"event_id"`
	HasRead   bool      `json:"has_read" bson:"has_read"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	DeletedAt time.Time `json:"deleted_at" bson:"deleted_at,omitempty"`
}

func (coll *NotificationCollection) GetMany(filter bson.D) ([]Notification, error) {
	var notifications []Notification

	cursor, err := coll.collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var Notification Notification
		err := cursor.Decode(&Notification)
		if err != nil {
			return nil, err
		}

		notifications = append(notifications, Notification)
	}

	return notifications, nil
}

func (coll *NotificationCollection) GetManyByCampusID(campusID int) ([]Notification, error) {
	filter := bson.D{
		{Key: "campus_ids", Value: campusID},
		{Key: "begin_at", Value: bson.D{
			{Key: "$gte", Value: time.Now()},
		}},
	}

	var notifications []Notification

	cursor, err := coll.collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var Notification Notification
		err := cursor.Decode(&Notification)
		if err != nil {
			return nil, err
		}

		notifications = append(notifications, Notification)
	}

	return notifications, nil
}

func (coll *NotificationCollection) UpdateOneByFilter(filter primitive.D, update primitive.D) (*mongo.UpdateResult, error) {
	return coll.collection.UpdateOne(context.TODO(), filter, update)
}

func (coll *NotificationCollection) GetOneByFilter(filter primitive.D) (*Notification, error) {
	var Notification Notification
	err := coll.collection.FindOne(context.TODO(), filter).Decode(&Notification)
	if err != nil {
		return nil, err
	}

	return &Notification, nil
}

func (coll *NotificationCollection) InsertOne(e Notification) (*mongo.InsertOneResult, error) {
	return coll.collection.InsertOne(context.TODO(), e)
}

func (coll *NotificationCollection) InsertMany(notifications []Notification) (*mongo.InsertManyResult, error) {
	var docs []interface{}
	for _, Notification := range notifications {
		docs = append(docs, Notification)
	}

	return coll.collection.InsertMany(context.TODO(), docs)
}
