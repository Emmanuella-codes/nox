package messages

import "github.com/emmanuella-codes/nox/shared"

const (
	Invalid          shared.PipeMessage = "invalid_admission_payload"
	Forbidden        shared.PipeMessage = "admission_access_denied"
	NotFound         shared.PipeMessage = "admission_not_found"
	AlreadyCheckedIn shared.PipeMessage = "already_checked_in"
	InvalidAdmission shared.PipeMessage = "invalid_admission"
	CheckedIn        shared.PipeMessage = "admission_checked_in"
	Created          shared.PipeMessage = "admission_created"
	Listed           shared.PipeMessage = "admissions_listed"
	Reconciled       shared.PipeMessage = "admissions_reconciled"
	Internal         shared.PipeMessage = "internal_error"
)
