package migrations

import (
	"context"
	"time"

	"github.com/google/uuid"
	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	_ = migrate.Register(
		func(db *mongo.Database) error {
			/*
				opt := options.Index().SetName("my-index")
				keys := bson.D{
					{
						"my-key",
						1,
					},
				}
				model := mongo.IndexModel{Keys: keys, Options: opt}
			*/
			err := db.CreateCollection(
				context.TODO(),
				"action",
			)
			if err != nil {
				return err
			}

			id, err := uuid.NewUUID()
			if err != nil {
				return err
			}

			_, err = db.Collection("action").InsertOne(
				context.TODO(),
				bson.M{
					"_id":              id.String(),
					"name":             "pochesatso",
					"priority":         0,
					"scheduled_at":     time.Now().UTC(),
					"is_cancelled":     false,
					"last_executed_at": nil,
				},
			)
			if err != nil {
				return err
			}

			return nil
		},
		func(db *mongo.Database) error {
			return db.Collection("action").Drop(context.TODO())
		},
	)
}
