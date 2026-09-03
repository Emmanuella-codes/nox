package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	searchrepo "github.com/emmanuella-codes/nox/repositories/search"
	"github.com/google/uuid"
)

type viewerSuppressionState struct {
	excludedUsers      map[uuid.UUID]bool
	mutedUsers         map[uuid.UUID]bool
	suppressedPersonas map[uuid.UUID]bool
	suppressedPosts    map[uuid.UUID]bool
	suppressedEvents   map[uuid.UUID]bool
	suppressedSets     map[uuid.UUID]bool
}

func (p *SearchPipe) applyViewerSuppression(ctx context.Context, viewerPersonaID uuid.UUID, results *searchrepo.Results) error {
	if p.preferenceRepo == nil || p.personaRepo == nil || results == nil {
		return nil
	}
	state, err := p.loadViewerSuppressionState(ctx, viewerPersonaID)
	if err != nil {
		return err
	}
	results.Personas = filterPersonas(results.Personas, state)
	results.Posts = filterPosts(results.Posts, state)
	results.Sets = filterSets(results.Sets, state)
	results.Events, err = p.filterEvents(ctx, results.Events, state)
	if err != nil {
		return err
	}
	return nil
}

func (p *SearchPipe) loadViewerSuppressionState(ctx context.Context, viewerPersonaID uuid.UUID) (*viewerSuppressionState, error) {
	viewer, err := p.personaRepo.FindPersonaByID(ctx, viewerPersonaID)
	if err != nil {
		return nil, err
	}
	excludedUsers, err := p.preferenceRepo.FindExcludedUserIDs(ctx, viewer.UserID)
	if err != nil {
		return nil, err
	}
	mutedUsers, err := p.preferenceRepo.FindMutedUserIDs(ctx, viewer.UserID)
	if err != nil {
		return nil, err
	}
	suppressedPersonas, err := p.preferenceRepo.FindSuppressedTargetIDs(ctx, viewer.UserID, models.PersonaSuppressionTargetType)
	if err != nil {
		return nil, err
	}
	suppressedPosts, err := p.preferenceRepo.FindSuppressedTargetIDs(ctx, viewer.UserID, models.PostSuppressionTargetType)
	if err != nil {
		return nil, err
	}
	suppressedEvents, err := p.preferenceRepo.FindSuppressedTargetIDs(ctx, viewer.UserID, models.EventSuppressionTargetType)
	if err != nil {
		return nil, err
	}
	suppressedSets, err := p.preferenceRepo.FindSuppressedTargetIDs(ctx, viewer.UserID, models.SetSuppressionTargetType)
	if err != nil {
		return nil, err
	}
	return &viewerSuppressionState{
		excludedUsers:      excludedUsers,
		mutedUsers:         mutedUsers,
		suppressedPersonas: suppressedPersonas,
		suppressedPosts:    suppressedPosts,
		suppressedEvents:   suppressedEvents,
		suppressedSets:     suppressedSets,
	}, nil
}

func filterPersonas(personas []*models.Persona, state *viewerSuppressionState) []*models.Persona {
	filtered := make([]*models.Persona, 0, len(personas))
	for _, persona := range personas {
		if persona == nil || state.excludedUsers[persona.UserID] || state.mutedUsers[persona.UserID] || state.suppressedPersonas[persona.ID] {
			continue
		}
		filtered = append(filtered, persona)
	}
	return filtered
}

func filterPosts(posts []*searchrepo.PostResult, state *viewerSuppressionState) []*searchrepo.PostResult {
	filtered := make([]*searchrepo.PostResult, 0, len(posts))
	for _, result := range posts {
		if result == nil || result.Post == nil {
			continue
		}
		if state.excludedUsers[result.Post.AuthorUserID] || state.mutedUsers[result.Post.AuthorUserID] || state.suppressedPosts[result.Post.ID] {
			continue
		}
		if result.Post.PersonaID != nil && state.suppressedPersonas[*result.Post.PersonaID] {
			continue
		}
		if result.Post.EventID != nil && state.suppressedEvents[*result.Post.EventID] {
			continue
		}
		filtered = append(filtered, result)
	}
	return filtered
}

func filterSets(sets []*models.Set, state *viewerSuppressionState) []*models.Set {
	filtered := make([]*models.Set, 0, len(sets))
	for _, set := range sets {
		if set == nil || state.excludedUsers[set.AuthorUserID] || state.mutedUsers[set.AuthorUserID] || state.suppressedSets[set.ID] || state.suppressedPersonas[set.PersonaID] {
			continue
		}
		filtered = append(filtered, set)
	}
	return filtered
}

func (p *SearchPipe) filterEvents(ctx context.Context, events []*models.Event, state *viewerSuppressionState) ([]*models.Event, error) {
	filtered := make([]*models.Event, 0, len(events))
	for _, event := range events {
		if event == nil || state.suppressedEvents[event.ID] || state.suppressedPersonas[event.OrganizerID] {
			continue
		}
		organizer, err := p.personaRepo.FindPersonaByID(ctx, event.OrganizerID)
		if err != nil {
			return nil, err
		}
		if state.excludedUsers[organizer.UserID] || state.mutedUsers[organizer.UserID] {
			continue
		}
		filtered = append(filtered, event)
	}
	return filtered, nil
}
