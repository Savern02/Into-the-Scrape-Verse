import Link from "next/link";
import { ThemeSwitcher } from "@/components/theme-switcher";

/**
 * Provenance belongs in the footer of a price comparison site. Anyone acting
 * on these numbers deserves to know they were scraped, when, and from where —
 * a price with no stated source is a price you can't audit.
 */
export function Footer() {
  return (
    <footer className="w-full border-t mt-24">
      <div className="w-full max-w-5xl mx-auto px-5 py-10 flex flex-col gap-8">
        <div className="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-8">
          <div className="flex flex-col gap-2 max-w-sm">
            <span className="font-semibold text-sm">Cheap Chick</span>
            <p className="text-xs text-muted-foreground text-pretty leading-relaxed">
              Per-unit pricing for classroom supplies. Collected from public
              supplier listings, normalized for pack size, and refreshed on a
              schedule.
            </p>
          </div>

          <nav className="flex flex-col gap-2 text-xs">
            <span className="uppercase tracking-wide text-muted-foreground">
              Browse
            </span>
            <Link
              href="/supplies"
              className="hover:underline underline-offset-4"
            >
              All supplies
            </Link>
            <Link href="/" className="hover:underline underline-offset-4">
              Home
            </Link>
          </nav>

          <div className="flex flex-col gap-2 text-xs">
            <span className="uppercase tracking-wide text-muted-foreground">
              Built with
            </span>
            <a
              href="https://brightdata.com"
              target="_blank"
              rel="noopener noreferrer"
              className="hover:underline underline-offset-4"
            >
              Bright Data Scraper Studio
            </a>
            <span className="text-muted-foreground">Go · Postgres · R2</span>
          </div>
        </div>

        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 pt-6 border-t text-xs text-muted-foreground">
          <p className="text-pretty">
            Prices are scraped from publicly available listings and may lag the
            retailer. Confirm before purchasing. Not affiliated with any
            retailer or school district.
          </p>
          <ThemeSwitcher />
        </div>
      </div>
    </footer>
  );
}
