import { Image } from 'lucide-react';
import { resolveImageUrl } from '../utils/helpers';
import type { Post } from '../types';

interface NewsCardProps {
  post: Post;
  navigateToNewsDetail: (slug: string) => void;
  viewMode?: 'grid' | 'list';
}

export default function NewsCard({ post, navigateToNewsDetail, viewMode = 'grid' }: NewsCardProps) {
  if (viewMode === 'list') {
    return (
      <div className="bg-white border border-slate-200 rounded-2xl sm:rounded-3xl overflow-hidden shadow-sm hover:shadow-md transition-all flex flex-row h-28 sm:h-44 hover:-translate-y-0.5 text-left group">
        {/* Left Side: Image */}
        <div
          onClick={() => navigateToNewsDetail(post.slug)}
          className="w-24 sm:w-52 bg-slate-200 relative shrink-0 overflow-hidden cursor-pointer group/img"
        >
          {(post.featured_image || post.cover_image_url) ? (
            <img
              src={resolveImageUrl(post.featured_image || post.cover_image_url)}
              alt={post.title}
              className="w-full h-full object-cover object-center transition-transform duration-350 group-hover/img:scale-105"
            />
          ) : (
            <div className="w-full h-full flex items-center justify-center text-slate-400">
              <Image className="h-6 sm:h-10 w-6 sm:w-10 opacity-40" />
            </div>
          )}
          <div className="absolute top-2 left-2 bg-[#002147] text-white text-[7px] sm:text-[9px] font-bold uppercase tracking-wider px-2 py-0.5 sm:px-2.5 sm:py-1 rounded-full z-10">
            {post.category || 'Berita'}
          </div>
        </div>

        {/* Right Side: Content */}
        <div className="p-3 sm:p-5 flex-1 flex flex-col justify-between min-w-0">
          <div className="space-y-1 sm:space-y-2">
            <h3
              onClick={() => navigateToNewsDetail(post.slug)}
              className="text-xs sm:text-sm font-extrabold text-[#002147] hover:text-amber-600 transition-colors line-clamp-2 leading-snug cursor-pointer"
            >
              {post.title}
            </h3>
            <p className="hidden sm:line-clamp-3 text-[10px] sm:text-[11px] text-slate-400 font-medium leading-relaxed font-sans">
              {post.excerpt}
            </p>
          </div>
          <div className="flex items-center justify-between border-t border-slate-100 pt-2">
            <span className="text-[8px] sm:text-[10px] text-slate-400 font-semibold uppercase tracking-wider">
              {new Date(post.created_at).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })}
            </span>
            <button
              onClick={() => navigateToNewsDetail(post.slug)}
              className="text-[8px] sm:text-[10px] font-extrabold text-amber-600 hover:text-[#002147] transition-colors uppercase tracking-wider cursor-pointer bg-transparent border-0 p-0"
            >
              Selengkapnya
            </button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div
      className="bg-white border border-slate-200 rounded-2xl sm:rounded-3xl overflow-hidden shadow-sm hover:shadow-md transition-all flex flex-col h-full hover:-translate-y-1 text-left"
    >
      <div
        onClick={() => navigateToNewsDetail(post.slug)}
        className="h-28 sm:h-44 bg-slate-200 relative shrink-0 overflow-hidden cursor-pointer group/img"
      >
        {(post.featured_image || post.cover_image_url) ? (
          <img
            src={resolveImageUrl(post.featured_image || post.cover_image_url)}
            alt={post.title}
            className="w-full h-full object-cover object-center transition-transform duration-350 group-hover/img:scale-105"
          />
        ) : (
          <div className="w-full h-full flex items-center justify-center text-slate-400">
            <Image className="h-8 w-8 sm:h-10 sm:w-10 opacity-40" />
          </div>
        )}
        <div className="absolute top-2 left-2 sm:top-3 sm:left-3 bg-[#002147] text-white text-[8px] sm:text-[9px] font-bold uppercase tracking-wider px-2 py-0.5 sm:px-2.5 sm:py-1 rounded-full z-10">
          {post.category || 'Berita PPID'}
        </div>
      </div>
      <div className="p-3 sm:p-5 flex-1 flex flex-col justify-between gap-3">
        <div className="space-y-1 sm:space-y-2">
          <h3
            onClick={() => navigateToNewsDetail(post.slug)}
            className="text-xs sm:text-sm font-extrabold text-[#002147] hover:text-amber-600 transition-colors line-clamp-2 leading-snug cursor-pointer"
          >
            {post.title}
          </h3>
          <p className="text-[10px] sm:text-[11px] text-slate-400 line-clamp-2 sm:line-clamp-3 font-medium leading-relaxed font-sans">
            {post.excerpt}
          </p>
        </div>
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-2 border-t border-slate-100 pt-2 sm:pt-3">
          <span className="text-[8px] sm:text-[10px] text-slate-400 font-semibold uppercase tracking-wider">
            {new Date(post.created_at).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })}
          </span>
          <button
            onClick={() => navigateToNewsDetail(post.slug)}
            className="text-[8px] sm:text-[10px] font-extrabold text-amber-600 hover:text-[#002147] transition-colors uppercase tracking-wider cursor-pointer bg-transparent border-0 p-0 text-left"
          >
            Selengkapnya
          </button>
        </div>
      </div>
    </div>
  );
}

