package pipes

import (
	"context"
	"errors"
	"strings"

	"github.com/emmanuella-codes/nox/admin/dtos"
	"github.com/emmanuella-codes/nox/admin/messages"
	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/repositories/admin"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (p *AdminPipe) ListCrewsPipe(ctx context.Context, actorID uuid.UUID, status string, limit, offset int) *shared.PipeRes[[]models.AdminCrewSummary] {
	if message := p.requireActiveAdmin(ctx, actorID); message != "" {
		return shared.PipeError[[]models.AdminCrewSummary](shared.CreatePipeMessage(message))
	}
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	rows, err := p.db.Query(ctx, `SELECT ec.id, ec.event_id, ec.name, ec.owner_persona_id, ec.status,
		COUNT(DISTINCT ecm.persona_id) FILTER (WHERE ecm.left_at IS NULL),
		COUNT(DISTINCT ecm.persona_id) FILTER (WHERE ecm.left_at IS NULL AND ecm.location_sharing_enabled),
		ec.expires_at, ec.created_at, ec.updated_at
		FROM event_crews ec LEFT JOIN event_crew_members ecm ON ecm.crew_id = ec.id
		WHERE ($1 = '' OR ec.status = $1) GROUP BY ec.id ORDER BY ec.created_at DESC LIMIT $2 OFFSET $3`, status, limit, offset)
	if err != nil {
		logInternalError(err, "crew.list")
		return shared.PipeError[[]models.AdminCrewSummary](messages.Internal_Error)
	}
	defer rows.Close()
	items := make([]models.AdminCrewSummary, 0)
	for rows.Next() {
		var item models.AdminCrewSummary
		if err := rows.Scan(&item.ID, &item.EventID, &item.Name, &item.OwnerPersonaID, &item.Status, &item.MemberCount, &item.SharingCount, &item.ExpiresAt, &item.CreatedAt, &item.UpdatedAt); err != nil {
			logInternalError(err, "crew.list_scan")
			return shared.PipeError[[]models.AdminCrewSummary](messages.Internal_Error)
		}
		items = append(items, item)
	}
	p.audit(ctx, actorID, "admin.crew.list", map[string]any{"status": status})
	return shared.PipeSuccess(messages.Admin_Crews_Loaded, &items)
}

func (p *AdminPipe) GetCrewPipe(ctx context.Context, actorID, crewID uuid.UUID) *shared.PipeRes[models.AdminCrewDetail] {
	if message := p.requireActiveAdmin(ctx, actorID); message != "" {
		return shared.PipeError[models.AdminCrewDetail](shared.CreatePipeMessage(message))
	}
	crew, err := p.findAdminCrew(ctx, crewID)
	if err != nil {
		return crewError[models.AdminCrewDetail](err, "crew.get")
	}
	members, err := p.listCrewMembers(ctx, crewID)
	if err != nil {
		logInternalError(err, "crew.members")
		return shared.PipeError[models.AdminCrewDetail](messages.Internal_Error)
	}
	result := models.AdminCrewDetail{Crew: *crew, Members: members}
	p.audit(ctx, actorID, "admin.crew.get", map[string]any{"crew_id": crewID.String()})
	return shared.PipeSuccess(messages.Admin_Crew_Loaded, &result)
}

func (p *AdminPipe) EndCrewSafetyPipe(ctx context.Context, actorID, crewID uuid.UUID, dto dtos.CrewSafetyActionDTO) *shared.PipeRes[models.AdminCrewSummary] {
	if message := p.requireActiveAdmin(ctx, actorID); message != "" {
		return shared.PipeError[models.AdminCrewSummary](shared.CreatePipeMessage(message))
	}
	dto.Reason = strings.TrimSpace(dto.Reason)
	if dto.Reason == "" {
		return shared.PipeError[models.AdminCrewSummary](messages.Invalid_Payload)
	}
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return shared.PipeError[models.AdminCrewSummary](messages.Internal_Error)
	}
	defer tx.Rollback(ctx)
	var status models.CrewStatus
	if err := tx.QueryRow(ctx, `UPDATE event_crews SET status = 'ended', updated_at = now() WHERE id = $1 RETURNING status`, crewID).Scan(&status); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return shared.PipeError[models.AdminCrewSummary](messages.Crew_Not_Found)
		}
		logInternalError(err, "crew.end")
		return shared.PipeError[models.AdminCrewSummary](messages.Internal_Error)
	}
	if _, err := tx.Exec(ctx, `UPDATE event_crew_members SET location_sharing_enabled = false, left_at = COALESCE(left_at, now()) WHERE crew_id = $1`, crewID); err != nil {
		return shared.PipeError[models.AdminCrewSummary](messages.Internal_Error)
	}
	if _, err := tx.Exec(ctx, `DELETE FROM event_crew_locations WHERE crew_id = $1`, crewID); err != nil {
		return shared.PipeError[models.AdminCrewSummary](messages.Internal_Error)
	}
	if err := tx.Commit(ctx); err != nil {
		return shared.PipeError[models.AdminCrewSummary](messages.Internal_Error)
	}
	crew, err := p.findAdminCrew(ctx, crewID)
	if err != nil {
		return crewError[models.AdminCrewSummary](err, "crew.end_get")
	}
	p.audit(ctx, actorID, "admin.crew.end", map[string]any{"crew_id": crewID.String(), "reason": dto.Reason})
	return shared.PipeSuccess(messages.Admin_Crew_Ended, crew)
}

func (p *AdminPipe) DisableCrewSharingPipe(ctx context.Context, actorID, crewID uuid.UUID, dto dtos.CrewSafetyActionDTO) *shared.PipeRes[any] {
	if message := p.requireActiveAdmin(ctx, actorID); message != "" {
		return shared.PipeError[any](shared.CreatePipeMessage(message))
	}
	dto.Reason = strings.TrimSpace(dto.Reason)
	if dto.Reason == "" {
		return shared.PipeError[any](messages.Invalid_Payload)
	}
	result, err := p.db.Exec(ctx, `UPDATE event_crew_members SET location_sharing_enabled = false WHERE crew_id = $1 AND left_at IS NULL`, crewID)
	if err != nil {
		logInternalError(err, "crew.disable_sharing")
		return shared.PipeError[any](messages.Internal_Error)
	}
	if result.RowsAffected() == 0 {
		if _, err := p.findAdminCrew(ctx, crewID); err != nil {
			return crewError[any](err, "crew.disable_sharing_find")
		}
	}
	if _, err := p.db.Exec(ctx, `DELETE FROM event_crew_locations WHERE crew_id = $1`, crewID); err != nil {
		return shared.PipeError[any](messages.Internal_Error)
	}
	p.audit(ctx, actorID, "admin.crew.disable_location_sharing", map[string]any{"crew_id": crewID.String(), "reason": dto.Reason})
	return shared.PipeSuccess[any](messages.Admin_Crew_Sharing_Disabled, nil)
}

func (p *AdminPipe) ListCrewLocationsPipe(ctx context.Context, actorID, crewID uuid.UUID, dto dtos.CrewLocationAccessDTO) *shared.PipeRes[[]models.EventCrewLocation] {
	if message := p.requireActiveAdmin(ctx, actorID); message != "" {
		return shared.PipeError[[]models.EventCrewLocation](shared.CreatePipeMessage(message))
	}
	dto.Reason = strings.TrimSpace(dto.Reason)
	if dto.Reason == "" {
		return shared.PipeError[[]models.EventCrewLocation](messages.Invalid_Payload)
	}
	if _, err := p.findAdminCrew(ctx, crewID); err != nil {
		return crewError[[]models.EventCrewLocation](err, "crew.locations_find")
	}
	rows, err := p.db.Query(ctx, `SELECT crew_id, user_id, persona_id, latitude, longitude, accuracy_meters, battery_level, recorded_at, expires_at FROM event_crew_locations WHERE crew_id = $1 AND expires_at > now() ORDER BY recorded_at DESC`, crewID)
	if err != nil {
		logInternalError(err, "crew.locations")
		return shared.PipeError[[]models.EventCrewLocation](messages.Internal_Error)
	}
	defer rows.Close()
	locations := make([]models.EventCrewLocation, 0)
	for rows.Next() {
		var item models.EventCrewLocation
		if err := rows.Scan(&item.CrewID, &item.UserID, &item.PersonaID, &item.Latitude, &item.Longitude, &item.AccuracyMeters, &item.BatteryLevel, &item.RecordedAt, &item.ExpiresAt); err != nil {
			return shared.PipeError[[]models.EventCrewLocation](messages.Internal_Error)
		}
		locations = append(locations, item)
	}
	meta := map[string]any{"crew_id": crewID.String(), "reason": dto.Reason}
	if dto.ReportID != "" {
		meta["report_id"] = dto.ReportID
	}
	p.audit(ctx, actorID, "admin.crew.locations.read", meta)
	return shared.PipeSuccess(messages.Admin_Crew_Locations_Loaded, &locations)
}

func (p *AdminPipe) findAdminCrew(ctx context.Context, crewID uuid.UUID) (*models.AdminCrewSummary, error) {
	var item models.AdminCrewSummary
	err := p.db.QueryRow(ctx, `SELECT ec.id, ec.event_id, ec.name, ec.owner_persona_id, ec.status, COUNT(DISTINCT ecm.persona_id) FILTER (WHERE ecm.left_at IS NULL), COUNT(DISTINCT ecm.persona_id) FILTER (WHERE ecm.left_at IS NULL AND ecm.location_sharing_enabled), ec.expires_at, ec.created_at, ec.updated_at FROM event_crews ec LEFT JOIN event_crew_members ecm ON ecm.crew_id = ec.id WHERE ec.id = $1 GROUP BY ec.id`, crewID).Scan(&item.ID, &item.EventID, &item.Name, &item.OwnerPersonaID, &item.Status, &item.MemberCount, &item.SharingCount, &item.ExpiresAt, &item.CreatedAt, &item.UpdatedAt)
	return &item, err
}

func (p *AdminPipe) listCrewMembers(ctx context.Context, crewID uuid.UUID) ([]models.AdminCrewMember, error) {
	rows, err := p.db.Query(ctx, `SELECT crew_id, user_id, persona_id, role, location_sharing_enabled, joined_at, left_at FROM event_crew_members WHERE crew_id = $1 ORDER BY joined_at`, crewID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]models.AdminCrewMember, 0)
	for rows.Next() {
		var item models.AdminCrewMember
		if err := rows.Scan(&item.CrewID, &item.UserID, &item.PersonaID, &item.Role, &item.LocationSharingEnabled, &item.JoinedAt, &item.LeftAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func crewError[T any](err error, operation string) *shared.PipeRes[T] {
	if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, admin.ErrCrewNotFound) {
		return shared.PipeError[T](messages.Crew_Not_Found)
	}
	logInternalError(err, operation)
	return shared.PipeError[T](messages.Internal_Error)
}
