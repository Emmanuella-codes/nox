"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { ChevronLeft } from "lucide-react";
import { FeedShell } from "@/src/components/user/feed/feed-shell";
import { TabBar } from "@/src/components/user/feed/tab-bar";
import { PostCard } from "@/src/components/user/feed/post-card";
import { Avatar } from "@/src/components/user/shared/avatar";
import { getPersona } from "@/src/utils/api/user/persona";
import { getPersonaPosts, likePost, unlikePost } from "@/src/utils/api/user/post";
import { followPersona, unfollowPersona, getFollowStatus } from "@/src/utils/api/user/follow";
import { useActivePersona } from "@/src/hooks/use-active-persona";
import { getAccessToken } from "@/src/utils/auth/session";
import type { Persona } from "@/src/types/api/persona";
import type { Post } from "@/src/types/api/post";

function formatCount(n: number): string {
  if (n >= 1000) return `${(n / 1000).toFixed(1)}k`;
  return String(n);
}

interface PersonaProfileScreenProps {
  personaID: string;
}

export function PersonaProfileScreen({ personaID }: PersonaProfileScreenProps) {
  const router = useRouter();
  const { activeID: viewerPersonaID, loading: personaLoading } = useActivePersona();
  const isOwnProfile = personaID === viewerPersonaID;

  const [persona, setPersona] = useState<Persona | null>(null);
  const [posts, setPosts] = useState<Post[]>([]);
  const [loading, setLoading] = useState(true);
  const [following, setFollowing] = useState(false);
  const [followLoading, setFollowLoading] = useState(false);
  const [message, setMessage] = useState("");

  useEffect(() => {
    if (personaLoading) return;

    setLoading(true);
    setMessage("");

    const token = viewerPersonaID ? getAccessToken() : undefined;
    Promise.all([
      getPersona(personaID),
      getPersonaPosts(personaID, token || undefined, viewerPersonaID || undefined),
    ])
      .then(async ([personaRes, postsRes]) => {
        const p = personaRes.data ?? null;
        setPersona(p);
        setPosts(postsRes.data ?? []);
        if (p) {
          setFollowing(Boolean(p.is_following));
          if (viewerPersonaID && viewerPersonaID !== p.id) {
            if (token) {
              try {
                const statusRes = await getFollowStatus(p.id, viewerPersonaID, token);
                setFollowing(Boolean(statusRes.data?.is_following));
              } catch {
                // Public profile still renders if status hydration fails.
              }
            }
          }
        }
      })
      .catch(() => setMessage("Could not load profile."))
      .finally(() => setLoading(false));
  }, [personaID, viewerPersonaID, personaLoading]);

  async function handleToggleFollow() {
    const token = getAccessToken();
    if (!token || !viewerPersonaID || !persona || followLoading) return;

    const nextFollowing = !following;
    setFollowing(nextFollowing);
    setFollowLoading(true);
    try {
      if (nextFollowing) {
        await followPersona(persona.id, viewerPersonaID, token);
      } else {
        await unfollowPersona(persona.id, viewerPersonaID, token);
      }
      setPersona((prev) =>
        prev ? { ...prev, follower_count: Math.max(0, prev.follower_count + (nextFollowing ? 1 : -1)) } : prev,
      );
    } catch {
      setFollowing(!nextFollowing);
    } finally {
      setFollowLoading(false);
    }
  }

  async function handleToggleLike(post: Post) {
    const token = getAccessToken();
    if (!token || !viewerPersonaID) return;

    const nextLiked = !post.is_liked;
    setPosts((prev) =>
      prev.map((p) =>
        p.id === post.id
          ? { ...p, is_liked: nextLiked, like_count: Math.max(0, p.like_count + (nextLiked ? 1 : -1)) }
          : p,
      ),
    );
    try {
      if (nextLiked) {
        await likePost(post.id, viewerPersonaID, token);
      } else {
        await unlikePost(post.id, viewerPersonaID, token);
      }
    } catch {
      setPosts((prev) =>
        prev.map((p) =>
          p.id === post.id
            ? { ...p, is_liked: !nextLiked, like_count: Math.max(0, p.like_count + (nextLiked ? -1 : 1)) }
            : p,
        ),
      );
    }
  }

  return (
    <FeedShell>
      <div className="flex-1 overflow-y-auto">
        {/* Cover strip */}
        <div className="h-24 w-full" style={{ background: "linear-gradient(135deg, #1a1028 0%, #0d0d14 100%)" }}>
          <button
            type="button"
            onClick={() => router.back()}
            className="m-3 flex size-8 items-center justify-center rounded-full transition"
            style={{ background: "rgba(0,0,0,.45)" }}
          >
            <ChevronLeft className="size-5 text-white" strokeWidth={1.8} />
          </button>
        </div>

        {/* Profile header */}
        <div className="px-4">
          <div className="flex items-end justify-between" style={{ marginTop: -28 }}>
            <div className="overflow-hidden rounded-[16px] border-2" style={{ borderColor: "var(--nox-bg)" }}>
              {loading || !persona ? (
                <div className="size-14 animate-pulse rounded-[14px] bg-(--nox-surface-alt)" />
              ) : (
                <Avatar id={persona.id} name={persona.display_name} size={56} square />
              )}
            </div>

            {!loading && persona && !isOwnProfile && viewerPersonaID && (
              <button
                type="button"
                onClick={handleToggleFollow}
                disabled={followLoading}
                className="mb-1 rounded-[10px] border px-4 py-1.5 text-[13px] font-semibold transition disabled:opacity-60"
                style={{
                  borderColor: following ? "var(--nox-accent-line)" : "var(--nox-border-strong)",
                  color: following ? "var(--nox-accent-ink)" : "var(--nox-ink)",
                  background: following ? "var(--nox-accent-soft)" : "transparent",
                }}
              >
                {following ? "following" : "follow"}
              </button>
            )}

            {isOwnProfile && (
              <div className="mb-1 flex gap-2">
                {persona?.category === "dj" ? (
                  <button
                    type="button"
                    onClick={() => router.push("/sets/create")}
                    className="rounded-[8px] border border-(--nox-accent-line) px-3 py-1.5 text-[12px] font-semibold text-(--nox-accent-ink) transition hover:border-(--nox-accent)"
                  >
                    new set
                  </button>
                ) : null}
                <button
                  type="button"
                  onClick={() => router.push("/settings")}
                  className="rounded-[8px] border border-(--nox-border-strong) px-3 py-1.5 text-[12px] font-semibold text-(--nox-ink) transition hover:border-(--nox-accent)"
                >
                  edit
                </button>
              </div>
            )}
          </div>

          {loading || !persona ? (
            <div className="mt-3 space-y-2">
              <div className="h-5 w-2/5 animate-pulse rounded bg-(--nox-surface-alt)" />
              <div className="h-3 w-1/4 animate-pulse rounded bg-(--nox-surface-alt)" />
            </div>
          ) : (
            <>
              <div className="mt-3">
                <h1 className="text-[20px] font-bold tracking-[-0.03em] text-(--nox-ink)">{persona.display_name}</h1>
                <p className="text-[13px] text-(--nox-ink-soft)">@{persona.handle}</p>
                <p className="mt-1 font-mono text-[10px] font-semibold uppercase tracking-[0.14em] text-(--nox-accent-ink)">
                  {persona.category}
                </p>
              </div>

              {persona.bio && (
                <p className="mt-2 text-[13px] leading-[1.55] text-(--nox-ink-mid)">{persona.bio}</p>
              )}

              {persona.genre_tags.length > 0 && (
                <div className="mt-2.5 flex flex-wrap gap-1.5">
                  {persona.genre_tags.map((tag) => (
                    <span
                      key={tag}
                      className="rounded-full px-2.5 py-1 font-mono text-[10px] font-medium lowercase"
                      style={{ background: "var(--nox-accent-soft)", color: "var(--nox-accent-ink)" }}
                    >
                      {tag}
                    </span>
                  ))}
                </div>
              )}

              <div className="mt-4 flex gap-5 border-b border-(--nox-divider) pb-4">
                <div>
                  <p className="text-[16px] font-bold text-(--nox-ink)">{formatCount(persona.post_count)}</p>
                  <p className="font-mono text-[10px] text-(--nox-ink-soft)">posts</p>
                </div>
                <button
                  type="button"
                  onClick={() => router.push(`/personas/${persona.id}/followers`)}
                  className="text-left transition hover:opacity-70"
                >
                  <p className="text-[16px] font-bold text-(--nox-ink)">{formatCount(persona.follower_count)}</p>
                  <p className="font-mono text-[10px] text-(--nox-ink-soft)">followers</p>
                </button>
                <button
                  type="button"
                  onClick={() => router.push(`/personas/${persona.id}/following`)}
                  className="text-left transition hover:opacity-70"
                >
                  <p className="text-[16px] font-bold text-(--nox-ink)">{formatCount(persona.following_count)}</p>
                  <p className="font-mono text-[10px] text-(--nox-ink-soft)">following</p>
                </button>
              </div>
            </>
          )}
        </div>

        {/* Posts */}
        {!message && (
          <div className="mt-1">
            {loading ? (
              <div>
                {[...Array(3)].map((_, i) => (
                  <div key={i} className="border-b border-(--nox-divider) px-4 py-4">
                    <div className="flex gap-3">
                      <div className="size-9 animate-pulse rounded-full bg-(--nox-surface-alt)" />
                      <div className="flex-1 space-y-2">
                        <div className="h-3 w-1/3 animate-pulse rounded bg-(--nox-surface-alt)" />
                        <div className="h-3 w-full animate-pulse rounded bg-(--nox-surface-alt)" />
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            ) : posts.length > 0 ? (
              posts.map((p) => (
                <PostCard
                  key={p.id}
                  post={p}
                  liked={p.is_liked}
                  onLike={handleToggleLike}
                  onClick={() => router.push(`/posts/${p.id}`)}
                />
              ))
            ) : (
              <p className="py-12 text-center text-[13px] text-(--nox-ink-soft)">no posts yet.</p>
            )}
          </div>
        )}

        {message && (
          <p className="px-4 py-10 text-center text-[13px] text-(--nox-ink-soft)">{message}</p>
        )}
      </div>

      <TabBar />
    </FeedShell>
  );
}
