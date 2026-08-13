import { CheckCircle2, Image as ImageIcon } from "lucide-react";
import { cn } from "@/lib/utils";

type PageThumbnailProps = {
  pageNumber: number;
  isActive: boolean;
  onClick: () => void;
  imgSrc?: string;
  isLoading?: boolean;
  grayscale?: boolean;
};

export function PageThumbnail({
  pageNumber,
  isActive,
  onClick,
  imgSrc,
  isLoading = false,
  grayscale = false,
}: PageThumbnailProps) {
  return (
    <button
      type="button"
      className="group w-full cursor-pointer space-y-2 text-left"
      onClick={onClick}
    >
      <div
        className={cn(
          "bg-surface-container-lowest rounded-lg p-1.5 shadow-sm transition-all duration-300 hover:shadow-md relative overflow-hidden aspect-3/4",
          isActive ? "border-2 border-primary" : "border-2 border-transparent"
        )}
      >
        {isLoading ? (
          <div className="w-full h-full bg-surface-container-high flex items-center justify-center animate-pulse">
            <ImageIcon className="w-8 h-8 text-outline-variant opacity-50" />
          </div>
        ) : (
          <img
            className={cn(
              "w-full h-full object-cover rounded shadow-inner",
              grayscale ? "grayscale-[0.2]" : "opacity-90"
            )}
            alt={`Document page ${pageNumber}`}
            src={imgSrc}
          />
        )}
        {isActive && (
          <div className="absolute inset-0 bg-primary/5 pointer-events-none" />
        )}
      </div>
      <div className="flex items-center justify-between px-1">
        <span
          className={cn(
            "text-xs font-bold",
            isActive ? "text-primary" : "text-on-surface-variant"
          )}
        >
          Page {pageNumber}
        </span>
        {isActive && <CheckCircle2 className="w-4 h-4 text-primary" />}
      </div>
    </button>
  );
}
