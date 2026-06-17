"use client";

export function SetSkeleton() {
  return (
    <div className="grid animate-pulse grid-cols-[144px_1fr] gap-3 border-b border-(--nox-divider) px-4 py-4">
      <div className="aspect-video rounded-[8px] bg-(--nox-surface-alt)" />
      <div className="min-w-0 space-y-2 py-0.5">
        <div className="h-[14px] w-3/4 rounded bg-(--nox-surface-alt)" />
        <div className="h-3 w-1/2 rounded bg-(--nox-surface-alt)" />
        <div className="mt-3 h-3 w-2/5 rounded bg-(--nox-surface-alt)" />
      </div>
    </div>
  );
}
