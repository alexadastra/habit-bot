package main

import (
	"context"
	"log"
	"os"

	"github.com/alexadastra/habit_bot/internal"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	log.Println("app started")
	// Create a context with a cancel function
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up the telegram bot API
	tgBot, err := tgbotapi.NewBotAPI(os.Getenv("BOT-TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	bot := internal.NewBot(tgBot)
	log.Println("bot created")

	storage, err := internal.NewStorage(os.Getenv("MONGO-DB-DSN"))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("storage created")

	// Set up the service
	service := internal.NewService(bot, storage)

	// Set up the worker pool
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates, err := tgBot.GetUpdatesChan(updateConfig)
	if err != nil {
		log.Fatal(err)
	}
	workerPool := internal.NewWorkerPool(updates, bot, service, 10)

	// Start the worker pool
	go workerPool.Start(ctx)
	defer workerPool.Stop()

	// Run the bot until the context is cancelled
	for {
		select {
		case <-ctx.Done():
			// Context cancelled, shutdown gracefully
		}
	}
}
