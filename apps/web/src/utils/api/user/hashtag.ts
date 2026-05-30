import type { Hashtag, HashtagDetail, HashtagPostsResponse } from "@/src/types/api/hashtag";
import { apiRequest } from "@/src/utils/api/api";

export function getTrendingHashtags(limit = 20) {
  return apiRequest<Hashtag[]>(`/hashtags/trending?limit=${encodeURIComponent(String(limit))}`);
}

export function getHashtag(tag: string) {
  return apiRequest<HashtagDetail>(`/hashtags/${encodeURIComponent(tag)}`);
}

export function getHashtagPosts(tag: string, limit = 20) {
  const params = new URLSearchParams({ limit: String(limit) });
  return apiRequest<HashtagPostsResponse>(
    `/hashtags/${encodeURIComponent(tag)}/posts?${params.toString()}`,
  );
}
