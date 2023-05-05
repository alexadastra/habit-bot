package service

import (
	"context"
	"log"

	"github.com/alexadastra/habit_bot/internal/models"
)

// the case where something went wrong and we should notify the user about that
// should work differently. maybe with some more user-friendly error set
func (s *Service) handleCheckinCommand(ctx context.Context, command models.UserCommand) error {
	if err := s.actionsStorage.AddCheckin(
		ctx,
		models.CheckinEvent{
			UserID:    command.UserID,
			CreatedAt: command.SentAt,
		},
	); err != nil {
		log.Printf("Error storing checkin: %v", err)
		return s.sendMessage(command.UserID, checkinFailedErrorMessage)
	}

	return s.sendMessage(command.UserID, "Successfully checked in!")
}
