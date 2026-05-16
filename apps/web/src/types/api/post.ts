export interface PostAuthor {
  mode: "public" | "anonymous";
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
  post_type: "text" | "image" | "set" | "event_tag";
  media_url?: string;
  media_type?: "image" | "youtube" | "soundcloud";
  location?: string;
  like_count: number;
  comment_count: number;
  repost_count: number;
  is_repost: boolean;
  repost_of?: string;
  created_at: string;
}

export interface CreatePostRequest {
  persona_id?: string;
  posting_mode: "public" | "anonymous";
  event_id?: string;
  body: string;
  post_type: "text" | "image" | "set" | "event_tag";
  media_url?: string;
  media_type?: "image" | "youtube" | "soundcloud";
  location?: string;
}
