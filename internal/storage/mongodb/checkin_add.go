package mongodb

import (
	"context"

	"github.com/alexadastra/habit_bot/internal/models"
	"go.mongodb.org/mongo-driver/bson"
)

func (s *Storage) AddCheckin(ctx context.Context, checkinMessage models.CheckinEvent) error {
	_, err := s.checkinColl.InsertOne(
		ctx,
		bson.M{
			checkinUserIDCollumnName:    checkinMessage.UserID,
			checkinTimestampCollumnName: checkinMessage.CreatedAt,
		},
	)
	return err
}
