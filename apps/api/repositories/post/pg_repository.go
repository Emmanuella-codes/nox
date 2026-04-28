package post

import (
	"context"
	"errors"

	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/post/dtos"
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

func (r *pgRepository) CreatePost(ctx context.Context, personaID uuid.UUID, dto dtos.CreatePostDTO) (*models.Post, error) {
	var eventID *uuid.UUID
	if dto.EventID != uuid.Nil {
		eventID = &dto.EventID
	}

	row := r.db.QueryRow(ctx, `
		INSERT INTO posts (
			persona_id, event_id, body, post_type, media_url, media_type, location
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, persona_id, event_id, body, post_type, COALESCE(media_url, ''), COALESCE(media_type, ''),
		          COALESCE(location, ''), like_count, comment_count, repost_count, is_repost, repost_of, created_at
	`, personaID, eventID, dto.Body, dto.PostType, emptyToNil(dto.MediaURL), emptyToNil(string(dto.MediaType)), emptyToNil(dto.Location))

	post, err := scanPost(row)
	if err != nil {
		return nil, mapPostError(err)
	}

	return post, nil
}

func (r *pgRepository) FindPostByID(ctx context.Context, postID uuid.UUID) (*models.Post, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, persona_id, event_id, body, post_type, COALESCE(media_url, ''), COALESCE(media_type, ''),
		       COALESCE(location, ''), like_count, comment_count, repost_count, is_repost, repost_of, created_at
		FROM posts
		WHERE id = $1
	`, postID)

	post, err := scanPost(row)
	if err != nil {
		return nil, mapPostError(err)
	}

	return post, nil
}

func (r *pgRepository) FindPostsByPersonaID(ctx context.Context, personaID uuid.UUID, limit int) ([]*models.Post, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, persona_id, event_id, body, post_type, COALESCE(media_url, ''), COALESCE(media_type, ''),
		       COALESCE(location, ''), like_count, comment_count, repost_count, is_repost, repost_of, created_at
		FROM posts
		WHERE persona_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`, personaID, normalizeLimit(limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		post, err := scanPost(rows)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *pgRepository) FindFeedPosts(ctx context.Context, personaID uuid.UUID, limit int) ([]*models.Post, error) {
	rows, err := r.db.Query(ctx, `
		SELECT p.id, p.persona_id, p.event_id, p.body, p.post_type, COALESCE(p.media_url, ''), COALESCE(p.media_type, ''),
		       COALESCE(p.location, ''), p.like_count, p.comment_count, p.repost_count, p.is_repost, p.repost_of, p.created_at
		FROM posts p
		JOIN personas pe ON pe.id = p.persona_id
		WHERE p.persona_id = $1 OR pe.persona_type = 'visible'
		ORDER BY p.created_at DESC
		LIMIT $2
	`, personaID, normalizeLimit(limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		post, err := scanPost(rows)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *pgRepository) DeletePost(ctx context.Context, postID uuid.UUID) error {
	commandTag, err := r.db.Exec(ctx, `DELETE FROM posts WHERE id = $1`, postID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return ErrPostNotFound
	}
	return nil
}

type postScanner interface {
	Scan(dest ...any) error
}

func scanPost(scanner postScanner) (*models.Post, error) {
	var post models.Post
	var eventID uuid.NullUUID
	var repostOf uuid.NullUUID
	err := scanner.Scan(
		&post.ID,
		&post.PersonaID,
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
	)
	if err != nil {
		return nil, err
	}
	if eventID.Valid {
		post.EventID = &eventID.UUID
	}
	if repostOf.Valid {
		post.RepostOf = &repostOf.UUID
	}
	return &post, nil
}

func mapPostError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrPostNotFound
	}
	return err
}

func normalizeLimit(limit int) int {
	if limit <= 0 {
		return 20
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
