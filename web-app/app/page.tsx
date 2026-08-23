import { Hero } from "@/components/main/hero";
import {Header} from "@/components/main/header"
import { ThemeSwitcher } from "@/components/theme-switcher";
import {YellowBar} from "@/components/main/yellow-bar";
import { Suspense } from "react";

export default function Home() {
  return (
    <main className="min-h-screen flex flex-col items-center">
      <YellowBar side="left"/>
      <YellowBar side="right"/>
      <div className="flex-1 w-full flex flex-col items-center ">
          <Header />
        <div className="flex-1 flex w-full flex-col gap-20 max-w-6xl p-10  text-white">
          <Hero />
          <main className="flex-1 flex flex-col gap-6 px-4">
          </main>
        </div>

        <footer className="w-full flex items-center justify-center border-t mx-auto text-center text-xs gap-8 py-16">
          <ThemeSwitcher />
        </footer>
      </div>
    </main>
  );
}
