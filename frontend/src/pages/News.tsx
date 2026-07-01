import { useState } from 'react';
import { Newspaper, LayoutGrid, List } from 'lucide-react';
import type { Post } from '../types';
import NewsCard from '../components/NewsCard';


interface NewsProps {
  posts: Post[];
  isPostsLoading: boolean;
  newsFilterCategory: string;
  newsFilterSearch: string;
  setNewsFilterCategory: (cat: string) => void;
  setNewsFilterSearch: (search: string) => void;
  navigateToHome: () => void;
  navigateToNewsDetail: (slug: string) => void;
}

export default function News({
  posts,
  isPostsLoading,
  newsFilterCategory,
  newsFilterSearch,
  setNewsFilterCategory,
  setNewsFilterSearch,
  navigateToHome,
  navigateToNewsDetail
}: NewsProps) {
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid');

  return (
    <div className="container mx-auto px-4 py-12 max-w-6xl text-left flex-1 space-y-6">

      <div className="space-y-2">
        <button
          onClick={navigateToHome}
          className="text-xs font-bold text-slate-500 hover:text-[#002147] transition-colors mb-2 uppercase tracking-wider cursor-pointer"
        >
          ← Beranda
        </button>
        {/* Header Banner */}
        <div className="bg-[#002147] text-white rounded-3xl p-8 lg:p-12 relative overflow-hidden shadow-lg mb-8">
          <div className="absolute right-0 bottom-0 opacity-10 translate-x-12 translate-y-12">
            <Newspaper className="h-64 w-64" />
          </div>
          <div className="relative z-10 space-y-3">
            <span className="bg-amber-400 text-[#002147] text-[10px] font-extrabold px-3 py-1 rounded-full uppercase tracking-wider">
              Informasi Warta
            </span>
            <h1 className="text-3xl lg:text-4xl font-extrabold leading-tight">Daftar Berita PPID</h1>
            <p className="text-xs lg:text-sm text-slate-200 leading-relaxed font-medium">
              Rilis informasi terkini, pemberitahuan layanan, serta liputan warta keterbukaan informasi di lingkungan kampus Universitas Perintis Indonesia.
            </p>
          </div>
        </div>
      </div>

      {(newsFilterSearch || newsFilterCategory) && (
        <div className="flex items-center gap-3 bg-amber-50 border border-amber-200 rounded-2xl p-4 text-xs font-semibold text-[#002147] animate-in fade-in duration-200">
          <span>
            Menampilkan berita dengan
            {newsFilterCategory && <> kategori <strong className="capitalize">"{newsFilterCategory}"</strong></>}
            {newsFilterSearch && <>{newsFilterCategory && ' dan'} kata kunci <strong>"{newsFilterSearch}"</strong></>}
          </span>
          <button
            onClick={() => { setNewsFilterSearch(''); setNewsFilterCategory(''); }}
            className="ml-auto px-3.5 py-1.5 bg-amber-450 hover:bg-amber-500 text-[#002147] font-bold rounded-xl transition-colors cursor-pointer shadow-sm text-[10px] uppercase tracking-wider"
          >
            Hapus Filter
          </button>
        </div>
      )}

      {/* View Switcher Bar */}
      <div className="flex items-center justify-between border-b border-slate-100 pb-3">
        <h2 className="text-xs font-black text-[#002147] uppercase tracking-wider">
          Arsip Warta Berita
        </h2>
        
        {/* Toggle Buttons */}
        <div className="bg-slate-100 p-1 rounded-xl flex items-center gap-1 border border-slate-200 shadow-sm shrink-0">
          <button
            type="button"
            onClick={() => setViewMode('grid')}
            className={`p-1.5 rounded-lg transition-all cursor-pointer border-0 ${
              viewMode === 'grid'
                ? 'bg-[#002147] text-white shadow-sm'
                : 'text-slate-505 hover:text-slate-800 bg-transparent'
            }`}
            title="Tampilan Kolom (Grid)"
          >
            <LayoutGrid className="h-4 w-4" />
          </button>
          <button
            type="button"
            onClick={() => setViewMode('list')}
            className={`p-1.5 rounded-lg transition-all cursor-pointer border-0 ${
              viewMode === 'list'
                ? 'bg-[#002147] text-white shadow-sm'
                : 'text-slate-505 hover:text-slate-800 bg-transparent'
            }`}
            title="Tampilan Daftar (List)"
          >
            <List className="h-4 w-4" />
          </button>
        </div>
      </div>

      {isPostsLoading ? (
        <div className={viewMode === 'grid' ? "grid grid-cols-2 md:grid-cols-3 gap-3 sm:gap-6" : "grid grid-cols-1 gap-4"}>
          {[1, 2, 3, 4, 5, 6].map(i => (
            <div key={i} className={`bg-slate-100 rounded-3xl animate-pulse ${viewMode === 'list' ? 'h-28 sm:h-44 w-full' : 'h-48 sm:h-64'}`} />
          ))}
        </div>
      ) : posts.length > 0 ? (
        <div className={viewMode === 'grid' ? "grid grid-cols-2 md:grid-cols-3 gap-3 sm:gap-6" : "grid grid-cols-1 gap-4"}>
          {posts.map((post) => (
            <NewsCard
              key={post.id}
              post={post}
              viewMode={viewMode}
              navigateToNewsDetail={navigateToNewsDetail}
            />
          ))}
        </div>
      ) : (
        <div className="py-16 text-center text-slate-400 text-xs font-medium">

          Belum ada berita dirilis.
        </div>
      )}
    </div>
  );
}
