package admin

import (
	"context"
	"errors"

	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *pgRepository) ListManagedUsers(ctx context.Context) ([]models.AdminManagedUser, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			u.id,
			u.fullname,
			u.email,
			u.email_verified,
			u.email_verified_at,
			u.status,
			am.role,
			am.is_active,
			u.created_at,
			u.updated_at
		FROM users u
		LEFT JOIN admin_memberships am ON am.user_id = u.id
		ORDER BY u.created_at DESC, u.email ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []models.AdminManagedUser{}
	for rows.Next() {
		user, err := scanManagedUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, *user)
	}

	return users, rows.Err()
}

func (r *pgRepository) FindManagedUserByID(ctx context.Context, userID uuid.UUID) (*models.AdminManagedUser, error) {
	user, err := scanManagedUser(r.db.QueryRow(ctx, `
		SELECT
			u.id,
			u.fullname,
			u.email,
			u.email_verified,
			u.email_verified_at,
			u.status,
			am.role,
			am.is_active,
			u.created_at,
			u.updated_at
		FROM users u
		LEFT JOIN admin_memberships am ON am.user_id = u.id
		WHERE u.id = $1
	`, userID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (r *pgRepository) UpdateUserStatus(ctx context.Context, userID uuid.UUID, status models.UserStatus) (*models.AdminManagedUser, error) {
	commandTag, err := r.db.Exec(ctx, `
		UPDATE users
		SET status = $2,
			updated_at = now()
		WHERE id = $1
	`, userID, status)
	if err != nil {
		return nil, err
	}
	if commandTag.RowsAffected() == 0 {
		return nil, ErrUserNotFound
	}

	return r.FindManagedUserByID(ctx, userID)
}

func scanManagedUser(row interface {
	Scan(dest ...any) error
}) (*models.AdminManagedUser, error) {
	user := &models.AdminManagedUser{}
	var role *models.AdminRole
	var adminActive *bool

	err := row.Scan(
		&user.ID,
		&user.Fullname,
		&user.Email,
		&user.EmailVerified,
		&user.EmailVerifiedAt,
		&user.Status,
		&role,
		&adminActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	user.Status = models.NormalizeUserStatus(user.Status)
	user.AdminRole = role
	user.AdminIsActive = adminActive
	return user, nil
}
