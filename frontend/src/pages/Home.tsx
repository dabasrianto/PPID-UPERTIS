import { useState, useEffect, useRef, useMemo } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import {
  Landmark, FileText, Clock, Users, CheckCircle, Check, MapPin, Mail, Phone, Search, ArrowRight
} from 'lucide-react';
import { resolveImageUrl } from '../utils/helpers';
import type { Post, DownloadItem } from '../types';
import NewsCard from '../components/NewsCard';
import DownloadTable from '../components/DownloadTable';

interface HomeProps {
  siteConfig: any;
  heroImages: string[];
  heroIdx: number;
  prevSlide: () => void;
  nextSlide: () => void;
  posts: Post[];
  isPostsLoading: boolean;
  dbDownloads: DownloadItem[];
  navigateToNews: () => void;
  navigateToPage: (slug: string) => void;
  navigateToNewsDetail: (slug: string) => void;
  handleNavigation: (href: string) => void;
  incrementDownloadCount: (id: string, fileUrl: string) => void;
}

export default function Home({
  siteConfig,
  heroImages,
  heroIdx,
  prevSlide,
  nextSlide,
  posts,
  isPostsLoading,
  dbDownloads,
  navigateToNews,
  navigateToPage,
  navigateToNewsDetail,
  handleNavigation,
  incrementDownloadCount
}: HomeProps) {
  const [aboutTab, setAboutTab] = useState<string>('sejarah');
  const [aboutWidth, setAboutWidth] = useState<number>(540);
  const aboutRef = useRef<HTMLDivElement>(null);
  const [activeTab, setActiveTab] = useState<string>('berkala');
  const [searchTerm, setSearchTerm] = useState<string>('');

  useEffect(() => {
    if (!aboutRef.current) return;

    const handleResize = () => {
      if (aboutRef.current) {
        setAboutWidth(aboutRef.current.offsetWidth);
      }
    };

    handleResize();
    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, []);

  // Compute DIP Data tree
  const dipData = useMemo(() => {
    const sections = [
      {
        id: 'berkala',
        title: 'Informasi Berkala',
        desc: 'Informasi publik yang wajib diperbarui dan disediakan secara berkala sesuai UU KIP Pasal 9.',
        dbKey: 'ppid-berkala',
        items: [
          { label: 'Profil UPERTIS', href: 'profil', type: 'Page' },
          { label: 'Visi Misi UPERTIS', href: 'visi-misi', type: 'Page' },
          { label: 'Tugas Fungsi dan Ruang Lingkup Kegiatan', href: 'tugas-dan-fungsi', type: 'Page' },
          { label: 'Struktur Organisasi UPERTIS', href: 'struktur-organisasi-2', type: 'Page' }
        ]
      },
      {
        id: 'setiap-saat',
        title: 'Tersedia Setiap Saat',
        desc: 'Informasi publik yang wajib disediakan dan siap diberikan setiap saat dibutuhkan sesuai UU KIP Pasal 11.',
        dbKey: 'ppid-setiap-saat',
        items: [
          { label: 'Daftar Informasi Publik (DIP) UPERTIS', href: 'https://ppid.upertis.ac.id/download/', type: 'External' }
        ]
      },
      {
        id: 'serta-merta',
        title: 'Informasi Serta Merta',
        desc: 'Informasi publik yang wajib diumumkan segera demi keselamatan dan kesehatan masyarakat (UU KIP Pasal 10).',
        dbKey: 'ppid-serta-merta',
        items: [
          { label: 'Panduan Mitigasi & Evakuasi Bencana Alam di Lingkungan Kampus', href: '#', type: 'Info' }
        ]
      },
      {
        id: 'dikecualikan',
        title: 'Informasi Dikecualikan',
        desc: 'Informasi publik yang bersifat rahasia dan dibatasi ketat sesuai UU KIP Pasal 17.',
        dbKey: 'ppid-dikecualikan',
        items: [
          { label: 'SK Daftar Informasi Dikecualikan UPERTIS', href: 'https://ppid.upertis.ac.id/download/daftar-informasi-dikecualikan-informasi-setiap-saat/', type: 'External' }
        ]
      }
    ];

    return sections.map(section => {
      const dbItems = dbDownloads
        .filter(item => item.category === section.dbKey || item.category === section.id)
        .map(item => ({
          label: item.title,
          href: resolveImageUrl(item.file_url),
          type: 'DownloadItem',
          downloadsCount: item.downloads_count,
          id: item.id
        }));
      return {
        ...section,
        items: [...section.items, ...dbItems]
      };
    });
  }, [dbDownloads]);

  const currentDipSection = useMemo(() => {
    return dipData.find(sec => sec.id === activeTab);
  }, [dipData, activeTab]);

  const filteredDipItems = useMemo(() => {
    if (!currentDipSection) return [];
    return currentDipSection.items.filter(item =>
      item.label.toLowerCase().includes(searchTerm.toLowerCase())
    );
  }, [currentDipSection, searchTerm]);

  return (
    <div className="flex-1 flex flex-col">
      {/* Premium Welcome Hero section */}
      <section className="relative overflow-hidden bg-[#002147] text-white py-24 lg:py-36 px-4 text-left">
        {/* Background Slider */}
        {heroImages.length > 0 && (
          <div className="absolute inset-0 overflow-hidden z-0">
            <AnimatePresence mode="popLayout">
              {(() => {
                const animType = siteConfig?.settings?.hero_animation || 'shutter-3d';
                const opacityVal = siteConfig?.settings?.hero_image_opacity !== undefined
                  ? parseFloat(siteConfig.settings.hero_image_opacity)
                  : 0.25;

                const slideImgUrl = resolveImageUrl(
                  typeof heroImages[heroIdx] === 'string'
                    ? (heroImages[heroIdx] as string)
                    : (heroImages[heroIdx] as any)?.image || ''
                );

                // --- 1. SHUTTER 3D (Grid Explode) ---
                if (animType === 'shutter-3d') {
                  return (
                    <div className="absolute inset-0 overflow-hidden" style={{ perspective: 1000 }}>
                      {Array.from({ length: 24 }).map((_, index) => {
                        const r = Math.floor(index / 6);
                        const c = index % 6;
                        const left = (c * 100) / 6;
                        const right = ((c + 1) * 100) / 6;
                        const top = (r * 100) / 4;
                        const bottom = ((r + 1) * 100) / 4;
                        const clipPath = `polygon(${left}% ${top}%, ${right}% ${top}%, ${right}% ${bottom}%, ${left}% ${bottom}%)`;
                        const delay = (c + r) * 0.04;

                        return (
                          <motion.div
                            key={`${heroIdx}-${r}-${c}`}
                            style={{
                              position: "absolute",
                              inset: 0,
                              clipPath: clipPath,
                              transformOrigin: `${left + 100 / 12}% ${top + 100 / 8}%`,
                            }}
                            initial={{
                              opacity: 0,
                              rotateX: 45,
                              rotateY: 45,
                              z: -150,
                              scale: 0.85,
                            }}
                            animate={{
                              opacity: opacityVal,
                              rotateX: 0,
                              rotateY: 0,
                              z: 0,
                              scale: 1,
                            }}
                            exit={{
                              opacity: 0,
                              rotateX: -45,
                              rotateY: -45,
                              z: -150,
                              scale: 0.85,
                            }}
                            transition={{
                              duration: 0.6,
                              delay: delay,
                              ease: [0.25, 1, 0.5, 1],
                            }}
                            className="absolute inset-0 w-full h-full pointer-events-none"
                          >
                            <img
                              src={slideImgUrl}
                              alt=""
                              className="w-full h-full object-cover"
                            />
                          </motion.div>
                        );
                      })}
                    </div>
                  );
                }

                // --- 2. FADE BLUR (Cinematic Soft Fade & Blur) ---
                if (animType === 'fade-blur') {
                  return (
                    <motion.div
                      key={`fade-blur-${heroIdx}`}
                      initial={{ opacity: 0, filter: 'blur(15px)', scale: 1.05 }}
                      animate={{ opacity: opacityVal, filter: 'blur(0px)', scale: 1 }}
                      exit={{ opacity: 0, filter: 'blur(15px)', scale: 0.95 }}
                      transition={{ duration: 0.8, ease: 'easeInOut' }}
                      className="absolute inset-0 w-full h-full pointer-events-none"
                    >
                      <img
                        src={slideImgUrl}
                        alt=""
                        className="w-full h-full object-cover"
                      />
                    </motion.div>
                  );
                }

                // --- 3. PARALLAX SLIDE (Horizontal Slide Over) ---
                return (
                  <motion.div
                    key={`parallax-slide-${heroIdx}`}
                    initial={animType === 'ken-burns' 
                      ? { opacity: 0, scale: 1.15 } 
                      : animType === 'split-diagonal'
                      ? { opacity: 0, clipPath: 'polygon(0 0, 0 0, 0 100%, 0% 100%)' }
                      : { x: '100%', opacity: 0 }
                    }
                    animate={animType === 'ken-burns'
                      ? { 
                          opacity: opacityVal, 
                          scale: 1, 
                          transition: { scale: { duration: 6, ease: 'linear' }, opacity: { duration: 0.8 } } 
                        }
                      : animType === 'split-diagonal'
                      ? {
                          opacity: opacityVal,
                          clipPath: 'polygon(0 0, 100% 0, 100% 100%, 0 100%)',
                          transition: { duration: 0.85, ease: [0.16, 1, 0.3, 1] }
                        }
                      : { 
                          x: '0%', 
                          opacity: opacityVal, 
                          transition: { type: 'spring', damping: 25, stiffness: 120 } 
                        }
                    }
                    exit={animType === 'ken-burns'
                      ? { opacity: 0, transition: { duration: 0.8 } }
                      : animType === 'split-diagonal'
                      ? { 
                          opacity: 0, 
                          clipPath: 'polygon(100% 0, 100% 0, 100% 100%, 100% 100%)', 
                          transition: { duration: 0.75, ease: [0.16, 1, 0.3, 1] } 
                        }
                      : { 
                          x: '-100%', 
                          opacity: 0, 
                          transition: { duration: 0.5, ease: 'easeInOut' } 
                        }
                    }
                    className="absolute inset-0 w-full h-full pointer-events-none"
                  >
                    <img
                      src={slideImgUrl}
                      alt=""
                      className="w-full h-full object-cover"
                    />
                  </motion.div>
                );
              })()}
            </AnimatePresence>
          </div>
        )}

        {/* Overlay vignette */}
        <div
          className="absolute inset-0 z-5 pointer-events-none"
          style={{
            backgroundColor: siteConfig?.settings?.hero_overlay_color || '#002147',
            opacity: siteConfig?.settings?.hero_image_opacity !== undefined
              ? (1 - parseFloat(siteConfig.settings.hero_image_opacity))
              : 0.8,
          }}
        />

        {/* Navigation Arrows */}
        {heroImages.length > 1 && (
          <>
            <button
              type="button"
              onClick={prevSlide}
              className="absolute left-4 top-1/2 -translate-y-1/2 z-20 flex h-10 w-10 items-center justify-center rounded-full bg-white/10 text-white hover:bg-white/20 transition-all hover:scale-105 cursor-pointer border border-white/15"
              aria-label="Sebelumnya"
            >
              —
            </button>
            <button
              type="button"
              onClick={nextSlide}
              className="absolute right-4 top-1/2 -translate-y-1/2 z-20 flex h-10 w-10 items-center justify-center rounded-full bg-white/10 text-white hover:bg-white/20 transition-all hover:scale-105 cursor-pointer border border-white/15"
              aria-label="Selanjutnya"
            >
              &gt;
            </button>
          </>
        )}

        <div className="container mx-auto max-w-6xl relative z-10 grid md:grid-cols-12 gap-8 items-center">
          {(() => {
            const currentSlide = heroImages[heroIdx] || {};
            const isSlideObj = typeof currentSlide === 'object' && currentSlide !== null;
            const slideSubtitle = (isSlideObj && (currentSlide as any).subtitle) || siteConfig?.settings?.hero_subtitle || 'Portal Transparansi Publik';
            const slideTitle = (isSlideObj && (currentSlide as any).title) || siteConfig?.settings?.hero_title || 'Portal Keterbukaan Informasi Publik';
            const slideDesc = (isSlideObj && (currentSlide as any).description) || siteConfig?.settings?.hero_description || 'Selamat datang di PPID UPERTIS. Akses informasi resmi dan dokumentasi transparansi pelayanan kampus kami secara bebas, cepat, dan mudah.';
            const slideBtn1Text = (isSlideObj && (currentSlide as any).btn1_text) || siteConfig?.settings?.hero_btn1_text || 'Permohonan Informasi';
            const slideBtn1Page = (isSlideObj && (currentSlide as any).btn1_page) || siteConfig?.settings?.hero_btn1_page || 'permohonan-informasi';
            const slideBtn2Text = (isSlideObj && (currentSlide as any).btn2_text) || siteConfig?.settings?.hero_btn2_text || 'Cari Dokumen Publik';
            const slideBtn2Page = (isSlideObj && (currentSlide as any).btn2_page) || siteConfig?.settings?.hero_btn2_page || 'regulasi';

            return (
              <motion.div
                key={heroIdx}
                initial={{ opacity: 0, y: 15 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -15 }}
                transition={{ duration: 0.5, ease: 'easeOut' }}
                className="md:col-span-8 space-y-6"
              >
                <span className="bg-amber-400 text-[#002147] text-[10px] font-extrabold px-3 py-1.5 rounded-full uppercase tracking-widest inline-flex items-center gap-1 shadow-sm">
                  <Landmark className="h-3.5 w-3.5" /> {slideSubtitle}
                </span>
                <h1 className="text-4xl lg:text-5xl font-extrabold leading-tight text-white">
                  {slideTitle}
                </h1>
                <p className="text-sm lg:text-base text-slate-200 font-medium leading-relaxed max-w-2xl">
                  {slideDesc}
                </p>
                <div className="flex flex-wrap gap-4 pt-2">
                  <button
                    onClick={() => navigateToPage(slideBtn1Page)}
                    className="px-6 py-3 bg-amber-400 hover:bg-amber-500 text-[#002147] font-bold text-xs uppercase tracking-wider rounded-xl transition-all shadow-md hover:shadow-lg hover:-translate-y-0.5 cursor-pointer inline-flex items-center gap-1.5"
                  >
                    <FileText className="h-4 w-4" /> {slideBtn1Text}
                  </button>
                  <button
                    onClick={() => navigateToPage(slideBtn2Page)}
                    className="px-6 py-3 bg-white/10 hover:bg-white/20 border border-white/20 text-white font-bold text-xs uppercase tracking-wider rounded-xl transition-all hover:-translate-y-0.5 cursor-pointer inline-flex items-center gap-1.5"
                  >
                    <Search className="h-4 w-4" /> {slideBtn2Text}
                  </button>
                </div>
              </motion.div>
            );
          })()}

          <div className="md:col-span-4 bg-[#001733]/60 border border-white/10 rounded-3xl p-6 backdrop-blur-md space-y-4">
            <h3 className="text-sm font-extrabold text-amber-400 uppercase tracking-wider flex items-center gap-1.5">
              <Clock className="h-4 w-4" /> Jam Pelayanan PPID
            </h3>
            <div className="space-y-2 text-xs font-semibold">
              {(siteConfig?.settings?.service_hours && Array.isArray(siteConfig.settings.service_hours) && siteConfig.settings.service_hours.length > 0
                ? siteConfig.settings.service_hours
                : [
                  { day: 'Senin – Kamis', time: '08:35 – 16:00 WIB', closed: false },
                  { day: 'Jumat', time: '08:30 – 12:30 WIB', closed: false },
                  { day: 'Sabtu, Minggu & Hari Besar', time: 'Tutup / Libur', closed: true },
                ]
              ).map((sh: any, idx: number) => (
                <div key={idx} className={`flex justify-between ${idx < 2 ? 'border-b border-white/5 pb-2' : 'pb-1'}`}>
                  <span className="text-slate-300">{sh.day}</span>
                  {sh.closed ? (
                    <span className="text-red-400 font-bold uppercase tracking-wider">{sh.time}</span>
                  ) : (
                    <span className="font-bold">{sh.time}</span>
                  )}
                </div>
              ))}
            </div>
            <p className="text-[10px] text-slate-400 leading-normal border-t border-white/5 pt-2 italic">
              {siteConfig?.settings?.service_location || 'Gedung A - Samping B'}
            </p>
          </div>
        </div>
      </section>

      {/* About PPID — Premium Industrie-Style Section */}
      <section className="py-20 lg:py-24 px-4 bg-white text-left overflow-hidden">
        <div className="container mx-auto max-w-6xl">
          <div className="grid lg:grid-cols-2 gap-12 lg:gap-16 items-center">
            {/* Left Column: Image wrapper with SVG mask and drop-shadow */}
            <div className="relative filter drop-shadow-[0_20px_25px_rgba(0,0,0,0.15)]">
              {/* SVG clipPath definition */}
              <svg width="0" height="0" className="absolute pointer-events-none" style={{ position: 'absolute', width: 0, height: 0 }}>
                <defs>
                  <clipPath id="about-image-clip" clipPathUnits="userSpaceOnUse">
                    <path d={`M 134,0 L ${aboutWidth - 40},0 A 40,40 0 0,1 ${aboutWidth},40 L ${aboutWidth},216 A 24,24 0 0,1 ${aboutWidth - 24},240 L ${aboutWidth - 220},240 A 40,40 0 0,0 ${aboutWidth - 260},280 L ${aboutWidth - 260},426 A 24,24 0 0,1 ${aboutWidth - 284},450 L 40,450 A 40,40 0 0,1 0,410 L 0,134 A 24,24 0 0,1 24,110 L 70,110 A 40,40 0 0,0 110,70 L 110,24 A 24,24 0 0,1 134,0 Z`} />
                  </clipPath>
                </defs>
              </svg>

              {/* Decorative dot grid at top-left */}
              <div className="absolute -left-8 -top-8 z-0 grid grid-cols-5 gap-2">
                {[...Array(25)].map((_, i) => (
                  <div key={i} className="w-2.5 h-2.5 rounded-full bg-amber-500/20" />
                ))}
              </div>

              {/* Image Wrapper (Clipped) */}
              <div
                ref={aboutRef}
                className="relative z-10 w-full h-[450px] group overflow-hidden"
                style={{
                  clipPath: 'url(#about-image-clip)',
                  WebkitClipPath: 'url(#about-image-clip)',
                }}
              >
                <img
                  src={
                    siteConfig?.settings?.about_image
                      ? resolveImageUrl(siteConfig.settings.about_image)
                      : siteConfig?.settings?.background_image
                        ? resolveImageUrl(siteConfig.settings.background_image)
                        : heroImages.length > 0
                          ? (typeof heroImages[0] === 'string' ? heroImages[0] : (heroImages[0] as any)?.image || '')
                          : 'https://images.unsplash.com/photo-1541339907198-e08756dedf3f?w=800'
                  }
                  alt="Tentang PPID UPERTIS"
                  className="w-full h-full object-cover transition-transform duration-700 ease-in-out group-hover:scale-105"
                />
                {/* Subtle gradient overlay on image */}
                <div className="absolute inset-0 bg-gradient-to-t from-[#002147]/20 to-transparent" />
              </div>

              {/* Stats overlay card */}
              <motion.div
                whileHover={{ y: -6, scale: 1.02 }}
                transition={{ type: 'spring', stiffness: 400, damping: 25 }}
                className="absolute bottom-5 right-5 w-[220px] h-[170px] bg-gradient-to-br from-[#002147] to-[#0b335c] rounded-[32px] p-6 text-white flex flex-col justify-center space-y-3.5 shadow-xl shadow-slate-900/20 z-20 border border-white/10"
              >
                <div className="flex items-center gap-3">
                  <div className="w-11 h-11 rounded-full bg-amber-400/20 border border-amber-400/30 flex items-center justify-center shrink-0 shadow-inner">
                    <Users className="h-5 w-5 text-amber-400" />
                  </div>
                  <span className="text-3xl font-black tracking-tight font-display text-amber-400">
                    {siteConfig?.settings?.about_stat_number || '2021'}
                  </span>
                </div>
                <div className="space-y-0.5 text-left">
                  <p className="text-xs font-bold text-amber-400 uppercase tracking-wider">
                    {siteConfig?.settings?.about_stat_label_accent || 'Tahun Berdiri'}
                  </p>
                  <p className="text-[11px] text-slate-200 font-medium leading-snug">
                    {siteConfig?.settings?.about_stat_label || 'PPID UPERTIS Melayani'}
                  </p>
                </div>
              </motion.div>
            </div>

            {/* Right Column: Content with tabs */}
            <div className="space-y-6">
              <span className="text-xs font-extrabold text-amber-500 uppercase tracking-[0.25em]">
                {siteConfig?.settings?.about_badge || 'Tentang PPID'}
              </span>

              <h2 className="text-3xl lg:text-4xl font-extrabold text-[#002147] leading-tight">
                {(() => {
                  const title = siteConfig?.settings?.about_heading || 'Melayani Keterbukaan Informasi sejak 2021.';
                  const parts = title.split(/(\*\*.*?\*\*)/g);
                  return parts.map((part: string, idx: number) => {
                    if (part.startsWith('**') && part.endsWith('**')) {
                      return <span key={idx} className="text-amber-500">{part.slice(2, -2)}</span>;
                    }
                    return <span key={idx}>{part}</span>;
                  });
                })()}
              </h2>

              <p className="text-sm text-slate-500 font-medium leading-relaxed">
                {siteConfig?.settings?.about_subtitle || 'Pejabat Pengelola Informasi dan Dokumentasi (PPID) Universitas Perintis Indonesia hadir sebagai wujud komitmen kampus dalam menerapkan keterbukaan informasi publik.'}
              </p>

              {/* Tab navigation */}
              <div className="flex gap-1 border-b border-slate-200">
                {[
                  { key: 'sejarah', label: 'Sejarah' },
                  { key: 'visi', label: 'Visi Kami' },
                  { key: 'misi', label: 'Misi Kami' },
                ].map((tab) => (
                  <button
                    key={tab.key}
                    onClick={() => setAboutTab(tab.key)}
                    className={`px-5 py-2.5 text-xs font-bold transition-all cursor-pointer relative ${aboutTab === tab.key
                      ? 'text-amber-600'
                      : 'text-slate-400 hover:text-slate-600'
                      }`}
                  >
                    {tab.label}
                    {aboutTab === tab.key && (
                      <motion.div
                        layoutId="aboutTabIndicator"
                        className="absolute bottom-0 left-0 right-0 h-[2px] bg-amber-500 rounded-full"
                        transition={{ type: 'spring', stiffness: 500, damping: 35 }}
                      />
                    )}
                  </button>
                ))}
              </div>

              {/* Tab content */}
              <AnimatePresence mode="wait">
                <motion.div
                  key={aboutTab}
                  initial={{ opacity: 0, y: 8 }}
                  animate={{ opacity: 1, y: 0 }}
                  exit={{ opacity: 0, y: -8 }}
                  transition={{ duration: 0.25 }}
                  className="text-sm text-slate-500 font-medium leading-relaxed min-h-[80px]"
                >
                  {aboutTab === 'sejarah' && (
                    <p>{siteConfig?.settings?.about_sejarah || 'PPID Universitas Perintis Indonesia dibentuk berdasarkan Surat Keputusan Rektor sebagai bagian dari pelaksanaan Undang-Undang No. 14 Tahun 2008 tentang Keterbukaan Informasi Publik. Sejak awal berdiri, PPID UPERTIS berkomitmen menyediakan informasi publik yang transparan, akuntabel, dan mudah diakses oleh seluruh pemangku kepentingan.'}</p>
                  )}
                  {aboutTab === 'visi' && (
                    <p>{siteConfig?.settings?.ppid_visi || 'Mewujudkan tata kelola informasi publik yang transparan, akuntabel, dan profesional di lingkungan Universitas Perintis Indonesia untuk mendukung good university governance.'}</p>
                  )}
                  {aboutTab === 'misi' && (
                    <ul className="space-y-2">
                      {(siteConfig?.settings?.ppid_misi
                        ? siteConfig.settings.ppid_misi.split('\n').filter((s: string) => s.trim())
                        : [
                          'Menyediakan akses informasi publik yang cepat, tepat, dan mudah dijangkau.',
                          'Mengelola dokumentasi informasi publik secara sistematis and terintegrasi.',
                          'Meningkatkan partisipasi masyarakat dalam pengawasan layanan pendidikan tinggi.',
                          'Mendorong budaya keterbukaan dan akuntabilitas di seluruh unit kerja kampus.'
                        ]
                      ).map((item: string, idx: number) => (
                        <li key={idx} className="flex items-start gap-2">
                          <CheckCircle className="h-4 w-4 text-amber-500 mt-0.5 shrink-0" />
                          <span>{item.replace(/^[-•\d.]\s*/, '')}</span>
                        </li>
                      ))}
                    </ul>
                  )}
                </motion.div>
              </AnimatePresence>

              {/* Feature bullets */}
              <div className="flex flex-wrap gap-x-8 gap-y-3 pt-2">
                {(siteConfig?.settings?.about_features
                  ? (Array.isArray(siteConfig.settings.about_features) ? siteConfig.settings.about_features : [])
                  : [
                    'Layanan Informasi Terbuka',
                    'Pengelolaan Data Transparan'
                  ]
                ).map((feature: string, idx: number) => (
                  <div key={idx} className="flex items-center gap-2">
                    <div className="w-5 h-5 rounded-full border-2 border-amber-500 flex items-center justify-center">
                      <Check className="h-3 w-3 text-amber-500" />
                    </div>
                    <span className="text-xs font-bold text-[#002147]">{feature}</span>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Selamat Datang – Layanan Utama Cards */}
      <section className="py-20 px-4 bg-white text-center">
        <div className="container mx-auto max-w-6xl space-y-12">

          {/* Header */}
          <div className="space-y-4 max-w-3xl mx-auto">
            <div className="flex items-center justify-center gap-2">
              <svg className="h-5 w-5 text-amber-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}><path strokeLinecap="round" strokeLinejoin="round" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0h6" /></svg>
              <span className="text-xs font-extrabold text-amber-500 uppercase tracking-[0.22em]">Layanan Informasi</span>
            </div>
            <h2 className="text-3xl lg:text-4xl font-extrabold text-[#002147] leading-tight">
              Selamat <span className="text-amber-500">Datang</span>
            </h2>
            <p className="text-sm text-slate-500 font-medium leading-relaxed">
              {siteConfig?.settings?.welcome_text ||
                'Layanan Informasi Publik UPERTIS disediakan untuk memudahkan publik mendapatkan informasi tentang UPERTIS. Publik berhak mengajukan informasi publik yang dikelola oleh UPERTIS sesuai ketentuan-ketentuan yang berlaku. UPERTIS melayani seluruh permohonan informasi melalui Layanan Informasi Publik secara daring maupun luring.'}
            </p>
          </div>

          {/* 3 Service Cards */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6 lg:gap-8">

            {/* Card 1 – Permohonan Informasi Publik */}
            <div className="group relative bg-white border border-slate-100 rounded-3xl p-8 text-left shadow-sm hover:shadow-xl hover:-translate-y-1.5 transition-all duration-300 overflow-hidden">
              {/* Decorative circle */}
              <div className="absolute -top-8 -right-8 w-28 h-28 rounded-full bg-amber-50 border border-amber-100/60 transition-all duration-300 group-hover:scale-125 group-hover:bg-amber-100/70" />
              {/* Icon */}
              <div className="relative z-10 w-14 h-14 rounded-2xl bg-amber-50 border border-amber-200/60 flex items-center justify-center mb-6 group-hover:bg-amber-500 transition-colors duration-300">
                <svg className="h-7 w-7 text-amber-500 group-hover:text-white transition-colors duration-300" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.6}>
                  <path strokeLinecap="round" strokeLinejoin="round" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                </svg>
              </div>
              <div className="relative z-10 space-y-3">
                <h3 className="text-base font-extrabold text-[#002147] leading-tight">{siteConfig?.settings?.service_card_1_title || 'Permohonan Informasi Publik'}</h3>
                <p className="text-xs text-slate-500 font-medium leading-relaxed">
                  {siteConfig?.settings?.service_card_1_desc || 'Ajukan permohonan informasi publik kepada PPID UPERTIS secara daring maupun luring dengan mudah dan cepat sesuai prosedur UU KIP.'}
                </p>
                <button
                  onClick={() => {
                    const link = siteConfig?.settings?.service_card_1_link || 'permohonan-informasi';
                    if (link.startsWith('http')) window.open(link, '_blank');
                    else navigateToPage(link);
                  }}
                  className="inline-flex items-center gap-1.5 text-[11px] font-extrabold text-[#002147] uppercase tracking-wider mt-2 group-hover:text-amber-600 transition-colors cursor-pointer border-0 bg-transparent p-0"
                >
                  Selengkapnya <ArrowRight className="h-3.5 w-3.5" />
                </button>
              </div>
            </div>

            {/* Card 2 – Keberatan Informasi Publik */}
            <div className="group relative bg-white border border-slate-100 rounded-3xl p-8 text-left shadow-sm hover:shadow-xl hover:-translate-y-1.5 transition-all duration-300 overflow-hidden">
              <div className="absolute -top-8 -right-8 w-28 h-28 rounded-full bg-amber-50 border border-amber-100/60 transition-all duration-300 group-hover:scale-125 group-hover:bg-amber-100/70" />
              <div className="relative z-10 w-14 h-14 rounded-2xl bg-amber-50 border border-amber-200/60 flex items-center justify-center mb-6 group-hover:bg-amber-500 transition-colors duration-300">
                <svg className="h-7 w-7 text-amber-500 group-hover:text-white transition-colors duration-300" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.6}>
                  <path strokeLinecap="round" strokeLinejoin="round" d="M3 6l3 1m0 0l-3 9a5.002 5.002 0 006.001 0M6 7l3 9M6 7l6-2m6 2l3-1m-3 1l-3 9a5.002 5.002 0 006.001 0M18 7l3 9m-3-9l-6-2m0-2v2m0 16V5m0 16H9m3 0h3" />
                </svg>
              </div>
              <div className="relative z-10 space-y-3">
                <h3 className="text-base font-extrabold text-[#002147] leading-tight">{siteConfig?.settings?.service_card_2_title || 'Keberatan Informasi Publik'}</h3>
                <p className="text-xs text-slate-500 font-medium leading-relaxed">
                  {siteConfig?.settings?.service_card_2_desc || 'Ajukan keberatan resmi jika permohonan informasi Anda tidak ditanggapi, ditolak, atau tidak sesuai dengan ketentuan yang berlaku.'}
                </p>
                <button
                  onClick={() => {
                    const link = siteConfig?.settings?.service_card_2_link || 'keberatan-informasi';
                    if (link.startsWith('http')) window.open(link, '_blank');
                    else navigateToPage(link);
                  }}
                  className="inline-flex items-center gap-1.5 text-[11px] font-extrabold text-[#002147] uppercase tracking-wider mt-2 group-hover:text-amber-600 transition-colors cursor-pointer border-0 bg-transparent p-0"
                >
                  Selengkapnya <ArrowRight className="h-3.5 w-3.5" />
                </button>
              </div>
            </div>

            {/* Card 3 – Pengaduan Layanan */}
            <div className="group relative bg-white border border-slate-100 rounded-3xl p-8 text-left shadow-sm hover:shadow-xl hover:-translate-y-1.5 transition-all duration-300 overflow-hidden">
              <div className="absolute -top-8 -right-8 w-28 h-28 rounded-full bg-amber-50 border border-amber-100/60 transition-all duration-300 group-hover:scale-125 group-hover:bg-amber-100/70" />
              <div className="relative z-10 w-14 h-14 rounded-2xl bg-amber-50 border border-amber-200/60 flex items-center justify-center mb-6 group-hover:bg-amber-500 transition-colors duration-300">
                <svg className="h-7 w-7 text-amber-500 group-hover:text-white transition-colors duration-300" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.6}>
                  <path strokeLinecap="round" strokeLinejoin="round" d="M18.364 5.636l-3.536 3.536m0 5.656l3.536 3.536M9.172 9.172L5.636 5.636m3.536 9.192l-3.536 3.536M21 12a9 9 0 11-18 0 9 9 0 0118 0zm-5 0a4 4 0 11-8 0 4 4 0 018 0z" />
                </svg>
              </div>
              <div className="relative z-10 space-y-3">
                <h3 className="text-base font-extrabold text-[#002147] leading-tight">{siteConfig?.settings?.service_card_3_title || 'Pengaduan Layanan'}</h3>
                <p className="text-xs text-slate-500 font-medium leading-relaxed">
                  {siteConfig?.settings?.service_card_3_desc || 'Sampaikan pengaduan terkait layanan informasi publik PPID UPERTIS melalui kanal resmi kami untuk penanganan yang transparan dan akuntabel.'}
                </p>
                <button
                  onClick={() => {
                    const link = siteConfig?.settings?.service_card_3_link || 'kontak';
                    if (link.startsWith('http')) window.open(link, '_blank');
                    else navigateToPage(link);
                  }}
                  className="inline-flex items-center gap-1.5 text-[11px] font-extrabold text-[#002147] uppercase tracking-wider mt-2 group-hover:text-amber-600 transition-colors cursor-pointer border-0 bg-transparent p-0"
                >
                  Selengkapnya <ArrowRight className="h-3.5 w-3.5" />
                </button>
              </div>
            </div>

          </div>
        </div>
      </section>


      {/* News Highlights section */}
      <section className="py-16 px-4 text-left">
        <div className="container mx-auto max-w-6xl space-y-8">
          <div className="flex items-end justify-between border-b border-slate-200 pb-4">
            <div className="space-y-1">
              <span className="text-xs font-bold text-amber-600 uppercase tracking-widest">Informasi Warta</span>
              <h2 className="text-2xl font-extrabold text-[#002147]">Berita Keterbukaan Informasi</h2>
            </div>
            <button
              onClick={navigateToNews}
              className="text-xs font-bold text-amber-600 hover:text-[#002147] transition-colors uppercase tracking-wider flex items-center gap-1 cursor-pointer"
            >
              Lihat Semua <ArrowRight className="h-4 w-4" />
            </button>
          </div>

          {isPostsLoading ? (
            <div className="grid md:grid-cols-3 gap-6">
              {[1, 2, 3].map((i) => (
                <div key={i} className="bg-slate-100 rounded-3xl h-64 animate-pulse" />
              ))}
            </div>
          ) : posts.length > 0 ? (
            <div className="grid md:grid-cols-3 gap-6">
              {posts.slice(0, 3).map((post) => (
                <NewsCard
                  key={post.id}
                  post={post}
                  navigateToNewsDetail={navigateToNewsDetail}
                />
              ))}
            </div>
          ) : (
            <div className="py-12 text-center text-slate-400 text-xs font-medium">
              Belum ada berita dipublikasikan.
            </div>
          )}
        </div>
      </section>

      {/* Physical Helpdesk Address section */}
      <section className="py-16 px-4 bg-slate-50 border-t border-slate-200/50 text-left">
        <div className="container mx-auto max-w-6xl grid md:grid-cols-2 gap-8 items-center">
          <div className="space-y-6">
            <span className="text-xs font-bold text-amber-600 uppercase tracking-widest">Desk Layanan PPID</span>
            <h2 className="text-2xl font-extrabold text-[#002147]">Layanan Informasi Fisik</h2>
            <p className="text-xs text-slate-400 font-medium leading-relaxed">
              Kami juga melayani permintaan informasi secara tatap muka (fisik). Anda dapat mendatangi desk pelayanan informasi kami di kantor rektorat universitas.
            </p>
            <div className="space-y-4">
              <div className="flex gap-3">
                <MapPin className="h-5 w-5 text-amber-500 shrink-0" />
                <div>
                  <span className="text-xs font-extrabold text-slate-800 block leading-tight">Alamat Gedung Rektorat</span>
                  <span className="text-[11px] text-slate-400 font-medium mt-0.5 block leading-normal">
                    {siteConfig?.settings?.rektorat_address || 'Kampus Utama Padang, Lantai 1 Gedung Rektorat Universitas Perintis Indonesia, Jl. Adinegoro Km. 17, Lubuk Buaya, Padang, Sumatera Barat.'}
                  </span>
                </div>
              </div>
              <div className="flex gap-3">
                <Mail className="h-5 w-5 text-amber-500 shrink-0" />
                <div>
                  <span className="text-xs font-extrabold text-slate-800 block leading-tight">Email Korespondensi</span>
                  <span className="text-[11px] text-slate-400 font-medium mt-0.5 block font-mono">
                    {siteConfig?.settings?.rektorat_email || 'ppidcare@upertis.ac.id'}
                  </span>
                </div>
              </div>
              <div className="flex gap-3">
                <Phone className="h-5 w-5 text-amber-500 shrink-0" />
                <div>
                  <span className="text-xs font-extrabold text-slate-800 block leading-tight">Hubungi WhatsApp Care</span>
                  <span className="text-[11px] text-slate-400 font-medium mt-0.5 block font-mono">
                    {siteConfig?.settings?.rektorat_phone || '+62 852 6355 7272'}
                  </span>
                </div>
              </div>
            </div>
          </div>

          <div className="bg-white border border-slate-200 p-4 rounded-[2rem] shadow-md h-[320px] overflow-hidden">
            <iframe
              title="Loket UPERTIS Map"
              src={siteConfig?.settings?.kampus1_map_url || 'https://maps.google.com/maps?q=Universitas+Perintis+Indonesia,+Jl.+Adinegoro,+Lubuk+Buaya,+Padang,+Sumatera+Barat&output=embed&hl=id&z=15'}
              width="100%"
              height="100%"
              style={{ border: 0, borderRadius: '1.5rem' }}
              allowFullScreen={true}
              loading="lazy"
              referrerPolicy="no-referrer-when-downgrade"
            />
          </div>
        </div>
      </section>
    </div>
  );
}
