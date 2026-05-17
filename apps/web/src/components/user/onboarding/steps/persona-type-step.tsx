"use client";

import { ArrowRight, Check, Ghost, User } from "lucide-react";

export type PersonaMode = "anonymous" | "visible";

interface PersonaTypeStepProps {
  selected: PersonaMode | null;
  onSelect: (mode: PersonaMode) => void;
  onContinue: () => void;
}

export function PersonaTypeStep({ selected, onSelect, onContinue }: PersonaTypeStepProps) {
  return (
    <div className="flex flex-1 flex-col px-5">
      <p className="font-mono text-[10px] font-medium uppercase tracking-[0.18em] text-(--nox-ink-soft)">
        Step 1 of 3
      </p>
      <h2 className="mt-2 text-[24px] font-bold leading-tight tracking-[-0.03em] text-(--nox-ink)">
        How do you want to show up?
      </h2>
      <p className="mt-2 text-[13px] leading-6 text-(--nox-ink-soft)">
        You can switch between modes anytime before posting.
      </p>

      <div className="mt-5 flex flex-col gap-3">
        <button
          type="button"
          onClick={() => onSelect("anonymous")}
          className="w-full rounded-[14px] border p-4 text-left transition"
          style={{
            borderColor: selected === "anonymous" ? "var(--nox-accent)" : "var(--nox-border)",
            borderWidth: selected === "anonymous" ? "1.5px" : "1px",
            background: selected === "anonymous" ? "var(--nox-accent-soft)" : "var(--nox-surface)",
            boxShadow:
              selected === "anonymous" ? "0 0 24px rgba(167,139,250,0.12)" : "none",
          }}
        >
          <div className="flex items-start gap-3">
            <div
              className="flex size-10 shrink-0 items-center justify-center rounded-[10px]"
              style={{
                background:
                  selected === "anonymous" ? "rgba(167,139,250,0.18)" : "var(--nox-surface-alt)",
                color: selected === "anonymous" ? "var(--nox-accent)" : "var(--nox-ink-mid)",
              }}
            >
              <Ghost className="size-5" strokeWidth={1.6} />
            </div>
            <div className="flex-1 min-w-0">
              <div className="flex items-center justify-between gap-2">
                <span className="text-[15px] font-bold text-(--nox-ink)">Anonymous mode</span>
                {selected === "anonymous" ? (
                  <span
                    className="flex size-5 shrink-0 items-center justify-center rounded-full"
                    style={{ background: "var(--nox-accent)" }}
                  >
                    <Check className="size-3 text-white" strokeWidth={2.5} />
                  </span>
                ) : (
                  <span className="size-5 shrink-0 rounded-full border border-(--nox-border-strong)" />
                )}
              </div>
              <p className="mt-1.5 text-[12px] leading-normal text-(--nox-ink-mid)">
                Post without exposing your persona. Each anonymous post is unlinked publicly.
              </p>
              <p
                className="mt-2 font-mono text-[9.5px] font-semibold uppercase tracking-[0.16em]"
                style={{ color: "var(--nox-accent)" }}
              >
                Recommended
              </p>
            </div>
          </div>
        </button>

        <button
          type="button"
          onClick={() => onSelect("visible")}
          className="w-full rounded-[14px] border p-4 text-left transition"
          style={{
            borderColor: selected === "visible" ? "var(--nox-accent)" : "var(--nox-border)",
            borderWidth: selected === "visible" ? "1.5px" : "1px",
            background: selected === "visible" ? "var(--nox-accent-soft)" : "var(--nox-surface)",
            boxShadow:
              selected === "visible" ? "0 0 24px rgba(167,139,250,0.12)" : "none",
          }}
        >
          <div className="flex items-start gap-3">
            <div
              className="flex size-10 shrink-0 items-center justify-center rounded-[10px]"
              style={{
                background:
                  selected === "visible" ? "rgba(167,139,250,0.18)" : "var(--nox-surface-alt)",
                color: selected === "visible" ? "var(--nox-accent)" : "var(--nox-ink-mid)",
              }}
            >
              <User className="size-5" strokeWidth={1.6} />
            </div>
            <div className="flex-1 min-w-0">
              <div className="flex items-center justify-between gap-2">
                <span className="text-[15px] font-bold text-(--nox-ink)">Visible persona</span>
                {selected === "visible" ? (
                  <span
                    className="flex size-5 shrink-0 items-center justify-center rounded-full"
                    style={{ background: "var(--nox-accent)" }}
                  >
                    <Check className="size-3 text-white" strokeWidth={2.5} />
                  </span>
                ) : (
                  <span className="size-5 shrink-0 rounded-full border border-(--nox-border-strong)" />
                )}
              </div>
              <p className="mt-1.5 text-[12px] leading-normal text-(--nox-ink-mid)">
                Post as your DJ alias or public name. Build a following, list events, and archive sets.
              </p>
            </div>
          </div>
        </button>
      </div>

      <button
        type="button"
        disabled={!selected}
        onClick={onContinue}
        className="mt-5 flex min-h-12 w-full items-center justify-center gap-2 rounded-[10px] border border-(--nox-accent) bg-(--nox-accent) px-4 py-3 text-[15px] font-semibold text-white shadow-[0_0_24px_rgba(167,139,250,0.24)] transition hover:brightness-110 disabled:cursor-not-allowed disabled:opacity-40"
      >
        Continue
        <ArrowRight className="size-4" strokeWidth={1.8} />
      </button>
    </div>
  );
}
