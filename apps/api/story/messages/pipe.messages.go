package messages

import "github.com/emmanuella-codes/nox/shared"

const (
	Invalid_Payload                shared.PipeMessage = "invalid_payload"
	Invalid_Story                  shared.PipeMessage = "invalid_story"
	Story_Not_Found                shared.PipeMessage = "story_not_found"
	Story_Item_Not_Found           shared.PipeMessage = "story_item_not_found"
	Event_Not_Found                shared.PipeMessage = "event_not_found"
	Persona_Not_Found              shared.PipeMessage = "persona_not_found"
	Media_Asset_Not_Found          shared.PipeMessage = "media_asset_not_found"
	Forbidden                      shared.PipeMessage = "forbidden"
	Story_Duration_Limit_Exceeded  shared.PipeMessage = "story_duration_limit_exceeded"
	Internal_Error                 shared.PipeMessage = "internal_error"
	Story_Created                  shared.PipeMessage = "story_created_successfully"
	Story_Fetched                  shared.PipeMessage = "story_fetched_successfully"
	Stories_Listed                 shared.PipeMessage = "stories_listed_successfully"
	Story_Deleted                  shared.PipeMessage = "story_deleted_successfully"
	Story_Item_Added               shared.PipeMessage = "story_item_added_successfully"
	Story_Item_Deleted             shared.PipeMessage = "story_item_deleted_successfully"
	Story_Items_Listed             shared.PipeMessage = "story_items_listed_successfully"
	Event_Highlight_Story_Added    shared.PipeMessage = "event_highlight_story_added_successfully"
	Event_Highlight_Stories_Listed shared.PipeMessage = "event_highlight_stories_listed_successfully"
	Event_Highlight_Story_Removed  shared.PipeMessage = "event_highlight_story_removed_successfully"
)
