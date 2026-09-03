package messages

import "github.com/emmanuella-codes/nox/shared"

const (
	Invalid_Payload        shared.PipeMessage = "invalid_payload"
	Conversation_Not_Found shared.PipeMessage = "conversation_not_found"
	Message_Not_Found      shared.PipeMessage = "message_not_found"
	Persona_Not_Found      shared.PipeMessage = "persona_not_found"
	Forbidden              shared.PipeMessage = "forbidden"
	Internal_Error         shared.PipeMessage = "internal_error"

	Conversation_Created shared.PipeMessage = "conversation_created_successfully"
	Conversations_Listed shared.PipeMessage = "conversations_listed_successfully"
	Conversation_Fetched shared.PipeMessage = "conversation_fetched_successfully"
	Message_Sent         shared.PipeMessage = "message_sent_successfully"
	Messages_Listed      shared.PipeMessage = "messages_listed_successfully"
	Conversation_Read    shared.PipeMessage = "conversation_marked_read_successfully"
	Members_Added        shared.PipeMessage = "members_added_successfully"
	Member_Removed       shared.PipeMessage = "member_removed_successfully"
	Member_Left          shared.PipeMessage = "member_left_successfully"
	Member_Role_Updated  shared.PipeMessage = "member_role_updated_successfully"
	Message_Deleted      shared.PipeMessage = "message_deleted_successfully"
	Message_Updated      shared.PipeMessage = "message_updated_successfully"
	Typing_Updated       shared.PipeMessage = "typing_updated_successfully"
)
