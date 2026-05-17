"use client";

import { useEffect, useRef, useState } from "react";
import { ArrowRight, MailOpen, RotateCcw } from "lucide-react";
import { ApiRequestError } from "@/src/utils/api/api";
import { login, resendVerification, verifyEmail } from "@/src/utils/api/user/auth";
import { useRouter } from "next/navigation";

interface VerifyEmailFormProps {
  email: string;
  password: string;
  expiresInSeconds: number;
}

export function VerifyEmailForm({ email, password, expiresInSeconds }: VerifyEmailFormProps) {
  const router = useRouter();
  const [digits, setDigits] = useState<string[]>(Array(6).fill(""));
  const [status, setStatus] = useState<"idle" | "loading" | "error">("idle");
  const [message, setMessage] = useState("");
  const [resendStatus, setResendStatus] = useState<"idle" | "loading" | "sent">("idle");
  const [countdown, setCountdown] = useState(expiresInSeconds);
  const inputRefs = useRef<(HTMLInputElement | null)[]>([]);

  useEffect(() => {
    inputRefs.current[0]?.focus();
  }, []);

  useEffect(() => {
    if (countdown <= 0) return;
    const id = setInterval(() => setCountdown((c) => c - 1), 1000);
    return () => clearInterval(id);
  }, [countdown]);

  function handleChange(index: number, value: string) {
    const char = value.replace(/\D/g, "").slice(-1);
    const next = [...digits];
    next[index] = char;
    setDigits(next);
    if (char && index < 5) {
      inputRefs.current[index + 1]?.focus();
    }
  }

  function handleKeyDown(index: number, event: React.KeyboardEvent<HTMLInputElement>) {
    if (event.key === "Backspace" && !digits[index] && index > 0) {
      inputRefs.current[index - 1]?.focus();
    }
  }

  function handlePaste(event: React.ClipboardEvent<HTMLInputElement>) {
    event.preventDefault();
    const pasted = event.clipboardData.getData("text").replace(/\D/g, "").slice(0, 6);
    const next = Array(6).fill("");
    for (let i = 0; i < pasted.length; i++) next[i] = pasted[i];
    setDigits(next);
    const focusIdx = Math.min(pasted.length, 5);
    inputRefs.current[focusIdx]?.focus();
  }

  const otp = digits.join("");
  const canSubmit = otp.length === 6 && status !== "loading";

  async function handleSubmit(event: React.FormEvent) {
    event.preventDefault();
    if (!canSubmit) return;
    setStatus("loading");
    setMessage("");
    try {
      await verifyEmail({ email, otp });
      const loginRes = await login({ email, password });
      if (loginRes.data?.tokens.access_token) {
        localStorage.setItem("nox_access_token", loginRes.data.tokens.access_token);
      }
      if (loginRes.data?.tokens.refresh_token) {
        localStorage.setItem("nox_refresh_token", loginRes.data.tokens.refresh_token);
      }
      if (loginRes.data?.user.id) {
        localStorage.setItem("nox_user_id", loginRes.data.user.id);
      }
      router.push("/onboarding");
    } catch (error) {
      setStatus("error");
      setMessage(error instanceof ApiRequestError ? error.message : "Verification failed.");
    }
  }

  async function handleResend() {
    setResendStatus("loading");
    try {
      const res = await resendVerification({ email });
      setCountdown(res.data?.expires_in_seconds ?? expiresInSeconds);
      setDigits(Array(6).fill(""));
      inputRefs.current[0]?.focus();
      setResendStatus("sent");
      setTimeout(() => setResendStatus("idle"), 3000);
    } catch {
      setResendStatus("idle");
    }
  }

  const mins = String(Math.floor(countdown / 60)).padStart(2, "0");
  const secs = String(countdown % 60).padStart(2, "0");

  return (
    <form onSubmit={handleSubmit} className="flex flex-1 flex-col">
      <div className="mb-5 flex size-12 items-center justify-center rounded-[12px] border border-(--nox-accent-line) bg-(--nox-accent-soft) text-(--nox-accent-ink)">
        <MailOpen className="size-6" strokeWidth={1.6} />
      </div>

      <p className="font-mono text-[10px] font-medium uppercase tracking-[0.18em] text-(--nox-ink-soft)">
        check your inbox
      </p>
      <h2 className="mt-3 text-[28px] font-bold leading-none tracking-[-0.04em] text-(--nox-ink)">
        Confirm your email
      </h2>
      <p className="mt-3 text-[13px] leading-6 text-(--nox-ink-mid)">
        We sent a 6-digit code to{" "}
        <span className="font-semibold text-(--nox-ink)">{email}</span>. Expires in{" "}
        <span className="font-mono text-(--nox-accent)">{mins}:{secs}</span>.
      </p>

      <div className="mt-6 flex justify-between gap-2">
        {digits.map((digit, i) => (
          <input
            key={i}
            ref={(el) => { inputRefs.current[i] = el; }}
            type="text"
            inputMode="numeric"
            maxLength={1}
            value={digit}
            onChange={(e) => handleChange(i, e.target.value)}
            onKeyDown={(e) => handleKeyDown(i, e)}
            onPaste={handlePaste}
            className="h-14 w-full rounded-[10px] border border-(--nox-border) bg-(--nox-surface) text-center font-mono text-[22px] font-semibold text-(--nox-ink) outline-none transition focus:border-(--nox-accent-line) focus:ring-1 focus:ring-(--nox-accent-line)"
          />
        ))}
      </div>

      {message ? (
        <p className="mt-4 rounded-[8px] border border-(--nox-danger) bg-(--nox-danger-soft) px-3 py-2 text-[12px] font-medium text-(--nox-danger)">
          {message}
        </p>
      ) : null}

      <button
        type="submit"
        disabled={!canSubmit}
        className="mt-5 flex min-h-12 w-full items-center justify-center gap-2 rounded-[10px] border border-(--nox-accent) bg-(--nox-accent) px-4 py-3 text-[15px] font-semibold text-white shadow-[0_0_24px_rgba(167,139,250,0.24)] transition hover:brightness-110 disabled:cursor-not-allowed disabled:opacity-50"
      >
        {status === "loading" ? "Verifying" : "Confirm code"}
        <ArrowRight className="size-4" strokeWidth={1.8} />
      </button>

      <button
        type="button"
        disabled={resendStatus === "loading" || countdown > 0}
        onClick={handleResend}
        className="mt-4 flex items-center justify-center gap-2 text-[12px] font-medium text-(--nox-ink-soft) transition hover:text-(--nox-ink) disabled:cursor-not-allowed disabled:opacity-40"
      >
        <RotateCcw className="size-3.5" strokeWidth={1.8} />
        {resendStatus === "sent"
          ? "Code resent"
          : countdown > 0
            ? `Resend in ${mins}:${secs}`
            : "Resend code"}
      </button>
    </form>
  );
}
