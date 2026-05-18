"use client";

import Link from "next/link";
import { ArrowRight } from "lucide-react";
import { Mono } from "@/src/components/landing/landing-primitives";
import type { Theme } from "@/src/components/landing/landing-theme";

interface LandingHeroProps {
  theme: Theme;
}

function MockPost({
  initials, color, handle, location, body, likes, comments, extra,
}: {
  initials: string; color: string; handle: string; location: string;
  body: string; likes: string; comments: string; extra?: string;
}) {
  return (
    <div className="mb-2.5 rounded-lg border border-(--l-border) bg-(--l-surface) p-3">
      <div className="mb-2.5 flex items-center gap-2">
        <div
          className="flex h-7 w-7 shrink-0 items-center justify-center rounded-full text-[11px] font-bold text-white"
          style={{ background: color }}
        >
          {initials}
        </div>
        <div className="flex-1">
          <div className="text-[11.5px] font-semibold">{handle}</div>
          <div className="font-mono text-[9px] text-(--l-ink-soft)">{location}</div>
        </div>
      </div>
      <p className="mb-2.5 text-[12.5px] leading-[1.45]">{body}</p>
      <div className="flex gap-3.5 text-[10.5px] font-medium text-(--l-ink-mid)">
        <span>{likes}</span>
        <span>{comments}</span>
        {extra && <span>{extra}</span>}
      </div>
    </div>
  );
}

function PhoneMockup() {
  return (
    <div className="landing-phone-wrap flex justify-center">
      <div
        className="relative h-[620px] w-[300px] overflow-hidden rounded-lg rotate-2 bg-(--l-bg) text-(--l-ink)"
        style={{
          border: "8px solid #0a0a0c",
          boxShadow:
            "0 30px 60px -20px rgba(0,0,0,.5), 0 0 0 1px var(--l-border-strong), 0 0 80px -10px var(--l-accent-soft)",
        }}
      >
        {/* Notch */}
        <div
          className="absolute left-1/2 top-2 z-3 h-[22px] w-[90px] -translate-x-1/2 rounded-lg"
          style={{ background: "#0a0a0c" }}
        />
        <div className="flex h-full flex-col bg-(--l-bg) px-[18px] pt-10 text-[12px]">
          <div className="flex items-center justify-between px-1.5 pb-3 text-[11px] font-semibold">
            <span>9:41</span>
            <span className="opacity-50">●●●</span>
          </div>
          <div className="flex shrink-0 items-baseline justify-between px-1.5 pb-3.5">
            <span className="text-[22px] font-bold">feed</span>
            <Mono>FRI · 02:13</Mono>
          </div>
          <div className="no-scrollbar min-h-0 flex-1 overflow-y-auto px-0.5 pb-[18px]">
            <MockPost
              initials="TA" color="#0f766e" handle="@tobe_tobey" location="lekki · 14m ago"
              body="sunrise playing till sunrise at rooftop house ☺️" likes="♡ 142" comments="◌ 23" extra="↗ share"
            />
            <MockPost
              initials="AN" color="#b45309" handle="anonymous" location="lekki · 2h ago"
              body="would definitely not be going for any event at green palace ever.. terrible everything" likes="♡ 112" comments="◌ 4"
            />
            <MockPost
              initials="KA" color="#6b5ab2" handle="kennysaysso" location="ikoyi · 1h ago"
              body="trust aniko to always mother 😩" likes="♡ 88" comments="◌ 14"
            />
            <MockPost
              initials="SW" color="#2563eb" handle="@princess_azula" location="ikeja · 2d ago"
              body="still thinking about that yanfss b2b axara ☺️" likes="♡ 130" comments="◌ 13" extra="↗ share"
            />
            <MockPost
              initials="SW" color="#939c1c" handle="@princess_azula" location="ikeja · 7h ago"
              body="you guys are eating good in that lagos 😩 raving in an helipad" likes="♡ 130" comments="◌ 13" extra="↗ share"
            />
          </div>
        </div>
      </div>
    </div>
  );
}

export function LandingHero({ theme }: LandingHeroProps) {
  return (
    <header className="relative overflow-hidden py-20 pb-[100px]">
      <div
        className="landing-hero-grid relative mx-auto grid max-w-[1240px] items-center gap-[80px] px-8"
        style={{ gridTemplateColumns: "1.15fr .85fr" }}
      >
        {/* Left */}
        <div>
          <span className="inline-flex items-center gap-2.5 rounded-lg border border-(--l-accent-line) bg-(--l-accent-soft) px-3 py-[7px] font-mono text-[10.5px] font-medium uppercase text-(--l-accent-ink)">
            <span className="inline-block h-1.5 w-1.5 rounded-full bg-(--l-accent)" />
            Lagos · Now playing
          </span>

          <h1 className="mt-6 text-[clamp(48px,7.4vw,96px)] font-bold leading-[0.96] text-(--l-ink)">
            the city after<br />
            midnight,{" "}
            <em className="not-italic text-(--l-accent)">unfiltered.</em>
          </h1>

          <p className="mt-6 max-w-[480px] text-[17px] leading-[1.55] text-(--l-ink-mid)">
            nox is where Lagos nightlife actually talks to itself — the people, the rooms, the DJs, the friends-of-friends. A noticeboard for the city you go out in, run by the people who go out in it.
          </p>

          <div className="mt-9 flex flex-wrap gap-3">
            <Link
              href="/auth"
              className={`inline-flex items-center gap-2 rounded-lg border border-(--l-accent) bg-(--l-accent) px-[22px] py-[15px] text-[14.5px] font-semibold text-white no-underline transition-colors${theme === "dark" ? " shadow-[0_10px_24px_rgba(0,0,0,0.22)]" : ""}`}
            >
              open nox <ArrowRight size={14} strokeWidth={1.8} />
            </Link>
            <a
              href="#what"
              className="inline-flex items-center gap-2 rounded-lg border border-(--l-border-strong) bg-transparent px-[22px] py-[15px] text-[14.5px] font-semibold text-(--l-ink) no-underline transition-colors"
            >
              see how it works
            </a>
          </div>
        </div>

        <PhoneMockup />
      </div>
    </header>
  );
}
