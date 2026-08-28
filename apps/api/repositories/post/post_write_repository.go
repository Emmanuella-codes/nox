package post

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/post/dtos"
	hashtag_repo "github.com/emmanuella-codes/nox/repositories/hashtag"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type execQuerier interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// CreatePost inserts one post row.
func (r *pgRepository) CreatePost(ctx context.Context, authorUserID uuid.UUID, dto dtos.CreatePostDTO) (*models.Post, error) {
	return createPost(ctx, r.db, authorUserID, dto)
}

// CreatePostWithHashtags inserts a post and associated hashtags.
func (r *pgRepository) CreatePostWithHashtags(ctx context.Context, authorUserID uuid.UUID, dto dtos.CreatePostDTO, tags []string) (*models.Post, error) {
	return r.CreatePostWithHashtagsAndMedia(ctx, authorUserID, dto, tags, nil)
}

// CreatePostWithHashtagsAndMedia inserts a post and associated hashtags and media.
func (r *pgRepository) CreatePostWithHashtagsAndMedia(ctx context.Context, authorUserID uuid.UUID, dto dtos.CreatePostDTO, tags []string, mediaAssetIDs []uuid.UUID) (*models.Post, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	post, err := createPost(ctx, tx, authorUserID, dto)
	if err != nil {
		return nil, err
	}
	if err := syncPostHashtags(ctx, tx, post.ID, tags); err != nil {
		return nil, err
	}
	if err := syncPostMediaAssets(ctx, tx, post.ID, mediaAssetIDs); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return post, nil
}

// DeletePost removes one post row.
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

// DeletePostWithHashtags removes one post and its hashtag links.
func (r *pgRepository) DeletePostWithHashtags(ctx context.Context, postID uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err := deletePostHashtags(ctx, tx, postID); err != nil {
		return err
	}
	commandTag, err := tx.Exec(ctx, `DELETE FROM posts WHERE id = $1`, postID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return ErrPostNotFound
	}
	return tx.Commit(ctx)
}

// createPost inserts one post row through the supplied query executor.
func createPost(ctx context.Context, db execQuerier, authorUserID uuid.UUID, dto dtos.CreatePostDTO) (*models.Post, error) {
	var eventID *uuid.UUID
	if dto.EventID != uuid.Nil {
		eventID = &dto.EventID
	}
	row := db.QueryRow(ctx, `
		INSERT INTO posts (
			author_user_id, persona_id, posting_mode, event_id, body, post_type, media_url, media_type, location
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, author_user_id, persona_id, posting_mode, event_id, body, post_type,
		          COALESCE(media_url, ''), COALESCE(media_type, ''), COALESCE(location, ''),
		          like_count, comment_count, repost_count, is_repost, repost_of, created_at
	`, authorUserID, dto.PersonaID, dto.PostingMode, eventID, dto.Body, dto.PostType, emptyToNil(dto.MediaURL), emptyToNil(string(dto.MediaType)), emptyToNil(dto.Location))
	post, err := scanPost(row)
	if err != nil {
		return nil, mapPostError(err)
	}
	return post, nil
}

// syncPostHashtags replaces hashtag links for a post.
func syncPostHashtags(ctx context.Context, db execQuerier, postID uuid.UUID, tags []string) error {
	if err := deletePostHashtags(ctx, db, postID); err != nil {
		return err
	}
	for _, tag := range tags {
		normalized := hashtag_repo.NormalizeTag(tag)
		if normalized == "" {
			continue
		}
		var hashtagID uuid.UUID
		if err := db.QueryRow(ctx, `
			INSERT INTO hashtags (tag, post_count)
			VALUES ($1, 0)
			ON CONFLICT (tag) DO UPDATE SET tag = EXCLUDED.tag
			RETURNING id
		`, normalized).Scan(&hashtagID); err != nil {
			return err
		}
		commandTag, err := db.Exec(ctx, `
			INSERT INTO post_hashtags (post_id, hashtag_id)
			VALUES ($1, $2)
			ON CONFLICT DO NOTHING
		`, postID, hashtagID)
		if err != nil {
			return err
		}
		if commandTag.RowsAffected() > 0 {
			if _, err := db.Exec(ctx, `UPDATE hashtags SET post_count = post_count + 1 WHERE id = $1`, hashtagID); err != nil {
				return err
			}
		}
	}
	return nil
}

// syncPostMediaAssets stores media links for a post in order.
func syncPostMediaAssets(ctx context.Context, db execQuerier, postID uuid.UUID, mediaAssetIDs []uuid.UUID) error {
	for position, mediaAssetID := range mediaAssetIDs {
		if mediaAssetID == uuid.Nil {
			continue
		}
		if _, err := db.Exec(ctx, `
			INSERT INTO post_media_assets (post_id, media_asset_id, position)
			VALUES ($1, $2, $3)
		`, postID, mediaAssetID, position); err != nil {
			return err
		}
	}
	return nil
}

// deletePostHashtags removes hashtag links and decrements tag counters.
func deletePostHashtags(ctx context.Context, db execQuerier, postID uuid.UUID) error {
	rows, err := db.Query(ctx, `SELECT hashtag_id FROM post_hashtags WHERE post_id = $1`, postID)
	if err != nil {
		return err
	}
	defer rows.Close()
	var hashtagIDs []uuid.UUID
	for rows.Next() {
		var hashtagID uuid.UUID
		if err := rows.Scan(&hashtagID); err != nil {
			return err
		}
		hashtagIDs = append(hashtagIDs, hashtagID)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	if len(hashtagIDs) == 0 {
		return nil
	}
	if _, err := db.Exec(ctx, `DELETE FROM post_hashtags WHERE post_id = $1`, postID); err != nil {
		return err
	}
	for _, hashtagID := range hashtagIDs {
		if _, err := db.Exec(ctx, `UPDATE hashtags SET post_count = GREATEST(post_count - 1, 0) WHERE id = $1`, hashtagID); err != nil {
			return err
		}
	}
	return nil
}
