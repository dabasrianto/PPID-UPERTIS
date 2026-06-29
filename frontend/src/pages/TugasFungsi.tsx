import { useState } from 'react';
import { 
  Briefcase, BookOpen, CheckCircle, Shield, FileText, Share2, 
  RefreshCw, BarChart2, ShieldAlert, FolderOpen, Compass, 
  HeartHandshake, Settings, ChevronRight, X
} from 'lucide-react';
import { resolveImageUrl, preprocessPostContent } from '../utils/helpers';
import type { PageData } from '../types';

interface TugasFungsiProps {
  pageData: PageData;
}

export default function TugasFungsi({ pageData }: TugasFungsiProps) {
  const [activeTab, setActiveTab] = useState<'tugas' | 'fungsi'>('tugas');
  const [openAccordion, setOpenAccordion] = useState<number | null>(0);
  const [lightboxImage, setLightboxImage] = useState<string | null>(null);

  // Helper to parse database JSON
  let parsedJson: any = null;
  try {
    parsedJson = JSON.parse(pageData.content || '');
  } catch (e) {
    parsedJson = null;
  }

  // Fallback defaults if database has legacy/empty content
  const defaultTugas = [
    { title: "Pengelolaan & Penyimpanan Arsip", desc: "Mengelola dan menyimpan dokumen serta arsip informasi publik agar terstruktur, aman, dan mudah diakses." },
    { title: "Penyediaan & Penerbitan Informasi", desc: "Menyediakan, memberikan, dan menerbitkan dokumen informasi publik secara proaktif kepada masyarakat luas." },
    { title: "Pengklasifikasian Informasi Publik", desc: "Melakukan klasifikasi berkala atas informasi (berkala, serta-merta, setiap saat, dan dikecualikan) sesuai regulasi." },
    { title: "Uji Konsekuensi Informasi Dikecualikan", desc: "Melakukan analisis dampak dan uji konsekuensi hukum yang ketat sebelum mengecualikan akses suatu informasi." },
    { title: "Penyusunan Daftar Informasi Publik (DIP)", desc: "Menyusun, memutakhirkan, dan mempublikasikan Daftar Informasi Publik (DIP) UPERTIS secara komprehensif." },
    { title: "Pengembangan Sistem Layanan Informasi", desc: "Membangun, mengoperasikan, dan mengawasi jalannya sistem IT maupun konvensional untuk layanan informasi publik." },
    { title: "Penyelesaian Sengketa Informasi Internal", desc: "Menyelesaikan setiap perselisihan permohonan informasi melalui jalur mediasi internal sebelum ke Komisi Informasi." },
    { title: "Pelaporan Berkala Kinerja Layanan", desc: "Membuat laporan rutin berkala pelaksanaan layanan informasi publik untuk diserahkan kepada Atasan PPID." }
  ];

  const defaultFungsi = [
    {
      title: "Fungsi Penyimpanan dan Pendokumentasian",
      items: [
        "Mengarsipkan setiap dokumen dan berkas informasi publik yang dihasilkan oleh seluruh unit kerja Universitas.",
        "Menyimpan basis data secara digital maupun fisik agar terawat, mudah dilacak, dan bebas risiko hilang."
      ]
    },
    {
      title: "Fungsi Penyediaan Informasi",
      items: [
        "Menyediakan akses informasi publik secara cepat, tepat waktu, transparan, dan menggunakan skema sederhana.",
        "Menyusun struktur Daftar Informasi Publik (DIP) yang selaras dengan ketentuan peraturan perundang-undangan."
      ]
    },
    {
      title: "Fungsi Pelayanan Informasi",
      items: [
        "Memberikan pemenuhan permohonan informasi kepada pemohon sesuai standar operating procedure (SOP) PPID.",
        "Menyediakan sarana fisik (meja layanan ppid) dan digital (situs web portal PPID, chatbot interaktif, desk online)."
      ]
    },
    {
      title: "Fungsi Koordinasi Internal & Eksternal",
      items: [
        "Berkoordinasi aktif dengan PPID Pelaksana di tingkat fakultas dan biro untuk menghimpun data secara berkala.",
        "Membangun komunikasi harmonis dengan Tim Pertimbangan Hukum, Rektorat (Atasan PPID), dan Petugas Informasi."
      ]
    },
    {
      title: "Fungsi Pengendalian dan Pengawasan",
      items: [
        "Mengontrol dan mengawasi seluruh rantai proses pelayanan dokumen informasi publik agar memenuhi standar mutu KIP.",
        "Melakukan evaluasi kepuasan layanan dan audit kepatuhan keterbukaan informasi di internal perguruan tinggi."
      ]
    }
  ];

  // Resolve dynamic vs fallback lists
  const displayTugas = parsedJson && Array.isArray(parsedJson.tugas) ? parsedJson.tugas : defaultTugas;
  const displayFungsi = parsedJson && Array.isArray(parsedJson.fungsi) ? parsedJson.fungsi : defaultFungsi;
  const displayTambahan = parsedJson ? (parsedJson.tambahan || '') : (pageData.content || '');

  // Helper to dynamically match icons
  const getTugasIcon = (idx: number) => {
    const icons = [
      <FolderOpen className="h-10 w-10 text-amber-500" />,
      <Share2 className="h-10 w-10 text-amber-500" />,
      <Settings className="h-10 w-10 text-amber-500" />,
      <ShieldAlert className="h-10 w-10 text-amber-500" />,
      <FileText className="h-10 w-10 text-amber-500" />,
      <RefreshCw className="h-10 w-10 text-amber-500" />,
      <HeartHandshake className="h-10 w-10 text-amber-500" />,
      <BarChart2 className="h-10 w-10 text-amber-500" />
    ];
    return icons[idx % icons.length];
  };

  const getFungsiIcon = (idx: number) => {
    const icons = [
      <FolderOpen className="h-5 w-5 text-blue-600" />,
      <Compass className="h-5 w-5 text-amber-600" />,
      <HeartHandshake className="h-5 w-5 text-teal-600" />,
      <Share2 className="h-5 w-5 text-indigo-600" />,
      <Shield className="h-5 w-5 text-emerald-600" />
    ];
    return icons[idx % icons.length];
  };

  return (
    <div className="space-y-10 animate-in fade-in duration-200 text-left w-full py-6">
      {/* Premium Header Banner */}
      <div className="bg-[#002147] text-white rounded-[2rem] p-10 lg:p-16 relative overflow-hidden shadow-xl border border-slate-800">
        <div className="absolute right-0 bottom-0 opacity-10 translate-x-12 translate-y-12">
          <Briefcase className="h-96 w-96 text-amber-400 animate-pulse" />
        </div>
        <div className="absolute top-0 left-0 w-full h-full bg-gradient-to-br from-transparent via-transparent to-amber-400/5 pointer-events-none" />
        <div className="relative z-10 space-y-4 max-w-3xl">
          <span className="bg-amber-400 text-[#002147] text-[10px] font-extrabold px-3.5 py-1.5 rounded-full uppercase tracking-widest font-mono">
            Tupoksi PPID
          </span>
          <h1 className="text-4xl lg:text-5xl font-black leading-tight tracking-tight">
            Tugas & Fungsi PPID
          </h1>
          <p className="text-xs lg:text-sm text-slate-350 leading-relaxed font-medium">
            Kedudukan hukum, wewenang operasional, dan komitmen pelayanan informasi publik Pejabat Pengelola Informasi & Dokumentasi UPERTIS.
          </p>
        </div>
      </div>

      {/* Luxury Tabs Switcher */}
      <div className="flex justify-center border-b border-slate-200 pb-px">
        <div className="inline-flex bg-slate-100 p-1.5 rounded-2xl border border-slate-200">
          <button
            onClick={() => setActiveTab('tugas')}
            className={`px-6 py-2.5 rounded-xl font-bold text-xs cursor-pointer transition-all border-0 ${
              activeTab === 'tugas' 
                ? 'bg-white text-[#002147] shadow-sm' 
                : 'bg-transparent text-slate-500 hover:text-[#002147]'
            }`}
          >
            Tugas Pokok PPID
          </button>
          <button
            onClick={() => setActiveTab('fungsi')}
            className={`px-6 py-2.5 rounded-xl font-bold text-xs cursor-pointer transition-all border-0 ${
              activeTab === 'fungsi' 
                ? 'bg-white text-[#002147] shadow-sm' 
                : 'bg-transparent text-slate-500 hover:text-[#002147]'
            }`}
          >
            Fungsi Strategis PPID
          </button>
        </div>
      </div>

      {/* Content Render Area */}
      <div className="max-w-7xl mx-auto px-4">
        {activeTab === 'tugas' ? (
          <div className="space-y-8 animate-in fade-in duration-300">
            <div className="text-center max-w-lg mx-auto mb-10 space-y-2">
              <h2 className="text-2xl font-black text-[#002147]">8 Tanggung Jawab Utama PPID</h2>
              <p className="text-xs text-slate-400 font-medium">Melaksanakan mandat UU No. 14 Tahun 2008 tentang Keterbukaan Informasi Publik di lingkungan perguruan tinggi.</p>
            </div>

            {/* Asymmetric Card Grid Layout (Matching Reference Screenshot) */}
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
              {displayTugas.map((tugas: any, idx: number) => (
                <div 
                  key={idx}
                  className="p-8 bg-white border border-slate-200/80 rounded-tl-[3.5rem] rounded-br-[3.5rem] rounded-tr-2xl rounded-bl-2xl hover:border-amber-400/60 hover:shadow-lg hover:shadow-slate-100 transition-all duration-300 flex flex-col justify-between min-h-[280px] relative overflow-hidden group text-left"
                >
                  {/* Decorative background number in top-right */}
                  <div className="absolute -top-1 -right-1 text-[5rem] font-black text-slate-100/70 group-hover:text-amber-500/10 pointer-events-none transition-colors duration-500 select-none font-mono leading-none tracking-tighter">
                    {idx + 1}
                  </div>
                  
                  <div className="space-y-5">
                    {/* Outline Orange Icon */}
                    <div className="h-12 w-12 flex items-center justify-start transition-transform duration-500 group-hover:scale-110 group-hover:rotate-3">
                      {getTugasIcon(idx)}
                    </div>
                    
                    <div className="space-y-2">
                      <h3 className="text-xs font-black text-slate-800 leading-snug tracking-tight group-hover:text-[#002147] transition-colors">
                        {tugas.title}
                      </h3>
                      <p className="text-[11px] text-slate-500 leading-relaxed font-medium">
                        {tugas.desc}
                      </p>
                    </div>
                  </div>

                  <div className="pt-6 mt-auto">
                    <span className="text-[9px] font-extrabold text-slate-400 group-hover:text-amber-600 tracking-widest uppercase transition-all duration-300 inline-flex items-center gap-1 cursor-pointer">
                      TUGAS POKOK {idx + 1} &rarr;
                    </span>
                  </div>
                </div>
              ))}
            </div>
          </div>
        ) : (
          <div className="max-w-7xl mx-auto space-y-6 animate-in fade-in duration-300">
            <div className="text-center max-w-lg mx-auto mb-8 space-y-2">
              <h2 className="text-2xl font-black text-[#002147]">Fungsi Layanan Informasi</h2>
              <p className="text-xs text-slate-400 font-medium">Pilar fungsional penjamin keterbukaan informasi publik secara akurat dan transparan.</p>
            </div>

            <div className="space-y-3">
              {displayFungsi.map((fungsi: any, idx: number) => {
                const isOpen = openAccordion === idx;
                return (
                  <div 
                    key={idx} 
                    className={`bg-white border rounded-[1.75rem] transition-all overflow-hidden ${
                      isOpen ? 'border-[#002147] shadow-sm' : 'border-slate-200 hover:border-slate-350'
                    }`}
                  >
                    {/* Header Accordion */}
                    <button
                      onClick={() => setOpenAccordion(isOpen ? null : idx)}
                      className="w-full flex items-center justify-between p-5 bg-transparent border-0 text-left cursor-pointer focus:outline-none"
                    >
                      <div className="flex items-center gap-3.5">
                        <div className={`h-9 w-9 rounded-xl flex items-center justify-center transition-all ${
                          isOpen ? 'bg-[#002147] text-white' : 'bg-slate-50 text-slate-650'
                        }`}>
                          {getFungsiIcon(idx)}
                        </div>
                        <span className="text-xs font-extrabold text-slate-800 tracking-tight leading-none">
                          {fungsi.title}
                        </span>
                      </div>
                      <ChevronRight className={`h-5 w-5 text-slate-400 transition-transform duration-300 ${
                        isOpen ? 'rotate-90 text-[#002147]' : ''
                      }`} />
                    </button>

                    {/* Content Accordion */}
                    {isOpen && (
                      <div className="px-5 pb-6 pt-1 border-t border-slate-100/60 bg-slate-50/40 text-left animate-in slide-in-from-top-1 duration-200">
                        <div className="space-y-3 pl-12">
                          {fungsi.items && fungsi.items.map((sub: string, sIdx: number) => (
                            <div key={sIdx} className="flex gap-2.5 items-start">
                              <CheckCircle className="h-4 w-4 text-emerald-500 shrink-0 mt-0.5" />
                              <span className="text-[11px] text-slate-600 font-medium leading-relaxed">
                                {sub}
                              </span>
                            </div>
                          ))}
                        </div>
                      </div>
                    )}
                  </div>
                );
              })}
            </div>
          </div>
        )}
      </div>

      {/* Landasan Hukum / Dinamik konten database di bawah */}
      {displayTambahan && displayTambahan.trim().length > 0 && (
        <div className="max-w-7xl mx-auto bg-slate-50 border border-slate-200 rounded-[2.25rem] p-8 lg:p-12 shadow-inner mt-10 px-4">
          <div className="flex items-center gap-2 pb-4 border-b border-slate-200 mb-6">
            <BookOpen className="h-5.5 w-5.5 text-amber-500" />
            <h3 className="text-xs font-extrabold text-[#002147] uppercase tracking-wider">Landasan Hukum & Keterangan Tambahan</h3>
          </div>
          
          {parsedJson && (parsedJson.image1 || parsedJson.image2 || parsedJson.image3) ? (
            <div className="grid grid-cols-1 lg:grid-cols-12 gap-10 items-center">
              {/* Kolom Kiri: Teks Landasan Hukum */}
              <div 
                className="lg:col-span-7 html-content text-xs lg:text-sm text-slate-655 leading-relaxed space-y-4 text-left w-full"
                dangerouslySetInnerHTML={{ __html: preprocessPostContent(displayTambahan) }}
              />

              {/* Kolom Kanan: Multi Image Bertumpuk Asimetris Overlap */}
              <div className="lg:col-span-5 flex justify-center py-8 lg:py-0">
                <div className="relative w-[340px] h-[340px] select-none group">
                  
                  {/* Gambar 1: Kiri/Belakang (Vertikal) */}
                  {parsedJson.image1 && (
                    <div 
                      onClick={() => setLightboxImage(resolveImageUrl(parsedJson.image1))}
                      className="absolute top-0 left-0 w-[190px] h-[240px] rounded-[2.5rem] overflow-hidden border-4 border-white shadow-lg z-10 transition-all duration-500 hover:scale-102 cursor-pointer hover:shadow-2xl hover:z-45 animate-float-slow group-hover:-translate-x-5 group-hover:-translate-y-2"
                    >
                      <img
                        src={resolveImageUrl(parsedJson.image1)}
                        alt="Landasan Hukum Gambar 1"
                        className="w-full h-full object-cover"
                      />
                    </div>
                  )}

                  {/* Gambar 2: Kanan Atas (Persegi Kecil) */}
                  {parsedJson.image2 && (
                    <div 
                      onClick={() => setLightboxImage(resolveImageUrl(parsedJson.image2))}
                      className="absolute top-6 left-[145px] w-[130px] h-[130px] rounded-[2rem] overflow-hidden border-4 border-white shadow-md z-20 transition-all duration-500 hover:scale-105 cursor-pointer hover:shadow-2xl hover:z-45 animate-float-medium group-hover:translate-x-5 group-hover:-translate-y-4"
                    >
                      <img
                        src={resolveImageUrl(parsedJson.image2)}
                        alt="Landasan Hukum Gambar 2"
                        className="w-full h-full object-cover"
                      />
                    </div>
                  )}

                  {/* Gambar 3: Tengah/Depan Bawah (Horizontal) */}
                  {parsedJson.image3 && (
                    <div 
                      onClick={() => setLightboxImage(resolveImageUrl(parsedJson.image3))}
                      className="absolute top-[135px] left-[65px] w-[230px] h-[180px] rounded-[2.5rem] overflow-hidden border-4 border-white shadow-xl z-30 transition-all duration-500 hover:scale-105 hover:shadow-2xl cursor-pointer hover:z-45 animate-float-fast group-hover:translate-y-4 group-hover:translate-x-2"
                    >
                      <img
                        src={resolveImageUrl(parsedJson.image3)}
                        alt="Landasan Hukum Gambar 3"
                        className="w-full h-full object-cover"
                      />
                    </div>
                  )}

                </div>
              </div>
            </div>
          ) : (
            <div
              className="html-content text-xs lg:text-sm text-slate-655 leading-relaxed space-y-4 text-left w-full"
              dangerouslySetInnerHTML={{ __html: preprocessPostContent(displayTambahan) }}
            />
          )}
        </div>
      )}

      {/* Lightbox Zoom Modal Overlay (Premium Backdrop Blur) */}
      {lightboxImage && (
        <div 
          className="fixed inset-0 bg-slate-900/80 backdrop-blur-md z-[9999] flex items-center justify-center p-4 animate-in fade-in duration-200"
          onClick={() => setLightboxImage(null)}
        >
          <div className="absolute top-6 right-6 z-50">
            <button 
              onClick={() => setLightboxImage(null)}
              className="p-3 bg-white/10 hover:bg-white/20 text-white rounded-full transition-colors border-0 cursor-pointer"
            >
              <X className="h-6 w-6" />
            </button>
          </div>
          <div 
            className="relative max-w-5xl max-h-[90vh] overflow-hidden rounded-2xl shadow-2xl border-4 border-white bg-slate-950 animate-in zoom-in-95 duration-250"
            onClick={(e) => e.stopPropagation()}
          >
            <img 
              src={lightboxImage} 
              alt="Zoomed Preview" 
              className="max-w-full max-h-[85vh] object-contain"
            />
          </div>
        </div>
      )}
    </div>
  );
}
