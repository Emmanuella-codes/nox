"use client";

import { useEffect, useState } from "react";
import { ChevronLeft, Eye, Heart, MessageCircle, Timer } from "lucide-react";
import { useRouter } from "next/navigation";
import { FeedShell } from "@/src/components/user/feed/feed-shell";
import { TabBar } from "@/src/components/user/feed/tab-bar";
import { Avatar } from "@/src/components/user/shared/avatar";
import { getSet } from "@/src/utils/api/user/set";
import type { Set } from "@/src/types/api/user/set";

interface SetDetailScreenProps {
  setID: string;
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

export function SetDetailScreen({ setID }: SetDetailScreenProps) {
  const router = useRouter();
  const [set, setSet] = useState<Set | null>(null);
  const [loading, setLoading] = useState(true);
  const [message, setMessage] = useState("");

  useEffect(() => {
    getSet(setID)
      .then((res) => setSet(res.data ?? null))
      .catch(() => setMessage("Could not load this set."))
      .finally(() => setLoading(false));
  }, [setID]);

  const media = set?.media_asset;
  const persona = set?.persona;
  const processingStatus = media?.processing_status;

  return (
    <FeedShell>
      <header className="flex items-center gap-3 border-b border-(--nox-divider) px-4 pt-[env(safe-area-inset-top,12px)] pb-3">
        <button
          type="button"
          onClick={() => router.back()}
          className="flex size-8 shrink-0 items-center justify-center rounded-full transition hover:bg-(--nox-surface-alt)"
        >
          <ChevronLeft className="size-5 text-(--nox-ink)" strokeWidth={1.8} />
        </button>
        <p className="min-w-0 truncate text-[16px] font-bold text-(--nox-ink)">
          {set ? set.title : "set"}
        </p>
      </header>

      <div className="flex-1 overflow-y-auto">
        {loading ? (
          <p className="px-4 py-8 text-[13px] text-(--nox-ink-soft)">Loading set...</p>
        ) : message || !set ? (
          <p className="px-4 py-8 text-[13px] text-(--nox-danger)">{message || "Set not found."}</p>
        ) : (
          <div>
            <div className="bg-black">
              {processingStatus === "pending" ? (
                <div className="flex aspect-video flex-col items-center justify-center gap-2">
                  <p className="text-[13px] font-semibold text-white/80">Processing video...</p>
                  <p className="text-[11px] text-white/50">Check back shortly.</p>
                </div>
              ) : processingStatus === "failed" ? (
                <div className="flex aspect-video items-center justify-center">
                  <p className="text-[13px] font-semibold text-red-400">Video processing failed.</p>
                </div>
              ) : media?.playback_url ? (
                <video
                  className="aspect-video w-full bg-black"
                  controls
                  poster={media.thumbnail_url || undefined}
                  preload="metadata"
                  src={media.playback_url}
                />
              ) : (
                <div className="flex aspect-video items-center justify-center text-[13px] text-white/70">
                  Video unavailable
                </div>
              )}
            </div>

            <div className="px-4 py-4">
              <div className="flex items-start justify-between gap-3">
                <h1 className="min-w-0 text-[22px] font-bold leading-tight tracking-[-0.03em] text-(--nox-ink)">
                  {set.title}
                </h1>
                <span className="flex shrink-0 items-center gap-1 rounded-[6px] bg-(--nox-surface-alt) px-2 py-1 font-mono text-[10px] font-semibold text-(--nox-ink-mid)">
                  <Timer className="size-3" strokeWidth={1.8} />
                  {formatDuration(set.duration_seconds)}
                </span>
              </div>

              <div className="mt-3 flex items-center gap-5">
                <span className="flex items-center gap-1.5 text-[12px] text-(--nox-ink-soft)">
                  <Eye className="size-3.5" strokeWidth={1.7} />
                  {formatCount(set.play_count)}
                </span>
                <span className="flex items-center gap-1.5 text-[12px] text-(--nox-ink-soft)">
                  <Heart className="size-3.5" strokeWidth={1.7} />
                  {formatCount(set.like_count)}
                </span>
                <span className="flex items-center gap-1.5 text-[12px] text-(--nox-ink-soft)">
                  <MessageCircle className="size-3.5" strokeWidth={1.7} />
                  {formatCount(set.comment_count)}
                </span>
              </div>

              {persona ? (
                <button
                  type="button"
                  onClick={() => router.push(`/personas/${persona.id}`)}
                  className="mt-4 flex items-center gap-2 transition hover:opacity-80"
                >
                  <Avatar id={persona.id} name={persona.display_name} size={34} />
                  <div className="text-left">
                    <p className="text-[13px] font-semibold text-(--nox-ink)">{persona.display_name}</p>
                    <p className="font-mono text-[10px] text-(--nox-ink-soft)">@{persona.handle}</p>
                  </div>
                </button>
              ) : null}

              {set.description ? (
                <p className="mt-4 text-[13px] leading-6 text-(--nox-ink-mid)">{set.description}</p>
              ) : null}

              {set.genre_tags.length > 0 ? (
                <div className="mt-4 flex flex-wrap gap-1.5">
                  {set.genre_tags.map((tag) => (
                    <span key={tag} className="rounded-full bg-(--nox-accent-soft) px-2.5 py-1 font-mono text-[10px] text-(--nox-accent-ink)">
                      {tag}
                    </span>
                  ))}
                </div>
              ) : null}
            </div>
          </div>
        )}
      </div>

      <TabBar />
    </FeedShell>
  );
}
