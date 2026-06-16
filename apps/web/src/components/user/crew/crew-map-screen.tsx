"use client";

import { ChevronLeft, MapPin, Users } from "lucide-react";
import { useRouter } from "next/navigation";
import { FeedShell } from "@/src/components/user/feed/feed-shell";
import { Avatar } from "@/src/components/user/shared/avatar";

interface CrewMapScreenProps { crewID: string }

const SAMPLE_MEMBERS = [
  { id: "m1", name: "Amirah Lagos", handle: "amirah_lagos", location: "near stage" },
  { id: "m2", name: "DJ Khalid", handle: "djkhalid", location: "backstage" },
  { id: "m3", name: "Scene Curator", handle: "scene_curator", location: "bar" },
];

export function CrewMapScreen({ crewID }: CrewMapScreenProps) {
  const router = useRouter();

  return (
    <FeedShell>
      <header className="flex items-center gap-3 border-b border-(--nox-divider) px-4 pt-[env(safe-area-inset-top,12px)] pb-3">
        <button type="button" onClick={() => router.back()} className="flex size-8 items-center justify-center rounded-full hover:bg-(--nox-surface-alt)">
          <ChevronLeft className="size-5 text-(--nox-ink)" strokeWidth={1.8} />
        </button>
        <h1 className="text-[18px] font-bold text-(--nox-ink)">live map</h1>
        <button type="button" onClick={() => router.push(`/crews/${crewID}/members`)}
          className="ml-auto flex items-center gap-1.5 rounded-[8px] border border-(--nox-border-strong) px-3 py-1.5 text-[12px] font-semibold text-(--nox-ink-mid) hover:border-(--nox-accent)">
          <Users className="size-3.5" strokeWidth={1.7} />
          {SAMPLE_MEMBERS.length}
        </button>
      </header>

      {/* Map placeholder */}
      <div className="relative mx-4 mt-4 h-64 overflow-hidden rounded-[12px] bg-(--nox-surface-alt)">
        <div className="flex size-full flex-col items-center justify-center gap-2">
          <MapPin className="size-8 text-(--nox-accent)" strokeWidth={1.5} />
          <p className="text-[13px] text-(--nox-ink-soft)">Live map coming soon</p>
          <p className="text-[11px] text-(--nox-ink-faint)">Member pins will appear here during the event</p>
        </div>
        {/* Sample pins */}
        {SAMPLE_MEMBERS.map((m, i) => (
          <div key={m.id} className="absolute flex flex-col items-center"
            style={{ top: `${20 + i * 25}%`, left: `${25 + i * 20}%` }}>
            <Avatar id={m.id} name={m.name} size={28} />
            <div className="mt-1 rounded-[4px] bg-black/60 px-1.5 py-0.5">
              <p className="font-mono text-[8px] text-white/80">@{m.handle}</p>
            </div>
          </div>
        ))}
      </div>

      <div className="mt-4 flex-1 overflow-y-auto px-4">
        <p className="mb-2 font-mono text-[10px] font-semibold uppercase tracking-[0.18em] text-(--nox-ink-soft)">crew nearby</p>
        {SAMPLE_MEMBERS.map((m) => (
          <div key={m.id} className="flex items-center gap-3 border-b border-(--nox-divider) py-3">
            <Avatar id={m.id} name={m.name} size={36} />
            <div className="flex-1">
              <p className="text-[13px] font-semibold text-(--nox-ink)">{m.name}</p>
              <p className="flex items-center gap-1 font-mono text-[10px] text-(--nox-ink-soft)">
                <MapPin className="size-2.5" strokeWidth={1.7} /> {m.location}
              </p>
            </div>
          </div>
        ))}
      </div>
    </FeedShell>
  );
}
