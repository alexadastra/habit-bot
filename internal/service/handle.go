package service

import (
	"context"
	"fmt"

	"github.com/alexadastra/habit_bot/internal/models"
	"github.com/pkg/errors"
)

func (s *Service) HandleCommand(ctx context.Context, command models.UserCommand) error {
	switch command.Command {
	case models.Start:
		return s.sendMessage(command.UserID, welcomeMessage)
	case models.Checkin:
		return s.handleCheckinCommand(ctx, command)
	case models.Gratitude:
		return s.handleGratitudeCommand(ctx, command)
	case models.Stats:
		return s.handleStats(ctx, command)
	default:
		return s.sendMessage(command.UserID, fmt.Sprintf("Invalid command: %s", command.Command))
	}
}

func (s *Service) HandleMessage(ctx context.Context, message models.UserMessage) error {
	state, err := s.statesStorage.Get(ctx, message.UserID)
	if err != nil {
		return s.sendMessage(message.UserID, stateFetchingFailedErrorMessage)
	}

	switch state {
	case models.GratitudeWaitingUserState:
		return s.handleGratitudeMessage(ctx, message)
	default:
		return s.sendMessage(message.UserID, invalidStateErrorMessage)
	}
}

func (s *Service) sendMessage(userID int64, text string) error {
	err := s.bot.SendMessage(userID, text)
	return errors.Wrap(err, "failed to store user data")
}
