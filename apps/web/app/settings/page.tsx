import { FeedShell } from "@/src/components/user/feed/feed-shell";
import { TabBar } from "@/src/components/user/feed/tab-bar";

export default function SettingsPage() {
  return (
    <FeedShell>
      <main className="flex flex-1 flex-col px-4 py-6">
        <p className="font-mono text-[10px] font-semibold uppercase tracking-[0.16em] text-(--nox-ink-soft)">
          settings
        </p>
        <h1 className="mt-2 text-[24px] font-bold tracking-[-0.03em] text-(--nox-ink)">
          Account settings
        </h1>
        <p className="mt-3 text-[13px] leading-6 text-(--nox-ink-soft)">
          Profile editing is coming soon. Your current persona and posts are available from profile.
        </p>
      </main>
      <TabBar />
    </FeedShell>
  );
}
