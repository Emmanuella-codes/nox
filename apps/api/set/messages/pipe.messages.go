package messages

import "github.com/emmanuella-codes/nox/shared"

const (
	Invalid_Payload     shared.PipeMessage = "invalid_payload"
	Invalid_Set         shared.PipeMessage = "invalid_set"
	Persona_Not_Found   shared.PipeMessage = "persona_not_found"
	Media_Not_Found     shared.PipeMessage = "media_asset_not_found"
	Set_Not_Found       shared.PipeMessage = "set_not_found"
	Comment_Not_Found   shared.PipeMessage = "comment_not_found"
	Media_In_Use        shared.PipeMessage = "media_asset_already_used"
	Forbidden           shared.PipeMessage = "forbidden"
	Internal_Error      shared.PipeMessage = "internal_error"
	Set_Created         shared.PipeMessage = "set_created_successfully"
	Set_Fetched         shared.PipeMessage = "set_fetched_successfully"
	Sets_Listed         shared.PipeMessage = "sets_listed_successfully"
	Set_Deleted         shared.PipeMessage = "set_deleted_successfully"
	Set_Liked           shared.PipeMessage = "set_liked_successfully"
	Set_Unliked         shared.PipeMessage = "set_unliked_successfully"
	Set_Play_Recorded   shared.PipeMessage = "set_play_recorded_successfully"
	Set_Commented       shared.PipeMessage = "set_comment_created_successfully"
	Set_Comments_Listed shared.PipeMessage = "set_comments_listed_successfully"
)
