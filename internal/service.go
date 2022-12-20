package internal

import (
	"fmt"
	"log"
	"time"
)

type Service struct {
	bot     *Bot
	storage *Storage
	states  map[int]string
}

func NewService(bot *Bot, storage *Storage) *Service {
	return &Service{
		bot:     bot,
		storage: storage,
		states:  make(map[int]string),
	}
}

func (s *Service) handleCommand(command string, userID int, text string) {
	switch command {
	case "checkin":
		err := s.storage.storeUserData(userID, time.Now())
		if err != nil {
			log.Printf("Error storing user data: %v", err)
			s.bot.SendMessage(int64(userID), "Error storing user data. Please try again.")
			return
		}
		s.bot.SendMessage(int64(userID), "Successfully checked in!")
	case "gratitude":
		s.states[userID] = "gratitude"
		s.bot.SendMessage(int64(userID), "What are you grateful for today?")
	default:
		s.bot.SendMessage(int64(userID), fmt.Sprintf("Invalid command: %s", command))
	}
}

func (s *Service) handleMessage(userID int, text string) {
	state, ok := s.states[userID]
	if !ok {
		s.bot.SendMessage(int64(userID), "Invalid state. Please try again.")
		return
	}

	switch state {
	case "gratitude":
		err := s.storage.storeGratitude(userID, text)
		if err != nil {
			log.Printf("Error storing gratitude: %v", err)
			s.bot.SendMessage(int64(userID), "Error storing gratitude. Please try again.")
			return
		}
		s.bot.SendMessage(int64(userID), "Successfully stored gratitude!")
		delete(s.states, userID)
	default:
		s.bot.SendMessage(int64(userID), "Invalid state. Please try again.")
	}
}
