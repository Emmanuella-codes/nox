package dtos

type CrewSafetyActionDTO struct {
	Reason string `json:"reason" validate:"required,max=500"`
}

type CrewLocationAccessDTO struct {
	Reason   string `json:"reason" validate:"required,max=500"`
	ReportID string `json:"report_id" validate:"omitempty,uuid4"`
}
