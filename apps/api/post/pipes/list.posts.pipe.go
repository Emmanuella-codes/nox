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

func (p *PostPipe) GetPostPipe(ctx context.Context, postID uuid.UUID) *shared.PipeRes[PostResponse] {
	post, err := p.postRepo.FindPostByID(ctx, postID)
	if err != nil {
		if err == post_repo.ErrPostNotFound {
			return pipeError[PostResponse](messages.Post_Not_Found)
		}
		return pipeInternalError[PostResponse](err, "post.get")
	}

	persona, pipeErr := p.publicPostPersona(ctx, post)
	if pipeErr != nil {
		return pipeErr
	}

	response := postResponse(post, persona)
	return pipeSuccess(messages.Post_Fetched, &response)
}

func (p *PostPipe) GetPostForViewerPipe(ctx context.Context, postID uuid.UUID, viewerPersonaID uuid.UUID) *shared.PipeRes[PostResponse] {
	res := p.GetPostPipe(ctx, postID)
	if !res.Success || res.Data == nil || p.likeRepo == nil {
		return res
	}

	liked, err := p.likeRepo.HasPostLike(ctx, viewerPersonaID, postID)
	if err != nil {
		return pipeInternalError[PostResponse](err, "post.viewer_like_status")
	}
	res.Data.IsLiked = liked
	return res
}

func (p *PostPipe) GetPersonaPostsPipe(ctx context.Context, personaID uuid.UUID, limit int) *shared.PipeRes[[]PostResponse] {
	if _, err := p.personaRepo.FindPersonaByID(ctx, personaID); err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return pipeError[[]PostResponse](messages.Persona_Not_Found)
		}
		return pipeInternalError[[]PostResponse](err, "post.find_persona")
	}

	posts, err := p.postRepo.FindPostsByPersonaID(ctx, personaID, limit)
	if err != nil {
		return pipeInternalError[[]PostResponse](err, "post.list_by_persona")
	}

	personas, pipeErr := p.publicPostPersonas(ctx, posts)
	if pipeErr != nil {
		return pipeErr
	}

	responses := postResponses(posts, personas)
	return pipeSuccess(messages.Posts_Listed, &responses)
}

func (p *PostPipe) GetFeedPipe(ctx context.Context, personaID uuid.UUID, limit int) *shared.PipeRes[[]PostResponse] {
	if _, err := p.personaRepo.FindPersonaByID(ctx, personaID); err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return pipeError[[]PostResponse](messages.Persona_Not_Found)
		}
		return pipeInternalError[[]PostResponse](err, "post.find_persona_for_feed")
	}

	posts, err := p.postRepo.FindFeedPosts(ctx, personaID, limit)
	if err != nil {
		return pipeInternalError[[]PostResponse](err, "post.feed")
	}

	personas, pipeErr := p.publicPostPersonas(ctx, posts)
	if pipeErr != nil {
		return pipeErr
	}

	responses := postResponses(posts, personas)
	if p.likeRepo != nil {
		liked, err := p.likeRepo.FindLikedPostIDs(ctx, personaID, postIDs(posts))
		if err != nil {
			return pipeInternalError[[]PostResponse](err, "post.feed_like_status")
		}
		for i := range responses {
			responses[i].IsLiked = liked[posts[i].ID]
		}
	}
	return pipeSuccess(messages.Feed_Listed, &responses)
}

func (p *PostPipe) FindViewerPersona(ctx context.Context, userID uuid.UUID, personaID uuid.UUID) (*models.Persona, shared.PipeMessage) {
	persona, err := p.personaRepo.FindPersonaByID(ctx, personaID)
	if err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return nil, messages.Persona_Not_Found
		}
		return nil, messages.Internal_Error
	}
	if persona.UserID != userID || persona.PersonaType != models.VisiblePersonaType {
		return nil, messages.Forbidden
	}
	return persona, ""
}

func (p *PostPipe) publicPostPersona(ctx context.Context, post *models.Post) (*models.Persona, *shared.PipeRes[PostResponse]) {
	if post.PostingMode != models.PublicPostingMode || post.PersonaID == nil {
		return nil, nil
	}

	persona, err := p.personaRepo.FindPersonaByID(ctx, *post.PersonaID)
	if err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return nil, pipeError[PostResponse](messages.Persona_Not_Found)
		}
		return nil, pipeInternalError[PostResponse](err, "post.find_public_persona")
	}

	return persona, nil
}

func (p *PostPipe) publicPostPersonas(ctx context.Context, posts []*models.Post) (map[string]*models.Persona, *shared.PipeRes[[]PostResponse]) {
	personas := make(map[string]*models.Persona)
	for _, post := range posts {
		if post.PostingMode != models.PublicPostingMode || post.PersonaID == nil {
			continue
		}
		personaID := post.PersonaID.String()
		if _, ok := personas[personaID]; ok {
			continue
		}

		persona, err := p.personaRepo.FindPersonaByID(ctx, *post.PersonaID)
		if err != nil {
			if err == persona_repo.ErrPersonaNotFound {
				return nil, pipeError[[]PostResponse](messages.Persona_Not_Found)
			}
			return nil, pipeInternalError[[]PostResponse](err, "post.find_public_persona")
		}
		personas[personaID] = persona
	}

	return personas, nil
}
