import type {
  CreateCrewRequest,
  Crew,
  CrewLocation,
  CrewMember,
  JoinCrewRequest,
  UpdateLocationRequest,
  UpdateSharingRequest,
} from "@/src/types/api/crew";
import { apiRequest } from "@/src/utils/api/api";

export function createCrew(eventID: string, payload: CreateCrewRequest, token: string) {
  return apiRequest<Crew, CreateCrewRequest>(`/events/${eventID}/crews`, {
    method: "POST",
    body: payload,
    token,
  });
}

export function listMyEventCrews(eventID: string, personaID: string, token: string) {
  return apiRequest<Crew[]>(`/events/${eventID}/crews/me?persona_id=${personaID}`, { token });
}

export function joinCrew(payload: JoinCrewRequest, token: string) {
  return apiRequest<Crew, JoinCrewRequest>("/crews/join", {
    method: "POST",
    body: payload,
    token,
  });
}

export function getCrew(crewID: string, personaID: string, token: string) {
  return apiRequest<Crew>(`/crews/${crewID}?persona_id=${personaID}`, { token });
}

export function leaveCrew(crewID: string, personaID: string, token: string) {
  return apiRequest(`/crews/${crewID}/leave?persona_id=${personaID}`, {
    method: "POST",
    token,
  });
}

export function endCrew(crewID: string, personaID: string, token: string) {
  return apiRequest<Crew>(`/crews/${crewID}/end?persona_id=${personaID}`, {
    method: "POST",
    token,
  });
}

export function updateLocationSharing(crewID: string, payload: UpdateSharingRequest, token: string) {
  return apiRequest<CrewMember, UpdateSharingRequest>(`/crews/${crewID}/location-sharing`, {
    method: "PATCH",
    body: payload,
    token,
  });
}

export function updateCrewLocation(crewID: string, payload: UpdateLocationRequest, token: string) {
  return apiRequest<CrewLocation, UpdateLocationRequest>(`/crews/${crewID}/location`, {
    method: "PUT",
    body: payload,
    token,
  });
}

export function listCrewLocations(crewID: string, personaID: string, token: string) {
  return apiRequest<CrewLocation[]>(`/crews/${crewID}/locations?persona_id=${personaID}`, { token });
}
