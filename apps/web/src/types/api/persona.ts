export interface Persona {
  id: string;
  handle: string;
  display_name: string;
  bio: string;
  avatar_url: string;
  cover_url: string;
  persona_type: "visible";
  category: "fan" | "dj" | "organizer" | "creator";
  genre_tags: string[];
  follower_count: number;
  following_count: number;
  is_following?: boolean;
  post_count: number;
  created_at: string;
  updated_at: string;
}

export interface FollowStatus {
  is_following: boolean;
}

export interface FollowListResponse {
  limit: number;
  offset: number;
  has_more: boolean;
  next_offset?: number;
  personas: Persona[];
}

export interface CreatePersonaRequest {
  handle: string;
  display_name: string;
  bio: string;
  avatar_url: string;
  cover_url: string;
  persona_type: "visible";
  category: "socialite" | "dj" | "organizer" | "creator";
  genre_tags: string[];
}
