package dtos

type InviteDTO struct {
	UserID string `json:"user_id"`
}
type RolesDTO struct {
	Roles []string `json:"roles"`
}
type ShiftDTO struct {
	StartsAt string `json:"starts_at"`
	EndsAt   string `json:"ends_at"`
}
