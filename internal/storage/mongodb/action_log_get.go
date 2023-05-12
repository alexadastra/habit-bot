package mongodb

import (
	"context"

	"github.com/alexadastra/habit_bot/internal/models"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	db_models "github.com/alexadastra/habit_bot/internal/storage/mongodb/models"
)

func (s *Storage) GetLastActionLog(ctx context.Context, actionID string) (*models.ActionLog, error) {
	opts := options.
		FindOne().
		SetSort(
			bson.D{
				{
					Key:   checkinTimestampCollumnName,
					Value: 1,
				},
			},
		)

	res := s.actionLogColl.FindOne(
		ctx,
		bson.M{
			actionLogActionLogIDColumnName: actionID,
		},
		opts,
	)

	if err := res.Err(); err != nil {
		return nil, errors.Wrap(err, "failed to get last action log")
	}

	actionLog := db_models.ActionLog{}
	if err := res.Decode(actionLog); err != nil {
		return nil, errors.Wrap(err, "failed to scan result into log")
	}

	return actionLog.ToDomain(), nil
}
