"use client";

import { ChevronLeft, LocateFixed, MapPin, MessageCircle, Users } from "lucide-react";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { useRouter } from "next/navigation";
import { FeedShell } from "@/src/components/user/feed/feed-shell";
import { Avatar } from "@/src/components/user/shared/avatar";
import type { Crew, CrewLocation } from "@/src/types/api/crew";
import { getCrew, listCrewLocations, updateCrewLocation, updateLocationSharing } from "@/src/utils/api/user/crew";
import { getAccessToken } from "@/src/utils/auth/session";
import { useActivePersona } from "@/src/hooks/use-active-persona";

interface CrewMapScreenProps { crewID: string }

function pinPosition(location: CrewLocation, index: number) {
  const latSeed = Math.abs(location.latitude * 1000);
  const lngSeed = Math.abs(location.longitude * 1000);
  return {
    top: `${16 + ((latSeed + index * 17) % 66)}%`,
    left: `${14 + ((lngSeed + index * 19) % 70)}%`,
  };
}

export function CrewMapScreen({ crewID }: CrewMapScreenProps) {
  const router = useRouter();
  const { activeID, loading: personaLoading } = useActivePersona();
  const [crew, setCrew] = useState<Crew | null>(null);
  const [locations, setLocations] = useState<CrewLocation[]>([]);
  const [loading, setLoading] = useState(true);
  const [sharing, setSharing] = useState(false);
  const [message, setMessage] = useState("");
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
    const enabled = !sharing;
    setMessage("");
    try {
      await updateLocationSharing(crewID, { persona_id: activeID, enabled }, token);
      setSharing(enabled);
      if (!enabled && watchRef.current !== null) {
        navigator.geolocation.clearWatch(watchRef.current);
        watchRef.current = null;
      }
      if (enabled) startLocationWatch();
    } catch {
      setMessage("Could not update location sharing.");
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

  return (
    <FeedShell>
      <header className="flex items-center gap-3 border-b border-(--nox-divider) px-4 pt-[env(safe-area-inset-top,12px)] pb-3">
        <button type="button" onClick={() => router.back()} className="flex size-8 items-center justify-center rounded-full hover:bg-(--nox-surface-alt)">
          <ChevronLeft className="size-5 text-(--nox-ink)" strokeWidth={1.8} />
        </button>
        <h1 className="text-[18px] font-bold text-(--nox-ink)">live map</h1>
        {crew?.conversation_id ? (
          <button type="button" onClick={() => router.push(`/messages/${crew.conversation_id}`)}
            className="ml-auto flex items-center gap-1.5 rounded-[8px] border border-(--nox-border-strong) px-3 py-1.5 text-[12px] font-semibold text-(--nox-ink-mid) hover:border-(--nox-accent)">
            <MessageCircle className="size-3.5" strokeWidth={1.7} />
            Chat
          </button>
        ) : null}
        <button type="button" onClick={() => router.push(`/crews/${crewID}/members`)}
          className="flex items-center gap-1.5 rounded-[8px] border border-(--nox-border-strong) px-3 py-1.5 text-[12px] font-semibold text-(--nox-ink-mid) hover:border-(--nox-accent)">
          <Users className="size-3.5" strokeWidth={1.7} />
          {memberCount}
        </button>
      </header>

      <div className="relative mx-4 mt-4 h-64 overflow-hidden rounded-[12px] bg-(--nox-surface-alt)">
        <div className="absolute inset-0 bg-[linear-gradient(90deg,rgba(255,255,255,0.06)_1px,transparent_1px),linear-gradient(rgba(255,255,255,0.06)_1px,transparent_1px)] bg-size-[36px_36px]" />
        <div className="absolute inset-0 flex flex-col items-center justify-center gap-2">
          <MapPin className="size-8 text-(--nox-accent)" strokeWidth={1.5} />
          <p className="text-[13px] text-(--nox-ink-soft)">{locations.length ? "Crew pins are live" : "No shared locations yet"}</p>
          <p className="text-[11px] text-(--nox-ink-faint)">Pins refresh every few seconds.</p>
        </div>
        {locations.map((location, index) => (
          <div key={location.persona_id} className="absolute flex flex-col items-center" style={pinPosition(location, index)}>
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
            <p className="mt-0.5 text-[11px] text-(--nox-ink-soft)">Share your live location during the event window.</p>
          </div>
          <button type="button" onClick={() => void handleSharingToggle()}
            className={`rounded-[8px] px-3 py-2 text-[12px] font-semibold text-white ${sharing ? "bg-(--nox-danger)" : "bg-(--nox-accent)"}`}>
            {sharing ? "Stop" : "Share"}
          </button>
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
