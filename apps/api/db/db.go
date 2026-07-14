package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func Connect(ctx context.Context, databaseUrl string) (*pgxpool.Pool, error) {
	databaseUrl = normalizeDatabaseURL(databaseUrl)

	cfg, err := pgxpool.ParseConfig(databaseUrl)
	if err != nil {
		return nil, err
	}

	cfg.MaxConns = 20
	cfg.MaxConnLifetime = 30 * time.Minute
	cfg.MaxConnIdleTime = 5 * time.Minute
	cfg.HealthCheckPeriod = 30 * time.Second

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("db: ping: %w", err)
	}

	log.Info().Msg("database connected")
	return pool, nil
}
