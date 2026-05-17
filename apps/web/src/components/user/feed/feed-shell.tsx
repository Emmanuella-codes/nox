"use client";

import type { ReactNode } from "react";
import { authThemeVars, darkAuthTheme } from "@/src/components/user/auth/auth-theme";

interface FeedShellProps {
  children: ReactNode;
}

export function FeedShell({ children }: FeedShellProps) {
  return (
    <main
      className="min-h-dvh w-full bg-(--nox-bg) text-(--nox-ink)"
      style={authThemeVars(darkAuthTheme)}
    >
      <div className="mx-auto flex min-h-dvh w-full max-w-[390px] flex-col">
        {children}
      </div>
    </main>
  );
}
