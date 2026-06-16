import { apiRequest } from "@/src/utils/api/api";
import type {
  Story,
  StoryItem,
  StoryListResponse,
  CreateStoryRequest,
  AddStoryItemRequest,
  EventHighlightStory,
} from "@/src/types/api/story";

export function getStory(storyID: string, viewerPersonaID?: string, token?: string) {
  const params = new URLSearchParams();
  if (viewerPersonaID) params.set("viewer_persona_id", viewerPersonaID);
  const q = params.toString();
  return apiRequest<Story>(`/stories/${storyID}${q ? `?${q}` : ""}`, { token });
}

export function listStoryItems(storyID: string, viewerPersonaID?: string, token?: string) {
  const params = new URLSearchParams();
  if (viewerPersonaID) params.set("viewer_persona_id", viewerPersonaID);
  const q = params.toString();
  return apiRequest<StoryItem[]>(`/stories/${storyID}/items${q ? `?${q}` : ""}`, { token });
}

export function listEventStories(
  eventID: string,
  limit = 20,
  offset = 0,
  viewerPersonaID?: string,
  token?: string,
) {
  const params = new URLSearchParams({ limit: String(limit), offset: String(offset) });
  if (viewerPersonaID) params.set("viewer_persona_id", viewerPersonaID);
  return apiRequest<StoryListResponse>(`/events/${eventID}/stories?${params}`, { token });
}

export function listPersonaStories(
  personaID: string,
  limit = 20,
  offset = 0,
  viewerPersonaID?: string,
  token?: string,
) {
  const params = new URLSearchParams({ limit: String(limit), offset: String(offset) });
  if (viewerPersonaID) params.set("viewer_persona_id", viewerPersonaID);
  return apiRequest<StoryListResponse>(`/personas/${personaID}/stories?${params}`, { token });
}

export function listEventHighlightStories(eventID: string, viewerPersonaID?: string, token?: string) {
  const params = new URLSearchParams();
  if (viewerPersonaID) params.set("viewer_persona_id", viewerPersonaID);
  const q = params.toString();
  return apiRequest<EventHighlightStory[]>(`/events/${eventID}/highlight-stories${q ? `?${q}` : ""}`, { token });
}

export function createStory(payload: CreateStoryRequest, token: string) {
  return apiRequest<Story, CreateStoryRequest>("/stories", { method: "POST", body: payload, token });
}

export function deleteStory(storyID: string, token: string) {
  return apiRequest<null>(`/stories/${storyID}`, { method: "DELETE", token });
}

export function addStoryItem(storyID: string, payload: AddStoryItemRequest, token: string) {
  return apiRequest<StoryItem, AddStoryItemRequest>(`/stories/${storyID}/items`, {
    method: "POST",
    body: payload,
    token,
  });
}

export function deleteStoryItem(storyID: string, itemID: string, token: string) {
  return apiRequest<null>(`/stories/${storyID}/items/${itemID}`, { method: "DELETE", token });
}

export function addEventHighlightStory(
  eventID: string,
  storyID: string,
  addedByPersonaID: string,
  token: string,
) {
  return apiRequest<EventHighlightStory, { story_id: string; added_by_persona_id: string }>(
    `/events/${eventID}/highlight-stories`,
    { method: "POST", body: { story_id: storyID, added_by_persona_id: addedByPersonaID }, token },
  );
}

export function removeEventHighlightStory(eventID: string, storyID: string, token: string) {
  return apiRequest<null>(`/events/${eventID}/highlight-stories/${storyID}`, {
    method: "DELETE",
    token,
  });
}
