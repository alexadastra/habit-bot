package migrations

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	if err := godotenv.Load("values.env"); err != nil {
		panic(err)
	}

	_ = migrate.Register(
		func(db *mongo.Database) error {
			return db.RunCommand(
				context.Background(),
				bson.D{
					primitive.E{Key: "createUser", Value: os.Getenv("MONGODB_USER")},
					primitive.E{Key: "pwd", Value: os.Getenv("MONGODB_PASSWORD")},
					primitive.E{
						Key: "roles",
						Value: []bson.M{
							{
								"role": "readWrite",
								"db":   os.Getenv("MONGODB_DATABASE"),
							},
						},
					},
				},
			).Err()
		},
		func(db *mongo.Database) error {
			return db.RunCommand(
				context.Background(),
				bson.D{
					primitive.E{Key: "dropUser", Value: "habitbotuser"},
				},
			).Err()
		},
	)
}

/*
example:
_ = migrate.Register(
	func(db *mongo.Database) error {
		opt := options.Index().SetName("my-index")
		keys := bson.D{
			{
				"my-key",
				1,
			},
		}
		model := mongo.IndexModel{Keys: keys, Options: opt}
		_, err := db.Collection("my-coll").
			Indexes().
			CreateOne(
				context.TODO(),
				model,
			)
		if err != nil {
			return err
		}

		return nil
	},
	func(db *mongo.Database) error {
		_, err := db.Collection("my-coll").
			Indexes().
			DropOne(
				context.TODO(),
				"my-index",
			)
		if err != nil {
			return err
		}
		return nil
	},
)
*/
