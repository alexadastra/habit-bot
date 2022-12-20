package internal

import (
	"context"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"golang.org/x/sync/errgroup"
)

type WorkerPool struct {
	updates    <-chan tgbotapi.Update
	bot        *Bot
	service    *Service
	group      *errgroup.Group
	cancel     context.CancelFunc
	numWorkers int
}

func NewWorkerPool(updates <-chan tgbotapi.Update, bot *Bot, service *Service, numWorkers int) *WorkerPool {
	return &WorkerPool{
		updates:    updates,
		bot:        bot,
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
			log.Println("worker stopped!")
			return ctx.Err()
		case update, ok := <-wp.updates:
			if !ok {
				return nil
			}
			log.Printf("Worker: Handling update %d", update.UpdateID)
			if err := wp.handleUpdate(ctx, update); err != nil {
				return err
			}
		}
	}
}

func (wp *WorkerPool) handleUpdate(ctx context.Context, update tgbotapi.Update) error {
	// Check if the context has been cancelled
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Handle commands
	if update.Message != nil && update.Message.IsCommand() {
		wp.service.handleCommand(update.Message.Command(), update.Message.From.ID, update.Message.Text)
	}

	// Handle messages that are not commands
	if update.Message != nil && !update.Message.IsCommand() {
		wp.service.handleMessage(update.Message.From.ID, update.Message.Text)
	}

	return nil
}

func (wp *WorkerPool) Stop() {
	wp.cancel()
	wp.group.Wait()
}
