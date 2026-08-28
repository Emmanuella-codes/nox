package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/post/messages"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

// FindViewerPersona verifies that a public profile belongs to the current user.
func (p *PostPipe) FindViewerPersona(ctx context.Context, userID uuid.UUID, personaID uuid.UUID) (*models.Persona, shared.PipeMessage) {
	persona, err := p.personaRepo.FindPersonaByID(ctx, personaID)
	if err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return nil, messages.Persona_Not_Found
		}
		return nil, messages.Internal_Error
	}
	if !persona.IsOwnedBy(userID) {
		return nil, messages.Forbidden
	}
	return persona, ""
}

// publicPostPersona fetches the public profile for a public post response.
func (p *PostPipe) publicPostPersona(ctx context.Context, post *models.Post) (*models.Persona, *shared.PipeRes[PostResponse]) {
	if post.PostingMode != models.PublicPostingMode || post.PersonaID == nil {
		return nil, nil
	}
	persona, err := p.personaRepo.FindPersonaByID(ctx, *post.PersonaID)
	if err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return nil, shared.PipeError[PostResponse](messages.Persona_Not_Found)
		}
		return nil, pipeInternalError[PostResponse](err, "post.find_public_persona")
	}
	return persona, nil
}

// publicPostPersonas fetches the public profiles required by a list of posts.
func (p *PostPipe) publicPostPersonas(ctx context.Context, posts []*models.Post) (map[string]*models.Persona, *shared.PipeRes[[]PostResponse]) {
	personas := make(map[string]*models.Persona)
	for _, post := range posts {
		if post.PostingMode != models.PublicPostingMode || post.PersonaID == nil {
			continue
		}
		personaID := post.PersonaID.String()
		if _, ok := personas[personaID]; ok {
			continue
		}
		persona, err := p.personaRepo.FindPersonaByID(ctx, *post.PersonaID)
		if err != nil {
			if err == persona_repo.ErrPersonaNotFound {
				return nil, shared.PipeError[[]PostResponse](messages.Persona_Not_Found)
			}
			return nil, pipeInternalError[[]PostResponse](err, "post.find_public_persona")
		}
		personas[personaID] = persona
	}
	return personas, nil
}
