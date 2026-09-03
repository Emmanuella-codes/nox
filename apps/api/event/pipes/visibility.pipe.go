package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/event/messages"
	"github.com/emmanuella-codes/nox/models"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

func (p *EventPipe) FindViewerPersona(ctx context.Context, userID uuid.UUID, personaID uuid.UUID) (*models.Persona, shared.PipeMessage) {
	persona, err := p.personaRepo.FindPersonaByID(ctx, personaID)
	if err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return nil, messages.Persona_Not_Found
		}
		return nil, messages.Internal_Error
	}
	if persona.UserID != userID {
		return nil, messages.Forbidden
	}
	return persona, ""
}

func (p *EventPipe) viewerCanSeeEvent(ctx context.Context, viewerPersonaID uuid.UUID, event *models.Event) (bool, error) {
	if p.preferenceRepo == nil {
		return true, nil
	}
	viewer, err := p.personaRepo.FindPersonaByID(ctx, viewerPersonaID)
	if err != nil {
		return false, err
	}
	organizer, err := p.personaRepo.FindPersonaByID(ctx, event.OrganizerID)
	if err != nil {
		return false, err
	}
	if viewer.UserID == organizer.UserID {
		return true, nil
	}
	blocked, err := p.preferenceRepo.IsBlockedBetween(ctx, viewer.UserID, organizer.UserID)
	if err != nil {
		return false, err
	}
	return !blocked, nil
}
