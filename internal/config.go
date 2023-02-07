package internal

import (
	"fmt"
	"os"
)

type Config struct {
	BotToken   string
	MongoDBDDN string
}

func NewConfigFromEnv() *Config {
	c := &Config{
		BotToken:   os.Getenv("BOT_TOKEN"),
		MongoDBDDN: fmt.Sprintf(
			"mongodb://%s:%s@%s/%s",
			os.Getenv("MONGODB_USER"),
			os.Getenv("MONGODB_PASSWORD"),
			os.Getenv("MONGODB_HOST"),
			os.Getenv("MONGODB_DATABASE"),
		),
	}

	return c
}
