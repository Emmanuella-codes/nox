"use client";

import { Avatar } from "@/src/components/user/shared/avatar";
import type { Persona } from "@/src/types/api/user/persona";

interface PersonaCardProps {
  persona: Persona;
  onPress?: () => void;
  showFollow?: boolean;
  isFollowing?: boolean;
  onFollow?: (persona: Persona) => void;
}

export function PersonaCard({
  persona,
  onPress,
  showFollow = false,
  isFollowing = false,
  onFollow,
}: PersonaCardProps) {
  return (
    <div
      className="flex items-center gap-3 px-4 py-3 transition active:bg-(--nox-surface)"
      style={{ cursor: onPress ? "pointer" : "default" }}
      onClick={onPress}
    >
      <Avatar id={persona.id} name={persona.display_name} size={42} />

      <div className="min-w-0 flex-1">
        <p className="text-[14px] font-semibold leading-tight text-(--nox-ink) truncate">
          {persona.display_name}
        </p>
        <p className="text-[12px] text-(--nox-ink-soft)">@{persona.handle}</p>
        {persona.genre_tags.length > 0 && (
          <div className="mt-1.5 flex flex-wrap gap-1">
            {persona.genre_tags.slice(0, 3).map((tag) => (
              <span
                key={tag}
                className="rounded-full px-2 py-0.5 font-mono text-[9.5px] font-medium lowercase"
                style={{
                  background: "var(--nox-accent-soft)",
                  color: "var(--nox-accent-ink)",
                }}
              >
                {tag}
              </span>
            ))}
          </div>
        )}
      </div>

      {showFollow && (
        <button
          type="button"
          onClick={(e) => {
            e.stopPropagation();
            onFollow?.(persona);
          }}
          className="shrink-0 rounded-full border px-3.5 py-1.5 text-[12px] font-semibold transition hover:border-(--nox-accent) hover:text-(--nox-accent)"
          style={{
            borderColor: isFollowing ? "var(--nox-accent-line)" : "var(--nox-border-strong)",
            color: isFollowing ? "var(--nox-accent-ink)" : "var(--nox-ink-mid)",
          }}
        >
          {isFollowing ? "following" : "follow"}
        </button>
      )}
    </div>
  );
}
