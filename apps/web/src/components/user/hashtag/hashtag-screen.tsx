"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { ChevronLeft, Hash } from "lucide-react";
import { FeedShell } from "@/src/components/user/feed/feed-shell";
import { PostCard } from "@/src/components/user/feed/post-card";
import { TabBar } from "@/src/components/user/feed/tab-bar";
import { getHashtag, getHashtagPosts } from "@/src/utils/api/user/hashtag";
import { getMyPersonas } from "@/src/utils/api/user/persona";
import { likePost, unlikePost } from "@/src/utils/api/user/post";
import { getAccessToken } from "@/src/utils/auth/session";
import type { HashtagDetail } from "@/src/types/api/hashtag";
import type { Post } from "@/src/types/api/post";

interface HashtagScreenProps {
  tag: string;
}

export function HashtagScreen({ tag }: HashtagScreenProps) {
  const router = useRouter();
  const normalizedTag = tag.replace(/^#/, "").toLowerCase();
  const [detail, setDetail] = useState<HashtagDetail | null>(null);
  const [posts, setPosts] = useState<Post[]>([]);
  const [viewerPersonaID, setViewerPersonaID] = useState("");
  const [loading, setLoading] = useState(true);
  const [message, setMessage] = useState("");

  useEffect(() => {
    async function load() {
      setLoading(true);
      setMessage("");
      try {
        const token = getAccessToken();
        if (token) {
          try {
            const personasRes = await getMyPersonas(token);
            setViewerPersonaID(personasRes.data?.[0]?.id ?? "");
          } catch {
            setViewerPersonaID("");
          }
        }

        const [detailRes, postsRes] = await Promise.all([
          getHashtag(normalizedTag),
          getHashtagPosts(normalizedTag, 30),
        ]);
        setDetail(detailRes.data ?? null);
        setPosts(postsRes.data?.posts ?? []);
      } catch {
        setMessage("Could not load this hashtag.");
      } finally {
        setLoading(false);
      }
    }

    load();
  }, [normalizedTag]);

  async function handleToggleLike(post: Post) {
    const token = getAccessToken();
    if (!token || !viewerPersonaID) return;

    const previousPosts = posts;
    const nextLiked = !post.is_liked;
    setPosts((current) =>
      current.map((item) =>
        item.id === post.id
          ? {
              ...item,
              is_liked: nextLiked,
              like_count: Math.max(0, item.like_count + (nextLiked ? 1 : -1)),
            }
          : item,
      ),
    );

    try {
      if (nextLiked) {
        await likePost(post.id, viewerPersonaID, token);
      } else {
        await unlikePost(post.id, viewerPersonaID, token);
      }
    } catch {
      setPosts(previousPosts);
    }
  }

  return (
    <FeedShell>
      <header className="border-b border-(--nox-divider) px-4 pt-[env(safe-area-inset-top,12px)] pb-4">
        <div className="flex items-center gap-3">
          <button
            type="button"
            onClick={() => router.back()}
            className="flex size-8 items-center justify-center rounded-full transition hover:bg-(--nox-surface-alt)"
          >
            <ChevronLeft className="size-5 text-(--nox-ink)" strokeWidth={1.8} />
          </button>
          <div className="flex size-9 items-center justify-center rounded-full bg-(--nox-surface-alt)">
            <Hash className="size-4 text-(--nox-accent-ink)" strokeWidth={1.8} />
          </div>
          <div className="min-w-0">
            <h1 className="truncate text-[20px] font-bold tracking-[-0.03em] text-(--nox-ink)">
              #{normalizedTag}
            </h1>
            <p className="text-[12px] text-(--nox-ink-soft)">
              {detail?.post_count ?? posts.length} post{(detail?.post_count ?? posts.length) === 1 ? "" : "s"}
            </p>
          </div>
        </div>
      </header>

      <div className="flex-1 overflow-y-auto">
        {loading ? (
          <p className="px-4 py-8 text-center text-[13px] text-(--nox-ink-soft)">Loading...</p>
        ) : message ? (
          <p className="px-4 py-8 text-center text-[13px] text-(--nox-danger)">{message}</p>
        ) : posts.length > 0 ? (
          <div className="pb-4">
            {posts.map((post) => (
              <PostCard
                key={post.id}
                post={post}
                liked={post.is_liked}
                onLike={handleToggleLike}
                onClick={() => router.push(`/posts/${post.id}`)}
              />
            ))}
          </div>
        ) : (
          <p className="px-4 py-8 text-center text-[13px] text-(--nox-ink-soft)">
            No posts under this tag yet.
          </p>
        )}
      </div>

      <TabBar />
    </FeedShell>
  );
}
