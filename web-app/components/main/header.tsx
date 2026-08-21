import Link from "next/link";
import Image from "next/image";
import { Suspense } from "react";
import { AuthButton } from "../auth/auth-button";

export function Header() {
  return (
    <nav className="w-full flex justify-center border-b border-b-foreground/10 h-16">
      <div className="w-full max-w-5xl flex justify-between items-center p-3 px-5 text-sm">
        <div className="flex gap-5 items-center font-semibold">
          <Link href="/" className="flex items-center gap-2">
            <Image
              src="/logos/logo.png"
              alt=""
              width={32}
              height={32}
              priority
            />
            <span>Cheap Chick</span>
          </Link>

          <Link
            href="/supplies"
            className="font-normal text-muted-foreground hover:text-foreground transition-colors"
          >
            Supplies
          </Link>
        </div>

        <Suspense fallback={<div className="h-8 w-32" />}>
          <AuthButton />
        </Suspense>
      </div>
    </nav>
  );
}
