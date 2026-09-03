package messages

import "github.com/emmanuella-codes/nox/shared"

const (
	// -
	Invalid_Payload           shared.PipeMessage = "invalid_payload"
	Handle_Required           shared.PipeMessage = "handle_required"
	Invalid_Persona_Type      shared.PipeMessage = "invalid_persona_type"
	Persona_Not_Found         shared.PipeMessage = "persona_not_found"
	Handle_Already_Taken      shared.PipeMessage = "handle_already_taken"
	Profile_Already_Exists    shared.PipeMessage = "profile_already_exists"
	Ghost_Persona_Already_Set shared.PipeMessage = "ghost_persona_already_exists"
	Forbidden                 shared.PipeMessage = "forbidden"
	Internal_Error            shared.PipeMessage = "internal_error"

	// +
	Persona_Created shared.PipeMessage = "persona_created_successfully"
	Persona_Updated shared.PipeMessage = "persona_updated_successfully"
	Persona_Fetched shared.PipeMessage = "persona_fetched_successfully"
	Personas_Listed shared.PipeMessage = "personas_listed_successfully"
)
