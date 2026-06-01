package story

import (
	"context"
	"errors"

	"github.com/emmanuella-codes/nox/models"
	storydtos "github.com/emmanuella-codes/nox/story/dtos"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrStoryNotFound              = errors.New("story not found")
	ErrStoryItemNotFound          = errors.New("story item not found")
	ErrStoryDurationLimitExceeded = errors.New("story duration limit exceeded")
	ErrEventHighlightNotFound     = errors.New("event highlight story not found")
)

type StoryRepository interface {
	CreateStory(ctx context.Context, ownerUserID uuid.UUID, dto storydtos.CreateStoryDTO) (*models.Story, error)
	FindStoryByID(ctx context.Context, storyID uuid.UUID) (*models.Story, error)
	FindStoriesByEventID(ctx context.Context, eventID uuid.UUID, limit int, offset int) ([]*models.Story, error)
	FindStoriesByOwnerPersonaID(ctx context.Context, personaID uuid.UUID, limit int, offset int) ([]*models.Story, error)
	DeleteStory(ctx context.Context, storyID uuid.UUID) error
	AddStoryItem(ctx context.Context, storyID uuid.UUID, contributorUserID uuid.UUID, durationSeconds int, anonymousLabel string, dto storydtos.AddStoryItemDTO) (*models.StoryItem, error)
	FindStoryItems(ctx context.Context, storyID uuid.UUID) ([]*models.StoryItem, error)
	DeleteStoryItem(ctx context.Context, storyID uuid.UUID, itemID uuid.UUID) (*models.StoryItem, error)
	AddEventHighlightStory(ctx context.Context, eventID uuid.UUID, dto storydtos.AddEventHighlightStoryDTO) (*models.EventHighlightStory, error)
	FindEventHighlightStories(ctx context.Context, eventID uuid.UUID) ([]*models.EventHighlightStory, error)
	RemoveEventHighlightStory(ctx context.Context, eventID uuid.UUID, storyID uuid.UUID) error
}

func NewStoryRepository(db *pgxpool.Pool) StoryRepository {
	return newPgRepository(db)
}
