import { ArrowRight, Copy, Download } from "lucide-react";

function ActionButton({
  icon,
  label,
}: {
  icon: React.ReactNode;
  label: string;
}) {
  return (
    <button
      type="button"
      className="flex items-center justify-center gap-2 px-4 py-2 border border-outline-variant/50 rounded-lg text-sm font-semibold hover:bg-surface-container-high transition-colors text-on-surface"
    >
      {icon}
      {label}
    </button>
  );
}

export function ExtractionBottomBar() {
  return (
    <footer className="border-t border-outline-variant/10 bg-surface/90 backdrop-blur-md px-4 py-4 shadow-[0_-4px_20px_rgba(0,0,0,0.02)] lg:px-6">
      <div className="flex flex-col gap-4 xl:flex-row xl:items-center xl:justify-between">
        <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
          <ActionButton
            icon={<Download className="w-4 h-4" />}
            label="Download File"
          />
          <ActionButton
            icon={<Copy className="w-4 h-4" />}
            label="Copy to Clipboard"
          />
        </div>

        <div className="flex items-center justify-center gap-3 bg-surface-container-high px-3 py-1.5 rounded-full border border-outline-variant/20">
          <label
            htmlFor="extraction-page-jump"
            className="text-xs font-bold text-on-surface-variant"
          >
            Go to page
          </label>
          <input
            id="extraction-page-jump"
            className="w-10 bg-surface-container-lowest border-none outline-none text-center text-sm font-bold rounded py-0.5 focus:ring-1 focus:ring-primary"
            placeholder="1"
            type="text"
          />
          <button
            type="button"
            aria-label="Go to page"
            className="w-7 h-7 flex items-center justify-center bg-primary text-on-primary rounded-full hover:bg-primary-dim transition-transform active:scale-90"
          >
            <ArrowRight className="w-4 h-4" />
          </button>
        </div>

        <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
          <button
            type="button"
            className="px-6 py-2 border border-outline-variant/50 rounded-lg text-sm font-semibold hover:bg-surface-container-high transition-colors text-on-surface"
          >
            Export as CSV
          </button>
          <button
            type="button"
            className="px-8 py-2 bg-primary text-on-primary rounded-lg text-sm font-bold hover:bg-primary-dim shadow-sm transition-all hover:shadow-md active:scale-95"
          >
            Save Edits
          </button>
        </div>
      </div>
    </footer>
  );
}
