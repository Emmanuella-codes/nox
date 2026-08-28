package pipes

import (
	"fmt"
	"time"

	"github.com/emmanuella-codes/nox/comment/messages"
	"github.com/emmanuella-codes/nox/models"
	comment_repo "github.com/emmanuella-codes/nox/repositories/comment"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	post_repo "github.com/emmanuella-codes/nox/repositories/post"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

type CommentPipe struct {
	commentRepo comment_repo.CommentRepository
	personaRepo persona_repo.PersonaRepository
	postRepo    post_repo.PostRepository
}

type CommentResponse struct {
	ID        string        `json:"id"`
	PostID    string        `json:"post_id"`
	Author    CommentAuthor `json:"author"`
	Body      string        `json:"body"`
	ParentID  *string       `json:"parent_id,omitempty"`
	LikeCount int           `json:"like_count"`
	CreatedAt time.Time     `json:"created_at"`
}

type CommentAuthor struct {
	Mode      models.PostingMode      `json:"mode"`
	Persona   *CommentPersonaAuthor   `json:"persona,omitempty"`
	Anonymous *CommentAnonymousAuthor `json:"anonymous,omitempty"`
}

type CommentPersonaAuthor struct {
	ID          string `json:"id"`
	Handle      string `json:"handle"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
}

type CommentAnonymousAuthor struct {
	Handle    string `json:"handle"`
	AvatarKey string `json:"avatar_key"`
}

// NewCommentPipe builds the comment orchestration layer from repositories.
func NewCommentPipe(commentRepo comment_repo.CommentRepository, personaRepo persona_repo.PersonaRepository, postRepo post_repo.PostRepository) *CommentPipe {
	return &CommentPipe{commentRepo: commentRepo, personaRepo: personaRepo, postRepo: postRepo}
}

// pipeInternalError maps internal comment errors to pipe responses.
func pipeInternalError[T any](err error, operation string) *shared.PipeRes[T] {
	return shared.PipeInternalError[T](err, "comment", operation, messages.Internal_Error)
}

// validPostingMode validates the supported comment publishing modes.
func validPostingMode(mode models.PostingMode) bool {
	return mode == models.PublicPostingMode || mode == models.AnonymousPostingMode
}

// commentResponse maps one comment model into the API response shape.
func commentResponse(comment *models.Comment, persona *models.Persona, identity *models.AnonymousThreadIdentity) CommentResponse {
	response := CommentResponse{
		ID:        comment.ID.String(),
		PostID:    comment.PostID.String(),
		Author:    CommentAuthor{Mode: comment.PostingMode},
		Body:      comment.Body,
		LikeCount: comment.LikeCount,
		CreatedAt: comment.CreatedAt,
	}
	if comment.ParentID != nil {
		parentID := comment.ParentID.String()
		response.ParentID = &parentID
	}
	if comment.PostingMode == models.PublicPostingMode && persona != nil {
		response.Author.Persona = &CommentPersonaAuthor{
			ID:          persona.ID.String(),
			Handle:      persona.Handle,
			DisplayName: persona.DisplayName,
			AvatarURL:   persona.AvatarURL,
		}
	}
	if comment.PostingMode == models.AnonymousPostingMode {
		response.Author.Anonymous = &CommentAnonymousAuthor{
			Handle:    anonymousHandleValue(identity),
			AvatarKey: anonymousAvatarValue(identity),
		}
	}
	return response
}

// anonymousHandle returns a generated anonymous handle for legacy tests and fallbacks.
func anonymousHandle() string {
	id := uuid.NewString()
	return fmt.Sprintf("ghost_%s", id[:8])
}
