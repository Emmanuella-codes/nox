import type { CSSProperties } from "react";

export type Theme = "dark" | "light";

const DARK: Record<string, string> = {
  "--l-bg": "#090909",
  "--l-bg-soft": "#111111",
  "--l-surface": "#151515",
  "--l-surface-alt": "#202020",
  "--l-border": "rgba(255,255,255,0.07)",
  "--l-border-strong": "rgba(255,255,255,0.14)",
  "--l-divider": "rgba(255,255,255,0.05)",
  "--l-ink": "#f2f2ef",
  "--l-ink-mid": "#b2b0aa",
  "--l-ink-soft": "#7b7973",
  "--l-ink-faint": "#4a4842",
  "--l-accent": "#14b8a6",
  "--l-accent-ink": "#5eead4",
  "--l-accent-soft": "rgba(20,184,166,0.13)",
  "--l-accent-line": "rgba(20,184,166,0.38)",
  "--l-gold": "#f2b84b",
  "--l-shadow": "0 1px 2px rgba(0,0,0,.4), 0 8px 32px rgba(0,0,0,.3)",
};

const LIGHT: Record<string, string> = {
  "--l-bg": "#f7f8f8",
  "--l-bg-soft": "#ffffff",
  "--l-surface": "#ffffff",
  "--l-surface-alt": "#eef1f1",
  "--l-border": "rgba(20,18,22,0.09)",
  "--l-border-strong": "rgba(20,18,22,0.16)",
  "--l-divider": "rgba(20,18,22,0.06)",
  "--l-ink": "#141414",
  "--l-ink-mid": "#555a57",
  "--l-ink-soft": "#858a87",
  "--l-ink-faint": "#b8bebb",
  "--l-accent": "#0f766e",
  "--l-accent-ink": "#0f766e",
  "--l-accent-soft": "#e2f5f1",
  "--l-accent-line": "rgba(15,118,110,0.3)",
  "--l-gold": "#9a6a10",
  "--l-shadow": "0 1px 2px rgba(20,18,22,.04), 0 8px 32px rgba(20,18,22,.06)",
};

export function themeVars(t: Theme): CSSProperties {
  return Object.fromEntries(Object.entries(t === "dark" ? DARK : LIGHT)) as CSSProperties;
}
