package mongodb

import (
	"context"

	"github.com/alexadastra/habit_bot/internal/models"
	"go.mongodb.org/mongo-driver/bson"
)

func (s *Storage) AddGratitude(ctx context.Context, gratitudeMessage models.GratitudeEvent) error {
	_, err := s.gratitudeColl.InsertOne(
		ctx,
		bson.M{
			gratitudeUserIDCollumnName:    gratitudeMessage.UserID,
			gratitudeTextCollumnName:      gratitudeMessage.Message,
			gratitudeTimestampCollumnName: gratitudeMessage.CreatedAt,
		},
	)
	return err
}
