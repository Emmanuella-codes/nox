package pipes

import (
	"context"

	"github.com/google/uuid"
)

// hydrateSetLikedState attaches viewer-specific liked state to searched sets.
func (p *SearchPipe) hydrateSetLikedState(ctx context.Context, viewerPersonaID uuid.UUID, sets []SearchSetResponse) error {
	if p.setRepo == nil || len(sets) == 0 {
		return nil
	}
	setIDs := make([]uuid.UUID, 0, len(sets))
	for _, set := range sets {
		setID, err := uuid.Parse(set.ID)
		if err != nil {
			return err
		}
		setIDs = append(setIDs, setID)
	}
	liked, err := p.setRepo.FindLikedSetIDs(ctx, viewerPersonaID, setIDs)
	if err != nil {
		return err
	}
	for i := range sets {
		sets[i].IsLiked = liked[setIDs[i]]
	}
	return nil
}
