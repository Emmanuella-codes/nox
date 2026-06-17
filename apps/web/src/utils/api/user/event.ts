import type { CreateEventRequest, Event } from "@/src/types/api/event";
import { apiRequest } from "@/src/utils/api/api";

export function getEvents() {
  return apiRequest<Event[]>("/events");
}

export function getEvent(eventID: string) {
  return apiRequest<Event>(`/events/${eventID}`);
}

export function createEvent(payload: CreateEventRequest, token: string) {
  return apiRequest<Event, CreateEventRequest>("/events", { method: "POST", body: payload, token });
}
