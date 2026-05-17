import type { AuthShellProps } from "@/src/types/components/auth";
import { authThemeVars } from "@/src/components/user/auth/auth-theme";

export function AuthShell({ theme, children, sidePanel }: AuthShellProps) {
  return (
    <main
      className="min-h-dvh w-full overflow-hidden bg-(--nox-bg) text-(--nox-ink)"
      style={authThemeVars(theme)}
    >
      <div className="mx-auto grid min-h-dvh w-full max-w-6xl items-center gap-8 px-4 py-6 md:grid-cols-[minmax(320px,390px)_1fr] md:px-8 lg:gap-12">
        <section className="mx-auto flex h-[min(680px,calc(100dvh-48px))] min-h-[620px] w-full max-w-[360px] flex-col overflow-hidden rounded-[22px] border border-(--nox-border-strong) bg-(--nox-bg-soft) shadow-(--nox-shadow) md:max-w-[360px]">
          {children}
        </section>
        <aside className="hidden min-h-[620px] md:block">{sidePanel}</aside>
      </div>
    </main>
  );
}
