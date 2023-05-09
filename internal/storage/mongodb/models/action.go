package models

import (
	"time"

	"github.com/alexadastra/habit_bot/internal/models"
)

type Action struct {
	ID             string    `bson:"_id"`
	Name           string    `bson:"name"`
	Priority       int64     `bson:"priority"`
	ScheduledAt    time.Time `bson:"scheduled_at"`
	IsCancelled    bool      `bson:"is_cancelled"`
	LastExecutedAt time.Time `bson:"last_executed_at,omitempty"`
}

func NewAction(a models.Action) Action {
	return Action{
		ID:             a.ID,
		Name:           a.Name,
		Priority:       a.Priority,
		ScheduledAt:    a.ScheduledAt,
		IsCancelled:    a.IsCancelled,
		LastExecutedAt: a.LastExecutedAt,
	}
}

func (a *Action) ToDomain() models.Action {
	return models.Action{
		ID:             a.ID,
		Name:           a.Name,
		Priority:       a.Priority,
		ScheduledAt:    a.ScheduledAt,
		IsCancelled:    a.IsCancelled,
		LastExecutedAt: a.LastExecutedAt,
	}
}
