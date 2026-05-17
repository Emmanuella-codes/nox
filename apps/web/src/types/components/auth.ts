import type { ComponentProps, CSSProperties, ReactNode } from "react";

export type AuthMode = "login" | "signup";

export interface AuthTheme {
  bg: string;
  bgSoft: string;
  surface: string;
  surfaceAlt: string;
  surfaceSunk: string;
  border: string;
  borderStrong: string;
  divider: string;
  ink: string;
  inkMid: string;
  inkSoft: string;
  accent: string;
  accentInk: string;
  accentSoft: string;
  accentLine: string;
  danger: string;
  dangerSoft: string;
  success: string;
  successSoft: string;
  gold: string;
  shadow: string;
}

export interface AuthShellProps {
  theme: AuthTheme;
  children: ReactNode;
  sidePanel: ReactNode;
}

export interface AuthFieldProps {
  id: string;
  label: string;
  type: string;
  value: string;
  placeholder: string;
  autoComplete: string;
  icon: ReactNode;
  action?: ReactNode;
  error?: string;
  onChange: (value: string) => void;
}

export interface LoginFormProps {
  theme: AuthTheme;
  onSubmit?: ComponentProps<"form">["onSubmit"];
  className?: string;
  style?: CSSProperties;
}

export interface SignupFormProps {
  theme: AuthTheme;
  onSubmit?: ComponentProps<"form">["onSubmit"];
  onSuccess?: (email: string, password: string, expiresInSeconds: number) => void;
  className?: string;
  style?: CSSProperties;
}
