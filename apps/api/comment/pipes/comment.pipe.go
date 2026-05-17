package pipes

import (
	"context"
	"strings"

	"github.com/emmanuella-codes/nox/comment/dtos"
	"github.com/emmanuella-codes/nox/comment/messages"
	"github.com/emmanuella-codes/nox/models"
	comment_repo "github.com/emmanuella-codes/nox/repositories/comment"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	post_repo "github.com/emmanuella-codes/nox/repositories/post"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

func (p *CommentPipe) CreateCommentPipe(ctx context.Context, userID uuid.UUID, postID uuid.UUID, dto dtos.CreateCommentDTO) *shared.PipeRes[models.Comment] {
	dto.Body = strings.TrimSpace(dto.Body)
	if _, err := p.postRepo.FindPostByID(ctx, postID); err != nil {
		if err == post_repo.ErrPostNotFound {
			return pipeError[models.Comment](messages.Post_Not_Found)
		}
		return pipeInternalError[models.Comment](err, "comment.find_post")
	}

	persona, err := p.personaRepo.FindPersonaByID(ctx, dto.PersonaID)
	if err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return pipeError[models.Comment](messages.Persona_Not_Found)
		}
		return pipeInternalError[models.Comment](err, "comment.find_persona")
	}
	if persona.UserID != userID || persona.PersonaType != models.VisiblePersonaType {
		return pipeError[models.Comment](messages.Forbidden)
	}

	comment, err := p.commentRepo.CreateComment(ctx, postID, dto)
	if err != nil {
		return pipeInternalError[models.Comment](err, "comment.create")
	}

	return pipeSuccess(messages.Comment_Created, comment)
}

func (p *CommentPipe) ListCommentsPipe(ctx context.Context, postID uuid.UUID, limit int) *shared.PipeRes[[]*models.Comment] {
	if _, err := p.postRepo.FindPostByID(ctx, postID); err != nil {
		if err == post_repo.ErrPostNotFound {
			return pipeError[[]*models.Comment](messages.Post_Not_Found)
		}
		return pipeInternalError[[]*models.Comment](err, "comment.find_post")
	}

	comments, err := p.commentRepo.FindCommentsByPostID(ctx, postID, limit)
	if err != nil {
		return pipeInternalError[[]*models.Comment](err, "comment.list")
	}

	return pipeSuccess(messages.Comments_Listed, &comments)
}

func (p *CommentPipe) DeleteCommentPipe(ctx context.Context, userID uuid.UUID, commentID uuid.UUID) *shared.PipeRes[any] {
	comment, err := p.commentRepo.FindCommentByID(ctx, commentID)
	if err != nil {
		if err == comment_repo.ErrCommentNotFound {
			return pipeError[any](messages.Comment_Not_Found)
		}
		return pipeInternalError[any](err, "comment.find_for_delete")
	}

	persona, err := p.personaRepo.FindPersonaByID(ctx, comment.PersonaID)
	if err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return pipeError[any](messages.Persona_Not_Found)
		}
		return pipeInternalError[any](err, "comment.find_persona_for_delete")
	}
	if persona.UserID != userID {
		return pipeError[any](messages.Forbidden)
	}

	if err := p.commentRepo.DeleteComment(ctx, commentID); err != nil {
		if err == comment_repo.ErrCommentNotFound {
			return pipeError[any](messages.Comment_Not_Found)
		}
		return pipeInternalError[any](err, "comment.delete")
	}

	return pipeSuccess[any](messages.Comment_Deleted, nil)
}
