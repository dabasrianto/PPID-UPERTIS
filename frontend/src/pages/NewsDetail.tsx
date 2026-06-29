import { useState } from 'react';
import { Newspaper, Search, Image } from 'lucide-react';
import { resolveImageUrl, preprocessPostContent } from '../utils/helpers';
import type { Post } from '../types';

interface NewsDetailProps {
  selectedPost: Post | null;
  isPageLoading: boolean;
  sidebarSearch: string;
  setSidebarSearch: (val: string) => void;
  posts: Post[];
  navigateToNews: () => void;
  setCurrentPage: (page: string) => void;
  setNewsFilterSearch: (search: string) => void;
  setNewsFilterCategory: (cat: string) => void;
}

export default function NewsDetail({
  selectedPost,
  isPageLoading,
  sidebarSearch,
  setSidebarSearch,
  posts,
  navigateToNews,
  setCurrentPage,
  setNewsFilterSearch,
  setNewsFilterCategory
}: NewsDetailProps) {
  const [copySuccess, setCopySuccess] = useState<boolean>(false);

  return (
    <div className="container mx-auto px-4 lg:px-8 py-12 lg:py-16 max-w-7xl text-left flex-1 space-y-6">
      <button
        onClick={navigateToNews}
        className="text-xs font-bold text-slate-500 hover:text-[#002147] transition-colors mb-2 uppercase tracking-wider cursor-pointer"
      >
        ← Kembali ke Berita
      </button>

      {isPageLoading ? (
        <div className="space-y-4 animate-pulse">
          <div className="h-8 bg-slate-200 rounded-lg w-3/4" />
          <div className="h-4 bg-slate-200 rounded-lg w-1/4" />
          <div className="h-80 bg-slate-200 rounded-2xl w-full" />
        </div>
      ) : selectedPost ? (
        <div className="space-y-6">
          {/* Premium Header Banner */}
          <div className="bg-[#002147] text-white rounded-3xl p-8 lg:p-12 relative overflow-hidden shadow-lg text-left">
            <div className="absolute right-0 bottom-0 opacity-10 translate-x-12 translate-y-12">
              <Newspaper className="h-64 w-64" />
            </div>
            <div className="relative z-10 space-y-4">
              <span className="bg-amber-400 text-[#002147] text-[10px] font-extrabold px-3 py-1 rounded-full uppercase tracking-wider inline-block">
                {selectedPost.category || 'Berita PPID'}
              </span>
              <h1 className="text-2xl lg:text-3xl font-extrabold leading-snug">
                {selectedPost.title}
              </h1>
              <div className="flex items-center gap-2 text-[10px] text-slate-350 font-semibold uppercase tracking-wider pt-2 border-t border-white/10">
                <span>Oleh Admin PPID</span>
                <span>•</span>
                <span>
                  {new Date(selectedPost.created_at).toLocaleDateString('id-ID', {
                    day: 'numeric',
                    month: 'long',
                    year: 'numeric',
                    hour: '2-digit',
                    minute: '2-digit'
                  })}
                </span>
              </div>
            </div>
          </div>

          {/* Grid Layout: Main Content (lg:col-span-2) + Sidebar (lg:col-span-1) */}
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
            {/* Main Content */}
            <div className="lg:col-span-2 space-y-6">
              <article className="bg-white border border-slate-200 rounded-3xl p-6 lg:p-10 shadow-sm space-y-6 text-left">
                {(selectedPost.featured_image || selectedPost.cover_image_url) && (
                  <div className="rounded-2xl overflow-hidden max-h-[450px] shadow-sm border border-slate-100">
                    <img
                      src={resolveImageUrl(selectedPost.featured_image || selectedPost.cover_image_url)}
                      alt={selectedPost.title}
                      className="w-full h-full object-cover"
                    />
                  </div>
                )}

                {/* Bagikan Berita (Social Share Buttons) */}
                <div className="flex flex-col sm:flex-row sm:items-center justify-between border-b border-slate-100 pb-5 gap-4">
                  <div className="flex items-center gap-2">
                    <span className="text-[11px] font-bold uppercase tracking-wider text-slate-400">Bagikan:</span>
                    <div className="flex items-center gap-2 flex-wrap">
                      {/* WhatsApp */}
                      <a
                        href={`https://api.whatsapp.com/send?text=${encodeURIComponent(selectedPost.title + ' - ' + window.location.href)}`}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="flex items-center justify-center p-2.5 rounded-full bg-emerald-50 text-emerald-600 hover:bg-emerald-600 hover:text-white transition-all shadow-sm border border-emerald-100 cursor-pointer"
                        title="Bagikan ke WhatsApp"
                      >
                        <svg className="h-4.5 w-4.5 fill-current" viewBox="0 0 24 24">
                          <path d="M17.472 14.382c-.022-.08-.078-.17-.189-.22-.11-.05-.66-.32-.76-.36-.1-.03-.18-.05-.25.05-.07.1-.29.36-.36.43-.07.07-.14.08-.25.03-.11-.05-.47-.17-.89-.55-.33-.29-.55-.65-.62-.76-.07-.11-.01-.17.05-.23.05-.05.11-.13.17-.2.06-.07.08-.12.12-.2.04-.07.02-.13-.01-.18-.03-.05-.25-.6-.35-.84-.1-.24-.2-.2-.25-.2-.05 0-.12-.01-.19-.01-.07 0-.18.03-.27.14-.1.1-.36.36-.36.88 0 .52.38 1.03.43 1.1.05.07.75 1.15 1.82 1.61.25.11.45.18.61.23.27.08.52.07.72.04.22-.03.66-.27.76-.53.1-.26.1-.49.07-.53-.03-.05-.09-.08-.2-.13zM12 2C6.477 2 2 6.477 2 12c0 2.01.59 3.88 1.61 5.45l-1.07 3.93 4.02-1.05C8.1 21.37 9.99 22 12 22c5.523 0 10-4.477 10-10S17.523 2 12 2zm0 18c-1.85 0-3.58-.53-5.06-1.46l-.36-.21-2.43.64.65-2.38-.24-.38C3.62 14.85 3 13.13 3 12c0-4.97 4.03-9 9-9s9 4.03 9 9-4.03 9-9 9z"/>
                        </svg>
                      </a>

                      {/* Facebook */}
                      <a
                        href={`https://www.facebook.com/sharer/sharer.php?u=${encodeURIComponent(window.location.href)}`}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="flex items-center justify-center p-2.5 rounded-full bg-blue-50 text-blue-600 hover:bg-blue-600 hover:text-white transition-all shadow-sm border border-blue-100 cursor-pointer"
                        title="Bagikan ke Facebook"
                      >
                        <svg className="h-4.5 w-4.5 fill-current" viewBox="0 0 24 24">
                          <path d="M9 8h-3v4h3v12h5v-12h3.642l.358-4h-4v-1.667c0-.955.192-1.333 1.115-1.333h2.885v-5h-3.808c-3.596 0-5.192 1.583-5.192 4.615v3.385z"/>
                        </svg>
                      </a>

                      {/* Twitter / X */}
                      <a
                        href={`https://twitter.com/intent/tweet?text=${encodeURIComponent(selectedPost.title)}&url=${encodeURIComponent(window.location.href)}`}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="flex items-center justify-center p-2.5 rounded-full bg-slate-50 text-slate-800 hover:bg-black hover:text-white transition-all shadow-sm border border-slate-200 cursor-pointer"
                        title="Bagikan ke Twitter / X"
                      >
                        <svg className="h-4.5 w-4.5 fill-current" viewBox="0 0 24 24">
                          <path d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-5.214-6.817L4.99 21.75H1.68l7.73-8.835L1.254 2.25H8.08l4.713 6.231zm-1.161 17.52h1.833L7.084 4.126H5.117z"/>
                        </svg>
                      </a>

                      {/* Telegram */}
                      <a
                        href={`https://t.me/share/url?url=${encodeURIComponent(window.location.href)}&text=${encodeURIComponent(selectedPost.title)}`}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="flex items-center justify-center p-2.5 rounded-full bg-sky-50 text-sky-600 hover:bg-sky-500 hover:text-white transition-all shadow-sm border border-sky-100 cursor-pointer"
                        title="Bagikan ke Telegram"
                      >
                        <svg className="h-4.5 w-4.5 fill-current" viewBox="0 0 24 24">
                          <path d="M22.05 3.577L1.137 11.64c-.287.113-.284.275-.05.347l5.372 1.677 1.996 6.136c.243.67.123.936.836.936.55 0 .794-.253 1.1-.55l2.64-2.565 5.5 4.062c1.01.557 1.74.27 1.99-.938l3.6-16.98c.368-1.477-.566-2.148-1.53-1.743z"/>
                        </svg>
                      </a>

                      {/* Salin Tautan */}
                      <button
                        onClick={() => {
                          navigator.clipboard.writeText(window.location.href);
                          setCopySuccess(true);
                          setTimeout(() => setCopySuccess(false), 2000);
                        }}
                        className="flex items-center justify-center p-2.5 rounded-full bg-slate-50 text-slate-600 hover:bg-[#002147] hover:text-white transition-all shadow-sm border border-slate-200 cursor-pointer"
                        title="Salin Tautan Berita"
                      >
                        <svg className="h-4.5 w-4.5 stroke-current fill-none" strokeWidth="2.5" viewBox="0 0 24 24" strokeLinecap="round" strokeLinejoin="round">
                          <path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"/>
                          <path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"/>
                        </svg>
                      </button>
                    </div>
                  </div>

                  {copySuccess && (
                    <span className="text-[10px] text-emerald-600 bg-emerald-50 border border-emerald-100 px-3 py-1.5 rounded-full font-bold self-start sm:self-center animate-in fade-in slide-in-from-top-1 duration-200">
                      ✓ Link berhasil disalin!
                    </span>
                  )}
                </div>

                <div
                  className="html-content text-xs text-slate-650 leading-relaxed space-y-4 font-medium whitespace-pre-wrap"
                  dangerouslySetInnerHTML={{ __html: preprocessPostContent(selectedPost.content || '') }}
                />
              </article>
            </div>

            {/* Sidebar */}
            <div className="lg:col-span-1 space-y-6">
              {/* Search Widget */}
              <div className="bg-white border border-slate-200 rounded-3xl p-6 shadow-sm space-y-4 text-left">
                <h3 className="text-xs font-bold text-[#002147] uppercase tracking-wider pb-1 border-b-2 border-amber-400 w-10">Cari</h3>
                <div className="relative flex items-center">
                  <input
                    type="text"
                    placeholder="Cari berita..."
                    value={sidebarSearch}
                    onChange={(e) => setSidebarSearch(e.target.value)}
                    onKeyDown={(e) => {
                      if (e.key === 'Enter') {
                        setCurrentPage('berita');
                        setNewsFilterSearch(sidebarSearch);
                        setNewsFilterCategory('');
                      }
                    }}
                    className="w-full pr-12 pl-4 py-2.5 bg-slate-50 border border-slate-200 rounded-full text-xs focus:outline-none focus:ring-2 focus:ring-amber-450/20 focus:border-amber-400 transition-all font-semibold text-slate-700"
                  />
                  <button
                    onClick={() => {
                      setCurrentPage('berita');
                      setNewsFilterSearch(sidebarSearch);
                      setNewsFilterCategory('');
                    }}
                    className="absolute right-1 p-2 bg-gradient-to-r from-amber-500 to-orange-500 hover:from-amber-600 hover:to-orange-600 text-white rounded-full transition-all shadow-md cursor-pointer flex items-center justify-center"
                  >
                    <Search className="h-3.5 w-3.5" />
                  </button>
                </div>
              </div>

              {/* Categories Widget */}
              <div className="bg-white border border-slate-200 rounded-3xl p-6 shadow-sm space-y-4 text-left">
                <div className="space-y-1">
                  <h3 className="text-xs font-bold text-slate-800 uppercase tracking-wider">Kategori Berita</h3>
                  <div className="h-0.5 bg-gradient-to-r from-amber-500 to-orange-500 w-12" />
                </div>
                <div className="space-y-1">
                  {[
                    { name: 'Semua Berita PPID', count: posts.length }
                  ].map((cat) => (
                    <button
                      key={cat.name}
                      onClick={() => {
                        setCurrentPage('berita');
                        setNewsFilterCategory('');
                        setNewsFilterSearch('');
                      }}
                      className="w-full flex items-center justify-between py-2 text-xs font-semibold text-[#002147] transition-all border-b border-slate-50 hover:translate-x-1 cursor-pointer"
                    >
                      <div className="flex items-center gap-1.5">
                        <span className="text-slate-400 font-bold">&rsaquo;</span>
                        <span>{cat.name}</span>
                      </div>
                      <span className="bg-slate-100 text-slate-500 text-[9px] font-bold px-2 py-0.5 rounded-md min-w-[20px] text-center">
                        {cat.count}
                      </span>
                    </button>
                  ))}
                </div>
              </div>
            </div>
          </div>
        </div>
      ) : (
        <div className="text-center text-slate-400 py-16 text-xs font-bold">
          Gagal memuat berita detail.
        </div>
      )}
    </div>
  );
}
