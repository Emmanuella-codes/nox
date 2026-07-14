import type {
  CreateMediaAssetRequest,
  CreateSetRequest,
  MediaAsset,
  Set,
  SetListResponse,
} from "@/src/types/api/user/set";
import { apiRequest } from "@/src/utils/api/api";

export function createSetVideoAsset(payload: CreateMediaAssetRequest, token: string) {
  return apiRequest<MediaAsset, CreateMediaAssetRequest>("/media-assets/set-video", {
    method: "POST",
    body: payload,
    token,
  });
}

export function createSet(payload: CreateSetRequest, token: string) {
  return apiRequest<Set, CreateSetRequest>("/sets", {
    method: "POST",
    body: payload,
    token,
  });
}

export function getSets(limit = 20, offset = 0) {
  return apiRequest<SetListResponse>(`/sets?limit=${encodeURIComponent(String(limit))}&offset=${encodeURIComponent(String(offset))}`);
}

export function getPersonaSets(personaID: string, limit = 20, offset = 0) {
  return apiRequest<SetListResponse>(
    `/sets/persona/${encodeURIComponent(personaID)}?limit=${encodeURIComponent(String(limit))}&offset=${encodeURIComponent(String(offset))}`,
  );
}

export function getSet(setID: string) {
  return apiRequest<Set>(`/sets/${encodeURIComponent(setID)}`);
}

export function deleteSet(setID: string, token: string) {
  return apiRequest<null>(`/sets/${encodeURIComponent(setID)}`, {
    method: "DELETE",
    token,
  });
}
