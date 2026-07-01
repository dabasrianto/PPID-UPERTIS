import React from 'react';
import { ChevronLeft, Plus, FileText, Edit, Trash2, Upload, X, Loader2 } from 'lucide-react';
import type { DownloadItem } from '../../types';

interface ManageDownloadsProps {
  editModalOpen: boolean;
  setEditModalOpen: (val: boolean) => void;
  editModalType: string;
  activeEditItem: any;
  adminEditTitle: string;
  setAdminEditTitle: (val: string) => void;
  adminEditCategory: string;
  setAdminEditCategory: (val: string) => void;
  adminEditFileUrl: string;
  setAdminEditFileUrl: (val: string) => void;
  adminEditDescription: string;
  setAdminEditDescription: (val: string) => void;
  handleSaveCrudItem: (e: React.FormEvent) => void;
  handleAdminUpload: (e: React.ChangeEvent<HTMLInputElement>) => void;
  openCreateModal: (type: string) => void;
  adminDownloads: DownloadItem[];
  openEditModal: (type: string, item: any) => void;
  handleDeleteCrudItem: (table: string, id: string) => void;
  API_BASE_URL: string;
  fetchAdminData: () => void;
  setAdminGlobalMessage: (msg: string) => void;
}

export default function ManageDownloads({
  editModalOpen,
  setEditModalOpen,
  editModalType,
  activeEditItem,
  adminEditTitle,
  setAdminEditTitle,
  adminEditCategory,
  setAdminEditCategory,
  adminEditFileUrl,
  setAdminEditFileUrl,
  adminEditDescription,
  setAdminEditDescription,
  handleSaveCrudItem,
  handleAdminUpload,
  openCreateModal,
  adminDownloads,
  openEditModal,
  handleDeleteCrudItem,
  API_BASE_URL,
  fetchAdminData,
  setAdminGlobalMessage
}: ManageDownloadsProps) {
  // Local state for bulk upload
  const [bulkModalOpen, setBulkModalOpen] = React.useState(false);
  const [bulkFiles, setBulkFiles] = React.useState<Array<{ file: File; title: string; description: string }>>([]);
  const [bulkCategory, setBulkCategory] = React.useState('ppid-berkala');
  const [isUploadingBulk, setIsUploadingBulk] = React.useState(false);
  const [bulkUploadError, setBulkUploadError] = React.useState('');
  // Local state for files within a single download document
  const [docFilesList, setDocFilesList] = React.useState<Array<{ name: string; url: string }>>([]);
  const [isUploadingSingleDocFiles, setIsUploadingSingleDocFiles] = React.useState(false);
  const [selectedFilterCategory, setSelectedFilterCategory] = React.useState<string>('all');


  // Sync parent state with local list when parent opens/changes
  React.useEffect(() => {
    if (!editModalOpen || editModalType !== 'download') {
      setDocFilesList([]);
      return;
    }
    
    if (adminEditFileUrl) {
      if (adminEditFileUrl.trim().startsWith('[')) {
        try {
          const parsed = JSON.parse(adminEditFileUrl);
          if (Array.isArray(parsed)) {
            setDocFilesList(parsed);
            return;
          }
        } catch (e) {
          console.error('Failed to parse file_url JSON:', e);
        }
      }
      
      // Fallback/Legacy: Wrap single URL in list format
      setDocFilesList([{ name: adminEditTitle || 'Dokumen', url: adminEditFileUrl }]);
    } else {
      setDocFilesList([]);
    }
  }, [editModalOpen, editModalType, adminEditFileUrl, adminEditTitle]);

  // Propagate local list changes back to the parent state as serialized JSON
  const updateParentFileUrl = (list: Array<{ name: string; url: string }>) => {
    setDocFilesList(list);
    setAdminEditFileUrl(JSON.stringify(list));
  };

  const handleBulkFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files;
    if (!files || files.length === 0) return;

    const list: Array<{ file: File; title: string; description: string }> = [];
    for (let i = 0; i < files.length; i++) {
      const file = files[i];
      const titleWithoutExt = file.name.replace(/\.[^/.]+$/, "");
      const formattedTitle = titleWithoutExt
        .replace(/[_-]+/g, ' ')
        .replace(/(^\w|\s\w)/g, m => m.toUpperCase());

      list.push({
        file,
        title: formattedTitle,
        description: ''
      });
    }

    setBulkFiles(list);
    setBulkCategory('ppid-berkala');
    setBulkModalOpen(true);
    setBulkUploadError('');
    e.target.value = '';
  };

  const handleSaveBulkDownloads = async (e: React.FormEvent) => {
    e.preventDefault();
    if (bulkFiles.length === 0) {
      setBulkUploadError('Silakan pilih minimal satu berkas.');
      return;
    }

    setIsUploadingBulk(true);
    setBulkUploadError('');
    const token = localStorage.getItem('auth_token');
    if (!token) {
      setBulkUploadError('Sesi login telah berakhir. Silakan login kembali.');
      setIsUploadingBulk(false);
      return;
    }

    try {
      // 1. Upload files to backend
      const formData = new FormData();
      formData.append('folder', 'downloads');
      bulkFiles.forEach(item => {
        formData.append('files', item.file);
      });

      const uploadRes = await fetch(`${API_BASE_URL}/admin/uploads`, {
        method: 'POST',
        headers: { Authorization: `Bearer ${token}` },
        body: formData
      });

      if (!uploadRes.ok) {
        const errData = await uploadRes.json();
        throw new Error(errData.error || 'Gagal mengunggah berkas ke server.');
      }

      const uploadData = await uploadRes.json();
      const uploadedList = uploadData.uploaded || [];
      if (uploadedList.length === 0) {
        throw new Error('Tidak ada berkas yang berhasil diunggah.');
      }

      // 2. Create download record for each file in database
      for (let i = 0; i < bulkFiles.length; i++) {
        const item = bulkFiles[i];
        const match = uploadedList.find((up: any) => up.original_name === item.file.name) || uploadedList[i];
        if (!match) continue;

        const body = {
          title: item.title,
          description: item.description,
          file_url: match.url,
          category: bulkCategory,
          active: true
        };

        const saveRes = await fetch(`${API_BASE_URL}/admin/downloads`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${token}`
          },
          body: JSON.stringify(body)
        });

        if (!saveRes.ok) {
          const errData = await saveRes.json();
          console.error(`Gagal menyimpan database untuk ${item.file.name}:`, errData.error);
        }
      }

      setAdminGlobalMessage(`${bulkFiles.length} berkas unduhan berhasil ditambahkan!`);
      setBulkModalOpen(false);
      fetchAdminData();
    } catch (err: any) {
      console.error(err);
      setBulkUploadError(err.message || 'Terjadi kesalahan saat mengunggah.');
    } finally {
      setIsUploadingBulk(false);
    }
  };

  const filteredDownloads = adminDownloads.filter(dl => 
    selectedFilterCategory === 'all' || dl.category === selectedFilterCategory
  );

  return (

    <>
      {editModalOpen && editModalType === 'download' ? (
        <div className="bg-white rounded-3xl border border-slate-200 p-8 shadow-sm space-y-6 text-left animate-in fade-in duration-200">
          <div className="flex items-center justify-between pb-4 border-b border-slate-100">
            <div>
              <h2 className="text-base font-extrabold text-[#002147]">
                {activeEditItem ? 'Edit Berkas Unduhan' : 'Tambah Berkas Unduhan Baru'}
              </h2>
              <span className="text-[11px] text-slate-400 font-medium block">
                Isi detail data berkas dokumen PPID dengan lengkap.
              </span>
            </div>
            <button
              type="button"
              onClick={() => setEditModalOpen(false)}
              className="px-4 py-2 bg-slate-100 hover:bg-slate-200 rounded-xl text-xs font-bold uppercase transition-all cursor-pointer flex items-center gap-1.5 border border-slate-250 text-slate-700"
            >
              <ChevronLeft className="h-4 w-4" /> Kembali
            </button>
          </div>

          <form onSubmit={handleSaveCrudItem} className="space-y-4">
            <div className="space-y-1">
              <label className="text-[10px] font-bold text-slate-600 uppercase tracking-wider block">
                Judul Dokumen Unduhan <span className="text-red-500">*</span>
              </label>
              <input
                type="text"
                required
                value={adminEditTitle}
                onChange={(e) => setAdminEditTitle(e.target.value)}
                className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 transition-all font-bold text-slate-808"
              />
            </div>

            <div className="space-y-4">
              <div className="space-y-1">
                <label className="text-[10px] font-bold text-slate-600 uppercase tracking-wider block">Kategori Informasi (DIP) <span className="text-red-500">*</span></label>
                <select
                  value={adminEditCategory}
                  onChange={(e) => setAdminEditCategory(e.target.value)}
                  className="w-full rounded-xl border border-slate-200 px-3 py-2.5 text-xs focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 transition-all bg-slate-50 font-bold text-slate-850"
                >
                  <option value="ppid-berkala">Informasi Berkala</option>
                  <option value="ppid-setiap-saat">Tersedia Setiap Saat</option>
                  <option value="ppid-serta-merta">Informasi Serta Merta</option>
                  <option value="ppid-dikecualikan">Informasi Dikecualikan</option>
                </select>
              </div>

              {/* Multiple files list editor */}
              <div className="space-y-3 pt-3 border-t border-slate-100 text-left">
                <div className="flex items-center justify-between pb-2 border-b border-slate-100">
                  <div>
                    <span className="text-[10px] font-extrabold text-[#002147] uppercase tracking-wider block">
                      Daftar Berkas Dokumen <span className="text-red-500">*</span>
                    </span>
                    <span className="text-[11px] text-slate-400 font-medium block">
                      Unggah berkas atau isi link URL dokumen (Bisa lebih dari 1 file dalam dokumen ini).
                    </span>
                  </div>
                  <div className="flex items-center gap-2">
                    <button
                      type="button"
                      onClick={() => updateParentFileUrl([...docFilesList, { name: '', url: '' }])}
                      className="px-3 py-1.5 bg-[#002147] hover:bg-[#003166] text-white rounded-xl text-[10px] font-bold uppercase tracking-wider inline-flex items-center gap-1 cursor-pointer border border-[#002147]/55"
                      disabled={isUploadingSingleDocFiles}
                    >
                      <Plus className="h-3.5 w-3.5" /> Tambah Berkas/Link
                    </button>
                    <label className="px-3 py-1.5 bg-emerald-600 hover:bg-emerald-700 text-white rounded-xl text-[10px] font-bold uppercase tracking-wider inline-flex items-center gap-1.5 cursor-pointer border border-emerald-600 shadow-sm transition-all">
                      {isUploadingSingleDocFiles ? (
                        <>
                          <Loader2 className="h-3.5 w-3.5 animate-spin" /> Uploading...
                        </>
                      ) : (
                        <>
                          <Upload className="h-3.5 w-3.5" /> Upload Berkas (Multi)
                        </>
                      )}
                      <input
                        type="file"
                        multiple
                        accept="application/pdf,image/*,application/msword,application/vnd.openxmlformats-officedocument.wordprocessingml.document"
                        onChange={async (e) => {
                          const files = e.target.files;
                          if (!files || files.length === 0) return;

                          setIsUploadingSingleDocFiles(true);
                          const token = localStorage.getItem('auth_token');
                          if (!token) return;

                          const formData = new FormData();
                          formData.append('folder', 'downloads');
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
                                const newFiles = uploadedList.map((item: any) => {
                                  const nameWithoutExt = item.original_name ? item.original_name.replace(/\.[^/.]+$/, "") : "Dokumen";
                                  const formattedName = nameWithoutExt
                                    .replace(/[_-]+/g, ' ')
                                    .replace(/(^\w|\s\w)/g, (m: string) => m.toUpperCase());
                                  return { name: formattedName, url: item.url };
                                });
                                updateParentFileUrl([...docFilesList, ...newFiles]);
                              }
                            }
                          } catch (err) {
                            console.error(err);
                          } finally {
                            setIsUploadingSingleDocFiles(false);
                            e.target.value = '';
                          }
                        }}
                        className="hidden"
                        disabled={isUploadingSingleDocFiles}
                      />
                    </label>
                  </div>
                </div>

                {/* List dynamic input rows */}
                <div className="space-y-4 max-h-[300px] overflow-y-auto pr-1">
                  {docFilesList.length > 0 ? (
                    docFilesList.map((fileItem, idx) => (
                      <div key={idx} className="p-4 border border-slate-200 rounded-2xl bg-slate-50/50 space-y-3 relative text-left">
                        <div className="absolute top-3.5 right-3.5">
                          <button
                            type="button"
                            onClick={() => updateParentFileUrl(docFilesList.filter((_, i) => i !== idx))}
                            className="p-1.5 text-slate-400 hover:text-red-650 hover:bg-red-50 rounded-lg transition-all border border-transparent"
                            title="Hapus Berkas"
                            disabled={isUploadingSingleDocFiles}
                          >
                            <Trash2 className="h-4 w-4" />
                          </button>
                        </div>

                        <div className="grid grid-cols-1 md:grid-cols-2 gap-3 pr-8 text-left">
                          <div className="space-y-1">
                            <label className="text-[9px] font-bold text-slate-500 uppercase block">Nama/Label Berkas <span className="text-red-500">*</span></label>
                            <input
                              type="text"
                              required
                              value={fileItem.name}
                              onChange={(e) => {
                                const newList = docFilesList.map((d, i) => i === idx ? { ...d, name: e.target.value } : d);
                                updateParentFileUrl(newList);
                              }}
                              className="w-full rounded-xl border border-slate-200 px-3 py-2.5 text-xs focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 bg-white font-bold text-slate-800"
                              placeholder="Contoh: Dokumen Utama / Lampiran I"
                            />
                          </div>

                          <div className="space-y-1">
                            <label className="text-[9px] font-bold text-slate-500 uppercase block">URL Berkas / Path <span className="text-red-500">*</span></label>
                            <div className="flex gap-2">
                              <input
                                type="text"
                                required
                                value={fileItem.url}
                                onChange={(e) => {
                                  const newList = docFilesList.map((d, i) => i === idx ? { ...d, url: e.target.value } : d);
                                  updateParentFileUrl(newList);
                                }}
                                className="flex-1 rounded-xl border border-slate-200 px-3 py-2.5 text-xs focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 bg-white font-mono font-medium text-slate-705"
                                placeholder="/uploads/..."
                              />
                              <label className="px-3 py-2 bg-slate-100 hover:bg-slate-200 hover:text-slate-800 rounded-xl text-xs font-bold transition-all cursor-pointer shrink-0 inline-flex items-center border border-slate-200 text-slate-600">
                                Upload
                                <input
                                  type="file"
                                  multiple
                                  onChange={async (e) => {
                                    const files = e.target.files;
                                    if (!files || files.length === 0) return;

                                    const token = localStorage.getItem('auth_token');
                                    if (!token) return;

                                    const formData = new FormData();
                                    formData.append('folder', 'downloads');
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
                                          const first = uploadedList[0];
                                          const rowName = fileItem.name || (first.original_name ? first.original_name.replace(/\.[^/.]+$/, "").replace(/[_-]+/g, ' ').replace(/(^\w|\s\w)/g, (m: string) => m.toUpperCase()) : "Dokumen");
                                          const updatedList = docFilesList.map((d, i) => i === idx ? { name: rowName, url: first.url } : d);
                                          
                                          if (uploadedList.length > 1) {
                                            const extra = uploadedList.slice(1).map((item: any) => {
                                              const nameWithoutExt = item.original_name ? item.original_name.replace(/\.[^/.]+$/, "") : "Dokumen";
                                              const formattedName = nameWithoutExt
                                                .replace(/[_-]+/g, ' ')
                                                .replace(/(^\w|\s\w)/g, (m: string) => m.toUpperCase());
                                              return { name: formattedName, url: item.url };
                                            });
                                            updatedList.push(...extra);
                                          }
                                          updateParentFileUrl(updatedList);
                                        }
                                      }
                                    } catch (err) {
                                      console.error(err);
                                    } finally {
                                      e.target.value = '';
                                    }
                                  }}
                                  className="hidden"
                                />
                              </label>
                            </div>
                          </div>

                        </div>
                      </div>
                    ))
                  ) : (
                    <div className="py-8 text-center text-slate-400 text-xs font-medium border border-dashed border-slate-200 rounded-2xl">
                      Belum ada berkas terunggah. Silakan klik tombol di atas.
                    </div>
                  )}
                </div>
              </div>

            </div>

            <div className="space-y-1">
              <label className="text-[10px] font-bold text-slate-600 uppercase tracking-wider block">Deskripsi Dokumen</label>
              <textarea
                value={adminEditDescription}
                onChange={(e) => setAdminEditDescription(e.target.value)}
                className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 transition-all bg-slate-50 h-32 font-medium resize-none text-slate-805"
              />
            </div>

            <div className="flex justify-end gap-3 pt-3 border-t border-slate-100">
              <button
                type="button"
                onClick={() => setEditModalOpen(false)}
                className="px-5 py-2.5 bg-slate-100 hover:bg-slate-200 rounded-xl text-xs font-bold uppercase transition-all cursor-pointer text-slate-600 border border-slate-200"
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
              <h2 className="text-base font-extrabold text-[#002147]">Daftar Berkas Unduhan</h2>
              <span className="text-[11px] text-slate-400 font-medium block">
                Tambahkan Surat Keputusan (SK), statuta, Rencana Strategis, atau dokumen PDF lainnya.
              </span>
            </div>
            <div className="flex items-center gap-2">
              <button
                onClick={() => openCreateModal('download')}
                className="px-4 py-2 bg-[#002147] hover:bg-[#003166] text-white rounded-xl text-xs font-bold uppercase tracking-wider inline-flex items-center gap-1 cursor-pointer border border-[#002147]/50 shadow-sm"
              >
                <Plus className="h-4 w-4" /> Tambah File
              </button>
              <label className="px-4 py-2 bg-emerald-600 hover:bg-emerald-700 text-white rounded-xl text-xs font-bold uppercase tracking-wider inline-flex items-center gap-1.5 cursor-pointer border border-emerald-600 shadow-sm transition-all">
                <Upload className="h-4 w-4" /> Bulk Upload
                <input
                  type="file"
                  multiple
                  accept="application/pdf,image/*,application/msword,application/vnd.openxmlformats-officedocument.wordprocessingml.document"
                  onChange={handleBulkFileSelect}
                  className="hidden"
                />
              </label>
            </div>
          </div>

          {/* Category Filter Tabs */}
          <div className="flex flex-wrap items-center gap-1.5 border-b border-slate-100 pb-3">
            {[
              { key: 'all', label: 'Semua Kategori' },
              { key: 'ppid-berkala', label: 'Informasi Berkala' },
              { key: 'ppid-setiap-saat', label: 'Tersedia Setiap Saat' },
              { key: 'ppid-serta-merta', label: 'Informasi Serta Merta' },
              { key: 'ppid-dikecualikan', label: 'Informasi Dikecualikan' }
            ].map((tab) => (
              <button
                key={tab.key}
                type="button"
                onClick={() => setSelectedFilterCategory(tab.key)}
                className={`px-3.5 py-2 rounded-xl text-[10px] font-black uppercase tracking-wider transition-all cursor-pointer border-0 ${
                  selectedFilterCategory === tab.key
                    ? 'bg-[#002147] text-white shadow-sm'
                    : 'bg-slate-100 hover:bg-slate-200 text-slate-650'
                }`}
              >
                {tab.label}
              </button>
            ))}
          </div>

          <div className="grid grid-cols-1 gap-3">
            {filteredDownloads.length === 0 ? (
              <div className="py-12 text-center text-slate-400 text-xs font-medium">
                {selectedFilterCategory === 'all' 
                  ? 'Belum ada berkas unduhan terdaftar.' 
                  : 'Tidak ada berkas unduhan di kategori ini.'}
              </div>
            ) : (
              filteredDownloads.map((dl) => (

                <div
                  key={dl.id}
                  className="flex flex-col gap-4 sm:flex-row sm:items-center bg-white border border-slate-200 rounded-xl p-4 shadow-sm hover:shadow-md transition-all group"
                >
                  {/* Download Info */}
                  <div className="flex items-center gap-3 flex-1 min-w-0">
                    <div className="h-14 w-14 rounded-lg bg-slate-100 flex items-center justify-center shrink-0 border border-slate-100">
                      <FileText className="h-5 w-5 text-slate-450" />
                    </div>
                    <div className="min-w-0 text-left">
                      <span className="text-[10px] text-slate-400 font-medium uppercase block">
                        {dl.category.replace('ppid-', '')}
                      </span>
                      <h4 className="text-sm font-bold text-slate-800 truncate mt-0.5" title={dl.title}>
                        {dl.title}
                      </h4>
                      <span className="text-[10px] text-slate-400 font-medium">{dl.downloads_count || 0} unduhan</span>
                    </div>
                  </div>

                  {/* Card Actions */}
                  <div className="flex items-center gap-2 sm:border-l sm:border-slate-200 sm:pl-3 shrink-0">
                    <button
                      onClick={() => openEditModal('download', dl)}
                      className="rounded-lg p-2 text-slate-500 hover:bg-slate-150 hover:text-[#002147] transition-colors cursor-pointer border border-transparent"
                      title="Edit Dokumen"
                    >
                      <Edit className="h-4 w-4" />
                    </button>
                    <button
                      onClick={() => handleDeleteCrudItem('downloads', dl.id)}
                      className="rounded-lg p-2 text-slate-500 hover:bg-red-50 hover:text-red-650 transition-colors cursor-pointer border border-transparent"
                      title="Hapus Dokumen"
                    >
                      <Trash2 className="h-4 w-4" />
                    </button>
                  </div>
                </div>
              ))
            )}
          </div>
        </div>
      )}

      {/* Bulk Upload Modal */}
      {bulkModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-sm overflow-y-auto animate-in fade-in duration-200">
          <div className="bg-white rounded-3xl border border-slate-200 p-8 shadow-2xl max-w-3xl w-full text-left space-y-6 animate-in zoom-in-95 duration-200 max-h-[90vh] flex flex-col">
            
            {/* Header */}
            <div className="flex items-center justify-between pb-4 border-b border-slate-100 shrink-0">
              <div>
                <h2 className="text-base font-extrabold text-[#002147]">
                  Bulk Upload Berkas Unduhan ({bulkFiles.length} Berkas)
                </h2>
                <span className="text-[11px] text-slate-400 font-medium block">
                  Unggah banyak dokumen sekaligus dan tentukan judul masing-masing berkas.
                </span>
              </div>
              <button
                type="button"
                onClick={() => setBulkModalOpen(false)}
                className="p-2 bg-slate-100 hover:bg-slate-200 rounded-xl text-slate-700 transition-colors cursor-pointer border border-slate-200"
                disabled={isUploadingBulk}
              >
                <X className="h-4 w-4" />
              </button>
            </div>

            {/* Error Message */}
            {bulkUploadError && (
              <div className="p-3.5 bg-red-50 border border-red-200 text-red-700 rounded-2xl text-xs font-semibold shrink-0">
                {bulkUploadError}
              </div>
            )}

            {/* Category Select */}
            <div className="space-y-1 shrink-0">
              <label className="text-[10px] font-bold text-slate-600 uppercase tracking-wider block">
                Kategori Informasi (DIP) Bersama <span className="text-red-500">*</span>
              </label>
              <select
                value={bulkCategory}
                onChange={(e) => setBulkCategory(e.target.value)}
                className="w-full rounded-xl border border-slate-200 px-3 py-2.5 text-xs focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 bg-slate-50 font-bold text-slate-805"
                disabled={isUploadingBulk}
              >
                <option value="ppid-berkala">Informasi Berkala</option>
                <option value="ppid-setiap-saat">Tersedia Setiap Saat</option>
                <option value="ppid-serta-merta">Informasi Serta Merta</option>
                <option value="ppid-dikecualikan">Informasi Dikecualikan</option>
              </select>
            </div>

            {/* Files List Table */}
            <div className="flex-1 overflow-y-auto min-h-[200px] border border-slate-250 rounded-2xl p-4 bg-slate-50/50 space-y-4">
              {bulkFiles.map((item, idx) => (
                <div key={idx} className="bg-white border border-slate-200 rounded-2xl p-4 shadow-sm space-y-3 relative">
                  
                  {/* Remove file button */}
                  <button
                    type="button"
                    onClick={() => {
                      const updated = bulkFiles.filter((_, i) => i !== idx);
                      setBulkFiles(updated);
                      if (updated.length === 0) setBulkModalOpen(false);
                    }}
                    className="absolute top-3.5 right-3.5 p-1.5 text-slate-400 hover:text-red-650 hover:bg-red-50 rounded-lg transition-all"
                    title="Hapus dari daftar"
                    disabled={isUploadingBulk}
                  >
                    <Trash2 className="h-4 w-4" />
                  </button>

                  <div className="flex items-center gap-2 pr-8 text-left">
                    <FileText className="h-5 w-5 text-slate-450 shrink-0" />
                    <span className="text-xs font-semibold text-slate-650 truncate animate-in duration-200" title={item.file.name}>
                      {item.file.name}
                    </span>
                    <span className="text-[10px] text-slate-400 font-medium shrink-0">
                      ({(item.file.size / (1024 * 1024)).toFixed(2)} MB)
                    </span>
                  </div>

                  <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                    <div className="space-y-1">
                      <label className="text-[9px] font-bold text-slate-500 uppercase block">Judul Dokumen Unduhan <span className="text-red-500">*</span></label>
                      <input
                        type="text"
                        required
                        value={item.title}
                        onChange={(e) => {
                          const updated = [...bulkFiles];
                          updated[idx].title = e.target.value;
                          setBulkFiles(updated);
                        }}
                        placeholder="Masukkan judul unduhan..."
                        className="w-full rounded-xl border border-slate-200 px-3 py-2 text-xs focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 bg-white font-bold text-slate-800"
                        disabled={isUploadingBulk}
                      />
                    </div>
                    <div className="space-y-1">
                      <label className="text-[9px] font-bold text-slate-500 uppercase block">Deskripsi Singkat (Opsional)</label>
                      <input
                        type="text"
                        value={item.description}
                        onChange={(e) => {
                          const updated = [...bulkFiles];
                          updated[idx].description = e.target.value;
                          setBulkFiles(updated);
                        }}
                        placeholder="Keterangan singkat berkas..."
                        className="w-full rounded-xl border border-slate-200 px-3 py-2 text-xs focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 bg-white font-medium text-slate-700"
                        disabled={isUploadingBulk}
                      />
                    </div>
                  </div>
                </div>
              ))}
            </div>

            {/* Footer Actions */}
            <div className="flex justify-end gap-3 pt-3 border-t border-slate-100 shrink-0">
              <button
                type="button"
                onClick={() => setBulkModalOpen(false)}
                className="px-5 py-2.5 bg-slate-100 hover:bg-slate-200 rounded-xl text-xs font-bold uppercase transition-all cursor-pointer text-slate-600 border border-slate-200"
                disabled={isUploadingBulk}
              >
                Batal
              </button>
              <button
                type="button"
                onClick={handleSaveBulkDownloads}
                className="px-6 py-2.5 bg-[#002147] hover:bg-amber-400 hover:text-[#002147] text-white rounded-xl text-xs font-bold uppercase transition-all cursor-pointer border border-[#002147]/50 inline-flex items-center gap-2"
                disabled={isUploadingBulk}
              >
                {isUploadingBulk ? (
                  <>
                    <Loader2 className="h-4 w-4 animate-spin" /> Proses Upload...
                  </>
                ) : (
                  'Simpan Semua'
                )}
              </button>
            </div>

          </div>
        </div>
      )}
    </>
  );
}
