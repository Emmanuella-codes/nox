import type { CSSProperties, ReactNode } from "react";

export function Mono({ children, className = "" }: { children: ReactNode; className?: string }) {
  return (
    <span className={`font-mono text-[10.5px] font-medium uppercase tracking-[0.18em] text-[var(--l-ink-soft)] ${className}`}>
      {children}
    </span>
  );
}

export function BrandDot() {
  return (
    <span
      className="inline-block size-2.5 rounded-full bg-[var(--l-accent)]"
      style={{ boxShadow: "0 0 12px var(--l-accent)" }}
    />
  );
}

export function SectionHead({ eyebrow, title, body }: { eyebrow: string; title: ReactNode; body: string }) {
  return (
    <div className="mb-14 max-w-[720px]">
      <Mono className="mb-3.5 block text-[var(--l-accent)]">{eyebrow}</Mono>
      <h2 className="text-[clamp(34px,4.4vw,54px)] font-bold leading-[1.04] tracking-tight text-[var(--l-ink)]">
        {title}
      </h2>
      <p className="mt-5 max-w-[580px] text-[17px] leading-[1.55] text-[var(--l-ink-mid)]">{body}</p>
    </div>
  );
}

// Inline style escape hatch — use sparingly for values Tailwind can't express
export function inlineStyle(style: CSSProperties) {
  return style;
}
