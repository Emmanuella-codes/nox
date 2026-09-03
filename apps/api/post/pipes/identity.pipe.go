package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/post/messages"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	preference_repo "github.com/emmanuella-codes/nox/repositories/preference"
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

func (p *PostPipe) blockedForViewer(ctx context.Context, viewerPersonaID uuid.UUID, targetUserID uuid.UUID) (bool, error) {
	if p.preferenceRepo == nil {
		return false, nil
	}
	viewer, err := p.personaRepo.FindPersonaByID(ctx, viewerPersonaID)
	if err != nil {
		return false, err
	}
	if viewer.UserID == targetUserID {
		return false, nil
	}
	return isBlockedBetweenUsers(ctx, p.preferenceRepo, viewer.UserID, targetUserID)
}

func isBlockedBetweenUsers(ctx context.Context, preferenceRepo preference_repo.PreferenceRepository, viewerUserID uuid.UUID, targetUserID uuid.UUID) (bool, error) {
	return preferenceRepo.IsBlockedBetween(ctx, viewerUserID, targetUserID)
}
