package notification

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	notification_dtos "github.com/emmanuella-codes/nox/notification/dtos"
	"github.com/google/uuid"
)

func (r *pgRepository) UpsertNotificationDevice(ctx context.Context, userID uuid.UUID, dto notification_dtos.UpsertNotificationDeviceDTO) (*models.NotificationDevice, error) {
	row := r.db.QueryRow(ctx, `
		INSERT INTO notification_devices (user_id, install_id, platform, push_token, app_version)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, install_id) DO UPDATE
		SET platform = EXCLUDED.platform,
		    push_token = EXCLUDED.push_token,
		    app_version = EXCLUDED.app_version,
		    last_seen_at = now(),
		    disabled_at = NULL,
		    updated_at = now()
		RETURNING id, user_id, install_id, platform, push_token, app_version, last_seen_at,
		          disabled_at, created_at, updated_at
	`, userID, dto.InstallID, dto.Platform, dto.PushToken, dto.AppVersion)
	return scanNotificationDevice(row)
}

func (r *pgRepository) FindNotificationDevices(ctx context.Context, userID uuid.UUID) ([]*models.NotificationDevice, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, install_id, platform, push_token, app_version, last_seen_at,
		       disabled_at, created_at, updated_at
		FROM notification_devices
		WHERE user_id = $1
		ORDER BY last_seen_at DESC, created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanNotificationDevices(rows)
}

func (r *pgRepository) DisableNotificationDevice(ctx context.Context, userID uuid.UUID, deviceID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
		UPDATE notification_devices
		SET disabled_at = COALESCE(disabled_at, now()),
		    updated_at = now()
		WHERE id = $1 AND user_id = $2
	`, deviceID, userID)
	return err
}

func (r *pgRepository) DisableNotificationDeviceByToken(ctx context.Context, pushToken string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE notification_devices
		SET disabled_at = COALESCE(disabled_at, now()),
		    updated_at = now()
		WHERE push_token = $1
	`, pushToken)
	return err
}
