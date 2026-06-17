"use client";

import { CalendarDays, MapPin, Ticket } from "lucide-react";
import { useRouter } from "next/navigation";
import { FeedShell } from "@/src/components/user/feed/feed-shell";
import { TabBar } from "@/src/components/user/feed/tab-bar";

// Tickets require a ticketing/payments API that is not yet available.
// This screen renders a placeholder until the backend is wired up.

function formatDate(iso: string) {
  const d = new Date(iso);
  return d.toLocaleDateString("en-GB", { weekday: "short", day: "numeric", month: "short", year: "numeric" });
}

const SAMPLE_TICKETS = [
  { id: "t1", event: "Warehouse Sessions Vol. 3", venue: "Eko Atlantic", location: "Lagos", date: new Date(Date.now() + 7 * 86400000).toISOString(), price: 5000, ref: "NOX-A1B2C3" },
  { id: "t2", event: "Afro Roots All Night", venue: "Hard Rock Café", location: "Lekki", date: new Date(Date.now() + 21 * 86400000).toISOString(), price: 0, ref: "NOX-D4E5F6" },
];

export function TicketsScreen() {
  const router = useRouter();

  return (
    <FeedShell>
      <header className="px-4 pt-[env(safe-area-inset-top,12px)] pb-3">
        <h1 className="text-[22px] font-bold tracking-[-0.03em] text-(--nox-ink)">my tickets</h1>
        <p className="mt-1 text-[12px] text-(--nox-ink-soft)">upcoming events you are attending</p>
      </header>

      <div className="flex-1 overflow-y-auto">
        {SAMPLE_TICKETS.map((ticket) => (
          <button key={ticket.id} type="button" onClick={() => router.push(`/events/${ticket.id}`)}
            className="w-full border-b border-(--nox-divider) px-4 py-4 text-left transition active:bg-(--nox-surface)">
            <div className="flex items-start justify-between gap-3">
              <div className="min-w-0 flex-1">
                <p className="truncate text-[15px] font-bold tracking-[-0.02em] text-(--nox-ink)">{ticket.event}</p>
                <div className="mt-1 flex items-center gap-1.5 text-[12px] text-(--nox-ink-soft)">
                  <MapPin className="size-3 shrink-0" strokeWidth={1.7} />
                  {ticket.venue} · {ticket.location}
                </div>
                <div className="mt-0.5 flex items-center gap-1.5 text-[12px] text-(--nox-ink-soft)">
                  <CalendarDays className="size-3 shrink-0" strokeWidth={1.7} />
                  {formatDate(ticket.date)}
                </div>
              </div>
              <span className="mt-0.5 font-mono text-[13px] font-semibold" style={{ color: ticket.price === 0 ? "var(--nox-success)" : "var(--nox-ink)" }}>
                {ticket.price === 0 ? "free" : `₦${ticket.price.toLocaleString()}`}
              </span>
            </div>
            <div className="mt-3 flex items-center gap-2 rounded-[8px] border border-(--nox-divider) bg-(--nox-surface-alt) px-3 py-2">
              <Ticket className="size-4 shrink-0 text-(--nox-ink-mid)" strokeWidth={1.7} />
              <span className="font-mono text-[11px] text-(--nox-ink-mid)">{ticket.ref}</span>
            </div>
          </button>
        ))}

        <div className="px-4 py-8 text-center">
          <p className="text-[12px] text-(--nox-ink-soft)">Ticket sync with Paystack coming soon.</p>
          <button type="button" onClick={() => router.push("/events")}
            className="mt-3 rounded-[8px] border border-(--nox-border-strong) px-4 py-2 text-[13px] font-semibold text-(--nox-ink) hover:border-(--nox-accent)">
            Browse events
          </button>
        </div>
      </div>

      <TabBar />
    </FeedShell>
  );
}
