package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/post/messages"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	post_repo "github.com/emmanuella-codes/nox/repositories/post"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

func (p *PostPipe) GetPostPipe(ctx context.Context, postID uuid.UUID) *shared.PipeRes[models.Post] {
	post, err := p.postRepo.FindPostByID(ctx, postID)
	if err != nil {
		if err == post_repo.ErrPostNotFound {
			return pipeError[models.Post](messages.Post_Not_Found)
		}
		return pipeInternalError[models.Post](err, "post.get")
	}

	return pipeSuccess(messages.Post_Fetched, post)
}

func (p *PostPipe) GetPersonaPostsPipe(ctx context.Context, personaID uuid.UUID, limit int) *shared.PipeRes[[]*models.Post] {
	if _, err := p.personaRepo.FindPersonaByID(ctx, personaID); err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return pipeError[[]*models.Post](messages.Persona_Not_Found)
		}
		return pipeInternalError[[]*models.Post](err, "post.find_persona")
	}

	posts, err := p.postRepo.FindPostsByPersonaID(ctx, personaID, limit)
	if err != nil {
		return pipeInternalError[[]*models.Post](err, "post.list_by_persona")
	}

	return pipeSuccess(messages.Posts_Listed, &posts)
}

func (p *PostPipe) GetFeedPipe(ctx context.Context, personaID uuid.UUID, limit int) *shared.PipeRes[[]*models.Post] {
	if _, err := p.personaRepo.FindPersonaByID(ctx, personaID); err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return pipeError[[]*models.Post](messages.Persona_Not_Found)
		}
		return pipeInternalError[[]*models.Post](err, "post.find_persona_for_feed")
	}

	posts, err := p.postRepo.FindFeedPosts(ctx, personaID, limit)
	if err != nil {
		return pipeInternalError[[]*models.Post](err, "post.feed")
	}

	return pipeSuccess(messages.Feed_Listed, &posts)
}
