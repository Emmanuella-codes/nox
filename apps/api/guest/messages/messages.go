package messages

import "github.com/emmanuella-codes/nox/shared"

const (
	Invalid   shared.PipeMessage = "invalid_guest_list_payload"
	Forbidden shared.PipeMessage = "guest_list_access_denied"
	NotFound  shared.PipeMessage = "guest_list_entry_not_found"
	Created   shared.PipeMessage = "guest_list_resource_created"
	Claimed   shared.PipeMessage = "guest_pass_claimed"
	CheckedIn shared.PipeMessage = "guest_checked_in"
	Conflict  shared.PipeMessage = "guest_list_capacity_or_status_conflict"
	Internal  shared.PipeMessage = "internal_error"
)
