package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/hashtag/messages"
	"github.com/emmanuella-codes/nox/models"
	post_pipes "github.com/emmanuella-codes/nox/post/pipes"
	hashtag_repo "github.com/emmanuella-codes/nox/repositories/hashtag"
	like_repo "github.com/emmanuella-codes/nox/repositories/like"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	post_repo "github.com/emmanuella-codes/nox/repositories/post"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

type HashtagPipe struct {
	repo        hashtag_repo.HashtagRepository
	personaRepo persona_repo.PersonaRepository
	likeRepo    like_repo.LikeRepository
	postRepo    post_repo.PostRepository
}

func NewHashtagPipe(repo hashtag_repo.HashtagRepository, deps ...any) *HashtagPipe {
	var personas persona_repo.PersonaRepository
	var likes like_repo.LikeRepository
	var posts post_repo.PostRepository
	for _, dep := range deps {
		switch typed := dep.(type) {
		case persona_repo.PersonaRepository:
			personas = typed
		case like_repo.LikeRepository:
			likes = typed
		case post_repo.PostRepository:
			posts = typed
		}
	}
	return &HashtagPipe{repo: repo, personaRepo: personas, likeRepo: likes, postRepo: posts}
}

func pipeInternalError[T any](err error, operation string) *shared.PipeRes[T] {
	return shared.PipeInternalError[T](err, "hashtag", operation, messages.Internal_Error)
}

func (p *HashtagPipe) postResponses(ctx context.Context, posts []*models.Post) ([]post_pipes.PostResponse, error) {
	personas := make(map[string]*models.Persona)
	if p.personaRepo != nil {
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
				return nil, err
			}
			personas[personaID] = persona
		}
	}

	tagsByPost, err := p.repo.FindTagsByPostIDs(ctx, postIDs(posts))
	if err != nil {
		return nil, err
	}

	mediaByPost := map[uuid.UUID][]*models.MediaAsset{}
	if p.postRepo != nil {
		mediaByPost, err = p.postRepo.FindMediaAssetsByPostIDs(ctx, postIDs(posts))
		if err != nil {
			return nil, err
		}
	}

	responses := make([]post_pipes.PostResponse, 0, len(posts))
	for _, post := range posts {
		identity, err := p.anonymousPostIdentity(ctx, post)
		if err != nil {
			return nil, err
		}
		response := postResponse(post, personas, identity)
		response.Hashtags = tagsByPost[post.ID]
		response.Media = mediaByPost[post.ID]
		if response.Media == nil {
			response.Media = []*models.MediaAsset{}
		}
		responses = append(responses, response)
	}
	return responses, nil
}

func (p *HashtagPipe) hydrateLikedState(ctx context.Context, viewerPersonaID uuid.UUID, responses []post_pipes.PostResponse) error {
	if p.likeRepo == nil || len(responses) == 0 {
		return nil
	}

	postIDs := make([]uuid.UUID, 0, len(responses))
	for _, response := range responses {
		postID, err := uuid.Parse(response.ID)
		if err != nil {
			return err
		}
		postIDs = append(postIDs, postID)
	}

	liked, err := p.likeRepo.FindLikedPostIDs(ctx, viewerPersonaID, postIDs)
	if err != nil {
		return err
	}
	for i := range responses {
		responses[i].IsLiked = liked[postIDs[i]]
	}
	return nil
}

func (p *HashtagPipe) FindViewerPersona(ctx context.Context, userID uuid.UUID, personaID uuid.UUID) (*models.Persona, shared.PipeMessage) {
	if p.personaRepo == nil {
		return nil, messages.Internal_Error
	}

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

func (p *HashtagPipe) anonymousPostIdentity(ctx context.Context, post *models.Post) (*models.AnonymousThreadIdentity, error) {
	if p.postRepo == nil || post.PostingMode != models.AnonymousPostingMode || post.PersonaID == nil {
		return nil, nil
	}
	return p.postRepo.EnsureAnonymousThreadIdentity(ctx, post.ID, post.AuthorUserID, *post.PersonaID, anonymousHandle())
}

func postResponse(post *models.Post, personas map[string]*models.Persona, identity *models.AnonymousThreadIdentity) post_pipes.PostResponse {
	response := post_pipes.PostResponse{
		ID:           post.ID.String(),
		Author:       post_pipes.PostAuthor{Mode: post.PostingMode},
		Body:         post.Body,
		PostType:     post.PostType,
		MediaURL:     post.MediaURL,
		MediaType:    post.MediaType,
		Location:     post.Location,
		LikeCount:    post.LikeCount,
		CommentCount: post.CommentCount,
		RepostCount:  post.RepostCount,
		IsRepost:     post.IsRepost,
		Hashtags:     []string{},
		Media:        []*models.MediaAsset{},
		CreatedAt:    post.CreatedAt,
	}
	if post.EventID != nil {
		eventID := post.EventID.String()
		response.EventID = &eventID
	}
	if post.RepostOf != nil {
		repostOf := post.RepostOf.String()
		response.RepostOf = &repostOf
	}
	if post.PostingMode == models.PublicPostingMode && post.PersonaID != nil {
		if persona := personas[post.PersonaID.String()]; persona != nil {
			response.Author.Persona = &post_pipes.PostPersonaAuthor{
				ID:          persona.ID.String(),
				Handle:      persona.Handle,
				DisplayName: persona.DisplayName,
				AvatarURL:   persona.AvatarURL,
			}
		}
	}
	if post.PostingMode == models.AnonymousPostingMode {
		handle := "anonymous"
		if identity != nil && identity.AnonymousHandle != "" {
			handle = identity.AnonymousHandle
		}
		response.Author.Anonymous = &post_pipes.PostAnonymousAuthor{Handle: handle}
	}
	return response
}

func anonymousHandle() string {
	id := uuid.NewString()
	return "ghost_" + id[:8]
}

func postIDs(posts []*models.Post) []uuid.UUID {
	ids := make([]uuid.UUID, 0, len(posts))
	for _, post := range posts {
		ids = append(ids, post.ID)
	}
	return ids
}
