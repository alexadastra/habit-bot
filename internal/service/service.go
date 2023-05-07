package service

import (
	"context"

	"github.com/alexadastra/habit_bot/internal/external/bot"
	"github.com/alexadastra/habit_bot/internal/models"
)

type StatsFetcher interface {
	FetchStats(context.Context, int64) (int64, int64, error)
}

type UserActionsStorage interface {
	AddCheckin(context.Context, models.CheckinEvent) error
	AddGratitude(context.Context, models.GratitudeEvent) error
}

type UserStatesStorage interface {
	Get(context.Context, int64) (models.UserState, error)
	Add(context.Context, int64, models.UserState) error
}

type Service struct {
	bot            *bot.Bot
	actionsStorage UserActionsStorage
	statesStorage  UserStatesStorage

	statsFetcher StatsFetcher
}

func NewService(
	bot *bot.Bot,
	actionsStorage UserActionsStorage,
	statesStorage UserStatesStorage,
	statsFetcher StatsFetcher,
) *Service {
	return &Service{
		bot:            bot,
		actionsStorage: actionsStorage,
		statesStorage:  statesStorage,
		statsFetcher:   statsFetcher,
	}
}
