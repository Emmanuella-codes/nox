package messages

import "github.com/emmanuella-codes/nox/shared"

const (
	Invalid_Payload   shared.PipeMessage = "invalid_payload"
	Invalid_Set       shared.PipeMessage = "invalid_set"
	Persona_Not_Found shared.PipeMessage = "persona_not_found"
	Media_Not_Found   shared.PipeMessage = "media_asset_not_found"
	Set_Not_Found     shared.PipeMessage = "set_not_found"
	Forbidden         shared.PipeMessage = "forbidden"
	Internal_Error    shared.PipeMessage = "internal_error"
	Set_Created       shared.PipeMessage = "set_created_successfully"
	Set_Fetched       shared.PipeMessage = "set_fetched_successfully"
	Sets_Listed       shared.PipeMessage = "sets_listed_successfully"
	Set_Deleted       shared.PipeMessage = "set_deleted_successfully"
)
