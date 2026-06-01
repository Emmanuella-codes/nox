import type { CreatePersonaRequest, Persona } from "@/src/types/api/persona";
import { apiRequest } from "@/src/utils/api/api";

export function createPersona(payload: CreatePersonaRequest, token: string) {
  return apiRequest<Persona, CreatePersonaRequest>("/personas", {
    method: "POST",
    body: payload,
    token,
  });
}

export function getMyPersonas(token: string) {
  return apiRequest<Persona[]>("/personas/me", { token });
}

export function getPersona(personaID: string) {
  return apiRequest<Persona>(`/personas/${personaID}`);
}
