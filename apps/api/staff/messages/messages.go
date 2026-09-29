package messages

import "github.com/emmanuella-codes/nox/shared"

const (
	Invalid   shared.PipeMessage = "invalid_staff_payload"
	Forbidden shared.PipeMessage = "staff_access_denied"
	NotFound  shared.PipeMessage = "staff_membership_not_found"
	Created   shared.PipeMessage = "staff_resource_created"
	Accepted  shared.PipeMessage = "staff_membership_accepted"
	Revoked   shared.PipeMessage = "staff_membership_revoked"
	Listed    shared.PipeMessage = "staff_members_listed"
	Internal  shared.PipeMessage = "internal_error"
)
