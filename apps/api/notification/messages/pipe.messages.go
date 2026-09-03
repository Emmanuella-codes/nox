package messages

import "github.com/emmanuella-codes/nox/shared"

const (
	Invalid_Payload        shared.PipeMessage = "invalid_payload"
	Notification_Not_Found shared.PipeMessage = "notification_not_found"
	Persona_Not_Found      shared.PipeMessage = "persona_not_found"
	Forbidden              shared.PipeMessage = "forbidden"
	Internal_Error         shared.PipeMessage = "internal_error"

	Notifications_Listed            shared.PipeMessage = "notifications_listed_successfully"
	Notification_Read               shared.PipeMessage = "notification_marked_read_successfully"
	Notifications_Read              shared.PipeMessage = "notifications_marked_read_successfully"
	Notification_Device_Upserted    shared.PipeMessage = "notification_device_upserted_successfully"
	Notification_Devices_Listed     shared.PipeMessage = "notification_devices_listed_successfully"
	Notification_Device_Removed     shared.PipeMessage = "notification_device_removed_successfully"
	Notification_Preferences_Listed shared.PipeMessage = "notification_preferences_listed_successfully"
	Notification_Preference_Updated shared.PipeMessage = "notification_preference_updated_successfully"
)
