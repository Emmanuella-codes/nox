package pipes

import (
	"context"
	"strings"

	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/post/dtos"
	"github.com/emmanuella-codes/nox/post/messages"
	hashtag_repo "github.com/emmanuella-codes/nox/repositories/hashtag"
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
	case models.PublicPostingMode, models.AnonymousPostingMode:
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
		if dto.PostingMode == models.PublicPostingMode {
			persona = foundPersona
		}
	}

	tags := hashtag_repo.ExtractTags(dto.Body)
	post, err := p.postRepo.CreatePostWithHashtags(ctx, userID, dto, tags)
	if err != nil {
		return pipeInternalError[PostResponse](err, "post.create")
	}

	var identity *models.AnonymousThreadIdentity
	if post.PostingMode == models.AnonymousPostingMode && post.PersonaID != nil {
		identity, err = p.postRepo.EnsureAnonymousThreadIdentity(ctx, post.ID, userID, *post.PersonaID, anonymousHandle())
		if err != nil {
			return pipeInternalError[PostResponse](err, "post.anonymous_identity")
		}
	}

	response := postResponse(post, persona, identity)
	response.Hashtags = tags
	return shared.PipeSuccess(messages.Post_Created, &response)
}
