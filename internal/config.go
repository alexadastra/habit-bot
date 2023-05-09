package internal

import (
	"fmt"
	"os"
)

type Config struct {
	TelegramBotCreds
	MongoDBCreds
	RedisCreds
}

type TelegramBotCreds struct {
	BotToken string
}

type MongoDBCreds struct {
	MongoDBDSN string
}

type RedisCreds struct {
	RedisAddr     string
	RedisPassword string
}

func NewConfigFromEnv() *Config {
	c := &Config{
		TelegramBotCreds: TelegramBotCreds{
			BotToken: os.Getenv("BOT_TOKEN"),
		},

		MongoDBCreds: MongoDBCreds{
			MongoDBDSN: fmt.Sprintf(
				"mongodb://%s:%s@%s/%s",
				os.Getenv("MONGODB_USER"),
				os.Getenv("MONGODB_PASSWORD"),
				os.Getenv("MONGODB_HOST"),
				os.Getenv("MONGODB_DATABASE"),
			),
		},

		RedisCreds: RedisCreds{
			RedisAddr:     os.Getenv("REDIS_ADDR"),
			RedisPassword: os.Getenv("REDIS_PASSWD"),
		},
	}

	return c
}
