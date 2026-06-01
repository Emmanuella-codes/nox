import type { Event } from "@/src/types/api/event";
import type { Hashtag } from "@/src/types/api/hashtag";
import type { Persona } from "@/src/types/api/persona";
import type { Post } from "@/src/types/api/post";

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
