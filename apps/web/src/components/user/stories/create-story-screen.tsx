"use client";

import { useEffect, useMemo, useState } from "react";
import { ChevronLeft } from "lucide-react";
import { useRouter, useSearchParams } from "next/navigation";
import { FeedShell } from "@/src/components/user/feed/feed-shell";
import { TabBar } from "@/src/components/user/feed/tab-bar";
import { AuthField } from "@/src/components/user/auth/auth-field";
import { useActivePersona } from "@/src/hooks/use-active-persona";
import { getEvents } from "@/src/utils/api/user/event";
import { createStory } from "@/src/utils/api/user/story";
import { getAccessToken } from "@/src/utils/auth/session";
import { ApiRequestError } from "@/src/utils/api/api";
import type { StoryContributionMode } from "@/src/types/api/story";
import type { Event } from "@/src/types/api/event";

export function CreateStoryScreen() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { activePersona, loading: personaLoading } = useActivePersona();

  const preselectedEventID = searchParams.get("event_id") ?? "";
  const [events, setEvents] = useState<Event[]>([]);
  const [eventID, setEventID] = useState(preselectedEventID);
  const [title, setTitle] = useState("");
  const [contributionMode, setContributionMode] = useState<StoryContributionMode>("public");
  const [status, setStatus] = useState<"idle" | "loading" | "error">("idle");
  const [message, setMessage] = useState("");

  useEffect(() => {
    if (preselectedEventID) return;
    void getEvents().then((res) => setEvents(res.data ?? []));
  }, [preselectedEventID]);

  const canSubmit = useMemo(
    () => Boolean(activePersona) && title.trim().length > 0 && eventID.length > 0 && status !== "loading",
    [activePersona, title, eventID, status],
  );

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    const token = getAccessToken();
    if (!token || !activePersona || !canSubmit) return;
    setStatus("loading");
    setMessage("");
    try {
      const res = await createStory(
        { event_id: eventID, owner_persona_id: activePersona.id, title: title.trim(), contribution_mode: contributionMode },
        token,
      );
      router.push(res.data?.id ? `/stories/${res.data.id}` : "/feed");
    } catch (err) {
      setStatus("error");
      setMessage(err instanceof ApiRequestError ? err.message : "Could not create story.");
    }
  }

  return (
    <FeedShell>
      <header className="flex items-center gap-3 border-b border-(--nox-divider) px-4 pt-[env(safe-area-inset-top,12px)] pb-3">
        <button type="button" onClick={() => router.back()} className="flex size-8 items-center justify-center rounded-full transition hover:bg-(--nox-surface-alt)">
          <ChevronLeft className="size-5 text-(--nox-ink)" strokeWidth={1.8} />
        </button>
        <div>
          <h1 className="text-[18px] font-bold text-(--nox-ink)">start a story</h1>
          <p className="text-[11px] text-(--nox-ink-soft)">tied to an event · expires in 24h</p>
        </div>
      </header>

      <form onSubmit={handleSubmit} className="flex-1 overflow-y-auto px-4 py-4">
        {personaLoading ? (
          <p className="text-[13px] text-(--nox-ink-soft)">Loading persona...</p>
        ) : !activePersona ? (
          <p className="text-[13px] text-(--nox-ink-soft)">Sign in to start a story.</p>
        ) : (
          <div className="grid gap-4">
            <AuthField id="story-title" label="title" type="text" value={title} placeholder="Warehouse Sessions" autoComplete="off" icon={<span />} onChange={setTitle} />

            {!preselectedEventID && (
              <div>
                <label className="mb-2 block font-mono text-[10px] font-medium uppercase tracking-[0.18em] text-(--nox-ink-soft)">event</label>
                <select
                  value={eventID}
                  onChange={(e) => setEventID(e.target.value)}
                  className="w-full rounded-[8px] border border-(--nox-border) bg-(--nox-surface) px-3 py-3 text-[14px] text-(--nox-ink) outline-none focus:border-(--nox-accent-line)"
                >
                  <option value="">Select an event</option>
                  {events.map((ev) => (
                    <option key={ev.id} value={ev.id}>{ev.title}</option>
                  ))}
                </select>
              </div>
            )}

            <div>
              <p className="mb-2 font-mono text-[10px] font-medium uppercase tracking-[0.18em] text-(--nox-ink-soft)">who can contribute</p>
              <div className="flex gap-2">
                {(["public", "private"] as StoryContributionMode[]).map((mode) => (
                  <button key={mode} type="button" onClick={() => setContributionMode(mode)}
                    className={`rounded-full border px-4 py-2 font-mono text-[10px] font-medium transition ${
                      contributionMode === mode
                        ? "border-(--nox-accent) bg-(--nox-accent-soft) text-(--nox-accent-ink)"
                        : "border-(--nox-border) text-(--nox-ink-soft)"
                    }`}>
                    {mode === "public" ? "everyone" : "followers only"}
                  </button>
                ))}
              </div>
            </div>
          </div>
        )}

        {message && (
          <p className="mt-4 rounded-[8px] border border-(--nox-danger) bg-(--nox-danger-soft) px-3 py-2 text-[12px] text-(--nox-danger)">
            {message}
          </p>
        )}

        {activePersona && (
          <button type="submit" disabled={!canSubmit}
            className="mt-5 w-full rounded-[8px] bg-(--nox-accent) py-3 text-[15px] font-semibold text-white transition disabled:opacity-50">
            {status === "loading" ? "Creating story..." : "Create story"}
          </button>
        )}
      </form>

      <TabBar />
    </FeedShell>
  );
}
