"use client";

import { ArrowRight } from "lucide-react";

const GENRE_OPTIONS = [
  "afrobeats",
  "amapiano",
  "afro-house",
  "afro-tech",
  "afro-soul",
  "alt-R&B",
  "R&B",
  "hip-hop",
  "highlife",
  "dancehall",
  "afropop",
  "electronic",
  "house",
  "drum & bass",
  "techno",
  "soul",
];

interface GenreStepProps {
  selected: string[];
  isGhost: boolean;
  loading: boolean;
  error: string;
  onChange: (tags: string[]) => void;
  onContinue: () => void;
  onBack: () => void;
}

export function GenreStep({
  selected,
  isGhost,
  loading,
  error,
  onChange,
  onContinue,
  onBack,
}: GenreStepProps) {
  function toggle(genre: string) {
    if (selected.includes(genre)) {
      onChange(selected.filter((g) => g !== genre));
    } else {
      onChange([...selected, genre]);
    }
  }

  return (
    <div className="flex flex-1 flex-col px-5">
      <p className="font-mono text-[10px] font-medium uppercase tracking-[0.18em] text-(--nox-ink-soft)">
        {isGhost ? "Step 2 of 2" : "Step 3 of 3"}
      </p>
      <h2 className="mt-2 text-[24px] font-bold leading-tight tracking-[-0.03em] text-(--nox-ink)">
        What's your scene?
      </h2>
      <p className="mt-2 text-[13px] leading-6 text-(--nox-ink-soft)">
        Pick the genres you're into. This tunes your feed and helps people find you.
      </p>

      <div className="mt-5 flex flex-wrap gap-2">
        {GENRE_OPTIONS.map((genre) => {
          const active = selected.includes(genre);
          return (
            <button
              key={genre}
              type="button"
              onClick={() => toggle(genre)}
              className="rounded-full border px-3 py-1.5 font-mono text-[11px] font-medium lowercase transition"
              style={{
                borderColor: active ? "var(--nox-accent-line)" : "var(--nox-border)",
                background: active ? "var(--nox-accent-soft)" : "transparent",
                color: active ? "var(--nox-accent-ink)" : "var(--nox-ink-mid)",
              }}
            >
              {genre}
            </button>
          );
        })}
      </div>

      {selected.length > 0 && (
        <p className="mt-3 font-mono text-[10px] text-(--nox-ink-soft)">
          {selected.length} selected
        </p>
      )}

      {error ? (
        <p className="mt-3 rounded-[8px] border border-(--nox-danger) bg-(--nox-danger-soft) px-3 py-2 text-[12px] font-medium text-(--nox-danger)">
          {error}
        </p>
      ) : null}

      <div className="mt-auto flex gap-3 pb-1">
        {!isGhost && (
          <button
            type="button"
            onClick={onBack}
            className="flex min-h-12 items-center justify-center rounded-[10px] border border-(--nox-border-strong) px-5 text-[14px] font-semibold text-(--nox-ink) transition hover:border-(--nox-ink)"
          >
            Back
          </button>
        )}
        <button
          type="button"
          disabled={loading}
          onClick={onContinue}
          className="flex flex-1 min-h-12 items-center justify-center gap-2 rounded-[10px] border border-(--nox-accent) bg-(--nox-accent) px-4 py-3 text-[15px] font-semibold text-white shadow-[0_0_24px_rgba(167,139,250,0.24)] transition hover:brightness-110 disabled:cursor-not-allowed disabled:opacity-50"
        >
          {loading ? "Setting up" : selected.length === 0 ? "Skip for now" : "Finish setup"}
          {!loading && <ArrowRight className="size-4" strokeWidth={1.8} />}
        </button>
      </div>
    </div>
  );
}
