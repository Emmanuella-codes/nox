"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { ChevronLeft } from "lucide-react";
import { FeedShell } from "@/src/components/user/feed/feed-shell";
import { TabBar } from "@/src/components/user/feed/tab-bar";
import { PersonaCard } from "@/src/components/user/discover/persona-card";
import { getFollowers, getFollowing, followPersona, unfollowPersona } from "@/src/utils/api/user/follow";
import { useActivePersona } from "@/src/hooks/use-active-persona";
import { getAccessToken } from "@/src/utils/auth/session";
import type { Persona } from "@/src/types/api/persona";

interface FollowListScreenProps {
  personaID: string;
  mode: "followers" | "following";
}

export function FollowListScreen({ personaID, mode }: FollowListScreenProps) {
  const router = useRouter();
  const [personas, setPersonas] = useState<Persona[]>([]);
  const [hasMore, setHasMore] = useState(false);
  const [nextOffset, setNextOffset] = useState(0);
  const [loading, setLoading] = useState(true);
  const [loadingMore, setLoadingMore] = useState(false);
  const [message, setMessage] = useState("");
  const { activeID: viewerPersonaID, loading: personaLoading } = useActivePersona();

  useEffect(() => {
    if (personaLoading) return;

    setLoading(true);
    setMessage("");

    const fn = mode === "followers" ? getFollowers : getFollowing;
    const token = viewerPersonaID ? getAccessToken() : undefined;
    fn(personaID, 20, 0, viewerPersonaID || undefined, token || undefined)
      .then((res) => {
        const data = res.data;
        if (!data) {
          setMessage("Could not load list.");
          return;
        }
        setPersonas(data.personas);
        setHasMore(data.has_more);
        setNextOffset(data.next_offset ?? 0);
      })
      .catch(() => setMessage("Could not load list."))
      .finally(() => setLoading(false));
  }, [personaID, mode, personaLoading, viewerPersonaID]);

  async function handleLoadMore() {
    if (!hasMore || loadingMore) return;
    setLoadingMore(true);
    const fn = mode === "followers" ? getFollowers : getFollowing;
    try {
      const token = viewerPersonaID ? getAccessToken() : undefined;
      const res = await fn(personaID, 20, nextOffset, viewerPersonaID || undefined, token || undefined);
      const data = res.data;
      if (data) {
        setPersonas((prev) => [...prev, ...data.personas]);
        setHasMore(data.has_more);
        setNextOffset(data.next_offset ?? 0);
      }
    } catch {
      // Keep the existing list visible if loading more fails.
    } finally {
      setLoadingMore(false);
    }
  }

  async function handleToggleFollow(persona: Persona) {
    const token = getAccessToken();
    if (!token || !viewerPersonaID || persona.id === viewerPersonaID) return;

    const nextFollowing = !persona.is_following;
    setPersonas((prev) =>
      prev.map((p) =>
        p.id === persona.id
          ? { ...p, is_following: nextFollowing, follower_count: Math.max(0, p.follower_count + (nextFollowing ? 1 : -1)) }
          : p,
      ),
    );
    try {
      if (nextFollowing) {
        await followPersona(persona.id, viewerPersonaID, token);
      } else {
        await unfollowPersona(persona.id, viewerPersonaID, token);
      }
    } catch {
      setPersonas((prev) =>
        prev.map((p) =>
          p.id === persona.id
            ? { ...p, is_following: !nextFollowing, follower_count: Math.max(0, p.follower_count + (nextFollowing ? -1 : 1)) }
            : p,
        ),
      );
    }
  }

  return (
    <FeedShell>
      <header className="flex items-center gap-3 border-b border-(--nox-divider) px-4 pt-[env(safe-area-inset-top,12px)] pb-4">
        <button
          type="button"
          onClick={() => router.back()}
          className="flex size-8 items-center justify-center rounded-full transition hover:bg-(--nox-surface-alt)"
        >
          <ChevronLeft className="size-5 text-(--nox-ink)" strokeWidth={1.8} />
        </button>
        <h1 className="text-[18px] font-bold tracking-[-0.03em] text-(--nox-ink)">{mode}</h1>
      </header>

      <div className="flex-1 overflow-y-auto">
        {loading ? (
          <div className="divide-y divide-(--nox-divider)">
            {[...Array(6)].map((_, i) => (
              <div key={i} className="flex items-center gap-3 px-4 py-3">
                <div className="size-10 rounded-full bg-(--nox-surface-alt) animate-pulse" />
                <div className="flex-1 space-y-2">
                  <div className="h-3 w-1/3 rounded bg-(--nox-surface-alt) animate-pulse" />
                  <div className="h-3 w-1/4 rounded bg-(--nox-surface-alt) animate-pulse" />
                </div>
              </div>
            ))}
          </div>
        ) : message ? (
          <p className="px-4 py-10 text-center text-[13px] text-(--nox-ink-soft)">{message}</p>
        ) : personas.length === 0 ? (
          <p className="px-4 py-10 text-center text-[13px] text-(--nox-ink-soft)">
            {mode === "followers" ? "No followers yet." : "Not following anyone yet."}
          </p>
        ) : (
          <div className="divide-y divide-(--nox-divider)">
            {personas.map((persona) => (
              <PersonaCard
                key={persona.id}
                persona={persona}
                onPress={() => router.push(`/personas/${persona.id}`)}
                showFollow={Boolean(viewerPersonaID && persona.id !== viewerPersonaID)}
                isFollowing={Boolean(persona.is_following)}
                onFollow={handleToggleFollow}
              />
            ))}
            {hasMore && (
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
            )}
          </div>
        )}
      </div>

      <TabBar />
    </FeedShell>
  );
}
