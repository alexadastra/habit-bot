package mongodb

import (
	"context"
	"time"

	"github.com/alexadastra/habit_bot/internal/models"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *Storage) GetGratitudeEvents(
	ctx context.Context,
	userID int64,
	from time.Time,
	to time.Time,
) (
	[]models.GratitudeEvent,
	error,
) {
	filter := bson.D{
		{
			Key: "$and",
			Value: bson.A{
				bson.D{
					{
						Key: gratitudeUserIDCollumnName,
						Value: bson.D{
							{Key: "$eq", Value: userID},
						},
					},
					{
						Key: gratitudeTimestampCollumnName,
						Value: bson.D{
							{Key: "$gte", Value: from},
						},
					},
					{
						Key: gratitudeTimestampCollumnName,
						Value: bson.D{
							{Key: "$lte", Value: to},
						},
					},
				},
			},
		},
	}

	sort := bson.D{{Key: gratitudeTimestampCollumnName, Value: 1}}
	opts := options.Find().SetSort(sort)

	cursor, err := s.gratitudeColl.Find(
		ctx,
		filter,
		opts,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get events from DB")
	}

	var result []GratitudeEvent
	if err = cursor.All(ctx, &result); err != nil {
		return nil, errors.Wrap(err, "failed to scan events into result")
	}

	return lo.Map(result, func(e GratitudeEvent, idx int) models.GratitudeEvent { return e.ToDomain() }), nil
}
