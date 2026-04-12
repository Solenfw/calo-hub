import React from 'react';
import { Search, Filter, ChevronLeft, ChevronRight, Download, Copy, Info, Plus } from 'lucide-react';
import { cn } from '@/lib/utils';

const catalogItems = [
  { code: 'KL-10-244-12', name: 'Mayo-Hegar Needle Holder', desc: '150mm, Tungsten Carbide Inserts' },
  { code: 'BB-77-901-00', name: 'Metzenbaum Scissors', desc: 'Curved, Standard Pattern, 14.5cm' },
  { code: 'ST-04-112-99', name: 'Bone Rongeur Friedman', desc: 'Slightly Curved, 140mm' },
  { code: 'KL-44-102-15', name: 'Adson Dressing Forceps', desc: 'Serrated, Straight, 120mm' },
  { code: 'BB-12-555-08', name: 'Scalpel Handle No. 4', desc: 'Solid, Stainless Steel' },
  { code: 'ST-99-001-22', name: 'Kelly Hemostatic Forceps', desc: 'Straight, 140mm' },
];

export function Catalog() {
  return (
    <div className="max-w-6xl mx-auto">
      {/* Page Header */}
      <div className="mb-12">
        <h1 className="font-manrope text-4xl font-extrabold tracking-tight text-on-surface mb-2">Online Catalog</h1>
        <p className="text-on-surface-variant font-body">Browse and search through the comprehensive clinical production database.</p>
      </div>

      {/* Filter & Search Toolbar */}
      <div className="bg-surface-container-lowest rounded-xl p-6 ambient-shadow mb-8 flex flex-col md:flex-row gap-6 items-end">
        <div className="flex-1 w-full">
          <label className="block text-[11px] font-bold uppercase tracking-widest text-outline-variant mb-3 ml-1">Search Products</label>
          <div className="relative group">
            <Search className="absolute left-4 top-1/2 -translate-y-1/2 text-outline-variant group-focus-within:text-primary transition-colors w-4 h-4" />
            <input 
              type="text" 
              className="w-full bg-surface-container-high border-none rounded-lg pl-12 pr-4 py-4 focus:ring-0 focus:border-b-2 focus:border-primary transition-all placeholder:text-outline-variant text-on-surface"
              placeholder="Enter product name or code (e.g. Scalpel, KN-400)"
            />
          </div>
        </div>
        <div className="w-full md:w-64">
          <label className="block text-[11px] font-bold uppercase tracking-widest text-outline-variant mb-3 ml-1">Brand Filter</label>
          <div className="relative">
            <select className="w-full appearance-none bg-surface-container-high border-none rounded-lg px-4 py-4 focus:ring-0 focus:border-b-2 focus:border-primary transition-all text-on-surface cursor-pointer">
              <option>All Brands</option>
              <option>KLS Martin</option>
              <option>B-Braun</option>
              <option>Stema</option>
            </select>
            <ChevronLeft className="absolute right-4 top-1/2 -translate-y-1/2 pointer-events-none text-outline-variant w-4 h-4 rotate-270" />
          </div>
        </div>
        <button className="cursor-pointer bg-primary text-on-primary px-8 py-4 rounded-lg font-bold flex items-center gap-2 hover:bg-primary-dim transition-all shadow-lg shadow-primary/10 active:scale-95">
          <Filter className="w-4 h-4" />
          Search
        </button>
      </div>

      {/* Grid Layout for Catalog Content */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-8">
        {/* Main Catalog Table */}
        <div className="lg:col-span-8 bg-surface-container-lowest rounded-xl overflow-hidden ambient-shadow">
          <div className="p-8 pb-0">
            <div className="flex justify-between items-center mb-6">
              <h3 className="font-manrope text-xl font-bold text-black">Catalog Entries</h3>
              <span className="text-[12px] font-bold text-primary bg-primary-container px-3 py-1 rounded-full">142 Results Found</span>
            </div>
          </div>
          <div className="overflow-x-auto">
            <table className="w-full text-left border-collapse">
              <thead>
                <tr className="bg-surface-container-low">
                  <th className="px-8 py-4 text-[11px] font-extrabold uppercase tracking-widest text-outline">The Code</th>
                  <th className="px-8 py-4 text-[11px] font-extrabold uppercase tracking-widest text-outline">The Product's Description</th>
                </tr>
              </thead>
              <tbody className="divide-y-0">
                {catalogItems.map((item) => (
                  <tr key={item.code} className="group hover:bg-surface-container-low transition-colors cursor-pointer">
                    <td className="px-8 py-6">
                      <span className="font-mono text-sm font-bold text-primary">{item.code}</span>
                    </td>
                    <td className="px-8 py-6">
                      <div className="flex flex-col">
                        <span className="font-medium text-on-surface">{item.name}</span>
                        <span className="text-xs text-on-surface-variant mt-1">{item.desc}</span>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          <div className="p-8 bg-surface-container-low/50 flex justify-center">
            <nav className="flex items-center gap-2">
              <button className="cursor-pointer w-10 h-10 flex items-center justify-center rounded-lg hover:bg-white transition-colors"><ChevronLeft className="w-4 h-4" /></button>
              <button className="cursor-pointer w-10 h-10 flex items-center justify-center rounded-lg bg-primary text-white font-bold">1</button>
              <button className="cursor-pointer w-10 h-10 flex items-center justify-center rounded-lg hover:bg-white transition-colors">2</button>
              <button className="cursor-pointer w-10 h-10 flex items-center justify-center rounded-lg hover:bg-white transition-colors">3</button>
              <button className="cursor-pointer w-10 h-10 flex items-center justify-center rounded-lg hover:bg-white transition-colors"><ChevronRight className="w-4 h-4" /></button>
            </nav>
          </div>
        </div>

        {/* Detail Side Column */}
        <div className="lg:col-span-4 flex flex-col gap-8">
          <div className="bg-surface-container-lowest p-8 rounded-xl ambient-shadow">
            <div className="aspect-square rounded-lg mb-6 overflow-hidden bg-surface-container-high">
              <img 
                src="https://picsum.photos/seed/forceps/400/400" 
                alt="Surgical Instrument Detail" 
                className="w-full h-full object-cover"
                referrerPolicy="no-referrer"
              />
            </div>
            <div className="flex items-center gap-2 mb-3">
              <span className="px-2 py-0.5 rounded bg-error-container/20 text-error font-bold text-[10px] uppercase">KLS Martin</span>
              <span className="px-2 py-0.5 rounded bg-on-tertiary-container/10 text-on-tertiary-container font-bold text-[10px] uppercase">Premium Class</span>
            </div>
            <h4 className="font-manrope text-2xl font-bold mb-2">Needle Holder</h4>
            <p className="text-on-surface-variant text-sm mb-6 leading-relaxed">Precision-engineered Mayo-Hegar pattern featuring gold-plated handles to signify tungsten carbide inserts for superior grip and durability.</p>
            <div className="space-y-4">
              <DetailRow label="Material" value="Hardened Steel" />
              <DetailRow label="Origin" value="Germany" />
              <DetailRow label="Sterility" value="Non-Sterile" />
            </div>
            <button className="cursor-pointer w-full mt-8 bg-secondary-container text-on-secondary-container py-4 rounded-lg font-bold hover:bg-secondary-fixed-dim transition-all active:scale-95">
              Request Technical Sheet
            </button>
          </div>

          <div className="bg-primary-container p-8 rounded-xl relative overflow-hidden group">
            <div className="relative z-10">
              <h5 className="text-primary font-manrope font-bold text-lg mb-2">Direct Integration</h5>
              <p className="text-primary/70 text-sm mb-6">Seamlessly export these codes directly to your medical inventory management system.</p>
              <button className="cursor-pointer inline-flex items-center gap-2 text-primary font-bold text-sm">
                Learn More
                <ArrowRight className="w-4 h-4" />
              </button>
            </div>
          </div>
        </div>
      </div>

      {/* FAB Action */}
      <div className="fixed bottom-8 right-8">
        <button className="cursor-pointer w-14 h-14 bg-primary text-on-primary rounded-full shadow-2xl flex items-center justify-center active:scale-95 transition-transform">
          <Plus className="w-6 h-6" />
        </button>
      </div>
    </div>
  );
}

function DetailRow({ label, value }: { label: string, value: string }) {
  return (
    <div className="flex justify-between items-center text-sm py-3 border-b border-outline-variant/10">
      <span className="text-outline">{label}</span>
      <span className="font-bold">{value}</span>
    </div>
  );
}

function ArrowRight(props: any) {
  return <ChevronRight {...props} />;
}
