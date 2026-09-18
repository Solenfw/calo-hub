"use client";

import React from 'react';
import { 
  BookOpen, 
  RefreshCw, 
  Scan, 
  CheckSquare, 
  ArrowRight, 
  Rocket, 
  FileText, 
  MoreVertical, 
} from 'lucide-react';
import { motion } from 'motion/react';
import { cn } from '@/lib/utils';
import Link from 'next/link';

export function Dashboard() {
  return (
    <div className="max-w-6xl mx-auto">
      {/* Welcome Hero Section */}
      <section className="mb-12">
        <motion.div 
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5 }}
          className="max-w-4xl"
        >
          <h1 className="text-5xl font-extrabold text-on-surface tracking-tight mb-4 leading-tight">
            Welcome back, <span className="text-primary">Dr. Cody</span>.
          </h1>
          <p className="text-on-surface-variant text-lg max-w-2xl leading-relaxed">
            Access your medical production suite. Efficiently manage catalogs, convert medical documentation, and streamline your clinical workflow with precision tools.
          </p>
        </motion.div>
      </section>

      {/* Bento Tool Grid */}
      <section className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-12">
        {/* Online Catalog Card */}
        <motion.div 
          whileHover={{ scale: 1.01 }}
          className="group lg:col-span-2 bg-surface-container-lowest p-8 rounded-xl ambient-shadow flex flex-col justify-between hover:bg-surface-container-low transition-all duration-300 overflow-hidden relative"
        >
          <div className="relative z-10">
            <div className="w-12 h-12 rounded-lg bg-primary-container text-primary flex items-center justify-center mb-6">
              <BookOpen className="w-6 h-6" />
            </div>
            <h2 className="text-2xl font-bold text-on-surface mb-2">Online Catalog</h2>
            <p className="text-on-surface-variant max-w-md">Browse and manage our extensive database of medical supplies and production components with real-time inventory tracking.</p>
          </div>
          <div className="mt-8 flex items-center gap-2 text-primary font-semibold text-sm group-hover:gap-4 transition-all">
            <Link href="/catalog">Open Catalog</Link>
            <ArrowRight className="w-4 h-4" />
          </div>
          <div className="absolute right-[-10%] bottom-[-10%] opacity-5 group-hover:opacity-10 transition-opacity">
            <BookOpen className="w-48 h-48" />
          </div>
        </motion.div>

        {/* Quick Convert Card */}
        <motion.div 
          whileHover={{ scale: 1.01 }}
          className="group bg-surface-container-lowest p-8 rounded-xl ambient-shadow flex flex-col justify-between hover:bg-surface-container-low transition-all duration-300"
        >
          <div>
            <div className="w-12 h-12 rounded-lg bg-tertiary-container text-tertiary flex items-center justify-center mb-6">
              <RefreshCw className="w-6 h-6" />
            </div>
            <h2 className="text-2xl font-bold text-on-surface mb-2">Quick Convert</h2>
            <p className="text-on-surface-variant">Rapidly transform medical units and technical specifications between international standards.</p>
          </div>
          <div className="mt-8 flex items-center gap-2 text-primary font-semibold text-sm group-hover:gap-4 transition-all">
            <Link href="/convert">Start Converting</Link>
            <ArrowRight className="w-4 h-4" />
          </div>
        </motion.div>

        {/* Extraction Card */}
        <motion.div 
          whileHover={{ scale: 1.01 }}
          className="group bg-surface-container-lowest p-8 rounded-xl ambient-shadow flex flex-col justify-between hover:bg-surface-container-low transition-all duration-300 relative overflow-hidden"
        >
          <div className="relative z-10">
            <div className="w-12 h-12 rounded-lg bg-secondary-container text-secondary flex items-center justify-center mb-6">
              <Scan className="w-6 h-6" />
            </div>
            <h2 className="text-2xl font-bold text-on-surface mb-2">Extraction</h2>
            <p className="text-on-surface-variant">Extract text and data from physical medical reports with our high-precision optical recognition engine.</p>
          </div>
          <div className="mt-8 flex items-center gap-2 text-primary font-semibold text-sm group-hover:gap-4 transition-all">
            <Link href="/extraction">Open Workbench</Link>
            <ArrowRight className="w-4 h-4" />
          </div>
          <div className="absolute top-0 right-0 p-4">
            <span className="px-3 py-1 bg-error-container/40 text-on-error-container text-[10px] font-bold rounded-full uppercase tracking-widest">New</span>
          </div>
        </motion.div>

        {/* Report Card */}
        <motion.div 
          whileHover={{ scale: 1.01 }}
            className="group lg:col-span-2 bg-primary text-on-primary p-8 rounded-xl shadow-[0_12px_40px_rgba(0,93,182,0.15)] flex items-center justify-between hover:bg-primary-dim transition-all duration-300 relative overflow-hidden"
        >
          <div className="max-w-md relative z-10">
            <div className="w-12 h-12 rounded-lg bg-white/20 flex items-center justify-center mb-6">
              <CheckSquare className="w-6 h-6" />
            </div>
            <h2 className="text-3xl font-bold mb-3">Report Studio</h2>
            <p className="text-white/80 text-lg leading-relaxed">Manage Martin reports, edit instrument lists, and export final PDF packs for your clinical documentation workflow.</p>
            <div className="mt-8 inline-flex items-center gap-4 bg-white text-primary px-6 py-3 rounded-lg font-bold text-sm hover:-translate-y-0.5 transition-transform">
              <Link href="/report">Open Report Studio</Link>
              <Rocket className="w-4 h-4" />
            </div>
          </div>
          <div className="hidden lg:block w-1/3 h-full relative">
            <div className="absolute inset-0 flex items-center justify-center">
              <div className="w-48 h-48 bg-white/10 rounded-full blur-3xl animate-pulse"></div>
            </div>
            <img 
              src="https://picsum.photos/seed/medical-lab/400/400" 
              alt="Medical Technology" 
              className="rounded-xl object-cover h-full w-full opacity-90"
              referrerPolicy="no-referrer"
            />
          </div>
        </motion.div>
      </section>

      {/* Secondary Info & Activity */}
      <section className="grid grid-cols-1 lg:grid-cols-4 gap-8">
        <div className="lg:col-span-3">
          <div className="flex items-center justify-between mb-8">
            <h3 className="text-xl font-bold tracking-tight">Recent Production Files</h3>
            <Link href="/files" className="cursor-pointer text-primary font-bold text-sm hover:underline">View All</Link>
          </div>
          <div className="space-y-4">
            <FileRow 
              name="Batch_Request_042.pdf" 
              info="Modified 2 hours ago • 4.2 MB" 
              status="Processed" 
              statusColor="bg-primary-container text-on-primary-container"
            />
            <FileRow 
              name="OCR_Scan_Lab_Results.docx" 
              info="Modified Yesterday • 1.1 MB" 
              status="Draft" 
              statusColor="bg-tertiary-container text-on-tertiary-container"
            />
          </div>
        </div>

        {/* Stats/Sidebar Info */}
        <div className="lg:col-span-1 space-y-6">
          <div className="p-6 bg-surface-container-low rounded-xl border-l-4 border-primary">
            <h4 className="text-sm font-bold text-primary uppercase tracking-widest mb-4">System Status</h4>
            <div className="flex items-center gap-3 mb-4">
              <div className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></div>
              <span className="text-sm font-medium">Production Node: Active</span>
            </div>
            <div className="flex items-center gap-3">
              <div className="w-2 h-2 rounded-full bg-emerald-500"></div>
              <span className="text-sm font-medium">Database: Online</span>
            </div>
          </div>
          
          <div className="p-6 bg-secondary-container/20 rounded-xl">
            <h4 className="text-sm font-bold text-on-secondary-container uppercase tracking-widest mb-4">Storage Usage</h4>
            <div className="w-full bg-slate-200 h-2 rounded-full mb-2">
              <div className="bg-primary h-full rounded-full w-[65%]"></div>
            </div>
            <p className="text-xs text-on-surface-variant font-medium">65.2 GB of 100 GB used</p>
          </div>
        </div>
      </section>
    </div>
  );
}

function FileRow({ name, info, status, statusColor }: { name: string, info: string, status: string, statusColor: string }) {
  return (
    <div className="flex items-center justify-between p-4 bg-white rounded-lg hover:bg-surface-container-low transition-colors group ambient-shadow">
      <div className="flex items-center gap-4">
        <div className="w-10 h-10 rounded-lg bg-surface-container-high flex items-center justify-center">
          <FileText className="text-on-surface-variant w-5 h-5" />
        </div>
        <div>
          <p className="font-bold text-on-surface">{name}</p>
          <p className="text-xs text-on-surface-variant">{info}</p>
        </div>
      </div>
      <div className="flex items-center gap-4">
        <span className={cn("px-3 py-1 text-[10px] font-bold rounded-full uppercase tracking-tighter", statusColor)}>
          {status}
        </span>
        <button className="cursor-pointer opacity-0 group-hover:opacity-100 p-2 hover:bg-white rounded-full transition-all">
          <MoreVertical className="w-4 h-4 text-on-surface-variant" />
        </button>
      </div>
    </div>
  );
}
