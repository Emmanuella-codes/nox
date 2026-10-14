import type { CSSProperties } from "react";

export type Theme = "dark";

const DARK: Record<string, string> = {
  "--l-bg": "#0b0b0a",
  "--l-surface": "#121211",
  "--l-surface-alt": "#181816",
  "--l-border": "#1d1d1b",
  "--l-border-strong": "#232321",
  "--l-divider": "#1d1d1b",
  "--l-ink": "#f2f1ec",
  "--l-ink-mid": "#9a9a94",
  "--l-ink-soft": "#8a8a84",
  "--l-ink-faint": "#6f6f69",
  "--l-accent": "#16b8a4",
  "--l-accent-ink": "#16b8a4",
  "--l-accent-soft": "rgba(22,184,164,0.1)",
  "--l-accent-line": "rgba(22,184,164,0.25)",
};

export function themeVars(): CSSProperties {
  return Object.fromEntries(Object.entries(DARK)) as CSSProperties;
}
