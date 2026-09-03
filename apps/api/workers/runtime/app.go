package runtime

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/emmanuella-codes/nox/config"
	"github.com/emmanuella-codes/nox/db"
	"github.com/emmanuella-codes/nox/repositories"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Options struct {
	ConnectRedis  bool
	RunMigrations bool
}

type App struct {
	Context context.Context
	Config  *config.Config
	DB      *pgxpool.Pool
	Redis   *redis.Client
	Repos   *repositories.Repositories
	stop    context.CancelFunc
}

func SignalContext() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
}

func Bootstrap(ctx context.Context, options Options) (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	if options.RunMigrations {
		if err := db.RunMigrations(cfg.DatabaseURL); err != nil {
			return nil, err
		}
	}
	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}
	var redisClient *redis.Client
	if options.ConnectRedis {
		redisClient, err = db.ConnectRedis(ctx, cfg.RedisURL)
		if err != nil {
			pool.Close()
			return nil, err
		}
	}
	childCtx, stop := context.WithCancel(ctx)
	return &App{
		Context: childCtx,
		Config:  cfg,
		DB:      pool,
		Redis:   redisClient,
		Repos:   repositories.Init(pool),
		stop:    stop,
	}, nil
}

func (a *App) Close() {
	if a == nil {
		return
	}
	if a.stop != nil {
		a.stop()
	}
	if a.Redis != nil {
		_ = a.Redis.Close()
	}
	if a.DB != nil {
		a.DB.Close()
	}
}
