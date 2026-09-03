package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/notification/messages"
	notification_repo "github.com/emmanuella-codes/nox/repositories/notification"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/emmanuella-codes/nox/shared/realtime"
	"github.com/google/uuid"
)

const timeFormat = "2006-01-02T15:04:05.999999999Z07:00"

type NotificationPipe struct {
	notificationRepo notification_repo.NotificationRepository
	personaRepo      persona_repo.PersonaRepository
	realtimeHub      *realtime.Hub
}

type ActorPersonaResponse struct {
	ID          string `json:"id"`
	Handle      string `json:"handle"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
}

type ActorAnonymousResponse struct {
	Handle    string `json:"handle"`
	AvatarKey string `json:"avatar_key"`
}

type NotificationActorResponse struct {
	Mode      models.PostingMode      `json:"mode"`
	Persona   *ActorPersonaResponse   `json:"persona,omitempty"`
	Anonymous *ActorAnonymousResponse `json:"anonymous,omitempty"`
}

type NotificationResponse struct {
	ID                         string                    `json:"id"`
	PersonaID                  string                    `json:"persona_id"`
	Actor                      NotificationActorResponse `json:"actor"`
	ConversationID             *string                   `json:"conversation_id,omitempty"`
	MessageID                  *string                   `json:"message_id,omitempty"`
	PostID                     *string                   `json:"post_id,omitempty"`
	CommentID                  *string                   `json:"comment_id,omitempty"`
	EventID                    *string                   `json:"event_id,omitempty"`
	StoryID                    *string                   `json:"story_id,omitempty"`
	StoryItemID                *string                   `json:"story_item_id,omitempty"`
	StoryContributionRequestID *string                   `json:"story_contribution_request_id,omitempty"`
	IsRead                     bool                      `json:"is_read"`
	ReadAt                     *string                   `json:"read_at,omitempty"`
	NotificationType           models.NotificationType   `json:"notification_type"`
	CreatedAt                  string                    `json:"created_at"`
}

type NotificationListResponse struct {
	Limit         int                    `json:"limit"`
	Offset        int                    `json:"offset"`
	HasMore       bool                   `json:"has_more"`
	NextOffset    *int                   `json:"next_offset,omitempty"`
	UnreadCount   int                    `json:"unread_count"`
	Notifications []NotificationResponse `json:"notifications"`
}

type NotificationUnreadCountResponse struct {
	UnreadCount int `json:"unread_count"`
}

func NewNotificationPipe(notificationRepo notification_repo.NotificationRepository, personaRepo persona_repo.PersonaRepository, deps ...any) *NotificationPipe {
	pipe := &NotificationPipe{notificationRepo: notificationRepo, personaRepo: personaRepo}
	for _, dep := range deps {
		if hub, ok := dep.(*realtime.Hub); ok {
			pipe.realtimeHub = hub
		}
	}
	return pipe
}

func pipeInternalError[T any](err error, operation string) *shared.PipeRes[T] {
	return shared.PipeInternalError[T](err, "notification", operation, messages.Internal_Error)
}

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

func (p *NotificationPipe) notificationResponse(ctx context.Context, notification *models.Notification) NotificationResponse {
	var conversationID *string
	var messageID *string
	var postID *string
	var commentID *string
	var eventID *string
	var storyID *string
	var storyItemID *string
	var storyContributionRequestID *string
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
	if notification.EventID != nil {
		value := notification.EventID.String()
		eventID = &value
	}
	if notification.StoryID != nil {
		value := notification.StoryID.String()
		storyID = &value
	}
	if notification.StoryItemID != nil {
		value := notification.StoryItemID.String()
		storyItemID = &value
	}
	if notification.StoryContributionRequestID != nil {
		value := notification.StoryContributionRequestID.String()
		storyContributionRequestID = &value
	}
	if notification.ReadAt != nil {
		value := notification.ReadAt.Format(timeFormat)
		readAt = &value
	}
	response := NotificationResponse{
		ID:                         notification.ID.String(),
		PersonaID:                  notification.RecipientPersonaID.String(),
		Actor:                      NotificationActorResponse{Mode: notification.ActorPostingMode},
		ConversationID:             conversationID,
		MessageID:                  messageID,
		PostID:                     postID,
		CommentID:                  commentID,
		EventID:                    eventID,
		StoryID:                    storyID,
		StoryItemID:                storyItemID,
		StoryContributionRequestID: storyContributionRequestID,
		IsRead:                     notification.IsRead,
		ReadAt:                     readAt,
		NotificationType:           notification.NotificationType,
		CreatedAt:                  notification.CreatedAt.Format(timeFormat),
	}
	if notification.ActorPostingMode == models.AnonymousPostingMode {
		response.Actor.Anonymous = &ActorAnonymousResponse{
			Handle:    notification.ActorAnonymousHandle,
			AvatarKey: notification.ActorAnonymousAvatarKey,
		}
		return response
	}
	if notification.ActorPersonaID == nil {
		return response
	}
	persona, err := p.personaRepo.FindPersonaByID(ctx, *notification.ActorPersonaID)
	if err == nil {
		response.Actor.Persona = &ActorPersonaResponse{
			ID:          persona.ID.String(),
			Handle:      persona.Handle,
			DisplayName: persona.DisplayName,
			AvatarURL:   persona.AvatarURL,
		}
	}
	return response
}

func (p *NotificationPipe) notificationResponses(ctx context.Context, notifications []*models.Notification) []NotificationResponse {
	responses := make([]NotificationResponse, 0, len(notifications))
	for _, notification := range notifications {
		responses = append(responses, p.notificationResponse(ctx, notification))
	}
	return responses
}
