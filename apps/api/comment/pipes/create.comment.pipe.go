package pipes

import (
	"context"
	"strings"

	"github.com/emmanuella-codes/nox/comment/dtos"
	"github.com/emmanuella-codes/nox/comment/messages"
	"github.com/emmanuella-codes/nox/models"
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
