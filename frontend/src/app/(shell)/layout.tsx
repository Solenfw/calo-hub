"use client";

import type { ReactNode } from "react";
import { usePathname } from "next/navigation";
import { motion, AnimatePresence } from "motion/react";
import { BottomNav, LeftSideNav, TopNav } from "@/components/navigation";

export default function ShellLayout({
  children,
}: Readonly<{
  children: ReactNode;
}>) {
  const pathname = usePathname();

  return (
    <div className="flex min-h-screen bg-background">
      <TopNav />
      <LeftSideNav />

      <main className="flex-1 md:ml-64 pt-24 px-8 lg:px-16 min-h-screen overflow-x-clip pb-24 md:pb-12">
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

      <BottomNav />
    </div>
  );
}
