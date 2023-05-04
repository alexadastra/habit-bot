package mongodb

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	databaseName = "habit-bot"

	checkinCollectionName       = "checkin"
	checkinUserIDCollumnName    = "user_id"
	checkinTimestampCollumnName = "timestamp"

	gratitudeCollectionName       = "gratitude"
	gratitudeUserIDCollumnName    = "user_id"
	gratitudeTextCollumnName      = "text"
	gratitudeTimestampCollumnName = "timestamp"
)

type Storage struct {
	client        *mongo.Client
	checkinColl   *mongo.Collection
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

	checkinColl := client.Database(databaseName).Collection(checkinCollectionName)
	gratitudeColl := client.Database(databaseName).Collection(gratitudeCollectionName)

	return &Storage{
		client:        client,
		checkinColl:   checkinColl,
		gratitudeColl: gratitudeColl,
	}, err
}
