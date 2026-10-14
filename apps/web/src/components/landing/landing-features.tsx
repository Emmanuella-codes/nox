"use client";

const ACCENT = "#16b8a4";

const SOCIAL_FEATURES = [
  { t: "Posts", d: "Text-first. Say it plainly." },
  { t: "Stories", d: "Photos and videos that expire per item after 24 hours." },
  { t: "Sets", d: "Share DJ sets with genre tags and audience engagement." },
  { t: "Hashtags", d: "Find posts by venue, genre, or scene." },
  { t: "Messages", d: "DMs for the after-party plan." },
  { t: "Profiles", d: "Your visible profile, posts, sets, and public activity." },
];

const EVENT_FEATURES = [
  { t: "Discovery", d: "Browse published events and upcoming nights out." },
  { t: "Ticketing", d: "Buy in-app, scan at the door." },
  { t: "Crews", d: "Temporary groups per event." },
  { t: "Safety", d: "Reports, admission checks, and temporary location sharing." },
];

function FeatureGrid({ items }: { items: { t: string; d: string }[] }) {
  return (
    <div
      style={{
        display: "grid",
        gridTemplateColumns: "repeat(2, minmax(0, 1fr))",
        gap: 1,
        background: "#232321",
        border: "1px solid #232321",
        borderRadius: 14,
        overflow: "hidden",
      }}
    >
      {items.map((f) => (
        <div
          key={f.t}
          style={{
            background: "#121211",
            padding: 18,
            display: "flex",
            flexDirection: "column",
            gap: 4,
          }}
        >
          <span style={{ fontWeight: 600, fontSize: 16 }}>{f.t}</span>
          <span style={{ color: "#8a8a84", fontSize: 14, lineHeight: 1.45 }}>{f.d}</span>
        </div>
      ))}
    </div>
  );
}

export function LandingTwoLayers() {
  return (
    <section id="how" style={{ borderTop: "1px solid #1d1d1b" }}>
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
          <h2
            style={
              {
                margin: 0,
                fontSize: "clamp(36px,4.6vw,64px)",
                lineHeight: 1,
                letterSpacing: "-.035em",
                fontWeight: 700,
                maxWidth: 760,
                textWrap: "balance",
              } as React.CSSProperties
            }
          >
            Two halves of one night out.
          </h2>
          <p style={{ margin: 0, maxWidth: 380, color: "#9a9a94", fontSize: 17, lineHeight: 1.55 }}>
            Talk about the scene. Then go be in it. Nox keeps both on the same account.
          </p>
        </div>
        <div
          style={{
            display: "grid",
            gridTemplateColumns: "repeat(auto-fit, minmax(min(100%, 420px), 1fr))",
            gap: 16,
          }}
        >
          <div
            id="social"
            style={{
              background: "#121211",
              border: "1px solid #232321",
              borderRadius: 24,
              padding: "clamp(28px,3vw,44px)",
              display: "flex",
              flexDirection: "column",
              gap: 28,
            }}
          >
            <span
              style={{
                fontFamily: "var(--font-jetbrains-mono), monospace",
                fontSize: 13,
                letterSpacing: ".12em",
                color: ACCENT,
              }}
            >
              01 · THE SOCIAL LAYER
            </span>
            <h3
              style={{
                margin: 0,
                fontSize: "clamp(28px,2.6vw,38px)",
                lineHeight: 1.05,
                letterSpacing: "-.025em",
                fontWeight: 600,
              }}
            >
              A text-first feed for what happened, what&apos;s happening, and who played.
            </h3>
            <FeatureGrid items={SOCIAL_FEATURES} />
          </div>
          <div
            style={{
              background: "#121211",
              border: "1px solid #232321",
              borderRadius: 24,
              padding: "clamp(28px,3vw,44px)",
              display: "flex",
              flexDirection: "column",
              gap: 28,
            }}
          >
            <span
              style={{
                fontFamily: "var(--font-jetbrains-mono), monospace",
                fontSize: 13,
                letterSpacing: ".12em",
                color: ACCENT,
              }}
            >
              02 · THE EVENT LAYER
            </span>
            <h3
              style={{
                margin: 0,
                fontSize: "clamp(28px,2.6vw,38px)",
                lineHeight: 1.05,
                letterSpacing: "-.025em",
                fontWeight: 600,
              }}
            >
              Find the party, buy the ticket, move with your people, and check in at the door.
            </h3>
            <FeatureGrid items={EVENT_FEATURES} />
          </div>
        </div>
      </div>
    </section>
  );
}
