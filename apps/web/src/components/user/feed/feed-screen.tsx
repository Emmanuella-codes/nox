"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { Ghost } from "lucide-react";
import { FeedShell } from "@/src/components/user/feed/feed-shell";
import { FeedTabs, type FeedTab } from "@/src/components/user/feed/feed-tabs";
import { PostCard } from "@/src/components/user/feed/post-card";
import { ComposeBar } from "@/src/components/user/feed/compose-bar";
import { TabBar } from "@/src/components/user/feed/tab-bar";
import { getPersonaFeed } from "@/src/utils/api/user/post";
import type { Post } from "@/src/types/api/post";

// Placeholder posts for loading / empty states
const PLACEHOLDER_POSTS: Post[] = [
  {
    id: "1",
    author: { mode: "anonymous" },
    body: "The 3am afro-house set at fabric last night was something else. #afrohouse #fabric",
    post_type: "text",
    like_count: 47,
    comment_count: 8,
    repost_count: 3,
    is_repost: false,
    created_at: new Date(Date.now() - 25 * 60_000).toISOString(),
  },
  {
    id: "2",
    author: {
      mode: "public",
      persona: {
        id: "p1",
        handle: "djkayode",
        display_name: "DJ Kayode",
        avatar_url: "",
      },
    },
    body: "New amapiano edit dropping Friday. First listen for everyone who shows up to the Lagos pop-up. #amapiano #lagos",
    post_type: "text",
    like_count: 123,
    comment_count: 22,
    repost_count: 14,
    is_repost: false,
    created_at: new Date(Date.now() - 2 * 3600_000).toISOString(),
  },
  {
    id: "3",
    author: { mode: "anonymous" },
    body: "Why do DJs still play the same afrobeats top-5 at every Shoreditch event? Where's the curation? #afrobeats",
    post_type: "text",
    like_count: 88,
    comment_count: 31,
    repost_count: 5,
    is_repost: false,
    created_at: new Date(Date.now() - 5 * 3600_000).toISOString(),
  },
  {
    id: "4",
    author: {
      mode: "public",
      persona: {
        id: "p2",
        handle: "nnekabeats",
        display_name: "Nneka Beats",
        avatar_url: "",
      },
    },
    body: "Soundcloud mix just hit 10k plays. Thank you all. Afro-soul only, always. #afrosoul",
    post_type: "text",
    like_count: 210,
    comment_count: 44,
    repost_count: 29,
    is_repost: false,
    created_at: new Date(Date.now() - 8 * 3600_000).toISOString(),
  },
];

export function FeedScreen() {
  const router = useRouter();
  const [tab, setTab] = useState<FeedTab>("for-you");
  const [posts, setPosts] = useState<Post[]>(PLACEHOLDER_POSTS);
  const [loading, setLoading] = useState(false);

  // Future: fetch real feed based on tab
  useEffect(() => {
    // Real fetch would go here when API endpoints are ready
  }, [tab]);

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
          ghost
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
          posts.map((post) => (
            <PostCard
              key={post.id}
              post={post}
              onClick={() => router.push(`/posts/${post.id}`)}
            />
          ))
        )}
      </div>

      {/* Compose bar */}
      <ComposeBar onClick={() => router.push("/compose")} />

      {/* Tab bar */}
      <TabBar />
    </FeedShell>
  );
}
