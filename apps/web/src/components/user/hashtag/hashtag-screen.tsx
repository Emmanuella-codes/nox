"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { ChevronLeft, Hash } from "lucide-react";
import { FeedShell } from "@/src/components/user/feed/feed-shell";
import { PostCard } from "@/src/components/user/feed/post-card";
import { TabBar } from "@/src/components/user/feed/tab-bar";
import { getHashtag, getHashtagPosts } from "@/src/utils/api/user/hashtag";
import { likePost, unlikePost } from "@/src/utils/api/user/post";
import { useActivePersona } from "@/src/hooks/use-active-persona";
import { getAccessToken } from "@/src/utils/auth/session";
import type { HashtagDetail } from "@/src/types/api/hashtag";
import type { Post } from "@/src/types/api/post";

interface HashtagScreenProps {
  tag: string;
}

export function HashtagScreen({ tag }: HashtagScreenProps) {
  const router = useRouter();
  const normalizedTag = tag.replace(/^#/, "").toLowerCase();
  const { activeID: viewerPersonaID, loading: personaLoading } = useActivePersona();

  const [detail, setDetail] = useState<HashtagDetail | null>(null);
  const [posts, setPosts] = useState<Post[]>([]);
  const [hasMore, setHasMore] = useState(false);
  const [nextOffset, setNextOffset] = useState(0);
  const [loading, setLoading] = useState(true);
  const [loadingMore, setLoadingMore] = useState(false);
  const [message, setMessage] = useState("");

  useEffect(() => {
    if (personaLoading) return;

    setLoading(true);
    setMessage("");

    const token = viewerPersonaID ? getAccessToken() || undefined : undefined;

    Promise.all([
      getHashtag(normalizedTag),
      getHashtagPosts(normalizedTag, 30, 0, viewerPersonaID || undefined, token),
    ])
      .then(([detailRes, postsRes]) => {
        setDetail(detailRes.data ?? null);
        const data = postsRes.data;
        if (data) {
          setPosts(data.posts);
          setHasMore(data.has_more);
          setNextOffset(data.next_offset ?? 0);
        }
      })
      .catch(() => setMessage("Could not load this hashtag."))
      .finally(() => setLoading(false));
  }, [normalizedTag, viewerPersonaID, personaLoading]);

  async function handleLoadMore() {
    if (!hasMore || loadingMore) return;
    setLoadingMore(true);
    const token = viewerPersonaID ? getAccessToken() || undefined : undefined;
    try {
      const res = await getHashtagPosts(normalizedTag, 30, nextOffset, viewerPersonaID || undefined, token);
      const data = res.data;
      if (data) {
        setPosts((prev) => [...prev, ...data.posts]);
        setHasMore(data.has_more);
        setNextOffset(data.next_offset ?? 0);
      }
    } catch { /* ignore */ }
    finally { setLoadingMore(false); }
  }

  async function handleToggleLike(post: Post) {
    const token = getAccessToken();
    if (!token || !viewerPersonaID) return;

    const nextLiked = !post.is_liked;
    setPosts((current) =>
      current.map((p) =>
        p.id === post.id
          ? { ...p, is_liked: nextLiked, like_count: Math.max(0, p.like_count + (nextLiked ? 1 : -1)) }
          : p,
      ),
    );
    try {
      if (nextLiked) await likePost(post.id, viewerPersonaID, token);
      else await unlikePost(post.id, viewerPersonaID, token);
    } catch {
      setPosts((current) =>
        current.map((p) =>
          p.id === post.id
            ? { ...p, is_liked: !nextLiked, like_count: Math.max(0, p.like_count + (nextLiked ? -1 : 1)) }
            : p,
        ),
      );
    }
  }

  const postCount = detail?.post_count ?? posts.length;

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
              {postCount} post{postCount === 1 ? "" : "s"}
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
            {hasMore && (
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
