import { Clock, Coffee, AlertCircle, ShieldCheck } from 'lucide-react';
import type { PageData } from '../types';

interface JadwalLayananProps {
  pageData: PageData;
}

export default function JadwalLayanan({ pageData }: JadwalLayananProps) {
  let parsed: any = {};
  try {
    parsed = JSON.parse(pageData.content || '{}');
  } catch (e) {
    parsed = {
      senin_kamis_kerja: '08:00 – 16:00 WIB',
      senin_kamis_istirahat: '12:00 – 13:30 WIB',
      jumat_kerja: '08:00 – 16:30 WIB',
      jumat_istirahat: '12:00 – 14:00 WIB',
      sabtu_minggu: 'Sistem Online Tetap Aktif 24/7 (Desk fisik Tutup)',
      offline_guide_1: 'Membawa kartu identitas resmi (KTP untuk perorangan, Akta/SK pendirian untuk instansi).',
      offline_guide_2: 'Loket Pelayanan berada di Lantai 1 Gedung Rektorat Kampus Utama UPERTIS Padang.',
      offline_guide_3: 'Petugas kami siap membantu pencatatan, peninjauan berkas, hingga cetak salinan dokumen.',
      online_guide_1: 'Gunakan menu Permohonan Informasi untuk mengisi formulir digital dan upload KTP.',
      online_guide_2: 'Simpan nomor tiket pengajuan untuk melacak tanggapan admin secara real-time.',
      online_guide_3: 'Jawaban atau link unduhan berkas digital akan dikirimkan langsung ke email/tiket Anda.',
      custom_remarks: pageData.content || ''
    };
  }

  return (
    <div className="space-y-8 animate-in fade-in duration-200 text-left w-full py-6">
      {/* Header Banner */}
      <div className="bg-[#002147] text-white rounded-3xl p-8 lg:p-12 relative overflow-hidden shadow-lg">
        <div className="absolute right-0 bottom-0 opacity-10 translate-x-12 translate-y-12">
          <Clock className="h-64 w-64 text-amber-400" />
        </div>
        <div className="relative z-10 space-y-3">
          <span className="bg-amber-400 text-[#002147] text-[10px] font-extrabold px-3 py-1 rounded-full uppercase tracking-wider">
            Waktu Operasional
          </span>
          <h1 className="text-3xl lg:text-4xl font-extrabold leading-tight">{pageData.title || 'Jadwal Layanan PPID'}</h1>
          <p className="text-xs lg:text-sm text-slate-200 leading-relaxed font-medium">
            Jam kerja operasional Desk Layanan PPID fisik di Kampus UPERTIS serta sistem penerimaan permohonan informasi secara daring.
          </p>
        </div>
      </div>

      {/* Schedule cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {/* Senin - Kamis Card */}
        <div className="bg-white border border-slate-200 rounded-[2rem] p-6 shadow-sm hover:shadow-md transition-all flex flex-col justify-between text-left space-y-4">
          <div className="space-y-2">
            <span className="bg-blue-50 text-blue-700 text-[9px] font-extrabold px-2.5 py-0.5 rounded-full uppercase tracking-wider border border-blue-100 inline-block">
              Senin - Kamis
            </span>
            <h3 className="text-sm font-extrabold text-slate-800">Hari Kerja Operasional</h3>
            <p className="text-[11px] text-slate-400 font-medium">Pelayanan penuh untuk kunjungan fisik dan verifikasi berkas secara offline.</p>
          </div>

          <div className="space-y-3 pt-2 border-t border-slate-50 text-xs">
            <div className="flex items-center gap-2.5 text-slate-700 font-bold">
              <Clock className="h-4.5 w-4.5 text-amber-500 shrink-0" />
              <span>{parsed.senin_kamis_kerja || '08:00 – 16:00 WIB'}</span>
            </div>
            <div className="flex items-center gap-2.5 text-slate-500 font-medium">
              <Coffee className="h-4.5 w-4.5 text-slate-400 shrink-0" />
              <span>Istirahat: {parsed.senin_kamis_istirahat || '12:00 – 13:30 WIB'}</span>
            </div>
          </div>
        </div>

        {/* Jumat Card */}
        <div className="bg-white border border-slate-200 rounded-[2rem] p-6 shadow-sm hover:shadow-md transition-all flex flex-col justify-between text-left space-y-4">
          <div className="space-y-2">
            <span className="bg-blue-50 text-blue-700 text-[9px] font-extrabold px-2.5 py-0.5 rounded-full uppercase tracking-wider border border-blue-100 inline-block">
              Hari Jumat
            </span>
            <h3 className="text-sm font-extrabold text-slate-800">Hari Kerja Terbatas</h3>
            <p className="text-[11px] text-slate-400 font-medium">Pelayanan dengan waktu istirahat ibadah salat jumat yang disesuaikan.</p>
          </div>

          <div className="space-y-3 pt-2 border-t border-slate-50 text-xs">
            <div className="flex items-center gap-2.5 text-slate-700 font-bold">
              <Clock className="h-4.5 w-4.5 text-amber-500 shrink-0" />
              <span>{parsed.jumat_kerja || '08:00 – 16:30 WIB'}</span>
            </div>
            <div className="flex items-center gap-2.5 text-slate-500 font-medium">
              <Coffee className="h-4.5 w-4.5 text-slate-400 shrink-0" />
              <span>Istirahat: {parsed.jumat_istirahat || '12:00 – 14:00 WIB'}</span>
            </div>
          </div>
        </div>

        {/* Sabtu - Minggu Card */}
        <div className="bg-slate-50/50 border border-slate-200 border-dashed rounded-[2rem] p-6 flex flex-col justify-between text-left space-y-4">
          <div className="space-y-2">
            <span className="bg-red-50 text-red-650 text-[9px] font-extrabold px-2.5 py-0.5 rounded-full uppercase tracking-wider border border-red-100 inline-block">
              Sabtu, Minggu & Hari Libur
            </span>
            <h3 className="text-sm font-extrabold text-slate-800">Tutup Operasional</h3>
            <p className="text-[11px] text-slate-400 font-medium">Desk fisik tutup, namun pengajuan online tetap dapat diakses 24 jam.</p>
          </div>

          <div className="space-y-3 pt-2 border-t border-slate-100 text-xs">
            <div className="flex items-center gap-2.5 text-slate-500 font-bold">
              <AlertCircle className="h-4.5 w-4.5 text-red-500 shrink-0" />
              <span>{parsed.sabtu_minggu || 'Sistem Online Tetap Aktif 24/7 (Desk fisik Tutup)'}</span>
            </div>
            <p className="text-[10px] text-slate-400 font-medium leading-relaxed">
              Pengajuan online di hari libur akan dicatat dan diverifikasi pada hari kerja berikutnya.
            </p>
          </div>
        </div>
      </div>

      {/* Guide Cards */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 items-start pt-4">
        {/* Offline Guide */}
        <div className="bg-white border border-slate-200 rounded-3xl p-6 lg:p-8 shadow-sm text-left space-y-4">
          <h3 className="text-xs font-bold text-slate-850 uppercase tracking-wider pb-1.5 border-b-2 border-amber-400 w-24">
            Layanan Fisik (Offline)
          </h3>
          <p className="text-xs text-slate-500 leading-relaxed font-medium">
            Pemohon informasi yang ingin berdiskusi langsung atau berkonsultasi mengenai kelengkapan berkas/dokumen fisik dapat mengunjungi sekretariat PPID:
          </p>
          <div className="p-4 bg-slate-50 rounded-2xl border border-slate-150 text-[11px] text-slate-655 font-medium leading-relaxed space-y-2.5">
            {parsed.offline_guide_1 && (
              <div className="flex items-start gap-2.5">
                <span className="h-2 w-2 rounded-full bg-amber-500 shrink-0 mt-1.5" />
                <span>{parsed.offline_guide_1}</span>
              </div>
            )}
            {parsed.offline_guide_2 && (
              <div className="flex items-start gap-2.5">
                <span className="h-2 w-2 rounded-full bg-amber-500 shrink-0 mt-1.5" />
                <span>{parsed.offline_guide_2}</span>
              </div>
            )}
            {parsed.offline_guide_3 && (
              <div className="flex items-start gap-2.5">
                <span className="h-2 w-2 rounded-full bg-amber-500 shrink-0 mt-1.5" />
                <span>{parsed.offline_guide_3}</span>
              </div>
            )}
          </div>
        </div>

        {/* Online Guide */}
        <div className="bg-white border border-slate-200 rounded-3xl p-6 lg:p-8 shadow-sm text-left space-y-4">
          <h3 className="text-xs font-bold text-slate-850 uppercase tracking-wider pb-1.5 border-b-2 border-amber-400 w-24">
            Layanan Daring (Online)
          </h3>
          <p className="text-xs text-slate-500 leading-relaxed font-medium">
            Pengajuan informasi publik dapat dilakukan kapan saja dan di mana saja secara digital melalui sistem pengajuan online:
          </p>
          <div className="p-4 bg-slate-50 rounded-2xl border border-slate-150 text-[11px] text-slate-655 font-medium leading-relaxed space-y-2.5">
            {parsed.online_guide_1 && (
              <div className="flex items-start gap-2.5">
                <span className="h-2 w-2 rounded-full bg-blue-500 shrink-0 mt-1.5" />
                <span>{parsed.online_guide_1}</span>
              </div>
            )}
            {parsed.online_guide_2 && (
              <div className="flex items-start gap-2.5">
                <span className="h-2 w-2 rounded-full bg-blue-500 shrink-0 mt-1.5" />
                <span>{parsed.online_guide_2}</span>
              </div>
            )}
            {parsed.online_guide_3 && (
              <div className="flex items-start gap-2.5">
                <span className="h-2 w-2 rounded-full bg-blue-500 shrink-0 mt-1.5" />
                <span>{parsed.online_guide_3}</span>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Regulation Time Box */}
      <div className="p-6 bg-gradient-to-br from-blue-50/50 to-indigo-50/20 border border-blue-150 rounded-[2.5rem] space-y-4 shadow-sm text-left">
        <div className="flex items-center gap-2">
          <ShieldCheck className="h-5 w-5 text-amber-500 shrink-0" />
          <h4 className="text-xs font-extrabold text-[#002147] uppercase tracking-wider">
            Jaminan Waktu Pelayanan Informasi (UU KIP)
          </h4>
        </div>
        <p className="text-xs text-slate-600 leading-relaxed font-medium">
          Sesuai dengan ketentuan regulasi UU No. 14 Tahun 2008 tentang Keterbukaan Informasi Publik, waktu pemrosesan permohonan informasi diatur sebagai berikut:
        </p>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 pt-1">
          <div className="p-4 bg-white/75 border border-slate-150 rounded-2xl space-y-1">
            <span className="text-[10px] text-slate-400 font-extrabold block">WAKTU PEMROSESAN UTAMA</span>
            <span className="text-xs font-extrabold text-slate-800">10 Hari Kerja</span>
            <p className="text-[10px] text-slate-500 leading-relaxed mt-1">PPID wajib memberikan pemberitahuan tertulis mengenai status dokumen (diterima/ditolak/dikecualikan) maksimal dalam 10 hari kerja.</p>
          </div>
          <div className="p-4 bg-white/75 border border-slate-150 rounded-2xl space-y-1">
            <span className="text-[10px] text-slate-400 font-extrabold block">WAKTU PERPANJANGAN MAKSIMAL</span>
            <span className="text-xs font-extrabold text-slate-800">+7 Hari Kerja</span>
            <p className="text-[10px] text-slate-500 leading-relaxed mt-1">PPID dapat memperpanjang waktu pemrosesan berkas paling lambat 7 hari kerja berikutnya disertai alasan tertulis yang sah.</p>
          </div>
        </div>
      </div>
    </div>
  );
}
