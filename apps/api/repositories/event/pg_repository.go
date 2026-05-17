package event

import (
	"context"
	"errors"

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
	row := r.db.QueryRow(ctx, `
		INSERT INTO events (title, venue, location, event_date, description, cover_url, ticket_url, price_ngn, genre_tags, organizer_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, title, venue, location, event_date, description, COALESCE(cover_url, ''),
		          COALESCE(ticket_url, ''), price_ngn, genre_tags, organizer_id, created_at
	`, dto.Title, dto.Venue, dto.Location, dto.EventDate, dto.Description, emptyToNil(dto.CoverURL), emptyToNil(dto.TicketURL), dto.Price, dto.GenreTags, dto.OrganizerID)

	event, err := scanEvent(row)
	if err != nil {
		return nil, mapEventError(err)
	}
	return event, nil
}

func (r *pgRepository) FindEventByID(ctx context.Context, eventID uuid.UUID) (*models.Event, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, title, venue, location, event_date, description, COALESCE(cover_url, ''),
		       COALESCE(ticket_url, ''), price_ngn, genre_tags, organizer_id, created_at
		FROM events
		WHERE id = $1
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
		       COALESCE(ticket_url, ''), price_ngn, genre_tags, organizer_id, created_at
		FROM events
		ORDER BY event_date ASC
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
