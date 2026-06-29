import { ShieldCheck } from 'lucide-react';
import { preprocessPostContent } from '../utils/helpers';
import type { PageData } from '../types';

interface MaklumatProps {
  pageData: PageData;
}

export default function Maklumat({ pageData }: MaklumatProps) {
  return (
    <div className="space-y-6 animate-in fade-in duration-150">
      <div className="border-b border-slate-100 pb-4 text-left">
        <h2 className="text-lg font-extrabold text-[#002147] flex items-center gap-2">
          <ShieldCheck className="h-5 w-5 text-amber-500" /> {pageData.title || 'Maklumat Pelayanan PPID'}
        </h2>
        {pageData.subtitle ? (
          <p className="text-xs text-slate-400 mt-0.5 font-medium">{pageData.subtitle}</p>
        ) : (
          <p className="text-xs text-slate-400 mt-0.5 font-medium">Ikrar janji pelayanan informasi publik yang akuntabel.</p>
        )}
      </div>

      <div className="relative border-4 border-double border-amber-200/80 rounded-3xl p-6 lg:p-8 bg-slate-50/50 text-center space-y-4 shadow-inner max-w-xl mx-auto animate-in fade-in zoom-in-95 duration-200">
        <h3 className="font-extrabold text-[#002147] tracking-wider text-xs uppercase">MAKLUMAT PELAYANAN INFORMASI</h3>
        <p className="text-slate-700 italic font-serif leading-relaxed text-xs max-w-md mx-auto">
          "Kami berkomitmen memberikan pelayanan informasi publik yang cepat, tepat, transparan, dan akuntabel sesuai dengan standar operasional prosedur demi mewujudkan keterbukaan informasi publik di lingkungan Universitas Perintis Indonesia."
        </p>
      </div>

      {pageData?.content && (
        <div
          className="html-content text-xs text-slate-650 leading-relaxed space-y-4 pt-4 border-t border-slate-100 text-left"
          dangerouslySetInnerHTML={{ __html: preprocessPostContent(pageData.content) }}
        />
      )}
    </div>
  );
}
