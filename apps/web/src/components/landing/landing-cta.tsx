"use client";

import Link from "next/link";

const ACCENT = "#16b8a4";

export function LandingCta() {
  return (
    <section id="join" style={{ borderTop: "1px solid #1d1d1b" }}>
      <div
        style={{
          maxWidth: 1320,
          margin: "0 auto",
          padding: "120px clamp(20px,4vw,56px)",
          display: "flex",
          flexDirection: "column",
          alignItems: "flex-start",
          gap: 36,
        }}
      >
        <h2
          style={
            {
              margin: 0,
              fontSize: "clamp(52px,9vw,148px)",
              lineHeight: 0.88,
              letterSpacing: "-.05em",
              fontWeight: 700,
            } as React.CSSProperties
          }
        >
          see you{" "}
          <span style={{ color: ACCENT }}>out there.</span>
        </h2>
        <div style={{ display: "flex", flexDirection: "column", gap: 12, alignItems: "flex-start" }}>
          <Link
            href="/auth"
            style={{
              cursor: "pointer",
              background: ACCENT,
              color: "#0b0b0a",
              fontFamily: "var(--font-space-grotesk), sans-serif",
              fontWeight: 600,
              fontSize: 17,
              padding: "17px 26px",
              borderRadius: 12,
              textDecoration: "none",
            }}
          >
            join the scene →
          </Link>
          <span style={{ color: "#8a8a84", fontSize: 15 }}>
            Create an account to explore posts, events, sets, messages, and crews.
          </span>
        </div>
      </div>
    </section>
  );
}
