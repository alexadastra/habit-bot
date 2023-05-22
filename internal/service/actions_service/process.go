package actions_service

import (
	"context"
	"log"
	"time"

	"github.com/alexadastra/habit_bot/internal/models"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (s *ActionsService) Process(ctx context.Context) {
	startedAt := time.Now().UTC()

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

	timeToExecute := time.Now().UTC().After(action.ScheduledAt)

	// if the execution time has not yet arrived, put the action back in the queue
	if !timeToExecute {
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

	err = s.execute(ctx, action)
	if err != nil {
		log.Println("Failed to execute action:", err)
		// TODO: retry?
	}

	doneAt := time.Now().UTC()

	if err := s.storage.AddActionLog(
		ctx,
		models.ActionLog{
			ID:               uuid.NewString(),
			ActionID:         action.ID,
			ExecutedAt:       doneAt,
			DurationMillisec: doneAt.Sub(startedAt).Milliseconds(),
			// TODO: remove hardcode
			Result: map[string]string{
				"status": "done",
			},
		},
	); err != nil {
		log.Println("Failed to add action log:", err)
	}

	if !action.IsRepeatable {
		return
	}

	// reset scheduledAt to the new execution time
	action.ScheduledAt, err = s.getNextExecutionTime(action)
	if err != nil {
		log.Println("Failed to reset action scheduled time:", err)
		return
	}

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

func (s *ActionsService) getNextExecutionTime(a models.Action) (time.Time, error) {
	cronExpr, err := s.parser.Parse(a.Crontab)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "failed to parse crontab")
	}

	return cronExpr.Next(time.Now().UTC()), nil
}
