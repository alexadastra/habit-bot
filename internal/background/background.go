package background

import (
	"context"
	"sync"
	"time"
)

type backgroundJob struct {
	duration time.Duration
	ticker   *time.Ticker

	callback func(ctx context.Context)

	isEnabled bool
	once      *sync.Once
	wg        *sync.WaitGroup
	cancel    context.CancelFunc
}

func newBackgroundJob(
	isEnabled bool,
	duration time.Duration,
	callback func(ctx context.Context),
) backgroundJob {
	return backgroundJob{
		isEnabled: isEnabled,
		duration:  duration,
		callback:  callback,
		wg:        &sync.WaitGroup{},
		once:      &sync.Once{},
	}
}

func (bj *backgroundJob) Start(ctx context.Context) {
	bj.once.Do(func() {
		// nolint:govet
		ctx, cancel := context.WithCancel(ctx)
		bj.cancel = cancel
		bj.wg.Add(1)

		bj.ticker = time.NewTicker(bj.duration)

		go func() {
			defer bj.wg.Done()
			bj.loop(ctx)
		}()
	})
}

func (bj *backgroundJob) Stop() error {
	bj.ticker.Stop()
	bj.cancel()
	bj.wg.Wait()
	bj.once = &sync.Once{}

	return nil
}

func (bj *backgroundJob) loop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			// TODO: log warning here
			return
		case <-bj.ticker.C:
			bj.callback(ctx)
		}
	}
}
