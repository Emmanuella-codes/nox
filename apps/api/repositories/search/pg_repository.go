package search

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgRepository struct {
	db *pgxpool.Pool
}

func newPgRepository(db *pgxpool.Pool) *pgRepository {
	return &pgRepository{db: db}
}

func (r *pgRepository) Search(ctx context.Context, query string, options Options) (*Results, error) {
	normalizedOptions := NormalizeOptions(options)
	fetchLimit := normalizedOptions.Limit + 1
	results := &Results{}
	var err error
	if normalizedOptions.Scope == "all" || normalizedOptions.Scope == "personas" {
		results.Personas, err = r.searchPersonas(ctx, query, fetchLimit, normalizedOptions.Offset)
		if err != nil {
			return nil, err
		}
	}
	if normalizedOptions.Scope == "all" || normalizedOptions.Scope == "posts" {
		results.Posts, err = r.searchPosts(ctx, query, fetchLimit, normalizedOptions.Offset)
		if err != nil {
			return nil, err
		}
	}
	if normalizedOptions.Scope == "all" || normalizedOptions.Scope == "events" {
		results.Events, err = r.searchEvents(ctx, query, fetchLimit, normalizedOptions.Offset)
		if err != nil {
			return nil, err
		}
	}
	if normalizedOptions.Scope == "all" || normalizedOptions.Scope == "hashtags" {
		results.Hashtags, err = r.searchHashtags(ctx, query, fetchLimit, normalizedOptions.Offset)
		if err != nil {
			return nil, err
		}
	}
	if normalizedOptions.Scope == "all" || normalizedOptions.Scope == "sets" {
		results.Sets, err = r.searchSets(ctx, query, fetchLimit, normalizedOptions.Offset)
		if err != nil {
			return nil, err
		}
	}
	results.HasMore = trimResults(results, normalizedOptions.Limit)
	return results, nil
}

func (r *pgRepository) searchPersonas(ctx context.Context, query string, limit int, offset int) ([]*models.Persona, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, handle, display_name, bio, avatar_url, cover_url, persona_type, category, genre_tags,
		       follower_count, following_count, post_count, created_at, updated_at
		FROM personas
		WHERE persona_type = 'visible'
		  AND (
		    handle ILIKE $1
		    OR display_name ILIKE $1
		    OR bio ILIKE $1
		    OR EXISTS (SELECT 1 FROM unnest(COALESCE(genre_tags, ARRAY[]::text[])) tag WHERE tag ILIKE $2)
		    OR similarity(handle, $3) > 0.25
		    OR similarity(display_name, $3) > 0.25
		  )
		ORDER BY
		  CASE WHEN lower(handle) = lower($3) THEN 0 ELSE 1 END,
		  CASE WHEN handle ILIKE $4 THEN 0 ELSE 1 END,
		  GREATEST(similarity(handle, $3), similarity(display_name, $3), similarity(COALESCE(bio, ''), $3)) DESC,
		  follower_count DESC,
		  created_at DESC,
		  id DESC
		LIMIT $5 OFFSET $6
	`, textMatchParam(query), tagMatchParam(query), query, prefixMatchParam(query), limit, offset)
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
			&persona.Category,
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

func (r *pgRepository) searchPosts(ctx context.Context, query string, limit int, offset int) ([]*PostResult, error) {
	rows, err := r.db.Query(ctx, `
		SELECT p.id, p.author_user_id, p.persona_id, p.posting_mode, p.event_id, p.body, p.post_type,
		       COALESCE(p.media_url, ''), COALESCE(p.media_type, ''), COALESCE(p.location, ''),
		       p.like_count, p.comment_count, p.repost_count, p.is_repost, p.repost_of, p.created_at,
		       pe.id, pe.user_id, COALESCE(pe.handle, ''), COALESCE(pe.display_name, ''),
		       COALESCE(pe.bio, ''), COALESCE(pe.avatar_url, ''), COALESCE(pe.cover_url, ''),
		       COALESCE(pe.persona_type, ''), COALESCE(pe.category, 'patron'), COALESCE(pe.genre_tags, ARRAY[]::text[]),
		       COALESCE(pe.follower_count, 0), COALESCE(pe.following_count, 0), COALESCE(pe.post_count, 0),
		       COALESCE(pe.created_at, now()), COALESCE(pe.updated_at, now())
		FROM posts p
		LEFT JOIN personas pe ON pe.id = p.persona_id
		WHERE (p.posting_mode = 'anonymous' OR pe.persona_type = 'visible')
		  AND (
		    p.body ILIKE $1
		    OR p.location ILIKE $1
		    OR EXISTS (SELECT 1 FROM unnest(COALESCE(pe.genre_tags, ARRAY[]::text[])) tag WHERE tag ILIKE $2)
		    OR similarity(p.body, $3) > 0.18
		    OR similarity(COALESCE(p.location, ''), $3) > 0.25
		  )
		ORDER BY
		  CASE WHEN p.body ILIKE $4 THEN 0 ELSE 1 END,
		  GREATEST(similarity(p.body, $3), similarity(COALESCE(p.location, ''), $3)) DESC,
		  (p.like_count + p.comment_count + p.repost_count) DESC,
		  p.created_at DESC,
		  p.id DESC
		LIMIT $5 OFFSET $6
	`, textMatchParam(query), tagMatchParam(query), query, prefixMatchParam(query), limit, offset)
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

func (r *pgRepository) searchEvents(ctx context.Context, query string, limit int, offset int) ([]*models.Event, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, title, venue, location, event_date, description, COALESCE(cover_url, ''),
		       COALESCE(ticket_url, ''), price_ngn, genre_tags, organizer_id, created_at
		FROM events
		WHERE title ILIKE $1
		   OR venue ILIKE $1
		   OR location ILIKE $1
		   OR description ILIKE $1
		   OR EXISTS (SELECT 1 FROM unnest(COALESCE(genre_tags, ARRAY[]::text[])) tag WHERE tag ILIKE $2)
		   OR similarity(title, $3) > 0.25
		   OR similarity(venue, $3) > 0.25
		   OR similarity(location, $3) > 0.25
		ORDER BY
		  CASE WHEN lower(title) = lower($3) THEN 0 ELSE 1 END,
		  CASE WHEN title ILIKE $4 THEN 0 ELSE 1 END,
		  GREATEST(similarity(title, $3), similarity(venue, $3), similarity(location, $3), similarity(COALESCE(description, ''), $3)) DESC,
		  CASE WHEN event_date >= now() THEN 0 ELSE 1 END,
		  event_date ASC,
		  id DESC
		LIMIT $5 OFFSET $6
	`, textMatchParam(query), tagMatchParam(query), query, prefixMatchParam(query), limit, offset)
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

func (r *pgRepository) searchHashtags(ctx context.Context, query string, limit int, offset int) ([]*models.Hashtag, error) {
	normalizedQuery := normalizeHashtagQuery(query)
	rows, err := r.db.Query(ctx, `
		SELECT id, tag, post_count, created_at
		FROM hashtags
		WHERE post_count > 0
		  AND (
		    tag ILIKE $1
		    OR similarity(tag, $2) > 0.25
		  )
		ORDER BY
		  CASE WHEN lower(tag) = lower($2) THEN 0 ELSE 1 END,
		  CASE WHEN tag ILIKE $3 THEN 0 ELSE 1 END,
		  similarity(tag, $2) DESC,
		  post_count DESC,
		  tag ASC
		LIMIT $4 OFFSET $5
	`, tagMatchParam(normalizedQuery), normalizedQuery, prefixMatchParam(normalizedQuery), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hashtags []*models.Hashtag
	for rows.Next() {
		var hashtag models.Hashtag
		if err := rows.Scan(&hashtag.ID, &hashtag.Tag, &hashtag.PostCount, &hashtag.CreatedAt); err != nil {
			return nil, err
		}
		hashtags = append(hashtags, &hashtag)
	}
	return hashtags, rows.Err()
}

// trimResults truncates grouped results to the requested limit and reports whether more rows exist.
func trimResults(results *Results, limit int) bool {
	hasMore := false
	if len(results.Personas) > limit {
		results.Personas = results.Personas[:limit]
		hasMore = true
	}
	if len(results.Posts) > limit {
		results.Posts = results.Posts[:limit]
		hasMore = true
	}
	if len(results.Events) > limit {
		results.Events = results.Events[:limit]
		hasMore = true
	}
	if len(results.Hashtags) > limit {
		results.Hashtags = results.Hashtags[:limit]
		hasMore = true
	}
	if len(results.Sets) > limit {
		results.Sets = results.Sets[:limit]
		hasMore = true
	}
	return hasMore
}
