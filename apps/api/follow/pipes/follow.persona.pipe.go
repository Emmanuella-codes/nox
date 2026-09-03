package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/follow/messages"
	"github.com/emmanuella-codes/nox/models"
	notification_repo "github.com/emmanuella-codes/nox/repositories/notification"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

func (p *FollowPipe) FollowPersonaPipe(ctx context.Context, userID uuid.UUID, followerPersonaID uuid.UUID, targetPersonaID uuid.UUID) *shared.PipeRes[any] {
	followerPersona, targetPersona, res := p.validateFollowAction(ctx, userID, followerPersonaID, targetPersonaID)
	if res != nil {
		return res
	}

	if err := p.followRepo.Follow(ctx, followerPersonaID, targetPersonaID); err != nil {
		return p.mapFollowError(err, "follow.create")
	}
	p.createFollowNotification(ctx, followerPersona, targetPersona)

	return shared.PipeSuccess[any](messages.Followed_Successfully, nil)
}

// createFollowNotification persists and publishes one follow notification.
func (p *FollowPipe) createFollowNotification(ctx context.Context, actor *models.Persona, target *models.Persona) {
	if p.notificationRepo == nil || target == nil || actor == nil || target.UserID == actor.UserID {
		return
	}
	created, err := p.notificationRepo.CreateNotifications(ctx, []notification_repo.CreateNotificationInput{{
		RecipientUserID:    target.UserID,
		RecipientPersonaID: target.ID,
		ActorPersonaID:     &actor.ID,
		ActorPostingMode:   models.PublicPostingMode,
		NotificationType:   models.FollowNotificationType,
	}})
	if err != nil || len(created) == 0 || p.notificationPublisher == nil {
		return
	}
	p.notificationPublisher.PublishCreatedNotification(ctx, created[0])
}
