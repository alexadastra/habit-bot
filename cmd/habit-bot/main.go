package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexadastra/habit_bot/internal"
	"github.com/alexadastra/habit_bot/internal/background"
	"github.com/alexadastra/habit_bot/internal/external/bot"
	"github.com/alexadastra/habit_bot/internal/service"
	"github.com/alexadastra/habit_bot/internal/service/actions_service"
	"github.com/alexadastra/habit_bot/internal/service/stats_service"
	"github.com/alexadastra/habit_bot/internal/storage/inmemory"
	"github.com/alexadastra/habit_bot/internal/storage/mongodb"
	"github.com/alexadastra/habit_bot/internal/storage/redis"
)

func main() {
	log.Println("app started")
	// Create a context with a cancel function
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := internal.NewConfigFromEnv()

	bot, err := bot.NewBot(config.BotToken)
	if err != nil {
		log.Fatal(err)
	}
	go func() { _ = bot.Start(ctx) }()
	defer bot.Stop()
	log.Println("bot created")

	mongoStorage, err := mongodb.NewStorage(ctx, config.MongoDBDSN)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("storage created")

	statsService := stats_service.NewStatsService(mongoStorage)

	queue := redis.NewRedisQueue(config.RedisCreds)
	log.Println("queue created")

	actionsService, err := actions_service.New(ctx, queue, mongoStorage)
	if err != nil {
		log.Fatal(err)
	}

	stateStorage := inmemory.NewInMemoryStateStorage()

	// Set up the service
	service := service.NewService(bot, mongoStorage, stateStorage, statsService)

	workerPool := internal.NewWorkerPool(
		bot.GetCommandsChan(),
		bot.GetMessagesChan(),
		service,
		10, // TODO: move to config
	)

	// Start the worker pool
	go workerPool.Start(ctx)
	defer func() { _ = workerPool.Stop() }()

	// Start notifications sending
	notifier := background.NewStatsNotifier(
		true, // TODO: move to config
		actionsService,
		time.Second, // TODO: move to config
	)

	go notifier.Start(ctx)
	defer func() { _ = notifier.Stop() }()

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
