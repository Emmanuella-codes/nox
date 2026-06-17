package post

import (
	"context"
	"errors"
	"time"

	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/post/dtos"
	hashtag_repo "github.com/emmanuella-codes/nox/repositories/hashtag"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgRepository struct {
	db *pgxpool.Pool
}

func newPgRepository(db *pgxpool.Pool) *pgRepository {
	return &pgRepository{db: db}
}

type execQuerier interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func (r *pgRepository) CreatePost(ctx context.Context, authorUserID uuid.UUID, dto dtos.CreatePostDTO) (*models.Post, error) {
	return createPost(ctx, r.db, authorUserID, dto)
}

func (r *pgRepository) CreatePostWithHashtags(ctx context.Context, authorUserID uuid.UUID, dto dtos.CreatePostDTO, tags []string) (*models.Post, error) {
	return r.CreatePostWithHashtagsAndMedia(ctx, authorUserID, dto, tags, nil)
}

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

func (r *pgRepository) FindPostsByPersonaID(ctx context.Context, personaID uuid.UUID, limit int) ([]*models.Post, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, author_user_id, persona_id, posting_mode, event_id, body, post_type,
		       COALESCE(media_url, ''), COALESCE(media_type, ''), COALESCE(location, ''),
		       like_count, comment_count, repost_count, is_repost, repost_of, created_at
		FROM posts
		WHERE persona_id = $1 AND posting_mode = 'public'
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
		SELECT p.id, p.author_user_id, p.persona_id, p.posting_mode, p.event_id, p.body, p.post_type,
		       COALESCE(p.media_url, ''), COALESCE(p.media_type, ''), COALESCE(p.location, ''),
		       p.like_count, p.comment_count, p.repost_count, p.is_repost, p.repost_of, p.created_at
		FROM posts p
		LEFT JOIN personas pe ON pe.id = p.persona_id
		WHERE p.posting_mode = 'anonymous'
		   OR p.persona_id = $1
		   OR (p.posting_mode = 'public' AND pe.persona_type = 'visible')
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

func (r *pgRepository) FindFollowingFeedPosts(ctx context.Context, personaID uuid.UUID, limit int) ([]*models.Post, error) {
	rows, err := r.db.Query(ctx, `
		SELECT p.id, p.author_user_id, p.persona_id, p.posting_mode, p.event_id, p.body, p.post_type,
		       COALESCE(p.media_url, ''), COALESCE(p.media_type, ''), COALESCE(p.location, ''),
		       p.like_count, p.comment_count, p.repost_count, p.is_repost, p.repost_of, p.created_at
		FROM persona_follows pf
		INNER JOIN posts p ON p.persona_id = pf.following_id
		INNER JOIN personas pe ON pe.id = p.persona_id
		WHERE pf.follower_id = $1
		  AND p.posting_mode = 'public'
		  AND pe.persona_type = 'visible'
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

func (r *pgRepository) EnsureAnonymousThreadIdentity(ctx context.Context, threadID uuid.UUID, userID uuid.UUID, personaID uuid.UUID, anonymousHandle string) (*models.AnonymousThreadIdentity, error) {
	row := r.db.QueryRow(ctx, `
		INSERT INTO anonymous_thread_identities (thread_id, user_id, persona_id, anonymous_handle)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (thread_id, persona_id) DO UPDATE
		SET anonymous_handle = anonymous_thread_identities.anonymous_handle
		RETURNING id, thread_id, user_id, persona_id, anonymous_handle, created_at
	`, threadID, userID, personaID, anonymousHandle)
	return scanAnonymousThreadIdentity(row)
}

func (r *pgRepository) FindAnonymousThreadIdentities(ctx context.Context, threadID uuid.UUID, personaIDs []uuid.UUID) (map[uuid.UUID]*models.AnonymousThreadIdentity, error) {
	identities := make(map[uuid.UUID]*models.AnonymousThreadIdentity)
	if len(personaIDs) == 0 {
		return identities, nil
	}

	rows, err := r.db.Query(ctx, `
		SELECT id, thread_id, user_id, persona_id, anonymous_handle, created_at
		FROM anonymous_thread_identities
		WHERE thread_id = $1 AND persona_id = ANY($2)
	`, threadID, personaIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		identity, err := scanAnonymousThreadIdentity(rows)
		if err != nil {
			return nil, err
		}
		identities[identity.PersonaID] = identity
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return identities, nil
}

func (r *pgRepository) FindMediaAssetsByPostIDs(ctx context.Context, postIDs []uuid.UUID) (map[uuid.UUID][]*models.MediaAsset, error) {
	assetsByPost := make(map[uuid.UUID][]*models.MediaAsset)
	if len(postIDs) == 0 {
		return assetsByPost, nil
	}
	rows, err := r.db.Query(ctx, `
		SELECT pma.post_id, ma.id, ma.owner_user_id, ma.owner_persona_id, ma.media_kind, ma.storage_key,
		       ma.playback_url, COALESCE(ma.thumbnail_url, ''), ma.mime_type, ma.duration_seconds,
		       ma.size_bytes, ma.processing_status, ma.created_at, ma.updated_at
		FROM post_media_assets pma
		INNER JOIN media_assets ma ON ma.id = pma.media_asset_id
		WHERE pma.post_id = ANY($1)
		ORDER BY pma.position ASC
	`, postIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var postID uuid.UUID
		var asset models.MediaAsset
		if err := rows.Scan(
			&postID,
			&asset.ID,
			&asset.OwnerUserID,
			&asset.OwnerPersonaID,
			&asset.MediaKind,
			&asset.StorageKey,
			&asset.PlaybackURL,
			&asset.ThumbnailURL,
			&asset.MimeType,
			&asset.DurationSeconds,
			&asset.SizeBytes,
			&asset.ProcessingStatus,
			&asset.CreatedAt,
			&asset.UpdatedAt,
		); err != nil {
			return nil, err
		}
		assetsByPost[postID] = append(assetsByPost[postID], &asset)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return assetsByPost, nil
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

func deletePostHashtags(ctx context.Context, db execQuerier, postID uuid.UUID) error {
	rows, err := db.Query(ctx, `
		SELECT hashtag_id
		FROM post_hashtags
		WHERE post_id = $1
	`, postID)
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

type postScanner interface {
	Scan(dest ...any) error
}

type anonymousThreadIdentityScanner interface {
	Scan(dest ...any) error
}

func scanAnonymousThreadIdentity(scanner anonymousThreadIdentityScanner) (*models.AnonymousThreadIdentity, error) {
	var identity models.AnonymousThreadIdentity
	var createdAt time.Time
	err := scanner.Scan(
		&identity.ID,
		&identity.ThreadID,
		&identity.UserID,
		&identity.PersonaID,
		&identity.AnonymousHandle,
		&createdAt,
	)
	if err != nil {
		return nil, err
	}
	identity.CreatedAt = createdAt
	return &identity, nil
}

func scanPost(scanner postScanner) (*models.Post, error) {
	var post models.Post
	var personaID uuid.NullUUID
	var eventID uuid.NullUUID
	var repostOf uuid.NullUUID
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
