package pipes

import (
	"context"
	"testing"

	"github.com/emmanuella-codes/nox/models"
	searchrepo "github.com/emmanuella-codes/nox/repositories/search"
	"github.com/google/uuid"
)

// TestSearchForwardsScopeAndMapsSetResults verifies scoped discovery search returns sets and normalized scope metadata.
func TestSearchForwardsScopeAndMapsSetResults(t *testing.T) {
	repo := &searchTestRepo{
		results: &searchrepo.Results{
			Sets: []*models.Set{{
				ID:          mustUUID("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
				PersonaID:   mustUUID("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"),
				Title:       "Late Night Set",
				Description: "Afro house",
			}},
		},
	}
	pipe := NewSearchPipe(repo, nil, nil)

	res := pipe.Search(context.Background(), "late", SearchOptions{Limit: 10, Scope: "sets"})
	if !res.Success {
		t.Fatalf("expected success, got %q", res.Message)
	}
	if repo.options.Scope != "sets" {
		t.Fatalf("expected scope sets, got %q", repo.options.Scope)
	}
	if res.Data.Scope != "sets" || len(res.Data.Sets) != 1 {
		t.Fatalf("expected one scoped set result, got %+v", res.Data)
	}
}

// mustUUID parses one static UUID fixture for tests.
func mustUUID(value string) uuid.UUID {
	id, err := uuid.Parse(value)
	if err != nil {
		panic(err)
	}
	return id
}
