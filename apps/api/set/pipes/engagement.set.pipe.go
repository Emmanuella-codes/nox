package pipes

import (
	"context"
	"strings"

	"github.com/emmanuella-codes/nox/models"
	notification_repo "github.com/emmanuella-codes/nox/repositories/notification"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	set_repo "github.com/emmanuella-codes/nox/repositories/set"
	"github.com/emmanuella-codes/nox/set/dtos"
	"github.com/emmanuella-codes/nox/set/messages"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

type setCommentLifecycleRepository interface {
	UpdateSetComment(context.Context, uuid.UUID, uuid.UUID, string) (*models.SetComment, error)
	DeleteSetComment(context.Context, uuid.UUID, uuid.UUID) error
	LikeSetComment(context.Context, uuid.UUID, uuid.UUID) error
	UnlikeSetComment(context.Context, uuid.UUID, uuid.UUID) error
}

func (p *SetPipe) LikeSetPipe(ctx context.Context, userID uuid.UUID, setID uuid.UUID, dto dtos.SetPersonaActionDTO) *shared.PipeRes[any] {
	if message := p.authorizedPersona(ctx, userID, dto.PersonaID); message != "" {
		return shared.PipeError[any](message)
	}
	if _, err := p.setRepo.FindSetByID(ctx, setID); err != nil {
		if err == set_repo.ErrSetNotFound {
			return shared.PipeError[any](messages.Set_Not_Found)
		}
		return pipeInternalError[any](err, "set.like_find")
	}
	if err := p.setRepo.LikeSet(ctx, dto.PersonaID, setID); err != nil {
		return pipeInternalError[any](err, "set.like")
	}
	p.createSetNotification(ctx, userID, dto.PersonaID, setID, models.LikeNotificationType, nil)
	return shared.PipeSuccess[any](messages.Set_Liked, nil)
}

func (p *SetPipe) UnlikeSetPipe(ctx context.Context, userID uuid.UUID, setID uuid.UUID, dto dtos.SetPersonaActionDTO) *shared.PipeRes[any] {
	if message := p.authorizedPersona(ctx, userID, dto.PersonaID); message != "" {
		return shared.PipeError[any](message)
	}
	if _, err := p.setRepo.FindSetByID(ctx, setID); err != nil {
		if err == set_repo.ErrSetNotFound {
			return shared.PipeError[any](messages.Set_Not_Found)
		}
		return pipeInternalError[any](err, "set.unlike_find")
	}
	if err := p.setRepo.UnlikeSet(ctx, dto.PersonaID, setID); err != nil {
		return pipeInternalError[any](err, "set.unlike")
	}
	return shared.PipeSuccess[any](messages.Set_Unliked, nil)
}

func (p *SetPipe) RecordSetPlayPipe(ctx context.Context, userID, setID uuid.UUID) *shared.PipeRes[any] {
	counted, err := p.setRepo.RecordSetPlay(ctx, setID, userID)
	if err != nil {
		if err == set_repo.ErrSetNotFound {
			return shared.PipeError[any](messages.Set_Not_Found)
		}
		return pipeInternalError[any](err, "set.play")
	}
	if !counted {
		return shared.PipeSuccess[any](messages.Set_Play_Already_Recorded, nil)
	}
	return shared.PipeSuccess[any](messages.Set_Play_Recorded, nil)
}

func (p *SetPipe) CreateSetCommentPipe(ctx context.Context, userID uuid.UUID, setID uuid.UUID, dto dtos.CreateSetCommentDTO) *shared.PipeRes[SetCommentResponse] {
	body := strings.TrimSpace(dto.Body)
	if body == "" || len(body) > 280 {
		return shared.PipeError[SetCommentResponse](messages.Invalid_Payload)
	}
	if message := p.authorizedPersona(ctx, userID, dto.PersonaID); message != "" {
		return shared.PipeError[SetCommentResponse](message)
	}
	if _, err := p.setRepo.FindSetByID(ctx, setID); err != nil {
		if err == set_repo.ErrSetNotFound {
			return shared.PipeError[SetCommentResponse](messages.Set_Not_Found)
		}
		return pipeInternalError[SetCommentResponse](err, "set.comment_find")
	}
	comment, err := p.setRepo.CreateSetComment(ctx, dto.PersonaID, setID, body, dto.ParentID)
	if err != nil {
		return pipeInternalError[SetCommentResponse](err, "set.comment")
	}
	p.createSetNotification(ctx, userID, dto.PersonaID, setID, models.CommentNotificationType, &comment.ID)
	if err := p.hydrateSetComments(ctx, []*models.SetComment{comment}); err != nil {
		return pipeInternalError[SetCommentResponse](err, "set.comment_hydrate")
	}
	response := setCommentResponse(comment)
	return shared.PipeSuccess(messages.Set_Commented, &response)
}

func (p *SetPipe) createSetNotification(ctx context.Context, actorUserID, actorPersonaID, setID uuid.UUID, kind models.NotificationType, commentID *uuid.UUID) {
	if p.notificationRepo == nil {
		return
	}
	set, err := p.setRepo.FindSetByID(ctx, setID)
	if err != nil {
		return
	}
	if set.AuthorUserID == actorUserID {
		return
	}
	ownerPersona, err := p.personaRepo.FindPersonaByID(ctx, set.PersonaID)
	if err != nil {
		return
	}
	input := notification_repo.CreateNotificationInput{RecipientUserID: set.AuthorUserID, RecipientPersonaID: ownerPersona.ID, ActorPersonaID: &actorPersonaID, ActorPostingMode: models.PublicPostingMode, SetID: &setID, NotificationType: kind}
	input.SetCommentID = commentID
	created, err := p.notificationRepo.CreateNotifications(ctx, []notification_repo.CreateNotificationInput{input})
	if err == nil && len(created) > 0 && p.notificationPublisher != nil {
		p.notificationPublisher.PublishCreatedNotification(ctx, created[0])
	}
}

func (p *SetPipe) ListSetCommentsPipe(ctx context.Context, setID uuid.UUID, limit int, offset int) *shared.PipeRes[SetCommentListResponse] {
	limit = normalizeLimit(limit)
	offset = normalizeOffset(offset)
	comments, err := p.setRepo.FindSetComments(ctx, setID, limit+1, offset)
	if err != nil {
		return pipeInternalError[SetCommentListResponse](err, "set.comments")
	}
	hasMore := len(comments) > limit
	if hasMore {
		comments = comments[:limit]
	}
	if err := p.hydrateSetComments(ctx, comments); err != nil {
		return pipeInternalError[SetCommentListResponse](err, "set.comments_hydrate")
	}
	responses := make([]SetCommentResponse, 0, len(comments))
	for _, comment := range comments {
		responses = append(responses, setCommentResponse(comment))
	}
	return shared.PipeSuccess(messages.Set_Comments_Listed, &SetCommentListResponse{
		Limit:      limit,
		Offset:     offset,
		HasMore:    hasMore,
		NextOffset: nextOffset(limit, offset, hasMore),
		Comments:   responses,
	})
}

func (p *SetPipe) UpdateSetCommentPipe(ctx context.Context, userID, commentID uuid.UUID, dto dtos.UpdateSetCommentDTO) *shared.PipeRes[SetCommentResponse] {
	repo, ok := p.setRepo.(setCommentLifecycleRepository)
	if !ok {
		return pipeInternalError[SetCommentResponse](set_repo.ErrSetCommentNotFound, "set.comment_update")
	}
	if message := p.authorizedPersona(ctx, userID, dto.PersonaID); message != "" {
		return shared.PipeError[SetCommentResponse](message)
	}
	body := strings.TrimSpace(dto.Body)
	if body == "" || len(body) > 280 {
		return shared.PipeError[SetCommentResponse](messages.Invalid_Payload)
	}
	comment, err := repo.UpdateSetComment(ctx, dto.PersonaID, commentID, body)
	if err != nil {
		return pipeInternalError[SetCommentResponse](err, "set.comment_update")
	}
	if err := p.hydrateSetComments(ctx, []*models.SetComment{comment}); err != nil {
		return pipeInternalError[SetCommentResponse](err, "set.comment_update_hydrate")
	}
	response := setCommentResponse(comment)
	return shared.PipeSuccess(messages.Set_Comment_Updated, &response)
}

func (p *SetPipe) DeleteSetCommentPipe(ctx context.Context, userID, commentID, personaID uuid.UUID) *shared.PipeRes[any] {
	repo, ok := p.setRepo.(setCommentLifecycleRepository)
	if !ok {
		return pipeInternalError[any](set_repo.ErrSetCommentNotFound, "set.comment_delete")
	}
	if message := p.authorizedPersona(ctx, userID, personaID); message != "" {
		return shared.PipeError[any](message)
	}
	if err := repo.DeleteSetComment(ctx, personaID, commentID); err != nil {
		return pipeInternalError[any](err, "set.comment_delete")
	}
	return shared.PipeSuccess[any](messages.Set_Comment_Deleted, nil)
}

func (p *SetPipe) SetCommentLikePipe(ctx context.Context, userID, commentID, personaID uuid.UUID, unlike bool) *shared.PipeRes[any] {
	repo, ok := p.setRepo.(setCommentLifecycleRepository)
	if !ok {
		return pipeInternalError[any](set_repo.ErrSetCommentNotFound, "set.comment_like")
	}
	if message := p.authorizedPersona(ctx, userID, personaID); message != "" {
		return shared.PipeError[any](message)
	}
	var err error
	if unlike {
		err = repo.UnlikeSetComment(ctx, personaID, commentID)
	} else {
		err = repo.LikeSetComment(ctx, personaID, commentID)
	}
	if err != nil {
		return pipeInternalError[any](err, "set.comment_like")
	}
	if unlike {
		return shared.PipeSuccess[any](messages.Set_Comment_Unliked, nil)
	}
	return shared.PipeSuccess[any](messages.Set_Comment_Liked, nil)
}

func (p *SetPipe) authorizedPersona(ctx context.Context, userID uuid.UUID, personaID uuid.UUID) shared.PipeMessage {
	persona, err := p.personaRepo.FindPersonaByID(ctx, personaID)
	if err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return messages.Persona_Not_Found
		}
		return messages.Internal_Error
	}
	if persona.UserID != userID {
		return messages.Forbidden
	}
	return ""
}
