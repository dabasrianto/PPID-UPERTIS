import React from 'react';
import { ChevronLeft, Plus, Trash2, FileText, Edit, ArrowUp, ArrowDown, Upload, Loader2, Info } from 'lucide-react';
import { resolveImageUrl, parseRegulasiHTML, defaultRegulasiData } from '../../utils/helpers';

interface PageDoc {
  title: string;
  description: string;
  file_url: string;
}

interface ManagePagesProps {
  editModalOpen: boolean;
  setEditModalOpen: (val: boolean) => void;
  editModalType: string;
  activeEditItem: any;
  adminEditTitle: string;
  setAdminEditTitle: (val: string) => void;
  adminEditSubtitle: string;
  setAdminEditSubtitle: (val: string) => void;
  adminEditSlug: string;
  setAdminEditSlug: (val: string) => void;
  adminEditCoverImage: string;
  setAdminEditCoverImage: (val: string) => void;
  handleSaveCrudItem: (e: React.FormEvent) => void;
  handleAdminUpload: (e: React.ChangeEvent<HTMLInputElement>) => void;
  openCreateModal: (type: string) => void;
  adminPages: any[];
  openEditModal: (type: string, item: any) => void;
  handleDeleteCrudItem: (table: string, id: string) => void;
  API_BASE_URL: string;

  // Jadwal Page states
  editJadwalSeninKamisKerja: string;
  setEditJadwalSeninKamisKerja: (val: string) => void;
  editJadwalSeninKamisIstirahat: string;
  setEditJadwalSeninKamisIstirahat: (val: string) => void;
  editJadwalJumatKerja: string;
  setEditJadwalJumatKerja: (val: string) => void;
  editJadwalPageIstirahatJumat: string;
  setEditJadwalPageIstirahatJumat: (val: string) => void;
  editJadwalSabtuMinggu: string;
  setEditJadwalSabtuMinggu: (val: string) => void;
  editJadwalOffline1: string;
  setEditJadwalOffline1: (val: string) => void;
  editJadwalOffline2: string;
  setEditJadwalOffline2: (val: string) => void;
  editJadwalOffline3: string;
  setEditJadwalOffline3: (val: string) => void;
  editJadwalOnline1: string;
  setEditJadwalOnline1: (val: string) => void;
  editJadwalOnline2: string;
  setEditJadwalOnline2: (val: string) => void;
  editJadwalOnline3: string;
  setEditJadwalOnline3: (val: string) => void;
  editJadwalCustomRemarks: string;
  setEditJadwalCustomRemarks: (val: string) => void;

  // Page Docs list state
  adminEditPageDocs: PageDoc[];
  setAdminEditPageDocs: React.Dispatch<React.SetStateAction<PageDoc[]>>;

  // Keberatan Steps States
  editKeberatanManualSteps: any[];
  setEditKeberatanManualSteps: (val: any[]) => void;
  editKeberatanOnlineSteps: any[];
  setEditKeberatanOnlineSteps: (val: any[]) => void;

  // Content string state
  adminEditContent: string;
  setAdminEditContent: (val: string) => void;

  // Profil Stats Cards
  editProfilStats: Array<{ value: string; label: string }>;
  setEditProfilStats: (val: Array<{ value: string; label: string }>) => void;
  editProfilImage1: string;
  setEditProfilImage1: (val: string) => void;
  editProfilImage2: string;
  setEditProfilImage2: (val: string) => void;
  editProfilImage3: string;
  setEditProfilImage3: (val: string) => void;
  editTugasList: Array<{ title: string; desc: string }>;
  setEditTugasList: (val: Array<{ title: string; desc: string }>) => void;
  editFungsiList: Array<{ title: string; items: string[] }>;
  setEditFungsiList: (val: Array<{ title: string; items: string[] }>) => void;
  editTugasImage1: string;
  setEditTugasImage1: (val: string) => void;
  editTugasImage2: string;
  setEditTugasImage2: (val: string) => void;
  editTugasImage3: string;
  setEditTugasImage3: (val: string) => void;

  // Sengketa Verification Cards
  editSengketaCards: Array<{ title: string; desc: string }>;
  setEditSengketaCards: (val: Array<{ title: string; desc: string }>) => void;

  // Sengketa Flow Text
  editSengketaFlow: {
    jalurA_title: string; jalurA_desc: string;
    jalurB_title: string; jalurB_desc: string;
    verifikasi_title: string; verifikasi_desc: string;
    help_title: string; help_desc: string;
    download_title: string; download_desc: string;
  };
  setEditSengketaFlow: (val: any) => void;
}

export default function ManagePages({
  editModalOpen,
  setEditModalOpen,
  editModalType,
  activeEditItem,
  adminEditTitle,
  setAdminEditTitle,
  adminEditSubtitle,
  setAdminEditSubtitle,
  adminEditSlug,
  setAdminEditSlug,
  adminEditCoverImage,
  setAdminEditCoverImage,
  handleSaveCrudItem,
  handleAdminUpload,
  openCreateModal,
  adminPages,
  openEditModal,
  handleDeleteCrudItem,
  API_BASE_URL,

  editJadwalSeninKamisKerja,
  setEditJadwalSeninKamisKerja,
  editJadwalSeninKamisIstirahat,
  setEditJadwalSeninKamisIstirahat,
  editJadwalJumatKerja,
  setEditJadwalJumatKerja,
  editJadwalPageIstirahatJumat,
  setEditJadwalPageIstirahatJumat,
  editJadwalSabtuMinggu,
  setEditJadwalSabtuMinggu,
  editJadwalOffline1,
  setEditJadwalOffline1,
  editJadwalOffline2,
  setEditJadwalOffline2,
  editJadwalOffline3,
  setEditJadwalOffline3,
  editJadwalOnline1,
  setEditJadwalOnline1,
  editJadwalOnline2,
  setEditJadwalOnline2,
  editJadwalOnline3,
  setEditJadwalOnline3,
  editJadwalCustomRemarks,
  setEditJadwalCustomRemarks,

  adminEditPageDocs,
  setAdminEditPageDocs,

  editKeberatanManualSteps,
  setEditKeberatanManualSteps,
  editKeberatanOnlineSteps,
  setEditKeberatanOnlineSteps,

  adminEditContent,
  setAdminEditContent,
  editProfilStats,
  setEditProfilStats,
  editProfilImage1,
  setEditProfilImage1,
  editProfilImage2,
  setEditProfilImage2,
  editProfilImage3,
  setEditProfilImage3,
  editTugasList,
  setEditTugasList,
  editFungsiList,
  setEditFungsiList,
  editTugasImage1,
  setEditTugasImage1,
  editTugasImage2,
  setEditTugasImage2,
  editTugasImage3,
  setEditTugasImage3,
  editSengketaCards,
  setEditSengketaCards,
  editSengketaFlow,
  setEditSengketaFlow
}: ManagePagesProps) {
  const [isUploadingDocs, setIsUploadingDocs] = React.useState(false);
  // Local state for Regulasi accordion editor
  const [regulasiGroups, setRegulasiGroups] = React.useState<any[]>([]);
  const [isUploadingRegulasiFile, setIsUploadingRegulasiFile] = React.useState<Record<string, boolean>>({});

  // Sync parent state with local list when parent opens/changes
  React.useEffect(() => {
    if (!editModalOpen || editModalType !== 'page' || adminEditSlug !== 'regulasi') {
      setRegulasiGroups([]);
      return;
    }

    if (adminEditContent) {
      const trimmed = adminEditContent.trim();
      if (trimmed.startsWith('[') || trimmed.startsWith('{')) {
        try {
          const parsed = JSON.parse(trimmed);
          if (Array.isArray(parsed) && parsed.length > 0) {
            setRegulasiGroups(parsed);
            return;
          }
        } catch (e) {
          console.error('Failed to parse regulasi JSON content:', e);
        }
      }

      // Text-based fallback: parse using parseRegulasiHTML
      const parsed = parseRegulasiHTML(adminEditContent);
      setRegulasiGroups(parsed.length > 0 ? parsed : defaultRegulasiData);
    } else {
      setRegulasiGroups(defaultRegulasiData);
    }
  }, [editModalOpen, editModalType, adminEditSlug, adminEditContent]);

  // Propagate local list changes back to the parent state as serialized JSON
  const updateRegulasiContent = (groups: any[]) => {
    setRegulasiGroups(groups);
    setAdminEditContent(JSON.stringify(groups));
  };
  // Local state for DIP pages text description sections
  interface DipSection {
    text: string;
    imageUrls: string[];
    imagePosition: 'left' | 'right';
  }
  const [dipSections, setDipSections] = React.useState<DipSection[]>([]);
  const [isUploadingDipSectionImage, setIsUploadingDipSectionImage] = React.useState<Record<string, boolean>>({});

  // Sync parent state with local dipSections state when modal opens/changes
  React.useEffect(() => {
    if (!editModalOpen || editModalType !== 'page' || !['informasi-publik-berkala', 'informasi-tersedia-setiap-saat', 'info-serta-merta', 'informasi-dikecualikan'].includes(adminEditSlug)) {
      setDipSections([]);
      return;
    }

    if (adminEditContent) {
      const trimmed = adminEditContent.trim();
      if (trimmed.startsWith('{')) {
        try {
          const parsed = JSON.parse(trimmed);
          if (parsed && typeof parsed === 'object') {
            if (Array.isArray(parsed.sections)) {
              const mapped = parsed.sections.map((sec: any) => ({
                text: sec.text || '',
                imageUrls: Array.isArray(sec.imageUrls) ? sec.imageUrls : (sec.imageUrl ? [sec.imageUrl] : []),
                imagePosition: sec.imagePosition || 'right'
              }));
              setDipSections(mapped);
              return;
            } else {
              // Legacy JSON format (has intro but no sections list)
              setDipSections([{ text: parsed.intro || '', imageUrls: [], imagePosition: 'right' }]);
              return;
            }
          }
        } catch (e) {
          console.error('Failed to parse DIP page content JSON:', e);
        }
      }
      // Plain text legacy format
      setDipSections([{ text: adminEditContent, imageUrls: [], imagePosition: 'right' }]);
    } else {
      setDipSections([]);
    }
  }, [editModalOpen, editModalType, adminEditSlug, adminEditContent]);

  // Propagate local dipSections changes back to the parent state as serialized JSON
  const updateDipSections = (sections: DipSection[]) => {
    setDipSections(sections);
    setAdminEditContent(JSON.stringify({ 
      intro: sections[0]?.text || '', 
      sections 
    }));
  };

  // Local state for Struktur Organisasi
  const [strukturAtasan, setStrukturAtasan] = React.useState('');
  const [strukturUtama, setStrukturUtama] = React.useState('');
  const [strukturPelaksana, setStrukturPelaksana] = React.useState<string[]>([]);
  const [strukturPertimbangan, setStrukturPertimbangan] = React.useState<string[]>([]);
  const [strukturPelayanan, setStrukturPelayanan] = React.useState<string[]>([]);
  const [strukturDesc, setStrukturDesc] = React.useState('');

  // Sync parent state with local states for Struktur Organisasi when modal opens/changes
  React.useEffect(() => {
    if (!editModalOpen || editModalType !== 'page' || adminEditSlug !== 'struktur-organisasi-2') {
      setStrukturAtasan('');
      setStrukturUtama('');
      setStrukturPelaksana([]);
      setStrukturPertimbangan([]);
      setStrukturPelayanan([]);
      setStrukturDesc('');
      return;
    }

    const defaultPelaksana = ["Biro Humas", "Biro Akademik", "Biro Umum", "Dekan Fakultas"];
    const defaultPertimbangan = ["ka P2AMIA", "ka LPPM", "ka P3TS", "Ka Prodi", "Ka. UPT"];
    const defaultPelayanan = ["Staff Humas", "LTIK"];

    if (adminEditContent) {
      const trimmed = adminEditContent.trim();
      if (trimmed.startsWith('{')) {
        try {
          const parsed = JSON.parse(trimmed);
          if (parsed && typeof parsed === 'object') {
            setStrukturAtasan(parsed.atasan || 'Rektor Universitas Perintis');
            setStrukturUtama(parsed.utama || 'Wakil Rektor 1 & 2 UPERTIS');
            setStrukturPelaksana(Array.isArray(parsed.pelaksana) ? parsed.pelaksana : defaultPelaksana);
            setStrukturPertimbangan(Array.isArray(parsed.pertimbangan) ? parsed.pertimbangan : defaultPertimbangan);
            setStrukturPelayanan(Array.isArray(parsed.pelayanan) ? parsed.pelayanan : defaultPelayanan);
            setStrukturDesc(parsed.desc || '');
            return;
          }
        } catch (e) {
          console.error('Failed to parse Struktur Organisasi content JSON:', e);
        }
      }
      // Fallback legacy format
      setStrukturAtasan('Rektor Universitas Perintis');
      setStrukturUtama('Wakil Rektor 1 & 2 UPERTIS');
      setStrukturPelaksana(defaultPelaksana);
      setStrukturPertimbangan(defaultPertimbangan);
      setStrukturPelayanan(defaultPelayanan);
      setStrukturDesc(adminEditContent);
    } else {
      setStrukturAtasan('Rektor Universitas Perintis');
      setStrukturUtama('Wakil Rektor 1 & 2 UPERTIS');
      setStrukturPelaksana(defaultPelaksana);
      setStrukturPertimbangan(defaultPertimbangan);
      setStrukturPelayanan(defaultPelayanan);
      setStrukturDesc('');
    }
  }, [editModalOpen, editModalType, adminEditSlug, adminEditContent]);

  // Propagate local states back to parent state as serialized JSON
  const updateStrukturContent = (fields: {
    atasan?: string;
    utama?: string;
    pelaksana?: string[];
    pertimbangan?: string[];
    pelayanan?: string[];
    desc?: string;
  }) => {
    const finalAtasan = fields.atasan !== undefined ? fields.atasan : strukturAtasan;
    const finalUtama = fields.utama !== undefined ? fields.utama : strukturUtama;
    const finalPelaksana = fields.pelaksana !== undefined ? fields.pelaksana : strukturPelaksana;
    const finalPertimbangan = fields.pertimbangan !== undefined ? fields.pertimbangan : strukturPertimbangan;
    const finalPelayanan = fields.pelayanan !== undefined ? fields.pelayanan : strukturPelayanan;
    const finalDesc = fields.desc !== undefined ? fields.desc : strukturDesc;

    if (fields.atasan !== undefined) setStrukturAtasan(fields.atasan);
    if (fields.utama !== undefined) setStrukturUtama(fields.utama);
    if (fields.pelaksana !== undefined) setStrukturPelaksana(fields.pelaksana);
    if (fields.pertimbangan !== undefined) setStrukturPertimbangan(fields.pertimbangan);
    if (fields.pelayanan !== undefined) setStrukturPelayanan(fields.pelayanan);
    if (fields.desc !== undefined) setStrukturDesc(fields.desc);

    setAdminEditContent(JSON.stringify({
      atasan: finalAtasan,
      utama: finalUtama,
      pelaksana: finalPelaksana,
      pertimbangan: finalPertimbangan,
      pelayanan: finalPelayanan,
      desc: finalDesc
    }));
  };

  // Local state for Visi Misi
  const [visiText, setVisiText] = React.useState('');
  const [misiList, setMisiList] = React.useState<string[]>([]);

  // Sync parent state with local states for Visi Misi when modal opens/changes
  React.useEffect(() => {
    if (!editModalOpen || editModalType !== 'page' || adminEditSlug !== 'visi-misi') {
      setVisiText('');
      setMisiList([]);
      return;
    }

    const defaultMisi = [
      "Menyediakan pelayanan informasi publik yang cepat, tepat waktu, dan akurat.",
      "Mengembangkan sistem pengelolaan dokumen berbasis teknologi informasi yang aman dan mudah diakses.",
      "Meningkatkan kapasitas sumber daya pengelola layanan informasi secara berkelanjutan.",
      "Mewujudkan tata kelola perguruan tinggi yang bersih, transparan, dan akuntabel."
    ];

    if (adminEditContent) {
      const trimmed = adminEditContent.trim();
      if (trimmed.startsWith('{')) {
        try {
          const parsed = JSON.parse(trimmed);
          if (parsed && typeof parsed === 'object') {
            setVisiText(parsed.visi || '');
            setMisiList(Array.isArray(parsed.misi) ? parsed.misi : defaultMisi);
            return;
          }
        } catch (e) {
          console.error('Failed to parse Visi Misi content JSON:', e);
        }
      }
      // Fallback legacy format
      setVisiText(adminEditContent);
      setMisiList(defaultMisi);
    } else {
      setVisiText("Menjadi Pejabat Pengelola Informasi dan Dokumentasi (PPID) yang unggul, terpercaya, dan transparan dalam pelayanan informasi publik di lingkungan Universitas Perintis Indonesia.");
      setMisiList(defaultMisi);
    }
  }, [editModalOpen, editModalType, adminEditSlug, adminEditContent]);

  // Propagate local states back to parent state as serialized JSON
  const updateVisiMisiContent = (fields: {
    visi?: string;
    misi?: string[];
  }) => {
    const finalVisi = fields.visi !== undefined ? fields.visi : visiText;
    const finalMisi = fields.misi !== undefined ? fields.misi : misiList;

    if (fields.visi !== undefined) setVisiText(fields.visi);
    if (fields.misi !== undefined) setMisiList(fields.misi);

    setAdminEditContent(JSON.stringify({
      visi: finalVisi,
      misi: finalMisi
    }));
  };



  const handleBulkUploadPageDocs = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files;
    if (!files || files.length === 0) return;

    setIsUploadingDocs(true);
    const token = localStorage.getItem('auth_token');
    if (!token) {
      alert('Sesi login telah berakhir. Silakan login kembali.');
      setIsUploadingDocs(false);
      return;
    }

    const formData = new FormData();
    formData.append('folder', 'documents');
    for (let i = 0; i < files.length; i++) {
      formData.append('files', files[i]);
    }

    try {
      const res = await fetch(`${API_BASE_URL}/admin/uploads`, {
        method: 'POST',
        headers: { Authorization: `Bearer ${token}` },
        body: formData
      });
      if (res.ok) {
        const data = await res.json();
        const uploadedList = data.uploaded || [];
        if (uploadedList.length > 0) {
          const newDocs = uploadedList.map((item: any) => {
            const titleWithoutExt = item.original_name ? item.original_name.replace(/\.[^/.]+$/, "") : "Dokumen";
            const formattedTitle = titleWithoutExt
              .replace(/[_-]+/g, ' ')
              .replace(/(^\w|\s\w)/g, (m: string) => m.toUpperCase());
            return {
              title: formattedTitle,
              description: '',
              file_url: item.url
            };
          });
          setAdminEditPageDocs(prev => [...prev, ...newDocs]);
        }
      } else {
        const data = await res.json();
        alert(data.error || 'Gagal mengunggah berkas.');
      }
    } catch (err) {
      console.error(err);
      alert('Error saat mengunggah berkas.');
    } finally {
      setIsUploadingDocs(false);
      e.target.value = '';
    }
  };

  return (
    <>
      {editModalOpen && editModalType === 'page' ? (
        <div className="bg-white rounded-3xl border border-slate-200 p-8 shadow-sm space-y-6 text-left animate-in fade-in duration-200">
          <div className="flex items-center justify-between pb-4 border-b border-slate-100">
            <div>
              <h2 className="text-base font-extrabold text-[#002147]">
                {activeEditItem ? 'Edit Halaman Publik' : 'Tambah Halaman Publik Baru'}
              </h2>
              <span className="text-[11px] text-slate-400 font-medium block">
                Isi konfigurasi detail konten halaman rujukan PPID.
              </span>
            </div>
            <button
              type="button"
              onClick={() => setEditModalOpen(false)}
              className="px-4 py-2 bg-slate-100 hover:bg-slate-200 rounded-xl text-xs font-bold uppercase transition-all cursor-pointer flex items-center gap-1.5 border border-slate-200 text-slate-600"
            >
              <ChevronLeft className="h-4 w-4" /> Kembali
            </button>
          </div>

          <form onSubmit={handleSaveCrudItem} className="space-y-4">
            <div className="space-y-1">
              <label className="text-[10px] font-bold text-slate-600 uppercase tracking-wider block">
                Judul Halaman <span className="text-red-500">*</span>
              </label>
              <input
                type="text"
                required
                value={adminEditTitle}
                onChange={(e) => setAdminEditTitle(e.target.value)}
                className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 transition-all font-bold text-slate-805"
              />
            </div>

            <div className="space-y-1">
              <label className="text-[10px] font-bold text-slate-600 uppercase tracking-wider block">Subjudul Halaman</label>
              <input
                type="text"
                value={adminEditSubtitle}
                onChange={(e) => setAdminEditSubtitle(e.target.value)}
                className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 transition-all font-bold text-slate-805"
              />
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-1">
                <label className="text-[10px] font-bold text-slate-600 uppercase tracking-wider block">Slug Rute Rujukan</label>
                <input
                  type="text"
                  placeholder="contoh-slug-halaman"
                  value={adminEditSlug}
                  onChange={(e) => setAdminEditSlug(e.target.value)}
                  className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 transition-all font-mono font-medium text-slate-805"
                />
              </div>
              <div className="space-y-1">
                <label className="text-[10px] font-bold text-slate-600 uppercase tracking-wider block">File Cover Image URL</label>
                <div className="flex gap-2">
                  <input
                    type="text"
                    placeholder="/uploads/..."
                    value={adminEditCoverImage}
                    onChange={(e) => setAdminEditCoverImage(e.target.value)}
                    className="flex-1 rounded-xl border border-slate-200 px-3 py-2.5 text-xs focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 transition-all font-mono font-medium text-slate-805"
                  />
                  <label className="px-3.5 py-2 bg-slate-100 hover:bg-slate-200 hover:text-slate-800 rounded-xl text-xs font-bold transition-all cursor-pointer shrink-0 inline-flex items-center border border-slate-200 text-slate-600">
                    Upload
                    <input
                      type="file"
                      onChange={handleAdminUpload}
                      className="hidden"
                    />
                  </label>
                </div>
              </div>
            </div>

            {/* Profil Page: Stats Cards Editor */}
            {adminEditSlug === 'profil' && (
              <div className="space-y-4 pt-2 border-t border-slate-100 text-left">
                <div className="flex items-center justify-between pb-2 border-b border-slate-100">
                  <div>
                    <span className="text-[10px] font-extrabold text-[#002147] uppercase tracking-wider block">
                      Kartu Statistik Profil
                    </span>
                    <span className="text-[11px] text-slate-400 font-medium block">
                      Kartu angka/statistik yang tampil di samping foto halaman Profil (contoh: 2021 · Tahun Berdiri).
                    </span>
                  </div>
                  <button
                    type="button"
                    onClick={() => setEditProfilStats([...editProfilStats, { value: '', label: '' }])}
                    className="px-3.5 py-2 bg-[#002147] hover:bg-[#003166] text-white rounded-xl text-[10px] font-bold uppercase tracking-wider inline-flex items-center gap-1 cursor-pointer border border-[#002147]/55"
                  >
                    <Plus className="h-3.5 w-3.5" /> Tambah Kartu
                  </button>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                  {editProfilStats.length > 0 ? (
                    editProfilStats.map((stat, idx) => (
                      <div key={idx} className="p-4 border border-slate-200 rounded-2xl bg-slate-50/50 space-y-2 relative">
                        <button
                          type="button"
                          onClick={() => setEditProfilStats(editProfilStats.filter((_, i) => i !== idx))}
                          className="absolute top-3 right-3 p-1.5 text-slate-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-all"
                          title="Hapus Kartu"
                        >
                          <Trash2 className="h-4 w-4" />
                        </button>
                        <div className="space-y-1 pr-8">
                          <label className="text-[9px] font-bold text-slate-500 uppercase block">Nilai / Angka</label>
                          <input
                            type="text"
                            value={stat.value}
                            onChange={(e) => {
                              const updated = [...editProfilStats];
                              updated[idx] = { ...updated[idx], value: e.target.value };
                              setEditProfilStats(updated);
                            }}
                            placeholder="Contoh: 2021"
                            className="w-full rounded-xl border border-slate-200 px-3 py-2 text-xs bg-white font-black text-slate-800 focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600"
                          />
                        </div>
                        <div className="space-y-1">
                          <label className="text-[9px] font-bold text-slate-500 uppercase block">Label / Keterangan</label>
                          <input
                            type="text"
                            value={stat.label}
                            onChange={(e) => {
                              const updated = [...editProfilStats];
                              updated[idx] = { ...updated[idx], label: e.target.value };
                              setEditProfilStats(updated);
                            }}
                            placeholder="Contoh: TAHUN BERDIRI"
                            className="w-full rounded-xl border border-slate-200 px-3 py-2 text-xs bg-white font-semibold text-slate-700 focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600"
                          />
                        </div>
                      </div>
                    ))
                  ) : (
                    <div className="col-span-2 py-8 text-center text-slate-400 text-xs font-medium border border-dashed border-slate-200 rounded-2xl">
                      Belum ada kartu statistik. Klik "Tambah Kartu" di atas.
                    </div>
                  )}
                </div>

                {/* Sejarah text area */}
                <div className="space-y-1 pt-2 border-t border-slate-100">
                  <label className="text-[10px] font-bold text-slate-600 uppercase tracking-wider block">Teks Sejarah Singkat</label>
                  <textarea
                    value={adminEditContent}
                    onChange={(e) => setAdminEditContent(e.target.value)}
                    placeholder="Tulis sejarah singkat PPID UPERTIS di sini..."
                    className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 transition-all bg-slate-50 h-24 font-medium resize-none text-slate-805"
                  />
                </div>

                {/* Asymmetrical stacked images section */}
                <div className="space-y-3 pt-4 border-t border-slate-100">
                  <div>
                    <span className="text-[10px] font-extrabold text-[#002147] uppercase tracking-wider block">
                      Unggah Multi-Image Sejarah (Tata Letak Bertumpuk)
                    </span>
                    <span className="text-[11px] text-slate-400 font-medium block mt-0.5">
                      Unggah tiga gambar asimetris bertumpuk untuk menghias bagian Sejarah Singkat pada halaman Profil.
                    </span>
                  </div>

                  <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                    {/* Gambar 1 (Vertikal) */}
                    <div className="border border-slate-200 rounded-2xl p-4 bg-slate-50/50 space-y-2">
                      <label className="text-[9px] font-bold text-slate-500 uppercase block">1. Gambar Kiri (Vertikal)</label>
                      <div className="flex gap-2">
                        <input
                          type="text"
                          placeholder="/uploads/..."
                          value={editProfilImage1}
                          onChange={(e) => setEditProfilImage1(e.target.value)}
                          className="flex-1 rounded-xl border border-slate-200 px-2.5 py-1.5 text-[10px] focus:outline-none bg-white font-mono"
                        />
                        <label className="px-3 py-1.5 bg-slate-200 hover:bg-slate-300 rounded-xl text-[10px] font-bold cursor-pointer inline-flex items-center text-slate-700">
                          Upload
                          <input
                            type="file"
                            onChange={async (e) => {
                              const file = e.target.files?.[0];
                              if (!file) return;
                              const formData = new FormData();
                              formData.append('files', file);
                              const uploadRes = await fetch(`${API_BASE_URL}/admin/uploads`, {
                                method: 'POST',
                                headers: { Authorization: `Bearer ${localStorage.getItem('auth_token')}` },
                                body: formData
                              });
                              const resData = await uploadRes.json();
                              const uploadedUrl = resData.url || resData.path || (resData.urls && resData.urls[0]) || (resData.uploaded && resData.uploaded[0] && resData.uploaded[0].url) || '';
                              if (uploadRes.ok && uploadedUrl) setEditProfilImage1(uploadedUrl);
                            }}
                            className="hidden"
                          />
                        </label>
                      </div>
                      {editProfilImage1 && (
                        <img src={resolveImageUrl(editProfilImage1)} className="h-16 w-full object-cover rounded-lg border border-slate-200 mt-1" />
                      )}
                    </div>

                    {/* Gambar 2 (Persegi) */}
                    <div className="border border-slate-200 rounded-2xl p-4 bg-slate-50/50 space-y-2">
                      <label className="text-[9px] font-bold text-slate-500 uppercase block">2. Gambar Kanan Atas (Persegi)</label>
                      <div className="flex gap-2">
                        <input
                          type="text"
                          placeholder="/uploads/..."
                          value={editProfilImage2}
                          onChange={(e) => setEditProfilImage2(e.target.value)}
                          className="flex-1 rounded-xl border border-slate-200 px-2.5 py-1.5 text-[10px] focus:outline-none bg-white font-mono"
                        />
                        <label className="px-3 py-1.5 bg-slate-200 hover:bg-slate-300 rounded-xl text-[10px] font-bold cursor-pointer inline-flex items-center text-slate-700">
                          Upload
                          <input
                            type="file"
                            onChange={async (e) => {
                              const file = e.target.files?.[0];
                              if (!file) return;
                              const formData = new FormData();
                              formData.append('files', file);
                              const uploadRes = await fetch(`${API_BASE_URL}/admin/uploads`, {
                                method: 'POST',
                                headers: { Authorization: `Bearer ${localStorage.getItem('auth_token')}` },
                                body: formData
                              });
                              const resData = await uploadRes.json();
                              const uploadedUrl = resData.url || resData.path || (resData.urls && resData.urls[0]) || (resData.uploaded && resData.uploaded[0] && resData.uploaded[0].url) || '';
                              if (uploadRes.ok && uploadedUrl) setEditProfilImage2(uploadedUrl);
                            }}
                            className="hidden"
                          />
                        </label>
                      </div>
                      {editProfilImage2 && (
                        <img src={resolveImageUrl(editProfilImage2)} className="h-16 w-full object-cover rounded-lg border border-slate-200 mt-1" />
                      )}
                    </div>

                    {/* Gambar 3 (Horizontal) */}
                    <div className="border border-slate-200 rounded-2xl p-4 bg-slate-50/50 space-y-2">
                      <label className="text-[9px] font-bold text-slate-500 uppercase block">3. Gambar Tengah Depan</label>
                      <div className="flex gap-2">
                        <input
                          type="text"
                          placeholder="/uploads/..."
                          value={editProfilImage3}
                          onChange={(e) => setEditProfilImage3(e.target.value)}
                          className="flex-1 rounded-xl border border-slate-200 px-2.5 py-1.5 text-[10px] focus:outline-none bg-white font-mono"
                        />
                        <label className="px-3 py-1.5 bg-slate-200 hover:bg-slate-300 rounded-xl text-[10px] font-bold cursor-pointer inline-flex items-center text-slate-700">
                          Upload
                          <input
                            type="file"
                            onChange={async (e) => {
                              const file = e.target.files?.[0];
                              if (!file) return;
                              const formData = new FormData();
                              formData.append('files', file);
                              const uploadRes = await fetch(`${API_BASE_URL}/admin/uploads`, {
                                method: 'POST',
                                headers: { Authorization: `Bearer ${localStorage.getItem('auth_token')}` },
                                body: formData
                              });
                              const resData = await uploadRes.json();
                              const uploadedUrl = resData.url || resData.path || (resData.urls && resData.urls[0]) || (resData.uploaded && resData.uploaded[0] && resData.uploaded[0].url) || '';
                              if (uploadRes.ok && uploadedUrl) setEditProfilImage3(uploadedUrl);
                            }}
                            className="hidden"
                          />
                        </label>
                      </div>
                      {editProfilImage3 && (
                        <img src={resolveImageUrl(editProfilImage3)} className="h-16 w-full object-cover rounded-lg border border-slate-200 mt-1" />
                      )}
                    </div>
                  </div>
                </div>
              </div>
            )}

            {/* --- EDITOR TUGAS & FUNGSI PPID --- */}
            {['tugas-dan-fungsi', 'tugas-fungsi'].includes(adminEditSlug) && (
              <div className="space-y-6 pt-2 border-t border-slate-100 text-left">
                <div>
                  <span className="text-[10px] font-extrabold text-[#002147] uppercase tracking-wider block">
                    PENGATURAN KARTU TUGAS & FUNGSI PPID
                  </span>
                  <span className="text-[11px] text-slate-400 font-medium block mt-0.5">
                    Modifikasi judul, penjelasan kartu tugas asimetris, dan pilar fungsi beserta sub-item keterangannya.
                  </span>
                </div>

                <div className="space-y-3">
                  <div className="flex items-center justify-between border-b border-slate-100 pb-2">
                    <span className="text-[10px] font-black text-slate-700 uppercase tracking-widest">
                      A. Tugas Pokok PPID (Daftar Kartu Asimetris)
                    </span>
                    <button
                      type="button"
                      onClick={() => setEditTugasList([...editTugasList, { title: 'Tugas Baru', desc: 'Penjelasan detail tugas baru...' }])}
                      className="px-3 py-1.5 bg-blue-50 hover:bg-blue-100 text-blue-700 rounded-xl text-[10px] font-bold transition-all border border-blue-200 cursor-pointer"
                    >
                      + Tambah Kartu Tugas
                    </button>
                  </div>

                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    {editTugasList.map((tugas, idx) => (
                      <div key={idx} className="p-4 bg-slate-50 border border-slate-200 rounded-2xl space-y-3 relative group">
                        <button
                          type="button"
                          onClick={() => setEditTugasList(editTugasList.filter((_, i) => i !== idx))}
                          className="absolute top-2 right-2 text-rose-500 hover:text-rose-700 text-[10px] font-bold bg-rose-50 hover:bg-rose-100 p-1.5 rounded-lg border border-rose-200 cursor-pointer"
                        >
                          Hapus
                        </button>
                        <div>
                          <label className="text-[8px] font-extrabold text-slate-400 uppercase block mb-1">Judul Tugas 0{idx + 1}</label>
                          <input
                            type="text"
                            value={tugas.title}
                            onChange={(e) => {
                              const updated = [...editTugasList];
                              updated[idx].title = e.target.value;
                              setEditTugasList(updated);
                            }}
                            className="w-full rounded-xl border border-slate-200 px-3 py-1.5 text-xs bg-white font-semibold text-slate-705 focus:outline-none"
                          />
                        </div>
                        <div>
                          <label className="text-[8px] font-extrabold text-slate-400 uppercase block mb-1">Deskripsi Penjelasan</label>
                          <textarea
                            value={tugas.desc}
                            onChange={(e) => {
                              const updated = [...editTugasList];
                              updated[idx].desc = e.target.value;
                              setEditTugasList(updated);
                            }}
                            className="w-full rounded-xl border border-slate-200 px-3 py-2 text-[11px] bg-white font-medium text-slate-500 focus:outline-none h-16 resize-none"
                          />
                        </div>
                      </div>
                    ))}
                  </div>
                </div>

                {/* Section B: 5 Fungsi PPID */}
                <div className="space-y-4 pt-4 border-t border-slate-100">
                  <div className="flex items-center justify-between border-b border-slate-100 pb-2">
                    <span className="text-[10px] font-black text-slate-700 uppercase tracking-widest">
                      B. Fungsi Strategis PPID (Daftar Accordion)
                    </span>
                    <button
                      type="button"
                      onClick={() => setEditFungsiList([...editFungsiList, { title: 'Fungsi Strategis Baru', items: ['Detail fungsi poin 1'] }])}
                      className="px-3 py-1.5 bg-blue-50 hover:bg-blue-100 text-blue-700 rounded-xl text-[10px] font-bold transition-all border border-blue-200 cursor-pointer"
                    >
                      + Tambah Pilar Fungsi
                    </button>
                  </div>

                  <div className="space-y-4">
                    {editFungsiList.map((fungsi, idx) => (
                      <div key={idx} className="p-5 bg-white border border-slate-200 rounded-2xl space-y-3 relative shadow-sm">
                        <button
                          type="button"
                          onClick={() => setEditFungsiList(editFungsiList.filter((_, i) => i !== idx))}
                          className="absolute top-3 right-3 text-rose-500 hover:text-rose-700 text-[10px] font-bold bg-rose-50 hover:bg-rose-100 p-1.5 rounded-lg border border-rose-200 cursor-pointer"
                        >
                          Hapus Pilar
                        </button>
                        
                        <div>
                          <label className="text-[8px] font-extrabold text-slate-400 uppercase block mb-1">Judul Pilar Fungsi 0{idx + 1}</label>
                          <input
                            type="text"
                            value={fungsi.title}
                            onChange={(e) => {
                              const updated = [...editFungsiList];
                              updated[idx].title = e.target.value;
                              setEditFungsiList(updated);
                            }}
                            className="w-2/3 rounded-xl border border-slate-200 px-3 py-1.5 text-xs bg-slate-50 font-bold text-slate-800 focus:outline-none"
                          />
                        </div>

                        {/* List items inside each Fungsi */}
                        <div className="space-y-2 pt-2 border-t border-slate-100/60">
                          <div className="flex items-center justify-between">
                            <label className="text-[8px] font-extrabold text-[#002147] uppercase tracking-wider block">Daftar Sub-Poin (Checklist)</label>
                            <button
                              type="button"
                              onClick={() => {
                                const updated = [...editFungsiList];
                                updated[idx].items.push('Deskripsi checklist baru');
                                setEditFungsiList(updated);
                              }}
                              className="text-[9px] font-bold text-blue-600 hover:underline cursor-pointer border-0 bg-transparent"
                            >
                              + Tambah Sub-Poin
                            </button>
                          </div>
                          
                          <div className="space-y-2">
                            {fungsi.items.map((itemStr, sIdx) => (
                              <div key={sIdx} className="flex gap-2 items-center">
                                <span className="text-[10px] font-bold text-slate-400 font-mono w-4 shrink-0 text-right">{sIdx + 1}.</span>
                                <input
                                  type="text"
                                  value={itemStr}
                                  onChange={(e) => {
                                    const updated = [...editFungsiList];
                                    updated[idx].items[sIdx] = e.target.value;
                                    setEditFungsiList(updated);
                                  }}
                                  className="flex-1 rounded-lg border border-slate-200 px-2 py-1 text-[11px] font-medium text-slate-600 focus:outline-none"
                                />
                                <button
                                  type="button"
                                  onClick={() => {
                                    const updated = [...editFungsiList];
                                    updated[idx].items = updated[idx].items.filter((_, i) => i !== sIdx);
                                    setEditFungsiList(updated);
                                  }}
                                  className="text-rose-500 hover:text-rose-700 text-xs px-2 border-0 bg-transparent cursor-pointer font-bold"
                                >
                                  ×
                                </button>
                              </div>
                            ))}
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>

                {/* Section C: Landasan Hukum */}
                <div className="space-y-1 pt-4 border-t border-slate-100">
                  <label className="text-[10px] font-bold text-slate-600 uppercase tracking-wider block">Landasan Hukum & Keterangan Tambahan (Rich Text / HTML)</label>
                  <textarea
                    value={adminEditContent}
                    onChange={(e) => setAdminEditContent(e.target.value)}
                    placeholder="Masukkan landasan hukum atau teks tambahan di sini..."
                    className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none bg-slate-50 h-28 font-medium text-slate-805"
                  />
                </div>

                {/* Section D: Upload Gambar & Multi-Image */}
                <div className="space-y-4 pt-4 border-t border-slate-100">
                  <div>
                    <span className="text-[10px] font-bold text-slate-600 uppercase tracking-wider block">
                      D. Upload Gambar & Multi-Image (Tumpukan Asimetris)
                    </span>
                    <span className="text-[11px] text-slate-400 font-medium block mt-0.5">
                      Unggah 1 gambar utama (kiri) dan 2 gambar pendukung (kanan atas & tengah depan) untuk ditampilkan bertumpuk asimetris di sebelah kanan landasan hukum.
                    </span>
                  </div>

                  <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                    {/* Gambar 1 (Vertikal) */}
                    <div className="border border-slate-200 rounded-2xl p-4 bg-slate-50/50 space-y-2">
                      <label className="text-[9px] font-bold text-slate-500 uppercase block">1. Gambar Kiri (Vertikal)</label>
                      <div className="flex gap-2">
                        <input
                          type="text"
                          placeholder="/uploads/..."
                          value={editTugasImage1}
                          onChange={(e) => setEditTugasImage1(e.target.value)}
                          className="flex-1 rounded-xl border border-slate-200 px-2.5 py-1.5 text-[10px] focus:outline-none bg-white font-mono"
                        />
                        <label className="px-3 py-1.5 bg-slate-200 hover:bg-slate-300 rounded-xl text-[10px] font-bold cursor-pointer inline-flex items-center text-slate-700">
                          Upload
                          <input
                            type="file"
                            onChange={async (e) => {
                              const file = e.target.files?.[0];
                              if (!file) return;
                              const formData = new FormData();
                              formData.append('files', file);
                              const uploadRes = await fetch(`${API_BASE_URL}/admin/uploads`, {
                                method: 'POST',
                                headers: { Authorization: `Bearer ${localStorage.getItem('auth_token')}` },
                                body: formData
                              });
                              const resData = await uploadRes.json();
                              const uploadedUrl = resData.url || resData.path || (resData.urls && resData.urls[0]) || (resData.uploaded && resData.uploaded[0] && resData.uploaded[0].url) || '';
                              if (uploadRes.ok && uploadedUrl) setEditTugasImage1(uploadedUrl);
                            }}
                            className="hidden"
                          />
                        </label>
                      </div>
                      {editTugasImage1 && (
                        <img src={resolveImageUrl(editTugasImage1)} className="h-16 w-full object-cover rounded-lg border border-slate-200 mt-1" />
                      )}
                    </div>

                    {/* Gambar 2 (Persegi) */}
                    <div className="border border-slate-200 rounded-2xl p-4 bg-slate-50/50 space-y-2">
                      <label className="text-[9px] font-bold text-slate-500 uppercase block">2. Gambar Kanan Atas (Persegi)</label>
                      <div className="flex gap-2">
                        <input
                          type="text"
                          placeholder="/uploads/..."
                          value={editTugasImage2}
                          onChange={(e) => setEditTugasImage2(e.target.value)}
                          className="flex-1 rounded-xl border border-slate-200 px-2.5 py-1.5 text-[10px] focus:outline-none bg-white font-mono"
                        />
                        <label className="px-3 py-1.5 bg-slate-200 hover:bg-slate-300 rounded-xl text-[10px] font-bold cursor-pointer inline-flex items-center text-slate-700">
                          Upload
                          <input
                            type="file"
                            onChange={async (e) => {
                              const file = e.target.files?.[0];
                              if (!file) return;
                              const formData = new FormData();
                              formData.append('files', file);
                              const uploadRes = await fetch(`${API_BASE_URL}/admin/uploads`, {
                                method: 'POST',
                                headers: { Authorization: `Bearer ${localStorage.getItem('auth_token')}` },
                                body: formData
                              });
                              const resData = await uploadRes.json();
                              const uploadedUrl = resData.url || resData.path || (resData.urls && resData.urls[0]) || (resData.uploaded && resData.uploaded[0] && resData.uploaded[0].url) || '';
                              if (uploadRes.ok && uploadedUrl) setEditTugasImage2(uploadedUrl);
                            }}
                            className="hidden"
                          />
                        </label>
                      </div>
                      {editTugasImage2 && (
                        <img src={resolveImageUrl(editTugasImage2)} className="h-16 w-full object-cover rounded-lg border border-slate-200 mt-1" />
                      )}
                    </div>

                    {/* Gambar 3 (Horizontal) */}
                    <div className="border border-slate-200 rounded-2xl p-4 bg-slate-50/50 space-y-2">
                      <label className="text-[9px] font-bold text-slate-500 uppercase block">3. Gambar Tengah Depan</label>
                      <div className="flex gap-2">
                        <input
                          type="text"
                          placeholder="/uploads/..."
                          value={editTugasImage3}
                          onChange={(e) => setEditTugasImage3(e.target.value)}
                          className="flex-1 rounded-xl border border-slate-200 px-2.5 py-1.5 text-[10px] focus:outline-none bg-white font-mono"
                        />
                        <label className="px-3 py-1.5 bg-slate-200 hover:bg-slate-300 rounded-xl text-[10px] font-bold cursor-pointer inline-flex items-center text-slate-700">
                          Upload
                          <input
                            type="file"
                            onChange={async (e) => {
                              const file = e.target.files?.[0];
                              if (!file) return;
                              const formData = new FormData();
                              formData.append('files', file);
                              const uploadRes = await fetch(`${API_BASE_URL}/admin/uploads`, {
                                method: 'POST',
                                headers: { Authorization: `Bearer ${localStorage.getItem('auth_token')}` },
                                body: formData
                              });
                              const resData = await uploadRes.json();
                              const uploadedUrl = resData.url || resData.path || (resData.urls && resData.urls[0]) || (resData.uploaded && resData.uploaded[0] && resData.uploaded[0].url) || '';
                              if (uploadRes.ok && uploadedUrl) setEditTugasImage3(uploadedUrl);
                            }}
                            className="hidden"
                          />
                        </label>
                      </div>
                      {editTugasImage3 && (
                        <img src={resolveImageUrl(editTugasImage3)} className="h-16 w-full object-cover rounded-lg border border-slate-200 mt-1" />
                      )}
                    </div>
                  </div>
                </div>
              </div>
            )}

            {/* --- EDITOR BAGAN STRUKTUR ORGANISASI PPID --- */}
            {adminEditSlug === 'struktur-organisasi-2' && (
              <div className="space-y-6 pt-2 border-t border-slate-100 text-left">
                <div>
                  <span className="text-[10px] font-extrabold text-[#002147] uppercase tracking-wider block">
                    PENGATURAN BAGAN STRUKTUR ORGANISASI PPID
                  </span>
                  <span className="text-[11px] text-slate-400 font-medium block mt-0.5">
                    Modifikasi bagan keanggotaan dan alur koordinasi Pejabat Pengelola Informasi.
                  </span>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div className="space-y-1 text-left">
                    <label className="text-[10px] font-bold text-slate-600 uppercase tracking-wider block">Atasan PPID</label>
                    <input
                      type="text"
                      required
                      value={strukturAtasan}
                      onChange={(e) => updateStrukturContent({ atasan: e.target.value })}
                      className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 transition-all bg-slate-50 font-bold text-slate-805"
                      placeholder="Atasan PPID (contoh: Rektor Universitas Perintis)"
                    />
                  </div>
                  <div className="space-y-1 text-left">
                    <label className="text-[10px] font-bold text-slate-600 uppercase tracking-wider block">PPID Utama</label>
                    <input
                      type="text"
                      required
                      value={strukturUtama}
                      onChange={(e) => updateStrukturContent({ utama: e.target.value })}
                      className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 transition-all bg-slate-50 font-bold text-slate-805"
                      placeholder="PPID Utama (contoh: Wakil Rektor 1 & 2 UPERTIS)"
                    />
                  </div>
                </div>

                {/* Column 1: PPID Pelaksana */}
                <div className="space-y-3 p-4 border border-slate-200 rounded-2xl bg-slate-50/50">
                  <div className="flex items-center justify-between border-b border-slate-100 pb-2">
                    <span className="text-[10px] font-black text-blue-700 uppercase tracking-widest">
                      1. PPID Pelaksana (Biro & Dekan)
                    </span>
                    <button
                      type="button"
                      onClick={() => updateStrukturContent({ pelaksana: [...strukturPelaksana, ''] })}
                      className="text-[9px] font-bold text-blue-600 hover:underline cursor-pointer border-0 bg-transparent font-sans"
                    >
                      + Tambah Item
                    </button>
                  </div>
                  <div className="space-y-2">
                    {strukturPelaksana.map((item, idx) => (
                      <div key={idx} className="flex gap-2 items-center">
                        <span className="text-[10px] font-bold text-slate-400 font-mono w-4 shrink-0 text-right">{idx + 1}.</span>
                        <input
                          type="text"
                          required
                          value={item}
                          onChange={(e) => {
                            const updated = [...strukturPelaksana];
                            updated[idx] = e.target.value;
                            updateStrukturContent({ pelaksana: updated });
                          }}
                          className="flex-1 rounded-lg border border-slate-200 px-3 py-1.5 text-xs focus:outline-none bg-white font-semibold text-slate-700"
                        />
                        <button
                          type="button"
                          onClick={() => {
                            const updated = strukturPelaksana.filter((_, i) => i !== idx);
                            updateStrukturContent({ pelaksana: updated });
                          }}
                          className="text-rose-500 hover:text-rose-700 text-xs px-2 border-0 bg-transparent cursor-pointer font-bold"
                        >
                          ×
                        </button>
                      </div>
                    ))}
                    {strukturPelaksana.length === 0 && (
                      <p className="text-[11px] text-slate-400 italic">Belum ada item ditambahkan.</p>
                    )}
                  </div>
                </div>

                {/* Column 2: Tim Pertimbangan */}
                <div className="space-y-3 p-4 border border-slate-200 rounded-2xl bg-slate-50/50">
                  <div className="flex items-center justify-between border-b border-slate-100 pb-2">
                    <span className="text-[10px] font-black text-amber-700 uppercase tracking-widest">
                      2. Tim Pertimbangan (Komite Pertimbangan)
                    </span>
                    <button
                      type="button"
                      onClick={() => updateStrukturContent({ pertimbangan: [...strukturPertimbangan, ''] })}
                      className="text-[9px] font-bold text-amber-600 hover:underline cursor-pointer border-0 bg-transparent font-sans"
                    >
                      + Tambah Item
                    </button>
                  </div>
                  <div className="space-y-2">
                    {strukturPertimbangan.map((item, idx) => (
                      <div key={idx} className="flex gap-2 items-center">
                        <span className="text-[10px] font-bold text-slate-400 font-mono w-4 shrink-0 text-right">{idx + 1}.</span>
                        <input
                          type="text"
                          required
                          value={item}
                          onChange={(e) => {
                            const updated = [...strukturPertimbangan];
                            updated[idx] = e.target.value;
                            updateStrukturContent({ pertimbangan: updated });
                          }}
                          className="flex-1 rounded-lg border border-slate-200 px-3 py-1.5 text-xs focus:outline-none bg-white font-semibold text-slate-700"
                        />
                        <button
                          type="button"
                          onClick={() => {
                            const updated = strukturPertimbangan.filter((_, i) => i !== idx);
                            updateStrukturContent({ pertimbangan: updated });
                          }}
                          className="text-rose-500 hover:text-rose-700 text-xs px-2 border-0 bg-transparent cursor-pointer font-bold"
                        >
                          ×
                        </button>
                      </div>
                    ))}
                    {strukturPertimbangan.length === 0 && (
                      <p className="text-[11px] text-slate-400 italic">Belum ada item ditambahkan.</p>
                    )}
                  </div>
                </div>

                {/* Column 3: Petugas Pelayanan */}
                <div className="space-y-3 p-4 border border-slate-200 rounded-2xl bg-slate-50/50">
                  <div className="flex items-center justify-between border-b border-slate-100 pb-2">
                    <span className="text-[10px] font-black text-emerald-700 uppercase tracking-widest">
                      3. Petugas Pelayanan (Staf Desk Pelayanan)
                    </span>
                    <button
                      type="button"
                      onClick={() => updateStrukturContent({ pelayanan: [...strukturPelayanan, ''] })}
                      className="text-[9px] font-bold text-emerald-600 hover:underline cursor-pointer border-0 bg-transparent font-sans"
                    >
                      + Tambah Item
                    </button>
                  </div>
                  <div className="space-y-2">
                    {strukturPelayanan.map((item, idx) => (
                      <div key={idx} className="flex gap-2 items-center">
                        <span className="text-[10px] font-bold text-slate-400 font-mono w-4 shrink-0 text-right">{idx + 1}.</span>
                        <input
                          type="text"
                          required
                          value={item}
                          onChange={(e) => {
                            const updated = [...strukturPelayanan];
                            updated[idx] = e.target.value;
                            updateStrukturContent({ pelayanan: updated });
                          }}
                          className="flex-1 rounded-lg border border-slate-200 px-3 py-1.5 text-xs focus:outline-none bg-white font-semibold text-slate-700"
                        />
                        <button
                          type="button"
                          onClick={() => {
                            const updated = strukturPelayanan.filter((_, i) => i !== idx);
                            updateStrukturContent({ pelayanan: updated });
                          }}
                          className="text-rose-500 hover:text-rose-700 text-xs px-2 border-0 bg-transparent cursor-pointer font-bold"
                        >
                          ×
                        </button>
                      </div>
                    ))}
                    {strukturPelayanan.length === 0 && (
                      <p className="text-[11px] text-slate-400 italic">Belum ada item ditambahkan.</p>
                    )}
                  </div>
                </div>

                {/* Section Description */}
                <div className="space-y-1">
                  <label className="text-[10px] font-bold text-slate-600 uppercase tracking-wider block">Penjelasan Tambahan (HTML/Markdown)</label>
                  <textarea
                    value={strukturDesc}
                    onChange={(e) => updateStrukturContent({ desc: e.target.value })}
                    placeholder="Masukkan penjelasan tambahan di bawah bagan jika ada..."
                    className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none bg-slate-50 h-28 font-medium text-slate-805"
                  />
                </div>
              </div>
            )}

            {/* --- EDITOR VISI & MISI PPID --- */}
            {adminEditSlug === 'visi-misi' && (
              <div className="space-y-6 pt-2 border-t border-slate-100 text-left animate-in fade-in duration-200">
                <div>
                  <span className="text-[10px] font-extrabold text-[#002147] uppercase tracking-wider block">
                    PENGATURAN VISI & MISI PPID
                  </span>
                  <span className="text-[11px] text-slate-400 font-medium block mt-0.5">
                    Modifikasi teks pernyataan Visi dan daftar poin Misi strategis PPID.
                  </span>
                </div>

                <div className="space-y-1 text-left">
                  <label className="text-[10px] font-bold text-slate-600 uppercase tracking-wider block">Pernyataan Visi PPID</label>
                  <textarea
                    required
                    value={visiText}
                    onChange={(e) => updateVisiMisiContent({ visi: e.target.value })}
                    className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 transition-all bg-slate-50 font-serif italic text-slate-805 h-20 resize-none font-medium"
                    placeholder="Masukkan teks visi PPID..."
                  />
                </div>

                <div className="space-y-3 p-4 border border-slate-200 rounded-2xl bg-slate-50/50">
                  <div className="flex items-center justify-between border-b border-slate-100 pb-2">
                    <span className="text-[10px] font-black text-blue-700 uppercase tracking-widest">
                      Daftar Misi PPID (Checklist Poin)
                    </span>
                    <button
                      type="button"
                      onClick={() => updateVisiMisiContent({ misi: [...misiList, ''] })}
                      className="text-[9px] font-bold text-blue-600 hover:underline cursor-pointer border-0 bg-transparent font-sans"
                    >
                      + Tambah Poin Misi
                    </button>
                  </div>
                  <div className="space-y-2">
                    {misiList.map((item, idx) => (
                      <div key={idx} className="flex gap-2 items-start">
                        <span className="text-[10px] font-bold text-slate-400 font-mono w-4 shrink-0 text-right mt-2">{idx + 1}.</span>
                        <textarea
                          required
                          value={item}
                          onChange={(e) => {
                            const updated = [...misiList];
                            updated[idx] = e.target.value;
                            updateVisiMisiContent({ misi: updated });
                          }}
                          className="flex-1 rounded-lg border border-slate-200 px-3 py-1.5 text-xs focus:outline-none bg-white font-semibold text-slate-700 h-14 resize-none"
                        />
                        <button
                          type="button"
                          onClick={() => {
                            const updated = misiList.filter((_, i) => i !== idx);
                            updateVisiMisiContent({ misi: updated });
                          }}
                          className="text-rose-500 hover:text-rose-700 text-xs px-2 border-0 bg-transparent cursor-pointer font-bold mt-2"
                        >
                          ×
                        </button>
                      </div>
                    ))}
                    {misiList.length === 0 && (
                      <p className="text-[11px] text-slate-400 italic">Belum ada item ditambahkan.</p>
                    )}
                  </div>
                </div>
              </div>
            )}

            {adminEditSlug === 'jadwal-layanan-informasi' ? (
              <div className="space-y-4 pt-2 border-t border-slate-100 text-left">
                <span className="text-[10px] font-extrabold text-[#002147] uppercase tracking-wider block mb-2">
                  PENGATURAN JADWAL OPERASIONAL (SLIDER STYLE)
                </span>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div className="space-y-1 text-left">
                    <label className="text-[10px] font-bold text-slate-600 uppercase tracking-wider block">Jam Kerja (Senin - Kamis)</label>
                    <input
                      type="text"
                      value={editJadwalSeninKamisKerja}
                      onChange={(e) => setEditJadwalSeninKamisKerja(e.target.value)}
                      className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 transition-all bg-slate-50 font-medium text-slate-805"
                    />
                  </div>
                  <div className="space-y-1 text-left">
                    <label className="text-[10px] font-bold text-slate-600 uppercase tracking-wider block">Jam Istirahat (Senin - Kamis)</label>
                    <input
                      type="text"
                      value={editJadwalSeninKamisIstirahat}
                      onChange={(e) => setEditJadwalSeninKamisIstirahat(e.target.value)}
                      className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 transition-all bg-slate-50 font-medium text-slate-805"
                    />
                  </div>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div className="space-y-1 text-left">
                    <label className="text-[10px] font-bold text-slate-600 uppercase tracking-wider block">Jam Kerja (Hari Jumat)</label>
                    <input
                      type="text"
                      value={editJadwalJumatKerja}
                      onChange={(e) => setEditJadwalJumatKerja(e.target.value)}
                      className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 transition-all bg-slate-50 font-medium text-slate-805"
                    />
                  </div>
                  <div className="space-y-1 text-left">
                    <label className="text-[10px] font-bold text-slate-600 uppercase tracking-wider block">Jam Istirahat (Hari Jumat)</label>
                    <input
                      type="text"
                      value={editJadwalPageIstirahatJumat}
                      onChange={(e) => setEditJadwalPageIstirahatJumat(e.target.value)}
                      className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 transition-all bg-slate-50 font-medium text-slate-805"
                    />
                  </div>
                </div>

                <div className="space-y-1 text-left">
                  <label className="text-[10px] font-bold text-slate-600 uppercase tracking-wider block">Status Sabtu - Minggu & Hari Libur</label>
                  <input
                    type="text"
                    value={editJadwalSabtuMinggu}
                    onChange={(e) => setEditJadwalSabtuMinggu(e.target.value)}
                    className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 transition-all bg-slate-50 font-medium text-slate-805"
                  />
                </div>

                <span className="text-[10px] font-extrabold text-[#002147] uppercase tracking-wider block pt-2 border-t border-slate-100">
                  PANDUAN OPERASIONAL
                </span>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
                  <div className="space-y-1 text-left">
                    <label className="text-[9px] font-bold text-slate-600 uppercase tracking-wider block">Panduan Offline 1</label>
                    <textarea
                      value={editJadwalOffline1}
                      onChange={(e) => setEditJadwalOffline1(e.target.value)}
                      className="w-full rounded-xl border border-slate-200 px-3 py-2 text-[11px] focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 transition-all bg-slate-50 h-16 resize-none text-slate-805"
                    />
                  </div>
                  <div className="space-y-1 text-left">
                    <label className="text-[9px] font-bold text-slate-600 uppercase tracking-wider block">Panduan Offline 2</label>
                    <textarea
                      value={editJadwalOffline2}
                      onChange={(e) => setEditJadwalOffline2(e.target.value)}
                      className="w-full rounded-xl border border-slate-200 px-3 py-2 text-[11px] focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 transition-all bg-slate-50 h-16 resize-none text-slate-805"
                    />
                  </div>
                  <div className="space-y-1 text-left">
                    <label className="text-[9px] font-bold text-slate-600 uppercase tracking-wider block">Panduan Offline 3</label>
                    <textarea
                      value={editJadwalOffline3}
                      onChange={(e) => setEditJadwalOffline3(e.target.value)}
                      className="w-full rounded-xl border border-slate-200 px-3 py-2 text-[11px] focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 transition-all bg-slate-50 h-16 resize-none text-slate-805"
                    />
                  </div>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
                  <div className="space-y-1 text-left">
                    <label className="text-[9px] font-bold text-slate-600 uppercase tracking-wider block">Panduan Online 1</label>
                    <textarea
                      value={editJadwalOnline1}
                      onChange={(e) => setEditJadwalOnline1(e.target.value)}
                      className="w-full rounded-xl border border-slate-200 px-3 py-2 text-[11px] focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 transition-all bg-slate-50 h-16 resize-none text-slate-805"
                    />
                  </div>
                  <div className="space-y-1 text-left">
                    <label className="text-[9px] font-bold text-slate-600 uppercase tracking-wider block">Panduan Online 2</label>
                    <textarea
                      value={editJadwalOnline2}
                      onChange={(e) => setEditJadwalOnline2(e.target.value)}
                      className="w-full rounded-xl border border-slate-200 px-3 py-2 text-[11px] focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 transition-all bg-slate-50 h-16 resize-none text-slate-805"
                    />
                  </div>
                  <div className="space-y-1 text-left">
                    <label className="text-[9px] font-bold text-slate-600 uppercase tracking-wider block">Panduan Online 3</label>
                    <textarea
                      value={editJadwalOnline3}
                      onChange={(e) => setEditJadwalOnline3(e.target.value)}
                      className="w-full rounded-xl border border-slate-200 px-3 py-2 text-[11px] focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 transition-all bg-slate-50 h-16 resize-none text-slate-805"
                    />
                  </div>
                </div>

                <div className="space-y-1 pt-2 border-t border-slate-100 text-left">
                  <label className="text-[10px] font-bold text-slate-600 uppercase tracking-wider block">Catatan / Informasi Tambahan</label>
                  <textarea
                    value={editJadwalCustomRemarks}
                    onChange={(e) => setEditJadwalCustomRemarks(e.target.value)}
                    className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 transition-all bg-slate-50 h-24 font-medium resize-none text-slate-805"
                  />
                </div>
              </div>
            ) : adminEditSlug === 'regulasi' ? (
              <div className="space-y-6 pt-2 border-t border-slate-100 text-left">
                <div className="flex items-center justify-between pb-2 border-b border-slate-100">
                  <div>
                    <span className="text-[10px] font-extrabold text-[#002147] uppercase tracking-wider block">
                      Pengaturan Regulasi Accordion
                    </span>
                    <span className="text-[11px] text-slate-400 font-medium block">
                      Kelola grup accordion regulasi dan daftar dokumen/SOP yang ada di dalamnya.
                    </span>
                  </div>
                  <button
                    type="button"
                    onClick={() => {
                      const newGroup = {
                        key: String.fromCharCode(65 + regulasiGroups.length), // A, B, C...
                        label: `${String.fromCharCode(65 + regulasiGroups.length)}. Kategori Baru`,
                        desc: 'Deskripsi kategori regulasi baru.',
                        items: []
                      };
                      updateRegulasiContent([...regulasiGroups, newGroup]);
                    }}
                    className="px-3.5 py-2 bg-[#002147] hover:bg-[#003166] text-white rounded-xl text-[10px] font-bold uppercase tracking-wider inline-flex items-center gap-1 cursor-pointer border border-[#002147]/55"
                  >
                    <Plus className="h-3.5 w-3.5" /> Tambah Grup
                  </button>
                </div>

                <div className="space-y-6 max-h-[600px] overflow-y-auto pr-1">
                  {regulasiGroups.length > 0 ? (
                    regulasiGroups.map((group, gIdx) => (
                      <div key={gIdx} className="p-5 border border-slate-200 rounded-3xl bg-slate-50/50 space-y-4 relative">
                        <div className="absolute top-4 right-4 flex items-center gap-2">
                          <button
                            type="button"
                            onClick={() => {
                              const newList = regulasiGroups.filter((_, i) => i !== gIdx).map((g, i) => ({
                                ...g,
                                key: String.fromCharCode(65 + i),
                                label: g.label.replace(/^[A-Z]\.\s+/, `${String.fromCharCode(65 + i)}. `)
                              }));
                              updateRegulasiContent(newList);
                            }}
                            className="p-1.5 text-slate-400 hover:text-red-650 hover:bg-red-50 rounded-lg transition-all"
                            title="Hapus Grup"
                          >
                            <Trash2 className="h-4 w-4" />
                          </button>
                        </div>

                        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 pr-10">
                          <div className="space-y-1">
                            <label className="text-[9px] font-bold text-slate-500 uppercase block">Key</label>
                            <input
                              type="text"
                              required
                              value={group.key}
                              onChange={(e) => {
                                const newList = regulasiGroups.map((g, i) => i === gIdx ? { ...g, key: e.target.value.toUpperCase() } : g);
                                updateRegulasiContent(newList);
                              }}
                              className="w-full rounded-xl border border-slate-200 px-3 py-2 text-xs font-bold text-slate-800 bg-white"
                            />
                          </div>
                          <div className="space-y-1 md:col-span-2">
                            <label className="text-[9px] font-bold text-slate-500 uppercase block">Nama Grup (Accordion Header) <span className="text-red-500">*</span></label>
                            <input
                              type="text"
                              required
                              value={group.label}
                              onChange={(e) => {
                                const newList = regulasiGroups.map((g, i) => i === gIdx ? { ...g, label: e.target.value } : g);
                                updateRegulasiContent(newList);
                              }}
                              className="w-full rounded-xl border border-slate-200 px-3 py-2 text-xs font-bold text-slate-800 bg-white"
                              placeholder="A. Regulasi Nasional"
                            />
                          </div>
                        </div>

                        <div className="space-y-1">
                          <label className="text-[9px] font-bold text-slate-500 uppercase block">Deskripsi Grup</label>
                          <textarea
                            value={group.desc}
                            onChange={(e) => {
                              const newList = regulasiGroups.map((g, i) => i === gIdx ? { ...g, desc: e.target.value } : g);
                              updateRegulasiContent(newList);
                            }}
                            className="w-full rounded-xl border border-slate-200 px-3 py-2 text-xs font-medium text-slate-800 bg-white h-16 resize-none"
                            placeholder="Deskripsi singkat mengenai regulasi grup ini..."
                          />
                        </div>

                        {/* Items list inside group */}
                        <div className="space-y-3 pt-3 border-t border-slate-200/60 text-left">
                          <div className="flex items-center justify-between">
                            <span className="text-[9px] font-extrabold text-[#002147] uppercase tracking-wider">
                              Daftar Regulasi / Dokumen di Grup Ini
                            </span>
                            <button
                              type="button"
                              onClick={() => {
                                const newItem = {
                                  no: String(group.items.length + 1).padStart(2, '0'),
                                  title: '',
                                  detail: '',
                                  fileUrl: ''
                                };
                                const newList = regulasiGroups.map((g, i) => i === gIdx ? { ...g, items: [...g.items, newItem] } : g);
                                updateRegulasiContent(newList);
                              }}
                              className="px-2.5 py-1 bg-[#002147] hover:bg-[#003166] text-white rounded-lg text-[9px] font-bold uppercase tracking-wider inline-flex items-center gap-1"
                            >
                              <Plus className="h-3 w-3" /> Tambah Baris
                            </button>
                          </div>

                          <div className="space-y-3">
                            {group.items && group.items.length > 0 ? (
                              group.items.map((item: any, iIdx: number) => (
                                <div key={iIdx} className="p-4 border border-slate-200 rounded-2xl bg-white space-y-3 relative text-left">
                                  <button
                                    type="button"
                                    onClick={() => {
                                      const updatedItems = group.items.filter((_: any, i: number) => i !== iIdx).map((item: any, idx: number) => ({
                                        ...item,
                                        no: String(idx + 1).padStart(2, '0')
                                      }));
                                      const newList = regulasiGroups.map((g, i) => i === gIdx ? { ...g, items: updatedItems } : g);
                                      updateRegulasiContent(newList);
                                    }}
                                    className="absolute top-3.5 right-3.5 p-1 text-slate-400 hover:text-red-650 rounded-md"
                                    title="Hapus Regulasi"
                                  >
                                    <Trash2 className="h-3.5 w-3.5" />
                                  </button>

                                  <div className="grid grid-cols-1 md:grid-cols-4 gap-3 pr-8">
                                    <div className="space-y-1">
                                      <label className="text-[8px] font-bold text-slate-400 uppercase block">No</label>
                                      <input
                                        type="text"
                                        required
                                        value={item.no || item.number || ''}
                                        onChange={(e) => {
                                          const updatedItems = group.items.map((it: any, i: number) => i === iIdx ? { ...it, no: e.target.value } : it);
                                          const newList = regulasiGroups.map((g, i) => i === gIdx ? { ...g, items: updatedItems } : g);
                                          updateRegulasiContent(newList);
                                        }}
                                        className="w-full rounded-lg border border-slate-200 px-2 py-1.5 text-xs font-bold text-slate-800 bg-slate-50"
                                      />
                                    </div>
                                    <div className="space-y-1 md:col-span-3">
                                      <label className="text-[8px] font-bold text-slate-400 uppercase block">Nama Regulasi/Aturan <span className="text-red-500">*</span></label>
                                      <input
                                        type="text"
                                        required
                                        value={item.title}
                                        onChange={(e) => {
                                          const updatedItems = group.items.map((it: any, i: number) => i === iIdx ? { ...it, title: e.target.value } : it);
                                          const newList = regulasiGroups.map((g, i) => i === gIdx ? { ...g, items: updatedItems } : g);
                                          updateRegulasiContent(newList);
                                        }}
                                        className="w-full rounded-lg border border-slate-200 px-2 py-1.5 text-xs font-bold text-slate-800"
                                        placeholder="Undang-Undang Nomor 14 Tahun 2008"
                                      />
                                    </div>
                                  </div>

                                  <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                                    <div className="space-y-1">
                                      <label className="text-[8px] font-bold text-slate-400 uppercase block">Detail/Tentang</label>
                                      <input
                                        type="text"
                                        value={item.detail || ''}
                                        onChange={(e) => {
                                          const updatedItems = group.items.map((it: any, i: number) => i === iIdx ? { ...it, detail: e.target.value } : it);
                                          const newList = regulasiGroups.map((g, i) => i === gIdx ? { ...g, items: updatedItems } : g);
                                          updateRegulasiContent(newList);
                                        }}
                                        className="w-full rounded-lg border border-slate-200 px-2.5 py-1.5 text-xs font-medium text-slate-800"
                                        placeholder="tentang Keterbukaan Informasi Publik"
                                      />
                                    </div>

                                    <div className="space-y-1">
                                      <label className="text-[8px] font-bold text-slate-400 uppercase block">File URL / Path PDF</label>
                                      <div className="flex gap-2">
                                        <input
                                          type="text"
                                          value={item.fileUrl || item.url || ''}
                                          onChange={(e) => {
                                            const updatedItems = group.items.map((it: any, i: number) => i === iIdx ? { ...it, fileUrl: e.target.value } : it);
                                            const newList = regulasiGroups.map((g, i) => i === gIdx ? { ...g, items: updatedItems } : g);
                                            updateRegulasiContent(newList);
                                          }}
                                          className="flex-1 rounded-lg border border-slate-200 px-2.5 py-1.5 text-xs font-mono font-medium text-slate-700 bg-slate-50"
                                          placeholder="/uploads/regulasi/..."
                                        />
                                        <label className="px-3 py-1 bg-slate-100 hover:bg-slate-200 rounded-lg text-xs font-bold cursor-pointer inline-flex items-center border border-slate-200">
                                          {isUploadingRegulasiFile[`${gIdx}-${iIdx}`] ? '...' : 'Upload'}
                                          <input
                                            type="file"
                                            onChange={async (e) => {
                                              const file = e.target.files?.[0];
                                              if (!file) return;

                                              setIsUploadingRegulasiFile(prev => ({ ...prev, [`${gIdx}-${iIdx}`]: true }));
                                              const token = localStorage.getItem('auth_token');
                                              const formData = new FormData();
                                              formData.append('folder', 'regulasi');
                                              formData.append('files', file);

                                              try {
                                                const res = await fetch(`${API_BASE_URL}/admin/uploads`, {
                                                  method: 'POST',
                                                  headers: { Authorization: `Bearer ${token}` },
                                                  body: formData
                                                });
                                                if (res.ok) {
                                                  const data = await res.json();
                                                  const url = data.url || (data.uploaded && data.uploaded[0] && data.uploaded[0].url) || '';
                                                  if (url) {
                                                    const updatedItems = group.items.map((it: any, i: number) => i === iIdx ? { ...it, fileUrl: url } : it);
                                                    const newList = regulasiGroups.map((g, i) => i === gIdx ? { ...g, items: updatedItems } : g);
                                                    updateRegulasiContent(newList);
                                                  }
                                                }
                                              } catch (err) {
                                                console.error(err);
                                              } finally {
                                                setIsUploadingRegulasiFile(prev => ({ ...prev, [`${gIdx}-${iIdx}`]: false }));
                                                e.target.value = '';
                                              }
                                            }}
                                            className="hidden"
                                          />
                                        </label>
                                      </div>
                                    </div>
                                  </div>

                                  {/* SOP List Section inside item */}
                                  <div className="space-y-2 pt-2 border-t border-slate-100">
                                    <div className="flex items-center justify-between">
                                      <label className="text-[8px] font-bold text-slate-400 uppercase inline-flex items-center gap-1.5">
                                        <input
                                          type="checkbox"
                                          checked={!!item.isSopList}
                                          onChange={(e) => {
                                            const updatedItems = group.items.map((it: any, i: number) => i === iIdx ? {
                                              ...it,
                                              isSopList: e.target.checked,
                                              sops: e.target.checked ? (it.sops || ['']) : []
                                            } : it);
                                            const newList = regulasiGroups.map((g, i) => i === gIdx ? { ...g, items: updatedItems } : g);
                                            updateRegulasiContent(newList);
                                          }}
                                          className="rounded border-slate-300 text-blue-600 focus:ring-blue-500"
                                        />
                                        Memiliki Daftar SOP?
                                      </label>
                                      {item.isSopList && (
                                        <button
                                          type="button"
                                          onClick={() => {
                                            const updatedItems = group.items.map((it: any, i: number) => i === iIdx ? {
                                              ...it,
                                              sops: [...(it.sops || []), '']
                                            } : it);
                                            const newList = regulasiGroups.map((g, i) => i === gIdx ? { ...g, items: updatedItems } : g);
                                            updateRegulasiContent(newList);
                                          }}
                                          className="text-[8px] font-bold text-[#002147] hover:underline cursor-pointer"
                                        >
                                          + SOP Baru
                                        </button>
                                      )}
                                    </div>

                                    {item.isSopList && item.sops && item.sops.length > 0 && (
                                      <div className="space-y-2 pl-4">
                                        {item.sops.map((sop: string, sIdx: number) => (
                                          <div key={sIdx} className="flex gap-2 items-center">
                                            <input
                                              type="text"
                                              required
                                              value={sop}
                                              onChange={(e) => {
                                                const updatedSops = item.sops.map((s: string, idx: number) => idx === sIdx ? e.target.value : s);
                                                const updatedItems = group.items.map((it: any, i: number) => i === iIdx ? { ...it, sops: updatedSops } : it);
                                                const newList = regulasiGroups.map((g, i) => i === gIdx ? { ...g, items: updatedItems } : g);
                                                updateRegulasiContent(newList);
                                              }}
                                              className="flex-1 rounded-lg border border-slate-200 px-2 py-1 text-[10px]"
                                              placeholder="Contoh: SOP Pelayanan Informasi Publik"
                                            />
                                            <button
                                              type="button"
                                              onClick={() => {
                                                const updatedSops = item.sops.filter((_: any, idx: number) => idx !== sIdx);
                                                const updatedItems = group.items.map((it: any, i: number) => i === iIdx ? { ...it, sops: updatedSops, isSopList: updatedSops.length > 0 } : it);
                                                const newList = regulasiGroups.map((g, i) => i === gIdx ? { ...g, items: updatedItems } : g);
                                                updateRegulasiContent(newList);
                                              }}
                                              className="text-slate-400 hover:text-red-650 cursor-pointer"
                                            >
                                              <Trash2 className="h-3 w-3" />
                                            </button>
                                          </div>
                                        ))}
                                      </div>
                                    )}
                                  </div>

                                </div>
                              ))
                            ) : (
                              <div className="text-center text-slate-400 py-4 border border-dashed border-slate-200 rounded-xl text-[10px]">
                                Belum ada berkas regulasi. Klik "+ Tambah Baris" di atas.
                              </div>
                            )}
                          </div>
                        </div>

                      </div>
                    ))
                  ) : (
                    <div className="py-12 text-center text-slate-400 text-xs border border-dashed border-slate-200 rounded-3xl">
                      Belum ada grup accordion regulasi. Klik "Tambah Grup" di atas.
                    </div>
                  )}
                </div>
              </div>
            ) : ['informasi-publik-berkala', 'informasi-tersedia-setiap-saat', 'info-serta-merta', 'informasi-dikecualikan'].includes(adminEditSlug) ? (
              <div className="space-y-4 pt-4 border-t border-slate-100 text-left">
                <div className="bg-blue-50 border border-blue-200 rounded-2xl p-6 flex items-start gap-4 shadow-sm mb-4">
                  <Info className="h-5 w-5 text-blue-600 shrink-0 mt-0.5" />
                  <div className="space-y-1">
                    <h4 className="text-xs font-bold text-blue-900 uppercase tracking-wide">Dokumen Disinkronkan Otomatis</h4>
                    <p className="text-[11px] text-blue-750 font-semibold leading-relaxed">
                      Dokumen lampiran untuk halaman ini disinkronkan otomatis secara real-time dari database menu <strong>Download</strong>.
                    </p>
                    <p className="text-[11px] text-slate-500 font-medium leading-relaxed mt-1">
                      Untuk menambah, mengubah, atau menghapus berkas pada halaman ini, silakan buka tab <strong>Download</strong> di panel kiri admin, lalu pilih kategori yang sesuai dengan halaman ini.
                    </p>
                  </div>
                </div>

                {/* Dynamic Sections Editor */}
                <div className="space-y-4 pt-3 border-t border-slate-100">
                  <div className="flex items-center justify-between pb-2">
                    <div>
                      <span className="text-[10px] font-extrabold text-[#002147] uppercase tracking-wider block">
                        Daftar Seksi Pengantar / Deskripsi Halaman
                      </span>
                      <span className="text-[11px] text-slate-400 font-medium block">
                        Kelola blok paragraf penjelasan, gambar pendukung, dan tata letak per seksi.
                      </span>
                    </div>
                    <button
                      type="button"
                      onClick={() => {
                        const newSec: DipSection = {
                          text: '',
                          imageUrls: [],
                          imagePosition: 'right'
                        };
                        updateDipSections([...dipSections, newSec]);
                      }}
                      className="px-3.5 py-2 bg-[#002147] hover:bg-[#003166] text-white rounded-xl text-[10px] font-bold uppercase tracking-wider inline-flex items-center gap-1 cursor-pointer border border-[#002147]/55"
                    >
                      <Plus className="h-3.5 w-3.5" /> Tambah Seksi
                    </button>
                  </div>

                  <div className="space-y-4 max-h-[400px] overflow-y-auto pr-1">
                    {dipSections.length > 0 ? (
                      dipSections.map((section, sIdx) => (
                        <div key={sIdx} className="p-4 border border-slate-200 rounded-2xl bg-slate-50/50 space-y-3 relative text-left">
                          
                          {/* Reorder and Delete controls */}
                          <div className="absolute top-3.5 right-3.5 flex items-center gap-1">
                            {sIdx > 0 && (
                              <button
                                type="button"
                                onClick={() => {
                                  const arr = [...dipSections];
                                  [arr[sIdx - 1], arr[sIdx]] = [arr[sIdx], arr[sIdx - 1]];
                                  updateDipSections(arr);
                                }}
                                className="p-1 text-slate-400 hover:text-blue-600 hover:bg-blue-50 rounded"
                                title="Geser ke atas"
                              >
                                <ArrowUp className="h-3.5 w-3.5" />
                              </button>
                            )}
                            {sIdx < dipSections.length - 1 && (
                              <button
                                type="button"
                                onClick={() => {
                                  const arr = [...dipSections];
                                  [arr[sIdx], arr[sIdx + 1]] = [arr[sIdx + 1], arr[sIdx]];
                                  updateDipSections(arr);
                                }}
                                className="p-1 text-slate-400 hover:text-blue-600 hover:bg-blue-50 rounded"
                                title="Geser ke bawah"
                              >
                                <ArrowDown className="h-3.5 w-3.5" />
                              </button>
                            )}
                            <button
                              type="button"
                              onClick={() => {
                                const arr = dipSections.filter((_, i) => i !== sIdx);
                                updateDipSections(arr);
                              }}
                              className="p-1 text-slate-400 hover:text-red-600 hover:bg-red-50 rounded"
                              title="Hapus Seksi"
                            >
                              <Trash2 className="h-3.5 w-3.5" />
                            </button>
                          </div>

                          <span className="text-[9px] font-extrabold text-[#002147] bg-[#002147]/5 px-2.5 py-1 rounded-lg">
                            Seksi Deskripsi #{sIdx + 1}
                          </span>

                          <div className="grid grid-cols-1 md:grid-cols-2 gap-6 pt-2">
                            {/* Text Area */}
                            <div className="space-y-1">
                              <label className="text-[9px] font-bold text-slate-500 uppercase block">Teks Deskripsi / Penjelasan <span className="text-red-500">*</span></label>
                              <textarea
                                required
                                value={section.text}
                                onChange={(e) => {
                                  const arr = dipSections.map((sec, i) => i === sIdx ? { ...sec, text: e.target.value } : sec);
                                  updateDipSections(arr);
                                }}
                                className="w-full rounded-xl border border-slate-200 px-3 py-2 text-xs font-medium text-slate-800 bg-white h-[360px] resize-none focus:outline-none"
                                placeholder="Masukkan penjelasan teks pengantar di sini..."
                              />
                            </div>

                            {/* Image / Layout Settings */}
                            <div className="space-y-4">
                              <div className="space-y-2">
                                <label className="text-[9px] font-bold text-slate-500 uppercase block">Daftar Gambar Pendukung (Maks. 3 Gambar untuk Mode Bertumpuk)</label>
                                <div className="space-y-2.5">
                                  {[0, 1, 2].map((imgIdx) => {
                                    const currentUrl = section.imageUrls[imgIdx] || '';
                                    const uploadKey = `${sIdx}-${imgIdx}`;

                                    return (
                                      <div key={imgIdx} className="space-y-1 bg-white border border-slate-200/80 rounded-xl p-2.5 shadow-sm">
                                        <div className="flex items-center justify-between">
                                          <span className="text-[8px] font-extrabold text-slate-400 uppercase">Slot Gambar 0{imgIdx + 1}</span>
                                          {currentUrl && (
                                            <button
                                              type="button"
                                              onClick={() => {
                                                const newUrls = [...section.imageUrls];
                                                newUrls.splice(imgIdx, 1); // remove
                                                const cleaned = newUrls.filter(Boolean);
                                                const arr = dipSections.map((sec, i) => i === sIdx ? { ...sec, imageUrls: cleaned } : sec);
                                                updateDipSections(arr);
                                              }}
                                              className="text-[9px] font-bold text-rose-500 hover:text-rose-700 cursor-pointer border-0 bg-transparent"
                                            >
                                              Hapus
                                            </button>
                                          )}
                                        </div>

                                        <div className="flex gap-2">
                                          <input
                                            type="text"
                                            value={currentUrl}
                                            onChange={(e) => {
                                              const newUrls = [...section.imageUrls];
                                              newUrls[imgIdx] = e.target.value;
                                              const arr = dipSections.map((sec, i) => i === sIdx ? { ...sec, imageUrls: newUrls } : sec);
                                              updateDipSections(arr);
                                            }}
                                            className="flex-1 rounded-lg border border-slate-200 px-2.5 py-1 text-[10px] font-mono font-medium text-slate-700 bg-slate-50 focus:outline-none"
                                            placeholder="https://... atau /uploads/..."
                                          />
                                          <label className="px-2.5 py-1 bg-slate-100 hover:bg-slate-200 hover:text-slate-800 rounded-lg text-[10px] font-bold transition-all cursor-pointer shrink-0 inline-flex items-center border border-slate-200 text-slate-600">
                                            {isUploadingDipSectionImage[uploadKey] ? '...' : 'Upload'}
                                            <input
                                              type="file"
                                              onChange={async (e) => {
                                                const file = e.target.files?.[0];
                                                if (!file) return;

                                                setIsUploadingDipSectionImage(prev => ({ ...prev, [uploadKey]: true }));
                                                const token = localStorage.getItem('auth_token');
                                                const formData = new FormData();
                                                formData.append('folder', 'pages');
                                                formData.append('files', file);

                                                try {
                                                  const res = await fetch(`${API_BASE_URL}/admin/uploads`, {
                                                    method: 'POST',
                                                    headers: { Authorization: `Bearer ${token}` },
                                                    body: formData
                                                  });
                                                  if (res.ok) {
                                                    const data = await res.json();
                                                    const url = data.url || (data.uploaded && data.uploaded[0] && data.uploaded[0].url) || '';
                                                    if (url) {
                                                      const newUrls = [...section.imageUrls];
                                                      newUrls[imgIdx] = url;
                                                      const arr = dipSections.map((sec, i) => i === sIdx ? { ...sec, imageUrls: newUrls } : sec);
                                                      updateDipSections(arr);
                                                    }
                                                  }
                                                } catch (err) {
                                                  console.error(err);
                                                } finally {
                                                  setIsUploadingDipSectionImage(prev => ({ ...prev, [uploadKey]: false }));
                                                  e.target.value = '';
                                                }
                                              }}
                                              className="hidden"
                                            />
                                          </label>
                                        </div>

                                        {currentUrl && (
                                          <div className="mt-1.5 flex justify-center">
                                            <img
                                              src={resolveImageUrl(currentUrl)}
                                              alt={`Preview ${imgIdx + 1}`}
                                              className="h-12 w-auto object-cover rounded-lg border border-slate-150"
                                            />
                                          </div>
                                        )}
                                      </div>
                                    );
                                  })}
                                </div>
                              </div>

                              <div className="space-y-1">
                                <label className="text-[9px] font-bold text-slate-500 uppercase block">Tata Letak / Posisi Gambar</label>
                                <select
                                  value={section.imagePosition}
                                  onChange={(e) => {
                                    const arr = dipSections.map((sec, i) => i === sIdx ? { ...sec, imagePosition: e.target.value as 'left' | 'right' } : sec);
                                    updateDipSections(arr);
                                  }}
                                  className="w-full rounded-xl border border-slate-200 px-3 py-2 text-xs font-bold text-slate-800 bg-white focus:outline-none"
                                >
                                  <option value="right">Gambar di Kanan (Teks di Kiri)</option>
                                  <option value="left">Gambar di Kiri (Teks di Kanan)</option>
                                </select>
                              </div>
                            </div>
                          </div>

                        </div>
                      ))
                    ) : (
                      <div className="py-8 text-center text-slate-400 text-xs font-medium border border-dashed border-slate-200 rounded-2xl">
                        Belum ada seksi deskripsi. Silakan klik "Tambah Seksi" di atas.
                      </div>
                    )}
                  </div>
                </div>
              </div>
            ) : ['zona-integrasi', 'keberatan-informasi', 'Permohonan-penyelesaian-sengketa', 'permohonan-penyelesaian-sengketa'].includes(adminEditSlug) ? (
              <div className="space-y-4 pt-2 border-t border-slate-100 text-left">
                <div className="flex items-center justify-between pb-2 border-b border-slate-100">
                  <div>
                    <span className="text-[10px] font-extrabold text-[#002147] uppercase tracking-wider block">
                      Daftar Dokumen Lampiran
                    </span>
                    <span className="text-[11px] text-slate-400 font-medium block">
                      Tambahkan file upload atau input link URL yang akan ditampilkan pada halaman rujukan.
                    </span>
                  </div>
                  <div className="flex items-center gap-2">
                    <button
                      type="button"
                      onClick={() => setAdminEditPageDocs(prev => [...prev, { title: '', description: '', file_url: '' }])}
                      className="px-3.5 py-2 bg-[#002147] hover:bg-[#003166] text-white rounded-xl text-[10px] font-bold uppercase tracking-wider inline-flex items-center gap-1 cursor-pointer border border-[#002147]/55"
                      disabled={isUploadingDocs}
                    >
                      <Plus className="h-3.5 w-3.5" /> Tambah Berkas/Link
                    </button>
                    <label className="px-3.5 py-2 bg-emerald-600 hover:bg-emerald-700 text-white rounded-xl text-[10px] font-bold uppercase tracking-wider inline-flex items-center gap-1.5 cursor-pointer border border-emerald-600 shadow-sm transition-all">
                      {isUploadingDocs ? (
                        <>
                          <Loader2 className="h-3.5 w-3.5 animate-spin" /> Uploading...
                        </>
                      ) : (
                        <>
                          <Upload className="h-3.5 w-3.5" /> Bulk Upload PDF
                        </>
                      )}
                      <input
                        type="file"
                        multiple
                        accept="application/pdf"
                        onChange={handleBulkUploadPageDocs}
                        className="hidden"
                        disabled={isUploadingDocs}
                      />
                    </label>
                  </div>
                </div>

                {/* Dynamic List Input Rows */}
                <div className="space-y-4 max-h-[400px] overflow-y-auto pr-1">
                  {adminEditPageDocs.length > 0 ? (
                    adminEditPageDocs.map((doc, idx) => (
                      <div key={idx} className="p-4 border border-slate-200 rounded-2xl bg-slate-50/50 space-y-3 relative text-left">
                        <div className="absolute top-3.5 right-3.5">
                          <button
                            type="button"
                            onClick={() => setAdminEditPageDocs(prev => prev.filter((_, i) => i !== idx))}
                            className="p-1.5 text-slate-400 hover:text-red-650 hover:bg-red-50 rounded-lg transition-all border border-transparent"
                            title="Hapus Baris"
                          >
                            <Trash2 className="h-4 w-4" />
                          </button>
                        </div>

                        <div className="grid grid-cols-1 md:grid-cols-2 gap-3 pr-8 text-left">
                          <div className="space-y-1">
                            <label className="text-[9px] font-bold text-slate-500 uppercase block">Nama/Judul Dokumen <span className="text-red-500">*</span></label>
                            <input
                              type="text"
                              required
                              value={doc.title}
                              onChange={(e) => setAdminEditPageDocs(prev => prev.map((d, i) => i === idx ? { ...d, title: e.target.value } : d))}
                              className="w-full rounded-xl border border-slate-200 px-3 py-2.5 text-xs focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 bg-white font-bold text-slate-800"
                              placeholder="Contoh: Salinan SK Kepengurusan"
                            />
                          </div>
                          <div className="space-y-1">
                            <label className="text-[9px] font-bold text-slate-500 uppercase block">Deskripsi Singkat (Opsional)</label>
                            <input
                              type="text"
                              value={doc.description}
                              onChange={(e) => setAdminEditPageDocs(prev => prev.map((d, i) => i === idx ? { ...d, description: e.target.value } : d))}
                              className="w-full rounded-xl border border-slate-200 px-3 py-2.5 text-xs focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 bg-white font-medium text-slate-700"
                              placeholder="Keterangan singkat tentang isi berkas"
                            />
                          </div>
                        </div>

                        <div className="space-y-1 text-left">
                          <label className="text-[9px] font-bold text-slate-500 uppercase block">Unggah File atau Input Link URL <span className="text-red-500">*</span></label>
                          <div className="flex gap-2">
                            <input
                              type="text"
                              required
                              value={doc.file_url}
                              onChange={(e) => setAdminEditPageDocs(prev => prev.map((d, i) => i === idx ? { ...d, file_url: e.target.value } : d))}
                              className="flex-1 rounded-xl border border-slate-200 px-3 py-2.5 text-xs focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 bg-white font-mono font-medium text-slate-700"
                              placeholder="https://example.com/file.pdf atau upload berkas"
                            />
                            <label className="px-3.5 py-2 bg-slate-100 hover:bg-slate-250 hover:text-slate-800 rounded-xl text-xs font-bold transition-all cursor-pointer shrink-0 inline-flex items-center border border-slate-200 text-slate-600">
                              Upload
                              <input
                                type="file"
                                multiple
                                onChange={async (e) => {
                                  const files = e.target.files;
                                  if (!files || files.length === 0) return;

                                  const formData = new FormData();
                                  formData.append('folder', 'documents');
                                  for (let i = 0; i < files.length; i++) {
                                    formData.append('files', files[i]);
                                  }

                                  try {
                                    const res = await fetch(`${API_BASE_URL}/admin/uploads`, {
                                      method: 'POST',
                                      headers: { Authorization: `Bearer ${localStorage.getItem('auth_token')}` },
                                      body: formData
                                    });
                                    if (res.ok) {
                                      const data = await res.json();
                                      const uploadedList = data.uploaded || [];
                                      if (uploadedList.length > 0) {
                                        // Replace current row with first file
                                        const firstMatch = uploadedList[0];
                                        setAdminEditPageDocs(prev => prev.map((d, i) => i === idx ? { 
                                          ...d, 
                                          file_url: firstMatch.url,
                                          title: d.title || (firstMatch.original_name ? firstMatch.original_name.replace(/\.[^/.]+$/, "").replace(/[_-]+/g, ' ').replace(/(^\w|\s\w)/g, (m: string) => m.toUpperCase()) : "Dokumen")
                                        } : d));

                                        // Append subsequent files as new rows
                                        if (uploadedList.length > 1) {
                                          const extraDocs = uploadedList.slice(1).map((item: any) => {
                                            const titleWithoutExt = item.original_name ? item.original_name.replace(/\.[^/.]+$/, "") : "Dokumen";
                                            const formattedTitle = titleWithoutExt
                                              .replace(/[_-]+/g, ' ')
                                              .replace(/(^\w|\s\w)/g, (m: string) => m.toUpperCase());
                                            return {
                                              title: formattedTitle,
                                              description: '',
                                              file_url: item.url
                                            };
                                          });
                                          setAdminEditPageDocs(prev => [...prev, ...extraDocs]);
                                        }
                                      }
                                    }
                                  } catch (err) {
                                    console.error(err);
                                  }
                                }}
                                className="hidden"
                              />
                            </label>
                          </div>
                        </div>
                      </div>
                    ))
                  ) : (
                    <div className="py-8 text-center text-slate-400 text-xs font-medium border border-dashed border-slate-200 rounded-2xl">
                      Belum ada lampiran dokumen/link. Klik "Tambah Berkas/Link" di atas.
                    </div>
                  )}
                </div>

                {/* Sengketa Verification Cards Editor */}
                {(adminEditSlug === 'Permohonan-penyelesaian-sengketa' || adminEditSlug === 'permohonan-penyelesaian-sengketa') && (
                  <div className="space-y-4 pt-4 border-t border-slate-150 text-left">
                    <div className="flex items-center justify-between">
                      <div>
                        <span className="text-[10px] font-extrabold text-[#002147] uppercase tracking-wider block">
                          Kartu Persyaratan Dokumen Verifikasi
                        </span>
                        <span className="text-[11px] text-slate-400 font-medium block">
                          Kelola daftar kartu persyaratan dokumen yang ditampilkan di halaman sengketa.
                        </span>
                      </div>
                      <button
                        type="button"
                        onClick={() => setEditSengketaCards([...editSengketaCards, { title: '', desc: '' }])}
                        className="px-3.5 py-2 bg-[#002147] hover:bg-[#003166] text-white rounded-xl text-[10px] font-bold uppercase tracking-wider inline-flex items-center gap-1 cursor-pointer border border-[#002147]/55"
                      >
                        <Plus className="h-3.5 w-3.5" /> Tambah Kartu
                      </button>
                    </div>

                    <div className="space-y-3 max-h-[400px] overflow-y-auto pr-1">
                      {editSengketaCards.length > 0 ? (
                        editSengketaCards.map((card, idx) => (
                          <div key={idx} className="p-4 border border-slate-200 rounded-2xl bg-slate-50/50 space-y-3 relative text-left">
                            <div className="absolute top-3.5 right-3.5 flex items-center gap-1">
                              {idx > 0 && (
                                <button
                                  type="button"
                                  onClick={() => {
                                    const arr = [...editSengketaCards];
                                    [arr[idx - 1], arr[idx]] = [arr[idx], arr[idx - 1]];
                                    setEditSengketaCards(arr);
                                  }}
                                  className="p-1.5 text-slate-400 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-all border border-transparent"
                                  title="Geser ke atas"
                                >
                                  <ArrowUp className="h-3.5 w-3.5" />
                                </button>
                              )}
                              {idx < editSengketaCards.length - 1 && (
                                <button
                                  type="button"
                                  onClick={() => {
                                    const arr = [...editSengketaCards];
                                    [arr[idx], arr[idx + 1]] = [arr[idx + 1], arr[idx]];
                                    setEditSengketaCards(arr);
                                  }}
                                  className="p-1.5 text-slate-400 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-all border border-transparent"
                                  title="Geser ke bawah"
                                >
                                  <ArrowDown className="h-3.5 w-3.5" />
                                </button>
                              )}
                              <button
                                type="button"
                                onClick={() => setEditSengketaCards(editSengketaCards.filter((_, i) => i !== idx))}
                                className="p-1.5 text-slate-400 hover:text-red-650 hover:bg-red-50 rounded-lg transition-all border border-transparent"
                                title="Hapus Kartu"
                              >
                                <Trash2 className="h-4 w-4" />
                              </button>
                            </div>

                            <div className="flex items-center gap-2 mb-1">
                              <span className="h-6 w-6 rounded-lg bg-amber-100 text-amber-700 text-[10px] font-extrabold flex items-center justify-center shrink-0">{idx + 1}</span>
                              <span className="text-[9px] font-bold text-slate-400 uppercase">Kartu Persyaratan #{idx + 1}</span>
                            </div>

                            <div className="space-y-2 pr-16">
                              <div className="space-y-1">
                                <label className="text-[9px] font-bold text-slate-500 uppercase block">Judul Kartu <span className="text-red-500">*</span></label>
                                <input
                                  type="text"
                                  value={card.title}
                                  onChange={(e) => setEditSengketaCards(editSengketaCards.map((c, i) => i === idx ? { ...c, title: e.target.value } : c))}
                                  className="w-full rounded-xl border border-slate-200 px-3 py-2.5 text-xs focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 bg-white font-bold text-slate-800"
                                  placeholder="Contoh: Bukti Surat Permohonan Informasi"
                                />
                              </div>
                              <div className="space-y-1">
                                <label className="text-[9px] font-bold text-slate-500 uppercase block">Deskripsi Singkat</label>
                                <textarea
                                  value={card.desc}
                                  onChange={(e) => setEditSengketaCards(editSengketaCards.map((c, i) => i === idx ? { ...c, desc: e.target.value } : c))}
                                  className="w-full rounded-xl border border-slate-200 px-3 py-2.5 text-xs focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 bg-white font-medium text-slate-700 h-16 resize-none"
                                  placeholder="Penjelasan singkat tentang dokumen ini"
                                />
                              </div>
                            </div>
                          </div>
                        ))
                      ) : (
                        <div className="py-8 text-center text-slate-400 text-xs font-medium border border-dashed border-slate-200 rounded-2xl">
                          Belum ada kartu persyaratan. Klik "Tambah Kartu" di atas, atau biarkan kosong untuk menggunakan kartu bawaan.
                        </div>
                      )}
                    </div>
                  </div>
                )}

                {/* Sengketa Flow Text Editor */}
                {(adminEditSlug === 'Permohonan-penyelesaian-sengketa' || adminEditSlug === 'permohonan-penyelesaian-sengketa') && (
                  <div className="space-y-4 pt-4 border-t border-slate-150 text-left">
                    <div>
                      <span className="text-[10px] font-extrabold text-[#002147] uppercase tracking-wider block">
                        Editor Teks Alur Pengajuan Sengketa
                      </span>
                      <span className="text-[11px] text-slate-400 font-medium block">
                        Ubah judul dan deskripsi setiap elemen pada diagram alur. Kosongkan untuk menggunakan teks bawaan.
                      </span>
                    </div>

                    {/* Jalur A */}
                    <div className="p-4 border border-amber-200 rounded-2xl bg-amber-50/30 space-y-3">
                      <div className="flex items-center gap-2">
                        <span className="h-5 w-5 rounded-md bg-amber-400 text-[#002147] text-[9px] font-extrabold flex items-center justify-center shrink-0">A</span>
                        <span className="text-[9px] font-bold text-amber-800 uppercase">Jalur A — Datang Langsung</span>
                      </div>
                      <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                        <div className="space-y-1">
                          <label className="text-[9px] font-bold text-slate-500 uppercase block">Judul</label>
                          <input
                            type="text"
                            value={editSengketaFlow.jalurA_title}
                            onChange={(e) => setEditSengketaFlow({ ...editSengketaFlow, jalurA_title: e.target.value })}
                            className="w-full rounded-xl border border-slate-200 px-3 py-2 text-xs focus:outline-none focus:ring-2 focus:ring-amber-500/20 focus:border-amber-500 bg-white font-bold text-slate-800"
                            placeholder="Datang Langsung (Manual)"
                          />
                        </div>
                        <div className="space-y-1">
                          <label className="text-[9px] font-bold text-slate-500 uppercase block">Deskripsi</label>
                          <textarea
                            value={editSengketaFlow.jalurA_desc}
                            onChange={(e) => setEditSengketaFlow({ ...editSengketaFlow, jalurA_desc: e.target.value })}
                            className="w-full rounded-xl border border-slate-200 px-3 py-2 text-xs focus:outline-none focus:ring-2 focus:ring-amber-500/20 focus:border-amber-500 bg-white font-medium text-slate-700 h-14 resize-none"
                            placeholder="Pemohon menyerahkan surat/berkas pengajuan sengketa secara langsung..."
                          />
                        </div>
                      </div>
                    </div>

                    {/* Jalur B */}
                    <div className="p-4 border border-blue-200 rounded-2xl bg-blue-50/30 space-y-3">
                      <div className="flex items-center gap-2">
                        <span className="h-5 w-5 rounded-md bg-blue-500 text-white text-[9px] font-extrabold flex items-center justify-center shrink-0">B</span>
                        <span className="text-[9px] font-bold text-blue-800 uppercase">Jalur B — Kirim Surat/Pos</span>
                      </div>
                      <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                        <div className="space-y-1">
                          <label className="text-[9px] font-bold text-slate-500 uppercase block">Judul</label>
                          <input
                            type="text"
                            value={editSengketaFlow.jalurB_title}
                            onChange={(e) => setEditSengketaFlow({ ...editSengketaFlow, jalurB_title: e.target.value })}
                            className="w-full rounded-xl border border-slate-200 px-3 py-2 text-xs focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 bg-white font-bold text-slate-800"
                            placeholder="Kirim Surat/Pos"
                          />
                        </div>
                        <div className="space-y-1">
                          <label className="text-[9px] font-bold text-slate-500 uppercase block">Deskripsi</label>
                          <textarea
                            value={editSengketaFlow.jalurB_desc}
                            onChange={(e) => setEditSengketaFlow({ ...editSengketaFlow, jalurB_desc: e.target.value })}
                            className="w-full rounded-xl border border-slate-200 px-3 py-2 text-xs focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 bg-white font-medium text-slate-700 h-14 resize-none"
                            placeholder="Pemohon mengirimkan berkas pengajuan sengketa melalui jasa pengiriman surat..."
                          />
                        </div>
                      </div>
                    </div>

                    {/* Verifikasi */}
                    <div className="p-4 border border-slate-200 rounded-2xl bg-slate-50/50 space-y-3">
                      <div className="flex items-center gap-2">
                        <span className="h-5 w-5 rounded-md bg-amber-500 text-white text-[9px] font-extrabold flex items-center justify-center shrink-0">✓</span>
                        <span className="text-[9px] font-bold text-slate-700 uppercase">Tahap Verifikasi</span>
                      </div>
                      <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                        <div className="space-y-1">
                          <label className="text-[9px] font-bold text-slate-500 uppercase block">Judul</label>
                          <input
                            type="text"
                            value={editSengketaFlow.verifikasi_title}
                            onChange={(e) => setEditSengketaFlow({ ...editSengketaFlow, verifikasi_title: e.target.value })}
                            className="w-full rounded-xl border border-slate-200 px-3 py-2 text-xs focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 bg-white font-bold text-slate-800"
                            placeholder="Tahap 2: Verifikasi Dokumen Pendukung oleh Petugas"
                          />
                        </div>
                        <div className="space-y-1">
                          <label className="text-[9px] font-bold text-slate-500 uppercase block">Deskripsi</label>
                          <textarea
                            value={editSengketaFlow.verifikasi_desc}
                            onChange={(e) => setEditSengketaFlow({ ...editSengketaFlow, verifikasi_desc: e.target.value })}
                            className="w-full rounded-xl border border-slate-200 px-3 py-2 text-xs focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 bg-white font-medium text-slate-700 h-14 resize-none"
                            placeholder="Petugas PPID / Komisi Informasi akan memverifikasi kelengkapan berkas wajib..."
                          />
                        </div>
                      </div>
                    </div>

                    {/* Panduan Box */}
                    <div className="p-4 border border-slate-200 rounded-2xl bg-slate-50/50 space-y-3">
                      <div className="flex items-center gap-2">
                        <span className="h-5 w-5 rounded-md bg-slate-200 text-slate-600 text-[9px] font-extrabold flex items-center justify-center shrink-0">?</span>
                        <span className="text-[9px] font-bold text-slate-700 uppercase">Kotak Panduan Lanjutan</span>
                      </div>
                      <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                        <div className="space-y-1">
                          <label className="text-[9px] font-bold text-slate-500 uppercase block">Judul</label>
                          <input
                            type="text"
                            value={editSengketaFlow.help_title}
                            onChange={(e) => setEditSengketaFlow({ ...editSengketaFlow, help_title: e.target.value })}
                            className="w-full rounded-xl border border-slate-200 px-3 py-2 text-xs focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 bg-white font-bold text-slate-800"
                            placeholder="Panduan Lanjutan & Sekretariat"
                          />
                        </div>
                        <div className="space-y-1">
                          <label className="text-[9px] font-bold text-slate-500 uppercase block">Deskripsi</label>
                          <textarea
                            value={editSengketaFlow.help_desc}
                            onChange={(e) => setEditSengketaFlow({ ...editSengketaFlow, help_desc: e.target.value })}
                            className="w-full rounded-xl border border-slate-200 px-3 py-2 text-xs focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 bg-white font-medium text-slate-700 h-14 resize-none"
                            placeholder="Kantor Komisi Informasi Provinsi Sumatera Barat berlokasi di..."
                          />
                        </div>
                      </div>
                    </div>

                    {/* Download Box */}
                    <div className="p-4 border border-[#002147]/20 rounded-2xl bg-[#002147]/5 space-y-3">
                      <div className="flex items-center gap-2">
                        <span className="h-5 w-5 rounded-md bg-[#002147] text-amber-400 text-[9px] font-extrabold flex items-center justify-center shrink-0">↓</span>
                        <span className="text-[9px] font-bold text-[#002147] uppercase">Kotak Download</span>
                      </div>
                      <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                        <div className="space-y-1">
                          <label className="text-[9px] font-bold text-slate-500 uppercase block">Judul</label>
                          <input
                            type="text"
                            value={editSengketaFlow.download_title}
                            onChange={(e) => setEditSengketaFlow({ ...editSengketaFlow, download_title: e.target.value })}
                            className="w-full rounded-xl border border-slate-200 px-3 py-2 text-xs focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 bg-white font-bold text-slate-800"
                            placeholder="Formulir Cetak"
                          />
                        </div>
                        <div className="space-y-1">
                          <label className="text-[9px] font-bold text-slate-500 uppercase block">Deskripsi</label>
                          <textarea
                            value={editSengketaFlow.download_desc}
                            onChange={(e) => setEditSengketaFlow({ ...editSengketaFlow, download_desc: e.target.value })}
                            className="w-full rounded-xl border border-slate-200 px-3 py-2 text-xs focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 bg-white font-medium text-slate-700 h-14 resize-none"
                            placeholder="Unduh berkas kelengkapan pengujian konsekuensi / formulir resmi..."
                          />
                        </div>
                      </div>
                    </div>
                  </div>
                )}

                {/* Paragraph Intro text */}
                <div className="space-y-1 pt-2 border-t border-slate-100 text-left">
                  <label className="text-[10px] font-bold text-slate-600 uppercase tracking-wider block">Paragraf Pengantar Halaman (Opsional)</label>
                  <textarea
                    value={adminEditContent}
                    onChange={(e) => setAdminEditContent(e.target.value)}
                    placeholder="Tulis kalimat penjelasan umum untuk ditaruh sebelum daftar lampiran"
                    className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 transition-all bg-slate-50 h-24 font-medium resize-none text-slate-805"
                  />
                </div>

                {/* Timeline Flowchart Steps Editor */}
                {adminEditSlug === 'keberatan-informasi' && (
                  <div className="space-y-4 pt-4 border-t border-slate-150 text-left">
                    <div className="flex justify-between items-center">
                      <div>
                        <span className="text-[10px] font-extrabold text-[#002147] uppercase tracking-wider block">
                          Editor Langkah Bagan Alur (Timeline Steps)
                        </span>
                        <span className="text-[11px] text-slate-400 font-medium block font-sans">
                          Ubah, tambah, hapus, atau geser posisi langkah pada bagan alur pengajuan keberatan.
                        </span>
                      </div>
                    </div>

                    {/* Manual steps */}
                    <div className="space-y-2">
                      <div className="flex justify-between items-center">
                        <h5 className="text-[9px] font-bold text-slate-700 bg-slate-100 px-3 py-1 rounded-md inline-block uppercase tracking-wider">
                          Langkah Luring / Manual ({editKeberatanManualSteps.length} Langkah)
                        </h5>
                        <button
                          type="button"
                          onClick={() => setEditKeberatanManualSteps([...editKeberatanManualSteps, { title: 'Langkah Baru', desc: 'Deskripsi langkah baru' }])}
                          className="px-2.5 py-1 bg-[#002147] hover:bg-[#00346c] text-white rounded-lg text-[9px] font-bold transition-all flex items-center gap-1 cursor-pointer"
                        >
                          <Plus className="h-3 w-3" /> Tambah Langkah
                        </button>
                      </div>
                      <div className="grid grid-cols-1 md:grid-cols-2 gap-4 max-h-[250px] overflow-y-auto pr-1 border border-slate-200/70 p-3 rounded-2xl bg-slate-50/30">
                        {editKeberatanManualSteps.map((step, idx) => (
                          <div key={idx} className="bg-white p-3 rounded-xl border border-slate-200 space-y-2 relative">
                            <div className="flex items-center justify-between border-b border-slate-100 pb-1.5 mb-1">
                              <span className="text-[9px] font-extrabold text-blue-600 block">Langkah {idx + 1}</span>
                              <div className="flex items-center gap-1">
                                <button
                                  type="button"
                                  disabled={idx === 0}
                                  onClick={() => {
                                    const updated = [...editKeberatanManualSteps];
                                    const hold = updated[idx];
                                    updated[idx] = updated[idx - 1];
                                    updated[idx - 1] = hold;
                                    setEditKeberatanManualSteps(updated);
                                  }}
                                  className="p-1 hover:bg-slate-100 rounded text-slate-400 hover:text-slate-600 disabled:opacity-40"
                                >
                                  <ArrowUp className="h-3 w-3" />
                                </button>
                                <button
                                  type="button"
                                  disabled={idx === editKeberatanManualSteps.length - 1}
                                  onClick={() => {
                                    const updated = [...editKeberatanManualSteps];
                                    const hold = updated[idx];
                                    updated[idx] = updated[idx + 1];
                                    updated[idx + 1] = hold;
                                    setEditKeberatanManualSteps(updated);
                                  }}
                                  className="p-1 hover:bg-slate-100 rounded text-slate-400 hover:text-slate-600 disabled:opacity-40"
                                >
                                  <ArrowDown className="h-3 w-3" />
                                </button>
                                <button
                                  type="button"
                                  onClick={() => {
                                    if (confirm('Hapus langkah luring ini?')) {
                                      setEditKeberatanManualSteps(editKeberatanManualSteps.filter((_, i) => i !== idx));
                                    }
                                  }}
                                  className="p-1 hover:bg-red-50 rounded text-slate-450 hover:text-red-600"
                                >
                                  <Trash2 className="h-3 w-3" />
                                </button>
                              </div>
                            </div>
                            <div className="space-y-1">
                              <label className="text-[9px] font-bold text-slate-500 uppercase block">Judul Langkah</label>
                              <input
                                type="text"
                                value={step.title}
                                onChange={(e) => {
                                  const updated = [...editKeberatanManualSteps];
                                  updated[idx] = { ...updated[idx], title: e.target.value };
                                  setEditKeberatanManualSteps(updated);
                                }}
                                className="w-full rounded-lg border border-slate-200 px-2 py-1 text-xs bg-slate-50 font-bold text-slate-800"
                              />
                            </div>
                            <div className="space-y-1">
                              <label className="text-[9px] font-bold text-slate-500 uppercase block">Penjelasan Langkah</label>
                              <textarea
                                value={step.desc}
                                onChange={(e) => {
                                  const updated = [...editKeberatanManualSteps];
                                  updated[idx] = { ...updated[idx], desc: e.target.value };
                                  setEditKeberatanManualSteps(updated);
                                }}
                                className="w-full rounded-lg border border-slate-200 px-2 py-1 text-[11px] bg-slate-50 text-slate-750 font-medium h-16 resize-none"
                              />
                            </div>
                          </div>
                        ))}
                      </div>
                    </div>

                    {/* Online steps */}
                    <div className="space-y-2 pt-2">
                      <div className="flex justify-between items-center">
                        <h5 className="text-[9px] font-bold text-slate-700 bg-slate-100 px-3 py-1 rounded-md inline-block uppercase tracking-wider">
                          Langkah Daring / Online ({editKeberatanOnlineSteps.length} Langkah)
                        </h5>
                        <button
                          type="button"
                          onClick={() => setEditKeberatanOnlineSteps([...editKeberatanOnlineSteps, { title: 'Langkah Baru', desc: 'Deskripsi langkah baru' }])}
                          className="px-2.5 py-1 bg-[#002147] hover:bg-[#00346c] text-white rounded-lg text-[9px] font-bold transition-all flex items-center gap-1 cursor-pointer"
                        >
                          <Plus className="h-3 w-3" /> Tambah Langkah
                        </button>
                      </div>
                      <div className="grid grid-cols-1 md:grid-cols-2 gap-4 max-h-[250px] overflow-y-auto pr-1 border border-slate-200/70 p-3 rounded-2xl bg-slate-50/30">
                        {editKeberatanOnlineSteps.map((step, idx) => (
                          <div key={idx} className="bg-white p-3 rounded-xl border border-slate-200 space-y-2 relative">
                            <div className="flex items-center justify-between border-b border-slate-100 pb-1.5 mb-1">
                              <span className="text-[9px] font-extrabold text-amber-600 block">Langkah {idx + 1}</span>
                              <div className="flex items-center gap-1">
                                <button
                                  type="button"
                                  disabled={idx === 0}
                                  onClick={() => {
                                    const updated = [...editKeberatanOnlineSteps];
                                    const hold = updated[idx];
                                    updated[idx] = updated[idx - 1];
                                    updated[idx - 1] = hold;
                                    setEditKeberatanOnlineSteps(updated);
                                  }}
                                  className="p-1 hover:bg-slate-100 rounded text-slate-400 hover:text-slate-600 disabled:opacity-40"
                                >
                                  <ArrowUp className="h-3 w-3" />
                                </button>
                                <button
                                  type="button"
                                  disabled={idx === editKeberatanOnlineSteps.length - 1}
                                  onClick={() => {
                                    const updated = [...editKeberatanOnlineSteps];
                                    const hold = updated[idx];
                                    updated[idx] = updated[idx + 1];
                                    updated[idx + 1] = hold;
                                    setEditKeberatanOnlineSteps(updated);
                                  }}
                                  className="p-1 hover:bg-slate-100 rounded text-slate-400 hover:text-slate-600 disabled:opacity-40"
                                >
                                  <ArrowDown className="h-3 w-3" />
                                </button>
                                <button
                                  type="button"
                                  onClick={() => {
                                    if (confirm('Hapus langkah daring ini?')) {
                                      setEditKeberatanOnlineSteps(editKeberatanOnlineSteps.filter((_, i) => i !== idx));
                                    }
                                  }}
                                  className="p-1 hover:bg-red-50 rounded text-slate-450 hover:text-red-600"
                                >
                                  <Trash2 className="h-3 w-3" />
                                </button>
                              </div>
                            </div>
                            <div className="space-y-1">
                              <label className="text-[9px] font-bold text-slate-500 uppercase block">Judul Langkah</label>
                              <input
                                type="text"
                                value={step.title}
                                onChange={(e) => {
                                  const updated = [...editKeberatanOnlineSteps];
                                  updated[idx] = { ...updated[idx], title: e.target.value };
                                  setEditKeberatanOnlineSteps(updated);
                                }}
                                className="w-full rounded-lg border border-slate-200 px-2 py-1 text-xs bg-slate-50 font-bold text-slate-800"
                              />
                            </div>
                            <div className="space-y-1">
                              <label className="text-[9px] font-bold text-slate-500 uppercase block">Penjelasan Langkah</label>
                              <textarea
                                value={step.desc}
                                onChange={(e) => {
                                  const updated = [...editKeberatanOnlineSteps];
                                  updated[idx] = { ...updated[idx], desc: e.target.value };
                                  setEditKeberatanOnlineSteps(updated);
                                }}
                                className="w-full rounded-lg border border-slate-200 px-2 py-1 text-[11px] bg-slate-50 text-slate-750 font-medium h-16 resize-none"
                              />
                            </div>
                          </div>
                        ))}
                      </div>
                    </div>
                  </div>
                )}
              </div>
            ) : ['profil', 'visi-misi', 'tugas-dan-fungsi', 'tugas-fungsi', 'struktur-organisasi-2'].includes(adminEditSlug) ? null : (
              <div className="space-y-1 text-left">
                <label className="text-[10px] font-bold text-slate-600 uppercase tracking-wider block">Konten Isi Halaman (HTML/Markdown)</label>
                <textarea
                  required
                  value={adminEditContent}
                  onChange={(e) => setAdminEditContent(e.target.value)}
                  className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 transition-all bg-slate-50 h-96 font-mono resize-none text-slate-805"
                />
              </div>
            )}

            <div className="flex justify-end gap-3 pt-3 border-t border-slate-100">
              <button
                type="button"
                onClick={() => setEditModalOpen(false)}
                className="px-5 py-2.5 bg-slate-100 hover:bg-slate-200 rounded-xl text-xs font-bold uppercase transition-all cursor-pointer border border-slate-200 text-slate-600"
              >
                Batal
              </button>
              <button
                type="submit"
                className="px-6 py-2.5 bg-[#002147] hover:bg-amber-400 hover:text-[#002147] text-white rounded-xl text-xs font-bold uppercase transition-all cursor-pointer border border-[#002147]/50"
              >
                Simpan Data
              </button>
            </div>
          </form>
        </div>
      ) : (
        <div className="bg-white rounded-3xl border border-slate-200 p-6 shadow-sm space-y-4">
          <div className="flex items-center justify-between">
            <div>
              <h2 className="text-base font-extrabold text-[#002147]">Daftar Halaman Publik</h2>
              <span className="text-[11px] text-slate-400 font-medium block">
                Kelola menu halaman utama (Visi Misi, Sejarah, Profil, Tugas Fungsi, dsb).
              </span>
            </div>
            <button
              onClick={() => openCreateModal('page')}
              className="px-4 py-2 bg-[#002147] hover:bg-[#003166] text-white rounded-xl text-xs font-bold uppercase tracking-wider inline-flex items-center gap-1 cursor-pointer border border-[#002147]/50"
            >
              <Plus className="h-4 w-4" /> Tambah Halaman
            </button>
          </div>

          {adminPages.length > 0 ? (
            <div className="grid grid-cols-1 gap-3 text-left">
              {adminPages.map((page) => (
                <div
                  key={page.id}
                  className="flex flex-col gap-4 sm:flex-row sm:items-center bg-white border border-slate-200 rounded-xl p-4 shadow-sm hover:shadow-md transition-all group"
                >
                  {/* Page Thumbnail */}
                  <div className="flex items-center gap-3 flex-1 min-w-0">
                    {page.cover_image_url ? (
                      <img
                        src={resolveImageUrl(page.cover_image_url)}
                        alt={page.title}
                        className="h-14 w-20 rounded-lg object-cover shrink-0 border border-slate-100"
                      />
                    ) : (
                      <div className="h-14 w-20 rounded-lg bg-slate-100 flex items-center justify-center shrink-0 border border-slate-200">
                        <FileText className="h-5 w-5 text-slate-400" />
                      </div>
                    )}
                    {/* Page Info */}
                    <div className="min-w-0 text-left">
                      <span className="text-[10px] text-slate-400 font-mono font-medium block">/page/{page.slug}</span>
                      <h4 className="text-sm font-bold text-slate-850 truncate mt-0.5" title={page.title}>
                        {page.title}
                      </h4>
                      {page.subtitle ? (
                        <p className="text-[11px] text-slate-400 font-medium truncate" title={page.subtitle}>
                          {page.subtitle}
                        </p>
                      ) : (
                        <p className="text-[11px] text-slate-400 italic">Tidak ada subjudul</p>
                      )}
                    </div>
                  </div>

                  {/* Card Actions */}
                  <div className="flex items-center gap-2 sm:border-l sm:border-slate-200 sm:pl-3 shrink-0">
                    <button
                      onClick={() => openEditModal('page', page)}
                      className="rounded-lg p-2 text-slate-500 hover:bg-slate-100 hover:text-[#002147] transition-colors cursor-pointer border border-transparent"
                      title="Edit Halaman"
                    >
                      <Edit className="h-4 w-4" />
                    </button>
                    <button
                      onClick={() => handleDeleteCrudItem('pages', page.id)}
                      className="rounded-lg p-2 text-slate-500 hover:bg-red-50 hover:text-red-600 transition-colors cursor-pointer border border-transparent"
                      title="Hapus Halaman"
                    >
                      <Trash2 className="h-4 w-4" />
                    </button>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <div className="py-12 text-center text-slate-400 text-xs font-medium bg-slate-50/50 border border-dashed border-slate-200 rounded-2xl">
              Belum ada halaman dinamis dibuat.
            </div>
          )}
        </div>
      )}
    </>
  );
}
