package migrations

import (
	"context"
	"time"

	"github.com/google/uuid"
	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
					"_id":            id.String(),
					"name":           "pochesatso",
					"priority":       0,
					"is_cancelled":   false,
					"crontab":        "* * * * * *",
					"is_repeatable": true,
					"scheduled_at":   time.Now().UTC(),
				},
			)
			if err != nil {
				return err
			}

			idxName := "action_log_action_id_index"
			f := false

			_, err = db.Collection("action_log").Indexes().CreateOne(
				context.TODO(),
				mongo.IndexModel{
					Keys: bson.M{
						"action_id": 1,
					},
					Options: &options.IndexOptions{
						Name:   &idxName,
						Unique: &f,
					},
				},
			)
			return err
		},
		func(db *mongo.Database) error {
			return db.Collection("action").Drop(context.TODO())
		},
	)
}
