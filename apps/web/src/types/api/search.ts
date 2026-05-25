import type { Event } from "@/src/types/api/event";
import type { Persona } from "@/src/types/api/persona";
import type { Post } from "@/src/types/api/post";

export interface SearchResponse {
  query: string;
  personas: Persona[];
  posts: Post[];
  events: Event[];
}
