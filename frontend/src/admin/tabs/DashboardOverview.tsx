import React from 'react';
import { 
  FileText, Newspaper, Download, MessageSquare, Landmark, Plus, ArrowRight, Clock, CheckCircle2, AlertCircle 
} from 'lucide-react';
import type { PermohonanTicket } from '../../types';

interface DashboardOverviewProps {
  adminStats: {
    total_pages: number;
    total_posts: number;
    total_downloads: number;
    total_permohonan: number;
  };
  adminPermohonans: PermohonanTicket[];
  setAdminActiveTab: (tab: string) => void;
}

export default function DashboardOverview({ 
  adminStats, 
  adminPermohonans = [], 
  setAdminActiveTab 
}: DashboardOverviewProps) {

  const formatDate = (dateStr: string) => {
    if (!dateStr) return '-';
    try {
      const d = new Date(dateStr);
      return d.toLocaleDateString('id-ID', {
        day: 'numeric',
        month: 'short',
        year: 'numeric'
      });
    } catch (e) {
      return dateStr;
    }
  };

  const getStatusBadge = (status: string) => {
    const s = status ? status.toLowerCase() : 'pending';
    if (s === 'approved' || s === 'disetujui' || s === 'completed') {
      return (
        <span className="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full bg-emerald-50 text-emerald-700 text-[9px] font-extrabold uppercase tracking-wide border border-emerald-100">
          <CheckCircle2 className="h-3 w-3" /> Disetujui
        </span>
      );
    }
    if (s === 'rejected' || s === 'ditolak') {
      return (
        <span className="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full bg-red-50 text-red-700 text-[9px] font-extrabold uppercase tracking-wide border border-red-100">
          <AlertCircle className="h-3 w-3" /> Ditolak
        </span>
      );
    }
    return (
      <span className="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full bg-amber-50 text-amber-700 text-[9px] font-extrabold uppercase tracking-wide border border-amber-100">
        <Clock className="h-3 w-3" /> Pending
      </span>
    );
  };

  const pendingPermohonans = adminPermohonans.filter(p => (p.status || '').toLowerCase() === 'pending');
  const recentPermohonans = adminPermohonans.slice(0, 5);

  return (
    <div className="space-y-6 text-left">
      {/* Welcome Banner */}
      <div className="bg-[#002147] text-white rounded-3xl p-6 lg:p-8 relative overflow-hidden shadow-sm border border-slate-800">
        <div className="absolute right-0 bottom-0 opacity-10 translate-x-6 translate-y-6">
          <Landmark className="h-48 w-48 text-amber-400" />
        </div>
        <div className="relative z-10 space-y-2">
          <span className="bg-amber-400 text-[#002147] text-[9px] font-black px-2.5 py-1 rounded-full uppercase tracking-wider inline-block">
            Dashboard Utama
          </span>
          <h1 className="text-xl lg:text-2xl font-black">Selamat Datang di Workspace PPID UPERTIS</h1>
          <p className="text-xs text-slate-350 font-medium max-w-2xl leading-relaxed">
            Kelola seluruh kebutuhan transparansi data dan regulasi di Universitas Perintis Indonesia. Gunakan menu navigasi di bilah sisi kiri untuk mengunggah dokumen baru, memposting artikel berita, memverifikasi permohonan informasi masuk, atau mengkonfigurasi sistem.
          </p>
        </div>
      </div>

      {/* Stat Cards */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
        {[
          { label: 'Halaman Statis', val: adminStats.total_pages, desc: 'Profil, regulasi & panduan', color: 'text-rose-500 bg-rose-50 border-rose-100', icon: FileText },
          { label: 'Post Berita KIP', val: adminStats.total_posts, desc: 'Publikasi artikel & warta', color: 'text-purple-500 bg-purple-50 border-purple-100', icon: Newspaper },
          { label: 'Berkas Unduhan', val: adminStats.total_downloads, desc: 'Dokumen & regulasi resmi', color: 'text-emerald-500 bg-emerald-50 border-emerald-100', icon: Download },
          { label: 'Total Permohonan', val: adminStats.total_permohonan, desc: `${pendingPermohonans.length} tiket perlu diproses`, color: 'text-teal-500 bg-teal-50 border-teal-100', icon: MessageSquare }
        ].map((card, idx) => {
          const Icon = card.icon;
          return (
            <div key={idx} className="bg-white p-5 rounded-2xl border border-slate-200 shadow-sm flex items-start gap-4 hover:shadow-md transition-shadow">
              <div className={`flex h-11 w-11 shrink-0 items-center justify-center rounded-xl border ${card.color}`}>
                <Icon className="h-5 w-5" />
              </div>
              <div className="space-y-1">
                <span className="text-[10px] font-extrabold text-slate-400 uppercase tracking-wider block">{card.label}</span>
                <span className="text-2xl font-black text-slate-800 block leading-none">{card.val}</span>
                <span className="text-[9px] text-slate-400 font-semibold block">{card.desc}</span>
              </div>
            </div>
          );
        })}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Recent Inquiries List */}
        <div className="lg:col-span-2 bg-white border border-slate-200 rounded-3xl p-6 shadow-sm space-y-4 flex flex-col justify-between">
          <div className="space-y-1">
            <div className="flex items-center justify-between">
              <h3 className="text-sm font-black text-[#002147] uppercase tracking-wider">Permohonan Informasi Terbaru</h3>
              <button 
                onClick={() => setAdminActiveTab('permohonan')}
                className="text-[10px] font-extrabold text-[#002147] hover:text-amber-500 uppercase tracking-wider flex items-center gap-1 cursor-pointer bg-transparent border-0"
              >
                Lihat Semua <ArrowRight className="h-3 w-3" />
              </button>
            </div>
            <p className="text-[11px] text-slate-400 font-medium font-sans">Tanggapan & kelola dokumen permohonan informasi publik dari pemohon.</p>
          </div>

          <div className="overflow-x-auto">
            {recentPermohonans.length > 0 ? (
              <table className="w-full text-left border-collapse text-xs font-semibold text-slate-650">
                <thead>
                  <tr className="border-b border-slate-100 text-[10px] uppercase text-slate-400">
                    <th className="py-2.5 font-extrabold">No. Tiket</th>
                    <th className="py-2.5 font-extrabold">Nama Pemohon</th>
                    <th className="py-2.5 font-extrabold">Tanggal</th>
                    <th className="py-2.5 font-extrabold text-center">Status</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-50">
                  {recentPermohonans.map((perm) => (
                    <tr key={perm.id} className="hover:bg-slate-50/50 transition-colors">
                      <td className="py-3 font-mono font-bold text-slate-805">{perm.ticket_number}</td>
                      <td className="py-3 font-bold text-[#002147]">{perm.name}</td>
                      <td className="py-3 text-slate-400 font-medium">{formatDate(perm.created_at)}</td>
                      <td className="py-3 text-center">{getStatusBadge(perm.status)}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            ) : (
              <div className="py-8 text-center text-slate-450 border border-dashed border-slate-150 rounded-2xl bg-slate-50/40">
                Belum ada berkas permohonan masuk.
              </div>
            )}
          </div>
        </div>

        {/* Shortcuts Panel */}
        <div className="bg-white border border-slate-200 rounded-3xl p-6 shadow-sm space-y-4 flex flex-col justify-between">
          <div className="space-y-1">
            <h3 className="text-sm font-black text-[#002147] uppercase tracking-wider">Aksi Cepat Admin</h3>
            <p className="text-[11px] text-slate-400 font-medium font-sans">Jalan pintas mengunggah konten portal PPID UPERTIS.</p>
          </div>

          <div className="space-y-2">
            {[
              { label: 'Tulis Berita Baru', tab: 'posts', desc: 'Rilis berita keterbukaan informasi', color: 'bg-[#002147]/5 border-[#002147]/10 text-[#002147] hover:bg-[#002147] hover:text-white' },
              { label: 'Unggah Berkas Baru', tab: 'downloads', desc: 'Tambahkan dokumen ke menu download', color: 'bg-emerald-50 border border-emerald-100 text-emerald-700 hover:bg-emerald-600 hover:text-white' },
              { label: 'Kelola Halaman Dinamis', tab: 'pages', desc: 'Edit regulasi, profil, & sengketa', color: 'bg-rose-50 border border-rose-100 text-rose-700 hover:bg-rose-600 hover:text-white' },
              { label: 'Pengaturan Menu Navigasi', tab: 'menu-manager', desc: 'Ubah link header dan dropdown', color: 'bg-sky-50 border border-sky-100 text-sky-700 hover:bg-sky-600 hover:text-white' }
            ].map((act, i) => (
              <button
                key={i}
                onClick={() => setAdminActiveTab(act.tab)}
                className={`w-full p-3 rounded-2xl flex flex-col items-start text-left cursor-pointer transition-all border ${act.color}`}
              >
                <span className="text-xs font-black uppercase tracking-wide leading-none">{act.label}</span>
                <span className="text-[9px] font-medium opacity-80 mt-1 block">{act.desc}</span>
              </button>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
