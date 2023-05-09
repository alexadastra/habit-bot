package mongodb

import (
	"context"

	"github.com/alexadastra/habit_bot/internal/models"
	db_models "github.com/alexadastra/habit_bot/internal/storage/mongodb/models"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
)

func (s *Storage) GetActionByID(ctx context.Context, id string) (models.Action, error) {
	var action db_models.Action
	if err := s.actionColl.
		FindOne(
			ctx,
			bson.M{
				"_id": id,
			},
		).
		Decode(&action); err != nil {
		return models.Action{}, err
	}

	return action.ToDomain(), nil
}

func (s *Storage) GetAllActions(ctx context.Context) ([]models.Action, error) {
	var actions []db_models.Action
	cursor, err := s.actionColl.
		Find(
			ctx,
			bson.D{},
		)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &actions); err != nil {
		return nil, err
	}

	return lo.Map(
		actions,
		func(action db_models.Action, idx int) models.Action {
			return action.ToDomain()
		},
	), nil
}
