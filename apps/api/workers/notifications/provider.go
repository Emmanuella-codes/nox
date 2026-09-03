package main

import (
	"context"

	"github.com/emmanuella-codes/nox/config"
	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/shared/push"
	"github.com/rs/zerolog/log"
)

type logProvider struct{}

func newProvider(cfg *config.Config) push.Provider {
	switch cfg.PushProvider {
	default:
		return logProvider{}
	}
}

func (logProvider) Send(ctx context.Context, device *models.NotificationDevice, payload push.Payload) error {
	log.Info().
		Str("device_id", device.ID.String()).
		Str("platform", string(device.Platform)).
		Str("install_id", device.InstallID).
		Str("target_path", payload.TargetPath).
		Int("badge", payload.Badge).
		Msg("notification push delivered")
	return nil
}
