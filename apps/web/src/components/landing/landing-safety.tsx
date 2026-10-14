"use client";

const ACCENT = "#16b8a4";

const ITEMS = [
  {
    t: "Report an event",
    d: "Flag unsafe venues, overcrowding, or bad security for review by the Nox team.",
  },
  {
    t: "Admission check-ins",
    d: "Use a QR credential or authorized staff search to check in at the event.",
  },
  {
    t: "Share location",
    d: "Opt in to temporary location updates visible to active crew members.",
  },
  {
    t: "Moderated anonymity",
    d: "Anonymous posts and comments stay separate from your public profile while remaining linked to your account for moderation.",
  },
];

export function LandingSafety() {
  return (
    <section id="safety" style={{ borderTop: "1px solid #1d1d1b" }}>
      <div
        style={{
          maxWidth: 1320,
          margin: "0 auto",
          padding: "96px clamp(20px,4vw,56px)",
          display: "flex",
          flexDirection: "column",
          gap: 48,
        }}
      >
        <div
          style={{
            display: "flex",
            justifyContent: "space-between",
            alignItems: "flex-end",
            gap: 24,
            flexWrap: "wrap",
          }}
        >
          <div style={{ display: "flex", flexDirection: "column", gap: 20, maxWidth: 760 }}>
            <span
              style={{
                fontFamily: "var(--font-jetbrains-mono), monospace",
                fontSize: 13,
                letterSpacing: ".12em",
                color: ACCENT,
              }}
            >
              SAFETY
            </span>
            <h2
              style={
                {
                  margin: 0,
                  fontSize: "clamp(36px,4.6vw,64px)",
                  lineHeight: 1,
                  letterSpacing: "-.035em",
                  fontWeight: 700,
                  textWrap: "balance",
                } as React.CSSProperties
              }
            >
              Built for the walk home, too.
            </h2>
          </div>
          <p style={{ margin: 0, maxWidth: 380, color: "#9a9a94", fontSize: 17, lineHeight: 1.55 }}>
            Tools that work at 5am on low battery, not just in the settings menu.
          </p>
        </div>
        <div
          style={{
            display: "grid",
            gridTemplateColumns: "repeat(auto-fit, minmax(min(100%, 260px), 1fr))",
            borderTop: "1px solid #232321",
          }}
        >
          {ITEMS.map((s) => (
            <div
              key={s.t}
              style={{
                padding: "28px 24px 8px 0",
                display: "flex",
                flexDirection: "column",
                gap: 10,
                borderBottom: "1px solid #232321",
              }}
            >
              <span style={{ fontSize: 22, fontWeight: 600, letterSpacing: "-.01em" }}>{s.t}</span>
              <span style={{ fontSize: 15, color: "#8a8a84", lineHeight: 1.5, paddingBottom: 20 }}>
                {s.d}
              </span>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
