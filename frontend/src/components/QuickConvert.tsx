import React from 'react';
import { RefreshCw, Upload, Info, Package as Inventory, ShoppingCart } from 'lucide-react';
import { cn } from '@/lib/utils';

export function QuickConvert() {
  return (
    <div className="max-w-4xl mx-auto">
      <section className="mb-12 text-center">
        <h2 className="text-4xl font-extrabold text-on-surface tracking-tight mb-3">Quick Convert</h2>
        <p className="text-on-surface-variant max-w-lg mx-auto">Instantly map medical product codes across different manufacturers with clinical precision and material matching.</p>
      </section>

      <div className="bg-surface-container-lowest rounded-xl p-10 ambient-shadow mb-8">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-8 mb-8">
          <div className="space-y-2">
            <label className="text-xs font-bold uppercase tracking-widest text-outline pl-1">Source Brand</label>
            <div className="relative">
              <select className="w-full bg-surface-container-high border-none rounded-lg py-4 px-4 appearance-none focus:ring-2 focus:ring-primary/20 text-on-surface transition-all">
                <option>Calo Medical Systems</option>
                <option>Heraeus Kulzer</option>
                <option>Straumann Group</option>
                <option>Dentsply Sirona</option>
              </select>
              <RefreshCw className="absolute right-4 top-4 pointer-events-none text-outline w-4 h-4" />
            </div>
          </div>
          <div className="space-y-2">
            <label className="text-xs font-bold uppercase tracking-widest text-outline pl-1">Target Brand</label>
            <div className="relative">
              <select className="w-full bg-surface-container-high border-none rounded-lg py-4 px-4 appearance-none focus:ring-2 focus:ring-primary/20 text-on-surface transition-all">
                <option>Global Bio-Tech</option>
                <option>Nobel Biocare</option>
                <option>Zimmer Biomet</option>
                <option>BioHorizons</option>
              </select>
              <RefreshCw className="absolute right-4 top-4 pointer-events-none text-outline w-4 h-4" />
            </div>
          </div>
        </div>

        <div className="space-y-2 mb-10">
          <label className="text-xs font-bold uppercase tracking-widest text-outline pl-1">Product Code</label>
          <input 
            className="w-full text-2xl font-medium bg-surface-container-high border-none border-b-2 border-transparent focus:border-primary rounded-lg py-6 px-6 focus:ring-0 placeholder:text-outline/50 text-on-surface transition-all" 
            placeholder="Paste product code (e.g., CMS-4092-TX)..." 
            type="text"
          />
        </div>

        <div className="flex flex-col md:flex-row gap-4 items-center justify-center">
          <button className="cursor-pointer bg-linear-to-br from-primary to-primary-dim text-on-primary px-12 py-4 rounded-lg font-bold text-lg shadow-lg hover:-translate-y-0.5 active:scale-95 transition-all flex items-center gap-3 w-full md:w-auto">
            <RefreshCw className="w-5 h-5" />
            Convert
          </button>
          <button className="cursor-pointer bg-secondary-container text-on-secondary-container px-8 py-4 rounded-lg font-semibold hover:bg-secondary-container/80 transition-all flex items-center gap-3 w-full md:w-auto justify-center">
            <Upload className="w-5 h-5" />
            Import CSV
          </button>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="md:col-span-2 bg-surface-container-lowest rounded-xl p-8 border border-primary/5 ambient-shadow">
          <div className="flex items-start justify-between mb-8">
            <div>
              <span className="text-[10px] font-bold uppercase tracking-widest text-primary mb-1 block">Conversion Result</span>
              <h3 className="text-3xl font-bold text-on-surface">GBT-X900-V2</h3>
            </div>
            <div className="bg-primary-container/40 px-3 py-1 rounded-full flex items-center gap-2">
              <span className="w-2 h-2 rounded-full bg-primary animate-pulse"></span>
              <span className="text-xs font-bold text-primary">98% Match</span>
            </div>
          </div>
          
          <div className="space-y-6">
            <ConvertRow label="Standard Length" value="14.5mm (±0.01)" />
            <ConvertRow label="Material Base" value="Grade 5 Titanium (Ti6Al4V ELI)" />
            <ConvertRow label="Implant Type" value="Conical Connection (Internal)" />
          </div>
        </div>

        <div className="space-y-6">
          <div className="bg-primary/5 rounded-xl p-6 relative overflow-hidden group">
            <div className="absolute top-0 right-0 w-32 h-32 bg-primary/10 rounded-full -mr-16 -mt-16 group-hover:scale-110 transition-transform duration-500"></div>
            <h4 className="text-sm font-bold text-primary mb-4 flex items-center gap-2">
              <Info className="w-4 h-4" />
              Technical Note
            </h4>
            <p className="text-xs leading-relaxed text-on-surface-variant">Target brand uses a proprietary surface blasting technique (SLA equivalent) which may differ from the source's acid-etched finish.</p>
          </div>

          <div className="bg-surface-container-high rounded-xl p-6">
            <h4 className="text-sm font-bold text-on-surface mb-4">Stock Availability</h4>
            <div className="flex items-center gap-4 mb-4">
              <div className="w-10 h-10 rounded bg-white flex items-center justify-center shadow-sm">
                <RefreshCw className="w-5 h-5 text-primary" />
              </div>
              <div>
                <p className="text-xs text-on-surface-variant">Central Warehouse</p>
                <p className="text-sm font-bold text-on-surface">1,240 Units</p>
              </div>
            </div>
            <button className="cursor-pointer w-full py-2 bg-on-surface text-surface rounded-lg text-xs font-bold hover:bg-on-surface/90 transition-all">Add to Cart</button>
          </div>
        </div>
      </div>

      <footer className="mt-20 py-8 text-center text-on-surface-variant text-xs opacity-50 border-t border-surface-container-high">
        <p>© 2024 Calo Hub Medical Production Systems. All conversions verified by AI-Assisted Clinical Validation.</p>
      </footer>
    </div>
  );
}

function ConvertRow({ label, value }: { label: string, value: string }) {
  return (
    <div className="flex items-center justify-between py-4 border-b border-surface-container last:border-0">
      <div className="flex items-center gap-3">
        <span className="text-on-surface-variant font-medium">{label}</span>
      </div>
      <span className="font-bold text-on-surface">{value}</span>
    </div>
  );
}
