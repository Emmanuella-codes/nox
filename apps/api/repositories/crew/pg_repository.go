package crew

import (
	"context"
	"errors"
	"time"

	"github.com/emmanuella-codes/nox/crew/dtos"
	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const maxCrewMembers = 25

type pgRepository struct {
	db *pgxpool.Pool
}

func newPgRepository(db *pgxpool.Pool) *pgRepository {
	return &pgRepository{db: db}
}

func (r *pgRepository) CreateCrew(ctx context.Context, ownerUserID uuid.UUID, eventID uuid.UUID, joinCode string, expiresAt time.Time, dto dtos.CreateCrewDTO) (*models.EventCrew, error) {
	visibility := dto.Visibility
	if visibility == "" {
		visibility = models.InviteCodeCrewVisibility
	}
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	conversation, err := insertCrewConversation(ctx, tx, dto.Name, dto.OwnerPersonaID)
	if err != nil {
		return nil, err
	}
	if _, err := insertConversationMember(ctx, tx, conversation.ID, ownerUserID, dto.OwnerPersonaID, models.ConversationMemberRoleAdmin); err != nil {
		return nil, err
	}

	row := tx.QueryRow(ctx, `
		INSERT INTO event_crews (event_id, conversation_id, owner_user_id, owner_persona_id, name, join_code, visibility, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, event_id, conversation_id, owner_user_id, owner_persona_id, name, join_code, visibility, status, expires_at, created_at, updated_at
	`, eventID, conversation.ID, ownerUserID, dto.OwnerPersonaID, dto.Name, joinCode, visibility, expiresAt)
	crew, err := scanCrew(row)
	if err != nil {
		if isUniqueViolation(err, "event_crews_join_code_key") {
			return nil, ErrCrewCodeTaken
		}
		return nil, err
	}
	if _, err := insertMember(ctx, tx, crew.ID, ownerUserID, dto.OwnerPersonaID, models.OwnerCrewMemberRole); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return crew, nil
}

func (r *pgRepository) FindCrewByID(ctx context.Context, crewID uuid.UUID) (*models.EventCrew, error) {
	row := r.db.QueryRow(ctx, crewSelectSQL()+" WHERE id = $1", crewID)
	return scanCrew(row)
}

func (r *pgRepository) FindCrewByJoinCode(ctx context.Context, joinCode string) (*models.EventCrew, error) {
	row := r.db.QueryRow(ctx, crewSelectSQL()+" WHERE join_code = $1", joinCode)
	return scanCrew(row)
}

func (r *pgRepository) FindEventCrewsForPersona(ctx context.Context, eventID uuid.UUID, personaID uuid.UUID, limit int, offset int) ([]*models.EventCrew, error) {
	rows, err := r.db.Query(ctx, crewSelectSQL()+`
		JOIN event_crew_members ecm ON ecm.crew_id = event_crews.id
		WHERE event_crews.event_id = $1 AND ecm.persona_id = $2 AND ecm.left_at IS NULL
		ORDER BY event_crews.updated_at DESC
		LIMIT $3 OFFSET $4
	`, eventID, personaID, normalizeLimit(limit), normalizeOffset(offset))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanCrews(rows)
}

func (r *pgRepository) FindCrewMembers(ctx context.Context, crewID uuid.UUID) ([]*models.EventCrewMember, error) {
	rows, err := r.db.Query(ctx, `
		SELECT crew_id, user_id, persona_id, role, location_sharing_enabled, joined_at, left_at
		FROM event_crew_members
		WHERE crew_id = $1 AND left_at IS NULL
		ORDER BY joined_at ASC
	`, crewID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMembers(rows)
}

func (r *pgRepository) FindCrewMember(ctx context.Context, crewID uuid.UUID, personaID uuid.UUID) (*models.EventCrewMember, error) {
	row := r.db.QueryRow(ctx, `
		SELECT crew_id, user_id, persona_id, role, location_sharing_enabled, joined_at, left_at
		FROM event_crew_members
		WHERE crew_id = $1 AND persona_id = $2 AND left_at IS NULL
	`, crewID, personaID)
	return scanMember(row)
}

func (r *pgRepository) JoinCrew(ctx context.Context, crew *models.EventCrew, persona *models.Persona) (*models.EventCrewMember, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	if _, err := tx.Exec(ctx, `SELECT pg_advisory_xact_lock(hashtext($1::text))`, crew.ID.String()); err != nil {
		return nil, err
	}
	var memberCount int
	if err := tx.QueryRow(ctx, `
		SELECT COUNT(*)::INT
		FROM event_crew_members
		WHERE crew_id = $1 AND left_at IS NULL
	`, crew.ID).Scan(&memberCount); err != nil {
		return nil, err
	}
	if memberCount >= maxCrewMembers {
		if _, err := r.findCrewMember(ctx, tx, crew.ID, persona.ID); err != nil {
			return nil, ErrCrewFull
		}
	}
	member, err := insertMember(ctx, tx, crew.ID, persona.UserID, persona.ID, models.MemberCrewMemberRole)
	if err != nil {
		return nil, err
	}
	if _, err := insertConversationMember(ctx, tx, crew.ConversationID, persona.UserID, persona.ID, models.ConversationMemberRoleMember); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return member, nil
}

func (r *pgRepository) LeaveCrew(ctx context.Context, crewID uuid.UUID, personaID uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	crew, err := r.findCrewByID(ctx, tx, crewID)
	if err != nil {
		return err
	}
	tag, err := tx.Exec(ctx, `
		UPDATE event_crew_members SET left_at = now(), location_sharing_enabled = false
		WHERE crew_id = $1 AND persona_id = $2 AND left_at IS NULL
	`, crewID, personaID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrCrewMemberNotFound
	}
	if _, err := tx.Exec(ctx, `
		UPDATE conversation_members SET left_at = now()
		WHERE conversation_id = $1 AND persona_id = $2 AND left_at IS NULL
	`, crew.ConversationID, personaID); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `DELETE FROM event_crew_locations WHERE crew_id = $1 AND persona_id = $2`, crewID, personaID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *pgRepository) EndCrew(ctx context.Context, crewID uuid.UUID) (*models.EventCrew, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	row := tx.QueryRow(ctx, `
		UPDATE event_crews SET status = 'ended', updated_at = now()
		WHERE id = $1
		RETURNING id, event_id, conversation_id, owner_user_id, owner_persona_id, name, join_code, visibility, status, expires_at, created_at, updated_at
	`, crewID)
	crew, err := scanCrew(row)
	if err != nil {
		return nil, err
	}
	if _, err := tx.Exec(ctx, `DELETE FROM event_crew_locations WHERE crew_id = $1`, crewID); err != nil {
		return nil, err
	}
	if _, err := tx.Exec(ctx, `
		UPDATE event_crew_members
		SET left_at = COALESCE(left_at, now()), location_sharing_enabled = false
		WHERE crew_id = $1
	`, crewID); err != nil {
		return nil, err
	}
	if _, err := tx.Exec(ctx, `
		UPDATE conversation_members
		SET left_at = COALESCE(left_at, now())
		WHERE conversation_id = $1
	`, crew.ConversationID); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return crew, nil
}

func (r *pgRepository) UpdateLocationSharing(ctx context.Context, crewID uuid.UUID, personaID uuid.UUID, enabled bool) (*models.EventCrewMember, error) {
	row := r.db.QueryRow(ctx, `
		UPDATE event_crew_members SET location_sharing_enabled = $3
		WHERE crew_id = $1 AND persona_id = $2 AND left_at IS NULL
		RETURNING crew_id, user_id, persona_id, role, location_sharing_enabled, joined_at, left_at
	`, crewID, personaID, enabled)
	member, err := scanMember(row)
	if err != nil {
		return nil, err
	}
	if !enabled {
		_ = r.DeleteCrewLocation(ctx, crewID, personaID)
	}
	return member, nil
}

func (r *pgRepository) UpsertCrewLocation(ctx context.Context, crewID uuid.UUID, userID uuid.UUID, dto dtos.UpdateLocationDTO, expiresAt time.Time) (*models.EventCrewLocation, error) {
	row := r.db.QueryRow(ctx, `
		INSERT INTO event_crew_locations (crew_id, user_id, persona_id, latitude, longitude, accuracy_meters, battery_level, recorded_at, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, now(), $8)
		ON CONFLICT (crew_id, persona_id)
		DO UPDATE SET latitude = EXCLUDED.latitude, longitude = EXCLUDED.longitude, accuracy_meters = EXCLUDED.accuracy_meters,
		              battery_level = EXCLUDED.battery_level, recorded_at = now(), expires_at = EXCLUDED.expires_at
		RETURNING crew_id, user_id, persona_id, latitude, longitude, accuracy_meters, battery_level, recorded_at, expires_at
	`, crewID, userID, dto.PersonaID, dto.Latitude, dto.Longitude, dto.AccuracyMeters, dto.BatteryLevel, expiresAt)
	return scanLocation(row)
}

func (r *pgRepository) FindActiveCrewLocations(ctx context.Context, crewID uuid.UUID) ([]*models.EventCrewLocation, error) {
	rows, err := r.db.Query(ctx, `
		SELECT crew_id, user_id, persona_id, latitude, longitude, accuracy_meters, battery_level, recorded_at, expires_at
		FROM event_crew_locations
		WHERE crew_id = $1 AND expires_at > now()
		ORDER BY recorded_at DESC
	`, crewID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanLocations(rows)
}

func (r *pgRepository) DeleteCrewLocation(ctx context.Context, crewID uuid.UUID, personaID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM event_crew_locations WHERE crew_id = $1 AND persona_id = $2`, crewID, personaID)
	return err
}

func (r *pgRepository) DeleteExpiredCrewLocations(ctx context.Context, now time.Time) (int64, error) {
	tag, err := r.db.Exec(ctx, `DELETE FROM event_crew_locations WHERE expires_at <= $1`, now)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}

func crewSelectSQL() string {
	return `SELECT id, event_id, conversation_id, owner_user_id, owner_persona_id, name, join_code, visibility, status, expires_at, created_at, updated_at FROM event_crews`
}

func (r *pgRepository) findCrewByID(ctx context.Context, db execQuerier, crewID uuid.UUID) (*models.EventCrew, error) {
	row := db.QueryRow(ctx, crewSelectSQL()+" WHERE id = $1", crewID)
	return scanCrew(row)
}

func (r *pgRepository) findCrewMember(ctx context.Context, db execQuerier, crewID uuid.UUID, personaID uuid.UUID) (*models.EventCrewMember, error) {
	row := db.QueryRow(ctx, `
		SELECT crew_id, user_id, persona_id, role, location_sharing_enabled, joined_at, left_at
		FROM event_crew_members
		WHERE crew_id = $1 AND persona_id = $2 AND left_at IS NULL
	`, crewID, personaID)
	return scanMember(row)
}

func mapCrewError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrCrewNotFound
	}
	return err
}

func isUniqueViolation(err error, constraint string) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == constraint
}

func normalizeLimit(limit int) int {
	if limit <= 0 {
		return 20
	}
	if limit > 100 {
		return 100
	}
	return limit
}

func normalizeOffset(offset int) int {
	if offset < 0 {
		return 0
	}
	return offset
}
