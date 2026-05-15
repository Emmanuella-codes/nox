import { Check, Ghost, MapPin, Music2, ShieldCheck } from "lucide-react";
import type { ReactNode } from "react";

export function AuthBrandPanel() {
  return (
    <div className="flex h-full flex-col justify-between rounded-[18px] border border-(--nox-border) bg-(--nox-surface) p-8 shadow-(--nox-shadow)">
      <div>
        <div className="mb-8 inline-flex items-center gap-2 rounded-full border border-(--nox-accent-line) bg-(--nox-accent-soft) px-3 py-1.5 font-mono text-[10px] font-medium uppercase tracking-[0.18em] text-(--nox-accent-ink)">
          <MapPin className="size-3.5" strokeWidth={1.7} />
          Lagos after dark
        </div>
        <h1 className="max-w-xl text-[48px] font-bold leading-[0.96] tracking-[-0.04em] text-(--nox-ink) lg:text-[56px]">
          nox
        </h1>
        <p className="mt-5 max-w-md text-[15px] leading-7 text-(--nox-ink-mid)">
          Find the room, keep the memory, post without dragging your whole name through the night.
        </p>
      </div>

      <div className="grid gap-3">
        <PanelRow icon={<Ghost className="size-4" />} label="anonymous posts stay unlinked in public" />
        <PanelRow icon={<Music2 className="size-4" />} label="sets, venues, and scene notes in one feed" />
        <PanelRow icon={<ShieldCheck className="size-4" />} label="account ownership stays private and protected" />
      </div>
    </div>
  );
}

function PanelRow({ icon, label }: { icon: ReactNode; label: string }) {
  return (
    <div className="flex items-center gap-3 rounded-[10px] border border-(--nox-border) bg-(--nox-surface-alt) px-4 py-3 text-[13px] font-medium text-(--nox-ink-mid)">
      <span className="text-(--nox-accent-ink)">{icon}</span>
      <span className="flex-1">{label}</span>
      <Check className="size-4 text-(--nox-success)" strokeWidth={1.8} />
    </div>
  );
}
