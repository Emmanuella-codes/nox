package post

import (
	"errors"
	"time"

	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type postScanner interface {
	Scan(dest ...any) error
}

type anonymousThreadIdentityScanner interface {
	Scan(dest ...any) error
}

// scanAnonymousThreadIdentity maps one anonymous identity row.
func scanAnonymousThreadIdentity(scanner anonymousThreadIdentityScanner) (*models.AnonymousThreadIdentity, error) {
	var identity models.AnonymousThreadIdentity
	var createdAt time.Time
	err := scanner.Scan(
		&identity.ID,
		&identity.ThreadID,
		&identity.UserID,
		&identity.PersonaID,
		&identity.AnonymousHandle,
		&identity.AnonymousAvatarKey,
		&createdAt,
	)
	if err != nil {
		return nil, err
	}
	identity.CreatedAt = createdAt
	return &identity, nil
}

// scanPost maps one post row.
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

// mapPostError maps not-found database errors to repository errors.
func mapPostError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrPostNotFound
	}
	return err
}

// normalizeLimit clamps post list limits.
func normalizeLimit(limit int) int {
	if limit <= 0 {
		return 20
	}
	if limit > 100 {
		return 100
	}
	return limit
}

// emptyToNil converts empty strings to nullable insert values.
func emptyToNil(value string) any {
	if value == "" {
		return nil
	}
	return value
}
