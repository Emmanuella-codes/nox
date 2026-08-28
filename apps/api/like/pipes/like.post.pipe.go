package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/like/dtos"
	"github.com/emmanuella-codes/nox/like/messages"
	"github.com/emmanuella-codes/nox/models"
	notification_repo "github.com/emmanuella-codes/nox/repositories/notification"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

func (p *LikePipe) LikePostPipe(ctx context.Context, userID uuid.UUID, postID uuid.UUID, dto dtos.LikePostDTO) *shared.PipeRes[any] {
	persona, post, res := p.validatePersonaAndPost(ctx, userID, postID, dto.PersonaID)
	if res != nil {
		return res
	}
	if err := p.likeRepo.LikePost(ctx, dto.PersonaID, postID); err != nil {
		return pipeInternalError[any](err, "like.post")
	}
	p.createLikeNotification(ctx, persona, post)
	return shared.PipeSuccess[any](messages.Post_Liked, nil)
}

// createLikeNotification persists and publishes one post-like notification.
func (p *LikePipe) createLikeNotification(ctx context.Context, actor *models.Persona, post *models.Post) {
	if p.notificationRepo == nil || post == nil || post.PersonaID == nil || post.AuthorUserID == actor.UserID {
		return
	}
	created, err := p.notificationRepo.CreateNotifications(ctx, []notification_repo.CreateNotificationInput{{
		RecipientUserID:    post.AuthorUserID,
		RecipientPersonaID: *post.PersonaID,
		ActorPersonaID:     &actor.ID,
		ActorPostingMode:   models.PublicPostingMode,
		PostID:             &post.ID,
		NotificationType:   models.LikeNotificationType,
	}})
	if err != nil || len(created) == 0 || p.notificationPublisher == nil {
		return
	}
	p.notificationPublisher.PublishCreatedNotification(ctx, created[0])
}
