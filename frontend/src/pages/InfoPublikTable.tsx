import {
  Clock, Calendar, ShieldAlert, ShieldCheck, X, Layers, FileDown, Info, FileText
} from 'lucide-react';
import DownloadTable from '../components/DownloadTable';
import { resolveImageUrl, preprocessPostContent } from '../utils/helpers';
import type { PageData } from '../types';

interface InfoPublikTableProps {
  activeSlug: string;
  pageData: PageData;
  searchTerm: string;
  setSearchTerm: (term: string) => void;
}

export default function InfoPublikTable({
  activeSlug,
  pageData,
  searchTerm,
  setSearchTerm
}: InfoPublikTableProps) {
  const rawContent = pageData?.content || '';
  let docs: Array<{ title: string; description?: string; file_url: string }> = [];
  let introText = '';
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
      isJson = true;
    }
  } catch (e) {
    isJson = false;
    introText = rawContent;
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

      {/* Custom Description text */}
      {introText && introText.trim().length > 0 && (
        <div className="bg-white border border-slate-200 rounded-3xl p-6 lg:p-8 shadow-sm">
          <h3 className="text-xs font-bold text-[#002147] uppercase tracking-wider mb-3">Pengantar & Ketentuan</h3>
          <div
            className="html-content text-xs lg:text-sm text-slate-655 leading-relaxed space-y-3"
            dangerouslySetInnerHTML={{ __html: preprocessPostContent(introText) }}
          />
        </div>
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
    </div>
  );
}
