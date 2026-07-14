"use client";

import type { AuthMode } from "@/src/types/components/user/auth";

interface AuthModeTabsProps {
  mode: AuthMode;
  onModeChange: (mode: AuthMode) => void;
}

export function AuthModeTabs({ mode, onModeChange }: AuthModeTabsProps) {
  return (
    <div className="grid grid-cols-2 rounded-[10px] border border-[var(--nox-border)] bg-[var(--nox-surface-sunk)] p-1">
      {(["login", "signup"] as AuthMode[]).map((item) => {
        const active = mode === item;
        return (
          <button
            key={item}
            type="button"
            onClick={() => onModeChange(item)}
            className={`min-h-9 rounded-[8px] text-[13px] font-semibold capitalize transition ${
              active
                ? "bg-[var(--nox-surface-alt)] text-[var(--nox-ink)] shadow-sm"
                : "text-[var(--nox-ink-soft)] hover:text-[var(--nox-ink-mid)]"
            }`}
          >
            {item}
          </button>
        );
      })}
    </div>
  );
}
