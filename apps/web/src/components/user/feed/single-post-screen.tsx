"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { ChevronLeft, Send, Ghost, Heart, MessageCircle, Repeat2 } from "lucide-react";
import { FeedShell } from "@/src/components/user/feed/feed-shell";
import { CommentItem } from "@/src/components/user/feed/comment-item";
import { getPost } from "@/src/utils/api/user/post";
import { getPostComments } from "@/src/utils/api/user/comment";
import type { Post } from "@/src/types/api/post";
import type { Comment } from "@/src/types/api/comment";

interface SinglePostScreenProps {
  postId: string;
}

function formatTime(isoString: string): string {
  const d = new Date(isoString);
  return d.toLocaleDateString("en-GB", {
    day: "numeric",
    month: "short",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

function formatCount(n: number): string {
  if (n >= 1000) return `${(n / 1000).toFixed(1)}k`;
  return String(n);
}

function extractHashtags(body: string) {
  const parts = body.split(/(#\w+)/g);
  return parts.map((p, i) =>
    p.startsWith("#") ? (
      <span key={i} style={{ color: "var(--nox-accent-ink)" }}>
        {p}
      </span>
    ) : (
      <span key={i}>{p}</span>
    ),
  );
}

export function SinglePostScreen({ postId }: SinglePostScreenProps) {
  const router = useRouter();
  const [post, setPost] = useState<Post | null>(null);
  const [comments, setComments] = useState<Comment[]>([]);
  const [commentBody, setCommentBody] = useState("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    async function load() {
      try {
        const [postRes, commentsRes] = await Promise.all([
          getPost(postId),
          getPostComments(postId),
        ]);
        setPost(postRes.data ?? null);
        setComments(commentsRes.data ?? []);
      } catch {
        setError("Could not load post.");
      } finally {
        setLoading(false);
      }
    }
    load();
  }, [postId]);

  const isAnon = post?.author.mode === "anonymous";
  const persona = post?.author.persona;

  return (
    <FeedShell>
      {/* Header */}
      <header className="flex items-center gap-3 border-b border-(--nox-divider) px-4 py-3 pt-[env(safe-area-inset-top,12px)]">
        <button
          type="button"
          onClick={() => router.back()}
          className="flex size-8 items-center justify-center rounded-full transition hover:bg-(--nox-surface-alt)"
        >
          <ChevronLeft className="size-5 text-(--nox-ink)" strokeWidth={1.8} />
        </button>
        <span className="text-[15px] font-semibold text-(--nox-ink)">post</span>
      </header>

      <div className="flex-1 overflow-y-auto">
        {loading && (
          <div className="space-y-3 px-4 py-5">
            <div className="flex gap-3">
              <div className="size-9 rounded-full bg-(--nox-surface-alt) animate-pulse" />
              <div className="flex-1 space-y-2">
                <div className="h-3 w-1/3 rounded bg-(--nox-surface-alt) animate-pulse" />
                <div className="h-3 w-full rounded bg-(--nox-surface-alt) animate-pulse" />
                <div className="h-3 w-4/5 rounded bg-(--nox-surface-alt) animate-pulse" />
              </div>
            </div>
          </div>
        )}

        {error && (
          <p className="px-4 py-5 text-[13px] text-(--nox-danger)">{error}</p>
        )}

        {post && (
          <>
            {/* Expanded post */}
            <div className="border-b border-(--nox-divider) px-4 py-4">
              <div className="flex items-center gap-3">
                {isAnon ? (
                  <div
                    className="flex size-10 items-center justify-center rounded-full"
                    style={{ background: "var(--nox-surface-alt)" }}
                  >
                    <Ghost className="size-5" strokeWidth={1.6} style={{ color: "var(--nox-ink-mid)" }} />
                  </div>
                ) : (
                  <div
                    className="flex size-10 items-center justify-center rounded-full text-[15px] font-bold"
                    style={{ background: "var(--nox-accent-soft)", color: "var(--nox-accent-ink)" }}
                  >
                    {persona?.display_name?.[0]?.toUpperCase() ?? "?"}
                  </div>
                )}
                <div>
                  {isAnon ? (
                    <span
                      className="flex items-center gap-1 rounded-full px-2 py-0.5 font-mono text-[9.5px] font-semibold uppercase tracking-[0.12em]"
                      style={{ background: "var(--nox-surface-alt)", color: "var(--nox-ink-mid)" }}
                    >
                      <Ghost className="size-2.5" strokeWidth={1.8} />
                      ghost
                    </span>
                  ) : (
                    <>
                      <p className="text-[14px] font-semibold text-(--nox-ink)">{persona?.display_name}</p>
                      <p className="text-[12px] text-(--nox-ink-soft)">@{persona?.handle}</p>
                    </>
                  )}
                </div>
              </div>

              <p className="mt-3 text-[16px] leading-[1.6] text-(--nox-ink)">
                {extractHashtags(post.body)}
              </p>

              <p className="mt-3 text-[11px] text-(--nox-ink-soft)">{formatTime(post.created_at)}</p>

              {/* Engagement row */}
              <div className="mt-3 flex items-center gap-6 border-t border-(--nox-divider) pt-3">
                <button
                  type="button"
                  className="flex items-center gap-1.5 text-[13px] text-(--nox-ink-soft) transition hover:text-(--nox-accent)"
                >
                  <Heart className="size-4" strokeWidth={1.7} />
                  {formatCount(post.like_count)}
                </button>
                <button
                  type="button"
                  className="flex items-center gap-1.5 text-[13px] text-(--nox-ink-soft) transition hover:text-(--nox-accent)"
                >
                  <MessageCircle className="size-4" strokeWidth={1.7} />
                  {formatCount(post.comment_count)}
                </button>
                <button
                  type="button"
                  className="flex items-center gap-1.5 text-[13px] text-(--nox-ink-soft) transition hover:text-(--nox-accent)"
                >
                  <Repeat2 className="size-4" strokeWidth={1.7} />
                  {formatCount(post.repost_count)}
                </button>
              </div>
            </div>

            {/* Comments header */}
            {comments.length > 0 && (
              <div className="px-4 py-3">
                <p className="font-mono text-[10px] font-semibold uppercase tracking-[0.16em] text-(--nox-ink-soft)">
                  {comments.length} comment{comments.length !== 1 ? "s" : ""}
                </p>
              </div>
            )}

            {/* Comments */}
            <div className="divide-y divide-(--nox-divider) px-4">
              {comments.map((comment) => (
                <CommentItem key={comment.id} comment={comment} isReply={!!comment.parent_id} />
              ))}
              {comments.length === 0 && !loading && (
                <p className="py-8 text-center text-[13px] text-(--nox-ink-soft)">
                  no comments yet. be first.
                </p>
              )}
            </div>
          </>
        )}
      </div>

      {/* Comment input */}
      <div className="border-t border-(--nox-divider) px-4 py-3 pb-[env(safe-area-inset-bottom,12px)]">
        <div className="flex items-center gap-3">
          <div
            className="flex size-8 shrink-0 items-center justify-center rounded-full"
            style={{ background: "var(--nox-surface-alt)" }}
          >
            <Ghost className="size-3.5" strokeWidth={1.6} style={{ color: "var(--nox-ink-mid)" }} />
          </div>
          <div className="flex flex-1 items-center gap-2 rounded-[10px] border border-(--nox-border) bg-(--nox-surface) px-3 py-2 transition focus-within:border-(--nox-accent-line)">
            <input
              type="text"
              value={commentBody}
              onChange={(e) => setCommentBody(e.target.value)}
              placeholder="add a comment…"
              className="flex-1 bg-transparent text-[13px] text-(--nox-ink) outline-none placeholder:text-(--nox-ink-soft)"
            />
            <button
              type="button"
              disabled={!commentBody.trim()}
              className="shrink-0 transition disabled:opacity-40"
              style={{ color: "var(--nox-accent)" }}
            >
              <Send className="size-4" strokeWidth={1.7} />
            </button>
          </div>
        </div>
      </div>
    </FeedShell>
  );
}
