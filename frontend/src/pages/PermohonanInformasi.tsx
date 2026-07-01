import React from 'react';
import { Users, ChevronDown, ExternalLink } from 'lucide-react';
import TicketStatusTracker from '../components/TicketStatusTracker';
import type { PermohonanTicket } from '../types';

interface PermohonanInformasiProps {
  siteConfig: any;
  permohonanSubTab: string;
  setPermohonanSubTab: (tab: string) => void;
  formSubmitSuccess: any;
  setFormSubmitSuccess: (val: any) => void;
  formApplicantType: string;
  setFormApplicantType: (val: string) => void;
  formIdentityNumber: string;
  setFormIdentityNumber: (val: string) => void;
  formName: string;
  setFormName: (val: string) => void;
  formEmail: string;
  setFormEmail: (val: string) => void;
  formPhone: string;
  setFormPhone: (val: string) => void;
  formAddress: string;
  setFormAddress: (val: string) => void;
  formAttachmentUrl: string;
  setFormAttachmentUrl: (val: string) => void;
  formIsUploading: boolean;
  formIsSubmitting: boolean;
  handlePermohonanSubmit: (e: React.FormEvent) => void;
  handleAttachmentUpload: (e: React.ChangeEvent<HTMLInputElement>) => void;
  formDetails: string;
  setFormDetails: (val: string) => void;
  formPurpose: string;
  setFormPurpose: (val: string) => void;
  formObtainMethod: string;
  setFormObtainMethod: (val: string) => void;
  formDeliveryMethod: string;
  setFormDeliveryMethod: (val: string) => void;
  formSubmitError: string;

  // Ticket Tracking
  trackTicketNumber: string;
  setTrackTicketNumber: (val: string) => void;
  trackResult: PermohonanTicket | null;
  isTrackLoading: boolean;
  trackError: string;
  handleTrackSubmit: (e: React.FormEvent) => void;
}

export default function PermohonanInformasi({
  siteConfig,
  permohonanSubTab,
  setPermohonanSubTab,
  formSubmitSuccess,
  setFormSubmitSuccess,
  formApplicantType,
  setFormApplicantType,
  formIdentityNumber,
  setFormIdentityNumber,
  formName,
  setFormName,
  formEmail,
  setFormEmail,
  formPhone,
  setFormPhone,
  formAddress,
  setFormAddress,
  formAttachmentUrl,
  setFormAttachmentUrl,
  formIsUploading,
  formIsSubmitting,
  handlePermohonanSubmit,
  handleAttachmentUpload,
  formDetails,
  setFormDetails,
  formPurpose,
  setFormPurpose,
  formSubmitError,
  formObtainMethod,
  setFormObtainMethod,
  formDeliveryMethod,
  setFormDeliveryMethod,

  // Ticket Tracking
  trackTicketNumber,
  setTrackTicketNumber,
  trackResult,
  isTrackLoading,
  trackError,
  handleTrackSubmit
}: PermohonanInformasiProps) {
  return (
    <div className="space-y-8 animate-in fade-in duration-200 text-left">
      {/* Banner / Header */}
      <div className="bg-[#002147] text-white rounded-3xl p-8 lg:p-12 text-left relative overflow-hidden shadow-lg">
        <div className="absolute right-0 bottom-0 opacity-10 translate-x-12 translate-y-12">
          <Users className="h-64 w-64" />
        </div>
        <div className="relative z-10 max-w-2xl space-y-3">
          <span className="bg-amber-400 text-[#002147] text-[10px] font-extrabold px-3 py-1 rounded-full uppercase tracking-wider">
            Layanan PPID
          </span>
          <h1 className="text-3xl lg:text-4xl font-extrabold leading-tight">Permohonan Informasi Publik</h1>
          <p className="text-xs lg:text-sm text-slate-200 leading-relaxed font-medium">
            Formulir pengajuan online, pelacakan tiket permohonan, serta panduan tata cara memperoleh informasi resmi.
          </p>
        </div>
      </div>

      {/* Main Grid: Sidebar + Content */}
      <div className="grid grid-cols-1 md:grid-cols-12 gap-8 items-start">
        {/* Sidebar Navigation */}
        <div className="md:col-span-4 bg-white border border-slate-200 rounded-3xl p-4 shadow-sm sticky top-24 space-y-1">
          <span className="text-[10px] font-bold text-slate-400 uppercase tracking-wider px-3 pb-2 block border-b border-slate-100 mb-2">Pilihan Menu</span>
          {[
            { key: 'form', label: 'Formulir Online' },
            { key: 'track', label: 'Cek Status Tiket' },
            { key: 'flow', label: 'Alur & Prosedur' }
          ].map((subTab) => (
            <button
              key={subTab.key}
              type="button"
              onClick={() => setPermohonanSubTab(subTab.key)}
              className={`w-full flex items-center justify-between p-3 rounded-2xl transition-all text-left cursor-pointer ${
                permohonanSubTab === subTab.key
                  ? 'bg-blue-50/80 text-[#002147]'
                  : 'hover:bg-slate-50 text-slate-600 hover:text-[#002147]'
              }`}
            >
              <span className="text-xs font-extrabold">{subTab.label}</span>
              <ChevronDown className="h-4 w-4 opacity-55" />
            </button>
          ))}
        </div>

        {/* Content Box */}
        <div className="md:col-span-8 bg-white border border-slate-200 rounded-3xl p-6 lg:p-8 shadow-sm text-left">
          {/* ─── TAB 1: FORMULIR ONLINE ─── */}
          {permohonanSubTab === 'form' && (
            <div className="space-y-6">
              {formSubmitSuccess && formSubmitSuccess.ticket_number ? (
                <div className="p-6 bg-gradient-to-br from-green-50 to-emerald-50/50 border border-green-200 rounded-3xl space-y-4 text-center max-w-md mx-auto my-4 shadow-sm animate-in zoom-in-95 duration-200">
                  <div className="h-12 w-12 bg-green-500 text-white rounded-full flex items-center justify-center mx-auto text-xl font-bold shadow-md">✓</div>
                  <div className="space-y-1.5">
                    <h3 className="font-extrabold text-base text-[#002147]">Pengajuan Berhasil!</h3>
                    <p className="text-xs text-slate-500 leading-relaxed font-medium">Permohonan informasi Anda telah terdaftar resmi di sistem PPID.</p>
                  </div>
                  <div className="p-4 bg-white rounded-2xl border border-green-150 inline-block">
                    <span className="text-[10px] text-slate-400 block font-mono font-bold tracking-wider">NOMOR TIKET REGISTRASI</span>
                    <span className="text-base font-mono font-extrabold text-slate-700 tracking-wider block mt-1">{formSubmitSuccess.ticket_number}</span>
                  </div>
                  <p className="text-[10px] text-slate-500 leading-relaxed">
                    Catat dan simpan nomor tiket di atas untuk melacak perkembangan permohonan Anda melalui tab <strong>Cek Status Tiket</strong>.
                  </p>
                  <button
                    type="button"
                    onClick={() => setFormSubmitSuccess(null)}
                    className="w-full bg-[#002147] text-white py-2.5 rounded-xl font-bold text-xs hover:bg-amber-400 hover:text-[#002147] transition-all cursor-pointer shadow-md border-0"
                  >
                    Buat Permohonan Baru
                  </button>
                </div>
              ) : (siteConfig?.settings?.permohonan_form_type || 'internal') === 'external' ? (
                <div className="text-center py-10 px-4 bg-slate-50 border border-slate-200/60 rounded-3xl space-y-5 max-w-xl mx-auto my-4 shadow-sm">
                  <div className="h-14 w-14 bg-[#002147]/5 rounded-2xl flex items-center justify-center text-[#002147] mx-auto shadow-sm">
                    <ExternalLink className="h-6 w-6 text-amber-500" />
                  </div>
                  <div className="space-y-1.5">
                    <h3 className="text-base font-extrabold text-[#002147]">Formulir Permohonan Eksternal</h3>
                    <p className="text-xs text-slate-500 leading-relaxed max-w-sm mx-auto">
                      PPID Universitas Perintis Indonesia menggunakan formulir eksternal untuk melayani pengajuan informasi secara daring. Silakan klik tautan di bawah ini untuk membuka formulir pendaftaran:
                    </p>
                  </div>
                  <div className="pt-2">
                    <a
                      href={siteConfig?.settings?.permohonan_link || 'https://forms.gle/b2N4iBRcCFwKEg61A'}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="inline-flex items-center gap-2 rounded-xl bg-amber-400 hover:bg-amber-500 text-[#002147] px-6 py-3.5 font-bold text-xs uppercase tracking-wider transition-all shadow-md hover:shadow-lg hover:-translate-y-0.5 cursor-pointer"
                    >
                      Buka Formulir Pendaftaran <ExternalLink className="h-3.5 w-3.5" />
                    </a>
                  </div>
                </div>
              ) : (
                <form onSubmit={handlePermohonanSubmit} className="space-y-5">
                  <div>
                    <h3 className="font-extrabold text-sm text-[#002147]">Isi Formulir Permohonan</h3>
                    <p className="text-[11px] text-slate-400 font-medium">Lengkapi identitas diri serta detail informasi yang Anda butuhkan.</p>
                  </div>

                  {formSubmitError && (
                    <div className="p-4 bg-red-50 border border-red-200 text-red-750 text-xs rounded-2xl font-bold">
                      {formSubmitError}
                    </div>
                  )}

                  <div className="grid md:grid-cols-2 gap-5 text-left">
                    <div className="space-y-1.5">
                      <label className="text-[10px] font-bold text-slate-600 uppercase tracking-wider block">Kategori Pemohon <span className="text-red-500">*</span></label>
                      <div className="flex gap-2">
                        {[
                          { key: 'perseorangan', label: 'Perseorangan' },
                          { key: 'lembaga', label: 'Lembaga / Instansi' }
                        ].map((opt) => (
                          <button
                            key={opt.key}
                            type="button"
                            onClick={() => setFormApplicantType(opt.key)}
                            className={`flex-1 py-2 rounded-xl text-xs font-bold border transition-colors cursor-pointer ${
                              formApplicantType === opt.key
                                ? 'bg-[#002147] border-[#002147] text-white shadow-sm'
                                : 'bg-slate-50 border-slate-200 text-slate-600 hover:bg-slate-100'
                            }`}
                          >
                            {opt.label}
                          </button>
                        ))}
                      </div>
                    </div>

                    <div className="space-y-1.5">
                      <label className="text-[10px] font-bold text-slate-600 uppercase tracking-wider block">Nomor Identitas (NIK KTP) <span className="text-red-500">*</span></label>
                      <input
                        type="text"
                        required
                        value={formIdentityNumber}
                        onChange={(e) => setFormIdentityNumber(e.target.value)}
                        placeholder="Contoh: 16 digit nomor NIK KTP Anda"
                        className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none focus:ring-2 focus:ring-[#002147]/20 focus:border-[#002147] transition-all bg-slate-50 font-medium text-slate-800"
                      />
                    </div>

                    <div className="space-y-1.5">
                      <label className="text-[10px] font-bold text-slate-600 uppercase tracking-wider block">Nama Lengkap Pemohon <span className="text-red-500">*</span></label>
                      <input
                        type="text"
                        required
                        value={formName}
                        onChange={(e) => setFormName(e.target.value)}
                        placeholder="Sesuai kartu identitas (KTP)"
                        className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none focus:ring-2 focus:ring-[#002147]/20 focus:border-[#002147] transition-all bg-slate-50 font-medium text-slate-805"
                      />
                    </div>

                    <div className="space-y-1.5">
                      <label className="text-[10px] font-bold text-slate-600 uppercase tracking-wider block">Alamat Email Aktif <span className="text-red-500">*</span></label>
                      <input
                        type="email"
                        required
                        value={formEmail}
                        onChange={(e) => setFormEmail(e.target.value)}
                        placeholder="Contoh: nama@domain.com"
                        className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none focus:ring-2 focus:ring-[#002147]/20 focus:border-[#002147] transition-all bg-slate-50 font-mono font-medium text-slate-805"
                      />
                    </div>

                    <div className="space-y-1.5">
                      <label className="text-[10px] font-bold text-slate-600 uppercase tracking-wider block">Nomor WhatsApp Aktif <span className="text-red-500">*</span></label>
                      <input
                        type="text"
                        required
                        value={formPhone}
                        onChange={(e) => setFormPhone(e.target.value)}
                        placeholder="Contoh: 081234567890"
                        className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none focus:ring-2 focus:ring-[#002147]/20 focus:border-[#002147] transition-all bg-slate-50 font-mono font-medium text-slate-850"
                      />
                    </div>

                    <div className="space-y-1.5">
                      <label className="text-[10px] font-bold text-slate-600 uppercase tracking-wider block">Unggah Foto KTP / Surat Kuasa <span className="text-red-500">*</span></label>
                      <div className="flex gap-2">
                        <input
                          type="text"
                          required
                          value={formAttachmentUrl}
                          onChange={(e) => setFormAttachmentUrl(e.target.value)}
                          placeholder="Pilih file scan KTP untuk diunggah..."
                          className="flex-1 rounded-xl border border-slate-200 px-3 py-2 text-xs focus:outline-none focus:ring-2 focus:ring-[#002147]/20 focus:border-[#002147] transition-all bg-slate-50 font-mono font-medium text-slate-805"
                        />
                        <label className="cursor-pointer shrink-0">
                          <div className="px-4 py-2.5 border border-slate-300 hover:border-[#002147] rounded-xl bg-white text-slate-600 hover:text-[#002147] text-xs font-bold transition-all shadow-sm">
                            {formIsUploading ? 'Proses...' : 'Upload'}
                          </div>
                          <input
                            type="file"
                            accept="image/*,application/pdf"
                            onChange={handleAttachmentUpload}
                            className="hidden"
                          />
                        </label>
                      </div>
                    </div>

                    <div className="space-y-1.5 md:col-span-2">
                      <label className="text-[10px] font-bold text-slate-600 uppercase tracking-wider block">Alamat Lengkap Pemohon <span className="text-red-500">*</span></label>
                      <textarea
                        required
                        rows={2}
                        value={formAddress}
                        onChange={(e) => setFormAddress(e.target.value)}
                        placeholder="Contoh: Jl. Adinegoro No.17, Lubuk Buaya, Koto Tangah, Kota Padang"
                        className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none focus:ring-2 focus:ring-[#002147]/20 focus:border-[#002147] transition-all bg-slate-50 font-medium text-slate-805 resize-none"
                      />
                    </div>
                  </div>

                  <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
                    <div className="space-y-1.5 text-left">
                      <label className="text-[10px] font-bold text-slate-400 uppercase tracking-wider block">Rincian Informasi Yang Diminta <span className="text-red-500">*</span></label>
                      <textarea
                        required
                        rows={4}
                        value={formDetails}
                        onChange={(e) => setFormDetails(e.target.value)}
                        placeholder="Tuliskan secara spesifik rincian data / dokumen yang ingin Anda dapatkan..."
                        className="w-full rounded-xl border border-slate-200 bg-slate-50 px-4 py-2.5 text-xs focus:outline-none focus:ring-2 focus:ring-[#002147]/20 focus:border-[#002147] transition-all font-medium text-slate-850"
                      />
                    </div>

                    <div className="space-y-1.5 text-left">
                      <label className="text-[10px] font-bold text-slate-400 uppercase tracking-wider block">Tujuan Penggunaan Informasi <span className="text-red-500">*</span></label>
                      <textarea
                        required
                        rows={4}
                        value={formPurpose}
                        onChange={(e) => setFormPurpose(e.target.value)}
                        placeholder="Jelaskan tujuan penggunaan informasi tersebut secara rinci (cth: untuk penelitian skripsi)..."
                        className="w-full rounded-xl border border-slate-200 bg-slate-50 px-4 py-2.5 text-xs focus:outline-none focus:ring-2 focus:ring-[#002147]/20 focus:border-[#002147] transition-all font-medium text-slate-850"
                      />
                    </div>
                  </div>

                  <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
                    <div className="space-y-1.5 text-left">
                      <label className="text-[10px] font-bold text-slate-600 uppercase tracking-wider block">Cara Memperoleh Informasi <span className="text-red-500">*</span></label>
                      <select
                        value={formObtainMethod}
                        onChange={(e) => setFormObtainMethod(e.target.value)}
                        className="w-full rounded-xl border border-slate-200 bg-slate-50 px-3.5 py-2.5 text-xs focus:outline-none focus:ring-2 focus:ring-[#002147]/20 focus:border-[#002147] font-bold text-slate-800"
                      >
                        <option value="melihat-membaca">Melihat / Membaca / Mendengarkan</option>
                        <option value="mendapatkan-salinan">Mendapatkan Salinan (Softcopy / Hardcopy)</option>
                      </select>
                    </div>

                    <div className="space-y-1.5 text-left">
                      <label className="text-[10px] font-bold text-slate-600 uppercase tracking-wider block">Cara Mengirimkan Informasi <span className="text-red-500">*</span></label>
                      <select
                        value={formDeliveryMethod}
                        onChange={(e) => setFormDeliveryMethod(e.target.value)}
                        className="w-full rounded-xl border border-slate-200 bg-slate-50 px-3.5 py-2.5 text-xs focus:outline-none focus:ring-2 focus:ring-[#002147]/20 focus:border-[#002147] font-bold text-slate-800"
                      >
                        <option value="mengambil-langsung">Mengambil Langsung</option>
                        <option value="email">Melalui Email</option>
                        <option value="whatsapp">Melalui WhatsApp</option>
                        <option value="kurir-pos">Jasa Kurir / Pos</option>
                      </select>
                    </div>
                  </div>

                  <div className="flex justify-end pt-2">
                    <button
                      type="submit"
                      disabled={formIsSubmitting}
                      className="px-8 py-3.5 bg-[#002147] hover:bg-amber-450 hover:text-[#002147] text-white text-xs font-bold rounded-xl uppercase tracking-wider shadow-md hover:shadow-lg transition-all cursor-pointer disabled:opacity-60 border-0"
                    >
                      {formIsSubmitting ? 'Mengirim Data...' : 'Kirim Permohonan'}
                    </button>
                  </div>
                </form>
              )}
            </div>
          )}

          {/* ─── TAB 2: CEK STATUS TIKET ─── */}
          {permohonanSubTab === 'track' && (
            <TicketStatusTracker
              trackTicketNumber={trackTicketNumber}
              setTrackTicketNumber={setTrackTicketNumber}
              trackResult={trackResult}
              isTrackLoading={isTrackLoading}
              trackError={trackError}
              handleTrackSubmit={handleTrackSubmit}
            />
          )}

          {/* ─── TAB 3: ALUR & PROSEDUR ─── */}
          {permohonanSubTab === 'flow' && (
            <div className="space-y-8">
              <div>
                <h3 className="font-extrabold text-sm text-[#002147] uppercase tracking-wider border-b border-slate-100 pb-2">Alur Memperoleh Informasi PPID</h3>
                <p className="text-[11px] text-slate-500 mt-1">Panduan langkah-langkah untuk mendapatkan salinan dokumen publik baik secara manual maupun online.</p>
              </div>

              {/* Alur Timeline Manual */}
              <div className="space-y-4">
                <h4 className="font-extrabold text-xs text-[#002147] bg-slate-100 px-3 py-1.5 rounded-lg inline-block">1. Pemohon Datang Langsung (Manual)</h4>
                <div className="relative border-l-2 border-[#002147]/10 pl-6 ml-3 space-y-6">
                  <div className="relative">
                    <div className="absolute -left-[31px] top-0.5 bg-[#002147] text-white h-4.5 w-4.5 rounded-full flex items-center justify-center font-bold text-[9px]">1</div>
                    <h5 className="font-bold text-slate-800 text-[11px]">Datang Langsung</h5>
                    <p className="text-slate-500 text-[10px] mt-0.5 leading-relaxed">Pemohon datang langsung ke sekretariat PPID UPERTIS dan mengajukan permohonan lisan/tertulis secara manual.</p>
                  </div>
                  <div className="relative">
                    <div className="absolute -left-[31px] top-0.5 bg-[#002147] text-white h-4.5 w-4.5 rounded-full flex items-center justify-center font-bold text-[9px]">2</div>
                    <h5 className="font-bold text-slate-800 text-[11px]">Pencatatan Identitas</h5>
                    <p className="text-slate-500 text-[10px] mt-0.5 leading-relaxed font-medium">Petugas mencatat identitas pemohon (KTP), memeriksa kesesuaian berkas, dan memproses maksud permohonan.</p>
                  </div>
                  <div className="relative">
                    <div className="absolute -left-[31px] top-0.5 bg-[#002147] text-white h-4.5 w-4.5 rounded-full flex items-center justify-center font-bold text-[9px]">3</div>
                    <h5 className="font-bold text-slate-800 text-[11px]">Verifikasi & Registrasi</h5>
                    <p className="text-slate-500 text-[10px] mt-0.5 leading-relaxed font-medium">Petugas memberikan nomor registrasi dan bukti terima fisik. PPID memproses permintaan maksimal 10 hari kerja (+7 hari kerja jika diperpanjang).</p>
                  </div>
                  <div className="relative">
                    <div className="absolute -left-[31px] top-0.5 bg-[#002147] text-white h-4.5 w-4.5 rounded-full flex items-center justify-center font-bold text-[9px]">4</div>
                    <h5 className="font-bold text-slate-800 text-[11px]">Jawaban & Serah Terima</h5>
                    <p className="text-slate-500 text-[10px] mt-0.5 leading-relaxed font-medium">Pemohon mengambil hasil salinan informasi di sekretariat PPID sesuai cara perolehan yang dipilih.</p>
                  </div>
                </div>
              </div>

              {/* Alur Timeline Online */}
              <div className="space-y-4 pt-4 border-t border-slate-100">
                <h4 className="font-extrabold text-xs text-[#002147] bg-slate-100 px-3 py-1.5 rounded-lg inline-block">2. Pemohon Secara Online</h4>
                <div className="relative border-l-2 border-[#002147]/10 pl-6 ml-3 space-y-6">
                  <div className="relative">
                    <div className="absolute -left-[31px] top-0.5 bg-amber-500 text-white h-4.5 w-4.5 rounded-full flex items-center justify-center font-bold text-[9px]">1</div>
                    <h5 className="font-bold text-slate-800 text-[11px]">Mengisi Formulir Online</h5>
                    <p className="text-slate-500 text-[10px] mt-0.5 leading-relaxed font-medium">Pemohon mengisi formulir digital pada tab <strong>Formulir Online</strong> di website PPID dan mengunggah pindaian KTP.</p>
                  </div>
                  <div className="relative">
                    <div className="absolute -left-[31px] top-0.5 bg-amber-500 text-white h-4.5 w-4.5 rounded-full flex items-center justify-center font-bold text-[9px]">2</div>
                    <h5 className="font-bold text-slate-800 text-[11px]">Mendapatkan Nomor Tiket</h5>
                    <p className="text-slate-500 text-[10px] mt-0.5 leading-relaxed font-medium">Sistem secara instan mengeluarkan Nomor Tiket unik sebagai identitas pelacakan.</p>
                  </div>
                  <div className="relative">
                    <div className="absolute -left-[31px] top-0.5 bg-amber-500 text-white h-4.5 w-4.5 rounded-full flex items-center justify-center font-bold text-[9px]">3</div>
                    <h5 className="font-bold text-slate-800 text-[11px]">Proses Verifikasi & Lacak</h5>
                    <p className="text-slate-500 text-[10px] mt-0.5 leading-relaxed font-medium">Petugas melakukan peninjauan. Pemohon dapat mengecek status secara mandiri di tab <strong>Cek Status Tiket</strong> kapan saja.</p>
                  </div>
                  <div className="relative">
                    <div className="absolute -left-[31px] top-0.5 bg-amber-500 text-white h-4.5 w-4.5 rounded-full flex items-center justify-center font-bold text-[9px]">4</div>
                    <h5 className="font-bold text-slate-800 text-[11px]">Jawaban Diterima</h5>
                    <p className="text-slate-500 text-[10px] mt-0.5 leading-relaxed font-medium">Setelah disetujui, jawaban admin atau link berkas softcopy langsung tertera pada halaman cek tiket.</p>
                  </div>
                </div>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
