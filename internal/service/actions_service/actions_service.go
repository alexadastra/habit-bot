package actions_service

import (
	"context"
	"time"

	"github.com/alexadastra/habit_bot/internal/models"
	"github.com/robfig/cron/v3"
)

type Queue interface {
	Push(ctx context.Context, id string, priority int64) error
	Pop(ctx context.Context, count int64) ([]string, error)
}

type ActionsStorage interface {
	GetActionByID(ctx context.Context, id string) (models.Action, error)
	GetAllActions(ctx context.Context) ([]models.Action, error)
	UpdateActionExecution(context.Context, string, time.Time) error

	AddActionLog(context.Context, models.ActionLog) error
}

// ActionsService represents a scheduler that uses a priority queue and an action storage
type ActionsService struct {
	queue   Queue
	storage ActionsStorage

	parser cron.Parser

	// TODO: turn these callbacks to interfaces or smth
	actions map[string]func(ctx context.Context) error
}

// New creates a new Scheduler with the given RedisQueue
func New(
	ctx context.Context,
	queue Queue,
	storage ActionsStorage,
	actions map[string]func(ctx context.Context) error,
) (
	*ActionsService,
	error,
) {
	service := &ActionsService{
		queue:   queue,
		storage: storage,
		parser:  cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor),
		actions: actions,
	}

	if err := service.startup(ctx); err != nil {
		return nil, err
	}

	return service, nil
}

func (s *ActionsService) startup(ctx context.Context) error {
	actions, err := s.storage.GetAllActions(ctx)
	if err != nil {
		return err
	}

	for _, action := range actions {
		if err := s.addAction(ctx, action); err != nil {
			return err
		}
	}

	return nil
}
