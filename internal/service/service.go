package service

import (
	"context"
	"fmt"
	"log"

	"github.com/alexadastra/habit_bot/internal/external/bot"
	"github.com/alexadastra/habit_bot/internal/models"
	"github.com/pkg/errors"
)

const (
	// internal errors
	checkinFailedErrorMessage       = "Error while saving 'checkin' action. Please try again."
	gratitudeFailedErrorMessage     = "Error while saving 'gratitude' action. Please try again."
	stateFetchingFailedErrorMessage = "Error while fetching the user state. Please try again."
	stateSettingFailedErrorMessage  = "Error while setting the user state. Please try again."

	// user-ish errors
	invalidStateErrorMessage = "Invalid state. Please try again."

	// successful messages
	welcomeMessage = "Welcome to Habit Bot by @alex_ad_astra!"
)

type UserActionsStorage interface {
	AddCheckin(context.Context, models.UserMessage) error
	AddGratitude(context.Context, models.UserMessage) error
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

func (s *Service) HandleCommand(command models.UserCommand) error {
	switch command.Command {
	case models.Start:
		return s.sendMessage(command.UserID, welcomeMessage)
	case models.Checkin:
		// the case where something went wrong and we should notify the user about that
		// should work differently. maybe with some more user-friendly error set
		if err := s.actionsStorage.StoreCheckin(context.Background(), command.UserMessage); err != nil {
			log.Printf("Error storing checkin: %v", err)
			return s.sendMessage(command.UserID, checkinFailedErrorMessage)
		}

		return s.sendMessage(command.UserID, "Successfully checked in!")
	case models.Gratitude:
		if err := s.statesStorage.SetByID(
			context.Background(),
			command.UserID,
			models.GratitudeWaitingUserState,
		); err != nil {
			log.Printf("Error setting user state: %v", err)
			return s.sendMessage(command.UserID, stateSettingFailedErrorMessage)
		}

		return s.sendMessage(command.UserID, "What are you grateful for today?")
	default:
		return s.sendMessage(command.UserID, fmt.Sprintf("Invalid command: %s", command.Command))
	}
}

func (s *Service) HandleMessage(message models.UserMessage) error {
	state, err := s.statesStorage.FetchByID(context.Background(), message.UserID)
	if err != nil {
		return s.sendMessage(message.UserID, stateFetchingFailedErrorMessage)
	}

	switch state {
	case models.GratitudeWaitingUserState:
		if err := s.actionsStorage.StoreGratitude(context.Background(), message); err != nil {
			log.Printf("Error storing gratitude: %v", err)
			return s.sendMessage(message.UserID, gratitudeFailedErrorMessage)
		}

		if err := s.statesStorage.SetByID(
			context.Background(),
			message.UserID,
			models.DefaultUserState,
		); err != nil {
			log.Printf("Error setting user state: %v", err)
			return s.sendMessage(message.UserID, stateSettingFailedErrorMessage)
		}

		return s.sendMessage(message.UserID, "Successfully stored gratitude!")
	default:
		return s.sendMessage(message.UserID, invalidStateErrorMessage)
	}
}

func (s *Service) sendMessage(userID int64, text string) error {
	err := s.bot.SendMessage(userID, text)
	return errors.Wrap(err, "failed to store user data")
}
