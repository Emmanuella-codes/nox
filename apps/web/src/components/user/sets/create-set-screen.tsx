"use client";

import type { FormEvent } from "react";
import { useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import { ChevronLeft } from "lucide-react";
import { FeedShell } from "@/src/components/user/feed/feed-shell";
import { TabBar } from "@/src/components/user/feed/tab-bar";
import { AuthField } from "@/src/components/user/auth/auth-field";
import { useActivePersona } from "@/src/hooks/use-active-persona";
import { getAccessToken } from "@/src/utils/auth/session";
import { createSet, createSetVideoAsset } from "@/src/utils/api/user/set";
import { ApiRequestError } from "@/src/utils/api/api";

function deriveStorageKey(url: string): string {
  try {
    return new URL(url).pathname.replace(/^\//, "");
  } catch {
    return "";
  }
}

export function CreateSetScreen() {
  const router = useRouter();
  const { activePersona, loading: personaLoading } = useActivePersona();
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [playbackURL, setPlaybackURL] = useState("");
  const [thumbnailURL, setThumbnailURL] = useState("");
  const [storageKey, setStorageKey] = useState("");
  const [duration, setDuration] = useState("");
  const [genres, setGenres] = useState("");
  const [status, setStatus] = useState<"idle" | "loading" | "error">("idle");
  const [message, setMessage] = useState("");

  const durationSeconds = Number(duration);
  const canCreate = activePersona?.category === "dj";
  const canSubmit = useMemo(
    () =>
      Boolean(canCreate) &&
      title.trim().length > 0 &&
      playbackURL.trim().length > 0 &&
      durationSeconds > 0 &&
      durationSeconds <= 900 &&
      status !== "loading",
    [canCreate, durationSeconds, playbackURL, status, title],
  );

  function handlePlaybackURLChange(v: string) {
    setPlaybackURL(v);
    setStorageKey(deriveStorageKey(v));
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const token = getAccessToken();
    if (!token || !activePersona) return;
    if (!canSubmit) {
      setStatus("error");
      setMessage("Complete the set details. Duration must be 900 seconds or less.");
      return;
    }

    setStatus("loading");
    setMessage("");
    try {
      const mediaRes = await createSetVideoAsset(
        {
          owner_persona_id: activePersona.id,
          storage_key: storageKey || deriveStorageKey(playbackURL.trim()),
          playback_url: playbackURL.trim(),
          thumbnail_url: thumbnailURL.trim() || undefined,
          mime_type: "video/mp4",
          duration_seconds: durationSeconds,
          size_bytes: 0,
        },
        token,
      );
      const mediaID = mediaRes.data?.id;
      if (!mediaID) throw new Error("media_asset_missing");

      const setRes = await createSet(
        {
          persona_id: activePersona.id,
          media_asset_id: mediaID,
          title: title.trim(),
          description: description.trim(),
          genre_tags: genres.split(",").map((g) => g.trim().toLowerCase()).filter(Boolean),
        },
        token,
      );
      router.push(setRes.data?.id ? `/sets/${setRes.data.id}` : "/sets");
    } catch (error) {
      setStatus("error");
      setMessage(error instanceof ApiRequestError ? error.message : "Could not create set.");
    }
  }

  return (
    <FeedShell>
      <header className="flex items-center gap-3 border-b border-(--nox-divider) px-4 pt-[env(safe-area-inset-top,12px)] pb-3">
        <button
          type="button"
          onClick={() => router.back()}
          className="flex size-8 items-center justify-center rounded-full transition hover:bg-(--nox-surface-alt)"
        >
          <ChevronLeft className="size-5 text-(--nox-ink)" strokeWidth={1.8} />
        </button>
        <div>
          <h1 className="text-[18px] font-bold text-(--nox-ink)">new set</h1>
          <p className="text-[11px] text-(--nox-ink-soft)">15 minutes max.</p>
        </div>
      </header>

      <form onSubmit={handleSubmit} className="flex-1 overflow-y-auto px-4 py-4">
        {personaLoading ? (
          <p className="text-[13px] text-(--nox-ink-soft)">Loading persona...</p>
        ) : !canCreate ? (
          <p className="rounded-[8px] border border-(--nox-border) bg-(--nox-surface) px-3 py-3 text-[13px] text-(--nox-ink-soft)">
            Only DJ personas can create sets.
          </p>
        ) : (
          <div className="grid gap-4">
            <AuthField id="set-title" label="title" type="text" value={title} placeholder="Late night amapiano" autoComplete="off" icon={<span />} onChange={setTitle} />
            <AuthField id="playback-url" label="playback url" type="url" value={playbackURL} placeholder="https://cdn.example.com/set.mp4" autoComplete="off" icon={<span />} onChange={handlePlaybackURLChange} />
            <AuthField id="duration" label="duration (seconds)" type="number" value={duration} placeholder="900" autoComplete="off" icon={<span />} onChange={setDuration} />
            <AuthField id="thumbnail" label="thumbnail url" type="url" value={thumbnailURL} placeholder="https://cdn.example.com/thumb.jpg" autoComplete="off" icon={<span />} onChange={setThumbnailURL} />
            <AuthField id="genres" label="genres" type="text" value={genres} placeholder="amapiano, afro-house" autoComplete="off" icon={<span />} onChange={setGenres} />

            <label className="block" htmlFor="set-description">
              <span className="mb-2 block font-mono text-[10px] font-medium uppercase tracking-[0.18em] text-(--nox-ink-soft)">description</span>
              <textarea
                id="set-description"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                rows={4}
                className="w-full resize-none rounded-[8px] border border-(--nox-border) bg-(--nox-surface) px-3 py-3 text-[14px] text-(--nox-ink) outline-none transition focus:border-(--nox-accent-line)"
              />
            </label>
          </div>
        )}

        {message ? (
          <p className="mt-4 rounded-[8px] border border-(--nox-danger) bg-(--nox-danger-soft) px-3 py-2 text-[12px] font-medium text-(--nox-danger)">
            {message}
          </p>
        ) : null}

        {canCreate ? (
          <button
            type="submit"
            disabled={!canSubmit}
            className="mt-5 min-h-12 w-full rounded-[8px] bg-(--nox-accent) px-4 py-3 text-[15px] font-semibold text-white transition hover:brightness-110 disabled:cursor-not-allowed disabled:opacity-50"
          >
            {status === "loading" ? "Creating set..." : "Create set"}
          </button>
        ) : null}
      </form>

      <TabBar />
    </FeedShell>
  );
}
