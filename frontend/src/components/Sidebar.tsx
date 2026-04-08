import React from 'react';
import { 
  Home, 
  BookOpen, 
  RefreshCw, 
  Scan, 
  CheckSquare, 
  Activity
} from 'lucide-react';
import { cn } from '@/lib/utils';

interface SidebarProps {
  activePage: string;
  onPageChange: (page: string) => void;
}

const navItems = [
  { id: 'home', label: 'Home', icon: Home },
  { id: 'catalog', label: 'Online Catalog', icon: BookOpen },
  { id: 'convert', label: 'Quick Convert', icon: RefreshCw },
  { id: 'ocr', label: 'OCR', icon: Scan },
  { id: 'eselector', label: 'eSelector', icon: CheckSquare },
];

export function Sidebar({ activePage, onPageChange }: SidebarProps) {
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
          const isActive = activePage === item.id;
          
          return (
            <button
              key={item.id}
              onClick={() => onPageChange(item.id)}
              className={cn(
                "cursor-pointer group flex items-center gap-4 px-6 py-3 ml-4 transition-all duration-200 hover:translate-x-1",
                isActive 
                  ? "bg-white dark:bg-slate-900 text-blue-700 dark:text-blue-400 rounded-l-full shadow-sm font-bold" 
                  : "text-slate-500 dark:text-slate-400 hover:text-blue-600 dark:hover:text-blue-300"
              )}
            >
              <Icon className={cn("w-5 h-5")} />
              <span className="font-inter text-sm font-medium">{item.label}</span>
            </button>
          );
        })}
      </nav>
    </aside>
  );
}
