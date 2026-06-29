import React, { useState } from 'react';
import {
  FileText, ShieldAlert, HelpCircle, MapPin, Mail, ClipboardList, FileCheck, Send, Clock, CheckCircle2, UserCheck, Globe, Download
} from 'lucide-react';
import DownloadTable from '../components/DownloadTable';
import SerpentineFlowchart from '../components/SerpentineFlowchart';
import { resolveImageUrl, preprocessPostContent } from '../utils/helpers';
import type { PageData } from '../types';

interface KeberatanInformasiProps {
  pageData: PageData;
  searchTerm: string;
  setSearchTerm: (term: string) => void;
}

export default function KeberatanInformasi({
  pageData,
  searchTerm,
  setSearchTerm
}: KeberatanInformasiProps) {
  const rawContent = pageData?.content || '';
  let docs: Array<{ title: string; description?: string; file_url: string }> = [];
  let introText = '';
  let isJson = false;
  let parsed: any = null;

  try {
    parsed = JSON.parse(rawContent);
    if (parsed && typeof parsed === 'object') {
      if (Array.isArray(parsed.docs)) {
        docs = parsed.docs;
      } else if (Array.isArray(parsed)) {
        docs = parsed;
      }
      introText = parsed.intro || '';
      isJson = true;
    }
  } catch (e) {
    isJson = false;
    introText = rawContent;
  }

  const [activeTab, setActiveTab] = useState<'manual' | 'online'>('manual');
  const [activeStepIdx, setActiveStepIdx] = useState<number>(0);

  const defaultManualSteps = [
    { title: "Pemohon Informasi", desc: "Pemohon datang langsung ke sekretariat PPID-UPERTIS untuk mengajukan keberatan atas informasi publik yang diterima." },
    { title: "Mengisi Form Keberatan", desc: "Petugas informasi menerima berkas pengajuan dan pemohon mengisi formulir keberatan dengan melampirkan identitas diri berupa KTP/SIM atau surat tugas dari institusi lainnya." },
    { title: "Menerima Tanda Terima", desc: "Petugas informasi memberikan tanda bukti resmi bahwa pemohon telah mengajukan keberatan atas layanan informasi publik." },
    { title: "Keberatan Di Proses", desc: "PPID Pelaksana meneruskan pengaduan ke Rektor/PPID Utama UPERTIS untuk diproses selambat-lambatnya 30 hari kerja." },
    { title: "Penyerahan Jawaban", desc: "Jawaban/tanggapan tertulis atas keberatan disampaikan kepada pemohon melalui PPID Pelaksana." },
    { title: "Layanan Selesai", desc: "Prosedur pengajuan keberatan informasi secara manual selesai dilakukan." }
  ];

  const defaultOnlineSteps = [
    { title: "Mengunjungi Website", desc: "Pemohon mengunjungi website resmi PPID Universitas Perintis Indonesia di ppid.upertis.ac.id." },
    { title: "Mengunduh Form", desc: "Pemohon mengunduh berkas formulir pengajuan keberatan (.doc) dari menu unduhan di website." },
    { title: "Mengirimkan Berkas", desc: "Pemohon mengirim formulir terisi beserta scan identitas (KTP/SIM) ke email resmi PPID UPERTIS di ppidcare@upertis.ac.id." },
    { title: "Menerima Tanda Terima", desc: "Petugas informasi mengirimkan konfirmasi tanda terima pengajuan keberatan via email." },
    { title: "Permohonan Di Proses", desc: "Berkas keberatan diteruskan ke Rektor UPERTIS untuk ditinjau dan diproses dalam kurun waktu 30 hari kerja." },
    { title: "Penyerahan Jawaban", desc: "Jawaban keputusan dari Rektor UPERTIS dikirimkan kepada pemohon melalui PPID Pelaksana." },
    { title: "Layanan Selesai", desc: "Prosedur pengajuan keberatan informasi secara online selesai dilakukan." }
  ];

  const rawManual = (isJson && (parsed as any).manual_steps && Array.isArray((parsed as any).manual_steps)) ? (parsed as any).manual_steps : defaultManualSteps;
  const rawOnline = (isJson && (parsed as any).online_steps && Array.isArray((parsed as any).online_steps)) ? (parsed as any).online_steps : defaultOnlineSteps;

  const stepIcons = [MapPin, ClipboardList, FileCheck, Send, Clock, UserCheck, CheckCircle2];
  const stepColors = [
    "border-blue-500 text-blue-600",
    "border-amber-500 text-amber-600",
    "border-teal-500 text-teal-600",
    "border-indigo-500 text-indigo-600",
    "border-purple-500 text-purple-600",
    "border-pink-500 text-pink-600",
    "border-emerald-500 text-emerald-600"
  ];

  const manualSteps = rawManual.map((step: any, idx: number) => {
    const iconIdx = idx % stepIcons.length;
    const colorIdx = idx % stepColors.length;
    return {
      step: idx + 1,
      title: step.title || `Langkah ${idx + 1}`,
      desc: step.desc || "",
      icon: stepIcons[iconIdx],
      color: stepColors[colorIdx],
      bullets: [step.desc || ""]
    };
  });

  const onlineSteps = rawOnline.map((step: any, idx: number) => {
    const iconIdx = idx % stepIcons.length;
    const colorIdx = idx % stepColors.length;
    return {
      step: idx + 1,
      title: step.title || `Langkah ${idx + 1}`,
      desc: step.desc || "",
      icon: stepIcons[iconIdx],
      color: stepColors[colorIdx],
      bullets: [step.desc || ""]
    };
  });

  // Fallback to normal rendering if not JSON
  if (!isJson) {
    return (
      <article className="space-y-6 text-left w-full">
        <div className="bg-[#002147] text-white rounded-3xl p-8 lg:p-12 relative overflow-hidden shadow-lg">
          <div className="absolute right-0 bottom-0 opacity-10 translate-x-12 translate-y-12">
            <ShieldAlert className="h-64 w-64 text-amber-400" />
          </div>
          <div className="relative z-10 space-y-3">
            <span className="bg-amber-400 text-[#002147] text-[10px] font-extrabold px-3 py-1 rounded-full uppercase tracking-wider">
              Layanan Keberatan
            </span>
            <h1 className="text-3xl lg:text-4xl font-extrabold leading-tight">{pageData?.title || 'Pengajuan Keberatan Informasi'}</h1>
          </div>
        </div>
        {pageData?.cover_image_url && (
          <img
            src={pageData.cover_image_url}
            alt={pageData.title}
            className="w-full h-72 lg:h-96 object-cover rounded-3xl shadow-sm border border-slate-200 mt-6"
          />
        )}
        <div
          className="border-t border-slate-200 pt-6 html-content text-sm text-slate-700 leading-relaxed space-y-4"
          dangerouslySetInnerHTML={{ __html: preprocessPostContent(rawContent) }}
        />
      </article>
    );
  }

  const filteredDocs = docs.filter(item => {
    return item.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
      (item.description && item.description.toLowerCase().includes(searchTerm.toLowerCase()));
  });

  const handleTabChange = (tab: 'manual' | 'online') => {
    setActiveTab(tab);
    setActiveStepIdx(0);
  };

  const steps = activeTab === 'manual' ? manualSteps : onlineSteps;
  const currentActiveStep = steps[activeStepIdx] || steps[0];

  return (
    <div className="space-y-8 animate-in fade-in duration-200 text-left w-full py-4">
      {/* Hero Banner */}
      <div className="bg-[#002147] text-white rounded-3xl p-8 lg:p-12 relative overflow-hidden shadow-lg border border-slate-800">
        <div className="absolute right-0 bottom-0 opacity-15 translate-x-12 translate-y-12">
          <ShieldAlert className="h-64 w-64 text-amber-400" />
        </div>
        <div className="relative z-10 space-y-4">
          <span className="bg-amber-400 text-[#002147] text-[10px] font-extrabold px-3 py-1.5 rounded-full uppercase tracking-wider inline-block">
            Layanan Keberatan
          </span>
          <h1 className="text-3xl lg:text-4xl font-extrabold leading-tight">
            {pageData?.title || 'Pengajuan Keberatan Informasi'}
          </h1>
          <p className="text-xs lg:text-sm text-slate-300 leading-relaxed font-medium max-w-2xl">
            {pageData?.subtitle || 'Mekanisme keberatan bagi pemohon informasi yang tidak puas dengan tanggapan, jangka waktu, atau pelayanan dari Pejabat Pengelola Informasi (PPID).'}
          </p>
        </div>
      </div>

      {/* Main Grid: Downloads on Left & Flowchart on Right */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 items-start">
        {/* Left Column: Download forms & Help */}
        <div className="lg:col-span-4 space-y-6">
          {/* Docs Table */}
          <div className="bg-white border border-slate-200 rounded-3xl p-6 shadow-sm space-y-4">
            <div className="flex items-center gap-2 border-b border-slate-100 pb-3">
              <FileText className="h-4.5 w-4.5 text-amber-500" />
              <h3 className="text-xs font-bold text-[#002147] uppercase tracking-wider">Formulir Unduhan</h3>
            </div>
            <DownloadTable
              items={filteredDocs.map(item => ({
                ...item,
                type: 'DownloadItem'
              }))}
              searchTerm={searchTerm}
              setSearchTerm={setSearchTerm}
              incrementDownloadCount={(id, fileUrl) => window.open(resolveImageUrl(fileUrl), '_blank')}
              gridColsClassName="grid grid-cols-1 gap-4"
            />
          </div>

          {/* Helpdesk Contact */}
          <div className="bg-slate-50 border border-slate-200 rounded-3xl p-6 flex items-start gap-4">
            <HelpCircle className="h-6 w-6 text-amber-500 shrink-0 mt-0.5" />
            <div className="space-y-1">
              <h4 className="text-xs font-extrabold text-[#002147]">Butuh Bantuan?</h4>
              <p className="text-[11px] text-slate-500 font-medium leading-relaxed">
                Jika mengalami kendala pengisian formulir, silakan datangi loket Desk Layanan PPID di Rektorat Lantai 1 atau kirim email pertanyaan ke ppidcare@upertis.ac.id.
              </p>
            </div>
          </div>
        </div>

        {/* Right Column: Flowchart Diagram & Explanations */}
        <div className="lg:col-span-8 space-y-6">
          <div className="bg-white border border-slate-200 rounded-3xl p-6 lg:p-8 shadow-sm space-y-8">
            <div className="flex flex-col sm:flex-row justify-between sm:items-center gap-4 border-b border-slate-100 pb-4">
              <div>
                <h3 className="text-xs font-bold text-[#002147] uppercase tracking-wider">Bagan Alur Pengajuan Keberatan</h3>
                <span className="text-[10px] text-slate-400 font-medium mt-0.5 block">Klik kotak bagan untuk melihat penjelasan langkah detail</span>
              </div>

              {/* Tab Selector */}
              <div className="flex bg-slate-100 p-1 rounded-xl shrink-0">
                <button
                  onClick={() => handleTabChange('manual')}
                  className={`px-4 py-2 rounded-lg text-[10px] font-extrabold transition-all cursor-pointer ${
                    activeTab === 'manual'
                      ? 'bg-[#002147] text-white shadow-sm'
                      : 'text-slate-500 hover:text-slate-800'
                  }`}
                >
                  Manual / Luring
                </button>
                <button
                  onClick={() => handleTabChange('online')}
                  className={`px-4 py-2 rounded-lg text-[10px] font-extrabold transition-all cursor-pointer ${
                    activeTab === 'online'
                      ? 'bg-[#002147] text-white shadow-sm'
                      : 'text-slate-500 hover:text-slate-800'
                  }`}
                >
                  Online / Daring
                </button>
              </div>
            </div>

            {/* FLOWCHART DIAGRAM WORKSPACE */}
            <SerpentineFlowchart
              steps={steps}
              activeStepIdx={activeStepIdx}
              setActiveStepIdx={setActiveStepIdx}
            />

            {/* STEP DETAIL VIEW PANEL */}
            <div className="bg-slate-50 border border-slate-200 rounded-3xl p-5 lg:p-6 text-left transition-all animate-in fade-in slide-in-from-bottom duration-250">
              <div className="flex items-center gap-2 mb-3 border-b border-slate-200/50 pb-2">
                <span className={`h-6 w-6 rounded-lg flex items-center justify-center border-2 ${currentActiveStep.color} bg-white text-xs`}>
                  {React.createElement(currentActiveStep.icon, { className: "h-3.5 w-3.5" })}
                </span>
                <h4 className="text-xs font-bold text-[#002147] uppercase tracking-wider">
                  Detail Penjelasan: {currentActiveStep.title} (Langkah {currentActiveStep.step})
                </h4>
              </div>

              {/* Explanatory bullets list parsed from DB structure */}
              <div className="space-y-3">
                {currentActiveStep.bullets.map((bullet: string, bIdx: number) => (
                  <div key={bIdx} className="flex gap-2 items-start text-xs lg:text-sm font-medium text-slate-655 leading-relaxed">
                    <span className="h-1.5 w-1.5 rounded-full bg-amber-500 shrink-0 mt-2" />
                    <span>{bullet}</span>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
