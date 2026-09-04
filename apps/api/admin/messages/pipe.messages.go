package messages

import "github.com/emmanuella-codes/nox/shared"

const (
	Invalid_Credentials  shared.PipeMessage = "invalid_credentials"
	Invalid_Token        shared.PipeMessage = "invalid_token"
	Invalid_Payload      shared.PipeMessage = "invalid_payload"
	Internal_Error       shared.PipeMessage = "internal_error"
	Admin_Access_Denied  shared.PipeMessage = "admin_access_denied"
	Admin_Logged_In      shared.PipeMessage = "admin_logged_in_successfully"
	Token_Refreshed      shared.PipeMessage = "token_refreshed_successfully"
	Admin_Logged_Out     shared.PipeMessage = "admin_logged_out_successfully"
	Admin_Profile_Loaded shared.PipeMessage = "admin_profile_loaded_successfully"
)
