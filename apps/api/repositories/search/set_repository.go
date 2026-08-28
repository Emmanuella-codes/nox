package search

import (
	"context"
	"time"

	"github.com/emmanuella-codes/nox/models"
)

// searchSets searches DJ sets with persona and media attached for discovery responses.
func (r *pgRepository) searchSets(ctx context.Context, query string, limit int, offset int) ([]*models.Set, error) {
	rows, err := r.db.Query(ctx, `
		SELECT s.id, s.author_user_id, s.persona_id, s.media_asset_id, s.title, s.description, s.genre_tags,
		       s.duration_seconds, s.like_count, s.comment_count, s.play_count, s.created_at, s.updated_at,
		       p.id, p.user_id, p.handle, p.display_name, COALESCE(p.avatar_url, ''), p.category,
		       m.id, m.owner_user_id, m.owner_persona_id, m.media_kind, m.playback_url, COALESCE(m.thumbnail_url, ''),
		       m.mime_type, m.duration_seconds, m.size_bytes, m.processing_status, m.created_at, m.updated_at
		FROM sets s
		INNER JOIN personas p ON p.id = s.persona_id
		INNER JOIN media_assets m ON m.id = s.media_asset_id
		WHERE p.persona_type = 'visible'
		  AND (
		    s.title ILIKE $1
		    OR s.description ILIKE $1
		    OR EXISTS (SELECT 1 FROM unnest(COALESCE(s.genre_tags, ARRAY[]::text[])) tag WHERE tag ILIKE $2)
		    OR similarity(s.title, $3) > 0.25
		    OR similarity(s.description, $3) > 0.18
		  )
		ORDER BY
		  CASE WHEN lower(s.title) = lower($3) THEN 0 ELSE 1 END,
		  CASE WHEN s.title ILIKE $4 THEN 0 ELSE 1 END,
		  GREATEST(similarity(s.title, $3), similarity(COALESCE(s.description, ''), $3)) DESC,
		  (s.play_count + s.like_count + s.comment_count) DESC,
		  s.created_at DESC,
		  s.id DESC
		LIMIT $5 OFFSET $6
	`, textMatchParam(query), tagMatchParam(query), query, prefixMatchParam(query), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	sets := make([]*models.Set, 0)
	for rows.Next() {
		set, err := scanSearchSet(rows)
		if err != nil {
			return nil, err
		}
		sets = append(sets, set)
	}
	return sets, rows.Err()
}

type searchSetScanner interface {
	Scan(dest ...any) error
}

// scanSearchSet scans one joined set search row into the set model with nested persona and media.
func scanSearchSet(scanner searchSetScanner) (*models.Set, error) {
	var set models.Set
	var persona models.Persona
	var media models.MediaAsset
	var personaCategory string
	var mediaCreatedAt time.Time
	var mediaUpdatedAt time.Time

	err := scanner.Scan(
		&set.ID,
		&set.AuthorUserID,
		&set.PersonaID,
		&set.MediaAssetID,
		&set.Title,
		&set.Description,
		&set.GenreTags,
		&set.DurationSeconds,
		&set.LikeCount,
		&set.CommentCount,
		&set.PlayCount,
		&set.CreatedAt,
		&set.UpdatedAt,
		&persona.ID,
		&persona.UserID,
		&persona.Handle,
		&persona.DisplayName,
		&persona.AvatarURL,
		&personaCategory,
		&media.ID,
		&media.OwnerUserID,
		&media.OwnerPersonaID,
		&media.MediaKind,
		&media.PlaybackURL,
		&media.ThumbnailURL,
		&media.MimeType,
		&media.DurationSeconds,
		&media.SizeBytes,
		&media.ProcessingStatus,
		&mediaCreatedAt,
		&mediaUpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	persona.Category = models.PersonaCategory(personaCategory)
	media.CreatedAt = mediaCreatedAt
	media.UpdatedAt = mediaUpdatedAt
	set.Persona = &persona
	set.MediaAsset = &media
	return &set, nil
}
