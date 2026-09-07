package dtos

import "time"

type FeatureSetDTO struct {
	ExpiresAt *time.Time `json:"expires_at"`
}

type ReviewStoryContributionDTO struct {
	Action string `json:"action" validate:"required"`
	Reason string `json:"reason" validate:"required,max=500"`
}
