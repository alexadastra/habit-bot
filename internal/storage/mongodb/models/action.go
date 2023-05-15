package models

import (
	"time"

	"github.com/alexadastra/habit_bot/internal/models"
)

type Action struct {
	ID           string    `bson:"_id"`
	Name         string    `bson:"name"`
	Priority     int64     `bson:"priority"`
	IsCancelled  bool      `bson:"is_cancelled"`
	Crontab      string    `bson:"crontab"`
	IsRepeatable bool      `bson:"is_repeatable"`
	ScheduledAt  time.Time `bson:"scheduled_at"`
}

func NewAction(a models.Action) Action {
	return Action{
		ID:           a.ID,
		Name:         a.Name,
		Priority:     a.Priority,
		IsCancelled:  a.IsCancelled,
		Crontab:      a.Crontab,
		IsRepeatable: a.IsRepeatable,
		ScheduledAt:  a.ScheduledAt,
	}
}

func (a *Action) ToDomain() models.Action {
	return models.Action{
		ID:          a.ID,
		IsCancelled: a.IsCancelled,

		ActionExecutionerProperties: models.ActionExecutionerProperties{
			Name:     a.Name,
			Priority: a.Priority,
		},

		ActionSchedulerProperties: models.ActionSchedulerProperties{
			Crontab:      a.Crontab,
			IsRepeatable: a.IsRepeatable,
			ScheduledAt:  a.ScheduledAt,
		},
	}
}
