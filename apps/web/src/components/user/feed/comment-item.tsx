"use client";

import { Heart, Ghost } from "lucide-react";
import type { Comment } from "@/src/types/api/user/comment";

interface CommentItemProps {
  comment: Comment;
  isReply?: boolean;
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

export function CommentItem({ comment, isReply = false }: CommentItemProps) {
  const displayName = comment.author.mode === "anonymous"
    ? comment.author.anonymous?.handle ?? "anonymous"
    : comment.author.persona?.display_name ?? comment.author.persona?.handle ?? "public commenter";

  return (
    <div className={`flex gap-3 py-3 ${isReply ? "pl-10" : ""}`}>
      {/* Avatar */}
      <div className="shrink-0">
        <div
          className="flex size-8 items-center justify-center rounded-full"
          style={{ background: "var(--nox-surface-alt)" }}
        >
          <Ghost className="size-3.5" strokeWidth={1.6} style={{ color: "var(--nox-ink-mid)" }} />
        </div>
      </div>

      {/* Content */}
      <div className="min-w-0 flex-1">
        <div className="flex items-center gap-2">
          <span className="text-[12px] font-semibold text-(--nox-ink)">
            {displayName}
          </span>
          <span className="text-[11px] text-(--nox-ink-soft)">
            {formatTime(comment.created_at)}
          </span>
        </div>
        <p className="mt-0.5 text-[13px] leading-[1.5] text-(--nox-ink)">{comment.body}</p>
        <div className="mt-2 flex items-center gap-4">
          <button
            type="button"
            className="flex items-center gap-1 text-[11px] text-(--nox-ink-soft) transition hover:text-(--nox-accent)"
          >
            <Heart className="size-3" strokeWidth={1.7} />
            {comment.like_count > 0 && comment.like_count}
          </button>
          <button type="button" className="text-[11px] text-(--nox-ink-soft) transition hover:text-(--nox-accent)">
            reply
          </button>
        </div>
      </div>
    </div>
  );
}
