package messages

import "github.com/emmanuella-codes/nox/shared"

const (
	Invalid_Query     shared.PipeMessage = "invalid_query"
	Persona_Not_Found shared.PipeMessage = "persona_not_found"
	Forbidden         shared.PipeMessage = "forbidden"
	Search_Listed     shared.PipeMessage = "search_listed_successfully"
	Internal_Error    shared.PipeMessage = "internal_error"
)
