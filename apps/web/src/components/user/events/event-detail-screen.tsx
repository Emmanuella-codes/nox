"use client";

import { useEffect, useState } from "react";
import { ChevronLeft, ExternalLink, MapPin, Plus, Star, Ticket } from "lucide-react";
import { useRouter } from "next/navigation";
import { FeedShell } from "@/src/components/user/feed/feed-shell";
import { TabBar } from "@/src/components/user/feed/tab-bar";
import { StoryCard } from "@/src/components/user/events/story-card";
import { useActivePersona } from "@/src/hooks/use-active-persona";
import { getEvent } from "@/src/utils/api/user/event";
import { listEventStories, listEventHighlightStories } from "@/src/utils/api/user/story";
import { getAccessToken } from "@/src/utils/auth/session";
import type { Event } from "@/src/types/api/user/event";
import type { Story, EventHighlightStory } from "@/src/types/api/user/story";

interface EventDetailScreenProps { eventID: string }

function formatDate(iso: string) {
  const d = new Date(iso);
  return {
    full: d.toLocaleDateString("en-GB", { weekday: "long", day: "numeric", month: "long", year: "numeric" }),
    time: d.toLocaleTimeString("en-GB", { hour: "2-digit", minute: "2-digit" }),
  };
}

export function EventDetailScreen({ eventID }: EventDetailScreenProps) {
  const router = useRouter();
  const { activePersona, loading: personaLoading } = useActivePersona();
  const [event, setEvent] = useState<Event | null>(null);
  const [highlights, setHighlights] = useState<EventHighlightStory[]>([]);
  const [stories, setStories] = useState<Story[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (personaLoading) return;
    const token = getAccessToken() ?? undefined;
    const viewerID = activePersona?.id;
    Promise.all([
      getEvent(eventID),
      listEventHighlightStories(eventID, viewerID, token),
      listEventStories(eventID, 20, 0, viewerID, token),
    ])
      .then(([evRes, hlRes, stRes]) => {
        setEvent(evRes.data ?? null);
        setHighlights(hlRes.data ?? []);
        setStories(stRes.data?.stories ?? []);
      })
      .finally(() => setLoading(false));
  }, [eventID, personaLoading, activePersona?.id]);

  const isOrganizer = event && activePersona?.id === event.organizer_id;

  if (loading) {
    return (
      <FeedShell>
        <header className="flex items-center gap-3 px-4 pt-[env(safe-area-inset-top,12px)] pb-3">
          <button type="button" onClick={() => router.back()} className="flex size-8 items-center justify-center"><ChevronLeft className="size-5 text-(--nox-ink)" strokeWidth={1.8} /></button>
          <div className="h-5 w-40 animate-pulse rounded bg-(--nox-surface-alt)" />
        </header>
        <div className="flex-1 space-y-3 overflow-y-auto px-4 py-4">
          <div className="h-48 animate-pulse rounded-[12px] bg-(--nox-surface-alt)" />
          <div className="h-4 w-3/4 animate-pulse rounded bg-(--nox-surface-alt)" />
          <div className="h-3 w-1/2 animate-pulse rounded bg-(--nox-surface-alt)" />
        </div>
        <TabBar />
      </FeedShell>
    );
  }

  if (!event) {
    return (
      <FeedShell>
        <header className="flex items-center gap-3 px-4 pt-[env(safe-area-inset-top,12px)] pb-3">
          <button type="button" onClick={() => router.back()} className="flex size-8 items-center justify-center"><ChevronLeft className="size-5 text-(--nox-ink)" strokeWidth={1.8} /></button>
        </header>
        <p className="px-4 py-8 text-[13px] text-(--nox-danger)">Event not found.</p>
        <TabBar />
      </FeedShell>
    );
  }

  const { full, time } = formatDate(event.event_date);

  return (
    <FeedShell>
      <header className="flex items-center gap-3 border-b border-(--nox-divider) px-4 pt-[env(safe-area-inset-top,12px)] pb-3">
        <button type="button" onClick={() => router.back()} className="flex size-8 shrink-0 items-center justify-center rounded-full hover:bg-(--nox-surface-alt)">
          <ChevronLeft className="size-5 text-(--nox-ink)" strokeWidth={1.8} />
        </button>
        <p className="min-w-0 truncate text-[16px] font-bold text-(--nox-ink)">{event.title}</p>
      </header>

      <div className="flex-1 overflow-y-auto">
        {event.cover_url && (
          // eslint-disable-next-line @next/next/no-img-element
          <img src={event.cover_url} alt="" className="h-48 w-full object-cover" />
        )}
        {!event.cover_url && (
          <div className="h-48 w-full" style={{ background: "linear-gradient(135deg, #1a1028 0%, #0d0d14 100%)" }} />
        )}

        <div className="px-4 py-4">
          <h1 className="text-[22px] font-bold tracking-[-0.03em] text-(--nox-ink)">{event.title}</h1>
          <div className="mt-2 flex items-center gap-1.5 text-[13px] text-(--nox-ink-soft)">
            <MapPin className="size-3.5 shrink-0" strokeWidth={1.7} />
            {event.venue}{event.location ? ` · ${event.location}` : ""}
          </div>
          <p className="mt-1 text-[13px] text-(--nox-ink-soft)">{full} · {time}</p>

          <div className="mt-3 flex items-center justify-between">
            <span className="font-mono text-[13px] font-semibold" style={{ color: event.price_ngn === 0 ? "var(--nox-success)" : "var(--nox-ink)" }}>
              {event.price_ngn === 0 ? "free" : `₦${event.price_ngn.toLocaleString()}`}
            </span>
            <div className="flex items-center gap-2">
              {isOrganizer && (
                <button type="button" onClick={() => router.push(`/events/${eventID}/highlights/manage`)}
                  className="flex items-center gap-1.5 rounded-[8px] border border-(--nox-border-strong) px-3 py-1.5 text-[12px] font-semibold text-(--nox-ink-mid) hover:border-(--nox-accent)">
                  <Star className="size-3" strokeWidth={1.7} /> highlights
                </button>
              )}
              {event.ticket_url && (
                <a href={event.ticket_url} target="_blank" rel="noopener noreferrer"
                  className="flex items-center gap-1.5 rounded-[8px] bg-(--nox-accent) px-3 py-1.5 text-[12px] font-semibold text-white">
                  <Ticket className="size-3" strokeWidth={1.8} /> get tickets <ExternalLink className="size-3" />
                </a>
              )}
            </div>
          </div>

          {event.genre_tags.length > 0 && (
            <div className="mt-3 flex flex-wrap gap-1.5">
              {event.genre_tags.map((tag) => (
                <span key={tag} className="rounded-full bg-(--nox-accent-soft) px-2.5 py-1 font-mono text-[10px] text-(--nox-accent-ink)">{tag}</span>
              ))}
            </div>
          )}

          {event.description && (
            <p className="mt-4 text-[13px] leading-6 text-(--nox-ink-mid)">{event.description}</p>
          )}
        </div>

        {highlights.length > 0 && (
          <div className="mt-1 border-t border-(--nox-divider) px-4 py-4">
            <p className="mb-3 font-mono text-[10px] font-semibold uppercase tracking-[0.18em] text-(--nox-ink-soft)">highlights</p>
            <div className="flex gap-3 overflow-x-auto [scrollbar-width:none]">
              {highlights.map((h) => h.story ? (
                <StoryCard key={h.id} story={h.story} compact onPress={() => router.push(`/stories/${h.story!.id}`)} />
              ) : null)}
            </div>
          </div>
        )}

        <div className="border-t border-(--nox-divider)">
          <div className="flex items-center justify-between px-4 py-4">
            <p className="font-mono text-[10px] font-semibold uppercase tracking-[0.18em] text-(--nox-ink-soft)">stories</p>
            {activePersona && (
              <button type="button" onClick={() => router.push(`/stories/create?event_id=${eventID}`)}
                className="flex items-center gap-1.5 rounded-[8px] bg-(--nox-accent) px-3 py-1.5 text-[12px] font-semibold text-white">
                <Plus className="size-3" strokeWidth={1.9} /> start a story
              </button>
            )}
          </div>
          {stories.length > 0 ? (
            stories.map((s) => (
              <StoryCard key={s.id} story={s} onPress={() => router.push(`/stories/${s.id}`)} />
            ))
          ) : (
            <p className="px-4 pb-6 text-[13px] text-(--nox-ink-soft)">No stories yet. Be the first.</p>
          )}
        </div>
      </div>

      <TabBar />
    </FeedShell>
  );
}
