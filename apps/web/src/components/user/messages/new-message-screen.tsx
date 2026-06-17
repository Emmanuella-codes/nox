"use client";

import { ChevronLeft, Search } from "lucide-react";
import { useState } from "react";
import { useRouter } from "next/navigation";
import { FeedShell } from "@/src/components/user/feed/feed-shell";
import { TabBar } from "@/src/components/user/feed/tab-bar";
import { Avatar } from "@/src/components/user/shared/avatar";

const SUGGESTIONS = [
  { id: "s1", name: "Amirah Lagos", handle: "amirah_lagos" },
  { id: "s2", name: "DJ Khalid", handle: "djkhalid" },
  { id: "s3", name: "Scene Curator", handle: "scene_curator" },
  { id: "s4", name: "Afro House Vibes", handle: "afrohousevibes" },
];

export function NewMessageScreen() {
  const router = useRouter();
  const [query, setQuery] = useState("");

  const filtered = query.trim()
    ? SUGGESTIONS.filter((s) => s.name.toLowerCase().includes(query.toLowerCase()) || s.handle.includes(query.toLowerCase()))
    : SUGGESTIONS;

  return (
    <FeedShell>
      <header className="flex items-center gap-3 border-b border-(--nox-divider) px-4 pt-[env(safe-area-inset-top,12px)] pb-3">
        <button type="button" onClick={() => router.back()} className="flex size-8 items-center justify-center rounded-full hover:bg-(--nox-surface-alt)">
          <ChevronLeft className="size-5 text-(--nox-ink)" strokeWidth={1.8} />
        </button>
        <h1 className="text-[18px] font-bold text-(--nox-ink)">new message</h1>
      </header>

      <div className="border-b border-(--nox-divider) px-4 py-3">
        <div className="flex items-center gap-2 rounded-[10px] border border-(--nox-border) bg-(--nox-surface) px-3 py-2.5">
          <Search className="size-4 shrink-0 text-(--nox-ink-soft)" strokeWidth={1.7} />
          <input value={query} onChange={(e) => setQuery(e.target.value)}
            placeholder="Search personas..."
            className="flex-1 bg-transparent text-[14px] text-(--nox-ink) outline-none placeholder:text-(--nox-ink-soft)" />
        </div>
      </div>

      <div className="flex-1 overflow-y-auto">
        <p className="px-4 pb-1 pt-4 font-mono text-[10px] font-semibold uppercase tracking-[0.18em] text-(--nox-ink-soft)">
          {query ? "results" : "suggested"}
        </p>
        {filtered.map((s) => (
          <button key={s.id} type="button" onClick={() => router.push(`/messages/${s.id}`)}
            className="flex w-full items-center gap-3 border-b border-(--nox-divider) px-4 py-3.5 text-left transition hover:bg-(--nox-surface)">
            <Avatar id={s.id} name={s.name} size={40} />
            <div>
              <p className="text-[14px] font-semibold text-(--nox-ink)">{s.name}</p>
              <p className="font-mono text-[11px] text-(--nox-ink-soft)">@{s.handle}</p>
            </div>
          </button>
        ))}
      </div>

      <TabBar />
    </FeedShell>
  );
}
