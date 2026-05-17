import { BatteryFull, Signal, Wifi } from "lucide-react";

export function AuthStatusBar() {
  return (
    <div className="flex h-[38px] shrink-0 items-center justify-between px-5 text-[12px] font-semibold text-(--nox-ink)">
      <span>9:41</span>
      <div className="flex items-center gap-1.5 text-(--nox-ink-mid)">
        <Signal className="size-3.5" strokeWidth={1.7} />
        <Wifi className="size-3.5" strokeWidth={1.7} />
        <BatteryFull className="size-4" strokeWidth={1.7} />
      </div>
    </div>
  );
}
