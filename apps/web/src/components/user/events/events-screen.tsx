"use client";

import { useEffect, useState } from "react";
import { FeedShell } from "@/src/components/user/feed/feed-shell";
import { TabBar } from "@/src/components/user/feed/tab-bar";
import { EventCard } from "@/src/components/user/events/event-card";
import { getEvents } from "@/src/utils/api/user/event";
import type { Event } from "@/src/types/api/event";

type EventFilter = "upcoming" | "this-week" | "this-month";

const FILTERS: { key: EventFilter; label: string }[] = [
  { key: "upcoming", label: "upcoming" },
  { key: "this-week", label: "this week" },
  { key: "this-month", label: "this month" },
];

const day = 86_400_000;

function filterEvents(events: Event[], filter: EventFilter): Event[] {
  const nowDate = new Date();
  return events.filter((e) => {
    const d = new Date(e.event_date);
    if (d < nowDate) return false;
    if (filter === "this-week") {
      const diff = d.getTime() - nowDate.getTime();
      return diff <= 7 * day;
    }
    if (filter === "this-month") {
      return d.getMonth() === nowDate.getMonth() && d.getFullYear() === nowDate.getFullYear();
    }
    return true;
  });
}

export function EventsScreen() {
  const [filter, setFilter] = useState<EventFilter>("upcoming");
  const [allEvents, setAllEvents] = useState<Event[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    async function loadEvents() {
      try {
        const res = await getEvents();
        setAllEvents(res.data ?? []);
      } catch {
        setError("Could not load events.");
      } finally {
        setLoading(false);
      }
    }
    loadEvents();
  }, []);

  const events = filterEvents(allEvents, filter);

  return (
    <FeedShell>
      {/* Header */}
      <header className="px-4 pt-[env(safe-area-inset-top,12px)] pb-3">
        <h1 className="text-[22px] font-bold tracking-[-0.03em] text-(--nox-ink)">events</h1>
      </header>

      {/* Filter tabs */}
      <div className="flex gap-2 border-b border-(--nox-divider) px-4 pb-3">
        {FILTERS.map(({ key, label }) => {
          const active = filter === key;
          return (
            <button
              key={key}
              type="button"
              onClick={() => setFilter(key)}
              className="rounded-full border px-3.5 py-1.5 font-mono text-[10.5px] font-medium transition"
              style={{
                borderColor: active ? "var(--nox-accent-line)" : "var(--nox-border)",
                background: active ? "var(--nox-accent-soft)" : "transparent",
                color: active ? "var(--nox-accent-ink)" : "var(--nox-ink-mid)",
              }}
            >
              {label}
            </button>
          );
        })}
      </div>

      {/* Event list */}
      <div className="flex-1 overflow-y-auto">
        {loading ? (
          <div className="px-4 py-8 text-[13px] text-(--nox-ink-soft)">Loading events...</div>
        ) : error ? (
          <p className="px-4 py-8 text-[13px] text-(--nox-danger)">{error}</p>
        ) : events.length > 0 ? (
          events.map((e) => <EventCard key={e.id} event={e} />)
        ) : (
          <div className="flex flex-col items-center justify-center py-16">
            <p className="text-[13px] text-(--nox-ink-soft)">no events {filter.replace("-", " ")}</p>
          </div>
        )}
      </div>

      <TabBar />
    </FeedShell>
  );
}
