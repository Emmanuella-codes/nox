import type { Persona } from "@/src/types/api/user/persona";

export interface MediaAsset {
  id: string;
  owner_persona_id: string;
  media_kind: "video";
  playback_url: string;
  thumbnail_url: string;
  mime_type: string;
  duration_seconds: number;
  size_bytes: number;
  processing_status: "pending" | "ready" | "failed";
  created_at: string;
  updated_at: string;
}

export interface Set {
  id: string;
  persona_id: string;
  media_asset_id: string;
  title: string;
  description: string;
  genre_tags: string[];
  duration_seconds: number;
  like_count: number;
  comment_count: number;
  play_count: number;
  created_at: string;
  updated_at: string;
  persona?: Persona;
  media_asset?: MediaAsset;
}

export interface CreateMediaAssetRequest {
  owner_persona_id: string;
  storage_key: string;
  playback_url: string;
  thumbnail_url?: string;
  mime_type: string;
  duration_seconds: number;
  size_bytes: number;
}

export interface CreateSetRequest {
  persona_id: string;
  media_asset_id: string;
  title: string;
  description: string;
  genre_tags: string[];
}

export interface SetListResponse {
  limit: number;
  offset: number;
  has_more: boolean;
  next_offset?: number;
  sets: Set[];
}
