"use client";

import { useMemo, useState, useEffect } from "react";
import { Plus } from "lucide-react";
import { useRouter } from "next/navigation";
import { FeedShell } from "@/src/components/user/feed/feed-shell";
import { TabBar } from "@/src/components/user/feed/tab-bar";
import { PersonaSwitcher } from "@/src/components/user/shared/persona-switcher";
import { SetCard } from "@/src/components/user/sets/set-card";
import { SetSkeleton } from "@/src/components/user/sets/set-skeleton";
import { useActivePersona } from "@/src/hooks/use-active-persona";
import { getSets } from "@/src/utils/api/user/set";
import type { Set } from "@/src/types/api/set";

export function SetsScreen() {
  const router = useRouter();
  const { personas, activePersona, loading: personaLoading, switchPersona } = useActivePersona();
  const [sets, setSets] = useState<Set[]>([]);
  const [hasMore, setHasMore] = useState(false);
  const [nextOffset, setNextOffset] = useState(0);
  const [loading, setLoading] = useState(true);
  const [loadingMore, setLoadingMore] = useState(false);
  const [message, setMessage] = useState("");
  const [activeGenre, setActiveGenre] = useState("");

  useEffect(() => {
    async function loadSets() {
      setLoading(true);
      setMessage("");
      try {
        const res = await getSets(20, 0);
        const data = res.data;
        setSets(data?.sets ?? []);
        setHasMore(Boolean(data?.has_more));
        setNextOffset(data?.next_offset ?? 0);
      } catch {
        setMessage("Could not load sets.");
      } finally {
        setLoading(false);
      }
    }
    void loadSets();
  }, []);

  async function handleLoadMore() {
    if (!hasMore || loadingMore) return;
    setLoadingMore(true);
    try {
      const res = await getSets(20, nextOffset);
      const data = res.data;
      if (data) {
        setSets((current) => [...current, ...data.sets]);
        setHasMore(data.has_more);
        setNextOffset(data.next_offset ?? 0);
      }
    } catch {
      // keep current list usable
    } finally {
      setLoadingMore(false);
    }
  }

  const genres = useMemo(() => {
    const seen = new Set<string>();
    for (const s of sets) {
      for (const g of s.genre_tags) seen.add(g);
    }
    return Array.from(seen).sort();
  }, [sets]);

  const filtered = activeGenre ? sets.filter((s) => s.genre_tags.includes(activeGenre)) : sets;
  const canCreate = activePersona?.category === "dj";

  return (
    <FeedShell>
      <header className="flex items-center justify-between px-4 pt-[env(safe-area-inset-top,12px)] pb-3">
        <div>
          <h1 className="text-[22px] font-bold tracking-[-0.03em] text-(--nox-ink)">sets</h1>
          <p className="mt-1 text-[12px] text-(--nox-ink-soft)">15-minute mixes from DJs.</p>
        </div>
        <div className="flex items-center gap-2">
          {!personaLoading && canCreate ? (
            <button
              type="button"
              onClick={() => router.push("/sets/create")}
              className="flex size-9 items-center justify-center rounded-[8px] bg-(--nox-accent) text-white transition hover:brightness-110"
              aria-label="Create set"
            >
              <Plus className="size-4" strokeWidth={1.9} />
            </button>
          ) : null}
          <PersonaSwitcher personas={personas} activePersona={activePersona} onSwitch={switchPersona} />
        </div>
      </header>

      {!loading && genres.length > 0 ? (
        <div className="flex gap-2 overflow-x-auto px-4 pb-3 [scrollbar-width:none]">
          <button
            type="button"
            onClick={() => setActiveGenre("")}
            className={`shrink-0 rounded-full border px-3 py-1 font-mono text-[10px] font-medium transition ${
              !activeGenre
                ? "border-(--nox-accent) bg-(--nox-accent-soft) text-(--nox-accent-ink)"
                : "border-(--nox-border) text-(--nox-ink-soft) hover:border-(--nox-accent)"
            }`}
          >
            all
          </button>
          {genres.map((g) => (
            <button
              key={g}
              type="button"
              onClick={() => setActiveGenre(g === activeGenre ? "" : g)}
              className={`shrink-0 rounded-full border px-3 py-1 font-mono text-[10px] font-medium transition ${
                activeGenre === g
                  ? "border-(--nox-accent) bg-(--nox-accent-soft) text-(--nox-accent-ink)"
                  : "border-(--nox-border) text-(--nox-ink-soft) hover:border-(--nox-accent)"
              }`}
            >
              {g}
            </button>
          ))}
        </div>
      ) : null}

      <div className="flex-1 overflow-y-auto">
        {loading ? (
          Array.from({ length: 5 }, (_, i) => <SetSkeleton key={i} />)
        ) : message ? (
          <p className="px-4 py-8 text-[13px] text-(--nox-danger)">{message}</p>
        ) : filtered.length > 0 ? (
          <>
            {filtered.map((set) => (
              <SetCard key={set.id} set={set} onPress={() => router.push(`/sets/${set.id}`)} />
            ))}
            {hasMore && !activeGenre ? (
              <div className="px-4 py-4">
                <button
                  type="button"
                  onClick={handleLoadMore}
                  disabled={loadingMore}
                  className="w-full rounded-[8px] border border-(--nox-border-strong) px-4 py-2.5 text-[13px] font-semibold text-(--nox-ink) transition hover:border-(--nox-accent) disabled:opacity-50"
                >
                  {loadingMore ? "Loading..." : "Load more"}
                </button>
              </div>
            ) : null}
          </>
        ) : sets.length === 0 && !personaLoading && canCreate ? (
          <div className="flex flex-col items-center px-6 py-16 text-center">
            <p className="text-[14px] font-semibold text-(--nox-ink)">No sets yet</p>
            <p className="mt-1 text-[12px] text-(--nox-ink-soft)">Share your first 15-minute mix.</p>
            <button
              type="button"
              onClick={() => router.push("/sets/create")}
              className="mt-5 rounded-[8px] bg-(--nox-accent) px-5 py-2.5 text-[13px] font-semibold text-white transition hover:brightness-110"
            >
              Create set
            </button>
          </div>
        ) : (
          <p className="px-4 py-12 text-center text-[13px] text-(--nox-ink-soft)">No sets yet.</p>
        )}
      </div>

      <TabBar />
    </FeedShell>
  );
}
