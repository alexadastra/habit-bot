package redis

import (
	"github.com/alexadastra/habit_bot/internal"
	redis "github.com/go-redis/redis/v8"
)

func newRedisClient(creds internal.RedisCreds) *redis.Client {
	return redis.NewClient(
		&redis.Options{
			Addr:     creds.RedisAddr,
			Password: creds.RedisPassword,
			DB:       0, // use default DB
		},
	)
}
