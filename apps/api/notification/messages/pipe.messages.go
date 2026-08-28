package messages

import "github.com/emmanuella-codes/nox/shared"

const (
	Invalid_Payload        shared.PipeMessage = "invalid_payload"
	Notification_Not_Found shared.PipeMessage = "notification_not_found"
	Persona_Not_Found      shared.PipeMessage = "persona_not_found"
	Forbidden              shared.PipeMessage = "forbidden"
	Internal_Error         shared.PipeMessage = "internal_error"

	Notifications_Listed shared.PipeMessage = "notifications_listed_successfully"
	Notification_Read    shared.PipeMessage = "notification_marked_read_successfully"
	Notifications_Read   shared.PipeMessage = "notifications_marked_read_successfully"
)
