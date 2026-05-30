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

	res := pipe.PostsByTagPipe(context.Background(), "#Amapiano", 10)
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

func (r *hashtagTestRepo) FindPostsByTag(ctx context.Context, tag string, limit int) ([]*models.Post, error) {
	if r.postsByTag == nil {
		return nil, nil
	}
	return r.postsByTag[tag], nil
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
