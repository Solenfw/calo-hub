import React from 'react';
import { CheckSquare, Search, Sliders, Heart, ShoppingCart, Download, FlaskConical as Science, Verified } from 'lucide-react';
import { cn } from '@/lib/utils';

export function ESelector() {
  return (
    <div className="max-w-6xl mx-auto">
      {/* Hero Section / Context */}
      <section className="mb-12">
        <div className="flex justify-between items-end">
          <div>
            <h2 className="text-4xl font-extrabold text-on-surface tracking-tight mb-2">eSelector Tool</h2>
            <p className="text-on-surface-variant max-w-xl">Configure specific requirements to filter our precision manufacturing database for compatible medical-grade components.</p>
          </div>
          <div className="text-right">
            <span className="text-[10px] font-bold text-primary tracking-widest uppercase">System Status</span>
            <div className="flex items-center gap-2 text-on-surface font-medium">
              <span className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
              Database v4.2.0 Connected
            </div>
          </div>
        </div>
      </section>

      <div className="grid grid-cols-12 gap-8">
        {/* Left Column: Input Requirements */}
        <div className="col-span-12 lg:col-span-5 flex flex-col gap-6">
          <div className="bg-surface-container-lowest p-8 rounded-xl ambient-shadow">
            <h3 className="text-lg font-bold mb-6 flex items-center gap-2">
              <Sliders className="text-primary w-5 h-5" />
              Requirement Configuration
            </h3>
            
            <div className="space-y-6">
              <SelectInput label="Material Type" options={['Titanium Grade 5 (ELI)', 'Stainless Steel 316L', 'Cobalt Chrome', 'PEEK-OPTIMA']} />
              <div className="space-y-2">
                <label className="text-xs font-bold text-outline-variant uppercase tracking-wider block">Total Length (mm)</label>
                <input className="w-full bg-surface-container-high border-none rounded-lg py-3 px-4 focus:ring-0 focus:border-b-2 focus:border-primary transition-all" placeholder="e.g. 120.5" type="text" />
              </div>
              <SelectInput label="Application Type" options={['Orthopedic Implants', 'Cardiovascular Support', 'Dental Abutments', 'Surgical Tooling']} />
              
              <div className="space-y-2">
                <label className="text-xs font-bold text-outline-variant uppercase tracking-wider block">Precision Tolerance</label>
                <div className="grid grid-cols-2 gap-4">
                  <ToleranceOption label="Standard (±0.05)" checked />
                  <ToleranceOption label="Ultra (±0.001)" />
                </div>
              </div>

              <button className="cursor-pointer w-full mt-4 bg-linear-to-br from-primary to-primary-dim text-on-primary font-bold py-4 rounded-lg flex items-center justify-center gap-3 shadow-lg shadow-primary/20 hover:-translate-y-0.5 active:translate-y-0 transition-all duration-300">
                Find Matching Codes
                <Search className="w-4 h-4" />
              </button>
            </div>
          </div>

          <div className="bg-surface-container-low p-6 rounded-xl relative overflow-hidden">
            <div className="relative z-10">
              <h4 className="font-bold text-on-surface-variant mb-2">Technical Guidance</h4>
              <p className="text-sm text-on-surface-variant/80">Selected codes are automatically verified against ISO 13485:2016 standards for medical device manufacturing.</p>
            </div>
            <Verified className="absolute -right-4 -bottom-4 w-24 h-24 text-primary/5 rotate-12" />
          </div>
        </div>

        {/* Right Column: Matching Results */}
        <div className="col-span-12 lg:col-span-7 flex flex-col gap-6">
          <div className="flex items-center justify-between mb-2">
            <h3 className="text-lg font-bold">Search Results <span className="ml-2 text-sm font-medium text-outline-variant bg-surface-container px-2 py-0.5 rounded-full">3 matches found</span></h3>
          </div>

          <div className="space-y-4">
            <ProductCard 
              code="CALO-TI-9082X" 
              name="Cortical Compression Bolt" 
              match="98%" 
              length="120.5 mm" 
              material="Ti-6Al-4V" 
              thread="M4 x 0.7" 
              image="https://picsum.photos/seed/bolt/200/200"
            />
            <ProductCard 
              code="CALO-SS-1142V" 
              name="Cannulated Screw (Standard)" 
              match="92%" 
              length="120.0 mm" 
              material="SS 316L" 
              thread="M3.5 x 0.6" 
              image="https://picsum.photos/seed/screw/200/200"
            />
            <ProductCard 
              code="CALO-TI-8874P" 
              name="Locking Plate Interface" 
              match="85%" 
              length="121.2 mm" 
              material="Ti Grade 2" 
              thread="N/A" 
              image="https://picsum.photos/seed/plate/200/200"
            />
          </div>

          <div className="mt-4 flex items-center justify-end gap-4">
            <button className="cursor-pointer flex items-center gap-2 px-6 py-2.5 rounded-lg border border-outline-variant text-on-surface-variant font-bold hover:bg-surface-container-low transition-colors">
              <Download className="w-4 h-4" />
              Export PDF Spec Sheet
            </button>
            <button className="cursor-pointer flex items-center gap-2 px-6 py-2.5 rounded-lg bg-secondary text-on-secondary font-bold hover:opacity-90 transition-opacity">
              <Science className="w-4 h-4" />
              Request Sample Kit
            </button>
          </div>
        </div>
      </div>

      <footer className="fixed bottom-4 right-8 flex gap-6 text-[10px] font-bold text-outline uppercase tracking-widest pointer-events-none opacity-40">
        <span>ISO 13485 CERTIFIED</span>
        <span>CLEANROOM CLASS 100</span>
        <span>FDA REGISTERED</span>
      </footer>
    </div>
  );
}

function SelectInput({ label, options }: { label: string, options: string[] }) {
  return (
    <div className="space-y-2">
      <label className="text-xs font-bold text-outline-variant uppercase tracking-wider block">{label}</label>
      <div className="relative">
        <select className="w-full bg-surface-container-high border-none rounded-lg py-3 px-4 appearance-none focus:ring-0 focus:border-b-2 focus:border-primary transition-all">
          {options.map(opt => <option key={opt}>{opt}</option>)}
        </select>
        <Sliders className="absolute right-3 top-1/2 -translate-y-1/2 pointer-events-none text-outline w-4 h-4" />
      </div>
    </div>
  );
}

function ToleranceOption({ label, checked = false }: { label: string, checked?: boolean }) {
  return (
    <label className="flex items-center gap-3 p-3 bg-surface rounded-lg cursor-pointer hover:bg-primary-container/20 transition-colors group">
      <input type="radio" name="tolerance" className="text-primary focus:ring-primary" defaultChecked={checked} />
      <span className="text-sm font-medium text-on-surface">{label}</span>
    </label>
  );
}

function ProductCard({ code, name, match, length, material, thread, image }: any) {
  return (
    <div className="group bg-surface-container-lowest p-1 rounded-xl ambient-shadow hover:shadow-md transition-shadow duration-300">
      <div className="flex items-center gap-6 p-5">
        <div className="w-24 h-24 bg-surface rounded-lg flex items-center justify-center overflow-hidden">
          <img src={image} alt={name} className="w-full h-full object-cover" referrerPolicy="no-referrer" />
        </div>
        <div className="flex-1">
          <div className="flex justify-between items-start">
            <div>
              <div className="flex items-center gap-2 mb-1">
                <span className="text-[10px] font-black bg-primary/10 text-primary px-2 py-0.5 rounded-full uppercase tracking-tighter">Code</span>
                <span className="text-lg font-extrabold text-on-surface font-manrope">{code}</span>
              </div>
              <h4 className="font-medium text-on-surface-variant text-sm">{name}</h4>
            </div>
            <div className="bg-primary-container/40 px-3 py-1.5 rounded-full flex items-center gap-2">
              <span className="text-xs font-bold text-primary">Match: {match}</span>
            </div>
          </div>
          <div className="mt-4 flex gap-6">
            <StatItem label="Length" value={length} />
            <StatItem label="Material" value={material} />
            <StatItem label="Thread" value={thread} />
          </div>
        </div>
        <div className="flex flex-col gap-2">
          <button className="cursor-pointer p-2 rounded-lg hover:bg-surface transition-colors"><Heart className="w-4 h-4 text-outline" /></button>
          <button className="cursor-pointer p-2 rounded-lg bg-primary/5 hover:bg-primary/10 text-primary transition-colors"><ShoppingCart className="w-4 h-4" /></button>
        </div>
      </div>
    </div>
  );
}

function StatItem({ label, value }: { label: string, value: string }) {
  return (
    <div className="flex flex-col">
      <span className="text-[10px] text-outline-variant font-bold uppercase">{label}</span>
      <span className="text-sm font-semibold">{value}</span>
    </div>
  );
}
