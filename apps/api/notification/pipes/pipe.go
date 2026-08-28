package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/notification/messages"
	notification_repo "github.com/emmanuella-codes/nox/repositories/notification"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

const timeFormat = "2006-01-02T15:04:05.999999999Z07:00"

type NotificationPipe struct {
	notificationRepo notification_repo.NotificationRepository
	personaRepo      persona_repo.PersonaRepository
}

type ActorPersonaResponse struct {
	ID          string `json:"id"`
	Handle      string `json:"handle"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
}

type NotificationResponse struct {
	ID               string                  `json:"id"`
	PersonaID        string                  `json:"persona_id"`
	ActorPersonaID   string                  `json:"actor_persona_id"`
	ActorPersona     *ActorPersonaResponse   `json:"actor_persona,omitempty"`
	ConversationID   *string                 `json:"conversation_id,omitempty"`
	MessageID        *string                 `json:"message_id,omitempty"`
	PostID           *string                 `json:"post_id,omitempty"`
	CommentID        *string                 `json:"comment_id,omitempty"`
	IsRead           bool                    `json:"is_read"`
	ReadAt           *string                 `json:"read_at,omitempty"`
	NotificationType models.NotificationType `json:"notification_type"`
	CreatedAt        string                  `json:"created_at"`
}

// NewNotificationPipe builds the notification orchestration layer from repositories.
func NewNotificationPipe(notificationRepo notification_repo.NotificationRepository, personaRepo persona_repo.PersonaRepository) *NotificationPipe {
	return &NotificationPipe{notificationRepo: notificationRepo, personaRepo: personaRepo}
}

// pipeInternalError maps internal notification errors to pipe responses.
func pipeInternalError[T any](err error, operation string) *shared.PipeRes[T] {
	return shared.PipeInternalError[T](err, "notification", operation, messages.Internal_Error)
}

// profilePersona validates that one persona belongs to the current user.
func (p *NotificationPipe) profilePersona(ctx context.Context, userID uuid.UUID, personaID uuid.UUID) (*models.Persona, shared.PipeMessage) {
	persona, err := p.personaRepo.FindPersonaByID(ctx, personaID)
	if err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return nil, messages.Persona_Not_Found
		}
		return nil, messages.Internal_Error
	}
	if persona.UserID != userID {
		return nil, messages.Forbidden
	}
	return persona, ""
}

// notificationResponse maps one notification into the API response shape.
func (p *NotificationPipe) notificationResponse(ctx context.Context, notification *models.Notification) NotificationResponse {
	var conversationID *string
	var messageID *string
	var postID *string
	var commentID *string
	var readAt *string
	if notification.ConversationID != nil {
		value := notification.ConversationID.String()
		conversationID = &value
	}
	if notification.MessageID != nil {
		value := notification.MessageID.String()
		messageID = &value
	}
	if notification.PostID != nil {
		value := notification.PostID.String()
		postID = &value
	}
	if notification.CommentID != nil {
		value := notification.CommentID.String()
		commentID = &value
	}
	if notification.ReadAt != nil {
		value := notification.ReadAt.Format(timeFormat)
		readAt = &value
	}
	response := NotificationResponse{
		ID:               notification.ID.String(),
		PersonaID:        notification.RecipientPersonaID.String(),
		ActorPersonaID:   notification.ActorPersonaID.String(),
		ConversationID:   conversationID,
		MessageID:        messageID,
		PostID:           postID,
		CommentID:        commentID,
		IsRead:           notification.IsRead,
		ReadAt:           readAt,
		NotificationType: notification.NotificationType,
		CreatedAt:        notification.CreatedAt.Format(timeFormat),
	}
	persona, err := p.personaRepo.FindPersonaByID(ctx, notification.ActorPersonaID)
	if err == nil {
		response.ActorPersona = &ActorPersonaResponse{
			ID:          persona.ID.String(),
			Handle:      persona.Handle,
			DisplayName: persona.DisplayName,
			AvatarURL:   persona.AvatarURL,
		}
	}
	return response
}

// notificationResponses maps many notifications into API response shape.
func (p *NotificationPipe) notificationResponses(ctx context.Context, notifications []*models.Notification) []NotificationResponse {
	responses := make([]NotificationResponse, 0, len(notifications))
	for _, notification := range notifications {
		responses = append(responses, p.notificationResponse(ctx, notification))
	}
	return responses
}
