import { PageThumbnail } from "./PageThumbnail";

type SourcePanelProps = {
  activePage: number;
  onPageSelect: (page: number) => void;
};

const pagePreviews = [
  {
    pageNumber: 1,
    grayscale: false,
    isLoading: false,
    imgSrc:
      "https://lh3.googleusercontent.com/aida-public/AB6AXuBud6m7d3T3KAk4VHYwalvw7G_szVMdSA-g4_sjA5_ZvC6xWDaA5snMe4Nn1rQ1XRyXvfS6JAIQI-NREgkE6WxjcBLZ8Z_wzeARyQwZnlnQn3WNP3QkpP5Rd5x5vaA1JiiM7IWnjTcdg3u593hsba9nX7_QNZIU2LXvNCkUQwYH4tF-N7kR0EwrC6ad0v3ztbhbCx0Jox6Sstkg655dJQw7DtrCaEpRkDG3Y5SYd7c7ZuAzj5kBFvuAjVdvLC8drSfSfLwwT2LNlRQ",
  },
  {
    pageNumber: 2,
    grayscale: true,
    isLoading: false,
    imgSrc:
      "https://lh3.googleusercontent.com/aida-public/AB6AXuB0UD_tGd6oVVozpxVkgZ8BY-mo1-UsnLzeu_XEeZ3QMBjgNrI4ubD8t-5cDEM6xm8cyd620B2abRPbBRSvcIDfq86riW9E19Sc5yb_PluRcdtiV9st_REw6K6deiJ4_Y0THddaID7XOSH5RrLI-p69LmugygTd77TeCsJvet0EIMla015ahoyvT5yzuKTtI-eUDHWQ9_4eU0qKqxrCkoNXOlTRMRXAkwySRSCLhfbsTfP_Qz55VARi9gTl_VzKAFKmgI8fcoCCF_o",
  },
  {
    pageNumber: 3,
    grayscale: false,
    isLoading: true,
    imgSrc: undefined,
  },
] as const;

export function SourcePanel({ activePage, onPageSelect }: SourcePanelProps) {
  return (
    <section className="min-h-0 flex flex-col bg-surface-container-low border-b border-outline-variant/10 lg:border-b-0 lg:border-r">
      <div className="p-4 bg-surface-container border-b border-outline-variant/10 flex items-center justify-between shrink-0">
        <div className="flex min-w-0 flex-col">
          <span className="text-xs font-bold text-on-surface-variant uppercase tracking-wider">
            Source View
          </span>
          <span className="text-sm font-semibold truncate max-w-45">
            report_042.pdf
          </span>
        </div>
        <span className="bg-surface-container-lowest text-on-secondary-container text-[10px] font-bold px-2 py-1 rounded border border-outline-variant/20 uppercase">
          12 pages
        </span>
      </div>

      <div className="grid grid-cols-2 gap-4 overflow-y-auto bg-surface-container-low/50 p-4 md:grid-cols-3 lg:block lg:space-y-6 lg:p-6 custom-scrollbar">
        {pagePreviews.map((page) => (
          <PageThumbnail
            key={page.pageNumber}
            pageNumber={page.pageNumber}
            isActive={activePage === page.pageNumber}
            onClick={() => onPageSelect(page.pageNumber)}
            imgSrc={page.imgSrc}
            isLoading={page.isLoading}
            grayscale={page.grayscale}
          />
        ))}
      </div>
    </section>
  );
}
