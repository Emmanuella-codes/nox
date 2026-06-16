"use client";

import { useState } from "react";
import { ChevronLeft } from "lucide-react";
import { useRouter } from "next/navigation";
import { FeedShell } from "@/src/components/user/feed/feed-shell";
import { TabBar } from "@/src/components/user/feed/tab-bar";

interface ToggleRowProps { label: string; sub?: string; value: boolean; onChange: (v: boolean) => void }
function ToggleRow({ label, sub, value, onChange }: ToggleRowProps) {
  return (
    <div className="flex items-center justify-between border-b border-(--nox-divider) px-4 py-4">
      <div className="mr-4 min-w-0">
        <p className="text-[14px] font-medium text-(--nox-ink)">{label}</p>
        {sub && <p className="mt-0.5 text-[11px] text-(--nox-ink-soft)">{sub}</p>}
      </div>
      <button type="button" onClick={() => onChange(!value)} role="switch" aria-checked={value}
        className={`relative h-6 w-11 shrink-0 rounded-full transition-colors ${value ? "bg-(--nox-accent)" : "bg-(--nox-surface-alt)"}`}>
        <span className={`absolute top-0.5 h-5 w-5 rounded-full bg-white shadow transition-transform ${value ? "translate-x-5" : "translate-x-0.5"}`} />
      </button>
    </div>
  );
}

export function PrivacyScreen() {
  const router = useRouter();
  const [anonDefault, setAnonDefault] = useState(false);
  const [storyFollowers, setStoryFollowers] = useState(false);
  const [crewLocation, setCrewLocation] = useState(true);
  const [searchable, setSearchable] = useState(true);

  return (
    <FeedShell>
      <header className="flex items-center gap-3 border-b border-(--nox-divider) px-4 pt-[env(safe-area-inset-top,12px)] pb-3">
        <button type="button" onClick={() => router.back()} className="flex size-8 items-center justify-center rounded-full hover:bg-(--nox-surface-alt)">
          <ChevronLeft className="size-5 text-(--nox-ink)" strokeWidth={1.8} />
        </button>
        <h1 className="text-[18px] font-bold text-(--nox-ink)">privacy</h1>
      </header>

      <div className="flex-1 overflow-y-auto">
        <p className="px-4 pb-1 pt-5 font-mono text-[10px] font-semibold uppercase tracking-[0.18em] text-(--nox-ink-soft)">posting</p>
        <ToggleRow label="default to anonymous" sub="Story clips will be anonymous unless you change it" value={anonDefault} onChange={setAnonDefault} />
        <ToggleRow label="story visibility" sub="New stories default to followers-only (instead of public)" value={storyFollowers} onChange={setStoryFollowers} />

        <p className="px-4 pb-1 pt-5 font-mono text-[10px] font-semibold uppercase tracking-[0.18em] text-(--nox-ink-soft)">crew & location</p>
        <ToggleRow label="share location in crew" sub="Allow crew members to see your pin on the live map" value={crewLocation} onChange={setCrewLocation} />

        <p className="px-4 pb-1 pt-5 font-mono text-[10px] font-semibold uppercase tracking-[0.18em] text-(--nox-ink-soft)">discovery</p>
        <ToggleRow label="appear in search" sub="Let other personas find you in search and discover" value={searchable} onChange={setSearchable} />

        <p className="px-4 py-6 text-[12px] text-(--nox-ink-soft)">
          Privacy controls are stored locally until the server-side preferences API is available.
        </p>
      </div>

      <TabBar />
    </FeedShell>
  );
}
