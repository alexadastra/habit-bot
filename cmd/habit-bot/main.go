package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/alexadastra/habit_bot/internal"
)

func main() {
	log.Println("app started")
	// Create a context with a cancel function
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := internal.NewConfigFromEnv()

	bot, err := internal.NewBot(config.BotToken)
	if err != nil {
		log.Fatal(err)
	}
	go func() { _ = bot.Start(ctx) }()
	defer bot.Stop()
	log.Println("bot created")

	storage, err := internal.NewStorage(ctx, config.MongoDBDDN)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("storage created")

	// Set up the service
	service := internal.NewService(bot, storage)

	workerPool := internal.NewWorkerPool(bot.GetCommandsChan(), bot.GetMessagesChan(), service, 10)

	// Start the worker pool
	go workerPool.Start(ctx)
	defer func() { _ = workerPool.Stop() }()

	// Run the bot until the context is cancelled
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sigChan:
		// System signal received.
		cancel()
	case <-ctx.Done():
		// Context cancelled.
	}
}
