"use client";

import { SectionHead } from "@/src/components/landing/landing-primitives";

const SCENES = [
  { tag: "#amapiano", title: "Yaba\nbasement", sub: "FRI · 134 there now", color: "#0f766e" },
  { tag: "#afro-house", title: "VI\nrooftop set", sub: "SAT · 280 saved", color: "#b45309" },
  { tag: "#mara", title: "Onikan\nafter hours", sub: "SUN · doors 02:00", color: "#525252" },
  { tag: "#gqom", title: "Lekki\nclub night", sub: "THU · monthly", color: "#2563eb" },
];

export function LandingScenes() {
  return (
    <section id="scenes" className="pb-[100px] pt-10">
      <div className="mx-auto max-w-[1240px] px-8">
        <SectionHead
          eyebrow="scenes on nox"
          title={<>The city has more<br />than one sound.</>}
          body="Browse the rooms by what's actually playing in them this week."
        />
        <div className="landing-showcase grid grid-cols-4 gap-3.5">
          {SCENES.map(({ tag, title, sub, color }) => (
            <div
              key={tag}
              className="flex flex-col justify-between overflow-hidden rounded-lg p-[18px] text-white"
              style={{ aspectRatio: "9/16", background: color }}
            >
              <span className="self-start rounded font-mono text-[10px] uppercase px-2 py-1" style={{ background: "rgba(0,0,0,.32)" }}>
                {tag}
              </span>
              <div>
                <div className="whitespace-pre-line text-[18px] md:text-[22px] font-bold leading-[1.1]">
                  {title}
                </div>
                <div className="mt-2 font-mono text-[10px] opacity-80">
                  {sub}
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
