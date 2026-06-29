import { motion, AnimatePresence } from 'framer-motion';
import { Image, Search, X } from 'lucide-react';
import { resolveImageUrl } from '../utils/helpers';
import type { GalleryItem } from '../types';

interface GaleriProps {
  galleries: GalleryItem[];
  activeLightboxImage: GalleryItem | null;
  setActiveLightboxImage: (item: GalleryItem | null) => void;
}

export default function Galeri({
  galleries,
  activeLightboxImage,
  setActiveLightboxImage
}: GaleriProps) {
  return (
    <div className="space-y-8 animate-in fade-in duration-200 text-left w-full py-6">
      {/* Header Banner */}
      <div className="bg-[#002147] text-white rounded-3xl p-8 lg:p-12 relative overflow-hidden shadow-lg">
        <div className="absolute right-0 bottom-0 opacity-10 translate-x-12 translate-y-12">
          <Image className="h-64 w-64" />
        </div>
        <div className="relative z-10 space-y-3">
          <span className="bg-amber-400 text-[#002147] text-[10px] font-extrabold px-3 py-1 rounded-full uppercase tracking-wider">
            Dokumentasi Media
          </span>
          <h1 className="text-3xl lg:text-4xl font-extrabold leading-tight">Galeri Kegiatan PPID</h1>
          <p className="text-xs lg:text-sm text-slate-200 leading-relaxed font-medium">
            Dokumentasi foto kegiatan layanan informasi, sosialisasi keterbukaan informasi publik, serta publikasi visual PPID Universitas Perintis Indonesia.
          </p>
        </div>
      </div>

      {/* Galleries Grid */}
      <div className="bg-white border border-slate-200 rounded-3xl p-6 lg:p-8 shadow-sm">
        {galleries.length > 0 ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {galleries.map((item) => (
              <div
                key={item.id}
                className="group relative border border-slate-100 hover:border-blue-100 hover:bg-blue-50/5 rounded-2xl overflow-hidden shadow-sm hover:shadow-md transition-all flex flex-col h-full"
              >
                {/* Media Container */}
                <div
                  onClick={() => setActiveLightboxImage(item)}
                  className="h-52 bg-slate-100 overflow-hidden relative shrink-0 cursor-pointer"
                >
                  {item.media_url ? (
                    <img
                      src={resolveImageUrl(item.media_url)}
                      alt={item.title}
                      className="w-full h-full object-cover transition-transform duration-700 ease-in-out group-hover:scale-105"
                    />
                  ) : (
                    <div className="w-full h-full flex items-center justify-center text-slate-400">
                      <Image className="h-10 w-10 opacity-30" />
                    </div>
                  )}
                  {/* Hover Search Icon Overlay */}
                  <div className="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity duration-300 flex items-center justify-center">
                    <div className="h-10 w-10 bg-white/20 border border-white/30 rounded-full flex items-center justify-center text-white backdrop-blur-sm shadow-md transition-transform duration-300 group-hover:scale-110">
                      <Search className="h-5 w-5" />
                    </div>
                  </div>
                </div>
                {/* Media Details */}
                <div className="p-4 flex-1 flex flex-col justify-between space-y-2 text-left">
                  <div className="space-y-1">
                    <h4 className="text-xs font-bold text-slate-800 line-clamp-1" title={item.title}>{item.title}</h4>
                    {item.description && (
                      <p className="text-[10px] text-slate-500 font-medium leading-relaxed line-clamp-2" title={item.description}>{item.description}</p>
                    )}
                  </div>
                  <div className="text-[9px] text-slate-400 font-bold border-t border-slate-50 pt-2 flex items-center justify-between">
                    <span className="uppercase tracking-wide text-blue-600 font-bold">{item.media_type || 'Gambar'}</span>
                    <span>
                      {item.created_at ? new Date(item.created_at).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' }) : ''}
                    </span>
                  </div>
                </div>
              </div>
            ))}
          </div>
        ) : (
          <div className="py-16 text-center text-slate-400 text-xs font-medium">
            Belum ada dokumentasi media di galeri PPID.
          </div>
        )}
      </div>

      {/* Lightbox Modal */}
      <AnimatePresence>
        {activeLightboxImage && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 z-[9999] bg-black/90 backdrop-blur-md flex items-center justify-center p-4 cursor-zoom-out"
            onClick={() => setActiveLightboxImage(null)}
          >
            <button
              onClick={() => setActiveLightboxImage(null)}
              className="absolute top-6 right-6 text-white/70 hover:text-white hover:bg-white/10 p-2.5 rounded-full transition-all cursor-pointer z-50 border border-white/10"
              title="Tutup"
            >
              <X className="h-6 w-6" />
            </button>

            <motion.div
              initial={{ scale: 0.9, y: 20 }}
              animate={{ scale: 1, y: 0 }}
              exit={{ scale: 0.9, y: 20 }}
              className="relative max-w-4xl w-full flex flex-col items-center gap-4 cursor-default"
              onClick={(e) => e.stopPropagation()}
            >
              <div className="relative rounded-2xl overflow-hidden bg-slate-900 shadow-2xl border border-white/5 max-h-[75vh]">
                <img
                  src={resolveImageUrl(activeLightboxImage.media_url)}
                  alt={activeLightboxImage.title}
                  className="max-w-full max-h-[70vh] object-contain"
                />
              </div>
              <div className="text-center text-white max-w-2xl px-4 space-y-1">
                <h3 className="text-sm font-extrabold tracking-wide uppercase text-amber-400">{activeLightboxImage.title}</h3>
                {activeLightboxImage.description && (
                  <p className="text-xs text-slate-350 font-medium leading-relaxed">{activeLightboxImage.description}</p>
                )}
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
