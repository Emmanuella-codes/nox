"use client";

import { ChevronLeft, ImagePlus, Send, Trash2, X } from "lucide-react";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { useRouter } from "next/navigation";
import { FeedShell } from "@/src/components/user/feed/feed-shell";
import { Avatar } from "@/src/components/user/shared/avatar";
import type { Conversation, Message } from "@/src/types/api/user/messaging";
import { confirmPostMediaUpload, initiatePostMediaUpload, uploadToCloudinary } from "@/src/utils/api/user/media";
import { deleteMessage, getConversation, listMessages, markConversationRead, sendMessage } from "@/src/utils/api/user/messaging";
import { getAccessToken } from "@/src/utils/auth/session";
import { formatDateTime } from "@/src/utils/format/date";
import { conversationHandle, conversationName, otherMemberID } from "@/src/utils/messaging/display";
import { useActivePersona } from "@/src/hooks/use-active-persona";

interface ThreadScreenProps {
  conversationID: string;
}

const PAGE_SIZE = 30;

export function ThreadScreen({ conversationID }: ThreadScreenProps) {
  const router = useRouter();
  const { activeID, loading: personaLoading } = useActivePersona();
  const [input, setInput] = useState("");
  const [conversation, setConversation] = useState<Conversation | null>(null);
  const [messages, setMessages] = useState<Message[]>([]);
  const [mediaFile, setMediaFile] = useState<File | null>(null);
  const [loading, setLoading] = useState(true);
  const [loadingOlder, setLoadingOlder] = useState(false);
  const [hasOlder, setHasOlder] = useState(false);
  const [sending, setSending] = useState(false);
  const [error, setError] = useState("");
  const fileInputRef = useRef<HTMLInputElement>(null);
  const token = useMemo(() => getAccessToken(), []);

  const loadThread = useCallback(async (silent: boolean) => {
    if (!silent) setLoading(true);
    setError("");
    try {
      const [conversationRes, messagesRes] = await Promise.all([
        getConversation(conversationID, token),
        listMessages(conversationID, activeID, token, PAGE_SIZE),
      ]);
      const loadedConversation = conversationRes.data ?? null;
      const loadedMessages = [...(messagesRes.data ?? [])].reverse();
      setConversation(loadedConversation);
      if (silent) {
        setMessages((current) => mergeMessages(current, loadedMessages));
      } else {
        setMessages(loadedMessages);
        setHasOlder((messagesRes.data ?? []).length === PAGE_SIZE);
      }
      const last = loadedMessages.at(-1);
      if (last) void markConversationRead(conversationID, { persona_id: activeID, message_id: last.id }, token);
    } catch {
      if (!silent) setError("Could not load conversation.");
    } finally {
      if (!silent) setLoading(false);
    }
  }, [activeID, conversationID, token]);

  useEffect(() => {
    if (personaLoading) return;
    if (!token) {
      router.replace("/auth");
      return;
    }
    if (!activeID) {
      void Promise.resolve().then(() => {
        setError("Choose a persona before opening messages.");
        setLoading(false);
      });
      return;
    }
    void loadThread(false);
  }, [activeID, conversationID, loadThread, personaLoading, router, token]);

  useEffect(() => {
    if (personaLoading || !token || !activeID) return;
    const interval = window.setInterval(() => {
      void loadThread(true);
    }, 5000);
    return () => window.clearInterval(interval);
  }, [activeID, loadThread, personaLoading, token]);

  async function handleSend() {
    const text = input.trim();
    if ((!text && !mediaFile) || !activeID || !token || sending) return;
    setSending(true);
    setError("");
    try {
      let mediaAssetID: string | undefined;
      let messageType: "text" | "image" | "video" = "text";
      if (mediaFile) {
        const mediaKind = mediaFile.type.startsWith("video/") ? "video" : "image";
        const uploadInit = await initiatePostMediaUpload(
          {
            owner_persona_id: activeID,
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
            owner_persona_id: activeID,
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
        messageType = mediaKind;
      }
      const res = await sendMessage(
        conversationID,
        {
          sender_persona_id: activeID,
          body: text,
          message_type: messageType,
          ...(mediaAssetID ? { media_asset_id: mediaAssetID } : {}),
        },
        token,
      );
      if (res.data) setMessages((current) => [...current, res.data as Message]);
      setInput("");
      setMediaFile(null);
    } catch {
      setError("Could not send message.");
    } finally {
      setSending(false);
    }
  }

  async function handleLoadOlder() {
    if (!activeID || !token || loadingOlder) return;
    setLoadingOlder(true);
    try {
      const res = await listMessages(conversationID, activeID, token, PAGE_SIZE, messages.length);
      const older = [...(res.data ?? [])].reverse();
      setMessages((current) => [...older, ...current]);
      setHasOlder((res.data ?? []).length === PAGE_SIZE);
    } catch {
      setError("Could not load older messages.");
    } finally {
      setLoadingOlder(false);
    }
  }

  async function handleDelete(messageID: string) {
    if (!token) return;
    try {
      const res = await deleteMessage(messageID, token);
      if (res.data) {
        setMessages((current) => current.map((message) => (message.id === messageID ? res.data as Message : message)));
      }
    } catch {
      setError("Could not delete message.");
    }
  }

  const title = conversation ? conversationName(conversation, activeID) : "message";
  const subtitle = conversation ? conversationHandle(conversation, activeID) : "conversation";
  const avatarID = conversation ? otherMemberID(conversation, activeID) || conversation.id : conversationID;

  return (
    <FeedShell>
      <header className="flex items-center gap-3 border-b border-(--nox-divider) px-4 pt-[env(safe-area-inset-top,12px)] pb-3">
        <button type="button" onClick={() => router.back()} className="flex size-8 items-center justify-center rounded-full hover:bg-(--nox-surface-alt)">
          <ChevronLeft className="size-5 text-(--nox-ink)" strokeWidth={1.8} />
        </button>
        <Avatar id={avatarID} name={title} size={32} />
        <div className="min-w-0">
          <p className="truncate text-[14px] font-bold text-(--nox-ink)">{title}</p>
          <p className="truncate font-mono text-[10px] text-(--nox-ink-soft)">{subtitle}</p>
        </div>
      </header>

      <div className="flex-1 space-y-3 overflow-y-auto px-4 py-4">
        {loading && <p className="py-8 text-center text-[12px] text-(--nox-ink-soft)">Loading conversation...</p>}
        {error && <p className="py-3 text-center text-[12px] text-(--nox-danger)">{error}</p>}
        {!loading && messages.length === 0 && (
          <p className="py-8 text-center text-[12px] text-(--nox-ink-soft)">Start the conversation.</p>
        )}
        {!loading && hasOlder && (
          <button
            type="button"
            onClick={() => void handleLoadOlder()}
            disabled={loadingOlder}
            className="mx-auto block rounded-full border border-(--nox-border) px-3 py-1.5 text-[11px] text-(--nox-ink-soft) disabled:opacity-50"
          >
            {loadingOlder ? "Loading..." : "Load older"}
          </button>
        )}
        {messages.map((message) => {
          const mine = message.sender_persona_id === activeID;
          return (
            <div key={message.id} className={`flex ${mine ? "justify-end" : "justify-start"}`}>
              <div className={`max-w-[78%] rounded-[8px] px-3.5 py-2.5 ${mine ? "bg-(--nox-accent) text-white" : "bg-(--nox-surface-alt) text-(--nox-ink)"}`}>
                {message.media && !message.deleted && (
                  <div className="mb-2 overflow-hidden rounded-[8px] bg-black/20">
                    {message.media.media_kind === "video" ? (
                      <video className="max-h-[260px] w-full object-cover" src={message.media.playback_url} controls playsInline />
                    ) : (
                      // eslint-disable-next-line @next/next/no-img-element
                      <img className="max-h-[260px] w-full object-cover" src={message.media.playback_url} alt="" />
                    )}
                  </div>
                )}
                {(message.body || message.deleted) && (
                  <p className="text-[13px] leading-snug">{message.deleted ? "Message deleted" : message.body}</p>
                )}
                <p className={`mt-1 font-mono text-[9px] ${mine ? "text-white/60" : "text-(--nox-ink-faint)"}`}>
                  {formatDateTime(message.created_at)}
                </p>
                {mine && !message.deleted && (
                  <button
                    type="button"
                    onClick={() => void handleDelete(message.id)}
                    className="mt-1 inline-flex items-center gap-1 font-mono text-[9px] text-white/60"
                  >
                    <Trash2 className="size-3" strokeWidth={1.8} />
                    delete
                  </button>
                )}
              </div>
            </div>
          );
        })}
      </div>

      <div className="border-t border-(--nox-divider) px-4 py-3 pb-[max(12px,env(safe-area-inset-bottom))]">
        {mediaFile && (
          <div className="mb-2 flex items-center justify-between rounded-[8px] border border-(--nox-border) px-3 py-2 text-[12px] text-(--nox-ink-soft)">
            <span className="truncate">{mediaFile.name}</span>
            <button type="button" onClick={() => setMediaFile(null)} className="ml-3 text-(--nox-ink)">
              <X className="size-4" strokeWidth={1.8} />
            </button>
          </div>
        )}
        <div className="flex items-center gap-2">
          <input
            ref={fileInputRef}
            type="file"
            accept="image/jpeg,image/png,image/webp,image/gif,video/mp4,video/webm,video/quicktime"
            className="hidden"
            onChange={(event) => setMediaFile(event.target.files?.[0] ?? null)}
          />
          <button
            type="button"
            onClick={() => fileInputRef.current?.click()}
            className="flex size-10 items-center justify-center rounded-full border border-(--nox-border) text-(--nox-ink-soft)"
          >
            <ImagePlus className="size-4" strokeWidth={1.8} />
          </button>
          <input
            value={input}
            onChange={(event) => setInput(event.target.value)}
            onKeyDown={(event) => event.key === "Enter" && void handleSend()}
            placeholder="Message..."
            className="flex-1 rounded-full border border-(--nox-border) bg-(--nox-surface) px-4 py-2.5 text-[14px] text-(--nox-ink) outline-none focus:border-(--nox-accent-line)"
          />
          <button
            type="button"
            onClick={() => void handleSend()}
            disabled={(!input.trim() && !mediaFile) || sending}
            className="flex size-10 items-center justify-center rounded-full bg-(--nox-accent) text-white disabled:opacity-40"
          >
            <Send className="size-4" strokeWidth={1.8} />
          </button>
        </div>
      </div>
    </FeedShell>
  );
}

function mergeMessages(current: Message[], incoming: Message[]) {
  const byID = new Map<string, Message>();
  current.forEach((message) => byID.set(message.id, message));
  incoming.forEach((message) => byID.set(message.id, message));
  return [...byID.values()].sort((a, b) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime());
}
