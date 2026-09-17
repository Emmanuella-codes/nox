package messages

import "github.com/emmanuella-codes/nox/shared"

const (
	Invalid     shared.PipeMessage = "invalid_ticket_payload"
	Forbidden   shared.PipeMessage = "forbidden"
	NotFound    shared.PipeMessage = "ticket_resource_not_found"
	Unavailable shared.PipeMessage = "ticket_inventory_unavailable"
	Internal    shared.PipeMessage = "internal_error"
	Created     shared.PipeMessage = "ticket_resource_created"
	Listed      shared.PipeMessage = "ticket_products_listed"
	Reserved    shared.PipeMessage = "ticket_inventory_reserved"
)
