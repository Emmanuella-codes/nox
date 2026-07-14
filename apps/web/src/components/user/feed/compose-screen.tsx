"use client";

import { useEffect, useMemo, useRef, useState } from "react";
import { useRouter } from "next/navigation";
import { Ghost, User, Image as ImageIcon, Music, MapPin, X, Tag } from "lucide-react";
import { FeedShell } from "@/src/components/user/feed/feed-shell";
import { createPost } from "@/src/utils/api/user/post";
import { confirmPostMediaUpload, initiatePostMediaUpload, uploadToCloudinary } from "@/src/utils/api/user/media";
import { getMyPersonas } from "@/src/utils/api/user/persona";
import { getAccessToken, getActivePersonaID, setActivePersonaID } from "@/src/utils/auth/session";
import { ApiRequestError } from "@/src/utils/api/api";
import type { Persona } from "@/src/types/api/user/persona";

type PostingMode = "anonymous" | "public";

const MAX_CHARS = 280;

export function ComposeScreen() {
  const router = useRouter();
  const [mode, setMode] = useState<PostingMode>("anonymous");
  const [body, setBody] = useState("");
  const [posting, setPosting] = useState(false);
  const [personas, setPersonas] = useState<Persona[]>([]);
  const [personaID, setPersonaID] = useState("");
  const [mediaFile, setMediaFile] = useState<File | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [error, setError] = useState("");

  useEffect(() => {
    async function loadPersonas() {
      const token = getAccessToken();
      if (!token) return;
      try {
        const res = await getMyPersonas(token);
        const visiblePersonas = res.data ?? [];
        const activePersonaID = getActivePersonaID();
        const selectedPersona = visiblePersonas.find((persona) => persona.id === activePersonaID) ?? visiblePersonas[0];
        setPersonas(visiblePersonas);
        if (selectedPersona) {
          setActivePersonaID(selectedPersona.id);
        }
        setPersonaID(selectedPersona?.id ?? "");
      } catch {
        setPersonas([]);
      }
    }
    loadPersonas();
  }, []);

  const remaining = MAX_CHARS - body.length;
  const token = useMemo(() => getAccessToken(), []);
  const canPost =
    body.trim().length > 0 &&
    remaining >= 0 &&
    !posting &&
    Boolean(token) &&
    Boolean(personaID);

  const ringColor =
    remaining < 0 ? "var(--nox-danger)" : remaining < 30 ? "var(--nox-gold)" : "var(--nox-accent)";

  async function handlePost() {
    if (!canPost) return;
    setPosting(true);
    setError("");

    try {
      let mediaAssetID: string | undefined;
      let postType: "text" | "image" | "video" = "text";
      if (mediaFile) {
        const mediaKind = mediaFile.type.startsWith("video/") ? "video" : "image";
        const uploadInit = await initiatePostMediaUpload(
          {
            owner_persona_id: personaID,
            media_kind: mediaKind,
            mime_type: mediaFile.type,
            size_bytes: mediaFile.size,
          },
          token,
        );
        if (!uploadInit.data) throw new Error("upload_signature_missing");
        const uploaded = await uploadToCloudinary(mediaFile, uploadInit.data);
        const confirmed = await confirmPostMediaUpload(
          {
            owner_persona_id: personaID,
            media_kind: mediaKind,
            public_id: uploaded.public_id,
            secure_url: uploaded.secure_url,
            thumbnail_url: uploaded.thumbnail_url,
            mime_type: mediaFile.type,
            duration_seconds: Math.ceil(uploaded.duration ?? 0),
            size_bytes: uploaded.bytes || mediaFile.size,
          },
          token,
        );
        mediaAssetID = confirmed.data?.id;
        postType = mediaKind;
      }
      await createPost(
        {
          body: body.trim(),
          posting_mode: mode,
          post_type: postType,
          persona_id: personaID,
          ...(mediaAssetID ? { media_asset_ids: [mediaAssetID] } : {}),
        },
        token,
      );
      router.push("/feed");
    } catch (err) {
      setPosting(false);
      setError(err instanceof ApiRequestError ? err.message : "Could not create post.");
    }
  }

  return (
    <FeedShell>
      {/* Header */}
      <header className="flex items-center justify-between border-b border-(--nox-divider) px-4 py-3 pt-[env(safe-area-inset-top,12px)]">
        <button
          type="button"
          onClick={() => router.back()}
          className="flex size-8 items-center justify-center rounded-full transition hover:bg-(--nox-surface-alt)"
        >
          <X className="size-5 text-(--nox-ink)" strokeWidth={1.8} />
        </button>
        <span className="text-[15px] font-semibold text-(--nox-ink)">new post</span>
        <button
          type="button"
          disabled={!canPost}
          onClick={handlePost}
          className="rounded-full px-4 py-1.5 text-[13px] font-semibold text-white transition hover:brightness-110 disabled:cursor-not-allowed disabled:opacity-50"
          style={{ background: "var(--nox-accent)" }}
        >
          {posting ? "posting…" : "post"}
        </button>
      </header>

      <div className="flex flex-1 flex-col overflow-y-auto px-4 py-4">
        {/* Persona indicator */}
        <div className="mb-4 flex items-center gap-2">
          <button
            type="button"
            onClick={() => setMode("anonymous")}
            className="flex items-center gap-1.5 rounded-full border px-3 py-1.5 font-mono text-[10px] font-semibold uppercase tracking-[0.12em] transition"
            style={{
              borderColor: mode === "anonymous" ? "var(--nox-accent)" : "var(--nox-border-strong)",
              background: mode === "anonymous" ? "var(--nox-accent-soft)" : "transparent",
              color: mode === "anonymous" ? "var(--nox-accent)" : "var(--nox-ink-mid)",
            }}
          >
            <Ghost className="size-3" strokeWidth={1.8} />
            anonymous
          </button>
          <button
            type="button"
            onClick={() => setMode("public")}
            className="flex items-center gap-1.5 rounded-full border px-3 py-1.5 font-mono text-[10px] font-semibold uppercase tracking-[0.12em] transition"
            style={{
              borderColor: mode === "public" ? "var(--nox-accent)" : "var(--nox-border-strong)",
              background: mode === "public" ? "var(--nox-accent-soft)" : "transparent",
              color: mode === "public" ? "var(--nox-accent)" : "var(--nox-ink-mid)",
            }}
          >
            <User className="size-3" strokeWidth={1.8} />
            public
          </button>
        </div>

        {mode === "public" && (
          <label className="mb-4 block">
            <span className="mb-2 block font-mono text-[10px] font-semibold uppercase tracking-[0.14em] text-(--nox-ink-soft)">
              post as
            </span>
            <select
              value={personaID}
              onChange={(event) => setPersonaID(event.target.value)}
              className="w-full rounded-[8px] border border-(--nox-border) bg-(--nox-surface) px-3 py-2 text-[13px] text-(--nox-ink) outline-none"
            >
              {personas.length === 0 ? (
                <option value="">Create a public persona first</option>
              ) : (
                personas.map((persona) => (
                  <option key={persona.id} value={persona.id}>
                    {persona.display_name} (@{persona.handle})
                  </option>
                ))
              )}
            </select>
          </label>
        )}
        {mode === "anonymous" && personas.length > 1 ? (
          <label className="mb-4 block">
            <span className="mb-2 block font-mono text-[10px] font-semibold uppercase tracking-[0.14em] text-(--nox-ink-soft)">
              persona used privately
            </span>
            <select
              value={personaID}
              onChange={(event) => setPersonaID(event.target.value)}
              className="w-full rounded-[8px] border border-(--nox-border) bg-(--nox-surface) px-3 py-2 text-[13px] text-(--nox-ink) outline-none"
            >
              {personas.map((persona) => (
                <option key={persona.id} value={persona.id}>
                  {persona.display_name} (@{persona.handle})
                </option>
              ))}
            </select>
          </label>
        ) : null}

        {/* Text area */}
        <textarea
          autoFocus
          value={body}
          onChange={(e) => setBody(e.target.value)}
          placeholder="what's the vibe?"
          rows={7}
          className="w-full resize-none bg-transparent text-[16px] leading-[1.6] text-(--nox-ink) outline-none placeholder:text-(--nox-ink-soft)"
        />

        {error ? (
          <p className="mt-3 rounded-[8px] border border-(--nox-danger) bg-(--nox-danger-soft) px-3 py-2 text-[12px] font-medium text-(--nox-danger)">
            {error}
          </p>
        ) : null}
        {mediaFile ? (
          <div className="mt-3 flex items-center justify-between rounded-[8px] border border-(--nox-border) bg-(--nox-surface) px-3 py-2 text-[12px] text-(--nox-ink-mid)">
            <span className="truncate">{mediaFile.name}</span>
            <button type="button" onClick={() => setMediaFile(null)} className="text-(--nox-danger)">
              remove
            </button>
          </div>
        ) : null}

        {/* Attachments row */}
        <div className="mt-auto border-t border-(--nox-divider) pt-3">
          <div className="flex items-center gap-1">
            <button
              type="button"
              onClick={() => fileInputRef.current?.click()}
              className="flex size-9 items-center justify-center rounded-full transition hover:bg-(--nox-surface-alt) text-(--nox-ink-soft)"
            >
              <ImageIcon className="size-4" strokeWidth={1.7} />
            </button>
            <input
              ref={fileInputRef}
              type="file"
              accept="image/jpeg,image/png,image/webp,image/gif,video/mp4,video/webm,video/quicktime"
              className="hidden"
              onChange={(event) => setMediaFile(event.target.files?.[0] ?? null)}
            />
            <button
              type="button"
              className="flex size-9 items-center justify-center rounded-full transition hover:bg-(--nox-surface-alt) text-(--nox-ink-soft)"
            >
              <Music className="size-4" strokeWidth={1.7} />
            </button>
            <button
              type="button"
              className="flex size-9 items-center justify-center rounded-full transition hover:bg-(--nox-surface-alt) text-(--nox-ink-soft)"
            >
              <Tag className="size-4" strokeWidth={1.7} />
            </button>
            <button
              type="button"
              className="flex size-9 items-center justify-center rounded-full transition hover:bg-(--nox-surface-alt) text-(--nox-ink-soft)"
            >
              <MapPin className="size-4" strokeWidth={1.7} />
            </button>

            {/* Char counter */}
            <div className="ml-auto flex items-center gap-2">
              <span
                className="font-mono text-[11px] font-medium"
                style={{ color: ringColor }}
              >
                {remaining}
              </span>
              <svg className="size-5 -rotate-90" viewBox="0 0 20 20">
                <circle cx="10" cy="10" r="8" fill="none" stroke="var(--nox-border)" strokeWidth="2" />
                <circle
                  cx="10"
                  cy="10"
                  r="8"
                  fill="none"
                  stroke={ringColor}
                  strokeWidth="2"
                  strokeDasharray={`${Math.min(2 * Math.PI * 8, (Math.max(0, MAX_CHARS - body.length) / MAX_CHARS) * 2 * Math.PI * 8)} ${2 * Math.PI * 8}`}
                  strokeLinecap="round"
                />
              </svg>
            </div>
          </div>
        </div>
      </div>
    </FeedShell>
  );
}
