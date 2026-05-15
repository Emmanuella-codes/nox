"use client";

import type { ComponentProps } from "react";
import { useMemo, useState } from "react";
import { ArrowRight, Eye, EyeOff, LockKeyhole, Mail, UserRound } from "lucide-react";
import { AuthField } from "@/src/components/user/auth/auth-field";
import type { SignupFormProps } from "@/src/types/components/auth";
import { signup } from "@/src/utils/api/user/auth";
import { ApiRequestError } from "@/src/utils/api/api";

export function SignupForm({ className }: SignupFormProps) {
  const [firstname, setFirstname] = useState("");
  const [lastname, setLastname] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [status, setStatus] = useState<"idle" | "loading" | "success" | "error">("idle");
  const [message, setMessage] = useState("");

  const canSubmit = useMemo(
    () =>
      firstname.trim().length > 0 &&
      lastname.trim().length > 0 &&
      email.trim().length > 0 &&
      password.trim().length >= 8 &&
      status !== "loading",
    [email, firstname, lastname, password, status],
  );

  const handleSubmit: NonNullable<ComponentProps<"form">["onSubmit"]> = async (event) => {
    event.preventDefault();
    if (!canSubmit) {
      setStatus("error");
      setMessage("Complete all fields. Password must be at least 8 characters.");
      return;
    }

    setStatus("loading");
    setMessage("");
    try {
      await signup({
        firstname: firstname.trim(),
        lastname: lastname.trim(),
        email: email.trim(),
        password,
      });
      setStatus("success");
      setMessage("Verification code sent. Check your email.");
    } catch (error) {
      setStatus("error");
      setMessage(error instanceof ApiRequestError ? error.message : "Unable to create account.");
    }
  };

  return (
    <form className={className} onSubmit={handleSubmit}>
      <div className="grid gap-4">
        <div className="grid grid-cols-2 gap-3">
          <AuthField
            id="firstname"
            label="first"
            type="text"
            value={firstname}
            placeholder="Ada"
            autoComplete="given-name"
            icon={<UserRound className="size-4" strokeWidth={1.7} />}
            onChange={setFirstname}
          />
          <AuthField
            id="lastname"
            label="last"
            type="text"
            value={lastname}
            placeholder="Lovelace"
            autoComplete="family-name"
            icon={<UserRound className="size-4" strokeWidth={1.7} />}
            onChange={setLastname}
          />
        </div>
        <AuthField
          id="signup-email"
          label="email"
          type="email"
          value={email}
          placeholder="you@example.com"
          autoComplete="email"
          icon={<Mail className="size-4" strokeWidth={1.7} />}
          onChange={setEmail}
        />
        <AuthField
          id="signup-password"
          label="password"
          type={showPassword ? "text" : "password"}
          value={password}
          placeholder="Minimum 8 characters"
          autoComplete="new-password"
          icon={<LockKeyhole className="size-4" strokeWidth={1.7} />}
          action={
            <button
              type="button"
              aria-label={showPassword ? "Hide password" : "Show password"}
              onClick={() => setShowPassword((value) => !value)}
              className="rounded-md p-1 text-[var(--nox-ink-soft)] transition hover:text-[var(--nox-ink)]"
            >
              {showPassword ? <EyeOff className="size-4" /> : <Eye className="size-4" />}
            </button>
          }
          onChange={setPassword}
        />
      </div>

      {message ? (
        <p
          className={`mt-4 rounded-[8px] border px-3 py-2 text-[12px] font-medium ${
            status === "success"
              ? "border-[var(--nox-success)] bg-[var(--nox-success-soft)] text-[var(--nox-success)]"
              : "border-[var(--nox-danger)] bg-[var(--nox-danger-soft)] text-[var(--nox-danger)]"
          }`}
        >
          {message}
        </p>
      ) : null}

      <button
        type="submit"
        disabled={!canSubmit}
        className="mt-5 flex min-h-12 w-full items-center justify-center gap-2 rounded-[10px] border border-[var(--nox-accent)] bg-[var(--nox-accent)] px-4 py-3 text-[15px] font-semibold text-white shadow-[0_0_24px_rgba(167,139,250,0.24)] transition hover:brightness-110 disabled:cursor-not-allowed disabled:opacity-50"
      >
        {status === "loading" ? "Creating account" : "Create account"}
        <ArrowRight className="size-4" strokeWidth={1.8} />
      </button>
    </form>
  );
}
