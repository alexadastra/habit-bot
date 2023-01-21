package main

import (
	"context"
	"log"
	"os"

	"github.com/alexadastra/habit_bot/internal"
)

func main() {
	log.Println("app started")
	// Create a context with a cancel function
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bot, err := internal.NewBot(os.Getenv("BOT-TOKEN"))
	if err != nil {
		log.Fatal(err)
	}
	go func() { _ = bot.Start(ctx) }()
	defer bot.Stop()
	log.Println("bot created")

	storage, err := internal.NewStorage(os.Getenv("MONGO-DB-DSN"))
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
	for {
		select {
		case <-ctx.Done():
			// Context cancelled, shutdown gracefully
		}
	}
}
