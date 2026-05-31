"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { FeedShell } from "@/src/components/user/feed/feed-shell";
import { FeedTabs, type FeedTab } from "@/src/components/user/feed/feed-tabs";
import { PostCard } from "@/src/components/user/feed/post-card";
import { ComposeBar } from "@/src/components/user/feed/compose-bar";
import { TabBar } from "@/src/components/user/feed/tab-bar";
import { PersonaSwitcher } from "@/src/components/user/shared/persona-switcher";
import { getPersonaFeed, getFollowingFeed, likePost, unlikePost } from "@/src/utils/api/user/post";
import { useActivePersona } from "@/src/hooks/use-active-persona";
import { getAccessToken } from "@/src/utils/auth/session";
import type { Post } from "@/src/types/api/post";

export function FeedScreen() {
  const router = useRouter();
  const { personas, activeID, activePersona, loading: personaLoading, switchPersona } = useActivePersona();
  const [tab, setTab] = useState<FeedTab>("for-you");
  const [posts, setPosts] = useState<Post[]>([]);
  const [loading, setLoading] = useState(true);
  const [message, setMessage] = useState("");
  const [likedPostIDs, setLikedPostIDs] = useState<Set<string>>(new Set());

  useEffect(() => {
    async function loadFeed() {
      if (personaLoading) return;

      if (!activePersona) {
        setMessage("Create a public persona to load the feed.");
        setLoading(false);
        return;
      }

      if (tab === "events") {
        router.push("/events");
        return;
      }

      if (tab === "sets") {
        setPosts([]);
        setMessage("Set archive coming soon.");
        setLoading(false);
        return;
      }

      const token = getAccessToken();
      if (!token) {
        setMessage("Sign in to load your feed.");
        setLoading(false);
        return;
      }

      setLoading(true);
      setMessage("");

      try {
        const res =
          tab === "following"
            ? await getFollowingFeed(activePersona.id, token)
            : await getPersonaFeed(activePersona.id, token);
        const nextPosts = res.data ?? [];
        setPosts(nextPosts);
        setLikedPostIDs(new Set(nextPosts.filter((p) => p.is_liked).map((p) => p.id)));
        if (nextPosts.length === 0 && tab === "following") {
          setMessage("Follow some people to see their posts here.");
        }
      } catch {
        setMessage("Could not load feed.");
      } finally {
        setLoading(false);
      }
    }

    void loadFeed();
  }, [activeID, tab, personaLoading, activePersona, router]);

  async function handleToggleLike(post: Post) {
    const token = getAccessToken();
    if (!token || !activePersona) return;

    const nextLikedIDs = new Set(likedPostIDs);
    const liked = nextLikedIDs.has(post.id);
    if (liked) {
      nextLikedIDs.delete(post.id);
    } else {
      nextLikedIDs.add(post.id);
    }
    setLikedPostIDs(nextLikedIDs);
    setPosts((current) =>
      current.map((item) =>
        item.id === post.id
          ? { ...item, like_count: Math.max(0, item.like_count + (liked ? -1 : 1)) }
          : item,
      ),
    );

    try {
      if (liked) {
        await unlikePost(post.id, activePersona.id, token);
      } else {
        await likePost(post.id, activePersona.id, token);
      }
    } catch {
      setLikedPostIDs(likedPostIDs);
      setPosts((current) =>
        current.map((item) =>
          item.id === post.id
            ? { ...item, like_count: Math.max(0, item.like_count + (liked ? 1 : -1)) }
            : item,
        ),
      );
    }
  }

  return (
    <FeedShell>
      {/* Top bar */}
      <header className="flex items-center justify-between px-4 pt-[env(safe-area-inset-top,12px)] pb-2">
        <span
          className="font-mono text-[22px] font-bold tracking-[-0.04em]"
          style={{ color: "var(--nox-accent)" }}
        >
          nox
        </span>
        <PersonaSwitcher
          personas={personas}
          activePersona={activePersona}
          onSwitch={switchPersona}
        />
      </header>

      {/* Feed tabs */}
      <FeedTabs active={tab} onChange={setTab} />

      {/* Post list */}
      <div className="flex-1 overflow-y-auto">
        {loading ? (
          <div className="flex flex-col gap-0">
            {[...Array(4)].map((_, i) => (
              <div key={i} className="border-b border-(--nox-divider) px-4 py-4">
                <div className="flex gap-3">
                  <div className="size-9 animate-pulse rounded-full bg-(--nox-surface-alt)" />
                  <div className="flex-1 space-y-2">
                    <div className="h-3 w-1/3 animate-pulse rounded bg-(--nox-surface-alt)" />
                    <div className="h-3 w-full animate-pulse rounded bg-(--nox-surface-alt)" />
                    <div className="h-3 w-4/5 animate-pulse rounded bg-(--nox-surface-alt)" />
                  </div>
                </div>
              </div>
            ))}
          </div>
        ) : posts.length > 0 ? (
          posts.map((post) => (
            <PostCard
              key={post.id}
              post={post}
              liked={likedPostIDs.has(post.id)}
              onLike={handleToggleLike}
              onClick={() => router.push(`/posts/${post.id}`)}
            />
          ))
        ) : (
          <p className="px-4 py-10 text-center text-[13px] text-(--nox-ink-soft)">
            {message || "No posts yet."}
          </p>
        )}
      </div>

      {/* Compose bar */}
      <ComposeBar onClick={() => router.push("/compose")} />

      {/* Tab bar */}
      <TabBar />
    </FeedShell>
  );
}
