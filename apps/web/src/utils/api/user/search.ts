import type { SearchResponse } from "@/src/types/api/search";
import { apiRequest } from "@/src/utils/api/api";

export function searchNox(query: string, limit = 10) {
  const params = new URLSearchParams({
    q: query,
    limit: String(limit),
  });
  return apiRequest<SearchResponse>(`/search?${params.toString()}`);
}
