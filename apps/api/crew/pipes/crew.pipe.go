package pipes

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/emmanuella-codes/nox/crew/dtos"
	"github.com/emmanuella-codes/nox/crew/messages"
	"github.com/emmanuella-codes/nox/models"
	crew_repo "github.com/emmanuella-codes/nox/repositories/crew"
	event_repo "github.com/emmanuella-codes/nox/repositories/event"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

func (p *CrewPipe) CreateCrewPipe(ctx context.Context, userID uuid.UUID, eventID uuid.UUID, dto dtos.CreateCrewDTO) *shared.PipeRes[CrewResponse] {
	dto.Name = strings.TrimSpace(dto.Name)
	if dto.Name == "" || (dto.Visibility != "" && dto.Visibility != models.InviteCodeCrewVisibility) {
		return shared.PipeError[CrewResponse](messages.Invalid_Payload)
	}
	persona, message := p.visiblePersona(ctx, userID, dto.OwnerPersonaID)
	if message != "" {
		return shared.PipeError[CrewResponse](message)
	}
	event, err := p.eventRepo.FindEventByID(ctx, eventID)
	if err != nil {
		if errors.Is(err, event_repo.ErrEventNotFound) {
			return shared.PipeError[CrewResponse](messages.Event_Not_Found)
		}
		return pipeInternalError[CrewResponse](err, "crew.find_event")
	}
	if time.Now().After(crewExpiresAt(event)) {
		return shared.PipeError[CrewResponse](messages.Forbidden)
	}

	var crew *models.EventCrew
	for i := 0; i < 5; i++ {
		crew, err = p.crewRepo.CreateCrew(ctx, userID, eventID, joinCode(), crewExpiresAt(event), dto)
		if !errors.Is(err, crew_repo.ErrCrewCodeTaken) {
			break
		}
	}
	if err != nil {
		return pipeInternalError[CrewResponse](err, "crew.create")
	}
	members := []*models.EventCrewMember{{CrewID: crew.ID, UserID: userID, PersonaID: persona.ID, Role: models.OwnerCrewMemberRole, JoinedAt: time.Now()}}
	response, err := p.crewResponse(ctx, crew, members, nil)
	if err != nil {
		return pipeInternalError[CrewResponse](err, "crew.response")
	}
	return shared.PipeSuccess(messages.Crew_Created, response)
}

func (p *CrewPipe) JoinCrewPipe(ctx context.Context, userID uuid.UUID, dto dtos.JoinCrewDTO) *shared.PipeRes[CrewResponse] {
	dto.JoinCode = strings.ToUpper(strings.TrimSpace(dto.JoinCode))
	persona, message := p.visiblePersona(ctx, userID, dto.PersonaID)
	if message != "" {
		return shared.PipeError[CrewResponse](message)
	}
	crew, err := p.crewRepo.FindCrewByJoinCode(ctx, dto.JoinCode)
	if err != nil {
		if errors.Is(err, crew_repo.ErrCrewNotFound) {
			return shared.PipeError[CrewResponse](messages.Crew_Not_Found)
		}
		return pipeInternalError[CrewResponse](err, "crew.find_code")
	}
	if crew.Visibility != models.InviteCodeCrewVisibility || !activeCrew(crew) {
		return shared.PipeError[CrewResponse](messages.Forbidden)
	}
	if _, err := p.crewRepo.JoinCrew(ctx, crew, persona); err != nil {
		if errors.Is(err, crew_repo.ErrCrewFull) {
			return shared.PipeError[CrewResponse](messages.Crew_Full)
		}
		return pipeInternalError[CrewResponse](err, "crew.join")
	}
	return p.GetCrewPipe(ctx, userID, crew.ID, persona.ID)
}

func (p *CrewPipe) ListMyEventCrewsPipe(ctx context.Context, userID uuid.UUID, eventID uuid.UUID, personaID uuid.UUID, limit int, offset int) *shared.PipeRes[[]CrewResponse] {
	if _, message := p.visiblePersona(ctx, userID, personaID); message != "" {
		return shared.PipeError[[]CrewResponse](message)
	}
	if _, err := p.eventRepo.FindEventByID(ctx, eventID); err != nil {
		if errors.Is(err, event_repo.ErrEventNotFound) {
			return shared.PipeError[[]CrewResponse](messages.Event_Not_Found)
		}
		return pipeInternalError[[]CrewResponse](err, "crew.list_event")
	}
	crews, err := p.crewRepo.FindEventCrewsForPersona(ctx, eventID, personaID, limit, offset)
	if err != nil {
		return pipeInternalError[[]CrewResponse](err, "crew.list")
	}
	responses := make([]CrewResponse, 0, len(crews))
	for _, crew := range crews {
		members, err := p.crewRepo.FindCrewMembers(ctx, crew.ID)
		if err != nil {
			return pipeInternalError[[]CrewResponse](err, "crew.list_members")
		}
		response, err := p.crewResponse(ctx, crew, members, nil)
		if err != nil {
			return pipeInternalError[[]CrewResponse](err, "crew.list_response")
		}
		responses = append(responses, *response)
	}
	return shared.PipeSuccess(messages.Crews_Listed, &responses)
}

func (p *CrewPipe) GetCrewPipe(ctx context.Context, userID uuid.UUID, crewID uuid.UUID, personaID uuid.UUID) *shared.PipeRes[CrewResponse] {
	if _, message := p.visiblePersona(ctx, userID, personaID); message != "" {
		return shared.PipeError[CrewResponse](message)
	}
	crew, err := p.crewRepo.FindCrewByID(ctx, crewID)
	if err != nil {
		if errors.Is(err, crew_repo.ErrCrewNotFound) {
			return shared.PipeError[CrewResponse](messages.Crew_Not_Found)
		}
		return pipeInternalError[CrewResponse](err, "crew.get")
	}
	if _, err := p.crewRepo.FindCrewMember(ctx, crewID, personaID); err != nil {
		return shared.PipeError[CrewResponse](messages.Forbidden)
	}
	members, err := p.crewRepo.FindCrewMembers(ctx, crewID)
	if err != nil {
		return pipeInternalError[CrewResponse](err, "crew.members")
	}
	locations, err := p.crewRepo.FindActiveCrewLocations(ctx, crewID)
	if err != nil {
		return pipeInternalError[CrewResponse](err, "crew.locations")
	}
	response, err := p.crewResponse(ctx, crew, members, locations)
	if err != nil {
		return pipeInternalError[CrewResponse](err, "crew.response")
	}
	return shared.PipeSuccess(messages.Crew_Fetched, response)
}

func (p *CrewPipe) LeaveCrewPipe(ctx context.Context, userID uuid.UUID, crewID uuid.UUID, personaID uuid.UUID) *shared.PipeRes[any] {
	member, message := p.requireMember(ctx, userID, crewID, personaID)
	if message != "" {
		return shared.PipeError[any](message)
	}
	if member.Role == models.OwnerCrewMemberRole {
		if _, err := p.crewRepo.EndCrew(ctx, crewID); err != nil {
			return pipeInternalError[any](err, "crew.owner_end")
		}
		return shared.PipeSuccess[any](messages.Crew_Ended, nil)
	}
	if err := p.crewRepo.LeaveCrew(ctx, crewID, personaID); err != nil {
		return pipeInternalError[any](err, "crew.leave")
	}
	return shared.PipeSuccess[any](messages.Crew_Left, nil)
}

func (p *CrewPipe) EndCrewPipe(ctx context.Context, userID uuid.UUID, crewID uuid.UUID, personaID uuid.UUID) *shared.PipeRes[CrewResponse] {
	member, message := p.requireMember(ctx, userID, crewID, personaID)
	if message != "" {
		return shared.PipeError[CrewResponse](message)
	}
	if member.Role != models.OwnerCrewMemberRole {
		return shared.PipeError[CrewResponse](messages.Forbidden)
	}
	crew, err := p.crewRepo.EndCrew(ctx, crewID)
	if err != nil {
		return pipeInternalError[CrewResponse](err, "crew.end")
	}
	response, err := p.crewResponse(ctx, crew, []*models.EventCrewMember{member}, nil)
	if err != nil {
		return pipeInternalError[CrewResponse](err, "crew.end_response")
	}
	return shared.PipeSuccess(messages.Crew_Ended, response)
}

func (p *CrewPipe) requireMember(ctx context.Context, userID uuid.UUID, crewID uuid.UUID, personaID uuid.UUID) (*models.EventCrewMember, shared.PipeMessage) {
	persona, message := p.visiblePersona(ctx, userID, personaID)
	if message != "" {
		return nil, message
	}
	member, err := p.crewRepo.FindCrewMember(ctx, crewID, persona.ID)
	if err != nil {
		return nil, messages.Forbidden
	}
	return member, ""
}
