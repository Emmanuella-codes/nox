import type { Post } from "@/src/types/api/user/post";

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
  limit: number;
  offset: number;
  has_more: boolean;
  next_offset?: number;
  posts: Post[];
}
