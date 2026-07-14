"use client";

import { useEffect, useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import { Search } from "lucide-react";
import { FeedShell } from "@/src/components/user/feed/feed-shell";
import { TabBar } from "@/src/components/user/feed/tab-bar";
import { PersonaCard } from "@/src/components/user/discover/persona-card";
import { PostCard } from "@/src/components/user/feed/post-card";
import { EventCard } from "@/src/components/user/events/event-card";
import { TrendingHashtags } from "@/src/components/user/discover/trending-hashtags";
import { searchNox } from "@/src/utils/api/user/search";
import { followPersona, unfollowPersona } from "@/src/utils/api/user/follow";
import { likePost, unlikePost } from "@/src/utils/api/user/post";
import { useActivePersona } from "@/src/hooks/use-active-persona";
import { getAccessToken } from "@/src/utils/auth/session";
import type { SearchResponse } from "@/src/types/api/user/search";
import type { Post } from "@/src/types/api/user/post";
import type { Persona } from "@/src/types/api/user/persona";

const GENRE_FILTERS = [
  "all", "afrobeats", "amapiano", "afro-house", "afro-tech", "alt-R&B", "hip-hop", "dancehall", "electronic",
];

export function DiscoverScreen() {
  const router = useRouter();
  const { activeID: viewerPersonaID } = useActivePersona();
  const [query, setQuery] = useState("");
  const [genre, setGenre] = useState("all");
  const [results, setResults] = useState<SearchResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [loadingMore, setLoadingMore] = useState(false);
  const searchTerm = useMemo(() => query.trim() || (genre === "all" ? "" : genre), [genre, query]);

  useEffect(() => {
    if (searchTerm.length < 2) {
      setResults(null);
      setLoading(false);
      return;
    }

    const controller = new AbortController();
    const timeout = window.setTimeout(async () => {
      setLoading(true);
      try {
        const token = viewerPersonaID ? getAccessToken() : undefined;
        const res = await searchNox(searchTerm, 10, token || undefined, viewerPersonaID || undefined);
        if (!controller.signal.aborted) setResults(res.data ?? null);
      } catch {
        if (!controller.signal.aborted) setResults(null);
      } finally {
        if (!controller.signal.aborted) setLoading(false);
      }
    }, 250);

    return () => { controller.abort(); window.clearTimeout(timeout); };
  }, [searchTerm, viewerPersonaID]);

  async function handleLoadMore() {
    if (!results?.has_more || results.next_offset === undefined || loadingMore) return;
    setLoadingMore(true);
    try {
      const token = viewerPersonaID ? getAccessToken() : undefined;
      const res = await searchNox(searchTerm, results.limit, token || undefined, viewerPersonaID || undefined, results.next_offset);
      if (res.data) {
        setResults({
          ...res.data,
          personas: [...results.personas, ...res.data.personas],
          posts: [...results.posts, ...res.data.posts],
          events: [...results.events, ...res.data.events],
          hashtags: [...(results.hashtags ?? []), ...(res.data.hashtags ?? [])],
        });
      }
    } catch { /* ignore */ }
    finally { setLoadingMore(false); }
  }

  async function handleToggleLike(post: Post) {
    const token = getAccessToken();
    if (!token || !viewerPersonaID || !results) return;
    const prev = results;
    const nextLiked = !post.is_liked;
    setResults({ ...results, posts: results.posts.map((p) => p.id === post.id ? { ...p, is_liked: nextLiked, like_count: Math.max(0, p.like_count + (nextLiked ? 1 : -1)) } : p) });
    try {
      if (nextLiked) await likePost(post.id, viewerPersonaID, token);
      else await unlikePost(post.id, viewerPersonaID, token);
    } catch { setResults(prev); }
  }

  async function handleToggleFollow(persona: Persona) {
    const token = getAccessToken();
    if (!token || !viewerPersonaID || !results || persona.id === viewerPersonaID) return;
    const prev = results;
    const nextFollowing = !persona.is_following;
    setResults({ ...results, personas: results.personas.map((p) => p.id === persona.id ? { ...p, is_following: nextFollowing, follower_count: Math.max(0, p.follower_count + (nextFollowing ? 1 : -1)) } : p) });
    try {
      if (nextFollowing) await followPersona(persona.id, viewerPersonaID, token);
      else await unfollowPersona(persona.id, viewerPersonaID, token);
    } catch { setResults(prev); }
  }

  const hasResults = Boolean(
    results &&
    (results.personas.length > 0 || results.posts.length > 0 || results.events.length > 0 || (results.hashtags?.length ?? 0) > 0),
  );

  return (
    <FeedShell>
      <header className="px-4 pt-[env(safe-area-inset-top,12px)] pb-3">
        <h1 className="mb-3 text-[22px] font-bold tracking-[-0.03em] text-(--nox-ink)">discover</h1>
        <div className="flex items-center gap-3 rounded-[12px] border border-(--nox-border) bg-(--nox-surface) px-3.5 py-3 transition focus-within:border-(--nox-accent-line)">
          <Search className="size-4 shrink-0 text-(--nox-ink-soft)" strokeWidth={1.7} />
          <input
            type="text"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="search artists, handles, scenes…"
            className="flex-1 bg-transparent text-[14px] text-(--nox-ink) outline-none placeholder:text-(--nox-ink-soft)"
          />
        </div>
      </header>

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
          <p className="px-4 py-8 text-center text-[13px] text-(--nox-ink-soft)">Searching...</p>
        ) : hasResults && results ? (
          <div className="pb-4">
            {(results.hashtags?.length ?? 0) > 0 && (
              <section className="border-b border-(--nox-divider) pb-3">
                <p className="px-4 pb-2 pt-3 font-mono text-[10px] font-semibold uppercase tracking-[0.16em] text-(--nox-ink-soft)">scenes</p>
                <div className="flex flex-wrap gap-2 px-4">
                  {results.hashtags!.map((h) => (
                    <button
                      key={h.id}
                      type="button"
                      onClick={() => router.push(`/hashtags/${encodeURIComponent(h.tag)}`)}
                      className="flex items-center gap-1 rounded-full border border-(--nox-border) px-3 py-1.5 font-mono text-[11px] text-(--nox-ink-mid) transition hover:border-(--nox-accent-line) hover:text-(--nox-accent-ink)"
                    >
                      <span style={{ color: "var(--nox-accent-ink)" }}>#</span>{h.tag}
                      <span className="ml-0.5 text-[10px] text-(--nox-ink-soft)">{h.post_count}</span>
                    </button>
                  ))}
                </div>
              </section>
            )}

            {results.personas.length > 0 && (
              <section>
                <p className="px-4 pb-2 pt-3 font-mono text-[10px] font-semibold uppercase tracking-[0.16em] text-(--nox-ink-soft)">artists</p>
                <div className="divide-y divide-(--nox-divider)">
                  {results.personas.map((persona) => (
                    <PersonaCard
                      key={persona.id}
                      persona={persona}
                      onPress={() => router.push(`/personas/${persona.id}`)}
                      showFollow={Boolean(viewerPersonaID && persona.id !== viewerPersonaID)}
                      isFollowing={Boolean(persona.is_following)}
                      onFollow={handleToggleFollow}
                    />
                  ))}
                </div>
              </section>
            )}

            {results.posts.length > 0 && (
              <section className="mt-1">
                <p className="px-4 pb-2 pt-2 font-mono text-[10px] font-semibold uppercase tracking-[0.16em] text-(--nox-ink-soft)">posts</p>
                {results.posts.map((post) => (
                  <PostCard key={post.id} post={post} liked={post.is_liked} onLike={handleToggleLike} onClick={() => router.push(`/posts/${post.id}`)} />
                ))}
              </section>
            )}

            {results.events.length > 0 && (
              <section className="mt-1">
                <p className="px-4 pb-2 pt-2 font-mono text-[10px] font-semibold uppercase tracking-[0.16em] text-(--nox-ink-soft)">events</p>
                {results.events.map((event) => (
                  <EventCard key={event.id} event={event} />
                ))}
              </section>
            )}

            {results.has_more && (
              <div className="px-4 pt-4">
                <button
                  type="button"
                  onClick={handleLoadMore}
                  disabled={loadingMore}
                  className="w-full rounded-[8px] border border-(--nox-border-strong) px-4 py-2.5 text-[13px] font-semibold text-(--nox-ink) transition hover:border-(--nox-accent) disabled:opacity-50"
                >
                  {loadingMore ? "Loading..." : "Load more"}
                </button>
              </div>
            )}

            <TrendingHashtags />
          </div>
        ) : (
          <TrendingHashtags />
        )}
      </div>

      <TabBar />
    </FeedShell>
  );
}
