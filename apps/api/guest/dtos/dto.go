package dtos

type ConfigureDTO struct {
	Capacity int `json:"capacity"`
}
type QuotaDTO struct {
	OwnerUserID string `json:"owner_user_id"`
	Label       string `json:"label"`
	Capacity    int    `json:"capacity"`
}
type GuestDTO struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Quantity int    `json:"quantity"`
}
