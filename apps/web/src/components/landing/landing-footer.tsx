"use client";

import { BrandDot } from "@/src/components/landing/landing-primitives";

interface FooterLink {
  label: string;
  href?: string;
}

const FOOTER_LINKS: { head: string; links: FooterLink[] }[] = [
  {
    head: "Company",
    links: [
      { label: "About", href: "#manifesto" },
      { label: "Contact", href: "mailto:scene@nox.app" },
    ],
  },
  {
    head: "The fine print",
    links: [
      { label: "Privacy & Terms" },
      { label: "Safety" },
      { label: "Community" },
    ],
  },
];

export function LandingFooter() {
  return (
    <footer className="border-t border-(--l-divider) pb-10 pt-[60px]">
      <div className="mx-auto max-w-[1240px] px-8">
        <div
          className="landing-foot grid items-start gap-10"
          style={{ gridTemplateColumns: "minmax(0,2fr) repeat(2,minmax(140px,1fr))" }}
        >
          <div>
            <div className="flex items-center gap-2.5 text-[22px] font-bold text-(--l-ink)">
              <BrandDot />
              <span>nox</span>
            </div>
            <p className="mt-3.5 max-w-[320px] text-[14px] leading-[1.55] text-(--l-ink-mid)">
              A noticeboard for the city after midnight. Built in Lagos. Run by people who go out.
            </p>
          </div>

          {FOOTER_LINKS.map(({ head, links }) => (
            <div key={head}>
              <h4 className="mb-4 font-mono text-[10px] uppercase text-(--l-ink-soft)">
                {head}
              </h4>
              <div className="grid gap-1.5">
                {links.map((link) =>
                  link.href ? (
                    <a
                      key={link.label}
                      href={link.href}
                      className="text-[14px] leading-normal text-(--l-ink-mid) no-underline"
                    >
                      {link.label}
                    </a>
                  ) : (
                    <span
                      key={link.label}
                      className="text-[14px] leading-normal text-(--l-ink-soft)"
                    >
                      {link.label}
                    </span>
                  ),
                )}
              </div>
            </div>
          ))}
        </div>

        <div className="landing-foot-bottom mt-14 flex justify-between gap-4 border-t border-(--l-divider) pt-6 font-mono text-[10.5px] uppercase text-(--l-ink-soft)">
          <span>© 2026 nox</span>
          <span>Lagos, NG</span>
        </div>
      </div>
    </footer>
  );
}
