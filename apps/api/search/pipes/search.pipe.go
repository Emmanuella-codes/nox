package pipes

import (
	"context"
	"strings"

	"github.com/emmanuella-codes/nox/models"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	searchmessages "github.com/emmanuella-codes/nox/search/messages"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

func (p *SearchPipe) Search(ctx context.Context, query string, limit int) *shared.PipeRes[SearchResponse] {
	return p.search(ctx, query, limit, nil)
}

func (p *SearchPipe) SearchForViewer(ctx context.Context, query string, limit int, viewerPersonaID uuid.UUID) *shared.PipeRes[SearchResponse] {
	return p.search(ctx, query, limit, &viewerPersonaID)
}

func (p *SearchPipe) FindViewerPersona(ctx context.Context, userID uuid.UUID, personaID uuid.UUID) (*models.Persona, shared.PipeMessage) {
	if p.personaRepo == nil {
		return nil, searchmessages.Internal_Error
	}

	persona, err := p.personaRepo.FindPersonaByID(ctx, personaID)
	if err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return nil, searchmessages.Persona_Not_Found
		}
		return nil, searchmessages.Internal_Error
	}
	if persona.UserID != userID || persona.PersonaType != models.VisiblePersonaType {
		return nil, searchmessages.Forbidden
	}
	return persona, ""
}

func (p *SearchPipe) search(ctx context.Context, query string, limit int, viewerPersonaID *uuid.UUID) *shared.PipeRes[SearchResponse] {
	query = strings.TrimSpace(query)
	if len(query) < 2 || len(query) > 80 {
		return pipeError[SearchResponse](searchmessages.Invalid_Query)
	}

	results, err := p.repo.Search(ctx, query, limit)
	if err != nil {
		return pipeInternalError[SearchResponse](err, "search.query")
	}

	response := SearchResponse{
		Query:    query,
		Personas: personaResponses(results.Personas),
		Posts:    postResponses(results.Posts),
		Events:   eventResponses(results.Events),
	}
	if viewerPersonaID != nil && p.likeRepo != nil {
		if err := p.hydrateLikedState(ctx, *viewerPersonaID, response.Posts); err != nil {
			return pipeInternalError[SearchResponse](err, "search.like_status")
		}
	}
	return pipeSuccess(searchmessages.Search_Listed, &response)
}

func (p *SearchPipe) hydrateLikedState(ctx context.Context, viewerPersonaID uuid.UUID, posts []SearchPostResponse) error {
	if len(posts) == 0 {
		return nil
	}

	postIDs := make([]uuid.UUID, 0, len(posts))
	for _, post := range posts {
		postID, err := uuid.Parse(post.ID)
		if err != nil {
			return err
		}
		postIDs = append(postIDs, postID)
	}

	liked, err := p.likeRepo.FindLikedPostIDs(ctx, viewerPersonaID, postIDs)
	if err != nil {
		return err
	}
	for i := range posts {
		posts[i].IsLiked = liked[postIDs[i]]
	}
	return nil
}
