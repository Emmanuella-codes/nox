"use client";

import { Bell, ChevronRight, CreditCard, Lock, LogOut, User } from "lucide-react";
import { useRouter } from "next/navigation";
import { FeedShell } from "@/src/components/user/feed/feed-shell";
import { TabBar } from "@/src/components/user/feed/tab-bar";
import { Avatar } from "@/src/components/user/shared/avatar";
import { useActivePersona } from "@/src/hooks/use-active-persona";
import { clearSession } from "@/src/utils/auth/session";

const SECTIONS = [
  {
    title: "account",
    items: [
      { label: "edit profile", icon: User, href: "/profile/edit" },
      { label: "notifications", icon: Bell, href: "/settings/notifications" },
      { label: "privacy", icon: Lock, href: "/settings/privacy" },
      { label: "payments & history", icon: CreditCard, href: "/settings/payments" },
    ],
  },
];

export function SettingsScreen() {
  const router = useRouter();
  const { activePersona, loading } = useActivePersona();

  function handleLogout() {
    clearSession();
    router.push("/auth");
  }

  return (
    <FeedShell>
      <header className="px-4 pt-[env(safe-area-inset-top,12px)] pb-3">
        <h1 className="text-[22px] font-bold tracking-[-0.03em] text-(--nox-ink)">settings</h1>
      </header>

      <div className="flex-1 overflow-y-auto">
        {!loading && activePersona && (
          <button type="button" onClick={() => router.push("/profile/edit")}
            className="flex w-full items-center gap-3 border-b border-(--nox-divider) px-4 py-4 transition hover:bg-(--nox-surface)">
            <Avatar id={activePersona.id} name={activePersona.display_name} size={44} square />
            <div className="min-w-0 flex-1 text-left">
              <p className="truncate text-[15px] font-bold text-(--nox-ink)">{activePersona.display_name}</p>
              <p className="font-mono text-[11px] text-(--nox-ink-soft)">@{activePersona.handle} · {activePersona.category}</p>
            </div>
            <ChevronRight className="size-4 shrink-0 text-(--nox-ink-faint)" strokeWidth={1.7} />
          </button>
        )}

        {SECTIONS.map((section) => (
          <div key={section.title}>
            <p className="px-4 pb-1 pt-5 font-mono text-[10px] font-semibold uppercase tracking-[0.18em] text-(--nox-ink-soft)">
              {section.title}
            </p>
            {section.items.map((item) => {
              const Icon = item.icon;
              return (
                <button key={item.href} type="button" onClick={() => router.push(item.href)}
                  className="flex w-full items-center gap-3 border-b border-(--nox-divider) px-4 py-4 transition hover:bg-(--nox-surface)">
                  <span className="flex size-8 items-center justify-center rounded-[8px] bg-(--nox-surface-alt)">
                    <Icon className="size-4 text-(--nox-ink-mid)" strokeWidth={1.7} />
                  </span>
                  <span className="flex-1 text-left text-[14px] font-medium text-(--nox-ink)">{item.label}</span>
                  <ChevronRight className="size-4 shrink-0 text-(--nox-ink-faint)" strokeWidth={1.7} />
                </button>
              );
            })}
          </div>
        ))}

        <div className="mt-4 px-4">
          <button type="button" onClick={handleLogout}
            className="flex w-full items-center gap-3 rounded-[10px] border border-(--nox-border) px-4 py-3.5 text-[14px] font-medium text-(--nox-danger) transition hover:bg-(--nox-danger-soft)">
            <LogOut className="size-4" strokeWidth={1.7} />
            log out
          </button>
        </div>

        <div className="py-8 text-center">
          <p className="font-mono text-[10px] text-(--nox-ink-faint)">nox · lagos nightlife · no filter</p>
        </div>
      </div>

      <TabBar />
    </FeedShell>
  );
}
