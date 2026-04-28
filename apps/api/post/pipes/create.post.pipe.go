package pipes

import (
	"context"
	"strings"

	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/post/dtos"
	"github.com/emmanuella-codes/nox/post/messages"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

func (p *PostPipe) CreatePostPipe(ctx context.Context, userID uuid.UUID, dto dtos.CreatePostDTO) *shared.PipeRes[models.Post] {
	persona, err := p.personaRepo.FindPersonaByID(ctx, dto.PersonaID)
	if err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return pipeError[models.Post](messages.Persona_Not_Found)
		}
		return pipeInternalError[models.Post](err, "post.find_persona")
	}
	if persona.UserID != userID {
		return pipeError[models.Post](messages.Forbidden)
	}

	dto.Body = strings.TrimSpace(dto.Body)
	dto.MediaURL = strings.TrimSpace(dto.MediaURL)
	dto.Location = strings.TrimSpace(dto.Location)

	post, err := p.postRepo.CreatePost(ctx, dto.PersonaID, dto)
	if err != nil {
		return pipeInternalError[models.Post](err, "post.create")
	}

	return pipeSuccess(messages.Post_Created, post)
}
