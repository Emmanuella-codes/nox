export function getAccessToken() {
  if (typeof window === "undefined") {
    return "";
  }

  return localStorage.getItem("nox_access_token") ?? "";
}

export function hasAccessToken() {
  return Boolean(getAccessToken());
}
