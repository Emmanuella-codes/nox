package pipes

import (
	"context"
	"crypto/rand"
	"errors"
	"math"
	"strings"
	"time"

	"github.com/emmanuella-codes/nox/crew/messages"
	"github.com/emmanuella-codes/nox/models"
	crew_repo "github.com/emmanuella-codes/nox/repositories/crew"
	event_repo "github.com/emmanuella-codes/nox/repositories/event"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

const (
	crewEventWindow  = 12 * time.Hour
	locationPingTTL  = 5 * time.Minute
	joinCodeAlphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
)

type CrewPipe struct {
	crewRepo    crew_repo.CrewRepository
	eventRepo   event_repo.EventRepository
	personaRepo persona_repo.PersonaRepository
}

func NewCrewPipe(crewRepo crew_repo.CrewRepository, eventRepo event_repo.EventRepository, personaRepo persona_repo.PersonaRepository) *CrewPipe {
	return &CrewPipe{crewRepo: crewRepo, eventRepo: eventRepo, personaRepo: personaRepo}
}

func pipeInternalError[T any](err error, operation string) *shared.PipeRes[T] {
	return shared.PipeInternalError[T](err, "crew", operation, messages.Internal_Error)
}

type CrewResponse struct {
	ID             string             `json:"id"`
	EventID        string             `json:"event_id"`
	ConversationID string             `json:"conversation_id"`
	OwnerPersonaID string             `json:"owner_persona_id"`
	Name           string             `json:"name"`
	JoinCode       string             `json:"join_code"`
	Visibility     string             `json:"visibility"`
	Status         string             `json:"status"`
	ExpiresAt      time.Time          `json:"expires_at"`
	Members        []MemberResponse   `json:"members"`
	Locations      []LocationResponse `json:"locations"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
}

type MemberResponse struct {
	PersonaID              string          `json:"persona_id"`
	Role                   string          `json:"role"`
	LocationSharingEnabled bool            `json:"location_sharing_enabled"`
	Persona                *PersonaSummary `json:"persona,omitempty"`
	JoinedAt               time.Time       `json:"joined_at"`
}

type PersonaSummary struct {
	ID          string `json:"id"`
	Handle      string `json:"handle"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
}

type LocationResponse struct {
	PersonaID      string          `json:"persona_id"`
	Latitude       float64         `json:"latitude"`
	Longitude      float64         `json:"longitude"`
	AccuracyMeters float64         `json:"accuracy_meters"`
	BatteryLevel   *int            `json:"battery_level,omitempty"`
	RecordedAt     time.Time       `json:"recorded_at"`
	ExpiresAt      time.Time       `json:"expires_at"`
	Persona        *PersonaSummary `json:"persona,omitempty"`
}

func (p *CrewPipe) visiblePersona(ctx context.Context, userID uuid.UUID, personaID uuid.UUID) (*models.Persona, shared.PipeMessage) {
	persona, err := p.personaRepo.FindPersonaByID(ctx, personaID)
	if err != nil {
		if errors.Is(err, persona_repo.ErrPersonaNotFound) {
			return nil, messages.Persona_Not_Found
		}
		return nil, messages.Internal_Error
	}
	if persona.UserID != userID || persona.PersonaType != models.VisiblePersonaType {
		return nil, messages.Forbidden
	}
	return persona, ""
}

func (p *CrewPipe) crewResponse(ctx context.Context, crew *models.EventCrew, members []*models.EventCrewMember, locations []*models.EventCrewLocation) (*CrewResponse, error) {
	personas, err := p.personasForCrew(ctx, members, locations)
	if err != nil {
		return nil, err
	}
	response := &CrewResponse{
		ID: crew.ID.String(), EventID: crew.EventID.String(), ConversationID: crew.ConversationID.String(), OwnerPersonaID: crew.OwnerPersonaID.String(),
		Name: crew.Name, JoinCode: crew.JoinCode, Visibility: string(crew.Visibility), Status: string(crew.Status),
		ExpiresAt: crew.ExpiresAt, Members: memberResponses(members, personas), Locations: locationResponses(locations, personas),
		CreatedAt: crew.CreatedAt, UpdatedAt: crew.UpdatedAt,
	}
	return response, nil
}

func (p *CrewPipe) personasForCrew(ctx context.Context, members []*models.EventCrewMember, locations []*models.EventCrewLocation) (map[uuid.UUID]*models.Persona, error) {
	personas := map[uuid.UUID]*models.Persona{}
	for _, member := range members {
		if _, ok := personas[member.PersonaID]; ok {
			continue
		}
		persona, err := p.personaRepo.FindPersonaByID(ctx, member.PersonaID)
		if err != nil {
			return nil, err
		}
		personas[member.PersonaID] = persona
	}
	for _, location := range locations {
		if _, ok := personas[location.PersonaID]; ok {
			continue
		}
		persona, err := p.personaRepo.FindPersonaByID(ctx, location.PersonaID)
		if err != nil {
			return nil, err
		}
		personas[location.PersonaID] = persona
	}
	return personas, nil
}

func memberResponses(members []*models.EventCrewMember, personas map[uuid.UUID]*models.Persona) []MemberResponse {
	responses := make([]MemberResponse, 0, len(members))
	for _, member := range members {
		responses = append(responses, MemberResponse{
			PersonaID: member.PersonaID.String(), Role: string(member.Role), LocationSharingEnabled: member.LocationSharingEnabled,
			Persona: personaSummary(personas[member.PersonaID]), JoinedAt: member.JoinedAt,
		})
	}
	return responses
}

func locationResponses(locations []*models.EventCrewLocation, personas map[uuid.UUID]*models.Persona) []LocationResponse {
	responses := make([]LocationResponse, 0, len(locations))
	for _, location := range locations {
		responses = append(responses, LocationResponse{
			PersonaID: location.PersonaID.String(), Latitude: location.Latitude, Longitude: location.Longitude,
			AccuracyMeters: location.AccuracyMeters, BatteryLevel: location.BatteryLevel, RecordedAt: location.RecordedAt,
			ExpiresAt: location.ExpiresAt, Persona: personaSummary(personas[location.PersonaID]),
		})
	}
	return responses
}

func personaSummary(persona *models.Persona) *PersonaSummary {
	if persona == nil {
		return nil
	}
	return &PersonaSummary{ID: persona.ID.String(), Handle: persona.Handle, DisplayName: persona.DisplayName, AvatarURL: persona.AvatarURL}
}

func activeCrew(crew *models.EventCrew) bool {
	return crew.Status == models.ActiveCrewStatus && time.Now().Before(crew.ExpiresAt)
}

func crewExpiresAt(event *models.Event) time.Time {
	return event.EventDate.Add(crewEventWindow)
}

func validLocation(lat, lng, accuracy float64) bool {
	return !math.IsNaN(lat) && !math.IsNaN(lng) && lat >= -90 && lat <= 90 && lng >= -180 && lng <= 180 && accuracy >= 0 && accuracy <= 5000
}

func joinCode() string {
	buf := make([]byte, 6)
	_, err := rand.Read(buf)
	if err != nil {
		return strings.ToUpper(uuid.NewString()[:6])
	}
	for i := range buf {
		buf[i] = joinCodeAlphabet[int(buf[i])%len(joinCodeAlphabet)]
	}
	return string(buf)
}
