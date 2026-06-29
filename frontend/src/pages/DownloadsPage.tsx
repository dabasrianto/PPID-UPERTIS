import { Download } from 'lucide-react';
import DownloadTable from '../components/DownloadTable';
import type { DownloadItem } from '../types';

interface DownloadsPageProps {
  dbDownloads: DownloadItem[];
  searchTerm: string;
  setSearchTerm: (term: string) => void;
  incrementDownloadCount: (id: string, fileUrl: string) => void;
  activeTab: string;
  setActiveTab: (tab: string) => void;
}

export default function DownloadsPage({
  dbDownloads,
  searchTerm,
  setSearchTerm,
  incrementDownloadCount,
  activeTab,
  setActiveTab
}: DownloadsPageProps) {
  return (
    <div className="space-y-8 animate-in fade-in duration-200 text-left w-full py-6">
      {/* Header Banner */}
      <div className="bg-[#002147] text-white rounded-3xl p-8 lg:p-12 relative overflow-hidden shadow-lg">
        <div className="absolute right-0 bottom-0 opacity-10 translate-x-12 translate-y-12">
          <Download className="h-64 w-64" />
        </div>
        <div className="relative z-10 space-y-3">
          <span className="bg-amber-400 text-[#002147] text-[10px] font-extrabold px-3 py-1 rounded-full uppercase tracking-wider">
            {activeTab === 'dikecualikan' ? 'Informasi Publik' : 'Pusat Unduhan'}
          </span>
          <h1 className="text-3xl lg:text-4xl font-extrabold leading-tight">
            {activeTab === 'dikecualikan' ? 'Informasi Dikecualikan' : 'Unduhan Dokumen Resmi'}
          </h1>
          <p className="text-xs lg:text-sm text-slate-200 leading-relaxed font-medium">
            {activeTab === 'dikecualikan' 
              ? 'Daftar dokumen keterbukaan informasi publik di lingkungan UPERTIS yang dikecualikan berdasarkan peraturan perundang-undangan dan pengujian konsekuensi.'
              : 'Unduh berkas Surat Keputusan (SK), regulasi, standar operasional prosedur (SOP), dan dokumen penting keterbukaan informasi publik lainnya.'}
          </p>
        </div>
      </div>

      <DownloadTable
        items={dbDownloads.filter(item => {
          const matchesCategory = activeTab === 'all' ||
            item.category === activeTab ||
            item.category === `ppid-${activeTab}` ||
            (activeTab === 'umum' && !['berkala', 'ppid-berkala', 'setiap-saat', 'ppid-setiap-saat', 'serta-merta', 'ppid-serta-merta', 'dikecualikan', 'ppid-dikecualikan'].includes(item.category));

          const matchesSearch = item.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
            (item.description && item.description.toLowerCase().includes(searchTerm.toLowerCase()));
          return matchesCategory && matchesSearch;
        })}
        searchTerm={searchTerm}
        setSearchTerm={setSearchTerm}
        incrementDownloadCount={incrementDownloadCount}
        activeTab={activeTab}
        setActiveTab={setActiveTab}
        categories={[
          { id: 'all', title: 'Semua Kategori' },
          { id: 'berkala', title: 'Informasi Berkala' },
          { id: 'setiap-saat', title: 'Tersedia Setiap Saat' },
          { id: 'serta-merta', title: 'Informasi Serta Merta' },
          { id: 'dikecualikan', title: 'Informasi Dikecualikan' },
          { id: 'umum', title: 'Umum' }
        ]}
      />
    </div>
  );
}
