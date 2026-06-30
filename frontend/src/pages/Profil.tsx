import { useState, useEffect } from 'react';
import {
  Users, History, Award, CheckCircle, Info, Target, Compass, Network, ShieldCheck, X
} from 'lucide-react';
import { resolveImageUrl, preprocessPostContent } from '../utils/helpers';
import type { PageData } from '../types';

interface ProfilProps {
  pageData: PageData;
}

export default function Profil({ pageData }: ProfilProps) {
  const [extraData, setExtraData] = useState<{
    visiMisi: PageData | null;
    maklumat: PageData | null;
    struktur: PageData | null;
  }>({ visiMisi: null, maklumat: null, struktur: null });
  const [loading, setLoading] = useState(true);
  const [lightboxImage, setLightboxImage] = useState<string | null>(null);

  useEffect(() => {
    const apiBase = import.meta.env.VITE_API_URL || '/api/v1';
    Promise.all([
      fetch(`${apiBase}/pages/visi-misi`).then(res => res.ok ? res.json() : null).catch(() => null),
      fetch(`${apiBase}/pages/maklumat`).then(res => res.ok ? res.json() : null).catch(() => null),
      fetch(`${apiBase}/pages/struktur-organisasi-2`).then(res => res.ok ? res.json() : null).catch(() => null),
    ]).then(([visiMisi, maklumat, struktur]) => {
      setExtraData({ visiMisi, maklumat, struktur });
      setLoading(false);
    });
  }, []);

  const scrollToSection = (id: string) => {
    const element = document.getElementById(id);
    if (element) {
      element.scrollIntoView({ behavior: 'smooth', block: 'start' });
    }
  };

  // Helper to check if string is JSON
  let parsedJson: any = null;
  try {
    parsedJson = JSON.parse(pageData.content || '');
  } catch (e) {
    parsedJson = null;
  }

  // Parse Struktur Organisasi JSON
  let ppidAtasan = 'Rektor Universitas Perintis';
  let ppidUtama = 'Wakil Rektor 1 & 2 UPERTIS';
  let ppidPelaksana = ["Biro Humas", "Biro Akademik", "Biro Umum", "Dekan Fakultas"];
  let ppidPertimbangan = ["ka P2AMIA", "ka LPPM", "ka P3TS", "Ka Prodi", "Ka. UPT"];
  let ppidPelayanan = ["Staff Humas", "LTIK"];
  let ppidDesc = '';

  if (extraData.struktur?.content) {
    const trimmed = extraData.struktur.content.trim();
    if (trimmed.startsWith('{')) {
      try {
        const parsed = JSON.parse(trimmed);
        if (parsed && typeof parsed === 'object') {
          ppidAtasan = parsed.atasan || ppidAtasan;
          ppidUtama = parsed.utama || ppidUtama;
          if (Array.isArray(parsed.pelaksana)) ppidPelaksana = parsed.pelaksana;
          if (Array.isArray(parsed.pertimbangan)) ppidPertimbangan = parsed.pertimbangan;
          if (Array.isArray(parsed.pelayanan)) ppidPelayanan = parsed.pelayanan;
          ppidDesc = parsed.desc || '';
        }
      } catch (e) {
        console.error('Failed to parse dynamic struktur JSON:', e);
        ppidDesc = extraData.struktur.content;
      }
    } else {
      ppidDesc = extraData.struktur.content;
    }
  }

  // Parse Visi Misi JSON
  let ppidVisi = '';
  let ppidMisi = [
    "Menyediakan pelayanan informasi publik yang cepat, tepat waktu, dan akurat.",
    "Mengembangkan sistem pengelolaan dokumen berbasis teknologi informasi yang aman dan mudah diakses.",
    "Meningkatkan kapasitas sumber daya pengelola layanan informasi secara berkelanjutan.",
    "Mewujudkan tata kelola perguruan tinggi yang bersih, transparan, dan akuntabel."
  ];

  if (extraData.visiMisi?.content) {
    const trimmed = extraData.visiMisi.content.trim();
    if (trimmed.startsWith('{')) {
      try {
        const parsed = JSON.parse(trimmed);
        if (parsed && typeof parsed === 'object') {
          ppidVisi = parsed.visi || '';
          if (Array.isArray(parsed.misi)) ppidMisi = parsed.misi;
        }
      } catch (e) {
        console.error('Failed to parse dynamic visi-misi JSON:', e);
        ppidVisi = extraData.visiMisi.content;
      }
    } else {
      ppidVisi = extraData.visiMisi.content;
    }
  } else {
    // Check if parent pageData (profil) has it (for legacy config)
    if (parsedJson?.visi) {
      ppidVisi = parsedJson.visi;
    } else {
      ppidVisi = "Menjadi Pejabat Pengelola Informasi dan Dokumentasi (PPID) yang unggul, terpercaya, dan transparan dalam pelayanan informasi publik di lingkungan Universitas Perintis Indonesia.";
    }
    if (parsedJson?.misi && Array.isArray(parsedJson.misi)) {
      ppidMisi = parsedJson.misi;
    }
  }



  return (
    <div className="space-y-12 animate-in fade-in duration-200 text-left w-full py-6">
      {/* Premium Hero Banner */}
      <div className="bg-[#002147] text-white rounded-[2rem] p-10 lg:p-16 relative overflow-hidden shadow-xl border border-slate-800">
        <div className="absolute right-0 bottom-0 opacity-10 translate-x-12 translate-y-12">
          <Users className="h-96 w-96 text-amber-400" />
        </div>
        <div className="absolute top-0 left-0 w-full h-full bg-gradient-to-br from-transparent via-transparent to-amber-400/5 pointer-events-none" />
        <div className="relative z-10 space-y-4 max-w-3xl">
          <span className="bg-amber-400 text-[#002147] text-[10px] font-extrabold px-3.5 py-1.5 rounded-full uppercase tracking-widest font-mono">
            Profil & Tata Pamong
          </span>
          <h1 className="text-4xl lg:text-5xl font-black leading-tight tracking-tight">
            Sekilas Profil PPID UPERTIS
          </h1>
          <p className="text-xs lg:text-sm text-slate-350 leading-relaxed font-medium">
            Komitmen transparansi informasi publik di bawah naungan Universitas Perintis Indonesia, menyajikan tata pamong organisasi yang akuntabel.
          </p>
        </div>
      </div>

      {/* Quick Navigation Anchor Links (Luxury sticky pill menu) */}
      <div className="sticky top-20 z-40 bg-white/80 backdrop-blur-md border border-slate-200 rounded-2xl p-2.5 flex flex-wrap gap-2 shadow-sm justify-center md:justify-start">
        <button
          onClick={() => scrollToSection('sejarah')}
          className="px-4 py-2 hover:bg-slate-50 text-[#002147] font-bold text-xs rounded-xl cursor-pointer transition-all border-0 bg-transparent"
        >
          Sejarah Singkat
        </button>
        <button
          onClick={() => scrollToSection('visi-misi')}
          className="px-4 py-2 hover:bg-slate-50 text-[#002147] font-bold text-xs rounded-xl cursor-pointer transition-all border-0 bg-transparent"
        >
          Visi & Misi
        </button>
        <button
          onClick={() => scrollToSection('struktur')}
          className="px-4 py-2 hover:bg-slate-50 text-[#002147] font-bold text-xs rounded-xl cursor-pointer transition-all border-0 bg-transparent"
        >
          Struktur Organisasi
        </button>
        <button
          onClick={() => scrollToSection('maklumat')}
          className="px-4 py-2 hover:bg-slate-50 text-[#002147] font-bold text-xs rounded-xl cursor-pointer transition-all border-0 bg-transparent"
        >
          Maklumat Pelayanan
        </button>
      </div>

      {/* --- SECTION 1: SEJARAH --- */}
      <section id="sejarah" className="scroll-mt-36 bg-white border border-slate-200 rounded-[2.5rem] p-8 lg:p-12 shadow-sm space-y-6">
        <div className="border-b border-slate-100 pb-4">
          <h2 className="text-xl font-extrabold text-[#002147] flex items-center gap-2.5">
            <History className="h-5.5 w-5.5 text-amber-500" /> Sejarah Singkat
          </h2>
          <p className="text-xs text-slate-400 mt-1 font-medium">Awal mula pendirian dan visi transparansi informasi PPID UPERTIS.</p>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-12 gap-10 items-center pt-4">
          {/* Kolom Kiri: Teks Sejarah & Kartu Statistik */}
          <div className="lg:col-span-7 space-y-6">
            <div className="text-xs text-slate-600 font-medium leading-relaxed space-y-4">
              {parsedJson && parsedJson.sejarah ? (
                <div dangerouslySetInnerHTML={{ __html: preprocessPostContent(parsedJson.sejarah) }} />
              ) : (
                <div dangerouslySetInnerHTML={{ __html: preprocessPostContent(pageData.content || '') }} />
              )}
            </div>

            {/* Stats Cards inline di bawah teks sejarah */}
            {parsedJson && Array.isArray(parsedJson.stats) && parsedJson.stats.length > 0 ? (
              <div className="grid grid-cols-2 sm:grid-cols-3 gap-4 pt-4 border-t border-slate-100">
                {parsedJson.stats.map((st: any, idx: number) => (
                  <div key={idx} className="p-4 bg-slate-50 border border-slate-200/80 rounded-2xl text-center space-y-0.5 hover:border-amber-300 transition-all hover:shadow-sm">
                    <span className="text-2xl font-black text-[#002147] block font-mono">{st.value}</span>
                    <span className="text-[9px] text-slate-400 font-extrabold uppercase tracking-widest block">{st.label}</span>
                  </div>
                ))}
              </div>
            ) : (
              <div className="p-4 bg-gradient-to-br from-blue-50/50 to-indigo-50/20 border border-blue-150 rounded-2xl space-y-2 pt-4 border-t border-slate-100">
                <span className="text-[10px] font-extrabold text-[#002147] uppercase tracking-widest block">Status Layanan</span>
                <p className="text-[11px] text-slate-500 leading-relaxed font-medium">
                  PPID berkomitmen penuh untuk memproses permohonan informasi publik secara transparan sesuai UU KIP No. 14 Tahun 2008.
                </p>
              </div>
            )}
          </div>

          {/* Kolom Kanan: Multi Image Bertumpuk Asimetris Overlap (Sesuai Screenshot) */}
          <div className="lg:col-span-5 flex justify-center py-8 lg:py-0">
            <div className="relative w-[340px] h-[340px] select-none group">
              
              {/* Gambar 1: Kiri/Belakang (Vertikal) */}
              <div 
                onClick={() => setLightboxImage(parsedJson && parsedJson.image1 ? resolveImageUrl(parsedJson.image1) : "https://images.unsplash.com/photo-1517245386807-bb43f82c33c4?auto=format&fit=crop&q=80&w=600")}
                className="absolute top-0 left-0 w-[190px] h-[240px] rounded-[2.5rem] overflow-hidden border-4 border-white shadow-lg z-10 transition-all duration-500 hover:scale-102 cursor-pointer hover:shadow-2xl hover:z-45 animate-float-slow group-hover:-translate-x-5 group-hover:-translate-y-2"
              >
                <img
                  src={parsedJson && parsedJson.image1 ? resolveImageUrl(parsedJson.image1) : "https://images.unsplash.com/photo-1517245386807-bb43f82c33c4?auto=format&fit=crop&q=80&w=600"}
                  alt="Sejarah Kampus UPERTIS"
                  className="w-full h-full object-cover"
                />
              </div>

              {/* Gambar 2: Kanan Atas (Persegi Kecil) */}
              <div 
                onClick={() => setLightboxImage(parsedJson && parsedJson.image2 ? resolveImageUrl(parsedJson.image2) : "https://images.unsplash.com/photo-1522202176988-66273c2fd55f?auto=format&fit=crop&q=80&w=400")}
                className="absolute top-6 left-[145px] w-[130px] h-[130px] rounded-[2rem] overflow-hidden border-4 border-white shadow-md z-20 transition-all duration-500 hover:scale-105 cursor-pointer hover:shadow-2xl hover:z-45 animate-float-medium group-hover:translate-x-5 group-hover:-translate-y-4"
              >
                <img
                  src={parsedJson && parsedJson.image2 ? resolveImageUrl(parsedJson.image2) : "https://images.unsplash.com/photo-1522202176988-66273c2fd55f?auto=format&fit=crop&q=80&w=400"}
                  alt="Aktivitas Pelayanan PPID"
                  className="w-full h-full object-cover"
                />
              </div>

              {/* Gambar 3: Tengah/Depan Bawah (Horizontal) */}
              <div 
                onClick={() => setLightboxImage(parsedJson && parsedJson.image3 ? resolveImageUrl(parsedJson.image3) : (pageData.cover_image_url ? resolveImageUrl(pageData.cover_image_url) : "https://images.unsplash.com/photo-1497366811353-6870744d04b2?auto=format&fit=crop&q=80&w=600"))}
                className="absolute top-[135px] left-[65px] w-[230px] h-[180px] rounded-[2.5rem] overflow-hidden border-4 border-white shadow-xl z-30 transition-all duration-500 hover:scale-105 hover:shadow-2xl cursor-pointer hover:z-45 animate-float-fast group-hover:translate-y-4 group-hover:translate-x-2"
              >
                <img
                  src={parsedJson && parsedJson.image3 ? resolveImageUrl(parsedJson.image3) : (pageData.cover_image_url ? resolveImageUrl(pageData.cover_image_url) : "https://images.unsplash.com/photo-1497366811353-6870744d04b2?auto=format&fit=crop&q=80&w=600")}
                  alt="Gedung Rektorat UPERTIS"
                  className="w-full h-full object-cover"
                />
              </div>

            </div>
          </div>
        </div>
      </section>

      {/* --- SECTION 2: VISI & MISI --- */}
      <section id="visi-misi" className="scroll-mt-36 grid grid-cols-1 lg:grid-cols-2 gap-8">
        {/* Visi Card */}
        <div className="bg-gradient-to-br from-[#002147] to-[#003166] text-white p-8 lg:p-12 rounded-[2.5rem] shadow-md border border-slate-800 space-y-6 flex flex-col justify-between">
          <div className="space-y-4">
            <span className="bg-amber-400 text-[#002147] text-[9px] font-extrabold px-3 py-1 rounded-full uppercase tracking-wider inline-block">
              Visi Utama
            </span>
            <h3 className="text-xl lg:text-2xl font-extrabold flex items-center gap-2">
              <Target className="h-6 w-6 text-amber-400" /> Visi PPID
            </h3>
            <p className="text-xs lg:text-sm text-slate-200 leading-relaxed font-medium font-serif italic pt-2">
              {ppidVisi ? `"${ppidVisi}"` : '"-visi belum diisi-"'}
            </p>
          </div>
        </div>

        {/* Misi Card */}
        <div className="bg-white border border-slate-200 p-8 lg:p-12 rounded-[2.5rem] shadow-sm space-y-6">
          <div className="space-y-1">
            <span className="bg-blue-50 text-blue-700 border border-blue-100 text-[9px] font-extrabold px-3 py-1 rounded-full uppercase tracking-wider inline-block">
              Misi Operasional
            </span>
            <h3 className="text-xl font-extrabold text-[#002147] flex items-center gap-2 pt-2">
              <Compass className="h-6 w-6 text-amber-500" /> Misi PPID
            </h3>
          </div>
          <div className="space-y-4 text-xs text-slate-655 font-medium">
            {ppidMisi.map((m: string, idx: number) => (
              <div key={idx} className="flex gap-3 items-start text-left">
                <span className="h-6 w-6 rounded-lg bg-slate-50 border border-slate-200 text-[#002147] font-black text-[10px] flex items-center justify-center shrink-0 mt-0.5 shadow-sm">
                  {idx + 1}
                </span>
                <span className="leading-relaxed">{m}</span>
              </div>
            ))}
            {ppidMisi.length === 0 && (
              <p className="text-[11px] text-slate-400 italic">Belum ada misi ditambahkan.</p>
            )}
          </div>
        </div>
      </section>

      {/* --- SECTION 3: STRUKTUR ORGANISASI --- */}
      <section id="struktur" className="scroll-mt-36 bg-white border border-slate-200 rounded-[2.5rem] p-8 lg:p-12 shadow-sm space-y-6">
        <div className="border-b border-slate-100 pb-4">
          <h2 className="text-xl font-extrabold text-[#002147] flex items-center gap-2.5">
            <Network className="h-5.5 w-5.5 text-amber-500" /> Struktur Organisasi PPID
          </h2>
          <p className="text-xs text-slate-400 mt-1 font-medium">Bagan keanggotaan dan alur koordinasi Pejabat Pengelola Informasi.</p>
        </div>

        <div className="p-8 bg-slate-50 border border-slate-150 rounded-[2rem] space-y-8 w-full shadow-inner relative overflow-hidden">
          {/* Level 1: Atasan PPID */}
          <div className="flex flex-col items-center relative">
            <div className="group relative bg-[#002147] text-white p-5 rounded-2xl shadow-md border border-amber-400 w-72 text-center transition-all hover:scale-105 hover:shadow-lg hover:border-amber-300">
              <span className="text-[9px] text-amber-400 font-extrabold uppercase tracking-widest font-mono block">ATASAN PPID</span>
              <h4 className="text-xs font-black mt-1 leading-snug">{ppidAtasan}</h4>
              <p className="text-[10px] text-slate-350 font-medium mt-1 leading-relaxed border-t border-white/10 pt-1.5">
                {ppidAtasan}
              </p>
            </div>
            
            {/* Connector Line 1-2 */}
            <div className="h-6 w-0.5 bg-slate-300 relative" />
          </div>

          {/* Level 2: PPID Utama */}
          <div className="flex flex-col items-center relative">
            <div className="group relative bg-[#002147] text-white p-5 rounded-2xl shadow-md border border-blue-400 w-72 text-center transition-all hover:scale-105 hover:shadow-lg hover:border-blue-300">
              <span className="text-[9px] text-blue-300 font-extrabold uppercase tracking-widest font-mono block">PPID UTAMA</span>
              <h4 className="text-xs font-black mt-1 leading-snug">{ppidUtama}</h4>
              <p className="text-[10px] text-slate-350 font-medium mt-1 leading-relaxed border-t border-white/10 pt-1.5">
                {ppidUtama}
              </p>
            </div>

            {/* Connector Line 2-3 with 3 branches */}
            <div className="w-full relative hidden md:flex flex-col items-center">
              {/* Vertical line down from PPID Utama */}
              <div className="h-6 w-0.5 bg-slate-300" />
              
              {/* Horizontal line branching left and right */}
              <div className="w-2/3 h-0.5 bg-slate-300 flex justify-between relative">
                {/* Left corner arrow to Pelaksana */}
                <div className="h-4 w-0.5 bg-slate-300" />
                {/* Right corner arrow to Pelayanan */}
                <div className="h-4 w-0.5 bg-slate-300" />
              </div>

              {/* Middle vertical line for Tim Pertimbangan */}
              <div className="absolute top-6 left-1/2 -translate-x-1/2 h-4 w-0.5 bg-slate-300" />
            </div>
          </div>

          {/* Level 3: The Three Columns (Pelaksana, Pertimbangan, Petugas Pelayanan) */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6 max-w-5xl mx-auto pt-2">
            
            {/* Column 1: PPID Pelaksana */}
            <div className="flex flex-col items-center">
              {/* Mobile-only connector */}
              <div className="h-4 w-0.5 bg-slate-300 md:hidden" />
              <div className="group bg-white border border-slate-200 p-5 rounded-2xl shadow-sm hover:shadow-md transition-all text-center w-full min-h-[160px] flex flex-col justify-between hover:border-slate-300">
                <div>
                  <span className="text-[8px] text-blue-600 font-extrabold uppercase tracking-wider font-mono block">PPID PELAKSANA</span>
                  <h4 className="text-xs font-black text-slate-800 mt-1">Biro & Dekan</h4>
                </div>
                <div className="text-[10px] text-slate-500 font-semibold border-t border-slate-100 pt-2.5 mt-2 space-y-1 text-left">
                  {ppidPelaksana.map((item, idx) => (
                    <div key={idx} className="flex items-center gap-1.5">
                      <span className="h-1 w-1 rounded-full bg-blue-500" />
                      <span>{item}</span>
                    </div>
                  ))}
                </div>
              </div>
            </div>

            {/* Column 2: Tim Pertimbangan */}
            <div className="flex flex-col items-center">
              {/* Mobile-only connector */}
              <div className="h-4 w-0.5 bg-slate-300 md:hidden" />
              <div className="group bg-white border border-slate-200 p-5 rounded-2xl shadow-sm hover:shadow-md transition-all text-center w-full min-h-[160px] flex flex-col justify-between hover:border-slate-300">
                <div>
                  <span className="text-[8px] text-amber-600 font-extrabold uppercase tracking-wider font-mono block">TIM PERTIMBANGAN</span>
                  <h4 className="text-xs font-black text-slate-800 mt-1">Komite Pertimbangan</h4>
                </div>
                <div className="text-[10px] text-slate-500 font-semibold border-t border-slate-100 pt-2.5 mt-2 space-y-1 text-left">
                  {ppidPertimbangan.map((item, idx) => (
                    <div key={idx} className="flex items-center gap-1.5">
                      <span className="h-1 w-1 rounded-full bg-amber-500" />
                      <span>{item}</span>
                    </div>
                  ))}
                </div>
              </div>
            </div>

            {/* Column 3: Petugas Pelayanan */}
            <div className="flex flex-col items-center">
              {/* Mobile-only connector */}
              <div className="h-4 w-0.5 bg-slate-300 md:hidden" />
              <div className="group bg-white border border-slate-200 p-5 rounded-2xl shadow-sm hover:shadow-md transition-all text-center w-full min-h-[160px] flex flex-col justify-between hover:border-slate-300">
                <div>
                  <span className="text-[8px] text-emerald-600 font-extrabold uppercase tracking-wider font-mono block">PETUGAS PELAYANAN</span>
                  <h4 className="text-xs font-black text-slate-800 mt-1">Staf Desk Pelayanan</h4>
                </div>
                <div className="text-[10px] text-slate-500 font-semibold border-t border-slate-100 pt-2.5 mt-2 space-y-1 text-left">
                  {ppidPelayanan.map((item, idx) => (
                    <div key={idx} className="flex items-center gap-1.5">
                      <span className="h-1 w-1 rounded-full bg-emerald-500" />
                      <span>{item}</span>
                    </div>
                  ))}
                </div>
              </div>
            </div>

          </div>
        </div>

        {ppidDesc && (
          <div
            className="html-content text-xs text-slate-600 leading-relaxed space-y-4 pt-6 border-t border-slate-100 max-w-2xl mx-auto"
            dangerouslySetInnerHTML={{ __html: preprocessPostContent(ppidDesc) }}
          />
        )}
      </section>

      {/* --- SECTION 4: MAKLUMAT PELAYANAN --- */}
      <section id="maklumat" className="scroll-mt-36 bg-white border border-slate-200 rounded-[2.5rem] p-8 lg:p-12 shadow-sm space-y-6">
        <div className="border-b border-slate-100 pb-4">
          <h2 className="text-xl font-extrabold text-[#002147] flex items-center gap-2.5">
            <ShieldCheck className="h-5.5 w-5.5 text-amber-500" /> Maklumat Pelayanan
          </h2>
          <p className="text-xs text-slate-400 mt-1 font-medium">Ikrar komitmen janji pelayanan informasi publik yang akuntabel.</p>
        </div>

        <div className="relative border-8 border-double border-amber-200/60 rounded-3xl p-8 lg:p-14 bg-gradient-to-br from-slate-50 to-amber-50/10 text-center space-y-5 shadow-inner max-w-2xl mx-auto animate-in fade-in zoom-in-95 duration-200">
          <div className="mx-auto h-12 w-12 rounded-full bg-amber-50 text-amber-600 border border-amber-200 flex items-center justify-center font-bold text-lg shadow-sm">
            ★
          </div>
          <h3 className="font-extrabold text-[#002147] tracking-widest text-xs uppercase font-mono">MAKLUMAT PELAYANAN INFORMASI</h3>
          <p className="text-slate-750 italic font-serif leading-relaxed text-xs lg:text-sm max-w-lg mx-auto">
            "Kami berkomitmen memberikan pelayanan informasi publik yang cepat, tepat, transparan, dan akuntabel sesuai dengan standar operasional prosedur demi mewujudkan keterbukaan informasi publik di lingkungan Universitas Perintis Indonesia."
          </p>
        </div>

        {extraData.maklumat?.content && (
          <div
            className="html-content text-xs text-slate-655 leading-relaxed space-y-4 pt-6 border-t border-slate-100 max-w-2xl mx-auto"
            dangerouslySetInnerHTML={{ __html: preprocessPostContent(extraData.maklumat.content) }}
          />
        )}
      </section>

      {/* Info Footer Note */}
      <div className="bg-slate-50 border border-slate-200 rounded-2xl p-5 flex items-start gap-3.5 text-left">
        <Info className="h-5 w-5 text-amber-500 shrink-0 mt-0.5" />
        <div className="space-y-1">
          <h4 className="text-xs font-bold text-[#002147]">Manajemen Halaman Profil Terpadu</h4>
          <p className="text-[11px] text-slate-500 font-medium leading-relaxed">
            Data pada setiap sub-bagian di atas terhubung langsung dengan sistem pengelolaan di Dashboard Admin. Anda tetap dapat mengedit halaman Profil, Visi Misi, Maklumat, dan Struktur Organisasi secara mandiri, dan sistem akan menyajikannya secara otomatis dalam satu layout terpadu ini.
          </p>
        </div>
      </div>

      {/* Lightbox Modal */}
      {lightboxImage && (
        <div
          className="fixed inset-0 z-[100] bg-slate-950/90 backdrop-blur-sm flex items-center justify-center p-4 animate-in fade-in duration-300"
          onClick={() => setLightboxImage(null)}
        >
          <button
            onClick={() => setLightboxImage(null)}
            className="absolute top-6 right-6 w-12 h-12 bg-white/10 hover:bg-white/20 text-white rounded-full flex items-center justify-center transition-colors cursor-pointer border-0"
          >
            <X className="h-6 w-6" />
          </button>
          <img
            src={lightboxImage}
            alt="Zoomed Sejarah"
            className="max-w-full max-h-[85vh] object-contain rounded-2xl shadow-2xl border border-white/10 animate-in zoom-in-95 duration-300"
          />
        </div>
      )}
    </div>
  );
}
