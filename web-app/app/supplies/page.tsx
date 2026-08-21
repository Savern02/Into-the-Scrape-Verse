import { Suspense } from "react";
import { Header } from "@/components/main/header";
import { Hero } from "@/components/main/hero";
import { Footer } from "@/components/main/footer";

export default function Home() {
  return (
    <main className="min-h-screen flex flex-col items-center">
      <div className="flex-1 w-full flex flex-col items-center">
        <Header />

        <div className="flex-1 w-full max-w-5xl px-5">
          {/* cacheComponents is on in next.config.ts, so any dynamic read has
              to sit behind a Suspense boundary or the build fails. */}
          <Suspense fallback={<HeroSkeleton />}>
            <Hero />
          </Suspense>
        </div>

        <Footer />
      </div>
    </main>
  );
}

function HeroSkeleton() {
  return (
    <div className="flex flex-col gap-8 py-12" aria-hidden>
      <div className="flex flex-col gap-4 max-w-2xl">
        <div className="h-3 w-40 rounded bg-muted animate-pulse" />
        <div className="h-12 w-full rounded bg-muted animate-pulse" />
        <div className="h-4 w-3/4 rounded bg-muted animate-pulse" />
      </div>
      <div className="flex gap-10">
        {Array.from({ length: 3 }).map((_, i) => (
          <div key={i} className="flex flex-col gap-2">
            <div className="h-3 w-24 rounded bg-muted animate-pulse" />
            <div className="h-8 w-20 rounded bg-muted animate-pulse" />
          </div>
        ))}
      </div>
    </div>
  );
}
