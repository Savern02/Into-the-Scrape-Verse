import { AuthButton } from "@/components/auth-button";
import { Hero } from "@/components/main/hero";
import {Header} from "@/components/main/header"
import { ThemeSwitcher } from "@/components/theme-switcher";
import { Suspense } from "react";
 import {YellowBar} from "@/components/YellowBar"; 

export default function Home() {
  return (
    <main className="min-h-screen flex flex-col items-center">
      <div className="flex-1 w-full flex flex-col items-center">
          <Header />
          <YellowBar side="right" />
        <div className="flex-1 flex w-full flex-col gap-20 max-w-5xl p-10 bg-seaweed-500 text-white">
          <Hero />
          <div className="flex-1 flex flex-col gap-6 px-4">
          </div>
        </div>

        <footer className="w-full flex items-center justify-center border-t mx-auto text-center text-xs gap-8 py-16">
          <ThemeSwitcher />
        </footer>
        <YellowBar side="left" />
      </div>
    </main>
  );
}
