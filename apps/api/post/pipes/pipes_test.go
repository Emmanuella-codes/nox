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

func TestCreatePostPipeSyncsExtractedHashtags(t *testing.T) {
	userID := uuid.New()
	hashtagRepo := &postTestHashtagRepo{}
	pipe := NewPostPipe(&postTestRepo{}, &postTestPersonaRepo{}, hashtagRepo)

	res := pipe.CreatePostPipe(context.Background(), userID, postdtos.CreatePostDTO{
		PostingMode: models.AnonymousPostingMode,
		Body:        "tonight at #Amapiano with #afro-house and #amapiano",
		PostType:    models.TextPostType,
	})
	if !res.Success {
		t.Fatalf("expected success, got %q", res.Message)
	}
	if len(hashtagRepo.syncedTags) != 2 {
		t.Fatalf("expected 2 synced tags, got %v", hashtagRepo.syncedTags)
	}
	if hashtagRepo.syncedTags[0] != "amapiano" || hashtagRepo.syncedTags[1] != "afro-house" {
		t.Fatalf("unexpected synced tags: %v", hashtagRepo.syncedTags)
	}
	if len(res.Data.Hashtags) != 2 {
		t.Fatalf("expected response hashtags, got %v", res.Data.Hashtags)
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
	hashtagRepo := &postTestHashtagRepo{}
	postRepo := &postTestRepo{
		posts: map[string]*models.Post{
			postID.String(): {
				ID:           postID,
				AuthorUserID: userID,
				PostingMode:  models.AnonymousPostingMode,
			},
		},
	}
	pipe := NewPostPipe(postRepo, &postTestPersonaRepo{}, hashtagRepo)

	res := pipe.DeletePostPipe(context.Background(), userID, postID)
	if !res.Success {
		t.Fatalf("expected delete success, got %q", res.Message)
	}
	if postRepo.deletedPostID != postID {
		t.Fatalf("expected deleted post %s, got %s", postID, postRepo.deletedPostID)
	}
	if hashtagRepo.deletedPostID != postID {
		t.Fatalf("expected hashtag cleanup for %s, got %s", postID, hashtagRepo.deletedPostID)
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

func TestGetPostPipeHydratesHashtags(t *testing.T) {
	postID := uuid.New()
	pipe := NewPostPipe(&postTestRepo{
		posts: map[string]*models.Post{
			postID.String(): {
				ID:          postID,
				PostingMode: models.AnonymousPostingMode,
				Body:        "#amapiano post",
				PostType:    models.TextPostType,
			},
		},
	}, &postTestPersonaRepo{}, &postTestHashtagRepo{
		tagsByPost: map[uuid.UUID][]string{postID: []string{"amapiano"}},
	})

	res := pipe.GetPostPipe(context.Background(), postID)
	if !res.Success {
		t.Fatalf("expected get success, got %q", res.Message)
	}
	if len(res.Data.Hashtags) != 1 || res.Data.Hashtags[0] != "amapiano" {
		t.Fatalf("expected hydrated hashtags, got %v", res.Data.Hashtags)
	}
}

func TestGetFeedPipeHydratesLikedState(t *testing.T) {
	viewerPersonaID := uuid.New()
	likedPostID := uuid.New()
	unlikedPostID := uuid.New()
	pipe := NewPostPipe(&postTestRepo{
		feedPosts: []*models.Post{
			{ID: likedPostID, PostingMode: models.AnonymousPostingMode, Body: "liked", PostType: models.TextPostType},
			{ID: unlikedPostID, PostingMode: models.AnonymousPostingMode, Body: "unliked", PostType: models.TextPostType},
		},
	}, &postTestPersonaRepo{
		personas: map[string]*models.Persona{
			viewerPersonaID.String(): {
				ID:          viewerPersonaID,
				UserID:      uuid.New(),
				PersonaType: models.VisiblePersonaType,
			},
		},
	}, &postTestLikeRepo{
		likedPostIDs: map[uuid.UUID]bool{likedPostID: true},
	})

	res := pipe.GetFeedPipe(context.Background(), viewerPersonaID, 20)
	if !res.Success {
		t.Fatalf("expected feed success, got %q", res.Message)
	}
	if len(*res.Data) != 2 {
		t.Fatalf("expected 2 posts, got %d", len(*res.Data))
	}
	if !(*res.Data)[0].IsLiked {
		t.Fatal("expected first post to be liked")
	}
	if (*res.Data)[1].IsLiked {
		t.Fatal("expected second post not to be liked")
	}
}

func TestGetPersonaPostsForViewerPipeHydratesLikedState(t *testing.T) {
	personaID := uuid.New()
	viewerPersonaID := uuid.New()
	likedPostID := uuid.New()
	unlikedPostID := uuid.New()
	pipe := NewPostPipe(&postTestRepo{
		personaPosts: []*models.Post{
			{ID: likedPostID, PersonaID: &personaID, PostingMode: models.PublicPostingMode, Body: "liked", PostType: models.TextPostType},
			{ID: unlikedPostID, PersonaID: &personaID, PostingMode: models.PublicPostingMode, Body: "unliked", PostType: models.TextPostType},
		},
	}, &postTestPersonaRepo{
		personas: map[string]*models.Persona{
			personaID.String(): {
				ID:          personaID,
				UserID:      uuid.New(),
				PersonaType: models.VisiblePersonaType,
			},
			viewerPersonaID.String(): {
				ID:          viewerPersonaID,
				UserID:      uuid.New(),
				PersonaType: models.VisiblePersonaType,
			},
		},
	}, &postTestLikeRepo{
		likedPostIDs: map[uuid.UUID]bool{likedPostID: true},
	})

	res := pipe.GetPersonaPostsForViewerPipe(context.Background(), personaID, viewerPersonaID, 20)
	if !res.Success {
		t.Fatalf("expected persona posts success, got %q", res.Message)
	}
	if len(*res.Data) != 2 {
		t.Fatalf("expected 2 posts, got %d", len(*res.Data))
	}
	if !(*res.Data)[0].IsLiked {
		t.Fatal("expected first post to be liked")
	}
	if (*res.Data)[1].IsLiked {
		t.Fatal("expected second post not to be liked")
	}
}

type postTestRepo struct {
	posts               map[string]*models.Post
	personaPosts        []*models.Post
	feedPosts           []*models.Post
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
	return r.personaPosts, nil
}

func (r *postTestRepo) FindFeedPosts(ctx context.Context, personaID uuid.UUID, limit int) ([]*models.Post, error) {
	return r.feedPosts, nil
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

type postTestLikeRepo struct {
	likedPostIDs map[uuid.UUID]bool
}

type postTestHashtagRepo struct {
	syncedPostID  uuid.UUID
	syncedTags    []string
	deletedPostID uuid.UUID
	tagsByPost    map[uuid.UUID][]string
}

func (r *postTestHashtagRepo) SyncPostHashtags(ctx context.Context, postID uuid.UUID, tags []string) error {
	r.syncedPostID = postID
	r.syncedTags = tags
	return nil
}

func (r *postTestHashtagRepo) DeletePostHashtags(ctx context.Context, postID uuid.UUID) error {
	r.deletedPostID = postID
	return nil
}

func (r *postTestHashtagRepo) FindTagsByPostIDs(ctx context.Context, postIDs []uuid.UUID) (map[uuid.UUID][]string, error) {
	tags := make(map[uuid.UUID][]string)
	for _, postID := range postIDs {
		if r.tagsByPost != nil {
			tags[postID] = r.tagsByPost[postID]
		}
	}
	return tags, nil
}

func (r *postTestHashtagRepo) FindTrending(ctx context.Context, limit int) ([]*models.Hashtag, error) {
	return nil, nil
}

func (r *postTestHashtagRepo) FindByTag(ctx context.Context, tag string) (*models.Hashtag, error) {
	return nil, nil
}

func (r *postTestHashtagRepo) FindPostsByTag(ctx context.Context, tag string, limit int) ([]*models.Post, error) {
	return nil, nil
}

func (r *postTestLikeRepo) LikePost(ctx context.Context, personaID uuid.UUID, postID uuid.UUID) error {
	return nil
}

func (r *postTestLikeRepo) UnlikePost(ctx context.Context, personaID uuid.UUID, postID uuid.UUID) error {
	return nil
}

func (r *postTestLikeRepo) HasPostLike(ctx context.Context, personaID uuid.UUID, postID uuid.UUID) (bool, error) {
	return r.likedPostIDs[postID], nil
}

func (r *postTestLikeRepo) FindLikedPostIDs(ctx context.Context, personaID uuid.UUID, postIDs []uuid.UUID) (map[uuid.UUID]bool, error) {
	liked := make(map[uuid.UUID]bool)
	for _, postID := range postIDs {
		if r.likedPostIDs[postID] {
			liked[postID] = true
		}
	}
	return liked, nil
}
