package guest

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/emmanuella-codes/nox/guest/dtos"
	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
	"time"
)

var ErrNotFound = errors.New("guest entry not found")
var ErrConflict = errors.New("guest list conflict")

type pgRepository struct{ db *pgxpool.Pool }

func NewGuestRepository(db *pgxpool.Pool) GuestRepository { return &pgRepository{db: db} }

func (r *pgRepository) Configure(ctx context.Context, userID, eventID uuid.UUID, dto dtos.ConfigureDTO) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var allocated int
	if err := tx.QueryRow(ctx, `SELECT COALESCE(SUM(capacity),0) FROM event_guest_quotas WHERE event_id=$1`, eventID).Scan(&allocated); err != nil {
		return err
	}
	if dto.Capacity < allocated {
		return ErrConflict
	}
	if _, err := tx.Exec(ctx, `INSERT INTO event_guest_settings(event_id,capacity) SELECT id,$3 FROM events WHERE id=$1 AND organizer_id=$2 ON CONFLICT(event_id) DO UPDATE SET capacity=EXCLUDED.capacity,updated_at=now()`, eventID, userID, dto.Capacity); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *pgRepository) CreateQuota(ctx context.Context, userID, eventID uuid.UUID, dto dtos.QuotaDTO) (*models.GuestQuota, error) {
	owner, err := uuid.Parse(dto.OwnerUserID)
	if err != nil {
		return nil, ErrConflict
	}
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	var capacity, allocated int
	if err := tx.QueryRow(ctx, `SELECT s.capacity FROM event_guest_settings s WHERE s.event_id=$1 AND EXISTS (SELECT 1 FROM events WHERE id=$1 AND organizer_id=$2) FOR UPDATE`, eventID, userID).Scan(&capacity); err != nil {
		return nil, ErrConflict
	}
	if err := tx.QueryRow(ctx, `SELECT COALESCE(SUM(capacity),0) FROM event_guest_quotas WHERE event_id=$1`, eventID).Scan(&allocated); err != nil {
		return nil, err
	}
	if allocated+dto.Capacity > capacity {
		return nil, ErrConflict
	}
	row := tx.QueryRow(ctx, `INSERT INTO event_guest_quotas(event_id,owner_user_id,label,capacity) VALUES($1,$2,$3,$4) RETURNING id,event_id,owner_user_id,label,capacity,reserved_count,status`, eventID, owner, strings.TrimSpace(dto.Label), dto.Capacity)
	item, err := scanQuota(row)
	if err != nil {
		return nil, ErrConflict
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return item, nil
}

func (r *pgRepository) AddGuest(ctx context.Context, userID, eventID, quotaID uuid.UUID, dto dtos.GuestDTO, invitationHashes []string) ([]*models.GuestEntry, error) {
	if strings.TrimSpace(dto.Email) == "" && strings.TrimSpace(dto.Phone) == "" {
		return nil, ErrConflict
	}
	if dto.Quantity <= 0 {
		dto.Quantity = 1
	}
	if dto.Quantity > 10 {
		return nil, ErrConflict
	}
	if len(invitationHashes) != dto.Quantity {
		return nil, ErrConflict
	}
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	var quota uuid.UUID
	row := tx.QueryRow(ctx, `UPDATE event_guest_quotas q SET reserved_count=q.reserved_count+$4 WHERE q.id=$1 AND q.event_id=$2 AND q.status='active' AND q.reserved_count+$4 <= q.capacity AND EXISTS (SELECT 1 FROM events e WHERE e.id=q.event_id AND (e.organizer_id=$3 OR EXISTS (SELECT 1 FROM event_guest_quotas h WHERE h.id=q.id AND h.owner_user_id=$3 AND h.status='active'))) RETURNING q.id`, quotaID, eventID, userID, dto.Quantity)
	if err := row.Scan(&quota); err != nil {
		return nil, ErrConflict
	}
	entries := make([]*models.GuestEntry, 0, dto.Quantity)
	for i := 0; i < dto.Quantity; i++ {
		row = tx.QueryRow(ctx, `INSERT INTO event_guest_entries(event_id,quota_id,added_by,name,email,phone,matched_user_id,invitation_hash,invitation_expires_at) SELECT $1,$2,$3,$4,NULLIF($5,''),NULLIF($6,''),u.id,$7,e.event_date+interval '1 day' FROM events e LEFT JOIN users u ON lower(u.email)=lower(NULLIF($5,'')) WHERE e.id=$1 RETURNING id,event_id,quota_id,name,COALESCE(email,''),COALESCE(phone,''),matched_user_id,status,created_at`, eventID, quota, userID, strings.TrimSpace(dto.Name), strings.TrimSpace(dto.Email), strings.TrimSpace(dto.Phone), invitationHashes[i])
		entry := &models.GuestEntry{}
		if err := row.Scan(&entry.ID, &entry.EventID, &entry.QuotaID, &entry.Name, &entry.Email, &entry.Phone, &entry.MatchedUserID, &entry.Status, &entry.CreatedAt); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return entries, nil
}

func (r *pgRepository) Claim(ctx context.Context, token string, userID *uuid.UUID) (*models.GuestEntry, error) {
	var entry models.GuestEntry
	credential := uuid.NewString()
	hash := hashToken(credential)
	var claimedUser any
	if userID != nil {
		claimedUser = *userID
	}
	row := r.db.QueryRow(ctx, `UPDATE event_guest_entries SET status='claimed',claimed_at=COALESCE(claimed_at,now()),updated_at=now() WHERE invitation_hash=$1 AND status='invited' AND invitation_expires_at>now() RETURNING id,event_id,quota_id,name,COALESCE(email,''),COALESCE(phone,''),status,created_at`, hashToken(token))
	if err := row.Scan(&entry.ID, &entry.EventID, &entry.QuotaID, &entry.Name, &entry.Email, &entry.Phone, &entry.Status, &entry.CreatedAt); err != nil {
		return nil, ErrNotFound
	}
	if _, err := r.db.Exec(ctx, `UPDATE event_guest_entries SET invitation_hash=NULL,credential_hash=$2,matched_user_id=COALESCE(matched_user_id,$3),updated_at=now() WHERE id=$1`, entry.ID, hash, claimedUser); err != nil {
		return nil, err
	}
	entry.Credential = credential
	return &entry, nil
}

func (r *pgRepository) Decline(ctx context.Context, token string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var quotaID uuid.UUID
	if err := tx.QueryRow(ctx, `UPDATE event_guest_entries SET status='declined',declined_at=now(),invitation_hash=NULL,updated_at=now() WHERE invitation_hash=$1 AND status='invited' AND invitation_expires_at>now() RETURNING quota_id`, hashToken(token)).Scan(&quotaID); err != nil {
		return ErrNotFound
	}
	if _, err := tx.Exec(ctx, `UPDATE event_guest_quotas SET reserved_count=GREATEST(reserved_count-1,0) WHERE id=$1`, quotaID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *pgRepository) Revoke(ctx context.Context, userID, entryID uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var quotaID uuid.UUID
	if err := tx.QueryRow(ctx, `UPDATE event_guest_entries g SET status='revoked',revoked_at=now(),invitation_hash=NULL,credential_hash=NULL,updated_at=now() FROM events e WHERE g.id=$1 AND (e.organizer_id=$2 OR EXISTS (SELECT 1 FROM event_guest_quotas q WHERE q.id=g.quota_id AND q.owner_user_id=$2)) AND g.status IN ('invited','claimed') RETURNING g.quota_id`, entryID, userID).Scan(&quotaID); err != nil {
		return ErrNotFound
	}
	if _, err := tx.Exec(ctx, `UPDATE event_guest_quotas SET reserved_count=GREATEST(reserved_count-1,0) WHERE id=$1`, quotaID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *pgRepository) Attach(ctx context.Context, userID, entryID uuid.UUID) error {
	result, err := r.db.Exec(ctx, `UPDATE event_guest_entries g SET matched_user_id=$2,updated_at=now() FROM users u WHERE g.id=$1 AND u.id=$2 AND g.status='claimed' AND g.matched_user_id IS NULL AND lower(u.email)=lower(g.email)`, entryID, userID)
	if err != nil || result.RowsAffected() != 1 {
		return ErrConflict
	}
	return nil
}

func (r *pgRepository) ExpireInvitations(ctx context.Context, now time.Time) (int64, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)
	result, err := tx.Exec(ctx, `WITH expired AS (UPDATE event_guest_entries SET status='expired',invitation_hash=NULL,updated_at=now() WHERE status='invited' AND invitation_expires_at <= $1 RETURNING quota_id), counts AS (SELECT quota_id,COUNT(*) total FROM expired GROUP BY quota_id) UPDATE event_guest_quotas q SET reserved_count=GREATEST(q.reserved_count-c.total,0) FROM counts c WHERE q.id=c.quota_id`, now)
	if err != nil {
		return 0, err
	}
	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func (r *pgRepository) MarkNoShows(ctx context.Context, now time.Time) (int64, error) {
	result, err := r.db.Exec(ctx, `UPDATE event_guest_entries g SET status='no_show',no_show_at=now(),updated_at=now() FROM events e WHERE g.status IN ('invited','claimed') AND g.event_id=e.id AND e.event_date < $1`, now)
	return result.RowsAffected(), err
}

func (r *pgRepository) CheckIn(ctx context.Context, eventID, staffID uuid.UUID, credential string) (*models.GuestEntry, error) {
	var entry models.GuestEntry
	row := r.db.QueryRow(ctx, `UPDATE event_guest_entries g SET status='checked_in',checked_in_at=now(),checked_in_by=$2,updated_at=now() FROM events e WHERE g.event_id=$1 AND g.credential_hash=$3 AND g.status='claimed' AND e.id=g.event_id AND (e.organizer_id=$2 OR EXISTS (SELECT 1 FROM event_guest_quotas q WHERE q.id=g.quota_id AND q.owner_user_id=$2)) RETURNING g.id,g.event_id,g.quota_id,g.name,COALESCE(g.email,''),COALESCE(g.phone,''),g.status,g.created_at`, eventID, staffID, hashToken(credential))
	if err := row.Scan(&entry.ID, &entry.EventID, &entry.QuotaID, &entry.Name, &entry.Email, &entry.Phone, &entry.Status, &entry.CreatedAt); err != nil {
		return nil, ErrConflict
	}
	return &entry, nil
}

func (r *pgRepository) Report(ctx context.Context, userID, eventID uuid.UUID) (*models.GuestListReport, error) {
	var report models.GuestListReport
	report.EventID = eventID
	if err := r.db.QueryRow(ctx, `SELECT s.capacity FROM event_guest_settings s JOIN events e ON e.id=s.event_id WHERE s.event_id=$1 AND e.organizer_id=$2`, eventID, userID).Scan(&report.Capacity); err != nil {
		return nil, ErrNotFound
	}
	rows, err := r.db.Query(ctx, `SELECT id,event_id,owner_user_id,label,capacity,reserved_count,status FROM event_guest_quotas WHERE event_id=$1 ORDER BY label,id`, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		quota, err := scanQuota(rows)
		if err != nil {
			return nil, err
		}
		report.Reserved += quota.ReservedCount
		report.Quotas = append(report.Quotas, quota)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	rows, err = r.db.Query(ctx, `SELECT status,COUNT(*) FROM event_guest_entries WHERE event_id=$1 GROUP BY status`, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	report.Statuses = map[string]int{}
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		report.Statuses[status] = count
	}
	return &report, rows.Err()
}

func hashToken(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}
func scanQuota(s interface{ Scan(...any) error }) (*models.GuestQuota, error) {
	var q models.GuestQuota
	err := s.Scan(&q.ID, &q.EventID, &q.OwnerUserID, &q.Label, &q.Capacity, &q.ReservedCount, &q.Status)
	return &q, err
}
