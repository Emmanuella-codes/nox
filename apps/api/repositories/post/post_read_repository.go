package post

import (
	"context"
	"time"

	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
)

// FindPostByID fetches one post by id.
func (r *pgRepository) FindPostByID(ctx context.Context, postID uuid.UUID) (*models.Post, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, author_user_id, persona_id, posting_mode, event_id, body, post_type,
		       COALESCE(media_url, ''), COALESCE(media_type, ''), COALESCE(location, ''),
		       like_count, comment_count, repost_count, is_repost, repost_of, created_at
		FROM posts
		WHERE id = $1
	`, postID)
	post, err := scanPost(row)
	if err != nil {
		return nil, mapPostError(err)
	}
	return post, nil
}

// FindPostsByPersonaID fetches public posts for a public profile.
func (r *pgRepository) FindPostsByPersonaID(ctx context.Context, personaID uuid.UUID, limit int) ([]*models.Post, error) {
	return r.findPosts(ctx, `
		SELECT id, author_user_id, persona_id, posting_mode, event_id, body, post_type,
		       COALESCE(media_url, ''), COALESCE(media_type, ''), COALESCE(location, ''),
		       like_count, comment_count, repost_count, is_repost, repost_of, created_at
		FROM posts
		WHERE persona_id = $1 AND posting_mode = 'public'
		ORDER BY created_at DESC
		LIMIT $2
	`, personaID, normalizeLimit(limit))
}

// FindPostsByAuthorUserID fetches all posts for one owner, including anonymous posts.
func (r *pgRepository) FindPostsByAuthorUserID(ctx context.Context, authorUserID uuid.UUID, limit int) ([]*models.Post, error) {
	return r.findPosts(ctx, `
		SELECT id, author_user_id, persona_id, posting_mode, event_id, body, post_type,
		       COALESCE(media_url, ''), COALESCE(media_type, ''), COALESCE(location, ''),
		       like_count, comment_count, repost_count, is_repost, repost_of, created_at
		FROM posts
		WHERE author_user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`, authorUserID, normalizeLimit(limit))
}

// FindFeedPosts fetches the mixed public and anonymous feed.
func (r *pgRepository) FindFeedPosts(ctx context.Context, personaID uuid.UUID, options FeedOptions) ([]*models.Post, error) {
	options = NormalizeFeedOptions(options)
	limit := options.Limit + 1
	cursorCreatedAt, cursorID := nullableFeedCursor(options.Cursor)
	return r.findPosts(ctx, `
		WITH viewer AS (
			SELECT user_id FROM personas WHERE id = $1
		)
		SELECT p.id, p.author_user_id, p.persona_id, p.posting_mode, p.event_id, p.body, p.post_type,
		       COALESCE(p.media_url, ''), COALESCE(p.media_type, ''), COALESCE(p.location, ''),
		       p.like_count, p.comment_count, p.repost_count, p.is_repost, p.repost_of, p.created_at
		FROM posts p
		LEFT JOIN persona_follows pf
		  ON pf.follower_id = $1 AND pf.following_id = p.persona_id
		CROSS JOIN viewer v
		WHERE (
			p.posting_mode = 'public'
			OR p.posting_mode = 'anonymous'
			OR p.author_user_id = v.user_id
		)
		  AND NOT EXISTS (
			SELECT 1
			FROM user_blocks ub
			WHERE (ub.blocker_user_id = v.user_id AND ub.blocked_user_id = p.author_user_id)
			   OR (ub.blocker_user_id = p.author_user_id AND ub.blocked_user_id = v.user_id)
		  )
		  AND NOT EXISTS (
			SELECT 1
			FROM user_mutes um
			WHERE um.user_id = v.user_id AND um.muted_user_id = p.author_user_id
		  )
		  AND (
			$2::timestamptz IS NULL
			OR p.created_at < $2
			OR (p.created_at = $2 AND p.id < $3)
		  )
		ORDER BY
		  CASE
		    WHEN p.author_user_id = v.user_id THEN 0
		    WHEN pf.follower_id IS NOT NULL THEN 1
		    WHEN p.posting_mode = 'public' THEN 2
		    ELSE 3
		  END ASC,
		  p.created_at DESC,
		  p.id DESC
		LIMIT $4
	`, personaID, cursorCreatedAt, cursorID, limit)
}

// FindFollowingFeedPosts fetches followed public-profile posts.
func (r *pgRepository) FindFollowingFeedPosts(ctx context.Context, personaID uuid.UUID, options FeedOptions) ([]*models.Post, error) {
	options = NormalizeFeedOptions(options)
	limit := options.Limit + 1
	cursorCreatedAt, cursorID := nullableFeedCursor(options.Cursor)
	return r.findPosts(ctx, `
		WITH viewer AS (
			SELECT user_id FROM personas WHERE id = $1
		)
		SELECT p.id, p.author_user_id, p.persona_id, p.posting_mode, p.event_id, p.body, p.post_type,
		       COALESCE(p.media_url, ''), COALESCE(p.media_type, ''), COALESCE(p.location, ''),
		       p.like_count, p.comment_count, p.repost_count, p.is_repost, p.repost_of, p.created_at
		FROM posts p
		LEFT JOIN persona_follows pf
		  ON pf.follower_id = $1 AND pf.following_id = p.persona_id
		CROSS JOIN viewer v
		WHERE p.posting_mode = 'public'
		  AND (pf.follower_id IS NOT NULL OR p.author_user_id = v.user_id)
		  AND NOT EXISTS (
			SELECT 1
			FROM user_blocks ub
			WHERE (ub.blocker_user_id = v.user_id AND ub.blocked_user_id = p.author_user_id)
			   OR (ub.blocker_user_id = p.author_user_id AND ub.blocked_user_id = v.user_id)
		  )
		  AND NOT EXISTS (
			SELECT 1
			FROM user_mutes um
			WHERE um.user_id = v.user_id AND um.muted_user_id = p.author_user_id
		  )
		  AND (
			$2::timestamptz IS NULL
			OR p.created_at < $2
			OR (p.created_at = $2 AND p.id < $3)
		  )
		ORDER BY
		  CASE WHEN p.author_user_id = v.user_id THEN 0 ELSE 1 END ASC,
		  p.created_at DESC,
		  p.id DESC
		LIMIT $4
	`, personaID, cursorCreatedAt, cursorID, limit)
}

// findPosts executes a repeated post list query.
func (r *pgRepository) findPosts(ctx context.Context, sql string, args ...any) ([]*models.Post, error) {
	rows, err := r.db.Query(ctx, sql, args...)
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

// nullableFeedCursor maps one optional feed cursor into nullable SQL arguments.
func nullableFeedCursor(cursor *FeedCursor) (*time.Time, *uuid.UUID) {
	if cursor == nil {
		return nil, nil
	}
	createdAt := cursor.CreatedAt
	id := cursor.ID
	return &createdAt, &id
}
