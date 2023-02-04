package internal

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage struct {
	client        *mongo.Client
	usersColl     *mongo.Collection
	gratitudeColl *mongo.Collection
}

func NewStorage(ctx context.Context, dsn string) (*Storage, error) {
	// Set up MongoDB client
	client, err := mongo.NewClient(options.Client().ApplyURI(dsn))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create client")
	}
	err = client.Connect(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to db")
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to ping")
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
