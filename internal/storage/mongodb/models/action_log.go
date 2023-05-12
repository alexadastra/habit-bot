package models

import (
	"time"

	"github.com/alexadastra/habit_bot/internal/models"
)

type ActionLog struct {
	ID               string      `bson:"_id"`
	ActionID         string      `bson:"action_id"`
	ExecutedAt       time.Time   `bson:"executed_at"`
	DurationMillisec int64       `bson:"duration"`
	Result           interface{} `bson:"result"`
}

func (a *ActionLog) ToDomain() *models.ActionLog {
	return &models.ActionLog{
		ID:               a.ID,
		ActionID:         a.ActionID,
		ExecutedAt:       a.ExecutedAt,
		DurationMillisec: a.DurationMillisec,
		Result:           a.Result,
	}
}
