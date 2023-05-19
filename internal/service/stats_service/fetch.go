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

func findWeekRange(t time.Time) (time.Time, time.Time) {
	tUTC := t.UTC()

	// timeZoneDiff = t - tUTC in nanoseconds (positive for zones to the east from Greenwich, negative - for zones to the west)
	_, locOffset := t.Zone()
	_, utcOffset := tUTC.Zone()
	timeZoneDiff := time.Duration((locOffset - utcOffset) * 1000000000)

	// extract the day of week (0 for Sunday, 1 for Monday, ..., 6 for Saturday)
	currentDay := int(tUTC.Weekday())

	// handle Sunday as i want the week to be started from Monday
	if currentDay == 0 {
		currentDay = 7
	}

	// whenever currentDay is, if we count now() - (currentDay - 1) days, it'll always point to Monday
	fromDayUTC := tUTC.Add(-1 * (time.Duration(currentDay) - 1) * 24 * time.Hour)

	// set time to UTC midnight
	fromDayUTC = fromDayUTC.Truncate(24 * time.Hour)

	// fromDayUTCDiff = tUTC - fromDay (at least zero)
	fromDayUTCDiff := tUTC.Sub(fromDayUTC)

	// extract found difference between current UTC and monday midnight UTC
	fromDayLoc := t.Add(-1 * fromDayUTCDiff)

	// extract difference between time zones, as we truncated the UTC time
	fromDayLoc = fromDayLoc.Add(-1 * timeZoneDiff)

	return fromDayLoc, fromDayLoc.Add(7 * 24 * time.Hour)
}
