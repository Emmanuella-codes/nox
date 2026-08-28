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

// TestCreatePostPipeCreatesAnonymousPostWithoutPublicIdentity keeps anonymous responses detached from profiles.
func TestCreatePostPipeCreatesAnonymousPostWithoutPublicIdentity(t *testing.T) {
	userID := uuid.New()
	postRepo := &postTestRepo{}
	personaID := uuid.New()
	personaRepo := &postTestPersonaRepo{
		personas: map[string]*models.Persona{
			personaID.String(): {
				ID:          personaID,
				UserID:      userID,
				PersonaType: models.VisiblePersonaType,
			},
		},
	}
	pipe := NewPostPipe(postRepo, personaRepo)

	res := pipe.CreatePostPipe(context.Background(), userID, postdtos.CreatePostDTO{
		PersonaID:   &personaID,
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
	if postRepo.createdDTO.PersonaID == nil || *postRepo.createdDTO.PersonaID != personaID {
		t.Fatal("expected anonymous post persona id to be retained internally")
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
	if res.Data.Author.Anonymous == nil || res.Data.Author.Anonymous.Handle == "" {
		t.Fatal("expected anonymous response to expose thread alias")
	}
	if res.Data.Author.Anonymous.AvatarKey == "" {
		t.Fatal("expected anonymous response to expose fallback avatar key")
	}
}

// TestCreatePostPipeSyncsExtractedHashtags persists normalized tags from the post body.
func TestCreatePostPipeSyncsExtractedHashtags(t *testing.T) {
	userID := uuid.New()
	personaID := uuid.New()
	postRepo := &postTestRepo{}
	pipe := NewPostPipe(postRepo, &postTestPersonaRepo{
		personas: map[string]*models.Persona{
			personaID.String(): {
				ID:          personaID,
				UserID:      userID,
				PersonaType: models.VisiblePersonaType,
			},
		},
	}, &postTestHashtagRepo{})

	res := pipe.CreatePostPipe(context.Background(), userID, postdtos.CreatePostDTO{
		PersonaID:   &personaID,
		PostingMode: models.AnonymousPostingMode,
		Body:        "tonight at #Amapiano with #afro-house and #amapiano",
		PostType:    models.TextPostType,
	})
	if !res.Success {
		t.Fatalf("expected success, got %q", res.Message)
	}
	if len(postRepo.createdTags) != 2 {
		t.Fatalf("expected 2 synced tags, got %v", postRepo.createdTags)
	}
	if postRepo.createdTags[0] != "amapiano" || postRepo.createdTags[1] != "afro-house" {
		t.Fatalf("unexpected synced tags: %v", postRepo.createdTags)
	}
	if len(res.Data.Hashtags) != 2 {
		t.Fatalf("expected response hashtags, got %v", res.Data.Hashtags)
	}
}

// TestCreatePostPipeRequiresPersonaForPublicPost enforces the single-profile author input.
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

// TestCreatePostPipeRejectsNonOwnedPublicPersona blocks writing through another user's profile.
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

// TestDeletePostPipeUsesAuthorUserID authorizes deletion by owning user instead of persona type.
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
	pipe := NewPostPipe(postRepo, &postTestPersonaRepo{}, &postTestHashtagRepo{})

	res := pipe.DeletePostPipe(context.Background(), userID, postID)
	if !res.Success {
		t.Fatalf("expected delete success, got %q", res.Message)
	}
	if postRepo.deletedPostID != postID {
		t.Fatalf("expected deleted post %s, got %s", postID, postRepo.deletedPostID)
	}
	if !postRepo.deletedWithHashtags {
		t.Fatal("expected atomic hashtag cleanup delete path")
	}
}

// TestGetPostPipeHidesAnonymousIdentity keeps profile data off anonymous post responses.
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

// TestGetPostPipeHydratesHashtags attaches stored hashtags to a post response.
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
		tagsByPost: map[uuid.UUID][]string{postID: {"amapiano"}},
	})

	res := pipe.GetPostPipe(context.Background(), postID)
	if !res.Success {
		t.Fatalf("expected get success, got %q", res.Message)
	}
	if len(res.Data.Hashtags) != 1 || res.Data.Hashtags[0] != "amapiano" {
		t.Fatalf("expected hydrated hashtags, got %v", res.Data.Hashtags)
	}
}

// TestGetFeedPipeHydratesLikedState marks feed items liked for the viewer profile.
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

// TestGetPersonaPostsForViewerPipeHydratesLikedState attaches like state on profile post listings.
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

// TestGetPersonaPostsForViewerPipeIncludesAnonymousPostsForOwner shows owner-only anonymous posts on their profile.
func TestGetPersonaPostsForViewerPipeIncludesAnonymousPostsForOwner(t *testing.T) {
	userID := uuid.New()
	personaID := uuid.New()
	anonymousPostID := uuid.New()
	pipe := NewPostPipe(&postTestRepo{
		authorPosts: []*models.Post{
			{ID: anonymousPostID, AuthorUserID: userID, PersonaID: &personaID, PostingMode: models.AnonymousPostingMode, Body: "hidden", PostType: models.TextPostType},
		},
	}, &postTestPersonaRepo{
		personas: map[string]*models.Persona{
			personaID.String(): {ID: personaID, UserID: userID, PersonaType: models.VisiblePersonaType},
		},
	})

	res := pipe.GetPersonaPostsForViewerPipe(context.Background(), personaID, personaID, 20)
	if !res.Success {
		t.Fatalf("expected owner profile success, got %q", res.Message)
	}
	if len(*res.Data) != 1 {
		t.Fatalf("expected 1 owner post, got %d", len(*res.Data))
	}
	if (*res.Data)[0].Author.Persona != nil {
		t.Fatal("expected owner anonymous post response not to expose persona")
	}
	if (*res.Data)[0].Author.Anonymous == nil || (*res.Data)[0].Author.Anonymous.Handle == "" {
		t.Fatal("expected owner anonymous post response to retain anonymous identity")
	}
}

type postTestRepo struct {
	posts               map[string]*models.Post
	personaPosts        []*models.Post
	authorPosts         []*models.Post
	feedPosts           []*models.Post
	followingFeedPosts  []*models.Post
	identities          map[uuid.UUID]*models.AnonymousThreadIdentity
	createdAuthorUserID uuid.UUID
	createdDTO          postdtos.CreatePostDTO
	createdTags         []string
	deletedPostID       uuid.UUID
	deletedWithHashtags bool
}

// CreatePost stores the post creation inputs for assertions.
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

// CreatePostWithHashtags stores synced tags before delegating to CreatePost.
func (r *postTestRepo) CreatePostWithHashtags(ctx context.Context, authorUserID uuid.UUID, dto postdtos.CreatePostDTO, tags []string) (*models.Post, error) {
	r.createdTags = tags
	return r.CreatePost(ctx, authorUserID, dto)
}

// CreatePostWithHashtagsAndMedia stores synced tags before delegating to CreatePost.
func (r *postTestRepo) CreatePostWithHashtagsAndMedia(ctx context.Context, authorUserID uuid.UUID, dto postdtos.CreatePostDTO, tags []string, mediaAssetIDs []uuid.UUID) (*models.Post, error) {
	r.createdTags = tags
	return r.CreatePost(ctx, authorUserID, dto)
}

// FindPostByID returns a post by id from the in-memory fixture.
func (r *postTestRepo) FindPostByID(ctx context.Context, postID uuid.UUID) (*models.Post, error) {
	post, ok := r.posts[postID.String()]
	if !ok {
		return nil, post_repo.ErrPostNotFound
	}
	return post, nil
}

// FindPostsByPersonaID returns public profile posts from the fixture.
func (r *postTestRepo) FindPostsByPersonaID(ctx context.Context, personaID uuid.UUID, limit int) ([]*models.Post, error) {
	return r.personaPosts, nil
}

// FindPostsByAuthorUserID returns all owner posts from the fixture.
func (r *postTestRepo) FindPostsByAuthorUserID(ctx context.Context, authorUserID uuid.UUID, limit int) ([]*models.Post, error) {
	return r.authorPosts, nil
}

// FindFeedPosts returns the feed posts from the fixture.
func (r *postTestRepo) FindFeedPosts(ctx context.Context, personaID uuid.UUID, limit int) ([]*models.Post, error) {
	return r.feedPosts, nil
}

// FindFollowingFeedPosts returns the following feed posts from the fixture.
func (r *postTestRepo) FindFollowingFeedPosts(ctx context.Context, personaID uuid.UUID, limit int) ([]*models.Post, error) {
	return r.followingFeedPosts, nil
}

// EnsureAnonymousThreadIdentity stores or reuses one anonymous identity per thread owner.
func (r *postTestRepo) EnsureAnonymousThreadIdentity(ctx context.Context, threadID uuid.UUID, userID uuid.UUID, personaID uuid.UUID, anonymousHandle string, anonymousAvatarKey string) (*models.AnonymousThreadIdentity, error) {
	if r.identities == nil {
		r.identities = map[uuid.UUID]*models.AnonymousThreadIdentity{}
	}
	if identity, ok := r.identities[threadID]; ok {
		return identity, nil
	}
	identity := &models.AnonymousThreadIdentity{
		ID:              uuid.New(),
		ThreadID:        threadID,
		UserID:          userID,
		PersonaID:       personaID,
		AnonymousHandle: anonymousHandle,
		AnonymousAvatarKey: anonymousAvatarKey,
	}
	r.identities[threadID] = identity
	return identity, nil
}

// FindAnonymousThreadIdentities returns anonymous identities keyed by owner user id.
func (r *postTestRepo) FindAnonymousThreadIdentities(ctx context.Context, threadID uuid.UUID, userIDs []uuid.UUID) (map[uuid.UUID]*models.AnonymousThreadIdentity, error) {
	identities := map[uuid.UUID]*models.AnonymousThreadIdentity{}
	if r.identities == nil {
		return identities, nil
	}
	for _, identity := range r.identities {
		if identity.ThreadID != threadID {
			continue
		}
		for _, userID := range userIDs {
			if identity.UserID == userID {
				identities[userID] = identity
			}
		}
	}
	return identities, nil
}

// FindMediaAssetsByPostIDs returns an empty media map for fixture-backed tests.
func (r *postTestRepo) FindMediaAssetsByPostIDs(ctx context.Context, postIDs []uuid.UUID) (map[uuid.UUID][]*models.MediaAsset, error) {
	return map[uuid.UUID][]*models.MediaAsset{}, nil
}

// DeletePost records the deleted post id for assertions.
func (r *postTestRepo) DeletePost(ctx context.Context, postID uuid.UUID) error {
	r.deletedPostID = postID
	return nil
}

// DeletePostWithHashtags records the hashtag-aware delete path for assertions.
func (r *postTestRepo) DeletePostWithHashtags(ctx context.Context, postID uuid.UUID) error {
	r.deletedWithHashtags = true
	return r.DeletePost(ctx, postID)
}

type postTestPersonaRepo struct {
	personas map[string]*models.Persona
}

// CreatePersona is unused in these tests.
func (r *postTestPersonaRepo) CreatePersona(ctx context.Context, userID uuid.UUID, dto dtos.CreatePersonaDTO) (*models.Persona, error) {
	return nil, nil
}

// FindPersonaByID returns a profile fixture by id.
func (r *postTestPersonaRepo) FindPersonaByID(ctx context.Context, personaID uuid.UUID) (*models.Persona, error) {
	persona, ok := r.personas[personaID.String()]
	if !ok {
		return nil, persona_repo.ErrPersonaNotFound
	}
	return persona, nil
}

// FindPersonasByUserID is unused in these tests.
func (r *postTestPersonaRepo) FindPersonasByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Persona, error) {
	return nil, nil
}

// FindPersonaByHandle is unused in these tests.
func (r *postTestPersonaRepo) FindPersonaByHandle(ctx context.Context, handle string) (*models.Persona, error) {
	return nil, persona_repo.ErrPersonaNotFound
}

// UpdatePersona is unused in these tests.
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

// SyncPostHashtags records synced tags for assertions.
func (r *postTestHashtagRepo) SyncPostHashtags(ctx context.Context, postID uuid.UUID, tags []string) error {
	r.syncedPostID = postID
	r.syncedTags = tags
	return nil
}

// DeletePostHashtags records deleted hashtag relations for assertions.
func (r *postTestHashtagRepo) DeletePostHashtags(ctx context.Context, postID uuid.UUID) error {
	r.deletedPostID = postID
	return nil
}

// FindTagsByPostIDs returns the configured tag fixtures.
func (r *postTestHashtagRepo) FindTagsByPostIDs(ctx context.Context, postIDs []uuid.UUID) (map[uuid.UUID][]string, error) {
	tags := make(map[uuid.UUID][]string)
	for _, postID := range postIDs {
		if r.tagsByPost != nil {
			tags[postID] = r.tagsByPost[postID]
		}
	}
	return tags, nil
}

// FindTrending is unused in these tests.
func (r *postTestHashtagRepo) FindTrending(ctx context.Context, limit int) ([]*models.Hashtag, error) {
	return nil, nil
}

// FindByTag is unused in these tests.
func (r *postTestHashtagRepo) FindByTag(ctx context.Context, tag string) (*models.Hashtag, error) {
	return nil, nil
}

// FindPostsByTag is unused in these tests.
func (r *postTestHashtagRepo) FindPostsByTag(ctx context.Context, tag string, limit int, offset int) ([]*models.Post, error) {
	return nil, nil
}

// Search is unused in these tests.
func (r *postTestHashtagRepo) Search(ctx context.Context, query string, limit int, offset int) ([]*models.Hashtag, error) {
	return nil, nil
}

// LikePost is unused in these tests.
func (r *postTestLikeRepo) LikePost(ctx context.Context, personaID uuid.UUID, postID uuid.UUID) error {
	return nil
}

// UnlikePost is unused in these tests.
func (r *postTestLikeRepo) UnlikePost(ctx context.Context, personaID uuid.UUID, postID uuid.UUID) error {
	return nil
}

// HasPostLike returns fixture-backed like state for one post.
func (r *postTestLikeRepo) HasPostLike(ctx context.Context, personaID uuid.UUID, postID uuid.UUID) (bool, error) {
	return r.likedPostIDs[postID], nil
}

// FindLikedPostIDs returns fixture-backed like state for many posts.
func (r *postTestLikeRepo) FindLikedPostIDs(ctx context.Context, personaID uuid.UUID, postIDs []uuid.UUID) (map[uuid.UUID]bool, error) {
	liked := make(map[uuid.UUID]bool)
	for _, postID := range postIDs {
		if r.likedPostIDs[postID] {
			liked[postID] = true
		}
	}
	return liked, nil
}
