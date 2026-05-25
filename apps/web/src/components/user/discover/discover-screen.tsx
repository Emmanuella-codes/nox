"use client";

import { useEffect, useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import { Search } from "lucide-react";
import { FeedShell } from "@/src/components/user/feed/feed-shell";
import { TabBar } from "@/src/components/user/feed/tab-bar";
import { PersonaCard } from "@/src/components/user/discover/persona-card";
import { PostCard } from "@/src/components/user/feed/post-card";
import { EventCard } from "@/src/components/user/events/event-card";
import { searchNox } from "@/src/utils/api/user/search";
import type { SearchResponse } from "@/src/types/api/search";

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
  const router = useRouter();
  const [query, setQuery] = useState("");
  const [genre, setGenre] = useState("all");
  const [results, setResults] = useState<SearchResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState("Search artists, posts, events, and genres.");
  const searchTerm = useMemo(() => query.trim() || (genre === "all" ? "" : genre), [genre, query]);

  useEffect(() => {
    if (searchTerm.length < 2) {
      setResults(null);
      setLoading(false);
      setMessage("Search artists, posts, events, and genres.");
      return;
    }

    const controller = new AbortController();
    const timeout = window.setTimeout(async () => {
      setLoading(true);
      setMessage("");
      try {
        const res = await searchNox(searchTerm, 10);
        if (!controller.signal.aborted) {
          setResults(res.data ?? null);
        }
      } catch {
        if (!controller.signal.aborted) {
          setResults(null);
          setMessage("Could not load search results.");
        }
      } finally {
        if (!controller.signal.aborted) {
          setLoading(false);
        }
      }
    }, 250);

    return () => {
      controller.abort();
      window.clearTimeout(timeout);
    };
  }, [searchTerm]);

  const hasResults = Boolean(
    results &&
      (results.personas.length > 0 || results.posts.length > 0 || results.events.length > 0),
  );

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
        {loading ? (
          <p className="px-4 py-8 text-center text-[13px] text-(--nox-ink-soft)">
            Searching...
          </p>
        ) : hasResults && results ? (
          <div className="pb-4">
            {results.personas.length > 0 && (
              <section>
                <p className="px-4 pb-2 font-mono text-[10px] font-semibold uppercase tracking-[0.16em] text-(--nox-ink-soft)">
                  artists
                </p>
                <div className="divide-y divide-(--nox-divider)">
                  {results.personas.map((persona) => (
                    <PersonaCard key={persona.id} persona={persona} />
                  ))}
                </div>
              </section>
            )}

            {results.posts.length > 0 && (
              <section className="mt-4">
                <p className="px-4 pb-2 font-mono text-[10px] font-semibold uppercase tracking-[0.16em] text-(--nox-ink-soft)">
                  posts
                </p>
                <div>
                  {results.posts.map((post) => (
                    <PostCard
                      key={post.id}
                      post={post}
                      onClick={() => router.push(`/posts/${post.id}`)}
                    />
                  ))}
                </div>
              </section>
            )}

            {results.events.length > 0 && (
              <section className="mt-4">
                <p className="px-4 pb-2 font-mono text-[10px] font-semibold uppercase tracking-[0.16em] text-(--nox-ink-soft)">
                  events
                </p>
                <div>
                  {results.events.map((event) => (
                    <EventCard key={event.id} event={event} />
                  ))}
                </div>
              </section>
            )}
          </div>
        ) : (
          <div>
          <p className="px-4 pb-2 font-mono text-[10px] font-semibold uppercase tracking-[0.16em] text-(--nox-ink-soft)">
            discover
          </p>
          <p className="px-4 py-8 text-center text-[13px] leading-6 text-(--nox-ink-soft)">
            {message || `No results for "${searchTerm}".`}
          </p>
        </div>
        )}
      </div>

      <TabBar />
    </FeedShell>
  );
}
