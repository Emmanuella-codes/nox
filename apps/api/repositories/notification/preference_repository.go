package notification

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
)

func (r *pgRepository) FindNotificationPreferences(ctx context.Context, personaID uuid.UUID) ([]*models.NotificationPreference, error) {
	rows, err := r.db.Query(ctx, `
		SELECT persona_id, notification_type, push_enabled, updated_at
		FROM notification_preferences
		WHERE persona_id = $1
		ORDER BY notification_type ASC
	`, personaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanNotificationPreferences(rows)
}

func (r *pgRepository) UpsertNotificationPreference(ctx context.Context, personaID uuid.UUID, notificationType models.NotificationType, pushEnabled bool) (*models.NotificationPreference, error) {
	row := r.db.QueryRow(ctx, `
		INSERT INTO notification_preferences (persona_id, notification_type, push_enabled)
		VALUES ($1, $2, $3)
		ON CONFLICT (persona_id, notification_type) DO UPDATE
		SET push_enabled = EXCLUDED.push_enabled,
		    updated_at = now()
		RETURNING persona_id, notification_type, push_enabled, updated_at
	`, personaID, notificationType, pushEnabled)
	return scanNotificationPreference(row)
}

func (r *pgRepository) PushEnabledForPersona(ctx context.Context, personaID uuid.UUID, notificationType models.NotificationType) (bool, error) {
	var enabled bool
	err := r.db.QueryRow(ctx, `
		SELECT COALESCE((
			SELECT push_enabled
			FROM notification_preferences
			WHERE persona_id = $1 AND notification_type = $2
		), TRUE)
	`, personaID, notificationType).Scan(&enabled)
	return enabled, err
}
