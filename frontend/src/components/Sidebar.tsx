"use client";

import React from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { 
  Home, 
  BookOpen, 
  RefreshCw, 
  Scan, 
  CheckSquare, 
  Activity
} from 'lucide-react';
import { cn } from '@/lib/utils';

const navItems = [
  { href: '/', label: 'Home', icon: Home },
  { href: '/catalog', label: 'Online Catalog', icon: BookOpen },
  { href: '/convert', label: 'Quick Convert', icon: RefreshCw },
  { href: '/ocr', label: 'OCR', icon: Scan },
  { href: '/eselector', label: 'eSelector', icon: CheckSquare },
];

function isNavActive(pathname: string, href: string) {
  if (href === '/') {
    return pathname === '/' || pathname === '/dashboard';
  }
  return pathname === href || pathname.startsWith(`${href}/`);
}

export function Sidebar() {
  const pathname = usePathname();

  return (
    <aside className="h-screen w-64 fixed left-0 top-0 pt-20 bg-slate-50 dark:bg-slate-950 flex-col gap-2 border-r-0 z-40 hidden md:flex">
      <div className="px-8 mb-8">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 bg-primary rounded-xl flex items-center justify-center">
            <Activity className="text-white w-6 h-6" />
          </div>
          <div>
            <h2 className="font-manrope text-sm font-bold text-blue-700 dark:text-blue-400 leading-none">Calo Hub</h2>
            <p className="text-[10px] uppercase tracking-widest text-outline-variant font-semibold mt-1">Medical Production</p>
          </div>
        </div>
      </div>
      
      <nav className="flex flex-col gap-1 pr-4">
        {navItems.map((item) => {
          const Icon = item.icon;
          const active = isNavActive(pathname, item.href);
          
          return (
            <Link
              key={item.href}
              href={item.href}
              className={cn(
                "cursor-pointer group flex items-center gap-4 px-6 py-3 ml-4 transition-all duration-200 hover:translate-x-1",
                active 
                  ? "bg-white dark:bg-slate-900 text-blue-700 dark:text-blue-400 rounded-l-full shadow-sm font-bold" 
                  : "text-slate-500 dark:text-slate-400 hover:text-blue-600 dark:hover:text-blue-300"
              )}
            >
              <Icon className={cn("w-5 h-5")} />
              <span className="font-inter text-sm font-medium">{item.label}</span>
            </Link>
          );
        })}
      </nav>
    </aside>
  );
}
