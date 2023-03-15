package mongodb

import (
	"time"

	"github.com/alexadastra/habit_bot/internal/models"
)

type CheckinEvent struct {
	UserID    int64     `bson:"user_id"`
	Timestamp time.Time `bson:"timestamp"`
}

type GratitudeEvent struct {
	UserID    int64     `bson:"user_id"`
	Text      string    `bson:"text"`
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

func NewGratitudeEvent(event models.GratitudeEvent) GratitudeEvent {
	return GratitudeEvent{
		UserID:    event.UserID,
		Text:      event.Message,
		Timestamp: event.CreatedAt,
	}
}

func (event *GratitudeEvent) ToDomain() models.GratitudeEvent {
	return models.GratitudeEvent{
		UserID:    event.UserID,
		Message:   event.Text,
		CreatedAt: event.Timestamp,
	}
}
