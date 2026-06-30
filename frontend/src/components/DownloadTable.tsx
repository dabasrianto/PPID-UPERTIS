import { Search, FileDown, Download, ArrowRight } from 'lucide-react';

interface DownloadItemData {
  id?: string;
  label?: string;
  title?: string;
  href?: string;
  file_url?: string;
  downloadsCount?: number;
  downloads_count?: number;
  description?: string;
  category?: string;
  type?: string;
}

interface DownloadTableProps {
  items: DownloadItemData[];
  searchTerm: string;
  setSearchTerm: (term: string) => void;
  incrementDownloadCount: (id: string, fileUrl: string) => void;
  handleNavigation?: (href: string) => void;
  activeTab?: string;
  setActiveTab?: (tab: string) => void;
  categories?: Array<{ id: string; title: string; desc?: string }>;
  title?: string;
  description?: string;
  gridColsClassName?: string;
}

export default function DownloadTable({
  items,
  searchTerm,
  setSearchTerm,
  incrementDownloadCount,
  handleNavigation,
  activeTab,
  setActiveTab,
  categories,
  title,
  description,
  gridColsClassName
}: DownloadTableProps) {
  return (
    <div className="space-y-6">
      {/* Optional Category Tabs */}
      {categories && setActiveTab && activeTab && (
        <div className="flex flex-wrap justify-center gap-2 border-b border-slate-200 pb-4">
          {categories.map((category) => (
            <button
              key={category.id}
              onClick={() => setActiveTab(category.id)}
              className={`px-5 py-2.5 rounded-full text-xs font-bold transition-all cursor-pointer ${
                activeTab === category.id
                  ? 'bg-[#002147] text-white shadow-md'
                  : 'bg-white border border-slate-200 text-slate-600 hover:text-[#002147]'
              }`}
            >
              {category.title}
            </button>
          ))}
        </div>
      )}

      {/* Search Input */}
      <div className="max-w-md relative">
        <input
          type="text"
          placeholder="Cari berkas dokumen..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          className="w-full pl-10 pr-4 py-2.5 bg-white border border-slate-200 rounded-xl text-xs focus:outline-none focus:ring-2 focus:ring-[#002147]/20 focus:border-[#002147] transition-all font-medium text-slate-800"
        />
        <Search className="absolute left-3.5 top-3 h-4 w-4 text-slate-400" />
      </div>

      {/* Main Grid Container */}
      <div className="bg-white border border-slate-200 rounded-3xl p-6 lg:p-8 shadow-sm">
        {title && (
          <h4 className="text-xs font-extrabold text-[#002147] mb-2 uppercase tracking-wider text-left">
            {title}
          </h4>
        )}
        {description && (
          <p className="text-[11px] text-slate-400 mb-6 font-medium leading-relaxed text-left">
            {description}
          </p>
        )}

        {items.length > 0 ? (
          <div className={gridColsClassName || "grid md:grid-cols-2 gap-4"}>
            {items.map((item, idx) => {
              const itemTitle = item.label || item.title || '';
              const rawFileUrl = item.href || item.file_url || '';
              
              // Parse files
              let parsedFilesList: Array<{ name: string; url: string }> = [];
              if (rawFileUrl.trim().startsWith('[')) {
                try {
                  const parsed = JSON.parse(rawFileUrl);
                  if (Array.isArray(parsed)) {
                    parsedFilesList = parsed;
                  }
                } catch (e) {
                  console.error('Failed to parse file_url JSON in DownloadTable:', e);
                }
              }

              if (parsedFilesList.length === 0 && rawFileUrl) {
                // Single file / legacy format
                parsedFilesList = [{ name: 'Unduh Berkas', url: rawFileUrl }];
              }

              const hasMultipleFiles = parsedFilesList.length > 1;
              const isDownloadItem = item.type === 'DownloadItem' || item.file_url;
              const count = item.downloadsCount !== undefined ? item.downloadsCount : item.downloads_count;

              return (
                <div
                  key={item.id || idx}
                  className="flex items-center justify-between p-4 border border-slate-100 hover:border-amber-250 hover:bg-amber-50/5 rounded-2xl transition-all"
                >
                  <div className="flex items-center gap-3.5 min-w-0">
                    <div className="h-9 w-9 rounded-xl bg-slate-50 text-[#002147] flex items-center justify-center shrink-0 border border-slate-200">
                      <FileDown className="h-4.5 w-4.5" />
                    </div>
                    <div className="min-w-0 text-left">
                      <span className="text-xs font-bold text-slate-800 block truncate" title={itemTitle}>
                        {itemTitle}
                      </span>
                      {item.description && (
                        <span className="text-[10px] text-slate-400 block truncate max-w-[250px] mb-0.5">
                          {item.description}
                        </span>
                      )}
                      <span className="text-[10px] text-slate-400 font-medium block mt-0.5 text-left">
                        {item.category && (
                          <>
                            Kategori: <span className="text-slate-650 capitalize font-semibold">{item.category.replace('ppid-', '')}</span>
                            {count !== undefined && <span className="mx-1.5">•</span>}
                          </>
                        )}
                        {count !== undefined ? (
                          <>
                            Unduhan: <strong className="text-slate-600">{count}</strong>
                          </>
                        ) : (
                          !item.category && (
                            <span className="text-[10px] text-amber-600 font-bold uppercase tracking-wide">
                              Halaman Portal
                            </span>
                          )
                        )}
                      </span>
                      {/* Render multiple download buttons inline if document has multiple files */}
                      {hasMultipleFiles && (
                        <div className="flex flex-wrap gap-2 mt-2">
                          {parsedFilesList.map((file, fIdx) => (
                            <a
                              key={fIdx}
                              href={file.url}
                              target="_blank"
                              rel="noopener noreferrer"
                              onClick={() => incrementDownloadCount(item.id || '', file.url)}
                              className="inline-flex items-center gap-1.5 px-3 py-1.5 bg-slate-50 hover:bg-amber-400 hover:text-[#002147] border border-slate-200 text-[#002147] hover:border-amber-300 rounded-xl text-[10px] font-bold uppercase transition-all cursor-pointer"
                            >
                              <Download className="h-3 w-3" /> {file.name || `Berkas ${fIdx + 1}`}
                            </a>
                          ))}
                        </div>
                      )}
                    </div>
                  </div>

                  {isDownloadItem ? (
                    !hasMultipleFiles && (
                      <button
                        onClick={() => incrementDownloadCount(item.id || '', parsedFilesList[0]?.url || '#')}
                        className="p-2 bg-[#002147] hover:bg-amber-500 text-white hover:text-[#002147] rounded-lg transition-all inline-flex items-center justify-center cursor-pointer shrink-0 border border-[#002147]/55"
                        title="Unduh Berkas"
                      >
                        <Download className="h-4 w-4" />
                      </button>
                    )
                  ) : (
                    handleNavigation && (
                      <button
                        onClick={() => handleNavigation(rawFileUrl)}
                        className="px-3.5 py-1.5 bg-slate-50 hover:bg-slate-100 text-slate-600 hover:text-[#002147] rounded-lg text-[10px] font-bold uppercase transition-all inline-flex items-center gap-1 cursor-pointer shrink-0 border border-slate-200/50"
                      >
                        <ArrowRight className="h-3.5 w-3.5" /> Buka
                      </button>
                    )
                  )}
                </div>
              );
            })}
          </div>
        ) : (
          <div className="py-12 text-center text-slate-400 text-xs font-medium">
            Tidak ada dokumen ditemukan untuk pencarian ini.
          </div>
        )}
      </div>
    </div>
  );
}
