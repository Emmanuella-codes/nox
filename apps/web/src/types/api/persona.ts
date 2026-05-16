export interface Persona {
  id: string;
  handle: string;
  display_name: string;
  bio: string;
  avatar_url: string;
  cover_url: string;
  persona_type: "visible";
  genre_tags: string[];
  follower_count: number;
  following_count: number;
  post_count: number;
  created_at: string;
  updated_at: string;
}

export interface CreatePersonaRequest {
  handle: string;
  display_name: string;
  bio: string;
  avatar_url: string;
  cover_url: string;
  persona_type: "visible";
  genre_tags: string[];
}
