"use client";

import { PenLine } from "lucide-react";

interface ComposeBarProps {
  onClick: () => void;
}

export function ComposeBar({ onClick }: ComposeBarProps) {
  return (
    <div className="border-t border-(--nox-divider) px-4 py-3">
      <button
        type="button"
        onClick={onClick}
        className="flex w-full items-center gap-3 rounded-[12px] border border-(--nox-border) bg-(--nox-surface) px-4 py-3 text-left transition hover:border-(--nox-border-strong)"
      >
        <PenLine className="size-4 shrink-0 text-(--nox-ink-soft)" strokeWidth={1.7} />
        <span className="text-[14px] text-(--nox-ink-soft)">what&apos;s the vibe?</span>
        <span
          className="ml-auto rounded-full px-3 py-1 text-[11px] font-semibold"
          style={{ background: "var(--nox-accent)", color: "#fff" }}
        >
          post
        </span>
      </button>
    </div>
  );
}
