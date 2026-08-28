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
	ErrStoryNotFound                    = errors.New("story not found")
	ErrStoryItemNotFound                = errors.New("story item not found")
	ErrStoryContributionRequestNotFound = errors.New("story contribution request not found")
	ErrStoryContributionRequestPending  = errors.New("story contribution request already pending")
	ErrStoryContributionRequestReviewed = errors.New("story contribution request already reviewed")
	ErrStoryDurationLimitExceeded       = errors.New("story duration limit exceeded")
	ErrStoryMediaInUse                  = errors.New("story media asset already used")
	ErrEventHighlightNotFound           = errors.New("event highlight story not found")
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
	ReorderStoryItem(ctx context.Context, storyID uuid.UUID, itemID uuid.UUID, position int) (*models.StoryItem, error)
	CreateStoryContributionRequest(ctx context.Context, storyID uuid.UUID, contributorUserID uuid.UUID, dto storydtos.CreateStoryContributionRequestDTO) (*models.StoryContributionRequest, error)
	FindStoryContributionRequestByID(ctx context.Context, storyID uuid.UUID, requestID uuid.UUID) (*models.StoryContributionRequest, error)
	FindStoryContributionRequests(ctx context.Context, storyID uuid.UUID, status *models.StoryContributionRequestStatus, limit int, offset int) ([]*models.StoryContributionRequest, error)
	AcceptStoryContributionRequest(ctx context.Context, story *models.Story, requestID uuid.UUID, reviewerPersonaID uuid.UUID, durationSeconds int) (*models.StoryContributionRequest, *models.StoryItem, error)
	RejectStoryContributionRequest(ctx context.Context, storyID uuid.UUID, requestID uuid.UUID, reviewerPersonaID uuid.UUID) (*models.StoryContributionRequest, error)
	AddEventHighlightStory(ctx context.Context, eventID uuid.UUID, dto storydtos.AddEventHighlightStoryDTO) (*models.EventHighlightStory, error)
	FindEventHighlightStories(ctx context.Context, eventID uuid.UUID) ([]*models.EventHighlightStory, error)
	ReorderEventHighlightStory(ctx context.Context, eventID uuid.UUID, storyID uuid.UUID, position int) (*models.EventHighlightStory, error)
	RemoveEventHighlightStory(ctx context.Context, eventID uuid.UUID, storyID uuid.UUID) error
}

func NewStoryRepository(db *pgxpool.Pool) StoryRepository {
	return newPgRepository(db)
}
