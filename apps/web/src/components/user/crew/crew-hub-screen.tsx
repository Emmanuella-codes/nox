"use client";

import type { FormEvent } from "react";
import { Check, ChevronLeft, Copy, Map, MessageCircle, Plus, Users } from "lucide-react";
import { useCallback, useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { FeedShell } from "@/src/components/user/feed/feed-shell";
import { TabBar } from "@/src/components/user/feed/tab-bar";
import { useActivePersona } from "@/src/hooks/use-active-persona";
import type { Crew } from "@/src/types/api/crew";
import { ApiRequestError } from "@/src/utils/api/api";
import { createCrew, joinCrew, listMyEventCrews } from "@/src/utils/api/user/crew";
import { getAccessToken } from "@/src/utils/auth/session";

interface CrewHubScreenProps { eventID: string }

function crewClosed(crew: Crew) {
  return crew.status === "ended" || new Date(crew.expires_at).getTime() <= Date.now();
}

export function CrewHubScreen({ eventID }: CrewHubScreenProps) {
  const router = useRouter();
  const { activeID, activePersona, loading: personaLoading } = useActivePersona();
  const [crews, setCrews] = useState<Crew[]>([]);
  const [crewName, setCrewName] = useState("");
  const [joinCode, setJoinCode] = useState("");
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState<"create" | "join" | "">("");
  const [copiedCrewID, setCopiedCrewID] = useState("");
  const [message, setMessage] = useState("");

  const loadCrews = useCallback(async () => {
    const token = getAccessToken();
    if (!token || !activeID) {
      setLoading(false);
      return;
    }
    setLoading(true);
    setMessage("");
    try {
      const res = await listMyEventCrews(eventID, activeID, token);
      setCrews(res.data ?? []);
    } catch (error) {
      setMessage(error instanceof ApiRequestError ? error.message : "Could not load crews.");
    } finally {
      setLoading(false);
    }
  }, [activeID, eventID]);

  useEffect(() => {
    if (personaLoading) return;
    if (!getAccessToken()) {
      router.replace("/auth");
      return;
    }
    void loadCrews();
  }, [loadCrews, personaLoading, router]);

  async function handleCreate(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const token = getAccessToken();
    if (!token || !activeID || submitting) return;
    setSubmitting("create");
    setMessage("");
    try {
      const res = await createCrew(
        eventID,
        {
          owner_persona_id: activeID,
          name: crewName.trim() || "Night crew",
          visibility: "invite_code",
        },
        token,
      );
      if (res.data?.id) router.push(`/crews/${res.data.id}`);
    } catch (error) {
      setMessage(error instanceof ApiRequestError ? error.message : "Could not create crew.");
    } finally {
      setSubmitting("");
    }
  }

  async function handleJoin() {
    const token = getAccessToken();
    if (!token || !activeID || joinCode.length < 6 || submitting) return;
    setSubmitting("join");
    setMessage("");
    try {
      const res = await joinCrew({ persona_id: activeID, join_code: joinCode }, token);
      if (res.data?.id) router.push(`/crews/${res.data.id}`);
    } catch (error) {
      setMessage(error instanceof ApiRequestError ? error.message : "Could not join crew.");
    } finally {
      setSubmitting("");
    }
  }

  async function handleCopyCode(crew: Crew) {
    try {
      await navigator.clipboard.writeText(crew.join_code);
      setCopiedCrewID(crew.id);
      window.setTimeout(() => setCopiedCrewID(""), 1800);
    } catch {
      setMessage("Could not copy invite code.");
    }
  }

  return (
    <FeedShell>
      <header className="flex items-center gap-3 border-b border-(--nox-divider) px-4 pt-[env(safe-area-inset-top,12px)] pb-3">
        <button type="button" onClick={() => router.back()} className="flex size-8 items-center justify-center rounded-full hover:bg-(--nox-surface-alt)">
          <ChevronLeft className="size-5 text-(--nox-ink)" strokeWidth={1.8} />
        </button>
        <div>
          <h1 className="text-[18px] font-bold text-(--nox-ink)">crew hub</h1>
          <p className="text-[11px] text-(--nox-ink-soft)">event {eventID.slice(0, 8)} · @{activePersona?.handle ?? "guest"}</p>
        </div>
      </header>

      <div className="flex-1 overflow-y-auto px-4 py-6 space-y-4">
        <form onSubmit={handleCreate} className="rounded-[12px] border border-(--nox-border) bg-(--nox-surface) p-5">
          <div className="flex items-center gap-4">
            <span className="flex size-10 items-center justify-center rounded-[10px] bg-(--nox-accent-soft)">
              <Plus className="size-5 text-(--nox-accent-ink)" strokeWidth={1.8} />
            </span>
            <div>
              <p className="text-[15px] font-bold text-(--nox-ink)">Create a crew</p>
              <p className="mt-0.5 text-[12px] text-(--nox-ink-soft)">Share a code with friends at this event.</p>
            </div>
          </div>
          <div className="mt-4 flex gap-2">
            <input
              value={crewName}
              onChange={(e) => setCrewName(e.target.value)}
              placeholder="Night crew"
              className="flex-1 rounded-[8px] border border-(--nox-border) bg-(--nox-surface-alt) px-3 py-2.5 text-[14px] text-(--nox-ink) outline-none focus:border-(--nox-accent-line)"
            />
            <button type="submit" disabled={!activeID || submitting === "create"}
              className="rounded-[8px] bg-(--nox-accent) px-4 py-2.5 text-[13px] font-semibold text-white disabled:opacity-40">
              {submitting === "create" ? "Creating" : "Create"}
            </button>
          </div>
        </form>

        <div className="rounded-[12px] border border-(--nox-border) bg-(--nox-surface) p-5">
          <p className="text-[15px] font-bold text-(--nox-ink)">Join with a code</p>
          <p className="mt-0.5 text-[12px] text-(--nox-ink-soft)">Got an invite? Enter the 6-character code</p>
          <div className="mt-3 flex gap-2">
            <input value={joinCode} onChange={(e) => setJoinCode(e.target.value.toUpperCase())}
              maxLength={6} placeholder="ABC123"
              className="flex-1 rounded-[8px] border border-(--nox-border) bg-(--nox-surface-alt) px-3 py-2.5 font-mono text-[15px] tracking-[0.2em] text-(--nox-ink) outline-none focus:border-(--nox-accent-line) uppercase" />
            <button type="button" onClick={() => void handleJoin()} disabled={joinCode.length < 6 || submitting === "join"}
              className="rounded-[8px] bg-(--nox-accent) px-4 py-2.5 text-[13px] font-semibold text-white disabled:opacity-40">
              {submitting === "join" ? "Joining" : "Join"}
            </button>
          </div>
        </div>

        {message && <p className="rounded-[8px] border border-(--nox-danger) bg-(--nox-danger-soft) px-3 py-2 text-[12px] text-(--nox-danger)">{message}</p>}

        <div className="space-y-3">
          <p className="font-mono text-[10px] font-semibold uppercase tracking-[0.18em] text-(--nox-ink-soft)">your event crews</p>
          {loading ? <p className="text-[12px] text-(--nox-ink-soft)">Loading crews...</p> : null}
          {!loading && crews.length === 0 ? <p className="text-[12px] text-(--nox-ink-soft)">Create or join a crew for this event.</p> : null}
          {crews.map((crew) => {
            const closed = crewClosed(crew);
            return (
            <div key={crew.id} className={`rounded-[12px] border p-4 ${closed ? "border-(--nox-divider) bg-(--nox-surface-alt)" : "border-(--nox-border) bg-(--nox-surface)"}`}>
              <div className="flex items-start justify-between gap-3">
                <div>
                  <p className="text-[15px] font-bold text-(--nox-ink)">{crew.name}</p>
                  <p className="mt-1 font-mono text-[10px] text-(--nox-ink-soft)">
                    {closed ? "crew closed" : `code ${crew.join_code}`} · {crew.members.length} members
                  </p>
                </div>
                <span className={`rounded-full px-2 py-1 font-mono text-[9px] font-semibold ${closed ? "bg-(--nox-surface) text-(--nox-ink-faint)" : "bg-(--nox-accent-soft) text-(--nox-accent-ink)"}`}>
                  {closed ? "closed" : "active"}
                </span>
              </div>
              {!closed && (
                <button type="button" onClick={() => void handleCopyCode(crew)}
                  className="mt-3 flex w-full items-center justify-center gap-2 rounded-[8px] border border-(--nox-border-strong) py-2 text-[12px] font-semibold text-(--nox-ink-mid)">
                  {copiedCrewID === crew.id ? <Check className="size-4" strokeWidth={1.6} /> : <Copy className="size-4" strokeWidth={1.6} />}
                  {copiedCrewID === crew.id ? "Copied" : "Copy invite code"}
                </button>
              )}
              <div className={`mt-3 grid ${closed ? "grid-cols-1" : "grid-cols-3"} gap-2`}>
                <button type="button" onClick={() => router.push(`/crews/${crew.id}`)}
                  className="flex items-center justify-center gap-2 rounded-[8px] border border-(--nox-border-strong) py-2 text-[12px] font-semibold text-(--nox-ink-mid)">
                  <Map className="size-4" strokeWidth={1.6} /> {closed ? "View crew" : "Live map"}
                </button>
                {!closed && (
                  <>
                    <button type="button" onClick={() => router.push(`/messages/${crew.conversation_id}`)}
                      className="flex items-center justify-center gap-2 rounded-[8px] border border-(--nox-border-strong) py-2 text-[12px] font-semibold text-(--nox-ink-mid)">
                      <MessageCircle className="size-4" strokeWidth={1.6} /> Chat
                    </button>
                    <button type="button" onClick={() => router.push(`/crews/${crew.id}/members`)}
                      className="flex items-center justify-center gap-2 rounded-[8px] border border-(--nox-border-strong) py-2 text-[12px] font-semibold text-(--nox-ink-mid)">
                      <Users className="size-4" strokeWidth={1.6} /> Members
                    </button>
                  </>
                )}
              </div>
            </div>
          )})}
        </div>

        <p className="py-2 text-center text-[12px] text-(--nox-ink-soft)">
          Location sharing expires when the event window closes.
        </p>
      </div>

      <TabBar />
    </FeedShell>
  );
}
