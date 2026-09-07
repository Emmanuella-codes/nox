package runtime

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	WorkerMediaKey         = "media"
	WorkerNotificationsKey = "notifications"
	WorkerStoriesKey       = "stories"
	WorkerCrewsKey         = "crews"
)

func WorkerHeartbeatKey(name string) string {
	return "worker:" + name + ":last_seen"
}

func WorkerHeartbeatTTL(interval time.Duration) time.Duration {
	if interval <= 0 {
		return time.Minute
	}
	if interval < 10*time.Second {
		return 30 * time.Second
	}
	return interval * 3
}

func RecordHeartbeat(ctx context.Context, redisClient *redis.Client, workerName string, ttl time.Duration) error {
	if redisClient == nil {
		return nil
	}
	return redisClient.Set(ctx, WorkerHeartbeatKey(workerName), time.Now().UTC().Format(time.RFC3339Nano), ttl).Err()
}
