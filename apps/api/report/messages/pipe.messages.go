package messages

import "github.com/emmanuella-codes/nox/shared"

const (
	ReportCreated        shared.PipeMessage = "report_created_successfully"
	InvalidReportTarget  shared.PipeMessage = "invalid_report_target"
	ReportAlreadyExists  shared.PipeMessage = "report_already_exists"
	CannotReportOwnData  shared.PipeMessage = "cannot_report_own_content"
	ReportTargetNotFound shared.PipeMessage = "report_target_not_found"
	InternalError        shared.PipeMessage = "internal_error"
	InvalidPayload       shared.PipeMessage = "invalid_payload"
	PersonaNotFound      shared.PipeMessage = "persona_not_found"
	Forbidden            shared.PipeMessage = "forbidden"
)
