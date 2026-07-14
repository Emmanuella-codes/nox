package messages

import "github.com/emmanuella-codes/nox/shared"

const (
	Invalid_Payload   shared.PipeMessage = "invalid_payload"
	Event_Not_Found   shared.PipeMessage = "event_not_found"
	Crew_Not_Found    shared.PipeMessage = "crew_not_found"
	Crew_Full         shared.PipeMessage = "crew_full"
	Persona_Not_Found shared.PipeMessage = "persona_not_found"
	Forbidden         shared.PipeMessage = "forbidden"
	Internal_Error    shared.PipeMessage = "internal_error"

	Crew_Created     shared.PipeMessage = "crew_created_successfully"
	Crew_Joined      shared.PipeMessage = "crew_joined_successfully"
	Crews_Listed     shared.PipeMessage = "crews_listed_successfully"
	Crew_Fetched     shared.PipeMessage = "crew_fetched_successfully"
	Crew_Left        shared.PipeMessage = "crew_left_successfully"
	Crew_Ended       shared.PipeMessage = "crew_ended_successfully"
	Sharing_Updated  shared.PipeMessage = "location_sharing_updated_successfully"
	Location_Updated shared.PipeMessage = "location_updated_successfully"
	Locations_Listed shared.PipeMessage = "locations_listed_successfully"
)
