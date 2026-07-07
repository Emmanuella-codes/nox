package pipes

import (
	"context"
	"time"

	"github.com/emmanuella-codes/nox/crew/dtos"
	"github.com/emmanuella-codes/nox/crew/messages"
	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

func (p *CrewPipe) UpdateLocationSharingPipe(ctx context.Context, userID uuid.UUID, crewID uuid.UUID, dto dtos.UpdateSharingDTO) *shared.PipeRes[MemberResponse] {
	if _, message := p.requireMember(ctx, userID, crewID, dto.PersonaID); message != "" {
		return shared.PipeError[MemberResponse](message)
	}
	crew, err := p.crewRepo.FindCrewByID(ctx, crewID)
	if err != nil {
		return pipeInternalError[MemberResponse](err, "crew.sharing_crew")
	}
	if !activeCrew(crew) {
		return shared.PipeError[MemberResponse](messages.Forbidden)
	}
	member, err := p.crewRepo.UpdateLocationSharing(ctx, crewID, dto.PersonaID, dto.Enabled)
	if err != nil {
		return pipeInternalError[MemberResponse](err, "crew.sharing")
	}
	personas, err := p.personasForCrew(ctx, []*models.EventCrewMember{member}, nil)
	if err != nil {
		return pipeInternalError[MemberResponse](err, "crew.sharing_persona")
	}
	response := memberResponses([]*models.EventCrewMember{member}, personas)[0]
	return shared.PipeSuccess(messages.Sharing_Updated, &response)
}

func (p *CrewPipe) UpdateLocationPipe(ctx context.Context, userID uuid.UUID, crewID uuid.UUID, dto dtos.UpdateLocationDTO) *shared.PipeRes[LocationResponse] {
	member, message := p.requireMember(ctx, userID, crewID, dto.PersonaID)
	if message != "" {
		return shared.PipeError[LocationResponse](message)
	}
	crew, err := p.crewRepo.FindCrewByID(ctx, crewID)
	if err != nil {
		return pipeInternalError[LocationResponse](err, "crew.location_crew")
	}
	if !activeCrew(crew) || !member.LocationSharingEnabled || !validLocation(dto.Latitude, dto.Longitude, dto.AccuracyMeters) {
		return shared.PipeError[LocationResponse](messages.Forbidden)
	}
	location, err := p.crewRepo.UpsertCrewLocation(ctx, crewID, userID, dto, time.Now().Add(locationPingTTL))
	if err != nil {
		return pipeInternalError[LocationResponse](err, "crew.location")
	}
	personas, err := p.personasForCrew(ctx, nil, []*models.EventCrewLocation{location})
	if err != nil {
		return pipeInternalError[LocationResponse](err, "crew.location_persona")
	}
	response := locationResponses([]*models.EventCrewLocation{location}, personas)[0]
	return shared.PipeSuccess(messages.Location_Updated, &response)
}

func (p *CrewPipe) ListLocationsPipe(ctx context.Context, userID uuid.UUID, crewID uuid.UUID, personaID uuid.UUID) *shared.PipeRes[[]LocationResponse] {
	if _, message := p.requireMember(ctx, userID, crewID, personaID); message != "" {
		return shared.PipeError[[]LocationResponse](message)
	}
	locations, err := p.crewRepo.FindActiveCrewLocations(ctx, crewID)
	if err != nil {
		return pipeInternalError[[]LocationResponse](err, "crew.list_locations")
	}
	personas, err := p.personasForCrew(ctx, nil, locations)
	if err != nil {
		return pipeInternalError[[]LocationResponse](err, "crew.location_personas")
	}
	response := locationResponses(locations, personas)
	return shared.PipeSuccess(messages.Locations_Listed, &response)
}
