"use client";

import { ChevronLeft, Send } from "lucide-react";
import { useState } from "react";
import { useRouter } from "next/navigation";
import { FeedShell } from "@/src/components/user/feed/feed-shell";
import { Avatar } from "@/src/components/user/shared/avatar";

interface ThreadScreenProps { conversationID: string }

const SAMPLE_MESSAGES = [
  { id: "m1", mine: false, text: "yo are you going to Warehouse Sessions?", time: "9:12 PM" },
  { id: "m2", mine: true, text: "yeah definitely. doors at 11?", time: "9:13 PM" },
  { id: "m3", mine: false, text: "midnight. DJ Khalid closes the night", time: "9:14 PM" },
  { id: "m4", mine: true, text: "perfect. see you on the floor 🔥", time: "9:14 PM" },
  { id: "m5", mine: false, text: "bring the energy 🎶", time: "9:15 PM" },
];

export function ThreadScreen({ conversationID }: ThreadScreenProps) {
  const router = useRouter();
  const [input, setInput] = useState("");
  const [messages, setMessages] = useState(SAMPLE_MESSAGES);

  function handleSend() {
    const text = input.trim();
    if (!text) return;
    setMessages((m) => [...m, { id: String(Date.now()), mine: true, text, time: "now" }]);
    setInput("");
  }

  return (
    <FeedShell>
      <header className="flex items-center gap-3 border-b border-(--nox-divider) px-4 pt-[env(safe-area-inset-top,12px)] pb-3">
        <button type="button" onClick={() => router.back()} className="flex size-8 items-center justify-center rounded-full hover:bg-(--nox-surface-alt)">
          <ChevronLeft className="size-5 text-(--nox-ink)" strokeWidth={1.8} />
        </button>
        <Avatar id={conversationID} name="Amirah Lagos" size={32} />
        <div>
          <p className="text-[14px] font-bold text-(--nox-ink)">Amirah Lagos</p>
          <p className="font-mono text-[10px] text-(--nox-ink-soft)">@amirah_lagos</p>
        </div>
      </header>

      <div className="flex-1 overflow-y-auto px-4 py-4 space-y-3">
        {messages.map((msg) => (
          <div key={msg.id} className={`flex ${msg.mine ? "justify-end" : "justify-start"}`}>
            <div className={`max-w-[78%] rounded-[14px] px-3.5 py-2.5 ${
              msg.mine
                ? "rounded-br-[4px] bg-(--nox-accent) text-white"
                : "rounded-bl-[4px] bg-(--nox-surface-alt) text-(--nox-ink)"
            }`}>
              <p className="text-[13px] leading-snug">{msg.text}</p>
              <p className={`mt-1 font-mono text-[9px] ${msg.mine ? "text-white/60" : "text-(--nox-ink-faint)"}`}>{msg.time}</p>
            </div>
          </div>
        ))}
      </div>

      <div className="border-t border-(--nox-divider) px-4 py-3 pb-[max(12px,env(safe-area-inset-bottom))]">
        <div className="flex items-center gap-2">
          <input value={input} onChange={(e) => setInput(e.target.value)}
            onKeyDown={(e) => e.key === "Enter" && handleSend()}
            placeholder="Message..."
            className="flex-1 rounded-full border border-(--nox-border) bg-(--nox-surface) px-4 py-2.5 text-[14px] text-(--nox-ink) outline-none focus:border-(--nox-accent-line)"
          />
          <button type="button" onClick={handleSend} disabled={!input.trim()}
            className="flex size-10 items-center justify-center rounded-full bg-(--nox-accent) text-white disabled:opacity-40">
            <Send className="size-4" strokeWidth={1.8} />
          </button>
        </div>
        <p className="mt-2 text-center text-[10px] text-(--nox-ink-faint)">Real-time messaging coming soon.</p>
      </div>
    </FeedShell>
  );
}
