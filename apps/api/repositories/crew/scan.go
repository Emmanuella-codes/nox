package crew

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type crewScanner interface {
	Scan(dest ...any) error
}

type crewRows interface {
	crewScanner
	Next() bool
	Err() error
}

type execQuerier interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func scanCrew(scanner crewScanner) (*models.EventCrew, error) {
	var crew models.EventCrew
	err := scanner.Scan(
		&crew.ID,
		&crew.EventID,
		&crew.ConversationID,
		&crew.OwnerUserID,
		&crew.OwnerPersonaID,
		&crew.Name,
		&crew.JoinCode,
		&crew.Visibility,
		&crew.Status,
		&crew.ExpiresAt,
		&crew.CreatedAt,
		&crew.UpdatedAt,
	)
	if err != nil {
		return nil, mapCrewError(err)
	}
	return &crew, nil
}

func scanCrews(rows crewRows) ([]*models.EventCrew, error) {
	crews := make([]*models.EventCrew, 0)
	for rows.Next() {
		crew, err := scanCrew(rows)
		if err != nil {
			return nil, err
		}
		crews = append(crews, crew)
	}
	return crews, rows.Err()
}

func insertCrewConversation(ctx context.Context, db execQuerier, title string, createdBy uuid.UUID) (*models.Conversation, error) {
	row := db.QueryRow(ctx, `
		INSERT INTO conversations (conversation_type, title, created_by)
		VALUES ($1, $2, $3)
		RETURNING id, conversation_type, title, created_by, last_message_id, created_at, updated_at
	`, models.GroupConversationType, title, createdBy)
	var conversation models.Conversation
	err := row.Scan(
		&conversation.ID,
		&conversation.ConversationType,
		&conversation.Title,
		&conversation.CreatedBy,
		&conversation.LastMessageID,
		&conversation.CreatedAt,
		&conversation.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &conversation, nil
}

func insertConversationMember(ctx context.Context, db execQuerier, conversationID uuid.UUID, userID uuid.UUID, personaID uuid.UUID, role models.ConversationMemberRole) (*models.ConversationMember, error) {
	row := db.QueryRow(ctx, `
		INSERT INTO conversation_members (conversation_id, user_id, persona_id, role)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (conversation_id, persona_id)
		DO UPDATE SET left_at = NULL
		RETURNING conversation_id, user_id, persona_id, role, last_read_message_id, joined_at, left_at
	`, conversationID, userID, personaID, role)
	var member models.ConversationMember
	err := row.Scan(
		&member.ConversationID,
		&member.UserID,
		&member.PersonaID,
		&member.Role,
		&member.LastReadMessageID,
		&member.JoinedAt,
		&member.LeftAt,
	)
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func insertMember(ctx context.Context, db execQuerier, crewID uuid.UUID, userID uuid.UUID, personaID uuid.UUID, role models.CrewMemberRole) (*models.EventCrewMember, error) {
	row := db.QueryRow(ctx, `
		INSERT INTO event_crew_members (crew_id, user_id, persona_id, role, left_at)
		VALUES ($1, $2, $3, $4, NULL)
		ON CONFLICT (crew_id, persona_id)
		DO UPDATE SET left_at = NULL
		RETURNING crew_id, user_id, persona_id, role, location_sharing_enabled, joined_at, left_at
	`, crewID, userID, personaID, role)
	return scanMember(row)
}

func scanMember(scanner crewScanner) (*models.EventCrewMember, error) {
	var member models.EventCrewMember
	err := scanner.Scan(
		&member.CrewID,
		&member.UserID,
		&member.PersonaID,
		&member.Role,
		&member.LocationSharingEnabled,
		&member.JoinedAt,
		&member.LeftAt,
	)
	if err != nil {
		return nil, mapCrewError(err)
	}
	return &member, nil
}

func scanMembers(rows crewRows) ([]*models.EventCrewMember, error) {
	members := make([]*models.EventCrewMember, 0)
	for rows.Next() {
		member, err := scanMember(rows)
		if err != nil {
			return nil, err
		}
		members = append(members, member)
	}
	return members, rows.Err()
}

func scanLocation(scanner crewScanner) (*models.EventCrewLocation, error) {
	var location models.EventCrewLocation
	err := scanner.Scan(
		&location.CrewID,
		&location.UserID,
		&location.PersonaID,
		&location.Latitude,
		&location.Longitude,
		&location.AccuracyMeters,
		&location.BatteryLevel,
		&location.RecordedAt,
		&location.ExpiresAt,
	)
	if err != nil {
		return nil, err
	}
	return &location, nil
}

func scanLocations(rows crewRows) ([]*models.EventCrewLocation, error) {
	locations := make([]*models.EventCrewLocation, 0)
	for rows.Next() {
		location, err := scanLocation(rows)
		if err != nil {
			return nil, err
		}
		locations = append(locations, location)
	}
	return locations, rows.Err()
}
