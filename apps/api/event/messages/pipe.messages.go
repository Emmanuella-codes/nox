package messages

import "github.com/emmanuella-codes/nox/shared"

const (
	Invalid_Payload   shared.PipeMessage = "invalid_payload"
	Event_Not_Found   shared.PipeMessage = "event_not_found"
	Persona_Not_Found shared.PipeMessage = "persona_not_found"
	Forbidden         shared.PipeMessage = "forbidden"
	Internal_Error    shared.PipeMessage = "internal_error"

	Event_Created shared.PipeMessage = "event_created_successfully"
	Event_Fetched shared.PipeMessage = "event_fetched_successfully"
	Events_Listed shared.PipeMessage = "events_listed_successfully"
)
