"use client";

import { useEffect, useRef, useState, type ChangeEvent } from "react";
import {
  DeleteIcon,
  Download,
  FileSpreadsheet,
  Plus,
  Save,
  Trash2,
  Upload,
} from "lucide-react";
import * as XLSX from "xlsx";
import {
  createMartinReport,
  deleteMartinReport,
  downloadImportTemplate,
  exportMartinReportPdf,
  getMartinReportProducts,
  listMartinReports,
  updateReportProducts,
} from "@/report/reportHandlers";
import type { MartinReportListResponse, MartinReportProductResponse } from "@/types";

type InstrumentRow = MartinReportProductResponse;

const createInstrumentRow = (rowNo: number): InstrumentRow => ({
  row_no: rowNo,
  code: "",
  description: "",
  image: null,
  quantity: 1,
});

export function Report() {
  const [reports, setReports] = useState<MartinReportListResponse[]>([]);
  const [selectedReportId, setSelectedReportId] = useState<number | null>(null);
  const [selectedReportName, setSelectedReportName] = useState("");
  const [instruments, setInstruments] = useState<InstrumentRow[]>([]);
  const [newReportName, setNewReportName] = useState("");
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const fileInputRef = useRef<HTMLInputElement | null>(null);

  useEffect(() => {
    void loadReports();
  }, []);

  const loadReports = async () => {
    const nextReports = await listMartinReports();
    setReports(nextReports);

    if (nextReports.length > 0) {
      const first = nextReports[0];
      setSelectedReportId(first.id);
      setSelectedReportName(first.name);
      const products = await getMartinReportProducts(first.id);
      setInstruments(products.map((product, index) => ({ ...product, row_no: product.row_no || index + 1 })));
    } else {
      setSelectedReportId(null);
      setSelectedReportName("");
      setInstruments([]);
    }

    setLoading(false);
  };

  const selectReport = async (report: MartinReportListResponse) => {
    setSelectedReportId(report.id);
    setSelectedReportName(report.name);
    const products = await getMartinReportProducts(report.id);
    setInstruments(products.map((product, index) => ({ ...product, row_no: product.row_no || index + 1 })));
  };

  const handleCreateReport = async () => {
    const name = newReportName.trim();
    if (!name) {
      setError("Please enter a report name before creating it.");
      return;
    }

    const created = await createMartinReport(name);
    if (!created) {
      setError("The report could not be created.");
      return;
    }

    setNewReportName("");
    setError(null);
    const nextReports = [...reports, created];
    setReports(nextReports);
    await selectReport(created);
  };

  const handleDeleteReport = async (reportId: number) => {
    const report = reports.find((entry) => entry.id === reportId);
    if (!report) {
      return;
    }

    const shouldDelete = window.confirm(`Delete "${report.name}"?`);
    if (!shouldDelete) {
      return;
    }

    const ok = await deleteMartinReport(reportId);
    if (!ok) {
      setError("The report could not be deleted.");
      return;
    }

    const nextReports = reports.filter((entry) => entry.id !== reportId);
    setReports(nextReports);
    if (selectedReportId === reportId) {
      setSelectedReportId(nextReports[0]?.id ?? null);
      setSelectedReportName(nextReports[0]?.name ?? "");
      setInstruments(nextReports[0] ? await getMartinReportProducts(nextReports[0].id) : []);
    }
    setError(null);
  };

  const normalizeInstruments = (rows: InstrumentRow[]) =>
    rows.map((row, index) => ({
      row_no: index + 1,
      code: row.code.trim(),
      description: row.description.trim(),
      image : row.image?.trim() || "",
      quantity: Number(row.quantity) > 0 ? Number(row.quantity) : 1,
    }));

  const handleSave = async () => {
    if (selectedReportId == null) {
      return;
    }

    const normalized = normalizeInstruments(instruments);
    setSaving(true);
    const ok = await updateReportProducts(selectedReportId, normalized);
    setSaving(false);

    if (!ok) {
      setError("The report changes could not be saved.");
      return;
    }

    setInstruments(normalized);
    setError(null);
  };

  const handleImport = async (event: ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file || selectedReportId == null) {
      return;
    }

    try {
      const workbook = XLSX.read(await file.arrayBuffer(), { type: "array" });
      const sheet = workbook.Sheets[workbook.SheetNames[0]];
      const rows = XLSX.utils.sheet_to_json<Record<string, unknown>>(sheet, { defval: "" });

      const imported: MartinReportProductResponse[] = rows.flatMap((row, index) => {
        const code = String(row.code ?? row.Code ?? row["Code"] ?? "").trim();
        const description = String(
          row.description ?? row.Description ?? row["Description"] ?? "",
        ).trim();
        const quantity = Number(row.quantity ?? row.Quantity ?? row.Qty ?? row.qty ?? row["qty"] ?? 1) || 1;

        if (!code || !description) {
          return [];
        }

        return [{
          row_no: index + 1,
          code,
          description,
          image: null,
          quantity,
        }];
      });

      if (imported.length === 0) {
        setError("No valid instrument rows were found in the uploaded file.");
        return;
      }

      const merged = normalizeInstruments([...instruments, ...imported]);
      setInstruments(merged);
      const ok = await updateReportProducts(selectedReportId, merged);
      if (!ok) {
        setError("The imported instruments could not be saved to the report.");
      } else {
        setError(null);
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : "The spreadsheet could not be parsed.");
    } finally {
      event.target.value = "";
    }
  };

  const handleDownloadTemplate = async () => {
    try {
      const blob = await downloadImportTemplate();
      const url = URL.createObjectURL(blob);
      const link = document.createElement("a");
      link.href = url;
      link.download = "martin-import-template.xlsx";
      link.click();
      URL.revokeObjectURL(url);
    } catch (err) {
      setError(err instanceof Error ? err.message : "The import template could not be downloaded.");
    }
  };

  const handleExport = async (report: MartinReportListResponse) => {
    try {
      const blob = await exportMartinReportPdf(report.name);
      const url = URL.createObjectURL(blob);
      const link = document.createElement("a");
      link.href = url;
      link.download = `${report.name.replace(/\s+/g, "_") || "report"}.pdf`;
      link.click();
      URL.revokeObjectURL(url);
    } catch (err) {
      setError(err instanceof Error ? err.message : "The report PDF could not be generated.");
    }
  };

  const handleInstrumentChange = (
    index: number,
    field: "code" | "description" | "image" | "quantity",
    value: string | number | null,
  ) => {
    setInstruments((prev) =>
      prev.map((row, rowIndex) => {
        if (rowIndex !== index) {
          return row;
        }

        return { ...row, [field]: value } as InstrumentRow;
      }),
    );
  };

  if (loading) {
    return <div className="max-w-7xl mx-auto text-on-surface">Loading report list...</div>;
  }

  return (
    <div className="max-w-7xl mx-auto space-y-6">
      <section className="bg-surface-container-lowest p-6 rounded-xl ambient-shadow">
        <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
          <div>
            <p className="text-xs font-bold tracking-[0.16em] uppercase text-primary">Martin reports</p>
            <h1 className="mt-2 text-3xl font-extrabold text-on-surface">Reports</h1>
          </div>

          <div className="flex w-full max-w-xl gap-3">
            <input
              value={newReportName}
              onChange={(event) => setNewReportName(event.target.value)}
              onKeyDown={(event) => {
                if (event.key === "Enter") {
                  void handleCreateReport();
                }
              }}
              placeholder="New report name"
              className="flex-1 rounded-lg border border-outline-variant bg-surface-container-high px-4 py-3 text-on-surface outline-none ring-0"
            />
            <button
              type="button"
              onClick={() => void handleCreateReport()}
              className="flex items-center gap-2 rounded-lg bg-primary px-4 py-3 font-bold text-on-primary"
            >
              <Plus className="h-4 w-4" />
              Create
            </button>
          </div>
        </div>

        {error ? <p className="mt-4 text-sm text-error">{error}</p> : null}

        <div className="mt-6 overflow-hidden rounded-xl border border-outline-variant/70">
          <table className="w-full border-collapse text-left text-sm text-on-surface">
            <thead className="bg-surface-container-high text-xs uppercase tracking-[0.12em] text-on-surface-variant">
              <tr>
                <th className="px-4 py-3">Name</th>
                <th className="px-4 py-3">Instruments</th>
                <th className="px-4 py-3 text-right">Actions</th>
              </tr>
            </thead>
            <tbody>
              {reports.length === 0 ? (
                <tr>
                  <td colSpan={3} className="px-4 py-6 text-center text-on-surface-variant">
                    No reports created yet.
                  </td>
                </tr>
              ) : (
                reports.map((report) => {
                  const reportInstrumentCount = report.id === selectedReportId ? instruments.length : 0;
                  const isSelected = report.id === selectedReportId;

                  return (
                    <tr
                      key={report.id}
                      className={isSelected ? "bg-primary-container/20" : "bg-surface-container-lowest hover:bg-surface-container-low"}
                    >
                      <td
                        className="cursor-pointer px-4 py-3 font-semibold"
                        onClick={() => void selectReport(report)}
                      >
                        {report.name}
                      </td>
                      <td className="px-4 py-3 text-on-surface-variant">
                        {reportInstrumentCount} rows
                      </td>
                      <td className="px-4 py-3">
                        <div className="flex items-center justify-end gap-2">
                          <button
                            type="button"
                            onClick={() => void handleDownloadTemplate()}
                            className="flex items-center gap-2 rounded-lg border border-outline-variant px-3 py-2 text-xs font-bold text-on-surface"
                          >
                            <Download className="h-3.5 w-3.5" />
                            Template
                          </button>
                          <button
                            type="button"
                            onClick={() => {
                              setSelectedReportId(report.id);
                              setSelectedReportName(report.name);
                              void handleExport(report);
                            }}
                            className="flex items-center gap-2 rounded-lg border border-outline-variant px-3 py-2 text-xs font-bold text-on-surface"
                          >
                            <Download className="h-3.5 w-3.5" />
                            Export
                          </button>
                          <button
                            type="button"
                            onClick={() => {
                              setSelectedReportId(report.id);
                              setSelectedReportName(report.name);
                              if (fileInputRef.current) {
                                fileInputRef.current.click();
                              }
                            }}
                            className="flex items-center gap-2 rounded-lg border border-outline-variant px-3 py-2 text-xs font-bold text-on-surface"
                          >
                            <Upload className="h-3.5 w-3.5" />
                            Import
                          </button>
                          <button
                            type="button"
                            onClick={() => void handleDeleteReport(report.id)}
                            className="rounded-lg p-2 text-on-surface-variant transition hover:bg-surface-container-high"
                            aria-label={`Delete ${report.name}`}
                          >
                            <DeleteIcon className="h-4 w-4" />
                          </button>
                        </div>
                      </td>
                    </tr>
                  );
                })
              )}
            </tbody>
          </table>
        </div>
      </section>

      <section className="bg-surface-container-lowest p-6 rounded-xl ambient-shadow">
        {selectedReportId == null ? (
          <div className="text-on-surface-variant">Select a report to review its instruments.</div>
        ) : (
          <>
            <div className="mb-4 flex items-center justify-between gap-4">
              <div>
                <p className="text-xs font-bold uppercase tracking-[0.12em] text-primary">Detail view</p>
                <h2 className="mt-1 text-2xl font-bold text-on-surface">{selectedReportName}</h2>
              </div>

              <button
                type="button"
                onClick={() => void handleSave()}
                disabled={saving}
                className="flex items-center gap-2 rounded-lg bg-primary px-4 py-3 font-bold text-on-primary disabled:opacity-60"
              >
                <Save className="h-4 w-4" />
                {saving ? "Saving..." : "Save"}
              </button>
            </div>

            <div className="mb-4 flex items-center gap-2 text-xs uppercase tracking-[0.12em] text-on-surface-variant">
              <FileSpreadsheet className="h-4 w-4" />
              Instrument list
            </div>

            <div className="overflow-hidden rounded-xl border border-outline-variant/70">
              <table className="w-full border-collapse text-sm text-on-surface">
                <thead className="bg-surface-container-high text-xs uppercase tracking-[0.12em] text-on-surface-variant">
                  <tr>
                    <th className="px-3 py-2">Code</th>
                    <th className="px-3 py-2">Description</th>
                    <th className="px-3 py-2">Qty</th>
                    <th className="px-3 py-2 text-right">Delete</th>
                  </tr>
                </thead>
                <tbody>
                  {instruments.length === 0 ? (
                    <tr>
                      <td colSpan={5} className="px-3 py-6 text-center text-on-surface-variant">
                        No instruments yet.
                      </td>
                    </tr>
                  ) : (
                    instruments.map((instrument, index) => (
                      <tr key={`${instrument.code}-${index}`} className="border-t border-outline-variant/50">
                        <td className="px-3 py-2">
                          <input
                            value={instrument.code ?? ""}
                            onChange={(event) => handleInstrumentChange(index, "code", event.target.value)}
                            className="w-full rounded-md border border-outline-variant bg-surface-container-high px-2 py-2 outline-none"
                          />
                        </td>
                        <td className="px-3 py-2">
                          <input
                            value={instrument.description ?? ""}
                            onChange={(event) => handleInstrumentChange(index, "description", event.target.value)}
                            className="w-full rounded-md border border-outline-variant bg-surface-container-high px-2 py-2 outline-none"
                          />
                        </td>
                        <td className="px-3 py-2">
                          <input
                            type="number"
                            min={1}
                            value={instrument.quantity ?? 1}
                            onChange={(event) => handleInstrumentChange(index, "quantity", Number(event.target.value) || 1)}
                            className="w-20 rounded-md border border-outline-variant bg-surface-container-high px-2 py-2 outline-none"
                          />
                        </td>
                        <td className="px-3 py-2 text-right">
                          <button
                            type="button"
                            onClick={() => setInstruments((prev) => prev.filter((_, rowIndex) => rowIndex !== index))}
                            className="rounded-lg p-2 text-error hover:bg-error/10"
                            aria-label="Remove instrument"
                          >
                            <Trash2 className="h-4 w-4" />
                          </button>
                        </td>
                      </tr>
                    ))
                  )}
                </tbody>
              </table>
            </div>

            <div className="mt-4 flex justify-end">
              <button
                type="button"
                onClick={() => setInstruments((prev) => [...prev, createInstrumentRow(prev.length + 1)])}
                className="flex items-center gap-2 rounded-lg border border-outline-variant px-4 py-2 font-bold text-on-surface"
              >
                <Plus className="h-4 w-4" />
                Add instrument
              </button>
            </div>
          </>
        )}
      </section>

      <input
        ref={fileInputRef}
        type="file"
        accept=".xlsx,.xls,.csv"
        className="hidden"
        onChange={(event) => {
          void handleImport(event);
        }}
      />
    </div>
  );
}
