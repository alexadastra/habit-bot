package models

import (
	"time"

	"github.com/alexadastra/habit_bot/internal/models"
)

type GratitudeEvent struct {
	UserID    int64     `bson:"user_id"`
	Text      string    `bson:"text"`
	Timestamp time.Time `bson:"timestamp"`
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
