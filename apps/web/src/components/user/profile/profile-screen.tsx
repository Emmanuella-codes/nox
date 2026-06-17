"use client";

import { useEffect, useState } from "react";
import { Settings } from "lucide-react";
import { useRouter } from "next/navigation";
import { FeedShell } from "@/src/components/user/feed/feed-shell";
import { TabBar } from "@/src/components/user/feed/tab-bar";
import { PostCard } from "@/src/components/user/feed/post-card";
import { Avatar } from "@/src/components/user/shared/avatar";
import { getMyPersonas } from "@/src/utils/api/user/persona";
import { getPersonaPosts, likePost, unlikePost } from "@/src/utils/api/user/post";
import { getAccessToken, getActivePersonaID, setActivePersonaID } from "@/src/utils/auth/session";
import type { Persona } from "@/src/types/api/persona";
import type { Post } from "@/src/types/api/post";

type ProfileTab = "posts" | "sets";

function formatCount(n: number): string {
  if (n >= 1000) return `${(n / 1000).toFixed(1)}k`;
  return String(n);
}

export function ProfileScreen() {
  const router = useRouter();
  const [persona, setPersona] = useState<Persona | null>(null);
  const [posts, setPosts] = useState<Post[]>([]);
  const [tab, setTab] = useState<ProfileTab>("posts");
  const [loading, setLoading] = useState(true);
  const [message, setMessage] = useState("");

  useEffect(() => {
    async function load() {
      try {
        const token = getAccessToken();
        if (!token) {
          setMessage("Sign in to view your profile.");
          return;
        }
        const res = await getMyPersonas(token);
        const personas = res.data ?? [];
        const activePersonaID = getActivePersonaID();
        const selectedPersona = personas.find((item) => item.id === activePersonaID) ?? personas[0] ?? null;
        setPersona(selectedPersona);
        if (!selectedPersona) {
          setMessage("Create a public persona to build your profile.");
          return;
        }
        if (selectedPersona) {
          setActivePersonaID(selectedPersona.id);
          const postsRes = await getPersonaPosts(selectedPersona.id, token, selectedPersona.id);
          setPosts(postsRes.data ?? []);
        }
      } catch {
        setMessage("Could not load profile.");
      } finally {
        setLoading(false);
      }
    }
    load();
  }, []);

  const display: Persona = persona ?? {
    id: "profile",
    handle: "yourhandle",
    display_name: "Your public persona",
    bio: message,
    avatar_url: "",
    cover_url: "",
    persona_type: "visible",
    category: "patron",
    genre_tags: ["afro-house", "amapiano", "electronic"],
    follower_count: 0,
    following_count: 0,
    post_count: 0,
    created_at: "",
    updated_at: "",
  };

  async function handleToggleLike(post: Post) {
    const token = getAccessToken();
    if (!token || !persona) return;

    const previousPosts = posts;
    const nextLiked = !post.is_liked;
    setPosts((current) =>
      current.map((item) =>
        item.id === post.id
          ? {
              ...item,
              is_liked: nextLiked,
              like_count: Math.max(0, item.like_count + (nextLiked ? 1 : -1)),
            }
          : item,
      ),
    );

    try {
      if (nextLiked) {
        await likePost(post.id, persona.id, token);
      } else {
        await unlikePost(post.id, persona.id, token);
      }
    } catch {
      setPosts(previousPosts);
    }
  }

  return (
    <FeedShell>
      <div className="flex-1 overflow-y-auto">
        {/* Cover strip */}
        <div
          className="h-24 w-full"
          style={{
            background: `linear-gradient(135deg, #1a1028 0%, #0d0d14 100%)`,
          }}
        />

        {/* Profile header */}
        <div className="px-4">
          {/* Avatar row */}
          <div className="flex items-end justify-between" style={{ marginTop: -28 }}>
            <div className="rounded-[16px] border-2 overflow-hidden" style={{ borderColor: "var(--nox-bg)" }}>
              <Avatar id={display.id} name={display.display_name} size={56} square />
            </div>

            <div className="flex items-center gap-2 pb-1">
              <button
                type="button"
                onClick={() => router.push("/settings")}
                className="flex size-8 items-center justify-center rounded-full border border-(--nox-border-strong) transition hover:border-(--nox-accent)"
              >
                <Settings className="size-4 text-(--nox-ink-mid)" strokeWidth={1.7} />
              </button>
              <button
                type="button"
                onClick={() => router.push(persona ? "/settings" : "/onboarding")}
                className="rounded-[10px] border border-(--nox-border-strong) px-4 py-1.5 text-[13px] font-semibold text-(--nox-ink) transition hover:border-(--nox-accent) hover:text-(--nox-accent)"
              >
                {persona ? "edit profile" : "create persona"}
              </button>
            </div>
          </div>

          {/* Name & handle */}
          <div className="mt-3">
            <h1 className="text-[20px] font-bold tracking-[-0.03em] text-(--nox-ink)">
              {display.display_name}
            </h1>
            <p className="text-[13px] text-(--nox-ink-soft)">@{display.handle}</p>
          </div>

          {/* Bio */}
          {display.bio && (
            <p className="mt-2 text-[13px] leading-[1.55] text-(--nox-ink-mid)">{display.bio}</p>
          )}

          {/* Genre tags */}
          {display.genre_tags.length > 0 && (
            <div className="mt-2.5 flex flex-wrap gap-1.5">
              {display.genre_tags.map((tag) => (
                <span
                  key={tag}
                  className="rounded-full px-2.5 py-1 font-mono text-[10px] font-medium lowercase"
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

          {/* Stats */}
          <div className="mt-4 flex gap-5 border-b border-(--nox-divider) pb-4">
            <div>
              <p className="text-[16px] font-bold text-(--nox-ink)">{formatCount(display.post_count)}</p>
              <p className="font-mono text-[10px] text-(--nox-ink-soft)">posts</p>
            </div>
            {[
              { label: "followers", value: display.follower_count, href: `/personas/${display.id}/followers` },
              { label: "following", value: display.following_count, href: `/personas/${display.id}/following` },
            ].map(({ label, value, href }) => (
              <button
                key={label}
                type="button"
                onClick={() => persona && router.push(href)}
                disabled={!persona}
                className="text-left transition hover:opacity-70 disabled:pointer-events-none"
              >
                <p className="text-[16px] font-bold text-(--nox-ink)">{formatCount(value)}</p>
                <p className="font-mono text-[10px] text-(--nox-ink-soft)">{label}</p>
              </button>
            ))}
          </div>

          {/* Profile tabs */}
          <div className="flex gap-0 pt-1">
            {(["posts", "sets"] as ProfileTab[]).map((t) => {
              const active = tab === t;
              return (
                <button
                  key={t}
                  type="button"
                  onClick={() => setTab(t)}
                  className="relative pb-2.5 pr-6 font-mono text-[11px] font-medium tracking-[0.06em] transition"
                  style={{ color: active ? "var(--nox-accent)" : "var(--nox-ink-soft)" }}
                >
                  {t}
                  {active && (
                    <span
                      className="absolute bottom-0 left-0 right-6 h-[1.5px] rounded-full"
                      style={{ background: "var(--nox-accent)" }}
                    />
                  )}
                </button>
              );
            })}
          </div>
        </div>

        {/* Content */}
        {tab === "posts" && (
          <div className="mt-1">
            {loading ? (
              <div className="space-y-0">
                {[...Array(3)].map((_, i) => (
                  <div key={i} className="border-b border-(--nox-divider) px-4 py-4">
                    <div className="flex gap-3">
                      <div className="size-9 rounded-full bg-(--nox-surface-alt) animate-pulse" />
                      <div className="flex-1 space-y-2">
                        <div className="h-3 w-1/3 rounded bg-(--nox-surface-alt) animate-pulse" />
                        <div className="h-3 w-full rounded bg-(--nox-surface-alt) animate-pulse" />
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            ) : posts.length > 0 ? (
              posts.map((p) => (
                <PostCard
                  key={p.id}
                  post={p}
                  liked={p.is_liked}
                  onLike={handleToggleLike}
                  onClick={() => router.push(`/posts/${p.id}`)}
                />
              ))
            ) : (
              <p className="py-12 text-center text-[13px] text-(--nox-ink-soft)">no posts yet.</p>
            )}
          </div>
        )}

        {tab === "sets" && (
          <div className="flex flex-col items-center justify-center py-16 px-4">
            <p className="text-[13px] text-(--nox-ink-soft)">set archive coming soon.</p>
          </div>
        )}
      </div>

      <TabBar />
    </FeedShell>
  );
}
