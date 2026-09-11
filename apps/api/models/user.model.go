package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusSuspended UserStatus = "suspended"
)

type User struct {
	ID              uuid.UUID    `json:"id"`
	Fullname        string       `json:"fullname"`
	Email           string       `json:"email"`
	Password        string       `json:"-"`
	EmailVerified   bool         `json:"email_verified"`
	EmailVerifiedAt sql.NullTime `json:"-"`
	Status          UserStatus   `json:"status"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
}

type AdminManagedUser struct {
	ID              uuid.UUID    `json:"id"`
	Fullname        string       `json:"fullname"`
	Email           string       `json:"email"`
	EmailVerified   bool         `json:"email_verified"`
	EmailVerifiedAt sql.NullTime `json:"-"`
	Status          UserStatus   `json:"status"`
	AdminRole       *AdminRole   `json:"admin_role,omitempty"`
	AdminIsActive   *bool        `json:"admin_is_active,omitempty"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
}

func NormalizeUserStatus(status UserStatus) UserStatus {
	if status == "" {
		return UserStatusActive
	}

	return status
}

func ValidUserStatus(status UserStatus) bool {
	switch status {
	case UserStatusActive, UserStatusSuspended:
		return true
	default:
		return false
	}
}
