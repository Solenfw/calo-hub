"use client"

import { useState } from 'react';
import { Search, Filter, ChevronLeft, ChevronRight, ChevronDown, Plus, X, ZoomIn, Copy, Check } from 'lucide-react';
import { cn } from '@/lib/utils';
import { KLSProduct, AesculapProduct } from '@/types';
import { searchAesculapProduct, searchKLSProduct, getKLSImages } from './productHandlers';

export function Catalog() {
  // Search & Filter States
  const [searchTerm, setSearchTerm] = useState('');
  const [brand, setBrand] = useState('All Brands');
  const [lang, setLang] = useState('EN');
  const [limit, setLimit] = useState(50);
  const [copied, setCopied] = useState(false);
  
  // Data States
  const [catalogItems, setCatalogItems] = useState<KLSProduct[] | AesculapProduct[] | null>(null);
  const [chosenItem, setChosenItem] = useState<KLSProduct | AesculapProduct | null>(null);
  
  // Media / Lightbox States
  const [currentImages, setCurrentImages] = useState<string[]>([]);
  const [currentImageIdx, setCurrentImageIdx] = useState(0);
  const [isImageLoading, setIsImageLoading] = useState(false);
  const [isExpanded, setIsExpanded] = useState(false);

  // Input Handlers
  const handleLimitChange = (e: React.ChangeEvent<HTMLSelectElement>) => setLimit(parseInt(e.target.value));
  const handleLangChange = (e: React.ChangeEvent<HTMLSelectElement>) => setLang(e.target.value);
  const handleBrandChange = (e: React.ChangeEvent<HTMLSelectElement>) => setBrand(e.target.value);
  const handleSearchTermChange = (e: React.ChangeEvent<HTMLInputElement>) => setSearchTerm(e.target.value);

  // Search Handler
  const handleSearch = async () => {
    let products: (KLSProduct | AesculapProduct)[] = [];

    if (brand === 'Martin') {
        products = await searchKLSProduct(searchTerm, limit) ?? [];
    } else if (brand === 'B-Braun') {
        products = await searchAesculapProduct(searchTerm, limit) ?? [];
    } else if (brand === 'All Brands') {
        const [kls, aesculap] = await Promise.all([
            searchKLSProduct(searchTerm, limit),
            searchAesculapProduct(searchTerm, limit),
        ]);
        products = [...(kls ?? []), ...(aesculap ?? [])];
    }

    setCatalogItems(products);
    setChosenItem(null);
    setCurrentImages([]);
    setCurrentImageIdx(0);
  };

  // Row Click Handler (Fetches images dynamically)
  const handleSelectItem = async (item: KLSProduct | AesculapProduct) => {
    setChosenItem(item);
    setCurrentImageIdx(0);
    setCurrentImages([]);
    
    // Handle Martin (up to 3 images fetched from API)
    if (item.brand === 'Martin') {
      setIsImageLoading(true);
      try {
        const klsMedia = await getKLSImages(item.code); 
        if (klsMedia) {
          // Filter out empty URLs to create a clean array of available images
          const validImages = [klsMedia.img1_url, klsMedia.img2_url, klsMedia.img3_url].filter(Boolean) as string[];
          setCurrentImages(validImages.length > 0 ? validImages : ['https://placehold.co/600x400?text=No+Image']);
        } else {
          setCurrentImages(['https://placehold.co/600x400?text=No+Image']);
        }
      } catch (error) {
        setCurrentImages(['https://placehold.co/600x400?text=No+Image']);
      } finally {
        setIsImageLoading(false);
      }
    } 
    // Handle Aesculap (1 image attached directly to product)
    else {
      const aesculapImg = (item as AesculapProduct).image;
      setCurrentImages(aesculapImg ? [aesculapImg] : ['/placeholder-image.png']);
    }
  };

  // Slider Control Helpers
  const nextImage = (e: React.MouseEvent) => {
    e.stopPropagation();
    setCurrentImageIdx(prev => prev === currentImages.length - 1 ? 0 : prev + 1);
  };
  
  const prevImage = (e: React.MouseEvent) => {
    e.stopPropagation();
    setCurrentImageIdx(prev => prev === 0 ? currentImages.length - 1 : prev - 1);
  };

  return (
    <div className="max-w-6xl mx-auto relative">
      
      {/* Page Header */}
      <div className="mb-12">
        <h1 className="font-manrope text-4xl font-extrabold tracking-tight text-on-surface mb-2">Online Catalog</h1>
        <p className="text-on-surface-variant font-body">Browse and search through the comprehensive clinical production database.</p>
      </div>

      {/* Filter & Search Toolbar */}
      <div className="bg-surface-container-lowest rounded-xl p-6 ambient-shadow mb-8 flex flex-col md:flex-row gap-4 items-end">
        <div className="flex-1 w-full">
          <label className="block text-[11px] font-bold uppercase tracking-widest text-outline-variant mb-2 ml-1">Search Products</label>
          <div className="relative group">
            <Search className="absolute left-4 top-1/2 -translate-y-1/2 text-outline-variant group-focus-within:text-primary transition-colors w-4 h-4" />
            <input 
              value={searchTerm}
              onChange={handleSearchTermChange}
              onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
              type="text" 
              className="w-full bg-surface-container-high border-none rounded-lg pl-12 pr-4 py-3.5 focus:ring-0 focus:border-b-2 focus:border-primary transition-all placeholder:text-outline-variant text-on-surface"
              placeholder="Enter product name or code (e.g. KN-400)"
            />
          </div>
        </div>
        
        <div className="w-full md:w-48">
          <label className="block text-[11px] font-bold uppercase tracking-widest text-outline-variant mb-2 ml-1">Brand Filter</label>
          <div className="relative">
            <select 
              value={brand}
              onChange={handleBrandChange}
              className="w-full appearance-none bg-surface-container-high border-none rounded-lg pl-4 pr-10 py-3.5 focus:ring-0 focus:border-b-2 focus:border-primary transition-all text-on-surface cursor-pointer"
            >
              <option>All Brands</option>
              <option>Martin</option>
              <option>B-Braun</option>
              <option>Stema</option>
            </select>
            <ChevronDown className="absolute right-4 top-1/2 -translate-y-1/2 pointer-events-none text-outline-variant w-4 h-4" />
          </div>
        </div>

        <div className="w-full md:w-32">
          <label className="block text-[11px] font-bold uppercase tracking-widest text-outline-variant mb-2 ml-1">Language</label>
          <div className="relative">
            <select 
              value={lang}
              onChange={handleLangChange}
              className="w-full appearance-none bg-surface-container-high border-none rounded-lg pl-4 pr-10 py-3.5 focus:ring-0 focus:border-b-2 focus:border-primary transition-all text-on-surface cursor-pointer"
            >
              <option>EN</option>
              <option>VN</option>
            </select>
            <ChevronDown className="absolute right-4 top-1/2 -translate-y-1/2 pointer-events-none text-outline-variant w-4 h-4" />
          </div>
        </div>

        <div className="w-full md:w-24">
          <label className="block text-[11px] font-bold uppercase tracking-widest text-outline-variant mb-2 ml-1">Limit</label>
          <div className="relative">
            <select 
              value={limit}
              onChange={handleLimitChange}
              className="w-full appearance-none bg-surface-container-high border-none rounded-lg pl-4 pr-10 py-3.5 focus:ring-0 focus:border-b-2 focus:border-primary transition-all text-on-surface cursor-pointer"
            >
              <option>50</option>
              <option>100</option>
              <option>200</option>
            </select>
            <ChevronDown className="absolute right-4 top-1/2 -translate-y-1/2 pointer-events-none text-outline-variant w-4 h-4" />
          </div>
        </div>

        <button 
          onClick={handleSearch}
          className="cursor-pointer bg-primary text-on-primary px-8 py-3.5 rounded-lg font-bold flex items-center justify-center gap-2 hover:bg-primary-dim transition-all shadow-lg shadow-primary/10 active:scale-95 h-12 w-full md:w-auto"
        >
          <Filter className="w-4 h-4" />
          Search
        </button>
      </div>

      {/* Grid Layout for Catalog Content */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 relative items-start">
        
        {/* Main Catalog Table */}
        <div className="lg:col-span-8 bg-surface-container-lowest rounded-xl overflow-hidden ambient-shadow h-fit">
          <div className="p-8 pb-4">
            <div className="flex justify-between items-center mb-2">
              <h3 className="font-manrope text-xl font-bold text-black">Catalog Entries</h3>
              <span className="text-[12px] font-bold text-primary bg-primary-container px-3 py-1 rounded-full">
                { catalogItems ? `${catalogItems.length} items` : 'No results' }
              </span>
            </div>
          </div>
          
          <div className="overflow-x-auto">
            <table className="w-full text-left border-collapse">
              <thead>
                <tr className="bg-surface-container-low border-y border-outline-variant/20">
                  <th className="px-8 py-4 text-[11px] font-extrabold uppercase tracking-widest text-outline">Code</th>
                  <th className="px-8 py-4 text-[11px] font-extrabold uppercase tracking-widest text-outline">Description</th>
                  {(brand === 'B-Braun' || brand === 'All Brands') && (
                    <th className="px-8 py-4 text-[11px] font-extrabold uppercase tracking-widest text-outline">Alt Code</th>
                  )}
                </tr>
              </thead>
              <tbody className="divide-y divide-outline-variant/10">
                {catalogItems?.map((item) => {
                  const isSelected = chosenItem?.code === item.code;
                  return (
                    <tr 
                      onClick={() => handleSelectItem(item)}
                      key={item.code} 
                      className={cn(
                        "group transition-colors cursor-pointer",
                        isSelected ? "bg-primary-container/30" : "hover:bg-surface-container-low"
                      )}
                    >
                      <td className="px-8 py-5">
                        <span className="font-mono text-sm font-bold text-primary">{item.code}</span>
                      </td>
                      <td className="px-8 py-5">
                        <div className="flex flex-col">
                          <span className="font-medium text-on-surface">
                            {lang === 'EN' ? item.eng_desc : item.viet_desc}
                          </span>
                        </div>
                      </td>
                      {(brand === 'B-Braun' || brand === 'All Brands') && (
                        <td className="px-8 py-5">
                          <span className="font-mono text-sm font-medium text-outline">
                            {(item as AesculapProduct).alternative_code || '-'}
                          </span>
                        </td>
                      )}
                    </tr>
                  )
                })}
              </tbody>
            </table>
            
            {!catalogItems?.length && (
              <div className="p-12 text-center text-outline-variant">
                <Search className="w-12 h-12 mx-auto mb-4 opacity-20" />
                <p>Use the search bar above to find clinical instruments.</p>
              </div>
            )}
          </div>
        </div>

        {/* Detail Side Column */}
        <div className="lg:col-span-4 flex flex-col gap-8 self-start">
          <div className="bg-surface-container-lowest p-8 rounded-xl ambient-shadow sticky top-0">
            
            {/* Image Box */}
            <div 
              onClick={() => { if (chosenItem && !isImageLoading) setIsExpanded(true); }}
              className="relative aspect-square rounded-lg mb-6 overflow-hidden bg-white border border-outline-variant/30 flex items-center justify-center group cursor-zoom-in"
            >
              {chosenItem ? (
                isImageLoading ? (
                  // Loading State
                  <div className="flex flex-col items-center text-outline-variant">
                    <div className="w-8 h-8 border-4 border-outline-variant/30 border-t-primary rounded-full animate-spin mb-4" />
                    <span className="text-[10px] font-bold uppercase tracking-widest">Loading Media...</span>
                  </div>
                ) : (
                  <>
                    {/* Main Image View */}
                    <img 
                      src={currentImages[currentImageIdx]}
                      alt={`${chosenItem.code} Detail`} 
                      className="w-full h-full object-contain p-4 mix-blend-darken transition-transform duration-300 group-hover:scale-105"
                      referrerPolicy="no-referrer"
                    />
                    
                    {/* Expand Hover Overlay */}
                    <div className="absolute inset-0 bg-black/0 group-hover:bg-black/5 transition-colors duration-300 pointer-events-none flex items-center justify-center">
                      <div className="bg-white/80 p-2 rounded-full opacity-0 group-hover:opacity-100 transition-opacity shadow-sm drop-shadow-md">
                        <ZoomIn className="w-5 h-5 text-on-surface" />
                      </div>
                    </div>

                    {/* MINI SLIDER (Only renders if > 1 image exists) */}
                    {currentImages.length > 1 && (
                      <>
                        <button onClick={prevImage} className="absolute left-2 top-1/2 -translate-y-1/2 bg-white/90 shadow-md p-1.5 rounded-full opacity-0 group-hover:opacity-100 transition-opacity hover:bg-surface-container-highest cursor-pointer">
                          <ChevronLeft className="w-4 h-4 text-on-surface" />
                        </button>
                        
                        <button onClick={nextImage} className="absolute right-2 top-1/2 -translate-y-1/2 bg-white/90 shadow-md p-1.5 rounded-full opacity-0 group-hover:opacity-100 transition-opacity hover:bg-surface-container-highest cursor-pointer">
                          <ChevronRight className="w-4 h-4 text-on-surface" />
                        </button>
                        
                        {/* Pagination Dots */}
                        <div className="absolute bottom-3 left-1/2 -translate-x-1/2 flex gap-1.5 z-10">
                          {currentImages.map((_, idx) => (
                            <div 
                              key={idx} 
                              className={cn("h-1.5 rounded-full transition-all", currentImageIdx === idx ? "bg-primary w-4" : "bg-outline-variant/60 w-1.5" )} 
                            />
                          ))}
                        </div>
                      </>
                    )}
                  </>
                )
              ) : (
                // Empty State
                <div className="text-outline-variant flex flex-col items-center gap-2">
                  <div className="w-16 h-16 rounded-full bg-surface-container flex items-center justify-center mb-2">
                    <Search className="w-6 h-6 opacity-50" />
                  </div>
                  <span className="text-sm font-medium">Select a product</span>
                </div>
              )}
            </div>

            {/* Product Details */}
            <div className="min-h-55">
              {chosenItem ? (
                <>
                  <div className="flex items-center gap-2 mb-3">
                    <span className="px-2 py-0.5 rounded bg-error-container/20 text-error font-bold text-[10px] uppercase">
                      {chosenItem.brand}
                    </span>
                  </div>
                  <h4 className="inline-block font-manrope text-3xl mb-2 font-bold text-black">{chosenItem.code}</h4>
                  <button
                    onClick={() => {
                      navigator.clipboard.writeText(chosenItem.code);
                      setCopied(true);
                      setTimeout(() => setCopied(false), 1500);
                    }}
                    className="inline-block p-1.5 rounded hover:bg-surface-variant transition-colors text-outline-variant hover:text-on-surface"
                    title="Copy code"
                  >
                    {copied ? <Check size={16} className="text-green-500" /> : <Copy size={16} />}
                  </button>
                  <p className="text-on-surface-variant text-sm mb-6 leading-relaxed line-clamp-3">
                    {lang === 'EN' ? chosenItem.eng_desc : chosenItem.viet_desc}
                  </p>
                  <div className="space-y-4">
                    <DetailRow label="Material" value="Hardened Steel" />
                    <DetailRow label="Sterility" value="Non-Sterile" />
                  </div>
                </>
              ) : (
                <div className="h-full flex flex-col items-center justify-center text-center text-outline-variant mt-10">
                  <p className="text-sm">Product details and specifications <br/> will appear here.</p>
                </div>
              )}
            </div>

            <button 
              onClick={(e) => {
                if (!chosenItem) {
                  e.preventDefault();
                  return; 
                }
              }}
              aria-disabled={!chosenItem}
              className="cursor-pointer disabled:cursor-not-allowed disabled:opacity-50 w-full mt-8 bg-secondary-container text-on-secondary-container py-4 rounded-lg font-bold hover:bg-secondary-fixed-dim transition-all active:scale-95"
            >
              Request Technical Sheet
            </button>
          </div>

          {/* Integration Ad */}
          <div className="bg-primary-container p-8 rounded-xl relative overflow-hidden group">
            <div className="relative z-10">
              <h5 className="text-primary font-manrope font-bold text-lg mb-2">Direct Integration</h5>
              <p className="text-primary/70 text-sm mb-6">Seamlessly export these codes directly to your medical inventory management system.</p>
              <button className="cursor-pointer inline-flex items-center gap-2 text-primary font-bold text-sm hover:underline">
                Learn More
                <ChevronRight className="w-4 h-4" />
              </button>
            </div>
          </div>
        </div>

      </div>

      {/* Floating Action Button */}
      <div className="fixed bottom-8 right-8 z-40">
        <button className="cursor-pointer w-14 h-14 bg-primary text-on-primary rounded-full shadow-2xl flex items-center justify-center hover:scale-105 active:scale-95 transition-transform">
          <Plus className="w-6 h-6" />
        </button>
      </div>

      {/* FULLSCREEN LIGHTBOX MODAL */}
      {isExpanded && chosenItem && (
        <div className="fixed inset-0 z-100 bg-black/95 backdrop-blur-sm flex items-center justify-center">
          
          {/* Top-Left Close Button */}
          <button 
            onClick={() => setIsExpanded(false)}
            className="absolute top-6 left-6 text-white/70 hover:text-white bg-white/10 hover:bg-white/20 p-3 rounded-full transition-colors cursor-pointer z-50 flex items-center justify-center"
          >
            <X className="w-6 h-6" />
          </button>

          {/* Expanded Image Container */}
          <div className="relative w-full max-w-6xl h-[85vh] flex flex-col items-center justify-center p-4 md:p-12">
            <img 
              src={currentImages[currentImageIdx]} 
              alt="Expanded View" 
              className="w-full h-full object-contain drop-shadow-2xl"
              referrerPolicy="no-referrer"
            />
            
            {/* FULLSCREEN SLIDER (Only shows if > 1 image) */}
            {currentImages.length > 1 && (
              <>
                {/* Left Arrow */}
                <button 
                  onClick={prevImage} 
                  className="absolute left-4 md:left-12 top-1/2 -translate-y-1/2 bg-white/10 hover:bg-white/20 text-white p-4 rounded-full transition-colors cursor-pointer backdrop-blur-md"
                >
                  <ChevronLeft className="w-8 h-8" />
                </button>
                
                {/* Right Arrow */}
                <button 
                  onClick={nextImage} 
                  className="absolute right-4 md:right-12 top-1/2 -translate-y-1/2 bg-white/10 hover:bg-white/20 text-white p-4 rounded-full transition-colors cursor-pointer backdrop-blur-md"
                >
                  <ChevronRight className="w-8 h-8" />
                </button>

                {/* Fullscreen Pagination Dots */}
                <div className="absolute bottom-0 left-1/2 -translate-x-1/2 flex gap-3">
                  {currentImages.map((_, idx) => (
                    <button 
                      key={idx}
                      onClick={(e) => { e.stopPropagation(); setCurrentImageIdx(idx); }}
                      className={cn(
                        "h-2 rounded-full transition-all cursor-pointer", 
                        currentImageIdx === idx ? "bg-white w-8" : "bg-white/40 hover:bg-white/70 w-2"
                      )}
                    />
                  ))}
                </div>
              </>
            )}
          </div>
          
          {/* Modal Bottom-Left Product Info */}
          <div className="absolute bottom-8 left-8 md:left-12 text-white pointer-events-none">
            <span className="bg-primary px-3 py-1 text-xs font-bold uppercase rounded-md mb-2 inline-block">
              {chosenItem.brand}
            </span>
            <h2 className="text-3xl md:text-5xl font-manrope font-bold">{chosenItem.code}</h2>
          </div>
          
        </div>
      )}

    </div>
  );
}

function DetailRow({ label, value }: { label: string, value: string }) {
  return (
    <div className="flex justify-between items-center text-sm py-3 border-b border-outline-variant/10 last:border-0">
      <span className="text-outline">{label}</span>
      <span className="font-bold text-on-surface">{value}</span>
    </div>
  );
}