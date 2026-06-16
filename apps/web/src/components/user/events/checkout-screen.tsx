"use client";

import { useEffect, useState } from "react";
import { ChevronLeft, Lock } from "lucide-react";
import { useRouter } from "next/navigation";
import { FeedShell } from "@/src/components/user/feed/feed-shell";
import { TabBar } from "@/src/components/user/feed/tab-bar";
import { getEvent } from "@/src/utils/api/user/event";
import type { Event } from "@/src/types/api/event";

interface CheckoutScreenProps { eventID: string }

export function CheckoutScreen({ eventID }: CheckoutScreenProps) {
  const router = useRouter();
  const [event, setEvent] = useState<Event | null>(null);
  const [loading, setLoading] = useState(true);
  const [quantity, setQuantity] = useState(1);

  useEffect(() => {
    getEvent(eventID).then((res) => setEvent(res.data ?? null)).finally(() => setLoading(false));
  }, [eventID]);

  const unitPrice = event?.price_ngn ?? 0;
  const total = unitPrice * quantity;
  const isFree = unitPrice === 0;

  return (
    <FeedShell>
      <header className="flex items-center gap-3 border-b border-(--nox-divider) px-4 pt-[env(safe-area-inset-top,12px)] pb-3">
        <button type="button" onClick={() => router.back()} className="flex size-8 items-center justify-center rounded-full hover:bg-(--nox-surface-alt)">
          <ChevronLeft className="size-5 text-(--nox-ink)" strokeWidth={1.8} />
        </button>
        <h1 className="text-[18px] font-bold text-(--nox-ink)">checkout</h1>
      </header>

      <div className="flex-1 overflow-y-auto px-4 py-4">
        {loading ? (
          <div className="space-y-3">
            <div className="h-5 w-2/3 animate-pulse rounded bg-(--nox-surface-alt)" />
            <div className="h-4 w-1/2 animate-pulse rounded bg-(--nox-surface-alt)" />
          </div>
        ) : !event ? (
          <p className="text-[13px] text-(--nox-danger)">Event not found.</p>
        ) : (
          <div className="space-y-5">
            <div className="rounded-[10px] border border-(--nox-border) bg-(--nox-surface) p-4">
              <p className="text-[15px] font-bold text-(--nox-ink)">{event.title}</p>
              <p className="mt-1 text-[12px] text-(--nox-ink-soft)">{event.venue} · {event.location}</p>
              <p className="mt-1 font-mono text-[13px] font-semibold text-(--nox-ink)">
                {isFree ? "Free" : `₦${unitPrice.toLocaleString()} per ticket`}
              </p>
            </div>

            {!isFree && (
              <div>
                <p className="mb-2 font-mono text-[10px] font-medium uppercase tracking-[0.18em] text-(--nox-ink-soft)">quantity</p>
                <div className="flex items-center gap-3">
                  <button type="button" onClick={() => setQuantity((q) => Math.max(1, q - 1))}
                    className="flex size-9 items-center justify-center rounded-full border border-(--nox-border-strong) text-[18px] font-bold text-(--nox-ink) hover:border-(--nox-accent)">−</button>
                  <span className="w-8 text-center text-[16px] font-bold text-(--nox-ink)">{quantity}</span>
                  <button type="button" onClick={() => setQuantity((q) => Math.min(10, q + 1))}
                    className="flex size-9 items-center justify-center rounded-full border border-(--nox-border-strong) text-[18px] font-bold text-(--nox-ink) hover:border-(--nox-accent)">+</button>
                </div>
              </div>
            )}

            <div className="rounded-[10px] border border-(--nox-divider) bg-(--nox-surface-alt) px-4 py-3">
              <div className="flex justify-between text-[13px] text-(--nox-ink-mid)">
                <span>{quantity} × {isFree ? "free" : `₦${unitPrice.toLocaleString()}`}</span>
                <span className="font-semibold text-(--nox-ink)">{isFree ? "₦0" : `₦${total.toLocaleString()}`}</span>
              </div>
            </div>

            <div className="rounded-[10px] border border-(--nox-border) bg-(--nox-surface) px-4 py-4">
              <p className="mb-3 font-mono text-[10px] font-medium uppercase tracking-[0.18em] text-(--nox-ink-soft)">payment via paystack</p>
              <p className="text-[12px] text-(--nox-ink-soft)">You will be redirected to Paystack to complete your payment securely.</p>
            </div>

            {event.ticket_url ? (
              <a href={event.ticket_url} target="_blank" rel="noopener noreferrer"
                className="flex w-full items-center justify-center gap-2 rounded-[8px] bg-(--nox-accent) py-3 text-[15px] font-semibold text-white">
                <Lock className="size-4" strokeWidth={1.8} />
                {isFree ? "Reserve free ticket" : `Pay ₦${total.toLocaleString()}`}
              </a>
            ) : (
              <p className="text-center text-[13px] text-(--nox-ink-soft)">Ticket purchase link not available yet.</p>
            )}
          </div>
        )}
      </div>

      <TabBar />
    </FeedShell>
  );
}
