package service

import (
	"context"
	"fmt"
	"log"

	"github.com/alexadastra/habit_bot/internal/models"
)

func (s *Service) handleStats(ctx context.Context, command models.UserCommand) error {
	checkins, gratitudes, err := s.statsFetcher.FetchStats(ctx, command.UserID)

	if err != nil {
		log.Println(err)
		return s.sendMessage(
			command.UserID,
			statsFetchingFailedErrorMessage,
		)
	}

	return s.sendMessage(
		command.UserID,
		fmt.Sprintf(
			"Great job! This week you've tracked down %d checkin events and %d gratitude events!",
			checkins,
			gratitudes,
		),
	)
}
