"use client";

import { useRouter } from "next/navigation";
import { FeedShell } from "@/src/components/user/feed/feed-shell";
import { TabBar } from "@/src/components/user/feed/tab-bar";
import { Avatar } from "@/src/components/user/shared/avatar";
import { Edit } from "lucide-react";

const SAMPLE = [
  { id: "c1", name: "Amirah Lagos", handle: "amirah_lagos", lastMsg: "See you on the floor 🔥", time: "2m", unread: 2 },
  { id: "c2", name: "DJ Khalid", handle: "djkhalid", lastMsg: "Set starts at midnight bro", time: "14m", unread: 0 },
  { id: "c3", name: "Warehouse Sessions", handle: "warehouse_nox", lastMsg: "Thanks for coming through!", time: "1h", unread: 0 },
  { id: "c4", name: "Scene Curator", handle: "scene_curator", lastMsg: "Are you going Saturday?", time: "3h", unread: 1 },
];

export function MessagesScreen() {
  const router = useRouter();

  return (
    <FeedShell>
      <header className="flex items-center justify-between px-4 pt-[env(safe-area-inset-top,12px)] pb-3">
        <h1 className="text-[22px] font-bold tracking-[-0.03em] text-(--nox-ink)">messages</h1>
        <button type="button" onClick={() => router.push("/messages/new")}
          className="flex size-8 items-center justify-center rounded-full border border-(--nox-border-strong) hover:border-(--nox-accent)">
          <Edit className="size-4 text-(--nox-ink-mid)" strokeWidth={1.7} />
        </button>
      </header>

      <div className="flex-1 overflow-y-auto">
        {SAMPLE.map((c) => (
          <button key={c.id} type="button" onClick={() => router.push(`/messages/${c.id}`)}
            className="flex w-full items-center gap-3 border-b border-(--nox-divider) px-4 py-4 text-left transition hover:bg-(--nox-surface)">
            <Avatar id={c.id} name={c.name} size={44} />
            <div className="min-w-0 flex-1">
              <p className="truncate text-[14px] font-semibold text-(--nox-ink)">{c.name}</p>
              <p className="truncate text-[12px] text-(--nox-ink-soft)">{c.lastMsg}</p>
            </div>
            <div className="flex shrink-0 flex-col items-end gap-1">
              <span className="font-mono text-[10px] text-(--nox-ink-faint)">{c.time}</span>
              {c.unread > 0 && (
                <span className="flex size-5 items-center justify-center rounded-full bg-(--nox-accent) font-mono text-[10px] font-semibold text-white">
                  {c.unread}
                </span>
              )}
            </div>
          </button>
        ))}
        <p className="px-4 py-8 text-center text-[12px] text-(--nox-ink-soft)">Messaging API coming soon.</p>
      </div>

      <TabBar />
    </FeedShell>
  );
}
