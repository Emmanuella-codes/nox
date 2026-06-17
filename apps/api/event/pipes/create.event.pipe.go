package pipes

import (
	"context"
	"strings"

	"github.com/emmanuella-codes/nox/event/dtos"
	"github.com/emmanuella-codes/nox/event/messages"
	"github.com/emmanuella-codes/nox/models"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

func (p *EventPipe) CreateEventPipe(ctx context.Context, userID uuid.UUID, dto dtos.CreateEventDTO) *shared.PipeRes[models.Event] {
	dto.Title = strings.TrimSpace(dto.Title)
	dto.Venue = strings.TrimSpace(dto.Venue)
	dto.Location = strings.TrimSpace(dto.Location)
	dto.Description = strings.TrimSpace(dto.Description)
	dto.CoverURL = strings.TrimSpace(dto.CoverURL)
	dto.TicketURL = strings.TrimSpace(dto.TicketURL)

	persona, err := p.personaRepo.FindPersonaByID(ctx, dto.OrganizerID)
	if err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return shared.PipeError[models.Event](messages.Persona_Not_Found)
		}
		return pipeInternalError[models.Event](err, "event.find_organizer")
	}
	if persona.UserID != userID || !canCreateEvent(persona) {
		return shared.PipeError[models.Event](messages.Forbidden)
	}

	event, err := p.eventRepo.CreateEvent(ctx, dto)
	if err != nil {
		return pipeInternalError[models.Event](err, "event.create")
	}
	return shared.PipeSuccess(messages.Event_Created, event)
}

func canCreateEvent(persona *models.Persona) bool {
	return persona.PersonaType == models.VisiblePersonaType &&
		(persona.Category == models.DJPersonaCategory ||
			persona.Category == models.OrganizerPersonaCategory)
}
