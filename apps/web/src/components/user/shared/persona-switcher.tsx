"use client";

import { useState } from "react";
import { X } from "lucide-react";
import { Avatar } from "@/src/components/user/shared/avatar";
import type { Persona } from "@/src/types/api/user/persona";

interface PersonaSwitcherProps {
  personas: Persona[];
  activePersona: Persona | null;
  onSwitch: (id: string) => void;
}

export function PersonaSwitcher({ personas, activePersona, onSwitch }: PersonaSwitcherProps) {
  const [open, setOpen] = useState(false);

  if (!activePersona) return null;

  function handleSwitch(id: string) {
    onSwitch(id);
    setOpen(false);
  }

  return (
    <>
      <button
        type="button"
        onClick={() => personas.length > 1 ? setOpen(true) : undefined}
        className="flex items-center gap-2 rounded-full border border-(--nox-border-strong) px-2.5 py-1.5 transition active:bg-(--nox-surface)"
      >
        <Avatar id={activePersona.id} name={activePersona.display_name} size={18} />
        <span className="font-mono text-[10px] font-semibold text-(--nox-ink-mid)">
          @{activePersona.handle}
        </span>
      </button>

      {open && (
        <div
          className="fixed inset-0 z-50 flex flex-col justify-end"
          style={{ background: "rgba(0,0,0,.6)" }}
          onClick={() => setOpen(false)}
        >
          <div
            className="mx-auto w-full max-w-[390px] rounded-t-[20px] bg-(--nox-bg) pb-[env(safe-area-inset-bottom,16px)]"
            onClick={(e) => e.stopPropagation()}
          >
            <div className="flex items-center justify-between border-b border-(--nox-divider) px-4 py-4">
              <span className="font-mono text-[11px] font-semibold uppercase tracking-[0.12em] text-(--nox-ink-soft)">
                switch persona
              </span>
              <button
                type="button"
                onClick={() => setOpen(false)}
                className="flex size-7 items-center justify-center rounded-full hover:bg-(--nox-surface-alt)"
              >
                <X className="size-4 text-(--nox-ink-mid)" strokeWidth={1.8} />
              </button>
            </div>
            {personas.map((p) => (
              <button
                key={p.id}
                type="button"
                onClick={() => handleSwitch(p.id)}
                className="flex w-full items-center gap-3 px-4 py-3 text-left transition hover:bg-(--nox-surface)"
              >
                <Avatar id={p.id} name={p.display_name} size={40} />
                <div className="min-w-0 flex-1">
                  <p className="truncate text-[14px] font-semibold text-(--nox-ink)">{p.display_name}</p>
                  <p className="text-[12px] text-(--nox-ink-soft)">@{p.handle}</p>
                </div>
                {p.id === activePersona.id && (
                  <span
                    className="shrink-0 rounded-full px-2.5 py-1 font-mono text-[10px] font-medium"
                    style={{ background: "var(--nox-accent-soft)", color: "var(--nox-accent-ink)" }}
                  >
                    active
                  </span>
                )}
              </button>
            ))}
          </div>
        </div>
      )}
    </>
  );
}
