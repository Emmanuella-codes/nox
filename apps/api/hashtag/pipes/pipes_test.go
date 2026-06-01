package pipes

import (
	"context"
	"testing"
	"time"

	"github.com/emmanuella-codes/nox/hashtag/messages"
	"github.com/emmanuella-codes/nox/models"
	personadtos "github.com/emmanuella-codes/nox/persona/dtos"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	"github.com/google/uuid"
)

func TestTrendingPipeReturnsHashtags(t *testing.T) {
	repo := &hashtagTestRepo{
		trending: []*models.Hashtag{{ID: uuid.New(), Tag: "amapiano", PostCount: 3, CreatedAt: time.Now()}},
	}
	pipe := NewHashtagPipe(repo)

	res := pipe.TrendingPipe(context.Background(), 10)
	if !res.Success {
		t.Fatalf("expected success, got %q", res.Message)
	}
	if len(*res.Data) != 1 || (*res.Data)[0].Tag != "amapiano" {
		t.Fatalf("expected trending hashtag, got %v", res.Data)
	}
}

func TestGetHashtagPipeNormalizesTag(t *testing.T) {
	repo := &hashtagTestRepo{
		byTag: map[string]*models.Hashtag{
			"amapiano": {ID: uuid.New(), Tag: "amapiano", PostCount: 2},
		},
	}
	pipe := NewHashtagPipe(repo)

	res := pipe.GetHashtagPipe(context.Background(), "#Amapiano")
	if !res.Success {
		t.Fatalf("expected success, got %q", res.Message)
	}
	if res.Data.Tag != "amapiano" || res.Data.PostCount != 2 {
		t.Fatalf("unexpected hashtag detail: %+v", res.Data)
	}
}

func TestGetHashtagPipeRejectsInvalidTag(t *testing.T) {
	pipe := NewHashtagPipe(&hashtagTestRepo{})

	res := pipe.GetHashtagPipe(context.Background(), "---")
	if res.Message != messages.Invalid_Tag {
		t.Fatalf("expected %q, got %q", messages.Invalid_Tag, res.Message)
	}
}

func TestGetHashtagPipeReturnsNotFoundForUnknownTag(t *testing.T) {
	pipe := NewHashtagPipe(&hashtagTestRepo{})

	res := pipe.GetHashtagPipe(context.Background(), "missing")
	if res.Message != messages.Hashtag_Not_Found {
		t.Fatalf("expected %q, got %q", messages.Hashtag_Not_Found, res.Message)
	}
}

func TestPostsByTagPipeReturnsPostResponses(t *testing.T) {
	personaID := uuid.New()
	postID := uuid.New()
	repo := &hashtagTestRepo{
		postsByTag: map[string][]*models.Post{
			"amapiano": {
				{
					ID:          postID,
					PersonaID:   &personaID,
					PostingMode: models.PublicPostingMode,
					Body:        "#amapiano tonight",
					PostType:    models.TextPostType,
					CreatedAt:   time.Now(),
				},
			},
		},
		tagsByPost: map[uuid.UUID][]string{postID: []string{"amapiano"}},
	}
	personas := &hashtagTestPersonaRepo{
		personas: map[string]*models.Persona{
			personaID.String(): {
				ID:          personaID,
				Handle:      "ada",
				DisplayName: "Ada",
				PersonaType: models.VisiblePersonaType,
			},
		},
	}
	pipe := NewHashtagPipe(repo, personas)

	res := pipe.PostsByTagPipe(context.Background(), "#Amapiano", 10, 0)
	if !res.Success {
		t.Fatalf("expected success, got %q", res.Message)
	}
	if res.Data.Tag != "amapiano" || len(res.Data.Posts) != 1 {
		t.Fatalf("unexpected posts response: %+v", res.Data)
	}
	post := res.Data.Posts[0]
	if post.Author.Persona == nil || post.Author.Persona.Handle != "ada" {
		t.Fatalf("expected public persona author, got %+v", post.Author)
	}
	if len(post.Hashtags) != 1 || post.Hashtags[0] != "amapiano" {
		t.Fatalf("expected post hashtags, got %v", post.Hashtags)
	}
}

func TestPostsByTagPipeReturnsPaginationMetadata(t *testing.T) {
	posts := []*models.Post{
		{ID: uuid.New(), PostingMode: models.AnonymousPostingMode, Body: "one", PostType: models.TextPostType},
		{ID: uuid.New(), PostingMode: models.AnonymousPostingMode, Body: "two", PostType: models.TextPostType},
	}
	repo := &hashtagTestRepo{
		postsByTag: map[string][]*models.Post{"amapiano": posts},
	}
	pipe := NewHashtagPipe(repo)

	res := pipe.PostsByTagPipe(context.Background(), "amapiano", 1, 4)
	if !res.Success {
		t.Fatalf("expected success, got %q", res.Message)
	}
	if len(res.Data.Posts) != 1 {
		t.Fatalf("expected one returned post, got %d", len(res.Data.Posts))
	}
	if res.Data.Limit != 1 || res.Data.Offset != 4 || !res.Data.HasMore {
		t.Fatalf("unexpected pagination metadata: %+v", res.Data)
	}
	if res.Data.NextOffset == nil || *res.Data.NextOffset != 5 {
		t.Fatalf("expected next offset 5, got %v", res.Data.NextOffset)
	}
}

func TestPostsByTagForViewerPipeHydratesLikedState(t *testing.T) {
	viewerPersonaID := uuid.New()
	likedPostID := uuid.New()
	unlikedPostID := uuid.New()
	repo := &hashtagTestRepo{
		postsByTag: map[string][]*models.Post{
			"amapiano": {
				{ID: likedPostID, PostingMode: models.AnonymousPostingMode, Body: "liked", PostType: models.TextPostType},
				{ID: unlikedPostID, PostingMode: models.AnonymousPostingMode, Body: "open", PostType: models.TextPostType},
			},
		},
	}
	likes := &hashtagTestLikeRepo{likedPostIDs: map[uuid.UUID]bool{likedPostID: true}}
	pipe := NewHashtagPipe(repo, likes)

	res := pipe.PostsByTagForViewerPipe(context.Background(), "amapiano", 10, 0, viewerPersonaID)
	if !res.Success {
		t.Fatalf("expected success, got %q", res.Message)
	}
	if !res.Data.Posts[0].IsLiked {
		t.Fatal("expected first post to be liked")
	}
	if res.Data.Posts[1].IsLiked {
		t.Fatal("expected second post not to be liked")
	}
}

type hashtagTestRepo struct {
	trending   []*models.Hashtag
	byTag      map[string]*models.Hashtag
	postsByTag map[string][]*models.Post
	tagsByPost map[uuid.UUID][]string
}

func (r *hashtagTestRepo) SyncPostHashtags(ctx context.Context, postID uuid.UUID, tags []string) error {
	return nil
}

func (r *hashtagTestRepo) DeletePostHashtags(ctx context.Context, postID uuid.UUID) error {
	return nil
}

func (r *hashtagTestRepo) FindTagsByPostIDs(ctx context.Context, postIDs []uuid.UUID) (map[uuid.UUID][]string, error) {
	tags := make(map[uuid.UUID][]string)
	for _, postID := range postIDs {
		tags[postID] = r.tagsByPost[postID]
	}
	return tags, nil
}

func (r *hashtagTestRepo) FindTrending(ctx context.Context, limit int) ([]*models.Hashtag, error) {
	return r.trending, nil
}

func (r *hashtagTestRepo) FindByTag(ctx context.Context, tag string) (*models.Hashtag, error) {
	if r.byTag == nil {
		return nil, nil
	}
	return r.byTag[tag], nil
}

func (r *hashtagTestRepo) FindPostsByTag(ctx context.Context, tag string, limit int, offset int) ([]*models.Post, error) {
	if r.postsByTag == nil {
		return nil, nil
	}
	return r.postsByTag[tag], nil
}

func (r *hashtagTestRepo) Search(ctx context.Context, query string, limit int, offset int) ([]*models.Hashtag, error) {
	return nil, nil
}

type hashtagTestPersonaRepo struct {
	personas map[string]*models.Persona
}

func (r *hashtagTestPersonaRepo) CreatePersona(ctx context.Context, userID uuid.UUID, dto personadtos.CreatePersonaDTO) (*models.Persona, error) {
	return nil, nil
}

func (r *hashtagTestPersonaRepo) FindPersonaByID(ctx context.Context, personaID uuid.UUID) (*models.Persona, error) {
	persona, ok := r.personas[personaID.String()]
	if !ok {
		return nil, persona_repo.ErrPersonaNotFound
	}
	return persona, nil
}

func (r *hashtagTestPersonaRepo) FindPersonasByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Persona, error) {
	return nil, nil
}

func (r *hashtagTestPersonaRepo) FindPersonaByHandle(ctx context.Context, handle string) (*models.Persona, error) {
	return nil, persona_repo.ErrPersonaNotFound
}

func (r *hashtagTestPersonaRepo) UpdatePersona(ctx context.Context, personaID uuid.UUID, dto personadtos.UpdatePersonaDTO) (*models.Persona, error) {
	return nil, nil
}

type hashtagTestLikeRepo struct {
	likedPostIDs map[uuid.UUID]bool
}

func (r *hashtagTestLikeRepo) LikePost(ctx context.Context, personaID uuid.UUID, postID uuid.UUID) error {
	return nil
}

func (r *hashtagTestLikeRepo) UnlikePost(ctx context.Context, personaID uuid.UUID, postID uuid.UUID) error {
	return nil
}

func (r *hashtagTestLikeRepo) HasPostLike(ctx context.Context, personaID uuid.UUID, postID uuid.UUID) (bool, error) {
	return r.likedPostIDs[postID], nil
}

func (r *hashtagTestLikeRepo) FindLikedPostIDs(ctx context.Context, personaID uuid.UUID, postIDs []uuid.UUID) (map[uuid.UUID]bool, error) {
	liked := make(map[uuid.UUID]bool)
	for _, postID := range postIDs {
		if r.likedPostIDs[postID] {
			liked[postID] = true
		}
	}
	return liked, nil
}
