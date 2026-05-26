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

func (p *PostPipe) CreatePostPipe(ctx context.Context, userID uuid.UUID, dto dtos.CreatePostDTO) *shared.PipeRes[PostResponse] {
	dto.Body = strings.TrimSpace(dto.Body)
	dto.MediaURL = strings.TrimSpace(dto.MediaURL)
	dto.Location = strings.TrimSpace(dto.Location)

	if !validPostingMode(dto.PostingMode) {
		return shared.PipeError[PostResponse](messages.Invalid_Posting_Mode)
	}

	var persona *models.Persona
	switch dto.PostingMode {
	case models.PublicPostingMode:
		if dto.PersonaID == nil || *dto.PersonaID == uuid.Nil {
			return shared.PipeError[PostResponse](messages.Persona_Required)
		}

		foundPersona, err := p.personaRepo.FindPersonaByID(ctx, *dto.PersonaID)
		if err != nil {
			if err == persona_repo.ErrPersonaNotFound {
				return shared.PipeError[PostResponse](messages.Persona_Not_Found)
			}
			return pipeInternalError[PostResponse](err, "post.find_persona")
		}
		if foundPersona.UserID != userID || foundPersona.PersonaType != models.VisiblePersonaType {
			return shared.PipeError[PostResponse](messages.Forbidden)
		}
		persona = foundPersona
	case models.AnonymousPostingMode:
		dto.PersonaID = nil
	}

	post, err := p.postRepo.CreatePost(ctx, userID, dto)
	if err != nil {
		return pipeInternalError[PostResponse](err, "post.create")
	}

	response := postResponse(post, persona)
	return shared.PipeSuccess(messages.Post_Created, &response)
}
