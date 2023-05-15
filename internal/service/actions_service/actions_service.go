package actions_service

import (
	"context"
	"log"
	"time"

	"github.com/alexadastra/habit_bot/internal/models"
)

type Queue interface {
	Push(ctx context.Context, id string, priority int64) error
	Pop(ctx context.Context, count int64) ([]string, error)
}

type ActionsStorage interface {
	GetActionByID(ctx context.Context, id string) (models.Action, error)
	GetAllActions(ctx context.Context) ([]models.Action, error)
	UpdateActionExecution(context.Context, string, time.Time, time.Time) error
}

// ActionsService represents a scheduler that uses a priority queue and an action storage
type ActionsService struct {
	queue   Queue
	storage ActionsStorage
}

// New creates a new Scheduler with the given RedisQueue
func New(ctx context.Context, queue Queue, storage ActionsStorage) (*ActionsService, error) {
	service := &ActionsService{
		queue:   queue,
		storage: storage,
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

func (s *ActionsService) Process(ctx context.Context) {
	// get the next action from the queue
	actionIDs, err := s.queue.Pop(ctx, 1)
	if err != nil {
		log.Println("Failed to get action id from queue:", err)
		return
	}

	if len(actionIDs) == 0 {
		return
	}

	action, err := s.storage.GetActionByID(ctx, actionIDs[0])
	if err != nil {
		log.Println("Failed to get action from db:", err)
		return
	}

	// drop execution if the action has been cancelled
	if action.IsCancelled {
		return
	}

	// get the new execution time of the action
	now := time.Now().UTC()

	scheduledAt := action.GetNextExecutionTime()

	// if the execution time has not yet arrived, put the action back in the queue
	if (action.LastExecutedAt != time.Time{}) && now.Before(scheduledAt) {
		err = s.addAction(ctx, action)
		if err != nil {
			log.Println("Failed to push action back to queue:", err)
		}
		/*
			when dealing with the latest action, that will happen in few hours, we might add
			time.Sleep(time.Until(action.ScheduledAt))
			for so we won't pull and push the same action every second. but we won't, because
			if some kind of newer event appear in the queue, the scheduler woun't know about it
			and stuck with the older one
		*/
		return
	}

	err = s.execute(action)
	if err != nil {
		log.Println("Failed to execute action:", err)
		// TODO: retry?
	}

	if !action.IsRepeatable {
		return
	}

	action.LastExecutedAt = now

	if err := s.storage.UpdateActionExecution(ctx, action.ID, action.ScheduledAt); err != nil {
		log.Println("Failed to update action scheduled time", err)
		return
	}

	if err := s.addAction(ctx, action); err != nil {
		log.Println("Failed to add action:", err)
	}
}

// addAction adds an action to the queue
func (s *ActionsService) addAction(ctx context.Context, action models.Action) error {
	return s.queue.Push(ctx, action.ID, getPriority(action))
}

// getPriority returns the priority the action id will be pushed to the queue
// currently we pull by the least scheduling time,
// but later may add prioritisation for some types of actions
func getPriority(action models.Action) int64 {
	return action.ScheduledAt.Unix()
}

func (s *ActionsService) execute(action models.Action) error {
	// TODO: actually do something
	log.Printf("doing action '%s' with the title '%s'!", action.ID, action.Name)

	return nil
}
