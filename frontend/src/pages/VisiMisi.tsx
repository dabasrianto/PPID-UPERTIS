import { Target } from 'lucide-react';
import { preprocessPostContent } from '../utils/helpers';
import type { PageData } from '../types';

interface VisiMisiProps {
  pageData: PageData;
}

export default function VisiMisi({ pageData }: VisiMisiProps) {
  return (
    <div className="space-y-6 animate-in fade-in duration-150 text-left">
      <div className="border-b border-slate-100 pb-4">
        <h2 className="text-lg font-extrabold text-[#002147] flex items-center gap-2">
          <Target className="h-5 w-5 text-amber-500" /> {pageData.title || 'Visi & Misi PPID'}
        </h2>
        {pageData.subtitle ? (
          <p className="text-xs text-slate-400 mt-0.5 font-medium">{pageData.subtitle}</p>
        ) : (
          <p className="text-xs text-slate-400 mt-0.5 font-medium">Arah panduan dan misi strategis pelayanan informasi.</p>
        )}
      </div>

      {pageData?.content && (
        <div
          className="html-content text-xs text-slate-650 leading-relaxed space-y-4 pt-2"
          dangerouslySetInnerHTML={{ __html: preprocessPostContent(pageData.content) }}
        />
      )}
    </div>
  );
}
