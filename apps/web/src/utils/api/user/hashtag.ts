import type { Hashtag, HashtagDetail, HashtagPostsResponse } from "@/src/types/api/user/hashtag";
import { apiRequest } from "@/src/utils/api/api";

export function getTrendingHashtags(limit = 20) {
  return apiRequest<Hashtag[]>(`/hashtags/trending?limit=${encodeURIComponent(String(limit))}`);
}

export function getHashtag(tag: string) {
  return apiRequest<HashtagDetail>(`/hashtags/${encodeURIComponent(tag)}`);
}

export function getHashtagPosts(
  tag: string,
  limit = 20,
  offset = 0,
  viewerPersonaID?: string,
  token?: string,
) {
  const params = new URLSearchParams({ limit: String(limit), offset: String(offset) });
  if (viewerPersonaID) params.set("viewer_persona_id", viewerPersonaID);
  return apiRequest<HashtagPostsResponse>(
    `/hashtags/${encodeURIComponent(tag)}/posts?${params.toString()}`,
    { token },
  );
}
