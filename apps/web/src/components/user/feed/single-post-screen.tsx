"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { ChevronLeft, Send, Ghost, Heart, MessageCircle, Repeat2 } from "lucide-react";
import { Avatar } from "@/src/components/user/shared/avatar";
import { FeedShell } from "@/src/components/user/feed/feed-shell";
import { CommentItem } from "@/src/components/user/feed/comment-item";
import { getPost, getPostForViewer, likePost, unlikePost } from "@/src/utils/api/user/post";
import { createComment, getPostComments } from "@/src/utils/api/user/comment";
import { getMyPersonas } from "@/src/utils/api/user/persona";
import { getAccessToken, getActivePersonaID, setActivePersonaID } from "@/src/utils/auth/session";
import { ApiRequestError } from "@/src/utils/api/api";
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

function extractHashtags(body: string, onTagClick: (tag: string) => void) {
  const parts = body.split(/(#[A-Za-z0-9_][A-Za-z0-9_-]{0,49})/g);
  return parts.map((p, i) =>
    p.startsWith("#") ? (
      <button
        key={i}
        type="button"
        className="font-medium transition hover:underline"
        style={{ color: "var(--nox-accent-ink)" }}
        onClick={() => onTagClick(p.slice(1).toLowerCase())}
      >
        {p}
      </button>
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
  const [actionError, setActionError] = useState("");
  const [viewerPersonaID, setViewerPersonaID] = useState("");
  const [liked, setLiked] = useState(false);
  const [submittingComment, setSubmittingComment] = useState(false);
  const [togglingLike, setTogglingLike] = useState(false);

  useEffect(() => {
    async function load() {
      try {
        const token = getAccessToken();
        let personaID = "";
        if (token) {
          const personasRes = await getMyPersonas(token);
          const personas = personasRes.data ?? [];
          const activePersonaID = getActivePersonaID();
          const selectedPersona = personas.find((persona) => persona.id === activePersonaID) ?? personas[0];
          if (selectedPersona) {
            setActivePersonaID(selectedPersona.id);
          }
          personaID = selectedPersona?.id ?? "";
          setViewerPersonaID(personaID);
        }

        const [postRes, commentsRes] = await Promise.all([
          personaID ? getPostForViewer(postId, personaID, token) : getPost(postId),
          getPostComments(postId),
        ]);
        setPost(postRes.data ?? null);
        setLiked(Boolean(postRes.data?.is_liked));
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
  const canAct = Boolean(getAccessToken() && viewerPersonaID);

  async function handleToggleLike() {
    if (!post || togglingLike || !canAct) return;
    const token = getAccessToken();
    const previousPost = post;
    const nextLiked = !liked;

    setTogglingLike(true);
    setActionError("");
    setLiked(nextLiked);
    setPost({
      ...post,
      like_count: Math.max(0, post.like_count + (nextLiked ? 1 : -1)),
    });

    try {
      if (nextLiked) {
        await likePost(post.id, viewerPersonaID, token);
      } else {
        await unlikePost(post.id, viewerPersonaID, token);
      }
    } catch (err) {
      setLiked(!nextLiked);
      setPost(previousPost);
      setActionError(err instanceof ApiRequestError ? err.message : "Could not update like.");
    } finally {
      setTogglingLike(false);
    }
  }

  async function handleCreateComment() {
    if (!post || !commentBody.trim() || submittingComment || !canAct) return;
    const token = getAccessToken();
    setSubmittingComment(true);
    setActionError("");

    try {
      const res = await createComment(
        post.id,
        { persona_id: viewerPersonaID, body: commentBody.trim() },
        token,
      );
      if (res.data) {
        setComments((current) => [...current, res.data as Comment]);
        setPost({ ...post, comment_count: post.comment_count + 1 });
      }
      setCommentBody("");
    } catch (err) {
      setActionError(err instanceof ApiRequestError ? err.message : "Could not add comment.");
    } finally {
      setSubmittingComment(false);
    }
  }

  return (
    <FeedShell>
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
                  <Avatar id={persona?.id ?? "anon"} name={persona?.display_name ?? "?"} size={40} />
                )}
                <div>
                  {isAnon ? (
                    <span
                      className="flex items-center gap-1 rounded-full px-2 py-0.5 font-mono text-[9.5px] font-semibold uppercase tracking-[0.12em]"
                      style={{ background: "var(--nox-surface-alt)", color: "var(--nox-ink-mid)" }}
                    >
                      <Ghost className="size-2.5" strokeWidth={1.8} />
                      anonymous
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
                {extractHashtags(post.body, (tag) => router.push(`/hashtags/${encodeURIComponent(tag)}`))}
              </p>

              <p className="mt-3 text-[11px] text-(--nox-ink-soft)">{formatTime(post.created_at)}</p>

              <div className="mt-3 flex items-center gap-6 border-t border-(--nox-divider) pt-3">
                <button
                  type="button"
                  onClick={handleToggleLike}
                  disabled={!canAct || togglingLike}
                  className="flex items-center gap-1.5 text-[13px] text-(--nox-ink-soft) transition hover:text-(--nox-accent)"
                  style={{ color: liked ? "var(--nox-accent)" : undefined }}
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
        {actionError ? (
          <p className="mb-2 text-[12px] font-medium text-(--nox-danger)">{actionError}</p>
        ) : null}
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
              disabled={!commentBody.trim() || !canAct || submittingComment}
              onClick={handleCreateComment}
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
