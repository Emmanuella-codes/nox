"use client";

import { useState } from "react";
import { Search } from "lucide-react";
import { FeedShell } from "@/src/components/user/feed/feed-shell";
import { TabBar } from "@/src/components/user/feed/tab-bar";

const GENRE_FILTERS = [
  "all",
  "afrobeats",
  "amapiano",
  "afro-house",
  "afro-tech",
  "alt-R&B",
  "hip-hop",
  "dancehall",
  "electronic",
];

export function DiscoverScreen() {
  const [query, setQuery] = useState("");
  const [genre, setGenre] = useState("all");

  return (
    <FeedShell>
      {/* Top bar */}
      <header className="px-4 pt-[env(safe-area-inset-top,12px)] pb-3">
        <h1
          className="mb-3 text-[22px] font-bold tracking-[-0.03em] text-(--nox-ink)"
        >
          discover
        </h1>

        {/* Search */}
        <div className="flex items-center gap-3 rounded-[12px] border border-(--nox-border) bg-(--nox-surface) px-3.5 py-3 transition focus-within:border-(--nox-accent-line)">
          <Search className="size-4 shrink-0 text-(--nox-ink-soft)" strokeWidth={1.7} />
          <input
            type="text"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="search artists, handles, genres…"
            className="flex-1 bg-transparent text-[14px] text-(--nox-ink) outline-none placeholder:text-(--nox-ink-soft)"
          />
        </div>
      </header>

      {/* Genre filter chips */}
      <div className="flex gap-2 overflow-x-auto px-4 pb-3 [scrollbar-width:none] [&::-webkit-scrollbar]:hidden">
        {GENRE_FILTERS.map((g) => {
          const active = genre === g;
          return (
            <button
              key={g}
              type="button"
              onClick={() => setGenre(g)}
              className="shrink-0 rounded-full border px-3 py-1.5 font-mono text-[10.5px] font-medium lowercase transition"
              style={{
                borderColor: active ? "var(--nox-accent-line)" : "var(--nox-border)",
                background: active ? "var(--nox-accent-soft)" : "transparent",
                color: active ? "var(--nox-accent-ink)" : "var(--nox-ink-mid)",
              }}
            >
              {g}
            </button>
          );
        })}
      </div>

      <div className="flex-1 overflow-y-auto">
        <div>
          <p className="px-4 pb-2 font-mono text-[10px] font-semibold uppercase tracking-[0.16em] text-(--nox-ink-soft)">
            artists
          </p>
          <p className="px-4 py-8 text-center text-[13px] leading-6 text-(--nox-ink-soft)">
            Discovery is opening soon. Feed and events are ready now.
            {query || genre !== "all" ? " Your filters will stay here for this visit." : ""}
          </p>
        </div>
      </div>

      <TabBar />
    </FeedShell>
  );
}
