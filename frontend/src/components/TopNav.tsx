import React from 'react';
import { Search, Settings, Bell } from 'lucide-react';

export function TopNav() {
  return (
    <header className="fixed top-0 w-full z-50 glass-nav bg-white/80 dark:bg-slate-900/80 shadow-[0_12px_40px_rgba(0,93,182,0.06)] flex justify-between items-center px-8 h-16">
      <div className="flex items-center gap-4">
        <span className="font-manrope tracking-tight font-bold text-2xl text-blue-700 dark:text-blue-400">Calo Hub</span>
      </div>
      
      <div className="flex items-center gap-6">
        <div className="hidden md:flex items-center bg-surface-container-high rounded-full px-4 py-1.5 gap-2">
          <Search className="text-outline w-4 h-4" />
          <input 
            type="text" 
            placeholder="Global search..." 
            className="bg-transparent border-none focus:ring-0 text-sm w-48 placeholder:text-on-surface-variant"
          />
        </div>
        
        <div className="flex items-center gap-4">
          <button className="active:scale-95 transition-transform text-outline hover:text-primary">
            <Bell className="w-5 h-5" />
          </button>
          <button className="active:scale-95 transition-transform text-outline hover:text-primary">
            <Settings className="w-5 h-5" />
          </button>
          
          <div className="flex items-center gap-3 pl-2 border-l border-outline-variant/20">
            <div className="text-right hidden sm:block">
              <p className="text-xs font-bold text-on-surface">Dr. Vincent Cody</p>
              <p className="text-[10px] text-on-surface-variant font-medium">Clinical Lead</p>
            </div>
            <img 
              src="https://picsum.photos/seed/doctor/100/100" 
              alt="User Profile" 
              className="w-8 h-8 rounded-full object-cover ring-2 ring-primary-container"
              referrerPolicy="no-referrer"
            />
          </div>
        </div>
      </div>
    </header>
  );
}
