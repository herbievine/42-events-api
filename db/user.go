package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	UserID          int       `json:"user_id" bson:"user_id"`
	Login           string    `json:"login" bson:"login"`
	ImageURL        string    `json:"image_url" bson:"image_url"`
	CampusIDs       []int     `json:"campus_ids" bson:"campus_ids"`
	PrimaryCampusID int       `json:"primary_campus_id" bson:"primary_campus_id,omitempty"`
	LastSeen        time.Time `json:"last_seen" bson:"last_seen"`
	CreatedAt       time.Time `json:"created_at" bson:"created_at"`
}

func (coll *UserCollection) GetMany() ([]User, error) {
	var users []User

	cursor, err := coll.collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var user User
		err := cursor.Decode(&user)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (coll *UserCollection) GetManyByCampusID(campusID int) ([]User, error) {
	filter := bson.D{{Key: "campus_ids", Value: campusID}}

	var users []User

	cursor, err := coll.collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var user User
		err := cursor.Decode(&user)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (coll *UserCollection) GetOneByID(userID int) (*User, error) {
	filter := bson.D{{Key: "user_id", Value: userID}}

	var user User
	err := coll.collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (coll *UserCollection) InsertOne(u User) (*mongo.InsertOneResult, error) {
	return coll.collection.InsertOne(context.TODO(), u)
}

func (coll *UserCollection) InsertMany(users []User) (*mongo.InsertManyResult, error) {
	var docs []interface{}
	for _, user := range users {
		docs = append(docs, user)
	}

	return coll.collection.InsertMany(context.TODO(), docs)
}
