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

func (p *CommentPipe) CreateCommentPipe(ctx context.Context, userID uuid.UUID, postID uuid.UUID, dto dtos.CreateCommentDTO) *shared.PipeRes[CommentResponse] {
	dto.Body = strings.TrimSpace(dto.Body)
	if dto.PostingMode == "" {
		dto.PostingMode = models.PublicPostingMode
	}
	if !validPostingMode(dto.PostingMode) {
		return shared.PipeError[CommentResponse](messages.Invalid_Payload)
	}

	if _, err := p.postRepo.FindPostByID(ctx, postID); err != nil {
		if err == post_repo.ErrPostNotFound {
			return shared.PipeError[CommentResponse](messages.Post_Not_Found)
		}
		return pipeInternalError[CommentResponse](err, "comment.find_post")
	}

	persona, err := p.personaRepo.FindPersonaByID(ctx, dto.PersonaID)
	if err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return shared.PipeError[CommentResponse](messages.Persona_Not_Found)
		}
		return pipeInternalError[CommentResponse](err, "comment.find_persona")
	}
	if persona.UserID != userID || persona.PersonaType != models.VisiblePersonaType {
		return shared.PipeError[CommentResponse](messages.Forbidden)
	}

	comment, err := p.commentRepo.CreateComment(ctx, postID, dto)
	if err != nil {
		return pipeInternalError[CommentResponse](err, "comment.create")
	}

	var publicPersona *models.Persona
	var identity *models.AnonymousThreadIdentity
	if comment.PostingMode == models.PublicPostingMode {
		publicPersona = persona
	} else {
		identity, err = p.postRepo.EnsureAnonymousThreadIdentity(ctx, postID, userID, dto.PersonaID, anonymousHandle())
		if err != nil {
			return pipeInternalError[CommentResponse](err, "comment.anonymous_identity")
		}
	}

	response := commentResponse(comment, publicPersona, identity)
	return shared.PipeSuccess(messages.Comment_Created, &response)
}
