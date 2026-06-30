import React, { useState } from 'react';
import {
  Clock, Calendar, ShieldAlert, ShieldCheck, X, Layers, FileDown, Info, FileText,
  Building, CheckCircle, ClipboardList, Users, User, BookOpen, Book, Link, TrendingUp, Briefcase, HelpCircle
} from 'lucide-react';
import DownloadTable from '../components/DownloadTable';
import { resolveImageUrl, preprocessPostContent } from '../utils/helpers';
import type { PageData } from '../types';

interface InfoPublikTableProps {
  activeSlug: string;
  pageData: PageData;
  searchTerm: string;
  setSearchTerm: (term: string) => void;
  dbDownloads?: any[];
}

const iconMap: Record<string, React.ComponentType<any>> = {
  building: Building,
  check: CheckCircle,
  clipboard: ClipboardList,
  users: Users,
  user: User,
  bookOpen: BookOpen,
  book: Book,
  calendar: Calendar,
  fileText: FileText,
  link: Link,
  trendingUp: TrendingUp,
  briefcase: Briefcase,
  help: HelpCircle
};

export default function InfoPublikTable({
  activeSlug,
  pageData,
  searchTerm,
  setSearchTerm,
  dbDownloads = []
}: InfoPublikTableProps) {
  const [activeLightboxImage, setActiveLightboxImage] = useState<string | null>(null);
  const rawContent = pageData?.content || '';
  let docs: Array<{ title: string; description?: string; file_url: string }> = [];
  let introText = '';
  let sections: Array<{ text: string; imageUrl?: string; imageUrls?: string[]; imagePosition?: 'left' | 'right' }> = [];
  let isJson = false;

  try {
    const parsed = JSON.parse(rawContent);
    if (parsed && typeof parsed === 'object') {
      if (Array.isArray(parsed.docs)) {
        docs = parsed.docs;
      } else if (Array.isArray(parsed)) {
        docs = parsed;
      }
      introText = parsed.intro || '';
      if (Array.isArray(parsed.sections)) {
        sections = parsed.sections;
      }
      isJson = true;
    }
  } catch (e) {
    isJson = false;
    introText = rawContent;
  }

  // Determine if it is a synchronized DIP category page
  const dipCategories: Record<string, string> = {
    'informasi-publik-berkala': 'ppid-berkala',
    'informasi-tersedia-setiap-saat': 'ppid-setiap-saat',
    'info-serta-merta': 'ppid-serta-merta',
    'informasi-dikecualikan': 'ppid-dikecualikan'
  };

  const targetCategory = dipCategories[activeSlug];
  if (targetCategory) {
    docs = dbDownloads
      .filter((item: any) => item.category === targetCategory || item.category === targetCategory.replace('ppid-', ''))
      .map((item: any) => ({
        id: item.id,
        title: item.title,
        description: item.description,
        file_url: item.file_url
      }));
    isJson = true; // Force structured table layout for synced categories
  }

  // Fallback to normal rendering if it is not JSON
  if (!isJson) {
    return (
      <article className="space-y-6 text-left w-full">
        <div className="bg-[#002147] text-white rounded-3xl p-8 lg:p-12 relative overflow-hidden shadow-lg">
          <div className="absolute right-0 bottom-0 opacity-10 translate-x-12 translate-y-12">
            <FileText className="h-64 w-64" />
          </div>
          <div className="relative z-10 space-y-3">
            <span className="bg-amber-400 text-[#002147] text-[10px] font-extrabold px-3 py-1 rounded-full uppercase tracking-wider">
              Informasi Publik
            </span>
            <h1 className="text-3xl lg:text-4xl font-extrabold leading-tight">{pageData?.title || 'Informasi Publik'}</h1>
          </div>
        </div>
        {pageData?.cover_image_url && (
          <img
            src={pageData.cover_image_url}
            alt={pageData.title}
            className="w-full h-72 lg:h-96 object-cover rounded-3xl shadow-sm border border-slate-200 mt-6"
          />
        )}
        <div
          className="border-t border-slate-200 pt-6 html-content text-sm text-slate-700 leading-relaxed space-y-4"
          dangerouslySetInnerHTML={{ __html: preprocessPostContent(rawContent) }}
        />
      </article>
    );
  }

  const filteredDocs = docs.filter(item => {
    const matchesSearch = item.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
      (item.description && item.description.toLowerCase().includes(searchTerm.toLowerCase()));
    return matchesSearch;
  });

  // Banner settings
  const bannerInfo = {
    'informasi-publik-berkala': {
      tag: 'Informasi Berkala',
      desc: 'Daftar informasi publik yang wajib disediakan dan diumumkan secara berkala sesuai dengan UU Nomor 14 Tahun 2008 Pasal 9.',
      icon: Clock
    },
    'informasi-tersedia-setiap-saat': {
      tag: 'Setiap Saat',
      desc: 'Daftar informasi publik yang wajib disediakan dan sedia setiap saat oleh Badan Publik sesuai UU KIP Pasal 11.',
      icon: Calendar
    },
    'info-serta-merta': {
      tag: 'Serta Merta',
      desc: 'Daftar informasi publik yang wajib diumumkan secara serta merta menyangkut hajat hidup orang banyak dan ketertiban umum.',
      icon: ShieldAlert
    },
    'informasi-dikecualikan': {
      tag: 'Dikecualikan',
      desc: 'Daftar informasi publik yang dikecualikan dan bersifat rahasia berdasarkan pengujian konsekuensi UU KIP Pasal 17.',
      icon: X
    },
    'zona-integrasi': {
      tag: 'Zona Integrasi',
      desc: 'Daftar dokumen keterbukaan informasi publik dan arsip program Zona Integrasi Wilayah Bebas Korupsi (WBK) UPERTIS.',
      icon: ShieldCheck
    }
  }[activeSlug as 'informasi-publik-berkala' | 'informasi-tersedia-setiap-saat' | 'info-serta-merta' | 'informasi-dikecualikan' | 'zona-integrasi'] || {
    tag: 'Informasi Publik',
    desc: 'Daftar informasi publik resmi Universitas Perintis Indonesia.',
    icon: Layers
  };

  const BannerIcon = bannerInfo.icon;

  return (
    <div className="space-y-8 animate-in fade-in duration-200 text-left w-full">
      {/* Hero Banner */}
      <div className="bg-[#002147] text-white rounded-3xl p-8 lg:p-12 relative overflow-hidden shadow-lg">
        <div className="absolute right-0 bottom-0 opacity-10 translate-x-12 translate-y-12">
          <BannerIcon className="h-64 w-64 text-amber-400" />
        </div>
        <div className="absolute top-0 left-0 w-full h-full bg-gradient-to-br from-transparent via-transparent to-amber-400/10 pointer-events-none" />
        <div className="relative z-10 space-y-4">
          <span className="bg-amber-400 text-[#002147] text-[10px] font-extrabold px-3 py-1 rounded-full uppercase tracking-wider inline-block">
            {bannerInfo.tag}
          </span>
          <h1 className="text-3xl lg:text-4xl font-extrabold leading-tight">{pageData?.title || 'Informasi Publik'}</h1>
          <p className="text-xs lg:text-sm text-slate-350 leading-relaxed font-medium max-w-2xl">
            {bannerInfo.desc}
          </p>
        </div>
      </div>

      {/* Custom Description sections */}
      {sections.length > 0 ? (
        <div className="space-y-6">
          {sections.map((section, sIdx) => {
            const images = Array.isArray(section.imageUrls)
              ? section.imageUrls
              : (section.imageUrl ? [section.imageUrl] : []);
            const hasImage = images.length > 0;
            const isImageLeft = section.imagePosition === 'left';

            return (
              <div
                key={sIdx}
                className="bg-white border border-slate-200 rounded-3xl p-6 lg:p-8 shadow-sm flex flex-col lg:flex-row items-center gap-8 text-left"
              >
                {/* Text Column */}
                <div
                  className={`w-full ${hasImage ? 'lg:w-1/2' : 'w-full'} space-y-3 ${
                    hasImage ? (isImageLeft ? 'lg:order-2' : 'lg:order-1') : ''
                  }`}
                >
                  {sIdx === 0 && !section.title && (
                    <h3 className="text-xs font-bold text-[#002147] uppercase tracking-wider mb-2">Pengantar & Ketentuan</h3>
                  )}
                  {section.title && (
                    <h3 className="text-xs font-extrabold text-[#002147] uppercase tracking-wider mb-2">
                      {section.title}
                    </h3>
                  )}
                  {section.subtitle && (
                    <p className="text-[11px] text-slate-400 mb-4 leading-relaxed font-medium">
                      {section.subtitle}
                    </p>
                  )}
                  {section.cards && section.cards.length > 0 && (
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mt-2 mb-4">
                      {section.cards.map((card: any, cIdx: number) => {
                        const CardIcon = iconMap[card.icon] || Link;
                        return (
                          <a
                            key={cIdx}
                            href={card.url}
                            className="flex items-center gap-3 p-4 bg-slate-50 hover:bg-blue-50 border border-slate-200 rounded-2xl transition-all group"
                          >
                            <div className="p-2.5 bg-white text-blue-700 rounded-xl shadow-sm group-hover:bg-[#002147] group-hover:text-white transition-all shrink-0">
                              <CardIcon className="w-4 h-4" />
                            </div>
                            <div>
                              <span className="text-xs font-bold text-slate-800 block">
                                {card.title}
                              </span>
                              <span className="text-[10px] text-slate-400 block mt-0.5 leading-normal">
                                {card.subtitle}
                              </span>
                            </div>
                          </a>
                        );
                      })}
                    </div>
                  )}
                  {section.text && section.text.trim().length > 0 && (
                    <div
                      className="html-content text-xs lg:text-sm text-slate-655 leading-relaxed space-y-3"
                      dangerouslySetInnerHTML={{ __html: preprocessPostContent(section.text) }}
                    />
                  )}
                </div>

                {/* Image Column */}
                {hasImage && (
                  <div
                    className={`w-full lg:w-1/2 flex justify-center py-6 lg:py-0 ${
                      isImageLeft ? 'lg:order-1' : 'lg:order-2'
                    }`}
                  >
                    {images.length === 1 ? (
                      <img
                        src={resolveImageUrl(images[0])}
                        alt={`Visual ${sIdx + 1}`}
                        onClick={() => setActiveLightboxImage(resolveImageUrl(images[0]))}
                        className="w-full h-auto max-h-[300px] object-cover rounded-2xl shadow-sm border border-slate-100 hover:scale-[1.02] cursor-pointer hover:shadow-md transition-all duration-305"
                      />
                    ) : (
                      <div className="relative w-[340px] h-[340px] select-none group flex shrink-0 justify-center items-center">
                        {/* Gambar 1: Kiri/Belakang (Vertikal) */}
                        {images[0] && (
                          <div 
                            onClick={() => setActiveLightboxImage(resolveImageUrl(images[0]))}
                            className="absolute top-0 left-0 w-[190px] h-[240px] rounded-[2.5rem] overflow-hidden border-4 border-white shadow-lg z-10 transition-all duration-500 hover:scale-102 cursor-pointer hover:shadow-2xl hover:z-40 animate-float-slow group-hover:-translate-x-3 group-hover:-translate-y-1"
                          >
                            <img
                              src={resolveImageUrl(images[0])}
                              alt="Visual 1"
                              className="w-full h-full object-cover"
                            />
                          </div>
                        )}

                        {/* Gambar 2: Kanan Atas (Persegi) */}
                        {images[1] && (
                          <div 
                            onClick={() => setActiveLightboxImage(resolveImageUrl(images[1]))}
                            className="absolute top-6 left-[145px] w-[130px] h-[130px] rounded-[2rem] overflow-hidden border-4 border-white shadow-md z-25 transition-all duration-500 hover:scale-105 cursor-pointer hover:shadow-2xl hover:z-40 animate-float-medium group-hover:translate-x-3 group-hover:-translate-y-2"
                          >
                            <img
                              src={resolveImageUrl(images[1])}
                              alt="Visual 2"
                              className="w-full h-full object-cover"
                            />
                          </div>
                        )}

                        {/* Gambar 3: Tengah/Depan Bawah (Horizontal) */}
                        {images[2] && (
                          <div 
                            onClick={() => setActiveLightboxImage(resolveImageUrl(images[2]))}
                            className="absolute top-[135px] left-[65px] w-[230px] h-[180px] rounded-[2.5rem] overflow-hidden border-4 border-white shadow-xl z-30 transition-all duration-500 hover:scale-105 hover:shadow-2xl cursor-pointer hover:z-40 animate-float-fast group-hover:translate-y-2 group-hover:translate-x-1"
                          >
                            <img
                              src={resolveImageUrl(images[2])}
                              alt="Visual 3"
                              className="w-full h-full object-cover"
                            />
                          </div>
                        )}
                      </div>
                    )}
                  </div>
                )}
              </div>
            );
          })}
        </div>
      ) : (
        introText && introText.trim().length > 0 && (
          <div className="bg-white border border-slate-200 rounded-3xl p-6 lg:p-8 shadow-sm">
            <h3 className="text-xs font-bold text-[#002147] uppercase tracking-wider mb-3">Pengantar & Ketentuan</h3>
            <div
              className="html-content text-xs lg:text-sm text-slate-655 leading-relaxed space-y-3"
              dangerouslySetInnerHTML={{ __html: preprocessPostContent(introText) }}
            />
          </div>
        )
      )}

      <DownloadTable
        items={filteredDocs.map(item => ({
          ...item,
          type: 'DownloadItem'
        }))}
        searchTerm={searchTerm}
        setSearchTerm={setSearchTerm}
        incrementDownloadCount={(id, fileUrl) => window.open(resolveImageUrl(fileUrl), '_blank')}
      />

      {/* Footer Legal Note */}
      <div className="bg-slate-50 border border-slate-200 rounded-2xl p-6 flex items-start gap-4 text-left">
        <Info className="h-5 w-5 text-amber-500 shrink-0 mt-0.5" />
        <div className="space-y-1">
          <h4 className="text-xs font-bold text-[#002147]">Ketentuan Publikasi Dokumen Resmi</h4>
          <p className="text-[11px] text-slate-500 font-medium leading-relaxed">
            Dokumen-dokumen di atas diperbarui secara berkala oleh pejabat pengelola informasi dan dokumentasi (PPID) Universitas Perintis Indonesia.
          </p>
        </div>
      </div>

      {/* Interactive Image Lightbox Modal */}
      {activeLightboxImage && (
        <div 
          onClick={() => setActiveLightboxImage(null)}
          className="fixed inset-0 bg-slate-950/80 backdrop-blur-md z-[9999] flex items-center justify-center p-4 transition-all duration-300 animate-in fade-in cursor-zoom-out"
        >
          <div className="relative max-w-4xl max-h-[85vh] overflow-hidden rounded-2xl border border-white/10 shadow-2xl animate-in zoom-in-95 duration-200" onClick={(e) => e.stopPropagation()}>
            <button 
              onClick={() => setActiveLightboxImage(null)}
              className="absolute top-4 right-4 bg-black/60 hover:bg-black/80 text-white rounded-full p-2.5 transition-colors border border-white/10 cursor-pointer z-50"
              title="Tutup Zoom"
            >
              <X className="h-5 w-5" />
            </button>
            <img 
              src={activeLightboxImage} 
              alt="Zoomed preview" 
              className="max-w-full max-h-[80vh] object-contain rounded-2xl"
            />
          </div>
        </div>
      )}
    </div>
  );
}
