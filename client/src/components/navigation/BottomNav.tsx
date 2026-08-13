"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { cn } from "@/lib/utils";
import { isNavigationActive, navigationItems } from "./navItems";

export function BottomNav() {
  const pathname = usePathname();

  return (
    <nav className="md:hidden fixed bottom-0 left-0 w-full bg-white/95 backdrop-blur-lg px-6 py-3 flex justify-around items-center z-50 border-t border-slate-100">
      {navigationItems.map((item) => {
        const Icon = item.icon;
        const active = isNavigationActive(pathname, item.href);

        return (
          <Link
            key={item.href}
            href={item.href}
            className={cn(
              "flex flex-col items-center gap-1 transition-colors",
              active ? "text-primary" : "text-slate-400"
            )}
          >
            <Icon className="h-5 w-5" strokeWidth={active ? 2.5 : 2} />
            <span className="text-[10px] font-bold">{item.mobileLabel}</span>
          </Link>
        );
      })}
    </nav>
  );
}
