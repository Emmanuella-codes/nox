package messages

import "github.com/emmanuella-codes/nox/shared"

const (
	Invalid_Payload   shared.PipeMessage = "invalid_payload"
	Post_Not_Found    shared.PipeMessage = "post_not_found"
	Persona_Not_Found shared.PipeMessage = "persona_not_found"
	Forbidden         shared.PipeMessage = "forbidden"
	Internal_Error    shared.PipeMessage = "internal_error"

	Post_Liked   shared.PipeMessage = "post_liked_successfully"
	Post_Unliked shared.PipeMessage = "post_unliked_successfully"
)
