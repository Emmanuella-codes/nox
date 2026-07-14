"use client";

import { ChevronLeft, Ticket } from "lucide-react";
import { useRouter } from "next/navigation";
import { FeedShell } from "@/src/components/user/feed/feed-shell";
import { TabBar } from "@/src/components/user/feed/tab-bar";
import type { PaymentHistoryItem } from "@/src/types/components/user/settings";

const SAMPLE_HISTORY: PaymentHistoryItem[] = [
  { id: "p1", event: "Warehouse Sessions Vol. 3", amount: 5000, qty: 1, date: "2026-06-08", ref: "NOX-A1B2C3", status: "success" },
  { id: "p2", event: "Afro Roots All Night", amount: 0, qty: 1, date: "2026-05-21", ref: "NOX-D4E5F6", status: "success" },
  { id: "p3", event: "Lagos Jazz Collective", amount: 3500, qty: 2, date: "2026-04-14", ref: "NOX-G7H8I9", status: "refunded" },
];

function formatDate(iso: string) {
  return new Date(iso).toLocaleDateString("en-GB", { day: "numeric", month: "short", year: "numeric" });
}

export function PaymentsScreen() {
  const router = useRouter();

  return (
    <FeedShell>
      <header className="flex items-center gap-3 border-b border-(--nox-divider) px-4 pt-[env(safe-area-inset-top,12px)] pb-3">
        <button type="button" onClick={() => router.back()} className="flex size-8 items-center justify-center rounded-full hover:bg-(--nox-surface-alt)">
          <ChevronLeft className="size-5 text-(--nox-ink)" strokeWidth={1.8} />
        </button>
        <h1 className="text-[18px] font-bold text-(--nox-ink)">payments & history</h1>
      </header>

      <div className="flex-1 overflow-y-auto">
        <p className="px-4 pb-1 pt-5 font-mono text-[10px] font-semibold uppercase tracking-[0.18em] text-(--nox-ink-soft)">transaction history</p>

        {SAMPLE_HISTORY.map((item) => (
          <div key={item.id} className="border-b border-(--nox-divider) px-4 py-4">
            <div className="flex items-start justify-between gap-3">
              <div className="min-w-0 flex-1">
                <p className="truncate text-[14px] font-semibold text-(--nox-ink)">{item.event}</p>
                <p className="mt-0.5 text-[11px] text-(--nox-ink-soft)">{formatDate(item.date)} · {item.qty} ticket{item.qty > 1 ? "s" : ""}</p>
                <p className="mt-1 font-mono text-[10px] text-(--nox-ink-faint)">{item.ref}</p>
              </div>
              <div className="shrink-0 text-right">
                <p className="font-mono text-[13px] font-semibold text-(--nox-ink)">
                  {item.amount === 0 ? "free" : `₦${(item.amount * item.qty).toLocaleString()}`}
                </p>
                <span className={`mt-1 inline-block rounded-full px-2 py-0.5 font-mono text-[9px] font-semibold ${
                  item.status === "success"
                    ? "bg-(--nox-success-soft) text-(--nox-success)"
                    : "bg-(--nox-danger-soft) text-(--nox-danger)"
                }`}>
                  {item.status}
                </span>
              </div>
            </div>
          </div>
        ))}

        <div className="px-4 py-6">
          <button type="button" onClick={() => router.push("/tickets")}
            className="flex w-full items-center justify-center gap-2 rounded-[8px] border border-(--nox-border-strong) py-2.5 text-[13px] font-semibold text-(--nox-ink) hover:border-(--nox-accent)">
            <Ticket className="size-4" strokeWidth={1.7} />
            View my tickets
          </button>
        </div>

        <p className="px-4 pb-8 text-center text-[12px] text-(--nox-ink-soft)">
          Full payment history syncs from Paystack once the integration is live.
        </p>
      </div>

      <TabBar />
    </FeedShell>
  );
}
