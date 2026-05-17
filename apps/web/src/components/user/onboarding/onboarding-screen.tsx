"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { CalendarDays, Grid2X2, Search, User } from "lucide-react";
import { OnboardingShell } from "@/src/components/user/onboarding/onboarding-shell";
import { OnboardingBrandPanel } from "@/src/components/user/onboarding/onboarding-brand-panel";
import { OnboardingProgress } from "@/src/components/user/onboarding/onboarding-progress";
import { PersonaTypeStep, type PersonaMode } from "@/src/components/user/onboarding/steps/persona-type-step";
import { PersonaSetupStep } from "@/src/components/user/onboarding/steps/persona-setup-step";
import { GenreStep } from "@/src/components/user/onboarding/steps/genre-step";
import { createPersona } from "@/src/utils/api/user/persona";
import { getAccessToken } from "@/src/utils/auth/session";
import { ApiRequestError } from "@/src/utils/api/api";

const navItems = [
  { label: "feed", icon: Grid2X2 },
  { label: "discover", icon: Search },
  { label: "events", icon: CalendarDays },
  { label: "profile", icon: User },
];

// Ghost onboarding: steps 0 → 1 (type → genres)
// Visible onboarding: steps 0 → 1 → 2 (type → setup → genres)
type Step = 0 | 1 | 2;

export function OnboardingScreen() {
  const router = useRouter();
  const [step, setStep] = useState<Step>(0);
  const [personaMode, setPersonaMode] = useState<PersonaMode | null>(null);

  const [handle, setHandle] = useState("");
  const [displayName, setDisplayName] = useState("");
  const [bio, setBio] = useState("");
  const [genreTags, setGenreTags] = useState<string[]>([]);

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const isGhost = personaMode === "ghost";
  const totalSteps = isGhost ? 2 : 3;

  async function handleFinish() {
    setLoading(true);
    setError("");

    try {
      if (personaMode === "visible") {
        const token = getAccessToken();
        await createPersona(
          {
            handle,
            display_name: displayName,
            bio,
            avatar_url: "",
            cover_url: "",
            persona_type: "visible",
            genre_tags: genreTags,
          },
          token,
        );
      }
      router.push("/feed");
    } catch (err) {
      setLoading(false);
      setError(err instanceof ApiRequestError ? err.message : "Something went wrong. Try again.");
    }
  }

  function advanceFromStep0() {
    if (!personaMode) return;
    setStep(1);
  }

  function advanceFromStep1() {
    // step 1 is persona-setup (visible only). ghost goes directly to genres at step 1
    setStep(2);
  }

  function renderStep() {
    if (step === 0) {
      return (
        <PersonaTypeStep
          selected={personaMode}
          onSelect={setPersonaMode}
          onContinue={advanceFromStep0}
        />
      );
    }

    if (step === 1 && !isGhost) {
      return (
        <PersonaSetupStep
          handle={handle}
          displayName={displayName}
          bio={bio}
          onChange={({ handle: h, displayName: d, bio: b }) => {
            setHandle(h);
            setDisplayName(d);
            setBio(b);
          }}
          onContinue={advanceFromStep1}
          onBack={() => setStep(0)}
        />
      );
    }

    // step 1 for ghost, step 2 for visible → genres
    return (
      <GenreStep
        selected={genreTags}
        isGhost={isGhost}
        loading={loading}
        error={error}
        onChange={setGenreTags}
        onContinue={handleFinish}
        onBack={() => setStep(isGhost ? 0 : 1)}
      />
    );
  }

  return (
    <OnboardingShell sidePanel={<OnboardingBrandPanel />}>
      <OnboardingProgress current={step} total={totalSteps} />

      <div className="flex flex-1 flex-col gap-5 overflow-y-auto py-5">
        {renderStep()}
      </div>

      <div className="border-t border-(--nox-divider)">
        <nav className="grid grid-cols-4 pt-2 pb-3">
          {navItems.map((item, index) => {
            const Icon = item.icon;
            return (
              <span
                key={item.label}
                className={`flex flex-col items-center gap-1 text-[9.5px] font-medium ${
                  index === 3 ? "text-(--nox-accent)" : "text-(--nox-ink-soft)"
                }`}
              >
                <Icon className="size-4" strokeWidth={1.7} />
                {item.label}
              </span>
            );
          })}
        </nav>
      </div>
    </OnboardingShell>
  );
}
