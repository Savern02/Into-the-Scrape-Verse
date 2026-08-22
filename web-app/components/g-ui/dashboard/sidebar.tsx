"use client"

import Link from "next/link"
import { usePathname } from "next/navigation"
import { LogoutButton } from "@/components/auth/logout-button"
import {
  LayoutDashboard,
  ShoppingCart,
  BarChart3,
  Settings,
} from "lucide-react"

import { Button } from "@/components/ui/button"
import { cn } from "@/lib/utils"

const links = [
  {
    name: "Overview",
    href: "/dashboard",
    icon: LayoutDashboard,
  },
  {
    name: "Products",
    href: "/dashboard/products",
    icon: ShoppingCart,
  },
  {
    name: "Comparisons",
    href: "/dashboard/comparisons",
    icon: BarChart3,
  },
  {
    name: "Settings",
    href: "/dashboard/settings",
    icon: Settings,
  },
]

export function Sidebar() {
  const pathname = usePathname()

  return (
    <aside className="flex h-full w-64 flex-col border-r bg-background p-4">
      <div className="mb-6 px-2">
        <h2 className="text-lg font-semibold">
          Dashboard
        </h2>
      </div>

      <nav className="flex flex-col gap-2">
        {links.map((link) => {
          const Icon = link.icon

          const active =
            pathname === link.href ||
            pathname.startsWith(`${link.href}/`)

          return (
            <Button
              key={link.href}
              variant={active ? "secondary" : "ghost"}
              className={cn(
                "w-full justify-start gap-3",
                active && "font-semibold"
              )}
              asChild
            >
              <Link href={link.href}>
                <Icon className="h-4 w-4" />
                {link.name}
              </Link>
            </Button>
          )
        })}
      </nav>
      <div className="mt-auto pt-4">
          <Button asChild variant="outline">
            <Link href="/">
              Home
            </Link>
          </Button>
        <LogoutButton />
      </div>
    </aside>
  )
}
