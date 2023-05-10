package background

import (
	"context"
	"time"

	"github.com/alexadastra/habit_bot/internal/service/actions_service"
)

func NewStatsNotifier(
	isEnabled bool,
	service *actions_service.ActionsService,
	tickerDuration time.Duration,
) backgroundJob {
	handleFunc := func(ctx context.Context) {
		service.Process(ctx)
	}

	return newBackgroundJob(
		isEnabled,
		tickerDuration,
		handleFunc,
	)
}
