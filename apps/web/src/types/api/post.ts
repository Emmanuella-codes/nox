export interface PostAuthor {
  mode: "public" | "anonymous";
  anonymous?: {
    handle: string;
  };
  persona?: {
    id: string;
    handle: string;
    display_name: string;
    avatar_url: string;
  };
}

export interface Post {
  id: string;
  author: PostAuthor;
  event_id?: string;
  body: string;
  post_type: "text" | "image" | "video" | "set" | "event_tag";
  media_url?: string;
  media_type?: "image" | "youtube" | "soundcloud";
  location?: string;
  like_count: number;
  comment_count: number;
  repost_count: number;
  is_liked: boolean;
  is_repost: boolean;
  repost_of?: string;
  hashtags: string[];
  media: {
    id: string;
    owner_persona_id: string;
    media_kind: "image" | "video";
    playback_url: string;
    thumbnail_url: string;
    mime_type: string;
    duration_seconds: number;
    size_bytes: number;
    processing_status: "pending" | "ready" | "failed";
    created_at: string;
    updated_at: string;
  }[];
  created_at: string;
}

export interface CreatePostRequest {
  persona_id?: string;
  posting_mode: "public" | "anonymous";
  event_id?: string;
  body: string;
  post_type: "text" | "image" | "video" | "set" | "event_tag";
  media_asset_ids?: string[];
  media_url?: string;
  media_type?: "image" | "youtube" | "soundcloud";
  location?: string;
}
