package pipes

import (
	"context"
	"strings"

	"github.com/emmanuella-codes/nox/comment/dtos"
	"github.com/emmanuella-codes/nox/comment/messages"
	"github.com/emmanuella-codes/nox/models"
	notification_repo "github.com/emmanuella-codes/nox/repositories/notification"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	post_repo "github.com/emmanuella-codes/nox/repositories/post"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/emmanuella-codes/nox/shared/anonymousidentity"
	"github.com/google/uuid"
)

// CreateCommentPipe validates the author profile and persists a new comment.
func (p *CommentPipe) CreateCommentPipe(ctx context.Context, userID uuid.UUID, postID uuid.UUID, dto dtos.CreateCommentDTO) *shared.PipeRes[CommentResponse] {
	dto.Body = strings.TrimSpace(dto.Body)
	dto.AuthorUserID = userID
	if dto.PostingMode == "" {
		dto.PostingMode = models.PublicPostingMode
	}
	if !validPostingMode(dto.PostingMode) {
		return shared.PipeError[CommentResponse](messages.Invalid_Payload)
	}
	post, err := p.postRepo.FindPostByID(ctx, postID)
	if err != nil {
		if err == post_repo.ErrPostNotFound {
			return shared.PipeError[CommentResponse](messages.Post_Not_Found)
		}
		return pipeInternalError[CommentResponse](err, "comment.find_post")
	}
	profile, err := p.personaRepo.FindPersonaByID(ctx, dto.PersonaID)
	if err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return shared.PipeError[CommentResponse](messages.Persona_Not_Found)
		}
		return pipeInternalError[CommentResponse](err, "comment.find_persona")
	}
	if !profile.IsOwnedBy(userID) {
		return shared.PipeError[CommentResponse](messages.Forbidden)
	}
	comment, err := p.commentRepo.CreateComment(ctx, postID, dto)
	if err != nil {
		return pipeInternalError[CommentResponse](err, "comment.create")
	}
	var publicProfile *models.Persona
	var identity *models.AnonymousThreadIdentity
	if comment.PostingMode == models.PublicPostingMode {
		publicProfile = profile
	} else {
		identity, err = p.postRepo.EnsureAnonymousThreadIdentity(
			ctx,
			postID,
			userID,
			dto.PersonaID,
			anonymousidentity.GenerateHandle(),
			anonymousidentity.GenerateAvatarKey(),
		)
		if err != nil {
			return pipeInternalError[CommentResponse](err, "comment.anonymous_identity")
		}
	}
	response := commentResponse(comment, publicProfile, identity)
	p.createCommentNotification(ctx, post, profile, comment, identity)
	return shared.PipeSuccess(messages.Comment_Created, &response)
}

// createCommentNotification persists and publishes one post-comment notification.
func (p *CommentPipe) createCommentNotification(ctx context.Context, post *models.Post, actor *models.Persona, comment *models.Comment, identity *models.AnonymousThreadIdentity) {
	if p.notificationRepo == nil || post == nil || post.PersonaID == nil || actor == nil || comment == nil || post.AuthorUserID == actor.UserID {
		return
	}
	input := notification_repo.CreateNotificationInput{
		RecipientUserID:    post.AuthorUserID,
		RecipientPersonaID: *post.PersonaID,
		PostID:             &post.ID,
		CommentID:          &comment.ID,
		NotificationType:   models.CommentNotificationType,
	}
	if comment.PostingMode == models.AnonymousPostingMode {
		input.ActorPostingMode = models.AnonymousPostingMode
		input.ActorAnonymousHandle = anonymousHandleValue(identity)
		input.ActorAnonymousAvatarKey = anonymousAvatarValue(identity)
	} else {
		input.ActorPersonaID = &actor.ID
		input.ActorPostingMode = models.PublicPostingMode
	}
	created, err := p.notificationRepo.CreateNotifications(ctx, []notification_repo.CreateNotificationInput{input})
	if err != nil || len(created) == 0 || p.notificationPublisher == nil {
		return
	}
	p.notificationPublisher.PublishCreatedNotification(ctx, created[0])
}
