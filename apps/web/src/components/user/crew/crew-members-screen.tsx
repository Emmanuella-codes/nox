"use client";

import { ChevronLeft, MapPin, MessageCircle, Power, UserMinus } from "lucide-react";
import { useCallback, useEffect, useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import { FeedShell } from "@/src/components/user/feed/feed-shell";
import { Avatar } from "@/src/components/user/shared/avatar";
import type { Crew } from "@/src/types/api/crew";
import { endCrew, getCrew, leaveCrew } from "@/src/utils/api/user/crew";
import { getAccessToken } from "@/src/utils/auth/session";
import { useActivePersona } from "@/src/hooks/use-active-persona";

interface CrewMembersScreenProps { crewID: string }

export function CrewMembersScreen({ crewID }: CrewMembersScreenProps) {
  const router = useRouter();
  const { activeID, loading: personaLoading } = useActivePersona();
  const [crew, setCrew] = useState<Crew | null>(null);
  const [loading, setLoading] = useState(true);
  const [acting, setActing] = useState(false);
  const [message, setMessage] = useState("");
  const token = useMemo(() => getAccessToken(), []);

  const loadCrew = useCallback(async () => {
    if (!token || !activeID) return;
    setLoading(true);
    try {
      const res = await getCrew(crewID, activeID, token);
      setCrew(res.data ?? null);
    } catch {
      setMessage("Could not load crew members.");
    } finally {
      setLoading(false);
    }
  }, [activeID, crewID, token]);

  useEffect(() => {
    if (personaLoading) return;
    if (!token) {
      router.replace("/auth");
      return;
    }
    if (!activeID) {
      setMessage("Choose a persona before opening crew members.");
      setLoading(false);
      return;
    }
    void loadCrew();
  }, [activeID, loadCrew, personaLoading, router, token]);

  async function handleLeaveOrEnd() {
    if (!token || !activeID || !crew || acting) return;
    setActing(true);
    setMessage("");
    try {
      const currentMember = crew.members.find((member) => member.persona_id === activeID);
      if (currentMember?.role === "owner") {
        await endCrew(crewID, activeID, token);
      } else {
        await leaveCrew(crewID, activeID, token);
      }
      router.push(`/events/${crew.event_id}/crew`);
    } catch {
      setMessage("Could not update crew membership.");
    } finally {
      setActing(false);
    }
  }

  const members = crew?.members ?? [];
  const currentMember = members.find((member) => member.persona_id === activeID);
  const isOwner = currentMember?.role === "owner";

  return (
    <FeedShell>
      <header className="flex items-center gap-3 border-b border-(--nox-divider) px-4 pt-[env(safe-area-inset-top,12px)] pb-3">
        <button type="button" onClick={() => router.back()} className="flex size-8 items-center justify-center rounded-full hover:bg-(--nox-surface-alt)">
          <ChevronLeft className="size-5 text-(--nox-ink)" strokeWidth={1.8} />
        </button>
        <div>
          <h1 className="text-[18px] font-bold text-(--nox-ink)">crew members</h1>
          <p className="text-[11px] text-(--nox-ink-soft)">crew {crewID.slice(0, 8)} · {members.length} members</p>
        </div>
        {crew?.conversation_id ? (
          <button type="button" onClick={() => router.push(`/messages/${crew.conversation_id}`)}
            className="ml-auto flex items-center gap-1.5 rounded-[8px] border border-(--nox-border-strong) px-3 py-1.5 text-[12px] font-semibold text-(--nox-ink-mid)">
            <MessageCircle className="size-3.5" strokeWidth={1.7} />
            Chat
          </button>
        ) : null}
      </header>

      <div className="flex-1 overflow-y-auto">
        {loading && <p className="px-4 py-8 text-center text-[12px] text-(--nox-ink-soft)">Loading members...</p>}
        {message && <p className="mx-4 mt-4 rounded-[8px] border border-(--nox-danger) bg-(--nox-danger-soft) px-3 py-2 text-[12px] text-(--nox-danger)">{message}</p>}
        {members.map((member) => {
          const name = member.persona?.display_name ?? member.persona?.handle ?? "crew member";
          return (
          <div key={member.persona_id} className="flex items-center gap-3 border-b border-(--nox-divider) px-4 py-4">
            <Avatar id={member.persona_id} name={name} size={42} />
            <div className="min-w-0 flex-1">
              <div className="flex items-center gap-2">
                <p className="truncate text-[14px] font-semibold text-(--nox-ink)">{name}</p>
                {member.role === "owner" && (
                  <span className="shrink-0 rounded-full bg-(--nox-accent-soft) px-2 py-0.5 font-mono text-[8px] font-semibold text-(--nox-accent-ink)">host</span>
                )}
              </div>
              <p className="font-mono text-[10px] text-(--nox-ink-soft)">@{member.persona?.handle ?? "crew"}</p>
            </div>
            <div className="shrink-0 text-right">
              {member.location_sharing_enabled ? (
                <p className="flex items-center gap-1 text-[11px] text-(--nox-ink-soft)">
                  <MapPin className="size-3 text-(--nox-success)" strokeWidth={1.7} />
                  sharing
                </p>
              ) : (
                <p className="font-mono text-[10px] text-(--nox-ink-faint)">not sharing</p>
              )}
            </div>
          </div>
        )})}
        {crew && (
          <div className="px-4 py-5">
            <button type="button" onClick={() => void handleLeaveOrEnd()} disabled={acting}
              className="flex w-full items-center justify-center gap-2 rounded-[8px] border border-(--nox-danger) px-4 py-3 text-[13px] font-semibold text-(--nox-danger) disabled:opacity-50">
              {isOwner ? <Power className="size-4" strokeWidth={1.7} /> : <UserMinus className="size-4" strokeWidth={1.7} />}
              {acting ? "Updating..." : isOwner ? "End crew" : "Leave crew"}
            </button>
          </div>
        )}
        <p className="px-4 py-6 text-center text-[12px] text-(--nox-ink-soft)">Location sharing expires when the event window closes.</p>
      </div>
    </FeedShell>
  );
}
