"use client";

const ACCENT = "#16b8a4";

const ALL_EVENTS = [
  { name: "Sunrise Sessions", venue: "Rooftop House, Lekki", time: "FRI · 11PM–6AM", price: "from ₦15,000", tag: "AFRO HOUSE" },
  { name: "Basement 004", venue: "Undisclosed, VI", time: "FRI · 1AM", price: "from ₦10,000", tag: "TECHNO" },
  { name: "Aniko All Night", venue: "The Warehouse, Ikoyi", time: "SAT · 10PM", price: "from ₦20,000", tag: "EXTENDED SET" },
];

export function LandingEvents() {
  const events = ALL_EVENTS;

  return (
    <section id="events" style={{ borderTop: "1px solid #1d1d1b" }}>
      <div
        style={{
          maxWidth: 1320,
          margin: "0 auto",
          padding: "96px clamp(20px,4vw,56px)",
          display: "flex",
          flexDirection: "column",
          gap: 40,
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
          <div style={{ display: "flex", flexDirection: "column", gap: 20 }}>
            <span
              style={{
                fontFamily: "var(--font-jetbrains-mono), monospace",
                fontSize: 13,
                letterSpacing: ".12em",
                color: ACCENT,
              }}
            >
              SAMPLE EVENTS + TICKETING
            </span>
            <h2
              style={
                {
                  margin: 0,
                  fontSize: "clamp(36px,4.6vw,64px)",
                  lineHeight: 1,
                  letterSpacing: "-.035em",
                  fontWeight: 700,
                } as React.CSSProperties
              }
            >
              What&apos;s on.
            </h2>
          </div>
          <span style={{ maxWidth: 300, color: "#8a8a84", fontSize: 15, lineHeight: 1.45 }}>
            Examples of published events and ticket details.
          </span>
        </div>
        <div
          style={{
            display: "grid",
            gridTemplateColumns: "repeat(auto-fill, minmax(min(100%, 280px), 1fr))",
            gap: 16,
          }}
        >
          {events.map((e) => (
            <article
              key={e.name}
              style={{
                background: "#121211",
                border: "1px solid #232321",
                borderRadius: 20,
                overflow: "hidden",
                display: "flex",
                flexDirection: "column",
              }}
            >
              <div
                style={{
                  aspectRatio: "4/5",
                  background:
                    "repeating-linear-gradient(135deg,#181816 0 10px,#1e1e1c 10px 20px)",
                  display: "flex",
                  alignItems: "flex-end",
                  justifyContent: "space-between",
                  padding: 14,
                  fontFamily: "var(--font-jetbrains-mono), monospace",
                  fontSize: 11,
                  color: "#6f6f69",
                }}
              >
                <span>sample event</span>
                <span
                  style={{
                    color: "#0b0b0a",
                    background: "#f2f1ec",
                    padding: "4px 8px",
                    borderRadius: 6,
                  }}
                >
                  {e.tag}
                </span>
              </div>
              <div
                style={{
                  padding: 18,
                  display: "flex",
                  flexDirection: "column",
                  gap: 6,
                  flex: 1,
                }}
              >
                <span
                  style={{
                    fontFamily: "var(--font-jetbrains-mono), monospace",
                    fontSize: 12,
                    color: ACCENT,
                  }}
                >
                  {e.time}
                </span>
                <span style={{ fontSize: 20, fontWeight: 600, letterSpacing: "-.01em" }}>
                  {e.name}
                </span>
                <span style={{ fontSize: 14, color: "#8a8a84" }}>{e.venue}</span>
                <div
                  style={{
                    display: "flex",
                    justifyContent: "space-between",
                    alignItems: "center",
                    marginTop: "auto",
                    paddingTop: 14,
                  }}
                >
                <span style={{ fontWeight: 600 }}>{e.price}</span>
                  <span
                    style={{
                      border: "1px solid #2c2c29",
                      padding: "8px 14px",
                      borderRadius: 9,
                      fontSize: 14,
                      cursor: "pointer",
                    }}
                  >
                    ticket details
                  </span>
                </div>
              </div>
            </article>
          ))}
        </div>
      </div>
    </section>
  );
}
