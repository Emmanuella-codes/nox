"use client";

import { ChevronLeft, LocateFixed, MapPin, MessageCircle, Users } from "lucide-react";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { useRouter } from "next/navigation";
import { FeedShell } from "@/src/components/user/feed/feed-shell";
import { Avatar } from "@/src/components/user/shared/avatar";
import type { Crew, CrewLocation } from "@/src/types/api/user/crew";
import { getCrew, listCrewLocations, updateCrewLocation, updateLocationSharing } from "@/src/utils/api/user/crew";
import { getAccessToken } from "@/src/utils/auth/session";
import { useActivePersona } from "@/src/hooks/use-active-persona";

interface CrewMapScreenProps { crewID: string }

function crewClosed(crew: Crew | null) {
  return !crew || crew.status === "ended" || new Date(crew.expires_at).getTime() <= Date.now();
}

function mapBounds(locations: CrewLocation[]) {
  if (locations.length === 0) return null;
  const lats = locations.map((location) => location.latitude);
  const lngs = locations.map((location) => location.longitude);
  const minLat = Math.min(...lats);
  const maxLat = Math.max(...lats);
  const minLng = Math.min(...lngs);
  const maxLng = Math.max(...lngs);
  const latPad = Math.max((maxLat - minLat) * 0.25, 0.002);
  const lngPad = Math.max((maxLng - minLng) * 0.25, 0.002);
  return {
    minLat: minLat - latPad,
    maxLat: maxLat + latPad,
    minLng: minLng - lngPad,
    maxLng: maxLng + lngPad,
  };
}

function pinPosition(location: CrewLocation, bounds: NonNullable<ReturnType<typeof mapBounds>>) {
  const lngRange = bounds.maxLng - bounds.minLng || 1;
  const latRange = bounds.maxLat - bounds.minLat || 1;
  return {
    top: `${((bounds.maxLat - location.latitude) / latRange) * 100}%`,
    left: `${((location.longitude - bounds.minLng) / lngRange) * 100}%`,
  };
}

function mapURL(bounds: NonNullable<ReturnType<typeof mapBounds>>) {
  const markerLat = (bounds.minLat + bounds.maxLat) / 2;
  const markerLng = (bounds.minLng + bounds.maxLng) / 2;
  const bbox = `${bounds.minLng},${bounds.minLat},${bounds.maxLng},${bounds.maxLat}`;
  return `https://www.openstreetmap.org/export/embed.html?bbox=${encodeURIComponent(bbox)}&layer=mapnik&marker=${markerLat},${markerLng}`;
}

export function CrewMapScreen({ crewID }: CrewMapScreenProps) {
  const router = useRouter();
  const { activeID, loading: personaLoading } = useActivePersona();
  const [crew, setCrew] = useState<Crew | null>(null);
  const [locations, setLocations] = useState<CrewLocation[]>([]);
  const [loading, setLoading] = useState(true);
  const [sharing, setSharing] = useState(false);
  const [message, setMessage] = useState("");
  const [locating, setLocating] = useState(false);
  const watchRef = useRef<number | null>(null);
  const token = useMemo(() => getAccessToken(), []);

  const loadCrew = useCallback(async () => {
    if (!token || !activeID) return;
    setLoading(true);
    try {
      const res = await getCrew(crewID, activeID, token);
      const loaded = res.data ?? null;
      setCrew(loaded);
      setLocations(loaded?.locations ?? []);
      setSharing(Boolean(loaded?.members.find((m) => m.persona_id === activeID)?.location_sharing_enabled));
    } catch {
      setMessage("Could not load crew.");
    } finally {
      setLoading(false);
    }
  }, [activeID, crewID, token]);

  const loadLocations = useCallback(async () => {
    if (!token || !activeID) return;
    try {
      const res = await listCrewLocations(crewID, activeID, token);
      setLocations(res.data ?? []);
    } catch {
      setMessage("Could not refresh locations.");
    }
  }, [activeID, crewID, token]);

  useEffect(() => {
    if (personaLoading) return;
    if (!token) {
      router.replace("/auth");
      return;
    }
    if (!activeID) {
      setMessage("Choose a persona before opening crew mode.");
      setLoading(false);
      return;
    }
    void loadCrew();
  }, [activeID, loadCrew, personaLoading, router, token]);

  useEffect(() => {
    if (!token || !activeID) return;
    const interval = window.setInterval(() => void loadLocations(), 5000);
    return () => window.clearInterval(interval);
  }, [activeID, loadLocations, token]);

  useEffect(() => {
    return () => {
      if (watchRef.current !== null) navigator.geolocation.clearWatch(watchRef.current);
    };
  }, []);

  async function handleSharingToggle() {
    if (!token || !activeID) return;
    setMessage("");
    if (sharing) {
      try {
        await updateLocationSharing(crewID, { persona_id: activeID, enabled: false }, token);
        setSharing(false);
        if (watchRef.current !== null) {
          navigator.geolocation.clearWatch(watchRef.current);
          watchRef.current = null;
        }
      } catch {
        setMessage("Could not update location sharing.");
      }
      return;
    }
    if (closed) {
      setMessage("This crew is closed.");
      return;
    }
    if (!("geolocation" in navigator)) {
      setMessage("Location is not available in this browser.");
      return;
    }
    setLocating(true);
    navigator.geolocation.getCurrentPosition(
      (position) => {
        void enableSharing(position);
      },
      () => {
        setLocating(false);
        setMessage("Allow location access to share with your crew.");
      },
      { enableHighAccuracy: true, maximumAge: 10000, timeout: 15000 },
    );
  }

  async function enableSharing(position: GeolocationPosition) {
    if (!token || !activeID) return;
    try {
      await updateLocationSharing(crewID, { persona_id: activeID, enabled: true }, token);
      const res = await updateCrewLocation(
        crewID,
        {
          persona_id: activeID,
          latitude: position.coords.latitude,
          longitude: position.coords.longitude,
          accuracy_meters: position.coords.accuracy,
        },
        token,
      );
      if (res.data) {
        setLocations((current) => [res.data as CrewLocation, ...current.filter((item) => item.persona_id !== activeID)]);
      }
      setSharing(true);
      startLocationWatch();
    } catch {
      setMessage("Could not update location sharing.");
    } finally {
      setLocating(false);
    }
  }

  function startLocationWatch() {
    if (!token || !activeID) return;
    if (!("geolocation" in navigator)) {
      setMessage("Location is not available in this browser.");
      return;
    }
    watchRef.current = navigator.geolocation.watchPosition(
      (position) => {
        void updateCrewLocation(
          crewID,
          {
            persona_id: activeID,
            latitude: position.coords.latitude,
            longitude: position.coords.longitude,
            accuracy_meters: position.coords.accuracy,
          },
          token,
        ).then((res) => {
          if (res.data) setLocations((current) => [res.data as CrewLocation, ...current.filter((item) => item.persona_id !== activeID)]);
        }).catch(() => setMessage("Could not share your current location."));
      },
      () => setMessage("Allow location access to share with your crew."),
      { enableHighAccuracy: true, maximumAge: 10000, timeout: 15000 },
    );
  }

  const memberCount = crew?.members.length ?? 0;
  const closed = crewClosed(crew);
  const bounds = mapBounds(locations);

  return (
    <FeedShell>
      <header className="flex items-center gap-3 border-b border-(--nox-divider) px-4 pt-[env(safe-area-inset-top,12px)] pb-3">
        <button type="button" onClick={() => router.back()} className="flex size-8 items-center justify-center rounded-full hover:bg-(--nox-surface-alt)">
          <ChevronLeft className="size-5 text-(--nox-ink)" strokeWidth={1.8} />
        </button>
        <h1 className="text-[18px] font-bold text-(--nox-ink)">live map</h1>
        {crew?.conversation_id && !closed ? (
          <button type="button" onClick={() => router.push(`/messages/${crew.conversation_id}`)}
            className="ml-auto flex items-center gap-1.5 rounded-[8px] border border-(--nox-border-strong) px-3 py-1.5 text-[12px] font-semibold text-(--nox-ink-mid) hover:border-(--nox-accent)">
            <MessageCircle className="size-3.5" strokeWidth={1.7} />
            Chat
          </button>
        ) : null}
        {!closed && <button type="button" onClick={() => router.push(`/crews/${crewID}/members`)}
          className="flex items-center gap-1.5 rounded-[8px] border border-(--nox-border-strong) px-3 py-1.5 text-[12px] font-semibold text-(--nox-ink-mid) hover:border-(--nox-accent)">
          <Users className="size-3.5" strokeWidth={1.7} />
          {memberCount}
        </button>}
      </header>

      <div className="relative mx-4 mt-4 h-64 overflow-hidden rounded-[12px] bg-(--nox-surface-alt)">
        {bounds ? (
          <iframe title="Crew live map" src={mapURL(bounds)} className="absolute inset-0 size-full border-0 grayscale-[0.2]" loading="lazy" />
        ) : (
          <div className="absolute inset-0 bg-[linear-gradient(90deg,rgba(255,255,255,0.06)_1px,transparent_1px),linear-gradient(rgba(255,255,255,0.06)_1px,transparent_1px)] bg-size-[36px_36px]" />
        )}
        <div className="pointer-events-none absolute inset-0 bg-black/10" />
        {!bounds && (
          <div className="absolute inset-0 flex flex-col items-center justify-center gap-2">
            <MapPin className="size-8 text-(--nox-accent)" strokeWidth={1.5} />
            <p className="text-[13px] text-(--nox-ink-soft)">{closed ? "Crew is closed" : "No shared locations yet"}</p>
            <p className="text-[11px] text-(--nox-ink-faint)">Pins refresh every few seconds.</p>
          </div>
        )}
        {bounds && locations.map((location) => (
          <div key={location.persona_id} className="absolute flex -translate-x-1/2 -translate-y-1/2 flex-col items-center" style={pinPosition(location, bounds)}>
            <Avatar id={location.persona_id} name={location.persona?.display_name ?? location.persona?.handle ?? "crew"} size={28} />
            <div className="mt-1 rounded-[4px] bg-black/60 px-1.5 py-0.5">
              <p className="font-mono text-[8px] text-white/80">@{location.persona?.handle ?? "crew"}</p>
            </div>
          </div>
        ))}
      </div>

      <div className="mx-4 mt-3 rounded-[12px] border border-(--nox-border) bg-(--nox-surface) p-4">
        <div className="flex items-center justify-between gap-3">
          <div>
            <p className="text-[14px] font-bold text-(--nox-ink)">{crew?.name ?? "Crew mode"}</p>
            <p className="mt-0.5 text-[11px] text-(--nox-ink-soft)">
              {closed ? "This crew is closed. Location sharing is off." : "Share your live location during the event window."}
            </p>
          </div>
          {!closed && (
            <button type="button" onClick={() => void handleSharingToggle()} disabled={locating}
              className={`rounded-[8px] px-3 py-2 text-[12px] font-semibold text-white disabled:opacity-50 ${sharing ? "bg-(--nox-danger)" : "bg-(--nox-accent)"}`}>
              {locating ? "Locating" : sharing ? "Stop" : "Share"}
            </button>
          )}
        </div>
        {message && <p className="mt-3 text-[12px] text-(--nox-danger)">{message}</p>}
      </div>

      <div className="mt-4 flex-1 overflow-y-auto px-4">
        <p className="mb-2 font-mono text-[10px] font-semibold uppercase tracking-[0.18em] text-(--nox-ink-soft)">crew nearby</p>
        {loading && <p className="py-6 text-center text-[12px] text-(--nox-ink-soft)">Loading crew...</p>}
        {!loading && locations.length === 0 && <p className="py-6 text-center text-[12px] text-(--nox-ink-soft)">No one is sharing location.</p>}
        {locations.map((location) => (
          <div key={location.persona_id} className="flex items-center gap-3 border-b border-(--nox-divider) py-3">
            <Avatar id={location.persona_id} name={location.persona?.display_name ?? location.persona?.handle ?? "crew"} size={36} />
            <div className="flex-1">
              <p className="text-[13px] font-semibold text-(--nox-ink)">{location.persona?.display_name ?? location.persona?.handle ?? "crew member"}</p>
              <p className="flex items-center gap-1 font-mono text-[10px] text-(--nox-ink-soft)">
                <LocateFixed className="size-2.5" strokeWidth={1.7} /> {Math.round(location.accuracy_meters)}m accuracy
              </p>
            </div>
          </div>
        ))}
      </div>
    </FeedShell>
  );
}
