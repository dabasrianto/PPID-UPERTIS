import { useState, useEffect, useMemo, useRef, useCallback } from 'react';
import type { Post, PageData, DownloadItem, GalleryItem, PermohonanTicket, UserProfile } from './types';
import { resolveImageUrl } from './utils/helpers';

// Shared Layout Components
import Navbar from './components/Navbar';
import Footer from './components/Footer';
import ChatWidget from './components/ChatWidget';

// Public Pages
import Home from './pages/Home';
import News from './pages/News';
import NewsDetail from './pages/NewsDetail';
import DynamicPage from './pages/DynamicPage';

// Admin Core & Pages/Tabs
import AdminLogin from './admin/AdminLogin';
import AdminLayout from './admin/AdminLayout';
import DashboardOverview from './admin/tabs/DashboardOverview';
import ManageSlider from './admin/tabs/ManageSlider';
import ManageDownloads from './admin/tabs/ManageDownloads';
import ManageGallery from './admin/tabs/ManageGallery';
import ManagePages from './admin/tabs/ManagePages';
import ManageNews from './admin/tabs/ManageNews';
import ManagePermohonan from './admin/tabs/ManagePermohonan';
import PortalSettings from './admin/tabs/PortalSettings';
import MenuManager from './admin/tabs/MenuManager';

const API_BASE_URL = import.meta.env.VITE_API_URL || '/api/v1';

const migrateMenuStructure = (menu: any[]): any[] => {
  if (!menu || !Array.isArray(menu) || menu.length === 0) {
    return [];
  }

  // Check if we need migration (e.g. if none of the items has label "Home" or href "/")
  const hasHome = menu.some(item => 
    (item.label && (item.label.toLowerCase() === 'home' || item.label.toLowerCase() === 'beranda')) ||
    (item.href === '/' || item.href === 'home')
  );

  const hasNews = menu.some(item =>
    (item.label && item.label.toLowerCase().includes('berita')) ||
    (item.href === 'berita')
  );

  if (hasHome && hasNews) {
    return menu;
  }

  const migrated: any[] = [];

  // 1. Prepend Home if missing
  if (!hasHome) {
    migrated.push({ label: 'Home', href: '/', type: 'link' });
  }

  // 2. Add the existing items (normalize old group format to dropdown type)
  menu.forEach(item => {
    if (item.group && !item.label) {
      migrated.push({
        label: item.group,
        type: 'dropdown',
        items: item.items || []
      });
    } else {
      migrated.push(item);
    }
  });

  // 3. Append News, Gallery, Download, Hubungi Kami if missing
  const hasGallery = menu.some(item => item.href === 'galeri' || (item.label && item.label.toLowerCase() === 'galeri'));
  const hasDownload = menu.some(item => item.href === 'download' || (item.label && item.label.toLowerCase() === 'download'));
  const hasContact = menu.some(item => item.href === 'kontak' || (item.label && item.label.toLowerCase().includes('hubungi')));

  if (!hasNews) {
    migrated.push({ label: 'Berita PPID', href: 'berita', type: 'link' });
  }
  if (!hasGallery) {
    migrated.push({ label: 'Galeri', href: 'galeri', type: 'link' });
  }
  if (!hasDownload) {
    migrated.push({ label: 'Download', href: 'download', type: 'link' });
  }
  if (!hasContact) {
    migrated.push({ label: 'Hubungi Kami', href: 'kontak', type: 'link', isHighlight: true });
  }

  return migrated;
};

const defaultManualSteps = [
  { title: "Pemohon Informasi", desc: "Pemohon datang langsung ke sekretariat PPID-UPERTIS untuk mengajukan keberatan atas informasi publik yang diterima." },
  { title: "Mengisi Form Keberatan", desc: "Petugas informasi menerima berkas pengajuan dan pemohon mengisi formulir keberatan dengan melampirkan identitas diri berupa KTP/SIM atau surat tugas dari institusi lainnya." },
  { title: "Menerima Tanda Terima", desc: "Petugas informasi memberikan tanda bukti resmi bahwa pemohon telah mengajukan keberatan atas layanan informasi publik." },
  { title: "Keberatan Di Proses", desc: "PPID Pelaksana meneruskan pengaduan ke Rektor/PPID Utama UPERTIS untuk diproses selambat-lambatnya 30 hari kerja." },
  { title: "Penyerahan Jawaban", desc: "Jawaban/tanggapan tertulis atas keberatan disampaikan kepada pemohon melalui PPID Pelaksana." },
  { title: "Layanan Selesai", desc: "Prosedur pengajuan keberatan informasi secara manual selesai dilakukan." }
];

const defaultOnlineSteps = [
  { title: "Mengunjungi Website", desc: "Pemohon mengunjungi website resmi PPID Universitas Perintis Indonesia di ppid.upertis.ac.id." },
  { title: "Mengunduh Form", desc: "Pemohon mengunduh berkas formulir pengajuan keberatan (.doc) dari menu unduhan di website." },
  { title: "Mengirimkan Berkas", desc: "Pemohon mengirim formulir terisi beserta scan identitas (KTP/SIM) ke email resmi PPID UPERTIS di ppidcare@upertis.ac.id." },
  { title: "Menerima Tanda Terima", desc: "Petugas informasi mengirimkan konfirmasi tanda terima pengajuan keberatan via email." },
  { title: "Permohonan Di Proses", desc: "Berkas keberatan diteruskan ke Rektor UPERTIS untuk ditinjau dan diproses dalam kurun waktu 30 hari kerja." },
  { title: "Penyerahan Jawaban", desc: "Jawaban keputusan dari Rektor UPERTIS dikirimkan kepada pemohon melalui PPID Pelaksana." },
  { title: "Layanan Selesai", desc: "Prosedur pengajuan keberatan informasi secara online selesai dilakukan." }
];

const defaultTugasList = [
  { title: "Pengelolaan & Penyimpanan Arsip", desc: "Mengelola dan menyimpan dokumen serta arsip informasi publik agar terstruktur, aman, dan mudah diakses." },
  { title: "Penyediaan & Penerbitan Informasi", desc: "Menyediakan, memberikan, dan menerbitkan dokumen informasi publik secara proaktif kepada masyarakat luas." },
  { title: "Pengklasifikasian Informasi Publik", desc: "Melakukan klasifikasi berkala atas informasi (berkala, serta-merta, setiap saat, dan dikecualikan) sesuai regulasi." },
  { title: "Uji Konsekuensi Informasi Dikecualikan", desc: "Melakukan analisis dampak dan uji konsekuensi hukum yang ketat sebelum mengecualikan akses suatu informasi." },
  { title: "Penyusunan Daftar Informasi Publik (DIP)", desc: "Menyusun, memutakhirkan, dan mempublikasikan Daftar Informasi Publik (DIP) UPERTIS secara komprehensif." },
  { title: "Pengembangan Sistem Layanan Informasi", desc: "Membangun, mengoperasikan, dan mengawasi jalannya sistem IT maupun konvensional untuk layanan informasi publik." },
  { title: "Penyelesaian Sengketa Informasi Internal", desc: "Menyelesaikan setiap perselisihan permohonan informasi melalui jalur mediasi internal sebelum ke Komisi Informasi." },
  { title: "Pelaporan Berkala Kinerja Layanan", desc: "Membuat laporan rutin berkala pelaksanaan layanan informasi publik untuk diserahkan kepada Atasan PPID." }
];

const defaultFungsiList = [
  {
    title: "Fungsi Penyimpanan dan Pendokumentasian",
    items: [
      "Mengarsipkan setiap dokumen dan berkas informasi publik yang dihasilkan oleh seluruh unit kerja Universitas.",
      "Menyimpan basis data secara digital maupun fisik agar terawat, mudah dilacak, dan bebas risiko hilang."
    ]
  },
  {
    title: "Fungsi Penyediaan Informasi",
    items: [
      "Menyediakan akses informasi publik secara cepat, tepat waktu, transparan, dan menggunakan skema sederhana.",
      "Menyusun struktur Daftar Informasi Publik (DIP) yang selaras dengan ketentuan peraturan perundang-undangan."
    ]
  },
  {
    title: "Fungsi Pelayanan Informasi",
    items: [
      "Memberikan pemenuhan permohonan informasi kepada pemohon sesuai standar operating procedure (SOP) PPID.",
      "Menyediakan sarana fisik (meja layanan ppid) dan digital (situs web portal PPID, chatbot interaktif, desk online)."
    ]
  },
  {
    title: "Fungsi Koordinasi Internal & Eksternal",
    items: [
      "Berkoordinasi aktif dengan PPID Pelaksana di tingkat fakultas dan biro untuk menghimpun data secara berkala.",
      "Membangun komunikasi harmonis dengan Tim Pertimbangan Hukum, Rektorat (Atasan PPID), dan Petugas Informasi."
    ]
  },
  {
    title: "Fungsi Pengendalian dan Pengawasan",
    items: [
      "Mengontrol dan mengawasi seluruh rantai proses pelayanan dokumen informasi publik agar memenuhi standar mutu KIP.",
      "Melakukan evaluasi kepuasan layanan dan audit kepatuhan keterbukaan informasi di internal perguruan tinggi."
    ]
  }
];

function App() {
  const [currentPage, setCurrentPage] = useState<string>('home'); // 'home', 'page/:slug', 'berita', 'berita-detail', 'admin'
  const [activeSlug, setActiveSlug] = useState<string>('');
  const [pageData, setPageData] = useState<PageData | null>(null);
  const [isPageLoading, setIsPageLoading] = useState<boolean>(false);

  // Dynamic Site Config
  const [siteConfig, setSiteConfig] = useState<any>(null);
  const [isConfigLoading, setIsConfigLoading] = useState<boolean>(true);

  // Dynamic Content Data
  const [posts, setPosts] = useState<Post[]>([]);
  const [selectedPost, setSelectedPost] = useState<Post | null>(null);
  const [isPostsLoading, setIsPostsLoading] = useState<boolean>(false);
  const [dbDownloads, setDbDownloads] = useState<DownloadItem[]>([]);
  const [galleries, setGalleries] = useState<GalleryItem[]>([]);
  const [activeLightboxImage, setActiveLightboxImage] = useState<GalleryItem | null>(null);
  const [newsFilterCategory, setNewsFilterCategory] = useState<string>('');
  const [newsFilterSearch, setNewsFilterSearch] = useState<string>('');
  const [sidebarSearch, setSidebarSearch] = useState<string>('');
  const [contactName, setContactName] = useState<string>('');
  const [contactEmail, setContactEmail] = useState<string>('');
  const [contactPhone, setContactPhone] = useState<string>('');
  const [contactSubject, setContactSubject] = useState<string>('');
  const [contactMessage, setContactMessage] = useState<string>('');
  const [contactIsSubmitting, setContactIsSubmitting] = useState<boolean>(false);
  const [contactSuccess, setContactSuccess] = useState<string | null>(null);
  const [contactError, setContactError] = useState<string | null>(null);
  const [publicPages, setPublicPages] = useState<PageData[]>([]);

  // Navigation menu state
  const [mobileMenuOpen, setMobileMenuOpen] = useState<boolean>(false);
  const [activeMenuDropdown, setActiveMenuDropdown] = useState<string | null>(null);

  // Scroll-to-top visibility
  const [showScrollTop, setShowScrollTop] = useState(false);
  useEffect(() => {
    const onScroll = () => setShowScrollTop(window.scrollY > 400);
    window.addEventListener('scroll', onScroll, { passive: true });
    return () => window.removeEventListener('scroll', onScroll);
  }, []);
  const handleScrollToTop = useCallback(() => {
    window.scrollTo({ top: 0, behavior: 'smooth' });
  }, []);

  // Slideshow Logic for Welcome section
// Slideshow Logic for Welcome section
const heroImages = useMemo(() => {
  const s = siteConfig?.settings || {};
  if (s.hero_images && Array.isArray(s.hero_images) && s.hero_images.length > 0) {
    return s.hero_images;
  }
  if (s.hero_image) {
    return [s.hero_image];
  }
  // No hero images configured
  return [];
}, [siteConfig]);

  const [heroIdx, setHeroIdx] = useState(0);

  useEffect(() => {
    if (heroImages.length <= 1) return;
    const timer = setInterval(() => {
      setHeroIdx((idx) => (idx + 1) % heroImages.length);
    }, 6000);
    return () => clearInterval(timer);
  }, [heroImages.length]);

  const nextSlide = () => {
    if (heroImages.length > 1) {
      setHeroIdx((idx) => (idx + 1) % heroImages.length);
    }
  };

  const prevSlide = () => {
    if (heroImages.length > 1) {
      setHeroIdx((idx) => (idx - 1 + heroImages.length) % heroImages.length);
    }
  };

  // DIP Tabs
  const [activeTab, setActiveTab] = useState<string>('berkala');
  const [searchTerm, setSearchTerm] = useState<string>('');

  // Form Permohonan Informasi State
  const [formApplicantType, setFormApplicantType] = useState<string>('perseorangan');
  const [formName, setFormName] = useState<string>('');
  const [formIdentityNumber, setFormIdentityNumber] = useState<string>('');
  const [formEmail, setFormEmail] = useState<string>('');
  const [formPhone, setFormPhone] = useState<string>('');
  const [formAddress, setFormAddress] = useState<string>('');
  const [formDetails, setFormDetails] = useState<string>('');
  const [formPurpose, setFormPurpose] = useState<string>('');
  const [formObtainMethod, setFormObtainMethod] = useState<string>('melihat-membaca');
  const [formDeliveryMethod, setFormDeliveryMethod] = useState<string>('mengambil-langsung');
  const [formAttachmentUrl, setFormAttachmentUrl] = useState<string>('');
  const [formIsUploading, setFormIsUploading] = useState<boolean>(false);
  const [formIsSubmitting, setFormIsSubmitting] = useState<boolean>(false);
  const [formSubmitSuccess, setFormSubmitSuccess] = useState<any>(null);
  const [formSubmitError, setFormSubmitError] = useState<string>('');

  // Track Permohonan Ticket State
  const [trackTicketNumber, setTrackTicketNumber] = useState<string>('');
  const [trackResult, setTrackResult] = useState<PermohonanTicket | null>(null);
  const [isTrackLoading, setIsTrackLoading] = useState<boolean>(false);
  const [trackError, setTrackError] = useState<string>('');
  const [permohonanSubTab, setPermohonanSubTab] = useState<string>('form'); // 'form', 'track', 'flow'

  // Accordion state for Regulasi KIP
  const [expandedRegulasi, setExpandedRegulasi] = useState<Record<string, boolean>>({
    A: true,
    B: false,
    C: false
  });

  const toggleRegulasi = (key: string) => {
    setExpandedRegulasi(prev => ({
      ...prev,
      [key]: !prev[key]
    }));
  };

  // About PPID section tab state
  const [aboutTab, setAboutTab] = useState<string>('sejarah');

  // State to track about image width for dynamic SVG masking
  const [aboutWidth, setAboutWidth] = useState<number>(540);
  const aboutRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (currentPage !== 'home' || !aboutRef.current) return;

    const handleResize = () => {
      if (aboutRef.current) {
        setAboutWidth(aboutRef.current.offsetWidth);
      }
    };

    // Initial measurement
    handleResize();

    // Setup resize listener
    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, [currentPage]);

  // Admin Portal state
  const [adminUser, setAdminUser] = useState<UserProfile | null>(null);
  const [adminEmail, setAdminEmail] = useState<string>('');
  const [adminPassword, setAdminPassword] = useState<string>('');
  const [showPassword, setShowPassword] = useState<boolean>(false);
  const [adminLoginLoading, setAdminLoginLoading] = useState<boolean>(false);
  const [adminLoginError, setAdminLoginError] = useState<string>('');
  const [adminActiveTab, setAdminActiveTab] = useState<string>('dashboard'); // 'dashboard', 'downloads', 'gallery', 'pages', 'posts', 'messages', 'permohonan', 'settings'
  const [adminGlobalMessage, setAdminGlobalMessage] = useState<string>('');

  // Admin Lists
  const [adminPages, setAdminPages] = useState<PageData[]>([]);
  const [adminPosts, setAdminPosts] = useState<Post[]>([]);
  const [adminDownloads, setAdminDownloads] = useState<DownloadItem[]>([]);
  const [adminGalleries, setAdminGalleries] = useState<GalleryItem[]>([]);
  const [adminPermohonans, setAdminPermohonans] = useState<PermohonanTicket[]>([]);
  const [adminStats, setAdminStats] = useState<any>({ total_pages: 0, total_posts: 0, total_downloads: 0, total_permohonan: 0, pending_permohonan: 0 });

  // Admin Selected Items for Edit Modal
  const [activeEditItem, setActiveEditItem] = useState<any>(null);
  const [editModalOpen, setEditModalOpen] = useState<boolean>(false);
  const [editModalType, setEditModalType] = useState<string>(''); // 'page', 'post', 'download', 'gallery', 'permohonan_review'

  // Form values for admin CRUD
  const [adminEditTitle, setAdminEditTitle] = useState<string>('');
  const [adminEditSubtitle, setAdminEditSubtitle] = useState<string>('');
  const [adminEditContent, setAdminEditContent] = useState<string>('');
  const [adminEditSlug, setAdminEditSlug] = useState<string>('');
  const [adminEditCoverImage, setAdminEditCoverImage] = useState<string>('');
  const [adminEditCategory, setAdminEditCategory] = useState<string>('');
  const [adminEditFileUrl, setAdminEditFileUrl] = useState<string>('');
  const [adminEditDescription, setAdminEditDescription] = useState<string>('');
  const [adminEditMediaType, setAdminEditMediaType] = useState<string>('image');
  const [adminEditIsUploading, setAdminEditIsUploading] = useState<boolean>(false);

  // Admin configuration form fields (Branding logo settings)
  const [portalName, setPortalName] = useState<string>('');
  const [portalDesc, setPortalDesc] = useState<string>('');
  const [portalLogo, setPortalLogo] = useState<string>('');
  const [portalFavicon, setPortalFavicon] = useState<string>('');
  const [portalWelcomeText, setPortalWelcomeText] = useState<string>('');
  const [portalSkRektor, setPortalSkRektor] = useState<string>('');
  const [portalPermohonanLink, setPortalPermohonanLink] = useState<string>('');
  const [portalPermohonanFormType, setPortalPermohonanFormType] = useState<string>('internal');
  const [portalKeberatanLink, setPortalKeberatanLink] = useState<string>('');
  const [portalAboutStatNumber, setPortalAboutStatNumber] = useState<string>('2021');
  const [portalAboutStatLabelAccent, setPortalAboutStatLabelAccent] = useState<string>('Tahun Berdiri');
  const [portalAboutStatLabel, setPortalAboutStatLabel] = useState<string>('PPID UPERTIS Melayani');
  // Service cards (Selamat Datang section)
  const [portalCard1Title, setPortalCard1Title] = useState<string>('Permohonan Informasi Publik');
  const [portalCard1Desc, setPortalCard1Desc] = useState<string>('Ajukan permohonan informasi publik kepada PPID UPERTIS secara daring maupun luring dengan mudah dan cepat sesuai prosedur UU KIP.');
  const [portalCard2Title, setPortalCard2Title] = useState<string>('Keberatan Informasi Publik');
  const [portalCard2Desc, setPortalCard2Desc] = useState<string>('Ajukan keberatan resmi jika permohonan informasi Anda tidak ditanggapi, ditolak, atau tidak sesuai dengan ketentuan yang berlaku.');
  const [portalCard3Title, setPortalCard3Title] = useState<string>('Pengaduan Layanan');
  const [portalCard3Desc, setPortalCard3Desc] = useState<string>('Sampaikan pengaduan terkait layanan informasi publik PPID UPERTIS melalui kanal resmi kami untuk penanganan yang transparan dan akuntabel.');
  const [portalCard1Link, setPortalCard1Link] = useState<string>('permohonan-informasi');
  const [portalCard2Link, setPortalCard2Link] = useState<string>('keberatan-informasi');
  const [portalCard3Link, setPortalCard3Link] = useState<string>('kontak');
  const [portalFaqs, setPortalFaqs] = useState<{ question: string; answer: string }[]>([]);
  const [portalPengaduanLink, setPortalPengaduanLink] = useState<string>('');

  const [portalVisi, setPortalVisi] = useState<string>('');
  const [portalMisi, setPortalMisi] = useState<string>('');
  const [portalBgImage, setPortalBgImage] = useState<string>('');
  const [portalFontFamily, setPortalFontFamily] = useState<string>('DM Sans');
  const [portalFontSize, setPortalFontSize] = useState<string>('normal');
  const [portalRektoratEmail, setPortalRektoratEmail] = useState<string>('');
  const [portalRektoratPhone, setPortalRektoratPhone] = useState<string>('');
  const [portalRektoratAddress, setPortalRektoratAddress] = useState<string>('');
  const [portalKampus2Address, setPortalKampus2Address] = useState<string>('');
  const [portalKampus1MapUrl, setPortalKampus1MapUrl] = useState<string>('');
  const portalJadwalSabtuMingguDefault = 'Sistem Online Tetap Aktif 24/7 (Desk fisik Tutup)';
  const [portalKampus2MapUrl, setPortalKampus2MapUrl] = useState<string>('');
  const [portalJadwalSeninKamis, setPortalJadwalSeninKamis] = useState<string>('08:00 – 16:00 WIB');
  const [portalIstirahatSeninKamis, setPortalIstirahatSeninKamis] = useState<string>('12:00 – 13:30 WIB');
  const [portalJadwalJumat, setPortalJadwalJumat] = useState<string>('08:00 – 16:30 WIB');
  const [portalIstirahatJumat, setPortalIstirahatJumat] = useState<string>('12:00 – 14:00 WIB');
  const [portalJadwalSabtuMinggu, setPortalJadwalSabtuMinggu] = useState<string>(portalJadwalSabtuMingguDefault);
  const [portalIsSaving, setPortalIsSaving] = useState<boolean>(false);

  // Edit Page schedule fields states (for jadwal-layanan-informasi slug)
  const [editJadwalSeninKamisKerja, setEditJadwalSeninKamisKerja] = useState<string>('08:00 – 16:00 WIB');
  const [editJadwalSeninKamisIstirahat, setEditJadwalSeninKamisIstirahat] = useState<string>('12:00 – 13:30 WIB');
  const [editJadwalJumatKerja, setEditJadwalJumatKerja] = useState<string>('08:00 – 16:30 WIB');
  const [editJadwalPageIstirahatJumat, setEditJadwalPageIstirahatJumat] = useState<string>('12:00 – 14:00 WIB');
  const [editJadwalSabtuMinggu, setEditJadwalSabtuMinggu] = useState<string>('Sistem Online Tetap Aktif 24/7 (Desk fisik Tutup)');
  const [editJadwalOffline1, setEditJadwalOffline1] = useState<string>('Membawa kartu identitas resmi (KTP untuk perorangan, Akta/SK pendirian untuk instansi).');
  const [editJadwalOffline2, setEditJadwalOffline2] = useState<string>('Loket Pelayanan berada di Lantai 1 Gedung Rektorat Kampus Utama UPERTIS Padang.');
  const [editJadwalOffline3, setEditJadwalOffline3] = useState<string>('Petugas kami siap membantu pencatatan, peninjauan berkas, hingga cetak salinan dokumen.');
  const [editJadwalOnline1, setEditJadwalOnline1] = useState<string>('Gunakan menu Permohonan Informasi untuk mengisi formulir digital dan upload KTP.');
  const [editJadwalOnline2, setEditJadwalOnline2] = useState<string>('Simpan nomor tiket pengajuan untuk melacak tanggapan admin secara real-time.');
  const [editJadwalOnline3, setEditJadwalOnline3] = useState<string>('Jawaban atau link unduhan berkas digital akan dikirimkan langsung ke email/tiket Anda.');
  const [editJadwalCustomRemarks, setEditJadwalCustomRemarks] = useState<string>('');
  const [adminEditPageDocs, setAdminEditPageDocs] = useState<Array<{ title: string; description: string; file_url: string }>>([]);
  const [editKeberatanManualSteps, setEditKeberatanManualSteps] = useState<any[]>(defaultManualSteps);
  const [editKeberatanOnlineSteps, setEditKeberatanOnlineSteps] = useState<any[]>(defaultOnlineSteps);
  const [editSengketaCards, setEditSengketaCards] = useState<Array<{ title: string; desc: string }>>([]);
  const [editSengketaFlow, setEditSengketaFlow] = useState<{
    jalurA_title: string; jalurA_desc: string;
    jalurB_title: string; jalurB_desc: string;
    verifikasi_title: string; verifikasi_desc: string;
    help_title: string; help_desc: string;
    download_title: string; download_desc: string;
  }>({
    jalurA_title: '', jalurA_desc: '',
    jalurB_title: '', jalurB_desc: '',
    verifikasi_title: '', verifikasi_desc: '',
    help_title: '', help_desc: '',
    download_title: '', download_desc: ''
  });
  const [editProfilStats, setEditProfilStats] = useState<Array<{ value: string; label: string }>>([]);
  const [editProfilImage1, setEditProfilImage1] = useState<string>('');
  const [editProfilImage2, setEditProfilImage2] = useState<string>('');
  const [editProfilImage3, setEditProfilImage3] = useState<string>('');
  const [editTugasImage1, setEditTugasImage1] = useState<string>('');
  const [editTugasImage2, setEditTugasImage2] = useState<string>('');
  const [editTugasImage3, setEditTugasImage3] = useState<string>('');
  const [editTugasList, setEditTugasList] = useState<Array<{ title: string; desc: string }>>([]);
  const [editFungsiList, setEditFungsiList] = useState<Array<{ title: string; items: string[] }>>([]);

  // Permohonan Action State
  const [permohonanActionStatus, setPermohonanActionStatus] = useState<string>('pending');
  const [permohonanActionResponse, setPermohonanActionResponse] = useState<string>('');
  const [permohonanActionLoading, setPermohonanActionLoading] = useState<boolean>(false);

  // Auto-clear adminGlobalMessage after 3 seconds
  useEffect(() => {
    if (adminGlobalMessage) {
      const timer = setTimeout(() => {
        setAdminGlobalMessage('');
      }, 3000);
      return () => clearTimeout(timer);
    }
  }, [adminGlobalMessage]);

  // Fetch Site Config (Menu, Accent, Logo, Welcome Settings)
  useEffect(() => {
    setIsConfigLoading(true);
    const host = window.location.hostname;
    fetch(`${API_BASE_URL}/site-config?host=${host}`)
      .then((res) => res.json())
      .then((data) => {
        if (data && !data.is_main) {
          setSiteConfig(data);
          loadConfigFields(data);
        } else {
          // Fallback to fetch config via ppid.localhost slug explicitly
          fetch(`${API_BASE_URL}/site-config?host=ppid.localhost`)
            .then((r) => r.json())
            .then((d) => {
              if (d && !d.is_main) {
                setSiteConfig(d);
                loadConfigFields(d);
              }
            });
        }
      })
      .catch((err) => console.error('Error fetching site config:', err))
      .finally(() => setIsConfigLoading(false));
  }, []);

  // Fetch public page list
  useEffect(() => {
    fetch(`${API_BASE_URL}/pages`)
      .then(res => res.json())
      .then(data => {
        if (Array.isArray(data)) {
          setPublicPages(data);
        }
      })
      .catch(err => console.error(err));
  }, []);

  // Fetch public news with filters
  useEffect(() => {
    setIsPostsLoading(true);
    let url = `${API_BASE_URL}/posts?limit=12`;
    if (newsFilterSearch) {
      url += `&search=${encodeURIComponent(newsFilterSearch)}`;
    }
    fetch(url)
      .then((res) => res.json())
      .then((data) => {
        if (Array.isArray(data)) {
          setPosts(data);
        } else if (data && Array.isArray(data.posts)) {
          setPosts(data.posts);
        }
      })
      .catch((err) => console.error('Error fetching posts:', err))
      .finally(() => setIsPostsLoading(false));
  }, [newsFilterSearch]);

  // Fetch public downloads
  useEffect(() => {
    fetch(`${API_BASE_URL}/downloads`)
      .then((res) => res.json())
      .then((data) => {
        if (Array.isArray(data)) {
          setDbDownloads(data);
        }
      })
      .catch((err) => console.error('Error fetching downloads:', err));
  }, []);

  // Fetch public galleries
  useEffect(() => {
    fetch(`${API_BASE_URL}/gallery`)
      .then((res) => res.json())
      .then((data) => {
        if (Array.isArray(data)) {
          setGalleries(data);
        }
      })
      .catch((err) => console.error('Error fetching galleries:', err));
  }, []);

  // Check login on mount
  useEffect(() => {
    const token = localStorage.getItem('auth_token');
    if (token) {
      fetch(`${API_BASE_URL}/auth/me`, {
        headers: { Authorization: `Bearer ${token}` }
      })
        .then(res => {
          if (res.ok) return res.json();
          throw new Error('Invalid token');
        })
        .then(user => {
          if (user && (user.role === 'admin' || user.role === 'faculty_admin' || user.role === 'ppid_admin')) {
            setAdminUser(user);
            fetchAdminData();
          } else {
            localStorage.removeItem('auth_token');
          }
        })
        .catch(() => localStorage.removeItem('auth_token'));
    }
  }, []);

  // SPA HTML5 Pathname Router for SEO-Friendly Clean URLs
  useEffect(() => {
    const handleLocationChange = () => {
      const path = window.location.pathname;

      if (path === '/' || path === '') {
        setCurrentPage('home');
        setActiveSlug('');
        setPageData(null);
      } else if (path === '/berita' || path === '/berita/') {
        setCurrentPage('berita');
        setActiveSlug('');
        setPageData(null);
        setNewsFilterSearch('');
        setNewsFilterCategory('');
        setSidebarSearch('');
      } else if (path.startsWith('/berita/')) {
        const slug = path.replace('/berita/', '');
        navigateToNewsDetail(slug);
      } else if (path.startsWith('/halaman/')) {
        const slug = path.replace('/halaman/', '');
        navigateToPage(slug);
      } else if (path === '/admin' || path === '/admin/') {
        setCurrentPage('admin');
        setActiveSlug('');
        setPageData(null);
      }
    };

    // Run on initial load
    handleLocationChange();

    window.addEventListener('popstate', handleLocationChange);
    return () => window.removeEventListener('popstate', handleLocationChange);
  }, []);

  // Sync state values when config loads
  const loadConfigFields = (config: any) => {
    setPortalName(config.name || '');
    setPortalDesc(config.description || '');
    setPortalLogo(config.logo_url || '');
    setPortalFavicon(config.favicon_url || '');
    if (config.settings) {
      setPortalWelcomeText(config.settings.welcome_text || '');
      setPortalSkRektor(config.settings.sk_rektor || '');
      setPortalPermohonanLink(config.settings.permohonan_link || '');
      setPortalPermohonanFormType(config.settings.permohonan_form_type || 'internal');
      setPortalKeberatanLink(config.settings.keberatan_link || '');
      setPortalPengaduanLink(config.settings.pengaduan_link || '');

      setPortalVisi(config.settings.ppid_visi || '');
      setPortalMisi(config.settings.ppid_misi || '');
      setPortalBgImage(config.settings.background_image || '');
      setPortalFontFamily(config.settings.font_family || 'DM Sans');
      setPortalFontSize(config.settings.font_size || 'normal');
      setPortalRektoratEmail(config.settings.rektorat_email || '');
      setPortalRektoratPhone(config.settings.rektorat_phone || '');
      setPortalRektoratAddress(config.settings.rektorat_address || '');
      setPortalKampus2Address(config.settings.kampus2_address || '');
      setPortalKampus1MapUrl(config.settings.kampus1_map_url || '');
      setPortalKampus2MapUrl(config.settings.kampus2_map_url || '');
      setPortalJadwalSeninKamis(config.settings.jadwal_senin_kamis || '08:00 – 16:00 WIB');
      setPortalIstirahatSeninKamis(config.settings.istirahat_senin_kamis || '12:00 – 13:30 WIB');
      setPortalJadwalJumat(config.settings.jadwal_jumat || '08:00 – 16:30 WIB');
      setPortalIstirahatJumat(config.settings.istirahat_jumat || '12:00 – 14:00 WIB');
      setPortalJadwalSabtuMinggu(config.settings.jadwal_sabtu_minggu || portalJadwalSabtuMingguDefault);
      setPortalAboutStatNumber(config.settings.about_stat_number || '2021');
      setPortalAboutStatLabelAccent(config.settings.about_stat_label_accent || 'Tahun Berdiri');
      setPortalAboutStatLabel(config.settings.about_stat_label || 'PPID UPERTIS Melayani');
      setPortalCard1Title(config.settings.service_card_1_title || 'Permohonan Informasi Publik');
      setPortalCard1Desc(config.settings.service_card_1_desc || 'Ajukan permohonan informasi publik kepada PPID UPERTIS secara daring maupun luring dengan mudah dan cepat sesuai prosedur UU KIP.');
      setPortalCard2Title(config.settings.service_card_2_title || 'Keberatan Informasi Publik');
      setPortalCard2Desc(config.settings.service_card_2_desc || 'Ajukan keberatan resmi jika permohonan informasi Anda tidak ditanggapi, ditolak, atau tidak sesuai dengan ketentuan yang berlaku.');
      setPortalCard3Title(config.settings.service_card_3_title || 'Pengaduan Layanan');
      setPortalCard3Desc(config.settings.service_card_3_desc || 'Sampaikan pengaduan terkait layanan informasi publik PPID UPERTIS melalui kanal resmi kami untuk penanganan yang transparan dan akuntabel.');
      setPortalCard1Link(config.settings.service_card_1_link || 'permohonan-informasi');
      setPortalCard2Link(config.settings.service_card_2_link || 'keberatan-informasi');
      setPortalCard3Link(config.settings.service_card_3_link || 'kontak');
      setPortalFaqs(config.settings.faqs || [
        {
          question: "Bagaimana cara mengajukan permohonan informasi?",
          answer: "Anda dapat mengajukan secara online melalui menu 'Permohonan Informasi' di portal ini, atau datang langsung ke Desk Layanan PPID UPERTIS di Gedung Rektorat Lantai 1. Siapkan identitas diri seperti KTP (perorangan) atau Akta Pendirian (organisasi)."
        },
        {
          question: "Berapa lama waktu proses permohonan informasi?",
          answer: "Sesuai dengan UU KIP No. 14 Tahun 2008, PPID akan memberikan jawaban/tanggapan dalam waktu 10 hari kerja sejak permohonan terdaftar, dan dapat diperpanjang paling lambat 7 hari kerja berikutnya."
        },
        {
          question: "Bagaimana jika permohonan informasi saya ditolak?",
          answer: "Jika permohonan ditolak atau tanggapan kurang memuaskan, Anda dapat mengajukan Keberatan Informasi secara online melalui menu 'Keberatan Informasi' di portal ini dalam waktu maksimal 30 hari kerja."
        },
        {
          question: "Apakah layanan informasi ini dipungut biaya?",
          answer: "Layanan permohonan informasi publik di PPID UPERTIS sepenuhnya GRATIS. Jika ada biaya penggandaan berkas fisik atau pengiriman dokumen, biaya tersebut ditanggung oleh pemohon."
        },
        {
          question: "Kapan jadwal operasional Desk Layanan PPID?",
          answer: "Desk fisik PPID buka Senin s/d Kamis (08:00 - 16:00 WIB) dan Jumat (08:00 - 16:30 WIB). Untuk layanan online melalui portal ini aktif 24/7."
        }
      ]);
    }
  };

  // Fetch admin dashboard details
  const fetchAdminData = () => {
    const token = localStorage.getItem('auth_token');
    if (!token) return;

    // FetchStats
    fetch(`${API_BASE_URL}/admin/dashboard`, { headers: { Authorization: `Bearer ${token}` } })
      .then(res => res.json())
      .then(data => setAdminStats(data))
      .catch(err => console.error(err));

    // Fetch lists
    fetch(`${API_BASE_URL}/admin/pages`, { headers: { Authorization: `Bearer ${token}` } })
      .then(res => res.json())
      .then(data => {
        if (Array.isArray(data)) {
          const ppidPageSlugs = ['profil', 'visi-misi', 'maklumat', 'struktur-organisasi-2', 'tugas-dan-fungsi', 'regulasi', 'informasi-publik-berkala', 'informasi-tersedia-setiap-saat', 'info-serta-merta', 'zona-integrasi', 'jadwal-layanan-informasi', 'informasi-dikecualikan', 'keberatan-informasi', 'Permohonan-penyelesaian-sengketa', 'permohonan-penyelesaian-sengketa'];
          setAdminPages(data.filter(p => ppidPageSlugs.includes(p.slug)));
        }
      })
      .catch(err => console.error(err));

    fetch(`${API_BASE_URL}/admin/posts`, { headers: { Authorization: `Bearer ${token}` } })
      .then(res => res.json())
      .then(data => {
        if (Array.isArray(data)) {
          setAdminPosts(data);
        }
      })
      .catch(err => console.error(err));

    fetch(`${API_BASE_URL}/admin/downloads`, { headers: { Authorization: `Bearer ${token}` } })
      .then(res => res.json())
      .then(data => {
        if (Array.isArray(data)) {
          const ppidCategories = ['ppid-berkala', 'ppid-setiap-saat', 'ppid-serta-merta', 'ppid-dikecualikan'];
          setAdminDownloads(data.filter(d => ppidCategories.includes(d.category)));
        }
      })
      .catch(err => console.error(err));

    fetch(`${API_BASE_URL}/admin/permohonan_informasi`, { headers: { Authorization: `Bearer ${token}` } })
      .then(res => res.json())
      .then(data => Array.isArray(data) && setAdminPermohonans(data))
      .catch(err => console.error(err));

    fetch(`${API_BASE_URL}/admin/gallery`, { headers: { Authorization: `Bearer ${token}` } })
      .then(res => res.json())
      .then(data => Array.isArray(data) && setAdminGalleries(data))
      .catch(err => console.error(err));
  };

  // Increment download count API helper
  const incrementDownloadCount = async (id: string, url: string) => {
    setDbDownloads(prev => prev.map(item => item.id === id ? { ...item, downloads_count: item.downloads_count + 1 } : item));
    try {
      await fetch(`${API_BASE_URL}/downloads/${id}/increment`, { method: 'POST' });
    } catch (e) {
      console.error(e);
    }
    window.open(url, '_blank');
  };

  // Menu lists definitions
  const defaultMenuGroups = [
    { label: 'Home', href: '/', type: 'link' },
    {
      label: 'Tentang',
      type: 'dropdown',
      items: [
        { label: 'Profil PPID', href: 'profil' },
        { label: 'Visi & Misi', href: 'visi-misi' },
        { label: 'Maklumat Pelayanan', href: 'maklumat' },
        { label: 'Struktur Organisasi', href: 'struktur-organisasi-2' },
        { label: 'Tugas & Fungsi PPID', href: 'tugas-dan-fungsi' },
        { label: 'Regulasi KIP', href: 'regulasi' }
      ]
    },
    {
      label: 'Info Publik',
      type: 'dropdown',
      items: [
        { label: 'Info Publik Berkala', href: 'informasi-publik-berkala' },
        { label: 'Info Tersedia Setiap Saat', href: 'informasi-tersedia-setiap-saat' },
        { label: 'Info Serta Merta', href: 'info-serta-merta' },
        { label: 'Zona Integrasi', href: 'zona-integrasi' }
      ]
    },
    {
      label: 'Layanan',
      type: 'dropdown',
      items: [
        { label: 'Jadwal Layanan', href: 'jadwal-layanan-informasi' },
        {
          label: 'Permohonan Informasi',
          href: siteConfig?.settings?.permohonan_form_type === 'external'
            ? (siteConfig?.settings?.permohonan_link || 'https://forms.gle/b2N4iBRcCFwKEg61A')
            : 'permohonan-informasi'
        },
        { label: 'Pengajuan Keberatan', href: 'keberatan-informasi' },
        { label: 'Pengaduan Layanan', href: siteConfig?.settings?.pengaduan_link || 'https://lapor.go.id/', isExternal: true },
        { label: 'Informasi Dikecualikan', href: 'informasi-dikecualikan' }
      ]
    },
    { label: 'Berita PPID', href: 'berita', type: 'link' },
    { label: 'Galeri', href: 'galeri', type: 'link' },
    { label: 'Download', href: 'download', type: 'link' },
    { label: 'Hubungi Kami', href: 'kontak', type: 'link', isHighlight: true }
  ];

  const menuGroups = siteConfig?.menu && Array.isArray(siteConfig.menu) && siteConfig.menu.length > 0
    ? migrateMenuStructure(siteConfig.menu)
    : defaultMenuGroups;

  // DIP Categories items resolving
  const getDIPData = () => {
    const sections = [
      {
        id: 'berkala',
        title: 'Informasi Berkala',
        desc: 'Informasi publik yang wajib diperbarui dan disediakan secara berkala sesuai UU KIP Pasal 9.',
        dbKey: 'ppid-berkala',
        items: [
          { label: 'Profil UPERTIS', href: 'profil', type: 'Page' },
          { label: 'Visi Misi UPERTIS', href: 'visi-misi', type: 'Page' },
          { label: 'Tugas Fungsi dan Ruang Lingkup Kegiatan', href: 'tugas-dan-fungsi', type: 'Page' },
          { label: 'Struktur Organisasi UPERTIS', href: 'struktur-organisasi-2', type: 'Page' }
        ]
      },
      {
        id: 'setiap-saat',
        title: 'Tersedia Setiap Saat',
        desc: 'Informasi publik yang wajib disediakan dan siap diberikan setiap saat dibutuhkan sesuai UU KIP Pasal 11.',
        dbKey: 'ppid-setiap-saat',
        items: [
          { label: 'Daftar Informasi Publik (DIP) UPERTIS', href: 'https://ppid.upertis.ac.id/download/', type: 'External' }
        ]
      },
      {
        id: 'serta-merta',
        title: 'Informasi Serta Merta',
        desc: 'Informasi publik yang wajib diumumkan segera demi keselamatan dan kesehatan masyarakat (UU KIP Pasal 10).',
        dbKey: 'ppid-serta-merta',
        items: [
          { label: 'Panduan Mitigasi & Evakuasi Bencana Alam di Lingkungan Kampus', href: '#', type: 'Info' }
        ]
      },
      {
        id: 'dikecualikan',
        title: 'Informasi Dikecualikan',
        desc: 'Informasi publik yang bersifat rahasia dan dibatasi ketat sesuai UU KIP Pasal 17.',
        dbKey: 'ppid-dikecualikan',
        items: [
          { label: 'SK Daftar Informasi Dikecualikan UPERTIS', href: 'https://ppid.upertis.ac.id/download/daftar-informasi-dikecualikan-informasi-setiap-saat/', type: 'External' }
        ]
      }
    ];

    return sections.map(section => {
      const dbItems = dbDownloads
        .filter(item => item.category === section.dbKey || item.category === section.id)
        .map(item => ({
          label: item.title,
          href: resolveImageUrl(item.file_url),
          type: 'DownloadItem',
          downloadsCount: item.downloads_count,
          id: item.id
        }));
      return {
        ...section,
        items: [...section.items, ...dbItems]
      };
    });
  };

  const dipData = getDIPData();
  const currentDipSection = dipData.find(sec => sec.id === activeTab);
  const filteredDipItems = currentDipSection
    ? currentDipSection.items.filter(item => item.label.toLowerCase().includes(searchTerm.toLowerCase()))
    : [];

  // Routing navigation helpers
  const handleNavigation = (href: string) => {
    setMobileMenuOpen(false);
    if (href === '/' || href === '') {
      navigateToHome();
    } else if (href === '/posts' || href === '/berita') {
      navigateToNews();
    } else if (href.startsWith('/halaman/')) {
      navigateToPage(href.replace('/halaman/', ''));
    } else if (href.startsWith('/berita/')) {
      navigateToNewsDetail(href.replace('/berita/', ''));
    } else if (href.startsWith('http://') || href.startsWith('https://')) {
      window.open(href, '_blank', 'noopener,noreferrer');
    } else {
      navigateToPage(href.replace('/', ''));
    }
  };

  const navigateToPage = (slug: string) => {
    setIsPageLoading(true);
    setCurrentPage(`page/${slug}`);
    setActiveSlug(slug);
    setSearchTerm('');
    window.scrollTo({ top: 0, behavior: 'smooth' });

    if (slug === 'informasi-dikecualikan') {
      setActiveTab('dikecualikan');
    } else if (slug === 'informasi-publik-berkala') {
      setActiveTab('berkala');
    } else if (slug === 'informasi-tersedia-setiap-saat') {
      setActiveTab('setiap-saat');
    } else if (slug === 'info-serta-merta') {
      setActiveTab('serta-merta');
    }

    if (window.location.pathname !== `/halaman/${slug}`) {
      window.history.pushState(null, '', `/halaman/${slug}`);
    }

    fetch(`${API_BASE_URL}/pages/${slug}`)
      .then(res => {
        if (!res.ok) throw new Error('Halaman tidak ditemukan');
        return res.json();
      })
      .then(data => {
        setPageData(data);
      })
      .catch((err) => {
        console.error(err);
        setPageData({
          id: 'notfound',
          title: slug.split('-').map(w => w.charAt(0).toUpperCase() + w.slice(1)).join(' '),
          content: `<p class="text-slate-500">Konten halaman untuk <strong>${slug}</strong> sedang dipersiapkan oleh administrator portal.</p>`,
          slug: slug
        });
      })
      .finally(() => setIsPageLoading(false));
  };

  const navigateToHome = () => {
    setCurrentPage('home');
    setActiveSlug('');
    setPageData(null);
    window.scrollTo({ top: 0, behavior: 'smooth' });
    if (window.location.pathname !== '/') {
      window.history.pushState(null, '', '/');
    }
  };

  const navigateToNews = () => {
    setCurrentPage('berita');
    setActiveSlug('');
    setPageData(null);
    setNewsFilterSearch('');
    setNewsFilterCategory('');
    setSidebarSearch('');
    window.scrollTo({ top: 0, behavior: 'smooth' });
    if (window.location.pathname !== '/berita') {
      window.history.pushState(null, '', '/berita');
    }
  };

  const navigateToNewsDetail = (slugOrId: string) => {
    setIsPageLoading(true);
    setCurrentPage('berita-detail');
    window.scrollTo({ top: 0, behavior: 'smooth' });
    if (window.location.pathname !== `/berita/${slugOrId}`) {
      window.history.pushState(null, '', `/berita/${slugOrId}`);
    }
    fetch(`${API_BASE_URL}/posts/${slugOrId}`)
      .then(res => res.json())
      .then(data => setSelectedPost(data))
      .catch(err => console.error(err))
      .finally(() => setIsPageLoading(false));
  };

  // Submit online permohonan form API call
  const handlePermohonanSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setFormSubmitError('');
    setFormSubmitSuccess(null);
    setFormIsSubmitting(true);

    const payload = {
      applicant_type: formApplicantType,
      name: formName,
      identity_number: formIdentityNumber,
      email: formEmail,
      phone: formPhone,
      address: formAddress,
      details: formDetails,
      purpose: formPurpose,
      obtain_method: formObtainMethod,
      delivery_method: formDeliveryMethod,
      attachment_url: formAttachmentUrl
    };

    try {
      const res = await fetch(`${API_BASE_URL}/permohonan-informasi`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
      });
      const data = await res.json();
      if (res.ok) {
        setFormSubmitSuccess(data);
        // Reset form
        setFormName('');
        setFormIdentityNumber('');
        setFormEmail('');
        setFormPhone('');
        setFormAddress('');
        setFormDetails('');
        setFormPurpose('');
        setFormAttachmentUrl('');
      } else {
        setFormSubmitError(data.error || 'Terjadi kesalahan saat mengirim permohonan.');
      }
    } catch (err) {
      console.error(err);
      setFormSubmitError('Gagal menghubungi server. Silakan coba kembali.');
    } finally {
      setFormIsSubmitting(false);
    }
  };

  // Submit contact message form call
  const handleContactSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setContactError(null);
    setContactSuccess(null);
    setContactIsSubmitting(true);

    try {
      const res = await fetch(`${API_BASE_URL}/contact`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          name: contactName,
          email: contactEmail,
          phone: contactPhone,
          subject: contactSubject,
          message: contactMessage
        })
      });
      const data = await res.json();
      if (res.ok) {
        setContactSuccess('Pesan Anda berhasil terkirim. Terima kasih telah menghubungi PPID UPERTIS!');
        setContactName('');
        setContactEmail('');
        setContactPhone('');
        setContactSubject('');
        setContactMessage('');
      } else {
        setContactError(data.error || 'Gagal mengirim pesan. Silakan coba lagi.');
      }
    } catch (err) {
      console.error('Contact submit error:', err);
      setContactError('Terjadi kesalahan koneksi server. Silakan coba beberapa saat lagi.');
    } finally {
      setContactIsSubmitting(false);
    }
  };

  // Upload attachment document file call
  const handleAttachmentUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    setFormIsUploading(true);
    setFormSubmitError('');

    const formData = new FormData();
    formData.append('file', file);

    try {
      const res = await fetch(`${API_BASE_URL}/permohonan-informasi/upload`, {
        method: 'POST',
        body: formData
      });
      const data = await res.json();
      if (res.ok && data.url) {
        setFormAttachmentUrl(data.url);
      } else {
        setFormSubmitError(data.error || 'Gagal mengunggah berkas.');
      }
    } catch (err) {
      console.error(err);
      setFormSubmitError('Koneksi terputus saat mengunggah.');
    } finally {
      setFormIsUploading(false);
    }
  };

  // Track Ticket Status call
  const handleTrackSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setTrackError('');
    setTrackResult(null);
    setIsTrackLoading(true);

    try {
      const res = await fetch(`${API_BASE_URL}/permohonan-informasi/status/${trackTicketNumber.trim()}`);
      const data = await res.json();
      if (res.ok) {
        setTrackResult(data);
      } else {
        setTrackError(data.error || 'Nomor tiket tidak ditemukan.');
      }
    } catch (err) {
      console.error(err);
      setTrackError('Gagal melakukan pelacakan tiket.');
    } finally {
      setIsTrackLoading(false);
    }
  };

  // Admin Login handler
  const handleAdminLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setAdminLoginError('');
    setAdminLoginLoading(true);

    try {
      const res = await fetch(`${API_BASE_URL}/auth/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email: adminEmail, password: adminPassword })
      });
      const data = await res.json();
      if (!res.ok) throw new Error(data.error || 'Email atau password salah');

      if (data.token) {
        localStorage.setItem('auth_token', data.token);
        const userRes = await fetch(`${API_BASE_URL}/auth/me`, {
          headers: { Authorization: `Bearer ${data.token}` }
        });
        if (userRes.ok) {
          const user = await userRes.json();
          if (user && (user.role === 'admin' || user.role === 'faculty_admin' || user.role === 'ppid_admin')) {
            setAdminUser(user);
            fetchAdminData();
          } else {
            localStorage.removeItem('auth_token');
            setAdminLoginError('Anda tidak memiliki hak akses admin PPID');
          }
        } else {
          localStorage.removeItem('auth_token');
          setAdminLoginError('Verifikasi profil akun gagal');
        }
      }
    } catch (err: any) {
      setAdminLoginError(err.message || 'Terjadi kesalahan sistem');
    } finally {
      setAdminLoginLoading(false);
    }
  };

  // Helper for admin file uploading
  const handleAdminUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    setAdminEditIsUploading(true);
    const token = localStorage.getItem('auth_token');
    if (!token) return;

    const formData = new FormData();
    formData.append('files', file);

    try {
      const res = await fetch(`${API_BASE_URL}/admin/uploads`, {
        method: 'POST',
        headers: { Authorization: `Bearer ${token}` },
        body: formData
      });
      const data = await res.json();
      const uploadedUrl = data.url || data.path || (data.urls && data.urls[0]) || (data.uploaded && data.uploaded[0] && data.uploaded[0].url) || '';
      if (res.ok && uploadedUrl) {
        setAdminEditFileUrl(uploadedUrl);
        if (editModalType === 'post' || editModalType === 'page') {
          setAdminEditCoverImage(uploadedUrl);
        }
        setAdminGlobalMessage('File berhasil diunggah!');
      } else {
        alert(data.error || 'Gagal mengunggah file.');
      }
    } catch (e) {
      console.error(e);
      alert('Error saat mengunggah berkas.');
    } finally {
      setAdminEditIsUploading(false);
    }
  };

  // Save admin settings form
  const handleSaveSettings = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!portalName.trim()) {
      alert('Nama portal wajib diisi');
      return;
    }
    const token = localStorage.getItem('auth_token');
    if (!token) return;
    setPortalIsSaving(true);

    const siteId = siteConfig?.id || 'ppid';
    const payload = {
      ...siteConfig,
      name: portalName,
      description: portalDesc,
      logo_url: portalLogo,
      favicon_url: portalFavicon,
      settings: {
        ...siteConfig?.settings,
        welcome_text: portalWelcomeText,
        sk_rektor: portalSkRektor,
        permohonan_link: portalPermohonanLink,
        permohonan_form_type: portalPermohonanFormType,
        keberatan_link: portalKeberatanLink,
        pengaduan_link: portalPengaduanLink,

        ppid_visi: portalVisi,
        ppid_misi: portalMisi,
        background_image: portalBgImage,
        font_family: portalFontFamily,
        font_size: portalFontSize,
        rektorat_email: portalRektoratEmail,
        rektorat_phone: portalRektoratPhone,
        rektorat_address: portalRektoratAddress,
        kampus2_address: portalKampus2Address,
        kampus1_map_url: portalKampus1MapUrl,
        kampus2_map_url: portalKampus2MapUrl,
        jadwal_senin_kamis: portalJadwalSeninKamis,
        istirahat_senin_kamis: portalIstirahatSeninKamis,
        jadwal_jumat: portalJadwalJumat,
        istirahat_jumat: portalIstirahatJumat,
        jadwal_sabtu_minggu: portalJadwalSabtuMinggu,
        about_stat_number: portalAboutStatNumber,
        about_stat_label_accent: portalAboutStatLabelAccent,
        about_stat_label: portalAboutStatLabel,
        service_card_1_title: portalCard1Title,
        service_card_1_desc: portalCard1Desc,
        service_card_2_title: portalCard2Title,
        service_card_2_desc: portalCard2Desc,
        service_card_3_title: portalCard3Title,
        service_card_3_desc: portalCard3Desc,
        service_card_1_link: portalCard1Link,
        service_card_2_link: portalCard2Link,
        service_card_3_link: portalCard3Link,
        faqs: portalFaqs
      }
    };

    try {
      const res = await fetch(`${API_BASE_URL}/admin/sites/${siteId}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`
        },
        body: JSON.stringify(payload)
      });
      if (res.ok) {
        setAdminGlobalMessage('Pengaturan portal berhasil disimpan!');
        // Refresh site config
        const host = window.location.hostname;
        const refreshRes = await fetch(`${API_BASE_URL}/site-config?host=${host}`);
        if (refreshRes.ok) {
          const freshData = await refreshRes.json();
          if (freshData && !freshData.is_main) {
            setSiteConfig(freshData);
            loadConfigFields(freshData);
          } else {
            const fallbackRes = await fetch(`${API_BASE_URL}/site-config?host=ppid.localhost`);
            if (fallbackRes.ok) {
              const fallbackData = await fallbackRes.json();
              if (fallbackData && !fallbackData.is_main) {
                setSiteConfig(fallbackData);
                loadConfigFields(fallbackData);
              }
            }
          }
        }
      } else {
        alert(`Gagal menyimpan pengaturan: ${(await res.json()).error}`);
      }
    } catch (err) {
      console.error(err);
      alert('Terjadi kesalahan jaringan.');
    } finally {
      setPortalIsSaving(false);
    }
  };

  // Admin CRUD logic trigger (Add/Edit)
  const openCreateModal = (type: string) => {
    setActiveEditItem(null);
    setEditModalType(type);
    setAdminEditTitle('');
    setAdminEditSubtitle('');
    setAdminEditContent('');
    setAdminEditSlug('');
    setAdminEditPageDocs([]);
    setAdminEditCoverImage('');
    setAdminEditCategory(type === 'download' ? 'ppid-berkala' : '');
    setAdminEditFileUrl('');
    setAdminEditDescription('');
    setAdminEditMediaType('image');

    // Reset schedule states
    setEditJadwalSeninKamisKerja('08:00 – 16:00 WIB');
    setEditJadwalSeninKamisIstirahat('12:00 – 13:30 WIB');
    setEditJadwalJumatKerja('08:00 – 16:30 WIB');
    setEditJadwalPageIstirahatJumat('12:00 – 14:00 WIB');
    setEditJadwalSabtuMinggu('Sistem Online Tetap Aktif 24/7 (Desk fisik Tutup)');
    setEditJadwalOffline1('Membawa kartu identitas resmi (KTP untuk perorangan, Akta/SK pendirian untuk instansi).');
    setEditJadwalOffline2('Loket Pelayanan berada di Lantai 1 Gedung Rektorat Kampus Utama UPERTIS Padang.');
    setEditJadwalOffline3('Petugas kami siap membantu pencatatan, peninjauan berkas, hingga cetak salinan dokumen.');
    setEditJadwalOnline1('Gunakan menu Permohonan Informasi untuk mengisi formulir digital dan upload KTP.');
    setEditJadwalOnline2('Simpan nomor tiket pengajuan untuk melacak tanggapan admin secara real-time.');
    setEditJadwalOnline3('Jawaban atau link unduhan berkas digital akan dikirimkan langsung ke email/tiket Anda.');
    setEditJadwalCustomRemarks('');

    setEditModalOpen(true);
  };

  const openEditModal = (type: string, item: any) => {
    setActiveEditItem(item);
    setEditModalType(type);
    setAdminEditTitle(item.title || '');
    setAdminEditSubtitle(item.subtitle || '');
    setAdminEditContent(item.content || '');
    setAdminEditSlug(item.slug || '');
    setAdminEditCoverImage(item.cover_image_url || '');
    setAdminEditCategory(item.category || '');
    setAdminEditFileUrl(item.file_url || item.media_url || '');
    setAdminEditDescription(item.description || '');
    setAdminEditMediaType(item.media_type || 'image');

    // Parse JSON if editing Jadwal Layanan page
    if (type === 'page' && item.slug === 'jadwal-layanan-informasi') {
      try {
        const data = JSON.parse(item.content);
        setEditJadwalSeninKamisKerja(data.senin_kamis_kerja || '08:00 – 16:00 WIB');
        setEditJadwalSeninKamisIstirahat(data.senin_kamis_istirahat || '12:00 – 13:30 WIB');
        setEditJadwalJumatKerja(data.jumat_kerja || '08:00 – 16:30 WIB');
        setEditJadwalPageIstirahatJumat(data.jumat_istirahat || '12:00 – 14:00 WIB');
        setEditJadwalSabtuMinggu(data.sabtu_minggu || 'Sistem Online Tetap Aktif 24/7 (Desk fisik Tutup)');
        setEditJadwalOffline1(data.offline_guide_1 || 'Membawa kartu identitas resmi (KTP untuk perorangan, Akta/SK pendirian untuk instansi).');
        setEditJadwalOffline2(data.offline_guide_2 || 'Loket Pelayanan berada di Lantai 1 Gedung Rektorat Kampus Utama UPERTIS Padang.');
        setEditJadwalOffline3(data.offline_guide_3 || 'Petugas kami siap membantu pencatatan, peninjauan berkas, hingga cetak salinan dokumen.');
        setEditJadwalOnline1(data.online_guide_1 || 'Gunakan menu Permohonan Informasi untuk mengisi formulir digital dan upload KTP.');
        setEditJadwalOnline2(data.online_guide_2 || 'Simpan nomor tiket pengajuan untuk melacak tanggapan admin secara real-time.');
        setEditJadwalOnline3(data.online_guide_3 || 'Jawaban atau link unduhan berkas digital akan dikirimkan langsung ke email/tiket Anda.');
        setEditJadwalCustomRemarks(data.custom_remarks || '');
      } catch (e) {
        setEditJadwalSeninKamisKerja('08:00 – 16:00 WIB');
        setEditJadwalSeninKamisIstirahat('12:00 – 13:30 WIB');
        setEditJadwalJumatKerja('08:00 – 16:30 WIB');
        setEditJadwalPageIstirahatJumat('12:00 – 14:00 WIB');
        setEditJadwalSabtuMinggu('Sistem Online Tetap Aktif 24/7 (Desk fisik Tutup)');
        setEditJadwalOffline1('Membawa kartu identitas resmi (KTP untuk perorangan, Akta/SK pendirian untuk instansi).');
        setEditJadwalOffline2('Loket Pelayanan berada di Lantai 1 Gedung Rektorat Kampus Utama UPERTIS Padang.');
        setEditJadwalOffline3('Petugas kami siap membantu pencatatan, peninjauan berkas, hingga cetak salinan dokumen.');
        setEditJadwalOnline1('Gunakan menu Permohonan Informasi untuk mengisi formulir digital dan upload KTP.');
        setEditJadwalOnline2('Simpan nomor tiket pengajuan untuk melacak tanggapan admin secara real-time.');
        setEditJadwalOnline3('Jawaban atau link unduhan berkas digital akan dikirimkan langsung ke email/tiket Anda.');
        setEditJadwalCustomRemarks(item.content ? item.content.replace(/<[^>]*>/g, '') : '');
      }
    } else if (type === 'page' && item.slug === 'profil') {
      // Parse profil stats from JSON content
      try {
        const parsed = JSON.parse(item.content);
        if (parsed && Array.isArray(parsed.stats) && parsed.stats.length > 0) {
          setEditProfilStats(parsed.stats);
        } else {
          setEditProfilStats([
            { value: '2021', label: 'TAHUN BERDIRI' },
            { value: '14', label: 'TAHUN UU KIP' },
            { value: 'PPID', label: 'UPERTIS MELAYANI' }
          ]);
        }
        setEditProfilImage1(parsed.image1 || '');
        setEditProfilImage2(parsed.image2 || '');
        setEditProfilImage3(parsed.image3 || '');
        setAdminEditContent(parsed.sejarah || item.content || '');
      } catch (e) {
        setEditProfilStats([{ value: '2021', label: 'TAHUN BERDIRI' }]);
        setEditProfilImage1('');
        setEditProfilImage2('');
        setEditProfilImage3('');
        setAdminEditContent(item.content || '');
      }
    } else if (type === 'page' && item.slug === 'tugas-dan-fungsi') {
      try {
        const parsed = JSON.parse(item.content);
        if (parsed && typeof parsed === 'object') {
          setEditTugasList(parsed.tugas || defaultTugasList);
          setEditFungsiList(parsed.fungsi || defaultFungsiList);
          setEditTugasImage1(parsed.image1 || '');
          setEditTugasImage2(parsed.image2 || '');
          setEditTugasImage3(parsed.image3 || '');
          setAdminEditContent(parsed.tambahan || '');
        } else {
          setEditTugasList(defaultTugasList);
          setEditFungsiList(defaultFungsiList);
          setEditTugasImage1('');
          setEditTugasImage2('');
          setEditTugasImage3('');
          setAdminEditContent(item.content || '');
        }
      } catch (e) {
        setEditTugasList(defaultTugasList);
        setEditFungsiList(defaultFungsiList);
        setEditTugasImage1('');
        setEditTugasImage2('');
        setEditTugasImage3('');
        setAdminEditContent(item.content || '');
      }
      setAdminEditPageDocs([]);
      setEditKeberatanManualSteps(defaultManualSteps);
      setEditKeberatanOnlineSteps(defaultOnlineSteps);
    } else if (type === 'page' && ['informasi-publik-berkala', 'informasi-tersedia-setiap-saat', 'info-serta-merta', 'informasi-dikecualikan', 'zona-integrasi', 'keberatan-informasi', 'Permohonan-penyelesaian-sengketa', 'permohonan-penyelesaian-sengketa'].includes(item.slug)) {
      try {
        const parsed = JSON.parse(item.content);
        if (parsed && typeof parsed === 'object') {
          if (Array.isArray(parsed.docs)) {
            setAdminEditPageDocs(parsed.docs);
          } else if (Array.isArray(parsed)) {
            setAdminEditPageDocs(parsed);
          } else {
            setAdminEditPageDocs([]);
          }
          if (item.slug === 'keberatan-informasi') {
            setEditKeberatanManualSteps(parsed.manual_steps || defaultManualSteps);
            setEditKeberatanOnlineSteps(parsed.online_steps || defaultOnlineSteps);
            setEditSengketaCards([]);
          } else if (item.slug === 'Permohonan-penyelesaian-sengketa' || item.slug === 'permohonan-penyelesaian-sengketa') {
            setEditSengketaCards(Array.isArray(parsed.cards) ? parsed.cards : []);
            setEditSengketaFlow(parsed.flow && typeof parsed.flow === 'object' ? {
              jalurA_title: parsed.flow.jalurA_title || '',
              jalurA_desc: parsed.flow.jalurA_desc || '',
              jalurB_title: parsed.flow.jalurB_title || '',
              jalurB_desc: parsed.flow.jalurB_desc || '',
              verifikasi_title: parsed.flow.verifikasi_title || '',
              verifikasi_desc: parsed.flow.verifikasi_desc || '',
              help_title: parsed.flow.help_title || '',
              help_desc: parsed.flow.help_desc || '',
              download_title: parsed.flow.download_title || '',
              download_desc: parsed.flow.download_desc || ''
            } : {
              jalurA_title: '', jalurA_desc: '',
              jalurB_title: '', jalurB_desc: '',
              verifikasi_title: '', verifikasi_desc: '',
              help_title: '', help_desc: '',
              download_title: '', download_desc: ''
            });
            setEditKeberatanManualSteps(defaultManualSteps);
            setEditKeberatanOnlineSteps(defaultOnlineSteps);
          } else {
            setEditKeberatanManualSteps(defaultManualSteps);
            setEditKeberatanOnlineSteps(defaultOnlineSteps);
            setEditSengketaCards([]);
          }
        } else {
          setAdminEditPageDocs([]);
          setAdminEditContent(item.content || '');
          setEditKeberatanManualSteps(defaultManualSteps);
          setEditKeberatanOnlineSteps(defaultOnlineSteps);
        }
      } catch (e) {
        setAdminEditPageDocs([]);
        setAdminEditContent(item.content || '');
        setEditKeberatanManualSteps(defaultManualSteps);
        setEditKeberatanOnlineSteps(defaultOnlineSteps);
      }
    } else {
      setAdminEditPageDocs([]);
      setEditProfilStats([]);
      setEditProfilImage1('');
      setEditProfilImage2('');
      setEditProfilImage3('');
      setEditTugasImage1('');
      setEditTugasImage2('');
      setEditTugasImage3('');
      setEditTugasList([]);
      setEditFungsiList([]);
      setEditKeberatanManualSteps(defaultManualSteps);
      setEditKeberatanOnlineSteps(defaultOnlineSteps);
    }

    setEditModalOpen(true);
  };

  // Executing CRUD API commands (Insert/Update)
  const handleSaveCrudItem = async (e: React.FormEvent) => {
    e.preventDefault();
    const token = localStorage.getItem('auth_token');
    if (!token) return;

    let url = '';
    let method = 'POST';
    let body: any = {};
    let table = '';

    if (editModalType === 'page') {
      table = 'pages';
      let pageContent = adminEditContent;
      if (adminEditSlug === 'jadwal-layanan-informasi') {
        pageContent = JSON.stringify({
          senin_kamis_kerja: editJadwalSeninKamisKerja,
          senin_kamis_istirahat: editJadwalSeninKamisIstirahat,
          jumat_kerja: editJadwalJumatKerja,
          jumat_istirahat: editJadwalPageIstirahatJumat,
          sabtu_minggu: editJadwalSabtuMinggu,
          offline_guide_1: editJadwalOffline1,
          offline_guide_2: editJadwalOffline2,
          offline_guide_3: editJadwalOffline3,
          online_guide_1: editJadwalOnline1,
          online_guide_2: editJadwalOnline2,
          online_guide_3: editJadwalOnline3,
          custom_remarks: editJadwalCustomRemarks
        });
      } else if (adminEditSlug === 'profil') {
        // Serialize profil page with sejarah text + stats cards
        pageContent = JSON.stringify({
          sejarah: adminEditContent,
          stats: editProfilStats,
          image1: editProfilImage1,
          image2: editProfilImage2,
          image3: editProfilImage3
        });
      } else if (adminEditSlug === 'tugas-dan-fungsi') {
        pageContent = JSON.stringify({
          tugas: editTugasList,
          fungsi: editFungsiList,
          tambahan: adminEditContent,
          image1: editTugasImage1,
          image2: editTugasImage2,
          image3: editTugasImage3
        });
      } else if (['informasi-publik-berkala', 'informasi-tersedia-setiap-saat', 'info-serta-merta', 'informasi-dikecualikan', 'zona-integrasi', 'keberatan-informasi', 'Permohonan-penyelesaian-sengketa', 'permohonan-penyelesaian-sengketa'].includes(adminEditSlug)) {
        if (adminEditSlug === 'keberatan-informasi') {
          pageContent = JSON.stringify({
            intro: adminEditContent,
            docs: adminEditPageDocs,
            manual_steps: editKeberatanManualSteps,
            online_steps: editKeberatanOnlineSteps
          });
        } else if (adminEditSlug === 'Permohonan-penyelesaian-sengketa' || adminEditSlug === 'permohonan-penyelesaian-sengketa') {
          // Filter out empty flow fields so only user-edited fields override defaults
          const flowToSave: Record<string, string> = {};
          Object.entries(editSengketaFlow).forEach(([k, v]) => {
            if (v && v.trim()) flowToSave[k] = v;
          });
          pageContent = JSON.stringify({
            intro: adminEditContent,
            docs: adminEditPageDocs,
            cards: editSengketaCards,
            flow: Object.keys(flowToSave).length > 0 ? flowToSave : undefined
          });
        } else {
          pageContent = JSON.stringify({
            intro: adminEditContent,
            docs: adminEditPageDocs
          });
        }
      }
      body = {
        title: adminEditTitle,
        subtitle: adminEditSubtitle,
        content: pageContent,
        slug: adminEditSlug || adminEditTitle.toLowerCase().replace(/[^a-z0-9]+/g, '-'),
        cover_image_url: adminEditCoverImage
      };
    } else if (editModalType === 'post') {
      table = 'posts';
      body = {
        title: adminEditTitle,
        excerpt: adminEditDescription || adminEditContent.substring(0, 150),
        content: adminEditContent,
        slug: adminEditSlug || adminEditTitle.toLowerCase().replace(/[^a-z0-9]+/g, '-'),
        cover_image_url: adminEditCoverImage,
        status: 'published',
        category: 'PPID'
      };
    } else if (editModalType === 'download') {
      table = 'downloads';
      body = {
        title: adminEditTitle,
        description: adminEditDescription,
        file_url: adminEditFileUrl,
        category: adminEditCategory || 'ppid-berkala',
        downloads_count: activeEditItem ? activeEditItem.downloads_count : 0
      };
    } else if (editModalType === 'gallery') {
      table = 'gallery';
      body = {
        title: adminEditTitle,
        description: adminEditDescription,
        media_url: adminEditFileUrl,
        media_type: adminEditMediaType,
        active: true
      };
    }

    if (activeEditItem) {
      method = 'PUT';
      url = `${API_BASE_URL}/admin/${table}/${activeEditItem.id}`;
    } else {
      url = `${API_BASE_URL}/admin/${table}`;
    }

    try {
      const res = await fetch(url, {
        method,
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`
        },
        body: JSON.stringify(body)
      });
      if (res.ok) {
        setEditModalOpen(false);
        setAdminGlobalMessage(`Data ${editModalType} berhasil disimpan!`);
        fetchAdminData();
      } else {
        alert(`Gagal menyimpan: ${(await res.json()).error}`);
      }
    } catch (err) {
      console.error(err);
      alert('Terjadi kesalahan jaringan.');
    }
  };

  // CRUD Delete function
  const handleDeleteCrudItem = async (type: string, id: string) => {
    if (!confirm('Apakah Anda yakin ingin menghapus data ini?')) return;
    const token = localStorage.getItem('auth_token');
    if (!token) return;

    let table = type;
    if (type === 'gallery') table = 'gallery';

    try {
      const res = await fetch(`${API_BASE_URL}/admin/${table}/${id}`, {
        method: 'DELETE',
        headers: { Authorization: `Bearer ${token}` }
      });
      if (res.ok) {
        setAdminGlobalMessage(`Item berhasil dihapus!`);
        fetchAdminData();
      } else {
        alert('Gagal menghapus item.');
      }
    } catch (err) {
      console.error(err);
      alert('Terjadi kesalahan jaringan.');
    }
  };

  // Handle Review Permohonan Action (Status change)
  const handlePermohonanActionSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!activeEditItem) return;
    const token = localStorage.getItem('auth_token');
    if (!token) return;

    setPermohonanActionLoading(true);
    const body = {
      status: permohonanActionStatus,
      admin_response: permohonanActionResponse
    };

    try {
      const res = await fetch(`${API_BASE_URL}/admin/permohonan_informasi/${activeEditItem.id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`
        },
        body: JSON.stringify(body)
      });
      if (res.ok) {
        setEditModalOpen(false);
        setAdminGlobalMessage('Permohonan berhasil diperbarui!');
        fetchAdminData();
      } else {
        alert('Gagal memproses permohonan.');
      }
    } catch (err) {
      console.error(err);
      alert('Terjadi kesalahan koneksi.');
setPermohonanActionLoading(false);
    }
  };

  // Dynamic Google Font Injection
  useEffect(() => {
    if (!portalFontFamily) return;
    const fontId = 'dynamic-google-font';
    let linkElement = document.getElementById(fontId) as HTMLLinkElement;
    if (!linkElement) {
      linkElement = document.createElement('link');
      linkElement.id = fontId;
      linkElement.rel = 'stylesheet';
      document.head.appendChild(linkElement);
    }
    const formattedFont = portalFontFamily.replace(/\s+/g, '+');
    linkElement.href = `https://fonts.googleapis.com/css2?family=${formattedFont}:wght@300;400;550;600;700;800;900&display=swap`;
  }, [portalFontFamily]);

  // Determine text size tailwind wrapper class
  const getFontSizeClass = () => {
    if (portalFontSize === 'medium') return 'text-[15.5px] [&_p]:text-[15.5px] [&_span]:text-[14.5px]';
    if (portalFontSize === 'large') return 'text-[17px] [&_p]:text-[17px] [&_span]:text-[16px]';
    return ''; // Normal (Default)
  };

  return (
    <div
      style={{
        fontFamily: portalFontFamily ? `"${portalFontFamily}", var(--font-sans)` : undefined,
        ...(siteConfig?.settings?.background_image ? {
          backgroundImage: `url(${resolveImageUrl(siteConfig.settings.background_image)})`,
          backgroundSize: 'cover',
          backgroundPosition: 'center',
          backgroundAttachment: 'fixed',
          backgroundRepeat: 'no-repeat'
        } : {})
      }}
      className={`min-h-screen bg-slate-100 text-slate-900 selection:bg-[#002147] selection:text-white flex flex-col font-sans ${getFontSizeClass()}`}
    >
      <div className="w-full max-w-[1440px] mx-auto bg-white shadow-2xl flex flex-col min-h-screen relative animate-in fade-in duration-300">
        
        {/* --- GLOBAL HEADER (NAVBAR) --- */}
        <Navbar
          siteConfig={siteConfig}
          currentPage={currentPage}
          adminUser={adminUser}
          menuGroups={menuGroups}
          navigateToHome={navigateToHome}
          navigateToNews={navigateToNews}
          handleNavigation={handleNavigation}
          setCurrentPage={setCurrentPage}
          setAdminActiveTab={setAdminActiveTab}
        />

        {/* --- MAIN PAGE VIEW CONTENT --- */}
        <main className="flex-1 pb-20 lg:pb-0">
          {currentPage === 'home' && (
            <Home
              siteConfig={siteConfig}
              heroImages={heroImages}
              heroIdx={heroIdx}
              prevSlide={prevSlide}
              nextSlide={nextSlide}
              posts={posts}
              isPostsLoading={isPostsLoading}
              dbDownloads={dbDownloads}
              navigateToNews={navigateToNews}
              navigateToPage={navigateToPage}
              navigateToNewsDetail={navigateToNewsDetail}
              handleNavigation={handleNavigation}
              incrementDownloadCount={incrementDownloadCount}
            />
          )}

          {currentPage === 'berita' && (
            <News
              posts={posts}
              isPostsLoading={isPostsLoading}
              newsFilterCategory={newsFilterCategory}
              newsFilterSearch={newsFilterSearch}
              setNewsFilterCategory={setNewsFilterCategory}
              setNewsFilterSearch={setNewsFilterSearch}
              navigateToNewsDetail={navigateToNewsDetail}
              navigateToHome={navigateToHome}
            />
          )}

          {currentPage === 'berita-detail' && (
            <NewsDetail
              selectedPost={selectedPost}
              isPageLoading={isPageLoading}
              sidebarSearch={sidebarSearch}
              setSidebarSearch={setSidebarSearch}
              posts={posts}
              navigateToNews={navigateToNews}
              setCurrentPage={setCurrentPage}
              setNewsFilterSearch={setNewsFilterSearch}
              setNewsFilterCategory={setNewsFilterCategory}
            />
          )}

          {currentPage.startsWith('page/') && (
            <DynamicPage
              activeSlug={activeSlug}
              pageData={pageData}
              isPageLoading={isPageLoading}
              siteConfig={siteConfig}
              navigateToHome={navigateToHome}
              navigateToPage={navigateToPage}
              handleNavigation={handleNavigation}
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
              expandedRegulasi={expandedRegulasi}
              toggleRegulasi={toggleRegulasi}
              activeTab={activeTab}
              setActiveTab={setActiveTab}
              searchTerm={searchTerm}
              setSearchTerm={setSearchTerm}
              dbDownloads={dbDownloads}
              incrementDownloadCount={incrementDownloadCount}
              galleries={galleries}
              activeLightboxImage={activeLightboxImage}
              setActiveLightboxImage={setActiveLightboxImage}
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
          )}

          {currentPage === 'admin' && (
            !adminUser ? (
              <AdminLogin
                adminEmail={adminEmail}
                setAdminEmail={setAdminEmail}
                adminPassword={adminPassword}
                setAdminPassword={setAdminPassword}
                showPassword={showPassword}
                setShowPassword={setShowPassword}
                adminLoginLoading={adminLoginLoading}
                adminLoginError={adminLoginError}
                handleAdminLogin={handleAdminLogin}
              />
            ) : (
              <AdminLayout
                adminUser={adminUser}
                setAdminUser={setAdminUser}
                adminActiveTab={adminActiveTab}
                setAdminActiveTab={setAdminActiveTab}
                fetchAdminData={fetchAdminData}
                setCurrentPage={setCurrentPage}
                adminGlobalMessage={adminGlobalMessage}
                setAdminGlobalMessage={setAdminGlobalMessage}
                siteConfig={siteConfig}
              >
                {adminActiveTab === 'dashboard' && (
                  <DashboardOverview
                    adminStats={adminStats}
                    adminPermohonans={adminPermohonans}
                    setAdminActiveTab={setAdminActiveTab}
                  />
                )}

                {adminActiveTab === 'slider' && (
                  <ManageSlider
                    siteConfig={siteConfig}
                    heroImages={heroImages}
                    API_BASE_URL={API_BASE_URL}
                    setAdminGlobalMessage={setAdminGlobalMessage}
                    setSiteConfig={setSiteConfig}
                  />
                )}

                {adminActiveTab === 'downloads' && (
                  <ManageDownloads
                    editModalOpen={editModalOpen}
                    setEditModalOpen={setEditModalOpen}
                    editModalType={editModalType}
                    activeEditItem={activeEditItem}
                    adminEditTitle={adminEditTitle}
                    setAdminEditTitle={setAdminEditTitle}
                    adminEditDescription={adminEditDescription}
                    setAdminEditDescription={setAdminEditDescription}
                    adminEditFileUrl={adminEditFileUrl}
                    setAdminEditFileUrl={setAdminEditFileUrl}
                    adminEditCategory={adminEditCategory}
                    setAdminEditCategory={setAdminEditCategory}
                    handleSaveCrudItem={handleSaveCrudItem}
                    handleAdminUpload={handleAdminUpload}
                    openCreateModal={openCreateModal}
                    adminDownloads={adminDownloads}
                    openEditModal={openEditModal}
                    handleDeleteCrudItem={handleDeleteCrudItem}
                  />
                )}

                {adminActiveTab === 'gallery' && (
                  <ManageGallery
                    editModalOpen={editModalOpen}
                    setEditModalOpen={setEditModalOpen}
                    editModalType={editModalType}
                    activeEditItem={activeEditItem}
                    adminEditTitle={adminEditTitle}
                    setAdminEditTitle={setAdminEditTitle}
                    adminEditDescription={adminEditDescription}
                    setAdminEditDescription={setAdminEditDescription}
                    adminEditFileUrl={adminEditFileUrl}
                    setAdminEditFileUrl={setAdminEditFileUrl}
                    adminEditMediaType={adminEditMediaType}
                    setAdminEditMediaType={setAdminEditMediaType}
                    handleSaveCrudItem={handleSaveCrudItem}
                    handleAdminUpload={handleAdminUpload}
                    openCreateModal={openCreateModal}
                    adminGalleries={adminGalleries}
                    openEditModal={openEditModal}
                    handleDeleteCrudItem={handleDeleteCrudItem}
                  />
                )}

                {adminActiveTab === 'pages' && (
                  <ManagePages
                    editModalOpen={editModalOpen}
                    setEditModalOpen={setEditModalOpen}
                    editModalType={editModalType}
                    activeEditItem={activeEditItem}
                    adminEditTitle={adminEditTitle}
                    setAdminEditTitle={setAdminEditTitle}
                    adminEditSubtitle={adminEditSubtitle}
                    setAdminEditSubtitle={setAdminEditSubtitle}
                    adminEditSlug={adminEditSlug}
                    setAdminEditSlug={setAdminEditSlug}
                    adminEditCoverImage={adminEditCoverImage}
                    setAdminEditCoverImage={setAdminEditCoverImage}
                    handleSaveCrudItem={handleSaveCrudItem}
                    handleAdminUpload={handleAdminUpload}
                    openCreateModal={openCreateModal}
                    adminPages={adminPages}
                    openEditModal={openEditModal}
                    handleDeleteCrudItem={handleDeleteCrudItem}
                    API_BASE_URL={API_BASE_URL}
                    editJadwalSeninKamisKerja={editJadwalSeninKamisKerja}
                    setEditJadwalSeninKamisKerja={setEditJadwalSeninKamisKerja}
                    editJadwalSeninKamisIstirahat={editJadwalSeninKamisIstirahat}
                    setEditJadwalSeninKamisIstirahat={setEditJadwalSeninKamisIstirahat}
                    editJadwalJumatKerja={editJadwalJumatKerja}
                    setEditJadwalJumatKerja={setEditJadwalJumatKerja}
                    editJadwalPageIstirahatJumat={editJadwalPageIstirahatJumat}
                    setEditJadwalPageIstirahatJumat={setEditJadwalPageIstirahatJumat}
                    editJadwalSabtuMinggu={editJadwalSabtuMinggu}
                    setEditJadwalSabtuMinggu={setEditJadwalSabtuMinggu}
                    editJadwalOffline1={editJadwalOffline1}
                    setEditJadwalOffline1={setEditJadwalOffline1}
                    editJadwalOffline2={editJadwalOffline2}
                    setEditJadwalOffline2={setEditJadwalOffline2}
                    editJadwalOffline3={editJadwalOffline3}
                    setEditJadwalOffline3={setEditJadwalOffline3}
                    editJadwalOnline1={editJadwalOnline1}
                    setEditJadwalOnline1={setEditJadwalOnline1}
                    editJadwalOnline2={editJadwalOnline2}
                    setEditJadwalOnline2={setEditJadwalOnline2}
                    editJadwalOnline3={editJadwalOnline3}
                    setEditJadwalOnline3={setEditJadwalOnline3}
                    editJadwalCustomRemarks={editJadwalCustomRemarks}
                    setEditJadwalCustomRemarks={setEditJadwalCustomRemarks}
                    adminEditPageDocs={adminEditPageDocs}
                    setAdminEditPageDocs={setAdminEditPageDocs}
                    adminEditContent={adminEditContent}
                    setAdminEditContent={setAdminEditContent}
                    editKeberatanManualSteps={editKeberatanManualSteps}
                    setEditKeberatanManualSteps={setEditKeberatanManualSteps}
                    editKeberatanOnlineSteps={editKeberatanOnlineSteps}
                    setEditKeberatanOnlineSteps={setEditKeberatanOnlineSteps}
                    editProfilStats={editProfilStats}
                    setEditProfilStats={setEditProfilStats}
                    editProfilImage1={editProfilImage1}
                    setEditProfilImage1={setEditProfilImage1}
                    editProfilImage2={editProfilImage2}
                    setEditProfilImage2={setEditProfilImage2}
                    editProfilImage3={editProfilImage3}
                    setEditProfilImage3={setEditProfilImage3}
                    editTugasList={editTugasList}
                    setEditTugasList={setEditTugasList}
                    editFungsiList={editFungsiList}
                    setEditFungsiList={setEditFungsiList}
                    editTugasImage1={editTugasImage1}
                    setEditTugasImage1={setEditTugasImage1}
                    editTugasImage2={editTugasImage2}
                    setEditTugasImage2={setEditTugasImage2}
                    editTugasImage3={editTugasImage3}
                    setEditTugasImage3={setEditTugasImage3}
                    editSengketaCards={editSengketaCards}
                    setEditSengketaCards={setEditSengketaCards}
                    editSengketaFlow={editSengketaFlow}
                    setEditSengketaFlow={setEditSengketaFlow}
                  />
                )}

                {adminActiveTab === 'posts' && (
                  <ManageNews
                    editModalOpen={editModalOpen}
                    setEditModalOpen={setEditModalOpen}
                    editModalType={editModalType}
                    activeEditItem={activeEditItem}
                    adminEditTitle={adminEditTitle}
                    setAdminEditTitle={setAdminEditTitle}
                    adminEditSlug={adminEditSlug}
                    setAdminEditSlug={setAdminEditSlug}
                    adminEditCoverImage={adminEditCoverImage}
                    setAdminEditCoverImage={setAdminEditCoverImage}
                    adminEditDescription={adminEditDescription}
                    setAdminEditDescription={setAdminEditDescription}
                    adminEditContent={adminEditContent}
                    setAdminEditContent={setAdminEditContent}
                    handleSaveCrudItem={handleSaveCrudItem}
                    handleAdminUpload={handleAdminUpload}
                    openCreateModal={openCreateModal}
                    adminPosts={adminPosts}
                    openEditModal={openEditModal}
                    handleDeleteCrudItem={handleDeleteCrudItem}
                  />
                )}

                {adminActiveTab === 'permohonan' && (
                  <ManagePermohonan
                    editModalOpen={editModalOpen}
                    setEditModalOpen={setEditModalOpen}
                    editModalType={editModalType}
                    activeEditItem={activeEditItem}
                    setActiveEditItem={setActiveEditItem}
                    permohonanActionStatus={permohonanActionStatus}
                    setPermohonanActionStatus={setPermohonanActionStatus}
                    permohonanActionResponse={permohonanActionResponse}
                    setPermohonanActionResponse={setPermohonanActionResponse}
                    permohonanActionLoading={permohonanActionLoading}
                    handlePermohonanActionSubmit={handlePermohonanActionSubmit}
                    adminPermohonans={adminPermohonans}
                    setEditModalType={setEditModalType}
                  />
                )}

                {adminActiveTab === 'settings' && (
                  <PortalSettings
                    portalName={portalName}
                    setPortalName={setPortalName}
                    portalDesc={portalDesc}
                    setPortalDesc={setPortalDesc}
                    portalLogo={portalLogo}
                    setPortalLogo={setPortalLogo}
                    portalBgImage={portalBgImage}
                    setPortalBgImage={setPortalBgImage}
                    portalFontFamily={portalFontFamily}
                    setPortalFontFamily={setPortalFontFamily}
                    portalFontSize={portalFontSize}
                    setPortalFontSize={setPortalFontSize}
                    portalFavicon={portalFavicon}
                    setPortalFavicon={setPortalFavicon}
                    portalWelcomeText={portalWelcomeText}
                    setPortalWelcomeText={setPortalWelcomeText}
                    portalSkRektor={portalSkRektor}
                    setPortalSkRektor={setPortalSkRektor}
                    portalKeberatanLink={portalKeberatanLink}
                    setPortalKeberatanLink={setPortalKeberatanLink}
                    portalAboutStatNumber={portalAboutStatNumber}
                    setPortalAboutStatNumber={setPortalAboutStatNumber}
                    portalAboutStatLabelAccent={portalAboutStatLabelAccent}
                    setPortalAboutStatLabelAccent={setPortalAboutStatLabelAccent}
                    portalAboutStatLabel={portalAboutStatLabel}
                    setPortalAboutStatLabel={setPortalAboutStatLabel}
                    portalCard1Title={portalCard1Title}
                    setPortalCard1Title={setPortalCard1Title}
                    portalCard1Desc={portalCard1Desc}
                    setPortalCard1Desc={setPortalCard1Desc}
                    portalCard2Title={portalCard2Title}
                    setPortalCard2Title={setPortalCard2Title}
                    portalCard2Desc={portalCard2Desc}
                    setPortalCard2Desc={setPortalCard2Desc}
                    portalCard3Title={portalCard3Title}
                    setPortalCard3Title={setPortalCard3Title}
                    portalCard3Desc={portalCard3Desc}
                    setPortalCard3Desc={setPortalCard3Desc}
                    portalCard1Link={portalCard1Link}
                    setPortalCard1Link={setPortalCard1Link}
                    portalCard2Link={portalCard2Link}
                    setPortalCard2Link={setPortalCard2Link}
                    portalCard3Link={portalCard3Link}
                    setPortalCard3Link={setPortalCard3Link}
                    portalFaqs={portalFaqs}
                    setPortalFaqs={setPortalFaqs}
                    portalPermohonanLink={portalPermohonanLink}
                    setPortalPermohonanLink={setPortalPermohonanLink}
                    portalPermohonanFormType={portalPermohonanFormType}
                    setPortalPermohonanFormType={setPortalPermohonanFormType}
                    portalPengaduanLink={portalPengaduanLink}
                    setPortalPengaduanLink={setPortalPengaduanLink}
                    portalRektoratEmail={portalRektoratEmail}
                    setPortalRektoratEmail={setPortalRektoratEmail}
                    portalRektoratPhone={portalRektoratPhone}
                    setPortalRektoratPhone={setPortalRektoratPhone}
                    portalRektoratAddress={portalRektoratAddress}
                    setPortalRektoratAddress={setPortalRektoratAddress}
                    portalKampus2Address={portalKampus2Address}
                    setPortalKampus2Address={setPortalKampus2Address}
                    portalKampus1MapUrl={portalKampus1MapUrl}
                    setPortalKampus1MapUrl={setPortalKampus1MapUrl}
                    portalKampus2MapUrl={portalKampus2MapUrl}
                    setPortalKampus2MapUrl={setPortalKampus2MapUrl}
                    portalJadwalSeninKamis={portalJadwalSeninKamis}
                    setPortalJadwalSeninKamis={setPortalJadwalSeninKamis}
                    portalIstirahatSeninKamis={portalIstirahatSeninKamis}
                    setPortalIstirahatSeninKamis={setPortalIstirahatSeninKamis}
                    portalJadwalJumat={portalJadwalJumat}
                    setPortalJadwalJumat={setPortalJadwalJumat}
                    portalIstirahatJumat={portalIstirahatJumat}
                    setPortalIstirahatJumat={setPortalIstirahatJumat}
                    portalJadwalSabtuMinggu={portalJadwalSabtuMinggu}
                    setPortalJadwalSabtuMinggu={setPortalJadwalSabtuMinggu}
                    portalIsSaving={portalIsSaving}
                    handleSaveSettings={handleSaveSettings}
                    API_BASE_URL={API_BASE_URL}
                  />
                )}

                {adminActiveTab === 'menu-manager' && (
                  <MenuManager
                    siteConfig={siteConfig}
                    setSiteConfig={setSiteConfig}
                    adminPages={adminPages}
                    API_BASE_URL={API_BASE_URL}
                    setAdminGlobalMessage={setAdminGlobalMessage}
                    fetchAdminData={fetchAdminData}
                  />
                )}
              </AdminLayout>
            )
          )}
        </main>

        {/* --- GLOBAL FOOTER --- */}
        <Footer
          navigateToHome={navigateToHome}
          navigateToPage={navigateToPage}
          setCurrentPage={setCurrentPage}
        />
      </div>

      {/* === FLOATING WIDGETS === */}

      {/* FAQ Chatbot Widget */}
      <ChatWidget faqs={portalFaqs} rektoratPhone={siteConfig?.settings?.rektorat_phone} />

      {/* Scroll to Top Button */}
      <button
        onClick={handleScrollToTop}
        title="Kembali ke atas"
        className={`fixed bottom-20 lg:bottom-5 left-5 z-50 w-[52px] h-[52px] rounded-full bg-[#002147] text-white shadow-lg shadow-slate-900/30 flex items-center justify-center hover:bg-amber-500 hover:shadow-amber-400/30 hover:shadow-xl active:scale-95 transition-all duration-300 ${
          showScrollTop ? 'opacity-100 translate-y-0 pointer-events-auto' : 'opacity-0 translate-y-4 pointer-events-none'
        }`}
      >
        <svg viewBox="0 0 24 24" className="w-5 h-5 fill-none stroke-current" strokeWidth={2.5} strokeLinecap="round" strokeLinejoin="round">
          <polyline points="18 15 12 9 6 15" />
        </svg>
      </button>

    </div>
  );
}

export default App;
