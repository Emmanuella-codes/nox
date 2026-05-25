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
	pipe := NewSearchPipe(&searchTestRepo{})

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
	})

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

type searchTestRepo struct {
	results *searchrepo.Results
}

func (r *searchTestRepo) Search(ctx context.Context, query string, limit int) (*searchrepo.Results, error) {
	if r.results == nil {
		return &searchrepo.Results{}, nil
	}
	return r.results, nil
}
