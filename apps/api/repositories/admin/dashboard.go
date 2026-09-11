package admin

import (
	"context"
	"time"

	"github.com/emmanuella-codes/nox/models"
)

func (r *pgRepository) Dashboard(ctx context.Context) (*models.AdminDashboard, error) {
	dashboard := &models.AdminDashboard{}
	err := r.db.QueryRow(ctx, `
		SELECT
			(SELECT COUNT(*) FROM users),
			(SELECT COUNT(*) FROM users WHERE created_at >= CURRENT_DATE),
			(SELECT COUNT(*) FROM users WHERE created_at >= now() - INTERVAL '7 days'),
			(SELECT COUNT(*) FROM users WHERE status = 'suspended'),
			(SELECT COUNT(*) FROM posts WHERE moderation_status = 'active'),
			(SELECT COUNT(*) FROM comments WHERE moderation_status = 'active'),
			(SELECT COUNT(*) FROM personas WHERE moderation_status = 'active'),
			(SELECT COUNT(*) FROM events WHERE moderation_status = 'active'),
			(SELECT COUNT(*) FROM sets WHERE moderation_status = 'active'),
			(SELECT COUNT(*) FROM stories WHERE moderation_status = 'active' AND expires_at > now()),
			(SELECT COUNT(*) FROM reports WHERE status = 'open'),
			(SELECT COUNT(*) FROM reports WHERE status = 'reviewing'),
			(SELECT COUNT(*) FROM moderation_actions WHERE created_at >= now() - INTERVAL '24 hours'),
			(SELECT COUNT(*) FROM notification_outbox WHERE status IN ('failed', 'dead')),
			(SELECT COUNT(*) FROM media_assets WHERE processing_status = 'pending')
	`).Scan(&dashboard.Users.Total, &dashboard.Users.NewToday, &dashboard.Users.NewThisWeek,
		&dashboard.Users.Suspended, &dashboard.Content.Posts, &dashboard.Content.Comments,
		&dashboard.Content.Personas, &dashboard.Content.Events, &dashboard.Content.Sets,
		&dashboard.Content.ActiveStories, &dashboard.Moderation.OpenReports,
		&dashboard.Moderation.ReviewingReports, &dashboard.Moderation.RecentActions,
		&dashboard.Operations.FailedNotifications, &dashboard.Operations.PendingMedia)
	if err != nil {
		return nil, err
	}
	dashboard.GeneratedAt = time.Now().UTC()
	return dashboard, nil
}
