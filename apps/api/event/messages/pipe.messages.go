package messages

import "github.com/emmanuella-codes/nox/shared"

const (
	Invalid_Payload   shared.PipeMessage = "invalid_payload"
	Event_Not_Found   shared.PipeMessage = "event_not_found"
	Persona_Not_Found shared.PipeMessage = "persona_not_found"
	Forbidden         shared.PipeMessage = "forbidden"
	Internal_Error    shared.PipeMessage = "internal_error"

	Event_Created      shared.PipeMessage = "event_created_successfully"
	Event_Fetched      shared.PipeMessage = "event_fetched_successfully"
	Events_Listed      shared.PipeMessage = "events_listed_successfully"
	Event_Updated      shared.PipeMessage = "event_updated_successfully"
	Event_Published    shared.PipeMessage = "event_published_successfully"
	Event_Closed       shared.PipeMessage = "event_moved_to_follow_up_successfully"
	Event_Archived     shared.PipeMessage = "event_archived_successfully"
	Invalid_Transition shared.PipeMessage = "invalid_event_transition"
)
