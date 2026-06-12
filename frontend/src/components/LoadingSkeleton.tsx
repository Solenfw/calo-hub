export function LoadingSkeleton() {
  return (
    <div className="min-h-[60vh] w-full flex items-start justify-center pt-16">
      <div className="w-full max-w-3xl px-6">
        <div className="animate-pulse space-y-6">
          <div className="h-8 w-56 rounded-md bg-slate-200" />
          <div className="space-y-3">
            <div className="h-4 w-full rounded bg-slate-200" />
            <div className="h-4 w-11/12 rounded bg-slate-200" />
            <div className="h-4 w-9/12 rounded bg-slate-200" />
          </div>
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <div className="h-28 rounded-lg bg-slate-200" />
            <div className="h-28 rounded-lg bg-slate-200" />
          </div>
        </div>
      </div>
    </div>
  );
}

