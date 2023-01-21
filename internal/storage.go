package internal

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage struct {
	client        *mongo.Client
	usersColl     *mongo.Collection
	gratitudeColl *mongo.Collection
}

func NewStorage(dsn string) (*Storage, error) {
	// Set up MongoDB client
	client, err := mongo.NewClient(options.Client().ApplyURI(dsn))
	if err != nil {
		return nil, err
	}
	err = client.Connect(context.TODO())
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	usersColl := client.Database("habit-bot").Collection("users")
	gratitudeColl := client.Database("habit-bot").Collection("gratitude")

	return &Storage{
		client:        client,
		usersColl:     usersColl,
		gratitudeColl: gratitudeColl,
	}, err
}

func (s *Storage) storeUserData(userID int, timestamp time.Time) error {
	_, err := s.usersColl.InsertOne(context.TODO(), bson.M{"user_id": userID, "timestamp": timestamp})
	return err
}

func (s *Storage) storeGratitude(userID int, text string) error {
	_, err := s.gratitudeColl.InsertOne(context.TODO(), bson.M{"user_id": userID, "text": text, "timestamp": time.Now()})
	return err
}
