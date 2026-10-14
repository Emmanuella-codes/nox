"use client";

const ACCENT = "#16b8a4";

const STEPS = [
  { n: "01", t: "Make a crew", d: "Create a temporary crew for an event and share the join code." },
  { n: "02", t: "Coordinate", d: "Share meet points, rides, and temporary location updates when you choose." },
  { n: "03", t: "Check in", d: "Use your event credential at the door and keep your crew updated in chat." },
  { n: "04", t: "It expires", d: "The crew closes after the event window. Clean slate." },
];

export function LandingCrews() {
  return (
    <section id="crews" style={{ borderTop: "1px solid #1d1d1b", background: "#0e0e0d" }}>
      <div
        style={{
          maxWidth: 1320,
          margin: "0 auto",
          padding: "96px clamp(20px,4vw,56px)",
          display: "grid",
          gridTemplateColumns: "repeat(auto-fit, minmax(min(100%, 440px), 1fr))",
          gap: 56,
          alignItems: "center",
        }}
      >
        <div style={{ display: "flex", flexDirection: "column", gap: 24 }}>
          <span
            style={{
              fontFamily: "var(--font-jetbrains-mono), monospace",
              fontSize: 13,
              letterSpacing: ".12em",
              color: ACCENT,
            }}
          >
            EVENT CREWS
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
            A group chat that ends when the night does.
          </h2>
          <p style={{ margin: 0, maxWidth: 480, color: "#9a9a94", fontSize: 18, lineHeight: 1.55 }}>
            Spin up a crew for one event. Share the meet point, coordinate rides, and find each
            other by the speakers. When the event window ends, the crew closes — no forever-chat
            with people you met at 4am.
          </p>
        </div>
        <div
          style={{ display: "grid", gridTemplateColumns: "repeat(2, minmax(0, 1fr))", gap: 12 }}
        >
          {STEPS.map((s) => (
            <div
              key={s.n}
              style={{
                background: "#141413",
                border: "1px solid #232321",
                borderRadius: 18,
                padding: 22,
                display: "flex",
                flexDirection: "column",
                gap: 28,
              }}
            >
              <span
                style={{
                  fontFamily: "var(--font-jetbrains-mono), monospace",
                  fontSize: 13,
                  color: "#6f6f69",
                }}
              >
                {s.n}
              </span>
              <div style={{ display: "flex", flexDirection: "column", gap: 6 }}>
                <span style={{ fontSize: 19, fontWeight: 600 }}>{s.t}</span>
                <span style={{ fontSize: 14, color: "#8a8a84", lineHeight: 1.45 }}>{s.d}</span>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
