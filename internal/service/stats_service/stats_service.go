package stats_service

import (
	"context"
	"time"

	"github.com/alexadastra/habit_bot/internal/models"
)

type EventsStorage interface {
	GetCheckinEvents(context.Context, int64, time.Time, time.Time) ([]models.CheckinEvent, error)
	GetGratitudeEvents(context.Context, int64, time.Time, time.Time) ([]models.GratitudeEvent, error)
}

type StatsService struct {
	storage EventsStorage
}

func NewStatsService(storage EventsStorage) *StatsService {
	return &StatsService{storage: storage}
}
