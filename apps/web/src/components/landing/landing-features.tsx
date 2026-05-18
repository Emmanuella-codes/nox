"use client";

import { Grid2X2, CalendarDays, AudioLines, Star, User, Hash } from "lucide-react";
import { SectionHead } from "@/src/components/landing/landing-primitives";
import type { LucideIcon } from "lucide-react";

interface Feature {
  icon: LucideIcon;
  title: string;
  body: string;
  tag: string;
}

const FEATURES: Feature[] = [
  {
    icon: Grid2X2,
    title: "An unfiltered feed",
    body: "Real posts from real people on the floor — the gist, the gossip, the recommendations. Nothing smoothed over, nothing sponsored.",
    tag: "the feed",
  },
  {
    icon: CalendarDays,
    title: "Tonight, not next month",
    body: "An events board built for the way Lagos actually decides — one nudge at 11pm. Filter by sound, by area, by who's already on the floor.",
    tag: "events",
  },
  {
    icon: AudioLines,
    title: "Sets you can scrub",
    body: "Recorded sets attach to the night they were played. Tap to scrub the waveform, save the cue, pull the tracklist.",
    tag: "set archive",
  },
  {
    icon: Star,
    title: "Venue reputation, honestly",
    body: "Door policy, sound quality, women's review, security — rated by the people who actually got in last weekend.",
    tag: "venues",
  },
  {
    icon: User,
    title: "Post your way",
    body: "Some takes you want signed. Some you don't. Choose to post as yourself or anonymously, per post — not a default someone else chose for you.",
    tag: "profile",
  },
  {
    icon: Hash,
    title: "Scene archive",
    body: "Hashtags as living tags — #amapiano-vi, #afro-house-thursdays, #onikan-after-hours. Every scene gets a page, and a memory.",
    tag: "archive",
  },
];

function FeatureCard({ icon: Icon, title, body, tag }: Feature) {
  return (
    <div className="flex min-h-[240px] flex-col gap-3.5 bg-[var(--l-bg)] px-8 py-9">
      <div className="flex h-11 w-11 items-center justify-center rounded-lg border border-[var(--l-accent-line)] bg-[var(--l-accent-soft)] text-[var(--l-accent)]">
        <Icon size={20} strokeWidth={1.6} />
      </div>
      <h3 className="text-[20px] font-semibold text-[var(--l-ink)]">{title}</h3>
      <p className="text-[14.5px] leading-[1.55] text-[var(--l-ink-mid)]">{body}</p>
      <span className="mt-auto font-mono text-[10.5px] font-medium uppercase text-[var(--l-ink-soft)]">
        {tag}
      </span>
    </div>
  );
}

export function LandingFeatures() {
  return (
    <section id="what" className="relative py-[100px]">
      <div className="mx-auto max-w-[1240px] px-8">
        <SectionHead
          eyebrow="what nox does"
          title={<>Built for the people,<br />not the feed.</>}
          body="Every other social app smooths the city out. nox keeps the texture — the inside jokes, the slow nights, the friend-of-a-friend who runs the door, the rooms only forty people know about."
        />
        <div
          className="landing-features grid grid-cols-3 gap-px overflow-hidden rounded-lg border border-[var(--l-border)] bg-[var(--l-border)]"
        >
          {FEATURES.map((f) => (
            <FeatureCard key={f.tag} {...f} />
          ))}
        </div>
      </div>
    </section>
  );
}
