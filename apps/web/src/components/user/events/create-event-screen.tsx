"use client";

import { useMemo, useState } from "react";
import { ChevronLeft } from "lucide-react";
import { useRouter } from "next/navigation";
import { FeedShell } from "@/src/components/user/feed/feed-shell";
import { TabBar } from "@/src/components/user/feed/tab-bar";
import { AuthField } from "@/src/components/user/auth/auth-field";
import { useActivePersona } from "@/src/hooks/use-active-persona";
import { createEvent } from "@/src/utils/api/user/event";
import { getAccessToken } from "@/src/utils/auth/session";
import { ApiRequestError } from "@/src/utils/api/api";

export function CreateEventScreen() {
  const router = useRouter();
  const { activePersona, loading: personaLoading } = useActivePersona();
  const [title, setTitle] = useState("");
  const [venue, setVenue] = useState("");
  const [location, setLocation] = useState("");
  const [eventDate, setEventDate] = useState("");
  const [description, setDescription] = useState("");
  const [coverURL, setCoverURL] = useState("");
  const [ticketURL, setTicketURL] = useState("");
  const [price, setPrice] = useState("0");
  const [genres, setGenres] = useState("");
  const [status, setStatus] = useState<"idle" | "loading" | "error">("idle");
  const [message, setMessage] = useState("");

  const canCreate = activePersona?.category === "organizer";
  const canSubmit = useMemo(
    () =>
      canCreate &&
      title.trim().length > 0 &&
      venue.trim().length > 0 &&
      location.trim().length > 0 &&
      eventDate.length > 0 &&
      description.trim().length > 0 &&
      status !== "loading",
    [canCreate, title, venue, location, eventDate, description, status],
  );

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    const token = getAccessToken();
    if (!token || !activePersona || !canSubmit) return;
    setStatus("loading");
    setMessage("");
    try {
      const res = await createEvent(
        {
          title: title.trim(),
          venue: venue.trim(),
          location: location.trim(),
          event_date: new Date(eventDate).toISOString(),
          description: description.trim(),
          cover_url: coverURL.trim() || undefined,
          ticket_url: ticketURL.trim() || undefined,
          price_ngn: Number(price) || 0,
          genre_tags: genres.split(",").map((g) => g.trim().toLowerCase()).filter(Boolean),
          organizer_id: activePersona.id,
        },
        token,
      );
      router.push(res.data?.id ? `/events/${res.data.id}` : "/events");
    } catch (err) {
      setStatus("error");
      setMessage(err instanceof ApiRequestError ? err.message : "Could not create event.");
    }
  }

  return (
    <FeedShell>
      <header className="flex items-center gap-3 border-b border-(--nox-divider) px-4 pt-[env(safe-area-inset-top,12px)] pb-3">
        <button type="button" onClick={() => router.back()} className="flex size-8 items-center justify-center rounded-full hover:bg-(--nox-surface-alt)">
          <ChevronLeft className="size-5 text-(--nox-ink)" strokeWidth={1.8} />
        </button>
        <div>
          <h1 className="text-[18px] font-bold text-(--nox-ink)">create event</h1>
          <p className="text-[11px] text-(--nox-ink-soft)">organizer only</p>
        </div>
      </header>

      <form onSubmit={handleSubmit} className="flex-1 overflow-y-auto px-4 py-4">
        {personaLoading ? (
          <p className="text-[13px] text-(--nox-ink-soft)">Loading persona...</p>
        ) : !canCreate ? (
          <p className="rounded-[8px] border border-(--nox-border) bg-(--nox-surface) px-3 py-3 text-[13px] text-(--nox-ink-soft)">
            Only organizer personas can create events. Switch to an organizer persona to continue.
          </p>
        ) : (
          <div className="grid gap-4">
            <AuthField id="ev-title" label="title" type="text" value={title} placeholder="Warehouse Sessions Vol. 3" autoComplete="off" icon={<span />} onChange={setTitle} />
            <AuthField id="ev-venue" label="venue" type="text" value={venue} placeholder="Eko Atlantic, Lagos" autoComplete="off" icon={<span />} onChange={setVenue} />
            <AuthField id="ev-location" label="location" type="text" value={location} placeholder="Lagos, NG" autoComplete="off" icon={<span />} onChange={setLocation} />
            <AuthField id="ev-date" label="date & time" type="datetime-local" value={eventDate} placeholder="" autoComplete="off" icon={<span />} onChange={setEventDate} />
            <AuthField id="ev-price" label="ticket price (₦)" type="number" value={price} placeholder="0" autoComplete="off" icon={<span />} onChange={setPrice} />
            <AuthField id="ev-cover" label="cover image url" type="url" value={coverURL} placeholder="https://..." autoComplete="off" icon={<span />} onChange={setCoverURL} />
            <AuthField id="ev-ticket" label="ticket url" type="url" value={ticketURL} placeholder="https://paystack.com/..." autoComplete="off" icon={<span />} onChange={setTicketURL} />
            <AuthField id="ev-genres" label="genres" type="text" value={genres} placeholder="afro-house, amapiano" autoComplete="off" icon={<span />} onChange={setGenres} />
            <label htmlFor="ev-desc" className="block">
              <span className="mb-2 block font-mono text-[10px] font-medium uppercase tracking-[0.18em] text-(--nox-ink-soft)">description</span>
              <textarea id="ev-desc" value={description} rows={4} onChange={(e) => setDescription(e.target.value)}
                className="w-full resize-none rounded-[8px] border border-(--nox-border) bg-(--nox-surface) px-3 py-3 text-[14px] text-(--nox-ink) outline-none focus:border-(--nox-accent-line)" />
            </label>
          </div>
        )}

        {message && (
          <p className="mt-4 rounded-[8px] border border-(--nox-danger) bg-(--nox-danger-soft) px-3 py-2 text-[12px] text-(--nox-danger)">{message}</p>
        )}

        {canCreate && (
          <button type="submit" disabled={!canSubmit}
            className="mt-5 w-full rounded-[8px] bg-(--nox-accent) py-3 text-[15px] font-semibold text-white disabled:opacity-50">
            {status === "loading" ? "Creating event..." : "Create event"}
          </button>
        )}
      </form>

      <TabBar />
    </FeedShell>
  );
}
