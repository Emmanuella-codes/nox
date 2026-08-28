"use client";

import { Mono } from "@/src/components/landing/landing-primitives";

export function LandingManifesto() {
  return (
    <section
      id="manifesto"
      className="border-y border-(--l-divider) py-30 text-center"
    >
      <div className="mx-auto max-w-310 px-8">
        <p className="mx-auto max-w-225 text-[clamp(32px,4.2vw,52px)] font-medium leading-[1.18] text-(--l-ink)">
          We didn&apos;t build a network. We built{" "}
          <span className="text-(--l-accent)">a noticeboard</span> for
          the city you go out in — honest, late, and a little messy on purpose.
        </p>
        <Mono className="mt-8 inline-block">
          — nox, founders&apos; note
        </Mono>
      </div>
    </section>
  );
}
