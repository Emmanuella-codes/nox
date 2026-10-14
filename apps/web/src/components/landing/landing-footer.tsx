"use client";

const FOOTER_LINKS = [
  { label: "safety", href: "#safety" },
  { label: "for DJs and organizers", href: "/auth" },
  { label: "open nox", href: "/auth" },
];

export function LandingFooter() {
  return (
    <footer style={{ borderTop: "1px solid #1d1d1b" }}>
      <div
        className="landing-foot"
        style={{
          maxWidth: 1320,
          margin: "0 auto",
          padding: "28px clamp(20px,4vw,56px)",
          display: "flex",
          justifyContent: "space-between",
          gap: 20,
          flexWrap: "wrap",
          fontSize: 14,
          color: "#6f6f69",
        }}
      >
        <span>© 2026 nox · made in lagos</span>
        <div style={{ display: "flex", gap: 24, flexWrap: "wrap" }}>
          {FOOTER_LINKS.map(({ label, href }) => (
            <a key={label} href={href} style={{ color: "#6f6f69", textDecoration: "none" }}>
              {label}
            </a>
          ))}
        </div>
      </div>
    </footer>
  );
}
