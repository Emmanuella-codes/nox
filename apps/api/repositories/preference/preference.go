package preference

import (
	"context"
	"errors"

	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrPreferenceAlreadyExists = errors.New("preference already exists")
	ErrPreferenceNotFound      = errors.New("preference not found")
	ErrSelfTarget              = errors.New("self target not allowed")
	ErrInvalidTargetType       = errors.New("invalid discovery suppression target type")
)

type PreferenceRepository interface {
	BlockUser(ctx context.Context, blockerUserID uuid.UUID, blockedUserID uuid.UUID) error
	UnblockUser(ctx context.Context, blockerUserID uuid.UUID, blockedUserID uuid.UUID) error
	MuteUser(ctx context.Context, userID uuid.UUID, mutedUserID uuid.UUID) error
	UnmuteUser(ctx context.Context, userID uuid.UUID, mutedUserID uuid.UUID) error
	AddDiscoverySuppression(ctx context.Context, userID uuid.UUID, targetType models.DiscoverySuppressionTargetType, targetID uuid.UUID) error
	RemoveDiscoverySuppression(ctx context.Context, userID uuid.UUID, targetType models.DiscoverySuppressionTargetType, targetID uuid.UUID) error
	IsBlockedBetween(ctx context.Context, userA uuid.UUID, userB uuid.UUID) (bool, error)
	FindExcludedUserIDs(ctx context.Context, userID uuid.UUID) (map[uuid.UUID]bool, error)
	FindMutedUserIDs(ctx context.Context, userID uuid.UUID) (map[uuid.UUID]bool, error)
	FindSuppressedTargetIDs(ctx context.Context, userID uuid.UUID, targetType models.DiscoverySuppressionTargetType) (map[uuid.UUID]bool, error)
	DeleteFollowRelationsBetweenUsers(ctx context.Context, userA uuid.UUID, userB uuid.UUID) error
}

func NewPreferenceRepository(db *pgxpool.Pool) PreferenceRepository {
	return newPgRepository(db)
}
