package staff

import (
	"context"
	"errors"
	"time"

	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/staff/dtos"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("staff resource not found")
var ErrForbidden = errors.New("staff permission denied")

type pgRepository struct {
	db *pgxpool.Pool
}

func newPgRepository(db *pgxpool.Pool) *pgRepository { return &pgRepository{db: db} }

func (r *pgRepository) Invite(ctx context.Context, actor, eventID, userID uuid.UUID) (*models.EventStaffMembership, error) {
	if err := r.manage(ctx, actor, eventID); err != nil {
		return nil, err
	}
	var owner uuid.UUID
	if err := r.db.QueryRow(ctx, `SELECT organizer_id FROM events WHERE id=$1`, eventID).Scan(&owner); err != nil || owner == userID {
		return nil, ErrForbidden
	}
	row := r.db.QueryRow(ctx, `INSERT INTO event_staff_memberships(event_id,user_id,invited_by) VALUES($1,$2,$3) ON CONFLICT(event_id,user_id) DO UPDATE SET status='invited',invited_by=$3,updated_at=now() RETURNING id,event_id,user_id,status,accepted_at`, eventID, userID, actor)
	item, err := scanMembership(row)
	if err != nil {
		return nil, ErrNotFound
	}
	_ = r.audit(ctx, eventID, item.ID, actor, "staff.invited")
	return item, nil
}
func (r *pgRepository) AssignRoles(ctx context.Context, actor, membershipID uuid.UUID, roles []string) error {
	var eventID uuid.UUID
	if err := r.db.QueryRow(ctx, `SELECT event_id FROM event_staff_memberships WHERE id=$1`, membershipID).Scan(&eventID); err != nil {
		return ErrNotFound
	}
	if err := r.manage(ctx, actor, eventID); err != nil {
		return err
	}
	for _, role := range roles {
		if !validRole(role) {
			return ErrForbidden
		}
	}
	var ownerID, targetID uuid.UUID
	if err := r.db.QueryRow(ctx, `SELECT e.organizer_id,m.user_id FROM events e JOIN event_staff_memberships m ON m.event_id=e.id WHERE m.id=$1 AND e.status NOT IN ('closed','archived')`, membershipID).Scan(&ownerID, &targetID); err != nil {
		return ErrForbidden
	}
	if containsRole(roles, "owner") && ownerID != actor {
		return ErrForbidden
	}
	if targetID == ownerID && !containsRole(roles, "owner") {
		return ErrForbidden
	}
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err = tx.Exec(ctx, `DELETE FROM event_staff_roles WHERE membership_id=$1`, membershipID); err != nil {
		return err
	}
	for _, role := range roles {
		if _, err = tx.Exec(ctx, `INSERT INTO event_staff_roles(membership_id,role) VALUES($1,$2)`, membershipID, role); err != nil {
			return err
		}
	}
	if err = tx.Commit(ctx); err != nil {
		return err
	}
	return r.audit(ctx, eventID, membershipID, actor, "staff.roles_assigned")
}
func (r *pgRepository) Accept(ctx context.Context, userID, membershipID uuid.UUID) (*models.EventStaffMembership, error) {
	row := r.db.QueryRow(ctx, `UPDATE event_staff_memberships SET status='active',accepted_at=now(),updated_at=now() WHERE id=$1 AND user_id=$2 AND status='invited' RETURNING id,event_id,user_id,status,accepted_at`, membershipID, userID)
	item, err := scanMembership(row)
	if err != nil {
		return nil, ErrNotFound
	}
	_ = r.audit(ctx, item.EventID, item.ID, userID, "staff.accepted")
	return item, nil
}
func (r *pgRepository) Revoke(ctx context.Context, actor, membershipID uuid.UUID) error {
	var eventID uuid.UUID
	if err := r.db.QueryRow(ctx, `SELECT event_id FROM event_staff_memberships WHERE id=$1`, membershipID).Scan(&eventID); err != nil {
		return ErrNotFound
	}
	if err := r.manage(ctx, actor, eventID); err != nil {
		return err
	}
	if _, err := r.db.Exec(ctx, `UPDATE event_staff_memberships SET status='revoked',updated_at=now() WHERE id=$1 AND status<>'revoked'`, membershipID); err != nil {
		return err
	}
	return r.audit(ctx, eventID, membershipID, actor, "staff.revoked")
}
func (r *pgRepository) CreateShift(ctx context.Context, actor, membershipID uuid.UUID, dto dtos.ShiftDTO) (*models.EventStaffShift, error) {
	start, end, err := parseTimes(dto.StartsAt, dto.EndsAt)
	if err != nil {
		return nil, err
	}
	var eventID uuid.UUID
	if err = r.db.QueryRow(ctx, `SELECT event_id FROM event_staff_memberships WHERE id=$1 AND status='active'`, membershipID).Scan(&eventID); err != nil {
		return nil, ErrNotFound
	}
	if err = r.manage(ctx, actor, eventID); err != nil {
		return nil, err
	}
	row := r.db.QueryRow(ctx, `INSERT INTO event_staff_shifts(membership_id,starts_at,ends_at) VALUES($1,$2,$3) RETURNING id,membership_id,starts_at,ends_at,status`, membershipID, start, end)
	item := &models.EventStaffShift{}
	if err = row.Scan(&item.ID, &item.MembershipID, &item.StartsAt, &item.EndsAt, &item.Status); err != nil {
		return nil, err
	}
	_ = r.audit(ctx, eventID, membershipID, actor, "staff.shift_created")
	return item, nil
}
func (r *pgRepository) List(ctx context.Context, actor, eventID uuid.UUID) ([]*models.EventStaffMembership, error) {
	if err := r.manage(ctx, actor, eventID); err != nil {
		return nil, err
	}
	rows, err := r.db.Query(ctx, `SELECT m.id,m.event_id,m.user_id,m.status,m.accepted_at,COALESCE(array_agg(r.role) FILTER(WHERE r.role IS NOT NULL),'{}') FROM event_staff_memberships m LEFT JOIN event_staff_roles r ON r.membership_id=m.id WHERE m.event_id=$1 GROUP BY m.id ORDER BY m.created_at`, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []*models.EventStaffMembership
	for rows.Next() {
		item := &models.EventStaffMembership{}
		if err := rows.Scan(&item.ID, &item.EventID, &item.UserID, &item.Status, &item.AcceptedAt, &item.Roles); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}
func (r *pgRepository) Authorize(ctx context.Context, userID, eventID uuid.UUID, permission string, at time.Time) error {
	var allowed bool
	err := r.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM event_staff_memberships m JOIN event_staff_roles r ON r.membership_id=m.id JOIN events e ON e.id=m.event_id WHERE m.event_id=$2 AND m.user_id=$1 AND m.status='active' AND e.status NOT IN ('closed','archived') AND (r.role='owner' OR r.role='manager' OR ($3='gate' AND r.role='gate') OR ($3='safety' AND r.role IN ('safety_lead','safety_officer')) OR ($3='finance' AND r.role='finance')) AND (r.role IN ('owner','manager','read_only') OR EXISTS(SELECT 1 FROM event_staff_shifts s WHERE s.membership_id=m.id AND s.status IN ('scheduled','active') AND s.starts_at <= $4 AND s.ends_at >= $4))`, userID, eventID, permission, at).Scan(&allowed)
	if err != nil {
		return err
	}
	if !allowed {
		return ErrForbidden
	}
	return nil
}
func (r *pgRepository) manage(ctx context.Context, userID, eventID uuid.UUID) error {
	var ok bool
	err := r.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM event_staff_memberships m JOIN event_staff_roles r ON r.membership_id=m.id JOIN events e ON e.id=m.event_id WHERE m.event_id=$2 AND m.user_id=$1 AND m.status='active' AND e.status NOT IN ('closed','archived') AND r.role IN ('owner','manager'))`, userID, eventID).Scan(&ok)
	if err != nil {
		return err
	}
	if !ok {
		return ErrForbidden
	}
	return nil
}
func (r *pgRepository) audit(ctx context.Context, eventID, membershipID, actor uuid.UUID, action string) error {
	_, err := r.db.Exec(ctx, `INSERT INTO event_staff_audit_logs(event_id,membership_id,actor_user_id,action) VALUES($1,$2,$3,$4)`, eventID, membershipID, actor, action)
	return err
}
func validRole(role string) bool {
	switch role {
	case "owner", "manager", "safety_lead", "safety_officer", "gate", "finance", "read_only":
		return true
	}
	return false
}
func containsRole(roles []string, wanted string) bool {
	for _, role := range roles {
		if role == wanted {
			return true
		}
	}
	return false
}
func parseTimes(start, end string) (time.Time, time.Time, error) {
	a, e := time.Parse(time.RFC3339, start)
	if e != nil {
		return time.Time{}, time.Time{}, e
	}
	b, e := time.Parse(time.RFC3339, end)
	return a, b, e
}
func scanMembership(s interface{ Scan(...any) error }) (*models.EventStaffMembership, error) {
	item := &models.EventStaffMembership{}
	err := s.Scan(&item.ID, &item.EventID, &item.UserID, &item.Status, &item.AcceptedAt)
	return item, err
}
