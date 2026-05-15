package pipes

import (
	"context"
	"testing"

	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/persona/dtos"
	"github.com/emmanuella-codes/nox/persona/messages"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	"github.com/google/uuid"
)

func TestCreatePersonaPipeRejectsGhostPersona(t *testing.T) {
	repo := &personaTestRepo{}
	pipe := NewPersonaPipe(repo)

	res := pipe.CreatePersonaPipe(context.Background(), uuid.New(), dtos.CreatePersonaDTO{
		Handle:      "ghost_ada",
		DisplayName: "Ghost Ada",
		Bio:         "bio",
		AvatarURL:   "https://example.com/avatar.png",
		CoverURL:    "https://example.com/cover.png",
		PersonaType: models.GhostPersonaType,
		GenreTags:   []string{"house"},
	})
	if res.Message != messages.Invalid_Persona_Type {
		t.Fatalf("expected %q, got %q", messages.Invalid_Persona_Type, res.Message)
	}
}

func TestCreatePersonaPipeRequiresVisibleHandle(t *testing.T) {
	repo := &personaTestRepo{}
	pipe := NewPersonaPipe(repo)

	res := pipe.CreatePersonaPipe(context.Background(), uuid.New(), dtos.CreatePersonaDTO{
		DisplayName: "Visible Ada",
		Bio:         "bio",
		AvatarURL:   "https://example.com/avatar.png",
		CoverURL:    "https://example.com/cover.png",
		PersonaType: models.VisiblePersonaType,
		GenreTags:   []string{"house"},
	})
	if res.Message != messages.Handle_Required {
		t.Fatalf("expected %q, got %q", messages.Handle_Required, res.Message)
	}
}

func TestCreatePersonaPipePreservesVisibleHandle(t *testing.T) {
	repo := &personaTestRepo{}
	pipe := NewPersonaPipe(repo)

	res := pipe.CreatePersonaPipe(context.Background(), uuid.New(), dtos.CreatePersonaDTO{
		Handle:      "  MyHandle  ",
		DisplayName: "Visible Ada",
		Bio:         "bio",
		AvatarURL:   "https://example.com/avatar.png",
		CoverURL:    "https://example.com/cover.png",
		PersonaType: models.VisiblePersonaType,
		GenreTags:   []string{"house"},
	})
	if !res.Success {
		t.Fatalf("expected success, got %q", res.Message)
	}
	if repo.lastCreateDTO.Handle != "myhandle" {
		t.Fatalf("expected normalized visible handle, got %q", repo.lastCreateDTO.Handle)
	}
}

type personaTestRepo struct {
	lastCreateDTO dtos.CreatePersonaDTO
}

func (r *personaTestRepo) CreatePersona(ctx context.Context, userID uuid.UUID, dto dtos.CreatePersonaDTO) (*models.Persona, error) {
	r.lastCreateDTO = dto
	return &models.Persona{
		ID:          uuid.New(),
		UserID:      userID,
		Handle:      dto.Handle,
		DisplayName: dto.DisplayName,
		Bio:         dto.Bio,
		AvatarURL:   dto.AvatarURL,
		CoverURL:    dto.CoverURL,
		PersonaType: dto.PersonaType,
		GenreTags:   dto.GenreTags,
	}, nil
}

func (r *personaTestRepo) FindPersonaByID(ctx context.Context, personaID uuid.UUID) (*models.Persona, error) {
	return nil, persona_repo.ErrPersonaNotFound
}

func (r *personaTestRepo) FindPersonasByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Persona, error) {
	return nil, nil
}

func (r *personaTestRepo) FindPersonaByHandle(ctx context.Context, handle string) (*models.Persona, error) {
	return nil, persona_repo.ErrPersonaNotFound
}

func (r *personaTestRepo) UpdatePersona(ctx context.Context, personaID uuid.UUID, dto dtos.UpdatePersonaDTO) (*models.Persona, error) {
	return nil, nil
}
