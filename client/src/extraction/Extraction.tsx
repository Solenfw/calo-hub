"use client";

import { useState } from "react";
import { FileText } from "lucide-react";
import { ExtractionBottomBar } from "./components/ExtractionBottomBar";
import { OutputPanel } from "./components/OutputPanel";
import { SourcePanel } from "./components/SourcePanel";

export function Extraction() {
  const [activePage, setActivePage] = useState(1);

  return (
    <div className="mx-auto flex h-[calc(100vh-8rem)] min-h-176 max-w-7xl flex-col overflow-hidden rounded-xl bg-surface-container-lowest ambient-shadow">
      <div className="px-4 py-4 bg-surface flex flex-col gap-3 shrink-0 border-b border-outline-variant/10 md:flex-row md:items-center md:justify-between lg:px-6">
        <div className="flex min-w-0 items-center gap-3">
          <FileText className="w-5 h-5 shrink-0 text-primary" />
          <h1 className="text-lg font-bold text-on-surface">
            Extraction Workbench
          </h1>
          <span className="hidden text-outline-variant sm:inline">/</span>
          <span className="truncate text-on-surface-variant font-medium">
            report_042.pdf
          </span>
        </div>
        <span className="w-fit bg-primary-container/40 text-on-primary-container text-xs px-3 py-1 rounded-full font-semibold border border-primary/10">
          Active Session
        </span>
      </div>

      <div className="grid min-h-0 flex-1 grid-cols-1 overflow-hidden lg:grid-cols-[minmax(18rem,0.4fr)_minmax(0,0.6fr)]">
        <SourcePanel activePage={activePage} onPageSelect={setActivePage} />
        <OutputPanel activePage={activePage} />
      </div>

      <ExtractionBottomBar />
    </div>
  );
}
