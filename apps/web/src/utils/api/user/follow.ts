import type { FollowListResponse, FollowStatus } from "@/src/types/api/persona";
import { apiRequest } from "@/src/utils/api/api";

interface FollowPayload {
  persona_id: string;
}

export function followPersona(targetPersonaID: string, viewerPersonaID: string, token: string) {
  return apiRequest<null, FollowPayload>(`/personas/${targetPersonaID}/follow`, {
    method: "POST",
    body: { persona_id: viewerPersonaID },
    token,
  });
}

export function unfollowPersona(targetPersonaID: string, viewerPersonaID: string, token: string) {
  return apiRequest<null, FollowPayload>(`/personas/${targetPersonaID}/follow`, {
    method: "DELETE",
    body: { persona_id: viewerPersonaID },
    token,
  });
}

export function getFollowStatus(targetPersonaID: string, viewerPersonaID: string, token: string) {
  const params = new URLSearchParams({ persona_id: viewerPersonaID });
  return apiRequest<FollowStatus>(`/personas/${targetPersonaID}/follow-status?${params.toString()}`, {
    token,
  });
}

export function getFollowers(personaID: string, limit = 20, offset = 0) {
  const params = new URLSearchParams({ limit: String(limit), offset: String(offset) });
  return apiRequest<FollowListResponse>(`/personas/${personaID}/followers?${params.toString()}`);
}

export function getFollowing(personaID: string, limit = 20, offset = 0) {
  const params = new URLSearchParams({ limit: String(limit), offset: String(offset) });
  return apiRequest<FollowListResponse>(`/personas/${personaID}/following?${params.toString()}`);
}
