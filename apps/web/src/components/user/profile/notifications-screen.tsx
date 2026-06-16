"use client";

import { Bell, Heart, MessageCircle, UserPlus } from "lucide-react";
import { FeedShell } from "@/src/components/user/feed/feed-shell";
import { TabBar } from "@/src/components/user/feed/tab-bar";
import { Avatar } from "@/src/components/user/shared/avatar";

const SAMPLE = [
  { id: "1", type: "like", actor: "djkhalid", actorID: "a1", text: "liked your post", time: "2m ago" },
  { id: "2", type: "follow", actor: "amirah_lagos", actorID: "a2", text: "started following you", time: "14m ago" },
  { id: "3", type: "comment", actor: "warehouse_nox", actorID: "a3", text: "commented on your post", time: "1h ago" },
  { id: "4", type: "like", actor: "afrohousevibes", actorID: "a4", text: "liked your story clip", time: "3h ago" },
  { id: "5", type: "follow", actor: "scene_curator", actorID: "a5", text: "started following you", time: "1d ago" },
];

function NotifIcon({ type }: { type: string }) {
  if (type === "like") return <Heart className="size-3.5 text-(--nox-danger)" strokeWidth={1.7} />;
  if (type === "follow") return <UserPlus className="size-3.5 text-(--nox-accent)" strokeWidth={1.7} />;
  if (type === "comment") return <MessageCircle className="size-3.5 text-(--nox-ink-mid)" strokeWidth={1.7} />;
  return <Bell className="size-3.5 text-(--nox-ink-mid)" strokeWidth={1.7} />;
}

export function NotificationsScreen() {
  return (
    <FeedShell>
      <header className="px-4 pt-[env(safe-area-inset-top,12px)] pb-3">
        <h1 className="text-[22px] font-bold tracking-[-0.03em] text-(--nox-ink)">notifications</h1>
      </header>

      <div className="flex-1 overflow-y-auto">
        {SAMPLE.map((n) => (
          <div key={n.id} className="flex items-center gap-3 border-b border-(--nox-divider) px-4 py-3.5">
            <div className="relative shrink-0">
              <Avatar id={n.actorID} name={n.actor} size={38} />
              <span className="absolute -right-1 -bottom-1 flex size-5 items-center justify-center rounded-full bg-(--nox-surface) border border-(--nox-divider)">
                <NotifIcon type={n.type} />
              </span>
            </div>
            <div className="min-w-0 flex-1">
              <p className="text-[13px] text-(--nox-ink)">
                <span className="font-semibold">@{n.actor}</span>{" "}{n.text}
              </p>
            </div>
            <span className="font-mono text-[10px] text-(--nox-ink-faint) shrink-0">{n.time}</span>
          </div>
        ))}

        <p className="px-4 py-8 text-center text-[12px] text-(--nox-ink-soft)">
          Real-time notifications coming soon.
        </p>
      </div>

      <TabBar />
    </FeedShell>
  );
}
