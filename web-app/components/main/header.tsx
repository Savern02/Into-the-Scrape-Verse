import Link from "next/link";
import Image from "next/image";
import { AuthButton } from "../auth/auth-button.tsx"
import { Suspense } from 'react';

export function Header() {
  return (
        <nav className="w-full flex justify-center border-b border-b-foreground/10 h-16">
          <div className="w-full max-w-5xl flex justify-between items-center p-3 px-5 text-sm">
            <div className="flex gap-5 items-center font-semibold">
              <Image
                src="/logos/logo.png"
                alt="Cheap Chick"
                width={50}
                height={50}
              />
              <Link href={"/"}>Cheap Chick</Link>
              <div className="flex items-center gap-2">
              </div>
            </div>
            <div>
           <Suspense fallback={<div>Loading...</div>}>
             <AuthButton />
            </Suspense>
            </div>
          </div>
        </nav>
  );
}



