export interface CrewPersona {
  id: string;
  handle: string;
  display_name: string;
  avatar_url: string;
}

export interface CrewMember {
  persona_id: string;
  role: "owner" | "member";
  location_sharing_enabled: boolean;
  persona?: CrewPersona;
  joined_at: string;
}

export interface CrewLocation {
  persona_id: string;
  latitude: number;
  longitude: number;
  accuracy_meters: number;
  battery_level?: number;
  recorded_at: string;
  expires_at: string;
  persona?: CrewPersona;
}

export interface Crew {
  id: string;
  event_id: string;
  conversation_id: string;
  owner_persona_id: string;
  name: string;
  join_code: string;
  visibility: "private" | "invite_code";
  status: "active" | "ended";
  expires_at: string;
  members: CrewMember[];
  locations: CrewLocation[];
  created_at: string;
  updated_at: string;
}

export interface CreateCrewRequest {
  owner_persona_id: string;
  name: string;
  visibility: "private" | "invite_code";
}

export interface JoinCrewRequest {
  persona_id: string;
  join_code: string;
}

export interface UpdateSharingRequest {
  persona_id: string;
  enabled: boolean;
}

export interface UpdateLocationRequest {
  persona_id: string;
  latitude: number;
  longitude: number;
  accuracy_meters: number;
  battery_level?: number;
}
