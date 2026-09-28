package messages

import "github.com/emmanuella-codes/nox/shared"

const (
	Invalid   shared.PipeMessage = "invalid_guest_list_payload"
	Forbidden shared.PipeMessage = "guest_list_access_denied"
	NotFound  shared.PipeMessage = "guest_list_entry_not_found"
	Created   shared.PipeMessage = "guest_list_resource_created"
	Claimed   shared.PipeMessage = "guest_pass_claimed"
	CheckedIn shared.PipeMessage = "guest_checked_in"
	Declined  shared.PipeMessage = "guest_invitation_declined"
	Revoked   shared.PipeMessage = "guest_pass_revoked"
	Attached  shared.PipeMessage = "guest_pass_attached"
	Reported  shared.PipeMessage = "guest_list_reported"
	Conflict  shared.PipeMessage = "guest_list_capacity_or_status_conflict"
	Internal  shared.PipeMessage = "internal_error"
)
