"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { Ghost, User, Image, Music, MapPin, X, Tag } from "lucide-react";
import { FeedShell } from "@/src/components/user/feed/feed-shell";

type PostingMode = "ghost" | "visible";

const MAX_CHARS = 280;

export function ComposeScreen() {
  const router = useRouter();
  const [mode, setMode] = useState<PostingMode>("ghost");
  const [body, setBody] = useState("");
  const [posting, setPosting] = useState(false);

  const remaining = MAX_CHARS - body.length;
  const canPost = body.trim().length > 0 && remaining >= 0 && !posting;

  const ringColor =
    remaining < 0 ? "var(--nox-danger)" : remaining < 30 ? "var(--nox-gold)" : "var(--nox-accent)";

  function handlePost() {
    if (!canPost) return;
    setPosting(true);
    // Future: call createPost API
    setTimeout(() => {
      router.back();
    }, 500);
  }

  return (
    <FeedShell>
      {/* Header */}
      <header className="flex items-center justify-between border-b border-(--nox-divider) px-4 py-3 pt-[env(safe-area-inset-top,12px)]">
        <button
          type="button"
          onClick={() => router.back()}
          className="flex size-8 items-center justify-center rounded-full transition hover:bg-(--nox-surface-alt)"
        >
          <X className="size-5 text-(--nox-ink)" strokeWidth={1.8} />
        </button>
        <span className="text-[15px] font-semibold text-(--nox-ink)">new post</span>
        <button
          type="button"
          disabled={!canPost}
          onClick={handlePost}
          className="rounded-full px-4 py-1.5 text-[13px] font-semibold text-white transition hover:brightness-110 disabled:cursor-not-allowed disabled:opacity-50"
          style={{ background: "var(--nox-accent)" }}
        >
          {posting ? "posting…" : "post"}
        </button>
      </header>

      <div className="flex flex-1 flex-col overflow-y-auto px-4 py-4">
        {/* Persona indicator */}
        <div className="mb-4 flex items-center gap-2">
          <button
            type="button"
            onClick={() => setMode("ghost")}
            className="flex items-center gap-1.5 rounded-full border px-3 py-1.5 font-mono text-[10px] font-semibold uppercase tracking-[0.12em] transition"
            style={{
              borderColor: mode === "ghost" ? "var(--nox-accent)" : "var(--nox-border-strong)",
              background: mode === "ghost" ? "var(--nox-accent-soft)" : "transparent",
              color: mode === "ghost" ? "var(--nox-accent)" : "var(--nox-ink-mid)",
            }}
          >
            <Ghost className="size-3" strokeWidth={1.8} />
            ghost
          </button>
          <button
            type="button"
            onClick={() => setMode("visible")}
            className="flex items-center gap-1.5 rounded-full border px-3 py-1.5 font-mono text-[10px] font-semibold uppercase tracking-[0.12em] transition"
            style={{
              borderColor: mode === "visible" ? "var(--nox-accent)" : "var(--nox-border-strong)",
              background: mode === "visible" ? "var(--nox-accent-soft)" : "transparent",
              color: mode === "visible" ? "var(--nox-accent)" : "var(--nox-ink-mid)",
            }}
          >
            <User className="size-3" strokeWidth={1.8} />
            visible
          </button>
        </div>

        {/* Text area */}
        <textarea
          autoFocus
          value={body}
          onChange={(e) => setBody(e.target.value)}
          placeholder="what's the vibe?"
          rows={7}
          className="w-full resize-none bg-transparent text-[16px] leading-[1.6] text-(--nox-ink) outline-none placeholder:text-(--nox-ink-soft)"
        />

        {/* Attachments row */}
        <div className="mt-auto border-t border-(--nox-divider) pt-3">
          <div className="flex items-center gap-1">
            <button
              type="button"
              className="flex size-9 items-center justify-center rounded-full transition hover:bg-(--nox-surface-alt) text-(--nox-ink-soft)"
            >
              <Image className="size-4" strokeWidth={1.7} />
            </button>
            <button
              type="button"
              className="flex size-9 items-center justify-center rounded-full transition hover:bg-(--nox-surface-alt) text-(--nox-ink-soft)"
            >
              <Music className="size-4" strokeWidth={1.7} />
            </button>
            <button
              type="button"
              className="flex size-9 items-center justify-center rounded-full transition hover:bg-(--nox-surface-alt) text-(--nox-ink-soft)"
            >
              <Tag className="size-4" strokeWidth={1.7} />
            </button>
            <button
              type="button"
              className="flex size-9 items-center justify-center rounded-full transition hover:bg-(--nox-surface-alt) text-(--nox-ink-soft)"
            >
              <MapPin className="size-4" strokeWidth={1.7} />
            </button>

            {/* Char counter */}
            <div className="ml-auto flex items-center gap-2">
              <span
                className="font-mono text-[11px] font-medium"
                style={{ color: ringColor }}
              >
                {remaining}
              </span>
              <svg className="size-5 -rotate-90" viewBox="0 0 20 20">
                <circle cx="10" cy="10" r="8" fill="none" stroke="var(--nox-border)" strokeWidth="2" />
                <circle
                  cx="10"
                  cy="10"
                  r="8"
                  fill="none"
                  stroke={ringColor}
                  strokeWidth="2"
                  strokeDasharray={`${Math.min(2 * Math.PI * 8, (Math.max(0, MAX_CHARS - body.length) / MAX_CHARS) * 2 * Math.PI * 8)} ${2 * Math.PI * 8}`}
                  strokeLinecap="round"
                />
              </svg>
            </div>
          </div>
        </div>
      </div>
    </FeedShell>
  );
}
