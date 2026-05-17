interface OnboardingProgressProps {
  current: number;
  total: number;
}

export function OnboardingProgress({ current, total }: OnboardingProgressProps) {
  return (
    <div className="flex items-center gap-3 px-5 pt-5">
      <div className="flex flex-1 gap-1.5">
        {Array.from({ length: total }, (_, i) => (
          <div
            key={i}
            className="h-[3px] flex-1 rounded-full transition-all duration-300"
            style={{
              background:
                i < current
                  ? "var(--nox-accent)"
                  : i === current
                    ? "var(--nox-accent)"
                    : "var(--nox-border-strong)",
              opacity: i === current ? 1 : i < current ? 0.6 : 1,
            }}
          />
        ))}
      </div>
      <span className="font-mono text-[10px] font-medium tracking-[0.1em] text-(--nox-ink-soft)">
        {current + 1}/{total}
      </span>
    </div>
  );
}
