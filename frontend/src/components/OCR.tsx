import React from 'react';
import { Scan, Upload, FileText, Download, Copy, CheckCircle, Lock, EyeOff, HelpCircle } from 'lucide-react';
import { cn } from '@/lib/utils';

export function OCR() {
  return (
    <div className="max-w-350 mx-auto">
      <header className="mb-12">
        <h1 className="text-4xl font-extrabold text-on-surface tracking-tight mb-2">Optical Character Recognition</h1>
        <p className="text-on-surface-variant max-w-2xl leading-relaxed">Precision-engineered extraction for clinical documentation. Transform medical forms, lab results, and patient charts into structured, actionable data.</p>
      </header>

      <div className="grid grid-cols-12 gap-8">
        {/* Left Column: Upload Area */}
        <div className="col-span-12 lg:col-span-7 space-y-8">
          {/* Dropzone Card */}
          <div className="bg-surface-container-lowest rounded-xl p-1 ambient-shadow group">
            <div className="border-2 border-dashed border-outline-variant/30 rounded-lg p-16 flex flex-col items-center justify-center transition-all duration-300 group-hover:bg-primary-container/10 group-hover:border-primary/40">
              <div className="w-16 h-16 bg-primary-container rounded-full flex items-center justify-center mb-6 text-primary">
                <Upload className="w-8 h-8" />
              </div>
              <h3 className="text-xl font-bold text-on-surface mb-2">Upload clinical document</h3>
              <p className="text-on-surface-variant text-center mb-8">Drag and drop your medical PDF or high-resolution image here.<br/>Maximum file size 25MB.</p>
              <button className="bg-primary text-on-primary px-8 py-3 rounded-lg font-semibold flex items-center gap-2 hover:bg-primary-dim hover:-translate-y-0.5 transition-all active:scale-95">
                Browse Files
              </button>
            </div>
          </div>

          {/* Export Options */}
          <div className="bg-surface-container-low rounded-xl p-8">
            <h4 className="text-headline-sm uppercase tracking-widest text-xs font-bold text-on-surface-variant mb-6">Select Export Format</h4>
            <div className="grid grid-cols-3 gap-4">
              <FormatOption icon={<FileText className="w-5 h-5" />} label="Table" active />
              <FormatOption icon={<FileText className="w-5 h-5" />} label="CSV" />
              <FormatOption icon={<FileText className="w-5 h-5" />} label="Raw Text" />
            </div>
          </div>

          {/* Action Button */}
          <button className="w-full bg-primary py-5 rounded-xl text-on-primary font-bold text-lg shadow-lg hover:shadow-primary/20 hover:-translate-y-1 transition-all flex items-center justify-center gap-3">
            <Scan className="w-6 h-6" />
            Perform OCR Extraction
          </button>
        </div>

        {/* Right Column: Context & Preview */}
        <div className="col-span-12 lg:col-span-5 space-y-8">
          <div className="grid grid-cols-1 gap-6">
            <InfoCard 
              icon={<CheckCircle className="w-5 h-5" />} 
              title="Clinical Grade Accuracy" 
              desc="Our neural networks are specifically trained on medical terminology and handwritten physician notes."
              color="bg-blue-100 text-blue-700"
            />
            <InfoCard 
              icon={<Lock className="w-5 h-5" />} 
              title="HIPAA Compliant" 
              desc="Data is processed in an isolated environment and automatically purged after 15 minutes of inactivity."
              color="bg-green-100 text-green-700"
            />

            {/* Source Preview Placeholder */}
            <div className="bg-surface-container-lowest rounded-xl p-4 shadow-sm overflow-hidden border border-outline-variant/10">
              <h4 className="text-xs font-bold text-on-surface-variant uppercase tracking-wider mb-4 px-2">Recent Source Preview</h4>
              <div className="aspect-3/4 bg-surface rounded-lg relative overflow-hidden">
                <img 
                  src="https://picsum.photos/seed/medical-form/400/600" 
                  alt="Sample Medical Document" 
                  className="w-full h-full object-cover opacity-50 grayscale"
                  referrerPolicy="no-referrer"
                />
                <div className="absolute inset-0 flex flex-col items-center justify-center bg-white/40 backdrop-blur-[2px]">
                  <EyeOff className="text-slate-400 w-10 h-10 mb-2" />
                  <span className="text-xs font-bold text-slate-500">Awaiting Upload</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Extraction Results Section */}
      <section className="mt-20">
        <div className="flex items-center justify-between mb-8">
          <div>
            <h2 className="text-2xl font-bold text-on-surface">Extraction Results</h2>
            <p className="text-on-surface-variant">Validated data from lab_report_042.pdf</p>
          </div>
          <div className="flex gap-3">
            <button className="bg-surface-container-high px-4 py-2 rounded-lg text-sm font-semibold hover:bg-surface-container-highest transition-colors flex items-center gap-2">
              <Download className="w-4 h-4" /> Download File
            </button>
            <button className="bg-surface-container-high px-4 py-2 rounded-lg text-sm font-semibold hover:bg-surface-container-highest transition-colors flex items-center gap-2">
              <Copy className="w-4 h-4" /> Copy to Clipboard
            </button>
          </div>
        </div>
        
        <div className="bg-surface-container-lowest rounded-xl overflow-hidden shadow-sm">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="bg-surface-container-low border-b-0">
                <th className="px-8 py-4 font-bold text-xs uppercase tracking-wider text-on-surface-variant">Parameter</th>
                <th className="px-8 py-4 font-bold text-xs uppercase tracking-wider text-on-surface-variant">Extracted Value</th>
                <th className="px-8 py-4 font-bold text-xs uppercase tracking-wider text-on-surface-variant">Reference Range</th>
                <th className="px-8 py-4 font-bold text-xs uppercase tracking-wider text-on-surface-variant">Confidence</th>
              </tr>
            </thead>
            <tbody className="divide-y-0">
              <ResultRow parameter="Hemoglobin A1c" value="5.7 %" range="4.8 - 5.6" confidence="98% Match" confidenceColor="bg-green-100/40 text-green-700" />
              <ResultRow parameter="Glucose, Fasting" value="94 mg/dL" range="65 - 99" confidence="100% Match" confidenceColor="bg-green-100/40 text-green-700" />
              <ResultRow parameter="Creatinine" value="0.88 mg/dL" range="0.60 - 1.30" confidence="Check Manual" confidenceColor="bg-error-container/40 text-error" />
            </tbody>
          </table>
        </div>
      </section>

      <button className="fixed bottom-8 right-8 bg-primary text-on-primary w-14 h-14 rounded-full shadow-2xl flex items-center justify-center hover:scale-110 active:scale-95 transition-all group">
        <HelpCircle className="w-6 h-6" />
      </button>
    </div>
  );
}

function FormatOption({ icon, label, active = false }: { icon: React.ReactNode, label: string, active?: boolean }) {
  return (
    <label className={cn(
      "relative flex flex-col items-center p-4 bg-surface-container-lowest rounded-lg cursor-pointer hover:bg-white transition-all shadow-sm ring-2 ring-transparent",
      active && "ring-primary"
    )}>
      <input type="radio" name="format" className="absolute opacity-0" defaultChecked={active} />
      <div className={cn("mb-2", active ? "text-primary" : "text-secondary")}>
        {icon}
      </div>
      <span className="font-bold text-sm">{label}</span>
    </label>
  );
}

function InfoCard({ icon, title, desc, color }: { icon: React.ReactNode, title: string, desc: string, color: string }) {
  return (
    <div className="bg-surface-container-low rounded-xl p-6">
      <div className="flex items-start gap-4">
        <div className={cn("p-2 rounded-lg", color)}>
          {icon}
        </div>
        <div>
          <h4 className="font-bold text-on-surface mb-1">{title}</h4>
          <p className="text-sm text-on-surface-variant leading-relaxed">{desc}</p>
        </div>
      </div>
    </div>
  );
}

function ResultRow({ parameter, value, range, confidence, confidenceColor }: { parameter: string, value: string, range: string, confidence: string, confidenceColor: string }) {
  return (
    <tr className="hover:bg-surface-container-low transition-colors">
      <td className="px-8 py-6 font-semibold text-on-surface">{parameter}</td>
      <td className="px-8 py-6 text-on-surface-variant">{value}</td>
      <td className="px-8 py-6 text-on-surface-variant">{range}</td>
      <td className="px-8 py-6">
        <span className={cn("px-3 py-1 rounded-full text-xs font-bold", confidenceColor)}>
          {confidence}
        </span>
      </td>
    </tr>
  );
}
