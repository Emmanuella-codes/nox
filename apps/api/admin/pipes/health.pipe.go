package pipes

import (
	"context"
	"time"

	"github.com/emmanuella-codes/nox/admin/messages"
	"github.com/emmanuella-codes/nox/shared"
	worker_runtime "github.com/emmanuella-codes/nox/workers/runtime"
)

type HealthResponse struct {
	Status      string                  `json:"status"`
	Environment string                  `json:"environment"`
	Time        time.Time               `json:"time"`
	Database    ComponentHealth         `json:"database"`
	Redis       ComponentHealth         `json:"redis"`
	Workers     map[string]WorkerHealth `json:"workers"`
}

type ComponentHealth struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

type WorkerHealth struct {
	Status   string     `json:"status"`
	LastSeen *time.Time `json:"last_seen,omitempty"`
}

func (p *AdminPipe) HealthPipe(ctx context.Context) *shared.PipeRes[HealthResponse] {
	response := HealthResponse{
		Status:      "ok",
		Environment: p.cfg.Environment,
		Time:        time.Now().UTC(),
		Database:    ComponentHealth{Status: "ok"},
		Redis:       ComponentHealth{Status: "ok"},
		Workers: map[string]WorkerHealth{
			worker_runtime.WorkerMediaKey:         {Status: "offline"},
			worker_runtime.WorkerNotificationsKey: {Status: "offline"},
			worker_runtime.WorkerStoriesKey:       {Status: "offline"},
		},
	}

	dbErr := error(nil)
	if p.db == nil {
		dbErr = context.Canceled
	} else {
		dbErr = p.db.Ping(ctx)
	}
	if dbErr != nil {
		response.Status = "degraded"
		response.Database.Status = "down"
		if p.db == nil {
			response.Database.Error = "database_unavailable"
		} else {
			response.Database.Error = dbErr.Error()
		}
	}

	if p.redis == nil {
		response.Status = "degraded"
		response.Redis.Status = "down"
		response.Redis.Error = "redis_unavailable"
		return shared.PipeSuccess(messages.Admin_Health_Loaded, &response)
	}

	if _, err := p.redis.Ping(ctx).Result(); err != nil {
		response.Status = "degraded"
		response.Redis.Status = "down"
		response.Redis.Error = err.Error()
		return shared.PipeSuccess(messages.Admin_Health_Loaded, &response)
	}

	p.loadWorkerHealth(ctx, &response, worker_runtime.WorkerMediaKey, p.cfg.MediaCleanupInterval)
	p.loadWorkerHealth(ctx, &response, worker_runtime.WorkerNotificationsKey, p.cfg.PushWorkerPollInterval)
	p.loadWorkerHealth(ctx, &response, worker_runtime.WorkerStoriesKey, p.cfg.StoryCleanupInterval)
	p.loadWorkerHealth(ctx, &response, worker_runtime.WorkerCrewsKey, p.cfg.CrewCleanupInterval)
	p.loadWorkerHealth(ctx, &response, worker_runtime.WorkerModerationKey, p.cfg.ModerationCleanupInterval)

	if response.Database.Status == "down" && response.Redis.Status == "down" {
		response.Status = "down"
	}

	return shared.PipeSuccess(messages.Admin_Health_Loaded, &response)
}

func (p *AdminPipe) loadWorkerHealth(ctx context.Context, response *HealthResponse, workerName string, interval time.Duration) {
	rawValue, err := p.redis.Get(ctx, worker_runtime.WorkerHeartbeatKey(workerName)).Result()
	if err != nil {
		response.Status = "degraded"
		response.Workers[workerName] = WorkerHealth{Status: "offline"}
		return
	}

	lastSeen, err := time.Parse(time.RFC3339Nano, rawValue)
	if err != nil {
		response.Status = "degraded"
		response.Workers[workerName] = WorkerHealth{Status: "offline"}
		return
	}

	workerStatus := "online"
	if time.Since(lastSeen) > worker_runtime.WorkerHeartbeatTTL(interval) {
		workerStatus = "stale"
		response.Status = "degraded"
	}
	response.Workers[workerName] = WorkerHealth{
		Status:   workerStatus,
		LastSeen: &lastSeen,
	}
}
