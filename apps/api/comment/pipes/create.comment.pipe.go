package pipes

import (
	"context"
	"strings"

	"github.com/emmanuella-codes/nox/comment/dtos"
	"github.com/emmanuella-codes/nox/comment/messages"
	"github.com/emmanuella-codes/nox/models"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	post_repo "github.com/emmanuella-codes/nox/repositories/post"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/emmanuella-codes/nox/shared/anonymousidentity"
	"github.com/google/uuid"
)

// CreateCommentPipe validates the author profile and persists a new comment.
func (p *CommentPipe) CreateCommentPipe(ctx context.Context, userID uuid.UUID, postID uuid.UUID, dto dtos.CreateCommentDTO) *shared.PipeRes[CommentResponse] {
	dto.Body = strings.TrimSpace(dto.Body)
	dto.AuthorUserID = userID
	if dto.PostingMode == "" {
		dto.PostingMode = models.PublicPostingMode
	}
	if !validPostingMode(dto.PostingMode) {
		return shared.PipeError[CommentResponse](messages.Invalid_Payload)
	}
	if _, err := p.postRepo.FindPostByID(ctx, postID); err != nil {
		if err == post_repo.ErrPostNotFound {
			return shared.PipeError[CommentResponse](messages.Post_Not_Found)
		}
		return pipeInternalError[CommentResponse](err, "comment.find_post")
	}
	profile, err := p.personaRepo.FindPersonaByID(ctx, dto.PersonaID)
	if err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return shared.PipeError[CommentResponse](messages.Persona_Not_Found)
		}
		return pipeInternalError[CommentResponse](err, "comment.find_persona")
	}
	if !profile.IsOwnedBy(userID) {
		return shared.PipeError[CommentResponse](messages.Forbidden)
	}
	comment, err := p.commentRepo.CreateComment(ctx, postID, dto)
	if err != nil {
		return pipeInternalError[CommentResponse](err, "comment.create")
	}
	var publicProfile *models.Persona
	var identity *models.AnonymousThreadIdentity
	if comment.PostingMode == models.PublicPostingMode {
		publicProfile = profile
	} else {
		identity, err = p.postRepo.EnsureAnonymousThreadIdentity(
			ctx,
			postID,
			userID,
			dto.PersonaID,
			anonymousidentity.GenerateHandle(),
			anonymousidentity.GenerateAvatarKey(),
		)
		if err != nil {
			return pipeInternalError[CommentResponse](err, "comment.anonymous_identity")
		}
	}
	response := commentResponse(comment, publicProfile, identity)
	return shared.PipeSuccess(messages.Comment_Created, &response)
}
