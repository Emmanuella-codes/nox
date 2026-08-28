package search

import (
	"time"

	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
)

type postResultScanner interface {
	Scan(dest ...any) error
}

// scanPostResult scans one joined post search row into the post result shape.
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
		&persona.Category,
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
	if joinedPersonaID.Valid {
		persona.ID = joinedPersonaID.UUID
		persona.UserID = joinedPersonaUserID.UUID
		persona.CreatedAt = joinedPersonaCreatedAt
		persona.UpdatedAt = joinedPersonaUpdatedAt
		return &PostResult{Post: &post, Persona: &persona}, nil
	}
	return &PostResult{Post: &post}, nil
}
