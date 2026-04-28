package messages

import "github.com/emmanuella-codes/nox/shared"

const (
	Invalid_Payload   shared.PipeMessage = "invalid_payload"
	Post_Not_Found    shared.PipeMessage = "post_not_found"
	Persona_Not_Found shared.PipeMessage = "persona_not_found"
	Forbidden         shared.PipeMessage = "forbidden"
	Internal_Error    shared.PipeMessage = "internal_error"

	Post_Created shared.PipeMessage = "post_created_successfully"
	Post_Fetched shared.PipeMessage = "post_fetched_successfully"
	Posts_Listed shared.PipeMessage = "posts_listed_successfully"
	Feed_Listed  shared.PipeMessage = "feed_listed_successfully"
	Post_Deleted shared.PipeMessage = "post_deleted_successfully"
)
