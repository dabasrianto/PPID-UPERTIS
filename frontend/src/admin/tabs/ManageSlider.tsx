import React from 'react';
import {
  Layers, ChevronUp, ChevronDown, Trash2, Upload, Save, FileText, Search, Plus
} from 'lucide-react';
import { resolveImageUrl } from '../../utils/helpers';

interface ManageSliderProps {
  heroImages: any[];
  siteConfig: any;
  setSiteConfig: React.Dispatch<React.SetStateAction<any>>;
  setAdminGlobalMessage: (msg: string) => void;
  API_BASE_URL: string;
}

export default function ManageSlider({
  heroImages,
  siteConfig,
  setSiteConfig,
  setAdminGlobalMessage,
  API_BASE_URL
}: ManageSliderProps) {
  return (
    <div className="bg-white rounded-3xl border border-slate-200 p-6 shadow-sm space-y-6">
      <div>
        <h2 className="text-base font-extrabold text-[#002147] flex items-center gap-2">
          <Layers className="h-5 w-5 text-amber-500" /> Slider Beranda
        </h2>
        <span className="text-[11px] text-slate-400 font-medium block">
          Kelola gambar, judul, deskripsi, dan tombol navigasi secara individual untuk masing-masing slide beranda.
        </span>
      </div>

      {/* Current Slider List */}
      <div className="space-y-6">
        <h3 className="text-xs font-bold text-slate-600 uppercase tracking-wider">Daftar Slide Aktif ({heroImages.length} Slide)</h3>
        {heroImages.length === 0 ? (
          <div className="bg-slate-50 border border-dashed border-slate-300 rounded-2xl p-8 text-center">
            <Layers className="h-8 w-8 text-slate-300 mx-auto mb-2" />
            <p className="text-xs text-slate-400">Belum ada slide. Silakan unggah gambar di bawah.</p>
          </div>
        ) : (
          <div className="space-y-4">
            {heroImages.map((slide: any, idx: number) => {
              const isString = typeof slide === 'string';
              const imageUrl = isString ? slide : slide?.image || '';
              const slideTitle = isString ? '' : slide?.title || '';
              const slideSubtitle = isString ? '' : slide?.subtitle || '';
              const slideDesc = isString ? '' : slide?.description || '';
              const slideBtn1Text = isString ? '' : slide?.btn1_text || '';
              const slideBtn1Page = isString ? '' : slide?.btn1_page || '';
              const slideBtn2Text = isString ? '' : slide?.btn2_text || '';
              const slideBtn2Page = isString ? '' : slide?.btn2_page || '';

              const updateField = (field: string, val: string) => {
                const newImages = [...heroImages];
                newImages[idx] = isString
                  ? { image: imageUrl, [field]: val }
                  : { ...slide, [field]: val };
                setSiteConfig((prev: any) => ({
                  ...prev,
                  settings: { ...prev?.settings, hero_images: newImages }
                }));
              };

              return (
                <div key={idx} className="bg-slate-50 border border-slate-200 rounded-2xl p-5 relative space-y-4 text-left">
                  <div className="flex flex-col lg:flex-row gap-5">
                    {/* Left: Thumbnail & Actions */}
                    <div className="w-full lg:w-44 shrink-0 space-y-2">
                      <div className="relative rounded-xl overflow-hidden border border-slate-200 shadow-sm h-28 bg-white">
                        <img
                          src={resolveImageUrl(imageUrl)}
                          alt={`Slide ${idx + 1}`}
                          className="w-full h-full object-cover"
                          onError={(e) => { e.currentTarget.src = '/logo.png'; }}
                        />
                        <div className="absolute top-2 left-2 bg-[#002147]/80 text-white text-[9px] font-bold px-2 py-0.5 rounded-full">
                          Slide {idx + 1}
                        </div>
                      </div>
                      <div className="flex items-center gap-1.5 justify-center">
                        <button
                          type="button"
                          disabled={idx === 0}
                          onClick={() => {
                            if (idx === 0) return;
                            const newImages = [...heroImages];
                            const temp = newImages[idx];
                            newImages[idx] = newImages[idx - 1];
                            newImages[idx - 1] = temp;
                            setSiteConfig((prev: any) => ({
                              ...prev,
                              settings: { ...prev?.settings, hero_images: newImages }
                            }));
                          }}
                          className="p-1.5 bg-white border border-slate-200 text-slate-500 hover:text-[#002147] hover:bg-slate-100 rounded-lg text-xs disabled:opacity-40 disabled:hover:bg-white cursor-pointer"
                          title="Naikkan posisi"
                        >
                          <ChevronUp className="h-4 w-4" />
                        </button>
                        <button
                          type="button"
                          disabled={idx === heroImages.length - 1}
                          onClick={() => {
                            if (idx === heroImages.length - 1) return;
                            const newImages = [...heroImages];
                            const temp = newImages[idx];
                            newImages[idx] = newImages[idx + 1];
                            newImages[idx + 1] = temp;
                            setSiteConfig((prev: any) => ({
                              ...prev,
                              settings: { ...prev?.settings, hero_images: newImages }
                            }));
                          }}
                          className="p-1.5 bg-white border border-slate-200 text-slate-500 hover:text-[#002147] hover:bg-slate-100 rounded-lg text-xs disabled:opacity-40 disabled:hover:bg-white cursor-pointer"
                          title="Turunkan posisi"
                        >
                          <ChevronDown className="h-4 w-4" />
                        </button>
                        <button
                          type="button"
                          onClick={async () => {
                            const newImages = heroImages.filter((_: any, i: number) => i !== idx);
                            const token = localStorage.getItem('auth_token');
                            if (!token) return;
                            const siteId = siteConfig?.id || 'ppid';
                            const updatedConfig = {
                              ...siteConfig,
                              settings: {
                                ...siteConfig?.settings,
                                hero_images: newImages,
                                hero_image: newImages[0] || '',
                              },
                            };
                            try {
                              const res = await fetch(`${API_BASE_URL}/admin/sites/${siteId}`, {
                                method: 'PUT',
                                headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
                                body: JSON.stringify(updatedConfig),
                              });
                              if (res.ok) {
                                setSiteConfig(updatedConfig);
                                setAdminGlobalMessage('Slide berhasil dihapus!');
                              }
                            } catch (err) {
                              console.error(err);
                            }
                          }}
                          className="p-1.5 bg-red-50 hover:bg-red-100 text-red-500 rounded-lg text-xs cursor-pointer border border-red-200"
                          title="Hapus slide"
                        >
                          <Trash2 className="h-4 w-4" />
                        </button>
                      </div>
                    </div>

                    {/* Right: Content Forms */}
                    <div className="flex-1 grid grid-cols-1 md:grid-cols-2 gap-3 text-left font-sans">
                      <div className="space-y-1.5 md:col-span-2">
                        <label className="text-[10px] font-bold text-slate-500 uppercase tracking-wider block">Judul Slide (Hero Title)</label>
                        <input
                          type="text"
                          value={slideTitle}
                          onChange={(e) => updateField('title', e.target.value)}
                          className="w-full text-xs border border-slate-300 bg-white rounded-xl px-3 py-2 focus:ring-2 focus:ring-amber-400 focus:outline-none font-bold"
                          placeholder="Contoh: Portal Keterbukaan Informasi Publik"
                        />
                      </div>

                      <div className="space-y-1.5">
                        <label className="text-[10px] font-bold text-slate-500 uppercase tracking-wider block">Sub-judul (Badge Kuning)</label>
                        <input
                          type="text"
                          value={slideSubtitle}
                          onChange={(e) => updateField('subtitle', e.target.value)}
                          className="w-full text-xs border border-slate-300 bg-white rounded-xl px-3 py-2 focus:ring-2 focus:ring-amber-400 focus:outline-none"
                          placeholder="Contoh: Portal Transparansi Publik"
                        />
                      </div>

                      <div className="space-y-1.5">
                        <label className="text-[10px] font-bold text-slate-500 uppercase tracking-wider block">Deskripsi Slide</label>
                        <textarea
                          value={slideDesc}
                          onChange={(e) => updateField('description', e.target.value)}
                          className="w-full text-xs border border-slate-300 bg-white rounded-xl px-3 py-1.5 focus:ring-2 focus:ring-amber-400 focus:outline-none h-11 resize-none font-medium text-slate-805"
                          placeholder="Contoh: Akses dokumen dan data resmi kampus dengan mudah..."
                        />
                      </div>

                      {/* Buttons Config inside slide */}
                      <div className="bg-white border border-slate-200 rounded-xl p-3 space-y-2">
                        <span className="text-[9px] font-extrabold text-[#002147] uppercase block text-left">Tombol Utama (Kiri)</span>
                        <div className="grid grid-cols-2 gap-2">
                          <input
                            type="text"
                            value={slideBtn1Text}
                            onChange={(e) => updateField('btn1_text', e.target.value)}
                            className="w-full text-[11px] border border-slate-300 rounded-lg px-2.5 py-1.5 focus:outline-none"
                            placeholder="Label Tombol"
                          />
                          <input
                            type="text"
                            value={slideBtn1Page}
                            onChange={(e) => updateField('btn1_page', e.target.value)}
                            className="w-full text-[11px] border border-slate-300 rounded-lg px-2.5 py-1.5 focus:outline-none font-mono"
                            placeholder="Tujuan (Page)"
                          />
                        </div>
                      </div>

                      <div className="bg-white border border-slate-200 rounded-xl p-3 space-y-2">
                        <span className="text-[9px] font-extrabold text-[#002147] uppercase block text-left">Tombol Sekunder (Kanan)</span>
                        <div className="grid grid-cols-2 gap-2">
                          <input
                            type="text"
                            value={slideBtn2Text}
                            onChange={(e) => updateField('btn2_text', e.target.value)}
                            className="w-full text-[11px] border border-slate-300 rounded-lg px-2.5 py-1.5 focus:outline-none"
                            placeholder="Label Tombol"
                          />
                          <input
                            type="text"
                            value={slideBtn2Page}
                            onChange={(e) => updateField('btn2_page', e.target.value)}
                            className="w-full text-[11px] border border-slate-300 rounded-lg px-2.5 py-1.5 focus:outline-none font-mono"
                            placeholder="Tujuan (Page)"
                          />
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </div>

      {/* Upload New Slider Image */}
      <div className="border-t border-slate-100 pt-4 space-y-3">
        <h3 className="text-xs font-bold text-slate-600 uppercase tracking-wider text-left">Tambah Slide Baru</h3>
        <div className="flex items-center gap-3">
          <label className="flex-1 flex items-center gap-2 px-4 py-3 bg-slate-50 border border-dashed border-slate-300 rounded-xl cursor-pointer hover:bg-blue-50 hover:border-blue-300 transition-all">
            <Upload className="h-4 w-4 text-slate-400" />
            <span className="text-xs text-slate-500 font-medium">Unggah Gambar untuk Slide Baru...</span>
            <input
              type="file"
              accept="image/*"
              className="hidden"
              onChange={async (e) => {
                const file = e.target.files?.[0];
                if (!file) return;
                const token = localStorage.getItem('auth_token');
                if (!token) return;

                setAdminGlobalMessage('Mengunggah gambar slide baru...');

                const formData = new FormData();
                formData.append('file', file);
                formData.append('type', 'sites');

                try {
                  const uploadRes = await fetch(`${API_BASE_URL}/admin/uploads`, {
                    method: 'POST',
                    headers: { Authorization: `Bearer ${token}` },
                    body: formData,
                  });

                  if (!uploadRes.ok) {
                    setAdminGlobalMessage('Gagal mengunggah gambar!');
                    return;
                  }

                  const uploadData = await uploadRes.json();
                  const newImageUrl = uploadData.url || uploadData.path || (uploadData.urls && uploadData.urls[0]) || (uploadData.uploaded && uploadData.uploaded[0] && uploadData.uploaded[0].url) || '';

                  const newSlideObj = {
                    image: newImageUrl,
                    title: '',
                    subtitle: '',
                    description: '',
                    btn1_text: '',
                    btn1_page: '',
                    btn2_text: '',
                    btn2_page: '',
                  };

                  const newImages = [...heroImages, newSlideObj];
                  const siteId = siteConfig?.id || 'ppid';
                  const updatedConfig = {
                    ...siteConfig,
                    settings: {
                      ...siteConfig?.settings,
                      hero_images: newImages,
                      hero_image: typeof newImages[0] === 'string' ? newImages[0] : newImages[0]?.image || '',
                    },
                  };

                  const saveRes = await fetch(`${API_BASE_URL}/admin/sites/${siteId}`, {
                    method: 'PUT',
                    headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
                    body: JSON.stringify(updatedConfig),
                  });

                  if (saveRes.ok) {
                    setSiteConfig(updatedConfig);
                    setAdminGlobalMessage('Slide baru berhasil ditambahkan!');
                  } else {
                    setAdminGlobalMessage('Gagal menyimpan slide baru.');
                  }
                } catch (err) {
                  console.error(err);
                  setAdminGlobalMessage('Terjadi kesalahan saat mengunggah.');
                }

                e.target.value = '';
              }}
            />
          </label>
        </div>
      </div>

      {/* Global Save Slider Button */}
      <div className="border-t border-slate-100 pt-4 flex justify-end">
        <button
          type="button"
          onClick={async () => {
            const token = localStorage.getItem('auth_token');
            if (!token) return;
            const siteId = siteConfig?.id || 'ppid';
            try {
              const res = await fetch(`${API_BASE_URL}/admin/sites/${siteId}`, {
                method: 'PUT',
                headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
                body: JSON.stringify(siteConfig),
              });
              if (res.ok) {
                setAdminGlobalMessage('Semua konten slide berhasil disimpan!');
              } else {
                setAdminGlobalMessage('Gagal menyimpan data slide.');
              }
            } catch (err) {
              console.error(err);
              setAdminGlobalMessage('Terjadi kesalahan.');
            }
          }}
          className="px-6 py-3 bg-[#002147] hover:bg-[#00346e] text-white text-xs font-bold uppercase rounded-xl shadow-md transition-all flex items-center gap-2 cursor-pointer border border-[#002147]/50"
        >
          <Save className="h-4 w-4 text-amber-400" /> Simpan Semua Konten Slide
        </button>
      </div>

      {/* Overlay Settings */}
      <div className="border-t border-slate-100 pt-4 space-y-4 text-left">
        <h3 className="text-xs font-bold text-slate-600 uppercase tracking-wider">Pengaturan Overlay Slider</h3>
        <p className="text-[11px] text-slate-400">Atur warna dan transparansi overlay yang menutupi gambar slider.</p>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
          {/* Overlay Color */}
          <div className="space-y-2 text-left">
            <label className="text-xs font-bold text-slate-700 block">Warna Overlay</label>
            <div className="flex items-center gap-3">
              <input
                type="color"
                value={siteConfig?.settings?.hero_overlay_color || '#002147'}
                onChange={(e) => {
                  setSiteConfig((prev: any) => ({
                    ...prev,
                    settings: { ...prev?.settings, hero_overlay_color: e.target.value },
                  }));
                }}
                className="w-12 h-10 rounded-lg border border-slate-300 cursor-pointer p-0.5"
              />
              <input
                type="text"
                value={siteConfig?.settings?.hero_overlay_color || '#002147'}
                onChange={(e) => {
                  setSiteConfig((prev: any) => ({
                    ...prev,
                    settings: { ...prev?.settings, hero_overlay_color: e.target.value },
                  }));
                }}
                className="flex-1 text-xs border border-slate-300 rounded-xl px-3 py-2.5 font-mono focus:ring-2 focus:ring-amber-400 focus:outline-none text-slate-805"
                placeholder="#002147"
              />
            </div>
          </div>

          {/* Overlay Opacity */}
          <div className="space-y-2 text-left">
            <label className="text-xs font-bold text-slate-700 block">
              Transparansi Gambar: <span className="text-amber-500 font-extrabold">{Math.round((siteConfig?.settings?.hero_image_opacity !== undefined ? parseFloat(siteConfig.settings.hero_image_opacity) : 0.2) * 100)}%</span>
            </label>
            <input
              type="range"
              min="0"
              max="1"
              step="0.05"
              value={siteConfig?.settings?.hero_image_opacity !== undefined ? siteConfig.settings.hero_image_opacity : '0.2'}
              onChange={(e) => {
                setSiteConfig((prev: any) => ({
                  ...prev,
                  settings: { ...prev?.settings, hero_image_opacity: e.target.value },
                }));
              }}
              className="w-full h-2 bg-slate-200 rounded-lg appearance-none cursor-pointer accent-amber-500"
            />
            <div className="flex justify-between text-[9px] text-slate-400 font-medium">
              <span>0% (Gelap penuh)</span>
              <span>50%</span>
              <span>100% (Tanpa overlay)</span>
            </div>
          </div>
        </div>

        {/* Live Preview */}
        <div className="space-y-2 text-left">
          <label className="text-xs font-bold text-slate-700 block">Preview Overlay</label>
          <div className="relative rounded-2xl overflow-hidden h-40 border border-slate-200">
            <img
              src={heroImages.length > 0
                ? resolveImageUrl(
                  typeof heroImages[0] === 'string'
                    ? heroImages[0]
                    : heroImages[0]?.image || ''
                )
                : '/logo.png'}
              alt="Preview"
              className="w-full h-full object-cover"
            />
            <div
              className="absolute inset-0"
              style={{
                backgroundColor: siteConfig?.settings?.hero_overlay_color || '#002147',
                opacity: siteConfig?.settings?.hero_image_opacity !== undefined
                  ? (1 - parseFloat(siteConfig.settings.hero_image_opacity))
                  : 0.8,
              }}
            />
            <div className="absolute inset-0 flex items-center justify-center z-10">
              <span className="text-white text-sm font-extrabold drop-shadow-lg">Contoh Teks di Atas Slider</span>
            </div>
          </div>
        </div>

        {/* Save Button */}
        <button
          type="button"
          onClick={async () => {
            const token = localStorage.getItem('auth_token');
            if (!token) return;
            const siteId = siteConfig?.id || 'ppid';
            try {
              const res = await fetch(`${API_BASE_URL}/admin/sites/${siteId}`, {
                method: 'PUT',
                headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
                body: JSON.stringify(siteConfig),
              });
              if (res.ok) {
                setAdminGlobalMessage('Pengaturan overlay berhasil disimpan!');
              } else {
                setAdminGlobalMessage('Gagal menyimpan pengaturan overlay.');
              }
            } catch (err) {
              console.error(err);
              setAdminGlobalMessage('Terjadi kesalahan.');
            }
          }}
          className="px-6 py-2.5 bg-amber-400 hover:bg-amber-500 text-[#002147] text-xs font-extrabold rounded-xl shadow-sm transition-all flex items-center gap-2 cursor-pointer border border-amber-400/50"
        >
          <Save className="h-4 w-4" /> Simpan Pengaturan Overlay
        </button>
      </div>

      {/* Slider Animation Settings */}
      <div className="border-t border-slate-100 pt-4 space-y-4 text-left">
        <h3 className="text-xs font-bold text-slate-600 uppercase tracking-wider">Efek Animasi Slider (Slider Transitions)</h3>
        <p className="text-[11px] text-slate-400">Pilih gaya transisi animasi modern dan kekinian saat slide berganti di halaman beranda.</p>

        <div className="grid grid-cols-2 md:grid-cols-5 gap-3">
          {[
            { key: 'shutter-3d', label: '3D Shutter Grid', desc: 'Pecahan grid shutter 3D yang meledak keluar' },
            { key: 'fade-blur', label: 'Cinematic Fade Blur', desc: 'Transisi blur sinematik dengan perpindahan lembut' },
            { key: 'parallax-slide', label: 'Parallax Slide', desc: 'Geser horizontal dengan kecepatan layer berbeda' },
            { key: 'ken-burns', label: 'Ken Burns Zoom', desc: 'Zoom-in perlahan dipadu dengan fade klasik' },
            { key: 'split-diagonal', label: 'Diagonal Split', desc: 'Slide miring diagonal yang memotong layar kreatif' }
          ].map((anim) => {
            const isSelected = (siteConfig?.settings?.hero_animation || 'shutter-3d') === anim.key;
            return (
              <button
                key={anim.key}
                type="button"
                onClick={() => {
                  setSiteConfig((prev: any) => ({
                    ...prev,
                    settings: { ...prev?.settings, hero_animation: anim.key }
                  }));
                }}
                className={`p-3.5 rounded-2xl border text-left flex flex-col justify-between h-28 cursor-pointer transition-all ${
                  isSelected
                    ? 'border-[#002147] bg-[#002147]/5 ring-2 ring-[#002147]/10'
                    : 'border-slate-200 bg-white hover:bg-slate-50'
                }`}
              >
                <span className={`text-[10px] font-extrabold uppercase tracking-wide ${
                  isSelected ? 'text-[#002147]' : 'text-slate-700'
                }`}>
                  {anim.label}
                </span>
                <span className="text-[8.5px] text-slate-400 font-semibold leading-relaxed mt-2 block">
                  {anim.desc}
                </span>
              </button>
            );
          })}
        </div>

        <button
          type="button"
          onClick={async () => {
            const token = localStorage.getItem('auth_token');
            if (!token) return;
            const siteId = siteConfig?.id || 'ppid';
            try {
              const res = await fetch(`${API_BASE_URL}/admin/sites/${siteId}`, {
                method: 'PUT',
                headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
                body: JSON.stringify(siteConfig),
              });
              if (res.ok) {
                setAdminGlobalMessage('Efek transisi slider berhasil disimpan!');
              } else {
                setAdminGlobalMessage('Gagal menyimpan efek transisi slider.');
              }
            } catch (err) {
              console.error(err);
              setAdminGlobalMessage('Terjadi kesalahan.');
            }
          }}
          className="px-6 py-2.5 bg-amber-400 hover:bg-amber-500 text-[#002147] text-xs font-extrabold rounded-xl shadow-sm transition-all flex items-center gap-2 cursor-pointer border border-amber-400/50"
        >
          <Save className="h-4 w-4" /> Simpan Efek Transisi
        </button>
      </div>

      {/* Jam Pelayanan Settings */}
      <div className="border-t border-slate-100 pt-4 space-y-4 text-left">
        <h3 className="text-xs font-bold text-slate-600 uppercase tracking-wider">Jam Pelayanan PPID</h3>
        <p className="text-[11px] text-slate-400">Atur jadwal hari dan jam pelayanan yang ditampilkan di hero section.</p>

        {/* Current Entries */}
        <div className="space-y-2">
          {(siteConfig?.settings?.service_hours && Array.isArray(siteConfig.settings.service_hours) && siteConfig.settings.service_hours.length > 0
            ? siteConfig.settings.service_hours
            : [
              { day: 'Senin – Kamis', time: '08:35 – 16:00 WIB', closed: false },
              { day: 'Jumat', time: '08:30 – 12:30 WIB', closed: false },
              { day: 'Sabtu, Minggu & Hari Besar', time: 'Tutup / Libur', closed: true },
            ]
          ).map((sh: any, idx: number) => (
            <div key={idx} className="flex items-center gap-2 bg-slate-50 border border-slate-200 rounded-xl px-3 py-2.5 text-left">
              <input
                type="text"
                value={sh.day}
                onChange={(e) => {
                  const hours = [...(siteConfig?.settings?.service_hours || [
                    { day: 'Senin – Kamis', time: '08:35 – 16:00 WIB', closed: false },
                    { day: 'Jumat', time: '08:30 – 12:30 WIB', closed: false },
                    { day: 'Sabtu, Minggu & Hari Besar', time: 'Tutup / Libur', closed: true },
                  ])];
                  hours[idx] = { ...hours[idx], day: e.target.value };
                  setSiteConfig((prev: any) => ({ ...prev, settings: { ...prev?.settings, service_hours: hours } }));
                }}
                className="flex-1 text-xs border border-slate-300 rounded-lg px-2.5 py-2 focus:ring-2 focus:ring-amber-400 focus:outline-none"
                placeholder="Contoh: Senin – Kamis"
              />
              <input
                type="text"
                value={sh.time}
                onChange={(e) => {
                  const hours = [...(siteConfig?.settings?.service_hours || [
                    { day: 'Senin – Kamis', time: '08:35 – 16:00 WIB', closed: false },
                    { day: 'Jumat', time: '08:30 – 12:30 WIB', closed: false },
                    { day: 'Sabtu, Minggu & Hari Besar', time: 'Tutup / Libur', closed: true },
                  ])];
                  hours[idx] = { ...hours[idx], time: e.target.value };
                  setSiteConfig((prev: any) => ({ ...prev, settings: { ...prev?.settings, service_hours: hours } }));
                }}
                className="flex-1 text-xs border border-slate-300 rounded-lg px-2.5 py-2 focus:ring-2 focus:ring-amber-400 focus:outline-none"
                placeholder="Contoh: 08:00 – 16:00 WIB"
              />
              <label className="flex items-center gap-1.5 text-[10px] text-slate-500 font-bold cursor-pointer whitespace-nowrap">
                <input
                  type="checkbox"
                  checked={!!sh.closed}
                  onChange={(e) => {
                    const hours = [...(siteConfig?.settings?.service_hours || [
                      { day: 'Senin – Kamis', time: '08:35 – 16:00 WIB', closed: false },
                      { day: 'Jumat', time: '08:30 – 12:30 WIB', closed: false },
                      { day: 'Sabtu, Minggu & Hari Besar', time: 'Tutup / Libur', closed: true },
                    ])];
                    hours[idx] = { ...hours[idx], closed: e.target.checked };
                    setSiteConfig((prev: any) => ({ ...prev, settings: { ...prev?.settings, service_hours: hours } }));
                  }}
                  className="accent-red-500"
                />
                Tutup
              </label>
              <button
                type="button"
                onClick={() => {
                  const hours = [...(siteConfig?.settings?.service_hours || [
                    { day: 'Senin – Kamis', time: '08:35 – 16:00 WIB', closed: false },
                    { day: 'Jumat', time: '08:30 – 12:30 WIB', closed: false },
                    { day: 'Sabtu, Minggu & Hari Besar', time: 'Tutup / Libur', closed: true },
                  ])];
                  hours.splice(idx, 1);
                  setSiteConfig((prev: any) => ({ ...prev, settings: { ...prev?.settings, service_hours: hours } }));
                }}
                className="p-1.5 text-red-400 hover:text-red-650 hover:bg-red-55 rounded-lg transition-all"
                title="Hapus baris"
              >
                <Trash2 className="h-3.5 w-3.5" />
              </button>
            </div>
          ))}
        </div>

        {/* Add Row */}
        <button
          type="button"
          onClick={() => {
            const hours = [...(siteConfig?.settings?.service_hours || [
              { day: 'Senin – Kamis', time: '08:35 – 16:00 WIB', closed: false },
              { day: 'Jumat', time: '08:30 – 12:30 WIB', closed: false },
              { day: 'Sabtu, Minggu & Hari Besar', time: 'Tutup / Libur', closed: true },
            ])];
            hours.push({ day: '', time: '', closed: false });
            setSiteConfig((prev: any) => ({ ...prev, settings: { ...prev?.settings, service_hours: hours } }));
          }}
          className="px-4 py-2 bg-slate-100 hover:bg-slate-200 text-slate-600 text-[10px] font-bold rounded-xl transition-all flex items-center gap-1.5 cursor-pointer"
        >
          <Plus className="h-3.5 w-3.5" /> Tambah Baris Jadwal
        </button>

        {/* Service Location */}
        <div className="space-y-1.5 text-left">
          <label className="text-xs font-bold text-slate-700 block">Lokasi Pelayanan</label>
          <input
            type="text"
            value={siteConfig?.settings?.service_location || 'Gedung A - Samping B'}
            onChange={(e) => {
              setSiteConfig((prev: any) => ({ ...prev, settings: { ...prev?.settings, service_location: e.target.value } }));
            }}
            className="w-full text-xs border border-slate-300 rounded-xl px-3 py-2.5 focus:ring-2 focus:ring-amber-400 focus:outline-none"
            placeholder="Contoh: Gedung A - Samping B"
          />
        </div>

        {/* Save Button */}
        <button
          type="button"
          onClick={async () => {
            const token = localStorage.getItem('auth_token');
            if (!token) return;
            const siteId = siteConfig?.id || 'ppid';
            try {
              const res = await fetch(`${API_BASE_URL}/admin/sites/${siteId}`, {
                method: 'PUT',
                headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
                body: JSON.stringify(siteConfig),
              });
              if (res.ok) {
                setAdminGlobalMessage('Jam pelayanan berhasil disimpan!');
              } else {
                setAdminGlobalMessage('Gagal menyimpan jam pelayanan.');
              }
            } catch (err) {
              console.error(err);
              setAdminGlobalMessage('Terjadi kesalahan.');
            }
          }}
          className="px-6 py-2.5 bg-amber-400 hover:bg-amber-500 text-[#002147] text-xs font-extrabold rounded-xl shadow-sm transition-all flex items-center gap-2 cursor-pointer border border-amber-400/50"
        >
          <Save className="h-4 w-4" /> Simpan Jam Pelayanan
        </button>
      </div>

      {/* Slider Buttons Settings */}
      <div className="border-t border-slate-100 pt-4 space-y-4 text-left">
        <h3 className="text-xs font-bold text-slate-600 uppercase tracking-wider">Pengaturan Tombol Slider</h3>
        <p className="text-[11px] text-slate-400">Atur teks label dan tujuan navigasi untuk 2 tombol di slider utama.</p>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
          {/* Tombol Utama */}
          <div className="bg-slate-50 border border-slate-200 rounded-2xl p-4 space-y-3">
            <h4 className="text-xs font-bold text-[#002147] flex items-center gap-1.5 text-left">
              <FileText className="h-4 w-4 text-amber-500" /> Tombol Utama (Kiri)
            </h4>
            <div className="space-y-1.5 text-left">
              <label className="text-[10px] font-bold text-slate-500 uppercase tracking-wider">Teks Label</label>
              <input
                type="text"
                value={siteConfig?.settings?.hero_btn1_text || ''}
                onChange={(e) => {
                  setSiteConfig((prev: any) => ({
                    ...prev,
                    settings: { ...prev?.settings, hero_btn1_text: e.target.value },
                  }));
                }}
                className="w-full text-xs border border-slate-300 rounded-xl px-3 py-2 focus:ring-2 focus:ring-amber-400 focus:outline-none"
                placeholder="Default: Permohonan Informasi"
              />
            </div>
            <div className="space-y-1.5 text-left">
              <label className="text-[10px] font-bold text-slate-500 uppercase tracking-wider">Tujuan Navigasi (Slug Halaman / Rute)</label>
              <input
                type="text"
                value={siteConfig?.settings?.hero_btn1_page || ''}
                onChange={(e) => {
                  setSiteConfig((prev: any) => ({
                    ...prev,
                    settings: { ...prev?.settings, hero_btn1_page: e.target.value },
                  }));
                }}
                className="w-full text-xs border border-slate-300 rounded-xl px-3 py-2 focus:ring-2 focus:ring-amber-400 focus:outline-none font-mono"
                placeholder="Default: permohonan-informasi"
              />
            </div>
          </div>

          {/* Tombol Sekunder */}
          <div className="bg-slate-50 border border-slate-200 rounded-2xl p-4 space-y-3">
            <h4 className="text-xs font-bold text-[#002147] flex items-center gap-1.5 text-left">
              <Search className="h-4 w-4 text-amber-500" /> Tombol Sekunder (Kanan)
            </h4>
            <div className="space-y-1.5 text-left">
              <label className="text-[10px] font-bold text-slate-500 uppercase tracking-wider">Teks Label</label>
              <input
                type="text"
                value={siteConfig?.settings?.hero_btn2_text || ''}
                onChange={(e) => {
                  setSiteConfig((prev: any) => ({
                    ...prev,
                    settings: { ...prev?.settings, hero_btn2_text: e.target.value },
                  }));
                }}
                className="w-full text-xs border border-slate-300 rounded-xl px-3 py-2 focus:ring-2 focus:ring-amber-400 focus:outline-none"
                placeholder="Default: Cari Dokumen Publik"
              />
            </div>
            <div className="space-y-1.5 text-left">
              <label className="text-[10px] font-bold text-slate-500 uppercase tracking-wider">Tujuan Navigasi (Slug Halaman / Rute)</label>
              <input
                type="text"
                value={siteConfig?.settings?.hero_btn2_page || ''}
                onChange={(e) => {
                  setSiteConfig((prev: any) => ({
                    ...prev,
                    settings: { ...prev?.settings, hero_btn2_page: e.target.value },
                  }));
                }}
                className="w-full text-xs border border-slate-300 rounded-xl px-3 py-2 focus:ring-2 focus:ring-amber-400 focus:outline-none font-mono"
                placeholder="Default: regulasi"
              />
            </div>
          </div>
        </div>

        {/* Save Button */}
        <button
          type="button"
          onClick={async () => {
            const token = localStorage.getItem('auth_token');
            if (!token) return;
            const siteId = siteConfig?.id || 'ppid';
            try {
              const res = await fetch(`${API_BASE_URL}/admin/sites/${siteId}`, {
                method: 'PUT',
                headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
                body: JSON.stringify(siteConfig),
              });
              if (res.ok) {
                setAdminGlobalMessage('Pengaturan tombol slider berhasil disimpan!');
              } else {
                setAdminGlobalMessage('Gagal menyimpan pengaturan tombol slider.');
              }
            } catch (err) {
              console.error(err);
              setAdminGlobalMessage('Terjadi kesalahan.');
            }
          }}
          className="px-6 py-2.5 bg-amber-400 hover:bg-amber-500 text-[#002147] text-xs font-extrabold rounded-xl shadow-sm transition-all flex items-center gap-2 cursor-pointer border border-amber-400/50"
        >
          <Save className="h-4 w-4" /> Simpan Pengaturan Tombol
        </button>
      </div>
    </div>
  );
}
