package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	media_repo "github.com/emmanuella-codes/nox/repositories/media"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	set_repo "github.com/emmanuella-codes/nox/repositories/set"
	"github.com/emmanuella-codes/nox/set/messages"
	"github.com/emmanuella-codes/nox/shared"
)

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
