package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	notification_repo "github.com/emmanuella-codes/nox/repositories/notification"
)

// createNotifications persists story notifications and broadcasts any created rows.
func (p *StoryPipe) createNotifications(ctx context.Context, inputs []notification_repo.CreateNotificationInput) error {
	if p.notificationRepo == nil || len(inputs) == 0 {
		return nil
	}
	notifications, err := p.notificationRepo.CreateNotifications(ctx, inputs)
	if err != nil {
		return err
	}
	if p.notificationPublisher == nil {
		return nil
	}
	for _, notification := range notifications {
		p.notificationPublisher.PublishCreatedNotification(ctx, notification)
	}
	return nil
}

// notifyStoryContributionRequest publishes one pending contribution request notification to the story owner.
func (p *StoryPipe) notifyStoryContributionRequest(ctx context.Context, story *models.Story, contributor *models.Persona, request *models.StoryContributionRequest) error {
	if story.OwnerUserID == contributor.UserID {
		return nil
	}
	return p.createNotifications(ctx, []notification_repo.CreateNotificationInput{{
		RecipientUserID:            story.OwnerUserID,
		RecipientPersonaID:         story.OwnerPersonaID,
		ActorPersonaID:             &contributor.ID,
		ActorPostingMode:           models.PublicPostingMode,
		StoryID:                    &story.ID,
		StoryContributionRequestID: &request.ID,
		NotificationType:           models.StoryContributionRequestNotificationType,
	}})
}

// notifyStoryContributionDecision publishes one accepted or rejected contribution decision back to the contributor.
func (p *StoryPipe) notifyStoryContributionDecision(ctx context.Context, story *models.Story, contributor *models.Persona, notificationType models.NotificationType, request *models.StoryContributionRequest) error {
	if story.OwnerUserID == contributor.UserID {
		return nil
	}
	return p.createNotifications(ctx, []notification_repo.CreateNotificationInput{{
		RecipientUserID:            contributor.UserID,
		RecipientPersonaID:         contributor.ID,
		ActorPersonaID:             &story.OwnerPersonaID,
		ActorPostingMode:           models.PublicPostingMode,
		StoryID:                    &story.ID,
		StoryContributionRequestID: &request.ID,
		NotificationType:           notificationType,
	}})
}

// notifyStoryHighlightChange publishes one highlight add or remove notification to the story owner.
func (p *StoryPipe) notifyStoryHighlightChange(ctx context.Context, story *models.Story, actor *models.Persona, notificationType models.NotificationType) error {
	if story.OwnerUserID == actor.UserID {
		return nil
	}
	return p.createNotifications(ctx, []notification_repo.CreateNotificationInput{{
		RecipientUserID:    story.OwnerUserID,
		RecipientPersonaID: story.OwnerPersonaID,
		ActorPersonaID:     &actor.ID,
		ActorPostingMode:   models.PublicPostingMode,
		EventID:            &story.EventID,
		StoryID:            &story.ID,
		NotificationType:   notificationType,
	}})
}
