"use client";

const PALETTES = [
  ["#6b4dff", "#a78bfa"],
  ["#d4774a", "#f4a87a"],
  ["#3d8a8a", "#7ac4c4"],
  ["#b8842b", "#e6b95c"],
  ["#a8425e", "#d97aa0"],
  ["#3d5a99", "#7a9ad9"],
  ["#5a8a3d", "#a4cc7a"],
  ["#2a2a35", "#5a5a6a"],
] as const;

function paletteIndex(seed: string): number {
  let hash = 0;
  for (let i = 0; i < seed.length; i++) {
    hash = (hash * 31 + seed.charCodeAt(i)) >>> 0;
  }
  return hash % PALETTES.length;
}

interface AvatarProps {
  id: string;
  name: string;
  size?: number;
  square?: boolean;
  className?: string;
}

export function Avatar({ id, name, size = 36, square = false, className = "" }: AvatarProps) {
  const idx = paletteIndex(id);
  const [from, to] = PALETTES[idx];
  const gradId = `av-${id.slice(0, 8)}`;
  const r = square ? Math.round(size * 0.17) : size / 2;

  return (
    <svg
      width={size}
      height={size}
      viewBox={`0 0 ${size} ${size}`}
      aria-hidden="true"
      className={className}
    >
      <defs>
        <linearGradient id={gradId} x1="0" y1="0" x2="1" y2="1">
          <stop offset="0%" stopColor={from} />
          <stop offset="100%" stopColor={to} />
        </linearGradient>
      </defs>
      <rect x="0" y="0" width={size} height={size} rx={r} ry={r} fill={`url(#${gradId})`} />
      <text
        x={size / 2}
        y={size / 2 + size * 0.13}
        textAnchor="middle"
        fontSize={size * 0.38}
        fontWeight="700"
        fill="rgba(255,255,255,0.9)"
        fontFamily="Space Grotesk, Inter, sans-serif"
      >
        {name[0]?.toUpperCase() ?? "?"}
      </text>
    </svg>
  );
}
