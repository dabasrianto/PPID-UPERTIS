import React from 'react';

interface ContactFormProps {
  contactName: string;
  setContactName: (val: string) => void;
  contactEmail: string;
  setContactEmail: (val: string) => void;
  contactPhone: string;
  setContactPhone: (val: string) => void;
  contactSubject: string;
  setContactSubject: (val: string) => void;
  contactMessage: string;
  setContactMessage: (val: string) => void;
  contactIsSubmitting: boolean;
  contactSuccess: string | null;
  contactError: string | null;
  handleContactSubmit: (e: React.FormEvent) => void;
}

export default function ContactForm({
  contactName,
  setContactName,
  contactEmail,
  setContactEmail,
  contactPhone,
  setContactPhone,
  contactSubject,
  setContactSubject,
  contactMessage,
  setContactMessage,
  contactIsSubmitting,
  contactSuccess,
  contactError,
  handleContactSubmit
}: ContactFormProps) {
  return (
    <form onSubmit={handleContactSubmit} className="bg-white border border-slate-200 rounded-3xl p-6 lg:p-8 shadow-sm space-y-5 text-left w-full">
      <h3 className="text-xs font-bold text-slate-850 uppercase tracking-wider pb-1 border-b-2 border-amber-400 w-16 mb-2">
        Formulir
      </h3>

      {contactSuccess && (
        <div className="bg-green-50 border border-green-200 text-green-700 text-xs p-4 rounded-2xl font-semibold leading-relaxed">
          {contactSuccess}
        </div>
      )}

      {contactError && (
        <div className="bg-red-50 border border-red-200 text-red-700 text-xs p-4 rounded-2xl font-semibold leading-relaxed">
          {contactError}
        </div>
      )}

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div className="space-y-1.5 text-left">
          <label className="text-[10px] font-bold text-slate-600 uppercase tracking-wider block">
            Nama Lengkap <span className="text-red-500">*</span>
          </label>
          <input
            type="text"
            required
            value={contactName}
            onChange={(e) => setContactName(e.target.value)}
            className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none focus:ring-2 focus:ring-[#002147]/20 focus:border-[#002147] transition-all bg-slate-50 font-semibold text-slate-700"
          />
        </div>
        <div className="space-y-1.5 text-left">
          <label className="text-[10px] font-bold text-slate-600 uppercase tracking-wider block">
            Alamat Email <span className="text-red-500">*</span>
          </label>
          <input
            type="email"
            required
            value={contactEmail}
            onChange={(e) => setContactEmail(e.target.value)}
            className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none focus:ring-2 focus:ring-[#002147]/20 focus:border-[#002147] transition-all bg-slate-50 font-semibold text-slate-700"
          />
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div className="space-y-1.5 text-left">
          <label className="text-[10px] font-bold text-slate-600 uppercase tracking-wider block">
            Nomor Telepon/WA <span className="text-red-500">*</span>
          </label>
          <input
            type="text"
            required
            value={contactPhone}
            onChange={(e) => setContactPhone(e.target.value)}
            className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none focus:ring-2 focus:ring-[#002147]/20 focus:border-[#002147] transition-all bg-slate-50 font-semibold text-slate-700"
          />
        </div>
        <div className="space-y-1.5 text-left">
          <label className="text-[10px] font-bold text-slate-600 uppercase tracking-wider block">
            Subjek Pesan <span className="text-red-500">*</span>
          </label>
          <input
            type="text"
            required
            value={contactSubject}
            onChange={(e) => setContactSubject(e.target.value)}
            className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none focus:ring-2 focus:ring-[#002147]/20 focus:border-[#002147] transition-all bg-slate-50 font-semibold text-slate-700"
          />
        </div>
      </div>

      <div className="space-y-1.5 text-left">
        <label className="text-[10px] font-bold text-slate-600 uppercase tracking-wider block">
          Isi Pesan <span className="text-red-500">*</span>
        </label>
        <textarea
          required
          rows={5}
          value={contactMessage}
          onChange={(e) => setContactMessage(e.target.value)}
          className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none focus:ring-2 focus:ring-[#002147]/20 focus:border-[#002147] transition-all bg-slate-50 font-semibold text-slate-700 resize-none h-32 text-slate-800"
        />
      </div>

      <button
        type="submit"
        disabled={contactIsSubmitting}
        className="px-6 py-2.5 bg-[#002147] hover:bg-amber-450 text-white hover:text-[#002147] rounded-xl text-xs font-bold uppercase tracking-wider transition-all shadow-md flex items-center justify-center gap-2 cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed border border-[#002147]/20"
      >
        {contactIsSubmitting ? 'Mengirim...' : 'Kirim Pesan'}
      </button>
    </form>
  );
}
