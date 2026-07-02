"use client";

import { Edit } from "lucide-react";
import { useEffect, useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import { FeedShell } from "@/src/components/user/feed/feed-shell";
import { TabBar } from "@/src/components/user/feed/tab-bar";
import { Avatar } from "@/src/components/user/shared/avatar";
import type { Conversation } from "@/src/types/api/messaging";
import { listConversations } from "@/src/utils/api/user/messaging";
import { getAccessToken } from "@/src/utils/auth/session";
import { formatDateTime } from "@/src/utils/format/date";
import { conversationHandle, conversationName, lastMessagePreview } from "@/src/utils/messaging/display";
import { useActivePersona } from "@/src/hooks/use-active-persona";

export function MessagesScreen() {
  const router = useRouter();
	const { activeID, loading: personaLoading } = useActivePersona();
	const [conversations, setConversations] = useState<Conversation[]>([]);
	const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const token = useMemo(() => getAccessToken(), []);

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
    listConversations(activeID, token)
      .then((res) => setConversations(res.data ?? []))
      .catch(() => setError("Could not load conversations."))
      .finally(() => setLoading(false));
  }, [activeID, personaLoading, router, token]);

  return (
    <FeedShell>
      <header className="flex items-center justify-between px-4 pt-[env(safe-area-inset-top,12px)] pb-3">
        <h1 className="text-[22px] font-bold tracking-[-0.03em] text-(--nox-ink)">messages</h1>
        <button
          type="button"
          onClick={() => router.push("/messages/new")}
          className="flex size-8 items-center justify-center rounded-full border border-(--nox-border-strong) hover:border-(--nox-accent)"
        >
          <Edit className="size-4 text-(--nox-ink-mid)" strokeWidth={1.7} />
        </button>
      </header>

      <div className="flex-1 overflow-y-auto">
        {loading && <p className="px-4 py-8 text-center text-[12px] text-(--nox-ink-soft)">Loading messages...</p>}
        {error && <p className="px-4 py-8 text-center text-[12px] text-(--nox-danger)">{error}</p>}
        {!loading && !error && conversations.length === 0 && (
          <p className="px-4 py-8 text-center text-[12px] text-(--nox-ink-soft)">No messages yet.</p>
        )}
		{conversations.map((conversation) => {
			const name = conversationName(conversation, activeID);
          return (
            <button
              key={conversation.id}
              type="button"
              onClick={() => router.push(`/messages/${conversation.id}`)}
              className="flex w-full items-center gap-3 border-b border-(--nox-divider) px-4 py-4 text-left transition hover:bg-(--nox-surface)"
            >
              <Avatar id={conversation.id} name={name} size={44} />
              <div className="min-w-0 flex-1">
                <p className="truncate text-[14px] font-semibold text-(--nox-ink)">{name}</p>
                <p className="truncate text-[12px] text-(--nox-ink-soft)">{lastMessagePreview(conversation)}</p>
                <p className="truncate font-mono text-[10px] text-(--nox-ink-faint)">
									{conversationHandle(conversation, activeID)}
                </p>
              </div>
              <div className="flex shrink-0 flex-col items-end gap-1">
                <span className="font-mono text-[10px] text-(--nox-ink-faint)">{formatDateTime(conversation.updated_at)}</span>
                {conversation.unread_count > 0 && (
                  <span className="flex size-5 items-center justify-center rounded-full bg-(--nox-accent) font-mono text-[10px] font-semibold text-white">
                    {conversation.unread_count}
                  </span>
                )}
              </div>
            </button>
          );
        })}
      </div>

      <TabBar />
    </FeedShell>
  );
}
