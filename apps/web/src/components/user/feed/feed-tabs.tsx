"use client";

export type FeedTab = "for-you" | "following" | "sets" | "events";

const TAB_LABELS: { key: FeedTab; label: string }[] = [
  { key: "for-you", label: "for you" },
  { key: "following", label: "following" },
  { key: "sets", label: "sets" },
  { key: "events", label: "events" },
];

interface FeedTabsProps {
  active: FeedTab;
  onChange: (tab: FeedTab) => void;
}

export function FeedTabs({ active, onChange }: FeedTabsProps) {
  return (
    <div className="flex gap-0 border-b border-(--nox-divider) px-4">
      {TAB_LABELS.map(({ key, label }) => {
        const isActive = active === key;
        return (
          <button
            key={key}
            type="button"
            onClick={() => onChange(key)}
            className="relative pb-2.5 pr-5 font-mono text-[11px] font-medium tracking-[0.06em] transition"
            style={{ color: isActive ? "var(--nox-accent)" : "var(--nox-ink-soft)" }}
          >
            {label}
            {isActive && (
              <span
                className="absolute bottom-0 left-0 right-5 h-[1.5px] rounded-full"
                style={{ background: "var(--nox-accent)" }}
              />
            )}
          </button>
        );
      })}
    </div>
  );
}
