import React, { useState } from 'react';
import { FileText, ShieldAlert, ArrowDown, HelpCircle, MapPin, Mail, ClipboardCheck, CheckCircle2, Download } from 'lucide-react';
import { resolveImageUrl } from '../utils/helpers';
import type { PageData } from '../types';

interface SengketaInformasiProps {
  pageData: PageData | null;
}

interface FlowData {
  jalurA_title: string;
  jalurA_desc: string;
  jalurB_title: string;
  jalurB_desc: string;
  verifikasi_title: string;
  verifikasi_desc: string;
  help_title: string;
  help_desc: string;
  download_title: string;
  download_desc: string;
}

const defaultFlow: FlowData = {
  jalurA_title: 'Datang Langsung (Manual)',
  jalurA_desc: 'Pemohon menyerahkan surat/berkas pengajuan sengketa secara langsung ke Sekretariat Komisi Informasi Provinsi Sumatera Barat.',
  jalurB_title: 'Kirim Surat/Pos',
  jalurB_desc: 'Pemohon mengirimkan berkas pengajuan sengketa melalui jasa pengiriman surat (Pos/Kurir) resmi ditujukan ke Komisi Informasi Sumbar.',
  verifikasi_title: 'Tahap 2: Verifikasi Dokumen Pendukung oleh Petugas',
  verifikasi_desc: 'Petugas PPID / Komisi Informasi akan memverifikasi kelengkapan berkas wajib di bawah ini:',
  help_title: 'Panduan Lanjutan & Sekretariat',
  help_desc: 'Kantor Komisi Informasi Provinsi Sumatera Barat berlokasi di Kota Padang. Untuk kemudahan proses verifikasi, pemohon disarankan melampirkan salinan tanda terima keberatan yang lengkap agar berkas segera didaftarkan untuk persidangan ajudikasi non-litigasi.',
  download_title: 'Formulir Cetak',
  download_desc: 'Unduh berkas kelengkapan pengujian konsekuensi / formulir resmi pengajuan sengketa informasi.'
};

export default function SengketaInformasi({ pageData }: SengketaInformasiProps) {
  const [activeMethod, setActiveMethod] = useState<'langsung' | 'surat'>('langsung');

  const rawContent = pageData?.content || '';
  let docs: Array<{ title: string; description?: string; file_url: string }> = [];
  let introText = 'Prosedur resmi pengajuan sengketa informasi publik kepada Komisi Informasi Sumatera Barat jika keberatan pemohon ditolak atau tidak ditanggapi.';
  let parsedCards: Array<{ title: string; desc: string }> = [];
  let flow: FlowData = { ...defaultFlow };

  try {
    const parsed = JSON.parse(rawContent);
    if (parsed && typeof parsed === 'object') {
      if (Array.isArray(parsed.docs)) docs = parsed.docs;
      if (parsed.intro) introText = parsed.intro;
      if (Array.isArray(parsed.cards) && parsed.cards.length > 0) parsedCards = parsed.cards;
      if (parsed.flow && typeof parsed.flow === 'object') {
        flow = { ...defaultFlow, ...parsed.flow };
      }
    }
  } catch (e) {
    if (rawContent && !rawContent.trim().startsWith('{')) introText = rawContent;
  }

  const defaultCards = [
    { title: "Bukti Surat Permohonan Informasi & Dokumentasi", desc: "Salinan surat/bukti pengajuan awal permohonan informasi kepada Badan Publik yang bersangkutan." },
    { title: "Bukti Jawaban Permohonan Informasi", desc: "Salinan tanggapan/surat jawaban resmi (jika ada) dari Badan Publik atas permohonan informasi Anda." },
    { title: "Bukti Pengajuan Keberatan", desc: "Salinan/bukti tanda terima pengajuan surat keberatan resmi kepada Atasan PPID Badan Publik." },
    { title: "Bukti Jawaban Keberatan", desc: "Tanggapan tertulis atas keberatan (jika ada) yang dikeluarkan oleh Atasan PPID Badan Publik." },
    { title: "Bukti Identitas Pemohon", desc: "KTP/SIM yang sah untuk perorangan, atau Akta Pendirian/SK Organisasi jika mengatasnamakan lembaga." }
  ];

  const docRequirements = (parsedCards.length > 0 ? parsedCards : defaultCards).map((c, i) => ({
    no: i + 1, title: c.title, desc: c.desc
  }));

  return (
    <div className="space-y-8 animate-in fade-in duration-200 text-left w-full py-4">
      {/* Hero Banner */}
      <div className="bg-[#002147] text-white rounded-3xl p-8 lg:p-12 relative overflow-hidden shadow-lg border border-slate-800">
        <div className="absolute right-0 bottom-0 opacity-15 translate-x-12 translate-y-12">
          <ShieldAlert className="h-64 w-64 text-amber-400" />
        </div>
        <div className="relative z-10 space-y-4">
          <span className="bg-amber-400 text-[#002147] text-[10px] font-extrabold px-3 py-1.5 rounded-full uppercase tracking-wider inline-block">
            Prosedur KIP
          </span>
          <h1 className="text-3xl lg:text-4xl font-extrabold leading-tight">
            {pageData?.title || 'Permohonan Penyelesaian Sengketa'}
          </h1>
          <div 
            className="text-xs lg:text-sm text-slate-300 leading-relaxed font-medium max-w-2xl html-content"
            dangerouslySetInnerHTML={{ __html: introText }}
          />
        </div>
      </div>

      {/* Flowchart Section */}
      <div className="space-y-6">
        <div className="flex flex-col md:flex-row items-center justify-between border-b border-slate-100 pb-3">
          <h2 className="text-sm font-extrabold text-[#002147] uppercase tracking-wider flex items-center gap-2">
            <span className="h-2 w-2 rounded-full bg-amber-500 shrink-0" />
            Alur Pengajuan Sengketa
          </h2>
          <span className="text-[11px] text-slate-400 font-semibold mt-1 md:mt-0">
            Pilih metode pengajuan di bawah untuk panduan detail
          </span>
        </div>

        {/* Phase 1: Two Channels */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {/* Method A: Datang Langsung */}
          <button
            onClick={() => setActiveMethod('langsung')}
            className={`p-6 rounded-[2rem] border text-left transition-all relative overflow-hidden cursor-pointer ${
              activeMethod === 'langsung'
                ? 'bg-amber-50/50 border-amber-400 shadow-md'
                : 'bg-white border-slate-200 hover:border-slate-300 shadow-sm'
            }`}
          >
            <div className="absolute top-4 right-4 text-[10px] font-extrabold px-2.5 py-0.5 rounded-full bg-amber-400 text-[#002147]">
              Jalur A
            </div>
            <div className="flex items-start gap-4">
              <div className="h-10 w-10 rounded-2xl bg-amber-100 text-amber-700 flex items-center justify-center shrink-0">
                <MapPin className="h-5 w-5" />
              </div>
              <div className="space-y-1">
                <h3 className="text-xs font-extrabold text-[#002147]">{flow.jalurA_title}</h3>
                <p className="text-[10px] text-slate-500 font-medium leading-relaxed">
                  {flow.jalurA_desc}
                </p>
              </div>
            </div>
          </button>

          {/* Method B: Kirim Surat */}
          <button
            onClick={() => setActiveMethod('surat')}
            className={`p-6 rounded-[2rem] border text-left transition-all relative overflow-hidden cursor-pointer ${
              activeMethod === 'surat'
                ? 'bg-blue-50/50 border-blue-400 shadow-md'
                : 'bg-white border-slate-200 hover:border-slate-300 shadow-sm'
            }`}
          >
            <div className="absolute top-4 right-4 text-[10px] font-extrabold px-2.5 py-0.5 rounded-full bg-blue-500 text-white">
              Jalur B
            </div>
            <div className="flex items-start gap-4">
              <div className="h-10 w-10 rounded-2xl bg-blue-100 text-blue-750 flex items-center justify-center shrink-0">
                <Mail className="h-5 w-5" />
              </div>
              <div className="space-y-1">
                <h3 className="text-xs font-extrabold text-[#002147]">{flow.jalurB_title}</h3>
                <p className="text-[10px] text-slate-500 font-medium leading-relaxed">
                  {flow.jalurB_desc}
                </p>
              </div>
            </div>
          </button>
        </div>

        {/* Connection Arrow */}
        <div className="flex justify-center">
          <div className="flex flex-col items-center">
            <div className="h-8 w-0.5 bg-slate-200" />
            <div className="h-8 w-8 rounded-full bg-slate-100 flex items-center justify-center border border-slate-250 text-slate-500 shadow-sm">
              <ArrowDown className="h-4 w-4 animate-bounce" />
            </div>
            <div className="h-8 w-0.5 bg-slate-200" />
          </div>
        </div>

        {/* Phase 2: Verification Box */}
        <div className="bg-white border border-slate-200 rounded-[2.5rem] p-6 lg:p-8 shadow-sm space-y-6">
          <div className="flex items-start md:items-center gap-3.5 border-b border-slate-100 pb-4">
            <div className="h-11 w-11 rounded-2xl bg-amber-500 text-white flex items-center justify-center shrink-0 shadow-md shadow-amber-500/10">
              <ClipboardCheck className="h-5 w-5" />
            </div>
            <div>
              <h3 className="text-sm font-extrabold text-[#002147]">
                {flow.verifikasi_title}
              </h3>
              <p className="text-[11px] text-slate-400 font-medium">
                {flow.verifikasi_desc}
              </p>
            </div>
          </div>

          {/* Checklist Items */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {docRequirements.map((doc) => (
              <div key={doc.no} className="p-5 bg-slate-50 border border-slate-150 rounded-2xl flex flex-col justify-between text-left hover:bg-slate-100/50 transition-colors">
                <div className="space-y-2">
                  <div className="h-7 w-7 rounded-xl bg-white border border-slate-200 text-slate-700 font-extrabold text-[11px] flex items-center justify-center shadow-sm">
                    {doc.no}
                  </div>
                  <h4 className="text-xs font-extrabold text-slate-800 leading-snug">{doc.title}</h4>
                  <p className="text-[10px] text-slate-500 leading-relaxed font-medium">{doc.desc}</p>
                </div>
                <div className="mt-4 flex items-center gap-1.5 text-[9px] font-bold text-amber-700 uppercase tracking-wide">
                  <CheckCircle2 className="h-3.5 w-3.5 text-amber-500 shrink-0" /> Dokumen Wajib
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Info & Download Footer */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 pt-4">
        {/* Help box */}
        <div className="lg:col-span-2 bg-slate-50 border border-slate-200 rounded-3xl p-6 flex items-start gap-4">
          <HelpCircle className="h-6 w-6 text-amber-500 shrink-0 mt-0.5" />
          <div className="space-y-1">
            <h4 className="text-xs font-extrabold text-[#002147]">{flow.help_title}</h4>
            <p className="text-[11px] text-slate-500 font-medium leading-relaxed">
              {flow.help_desc}
            </p>
          </div>
        </div>

        {/* Download box */}
        <div className="bg-[#002147] text-white rounded-3xl p-6 flex flex-col justify-between space-y-4 shadow-md">
          <div className="space-y-1">
            <h4 className="text-xs font-extrabold text-amber-400 uppercase tracking-wider">{flow.download_title}</h4>
            <p className="text-[10px] text-slate-350 leading-relaxed font-medium">
              {flow.download_desc}
            </p>
          </div>
          <div className="space-y-2">
            {docs.length > 0 ? (
              docs.map((doc, dIdx) => (
                <a
                  key={dIdx}
                  href={resolveImageUrl(doc.file_url)}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="w-full flex items-center justify-center gap-2 py-2.5 rounded-2xl bg-amber-400 hover:bg-amber-500 text-[#002147] font-extrabold text-xs transition-colors border-0 cursor-pointer shadow-sm shadow-amber-400/20 text-center"
                >
                  <Download className="h-4 w-4" /> {doc.title || 'Unduh Berkas'}
                </a>
              ))
            ) : (
              <a
                href="/uploads/downloads/1782517011887620200.docx"
                className="w-full flex items-center justify-center gap-2 py-2.5 rounded-2xl bg-amber-400 hover:bg-amber-500 text-[#002147] font-extrabold text-xs transition-colors border-0 cursor-pointer shadow-sm shadow-amber-400/20"
              >
                <Download className="h-4 w-4" /> Unduh Form Sengketa
              </a>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
