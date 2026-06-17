"use client";

import { ChevronLeft, Map, Plus, Users } from "lucide-react";
import { useState } from "react";
import { useRouter } from "next/navigation";
import { FeedShell } from "@/src/components/user/feed/feed-shell";
import { TabBar } from "@/src/components/user/feed/tab-bar";
import { useActivePersona } from "@/src/hooks/use-active-persona";

interface CrewHubScreenProps { eventID: string }

export function CrewHubScreen({ eventID }: CrewHubScreenProps) {
  const router = useRouter();
  const { activePersona } = useActivePersona();
  const [joinCode, setJoinCode] = useState("");

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
        <button type="button" onClick={() => alert("Crew feature coming soon")}
          className="flex w-full items-center gap-4 rounded-[12px] border border-(--nox-border) bg-(--nox-surface) p-5 text-left transition hover:border-(--nox-accent)">
          <span className="flex size-10 items-center justify-center rounded-[10px] bg-(--nox-accent-soft)">
            <Plus className="size-5 text-(--nox-accent-ink)" strokeWidth={1.8} />
          </span>
          <div>
            <p className="text-[15px] font-bold text-(--nox-ink)">Create a crew</p>
            <p className="mt-0.5 text-[12px] text-(--nox-ink-soft)">Start a group — share a code to invite friends</p>
          </div>
        </button>

        <div className="rounded-[12px] border border-(--nox-border) bg-(--nox-surface) p-5">
          <p className="text-[15px] font-bold text-(--nox-ink)">Join with a code</p>
          <p className="mt-0.5 text-[12px] text-(--nox-ink-soft)">Got an invite? Enter the 6-character code</p>
          <div className="mt-3 flex gap-2">
            <input value={joinCode} onChange={(e) => setJoinCode(e.target.value.toUpperCase())}
              maxLength={6} placeholder="ABC123"
              className="flex-1 rounded-[8px] border border-(--nox-border) bg-(--nox-surface-alt) px-3 py-2.5 font-mono text-[15px] tracking-[0.2em] text-(--nox-ink) outline-none focus:border-(--nox-accent-line) uppercase" />
            <button type="button" disabled={joinCode.length < 6}
              className="rounded-[8px] bg-(--nox-accent) px-4 py-2.5 text-[13px] font-semibold text-white disabled:opacity-40">
              Join
            </button>
          </div>
        </div>

        <div className="flex gap-3">
          <button type="button" onClick={() => router.push(`/crews/demo`)}
            className="flex flex-1 flex-col items-center gap-2 rounded-[12px] border border-(--nox-border) bg-(--nox-surface) py-5 transition hover:border-(--nox-accent)">
            <Map className="size-6 text-(--nox-accent)" strokeWidth={1.6} />
            <span className="text-[12px] font-semibold text-(--nox-ink)">Live map</span>
          </button>
          <button type="button" onClick={() => router.push(`/crews/demo/members`)}
            className="flex flex-1 flex-col items-center gap-2 rounded-[12px] border border-(--nox-border) bg-(--nox-surface) py-5 transition hover:border-(--nox-accent)">
            <Users className="size-6 text-(--nox-accent)" strokeWidth={1.6} />
            <span className="text-[12px] font-semibold text-(--nox-ink)">Members</span>
          </button>
        </div>

        <p className="py-2 text-center text-[12px] text-(--nox-ink-soft)">
          Crew & live location features are coming soon. Location sharing expires when the event ends.
        </p>
      </div>

      <TabBar />
    </FeedShell>
  );
}
