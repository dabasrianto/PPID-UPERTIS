import React from 'react';
import { ChevronLeft, Plus, FileText, Edit, Trash2 } from 'lucide-react';
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
  handleDeleteCrudItem
}: ManageDownloadsProps) {
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

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-1">
                <label className="text-[10px] font-bold text-slate-600 uppercase tracking-wider block">Kategori Informasi (DIP)</label>
                <select
                  value={adminEditCategory}
                  onChange={(e) => setAdminEditCategory(e.target.value)}
                  className="w-full rounded-xl border border-slate-200 px-3 py-2.5 text-xs focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 transition-all bg-slate-50 font-bold text-slate-805"
                >
                  <option value="ppid-berkala">Informasi Berkala</option>
                  <option value="ppid-setiap-saat">Tersedia Setiap Saat</option>
                  <option value="ppid-serta-merta">Informasi Serta Merta</option>
                  <option value="ppid-dikecualikan">Informasi Dikecualikan</option>
                </select>
              </div>
              <div className="space-y-1">
                <label className="text-[10px] font-bold text-slate-600 uppercase tracking-wider block">
                  File Berkas Berkas (PDF / Gambar)
                </label>
                <div className="flex gap-2">
                  <input
                    type="text"
                    placeholder="/uploads/..."
                    required
                    value={adminEditFileUrl}
                    onChange={(e) => setAdminEditFileUrl(e.target.value)}
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
            <button
              onClick={() => openCreateModal('download')}
              className="px-4 py-2 bg-[#002147] hover:bg-[#003166] text-white rounded-xl text-xs font-bold uppercase tracking-wider inline-flex items-center gap-1 cursor-pointer border border-[#002147]/50"
            >
              <Plus className="h-4 w-4" /> Tambah File
            </button>
          </div>

          <div className="grid grid-cols-1 gap-3">
            {adminDownloads.length === 0 ? (
              <div className="py-12 text-center text-slate-400 text-xs font-medium">
                Belum ada berkas unduhan terdaftar.
              </div>
            ) : (
              adminDownloads.map((dl) => (
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
    </>
  );
}
