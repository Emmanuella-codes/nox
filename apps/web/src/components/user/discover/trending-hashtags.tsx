"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { Hash } from "lucide-react";
import { getTrendingHashtags } from "@/src/utils/api/user/hashtag";
import type { Hashtag } from "@/src/types/api/user/hashtag";

export function TrendingHashtags() {
  const router = useRouter();
  const [tags, setTags] = useState<Hashtag[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    getTrendingHashtags(16)
      .then((res) => setTags(res.data ?? []))
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  if (loading) {
    return (
      <div className="px-4 pt-4">
        <div className="mb-3 h-3 w-28 animate-pulse rounded bg-(--nox-surface-alt)" />
        <div className="flex flex-wrap gap-2">
          {[68, 90, 76, 102, 72, 88, 80, 64].map((w, i) => (
            <div
              key={i}
              className="h-8 animate-pulse rounded-full bg-(--nox-surface-alt)"
              style={{ width: w }}
            />
          ))}
        </div>
      </div>
    );
  }

  if (tags.length === 0) return null;

  return (
    <div className="px-4 pt-4">
      <div className="mb-3 flex items-center gap-2">
        <Hash className="size-3 text-(--nox-accent-ink)" strokeWidth={2} />
        <p className="font-mono text-[10px] font-semibold uppercase tracking-[0.16em] text-(--nox-ink-soft)">
          trending scenes
        </p>
      </div>
      <div className="flex flex-wrap gap-2">
        {tags.map((h) => (
          <button
            key={h.id}
            type="button"
            onClick={() => router.push(`/hashtags/${encodeURIComponent(h.tag)}`)}
            className="flex items-center gap-1 rounded-full border border-(--nox-border) px-3 py-1.5 font-mono text-[11px] text-(--nox-ink-mid) transition hover:border-(--nox-accent-line) hover:text-(--nox-accent-ink)"
          >
            <span style={{ color: "var(--nox-accent-ink)" }}>#</span>
            {h.tag}
            <span className="ml-0.5 text-[10px] text-(--nox-ink-soft)">{h.post_count}</span>
          </button>
        ))}
      </div>
    </div>
  );
}
