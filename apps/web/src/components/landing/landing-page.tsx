"use client";

import { LandingNav } from "@/src/components/landing/landing-nav";
import { LandingHero } from "@/src/components/landing/landing-hero";
import { LandingTwoLayers } from "@/src/components/landing/landing-features";
import { LandingAnonymity } from "@/src/components/landing/landing-scenes";
import { LandingEvents } from "@/src/components/landing/landing-manifesto";
import { LandingCrews } from "@/src/components/landing/landing-crews";
import { LandingSafety } from "@/src/components/landing/landing-safety";
import { LandingCta } from "@/src/components/landing/landing-cta";
import { LandingFooter } from "@/src/components/landing/landing-footer";
import { themeVars } from "@/src/components/landing/landing-theme";

export function LandingPage() {
  return (
    <div
      className="min-h-screen overflow-x-hidden antialiased"
      style={{
        ...themeVars(),
        background: "#0b0b0a",
        color: "#f2f1ec",
        fontFamily: "var(--font-space-grotesk), var(--font-inter), system-ui, sans-serif",
      }}
    >
      <LandingNav />
      <LandingHero />
      <LandingTwoLayers />
      <LandingAnonymity />
      <LandingEvents />
      <LandingCrews />
      <LandingSafety />
      <LandingCta />
      <LandingFooter />
    </div>
  );
}
