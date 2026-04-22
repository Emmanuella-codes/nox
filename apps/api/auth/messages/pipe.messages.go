package messages

import "github.com/emmanuella-codes/nox/shared"

const (
	// -
	User_Already_Exists    shared.PipeMessage = "user_already_exists"
	Invalid_Credentials    shared.PipeMessage = "invalid_credentials"
	Invalid_Token          shared.PipeMessage = "invalid_token"
	Invalid_Payload        shared.PipeMessage = "invalid_payload"
	Internal_Error         shared.PipeMessage = "internal_error"
	Email_Not_Verified     shared.PipeMessage = "email_not_verified"
	Email_Already_Verified shared.PipeMessage = "email_already_verified"
	Invalid_OTP            shared.PipeMessage = "invalid_otp"
	OTP_Expired            shared.PipeMessage = "otp_expired"

	// +
	User_Created      shared.PipeMessage = "user_created_successfully"
	User_Logged_In    shared.PipeMessage = "user_logged_in_successfully"
	Token_Refreshed   shared.PipeMessage = "token_refreshed_successfully"
	User_Logged_Out   shared.PipeMessage = "user_logged_out_successfully"
	Verification_Sent shared.PipeMessage = "verification_sent"
	Email_Verified    shared.PipeMessage = "email_verified_successfully"
)
