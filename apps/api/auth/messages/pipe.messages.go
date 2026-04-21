package messages

import "github.com/emmanuella-codes/nox/shared"

const (
	// -
	User_Already_Exists shared.PipeMessage = "user_already_exists"
	Invalid_Credentials shared.PipeMessage = "invalid_credentials"
	Invalid_Token       shared.PipeMessage = "invalid_token"
	Invalid_Payload     shared.PipeMessage = "invalid_payload"

	// +
	User_Created    shared.PipeMessage = "user_created_successfully"
	User_Logged_In  shared.PipeMessage = "user_logged_in_successfully"
	Token_Refreshed shared.PipeMessage = "token_refreshed_successfully"
	User_Logged_Out shared.PipeMessage = "user_logged_out_successfully"
)
