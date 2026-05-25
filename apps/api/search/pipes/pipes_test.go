package pipes

import (
	"context"
	"testing"

	"github.com/emmanuella-codes/nox/models"
	searchrepo "github.com/emmanuella-codes/nox/repositories/search"
	"github.com/emmanuella-codes/nox/search/messages"
	"github.com/google/uuid"
)

func TestSearchRejectsShortQuery(t *testing.T) {
	pipe := NewSearchPipe(&searchTestRepo{}, nil, nil)

	res := pipe.Search(context.Background(), "a", 10)
	if res.Message != messages.Invalid_Query {
		t.Fatalf("expected %q, got %q", messages.Invalid_Query, res.Message)
	}
}

func TestSearchReturnsGroupedResults(t *testing.T) {
	personaID := uuid.New()
	postID := uuid.New()
	eventID := uuid.New()
	pipe := NewSearchPipe(&searchTestRepo{
		results: &searchrepo.Results{
			Personas: []*models.Persona{{
				ID:          personaID,
				Handle:      "djkay",
				DisplayName: "DJ Kay",
				PersonaType: models.VisiblePersonaType,
			}},
			Posts: []*searchrepo.PostResult{{
				Post: &models.Post{
					ID:          postID,
					PostingMode: models.PublicPostingMode,
					PersonaID:   &personaID,
					Body:        "afro house tonight",
					PostType:    models.TextPostType,
				},
				Persona: &models.Persona{
					ID:          personaID,
					Handle:      "djkay",
					DisplayName: "DJ Kay",
					PersonaType: models.VisiblePersonaType,
				},
			}},
			Events: []*models.Event{{
				ID:    eventID,
				Title: "Afro house night",
			}},
		},
	}, nil, nil)

	res := pipe.Search(context.Background(), "afro", 10)
	if !res.Success {
		t.Fatalf("expected success, got %q", res.Message)
	}
	if len(res.Data.Personas) != 1 || len(res.Data.Posts) != 1 || len(res.Data.Events) != 1 {
		t.Fatalf("expected grouped results, got %+v", res.Data)
	}
	if res.Data.Posts[0].Author.Persona == nil {
		t.Fatal("expected public post persona author")
	}
}

func TestSearchForViewerHydratesLikedState(t *testing.T) {
	viewerPersonaID := uuid.New()
	likedPostID := uuid.New()
	unlikedPostID := uuid.New()
	pipe := NewSearchPipe(&searchTestRepo{
		results: &searchrepo.Results{
			Posts: []*searchrepo.PostResult{
				{
					Post: &models.Post{
						ID:          likedPostID,
						PostingMode: models.AnonymousPostingMode,
						Body:        "afro liked",
						PostType:    models.TextPostType,
					},
				},
				{
					Post: &models.Post{
						ID:          unlikedPostID,
						PostingMode: models.AnonymousPostingMode,
						Body:        "afro unliked",
						PostType:    models.TextPostType,
					},
				},
			},
		},
	}, &searchTestLikeRepo{
		likedPostIDs: map[uuid.UUID]bool{likedPostID: true},
	}, nil)

	res := pipe.SearchForViewer(context.Background(), "afro", 10, viewerPersonaID)
	if !res.Success {
		t.Fatalf("expected search success, got %q", res.Message)
	}
	if len(res.Data.Posts) != 2 {
		t.Fatalf("expected 2 posts, got %d", len(res.Data.Posts))
	}
	if !res.Data.Posts[0].IsLiked {
		t.Fatal("expected first post to be liked")
	}
	if res.Data.Posts[1].IsLiked {
		t.Fatal("expected second post not to be liked")
	}
}

type searchTestRepo struct {
	results *searchrepo.Results
}

func (r *searchTestRepo) Search(ctx context.Context, query string, limit int) (*searchrepo.Results, error) {
	if r.results == nil {
		return &searchrepo.Results{}, nil
	}
	return r.results, nil
}

type searchTestLikeRepo struct {
	likedPostIDs map[uuid.UUID]bool
}

func (r *searchTestLikeRepo) LikePost(ctx context.Context, personaID uuid.UUID, postID uuid.UUID) error {
	return nil
}

func (r *searchTestLikeRepo) UnlikePost(ctx context.Context, personaID uuid.UUID, postID uuid.UUID) error {
	return nil
}

func (r *searchTestLikeRepo) HasPostLike(ctx context.Context, personaID uuid.UUID, postID uuid.UUID) (bool, error) {
	return r.likedPostIDs[postID], nil
}

func (r *searchTestLikeRepo) FindLikedPostIDs(ctx context.Context, personaID uuid.UUID, postIDs []uuid.UUID) (map[uuid.UUID]bool, error) {
	liked := make(map[uuid.UUID]bool)
	for _, postID := range postIDs {
		if r.likedPostIDs[postID] {
			liked[postID] = true
		}
	}
	return liked, nil
}
