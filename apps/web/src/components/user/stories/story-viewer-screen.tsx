"use client";

import { useCallback, useEffect, useState } from "react";
import { ChevronLeft, Plus, X } from "lucide-react";
import { useRouter } from "next/navigation";
import { Sheet, SheetContent } from "@/components/ui/sheet";
import { StoryItemPlayer } from "@/src/components/user/stories/story-item-player";
import { Avatar } from "@/src/components/user/shared/avatar";
import { useActivePersona } from "@/src/hooks/use-active-persona";
import { getStory, addStoryItem, deleteStoryItem } from "@/src/utils/api/user/story";
import { createSetVideoAsset } from "@/src/utils/api/user/set";
import { getAccessToken } from "@/src/utils/auth/session";
import { ApiRequestError } from "@/src/utils/api/api";
import type { Story } from "@/src/types/api/story";

interface StoryViewerScreenProps {
  storyID: string;
}

export function StoryViewerScreen({ storyID }: StoryViewerScreenProps) {
  const router = useRouter();
  const { activePersona, loading: personaLoading } = useActivePersona();
  const [story, setStory] = useState<Story | null>(null);
  const [loading, setLoading] = useState(true);
  const [currentIndex, setCurrentIndex] = useState(0);
  const [addOpen, setAddOpen] = useState(false);
  const [clipURL, setClipURL] = useState("");
  const [postingMode, setPostingMode] = useState<"public" | "anonymous">("public");
  const [addStatus, setAddStatus] = useState<"idle" | "loading" | "error">("idle");
  const [addMsg, setAddMsg] = useState("");

  const load = useCallback(async () => {
    const token = getAccessToken() ?? undefined;
    const viewerID = activePersona?.id;
    const res = await getStory(storyID, viewerID, token);
    const data = res.data;
    if (data) setStory(data);
    setLoading(false);
  }, [storyID, activePersona?.id]);

  useEffect(() => {
    if (personaLoading) return;
    void load();
  }, [personaLoading, load]);

  const items = story?.items.slice().sort((a, b) => a.position - b.position) ?? [];
  const total = items.length;

  function handleNext() {
    setCurrentIndex((i) => Math.min(i + 1, total - 1));
  }

  function handlePrev() {
    setCurrentIndex((i) => Math.max(i - 1, 0));
  }

  async function handleDeleteItem(itemID: string) {
    const token = getAccessToken();
    if (!token) return;
    await deleteStoryItem(storyID, itemID, token);
    void load();
  }

  async function handleAddClip() {
    const token = getAccessToken();
    if (!token || !activePersona) return;
    setAddStatus("loading");
    setAddMsg("");
    try {
      const mediaRes = await createSetVideoAsset(
        {
          owner_persona_id: activePersona.id,
          storage_key: new URL(clipURL.trim()).pathname.replace(/^\//, ""),
          playback_url: clipURL.trim(),
          mime_type: "video/mp4",
          duration_seconds: 60,
          size_bytes: 0,
        },
        token,
      );
      const mediaID = mediaRes.data?.id;
      if (!mediaID) throw new Error("media_missing");
      await addStoryItem(storyID, { contributor_persona_id: activePersona.id, media_asset_id: mediaID, posting_mode: postingMode }, token);
      setAddOpen(false);
      setClipURL("");
      setAddStatus("idle");
      void load();
    } catch (err) {
      setAddStatus("error");
      setAddMsg(err instanceof ApiRequestError ? err.message : "Could not add clip.");
    }
  }

  if (loading) {
    return (
      <div className="flex size-full items-center justify-center bg-black">
        <p className="text-[13px] text-white/50">Loading...</p>
      </div>
    );
  }

  if (!story) {
    return (
      <div className="flex size-full flex-col items-center justify-center gap-4 bg-black">
        <p className="text-[13px] text-white/50">Story not found.</p>
        <button type="button" onClick={() => router.back()} className="text-[13px] text-white/70 underline">
          Go back
        </button>
      </div>
    );
  }

  const currentItem = items[currentIndex];
  const isOwner = activePersona?.id === story.owner.id;

  return (
    <div className="relative size-full bg-black overflow-hidden">
      {/* Progress bar */}
      <div className="absolute inset-x-0 top-0 z-20 flex gap-1 px-3 pt-[env(safe-area-inset-top,12px)]">
        {items.map((_, i) => (
          <div key={i} className="h-[3px] flex-1 rounded-full overflow-hidden bg-white/25">
            <div
              className="h-full rounded-full bg-white transition-all duration-200"
              style={{ width: i < currentIndex ? "100%" : i === currentIndex ? "50%" : "0%" }}
            />
          </div>
        ))}
        {items.length === 0 && <div className="h-[3px] flex-1 rounded-full bg-white/25" />}
      </div>

      {/* Header */}
      <div className="absolute inset-x-0 top-[calc(env(safe-area-inset-top,12px)+12px)] z-20 flex items-center justify-between px-3 pt-3">
        <button type="button" onClick={() => router.back()} className="flex size-8 items-center justify-center">
          <ChevronLeft className="size-5 text-white" strokeWidth={1.8} />
        </button>
        <div className="flex items-center gap-2">
          <Avatar id={story.owner.id} name={story.owner.display_name} size={28} />
          <p className="text-[13px] font-semibold text-white">{story.title}</p>
        </div>
        {story.can_contribute && activePersona ? (
          <button type="button" onClick={() => setAddOpen(true)} className="flex size-8 items-center justify-center rounded-full bg-white/15 backdrop-blur">
            <Plus className="size-4 text-white" strokeWidth={1.8} />
          </button>
        ) : <div className="size-8" />}
      </div>

      {/* Video */}
      {currentItem ? (
        <StoryItemPlayer item={currentItem} active={true} onEnded={handleNext} />
      ) : (
        <div className="flex size-full items-center justify-center">
          <p className="text-[13px] text-white/50">No clips yet.</p>
        </div>
      )}

      {/* Touch zones */}
      <button type="button" onClick={handlePrev} className="absolute inset-y-0 left-0 w-1/3" aria-label="Previous" />
      <button type="button" onClick={handleNext} className="absolute inset-y-0 right-0 w-1/3" aria-label="Next" />

      {/* Owner: delete current item */}
      {isOwner && currentItem && (
        <button type="button" onClick={() => handleDeleteItem(currentItem.id)}
          className="absolute bottom-10 right-4 z-20 flex items-center gap-1.5 rounded-full bg-black/50 px-3 py-1.5 backdrop-blur">
          <X className="size-3 text-white/70" strokeWidth={1.8} />
          <span className="text-[11px] text-white/70">remove clip</span>
        </button>
      )}

      {/* Add clip sheet */}
      <Sheet open={addOpen} onOpenChange={setAddOpen}>
        <SheetContent side="bottom" className="bg-(--nox-surface) px-5 pb-8">
          <p className="mb-4 text-[16px] font-bold tracking-[-0.02em] text-(--nox-ink)">Add your clip</p>
          <input
            className="w-full rounded-[8px] border border-(--nox-border) bg-(--nox-surface) px-3 py-3 text-[14px] text-(--nox-ink) outline-none focus:border-(--nox-accent-line)"
            placeholder="Video URL"
            type="url"
            value={clipURL}
            onChange={(e) => setClipURL(e.target.value)}
          />
          <div className="mt-3 flex gap-2">
            {(["public", "anonymous"] as const).map((mode) => (
              <button key={mode} type="button" onClick={() => setPostingMode(mode)}
                className={`rounded-full border px-3 py-1.5 font-mono text-[10px] font-medium transition ${postingMode === mode ? "border-(--nox-accent) bg-(--nox-accent-soft) text-(--nox-accent-ink)" : "border-(--nox-border) text-(--nox-ink-soft)"}`}>
                {mode}
              </button>
            ))}
          </div>
          {addMsg && <p className="mt-3 text-[12px] text-(--nox-danger)">{addMsg}</p>}
          <button type="button" onClick={handleAddClip} disabled={!clipURL.trim() || addStatus === "loading"}
            className="mt-4 w-full rounded-[8px] bg-(--nox-accent) py-3 text-[14px] font-semibold text-white disabled:opacity-50">
            {addStatus === "loading" ? "Adding..." : "Add clip"}
          </button>
        </SheetContent>
      </Sheet>
    </div>
  );
}
