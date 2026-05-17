import type { Event } from "@/src/types/api/event";
import { apiRequest } from "@/src/utils/api/api";

export function getEvents() {
  return apiRequest<Event[]>("/events");
}
