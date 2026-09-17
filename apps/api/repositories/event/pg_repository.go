package event

import (
	"context"
	"errors"
	"time"

	"github.com/emmanuella-codes/nox/event/dtos"
	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgRepository struct {
	db *pgxpool.Pool
}

func newPgRepository(db *pgxpool.Pool) *pgRepository {
	return &pgRepository{db: db}
}

func (r *pgRepository) CreateEvent(ctx context.Context, dto dtos.CreateEventDTO) (*models.Event, error) {
	endDate := dto.EventEndDate
	if endDate.IsZero() {
		endDate = dto.EventDate.Add(6 * time.Hour)
	}
	row := r.db.QueryRow(ctx, `
		INSERT INTO events (title, venue, location, event_date, event_end_date, description, cover_url, ticket_url, price_ngn, genre_tags, organizer_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, title, venue, location, event_date, description, COALESCE(cover_url, ''),
		          COALESCE(ticket_url, ''), price_ngn, genre_tags, organizer_id, created_at, event_end_date, status, ticket_sale_status, event_mode_status, updated_at
	`, dto.Title, dto.Venue, dto.Location, dto.EventDate, endDate, dto.Description, emptyToNil(dto.CoverURL), emptyToNil(dto.TicketURL), dto.Price, dto.GenreTags, dto.OrganizerID)

	event, err := scanEvent(row)
	if err != nil {
		return nil, mapEventError(err)
	}
	return event, nil
}

func (r *pgRepository) FindEventByID(ctx context.Context, eventID uuid.UUID) (*models.Event, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, title, venue, location, event_date, description, COALESCE(cover_url, ''),
		       COALESCE(ticket_url, ''), price_ngn, genre_tags, organizer_id, created_at, event_end_date, status, ticket_sale_status, event_mode_status, updated_at
		FROM events
		WHERE id = $1 AND moderation_status = 'active'
	`, eventID)

	event, err := scanEvent(row)
	if err != nil {
		return nil, mapEventError(err)
	}
	return event, nil
}

func (r *pgRepository) FindEvents(ctx context.Context, limit int) ([]*models.Event, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, title, venue, location, event_date, description, COALESCE(cover_url, ''),
		       COALESCE(ticket_url, ''), price_ngn, genre_tags, organizer_id, created_at, event_end_date, status, ticket_sale_status, event_mode_status, updated_at
		FROM events
		WHERE moderation_status = 'active'
		AND status IN ('published', 'pre_event', 'live', 'closing')
		ORDER BY event_date ASC, id ASC
		LIMIT $1
	`, normalizeLimit(limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*models.Event
	for rows.Next() {
		event, err := scanEvent(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, rows.Err()
}

func (r *pgRepository) UpdateEvent(ctx context.Context, userID, eventID uuid.UUID, dto dtos.UpdateEventDTO) (*models.Event, error) {
	row := r.db.QueryRow(ctx, `
		UPDATE events SET title=$3, venue=$4, location=$5, event_date=$6, event_end_date=$7, description=$8,
			cover_url=$9, ticket_url=$10, price_ngn=$11, genre_tags=$12, updated_at=now()
		WHERE id=$1 AND organizer_id=$2 AND status IN ('draft','published') AND (status='draft' OR event_date > now() + INTERVAL '24 hours')
		RETURNING id, title, venue, location, event_date, description, COALESCE(cover_url, ''), COALESCE(ticket_url, ''), price_ngn, genre_tags, organizer_id, created_at, event_end_date, status, ticket_sale_status, event_mode_status, updated_at
	`, eventID, userID, dto.Title, dto.Venue, dto.Location, dto.EventDate, dto.EventEndDate, dto.Description, emptyToNil(dto.CoverURL), emptyToNil(dto.TicketURL), dto.Price, dto.GenreTags)
	event, err := scanEvent(row)
	if err != nil {
		return nil, mapEventError(err)
	}
	return event, nil
}

func (r *pgRepository) PublishEvent(ctx context.Context, userID, eventID uuid.UUID) (*models.Event, error) {
	return r.transition(ctx, userID, eventID, "published", "organizer_published")
}

func (r *pgRepository) CloseEvent(ctx context.Context, userID, eventID uuid.UUID) (*models.Event, error) {
	return r.transition(ctx, userID, eventID, "follow_up", "organizer_closed")
}

func (r *pgRepository) ArchiveEvent(ctx context.Context, userID, eventID uuid.UUID) (*models.Event, error) {
	return r.transition(ctx, userID, eventID, "archived", "organizer_archived")
}

func (r *pgRepository) transition(ctx context.Context, userID, eventID uuid.UUID, next, reason string) (*models.Event, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	var previous string
	if err := tx.QueryRow(ctx, `SELECT status FROM events WHERE id=$1 AND organizer_id=$2 FOR UPDATE`, eventID, userID).Scan(&previous); err != nil {
		return nil, mapEventError(err)
	}
	if !validTransition(previous, next) {
		return nil, ErrInvalidEventTransition
	}
	row := tx.QueryRow(ctx, `UPDATE events SET status=$3, published_at=CASE WHEN $3='published' THEN COALESCE(published_at,now()) ELSE published_at END, closed_at=CASE WHEN $3='follow_up' THEN now() ELSE closed_at END, archived_at=CASE WHEN $3='archived' THEN now() ELSE archived_at END, updated_at=now() WHERE id=$1 AND organizer_id=$2 RETURNING id, title, venue, location, event_date, description, COALESCE(cover_url,''), COALESCE(ticket_url,''), price_ngn, genre_tags, organizer_id, created_at, event_end_date, status, ticket_sale_status, event_mode_status, updated_at`, eventID, userID, next)
	event, err := scanEvent(row)
	if err != nil {
		return nil, err
	}
	if _, err = tx.Exec(ctx, `INSERT INTO event_lifecycle_history (event_id, previous_status, next_status, actor_user_id, reason) VALUES ($1,$2,$3,$4,$5)`, eventID, previous, next, userID, reason); err != nil {
		return nil, err
	}
	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}
	return event, nil
}

func (r *pgRepository) AdvanceLifecycle(ctx context.Context, now time.Time, limit int) (int, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)
	rows, err := tx.Query(ctx, `SELECT id, status FROM events WHERE ((status='published' AND event_date <= $1 + INTERVAL '24 hours') OR (status='pre_event' AND event_date <= $1) OR (status='live' AND event_end_date <= $1) OR (status='closing' AND event_end_date + INTERVAL '1 hour' <= $1) OR (status='ended' AND updated_at <= $1 - INTERVAL '24 hours')) ORDER BY event_date ASC FOR UPDATE SKIP LOCKED LIMIT $2`, now, limit)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	count := 0
	for rows.Next() {
		var id uuid.UUID
		var previous string
		if err := rows.Scan(&id, &previous); err != nil {
			return 0, err
		}
		next := nextDueStatus(previous)
		if _, err = tx.Exec(ctx, `UPDATE events SET status=$2, event_mode_status=CASE WHEN $2='pre_event' THEN 'pre_event' WHEN $2 IN ('live','closing') THEN $2 WHEN $2 IN ('ended','follow_up','closed','archived') THEN 'ended' ELSE event_mode_status END, ended_at=CASE WHEN $2='ended' THEN COALESCE(ended_at,now()) ELSE ended_at END, updated_at=now() WHERE id=$1`, id, next); err != nil {
			return 0, err
		}
		if _, err = tx.Exec(ctx, `INSERT INTO event_lifecycle_history (event_id, previous_status, next_status, reason) VALUES ($1,$2,$3,'lifecycle_worker')`, id, previous, next); err != nil {
			return 0, err
		}
		count++
	}
	if err := rows.Err(); err != nil {
		return 0, err
	}
	if err = tx.Commit(ctx); err != nil {
		return 0, err
	}
	return count, nil
}

func validTransition(previous, next string) bool {
	return (previous == "draft" && next == "published") || ((previous == "published" || previous == "pre_event" || previous == "live" || previous == "closing" || previous == "ended") && next == "follow_up") || (previous == "follow_up" && next == "archived")
}
func nextDueStatus(status string) string {
	switch status {
	case "published":
		return "pre_event"
	case "pre_event":
		return "live"
	case "live":
		return "closing"
	case "closing":
		return "ended"
	case "ended":
		return "follow_up"
	}
	return status
}

type eventScanner interface {
	Scan(dest ...any) error
}

func scanEvent(scanner eventScanner) (*models.Event, error) {
	var event models.Event
	err := scanner.Scan(
		&event.ID,
		&event.Title,
		&event.Venue,
		&event.Location,
		&event.EventDate,
		&event.Description,
		&event.CoverURL,
		&event.TicketURL,
		&event.Price,
		&event.GenreTags,
		&event.OrganizerID,
		&event.CreatedAt,
		&event.EventEndDate,
		&event.Status,
		&event.TicketSaleStatus,
		&event.EventModeStatus,
		&event.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func mapEventError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrEventNotFound
	}
	return err
}

func normalizeLimit(limit int) int {
	if limit <= 0 {
		return 30
	}
	if limit > 100 {
		return 100
	}
	return limit
}

func emptyToNil(value string) any {
	if value == "" {
		return nil
	}
	return value
}
