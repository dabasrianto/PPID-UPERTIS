import { Network } from 'lucide-react';
import { preprocessPostContent } from '../utils/helpers';
import type { PageData } from '../types';

interface StrukturOrganisasiProps {
  pageData: PageData;
}

export default function StrukturOrganisasi({ pageData }: StrukturOrganisasiProps) {
  return (
    <div className="space-y-6 animate-in fade-in duration-150 text-left">
      <div className="border-b border-slate-100 pb-4">
        <h2 className="text-lg font-extrabold text-[#002147] flex items-center gap-2">
          <Network className="h-5 w-5 text-amber-500" /> {pageData.title || 'Struktur Organisasi PPID'}
        </h2>
        {pageData.subtitle ? (
          <p className="text-xs text-slate-400 mt-0.5 font-medium">{pageData.subtitle}</p>
        ) : (
          <p className="text-xs text-slate-400 mt-0.5 font-medium">Bagan keanggotaan dan alur koordinasi Pejabat Pengelola Informasi.</p>
        )}
      </div>

      <div className="p-4 bg-slate-50 border border-slate-150 rounded-3xl space-y-4">
        <div className="flex justify-center">
          <div className="bg-[#002147] text-white p-3 rounded-2xl shadow border border-white/10 text-center w-56">
            <div className="text-[8px] text-amber-300 font-extrabold uppercase tracking-widest font-mono">ATASAN PPID</div>
            <div className="text-[11px] font-extrabold mt-0.5 leading-snug">REKTOR UPERTIS</div>
          </div>
        </div>

        <div className="flex justify-center"><div className="h-6 w-0.5 bg-slate-350" /></div>

        <div className="flex justify-center">
          <div className="bg-amber-400 text-[#002147] p-3 rounded-2xl shadow text-center w-56">
            <div className="text-[8px] font-extrabold uppercase tracking-widest font-mono text-[#002147]/60">TIM PERTIMBANGAN</div>
            <div className="text-[11px] font-extrabold mt-0.5 leading-snug">WAKIL REKTOR / DEKAN</div>
          </div>
        </div>

        <div className="flex justify-center"><div className="h-6 w-0.5 bg-slate-350" /></div>

        <div className="flex justify-center">
          <div className="bg-[#002147]/5 text-[#002147] border-2 border-dashed border-[#002147]/30 p-3.5 rounded-2xl text-center w-64">
            <div className="text-[8px] text-slate-400 font-extrabold uppercase tracking-widest font-mono">PPID UTAMA</div>
            <div className="text-xs font-black mt-0.5 leading-normal">KEPALA HUMAS & PROTOKOLER</div>
          </div>
        </div>

        <div className="flex justify-center"><div className="h-6 w-0.5 bg-slate-355" /></div>

        <div className="grid grid-cols-2 gap-4 max-w-md mx-auto">
          <div className="bg-white border border-slate-200 p-3 rounded-2xl shadow-sm text-center">
            <div className="text-[8px] text-slate-400 font-extrabold uppercase font-mono">BIDANG LAYANAN</div>
            <div className="text-[10px] font-bold text-slate-700 mt-0.5">Petugas Informasi / Loket</div>
          </div>
          <div className="bg-white border border-slate-200 p-3 rounded-2xl shadow-sm text-center">
            <div className="text-[8px] text-slate-400 font-extrabold uppercase font-mono">PPID PELAKSANA</div>
            <div className="text-[10px] font-bold text-slate-700 mt-0.5">Fakultas / Unit Kerja</div>
          </div>
        </div>
      </div>

      {pageData?.content && (
        <div
          className="html-content text-xs text-slate-650 leading-relaxed space-y-4 pt-4 border-t border-slate-100"
          dangerouslySetInnerHTML={{ __html: preprocessPostContent(pageData.content) }}
        />
      )}
    </div>
  );
}
