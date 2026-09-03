package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/comment/messages"
	"github.com/emmanuella-codes/nox/models"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	post_repo "github.com/emmanuella-codes/nox/repositories/post"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

// ListCommentsPipe fetches comments for one post and hydrates public or anonymous author state.
func (p *CommentPipe) ListCommentsPipe(ctx context.Context, postID uuid.UUID, limit int) *shared.PipeRes[[]CommentResponse] {
	if _, err := p.postRepo.FindPostByID(ctx, postID); err != nil {
		if err == post_repo.ErrPostNotFound {
			return shared.PipeError[[]CommentResponse](messages.Post_Not_Found)
		}
		return pipeInternalError[[]CommentResponse](err, "comment.find_post")
	}
	comments, err := p.commentRepo.FindCommentsByPostID(ctx, postID, limit)
	if err != nil {
		return pipeInternalError[[]CommentResponse](err, "comment.list")
	}
	personas, err := p.commentPersonas(ctx, comments)
	if err != nil {
		return pipeInternalError[[]CommentResponse](err, "comment.personas")
	}
	identities, err := p.commentAnonymousIdentities(ctx, postID, comments)
	if err != nil {
		return pipeInternalError[[]CommentResponse](err, "comment.anonymous_identities")
	}
	responses := make([]CommentResponse, 0, len(comments))
	for _, comment := range comments {
		responses = append(responses, commentResponse(comment, personas[comment.PersonaID], identities[comment.AuthorUserID]))
	}
	return shared.PipeSuccess(messages.Comments_Listed, &responses)
}

// commentPersonas fetches the public profiles needed by public comments.
func (p *CommentPipe) commentPersonas(ctx context.Context, comments []*models.Comment) (map[uuid.UUID]*models.Persona, error) {
	personas := make(map[uuid.UUID]*models.Persona)
	for _, comment := range comments {
		if comment.PostingMode != models.PublicPostingMode {
			continue
		}
		if _, ok := personas[comment.PersonaID]; ok {
			continue
		}
		persona, err := p.personaRepo.FindPersonaByID(ctx, comment.PersonaID)
		if err != nil {
			if err == persona_repo.ErrPersonaNotFound {
				continue
			}
			return nil, err
		}
		personas[comment.PersonaID] = persona
	}
	return personas, nil
}

// commentAnonymousIdentities fetches thread-scoped anonymous identities keyed by author user id.
func (p *CommentPipe) commentAnonymousIdentities(ctx context.Context, threadID uuid.UUID, comments []*models.Comment) (map[uuid.UUID]*models.AnonymousThreadIdentity, error) {
	seen := map[uuid.UUID]bool{}
	var userIDs []uuid.UUID
	for _, comment := range comments {
		if comment.PostingMode != models.AnonymousPostingMode || seen[comment.AuthorUserID] {
			continue
		}
		seen[comment.AuthorUserID] = true
		userIDs = append(userIDs, comment.AuthorUserID)
	}
	return p.postRepo.FindAnonymousThreadIdentities(ctx, threadID, userIDs)
}

// anonymousHandleValue extracts the anonymous handle or returns the default label.
func anonymousHandleValue(identity *models.AnonymousThreadIdentity) string {
	if identity == nil || identity.AnonymousHandle == "" {
		return "anonymous"
	}
	return identity.AnonymousHandle
}

// anonymousAvatarValue extracts the anonymous avatar key or returns the default key.
func anonymousAvatarValue(identity *models.AnonymousThreadIdentity) string {
	if identity == nil || identity.AnonymousAvatarKey == "" {
		return "mask_01"
	}
	return identity.AnonymousAvatarKey
}
