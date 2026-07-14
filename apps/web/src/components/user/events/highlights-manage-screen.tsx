"use client";

import { useCallback, useEffect, useState } from "react";
import { ChevronLeft, Pin, PinOff } from "lucide-react";
import { useRouter } from "next/navigation";
import { FeedShell } from "@/src/components/user/feed/feed-shell";
import { TabBar } from "@/src/components/user/feed/tab-bar";
import { Avatar } from "@/src/components/user/shared/avatar";
import { useActivePersona } from "@/src/hooks/use-active-persona";
import { listEventStories, listEventHighlightStories, addEventHighlightStory, removeEventHighlightStory } from "@/src/utils/api/user/story";
import { getEvent } from "@/src/utils/api/user/event";
import { getAccessToken } from "@/src/utils/auth/session";
import type { Story, EventHighlightStory } from "@/src/types/api/user/story";

interface HighlightsManageScreenProps { eventID: string }

function formatDuration(s: number) {
  return `${Math.floor(s / 60)}:${String(s % 60).padStart(2, "0")}`;
}

export function HighlightsManageScreen({ eventID }: HighlightsManageScreenProps) {
  const router = useRouter();
  const { activePersona, loading: personaLoading } = useActivePersona();
  const [stories, setStories] = useState<Story[]>([]);
  const [highlights, setHighlights] = useState<EventHighlightStory[]>([]);
  const [loading, setLoading] = useState(true);
  const [toggling, setToggling] = useState<string | null>(null);
  const [isOrganizer, setIsOrganizer] = useState(false);

  const load = useCallback(async () => {
    const token = getAccessToken() ?? undefined;
    const viewerID = activePersona?.id;
    const [evRes, stRes, hlRes] = await Promise.all([
      getEvent(eventID),
      listEventStories(eventID, 50, 0, viewerID, token),
      listEventHighlightStories(eventID, viewerID, token),
    ]);
    setIsOrganizer(evRes.data?.organizer_id === activePersona?.id);
    setStories(stRes.data?.stories ?? []);
    setHighlights(hlRes.data ?? []);
    setLoading(false);
  }, [eventID, activePersona?.id]);

  useEffect(() => {
    if (personaLoading) return;
    void load();
  }, [personaLoading, load]);

  const highlightedIDs = new Set(highlights.map((h) => h.story_id));

  async function handleToggle(story: Story) {
    const token = getAccessToken();
    if (!token || !activePersona || toggling) return;
    setToggling(story.id);
    try {
      if (highlightedIDs.has(story.id)) {
        await removeEventHighlightStory(eventID, story.id, token);
      } else {
        await addEventHighlightStory(eventID, story.id, activePersona.id, token);
      }
      await load();
    } catch {
      // keep current state on error
    } finally {
      setToggling(null);
    }
  }

  return (
    <FeedShell>
      <header className="flex items-center gap-3 border-b border-(--nox-divider) px-4 pt-[env(safe-area-inset-top,12px)] pb-3">
        <button type="button" onClick={() => router.back()} className="flex size-8 items-center justify-center rounded-full hover:bg-(--nox-surface-alt)">
          <ChevronLeft className="size-5 text-(--nox-ink)" strokeWidth={1.8} />
        </button>
        <div>
          <h1 className="text-[18px] font-bold text-(--nox-ink)">highlight reel</h1>
          <p className="text-[11px] text-(--nox-ink-soft)">pin stories to the event page</p>
        </div>
      </header>

      <div className="flex-1 overflow-y-auto">
        {loading ? (
          Array.from({ length: 4 }, (_, i) => (
            <div key={i} className="flex animate-pulse items-center gap-3 border-b border-(--nox-divider) px-4 py-4">
              <div className="size-10 rounded-full bg-(--nox-surface-alt)" />
              <div className="flex-1 space-y-1.5">
                <div className="h-3.5 w-2/3 rounded bg-(--nox-surface-alt)" />
                <div className="h-3 w-1/3 rounded bg-(--nox-surface-alt)" />
              </div>
            </div>
          ))
        ) : !isOrganizer ? (
          <p className="px-4 py-8 text-[13px] text-(--nox-ink-soft)">Only the event organizer can manage highlights.</p>
        ) : stories.length === 0 ? (
          <p className="px-4 py-8 text-[13px] text-(--nox-ink-soft)">No stories for this event yet.</p>
        ) : (
          stories.map((story) => {
            const pinned = highlightedIDs.has(story.id);
            return (
              <div key={story.id} className="flex items-center gap-3 border-b border-(--nox-divider) px-4 py-4">
                <Avatar id={story.owner.id} name={story.owner.display_name} size={40} />
                <div className="min-w-0 flex-1">
                  <p className="truncate text-[14px] font-semibold text-(--nox-ink)">{story.title}</p>
                  <p className="text-[11px] text-(--nox-ink-soft)">
                    {story.owner.display_name} · {formatDuration(story.total_duration_seconds)} · {story.items.length} clips
                  </p>
                </div>
                <button type="button" disabled={toggling === story.id}
                  onClick={() => handleToggle(story)}
                  className={`flex items-center gap-1.5 rounded-[8px] border px-3 py-1.5 text-[12px] font-semibold transition disabled:opacity-50 ${
                    pinned
                      ? "border-(--nox-accent-line) bg-(--nox-accent-soft) text-(--nox-accent-ink)"
                      : "border-(--nox-border-strong) text-(--nox-ink-mid) hover:border-(--nox-accent)"
                  }`}>
                  {pinned ? <PinOff className="size-3" strokeWidth={1.8} /> : <Pin className="size-3" strokeWidth={1.8} />}
                  {pinned ? "unpin" : "pin"}
                </button>
              </div>
            );
          })
        )}
      </div>

      <TabBar />
    </FeedShell>
  );
}
