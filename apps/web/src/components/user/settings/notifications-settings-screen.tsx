"use client";

import { useState } from "react";
import { ChevronLeft } from "lucide-react";
import { useRouter } from "next/navigation";
import { FeedShell } from "@/src/components/user/feed/feed-shell";
import { TabBar } from "@/src/components/user/feed/tab-bar";

interface ToggleRowProps { label: string; sub?: string; value: boolean; onChange: (v: boolean) => void }
function ToggleRow({ label, sub, value, onChange }: ToggleRowProps) {
  return (
    <div className="flex items-center justify-between border-b border-(--nox-divider) px-4 py-4">
      <div>
        <p className="text-[14px] font-medium text-(--nox-ink)">{label}</p>
        {sub && <p className="mt-0.5 text-[11px] text-(--nox-ink-soft)">{sub}</p>}
      </div>
      <button type="button" onClick={() => onChange(!value)} role="switch" aria-checked={value}
        className={`relative h-6 w-11 rounded-full transition-colors ${value ? "bg-(--nox-accent)" : "bg-(--nox-surface-alt)"}`}>
        <span className={`absolute top-0.5 h-5 w-5 rounded-full bg-white shadow transition-transform ${value ? "translate-x-5" : "translate-x-0.5"}`} />
      </button>
    </div>
  );
}

export function NotificationsSettingsScreen() {
  const router = useRouter();
  const [likes, setLikes] = useState(true);
  const [follows, setFollows] = useState(true);
  const [comments, setComments] = useState(true);
  const [stories, setStories] = useState(false);
  const [events, setEvents] = useState(true);
  const [push, setPush] = useState(false);

  return (
    <FeedShell>
      <header className="flex items-center gap-3 border-b border-(--nox-divider) px-4 pt-[env(safe-area-inset-top,12px)] pb-3">
        <button type="button" onClick={() => router.back()} className="flex size-8 items-center justify-center rounded-full hover:bg-(--nox-surface-alt)">
          <ChevronLeft className="size-5 text-(--nox-ink)" strokeWidth={1.8} />
        </button>
        <h1 className="text-[18px] font-bold text-(--nox-ink)">notifications</h1>
      </header>

      <div className="flex-1 overflow-y-auto">
        <p className="px-4 pb-1 pt-5 font-mono text-[10px] font-semibold uppercase tracking-[0.18em] text-(--nox-ink-soft)">activity</p>
        <ToggleRow label="likes" sub="When someone likes your post or story" value={likes} onChange={setLikes} />
        <ToggleRow label="new followers" sub="When someone follows your persona" value={follows} onChange={setFollows} />
        <ToggleRow label="comments" sub="Replies and comments on your posts" value={comments} onChange={setComments} />
        <ToggleRow label="story contributions" sub="When someone adds a clip to your story" value={stories} onChange={setStories} />
        <ToggleRow label="event reminders" sub="24h before events you're attending" value={events} onChange={setEvents} />

        <p className="px-4 pb-1 pt-5 font-mono text-[10px] font-semibold uppercase tracking-[0.18em] text-(--nox-ink-soft)">push</p>
        <ToggleRow label="push notifications" sub="Browser / device push alerts" value={push} onChange={setPush} />

        <p className="px-4 py-6 text-[12px] text-(--nox-ink-soft)">
          Push notification delivery is coming soon. These preferences will be saved when the feature launches.
        </p>
      </div>

      <TabBar />
    </FeedShell>
  );
}
