import type { Event } from "@/src/types/api/user/event";
import type { Hashtag } from "@/src/types/api/user/hashtag";
import type { Persona } from "@/src/types/api/user/persona";
import type { Post } from "@/src/types/api/user/post";

export interface SearchResponse {
  query: string;
  limit: number;
  offset: number;
  has_more: boolean;
  next_offset?: number;
  personas: Persona[];
  posts: Post[];
  events: Event[];
  hashtags?: Hashtag[];
}
