"use client";

import { Clock, Users } from "lucide-react";
import { Avatar } from "@/src/components/user/shared/avatar";
import type { Story } from "@/src/types/api/story";

interface StoryCardProps {
  story: Story;
  onPress?: () => void;
  compact?: boolean;
}

function formatDuration(seconds: number) {
  const m = Math.floor(seconds / 60);
  const s = seconds % 60;
  return `${m}:${String(s).padStart(2, "0")}`;
}

export function StoryCard({ story, onPress, compact }: StoryCardProps) {
  const modeLabel = story.contribution_mode === "public" ? "open" : "followers";

  if (compact) {
    return (
      <button type="button" onClick={onPress}
        className="flex w-[140px] shrink-0 flex-col gap-2 rounded-[10px] border border-(--nox-border) bg-(--nox-surface) p-3 text-left transition active:bg-(--nox-surface-alt)">
        <Avatar id={story.owner.id} name={story.owner.display_name} size={32} />
        <p className="line-clamp-2 text-[12px] font-semibold leading-tight text-(--nox-ink)">{story.title}</p>
        <div className="flex items-center gap-2 text-(--nox-ink-soft)">
          <span className="flex items-center gap-0.5 font-mono text-[9px]">
            <Clock className="size-2.5" strokeWidth={1.7} />
            {formatDuration(story.total_duration_seconds)}
          </span>
          <span className="font-mono text-[9px]">{story.items.length} clips</span>
        </div>
      </button>
    );
  }

  return (
    <button type="button" onClick={onPress}
      className="flex w-full items-center gap-3 border-b border-(--nox-divider) px-4 py-3 text-left transition active:bg-(--nox-surface)">
      <Avatar id={story.owner.id} name={story.owner.display_name} size={38} />
      <div className="min-w-0 flex-1">
        <p className="truncate text-[14px] font-semibold text-(--nox-ink)">{story.title}</p>
        <p className="text-[11px] text-(--nox-ink-soft)">{story.owner.display_name}</p>
      </div>
      <div className="flex shrink-0 flex-col items-end gap-1">
        <span className="flex items-center gap-1 text-[11px] text-(--nox-ink-soft)">
          <Clock className="size-3" strokeWidth={1.7} />
          {formatDuration(story.total_duration_seconds)}
        </span>
        <span className="flex items-center gap-1 rounded-full bg-(--nox-surface-alt) px-2 py-0.5 font-mono text-[9px] text-(--nox-ink-mid)">
          <Users className="size-2.5" strokeWidth={1.7} />
          {modeLabel}
        </span>
      </div>
    </button>
  );
}
