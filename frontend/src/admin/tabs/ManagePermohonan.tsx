import React from 'react';
import { ChevronLeft, ExternalLink } from 'lucide-react';
import { resolveImageUrl } from '../../utils/helpers';
import type { PermohonanTicket } from '../../types';

interface ManagePermohonanProps {
  editModalOpen: boolean;
  setEditModalOpen: (val: boolean) => void;
  editModalType: string;
  activeEditItem: any;
  setActiveEditItem: (val: any) => void;
  permohonanActionStatus: string;
  setPermohonanActionStatus: (val: string) => void;
  permohonanActionResponse: string;
  setPermohonanActionResponse: (val: string) => void;
  permohonanActionLoading: boolean;
  handlePermohonanActionSubmit: (e: React.FormEvent) => void;
  adminPermohonans: PermohonanTicket[];
  setEditModalType: (val: string) => void;
}

export default function ManagePermohonan({
  editModalOpen,
  setEditModalOpen,
  editModalType,
  activeEditItem,
  setActiveEditItem,
  permohonanActionStatus,
  setPermohonanActionStatus,
  permohonanActionResponse,
  setPermohonanActionResponse,
  permohonanActionLoading,
  handlePermohonanActionSubmit,
  adminPermohonans,
  setEditModalType
}: ManagePermohonanProps) {
  return (
    <>
      {editModalOpen && editModalType === 'permohonan_review' && activeEditItem ? (
        <div className="bg-white rounded-3xl border border-slate-200 p-8 shadow-sm space-y-6 text-left animate-in fade-in duration-200">
          <div className="flex items-center justify-between pb-4 border-b border-slate-100">
            <div>
              <h2 className="text-base font-extrabold text-[#002147]">Tinjau & Tindak Lanjut Permohonan</h2>
              <span className="text-[11px] text-slate-400 font-medium block">Tinjau keselarasan data identitas pemohon dan kirim tanggapan resmi.</span>
            </div>
            <button
              type="button"
              onClick={() => setEditModalOpen(false)}
              className="px-4 py-2 bg-slate-100 hover:bg-slate-200 rounded-xl text-xs font-bold uppercase transition-all cursor-pointer flex items-center gap-1.5 border border-slate-200 text-slate-600"
            >
              <ChevronLeft className="h-4 w-4" /> Kembali
            </button>
          </div>

          <form onSubmit={handlePermohonanActionSubmit} className="space-y-4 text-xs font-medium">
            <div className="bg-slate-50 border border-slate-200 rounded-2xl p-4 space-y-2.5 leading-normal">
              <div className="grid grid-cols-3 gap-2 border-b border-slate-100 pb-2">
                <span className="text-slate-400 font-bold uppercase">No. Tiket:</span>
                <span className="col-span-2 font-mono font-bold text-slate-800">{activeEditItem.ticket_number}</span>
              </div>
              <div className="grid grid-cols-3 gap-2 border-b border-slate-100 pb-2">
                <span className="text-slate-400 font-bold uppercase">Nama Pemohon:</span>
                <span className="col-span-2 text-slate-700 font-bold">{activeEditItem.name}</span>
              </div>
              <div className="grid grid-cols-3 gap-2 border-b border-slate-100 pb-2">
                <span className="text-slate-400 font-bold uppercase">Kategori NIK:</span>
                <span className="col-span-2 text-slate-600 font-mono">{activeEditItem.identity_number}</span>
              </div>
              <div className="grid grid-cols-3 gap-2 border-b border-slate-100 pb-2">
                <span className="text-slate-400 font-bold uppercase">Alamat Email:</span>
                <span className="col-span-2 text-slate-600 font-mono">{activeEditItem.email}</span>
              </div>
              <div className="grid grid-cols-3 gap-2 border-b border-slate-100 pb-2">
                <span className="text-slate-400 font-bold uppercase">No. WhatsApp:</span>
                <span className="col-span-2 text-slate-600 font-mono">{activeEditItem.phone}</span>
              </div>
              <div className="grid grid-cols-3 gap-2 border-b border-slate-100 pb-2">
                <span className="text-slate-400 font-bold uppercase">Alamat Surat:</span>
                <span className="col-span-2 text-slate-600">{activeEditItem.address}</span>
              </div>
              <div className="grid grid-cols-3 gap-2 border-b border-slate-100 pb-2">
                <span className="text-slate-400 font-bold uppercase">Info yang Diminta:</span>
                <span className="col-span-2 text-slate-600 font-bold">{activeEditItem.details}</span>
              </div>
              <div className="grid grid-cols-3 gap-2 border-b border-slate-100 pb-2">
                <span className="text-slate-400 font-bold uppercase">Tujuan Penggunaan:</span>
                <span className="col-span-2 text-slate-600">{activeEditItem.purpose}</span>
              </div>
              <div className="grid grid-cols-3 gap-2 pb-1">
                <span className="text-slate-400 font-bold uppercase">Ktp Scan:</span>
                <span className="col-span-2">
                  {activeEditItem.attachment_url ? (
                    <a
                      href={resolveImageUrl(activeEditItem.attachment_url)}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="text-blue-600 hover:text-blue-800 underline font-bold uppercase tracking-wider inline-flex items-center gap-1"
                    >
                      Buka Dokumen Identitas <ExternalLink className="h-3 w-3" />
                    </a>
                  ) : (
                    <span className="text-slate-400 italic">Tidak melampirkan berkas KTP</span>
                  )}
                </span>
              </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 border-t border-slate-100 pt-4">
              <div className="space-y-1">
                <label className="text-[10px] font-bold text-slate-600 uppercase tracking-wider block">Perbarui Status</label>
                <select
                  value={permohonanActionStatus}
                  onChange={(e) => setPermohonanActionStatus(e.target.value)}
                  className="w-full rounded-xl border border-slate-200 px-3 py-2.5 text-xs focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 transition-all bg-slate-50 font-bold text-slate-805"
                >
                  <option value="pending">Pending (Mengantri)</option>
                  <option value="processing">Processing (Sedang Diproses)</option>
                  <option value="completed">Completed (Selesai/Disetujui)</option>
                  <option value="rejected">Rejected (Ditolak)</option>
                </select>
              </div>
            </div>

            <div className="space-y-1">
              <label className="text-[10px] font-bold text-slate-600 uppercase tracking-wider block">Tulis Tanggapan Resmi Petugas PPID</label>
              <textarea
                required
                placeholder="Tulis tautan berkas dokumen softcopy yang disetujui atau alasan resmi penolakan informasi..."
                value={permohonanActionResponse}
                onChange={(e) => setPermohonanActionResponse(e.target.value)}
                className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 transition-all bg-slate-50 h-28 resize-none text-slate-805"
              />
            </div>

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
                disabled={permohonanActionLoading}
                className="px-6 py-2.5 bg-[#002147] hover:bg-amber-400 hover:text-[#002147] text-white rounded-xl text-xs font-bold uppercase transition-all cursor-pointer border border-[#002147]/50"
              >
                {permohonanActionLoading ? 'Memproses...' : 'Simpan Tanggapan'}
              </button>
            </div>
          </form>
        </div>
      ) : (
        <div className="bg-white rounded-3xl border border-slate-200 p-6 shadow-sm space-y-4">
          <div>
            <h2 className="text-base font-extrabold text-[#002147]">Daftar Tiket Permohonan Masuk</h2>
            <span className="text-[11px] text-slate-400 font-medium block">
              Tinjau data identitas diri (KTP), kaji rincian tujuan pemohon, dan kirim tanggapan/status.
            </span>
          </div>

          <div className="overflow-x-auto border border-slate-100 rounded-2xl">
            <table className="w-full border-collapse text-xs">
              <thead>
                <tr className="bg-slate-50 border-b border-slate-100 text-left text-slate-500 font-bold uppercase tracking-wider">
                  <th className="px-4 py-3">No. Tiket</th>
                  <th className="px-4 py-3">Nama Pemohon</th>
                  <th className="px-4 py-3">Kategori</th>
                  <th className="px-4 py-3">Tanggal Masuk</th>
                  <th className="px-4 py-3 text-center">Status</th>
                  <th className="px-4 py-3 text-center">Aksi</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-50">
                {adminPermohonans.length === 0 ? (
                  <tr>
                    <td colSpan={6} className="py-12 text-center text-slate-400 text-xs font-medium">
                      Belum ada permohonan informasi masuk.
                    </td>
                  </tr>
                ) : (
                  adminPermohonans.map((perm) => (
                    <tr key={perm.id} className="hover:bg-slate-50/50">
                      <td className="px-4 py-3.5 font-mono font-bold text-slate-700">{perm.ticket_number}</td>
                      <td className="px-4 py-3.5 text-left">
                        <span className="font-bold text-slate-800 block">{perm.name}</span>
                        <span className="text-[9px] text-slate-400 font-medium block mt-0.5">{perm.email}</span>
                      </td>
                      <td className="px-4 py-3.5 capitalize text-left">{perm.applicant_type}</td>
                      <td className="px-4 py-3.5 text-slate-450 font-semibold text-left">
                        {new Date(perm.created_at).toLocaleDateString('id-ID', {
                          day: 'numeric',
                          month: 'short',
                          year: 'numeric'
                        })}
                      </td>
                      <td className="px-4 py-3.5 text-center">
                        <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-[10px] font-bold border capitalize ${perm.status === 'completed' ? 'bg-green-50 text-green-700 border-green-200' :
                            perm.status === 'processing' ? 'bg-blue-50 text-blue-700 border-blue-200' :
                              perm.status === 'rejected' ? 'bg-red-50 text-red-700 border-red-200' :
                                'bg-amber-50 text-amber-700 border-amber-200'
                          }`}>
                          {perm.status}
                        </span>
                      </td>
                      <td className="px-4 py-3.5 text-center">
                        <button
                          onClick={() => {
                            setActiveEditItem(perm);
                            setPermohonanActionStatus(perm.status || 'pending');
                            setPermohonanActionResponse(perm.admin_response || '');
                            setEditModalType('permohonan_review');
                            setEditModalOpen(true);
                          }}
                          className="px-3 py-1.5 bg-[#002147] hover:bg-amber-400 hover:text-[#002147] text-white rounded-lg font-bold text-[11px] transition-all cursor-pointer border border-[#002147]/50"
                        >
                          Tinjau
                        </button>
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>
        </div>
      )}
    </>
  );
}
