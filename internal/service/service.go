package service

import (
	"context"
	"time"

	"github.com/alexadastra/habit_bot/internal/external/bot"
	"github.com/alexadastra/habit_bot/internal/models"
)

type UserActionsStorage interface {
	AddCheckin(context.Context, models.CheckinEvent) error
	AddGratitude(context.Context, models.GratitudeEvent) error
	GetCheckinEvents(context.Context, int64, time.Time, time.Time) ([]models.CheckinEvent, error)
	GetGratitudeEvents(context.Context, int64, time.Time, time.Time) ([]models.GratitudeEvent, error)
}

type UserStatesStorage interface {
	Get(context.Context, int64) (models.UserState, error)
	Add(context.Context, int64, models.UserState) error
}

type Service struct {
	bot            *bot.Bot
	actionsStorage UserActionsStorage
	statesStorage  UserStatesStorage
}

func NewService(bot *bot.Bot, actionsStorage UserActionsStorage, statesStorage UserStatesStorage) *Service {
	return &Service{
		bot:            bot,
		actionsStorage: actionsStorage,
		statesStorage:  statesStorage,
	}
}
