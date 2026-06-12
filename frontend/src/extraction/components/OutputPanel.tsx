import { ChevronLeft, ChevronRight } from "lucide-react";

type OutputPanelProps = {
  activePage: number;
};

const extractedText = `CLINICAL ANALYSIS REPORT #42-B
DATE: OCTOBER 24, 2023
FACILITY: NORTH STAR DIAGNOSTIC HUB

PATIENT INFORMATION:
Name: [REDACTED]
Identifier: CALO-992-PX
Status: EVALUATION COMPLETE

DIAGNOSTIC SUMMARY:
The preliminary results indicate a consistent level of biomarker density in the upper quartiles. Compared to baseline report_031.pdf, there is a visible delta of +12% in operational efficiency.

OBSERVATIONS:
1. Neutrophil levels remain within the clinical threshold.
2. Platelet adhesion patterns show high-precision alignment.
3. No immediate secondary intervention required at this stage.

NOTES:
Review scheduled for follow-up in 14 days. Staff are advised to maintain current protocols and log all incremental shifts in the eSelector database.`;

export function OutputPanel({ activePage }: OutputPanelProps) {
  return (
    <section className="min-h-0 flex flex-col bg-surface-container-lowest overflow-hidden">
      <div className="p-4 bg-surface border-b border-outline-variant/10 flex flex-col gap-3 shrink-0 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex flex-col">
          <span className="text-xs font-bold text-on-surface-variant uppercase tracking-wider">
            Editable Output
          </span>
          <span className="text-sm font-semibold">Extracted Text</span>
        </div>
        <div className="flex items-center gap-3">
          <span className="text-xs font-medium text-on-surface-variant bg-surface-container-high px-2 py-1 rounded">
            Page {activePage} of 12
          </span>
          <div className="flex border border-outline-variant/20 rounded overflow-hidden">
            <button
              type="button"
              aria-label="Previous page"
              className="p-1 hover:bg-surface-container-high border-r border-outline-variant/20 transition-colors"
            >
              <ChevronLeft className="w-4 h-4 text-on-surface-variant" />
            </button>
            <button
              type="button"
              aria-label="Next page"
              className="p-1 hover:bg-surface-container-high transition-colors"
            >
              <ChevronRight className="w-4 h-4 text-on-surface-variant" />
            </button>
          </div>
        </div>
      </div>

      <div className="grow overflow-y-auto p-6 bg-surface/50 custom-scrollbar lg:p-12">
        <div className="max-w-2xl mx-auto space-y-12">
          <div className="relative group">
            <div className="absolute -left-8 top-0 h-full border-l-2 border-primary/20 group-focus-within:border-primary transition-colors" />
            <div className="mb-4 flex items-center gap-4">
              <div className="h-px grow bg-outline-variant/20" />
              <span className="text-[10px] font-bold text-primary uppercase tracking-[0.2em] bg-surface-container-lowest px-4">
                Extraction Page 01
              </span>
              <div className="h-px grow bg-outline-variant/20" />
            </div>
            <textarea
              className="w-full min-h-[25rem] border-none focus:ring-0 bg-transparent text-on-surface leading-loose font-body text-base resize-none outline-none"
              spellCheck="false"
              defaultValue={extractedText}
            />
          </div>

          <div className="relative group">
            <div className="absolute -left-8 top-0 h-full border-l-2 border-primary/20 group-focus-within:border-primary transition-colors" />
            <div className="mb-4 flex items-center gap-4">
              <div className="h-px grow bg-outline-variant/20" />
              <span className="text-[10px] font-bold text-on-surface-variant uppercase tracking-[0.2em] bg-surface-container-lowest px-4">
                Extraction Page 02
              </span>
              <div className="h-px grow bg-outline-variant/20" />
            </div>
            <textarea
              className="w-full min-h-[12.5rem] border-none focus:ring-0 bg-transparent text-on-surface/50 leading-loose font-body text-base resize-none italic outline-none"
              readOnly
              defaultValue="Loading extracted content for page 2...\nPlease select page in the left viewer to prioritize processing."
            />
          </div>
        </div>
      </div>
    </section>
  );
}
