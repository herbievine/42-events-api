package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Campus struct {
	CampusID  int    `json:"campus_id" bson:"campus_id"`
	Name      string `json:"name" bson:"name"`
	UserCount int    `json:"user_count" bson:"user_count"`
	City      string `json:"city" bson:"city"`
	Country   string `json:"country" bson:"country"`
}

func (coll *CampusCollection) GetMany() ([]Campus, error) {
	var campuses []Campus

	cursor, err := coll.collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var campus Campus
		err := cursor.Decode(&campus)
		if err != nil {
			return nil, err
		}

		campuses = append(campuses, campus)
	}

	return campuses, nil
}

func (coll *CampusCollection) GetOneByID(campusID int) (*Campus, error) {
	filter := bson.D{{Key: "campus_id", Value: campusID}}

	var campus Campus
	err := coll.collection.FindOne(context.TODO(), filter).Decode(&campus)
	if err != nil {
		return nil, err
	}

	return &campus, nil
}

func (coll *CampusCollection) InsertOne(c Campus) (*mongo.InsertOneResult, error) {
	return coll.collection.InsertOne(context.TODO(), c)
}

func (coll *CampusCollection) InsertMany(campuses []Campus) (*mongo.InsertManyResult, error) {
	var docs []interface{}
	for _, campus := range campuses {
		docs = append(docs, campus)
	}

	return coll.collection.InsertMany(context.TODO(), docs)
}
