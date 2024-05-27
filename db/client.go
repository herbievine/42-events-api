package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const uri = "mongodb+srv://me:11mfgj59bmNh0rh5@42-events.clgnegw.mongodb.net/?retryWrites=true&w=majority&appName=42-events"

type Client struct {
	client *mongo.Client
}

type UserCollection struct {
	collection *mongo.Collection
}

type EventCollection struct {
	collection *mongo.Collection
}

type CampusCollection struct {
	collection *mongo.Collection
}

type NotificationCollection struct {
	collection *mongo.Collection
}

func NewClient() (*Client, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	ctx := context.TODO()

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return &Client{client: client}, nil
}

func (c *Client) Close(ctx context.Context) error {
	return c.client.Disconnect(ctx)
}

func (c *Client) Users() *UserCollection {
	return &UserCollection{c.client.Database("42-events").Collection("users")}
}

func (c *Client) Events() *EventCollection {
	return &EventCollection{c.client.Database("42-events").Collection("events")}
}

func (c *Client) Campus() *CampusCollection {
	return &CampusCollection{c.client.Database("42-events").Collection("campus")}
}

func (c *Client) Notifications() *NotificationCollection {
	return &NotificationCollection{c.client.Database("42-events").Collection("notifications")}
}
