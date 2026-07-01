import { FileText } from 'lucide-react';
import { preprocessPostContent } from '../utils/helpers';
import type { PageData, DownloadItem, GalleryItem, PermohonanTicket } from '../types';

// Import sub-pages
import PermohonanInformasi from './PermohonanInformasi';
import Regulasi from './Regulasi';
import DownloadsPage from './DownloadsPage';
import Galeri from './Galeri';
import JadwalLayanan from './JadwalLayanan';
import Kontak from './Kontak';
import InfoPublikTable from './InfoPublikTable';
import Profil from './Profil';
import TugasFungsi from './TugasFungsi';
import KeberatanInformasi from './KeberatanInformasi';
import SengketaInformasi from './SengketaInformasi';
import ZonaIntegrasi from './ZonaIntegrasi';

interface DynamicPageProps {
  activeSlug: string;
  pageData: PageData | null;
  isPageLoading: boolean;
  siteConfig: any;
  navigateToHome: () => void;
  navigateToPage: (slug: string) => void;
  handleNavigation: (href: string) => void;

  // Permohonan
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

  // Regulasi
  expandedRegulasi: Record<string, boolean>;
  toggleRegulasi: (key: string) => void;

  // Download Directory
  activeTab: string;
  setActiveTab: (tab: string) => void;
  searchTerm: string;
  setSearchTerm: (term: string) => void;
  dbDownloads: DownloadItem[];
  incrementDownloadCount: (id: string, fileUrl: string) => void;

  // Gallery
  galleries: GalleryItem[];
  activeLightboxImage: GalleryItem | null;
  setActiveLightboxImage: (item: GalleryItem | null) => void;

  // Contact Form
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

export default function DynamicPage({
  activeSlug,
  pageData,
  isPageLoading,
  siteConfig,
  navigateToHome,
  navigateToPage,
  handleNavigation,

  // Permohonan
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
  formObtainMethod,
  setFormObtainMethod,
  formDeliveryMethod,
  setFormDeliveryMethod,
  formSubmitError,

  // Ticket Tracking
  trackTicketNumber,
  setTrackTicketNumber,
  trackResult,
  isTrackLoading,
  trackError,
  handleTrackSubmit,

  // Regulasi
  expandedRegulasi,
  toggleRegulasi,

  // Download Directory
  activeTab,
  setActiveTab,
  searchTerm,
  setSearchTerm,
  dbDownloads,
  incrementDownloadCount,

  // Gallery
  galleries,
  activeLightboxImage,
  setActiveLightboxImage,

  // Contact Form
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
}: DynamicPageProps) {
  return (
    <div className="container mx-auto px-4 lg:px-8 py-12 lg:py-16 text-left max-w-7xl">
      <button
        onClick={navigateToHome}
        className="inline-flex items-center gap-1 text-xs font-bold text-slate-500 hover:text-[#002147] transition-colors mb-6 uppercase tracking-wider cursor-pointer border-0 bg-transparent"
      >
        ← Kembali ke Beranda
      </button>

      {isPageLoading ? (
        <div className="space-y-4 animate-pulse">
          <div className="h-10 bg-slate-200 rounded-lg w-2/3" />
          <div className="h-4 bg-slate-200 rounded-lg w-1/3" />
          <div className="h-64 bg-slate-200 rounded-2xl w-full" />
          <div className="h-4 bg-slate-200 rounded-lg w-full" />
          <div className="h-4 bg-slate-200 rounded-lg w-5/6" />
        </div>
      ) : pageData ? (
        activeSlug === 'permohonan-informasi' ? (
          <PermohonanInformasi
            siteConfig={siteConfig}
            permohonanSubTab={permohonanSubTab}
            setPermohonanSubTab={setPermohonanSubTab}
            formSubmitSuccess={formSubmitSuccess}
            setFormSubmitSuccess={setFormSubmitSuccess}
            formApplicantType={formApplicantType}
            setFormApplicantType={setFormApplicantType}
            formIdentityNumber={formIdentityNumber}
            setFormIdentityNumber={setFormIdentityNumber}
            formName={formName}
            setFormName={setFormName}
            formEmail={formEmail}
            setFormEmail={setFormEmail}
            formPhone={formPhone}
            setFormPhone={setFormPhone}
            formAddress={formAddress}
            setFormAddress={setFormAddress}
            formAttachmentUrl={formAttachmentUrl}
            setFormAttachmentUrl={setFormAttachmentUrl}
            formIsUploading={formIsUploading}
            formIsSubmitting={formIsSubmitting}
            handlePermohonanSubmit={handlePermohonanSubmit}
            handleAttachmentUpload={handleAttachmentUpload}
            formDetails={formDetails}
            setFormDetails={setFormDetails}
            formPurpose={formPurpose}
            setFormPurpose={setFormPurpose}
            formObtainMethod={formObtainMethod}
            setFormObtainMethod={setFormObtainMethod}
            formDeliveryMethod={formDeliveryMethod}
            setFormDeliveryMethod={setFormDeliveryMethod}
            formSubmitError={formSubmitError}
            trackTicketNumber={trackTicketNumber}
            setTrackTicketNumber={setTrackTicketNumber}
            trackResult={trackResult}
            isTrackLoading={isTrackLoading}
            trackError={trackError}
            handleTrackSubmit={handleTrackSubmit}
          />
        ) : activeSlug === 'regulasi' ? (
          <Regulasi
            pageContent={pageData.content}
            expandedRegulasi={expandedRegulasi}
            toggleRegulasi={toggleRegulasi}
          />
        ) : activeSlug === 'download' ? (
          <DownloadsPage
            dbDownloads={dbDownloads}
            searchTerm={searchTerm}
            setSearchTerm={setSearchTerm}
            incrementDownloadCount={incrementDownloadCount}
            activeTab={activeTab}
            setActiveTab={setActiveTab}
          />
        ) : activeSlug === 'galeri' ? (
          <Galeri
            galleries={galleries}
            activeLightboxImage={activeLightboxImage}
            setActiveLightboxImage={setActiveLightboxImage}
          />
        ) : activeSlug === 'jadwal-layanan-informasi' ? (
          <JadwalLayanan pageData={pageData} />
        ) : activeSlug === 'kontak' ? (
          <Kontak
            siteConfig={siteConfig}
            contactSuccess={contactSuccess}
            contactError={contactError}
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
            handleContactSubmit={handleContactSubmit}
          />
        ) : activeSlug === 'informasi-publik-berkala' ||
          activeSlug === 'info-serta-merta' ||
          activeSlug === 'informasi-tersedia-setiap-saat' ||
          activeSlug === 'informasi-dikecualikan' ? (
          <InfoPublikTable
            activeSlug={activeSlug}
            pageData={pageData}
            searchTerm={searchTerm}
            setSearchTerm={setSearchTerm}
            dbDownloads={dbDownloads}
          />
        ) : activeSlug === 'zona-integrasi' ? (
          <ZonaIntegrasi pageData={pageData} />
        ) : ['tugas-dan-fungsi', 'tugas-fungsi'].includes(activeSlug) ? (
          <TugasFungsi pageData={pageData} />
        ) : activeSlug === 'Permohonan-penyelesaian-sengketa' || activeSlug === 'permohonan-penyelesaian-sengketa' ? (
          <SengketaInformasi pageData={pageData} />
        ) : activeSlug === 'keberatan-informasi' ? (
          <KeberatanInformasi
            pageData={pageData}
            searchTerm={searchTerm}
            setSearchTerm={setSearchTerm}
          />
        ) : ['profil', 'visi-misi', 'maklumat', 'struktur-organisasi-2'].includes(activeSlug) ? (
          <Profil pageData={pageData} />
        ) : (
          <article className="space-y-6 text-left">
            {/* Header Banner for static fallback pages */}
            <div className="bg-[#002147] text-white rounded-3xl p-8 lg:p-12 relative overflow-hidden shadow-lg">
              <div className="absolute right-0 bottom-0 opacity-10 translate-x-12 translate-y-12">
                <FileText className="h-64 w-64 text-amber-400" />
              </div>
              <div className="relative z-10 space-y-3">
                <span className="bg-amber-400 text-[#002147] text-[10px] font-extrabold px-3 py-1 rounded-full uppercase tracking-wider">
                  Informasi Publik
                </span>
                <h1 className="text-3xl lg:text-4xl font-extrabold leading-tight">{pageData.title}</h1>
                {pageData.subtitle && (
                  <p className="text-xs lg:text-sm text-slate-200 leading-relaxed font-medium">
                    {pageData.subtitle}
                  </p>
                )}
              </div>
            </div>

            {pageData.cover_image_url && (
              <img
                src={pageData.cover_image_url}
                alt={pageData.title}
                className="w-full h-72 lg:h-96 object-cover rounded-3xl shadow-sm border border-slate-200 mt-6"
              />
            )}
            <div
              className="border-t border-slate-200 pt-6 html-content text-sm text-slate-700 leading-relaxed space-y-4"
              dangerouslySetInnerHTML={{ __html: preprocessPostContent(pageData.content) }}
            />
          </article>
        )
      ) : (
        <div className="text-center py-12 text-slate-500">
          <p className="font-semibold text-sm">Halaman tidak ditemukan</p>
        </div>
      )}
    </div>
  );
}
