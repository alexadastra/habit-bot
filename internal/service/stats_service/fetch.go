package stats_service

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

func (s *StatsService) FetchStats(ctx context.Context, userID int64) (int64, int64, error) {
	to := time.Now()
	from := to.Add(-24 * 7 * time.Hour)

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
