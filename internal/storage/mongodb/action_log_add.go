package mongodb

import (
	"context"

	"github.com/alexadastra/habit_bot/internal/models"
	"go.mongodb.org/mongo-driver/bson"
)

func (s *Storage) AddActionLog(
	ctx context.Context,
	actionLog models.ActionLog,
) error {
	_, err := s.actionLogColl.InsertOne(
		ctx,
		bson.M{
			actionLogIDColumnName:          actionLog.ID,
			actionLogActionLogIDColumnName: actionLog.ActionID,
			actionLogExecutedAtCollumnName: actionLog.ExecutedAt,
			actionLogDurationCollumnName:   actionLog.DurationMillisec,
			actionLogResultCollumnName:     actionLog.Result,
		},
	)

	return err
}
