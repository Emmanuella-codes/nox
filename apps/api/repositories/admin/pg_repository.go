package admin

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgRepository struct {
	db *pgxpool.Pool
}

func newPgRepository(db *pgxpool.Pool) *pgRepository {
	return &pgRepository{db: db}
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
	identity := &models.AdminIdentity{}
	err := r.db.QueryRow(ctx, `
		SELECT u.id, u.fullname, u.email, am.role, am.is_active, am.created_at, am.updated_at
		FROM admin_memberships am
		INNER JOIN users u ON u.id = am.user_id
		WHERE am.user_id = $1
	`, userID).Scan(
		&identity.UserID,
		&identity.Fullname,
		&identity.Email,
		&identity.Role,
		&identity.IsActive,
		&identity.CreatedAt,
		&identity.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	identity.Scopes = models.AdminScopesForRole(identity.Role)
	return identity, nil
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
