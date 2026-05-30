import type { Post } from "@/src/types/api/post";

export interface Hashtag {
  id: string;
  tag: string;
  post_count: number;
  created_at: string;
}

export interface HashtagDetail {
  tag: string;
  post_count: number;
}

export interface HashtagPostsResponse {
  tag: string;
  posts: Post[];
}
