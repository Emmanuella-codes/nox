package hashtag

import (
	"context"
	"errors"

	"github.com/emmanuella-codes/nox/models"
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
}

type hashtagScanner interface {
	Scan(dest ...any) error
}

func (r *pgRepository) SyncPostHashtags(ctx context.Context, postID uuid.UUID, tags []string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := deletePostHashtags(ctx, tx, postID); err != nil {
		return err
	}

	for _, tag := range tags {
		normalized := NormalizeTag(tag)
		if normalized == "" {
			continue
		}

		var hashtagID uuid.UUID
		if err := tx.QueryRow(ctx, `
			INSERT INTO hashtags (tag, post_count)
			VALUES ($1, 0)
			ON CONFLICT (tag) DO UPDATE SET tag = EXCLUDED.tag
			RETURNING id
		`, normalized).Scan(&hashtagID); err != nil {
			return err
		}

		commandTag, err := tx.Exec(ctx, `
			INSERT INTO post_hashtags (post_id, hashtag_id)
			VALUES ($1, $2)
			ON CONFLICT DO NOTHING
		`, postID, hashtagID)
		if err != nil {
			return err
		}
		if commandTag.RowsAffected() > 0 {
			if _, err := tx.Exec(ctx, `UPDATE hashtags SET post_count = post_count + 1 WHERE id = $1`, hashtagID); err != nil {
				return err
			}
		}
	}

	return tx.Commit(ctx)
}

func (r *pgRepository) DeletePostHashtags(ctx context.Context, postID uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := deletePostHashtags(ctx, tx, postID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *pgRepository) FindTagsByPostIDs(ctx context.Context, postIDs []uuid.UUID) (map[uuid.UUID][]string, error) {
	tagsByPost := make(map[uuid.UUID][]string, len(postIDs))
	if len(postIDs) == 0 {
		return tagsByPost, nil
	}

	rows, err := r.db.Query(ctx, `
		SELECT ph.post_id, h.tag
		FROM post_hashtags ph
		INNER JOIN hashtags h ON h.id = ph.hashtag_id
		WHERE ph.post_id = ANY($1)
		ORDER BY ph.created_at ASC
	`, postIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var postID uuid.UUID
		var tag string
		if err := rows.Scan(&postID, &tag); err != nil {
			return nil, err
		}
		tagsByPost[postID] = append(tagsByPost[postID], tag)
	}
	return tagsByPost, rows.Err()
}

func (r *pgRepository) FindTrending(ctx context.Context, limit int) ([]*models.Hashtag, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, tag, post_count, created_at
		FROM hashtags
		WHERE post_count > 0
		ORDER BY post_count DESC, tag ASC
		LIMIT $1
	`, normalizeLimit(limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hashtags []*models.Hashtag
	for rows.Next() {
		hashtag, err := scanHashtag(rows)
		if err != nil {
			return nil, err
		}
		hashtags = append(hashtags, hashtag)
	}
	return hashtags, rows.Err()
}

func (r *pgRepository) FindByTag(ctx context.Context, tag string) (*models.Hashtag, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, tag, post_count, created_at
		FROM hashtags
		WHERE tag = $1
	`, NormalizeTag(tag))

	hashtag, err := scanHashtag(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return hashtag, err
}

func (r *pgRepository) FindPostsByTag(ctx context.Context, tag string, limit int) ([]*models.Post, error) {
	rows, err := r.db.Query(ctx, `
		SELECT p.id, p.author_user_id, p.persona_id, p.posting_mode, p.event_id, p.body, p.post_type,
		       COALESCE(p.media_url, ''), COALESCE(p.media_type, ''), COALESCE(p.location, ''),
		       p.like_count, p.comment_count, p.repost_count, p.is_repost, p.repost_of, p.created_at
		FROM hashtags h
		INNER JOIN post_hashtags ph ON ph.hashtag_id = h.id
		INNER JOIN posts p ON p.id = ph.post_id
		LEFT JOIN personas pe ON pe.id = p.persona_id
		WHERE h.tag = $1
		  AND (p.posting_mode = 'anonymous' OR pe.persona_type = 'visible')
		ORDER BY p.created_at DESC
		LIMIT $2
	`, NormalizeTag(tag), normalizeLimit(limit))
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
	return posts, rows.Err()
}

func deletePostHashtags(ctx context.Context, tx execQuerier, postID uuid.UUID) error {
	rows, err := tx.Query(ctx, `
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

	if _, err := tx.Exec(ctx, `DELETE FROM post_hashtags WHERE post_id = $1`, postID); err != nil {
		return err
	}
	for _, hashtagID := range hashtagIDs {
		if _, err := tx.Exec(ctx, `UPDATE hashtags SET post_count = GREATEST(post_count - 1, 0) WHERE id = $1`, hashtagID); err != nil {
			return err
		}
	}
	return nil
}

func scanHashtag(scanner hashtagScanner) (*models.Hashtag, error) {
	var hashtag models.Hashtag
	err := scanner.Scan(&hashtag.ID, &hashtag.Tag, &hashtag.PostCount, &hashtag.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &hashtag, nil
}

func scanPost(scanner hashtagScanner) (*models.Post, error) {
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

func normalizeLimit(limit int) int {
	if limit <= 0 {
		return 20
	}
	if limit > 100 {
		return 100
	}
	return limit
}
