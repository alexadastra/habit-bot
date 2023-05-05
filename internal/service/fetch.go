package service

import (
	"context"
	"fmt"
	"time"

	"github.com/alexadastra/habit_bot/internal/models"
)

func (s *Service) fetchStats(ctx context.Context, command models.UserCommand) error {
	to := time.Now()
	from := to.Add(-24 * 7 * time.Hour)

	checkins, err := s.actionsStorage.GetCheckinEvents(
		ctx,
		command.UserID,
		from,
		to,
	)
	if err != nil {
		return s.sendMessage(command.UserID, checkingFetchingFailedErrorMessage)
	}

	gratitudes, err := s.actionsStorage.GetGratitudeEvents(
		ctx,
		command.UserID,
		from,
		to,
	)
	if err != nil {
		return s.sendMessage(command.UserID, checkingFetchingFailedErrorMessage)
	}

	return s.sendMessage(
		command.UserID,
		fmt.Sprintf(
			"Great job! This week you've tracked down %d checkin events and %d gratitude events!",
			len(checkins),
			len(gratitudes),
		),
	)
}
