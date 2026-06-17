import type { MediaAsset } from "@/src/types/api/set";

export type StoryContributionMode = "public" | "private";
export type StoryPostingMode = "public" | "anonymous";

export interface StoryOwner {
  id: string;
  handle: string;
  display_name: string;
  avatar_url: string;
  category: string;
}

export interface StoryItem {
  id: string;
  story_id: string;
  media_asset?: MediaAsset;
  contributor?: StoryOwner;
  posting_mode: StoryPostingMode;
  anonymous_label?: string;
  duration_seconds: number;
  position: number;
  created_at: string;
}

export interface Story {
  id: string;
  event_id: string;
  owner: StoryOwner;
  title: string;
  contribution_mode: StoryContributionMode;
  total_duration_seconds: number;
  can_contribute: boolean;
  items: StoryItem[];
  expires_at: string;
  created_at: string;
  updated_at: string;
}

export interface StoryListResponse {
  limit: number;
  offset: number;
  has_more: boolean;
  next_offset?: number;
  stories: Story[];
}

export interface CreateStoryRequest {
  event_id: string;
  owner_persona_id: string;
  title: string;
  contribution_mode: StoryContributionMode;
  expires_at?: string;
}

export interface AddStoryItemRequest {
  contributor_persona_id: string;
  media_asset_id: string;
  posting_mode: StoryPostingMode;
}

export interface EventHighlightStory {
  id: string;
  event_id: string;
  story_id: string;
  added_by_persona_id: string;
  position: number;
  story?: Story;
  created_at: string;
}
