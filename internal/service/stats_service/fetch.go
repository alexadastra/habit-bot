package stats_service

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

func (s *StatsService) FetchStats(ctx context.Context, userID int64) (int64, int64, error) {
	from, to := findWeekRange(time.Now().UTC())

	checkins, err := s.storage.GetCheckinEvents(
		ctx,
		userID,
		from,
		to,
	)
	if err != nil {
		return 0, 0, errors.Wrap(err, "failed to fetch checkins")
	}

	gratitudes, err := s.storage.GetGratitudeEvents(
		ctx,
		userID,
		from,
		to,
	)
	if err != nil {
		return 0, 0, errors.Wrap(err, "failed to fetch gratitudes")
	}

	return int64(len(checkins)), int64(len(gratitudes)), nil
}

// suppose that t has UTC location
func findWeekRange(t time.Time) (time.Time, time.Time) {
	// extract the day of week (0 for Sunday, 1 for Monday, ..., 6 for Saturday)
	currentDay := int(t.Weekday())

	// handle Sunday as i want the week to be started from Monday
	if currentDay == 0 {
		currentDay = 7
	}

	// whenever currentDay is, if we count now() - (currentDay - 1) days, it'll always point to Monday
	fromDay := t.Add((time.Duration(currentDay) - 1) * 24 * time.Hour)

	fromDay = fromDay.Truncate(24 * time.Hour)

	return fromDay, fromDay.Add(7 * 24 * time.Hour)
}
