"use client";

import Link from "next/link";
import { ChevronRight, Mail } from "lucide-react";
import { Mono } from "@/src/components/landing/landing-primitives";

export function LandingCta() {
  return (
    <section id="start" className="py-[100px]">
      <div className="mx-auto max-w-[1240px] px-8">
        <div
          className="landing-cta relative grid items-center gap-12 overflow-hidden rounded-lg border border-[var(--l-border-strong)] bg-[var(--l-surface)] px-14 py-16"
          style={{ gridTemplateColumns: "1.4fr 1fr" }}
        >
          <div className="relative">
            <Mono className="mb-3.5 block text-[var(--l-accent)]">
              Doors open
            </Mono>
            <h2 className="text-[clamp(32px,3.8vw,44px)] font-bold leading-[1.04] text-[var(--l-ink)]">
              Open nox.<br />The first night is on us.
            </h2>
            <p className="mt-4 max-w-[460px] text-[16px] leading-[1.55] text-[var(--l-ink-mid)]">
              nox runs in your browser — no app, no invites. Just a name (or not), a friend, and somewhere to be later.
            </p>
          </div>

          <div className="relative flex flex-col gap-2.5">
            <Link
              href="/auth"
              className="flex items-center justify-between rounded-lg border border-[var(--l-ink)] bg-[var(--l-ink)] px-[22px] py-[18px] text-[14.5px] font-semibold text-[var(--l-bg)] no-underline transition-colors"
            >
              <span>open nox in your browser</span>
              <ChevronRight size={16} strokeWidth={1.8} />
            </Link>
            <a
              href="mailto:scene@nox.app"
              className="flex items-center justify-between rounded-lg border border-[var(--l-border-strong)] bg-transparent px-[22px] py-[18px] text-[14.5px] font-semibold text-[var(--l-ink)] no-underline transition-colors"
            >
              <span>get the email digest</span>
              <Mail size={16} strokeWidth={1.5} />
            </a>
            <Mono className="mt-2 self-end">
              works in any browser · no install
            </Mono>
          </div>
        </div>
      </div>
    </section>
  );
}
