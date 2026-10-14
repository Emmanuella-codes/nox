"use client";

import { useState } from "react";

const ACCENT = "#16b8a4";

export function LandingAnonymity() {
  const [anon, setAnon] = useState(false);

  const composer = anon
    ? { name: "anonymous", ini: "··", color: "#3a3a36", trackBg: ACCENT, knob: 23 }
    : { name: "@tobe_tobey", ini: "TA", color: "#0f7a6c", trackBg: "#2c2c29", knob: 3 };

  return (
    <section style={{ borderTop: "1px solid #1d1d1b" }}>
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
            PER-POST ANONYMITY
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
            One account. Post as you — or as no one.
          </h2>
          <p style={{ margin: 0, maxWidth: 480, color: "#9a9a94", fontSize: 18, lineHeight: 1.55 }}>
            Anonymity is a switch on a post or comment, not a burner account. Say the honest thing
            about the sound system or the door policy; your profile stays out of it. Accountability
            stays in — abuse can still be reviewed against the real account.
          </p>
        </div>

        <div
          style={{
            background: "#121211",
            border: "1px solid #232321",
            borderRadius: 24,
            padding: 24,
            display: "flex",
            flexDirection: "column",
            gap: 18,
          }}
        >
          <div style={{ display: "flex", alignItems: "center", gap: 12 }}>
            <span
              style={{
                width: 42,
                height: 42,
                borderRadius: "50%",
                background: composer.color,
                display: "grid",
                placeItems: "center",
                fontWeight: 700,
                fontSize: 14,
                transition: "background .3s",
              }}
            >
              {composer.ini}
            </span>
            <div style={{ display: "flex", flexDirection: "column" }}>
              <span
                style={{
                  fontSize: 12,
                  color: "#6f6f69",
                  fontFamily: "var(--font-jetbrains-mono), monospace",
                }}
              >
                POSTING AS
              </span>
              <span style={{ fontWeight: 600, fontSize: 17 }}>{composer.name}</span>
            </div>
          </div>

          <div
            style={{
              minHeight: 110,
              background: "#0e0e0d",
              border: "1px solid #232321",
              borderRadius: 14,
              padding: 16,
              fontSize: 18,
              lineHeight: 1.45,
            }}
          >
            the queue at the door was 90 minutes and the bar ran out of ice by 1am. good music
            though.
          </div>

          <div
            style={{
              display: "flex",
              justifyContent: "space-between",
              alignItems: "center",
              gap: 16,
              flexWrap: "wrap",
            }}
          >
            <button
              onClick={() => setAnon((a) => !a)}
              style={{
                display: "flex",
                alignItems: "center",
                gap: 12,
                background: "transparent",
                border: 0,
                color: "#f2f1ec",
                cursor: "pointer",
                fontFamily: "var(--font-space-grotesk), sans-serif",
                fontSize: 15,
                padding: 0,
              }}
            >
              <span
                style={{
                  width: 48,
                  height: 28,
                  borderRadius: 20,
                  background: composer.trackBg,
                  position: "relative",
                  transition: "background .25s",
                  display: "inline-block",
                  flexShrink: 0,
                }}
              >
                <span
                  style={{
                    position: "absolute",
                    top: 3,
                    left: composer.knob,
                    width: 22,
                    height: 22,
                    borderRadius: "50%",
                    background: "#f2f1ec",
                    transition: "left .25s",
                  }}
                />
              </span>
              post anonymously
            </button>
            <span
              style={{
                background: ACCENT,
                color: "#0b0b0a",
                fontWeight: 600,
                padding: "11px 20px",
                borderRadius: 10,
              }}
            >
              post
            </span>
          </div>
        </div>
      </div>
    </section>
  );
}
