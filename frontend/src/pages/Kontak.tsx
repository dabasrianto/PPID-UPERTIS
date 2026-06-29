import React from 'react';
import { Mail, MapPin, Phone, Clock } from 'lucide-react';
import ContactForm from '../components/ContactForm';

interface KontakProps {
  siteConfig: any;
  contactSuccess: string | null;
  contactError: string | null;
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
  handleContactSubmit: (e: React.FormEvent) => void;
}

export default function Kontak({
  siteConfig,
  contactSuccess,
  contactError,
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
  handleContactSubmit
}: KontakProps) {
  return (
    <div className="space-y-8 animate-in fade-in duration-200 text-left w-full py-6">
      {/* Header Banner */}
      <div className="bg-[#002147] text-white rounded-3xl p-8 lg:p-12 relative overflow-hidden shadow-lg">
        <div className="absolute right-0 bottom-0 opacity-10 translate-x-12 translate-y-12">
          <Mail className="h-64 w-64" />
        </div>
        <div className="relative z-10 space-y-3">
          <span className="bg-amber-400 text-[#002147] text-[10px] font-extrabold px-3 py-1 rounded-full uppercase tracking-wider">
            Hubungi PPID
          </span>
          <h1 className="text-3xl lg:text-4xl font-extrabold leading-tight">Kontak & Layanan Aduan</h1>
          <p className="text-xs lg:text-sm text-slate-200 leading-relaxed font-medium">
            Kirimkan pertanyaan, saran, aspirasi, atau pengaduan layanan informasi Anda secara langsung melalui formulir elektronik di bawah ini.
          </p>
        </div>
      </div>

      {/* Two Columns Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Left Column: Contact info */}
        <div className="lg:col-span-1 space-y-6">
          <div className="bg-white border border-slate-200 rounded-3xl p-6 shadow-sm space-y-6">
            <h3 className="text-xs font-bold text-slate-850 uppercase tracking-wider pb-1 border-b-2 border-amber-400 w-12">Informasi</h3>

            <div className="space-y-4">
              <div className="flex gap-3 text-left">
                <div className="h-9 w-9 rounded-xl bg-amber-50 text-amber-600 flex items-center justify-center shrink-0 border border-amber-100">
                  <MapPin className="h-4.5 w-4.5" />
                </div>
                <div className="space-y-0.5">
                  <span className="text-[10px] text-slate-400 font-extrabold uppercase tracking-wider block">Kantor Pelayanan</span>
                  <span className="text-xs font-bold text-slate-800 block leading-tight">
                    {siteConfig?.settings?.rektorat_address || 'Lobby Gedung Rektorat Kampus UPERTIS Padang'}
                  </span>
                </div>
              </div>

              <div className="flex gap-3 text-left">
                <div className="h-9 w-9 rounded-xl bg-amber-50 text-amber-600 flex items-center justify-center shrink-0 border border-amber-100">
                  <Mail className="h-4.5 w-4.5" />
                </div>
                <div className="space-y-0.5">
                  <span className="text-[10px] text-slate-400 font-extrabold uppercase tracking-wider block">Email Resmi</span>
                  <span className="text-xs font-bold text-slate-850 block">
                    {siteConfig?.settings?.rektorat_email || 'ppid@upertis.ac.id'}
                  </span>
                </div>
              </div>

              <div className="flex gap-3 text-left">
                <div className="h-9 w-9 rounded-xl bg-amber-50 text-amber-600 flex items-center justify-center shrink-0 border border-amber-100">
                  <Phone className="h-4.5 w-4.5" />
                </div>
                <div className="space-y-0.5">
                  <span className="text-[10px] text-slate-400 font-extrabold uppercase tracking-wider block">WhatsApp Desk</span>
                  <span className="text-xs font-bold text-slate-850 block">
                    {siteConfig?.settings?.rektorat_phone || '+62 821-7212-3211'}
                  </span>
                </div>
              </div>

              <div className="flex gap-3 text-left">
                <div className="h-9 w-9 rounded-xl bg-amber-50 text-amber-600 flex items-center justify-center shrink-0 border border-amber-100">
                  <Clock className="h-4.5 w-4.5" />
                </div>
                <div className="space-y-0.5">
                  <span className="text-[10px] text-slate-400 font-extrabold uppercase tracking-wider block">Jam Kerja</span>
                  <span className="text-xs font-bold text-slate-850 block">
                    {siteConfig?.settings?.service_hours && Array.isArray(siteConfig.settings.service_hours) && siteConfig.settings.service_hours.length > 0
                      ? siteConfig.settings.service_hours.filter((sh: any) => !sh.closed).map((sh: any) => `${sh.day} (${sh.time})`).join(', ')
                      : 'Senin – Jumat (08:00 – 16:00 WIB)'}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Right Column: Contact form */}
        <div className="lg:col-span-2">
          <ContactForm
            contactName={contactName}
            setContactName={setContactName}
            contactEmail={contactEmail}
            setContactEmail={setContactEmail}
            contactPhone={contactPhone}
            setContactPhone={setContactPhone}
            contactSubject={contactSubject}
            setContactSubject={setContactSubject}
            contactMessage={contactMessage}
            setContactMessage={setContactMessage}
            contactIsSubmitting={contactIsSubmitting}
            contactSuccess={contactSuccess}
            contactError={contactError}
            handleContactSubmit={handleContactSubmit}
          />
        </div>
      </div>

      {/* Maps Section */}
      <div className="space-y-4 pt-4">
        <div className="flex items-center gap-2 border-b border-slate-200 pb-3">
          <MapPin className="h-5 w-5 text-amber-500" />
          <h2 className="text-lg font-extrabold text-[#002147]">Peta Lokasi Kampus</h2>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {/* Kampus I Padang Map */}
          <div className="bg-white border border-slate-200 rounded-[2rem] p-4 shadow-sm space-y-3 flex flex-col h-full text-left">
            <div className="px-1.5 pt-1.5">
              <span className="text-[9px] font-extrabold text-amber-600 bg-amber-50 border border-amber-100 px-2 py-0.5 rounded-full uppercase tracking-wider inline-block mb-1.5">
                Kampus Utama (Padang)
              </span>
              <h4 className="text-xs font-bold text-slate-800 leading-tight">Kampus I UPERTIS</h4>
              <p className="text-[10px] text-slate-400 font-medium leading-relaxed mt-0.5">
                {siteConfig?.settings?.rektorat_address || 'Jl. Adinegoro No. 17, Lubuk Buaya, Koto Tangah, Kota Padang, Sumatera Barat 25173'}
              </p>
            </div>
            <div className="h-[500px] w-full rounded-2xl overflow-hidden border border-slate-100 relative bg-slate-50">
              <iframe
                title="Kampus Utama UPERTIS Padang Map"
                src={siteConfig?.settings?.kampus1_map_url || 'https://maps.google.com/maps?q=Universitas+Perintis+Indonesia+Kampus+I+Padang&t=&z=15&ie=UTF8&iwloc=&output=embed'}
                width="100%"
                height="100%"
                style={{ border: 0 }}
                allowFullScreen={true}
                loading="lazy"
              />
            </div>
          </div>

          {/* Kampus II Bukittinggi Map */}
          <div className="bg-white border border-slate-200 rounded-[2rem] p-4 shadow-sm space-y-3 flex flex-col h-full text-left">
            <div className="px-1.5 pt-1.5">
              <span className="text-[9px] font-extrabold text-amber-600 bg-amber-50 border border-amber-100 px-2 py-0.5 rounded-full uppercase tracking-wider inline-block mb-1.5">
                Kampus II (Bukittinggi)
              </span>
              <h4 className="text-xs font-bold text-slate-800 leading-tight">Kampus II UPERTIS</h4>
              <p className="text-[10px] text-slate-400 font-medium leading-relaxed mt-0.5">
                {siteConfig?.settings?.kampus2_address || 'Jl. Raya Bukittinggi - Padang Luar KM. 4, Bukittinggi, Sumatera Barat 26181'}
              </p>
            </div>
            <div className="h-[500px] w-full rounded-2xl overflow-hidden border border-slate-100 relative bg-slate-50">
              <iframe
                title="Kampus II UPERTIS Bukittinggi Map"
                src={siteConfig?.settings?.kampus2_map_url || 'https://maps.google.com/maps?q=Universitas+Perintis+Indonesia+Kampus+II+Bukittinggi&t=&z=15&ie=UTF8&iwloc=&output=embed'}
                width="100%"
                height="100%"
                style={{ border: 0 }}
                allowFullScreen={true}
                loading="lazy"
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
