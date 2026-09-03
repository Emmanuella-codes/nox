package preference

import (
	"context"
	"errors"

	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgRepository struct {
	db *pgxpool.Pool
}

func newPgRepository(db *pgxpool.Pool) *pgRepository {
	return &pgRepository{db: db}
}

func (r *pgRepository) BlockUser(ctx context.Context, blockerUserID uuid.UUID, blockedUserID uuid.UUID) error {
	if blockerUserID == blockedUserID {
		return ErrSelfTarget
	}
	_, err := r.db.Exec(ctx, `
		INSERT INTO user_blocks (blocker_user_id, blocked_user_id)
		VALUES ($1, $2)
	`, blockerUserID, blockedUserID)
	return mapPreferenceError(err)
}

func (r *pgRepository) UnblockUser(ctx context.Context, blockerUserID uuid.UUID, blockedUserID uuid.UUID) error {
	commandTag, err := r.db.Exec(ctx, `
		DELETE FROM user_blocks
		WHERE blocker_user_id = $1 AND blocked_user_id = $2
	`, blockerUserID, blockedUserID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return ErrPreferenceNotFound
	}
	return nil
}

func (r *pgRepository) MuteUser(ctx context.Context, userID uuid.UUID, mutedUserID uuid.UUID) error {
	if userID == mutedUserID {
		return ErrSelfTarget
	}
	_, err := r.db.Exec(ctx, `
		INSERT INTO user_mutes (user_id, muted_user_id)
		VALUES ($1, $2)
	`, userID, mutedUserID)
	return mapPreferenceError(err)
}

func (r *pgRepository) UnmuteUser(ctx context.Context, userID uuid.UUID, mutedUserID uuid.UUID) error {
	commandTag, err := r.db.Exec(ctx, `
		DELETE FROM user_mutes
		WHERE user_id = $1 AND muted_user_id = $2
	`, userID, mutedUserID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return ErrPreferenceNotFound
	}
	return nil
}

func (r *pgRepository) AddDiscoverySuppression(ctx context.Context, userID uuid.UUID, targetType models.DiscoverySuppressionTargetType, targetID uuid.UUID) error {
	if !validSuppressionTargetType(targetType) {
		return ErrInvalidTargetType
	}
	_, err := r.db.Exec(ctx, `
		INSERT INTO discovery_suppressions (user_id, target_type, target_id)
		VALUES ($1, $2, $3)
	`, userID, targetType, targetID)
	return mapPreferenceError(err)
}

func (r *pgRepository) RemoveDiscoverySuppression(ctx context.Context, userID uuid.UUID, targetType models.DiscoverySuppressionTargetType, targetID uuid.UUID) error {
	if !validSuppressionTargetType(targetType) {
		return ErrInvalidTargetType
	}
	commandTag, err := r.db.Exec(ctx, `
		DELETE FROM discovery_suppressions
		WHERE user_id = $1 AND target_type = $2 AND target_id = $3
	`, userID, targetType, targetID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return ErrPreferenceNotFound
	}
	return nil
}

func (r *pgRepository) IsBlockedBetween(ctx context.Context, userA uuid.UUID, userB uuid.UUID) (bool, error) {
	row := r.db.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM user_blocks
			WHERE (blocker_user_id = $1 AND blocked_user_id = $2)
			   OR (blocker_user_id = $2 AND blocked_user_id = $1)
		)
	`, userA, userB)
	var blocked bool
	if err := row.Scan(&blocked); err != nil {
		return false, err
	}
	return blocked, nil
}

func (r *pgRepository) FindExcludedUserIDs(ctx context.Context, userID uuid.UUID) (map[uuid.UUID]bool, error) {
	rows, err := r.db.Query(ctx, `
		SELECT blocked_user_id
		FROM user_blocks
		WHERE blocker_user_id = $1
		UNION
		SELECT blocker_user_id
		FROM user_blocks
		WHERE blocked_user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanIDMap(rows)
}

func (r *pgRepository) FindMutedUserIDs(ctx context.Context, userID uuid.UUID) (map[uuid.UUID]bool, error) {
	rows, err := r.db.Query(ctx, `
		SELECT muted_user_id
		FROM user_mutes
		WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanIDMap(rows)
}

func (r *pgRepository) FindSuppressedTargetIDs(ctx context.Context, userID uuid.UUID, targetType models.DiscoverySuppressionTargetType) (map[uuid.UUID]bool, error) {
	if !validSuppressionTargetType(targetType) {
		return nil, ErrInvalidTargetType
	}
	rows, err := r.db.Query(ctx, `
		SELECT target_id
		FROM discovery_suppressions
		WHERE user_id = $1 AND target_type = $2
	`, userID, targetType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanIDMap(rows)
}

func (r *pgRepository) DeleteFollowRelationsBetweenUsers(ctx context.Context, userA uuid.UUID, userB uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
		DELETE FROM persona_follows
		WHERE (follower_id IN (SELECT id FROM personas WHERE user_id = $1) AND following_id IN (SELECT id FROM personas WHERE user_id = $2))
		   OR (follower_id IN (SELECT id FROM personas WHERE user_id = $2) AND following_id IN (SELECT id FROM personas WHERE user_id = $1))
	`, userA, userB)
	return err
}

func validSuppressionTargetType(targetType models.DiscoverySuppressionTargetType) bool {
	switch targetType {
	case models.PersonaSuppressionTargetType, models.PostSuppressionTargetType, models.EventSuppressionTargetType, models.SetSuppressionTargetType:
		return true
	default:
		return false
	}
}

func scanIDMap(rows pgx.Rows) (map[uuid.UUID]bool, error) {
	ids := map[uuid.UUID]bool{}
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids[id] = true
	}
	return ids, rows.Err()
}

func mapPreferenceError(err error) error {
	if err == nil {
		return nil
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return ErrPreferenceAlreadyExists
	}
	return err
}
