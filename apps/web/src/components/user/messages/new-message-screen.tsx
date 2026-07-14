"use client";

import { ChevronLeft, Search } from "lucide-react";
import { useEffect, useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import { FeedShell } from "@/src/components/user/feed/feed-shell";
import { TabBar } from "@/src/components/user/feed/tab-bar";
import { Avatar } from "@/src/components/user/shared/avatar";
import type { Persona } from "@/src/types/api/user/persona";
import { createDirectConversation, createGroupConversation } from "@/src/utils/api/user/messaging";
import { searchNox } from "@/src/utils/api/user/search";
import { getAccessToken } from "@/src/utils/auth/session";
import { useActivePersona } from "@/src/hooks/use-active-persona";

export function NewMessageScreen() {
  const router = useRouter();
  const { activeID, loading: personaLoading } = useActivePersona();
  const [mode, setMode] = useState<"direct" | "group">("direct");
  const [query, setQuery] = useState("");
  const [title, setTitle] = useState("");
  const [results, setResults] = useState<Persona[]>([]);
  const [selected, setSelected] = useState<Persona[]>([]);
  const [loading, setLoading] = useState(false);
  const [startingID, setStartingID] = useState("");
  const [error, setError] = useState("");
  const token = useMemo(() => getAccessToken(), []);

  useEffect(() => {
    const trimmed = query.trim();
    if (!trimmed || trimmed.length < 2 || !token) {
      setResults([]);
      return;
    }
    const timeout = window.setTimeout(() => {
      setLoading(true);
      setError("");
      searchNox(trimmed, 12, token, activeID)
        .then((res) => setResults((res.data?.personas ?? []).filter((persona) => persona.id !== activeID)))
        .catch(() => setError("Could not search personas."))
        .finally(() => setLoading(false));
    }, 250);
    return () => window.clearTimeout(timeout);
  }, [activeID, query, token]);

  useEffect(() => {
    if (!personaLoading && !token) router.replace("/auth");
  }, [personaLoading, router, token]);

  async function startDirectConversation(recipientPersonaID: string) {
    if (!activeID || !token) {
      setError("Choose a persona before starting a message.");
      return;
    }
    setStartingID(recipientPersonaID);
    setError("");
    try {
      const res = await createDirectConversation(
        { sender_persona_id: activeID, recipient_persona_id: recipientPersonaID },
        token,
      );
      if (res.data) router.push(`/messages/${res.data.id}`);
    } catch {
      setError("Could not start conversation.");
    } finally {
      setStartingID("");
    }
  }

  function toggleSelected(persona: Persona) {
    setSelected((current) => {
      if (current.some((item) => item.id === persona.id)) {
        return current.filter((item) => item.id !== persona.id);
      }
      return [...current, persona];
    });
  }

  async function startGroupConversation() {
    if (!activeID || !token) {
      setError("Choose a persona before starting a group.");
      return;
    }
    if (!title.trim() || selected.length === 0) {
      setError("Add a group name and at least one member.");
      return;
    }
    setStartingID("group");
    setError("");
    try {
      const res = await createGroupConversation(
        {
          creator_persona_id: activeID,
          title: title.trim(),
          member_persona_ids: selected.map((persona) => persona.id),
        },
        token,
      );
      if (res.data) router.push(`/messages/${res.data.id}`);
    } catch {
      setError("Could not start group.");
    } finally {
      setStartingID("");
    }
  }

  return (
    <FeedShell>
      <header className="flex items-center gap-3 border-b border-(--nox-divider) px-4 pt-[env(safe-area-inset-top,12px)] pb-3">
        <button type="button" onClick={() => router.back()} className="flex size-8 items-center justify-center rounded-full hover:bg-(--nox-surface-alt)">
          <ChevronLeft className="size-5 text-(--nox-ink)" strokeWidth={1.8} />
        </button>
        <h1 className="text-[18px] font-bold text-(--nox-ink)">new message</h1>
      </header>

      <div className="border-b border-(--nox-divider) px-4 py-3">
        <div className="mb-3 grid grid-cols-2 gap-2">
          {(["direct", "group"] as const).map((item) => (
            <button
              key={item}
              type="button"
              onClick={() => setMode(item)}
              className={`rounded-[8px] border px-3 py-2 text-[12px] font-semibold ${mode === item ? "border-(--nox-accent) text-(--nox-accent)" : "border-(--nox-border) text-(--nox-ink-soft)"}`}
            >
              {item}
            </button>
          ))}
        </div>
        {mode === "group" && (
          <input
            value={title}
            onChange={(event) => setTitle(event.target.value)}
            placeholder="Group name"
            className="mb-3 w-full rounded-[8px] border border-(--nox-border) bg-(--nox-surface) px-3 py-2.5 text-[14px] text-(--nox-ink) outline-none placeholder:text-(--nox-ink-soft)"
          />
        )}
        <div className="flex items-center gap-2 rounded-[8px] border border-(--nox-border) bg-(--nox-surface) px-3 py-2.5">
          <Search className="size-4 shrink-0 text-(--nox-ink-soft)" strokeWidth={1.7} />
          <input
            value={query}
            onChange={(event) => setQuery(event.target.value)}
            placeholder="Search personas..."
            className="flex-1 bg-transparent text-[14px] text-(--nox-ink) outline-none placeholder:text-(--nox-ink-soft)"
          />
        </div>
        {mode === "group" && selected.length > 0 && (
          <div className="mt-3 flex flex-wrap gap-2">
            {selected.map((persona) => (
              <button
                key={persona.id}
                type="button"
                onClick={() => toggleSelected(persona)}
                className="rounded-full border border-(--nox-border) px-3 py-1 text-[11px] text-(--nox-ink-soft)"
              >
                @{persona.handle}
              </button>
            ))}
          </div>
        )}
        {error && <p className="mt-2 text-[11px] text-(--nox-danger)">{error}</p>}
      </div>

      <div className="flex-1 overflow-y-auto">
        <p className="px-4 pb-1 pt-4 font-mono text-[10px] font-semibold uppercase tracking-[0.18em] text-(--nox-ink-soft)">
          {query.trim() ? "results" : "search"}
        </p>
        {loading && <p className="px-4 py-6 text-center text-[12px] text-(--nox-ink-soft)">Searching...</p>}
        {!loading && query.trim().length >= 2 && results.length === 0 && (
          <p className="px-4 py-6 text-center text-[12px] text-(--nox-ink-soft)">No personas found.</p>
        )}
        {results.map((persona) => (
          <button
            key={persona.id}
            type="button"
            disabled={startingID === persona.id}
            onClick={() => (mode === "group" ? toggleSelected(persona) : void startDirectConversation(persona.id))}
            className="flex w-full items-center gap-3 border-b border-(--nox-divider) px-4 py-3.5 text-left transition hover:bg-(--nox-surface) disabled:opacity-50"
          >
            <Avatar id={persona.id} name={persona.display_name || persona.handle} size={40} />
            <div>
              <p className="text-[14px] font-semibold text-(--nox-ink)">{persona.display_name || persona.handle}</p>
              <p className="font-mono text-[11px] text-(--nox-ink-soft)">@{persona.handle}</p>
            </div>
            {mode === "group" && selected.some((item) => item.id === persona.id) && (
              <span className="ml-auto rounded-full bg-(--nox-accent) px-2 py-1 font-mono text-[10px] text-white">added</span>
            )}
          </button>
        ))}
      </div>

      {mode === "group" && (
        <div className="border-t border-(--nox-divider) px-4 py-3">
          <button
            type="button"
            disabled={startingID === "group" || !title.trim() || selected.length === 0}
            onClick={() => void startGroupConversation()}
            className="w-full rounded-[8px] bg-(--nox-accent) px-4 py-3 text-[13px] font-semibold text-white disabled:opacity-40"
          >
            create group
          </button>
        </div>
      )}

      <TabBar />
    </FeedShell>
  );
}
