package mongodb

import (
	"context"

	"github.com/alexadastra/habit_bot/internal/models"
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

func (s *Storage) StoreCheckin(ctx context.Context, checkinMessage models.UserMessage) error {
	_, err := s.usersColl.InsertOne(
		ctx,
		bson.M{
			"user_id":   checkinMessage.UserID,
			"timestamp": checkinMessage.SentAt,
		},
	)
	return err
}

func (s *Storage) StoreGratitude(ctx context.Context, gratitudeMessage models.UserMessage) error {
	_, err := s.gratitudeColl.InsertOne(
		ctx,
		bson.M{
			"user_id":   gratitudeMessage.UserID,
			"text":      gratitudeMessage.Message,
			"timestamp": gratitudeMessage.SentAt,
		},
	)
	return err
}
