"use client";

import { Heart, MessageCircle, Repeat2, Ghost } from "lucide-react";
import { Avatar } from "@/src/components/user/shared/avatar";
import type { Post } from "@/src/types/api/post";

interface PostCardProps {
  post: Post;
  onClick?: () => void;
  liked?: boolean;
  onLike?: (post: Post) => void;
}

function formatTime(isoString: string): string {
  const diff = Date.now() - new Date(isoString).getTime();
  const mins = Math.floor(diff / 60_000);
  if (mins < 1) return "just now";
  if (mins < 60) return `${mins}m`;
  const hours = Math.floor(mins / 60);
  if (hours < 24) return `${hours}h`;
  return `${Math.floor(hours / 24)}d`;
}

function formatCount(n: number): string {
  if (n >= 1000) return `${(n / 1000).toFixed(1)}k`;
  return String(n);
}

function extractHashtags(body: string): { segments: { text: string; isTag: boolean }[] } {
  const parts = body.split(/(#\w+)/g);
  return {
    segments: parts.map((p) => ({ text: p, isTag: p.startsWith("#") })),
  };
}

export function PostCard({ post, onClick, liked = false, onLike }: PostCardProps) {
  const isAnon = post.author.mode === "anonymous";
  const persona = post.author.persona;
  const anonymousLabel = post.author.anonymous_label ?? "anonymous";
  const { segments } = extractHashtags(post.body);

  return (
    <article
      className="border-b border-(--nox-divider) px-4 py-4 transition active:bg-(--nox-surface)"
      style={{ cursor: onClick ? "pointer" : "default" }}
      onClick={onClick}
    >
      <div className="flex gap-3">
        {/* Avatar */}
        <div className="shrink-0">
          {isAnon ? (
            <div
              className="flex size-9 items-center justify-center rounded-full"
              style={{ background: "var(--nox-surface-alt)" }}
            >
              <Ghost className="size-4" strokeWidth={1.6} style={{ color: "var(--nox-ink-mid)" }} />
            </div>
          ) : (
            <Avatar id={persona?.id ?? "anon"} name={persona?.display_name ?? "?"} size={36} />
          )}
        </div>

        {/* Content */}
        <div className="min-w-0 flex-1">
          {/* Header row */}
          <div className="flex items-center gap-2">
            {isAnon ? (
              <span
                className="flex items-center gap-1 rounded-full px-2 py-0.5 font-mono text-[9.5px] font-semibold uppercase tracking-[0.12em]"
                style={{
                  background: "var(--nox-surface-alt)",
                  color: "var(--nox-ink-mid)",
                }}
              >
                <Ghost className="size-2.5" strokeWidth={1.8} />
                {anonymousLabel}
              </span>
            ) : (
              <>
                <span className="text-[13px] font-semibold text-(--nox-ink) truncate">
                  {persona?.display_name}
                </span>
                <span className="text-[11px] text-(--nox-ink-soft) shrink-0">
                  @{persona?.handle}
                </span>
              </>
            )}
            <span className="ml-auto shrink-0 text-[11px] text-(--nox-ink-soft)">
              {formatTime(post.created_at)}
            </span>
          </div>

          {/* Body */}
          <p className="mt-1.5 text-[14px] leading-[1.55] text-(--nox-ink)">
            {segments.map((seg, i) =>
              seg.isTag ? (
                <span key={i} style={{ color: "var(--nox-accent-ink)" }}>
                  {seg.text}
                </span>
              ) : (
                <span key={i}>{seg.text}</span>
              ),
            )}
          </p>

          {/* Media */}
          {post.media_url && post.media_type === "image" && (
            <div className="mt-2.5 overflow-hidden rounded-[10px] border border-(--nox-border)">
              {/* eslint-disable-next-line @next/next/no-img-element */}
              <img
                src={post.media_url}
                alt=""
                className="w-full object-cover"
                style={{ maxHeight: 260 }}
              />
            </div>
          )}

          {/* Engagement row */}
          <div className="mt-3 flex items-center gap-5">
            <button
              type="button"
              className="flex items-center gap-1.5 text-[12px] text-(--nox-ink-soft) transition hover:text-(--nox-accent)"
              style={{ color: liked ? "var(--nox-accent)" : undefined }}
              onClick={(e) => {
                e.stopPropagation();
                onLike?.(post);
              }}
            >
              <Heart className="size-3.5" strokeWidth={1.7} />
              {post.like_count > 0 && formatCount(post.like_count)}
            </button>
            <button
              type="button"
              className="flex items-center gap-1.5 text-[12px] text-(--nox-ink-soft) transition hover:text-(--nox-accent)"
              onClick={(e) => e.stopPropagation()}
            >
              <MessageCircle className="size-3.5" strokeWidth={1.7} />
              {post.comment_count > 0 && formatCount(post.comment_count)}
            </button>
            <button
              type="button"
              className="flex items-center gap-1.5 text-[12px] text-(--nox-ink-soft) transition hover:text-(--nox-accent)"
              onClick={(e) => e.stopPropagation()}
            >
              <Repeat2 className="size-3.5" strokeWidth={1.7} />
              {post.repost_count > 0 && formatCount(post.repost_count)}
            </button>
          </div>
        </div>
      </div>
    </article>
  );
}
