package messages

import "github.com/emmanuella-codes/nox/shared"

const (
	Invalid_Payload          shared.PipeMessage = "invalid_payload"
	Invalid_Media            shared.PipeMessage = "invalid_media"
	Media_Not_Found          shared.PipeMessage = "media_asset_not_found"
	Persona_Not_Found        shared.PipeMessage = "persona_not_found"
	Forbidden                shared.PipeMessage = "forbidden"
	Internal_Error           shared.PipeMessage = "internal_error"
	Media_Asset_Created      shared.PipeMessage = "media_asset_created_successfully"
	Media_Upload_Initiated   shared.PipeMessage = "media_upload_initiated_successfully"
	Media_Processing_Updated shared.PipeMessage = "media_processing_updated_successfully"
	Media_Cleanup_Completed  shared.PipeMessage = "media_cleanup_completed_successfully"
)
