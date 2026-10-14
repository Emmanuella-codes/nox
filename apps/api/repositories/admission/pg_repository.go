package admission

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/emmanuella-codes/nox/models"
	staff_repo "github.com/emmanuella-codes/nox/repositories/staff"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("admission not found")
var ErrForbidden = errors.New("admission access denied")
var ErrAlreadyCheckedIn = errors.New("already checked in")

type InvalidError struct{ Reason string }

func (e *InvalidError) Error() string { return e.Reason }

type pgRepository struct {
	db    *pgxpool.Pool
	staff staff_repo.StaffRepository
}

func newPgRepository(db *pgxpool.Pool, staff staff_repo.StaffRepository) *pgRepository {
	return &pgRepository{db: db, staff: staff}
}

func (r *pgRepository) CheckIn(ctx context.Context, userID, eventID uuid.UUID, credential, method string) (*models.CheckInResult, error) {
	if r.staff != nil {
		if err := r.staff.Authorize(ctx, userID, eventID, "gate", time.Now().UTC()); err != nil {
			return nil, err
		}
	}
	return r.check(ctx, userID, eventID, credential, method, "")
}
func (r *pgRepository) ManualCheckIn(ctx context.Context, userID, eventID, admissionID uuid.UUID) (*models.CheckInResult, error) {
	if r.staff != nil {
		if err := r.staff.Authorize(ctx, userID, eventID, "gate", time.Now().UTC()); err != nil {
			return nil, err
		}
	}
	return r.check(ctx, userID, eventID, "", "manual", admissionID.String())
}
func (r *pgRepository) OverrideCheckIn(ctx context.Context, userID, admissionID uuid.UUID, reason string) (*models.CheckInResult, error) {
	var eventID uuid.UUID
	if err := r.db.QueryRow(ctx, `SELECT event_id FROM event_admissions WHERE id=$1`, admissionID).Scan(&eventID); err != nil {
		return nil, ErrNotFound
	}
	if r.staff != nil {
		if err := r.staff.Authorize(ctx, userID, eventID, "management", time.Now().UTC()); err != nil {
			return nil, err
		}
	}
	return r.check(ctx, userID, uuid.Nil, "", "override", admissionID.String()+"|"+reason)
}
func (r *pgRepository) check(ctx context.Context, userID, eventID uuid.UUID, credential, method, selector string) (*models.CheckInResult, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	var admissionID, attendeeID, actualEvent uuid.UUID
	var status, credentialStatus, attendeeStatus string
	var eventStatus string
	var user *uuid.UUID
	var reason string
	if method == "override" {
		parts := strings.SplitN(selector, "|", 2)
		id, e := uuid.Parse(parts[0])
		if e != nil {
			return nil, ErrNotFound
		}
		admissionID = id
		reason = parts[1]
	} else if selector != "" {
		admissionID, _ = uuid.Parse(selector)
	}
	query := `SELECT a.id,a.attendee_id,a.event_id,a.status,c.status,e.status,at.status,at.user_id FROM event_admissions a JOIN event_attendees at ON at.id=a.attendee_id JOIN events e ON e.id=a.event_id JOIN admission_credentials c ON c.admission_id=a.id WHERE a.id=$1 FOR UPDATE`
	args := []any{admissionID}
	if method != "override" && selector == "" {
		query = `SELECT a.id,a.attendee_id,a.event_id,a.status,c.status,e.status,at.status,at.user_id FROM event_admissions a JOIN event_attendees at ON at.id=a.attendee_id JOIN events e ON e.id=a.event_id JOIN admission_credentials c ON c.admission_id=a.id WHERE a.event_id=$1 AND c.credential_hash=$2 FOR UPDATE`
		args = []any{eventID, hashToken(credential)}
	}
	if err := tx.QueryRow(ctx, query, args...).Scan(&admissionID, &attendeeID, &actualEvent, &status, &credentialStatus, &eventStatus, &attendeeStatus, &user); err != nil {
		return nil, ErrNotFound
	}
	if method != "override" {
		reason = ""
	}
	var existing uuid.UUID
	if err := tx.QueryRow(ctx, `SELECT id FROM event_check_ins WHERE admission_id=$1`, admissionID).Scan(&existing); err == nil {
		_, _ = tx.Exec(ctx, `INSERT INTO event_admission_audits(event_id,admission_id,actor_user_id,action) VALUES($1,$2,$3,'duplicate_scan')`, actualEvent, admissionID, userID)
		_ = tx.Commit(ctx)
		return nil, ErrAlreadyCheckedIn
	}
	if err != pgx.ErrNoRows {
		return nil, err
	}
	invalid := ""
	if method != "override" && eventID != actualEvent {
		invalid = "WRONG_EVENT"
	}
	if eventStatus == "closed" || eventStatus == "archived" {
		invalid = "EVENT_CLOSED"
	}
	if status != "active" {
		invalid = strings.ToUpper(status)
	}
	if credentialStatus != "active" {
		invalid = strings.ToUpper(credentialStatus)
	}
	if attendeeStatus != "eligible" {
		invalid = "RESTRICTED"
	}
	if user != nil {
		var blocked bool
		if err := tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM event_attendee_restrictions WHERE event_id=$1 AND user_id=$2 AND status='active')`, actualEvent, *user).Scan(&blocked); err != nil {
			return nil, err
		}
		if blocked {
			invalid = "RESTRICTED"
		}
	}
	if method != "override" && invalid != "" {
		_, _ = tx.Exec(ctx, `INSERT INTO event_admission_audits(event_id,admission_id,attendee_id,actor_user_id,action,reason) VALUES($1,$2,$3,$4,'rejected_scan',$5)`, actualEvent, admissionID, attendeeID, userID, invalid)
		_ = tx.Commit(ctx)
		return nil, &InvalidError{Reason: invalid}
	}
	if method == "override" && strings.TrimSpace(reason) == "" {
		return nil, &InvalidError{Reason: "OVERRIDE_REASON_REQUIRED"}
	}
	var checkInID uuid.UUID
	if err := tx.QueryRow(ctx, `INSERT INTO event_check_ins(event_id,admission_id,staff_user_id,method,reason,invalid_reason) VALUES($1,$2,$3,$4,$5,$6) RETURNING id`, actualEvent, admissionID, userID, method, reason, invalid).Scan(&checkInID); err != nil {
		return nil, err
	}
	if _, err := tx.Exec(ctx, `UPDATE admission_credentials SET status='used',used_at=now() WHERE admission_id=$1`, admissionID); err != nil {
		return nil, err
	}
	if _, err := tx.Exec(ctx, `INSERT INTO event_admission_audits(event_id,admission_id,attendee_id,actor_user_id,action,reason) VALUES($1,$2,$3,$4,$5,$6)`, actualEvent, admissionID, attendeeID, userID, "check_in_"+method, reason); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return &models.CheckInResult{ID: checkInID, AdmissionID: admissionID, Method: method, Status: "checked_in", CheckedInAt: time.Now().UTC()}, nil
}
func (r *pgRepository) Search(ctx context.Context, userID, eventID uuid.UUID, term string) ([]*models.AdmissionSearchResult, error) {
	if r.staff != nil {
		if err := r.staff.Authorize(ctx, userID, eventID, "gate", time.Now().UTC()); err != nil {
			return nil, err
		}
	}
	rows, err := r.db.Query(ctx, `SELECT a.id,at.display_name,a.source_type,a.status,COALESCE(at.email,''),COALESCE(at.phone,'') FROM event_admissions a JOIN event_attendees at ON at.id=a.attendee_id WHERE a.event_id=$1 AND (at.display_name ILIKE $2 OR at.email ILIKE $2 OR at.phone ILIKE $2) ORDER BY at.display_name LIMIT 50`, eventID, "%"+term+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []*models.AdmissionSearchResult
	for rows.Next() {
		item := &models.AdmissionSearchResult{}
		if err := rows.Scan(&item.AdmissionID, &item.DisplayName, &item.SourceType, &item.Status, &item.Email, &item.Phone); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}
func (r *pgRepository) Reconcile(ctx context.Context, userID, eventID uuid.UUID) (map[string]int, error) {
	if r.staff != nil {
		if err := r.staff.Authorize(ctx, userID, eventID, "management", time.Now().UTC()); err != nil {
			return nil, err
		}
	}
	rows, err := r.db.Query(ctx, `SELECT 'issued_admissions',COUNT(*) FROM event_admissions WHERE event_id=$1 UNION ALL SELECT 'successful_check_ins',COUNT(*) FROM event_check_ins WHERE event_id=$1 UNION ALL SELECT 'manual_check_ins',COUNT(*) FROM event_check_ins WHERE event_id=$1 AND method='manual' UNION ALL SELECT 'overrides',COUNT(*) FROM event_check_ins WHERE event_id=$1 AND method='override' UNION ALL SELECT 'duplicate_scans',COUNT(*) FROM event_admission_audits WHERE event_id=$1 AND action='duplicate_scan' UNION ALL SELECT 'rejected_scans',COUNT(*) FROM event_admission_audits WHERE event_id=$1 AND action='rejected_scan' UNION ALL SELECT 'revoked_admissions',COUNT(*) FROM event_admissions WHERE event_id=$1 AND status='revoked' UNION ALL SELECT 'refunded_admissions',COUNT(*) FROM event_admissions WHERE event_id=$1 AND status='refunded' UNION ALL SELECT 'no_shows',COUNT(*) FROM event_guest_entries WHERE event_id=$1 AND status='no_show'`, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := map[string]int{}
	for rows.Next() {
		var key string
		var count int
		if err := rows.Scan(&key, &count); err != nil {
			return nil, err
		}
		result[key] = count
	}
	return result, rows.Err()
}
func (r *pgRepository) EnsureGuest(ctx context.Context, eventID, sourceID uuid.UUID, credential, name, email, phone string) error {
	hash := hashToken(credential)
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var attendee uuid.UUID
	if err = tx.QueryRow(ctx, `INSERT INTO event_attendees(event_id,user_id,display_name,email,phone) SELECT $1,matched_user_id,name,email,phone FROM event_guest_entries WHERE id=$2 ON CONFLICT DO NOTHING RETURNING id`, eventID, sourceID).Scan(&attendee); err != nil {
		if err != pgx.ErrNoRows {
			return err
		}
		if err = tx.QueryRow(ctx, `SELECT id FROM event_attendees WHERE event_id=$1 AND display_name=$2 AND COALESCE(email,'')=COALESCE($3,'') ORDER BY created_at DESC LIMIT 1`, eventID, name, email).Scan(&attendee); err != nil {
			return err
		}
	}
	if _, err = tx.Exec(ctx, `INSERT INTO event_admissions(attendee_id,event_id,source_type,source_id) VALUES($1,$2,'guest',$3) ON CONFLICT DO NOTHING`, attendee, eventID, sourceID); err != nil {
		return err
	}
	var admission uuid.UUID
	if err = tx.QueryRow(ctx, `SELECT id FROM event_admissions WHERE source_type='guest' AND source_id=$1`, sourceID).Scan(&admission); err != nil {
		return err
	}
	if _, err = tx.Exec(ctx, `INSERT INTO admission_credentials(admission_id,credential_hash) VALUES($1,$2) ON CONFLICT(admission_id) DO UPDATE SET credential_hash=EXCLUDED.credential_hash,status='active'`, admission, hash); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *pgRepository) Create(ctx context.Context, actor, eventID uuid.UUID, sourceType string, sourceID, userID uuid.UUID, name, email, phone, credential string) (*models.Admission, error) {
	if sourceType != "complimentary" && sourceType != "staff" {
		return nil, ErrForbidden
	}
	if r.staff != nil {
		if err := r.staff.Authorize(ctx, actor, eventID, "management", time.Now().UTC()); err != nil {
			return nil, err
		}
	}
	if sourceID == uuid.Nil {
		sourceID = uuid.New()
	}
	if name == "" {
		name = "Attendee"
	}
	if credential == "" {
		credential = "nox_" + uuid.NewString()
	}
	var attendeeUser any
	if userID != uuid.Nil {
		attendeeUser = userID
	}
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	var attendeeID, admissionID uuid.UUID
	if err = tx.QueryRow(ctx, `INSERT INTO event_attendees(event_id,user_id,display_name,email,phone) VALUES($1,$2,$3,NULLIF($4,''),NULLIF($5,'')) RETURNING id`, eventID, attendeeUser, name, email, phone).Scan(&attendeeID); err != nil {
		return nil, err
	}
	if err = tx.QueryRow(ctx, `INSERT INTO event_admissions(attendee_id,event_id,source_type,source_id) VALUES($1,$2,$3,$4) RETURNING id`, attendeeID, eventID, sourceType, sourceID).Scan(&admissionID); err != nil {
		return nil, err
	}
	if _, err = tx.Exec(ctx, `INSERT INTO admission_credentials(admission_id,credential_hash) VALUES($1,$2)`, admissionID, hashToken(credential)); err != nil {
		return nil, err
	}
	if _, err = tx.Exec(ctx, `INSERT INTO event_admission_audits(event_id,admission_id,attendee_id,actor_user_id,action) VALUES($1,$2,$3,$4,'admission_created')`, eventID, admissionID, attendeeID, actor); err != nil {
		return nil, err
	}
	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}
	return &models.Admission{ID: admissionID, AttendeeID: attendeeID, EventID: eventID, SourceType: sourceType, SourceID: sourceID, Status: "active", DisplayName: name, Email: email, Phone: phone, Credential: credential}, nil
}
func hashToken(v string) string { s := sha256.Sum256([]byte(v)); return hex.EncodeToString(s[:]) }
