package dtos

type ManualCheckInDTO struct {
	AdmissionID string `json:"admission_id"`
}
type OverrideDTO struct {
	Reason string `json:"reason"`
}
type CreateAdmissionDTO struct {
	SourceType string `json:"source_type"`
	SourceID   string `json:"source_id"`
	UserID     string `json:"user_id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
}
