package messages

import "github.com/emmanuella-codes/nox/shared"

const (
	Blocked_Successfully          shared.PipeMessage = "blocked_successfully"
	Unblocked_Successfully        shared.PipeMessage = "unblocked_successfully"
	Muted_Successfully            shared.PipeMessage = "muted_successfully"
	Unmuted_Successfully          shared.PipeMessage = "unmuted_successfully"
	Discovery_Suppression_Added   shared.PipeMessage = "discovery_suppression_added_successfully"
	Discovery_Suppression_Removed shared.PipeMessage = "discovery_suppression_removed_successfully"
	Persona_Not_Found             shared.PipeMessage = "persona_not_found"
	Target_Not_Found              shared.PipeMessage = "target_not_found"
	Invalid_Payload               shared.PipeMessage = "invalid_payload"
	Forbidden                     shared.PipeMessage = "forbidden"
	Internal_Error                shared.PipeMessage = "internal_error"
)
