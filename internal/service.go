package internal

import (
	"fmt"
	"log"
	"time"
)

type Service struct {
	bot     *Bot
	storage *Storage
	states  map[int]string // TODO: replace with cache?
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
		if err := s.storage.storeUserData(userID, time.Now()); err != nil {
			log.Printf("Error storing user data: %v", err)
			if err := s.bot.SendMessage(int64(userID), "Error storing user data. Please try again."); err != nil {
				log.Printf("Error sending the message: %v", err)
			}
			return
		}
		if err := s.bot.SendMessage(int64(userID), "Successfully checked in!"); err != nil {
			log.Printf("Error sending the message: %v", err)
		}
	case "gratitude":
		s.states[userID] = "gratitude"
		if err := s.bot.SendMessage(int64(userID), "What are you grateful for today?"); err != nil {
			log.Printf("Error sending the message: %v", err)
		}
	default:
		if err := s.bot.SendMessage(int64(userID), fmt.Sprintf("Invalid command: %s", command)); err != nil {
			log.Printf("Error sending the message: %v", err)
		}
	}
}

func (s *Service) handleMessage(userID int, text string) {
	state, ok := s.states[userID]
	if !ok {
		if err := s.bot.SendMessage(int64(userID), "Invalid state. Please try again."); err != nil {
			log.Printf("Error sending the message: %v", err)
		}
		return
	}

	switch state {
	case "gratitude":
		err := s.storage.storeGratitude(userID, text)
		if err != nil {
			log.Printf("Error storing gratitude: %v", err)
			if err := s.bot.SendMessage(int64(userID), "Error storing gratitude. Please try again."); err != nil {
				log.Printf("Error sending the message: %v", err)
			}
			return
		}
		if err := s.bot.SendMessage(int64(userID), "Successfully stored gratitude!"); err != nil {
			log.Printf("Error sending the message: %v", err)
		}
		delete(s.states, userID)
	default:
		if err := s.bot.SendMessage(int64(userID), "Invalid state. Please try again."); err != nil {
			log.Printf("Error sending the message: %v", err)
		}
	}
}
