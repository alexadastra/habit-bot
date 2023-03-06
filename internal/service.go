package internal

import (
	"fmt"
	"log"
	"time"

	"github.com/alexadastra/habit_bot/internal/models"
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

func (s *Service) handleCommand(command models.UserCommand) {
	switch command.Command {
	case models.Checkin:
		if err := s.storage.storeUserData(command.UserID, time.Now()); err != nil {
			log.Printf("Error storing user data: %v", err)
			if err := s.bot.SendMessage(command.UserID, "Error storing user data. Please try again."); err != nil {
				log.Printf("Error sending the message: %v", err)
			}
			return
		}
		if err := s.bot.SendMessage(command.UserID, "Successfully checked in!"); err != nil {
			log.Printf("Error sending the message: %v", err)
		}
	case models.Gratitude:
		s.states[command.UserID] = "gratitude"
		if err := s.bot.SendMessage(command.UserID, "What are you grateful for today?"); err != nil {
			log.Printf("Error sending the message: %v", err)
		}
	default:
		if err := s.bot.SendMessage(command.UserID, fmt.Sprintf("Invalid command: %s", command.Command)); err != nil {
			log.Printf("Error sending the message: %v", err)
		}
	}
}

func (s *Service) handleMessage(message models.UserMessage) {
	state, ok := s.states[message.UserID]
	if !ok {
		if err := s.bot.SendMessage(message.UserID, "Invalid state. Please try again."); err != nil {
			log.Printf("Error sending the message: %v", err)
		}
		return
	}

	switch state {
	case "gratitude":
		err := s.storage.storeGratitude(message.UserID, message.Message)
		if err != nil {
			log.Printf("Error storing gratitude: %v", err)
			if err := s.bot.SendMessage(message.UserID, "Error storing gratitude. Please try again."); err != nil {
				log.Printf("Error sending the message: %v", err)
			}
			return
		}
		if err := s.bot.SendMessage(message.UserID, "Successfully stored gratitude!"); err != nil {
			log.Printf("Error sending the message: %v", err)
		}
		delete(s.states, message.UserID)
	default:
		if err := s.bot.SendMessage(message.UserID, "Invalid state. Please try again."); err != nil {
			log.Printf("Error sending the message: %v", err)
		}
	}
}
