import { Ghost, Music2, Sparkles, Users } from "lucide-react";
import type { ReactNode } from "react";

export function OnboardingBrandPanel() {
  return (
    <div className="flex h-full flex-col justify-between rounded-[18px] border border-(--nox-border) bg-(--nox-surface) p-8 shadow-(--nox-shadow)">
      <div>
        <div className="mb-8 inline-flex items-center gap-2 rounded-full border border-(--nox-accent-line) bg-(--nox-accent-soft) px-3 py-1.5 font-mono text-[10px] font-medium uppercase tracking-[0.18em] text-(--nox-accent-ink)">
          <Sparkles className="size-3.5" strokeWidth={1.7} />
          Setting up your nox
        </div>
        <h2 className="max-w-xl text-[48px] font-bold leading-[0.96] tracking-[-0.04em] text-(--nox-ink) lg:text-[56px]">
          your scene,<br />your rules.
        </h2>
        <p className="mt-5 max-w-md text-[15px] leading-7 text-(--nox-ink-mid)">
          Post publicly under your name or stay ghost. Switch modes per post — your identity always stays in your control.
        </p>
      </div>

      <div className="grid gap-3">
        <PanelRow
          icon={<Ghost className="size-4" />}
          label="Ghost posts are never tied to your real account publicly"
        />
        <PanelRow
          icon={<Users className="size-4" />}
          label="Visible persona builds your following and set archive"
        />
        <PanelRow
          icon={<Music2 className="size-4" />}
          label="Genre tags help the right crowd find your posts"
        />
      </div>
    </div>
  );
}

function PanelRow({ icon, label }: { icon: ReactNode; label: string }) {
  return (
    <div className="flex items-center gap-3 rounded-[10px] border border-(--nox-border) bg-(--nox-surface-alt) px-4 py-3 text-[13px] font-medium text-(--nox-ink-mid)">
      <span className="text-(--nox-accent-ink)">{icon}</span>
      <span className="flex-1">{label}</span>
    </div>
  );
}
