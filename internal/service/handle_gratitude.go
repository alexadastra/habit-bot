package service

import (
	"context"
	"log"

	"github.com/alexadastra/habit_bot/internal/models"
)

func (s *Service) handleGratitudeCommand(ctx context.Context, command models.UserCommand) error {
	if err := s.statesStorage.Add(
		ctx,
		command.UserID,
		models.GratitudeWaitingUserState,
	); err != nil {
		log.Printf("Error setting user state: %v", err)
		return s.sendMessage(command.UserID, stateSettingFailedErrorMessage)
	}

	return s.sendMessage(command.UserID, "What are you grateful for today?")
}

func (s *Service) handleGratitudeMessage(ctx context.Context, message models.UserMessage) error {
	if err := s.actionsStorage.AddGratitude(
		ctx,
		models.GratitudeEvent{
			UserID:    message.UserID,
			Message:   message.Message,
			CreatedAt: message.SentAt,
		},
	); err != nil {
		log.Printf("Error storing gratitude: %v", err)
		return s.sendMessage(message.UserID, gratitudeFailedErrorMessage)
	}

	if err := s.statesStorage.Add(
		ctx,
		message.UserID,
		models.DefaultUserState,
	); err != nil {
		log.Printf("Error setting user state: %v", err)
		return s.sendMessage(message.UserID, stateSettingFailedErrorMessage)
	}

	return s.sendMessage(message.UserID, "Successfully stored gratitude!")
}
