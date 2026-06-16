"use client";

import { useEffect, useState } from "react";
import { ChevronLeft } from "lucide-react";
import { useRouter } from "next/navigation";
import { FeedShell } from "@/src/components/user/feed/feed-shell";
import { TabBar } from "@/src/components/user/feed/tab-bar";
import { AuthField } from "@/src/components/user/auth/auth-field";
import { Avatar } from "@/src/components/user/shared/avatar";
import { useActivePersona } from "@/src/hooks/use-active-persona";
import { updatePersona } from "@/src/utils/api/user/persona";
import { getAccessToken } from "@/src/utils/auth/session";
import { ApiRequestError } from "@/src/utils/api/api";

export function EditProfileScreen() {
  const router = useRouter();
  const { activePersona, loading: personaLoading } = useActivePersona();
  const [displayName, setDisplayName] = useState("");
  const [bio, setBio] = useState("");
  const [avatarURL, setAvatarURL] = useState("");
  const [coverURL, setCoverURL] = useState("");
  const [genres, setGenres] = useState("");
  const [status, setStatus] = useState<"idle" | "loading" | "success" | "error">("idle");
  const [message, setMessage] = useState("");

  useEffect(() => {
    if (!activePersona) return;
    setDisplayName(activePersona.display_name);
    setBio(activePersona.bio ?? "");
    setAvatarURL(activePersona.avatar_url ?? "");
    setCoverURL(activePersona.cover_url ?? "");
    setGenres(activePersona.genre_tags?.join(", ") ?? "");
  }, [activePersona]);

  async function handleSave(e: React.FormEvent) {
    e.preventDefault();
    const token = getAccessToken();
    if (!token || !activePersona) return;
    setStatus("loading");
    setMessage("");
    try {
      await updatePersona(
        activePersona.id,
        {
          display_name: displayName.trim() || undefined,
          bio: bio.trim() || undefined,
          avatar_url: avatarURL.trim() || undefined,
          cover_url: coverURL.trim() || undefined,
          genre_tags: genres.split(",").map((g) => g.trim().toLowerCase()).filter(Boolean),
        },
        token,
      );
      setStatus("success");
      setMessage("Profile updated.");
    } catch (err) {
      setStatus("error");
      setMessage(err instanceof ApiRequestError ? err.message : "Could not save changes.");
    }
  }

  return (
    <FeedShell>
      <header className="flex items-center gap-3 border-b border-(--nox-divider) px-4 pt-[env(safe-area-inset-top,12px)] pb-3">
        <button type="button" onClick={() => router.back()} className="flex size-8 items-center justify-center rounded-full hover:bg-(--nox-surface-alt)">
          <ChevronLeft className="size-5 text-(--nox-ink)" strokeWidth={1.8} />
        </button>
        <h1 className="text-[18px] font-bold text-(--nox-ink)">edit profile</h1>
      </header>

      <form onSubmit={handleSave} className="flex-1 overflow-y-auto px-4 py-4">
        {personaLoading ? (
          <p className="text-[13px] text-(--nox-ink-soft)">Loading persona...</p>
        ) : !activePersona ? (
          <p className="text-[13px] text-(--nox-ink-soft)">No active persona.</p>
        ) : (
          <div className="grid gap-4">
            <div className="flex items-center gap-3">
              <Avatar id={activePersona.id} name={activePersona.display_name} size={52} square />
              <div>
                <p className="text-[14px] font-semibold text-(--nox-ink)">@{activePersona.handle}</p>
                <p className="font-mono text-[10px] text-(--nox-ink-soft)">{activePersona.category}</p>
              </div>
            </div>

            <AuthField id="ep-name" label="display name" type="text" value={displayName} placeholder="Your name" autoComplete="off" icon={<span />} onChange={setDisplayName} />
            <AuthField id="ep-avatar" label="avatar url" type="url" value={avatarURL} placeholder="https://..." autoComplete="off" icon={<span />} onChange={setAvatarURL} />
            <AuthField id="ep-cover" label="cover url" type="url" value={coverURL} placeholder="https://..." autoComplete="off" icon={<span />} onChange={setCoverURL} />
            <AuthField id="ep-genres" label="genres" type="text" value={genres} placeholder="amapiano, afro-house" autoComplete="off" icon={<span />} onChange={setGenres} />

            <label htmlFor="ep-bio" className="block">
              <span className="mb-2 block font-mono text-[10px] font-medium uppercase tracking-[0.18em] text-(--nox-ink-soft)">bio</span>
              <textarea id="ep-bio" value={bio} rows={3} onChange={(e) => setBio(e.target.value)}
                className="w-full resize-none rounded-[8px] border border-(--nox-border) bg-(--nox-surface) px-3 py-3 text-[14px] text-(--nox-ink) outline-none focus:border-(--nox-accent-line)" />
            </label>
          </div>
        )}

        {message && (
          <p className={`mt-4 rounded-[8px] border px-3 py-2 text-[12px] font-medium ${
            status === "success"
              ? "border-(--nox-success) bg-(--nox-success-soft) text-(--nox-success)"
              : "border-(--nox-danger) bg-(--nox-danger-soft) text-(--nox-danger)"
          }`}>{message}</p>
        )}

        {activePersona && (
          <button type="submit" disabled={status === "loading"}
            className="mt-5 w-full rounded-[8px] bg-(--nox-accent) py-3 text-[15px] font-semibold text-white disabled:opacity-50">
            {status === "loading" ? "Saving..." : "Save changes"}
          </button>
        )}
      </form>

      <TabBar />
    </FeedShell>
  );
}
