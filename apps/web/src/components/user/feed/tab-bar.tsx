"use client";

import { CalendarDays, Grid2X2, ListMusic, Search, User } from "lucide-react";
import { useRouter, usePathname } from "next/navigation";

const tabs = [
  { label: "feed", icon: Grid2X2, href: "/feed" },
  { label: "discover", icon: Search, href: "/discover" },
  { label: "sets", icon: ListMusic, href: "/sets" },
  { label: "events", icon: CalendarDays, href: "/events" },
  { label: "profile", icon: User, href: "/profile" },
];

export function TabBar() {
  const router = useRouter();
  const pathname = usePathname();

  return (
    <nav className="border-t border-(--nox-divider) bg-(--nox-bg)">
      <div className="grid grid-cols-5 pb-[env(safe-area-inset-bottom,12px)] pt-2">
        {tabs.map(({ label, icon: Icon, href }) => {
          const active = pathname.startsWith(href);
          return (
            <button
              key={label}
              type="button"
              onClick={() => router.push(href)}
              className="flex flex-col items-center gap-1 text-[9.5px] font-medium"
              style={{ color: active ? "var(--nox-accent)" : "var(--nox-ink-soft)" }}
            >
              <Icon className="size-4" strokeWidth={1.7} />
              {label}
            </button>
          );
        })}
      </div>
    </nav>
  );
}
