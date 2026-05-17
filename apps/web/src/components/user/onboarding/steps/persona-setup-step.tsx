"use client";

import { useMemo, useState } from "react";
import { ArrowRight, AtSign, Type, AlignLeft } from "lucide-react";
import { AuthField } from "@/src/components/user/auth/auth-field";

interface PersonaSetupStepProps {
  handle: string;
  displayName: string;
  bio: string;
  onChange: (fields: { handle: string; displayName: string; bio: string }) => void;
  onContinue: () => void;
  onBack: () => void;
}

export function PersonaSetupStep({
  handle,
  displayName,
  bio,
  onChange,
  onContinue,
  onBack,
}: PersonaSetupStepProps) {
  const [handleError, setHandleError] = useState("");

  const handleValue = handle.toLowerCase().replace(/[^a-z0-9_]/g, "");

  function validateHandle(value: string) {
    if (value.length < 3) return "At least 3 characters";
    if (value.length > 30) return "Max 30 characters";
    if (!/^[a-z0-9_]+$/.test(value)) return "Letters, numbers, and underscores only";
    return "";
  }

  function onHandleChange(raw: string) {
    const cleaned = raw.toLowerCase().replace(/[^a-z0-9_]/g, "");
    setHandleError(validateHandle(cleaned));
    onChange({ handle: cleaned, displayName, bio });
  }

  const canContinue = useMemo(
    () =>
      handle.trim().length >= 3 &&
      displayName.trim().length > 0 &&
      !handleError,
    [handle, displayName, handleError],
  );

  return (
    <div className="flex flex-1 flex-col px-5">
      <p className="font-mono text-[10px] font-medium uppercase tracking-[0.18em] text-(--nox-ink-soft)">
        Step 2 of 3
      </p>
      <h2 className="mt-2 text-[24px] font-bold leading-tight tracking-[-0.03em] text-(--nox-ink)">
        Set up your persona
      </h2>
      <p className="mt-2 text-[13px] leading-6 text-(--nox-ink-soft)">
        This is what the scene sees when you post publicly.
      </p>

      <div className="mt-5 flex flex-col gap-4">
        <div>
          <AuthField
            id="handle"
            label="handle"
            type="text"
            value={handleValue}
            placeholder="djkayode"
            autoComplete="off"
            icon={<AtSign className="size-4" strokeWidth={1.7} />}
            error={handleError}
            onChange={onHandleChange}
          />
          <p className="mt-1.5 font-mono text-[10px] text-(--nox-ink-soft)">
            nox.app/@{handleValue || "yourhandle"}
          </p>
        </div>

        <AuthField
          id="display-name"
          label="display name"
          type="text"
          value={displayName}
          placeholder="DJ Kayode"
          autoComplete="off"
          icon={<Type className="size-4" strokeWidth={1.7} />}
          onChange={(v) => onChange({ handle, displayName: v, bio })}
        />

        <label className="block" htmlFor="bio">
          <span className="mb-2 block font-mono text-[10px] font-medium uppercase tracking-[0.18em] text-(--nox-ink-soft)">
            bio{" "}
            <span className="normal-case tracking-normal text-(--nox-ink-faint)">(optional)</span>
          </span>
          <span className="flex items-start gap-3 rounded-[10px] border border-(--nox-border) bg-(--nox-surface) px-3.5 py-3 transition focus-within:border-(--nox-accent-line)">
            <AlignLeft className="mt-0.5 size-4 shrink-0 text-(--nox-ink-soft)" strokeWidth={1.7} />
            <textarea
              id="bio"
              value={bio}
              placeholder="Lagos DJ, Afro-House devotee…"
              maxLength={160}
              rows={3}
              onChange={(e) => onChange({ handle, displayName, bio: e.target.value })}
              className="min-w-0 flex-1 resize-none bg-transparent text-[14px] font-medium text-(--nox-ink) outline-none placeholder:text-(--nox-ink-soft)"
            />
          </span>
          <p className="mt-1 text-right font-mono text-[10px] text-(--nox-ink-faint)">
            {bio.length}/160
          </p>
        </label>
      </div>

      <div className="mt-auto flex gap-3 pb-1">
        <button
          type="button"
          onClick={onBack}
          className="flex min-h-12 items-center justify-center rounded-[10px] border border-(--nox-border-strong) px-5 text-[14px] font-semibold text-(--nox-ink) transition hover:border-(--nox-ink)"
        >
          Back
        </button>
        <button
          type="button"
          disabled={!canContinue}
          onClick={onContinue}
          className="flex flex-1 min-h-12 items-center justify-center gap-2 rounded-[10px] border border-(--nox-accent) bg-(--nox-accent) px-4 py-3 text-[15px] font-semibold text-white shadow-[0_0_24px_rgba(167,139,250,0.24)] transition hover:brightness-110 disabled:cursor-not-allowed disabled:opacity-40"
        >
          Continue
          <ArrowRight className="size-4" strokeWidth={1.8} />
        </button>
      </div>
    </div>
  );
}
