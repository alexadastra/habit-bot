package internal

import (
	"context"
	"log"

	"github.com/alexadastra/habit_bot/internal/models"
	"golang.org/x/sync/errgroup"
)

type WorkerPool struct {
	commands   <-chan models.UserCommand
	messages   <-chan models.UserMessage
	service    *Service
	group      *errgroup.Group
	cancel     context.CancelFunc
	numWorkers int
}

func NewWorkerPool(commands <-chan models.UserCommand, messages <-chan models.UserMessage, service *Service, numWorkers int) *WorkerPool {
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
			// TODO: log warning here
			log.Println("worker stopped!")
			return ctx.Err()
		case command, ok := <-wp.commands:
			if !ok {
				// TODO: log warning here
				return nil
			}
			// TODO: process error here
			wp.service.handleCommand(command.Command, command.Id, command.Message)
		case message, ok := <-wp.messages:
			if !ok {
				// TODO: log warning here
				return nil
			}
			// TODO: process error here
			wp.service.handleMessage(message.Id, message.Message)
		}
	}
}

func (wp *WorkerPool) Stop() error {
	wp.cancel()
	return wp.group.Wait()
}
