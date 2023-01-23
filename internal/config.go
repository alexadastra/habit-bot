package internal

import (
	"os"
)

type Config struct {
	BotToken   string
	MongoDBDDN string
}

func NewConfigFromEnv() *Config {
	c := &Config{
		BotToken:   os.Getenv("BOT_TOKEN"),
		MongoDBDDN: os.Getenv("MONGO_DB_DSN"),
	}

	return c
}
