"use client";

import Link from "next/link";

const ACCENT = "#16b8a4";

const NAV_LINKS = [
  { label: "feed", href: "#social" },
  { label: "events", href: "#events" },
  { label: "crews", href: "#crews" },
  { label: "safety", href: "#safety" },
];

export function LandingNav() {
  return (
    <header
      style={{
        position: "sticky",
        top: 0,
        zIndex: 20,
        background: "rgba(11,11,10,.82)",
        backdropFilter: "blur(14px)",
        borderBottom: "1px solid #1d1d1b",
      }}
    >
      <nav
        style={{
          maxWidth: 1320,
          margin: "0 auto",
          padding: "18px clamp(20px,4vw,56px)",
          display: "flex",
          alignItems: "center",
          justifyContent: "space-between",
          gap: 24,
        }}
      >
        <Link
          href="#top"
          style={{
            display: "flex",
            alignItems: "center",
            gap: 10,
            fontWeight: 700,
            fontSize: 22,
            letterSpacing: "-.03em",
            color: "#f2f1ec",
            textDecoration: "none",
          }}
        >
          <span
            style={{
              width: 30,
              height: 30,
              borderRadius: "50%",
              border: "1.5px solid #f2f1ec",
              display: "grid",
              placeItems: "center",
              fontSize: 15,
            }}
          >
            N
          </span>
          nox
        </Link>

        <div
          className="landing-nav-links"
          style={{ display: "flex", gap: 28, fontSize: 15, color: "#9a9a94" }}
        >
          {NAV_LINKS.map(({ label, href }) => (
            <a key={label} href={href} style={{ color: "#9a9a94", textDecoration: "none" }}>
              {label}
            </a>
          ))}
        </div>

        <Link
          href="/auth"
          style={{
            background: ACCENT,
            color: "#0b0b0a",
            fontWeight: 600,
            fontSize: 15,
            padding: "10px 18px",
            borderRadius: 10,
            textDecoration: "none",
          }}
        >
          open nox →
        </Link>
      </nav>
    </header>
  );
}
