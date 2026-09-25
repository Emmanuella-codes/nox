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

func (r *pgRepository) AddGuest(ctx context.Context, userID, eventID, quotaID uuid.UUID, dto dtos.GuestDTO, invitationHash string) (*models.GuestEntry, error) {
	if strings.TrimSpace(dto.Email) == "" && strings.TrimSpace(dto.Phone) == "" {
		return nil, ErrConflict
	}
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	var entry models.GuestEntry
	row := tx.QueryRow(ctx, `UPDATE event_guest_quotas q SET reserved_count=q.reserved_count+1 WHERE q.id=$1 AND q.event_id=$2 AND q.status='active' AND q.reserved_count < q.capacity AND EXISTS (SELECT 1 FROM events e WHERE e.id=q.event_id AND (e.organizer_id=$3 OR EXISTS (SELECT 1 FROM event_guest_quotas h WHERE h.id=q.id AND h.owner_user_id=$3 AND h.status='active'))) RETURNING q.id`, quotaID, eventID, userID)
	if err := row.Scan(&entry.QuotaID); err != nil {
		return nil, ErrConflict
	}
	row = tx.QueryRow(ctx, `INSERT INTO event_guest_entries(event_id,quota_id,added_by,name,email,phone,matched_user_id,invitation_hash,invitation_expires_at) SELECT $1,$2,$3,$4,NULLIF($5,''),NULLIF($6,''),u.id,$7,e.event_date+interval '1 day' FROM events e LEFT JOIN users u ON lower(u.email)=lower(NULLIF($5,'')) WHERE e.id=$1 RETURNING id,event_id,quota_id,name,COALESCE(email,''),COALESCE(phone,''),status,created_at`, eventID, quotaID, userID, strings.TrimSpace(dto.Name), strings.TrimSpace(dto.Email), strings.TrimSpace(dto.Phone), invitationHash)
	if err := row.Scan(&entry.ID, &entry.EventID, &entry.QuotaID, &entry.Name, &entry.Email, &entry.Phone, &entry.Status, &entry.CreatedAt); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return &entry, nil
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
	if _, err := r.db.Exec(ctx, `UPDATE event_guest_entries SET invitation_hash=$2,matched_user_id=COALESCE(matched_user_id,$3),updated_at=now() WHERE id=$1`, entry.ID, hash, claimedUser); err != nil {
		return nil, err
	}
	entry.Credential = credential
	return &entry, nil
}

func (r *pgRepository) CheckIn(ctx context.Context, eventID, staffID uuid.UUID, credential string) (*models.GuestEntry, error) {
	var entry models.GuestEntry
	row := r.db.QueryRow(ctx, `UPDATE event_guest_entries g SET status='checked_in',checked_in_at=now(),checked_in_by=$2,updated_at=now() FROM events e WHERE g.event_id=$1 AND g.invitation_hash=$3 AND g.status='claimed' AND e.id=g.event_id AND (e.organizer_id=$2 OR EXISTS (SELECT 1 FROM event_guest_quotas q WHERE q.id=g.quota_id AND q.owner_user_id=$2)) RETURNING g.id,g.event_id,g.quota_id,g.name,COALESCE(g.email,''),COALESCE(g.phone,''),g.status,g.created_at`, eventID, staffID, hashToken(credential))
	if err := row.Scan(&entry.ID, &entry.EventID, &entry.QuotaID, &entry.Name, &entry.Email, &entry.Phone, &entry.Status, &entry.CreatedAt); err != nil {
		return nil, ErrConflict
	}
	return &entry, nil
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
