"use client";

import { useEffect, useState } from "react";
import { LandingNav } from "@/src/components/landing/landing-nav";
import { LandingHero } from "@/src/components/landing/landing-hero";
import { LandingFeatures } from "@/src/components/landing/landing-features";
import { LandingScenes } from "@/src/components/landing/landing-scenes";
import { LandingManifesto } from "@/src/components/landing/landing-manifesto";
import { LandingCta } from "@/src/components/landing/landing-cta";
import { LandingFooter } from "@/src/components/landing/landing-footer";
import { themeVars, type Theme } from "@/src/components/landing/landing-theme";

export function LandingPage() {
  const [theme, setTheme] = useState<Theme>("dark");

  useEffect(() => {
    try {
      const stored = localStorage.getItem("nox-landing-theme");
      if (stored === "light" || stored === "dark") {
        queueMicrotask(() => setTheme(stored));
      }
    } catch {
      /* ignore */
    }
  }, []);

  function toggleTheme() {
    const next: Theme = theme === "dark" ? "light" : "dark";
    setTheme(next);
    try {
      localStorage.setItem("nox-landing-theme", next);
    } catch {
      /* ignore */
    }
  }

  return (
    <div
      className="min-h-screen overflow-x-hidden bg-[var(--l-bg)] text-[var(--l-ink)] antialiased"
      style={{
        ...themeVars(theme),
        fontFamily: "var(--font-space-grotesk), var(--font-inter), system-ui, sans-serif",
        transition: "background-color .35s ease, color .35s ease",
      }}
    >
      <LandingNav theme={theme} onToggleTheme={toggleTheme} />
      <LandingHero theme={theme} />
      <LandingFeatures />
      <LandingScenes />
      <LandingManifesto />
      <LandingCta />
      <LandingFooter />
    </div>
  );
}
