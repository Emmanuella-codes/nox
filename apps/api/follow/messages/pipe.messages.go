package messages

import "github.com/emmanuella-codes/nox/shared"

const (
	Invalid_Payload         shared.PipeMessage = "invalid_payload"
	Persona_Not_Found       shared.PipeMessage = "persona_not_found"
	Already_Following       shared.PipeMessage = "already_following"
	Not_Following           shared.PipeMessage = "not_following"
	Self_Follow_Not_Allowed shared.PipeMessage = "self_follow_not_allowed"
	Forbidden               shared.PipeMessage = "forbidden"
	Internal_Error          shared.PipeMessage = "internal_error"

	Followed_Successfully   shared.PipeMessage = "followed_successfully"
	Unfollowed_Successfully shared.PipeMessage = "unfollowed_successfully"
	Followers_Listed        shared.PipeMessage = "followers_listed_successfully"
	Following_Listed        shared.PipeMessage = "following_listed_successfully"
	Follow_Status_Fetched   shared.PipeMessage = "follow_status_fetched_successfully"
)
