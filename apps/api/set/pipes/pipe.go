package pipes

import (
	"context"
	"regexp"
	"sort"
	"strings"

	"github.com/emmanuella-codes/nox/models"
	media_repo "github.com/emmanuella-codes/nox/repositories/media"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	set_repo "github.com/emmanuella-codes/nox/repositories/set"
	"github.com/emmanuella-codes/nox/set/messages"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

var genreTagPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,31}$`)

type SetPipe struct {
	setRepo     set_repo.SetRepository
	mediaRepo   media_repo.MediaRepository
	personaRepo persona_repo.PersonaRepository
}

func NewSetPipe(setRepo set_repo.SetRepository, mediaRepo media_repo.MediaRepository, personaRepo persona_repo.PersonaRepository) *SetPipe {
	return &SetPipe{setRepo: setRepo, mediaRepo: mediaRepo, personaRepo: personaRepo}
}

func pipeInternalError[T any](err error, operation string) *shared.PipeRes[T] {
	return shared.PipeInternalError[T](err, "set", operation, messages.Internal_Error)
}

func isDJPersona(persona *models.Persona) bool {
	return persona.PersonaType == models.VisiblePersonaType && persona.Category == models.DJPersonaCategory
}

func validSetMedia(asset *models.MediaAsset) bool {
	return asset.MediaKind == models.VideoMediaKind &&
		asset.ProcessingStatus == models.ReadyMediaStatus &&
		asset.DurationSeconds > 0 &&
		asset.DurationSeconds <= 900
}

func normalizeGenreTags(tags []string) ([]string, bool) {
	seen := map[string]bool{}
	normalized := make([]string, 0, len(tags))
	for _, tag := range tags {
		tag = strings.ToLower(strings.TrimSpace(tag))
		tag = strings.TrimPrefix(tag, "#")
		tag = strings.ReplaceAll(tag, " ", "-")
		if tag == "" {
			continue
		}
		if !genreTagPattern.MatchString(tag) {
			return nil, false
		}
		if !seen[tag] {
			seen[tag] = true
			normalized = append(normalized, tag)
		}
	}
	sort.Strings(normalized)
	return normalized, len(normalized) > 0 && len(normalized) <= 10
}

func (p *SetPipe) hydrateSet(ctx context.Context, set *models.Set) error {
	persona, err := p.personaRepo.FindPersonaByID(ctx, set.PersonaID)
	if err != nil {
		return err
	}
	asset, err := p.mediaRepo.FindMediaAssetByID(ctx, set.MediaAssetID)
	if err != nil {
		return err
	}
	set.Persona = persona
	set.MediaAsset = asset
	return nil
}

func (p *SetPipe) hydrateSets(ctx context.Context, sets []*models.Set) error {
	for _, set := range sets {
		if err := p.hydrateSet(ctx, set); err != nil {
			return err
		}
	}
	return nil
}

func (p *SetPipe) hydrateSetComments(ctx context.Context, comments []*models.SetComment) error {
	for _, comment := range comments {
		persona, err := p.personaRepo.FindPersonaByID(ctx, comment.PersonaID)
		if err != nil {
			return err
		}
		comment.Persona = persona
	}
	return nil
}

func (p *SetPipe) setResponse(ctx context.Context, set *models.Set, viewerPersonaID *uuid.UUID) (*SetResponse, error) {
	liked := false
	if viewerPersonaID != nil {
		current, err := p.setRepo.HasSetLike(ctx, *viewerPersonaID, set.ID)
		if err != nil {
			return nil, err
		}
		liked = current
	}
	response := setResponse(set, liked)
	return &response, nil
}
