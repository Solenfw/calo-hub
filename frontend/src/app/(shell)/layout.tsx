"use client";

import type { ReactNode } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import {
  Home,
  BookOpen,
  RefreshCw,
  Scan,
  CheckSquare,
} from "lucide-react";
import { motion, AnimatePresence } from "motion/react";
import { TopNav } from "@/components/TopNav";
import { Sidebar } from "@/components/Sidebar";
import { cn } from "@/lib/utils";

const mobileNav = [
  { href: "/", label: "Home", icon: Home },
  { href: "/catalog", label: "Catalog", icon: BookOpen },
  { href: "/convert", label: "Convert", icon: RefreshCw },
  { href: "/ocr", label: "OCR", icon: Scan },
  { href: "/eselector", label: "Select", icon: CheckSquare },
] as const;

function mobileActive(pathname: string, href: string) {
  if (href === "/") {
    return pathname === "/" || pathname === "/dashboard";
  }
  return pathname === href || pathname.startsWith(`${href}/`);
}

export default function ShellLayout({
  children,
}: Readonly<{
  children: ReactNode;
}>) {
  const pathname = usePathname();

  return (
    <div className="flex min-h-screen bg-background">
      <TopNav />
      <Sidebar />

      <main className="flex-1 md:ml-64 pt-24 px-8 lg:px-16 min-h-screen overflow-x-hidden pb-24 md:pb-12">
        <AnimatePresence mode="wait">
          <motion.div
            key={pathname}
            initial={{ opacity: 0, x: 10 }}
            animate={{ opacity: 1, x: 0 }}
            exit={{ opacity: 0, x: -10 }}
            transition={{ duration: 0.3, ease: "easeInOut" }}
          >
            {children}
          </motion.div>
        </AnimatePresence>
      </main>

      <nav className="md:hidden fixed bottom-0 left-0 w-full bg-white/95 backdrop-blur-lg px-6 py-3 flex justify-around items-center z-50 border-t border-slate-100">
        {mobileNav.map((item) => {
          const Icon = item.icon;
          const active = mobileActive(pathname, item.href);
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
              <span className="text-[10px] font-bold">{item.label}</span>
            </Link>
          );
        })}
      </nav>
    </div>
  );
}
