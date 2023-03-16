package internal

import (
	"context"
	"log"

	"github.com/alexadastra/habit_bot/internal/models"
	"github.com/alexadastra/habit_bot/internal/service"
	"golang.org/x/sync/errgroup"
)

type WorkerPool struct {
	commands   <-chan models.UserCommand
	messages   <-chan models.UserMessage
	service    *service.Service
	group      *errgroup.Group
	cancel     context.CancelFunc
	numWorkers int
}

func NewWorkerPool(
	commands <-chan models.UserCommand,
	messages <-chan models.UserMessage,
	service *service.Service,
	numWorkers int,
) *WorkerPool {
	return &WorkerPool{
		commands:   commands,
		messages:   messages,
		service:    service,
		numWorkers: numWorkers,
	}
}

func (wp *WorkerPool) Start(ctx context.Context) {
	ctx, wp.cancel = context.WithCancel(ctx)
	wp.group, ctx = errgroup.WithContext(ctx)
	for i := 0; i < wp.numWorkers; i++ {
		id := i
		wp.group.Go(func() error {
			return wp.handleUpdates(ctx, id)
		})
	}
}

func (wp *WorkerPool) handleUpdates(ctx context.Context, id int) error {
	log.Printf("worker %d started", id)
	for {
		select {
		case <-ctx.Done():
			log.Printf("worker %d stopped: context cancelled", id)
			return ctx.Err()
		case command, ok := <-wp.commands:
			if !ok {
				log.Printf("worker %d stopped: commands channel closed", id)
				return nil
			}
			if err := wp.service.HandleCommand(ctx, command); err != nil {
				log.Printf("error while handling command: %s", err)
			}
		case message, ok := <-wp.messages:
			if !ok {
				log.Printf("worker %d stopped: messages channel closed", id)
				return nil
			}
			if err := wp.service.HandleMessage(ctx, message); err != nil {
				log.Printf("error while handling message: %s", err)
			}
		}
	}
}

func (wp *WorkerPool) Stop() error {
	wp.cancel()
	return wp.group.Wait()
}
