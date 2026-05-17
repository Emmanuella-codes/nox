package pipes

import (
	"context"
	"testing"

	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/persona/dtos"
	postdtos "github.com/emmanuella-codes/nox/post/dtos"
	"github.com/emmanuella-codes/nox/post/messages"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	post_repo "github.com/emmanuella-codes/nox/repositories/post"
	"github.com/google/uuid"
)

func TestCreatePostPipeCreatesAnonymousPostWithoutPublicIdentity(t *testing.T) {
	userID := uuid.New()
	postRepo := &postTestRepo{}
	personaRepo := &postTestPersonaRepo{}
	pipe := NewPostPipe(postRepo, personaRepo)
	leakyPersonaID := uuid.New()

	res := pipe.CreatePostPipe(context.Background(), userID, postdtos.CreatePostDTO{
		PersonaID:   &leakyPersonaID,
		PostingMode: models.AnonymousPostingMode,
		Body:        "anonymous body",
		PostType:    models.TextPostType,
	})
	if !res.Success {
		t.Fatalf("expected success, got %q", res.Message)
	}
	if postRepo.createdAuthorUserID != userID {
		t.Fatalf("expected author user id %s, got %s", userID, postRepo.createdAuthorUserID)
	}
	if postRepo.createdDTO.PersonaID != nil {
		t.Fatal("expected anonymous post persona id to be cleared before persistence")
	}
	if res.Data == nil {
		t.Fatal("expected post response")
	}
	if res.Data.Author.Mode != models.AnonymousPostingMode {
		t.Fatalf("expected anonymous author mode, got %q", res.Data.Author.Mode)
	}
	if res.Data.Author.Persona != nil {
		t.Fatal("expected anonymous response not to expose persona")
	}
}

func TestCreatePostPipeRequiresPersonaForPublicPost(t *testing.T) {
	pipe := NewPostPipe(&postTestRepo{}, &postTestPersonaRepo{})

	res := pipe.CreatePostPipe(context.Background(), uuid.New(), postdtos.CreatePostDTO{
		PostingMode: models.PublicPostingMode,
		Body:        "public body",
		PostType:    models.TextPostType,
	})
	if res.Message != messages.Persona_Required {
		t.Fatalf("expected %q, got %q", messages.Persona_Required, res.Message)
	}
}

func TestCreatePostPipeRejectsNonOwnedPublicPersona(t *testing.T) {
	personaID := uuid.New()
	pipe := NewPostPipe(&postTestRepo{}, &postTestPersonaRepo{
		personas: map[string]*models.Persona{
			personaID.String(): {
				ID:          personaID,
				UserID:      uuid.New(),
				PersonaType: models.VisiblePersonaType,
			},
		},
	})

	res := pipe.CreatePostPipe(context.Background(), uuid.New(), postdtos.CreatePostDTO{
		PersonaID:   &personaID,
		PostingMode: models.PublicPostingMode,
		Body:        "public body",
		PostType:    models.TextPostType,
	})
	if res.Message != messages.Forbidden {
		t.Fatalf("expected %q, got %q", messages.Forbidden, res.Message)
	}
}

func TestDeletePostPipeUsesAuthorUserID(t *testing.T) {
	userID := uuid.New()
	postID := uuid.New()
	postRepo := &postTestRepo{
		posts: map[string]*models.Post{
			postID.String(): {
				ID:           postID,
				AuthorUserID: userID,
				PostingMode:  models.AnonymousPostingMode,
			},
		},
	}
	pipe := NewPostPipe(postRepo, &postTestPersonaRepo{})

	res := pipe.DeletePostPipe(context.Background(), userID, postID)
	if !res.Success {
		t.Fatalf("expected delete success, got %q", res.Message)
	}
	if postRepo.deletedPostID != postID {
		t.Fatalf("expected deleted post %s, got %s", postID, postRepo.deletedPostID)
	}
}

func TestGetPostPipeHidesAnonymousIdentity(t *testing.T) {
	postID := uuid.New()
	userID := uuid.New()
	pipe := NewPostPipe(&postTestRepo{
		posts: map[string]*models.Post{
			postID.String(): {
				ID:           postID,
				AuthorUserID: userID,
				PostingMode:  models.AnonymousPostingMode,
				Body:         "anonymous body",
				PostType:     models.TextPostType,
			},
		},
	}, &postTestPersonaRepo{})

	res := pipe.GetPostPipe(context.Background(), postID)
	if !res.Success {
		t.Fatalf("expected get success, got %q", res.Message)
	}
	if res.Data.Author.Mode != models.AnonymousPostingMode {
		t.Fatalf("expected anonymous author mode, got %q", res.Data.Author.Mode)
	}
	if res.Data.Author.Persona != nil {
		t.Fatal("expected anonymous response not to expose persona")
	}
}

type postTestRepo struct {
	posts               map[string]*models.Post
	createdAuthorUserID uuid.UUID
	createdDTO          postdtos.CreatePostDTO
	deletedPostID       uuid.UUID
}

func (r *postTestRepo) CreatePost(ctx context.Context, authorUserID uuid.UUID, dto postdtos.CreatePostDTO) (*models.Post, error) {
	r.createdAuthorUserID = authorUserID
	r.createdDTO = dto

	return &models.Post{
		ID:           uuid.New(),
		AuthorUserID: authorUserID,
		PersonaID:    dto.PersonaID,
		PostingMode:  dto.PostingMode,
		Body:         dto.Body,
		PostType:     dto.PostType,
		MediaURL:     dto.MediaURL,
		MediaType:    dto.MediaType,
		Location:     dto.Location,
	}, nil
}

func (r *postTestRepo) FindPostByID(ctx context.Context, postID uuid.UUID) (*models.Post, error) {
	post, ok := r.posts[postID.String()]
	if !ok {
		return nil, post_repo.ErrPostNotFound
	}
	return post, nil
}

func (r *postTestRepo) FindPostsByPersonaID(ctx context.Context, personaID uuid.UUID, limit int) ([]*models.Post, error) {
	return nil, nil
}

func (r *postTestRepo) FindFeedPosts(ctx context.Context, personaID uuid.UUID, limit int) ([]*models.Post, error) {
	return nil, nil
}

func (r *postTestRepo) DeletePost(ctx context.Context, postID uuid.UUID) error {
	r.deletedPostID = postID
	return nil
}

type postTestPersonaRepo struct {
	personas map[string]*models.Persona
}

func (r *postTestPersonaRepo) CreatePersona(ctx context.Context, userID uuid.UUID, dto dtos.CreatePersonaDTO) (*models.Persona, error) {
	return nil, nil
}

func (r *postTestPersonaRepo) FindPersonaByID(ctx context.Context, personaID uuid.UUID) (*models.Persona, error) {
	persona, ok := r.personas[personaID.String()]
	if !ok {
		return nil, persona_repo.ErrPersonaNotFound
	}
	return persona, nil
}

func (r *postTestPersonaRepo) FindPersonasByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Persona, error) {
	return nil, nil
}

func (r *postTestPersonaRepo) FindPersonaByHandle(ctx context.Context, handle string) (*models.Persona, error) {
	return nil, persona_repo.ErrPersonaNotFound
}

func (r *postTestPersonaRepo) UpdatePersona(ctx context.Context, personaID uuid.UUID, dto dtos.UpdatePersonaDTO) (*models.Persona, error) {
	return nil, nil
}
