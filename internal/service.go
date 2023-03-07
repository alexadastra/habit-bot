package internal

import (
	"fmt"
	"log"
	"time"

	"github.com/alexadastra/habit_bot/internal/models"
	"github.com/pkg/errors"
)

const (
	storageErrorMessage         = "Error storing user data. Please try again."
	invalidStateErrorMessage    = "Invalid state. Please try again."
	gratitudeFailedErrorMessage = "Error storing gratitude. Please try again."
)

type Service struct {
	bot     *Bot
	storage *Storage
	states  map[int64]string // TODO: replace with cache?
}

func NewService(bot *Bot, storage *Storage) *Service {
	return &Service{
		bot:     bot,
		storage: storage,
		states:  make(map[int64]string),
	}
}

func (s *Service) handleCommand(command models.UserCommand) error {
	switch command.Command {
	case models.Checkin:
		// the case where something went wrong and we should notify the user about that
		// should work differently. maybe with some more user-friendly error set
		if err := s.storage.storeUserData(command.UserID, time.Now()); err != nil {
			log.Printf("Error storing checkin: %v", err)
			return s.sendMessage(command.UserID, storageErrorMessage)
		}

		return s.sendMessage(command.UserID, "Successfully checked in!")
	case models.Gratitude:
		s.states[command.UserID] = "gratitude"

		return s.sendMessage(command.UserID, "What are you grateful for today?")
	default:
		return s.sendMessage(command.UserID, fmt.Sprintf("Invalid command: %s", command.Command))
	}
}

func (s *Service) handleMessage(message models.UserMessage) error {
	state, ok := s.states[message.UserID]
	if !ok {
		return s.sendMessage(message.UserID, invalidStateErrorMessage)
	}

	switch state {
	case "gratitude":
		if err := s.storage.storeGratitude(message.UserID, message.Message); err != nil {
			log.Printf("Error storing gratitude: %v", err)
			return s.sendMessage(message.UserID, gratitudeFailedErrorMessage)
		}

		delete(s.states, message.UserID)

		return s.sendMessage(message.UserID, "Successfully stored gratitude!")
	default:
		return s.sendMessage(message.UserID, invalidStateErrorMessage)
	}
}

func (s *Service) sendMessage(userID int64, text string) error {
	err := s.bot.SendMessage(userID, text)
	return errors.Wrap(err, "failed to store user data")
}
