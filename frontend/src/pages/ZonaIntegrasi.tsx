import React from 'react';
import { 
  ShieldCheck, UserPlus, GraduationCap, BookOpen, Award, Users, 
  Wallet, Book, Package, Archive, FileSpreadsheet, ClipboardCheck,
  ExternalLink, Download, FileText, CheckCircle2
} from 'lucide-react';

interface ZonaIntegrasiProps {
  pageData: any;
}

export default function ZonaIntegrasi({ pageData }: ZonaIntegrasiProps) {
  // Parse documents from pageData if available
  let docs: any[] = [];
  try {
    if (pageData?.content) {
      if (pageData.content.startsWith('{')) {
        const parsed = JSON.parse(pageData.content);
        docs = parsed.docs || [];
      } else if (pageData.content.startsWith('[')) {
        docs = JSON.parse(pageData.content);
      }
    }
  } catch (e) {
    console.error('Error parsing documents', e);
  }

  // 11 Sistem Informasi UPERTIS berdasarkan website live ppid.upertis.ac.id
  const sistemInformasiList = [
    {
      title: 'Penerimaan Mahasiswa Baru (PMB)',
      desc: 'Mengelola seluruh alur penerimaan mahasiswa baru secara online, mulai dari pendaftaran, upload berkas, ujian CAT online, hingga kelulusan.',
      url: 'https://pmbonline.upertis.ac.id/',
      icon: UserPlus,
      color: 'from-amber-500/20 to-orange-500/5'
    },
    {
      title: 'Sistem Informasi Akademik (SIAKAD)',
      desc: 'Portal utama pelayanan akademik terintegrasi untuk pengisian KRS online, pemantauan KHS, jadwal kuliah, absensi, hingga pengelolaan data kelulusan.',
      url: 'https://upertis.siakadcloud.com/',
      icon: GraduationCap,
      color: 'from-blue-500/20 to-indigo-500/5'
    },
    {
      title: 'Learning Management System (LMS)',
      desc: 'Media pembelajaran daring interaktif UPERTIS yang menyatukan ruang kelas virtual, modul perkuliahan, pengumpulan tugas, kuis, dan ujian online.',
      url: 'https://upertis.siakadcloud.com/',
      icon: BookOpen,
      color: 'from-emerald-500/20 to-teal-500/5'
    },
    {
      title: 'Tracer Study & Karir Alumni',
      desc: 'Sistem penelusuran lulusan untuk memantau transisi karir alumni, survei keselarasan kurikulum kerja, serta database alumni terintegrasi.',
      url: 'https://upertis.siakadcloud.com/',
      icon: Award,
      color: 'from-rose-500/20 to-red-500/5'
    },
    {
      title: 'Sistem Informasi Kepegawaian (SIMPEG)',
      desc: 'Platform terpusat pengelolaan data administrasi kepegawaian, daftar presensi dosen/staf, kenaikan pangkat, hingga penilaian kinerja internal.',
      url: 'https://upertis.siakadcloud.com/',
      icon: Users,
      color: 'from-purple-500/20 to-violet-500/5'
    },
    {
      title: 'Sistem Informasi Keuangan (SITU)',
      desc: 'Sistem keuangan terpadu UPERTIS yang mengelola pembayaran perkuliahan mahasiswa, sistem penggajian (payroll), neraca akuntansi, dan rincian transaksi.',
      url: 'https://situ.upertis.ac.id/',
      icon: Wallet,
      color: 'from-cyan-500/20 to-sky-500/5'
    },
    {
      title: 'Perpustakaan Digital (SIPUSTAKA)',
      desc: 'Layanan perpustakaan online kampus untuk pencarian katalog buku fisika (OPAC), peminjaman e-book, serta inventarisasi sirkulasi buku.',
      url: 'https://perpustakaan.upertis.ac.id/',
      icon: Book,
      color: 'from-yellow-500/20 to-amber-500/5'
    },
    {
      title: 'Sistem Informasi Aset',
      desc: 'Mengelola data inventarisasi sarana prasarana, pengadaan aset baru kampus, jadwal pemeliharaan fasilitas, hingga pelaporan depresiasi barang.',
      url: 'https://situ.upertis.ac.id/',
      icon: Package,
      color: 'from-indigo-500/20 to-blue-500/5'
    },
    {
      title: 'Repositori Institusi (E-Prints)',
      desc: 'Arsip digital publikasi ilmiah civitas akademika UPERTIS, menyimpan naskah skripsi mahasiswa, laporan tesis, jurnal dosen, serta karya ilmiah lokal.',
      url: 'https://repo.upertis.ac.id/',
      icon: Archive,
      color: 'from-teal-500/20 to-emerald-500/5'
    },
    {
      title: 'Portal Jurnal Ilmiah (OJS)',
      desc: 'Platform publikasi dan review jurnal ilmiah resmi Universitas Perintis Indonesia yang memfasilitasi publikasi riset dosen dan mahasiswa.',
      url: 'https://jurnal.upertis.ac.id/',
      icon: FileSpreadsheet,
      color: 'from-violet-500/20 to-fuchsia-500/5'
    },
    {
      title: 'Sistem Penjaminan Mutu Internal (SPMI)',
      desc: 'Aplikasi pemantauan standar mutu tri dharma perguruan tinggi secara berkala untuk audit mutu internal, evaluasi program studi, dan akreditasi.',
      url: 'https://spmi.upertis.ac.id/',
      icon: ClipboardCheck,
      color: 'from-lime-500/20 to-green-500/5'
    }
  ];

  return (
    <div className="space-y-12 animate-in fade-in duration-300 text-left w-full">
      {/* Premium Hero Banner */}
      <div className="bg-[#002147] text-white rounded-3xl p-8 lg:p-12 relative overflow-hidden shadow-lg border border-white/5">
        <div className="absolute right-0 bottom-0 opacity-10 translate-x-12 translate-y-12">
          <ShieldCheck className="h-72 w-72 text-amber-400" />
        </div>
        <div className="absolute top-0 left-0 w-full h-full bg-gradient-to-br from-transparent via-transparent to-amber-400/10 pointer-events-none" />
        <div className="relative z-10 space-y-4">
          <span className="bg-amber-400 text-[#002147] text-[10px] font-extrabold px-3 py-1 rounded-full uppercase tracking-wider inline-flex items-center gap-1">
            <ShieldCheck className="h-3.5 w-3.5" /> Zona Integrasi (ZI)
          </span>
          <h1 className="text-3xl lg:text-4xl font-extrabold leading-tight">{pageData?.title || 'Zona Integrasi (ZI)'}</h1>
          <p className="text-xs lg:text-sm text-slate-300 leading-relaxed font-medium max-w-3xl">
            {pageData?.subtitle || 'Pembangunan Zona Integrasi Menuju Wilayah Bebas dari Korupsi (WBK) & Wilayah Birokrasi Bersih dan Melayani (WBBM) Universitas Perintis Indonesia.'}
          </p>
        </div>
      </div>

      {/* 3 Pilar Fokus Area */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {[
          { title: 'Birokrasi Bersih & Melayani', desc: 'Meningkatkan kualitas pelayanan publik secara cepat, ramah, transparan, dan akuntabel kepada seluruh civitas akademika dan masyarakat.' },
          { title: 'Bebas dari Korupsi', desc: 'Mencegah terjadinya praktik korupsi, kolusi, gratifikasi, serta nepotisme dengan memperkuat fungsi pengawasan internal di lingkungan kampus.' },
          { title: 'Digitalisasi Tata Laksana', desc: 'Mengurangi tatap muka pelayanan fisik dan beralih ke e-governance dengan menerapkan Sistem Informasi kampus yang andal dan terintegrasi.' }
        ].map((pilar, idx) => (
          <div key={idx} className="bg-white border border-slate-100 p-6 rounded-2xl shadow-sm hover:shadow-md transition-all space-y-2">
            <div className="w-8 h-8 rounded-full bg-amber-50 flex items-center justify-center text-amber-500">
              <CheckCircle2 className="h-5 w-5" />
            </div>
            <h4 className="text-sm font-bold text-[#002147]">{pilar.title}</h4>
            <p className="text-xs text-slate-500 leading-relaxed font-medium">{pilar.desc}</p>
          </div>
        ))}
      </div>

      {/* Grid Sistem Informasi Penunjang */}
      <div className="space-y-6">
        <div>
          <h3 className="text-lg font-extrabold text-[#002147] flex items-center gap-2">
            🚀 Sistem Informasi Penunjang Tata Laksana ZI
          </h3>
          <p className="text-xs text-slate-400 font-medium">
            Digitalisasi proses kerja UPERTIS dalam rangka mendukung kelancaran birokrasi, transparansi akademik, dan efisiensi mutu pelayanan.
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {sistemInformasiList.map((sis, idx) => {
            const IconComponent = sis.icon;
            return (
              <div 
                key={idx} 
                className="group relative bg-white border border-slate-100 rounded-3xl p-6 shadow-sm hover:shadow-xl hover:-translate-y-1.5 transition-all duration-300 flex flex-col justify-between overflow-hidden"
              >
                {/* Background pastel accent */}
                <div className={`absolute -top-12 -right-12 w-28 h-28 rounded-full bg-gradient-to-br ${sis.color} group-hover:scale-150 transition-all duration-500`} />

                <div className="relative z-10 space-y-4 text-left">
                  {/* Icon */}
                  <div className="w-12 h-12 rounded-2xl bg-amber-50 border border-amber-200/50 flex items-center justify-center text-amber-500 group-hover:bg-[#002147] group-hover:text-white transition-colors duration-300">
                    <IconComponent className="h-5 w-5" />
                  </div>
                  <div className="space-y-1.5">
                    <h4 className="text-sm font-extrabold text-[#002147] group-hover:text-amber-500 transition-colors">
                      {sis.title}
                    </h4>
                    <p className="text-xs text-slate-500 leading-relaxed font-medium">
                      {sis.desc}
                    </p>
                  </div>
                </div>

                <div className="relative z-10 pt-6">
                  <a
                    href={sis.url}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="inline-flex items-center gap-1.5 text-[10px] font-extrabold text-[#002147] uppercase tracking-wider bg-slate-50 border border-slate-200 hover:border-amber-400 px-4 py-2 rounded-xl transition-all hover:bg-amber-400 hover:text-[#002147] cursor-pointer"
                  >
                    Kunjungi Layanan <ExternalLink className="h-3 w-3" />
                  </a>
                </div>
              </div>
            );
          })}
        </div>
      </div>

      {/* Dokumen Lampiran Resmi / Regulasi ZI (Jika ada dari admin) */}
      {docs.length > 0 && (
        <div className="space-y-6 border-t border-slate-100 pt-10">
          <div>
            <h3 className="text-lg font-extrabold text-[#002147] flex items-center gap-2">
              📄 Dokumen & Regulasi Pendukung
            </h3>
            <p className="text-xs text-slate-400 font-medium">
              Unduh berkas dokumen, surat keputusan, atau panduan resmi terkait pembangunan Zona Integrasi UPERTIS.
            </p>
          </div>

          <div className="bg-white border border-slate-150 rounded-3xl overflow-hidden shadow-sm">
            <div className="overflow-x-auto">
              <table className="w-full text-left border-collapse">
                <thead>
                  <tr className="bg-slate-50 border-b border-slate-150 text-[10px] font-bold text-slate-500 uppercase tracking-wider">
                    <th className="px-6 py-4">Nama Dokumen / Deskripsi</th>
                    <th className="px-6 py-4 text-right w-44">Tindakan</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-100 text-xs">
                  {docs.map((doc, idx) => (
                    <tr key={idx} className="hover:bg-slate-50/50 transition-all">
                      <td className="px-6 py-4 space-y-1">
                        <div className="font-bold text-[#002147] flex items-center gap-2">
                          <FileText className="h-4 w-4 text-amber-500 shrink-0" />
                          {doc.title || 'Dokumen Tanpa Judul'}
                        </div>
                        {doc.description && (
                          <div className="text-[11px] text-slate-400 font-medium pl-6">
                            {doc.description}
                          </div>
                        )}
                      </td>
                      <td className="px-6 py-4 text-right">
                        <a
                          href={doc.file_url}
                          target="_blank"
                          rel="noopener noreferrer"
                          className="inline-flex items-center gap-1 px-3 py-1.5 bg-amber-400 hover:bg-amber-500 text-[#002147] font-bold rounded-lg transition-all text-[10px] uppercase tracking-wider shadow-sm hover:shadow"
                        >
                          <Download className="h-3 w-3" /> Unduh
                        </a>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
