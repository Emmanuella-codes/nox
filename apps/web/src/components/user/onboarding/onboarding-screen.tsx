"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { CalendarDays, Grid2X2, Search, User } from "lucide-react";
import { OnboardingShell } from "@/src/components/user/onboarding/onboarding-shell";
import { OnboardingBrandPanel } from "@/src/components/user/onboarding/onboarding-brand-panel";
import { OnboardingProgress } from "@/src/components/user/onboarding/onboarding-progress";
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

type Step = 0 | 1;
type PersonaCategory = "patron" | "dj" | "organizer";

export function OnboardingScreen() {
  const router = useRouter();
  const [step, setStep] = useState<Step>(0);

  const [handle, setHandle] = useState("");
  const [displayName, setDisplayName] = useState("");
  const [bio, setBio] = useState("");
  const [genreTags, setGenreTags] = useState<string[]>([]);
  const [category] = useState<PersonaCategory>(() => {
    if (typeof window === "undefined") return "patron";
    const stored = localStorage.getItem("nox_signup_category");
    return stored === "dj" || stored === "organizer" || stored === "patron" ? stored : "patron";
  });

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  async function handleFinish() {
    setLoading(true);
    setError("");

    try {
      const token = getAccessToken();
      await createPersona(
        {
          handle,
          display_name: displayName,
          bio,
          avatar_url: "",
          cover_url: "",
          persona_type: "visible",
          category,
          genre_tags: genreTags,
        },
        token,
      );
      router.push("/feed");
    } catch (err) {
      setLoading(false);
      setError(err instanceof ApiRequestError ? err.message : "Something went wrong. Try again.");
    }
  }

  function advanceFromStep0() {
    setStep(1);
  }

  function renderStep() {
    if (step === 0) {
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
          onContinue={advanceFromStep0}
          onBack={() => setStep(0)}
        />
      );
    }

    return (
      <GenreStep
        selected={genreTags}
        loading={loading}
        error={error}
        onChange={setGenreTags}
        onContinue={handleFinish}
        onBack={() => setStep(0)}
      />
    );
  }

  return (
    <OnboardingShell sidePanel={<OnboardingBrandPanel />}>
      <OnboardingProgress current={step} total={2} />

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
