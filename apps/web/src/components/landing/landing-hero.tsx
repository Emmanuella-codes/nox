"use client";

import { useState } from "react";
import Link from "next/link";

const ACCENT = "#16b8a4";

type PhoneTab = "feed" | "events" | "crew";

const POSTS = [
  { ini: "TA", who: "@tobe_tobey", meta: "lekki · 14m", color: "#0f7a6c", text: "sunrise playing till sunrise at rooftop house. come through", likes: 142, comments: 23 },
  { ini: "··", who: "anonymous", meta: "lekki · 2h", color: "#3a3a36", text: "would definitely not be going for any event at green palace ever. terrible everything", likes: 112, comments: 4 },
  { ini: "KA", who: "kennysaysso", meta: "ikoyi · 1h", color: "#6b4fc4", text: "trust aniko to always mother", likes: 88, comments: 14 },
  { ini: "SW", who: "@princess_azula", meta: "ikeja · 2d", color: "#2f5fe0", text: "still thinking about that b2b on sunday", likes: 61, comments: 9 },
];

const PHONE_EVENTS = [
  { name: "Sunrise Sessions", venue: "Rooftop House, Lekki", time: "FRI · 11PM–6AM", price: "from ₦15,000" },
  { name: "Basement 004", venue: "Undisclosed, VI", time: "FRI · 1AM", price: "from ₦10,000" },
  { name: "Amapiano on the Sand", venue: "Landmark Beach", time: "SAT · 4PM", price: "from ₦8,000" },
  { name: "Aniko All Night", venue: "The Warehouse, Ikoyi", time: "SAT · 10PM", price: "from ₦20,000" },
];

const CREW_MSGS = [
  { who: "dami", text: "at the gate, left side by the palm trees", mine: false },
  { who: "you", text: "omw, 5 mins", mine: true },
  { who: "ife", text: "guys there's traffic at the entrance", mine: false },
  { who: "dami", text: "@ife i was in that traffic, i had to walk 😩", mine: false },
  { who: "you", text: "who's playing now?", mine: true },
];

function PhoneScreen({ tab, onTabChange }: { tab: PhoneTab; onTabChange: (t: PhoneTab) => void }) {
  return (
    <div
      style={{
        position: "relative",
        height: "100%",
        borderRadius: 48,
        background: "#0a0a09",
        overflow: "hidden",
        padding: "14px 14px 0",
        display: "flex",
        flexDirection: "column",
        boxShadow: "inset 0 0 0 1px #000",
      }}
    >
      {/* Dynamic island */}
      <div
        style={{
          position: "absolute",
          top: 11,
          left: "50%",
          transform: "translateX(-50%)",
          width: "31%",
          height: 30,
          borderRadius: 20,
          background: "#000",
          zIndex: 2,
        }}
      >
        <span
          style={{
            position: "absolute",
            right: 12,
            top: 10,
            width: 10,
            height: 10,
            borderRadius: "50%",
            background: "#11161c",
            boxShadow: "inset 0 0 0 2px #0a0d10",
          }}
        />
      </div>

      {/* Status bar */}
      <div
        style={{
          display: "flex",
          justifyContent: "space-between",
          alignItems: "center",
          padding: "4px 18px 0 24px",
          height: 34,
          fontSize: 15,
          fontWeight: 600,
        }}
      >
        <span>9:41</span>
        <span style={{ display: "flex", alignItems: "center", gap: 6 }}>
          <span style={{ display: "flex", alignItems: "flex-end", gap: 2, height: 11 }}>
            {[4, 6, 8, 11].map((h) => (
              <span key={h} style={{ width: 3, height: h, background: "#f2f1ec", borderRadius: 1 }} />
            ))}
          </span>
          <span
            style={{
              position: "relative",
              width: 24,
              height: 12,
              border: "1.2px solid rgba(242,241,236,.5)",
              borderRadius: 4,
              padding: "1.5px",
            }}
          >
            <span
              style={{ display: "block", width: "62%", height: "100%", background: "#f2f1ec", borderRadius: 1.5 }}
            />
            <span
              style={{
                position: "absolute",
                right: -4,
                top: 3,
                width: 2,
                height: 4,
                background: "rgba(242,241,236,.5)",
                borderRadius: "0 1px 1px 0",
              }}
            />
          </span>
        </span>
      </div>

      {/* Tab bar */}
      <div
        style={{
          display: "flex",
          gap: 6,
          margin: "14px 4px 12px",
          background: "#191917",
          padding: 4,
          borderRadius: 12,
        }}
      >
        {(["feed", "events", "crew"] as PhoneTab[]).map((t) => (
          <button
            key={t}
            onClick={() => onTabChange(t)}
            style={{
              flex: 1,
              border: 0,
              cursor: "pointer",
              fontFamily: "var(--font-space-grotesk), sans-serif",
              fontSize: 14,
              fontWeight: 600,
              padding: "9px 0",
              borderRadius: 9,
              background: tab === t ? "#f2f1ec" : "transparent",
              color: tab === t ? "#0b0b0a" : "#8a8a84",
            }}
          >
            {t}
          </button>
        ))}
      </div>

      {/* Content */}
      <div
        style={{
          flex: 1,
          minHeight: 0,
          paddingBottom: 24,
          overflow: "hidden",
          display: "flex",
          flexDirection: "column",
          gap: 10,
          WebkitMaskImage: "linear-gradient(#000 85%, transparent)",
          maskImage: "linear-gradient(#000 85%, transparent)",
        }}
      >
        {tab === "feed" &&
          POSTS.map((p, i) => (
            <div
              key={i}
              style={{
                background: "#181816",
                border: "1px solid #242422",
                borderRadius: 16,
                padding: "14px 16px",
                display: "flex",
                flexDirection: "column",
                gap: 10,
              }}
            >
              <div style={{ display: "flex", gap: 10, alignItems: "center" }}>
                <span
                  style={{
                    width: 34,
                    height: 34,
                    borderRadius: "50%",
                    background: p.color,
                    display: "grid",
                    placeItems: "center",
                    fontSize: 12,
                    fontWeight: 700,
                    flexShrink: 0,
                  }}
                >
                  {p.ini}
                </span>
                <div style={{ display: "flex", flexDirection: "column" }}>
                  <span style={{ fontSize: 14, fontWeight: 600 }}>{p.who}</span>
                  <span
                    style={{
                      fontFamily: "var(--font-jetbrains-mono), monospace",
                      fontSize: 11,
                      color: "#6f6f69",
                    }}
                  >
                    {p.meta}
                  </span>
                </div>
              </div>
              <div style={{ fontSize: 15, lineHeight: 1.4 }}>{p.text}</div>
              <div style={{ display: "flex", gap: 16, fontSize: 12, color: "#8a8a84" }}>
                <span>♡ {p.likes}</span>
                <span>○ {p.comments}</span>
              </div>
            </div>
          ))}

        {tab === "events" &&
          PHONE_EVENTS.map((e, i) => (
            <div
              key={i}
              style={{
                display: "flex",
                gap: 12,
                background: "#181816",
                border: "1px solid #242422",
                borderRadius: 16,
                padding: 10,
              }}
            >
              <div
                style={{
                  width: 64,
                  height: 64,
                  flexShrink: 0,
                  borderRadius: 10,
                  background: "repeating-linear-gradient(135deg,#1f1f1d 0 6px,#262624 6px 12px)",
                }}
              />
              <div style={{ display: "flex", flexDirection: "column", gap: 3, minWidth: 0 }}>
                <span
                  style={{
                    fontFamily: "var(--font-jetbrains-mono), monospace",
                    fontSize: 11,
                    color: ACCENT,
                  }}
                >
                  {e.time}
                </span>
                <span style={{ fontSize: 15, fontWeight: 600 }}>{e.name}</span>
                <span style={{ fontSize: 12, color: "#8a8a84" }}>
                  {e.venue} · {e.price}
                </span>
              </div>
            </div>
          ))}

        {tab === "crew" && (
          <>
            <div
              style={{
                display: "flex",
                justifyContent: "space-between",
                alignItems: "center",
                padding: "4px 6px 8px",
                borderBottom: "1px solid #242422",
              }}
            >
              <div style={{ display: "flex", flexDirection: "column" }}>
                <span style={{ fontWeight: 600 }}>landmark crew · 6</span>
                <span
                  style={{
                    fontFamily: "var(--font-jetbrains-mono), monospace",
                    fontSize: 11,
                    color: "#8a8a84",
                  }}
                >
                  expires 11:00 sat
                </span>
              </div>
              <span
                style={{
                  fontFamily: "var(--font-jetbrains-mono), monospace",
                  fontSize: 11,
                  color: ACCENT,
                  border: "1px solid #1f3b37",
                  padding: "4px 8px",
                  borderRadius: 6,
                }}
              >
                LIVE
              </span>
            </div>
            {CREW_MSGS.map((m, i) => (
              <div
                key={i}
                style={{
                  alignSelf: m.mine ? "flex-end" : "flex-start",
                  maxWidth: "80%",
                  background: m.mine ? ACCENT : "#1e1e1c",
                  color: m.mine ? "#0b0b0a" : "#f2f1ec",
                  padding: "9px 13px",
                  borderRadius: 14,
                  fontSize: 14,
                  lineHeight: 1.35,
                }}
              >
                <div style={{ fontSize: 11, opacity: 0.7, marginBottom: 2 }}>{m.who}</div>
                {m.text}
              </div>
            ))}
          </>
        )}
      </div>

      {/* Home indicator */}
      <div
        style={{
          position: "absolute",
          bottom: 8,
          left: "50%",
          transform: "translateX(-50%)",
          width: "36%",
          height: 5,
          borderRadius: 3,
          background: "#f2f1ec",
          opacity: 0.9,
        }}
      />
    </div>
  );
}

function PhoneMockup({ tab, onTabChange }: { tab: PhoneTab; onTabChange: (t: PhoneTab) => void }) {
  return (
    <div className="landing-phone-wrap" style={{ display: "flex", justifyContent: "center" }}>
      <div
        style={{
          position: "relative",
          width: "min(100%, 370px)",
          aspectRatio: "9 / 19.3",
          transform: "rotate(3deg)",
          transition: "transform .5s",
          borderRadius: 58,
          padding: 11,
          background:
            "linear-gradient(145deg,#4a4a47 0%,#1c1c1b 18%,#2e2e2c 50%,#151514 82%,#3d3d3a 100%)",
          boxShadow:
            "inset 0 0 0 1px rgba(255,255,255,.12), inset 0 0 0 3px #0c0c0b, 0 50px 120px -30px rgba(0,0,0,.9), 0 30px 140px -40px #16b8a4",
        }}
      >
        {/* Side buttons */}
        <span style={{ position: "absolute", left: -3, top: "17%", width: 4, height: 30, borderRadius: "2px 0 0 2px", background: "#2e2e2c" }} />
        <span style={{ position: "absolute", left: -3, top: "25%", width: 4, height: 56, borderRadius: "2px 0 0 2px", background: "#2e2e2c" }} />
        <span style={{ position: "absolute", left: -3, top: "34%", width: 4, height: 56, borderRadius: "2px 0 0 2px", background: "#2e2e2c" }} />
        <span style={{ position: "absolute", right: -3, top: "28%", width: 4, height: 86, borderRadius: "0 2px 2px 0", background: "#2e2e2c" }} />
        <PhoneScreen tab={tab} onTabChange={onTabChange} />
      </div>
    </div>
  );
}

export function LandingHero() {
  const [tab, setTab] = useState<PhoneTab>("feed");

  return (
    <section
      id="top"
      style={{
        maxWidth: 1320,
        margin: "0 auto",
        padding: "clamp(48px,8vw,104px) clamp(20px,4vw,56px) 72px",
        display: "grid",
        gridTemplateColumns: "repeat(auto-fit, minmax(min(100%, 460px), 1fr))",
        gap: 56,
        alignItems: "center",
      }}
    >
      <div style={{ display: "flex", flexDirection: "column", gap: 32, minWidth: 0 }}>
        <h1
          style={
            {
              margin: 0,
              fontSize: "clamp(56px,8.4vw,128px)",
              lineHeight: 0.9,
              letterSpacing: "-.045em",
              fontWeight: 700,
              textWrap: "balance",
            } as React.CSSProperties
          }
        >
          the city after midnight,{" "}
          <span style={{ color: ACCENT }}>unfiltered.</span>
        </h1>
        <p
          style={{
            margin: 0,
            maxWidth: 540,
            fontSize: "clamp(18px,1.5vw,21px)",
            lineHeight: 1.55,
            color: "#a3a39c",
          }}
        >
          The feed, the flyers, the tickets and the group chat — in one place, run by the people who
          actually go out in Lagos.
        </p>
        <div style={{ display: "flex", gap: 12, flexWrap: "wrap" }}>
          <Link
            href="/auth"
            style={{
              background: ACCENT,
              color: "#0b0b0a",
              fontWeight: 600,
              fontSize: 17,
              padding: "17px 26px",
              borderRadius: 12,
              textDecoration: "none",
            }}
          >
            open nox →
          </Link>
          <a
            href="#how"
            style={{
              border: "1px solid #2c2c29",
              fontWeight: 500,
              fontSize: 17,
              padding: "17px 26px",
              borderRadius: 12,
              color: "#f2f1ec",
              textDecoration: "none",
            }}
          >
            see how it works
          </a>
        </div>
      </div>

      <PhoneMockup tab={tab} onTabChange={setTab} />
    </section>
  );
}
