import type { SearchResponse } from "@/src/types/api/search";
import { apiRequest } from "@/src/utils/api/api";

export function searchNox(query: string, limit = 10, token?: string, viewerPersonaID?: string) {
	const params = new URLSearchParams({
		q: query,
		limit: String(limit),
	});
	if (viewerPersonaID) {
		params.set("viewer_persona_id", viewerPersonaID);
	}
	return apiRequest<SearchResponse>(`/search?${params.toString()}`, { token });
}
