"use client";

import { useState } from "react";
import { CalendarDays, Ghost, Grid2X2, Search, User } from "lucide-react";
import { AuthBrandPanel } from "@/src/components/user/auth/auth-brand-panel";
import { AuthModeTabs } from "@/src/components/user/auth/auth-mode-tabs";
import { AuthShell } from "@/src/components/user/auth/auth-shell";
// import { AuthStatusBar } from "@/src/components/user/auth/auth-status-bar";
import { darkAuthTheme } from "@/src/components/user/auth/auth-theme";
import { LoginForm } from "@/src/components/user/auth/login-form";
import { SignupForm } from "@/src/components/user/auth/signup-form";
import type { AuthMode } from "@/src/types/components/auth";

const navItems = [
  { label: "feed", icon: Grid2X2 },
  { label: "discover", icon: Search },
  { label: "events", icon: CalendarDays },
  { label: "profile", icon: User },
];

export function AuthScreen() {
  const [mode, setMode] = useState<AuthMode>("login");
  const isLogin = mode === "login";

  return (
    <AuthShell theme={darkAuthTheme} sidePanel={<AuthBrandPanel />}>
      {/* <AuthStatusBar /> */}
      <div className="flex flex-1 flex-col px-5 pb-0">
        <header className="pt-7">
          <div className="mb-5 flex size-12 items-center justify-center rounded-[12px] border border-(--nox-accent-line) bg-(--nox-accent-soft) text-(--nox-accent-ink)">
            <Ghost className="size-6" strokeWidth={1.6} />
          </div>
          <p className="font-mono text-[10px] font-medium uppercase tracking-[0.18em] text-(--nox-ink-soft)">
            nox access
          </p>
          <h1 className="mt-3 text-[34px] font-bold leading-none tracking-[-0.04em] text-(--nox-ink)">
            {isLogin ? "Get back in" : "Join the feed"}
          </h1>
          <p className="mt-3 max-w-[280px] text-[14px] leading-6 text-(--nox-ink-mid)">
            {isLogin
              ? "Keep your public name clean and your anonymous posts unlinked."
              : "Start with a private account, then choose when to post publicly."}
          </p>
        </header>

        <AuthModeTabs mode={mode} onModeChange={setMode} />
        {isLogin ? (
          <LoginForm theme={darkAuthTheme} className="mt-5" />
        ) : (
          <SignupForm theme={darkAuthTheme} className="mt-5" />
        )}

        <div className="mt-auto border-t border-(--nox-divider) pt-3">
          <nav className="grid grid-cols-4 border-t border-(--nox-divider) pt-2 pb-3">
            {navItems.map((item, index) => {
              const Icon = item.icon;
              const active = index === 0;
              return (
                <span
                  key={item.label}
                  className={`flex flex-col items-center gap-1 text-[9.5px] font-medium ${
                    active ? "text-(--nox-accent)" : "text-(--nox-ink-soft)"
                  }`}
                >
                  <Icon className="size-4" strokeWidth={1.7} />
                  {item.label}
                </span>
              );
            })}
          </nav>
        </div>
      </div>
    </AuthShell>
  );
}
