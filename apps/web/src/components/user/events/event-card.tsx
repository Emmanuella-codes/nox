"use client";

import { MapPin, Ticket } from "lucide-react";
import type { Event } from "@/src/types/api/user/event";

interface EventCardProps {
  event: Event;
  onPress?: () => void;
}

function formatEventDate(isoString: string) {
  const d = new Date(isoString);
  return {
    day: d.toLocaleDateString("en-GB", { weekday: "short" }).toUpperCase(),
    date: d.getDate(),
    month: d.toLocaleDateString("en-GB", { month: "short" }).toUpperCase(),
    time: d.toLocaleTimeString("en-GB", { hour: "2-digit", minute: "2-digit" }),
  };
}

function formatPrice(ngn: number): string {
  if (ngn === 0) return "free";
  return `₦${ngn.toLocaleString()}`;
}

export function EventCard({ event, onPress }: EventCardProps) {
  const { day, date, month, time } = formatEventDate(event.event_date);

  return (
    <div
      className="flex gap-4 border-b border-(--nox-divider) px-4 py-4 transition active:bg-(--nox-surface)"
      style={{ cursor: onPress ? "pointer" : "default" }}
      onClick={onPress}
    >
      {/* Date badge */}
      <div
        className="flex w-12 shrink-0 flex-col items-center justify-start pt-0.5 rounded-[10px] py-2"
        style={{ background: "var(--nox-surface-alt)" }}
      >
        <span className="font-mono text-[9px] font-semibold uppercase tracking-[0.1em] text-(--nox-ink-soft)">
          {day}
        </span>
        <span className="text-[22px] font-bold leading-tight tracking-[-0.03em] text-(--nox-ink)">
          {date}
        </span>
        <span className="font-mono text-[9px] font-semibold uppercase tracking-[0.1em] text-(--nox-ink-soft)">
          {month}
        </span>
      </div>

      {/* Content */}
      <div className="min-w-0 flex-1">
        <p className="text-[15px] font-bold leading-tight tracking-[-0.02em] text-(--nox-ink) truncate">
          {event.title}
        </p>

        <div className="mt-1 flex items-center gap-1.5">
          <MapPin className="size-3 shrink-0 text-(--nox-ink-soft)" strokeWidth={1.7} />
          <span className="text-[12px] text-(--nox-ink-soft) truncate">
            {event.venue}
            {event.location ? ` · ${event.location}` : ""}
          </span>
        </div>

        <p className="mt-0.5 text-[12px] text-(--nox-ink-soft)">{time}</p>

        {event.genre_tags.length > 0 && (
          <div className="mt-2 flex flex-wrap gap-1">
            {event.genre_tags.slice(0, 3).map((tag) => (
              <span
                key={tag}
                className="rounded-full px-2 py-0.5 font-mono text-[9.5px] font-medium lowercase"
                style={{
                  background: "var(--nox-surface-alt)",
                  color: "var(--nox-ink-mid)",
                }}
              >
                {tag}
              </span>
            ))}
          </div>
        )}

        <div className="mt-3 flex items-center justify-between">
          <span
            className="font-mono text-[11px] font-semibold"
            style={{ color: event.price_ngn === 0 ? "var(--nox-success)" : "var(--nox-ink)" }}
          >
            {formatPrice(event.price_ngn)}
          </span>

          {event.ticket_url && (
            <button
              type="button"
              onClick={(e) => {
                e.stopPropagation();
                window.open(event.ticket_url, "_blank", "noopener");
              }}
              className="flex items-center gap-1.5 rounded-[8px] px-3 py-1.5 text-[12px] font-semibold text-white transition hover:brightness-110"
              style={{ background: "var(--nox-accent)" }}
            >
              <Ticket className="size-3" strokeWidth={1.8} />
              tickets
            </button>
          )}
        </div>
      </div>
    </div>
  );
}
