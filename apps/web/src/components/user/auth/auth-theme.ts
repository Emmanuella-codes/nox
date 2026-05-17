import type { CSSProperties } from "react";
import type { AuthTheme } from "@/src/types/components/auth";

export const darkAuthTheme: AuthTheme = {
  bg: "#0a0a0c",
  bgSoft: "#101013",
  surface: "#15151a",
  surfaceAlt: "#1c1c22",
  surfaceSunk: "#0e0e11",
  border: "rgba(255,255,255,0.07)",
  borderStrong: "rgba(255,255,255,0.14)",
  divider: "rgba(255,255,255,0.05)",
  ink: "#ededf2",
  inkMid: "#9a96a8",
  inkSoft: "#6b6878",
  accent: "#a78bfa",
  accentInk: "#c8b6ff",
  accentSoft: "rgba(167,139,250,0.14)",
  accentLine: "rgba(167,139,250,0.4)",
  danger: "#e07070",
  dangerSoft: "rgba(224,112,112,0.14)",
  success: "#7ac687",
  successSoft: "rgba(122,198,135,0.14)",
  gold: "#d4a347",
  shadow: "0 1px 2px rgba(0,0,0,.4), 0 8px 32px rgba(0,0,0,.3)",
};

export function authThemeVars(theme: AuthTheme) {
  return {
    "--nox-bg": theme.bg,
    "--nox-bg-soft": theme.bgSoft,
    "--nox-surface": theme.surface,
    "--nox-surface-alt": theme.surfaceAlt,
    "--nox-surface-sunk": theme.surfaceSunk,
    "--nox-border": theme.border,
    "--nox-border-strong": theme.borderStrong,
    "--nox-divider": theme.divider,
    "--nox-ink": theme.ink,
    "--nox-ink-mid": theme.inkMid,
    "--nox-ink-soft": theme.inkSoft,
    "--nox-accent": theme.accent,
    "--nox-accent-ink": theme.accentInk,
    "--nox-accent-soft": theme.accentSoft,
    "--nox-accent-line": theme.accentLine,
    "--nox-danger": theme.danger,
    "--nox-danger-soft": theme.dangerSoft,
    "--nox-success": theme.success,
    "--nox-success-soft": theme.successSoft,
    "--nox-gold": theme.gold,
    "--nox-shadow": theme.shadow,
  } as CSSProperties;
}
