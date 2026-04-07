"use client";

/**
 * @license
 * SPDX-License-Identifier: Apache-2.0
 */

import React, { useState } from 'react';
import { Sidebar } from '@/components/Sidebar';
import { TopNav } from '@/components/TopNav';
import { Dashboard } from '@/components/Dashboard';
import { Catalog } from '@/components/Catalog';
import { QuickConvert } from '@/components/QuickConvert';
import { OCR } from '@/components/OCR';
import { ESelector } from '@/components/ESelector';
import { motion, AnimatePresence } from 'motion/react';

export default function App() {
  const [activePage, setActivePage] = useState('home');

  const renderPage = () => {
    switch (activePage) {
      case 'home':
        return <Dashboard />;
      case 'catalog':
        return <Catalog />;
      case 'convert':
        return <QuickConvert />;
      case 'ocr':
        return <OCR />;
      case 'eselector':
        return <ESelector />;
      default:
        return <Dashboard />;
    }
  };

  return (
    <div className="flex min-h-screen bg-background">
      <TopNav />
      <Sidebar activePage={activePage} onPageChange={setActivePage} />
      
      <main className="flex-1 md:ml-64 pt-24 pb-12 px-8 lg:px-16 min-h-screen overflow-x-hidden">
        <AnimatePresence mode="wait">
          <motion.div
            key={activePage}
            initial={{ opacity: 0, x: 10 }}
            animate={{ opacity: 1, x: 0 }}
            exit={{ opacity: 0, x: -10 }}
            transition={{ duration: 0.3, ease: "easeInOut" }}
          >
            {renderPage()}
          </motion.div>
        </AnimatePresence>
      </main>

      {/* Bottom Nav for Mobile */}
      <nav className="md:hidden fixed bottom-0 left-0 w-full bg-white/95 backdrop-blur-lg px-6 py-3 flex justify-around items-center z-50 border-t border-slate-100">
        <MobileNavItem 
          active={activePage === 'home'} 
          onClick={() => setActivePage('home')} 
          label="Home" 
          icon="home" 
        />
        <MobileNavItem 
          active={activePage === 'catalog'} 
          onClick={() => setActivePage('catalog')} 
          label="Catalog" 
          icon="menu_book" 
        />
        <MobileNavItem 
          active={activePage === 'convert'} 
          onClick={() => setActivePage('convert')} 
          label="Convert" 
          icon="transform" 
        />
        <MobileNavItem 
          active={activePage === 'ocr'} 
          onClick={() => setActivePage('ocr')} 
          label="OCR" 
          icon="document_scanner" 
        />
        <MobileNavItem 
          active={activePage === 'eselector'} 
          onClick={() => setActivePage('eselector')} 
          label="Select" 
          icon="checklist_rtl" 
        />
      </nav>
    </div>
  );
}

function MobileNavItem({ active, onClick, label, icon }: any) {
  return (
    <button 
      onClick={onClick}
      className={`flex flex-col items-center gap-1 transition-colors ${active ? 'text-primary' : 'text-slate-400'}`}
    >
      <span className="material-symbols-outlined" style={{ fontVariationSettings: `'FILL' ${active ? 1 : 0}` }}>{icon}</span>
      <span className="text-[10px] font-bold">{label}</span>
    </button>
  );
}
