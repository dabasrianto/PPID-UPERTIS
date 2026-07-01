import React, { useState } from 'react';
import { Settings } from 'lucide-react';

interface PortalSettingsProps {
  portalName: string;
  setPortalName: (val: string) => void;
  portalDesc: string;
  setPortalDesc: (val: string) => void;
  portalLogo: string;
  setPortalLogo: (val: string) => void;
  portalBgImage: string;
  setPortalBgImage: (val: string) => void;
  portalFontFamily: string;
  setPortalFontFamily: (val: string) => void;
  portalFontSize: string;
  setPortalFontSize: (val: string) => void;
  portalFavicon: string;
  setPortalFavicon: (val: string) => void;
  portalWelcomeText: string;
  setPortalWelcomeText: (val: string) => void;
  portalSkRektor: string;
  setPortalSkRektor: (val: string) => void;
  portalKeberatanLink: string;
  setPortalKeberatanLink: (val: string) => void;
  portalAboutStatNumber: string;
  setPortalAboutStatNumber: (val: string) => void;
  portalAboutStatLabelAccent: string;
  setPortalAboutStatLabelAccent: (val: string) => void;
  portalAboutStatLabel: string;
  setPortalAboutStatLabel: (val: string) => void;
  portalCard1Title: string;
  setPortalCard1Title: (val: string) => void;
  portalCard1Desc: string;
  setPortalCard1Desc: (val: string) => void;
  portalCard2Title: string;
  setPortalCard2Title: (val: string) => void;
  portalCard2Desc: string;
  setPortalCard2Desc: (val: string) => void;
  portalCard3Title: string;
  setPortalCard3Title: (val: string) => void;
  portalCard3Desc: string;
  setPortalCard3Desc: (val: string) => void;
  portalCard1Link: string;
  setPortalCard1Link: (val: string) => void;
  portalCard2Link: string;
  setPortalCard2Link: (val: string) => void;
  portalCard3Link: string;
  setPortalCard3Link: (val: string) => void;
  portalFaqs: { question: string; answer: string }[];
  setPortalFaqs: React.Dispatch<React.SetStateAction<{ question: string; answer: string }[]>>;
  portalPermohonanLink: string;
  setPortalPermohonanLink: (val: string) => void;
  portalPermohonanFormType: string;
  setPortalPermohonanFormType: (val: string) => void;
  portalPengaduanLink: string;
  setPortalPengaduanLink: (val: string) => void;
  portalRektoratEmail: string;
  setPortalRektoratEmail: (val: string) => void;
  portalRektoratPhone: string;
  setPortalRektoratPhone: (val: string) => void;
  portalPlaystoreLink: string;
  setPortalPlaystoreLink: (val: string) => void;
  portalRektoratAddress: string;
  setPortalRektoratAddress: (val: string) => void;
  portalKampus2Address: string;
  setPortalKampus2Address: (val: string) => void;
  portalKampus1MapUrl: string;
  setPortalKampus1MapUrl: (val: string) => void;
  portalKampus2MapUrl: string;
  setPortalKampus2MapUrl: (val: string) => void;
  portalJadwalSeninKamis: string;
  setPortalJadwalSeninKamis: (val: string) => void;
  portalIstirahatSeninKamis: string;
  setPortalIstirahatSeninKamis: (val: string) => void;
  portalJadwalJumat: string;
  setPortalJadwalJumat: (val: string) => void;
  portalIstirahatJumat: string;
  setPortalIstirahatJumat: (val: string) => void;
  portalJadwalSabtuMinggu: string;
  setPortalJadwalSabtuMinggu: (val: string) => void;
  portalIsSaving: boolean;
  handleSaveSettings: (e: React.FormEvent) => void;
  API_BASE_URL: string;
}

export default function PortalSettings({
  portalName,
  setPortalName,
  portalDesc,
  setPortalDesc,
  portalLogo,
  setPortalLogo,
  portalBgImage,
  setPortalBgImage,
  portalFontFamily,
  setPortalFontFamily,
  portalFontSize,
  setPortalFontSize,
  portalFavicon,
  setPortalFavicon,
  portalWelcomeText,
  setPortalWelcomeText,
  portalSkRektor,
  setPortalSkRektor,
  portalKeberatanLink,
  setPortalKeberatanLink,
  portalAboutStatNumber,
  setPortalAboutStatNumber,
  portalAboutStatLabelAccent,
  setPortalAboutStatLabelAccent,
  portalAboutStatLabel,
  setPortalAboutStatLabel,
  portalCard1Title,
  setPortalCard1Title,
  portalCard1Desc,
  setPortalCard1Desc,
  portalCard2Title,
  setPortalCard2Title,
  portalCard2Desc,
  setPortalCard2Desc,
  portalCard3Title,
  setPortalCard3Title,
  portalCard3Desc,
  setPortalCard3Desc,
  portalCard1Link,
  setPortalCard1Link,
  portalCard2Link,
  setPortalCard2Link,
  portalCard3Link,
  setPortalCard3Link,
  portalFaqs,
  setPortalFaqs,
  portalPermohonanLink,
  setPortalPermohonanLink,
  portalPermohonanFormType,
  setPortalPermohonanFormType,
  portalPengaduanLink,
  setPortalPengaduanLink,
  portalRektoratEmail,
  setPortalRektoratEmail,
  portalRektoratPhone,
  setPortalRektoratPhone,
  portalPlaystoreLink,
  setPortalPlaystoreLink,
  portalRektoratAddress,
  setPortalRektoratAddress,
  portalKampus2Address,
  setPortalKampus2Address,
  portalKampus1MapUrl,
  setPortalKampus1MapUrl,
  portalKampus2MapUrl,
  setPortalKampus2MapUrl,
  portalJadwalSeninKamis,
  setPortalJadwalSeninKamis,
  portalIstirahatSeninKamis,
  setPortalIstirahatSeninKamis,
  portalJadwalJumat,
  setPortalJadwalJumat,
  portalIstirahatJumat,
  setPortalIstirahatJumat,
  portalJadwalSabtuMinggu,
  setPortalJadwalSabtuMinggu,
  portalIsSaving,
  handleSaveSettings,
  API_BASE_URL
}: PortalSettingsProps) {
  const [activeSubTab, setActiveSubTab] = useState<'branding' | 'widgets' | 'contact' | 'hours' | 'chatbot'>('branding');

  const tabs = [
    { id: 'branding', label: 'Branding & Utama' },
    { id: 'widgets', label: 'Layanan & Stats' },
    { id: 'contact', label: 'Kontak & Lokasi' },
    { id: 'hours', label: 'Jam Operasional' },
    { id: 'chatbot', label: 'Q&A Chatbot FAQ' }
  ] as const;

  return (
    <div className="bg-white rounded-3xl border border-slate-200 p-6 shadow-sm space-y-6">
      {/* Header and Sub Tab buttons */}
      <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4 border-b border-slate-100 pb-4">
        <div>
          <h2 className="text-base font-extrabold text-[#002147] flex items-center gap-2">
            <Settings className="h-5 w-5 text-amber-500" /> Pengaturan Portal PPID
          </h2>
          <span className="text-[11px] text-slate-400 font-medium block">
            Kelola branding, widget beranda, chatbot, kontak, dan jam operasional dengan mudah.
          </span>
        </div>

        <div className="flex flex-wrap gap-1 bg-slate-50 p-1 rounded-2xl border border-slate-200">
          {tabs.map((tab) => (
            <button
              key={tab.id}
              type="button"
              onClick={() => setActiveSubTab(tab.id)}
              className={`px-3 py-2 text-xs font-bold transition-all rounded-xl cursor-pointer ${
                activeSubTab === tab.id
                  ? 'bg-[#002147] text-white shadow-sm border border-[#002147]'
                  : 'text-slate-500 hover:text-slate-800 bg-transparent border-0'
              }`}
            >
              {tab.label}
            </button>
          ))}
        </div>
      </div>

      <form onSubmit={handleSaveSettings} className="space-y-6 max-w-4xl text-left">
        {/* TAB 1: BRANDING & UTAMA */}
        {activeSubTab === 'branding' && (
          <div className="space-y-5 animate-in fade-in duration-200">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-1">
                <label className="text-xs font-bold text-slate-700 uppercase tracking-wider block">
                  Nama Portal PPID <span className="text-red-500">*</span>
                </label>
                <input
                  type="text"
                  required
                  value={portalName}
                  onChange={(e) => setPortalName(e.target.value)}
                  className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 bg-slate-50 font-bold text-slate-800"
                />
              </div>
              <div className="space-y-1">
                <label className="text-xs font-bold text-slate-700 uppercase tracking-wider block">Deskripsi Portal</label>
                <input
                  type="text"
                  value={portalDesc}
                  onChange={(e) => setPortalDesc(e.target.value)}
                  className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 bg-slate-50 font-medium text-slate-700"
                />
              </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 border-t border-slate-100 pt-4">
              <div className="space-y-1">
                <label className="text-xs font-bold text-slate-700 uppercase tracking-wider block">Jenis Font (Google Font)</label>
                <select
                  value={portalFontFamily}
                  onChange={(e) => setPortalFontFamily(e.target.value)}
                  className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none bg-slate-50 font-bold text-slate-800"
                >
                  <option value="DM Sans">DM Sans (Default)</option>
                  <option value="Inter">Inter (Sangat Bersih)</option>
                  <option value="Plus Jakarta Sans">Plus Jakarta Sans (Modern & Elegan)</option>
                  <option value="Poppins">Poppins (Bulat & Friendly)</option>
                  <option value="Outfit">Outfit (Minimalis & Mewah)</option>
                  <option value="Roboto">Roboto (Sederhana & Rapi)</option>
                  <option value="Montserrat">Montserrat (Tegas & Tebal)</option>
                  <option value="Open Sans">Open Sans (Klasik)</option>
                </select>
              </div>

              <div className="space-y-1">
                <label className="text-xs font-bold text-slate-700 uppercase tracking-wider block">Ukuran Teks Dasar (Base Font Size)</label>
                <select
                  value={portalFontSize}
                  onChange={(e) => setPortalFontSize(e.target.value)}
                  className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none bg-slate-50 font-bold text-slate-800"
                >
                  <option value="normal">Normal (Kecil / Default)</option>
                  <option value="medium">Sedang (Lebih Nyaman Dibaca)</option>
                  <option value="large">Besar (Sangat Jelas / Ramah Pengguna)</option>
                </select>
              </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-4 border-t border-slate-100 pt-4">
              {/* Logo upload */}
              <div className="space-y-1">
                <label className="text-xs font-bold text-slate-700 uppercase tracking-wider block">Logo Utama Portal</label>
                <div className="flex gap-2">
                  <input
                    type="text"
                    placeholder="/uploads/..."
                    value={portalLogo}
                    onChange={(e) => setPortalLogo(e.target.value)}
                    className="flex-1 rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none bg-slate-50 font-mono text-slate-705"
                  />
                  <label className="px-3 py-2 bg-slate-200 hover:bg-slate-350 rounded-xl text-xs font-bold transition-all cursor-pointer inline-flex items-center text-slate-700 border border-slate-300">
                    Upload
                    <input
                      type="file"
                      onChange={async (e) => {
                        const file = e.target.files?.[0];
                        if (!file) return;
                        const formData = new FormData();
                        formData.append('files', file);
                        const t = localStorage.getItem('auth_token');
                        const uploadRes = await fetch(`${API_BASE_URL}/admin/uploads`, {
                          method: 'POST',
                          headers: { Authorization: `Bearer ${t}` },
                          body: formData
                        });
                        const resData = await uploadRes.json();
                        const uploadedUrl = resData.url || resData.path || (resData.urls && resData.urls[0]) || (resData.uploaded && resData.uploaded[0] && resData.uploaded[0].url) || '';
                        if (uploadRes.ok && uploadedUrl) setPortalLogo(uploadedUrl);
                      }}
                      className="hidden"
                    />
                  </label>
                </div>
              </div>

              {/* BG Image upload */}
              <div className="space-y-1">
                <label className="text-xs font-bold text-slate-700 uppercase tracking-wider block">Gambar Latar Belakang</label>
                <div className="flex gap-2">
                  <input
                    type="text"
                    placeholder="/uploads/..."
                    value={portalBgImage}
                    onChange={(e) => setPortalBgImage(e.target.value)}
                    className="flex-1 rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none bg-slate-50 font-mono text-slate-705"
                  />
                  <label className="px-3 py-2 bg-slate-200 hover:bg-slate-350 rounded-xl text-xs font-bold transition-all cursor-pointer inline-flex items-center text-slate-700 border border-slate-300">
                    Upload
                    <input
                      type="file"
                      onChange={async (e) => {
                        const file = e.target.files?.[0];
                        if (!file) return;
                        const formData = new FormData();
                        formData.append('files', file);
                        const t = localStorage.getItem('auth_token');
                        const uploadRes = await fetch(`${API_BASE_URL}/admin/uploads`, {
                          method: 'POST',
                          headers: { Authorization: `Bearer ${t}` },
                          body: formData
                        });
                        const resData = await uploadRes.json();
                        const uploadedUrl = resData.url || resData.path || (resData.urls && resData.urls[0]) || (resData.uploaded && resData.uploaded[0] && resData.uploaded[0].url) || '';
                        if (uploadRes.ok && uploadedUrl) setPortalBgImage(uploadedUrl);
                      }}
                      className="hidden"
                    />
                  </label>
                </div>
              </div>

              {/* Favicon URL */}
              <div className="space-y-1">
                <label className="text-xs font-bold text-slate-700 uppercase tracking-wider block">Link Favicon URL</label>
                <div className="flex items-center gap-2">
                  <input
                    type="text"
                    placeholder="/uploads/favicon.png"
                    value={portalFavicon}
                    onChange={(e) => setPortalFavicon(e.target.value)}
                    className="flex-1 rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none bg-slate-50 font-mono text-slate-700"
                  />
                  <label className="px-3 py-2 bg-slate-200 hover:bg-slate-355 rounded-xl text-xs font-bold transition-all cursor-pointer inline-flex items-center text-slate-700 border border-slate-300">
                    Upload
                    <input
                      type="file"
                      onChange={async (e) => {
                        const file = e.target.files?.[0];
                        if (!file) return;
                        const formData = new FormData();
                        formData.append('files', file);
                        const t = localStorage.getItem('auth_token');
                        const uploadRes = await fetch(`${API_BASE_URL}/admin/uploads`, {
                          method: 'POST',
                          headers: { Authorization: `Bearer ${t}` },
                          body: formData
                        });
                        const resData = await uploadRes.json();
                        const uploadedUrl = resData.url || resData.path || (resData.urls && resData.urls[0]) || (resData.uploaded && resData.uploaded[0] && resData.uploaded[0].url) || '';
                        if (uploadRes.ok && uploadedUrl) setPortalFavicon(uploadedUrl);
                      }}
                      className="hidden"
                    />
                  </label>
                </div>
              </div>
            </div>

            <div className="space-y-3 border-t border-slate-100 pt-4">
              <div className="space-y-1">
                <label className="text-xs font-bold text-slate-700 uppercase tracking-wider block">Teks Sambutan Selamat Datang</label>
                <textarea
                  value={portalWelcomeText}
                  onChange={(e) => setPortalWelcomeText(e.target.value)}
                  rows={3}
                  className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none bg-slate-50 font-medium text-slate-700 resize-none"
                />
              </div>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="space-y-1">
                  <label className="text-xs font-bold text-slate-700 uppercase tracking-wider block">SK Rektor tentang PPID</label>
                  <input
                    type="text"
                    value={portalSkRektor}
                    onChange={(e) => setPortalSkRektor(e.target.value)}
                    className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none bg-slate-50 font-medium text-slate-700"
                  />
                </div>
                <div className="space-y-1">
                  <label className="text-xs font-bold text-slate-700 uppercase tracking-wider block">Tipe Formulir Permohonan</label>
                  <select
                    value={portalPermohonanFormType}
                    onChange={(e) => setPortalPermohonanFormType(e.target.value)}
                    className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none bg-slate-50 font-bold text-slate-800"
                  >
                    <option value="internal">Formulir Internal PPID (Built-in)</option>
                    <option value="external">Formulir Eksternal (Google Form / Tautan Luar)</option>
                  </select>
                </div>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-3 gap-4 border-t border-slate-100 pt-4">
                <div className="space-y-1">
                  <label className="text-xs font-bold text-slate-700 uppercase tracking-wider block">Link Formulir Permohonan Eksternal</label>
                  <input
                    type="text"
                    value={portalPermohonanLink}
                    onChange={(e) => setPortalPermohonanLink(e.target.value)}
                    placeholder="https://forms.gle/..."
                    disabled={portalPermohonanFormType === 'internal'}
                    className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none bg-slate-50 font-mono text-[11px] disabled:opacity-50"
                  />
                </div>
                <div className="space-y-1">
                  <label className="text-xs font-bold text-slate-700 uppercase tracking-wider block">Link Pengajuan Keberatan (Tautan)</label>
                  <input
                    type="text"
                    value={portalKeberatanLink}
                    onChange={(e) => setPortalKeberatanLink(e.target.value)}
                    placeholder="https://forms.gle/..."
                    className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none bg-slate-50 font-mono text-[11px]"
                  />
                </div>
                <div className="space-y-1">
                  <label className="text-xs font-bold text-slate-700 uppercase tracking-wider block">Link Pengaduan Layanan (Tautan)</label>
                  <input
                    type="text"
                    value={portalPengaduanLink}
                    onChange={(e) => setPortalPengaduanLink(e.target.value)}
                    placeholder="https://lapor.go.id/..."
                    className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none bg-slate-50 font-mono text-[11px]"
                  />
                </div>
              </div>
            </div>
          </div>
        )}

        {/* TAB 2: LAYANAN & STATS */}
        {activeSubTab === 'widgets' && (
          <div className="space-y-6 animate-in fade-in duration-200">
            {/* Kartu Statistik */}
            <div className="space-y-3">
              <div>
                <h4 className="text-xs font-extrabold text-[#002147] uppercase tracking-wider">Kartu Statistik (Halaman Utama)</h4>
                <p className="text-[11px] text-slate-400 font-medium mt-0.5">Kartu melayang yang tampil di samping foto pada bagian Tentang/About di halaman utama.</p>
              </div>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4 bg-slate-50/50 p-4 border border-slate-150 rounded-2xl">
                <div className="space-y-1">
                  <label className="text-[10px] font-bold text-slate-700 uppercase tracking-wider block">Angka / Nilai Utama</label>
                  <input
                    type="text"
                    value={portalAboutStatNumber}
                    onChange={(e) => setPortalAboutStatNumber(e.target.value)}
                    placeholder="2021"
                    className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none bg-white font-black text-slate-800"
                  />
                </div>
                <div className="space-y-1">
                  <label className="text-[10px] font-bold text-slate-700 uppercase tracking-wider block">Label Aksen (Kecil, Atas)</label>
                  <input
                    type="text"
                    value={portalAboutStatLabelAccent}
                    onChange={(e) => setPortalAboutStatLabelAccent(e.target.value)}
                    placeholder="Tahun Berdiri"
                    className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none bg-white font-medium text-slate-800"
                  />
                </div>
                <div className="space-y-1">
                  <label className="text-[10px] font-bold text-slate-700 uppercase tracking-wider block">Teks Bawah (Deskripsi Singkat)</label>
                  <input
                    type="text"
                    value={portalAboutStatLabel}
                    onChange={(e) => setPortalAboutStatLabel(e.target.value)}
                    placeholder="PPID UPERTIS Melayani"
                    className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none bg-white font-medium text-slate-800"
                  />
                </div>
              </div>
            </div>

            {/* Kartu Layanan */}
            <div className="space-y-4 border-t border-slate-100 pt-4">
              <div>
                <h4 className="text-xs font-extrabold text-[#002147] uppercase tracking-wider">Kartu Layanan Utama (Section Selamat Datang)</h4>
                <p className="text-[11px] text-slate-400 font-medium mt-0.5">Ubah judul, deskripsi, dan tujuan link tiga kartu menu utama di beranda.</p>
              </div>

              {/* Card 1 */}
              <div className="border border-slate-250/60 rounded-2xl p-4 bg-slate-50/40 space-y-3">
                <span className="text-[10px] font-extrabold text-[#002147] uppercase tracking-wider block">1. Kartu Permohonan Informasi</span>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                  <div className="space-y-1">
                    <label className="text-[9px] font-bold text-slate-500 uppercase block">Judul</label>
                    <input
                      type="text"
                      value={portalCard1Title}
                      onChange={(e) => setPortalCard1Title(e.target.value)}
                      className="w-full rounded-xl border border-slate-200 px-3 py-2 text-xs focus:outline-none bg-white font-bold"
                    />
                  </div>
                  <div className="space-y-1">
                    <label className="text-[9px] font-bold text-slate-500 uppercase block">Link / Tujuan Tombol</label>
                    <input
                      type="text"
                      value={portalCard1Link}
                      onChange={(e) => setPortalCard1Link(e.target.value)}
                      className="w-full rounded-xl border border-slate-200 px-3 py-2 text-xs focus:outline-none bg-white font-mono text-[11px]"
                    />
                  </div>
                </div>
                <div className="space-y-1">
                  <label className="text-[9px] font-bold text-slate-500 uppercase block">Deskripsi Kartu</label>
                  <textarea
                    value={portalCard1Desc}
                    onChange={(e) => setPortalCard1Desc(e.target.value)}
                    rows={2}
                    className="w-full rounded-xl border border-slate-200 px-3 py-2 text-xs focus:outline-none bg-white resize-none"
                  />
                </div>
              </div>

              {/* Card 2 */}
              <div className="border border-slate-250/60 rounded-2xl p-4 bg-slate-50/40 space-y-3">
                <span className="text-[10px] font-extrabold text-[#002147] uppercase tracking-wider block">2. Kartu Keberatan Informasi</span>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                  <div className="space-y-1">
                    <label className="text-[9px] font-bold text-slate-500 uppercase block">Judul</label>
                    <input
                      type="text"
                      value={portalCard2Title}
                      onChange={(e) => setPortalCard2Title(e.target.value)}
                      className="w-full rounded-xl border border-slate-200 px-3 py-2 text-xs focus:outline-none bg-white font-bold"
                    />
                  </div>
                  <div className="space-y-1">
                    <label className="text-[9px] font-bold text-slate-500 uppercase block">Link / Tujuan Tombol</label>
                    <input
                      type="text"
                      value={portalCard2Link}
                      onChange={(e) => setPortalCard2Link(e.target.value)}
                      className="w-full rounded-xl border border-slate-200 px-3 py-2 text-xs focus:outline-none bg-white font-mono text-[11px]"
                    />
                  </div>
                </div>
                <div className="space-y-1">
                  <label className="text-[9px] font-bold text-slate-500 uppercase block">Deskripsi Kartu</label>
                  <textarea
                    value={portalCard2Desc}
                    onChange={(e) => setPortalCard2Desc(e.target.value)}
                    rows={2}
                    className="w-full rounded-xl border border-slate-200 px-3 py-2 text-xs focus:outline-none bg-white resize-none"
                  />
                </div>
              </div>

              {/* Card 3 */}
              <div className="border border-slate-250/60 rounded-2xl p-4 bg-slate-50/40 space-y-3">
                <span className="text-[10px] font-extrabold text-[#002147] uppercase tracking-wider block">3. Kartu Pengaduan Layanan</span>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                  <div className="space-y-1">
                    <label className="text-[9px] font-bold text-slate-500 uppercase block">Judul</label>
                    <input
                      type="text"
                      value={portalCard3Title}
                      onChange={(e) => setPortalCard3Title(e.target.value)}
                      className="w-full rounded-xl border border-slate-200 px-3 py-2 text-xs focus:outline-none bg-white font-bold"
                    />
                  </div>
                  <div className="space-y-1">
                    <label className="text-[9px] font-bold text-slate-500 uppercase block">Link / Tujuan Tombol</label>
                    <input
                      type="text"
                      value={portalCard3Link}
                      onChange={(e) => setPortalCard3Link(e.target.value)}
                      className="w-full rounded-xl border border-slate-200 px-3 py-2 text-xs focus:outline-none bg-white font-mono text-[11px]"
                    />
                  </div>
                </div>
                <div className="space-y-1">
                  <label className="text-[9px] font-bold text-slate-500 uppercase block">Deskripsi Kartu</label>
                  <textarea
                    value={portalCard3Desc}
                    onChange={(e) => setPortalCard3Desc(e.target.value)}
                    rows={2}
                    className="w-full rounded-xl border border-slate-200 px-3 py-2 text-xs focus:outline-none bg-white resize-none"
                  />
                </div>
              </div>
            </div>
          </div>
        )}

        {/* TAB 3: KONTAK & LOKASI */}
        {activeSubTab === 'contact' && (
          <div className="space-y-5 animate-in fade-in duration-200">
            <div>
              <h4 className="text-xs font-extrabold text-[#002147] uppercase tracking-wider">Kontak & Lokasi Kampus</h4>
              <p className="text-[11px] text-slate-400 font-medium mt-0.5">Atur email, desk whatsapp, alamat fisik, dan link google maps kampus.</p>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-1">
                <label className="text-[10px] font-bold text-slate-700 uppercase tracking-wider block">Email Korespondensi Resmi</label>
                <input
                  type="text"
                  placeholder="ppid@upertis.ac.id"
                  value={portalRektoratEmail}
                  onChange={(e) => setPortalRektoratEmail(e.target.value)}
                  className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none bg-slate-50 font-medium"
                />
              </div>
              <div className="space-y-1">
                <label className="text-[10px] font-bold text-slate-700 uppercase tracking-wider block">WhatsApp Desk / Hotline</label>
                <input
                  type="text"
                  placeholder="+62 821-xxxx-xxxx"
                  value={portalRektoratPhone}
                  onChange={(e) => setPortalRektoratPhone(e.target.value)}
                  className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none bg-slate-50 font-medium"
                />
              </div>
            </div>

            <div className="space-y-1">
              <label className="text-[10px] font-bold text-slate-700 uppercase tracking-wider block">Link Aplikasi Play Store (Google Play)</label>
              <input
                type="text"
                placeholder="https://play.google.com/store/apps/details?id=..."
                value={portalPlaystoreLink}
                onChange={(e) => setPortalPlaystoreLink(e.target.value)}
                className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none bg-slate-50 font-medium text-slate-800 font-mono text-[11px]"
              />
            </div>

            <div className="space-y-3">
              <div className="space-y-1">
                <label className="text-[10px] font-bold text-slate-700 uppercase tracking-wider block">Alamat Kantor Pelayanan (Kampus Utama Padang)</label>
                <input
                  type="text"
                  placeholder="Lobby Gedung Rektorat Kampus UPERTIS..."
                  value={portalRektoratAddress}
                  onChange={(e) => setPortalRektoratAddress(e.target.value)}
                  className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none bg-slate-50 font-medium text-slate-800"
                />
              </div>

              <div className="space-y-1">
                <label className="text-[10px] font-bold text-slate-700 uppercase tracking-wider block">Alamat Kampus II Bukittinggi</label>
                <input
                  type="text"
                  placeholder="Jl. Raya Bukittinggi - Padang Luar..."
                  value={portalKampus2Address}
                  onChange={(e) => setPortalKampus2Address(e.target.value)}
                  className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none bg-slate-50 font-medium text-slate-800"
                />
              </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 border-t border-slate-100 pt-4">
              <div className="space-y-1">
                <label className="text-[10px] font-bold text-slate-700 uppercase tracking-wider block">Google Maps Embed URL - Kampus Utama</label>
                <input
                  type="text"
                  placeholder="https://maps.google.com/maps?q=..."
                  value={portalKampus1MapUrl}
                  onChange={(e) => setPortalKampus1MapUrl(e.target.value)}
                  className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none bg-slate-50 font-mono text-[10px]"
                />
              </div>
              <div className="space-y-1">
                <label className="text-[10px] font-bold text-slate-700 uppercase tracking-wider block">Google Maps Embed URL - Kampus II</label>
                <input
                  type="text"
                  placeholder="https://maps.google.com/maps?q=..."
                  value={portalKampus2MapUrl}
                  onChange={(e) => setPortalKampus2MapUrl(e.target.value)}
                  className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none bg-slate-50 font-mono text-[10px]"
                />
              </div>
            </div>
          </div>
        )}

        {/* TAB 4: JAM OPERASIONAL */}
        {activeSubTab === 'hours' && (
          <div className="space-y-5 animate-in fade-in duration-200">
            <div>
              <h4 className="text-xs font-extrabold text-[#002147] uppercase tracking-wider">Jadwal Pelayanan Informasi (Jadwal Layanan)</h4>
              <p className="text-[11px] text-slate-400 font-medium mt-0.5">Kelola jam kerja operasional desk pelayanan fisik loket PPID.</p>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-1">
                <label className="text-[10px] font-bold text-slate-700 uppercase tracking-wider block">Jam Kerja (Senin - Kamis)</label>
                <input
                  type="text"
                  placeholder="08:00 – 16:00 WIB"
                  value={portalJadwalSeninKamis}
                  onChange={(e) => setPortalJadwalSeninKamis(e.target.value)}
                  className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none bg-slate-50 font-medium"
                />
              </div>
              <div className="space-y-1">
                <label className="text-[10px] font-bold text-slate-700 uppercase tracking-wider block">Jam Istirahat (Senin - Kamis)</label>
                <input
                  type="text"
                  placeholder="12:00 – 13:30 WIB"
                  value={portalIstirahatSeninKamis}
                  onChange={(e) => setPortalIstirahatSeninKamis(e.target.value)}
                  className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none bg-slate-50 font-medium"
                />
              </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 border-t border-slate-100 pt-4">
              <div className="space-y-1">
                <label className="text-[10px] font-bold text-slate-700 uppercase tracking-wider block">Jam Kerja (Hari Jumat)</label>
                <input
                  type="text"
                  placeholder="08:00 – 16:30 WIB"
                  value={portalJadwalJumat}
                  onChange={(e) => setPortalJadwalJumat(e.target.value)}
                  className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none bg-slate-50 font-medium"
                />
              </div>
              <div className="space-y-1">
                <label className="text-[10px] font-bold text-slate-700 uppercase tracking-wider block">Jam Istirahat (Hari Jumat)</label>
                <input
                  type="text"
                  placeholder="12:00 – 14:00 WIB"
                  value={portalIstirahatJumat}
                  onChange={(e) => setPortalIstirahatJumat(e.target.value)}
                  className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none bg-slate-50 font-medium"
                />
              </div>
            </div>

            <div className="space-y-1 border-t border-slate-100 pt-4">
              <label className="text-[10px] font-bold text-slate-700 uppercase tracking-wider block">Status Sabtu - Minggu & Hari Libur</label>
              <input
                type="text"
                placeholder="Sistem Online Tetap Aktif 24/7..."
                value={portalJadwalSabtuMinggu}
                onChange={(e) => setPortalJadwalSabtuMinggu(e.target.value)}
                className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none bg-slate-50 font-medium"
              />
            </div>
          </div>
        )}

        {/* TAB 5: CHATBOT FAQ */}
        {activeSubTab === 'chatbot' && (
          <div className="space-y-5 animate-in fade-in duration-200">
            <div>
              <h4 className="text-xs font-extrabold text-[#002147] uppercase tracking-wider">Pengaturan Q&A Chatbot FAQ</h4>
              <p className="text-[11px] text-slate-400 font-medium mt-0.5">Kelola daftar pertanyaan dan jawaban otomatis yang muncul pada asisten virtual FAQ di beranda.</p>
            </div>

            <div className="space-y-3">
              {portalFaqs.map((faq, idx) => (
                <div key={idx} className="border border-slate-200 rounded-2xl p-4 bg-slate-50/50 space-y-3 relative">
                  <div className="flex justify-between items-center">
                    <span className="text-[10px] font-black text-amber-600 uppercase tracking-wider">Pertanyaan #{idx + 1}</span>
                    <button
                      type="button"
                      onClick={() => {
                        const updated = portalFaqs.filter((_, i) => i !== idx);
                        setPortalFaqs(updated);
                      }}
                      className="text-[10px] font-extrabold text-red-500 hover:text-red-700 transition-colors uppercase tracking-wider bg-transparent border-0 cursor-pointer"
                    >
                      Hapus
                    </button>
                  </div>
                  <div className="space-y-2">
                    <div className="space-y-1">
                      <label className="text-[9px] font-bold text-slate-500 uppercase block">Pertanyaan</label>
                      <input
                        type="text"
                        value={faq.question}
                        onChange={(e) => {
                          const updated = [...portalFaqs];
                          updated[idx].question = e.target.value;
                          setPortalFaqs(updated);
                        }}
                        placeholder="Masukkan pertanyaan..."
                        className="w-full rounded-xl border border-slate-200 px-3 py-2 text-xs focus:outline-none bg-white font-bold"
                      />
                    </div>
                    <div className="space-y-1">
                      <label className="text-[9px] font-bold text-slate-500 uppercase block">Jawaban</label>
                      <textarea
                        value={faq.answer}
                        onChange={(e) => {
                          const updated = [...portalFaqs];
                          updated[idx].answer = e.target.value;
                          setPortalFaqs(updated);
                        }}
                        placeholder="Masukkan jawaban..."
                        rows={2}
                        className="w-full rounded-xl border border-slate-200 px-3 py-2 text-xs focus:outline-none bg-white resize-none"
                      />
                    </div>
                  </div>
                </div>
              ))}

              {portalFaqs.length === 0 && (
                <div className="py-8 text-center text-slate-400 border border-dashed border-slate-200 rounded-2xl text-xs font-semibold">
                  Belum ada pertanyaan terprogram. Silakan tambah di bawah.
                </div>
              )}
            </div>

            <div className="pt-2">
              <button
                type="button"
                onClick={() => {
                  setPortalFaqs(prev => [...prev, { question: '', answer: '' }]);
                }}
                className="inline-flex items-center gap-1.5 px-4 py-2.5 bg-amber-500 hover:bg-amber-600 text-[#002147] text-xs font-extrabold uppercase tracking-wider rounded-xl transition-all border-0 cursor-pointer"
              >
                + Tambah Pertanyaan Baru
              </button>
            </div>
          </div>
        )}

        {/* Global Save Button */}
        <div className="flex justify-end gap-3 pt-5 border-t border-slate-100">
          <button
            type="submit"
            disabled={portalIsSaving}
            className="px-8 py-3 bg-[#002147] hover:bg-amber-500 hover:text-[#002147] text-white text-xs font-bold uppercase transition-all rounded-xl shadow-md cursor-pointer border border-[#002147]/50"
          >
            {portalIsSaving ? 'Menyimpan...' : 'Simpan Pengaturan Portal'}
          </button>
        </div>
      </form>
    </div>
  );
}
