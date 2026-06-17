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
		responses = append(responses, commentResponse(comment, personas[comment.PersonaID], identities[comment.PersonaID]))
	}
	return shared.PipeSuccess(messages.Comments_Listed, &responses)
}

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

func (p *CommentPipe) commentAnonymousIdentities(ctx context.Context, threadID uuid.UUID, comments []*models.Comment) (map[uuid.UUID]*models.AnonymousThreadIdentity, error) {
	seen := map[uuid.UUID]bool{}
	var personaIDs []uuid.UUID
	for _, comment := range comments {
		if comment.PostingMode != models.AnonymousPostingMode || seen[comment.PersonaID] {
			continue
		}
		seen[comment.PersonaID] = true
		personaIDs = append(personaIDs, comment.PersonaID)
	}
	return p.postRepo.FindAnonymousThreadIdentities(ctx, threadID, personaIDs)
}
