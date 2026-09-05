package admin

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const adminUniqueViolationCode = "23505"

type pgRepository struct {
	db *pgxpool.Pool
}

func newPgRepository(db *pgxpool.Pool) *pgRepository {
	return &pgRepository{db: db}
}

var moderationTables = map[models.ModerationEntityType]string{
	models.ModerationEntityPersona:   "personas",
	models.ModerationEntityPost:      "posts",
	models.ModerationEntityComment:   "comments",
	models.ModerationEntityStory:     "stories",
	models.ModerationEntityStoryItem: "story_items",
	models.ModerationEntitySet:       "sets",
	models.ModerationEntityEvent:     "events",
}

func (r *pgRepository) Moderate(ctx context.Context, params ModerateParams) (*models.ModerationState, error) {
	table, ok := moderationTables[params.EntityType]
	if !ok {
		return nil, ErrModerationEntityNotFound
	}
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var previous models.ModerationStatus
	if err := tx.QueryRow(ctx, `SELECT moderation_status FROM `+table+` WHERE id = $1 FOR UPDATE`, params.EntityID).Scan(&previous); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrModerationEntityNotFound
		}
		return nil, err
	}
	row := tx.QueryRow(ctx, `UPDATE `+table+` SET moderation_status = $2, moderation_reason = $3,
		moderated_at = now(), moderated_by_user_id = $4 WHERE id = $1
		RETURNING moderation_status, moderation_reason, moderated_at, moderated_by_user_id`,
		params.EntityID, params.Status, params.Reason, params.AdminID)
	var state models.ModerationState
	if err := row.Scan(&state.Status, &state.Reason, &state.ModeratedAt, &state.ModeratedBy); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrModerationEntityNotFound
		}
		return nil, err
	}
	if _, err := tx.Exec(ctx, `INSERT INTO moderation_actions
		(entity_type, entity_id, previous_status, status, reason, admin_user_id)
		VALUES ($1, $2, $3, $4, $5, $6)`, params.EntityType, params.EntityID, previous, state.Status, state.Reason, params.AdminID); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	state.EntityType = params.EntityType
	state.EntityID = params.EntityID
	return &state, nil
}

func (r *pgRepository) ListModerationActions(ctx context.Context, params ListModerationActionsParams) ([]models.ModerationAction, error) {
	limit := params.Limit
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	offset := params.Offset
	if offset < 0 {
		offset = 0
	}
	rows, err := r.db.Query(ctx, `SELECT id, entity_type, entity_id, previous_status, status, reason, admin_user_id, created_at
		FROM moderation_actions
		WHERE ($1::text IS NULL OR entity_type = $1) AND ($2::text IS NULL OR status = $2)
		ORDER BY created_at DESC, id DESC LIMIT $3 OFFSET $4`, params.EntityType, params.Status, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	actions := make([]models.ModerationAction, 0)
	for rows.Next() {
		var action models.ModerationAction
		if err := rows.Scan(&action.ID, &action.EntityType, &action.EntityID, &action.PreviousStatus,
			&action.Status, &action.Reason, &action.AdminUserID, &action.CreatedAt); err != nil {
			return nil, err
		}
		actions = append(actions, action)
	}
	return actions, rows.Err()
}

func (r *pgRepository) ListIdentities(ctx context.Context) ([]models.AdminIdentity, error) {
	rows, err := r.db.Query(ctx, `
		SELECT u.id, u.fullname, u.email, am.role, am.is_active, am.created_at, am.updated_at
		FROM admin_memberships am
		INNER JOIN users u ON u.id = am.user_id
		ORDER BY am.created_at ASC, u.email ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	identities := []models.AdminIdentity{}
	for rows.Next() {
		identity, err := scanAdminIdentity(rows)
		if err != nil {
			return nil, err
		}
		identities = append(identities, *identity)
	}

	return identities, rows.Err()
}

func (r *pgRepository) FindMembershipByUserID(ctx context.Context, userID uuid.UUID) (*models.AdminMembership, error) {
	membership := &models.AdminMembership{}
	err := r.db.QueryRow(ctx, `
		SELECT user_id, role, is_active, created_at, updated_at
		FROM admin_memberships
		WHERE user_id = $1
	`, userID).Scan(
		&membership.UserID,
		&membership.Role,
		&membership.IsActive,
		&membership.CreatedAt,
		&membership.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return membership, nil
}

func (r *pgRepository) FindIdentityByUserID(ctx context.Context, userID uuid.UUID) (*models.AdminIdentity, error) {
	identity, err := scanAdminIdentity(r.db.QueryRow(ctx, `
		SELECT u.id, u.fullname, u.email, am.role, am.is_active, am.created_at, am.updated_at
		FROM admin_memberships am
		INNER JOIN users u ON u.id = am.user_id
		WHERE am.user_id = $1
	`, userID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return identity, nil
}

func (r *pgRepository) CreateMembership(ctx context.Context, userID uuid.UUID, role models.AdminRole) (*models.AdminIdentity, error) {
	_, err := r.db.Exec(ctx, `
		INSERT INTO admin_memberships (user_id, role)
		VALUES ($1, $2)
	`, userID, role)
	if err != nil {
		if isAdminUniqueViolation(err) {
			return nil, ErrMembershipAlreadyExists
		}
		return nil, err
	}

	return r.FindIdentityByUserID(ctx, userID)
}

func (r *pgRepository) UpdateMembershipRole(ctx context.Context, userID uuid.UUID, role models.AdminRole) (*models.AdminIdentity, error) {
	commandTag, err := r.db.Exec(ctx, `
		UPDATE admin_memberships
		SET role = $2,
			updated_at = now()
		WHERE user_id = $1
	`, userID, role)
	if err != nil {
		return nil, err
	}
	if commandTag.RowsAffected() == 0 {
		return nil, ErrMembershipNotFound
	}

	return r.FindIdentityByUserID(ctx, userID)
}

func (r *pgRepository) UpdateMembershipStatus(ctx context.Context, userID uuid.UUID, isActive bool) (*models.AdminIdentity, error) {
	commandTag, err := r.db.Exec(ctx, `
		UPDATE admin_memberships
		SET is_active = $2,
			updated_at = now()
		WHERE user_id = $1
	`, userID, isActive)
	if err != nil {
		return nil, err
	}
	if commandTag.RowsAffected() == 0 {
		return nil, ErrMembershipNotFound
	}

	return r.FindIdentityByUserID(ctx, userID)
}

func (r *pgRepository) CountActiveMembershipsByRole(ctx context.Context, role models.AdminRole) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM admin_memberships
		WHERE role = $1 AND is_active = TRUE
	`, role).Scan(&count)
	return count, err
}

func (r *pgRepository) CreateAuditLog(ctx context.Context, params CreateAuditLogParams) error {
	metadata := map[string]any{}
	if params.Metadata != nil {
		metadata = params.Metadata
	}

	encoded, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, `
		INSERT INTO admin_audit_logs (
			admin_user_id,
			action,
			target_user_id,
			request_id,
			ip_address,
			user_agent,
			metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7::jsonb)
	`, params.AdminUserID, params.Action, params.TargetUserID, params.RequestID, params.IPAddress, params.UserAgent, string(encoded))
	return err
}

func scanAdminIdentity(row interface {
	Scan(dest ...any) error
}) (*models.AdminIdentity, error) {
	identity := &models.AdminIdentity{}
	err := row.Scan(
		&identity.UserID,
		&identity.Fullname,
		&identity.Email,
		&identity.Role,
		&identity.IsActive,
		&identity.CreatedAt,
		&identity.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return identity, nil
}

func isAdminUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == adminUniqueViolationCode
}
