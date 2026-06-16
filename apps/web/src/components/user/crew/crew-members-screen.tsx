"use client";

import { ChevronLeft, MapPin } from "lucide-react";
import { useRouter } from "next/navigation";
import { FeedShell } from "@/src/components/user/feed/feed-shell";
import { Avatar } from "@/src/components/user/shared/avatar";

interface CrewMembersScreenProps { crewID: string }

const SAMPLE_MEMBERS = [
  { id: "m1", name: "You (Amirah Lagos)", handle: "amirah_lagos", role: "organizer", location: "near stage", sharing: true },
  { id: "m2", name: "DJ Khalid", handle: "djkhalid", role: "member", location: "backstage", sharing: true },
  { id: "m3", name: "Scene Curator", handle: "scene_curator", role: "member", location: "bar", sharing: true },
  { id: "m4", name: "Afro House Vibes", handle: "afrohousevibes", role: "member", location: null, sharing: false },
];

export function CrewMembersScreen({ crewID }: CrewMembersScreenProps) {
  const router = useRouter();

  return (
    <FeedShell>
      <header className="flex items-center gap-3 border-b border-(--nox-divider) px-4 pt-[env(safe-area-inset-top,12px)] pb-3">
        <button type="button" onClick={() => router.back()} className="flex size-8 items-center justify-center rounded-full hover:bg-(--nox-surface-alt)">
          <ChevronLeft className="size-5 text-(--nox-ink)" strokeWidth={1.8} />
        </button>
        <div>
          <h1 className="text-[18px] font-bold text-(--nox-ink)">crew members</h1>
          <p className="text-[11px] text-(--nox-ink-soft)">crew {crewID.slice(0, 8)} · {SAMPLE_MEMBERS.length} members</p>
        </div>
      </header>

      <div className="flex-1 overflow-y-auto">
        {SAMPLE_MEMBERS.map((m) => (
          <div key={m.id} className="flex items-center gap-3 border-b border-(--nox-divider) px-4 py-4">
            <Avatar id={m.id} name={m.name} size={42} />
            <div className="min-w-0 flex-1">
              <div className="flex items-center gap-2">
                <p className="truncate text-[14px] font-semibold text-(--nox-ink)">{m.name}</p>
                {m.role === "organizer" && (
                  <span className="shrink-0 rounded-full bg-(--nox-accent-soft) px-2 py-0.5 font-mono text-[8px] font-semibold text-(--nox-accent-ink)">host</span>
                )}
              </div>
              <p className="font-mono text-[10px] text-(--nox-ink-soft)">@{m.handle}</p>
            </div>
            <div className="shrink-0 text-right">
              {m.sharing && m.location ? (
                <p className="flex items-center gap-1 text-[11px] text-(--nox-ink-soft)">
                  <MapPin className="size-3 text-(--nox-success)" strokeWidth={1.7} />
                  {m.location}
                </p>
              ) : (
                <p className="font-mono text-[10px] text-(--nox-ink-faint)">not sharing</p>
              )}
            </div>
          </div>
        ))}
        <p className="px-4 py-6 text-center text-[12px] text-(--nox-ink-soft)">Location sharing expires when the event ends.</p>
      </div>
    </FeedShell>
  );
}
