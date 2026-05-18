"use client";

import Link from "next/link";
import { Sun, Moon } from "lucide-react";
import { BrandDot } from "@/src/components/landing/landing-primitives";
import type { Theme } from "@/src/components/landing/landing-theme";

interface LandingNavProps {
  theme: Theme;
  onToggleTheme: () => void;
}

const NAV_LINKS = [
  { label: "what", href: "#what" },
  { label: "scenes", href: "#scenes" },
  { label: "manifesto", href: "#manifesto" },
  { label: "start", href: "#start" },
];

export function LandingNav({ theme, onToggleTheme }: LandingNavProps) {
  return (
    <nav
      className="sticky top-0 z-50 border-b border-[var(--l-divider)] backdrop-blur-[18px]"
      style={{
        background: theme === "dark" ? "rgba(9,9,9,0.8)" : "rgba(247,248,248,0.85)",
      }}
    >
      <div className="mx-auto flex max-w-[1240px] items-center justify-between px-8 py-4">
        {/* Brand */}
        <div className="flex items-center gap-2.5 text-[22px] font-bold text-[var(--l-ink)]">
          <BrandDot />
          <span>nox</span>
        </div>

        {/* Links — hidden on mobile */}
        <div className="landing-nav-links flex gap-7">
          {NAV_LINKS.map(({ label, href }) => (
            <a
              key={label}
              href={href}
              className="text-[13.5px] font-medium text-[var(--l-ink-mid)] no-underline transition-colors hover:text-[var(--l-ink)]"
            >
              {label}
            </a>
          ))}
        </div>

        {/* Actions */}
        <div className="flex items-center gap-3.5">
          {/* Theme toggle */}
          <button
            onClick={onToggleTheme}
            aria-label="Toggle theme"
            className="relative flex h-7 w-14 cursor-pointer items-center rounded-lg border border-[var(--l-border)] bg-[var(--l-surface-alt)] p-0.5"
          >
            <span
              className="flex size-[22px] items-center justify-center rounded-full bg-[var(--l-ink)] text-[var(--l-bg)] transition-transform duration-300"
              style={{ transform: theme === "dark" ? "translateX(0)" : "translateX(28px)" }}
            >
              {theme === "dark" ? <Moon size={13} strokeWidth={1.8} /> : <Sun size={13} strokeWidth={1.8} />}
            </span>
          </button>

          <Link
            href="/auth"
            className="inline-flex items-center gap-2 rounded-lg border border-[var(--l-ink)] bg-[var(--l-ink)] px-[18px] py-[11px] text-[13.5px] font-semibold text-[var(--l-bg)] no-underline transition-colors hover:bg-[var(--l-accent)] hover:border-[var(--l-accent)] hover:text-white"
          >
            open nox
          </Link>
        </div>
      </div>
    </nav>
  );
}
