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

func (s *Storage) GetCheckinEvents(
	ctx context.Context,
	userID int64,
	from time.Time,
	to time.Time,
) (
	[]models.CheckinEvent,
	error,
) {
	filter := bson.D{
		{
			Key: "$and",
			Value: bson.A{
				bson.D{
					{
						Key: checkinUserIDCollumnName,
						Value: bson.D{
							{Key: "$eq", Value: userID},
						},
					},
					{
						Key: checkinTimestampCollumnName,
						Value: bson.D{
							{Key: "$gte", Value: from},
						},
					},
					{
						Key: checkinTimestampCollumnName,
						Value: bson.D{
							{Key: "$lte", Value: to},
						},
					},
				},
			},
		},
	}

	sort := bson.D{{Key: checkinTimestampCollumnName, Value: 1}}
	opts := options.Find().SetSort(sort)

	cursor, err := s.checkinColl.Find(
		ctx,
		filter,
		opts,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get events from DB")
	}

	var result []CheckinEvent
	if err = cursor.All(ctx, &result); err != nil {
		return nil, errors.Wrap(err, "failed to scan events into result")
	}

	return lo.Map(result, func(e CheckinEvent, idx int) models.CheckinEvent { return e.ToDomain() }), nil
}
