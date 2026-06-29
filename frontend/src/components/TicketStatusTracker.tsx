import React from 'react';
import { Search, Clock, CheckCircle, X, HelpCircle, FileText } from 'lucide-react';
import type { PermohonanTicket } from '../types';

interface TicketStatusTrackerProps {
  trackTicketNumber: string;
  setTrackTicketNumber: (val: string) => void;
  trackResult: PermohonanTicket | null;
  isTrackLoading: boolean;
  trackError: string;
  handleTrackSubmit: (e: React.FormEvent) => void;
}

export default function TicketStatusTracker({
  trackTicketNumber,
  setTrackTicketNumber,
  trackResult,
  isTrackLoading,
  trackError,
  handleTrackSubmit
}: TicketStatusTrackerProps) {
  return (
    <div className="space-y-6 text-left">
      <form onSubmit={handleTrackSubmit} className="bg-white border border-slate-200 p-6 rounded-3xl shadow-sm space-y-4">
        <div className="space-y-1">
          <label className="text-[10px] font-bold text-slate-600 uppercase tracking-wider block">
            Nomor Tiket Pengajuan <span className="text-red-500">*</span>
          </label>
          <div className="flex gap-2">
            <div className="relative flex-1">
              <input
                type="text"
                required
                placeholder="Contoh: PPID-17182938475"
                value={trackTicketNumber}
                onChange={(e) => setTrackTicketNumber(e.target.value)}
                className="w-full pl-10 pr-4 py-2.5 bg-slate-50 border border-slate-200 rounded-xl text-xs font-semibold text-slate-700 focus:outline-none focus:ring-2 focus:ring-[#002147]/20 focus:border-[#002147] transition-all"
              />
              <Search className="absolute left-3.5 top-3 h-4 w-4 text-slate-400" />
            </div>
            <button
              type="submit"
              disabled={isTrackLoading}
              className="px-5 py-2.5 bg-[#002147] hover:bg-amber-450 text-white hover:text-[#002147] text-xs font-bold uppercase tracking-wider rounded-xl transition-all shadow-md cursor-pointer disabled:opacity-50 border border-[#002147]/20"
            >
              {isTrackLoading ? 'Melacak...' : 'Lacak'}
            </button>
          </div>
        </div>

        {trackError && (
          <div className="text-xs text-red-650 bg-red-50 border border-red-150 p-3 rounded-xl font-bold">
            {trackError}
          </div>
        )}
      </form>

      {/* Track Result Details Display */}
      {trackResult && (
        <div className="bg-white border border-slate-200 rounded-3xl p-6 shadow-sm space-y-6 animate-in fade-in duration-200">
          <div className="flex flex-wrap items-center justify-between gap-4 border-b border-slate-100 pb-4">
            <div className="space-y-0.5">
              <span className="text-[9px] text-slate-400 font-extrabold uppercase font-mono tracking-wider">Nomor Tiket</span>
              <h4 className="text-sm font-extrabold text-[#002147] font-mono leading-none">{trackResult.ticket_number}</h4>
            </div>
            <div className="text-right">
              <span className="text-[9px] text-slate-400 font-extrabold uppercase font-mono tracking-wider block">Status Pengajuan</span>
              <span
                className={`inline-block px-3 py-1 rounded-full text-[10px] font-black uppercase tracking-wider mt-1 border ${
                  trackResult.status === 'approved'
                    ? 'bg-green-50 text-green-700 border-green-150'
                    : trackResult.status === 'rejected'
                    ? 'bg-red-50 text-red-700 border-red-150'
                    : 'bg-amber-50 text-amber-700 border-amber-150'
                }`}
              >
                {trackResult.status === 'approved'
                  ? 'Selesai / Disetujui'
                  : trackResult.status === 'rejected'
                  ? 'Ditolak'
                  : 'Dalam Proses'}
              </span>
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-xs font-medium">
            <div className="space-y-1">
              <span className="text-[10px] text-slate-400 font-bold block">NAMA PEMOHON</span>
              <p className="text-slate-800 font-bold">{trackResult.name}</p>
            </div>
            <div className="space-y-1">
              <span className="text-[10px] text-slate-400 font-bold block">TANGGAL PENGAJUAN</span>
              <p className="text-slate-800">
                {new Date(trackResult.created_at).toLocaleDateString('id-ID', {
                  day: 'numeric',
                  month: 'long',
                  year: 'numeric',
                  hour: '2-digit',
                  minute: '2-digit'
                })}
              </p>
            </div>
            <div className="space-y-1 md:col-span-2">
              <span className="text-[10px] text-slate-400 font-bold block">RINCIAN INFORMASI YANG DIBUTUHKAN</span>
              <p className="text-slate-700 bg-slate-50 p-3.5 rounded-xl border border-slate-100 italic whitespace-pre-wrap">
                "{trackResult.details}"
              </p>
            </div>
          </div>

          {/* Timeline Tracking Process */}
          <div className="pt-4 border-t border-slate-100 space-y-4">
            <h5 className="text-[10px] font-extrabold text-[#002147] uppercase tracking-wider">
              Linimasa Proses Pengajuan
            </h5>

            <div className="relative pl-6 space-y-6 before:absolute before:left-2 before:top-2 before:bottom-2 before:w-0.5 before:bg-slate-200">
              {/* Step 1: Submit */}
              <div className="relative text-xs">
                <div className="absolute -left-6 top-0.5 h-4.5 w-4.5 rounded-full bg-green-500 border-4 border-white flex items-center justify-center shadow-sm" />
                <div className="space-y-0.5">
                  <span className="font-extrabold text-slate-800 block">Berkas Diterima</span>
                  <span className="text-[10px] text-slate-400 font-medium block">
                    Tiket dibuat otomatis oleh sistem dan masuk antrean verifikasi.
                  </span>
                </div>
              </div>

              {/* Step 2: Processing */}
              <div className="relative text-xs">
                <div
                  className={`absolute -left-6 top-0.5 h-4.5 w-4.5 rounded-full border-4 border-white flex items-center justify-center shadow-sm ${
                    trackResult.status === 'processing' || trackResult.status === 'approved' || trackResult.status === 'rejected'
                      ? 'bg-amber-500'
                      : 'bg-slate-300'
                  }`}
                />
                <div className="space-y-0.5">
                  <span
                    className={`font-extrabold block ${
                      trackResult.status !== 'pending' ? 'text-slate-800' : 'text-slate-400'
                    }`}
                  >
                    Verifikasi Dokumen
                  </span>
                  <span className="text-[10px] text-slate-400 font-medium block">
                    Petugas PPID memeriksa validitas kartu identitas pemohon dan relevansi rincian informasi.
                  </span>
                </div>
              </div>

              {/* Step 3: Finished */}
              <div className="relative text-xs">
                <div
                  className={`absolute -left-6 top-0.5 h-4.5 w-4.5 rounded-full border-4 border-white flex items-center justify-center shadow-sm ${
                    trackResult.status === 'approved'
                      ? 'bg-green-500'
                      : trackResult.status === 'rejected'
                      ? 'bg-red-500'
                      : 'bg-slate-300'
                  }`}
                />
                <div className="space-y-0.5">
                  <span
                    className={`font-extrabold block ${
                      trackResult.status === 'approved'
                        ? 'text-green-700'
                        : trackResult.status === 'rejected'
                        ? 'text-red-700'
                        : 'text-slate-400'
                    }`}
                  >
                    {trackResult.status === 'rejected' ? 'Pengajuan Ditolak' : 'Tanggapan Diberikan'}
                  </span>
                  {trackResult.admin_response ? (
                    <div className="p-3.5 bg-blue-50 border border-blue-150 rounded-xl mt-1.5 space-y-1.5">
                      <span className="text-[9px] text-blue-800 font-extrabold uppercase tracking-wide flex items-center gap-1">
                        <FileText className="h-3 w-3" /> Catatan Resmi Petugas PPID
                      </span>
                      <p className="text-slate-700 leading-relaxed font-semibold italic">
                        "{trackResult.admin_response}"
                      </p>
                    </div>
                  ) : (
                    <span className="text-[10px] text-slate-400 font-medium block">
                      Tanggapan resmi atau tautan unduhan dokumen sedang disiapkan oleh petugas.
                    </span>
                  )}
                </div>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
