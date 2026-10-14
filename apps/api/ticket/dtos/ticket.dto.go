package dtos

import "github.com/google/uuid"

type CreateWaveDTO struct {
	Name         string `json:"name" validate:"required"`
	DisplayOrder int    `json:"display_order"`
	SaleStartAt  string `json:"sale_start_at"`
	SaleEndAt    string `json:"sale_end_at"`
	ReleaseMode  string `json:"release_mode"`
}
type CreateTierDTO struct {
	Name                 string `json:"name" validate:"required"`
	Description          string `json:"description"`
	PriceMinor           int64  `json:"price_minor"`
	Currency             string `json:"currency"`
	Inventory            int    `json:"inventory" validate:"required"`
	AttendeesPerPurchase int    `json:"attendees_per_purchase" validate:"required"`
	MaxPurchasesPerUser  int    `json:"max_purchases_per_user" validate:"required"`
}
type ReserveDTO struct {
	Quantity       int    `json:"quantity" validate:"required"`
	IdempotencyKey string `json:"idempotency_key" validate:"required"`
}
type WaveTierDTO struct {
	WaveID uuid.UUID `json:"wave_id" validate:"required"`
}
