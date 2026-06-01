package pipes

import (
	"context"
	"strings"

	"github.com/emmanuella-codes/nox/models"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	set_repo "github.com/emmanuella-codes/nox/repositories/set"
	"github.com/emmanuella-codes/nox/set/dtos"
	"github.com/emmanuella-codes/nox/set/messages"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

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

func (p *SetPipe) RecordSetPlayPipe(ctx context.Context, setID uuid.UUID) *shared.PipeRes[any] {
	if err := p.setRepo.IncrementPlayCount(ctx, setID); err != nil {
		if err == set_repo.ErrSetNotFound {
			return shared.PipeError[any](messages.Set_Not_Found)
		}
		return pipeInternalError[any](err, "set.play")
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
	if err := p.hydrateSetComments(ctx, []*models.SetComment{comment}); err != nil {
		return pipeInternalError[SetCommentResponse](err, "set.comment_hydrate")
	}
	response := setCommentResponse(comment)
	return shared.PipeSuccess(messages.Set_Commented, &response)
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
