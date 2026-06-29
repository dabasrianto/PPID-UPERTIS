import { BookOpen, ShieldCheck, ChevronDown, Check, ExternalLink, Download } from 'lucide-react';
import { parseRegulasiHTML } from '../utils/helpers';

interface RegulasiProps {
  pageContent: string;
  expandedRegulasi: Record<string, boolean>;
  toggleRegulasi: (key: string) => void;
}

export default function Regulasi({
  pageContent,
  expandedRegulasi,
  toggleRegulasi
}: RegulasiProps) {
  const regulasiToRender = parseRegulasiHTML(pageContent);

  return (
    <div className="space-y-8 animate-in fade-in duration-200 text-left w-full">
      {/* Header Banner */}
      <div className="bg-[#002147] text-white rounded-3xl p-8 lg:p-12 relative overflow-hidden shadow-lg">
        <div className="absolute right-0 bottom-0 opacity-10 translate-x-12 translate-y-12">
          <BookOpen className="h-64 w-64" />
        </div>
        <div className="relative z-10 space-y-3">
          <span className="bg-amber-400 text-[#002147] text-[10px] font-extrabold px-3 py-1 rounded-full uppercase tracking-wider">
            Landasan Hukum
          </span>
          <h1 className="text-3xl lg:text-4xl font-extrabold leading-tight">Regulasi KIP</h1>
          <p className="text-xs lg:text-sm text-slate-200 leading-relaxed font-medium">
            Landasan hukum, dasar regulasi, dan ketentuan Keterbukaan Informasi Publik di lingkungan Universitas Perintis Indonesia.
          </p>
        </div>
      </div>

      {/* Accordion List & Hak Pemohon (Full Width) */}
      <div className="space-y-6">
        <div className="space-y-4">
          {regulasiToRender.map((group: any) => {
            const isOpen = !!expandedRegulasi[group.key];
            return (
              <div key={group.key} className="border border-slate-200 rounded-3xl overflow-hidden bg-white shadow-sm transition-all duration-300">
                {/* Accordion Header */}
                <button
                  type="button"
                  onClick={() => toggleRegulasi(group.key)}
                  className={`w-full flex items-center justify-between p-5 text-left cursor-pointer transition-colors border-0 ${
                    isOpen ? 'bg-slate-50 border-b border-slate-100' : 'hover:bg-slate-50/50 bg-white'
                  }`}
                >
                  <div className="space-y-1 pr-4">
                    <h3 className="font-extrabold text-sm text-[#002147] tracking-wide flex items-center gap-2">
                      <span className="h-2 w-2 rounded-full bg-amber-500 shrink-0" />
                      {group.label}
                    </h3>
                    <p className="text-[11px] text-slate-400 font-medium">{group.desc}</p>
                  </div>
                  <div className={`h-8 w-8 rounded-xl bg-[#002147]/5 hover:bg-[#002147]/10 flex items-center justify-center text-[#002147] transition-transform duration-300 ${
                    isOpen ? 'rotate-180' : ''
                  }`}>
                    <ChevronDown className="h-4 w-4" />
                  </div>
                </button>

                {/* Accordion Content */}
                {isOpen && (
                  <div className="p-5 divide-y divide-slate-100 bg-white animate-in fade-in duration-200">
                    {(group.items || []).map((item: any, idx: number) => (
                      <div key={idx} className="py-4 first:pt-0 last:pb-0 flex flex-col md:flex-row md:items-center justify-between gap-4">
                        <div className="flex items-start gap-4">
                          <div className="h-7 w-7 bg-[#002147]/5 text-[#002147] font-extrabold text-xs rounded-lg flex items-center justify-center shrink-0 mt-0.5">
                            {item.no}
                          </div>
                          <div className="space-y-1">
                            <h4 className="text-xs font-extrabold text-slate-800 leading-snug">{item.title}</h4>
                            <p className="text-[11px] text-slate-500 leading-relaxed font-medium">{item.detail}</p>

                            {item.isSopList && item.sops && (
                              <div className="grid grid-cols-1 md:grid-cols-2 gap-2 mt-3 pt-2 border-t border-slate-100">
                                {item.sops.map((sop: string, sIdx: number) => (
                                  <div key={sIdx} className="flex items-center gap-2 text-[10px] text-slate-600 bg-slate-50 border border-slate-100 p-2 rounded-xl">
                                    <Check className="h-3 w-3 text-amber-500 shrink-0" />
                                    <span className="font-bold">{sop}</span>
                                  </div>
                                ))}
                              </div>
                            )}
                          </div>
                        </div>

                        {/* Action Buttons */}
                        <div className="shrink-0 flex items-center gap-2 mt-2 md:mt-0 md:pl-11">
                          {item.fileUrl && item.fileUrl !== '#' ? (
                            <>
                              <a
                                href={item.fileUrl}
                                target="_blank"
                                rel="noopener noreferrer"
                                className="inline-flex items-center gap-1.5 px-3.5 py-2 rounded-xl border border-slate-200 hover:border-amber-400 hover:bg-amber-50 text-[11px] font-bold text-slate-600 hover:text-amber-700 transition-colors shadow-sm cursor-pointer"
                              >
                                <ExternalLink className="h-3 w-3" /> Lihat PDF
                              </a>
                              <a
                                href={item.fileUrl}
                                download
                                className="inline-flex items-center gap-1.5 px-3.5 py-2 rounded-xl bg-[#002147] hover:bg-amber-440 hover:text-[#002147] text-[11px] font-bold text-white transition-all shadow-sm cursor-pointer"
                              >
                                <Download className="h-3 w-3" /> Unduh
                              </a>
                            </>
                          ) : (
                            <span className="text-[10px] bg-slate-50 text-slate-400 border border-slate-100 px-3 py-1.5 rounded-xl font-bold uppercase select-none">
                              Dokumen Kampus
                            </span>
                          )}
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            );
          })}
        </div>

        {/* Info Box: Hak-Hak Pemohon */}
        <div className="p-5 bg-gradient-to-br from-amber-50/50 to-orange-50/30 border border-amber-200/50 rounded-3xl space-y-3.5 mt-6">
          <h4 className="text-xs font-extrabold text-[#002147] uppercase tracking-wider flex items-center gap-1.5">
            <ShieldCheck className="h-4.5 w-4.5 text-amber-500" /> Hak Pemohon Informasi Publik
          </h4>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-3 text-[11px] text-slate-655 font-medium leading-relaxed">
            <div className="flex items-start gap-2 bg-white/60 p-3 rounded-2xl border border-slate-100">
              <span className="h-2 w-2 rounded-full bg-amber-400 shrink-0 mt-1.5" />
              <span>Hak melihat dan mengetahui informasi publik yang bersifat terbuka bagi masyarakat umum secara berkala.</span>
            </div>
            <div className="flex items-start gap-2 bg-white/60 p-3 rounded-2xl border border-slate-100">
              <span className="h-2 w-2 rounded-full bg-amber-400 shrink-0 mt-1.5" />
              <span>Hak menghadiri sidang Komisi Informasi yang terbuka untuk umum dalam proses penyelesaian sengketa.</span>
            </div>
            <div className="flex items-start gap-2 bg-white/60 p-3 rounded-2xl border border-slate-100">
              <span className="h-2 w-2 rounded-full bg-amber-400 shrink-0 mt-1.5" />
              <span>Hak mendapatkan salinan dokumen/informasi publik melalui permohonan tertulis resmi ke sekretariat PPID.</span>
            </div>
            <div className="flex items-start gap-2 bg-white/60 p-3 rounded-2xl border border-slate-100">
              <span className="h-2 w-2 rounded-full bg-amber-400 shrink-0 mt-1.5" />
              <span>Hak mengajukan keberatan kepada atasan PPID dan gugatan sengketa apabila hak informasi publiknya terlanggar.</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
