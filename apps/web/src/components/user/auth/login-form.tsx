"use client";

import type { ComponentProps } from "react";
import { useMemo, useState } from "react";
import { ArrowRight, Eye, EyeOff, LockKeyhole, Mail, ShieldCheck } from "lucide-react";
import { AuthField } from "@/src/components/user/auth/auth-field";
import type { LoginFormProps } from "@/src/types/components/auth";
import { login } from "@/src/utils/api/user/auth";
import { ApiRequestError } from "@/src/utils/api/api";

export function LoginForm({ className }: LoginFormProps) {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [remember, setRemember] = useState(true);
  const [status, setStatus] = useState<"idle" | "loading" | "success" | "error">("idle");
  const [message, setMessage] = useState("");

  const canSubmit = useMemo(
    () => email.trim().length > 0 && password.trim().length > 0 && status !== "loading",
    [email, password, status],
  );

  const handleSubmit: NonNullable<ComponentProps<"form">["onSubmit"]> = async (event) => {
    event.preventDefault();
    if (!canSubmit) {
      setStatus("error");
      setMessage("Email and password are required.");
      return;
    }

    setStatus("loading");
    setMessage("");
    try {
      await login({ email: email.trim(), password });
      setStatus("success");
      setMessage(remember ? "Session ready." : "Signed in for this visit.");
    } catch (error) {
      setStatus("error");
      setMessage(error instanceof ApiRequestError ? error.message : "Unable to sign in.");
    }
  };

  return (
    <form className={className} onSubmit={handleSubmit}>
      <div className="grid gap-4">
        <AuthField
          id="email"
          label="email"
          type="email"
          value={email}
          placeholder="you@example.com"
          autoComplete="email"
          icon={<Mail className="size-4" strokeWidth={1.7} />}
          onChange={setEmail}
        />
        <AuthField
          id="password"
          label="password"
          type={showPassword ? "text" : "password"}
          value={password}
          placeholder="Enter password"
          autoComplete="current-password"
          icon={<LockKeyhole className="size-4" strokeWidth={1.7} />}
          action={
            <button
              type="button"
              aria-label={showPassword ? "Hide password" : "Show password"}
              onClick={() => setShowPassword((value) => !value)}
              className="rounded-md p-1 text-(--nox-ink-soft) transition hover:text-(--nox-ink)"
            >
              {showPassword ? <EyeOff className="size-4" /> : <Eye className="size-4" />}
            </button>
          }
          onChange={setPassword}
        />
      </div>

      <div className="mt-4 flex items-center justify-between gap-3 text-[12px] font-medium">
        <label className="flex min-w-0 items-center gap-2 text-(--nox-ink-mid)">
          <input
            type="checkbox"
            checked={remember}
            onChange={(event) => setRemember(event.target.checked)}
            className="size-4 rounded border-(--nox-border) accent-(--nox-accent)"
          />
          <span>Keep me in</span>
        </label>
        <a className="text-(--nox-accent-ink) transition hover:text-(--nox-ink)" href="#">
          Forgot password
        </a>
      </div>

      {message ? (
        <p
          className={`mt-4 rounded-[8px] border px-3 py-2 text-[12px] font-medium ${
            status === "success"
              ? "border-(--nox-success) bg-(--nox-success-soft) text-(--nox-success)"
              : "border-(--nox-danger) bg-(--nox-danger-soft) text-(--nox-danger)"
          }`}
        >
          {message}
        </p>
      ) : null}

      <button
        type="submit"
        disabled={!canSubmit}
        className="mt-5 flex min-h-12 w-full items-center justify-center gap-2 rounded-[10px] border border-(--nox-accent) bg-(--nox-accent) px-4 py-3 text-[15px] font-semibold text-white shadow-[0_0_24px_rgba(167,139,250,0.24)] transition hover:brightness-110 disabled:cursor-not-allowed disabled:opacity-50"
      >
        {status === "loading" ? "Signing in" : "Login"}
        <ArrowRight className="size-4" strokeWidth={1.8} />
      </button>

      <div className="mt-4 flex items-center justify-center gap-2 text-center text-[12px] text-(--nox-ink-soft)">
        <ShieldCheck className="size-4 text-(--nox-accent-ink)" strokeWidth={1.7} />
        <span>Private account, clean public trail.</span>
      </div>
    </form>
  );
}
