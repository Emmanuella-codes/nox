"use client";

import { Eye, Heart, Play, Timer } from "lucide-react";
import { Avatar } from "@/src/components/user/shared/avatar";
import type { Set } from "@/src/types/api/set";

interface SetCardProps {
  set: Set;
  onPress?: () => void;
}

function formatDuration(seconds: number) {
  const minutes = Math.floor(seconds / 60);
  const rest = seconds % 60;
  return `${minutes}:${String(rest).padStart(2, "0")}`;
}

function formatCount(n: number) {
  if (n >= 1000) return `${(n / 1000).toFixed(1)}K`;
  return String(n);
}

export function SetCard({ set, onPress }: SetCardProps) {
  const persona = set.persona;
  const thumbnail = set.media_asset?.thumbnail_url;
  const processingStatus = set.media_asset?.processing_status;
  const isPending = processingStatus === "pending";
  const isFailed = processingStatus === "failed";

  return (
    <button
      type="button"
      onClick={onPress}
      className="grid w-full grid-cols-[144px_1fr] gap-3 border-b border-(--nox-divider) px-4 py-4 text-left transition active:bg-(--nox-surface)"
    >
      <div className="relative aspect-video overflow-hidden rounded-[8px] bg-(--nox-surface-alt)">
        {thumbnail && !isPending && !isFailed ? (
          // eslint-disable-next-line @next/next/no-img-element
          <img src={thumbnail} alt="" className="size-full object-cover" />
        ) : (
          <div className="flex size-full items-center justify-center">
            <Play className="size-6 text-(--nox-accent-ink)" strokeWidth={1.7} />
          </div>
        )}
        {isPending ? (
          <span className="absolute inset-x-0 bottom-0 bg-black/60 py-1 text-center font-mono text-[9px] font-semibold uppercase tracking-wide text-amber-300">
            Processing
          </span>
        ) : isFailed ? (
          <span className="absolute inset-x-0 bottom-0 bg-black/60 py-1 text-center font-mono text-[9px] font-semibold uppercase tracking-wide text-red-400">
            Failed
          </span>
        ) : (
          <span className="absolute right-1.5 bottom-1.5 flex items-center gap-1 rounded-[5px] bg-black/70 px-1.5 py-0.5 font-mono text-[9px] font-semibold text-white">
            <Timer className="size-2.5" strokeWidth={1.8} />
            {formatDuration(set.duration_seconds)}
          </span>
        )}
      </div>

      <div className="min-w-0 py-0.5">
        <p className="line-clamp-2 text-[14px] font-bold leading-snug tracking-[-0.02em] text-(--nox-ink)">
          {set.title}
        </p>

        <div className="mt-2 flex items-center gap-3">
          <span className="flex items-center gap-1 text-[11px] text-(--nox-ink-soft)">
            <Eye className="size-3" strokeWidth={1.7} />
            {formatCount(set.play_count)}
          </span>
          <span className="flex items-center gap-1 text-[11px] text-(--nox-ink-soft)">
            <Heart className="size-3" strokeWidth={1.7} />
            {formatCount(set.like_count)}
          </span>
        </div>

        {persona ? (
          <div className="mt-3 flex items-center gap-1.5">
            <Avatar id={persona.id} name={persona.display_name} size={20} />
            <p className="truncate text-[11px] font-semibold text-(--nox-ink-mid)">
              {persona.display_name}
            </p>
          </div>
        ) : null}

        {set.genre_tags.length > 0 ? (
          <div className="mt-2 flex flex-wrap gap-1">
            {set.genre_tags.slice(0, 3).map((tag) => (
              <span
                key={tag}
                className="rounded-full bg-(--nox-accent-soft) px-2 py-0.5 font-mono text-[9px] font-medium lowercase text-(--nox-accent-ink)"
              >
                {tag}
              </span>
            ))}
          </div>
        ) : null}
      </div>
    </button>
  );
}
