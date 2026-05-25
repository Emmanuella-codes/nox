"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { Ghost } from "lucide-react";
import { FeedShell } from "@/src/components/user/feed/feed-shell";
import { FeedTabs, type FeedTab } from "@/src/components/user/feed/feed-tabs";
import { PostCard } from "@/src/components/user/feed/post-card";
import { ComposeBar } from "@/src/components/user/feed/compose-bar";
import { TabBar } from "@/src/components/user/feed/tab-bar";
import { getPersonaFeed, likePost, unlikePost } from "@/src/utils/api/user/post";
import { getMyPersonas } from "@/src/utils/api/user/persona";
import { getAccessToken } from "@/src/utils/auth/session";
import type { Post } from "@/src/types/api/post";

export function FeedScreen() {
  const router = useRouter();
  const [tab, setTab] = useState<FeedTab>("for-you");
  const [posts, setPosts] = useState<Post[]>([]);
  const [loading, setLoading] = useState(true);
  const [message, setMessage] = useState("");
  const [viewerPersonaID, setViewerPersonaID] = useState("");
  const [likedPostIDs, setLikedPostIDs] = useState<Set<string>>(new Set());

  useEffect(() => {
    async function loadFeed() {
      setLoading(true);
      setMessage("");

      try {
        const token = getAccessToken();
        if (!token) {
          setMessage("Sign in to load your feed.");
          return;
        }

        const personasRes = await getMyPersonas(token);
        const primaryPersona = personasRes.data?.[0];
        if (!primaryPersona) {
          setMessage("Create a public persona to load the feed.");
          return;
        }
        setViewerPersonaID(primaryPersona.id);

        if (tab === "events") {
          router.push("/events");
          return;
        }

        if (tab === "sets" || tab === "following") {
          setPosts([]);
          setMessage(`${tab.replace("-", " ")} is not connected yet.`);
          return;
        }

        const feedRes = await getPersonaFeed(primaryPersona.id, token);
        const nextPosts = feedRes.data ?? [];
        setPosts(nextPosts);
        setLikedPostIDs(new Set(nextPosts.filter((post) => post.is_liked).map((post) => post.id)));
      } catch {
        setMessage("Could not load feed.");
      } finally {
        setLoading(false);
      }
    }

    loadFeed();
  }, [router, tab]);

  async function handleToggleLike(post: Post) {
    const token = getAccessToken();
    if (!token || !viewerPersonaID) return;

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
        await unlikePost(post.id, viewerPersonaID, token);
      } else {
        await likePost(post.id, viewerPersonaID, token);
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
        <button
          type="button"
          className="flex items-center gap-1.5 rounded-full border px-3 py-1.5 font-mono text-[10px] font-semibold uppercase tracking-[0.12em] transition"
          style={{
            borderColor: "var(--nox-border-strong)",
            color: "var(--nox-ink-mid)",
          }}
        >
          <Ghost className="size-3" strokeWidth={1.8} />
          anonymous
        </button>
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
                  <div className="size-9 rounded-full bg-(--nox-surface-alt) animate-pulse" />
                  <div className="flex-1 space-y-2">
                    <div className="h-3 w-1/3 rounded bg-(--nox-surface-alt) animate-pulse" />
                    <div className="h-3 w-full rounded bg-(--nox-surface-alt) animate-pulse" />
                    <div className="h-3 w-4/5 rounded bg-(--nox-surface-alt) animate-pulse" />
                  </div>
                </div>
              </div>
            ))}
          </div>
        ) : (
          posts.length > 0 ? (
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
          )
        )}
      </div>

      {/* Compose bar */}
      <ComposeBar onClick={() => router.push("/compose")} />

      {/* Tab bar */}
      <TabBar />
    </FeedShell>
  );
}
