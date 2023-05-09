package redis

import (
	"context"

	"github.com/alexadastra/habit_bot/internal"
	redis "github.com/go-redis/redis/v8"
)

const actionsQueueKey = "actions_queue"

// RedisQueue represents a priority queue backed by Redis
type RedisQueue struct {
	client *redis.Client
}

// NewRedisQueue creates a new RedisQueue with the given Redis client
func NewRedisQueue(creds internal.RedisCreds) *RedisQueue {
	client := newRedisClient(creds)

	return &RedisQueue{
		client: client,
	}
}

// Push adds an action to the queue with the given priority
func (q *RedisQueue) Push(ctx context.Context, id string, priority int64) error {
	_, err := q.client.
		ZAdd(
			ctx,
			actionsQueueKey,
			&redis.Z{
				Score:  float64(priority),
				Member: id,
			},
		).
		Result()
	return err
}

// Pop removes and returns the action with the highest priority from the queue
func (q *RedisQueue) Pop(ctx context.Context, count int64) ([]string, error) {
	res, err := q.client.
		ZRangeByScore(
			ctx,
			actionsQueueKey,
			&redis.ZRangeBy{
				Min:    "-inf",
				Max:    "+inf",
				Offset: 0,
				Count:  count,
			},
		).
		Result()
	if err != nil {
		return nil, err
	}

	// remove by id so we won't possibly affect entries that were added after the select
	for _, id := range res {
		_, err = q.client.
			ZRem(
				ctx,
				actionsQueueKey,
				id,
			).
			Result()
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}
