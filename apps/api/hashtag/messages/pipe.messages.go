package messages

import "github.com/emmanuella-codes/nox/shared"

const (
	Invalid_Tag       shared.PipeMessage = "invalid_hashtag"
	Hashtag_Not_Found shared.PipeMessage = "hashtag_not_found"
	Persona_Not_Found shared.PipeMessage = "persona_not_found"
	Forbidden         shared.PipeMessage = "forbidden"
	Internal_Error    shared.PipeMessage = "internal_error"

	Hashtag_Fetched shared.PipeMessage = "hashtag_fetched_successfully"
	Hashtags_Listed shared.PipeMessage = "hashtags_listed_successfully"
	Hashtag_Posts   shared.PipeMessage = "hashtag_posts_listed_successfully"
)
