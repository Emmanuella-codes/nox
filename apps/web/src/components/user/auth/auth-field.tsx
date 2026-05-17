"use client";

import type { AuthFieldProps } from "@/src/types/components/auth";

export function AuthField({
  id,
  label,
  type,
  value,
  placeholder,
  autoComplete,
  icon,
  action,
  error,
  onChange,
}: AuthFieldProps) {
  return (
    <label className="block" htmlFor={id}>
      <span className="mb-2 block font-mono text-[10px] font-medium uppercase tracking-[0.18em] text-(--nox-ink-soft)">
        {label}
      </span>
      <span className="flex min-h-12 items-center gap-3 rounded-[10px] border border-(--nox-border) bg-(--nox-surface) px-3.5 text-(--nox-ink) transition focus-within:border-(--nox-accent-line)">
        <span className="text-(--nox-ink-soft)">{icon}</span>
        <input
          id={id}
          type={type}
          value={value}
          placeholder={placeholder}
          autoComplete={autoComplete}
          onChange={(event) => onChange(event.target.value)}
          className="min-w-0 flex-1 bg-transparent py-3 text-[14px] font-medium text-(--nox-ink) outline-none placeholder:text-(--nox-ink-soft)"
        />
        {action}
      </span>
      {error ? (
        <span className="mt-2 block text-[12px] font-medium text-(--nox-danger)">
          {error}
        </span>
      ) : null}
    </label>
  );
}
