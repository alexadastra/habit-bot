package models

import (
	"time"

	"github.com/alexadastra/habit_bot/internal/models"
)

type CheckinEvent struct {
	UserID    int64     `bson:"user_id"`
	Timestamp time.Time `bson:"timestamp"`
}

func NewCheckinEvent(event models.CheckinEvent) CheckinEvent {
	return CheckinEvent{
		UserID:    event.UserID,
		Timestamp: event.CreatedAt,
	}
}

func (event CheckinEvent) ToDomain() models.CheckinEvent {
	return models.CheckinEvent{
		UserID:    event.UserID,
		CreatedAt: event.Timestamp,
	}
}
