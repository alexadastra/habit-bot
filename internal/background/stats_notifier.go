package background

import (
	"context"
	"time"

	"github.com/alexadastra/habit_bot/internal/service"
)

func NewStatsNotifier(
	isEnabled bool,
	service *service.Service,
	tickerDuration time.Duration,
) backgroundJob {
	handleFunc := func(ctx context.Context) {
		// TODO: add service notification sending here
	}

	return newBackgroundJob(
		isEnabled,
		tickerDuration,
		handleFunc,
	)
}
