"use client";

import { useEffect, useRef } from "react";
import { Avatar } from "@/src/components/user/shared/avatar";
import type { StoryItem } from "@/src/types/api/user/story";

interface StoryItemPlayerProps {
  item: StoryItem;
  active: boolean;
  onEnded: () => void;
}

export function StoryItemPlayer({ item, active, onEnded }: StoryItemPlayerProps) {
  const videoRef = useRef<HTMLVideoElement>(null);
  const src = item.media_asset?.playback_url ?? "";
  const poster = item.media_asset?.thumbnail_url ?? undefined;
  const isAnon = item.posting_mode === "anonymous";
  const contributor = item.contributor;

  useEffect(() => {
    const vid = videoRef.current;
    if (!vid) return;
    if (active) {
      vid.currentTime = 0;
      void vid.play().catch(() => undefined);
    } else {
      vid.pause();
    }
  }, [active]);

  return (
    <div className="relative size-full bg-black">
      {src ? (
        <video
          ref={videoRef}
          src={src}
          poster={poster}
          playsInline
          preload="auto"
          onEnded={onEnded}
          className="size-full object-cover"
        />
      ) : (
        <div className="flex size-full items-center justify-center">
          <p className="text-[13px] text-white/60">No video</p>
        </div>
      )}

      {/* Contributor overlay */}
      <div className="absolute bottom-0 inset-x-0 bg-gradient-to-t from-black/70 to-transparent px-4 pb-6 pt-16">
        {isAnon ? (
          <div className="flex items-center gap-2">
            <div className="flex size-8 items-center justify-center rounded-full bg-white/10">
              <span className="font-mono text-[10px] text-white/70">?</span>
            </div>
            <p className="font-mono text-[11px] text-white/70">{item.anonymous_label ?? "anonymous"}</p>
          </div>
        ) : contributor ? (
          <div className="flex items-center gap-2">
            <Avatar id={contributor.id} name={contributor.display_name} size={32} />
            <div>
              <p className="text-[13px] font-semibold text-white">{contributor.display_name}</p>
              <p className="font-mono text-[10px] text-white/60">@{contributor.handle}</p>
            </div>
          </div>
        ) : null}
      </div>
    </div>
  );
}
