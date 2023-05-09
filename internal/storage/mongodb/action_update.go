package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func (s *Storage) UpdateActionExecution(
	ctx context.Context,
	id string,
	executedAt time.Time,
	scheduledAt time.Time,
) error {
	_, err := s.actionColl.UpdateByID(
		ctx,
		id,
		bson.M{
			"$set": bson.M{
				actionLastExecutedAtColumnName: executedAt,
				actionSchaduledAtColumnName:    scheduledAt,
			},
		},
	)

	return err
}
