import type { SearchResponse } from "@/src/types/api/user/search";
import { apiRequest } from "@/src/utils/api/api";

export function searchNox(
	query: string,
	limit = 10,
	token?: string,
	viewerPersonaID?: string,
	offset = 0,
) {
	const params = new URLSearchParams({
		q: query,
		limit: String(limit),
		offset: String(offset),
	});
	if (viewerPersonaID) {
		params.set("viewer_persona_id", viewerPersonaID);
	}
	return apiRequest<SearchResponse>(`/search?${params.toString()}`, { token });
}
