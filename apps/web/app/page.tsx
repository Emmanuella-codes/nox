import Link from "next/link";

export default function Home() {
  return (
    <main className="flex min-h-dvh items-center justify-center bg-[#0a0a0c] px-6 text-[#ededf2]">
      <section className="w-full max-w-md">
        <p className="font-mono text-[10px] font-medium uppercase tracking-[0.18em] text-[#6b6878]">
          nox
        </p>
        <h1 className="mt-3 text-[44px] font-bold leading-none tracking-[-0.04em]">
          Lagos nightlife, no filter.
        </h1>
        <p className="mt-4 text-[15px] leading-7 text-[#9a96a8]">
          Public when it should be public. Anonymous when it should stay with the moment.
        </p>
        <Link
          href="/auth"
          className="mt-7 inline-flex min-h-11 items-center justify-center rounded-[8px] border border-[#a78bfa] bg-[#a78bfa] px-5 text-[14px] font-semibold text-white transition hover:brightness-110"
        >
          Continue
        </Link>
      </section>
    </main>
  );
}
