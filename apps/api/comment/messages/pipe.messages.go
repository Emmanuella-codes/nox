package messages

import "github.com/emmanuella-codes/nox/shared"

const (
	Invalid_Payload   shared.PipeMessage = "invalid_payload"
	Comment_Not_Found shared.PipeMessage = "comment_not_found"
	Post_Not_Found    shared.PipeMessage = "post_not_found"
	Persona_Not_Found shared.PipeMessage = "persona_not_found"
	Forbidden         shared.PipeMessage = "forbidden"
	Internal_Error    shared.PipeMessage = "internal_error"

	Comment_Created shared.PipeMessage = "comment_created_successfully"
	Comments_Listed shared.PipeMessage = "comments_listed_successfully"
	Comment_Deleted shared.PipeMessage = "comment_deleted_successfully"
)
