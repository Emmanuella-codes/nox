package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// GetJSON loads one JSON payload from Redis into the supplied output value.
func GetJSON(ctx context.Context, client *redis.Client, key string, out any) (bool, error) {
	if client == nil {
		return false, nil
	}
	value, err := client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, err
	}
	if err := json.Unmarshal([]byte(value), out); err != nil {
		return false, err
	}
	return true, nil
}

// SetJSON stores one JSON payload in Redis with the supplied TTL.
func SetJSON(ctx context.Context, client *redis.Client, key string, value any, ttl time.Duration) error {
	if client == nil {
		return nil
	}
	payload, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return client.Set(ctx, key, payload, ttl).Err()
}
