package search

import (
	"context"
	"time"

	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgRepository struct {
	db *pgxpool.Pool
}

func newPgRepository(db *pgxpool.Pool) *pgRepository {
	return &pgRepository{db: db}
}

func (r *pgRepository) Search(ctx context.Context, query string, limit int) (*Results, error) {
	normalizedLimit := normalizeLimit(limit)
	personas, err := r.searchPersonas(ctx, query, normalizedLimit)
	if err != nil {
		return nil, err
	}
	posts, err := r.searchPosts(ctx, query, normalizedLimit)
	if err != nil {
		return nil, err
	}
	events, err := r.searchEvents(ctx, query, normalizedLimit)
	if err != nil {
		return nil, err
	}
	return &Results{Personas: personas, Posts: posts, Events: events}, nil
}

func (r *pgRepository) searchPersonas(ctx context.Context, query string, limit int) ([]*models.Persona, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, handle, display_name, bio, avatar_url, cover_url, persona_type, genre_tags,
		       follower_count, following_count, post_count, created_at, updated_at
		FROM personas
		WHERE persona_type = 'visible'
		  AND (
		    handle ILIKE $1
		    OR display_name ILIKE $1
		    OR bio ILIKE $1
		    OR EXISTS (SELECT 1 FROM unnest(COALESCE(genre_tags, ARRAY[]::text[])) tag WHERE tag ILIKE $2)
		  )
		ORDER BY
		  CASE WHEN handle ILIKE $3 THEN 0 ELSE 1 END,
		  follower_count DESC,
		  created_at DESC
		LIMIT $4
	`, textMatchParam(query), tagMatchParam(query), query+"%", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var personas []*models.Persona
	for rows.Next() {
		var persona models.Persona
		if err := rows.Scan(
			&persona.ID,
			&persona.UserID,
			&persona.Handle,
			&persona.DisplayName,
			&persona.Bio,
			&persona.AvatarURL,
			&persona.CoverURL,
			&persona.PersonaType,
			&persona.GenreTags,
			&persona.FollowerCount,
			&persona.FollowingCount,
			&persona.PostCount,
			&persona.CreatedAt,
			&persona.UpdatedAt,
		); err != nil {
			return nil, err
		}
		personas = append(personas, &persona)
	}
	return personas, rows.Err()
}

func (r *pgRepository) searchPosts(ctx context.Context, query string, limit int) ([]*PostResult, error) {
	rows, err := r.db.Query(ctx, `
		SELECT p.id, p.author_user_id, p.persona_id, p.posting_mode, p.event_id, p.body, p.post_type,
		       COALESCE(p.media_url, ''), COALESCE(p.media_type, ''), COALESCE(p.location, ''),
		       p.like_count, p.comment_count, p.repost_count, p.is_repost, p.repost_of, p.created_at,
		       pe.id, pe.user_id, COALESCE(pe.handle, ''), COALESCE(pe.display_name, ''),
		       COALESCE(pe.bio, ''), COALESCE(pe.avatar_url, ''), COALESCE(pe.cover_url, ''),
		       COALESCE(pe.persona_type, ''), COALESCE(pe.genre_tags, ARRAY[]::text[]),
		       COALESCE(pe.follower_count, 0), COALESCE(pe.following_count, 0), COALESCE(pe.post_count, 0),
		       COALESCE(pe.created_at, now()), COALESCE(pe.updated_at, now())
		FROM posts p
		LEFT JOIN personas pe ON pe.id = p.persona_id
		WHERE (p.posting_mode = 'anonymous' OR pe.persona_type = 'visible')
		  AND (
		    p.body ILIKE $1
		    OR p.location ILIKE $1
		    OR EXISTS (SELECT 1 FROM unnest(COALESCE(pe.genre_tags, ARRAY[]::text[])) tag WHERE tag ILIKE $2)
		  )
		ORDER BY p.created_at DESC
		LIMIT $3
	`, textMatchParam(query), tagMatchParam(query), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*PostResult
	for rows.Next() {
		result, err := scanPostResult(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, rows.Err()
}

func (r *pgRepository) searchEvents(ctx context.Context, query string, limit int) ([]*models.Event, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, title, venue, location, event_date, description, COALESCE(cover_url, ''),
		       COALESCE(ticket_url, ''), price_ngn, genre_tags, organizer_id, created_at
		FROM events
		WHERE title ILIKE $1
		   OR venue ILIKE $1
		   OR location ILIKE $1
		   OR description ILIKE $1
		   OR EXISTS (SELECT 1 FROM unnest(COALESCE(genre_tags, ARRAY[]::text[])) tag WHERE tag ILIKE $2)
		ORDER BY
		  CASE WHEN event_date >= now() THEN 0 ELSE 1 END,
		  event_date ASC
		LIMIT $3
	`, textMatchParam(query), tagMatchParam(query), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*models.Event
	for rows.Next() {
		var event models.Event
		if err := rows.Scan(
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
		); err != nil {
			return nil, err
		}
		events = append(events, &event)
	}
	return events, rows.Err()
}

type postResultScanner interface {
	Scan(dest ...any) error
}

func scanPostResult(scanner postResultScanner) (*PostResult, error) {
	var post models.Post
	var persona models.Persona
	var personaID uuid.NullUUID
	var eventID uuid.NullUUID
	var repostOf uuid.NullUUID
	var joinedPersonaID uuid.NullUUID
	var joinedPersonaUserID uuid.NullUUID
	var joinedPersonaCreatedAt time.Time
	var joinedPersonaUpdatedAt time.Time

	err := scanner.Scan(
		&post.ID,
		&post.AuthorUserID,
		&personaID,
		&post.PostingMode,
		&eventID,
		&post.Body,
		&post.PostType,
		&post.MediaURL,
		&post.MediaType,
		&post.Location,
		&post.LikeCount,
		&post.CommentCount,
		&post.RepostCount,
		&post.IsRepost,
		&repostOf,
		&post.CreatedAt,
		&joinedPersonaID,
		&joinedPersonaUserID,
		&persona.Handle,
		&persona.DisplayName,
		&persona.Bio,
		&persona.AvatarURL,
		&persona.CoverURL,
		&persona.PersonaType,
		&persona.GenreTags,
		&persona.FollowerCount,
		&persona.FollowingCount,
		&persona.PostCount,
		&joinedPersonaCreatedAt,
		&joinedPersonaUpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if personaID.Valid {
		post.PersonaID = &personaID.UUID
	}
	if eventID.Valid {
		post.EventID = &eventID.UUID
	}
	if repostOf.Valid {
		post.RepostOf = &repostOf.UUID
	}

	var joinedPersona *models.Persona
	if joinedPersonaID.Valid {
		persona.ID = joinedPersonaID.UUID
		persona.UserID = joinedPersonaUserID.UUID
		persona.CreatedAt = joinedPersonaCreatedAt
		persona.UpdatedAt = joinedPersonaUpdatedAt
		joinedPersona = &persona
	}

	return &PostResult{Post: &post, Persona: joinedPersona}, nil
}
