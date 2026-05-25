package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/comment/messages"
	comment_repo "github.com/emmanuella-codes/nox/repositories/comment"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

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
