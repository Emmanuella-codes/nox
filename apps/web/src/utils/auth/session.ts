export function getAccessToken() {
  if (typeof window === "undefined") {
    return "";
  }

  return localStorage.getItem("nox_access_token") ?? "";
}

export function hasAccessToken() {
  return Boolean(getAccessToken());
}

export function getActivePersonaID() {
  if (typeof window === "undefined") {
    return "";
  }

  return localStorage.getItem("nox_active_persona_id") ?? "";
}

export function setActivePersonaID(personaID: string) {
  if (typeof window === "undefined") {
    return;
  }

  localStorage.setItem("nox_active_persona_id", personaID);
}

export function clearSession() {
  if (typeof window === "undefined") {
    return;
  }

  localStorage.removeItem("nox_access_token");
  localStorage.removeItem("nox_active_persona_id");
  localStorage.removeItem("nox_signup_category");
}
