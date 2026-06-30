import {
  Home, User, BookOpen, Briefcase, Newspaper, Image, Download, FileText,
  Target, Info, Users, FileCheck, Clock, CheckCircle, AlertTriangle, ShieldAlert
} from 'lucide-react';

export const API_BASE_URL = import.meta.env.VITE_API_URL || '/api/v1';

// Resolve Image Utility
export const resolveImageUrl = (path?: string) => {
  if (!path) return '';
  if (path.startsWith('http://') || path.startsWith('https://') || path.startsWith('data:')) {
    return path;
  }
  const cleanBase = API_BASE_URL.replace('/api/v1', '');
  return `${cleanBase}${path.startsWith('/') ? '' : '/'}${path}`;
};

// HTML cleanup renderer for static pages
export const preprocessPostContent = (htmlContent: string) => {
  if (!htmlContent) return '';
  let formatted = htmlContent.replace(/\r\n/g, '<br/>');
  const baseServer = API_BASE_URL.replace('/api/v1', '');
  formatted = formatted.replace(/src=["']\/uploads\//g, `src="${baseServer}/uploads/`);
  return formatted;
};

export const getHeaderNavIcon = (label: string) => {
  const l = label.toLowerCase();
  if (l === 'home' || l === 'beranda') return Home;
  if (l.includes('tentang')) return User;
  if (l.includes('info') || l.includes('publik')) return BookOpen;
  if (l.includes('layanan')) return Briefcase;
  if (l.includes('berita')) return Newspaper;
  if (l.includes('galeri')) return Image;
  if (l.includes('download')) return Download;
  return FileText;
};

// Metadata mappings for styling PPID menu items in layout
export const getPPIDMenuItemMeta = (title: string) => {
  const t = title.toLowerCase();
  if (t.includes('profil')) return { icon: FileText, color: 'text-blue-500 bg-blue-50 border-blue-100', desc: 'Profil singkat & bagan struktur PPID' };
  if (t.includes('visi') || t.includes('misi')) return { icon: Target, color: 'text-indigo-500 bg-indigo-50 border-indigo-100', desc: 'Visi, misi & arah gerak layanan PPID' };
  if (t.includes('maklumat')) return { icon: Info, color: 'text-emerald-500 bg-emerald-50 border-emerald-100', desc: 'Maklumat komitmen resmi standard pelayanan' };
  if (t.includes('struktur')) return { icon: Users, color: 'text-amber-500 bg-amber-50 border-amber-100', desc: 'Bagan & susunan keanggotaan PPID' };
  if (t.includes('tugas') || t.includes('fungsi')) return { icon: FileCheck, color: 'text-rose-500 bg-rose-50 border-rose-100', desc: 'Tanggung jawab & kewajiban tugas PPID' };
  if (t.includes('regulasi')) return { icon: BookOpen, color: 'text-purple-500 bg-purple-50 border-purple-100', desc: 'Undang-Undang & landasan hukum KIP' };
  if (t.includes('berkala')) return { icon: Clock, color: 'text-blue-500 bg-blue-50 border-blue-100', desc: 'Daftar informasi berkala yang dirilis rutin' };
  if (t.includes('setiap saat')) return { icon: CheckCircle, color: 'text-emerald-500 bg-emerald-50 border-emerald-100', desc: 'Daftar informasi yang wajib tersedia' };
  if (t.includes('serta merta') || t.includes('serta serta')) return { icon: AlertTriangle, color: 'text-amber-500 bg-amber-50 border-amber-100', desc: 'Informasi penting darurat bagi publik' };
  if (t.includes('zona')) return { icon: Info, color: 'text-indigo-500 bg-indigo-50 border-indigo-100', desc: 'Komitmen wilayah bebas korupsi (WBK)' };
  if (t.includes('jadwal')) return { icon: Clock, color: 'text-teal-500 bg-teal-50 border-teal-100', desc: 'Waktu operasional loket layanan informasi' };
  if (t.includes('permohonan')) return { icon: FileText, color: 'text-blue-500 bg-blue-50 border-blue-100', desc: 'Formulir online permohonan dokumen' };
  if (t.includes('keberatan')) return { icon: ShieldAlert, color: 'text-red-500 bg-red-50 border-red-100', desc: 'Pengajuan keberatan atas layanan informasi' };

  return { icon: FileText, color: 'text-slate-500 bg-slate-50 border-slate-100', desc: 'Halaman informasi resmi PPID UPERTIS' };
};

// Default Indonesian Regulasi KIP Documents for Accordion Layout
export const defaultRegulasiData = [
  {
    category: 'Regulasi Nasional',
    key: 'nasional',
    items: [
      {
        number: 'UU No. 14 Tahun 2008',
        title: 'Undang-Undang tentang Keterbukaan Informasi Publik (KIP)',
        url: 'https://peraturan.go.id/id/uu-no-14-tahun-2008'
      },
      {
        number: 'PP No. 61 Tahun 2010',
        title: 'Peraturan Pemerintah Pelaksanaan UU KIP',
        url: 'https://peraturan.go.id/id/pp-no-61-tahun-2010'
      },
      {
        number: 'Perki No. 1 Tahun 2021',
        title: 'Peraturan Komisi Informasi tentang Standar Layanan Informasi Publik',
        url: 'https://komisiinformasi.go.id/regulasi/detail/perki-no-1-tahun-2021'
      }
    ]
  },
  {
    category: 'Regulasi Internal UPERTIS',
    key: 'internal',
    items: [
      {
        number: 'SK Rektor No. 008A/2025',
        title: 'Surat Keputusan Rektor tentang Penetapan Pejabat Pengelola Informasi dan Dokumentasi (PPID) Universitas Perintis Indonesia',
        url: '#'
      },
      {
        number: 'SK PPID Utama No. 010B/2025',
        title: 'Penetapan Daftar Informasi Publik (DIP) UPERTIS yang Dimutakhirkan Tahun 2025',
        url: '#'
      },
      {
        number: 'SK PPID Utama No. 010C/2025',
        title: 'Penetapan Daftar Informasi Publik yang Dikecualikan Tahun 2025',
        url: '#'
      }
    ]
  },
  {
    category: 'Maklumat & SOP Pelayanan',
    key: 'maklumat',
    items: [
      {
        number: 'Maklumat Pelayanan',
        title: 'Pernyataan Kesanggupan Memberikan Pelayanan Informasi Sesuai Standar Layanan',
        url: '#'
      },
      {
        number: 'SOP Pelayanan',
        title: 'Standar Operasional Prosedur Pelayanan Informasi Publik PPID UPERTIS',
        url: '#'
      }
    ]
  }
];

export const parseRegulasiHTML = (htmlContent?: string) => {
  if (!htmlContent) return defaultRegulasiData;

  const trimmed = htmlContent.trim();
  if (trimmed.startsWith('[') || trimmed.startsWith('{')) {
    try {
      const parsed = JSON.parse(trimmed);
      if (Array.isArray(parsed)) {
        return parsed.length > 0 ? parsed : defaultRegulasiData;
      }
    } catch (e) {
      console.error('Failed to parse htmlContent JSON in parseRegulasiHTML:', e);
    }
  }

  let cleanText = htmlContent
    .replace(/<br\s*\/?>/gi, '\n')
    .replace(/<\/p>/gi, '\n')
    .replace(/<p[^>]*>/gi, '')
    .replace(/<\/div>/gi, '\n')
    .replace(/<div[^>]*>/gi, '')
    .replace(/&nbsp;/g, ' ')
    .replace(/\r\n/g, '\n');

  const lines = cleanText.split('\n').map(l => l.trim()).filter(Boolean);

  const groups: any[] = [];
  let currentGroup: any = null;
  let currentItem: any = null;

  lines.forEach(line => {
    const groupMatch = line.match(/^([A-C])\.\s+(.*)/i);
    if (groupMatch) {
      const key = groupMatch[1].toUpperCase();
      const label = groupMatch[0];
      const desc = key === 'A'
        ? 'Undang-Undang, Peraturan Pemerintah, dan Peraturan Komisi Informasi Nasional.'
        : key === 'B'
          ? 'Surat Keputusan Rektor, Ketetapan DIP PPID Utama, dan Standar Operasional Prosedur (SOP) Internal.'
          : 'Pedoman Pengelolaan Organisasi dan Administrasi resmi Universitas Perintis Indonesia.';

      currentGroup = {
        key,
        label,
        desc,
        items: []
      };
      groups.push(currentGroup);
      currentItem = null;
      return;
    }

    if (!currentGroup) {
      currentGroup = {
        key: 'A',
        label: 'A. Regulasi Nasional (Dasar Hukum)',
        desc: 'Undang-Undang, Peraturan Pemerintah, dan Peraturan Komisi Informasi Nasional.',
        items: []
      };
      groups.push(currentGroup);
    }

    const sopMatch = line.match(/^[\*\-]\s+(.*)/);
    if (sopMatch && currentItem) {
      currentItem.isSopList = true;
      currentItem.sops = [...(currentItem.sops || []), sopMatch[1].trim()];
      return;
    }

    const itemMatch = line.match(/^(\d+)\.\s+(.*)/);
    if (itemMatch) {
      const no = itemMatch[1].padStart(2, '0');
      let text = itemMatch[2].trim();
      let fileUrl = '#';

      const pipeSplit = text.split('|');
      if (pipeSplit.length > 1) {
        text = pipeSplit[0].trim();
        fileUrl = pipeSplit[1].trim();
      } else {
        const bracketMatch = text.match(/(https?:\/\/[^\s\)]+)/i);
        if (bracketMatch) {
          fileUrl = bracketMatch[1];
          text = text.replace(bracketMatch[0], '').replace(/\(\s*\)/g, '').trim();
        }
      }

      let title = text;
      let detail = '';

      const tentangIndex = text.toLowerCase().indexOf(' tentang ');
      if (tentangIndex !== -1) {
        title = text.substring(0, tentangIndex).trim();
        detail = text.substring(tentangIndex).trim();
      }

      currentItem = {
        no,
        title,
        detail,
        fileUrl
      };
      currentGroup.items.push(currentItem);
      return;
    }

    if (currentItem && line.length > 3) {
      currentItem.detail = currentItem.detail ? `${currentItem.detail} ${line}` : line;
    }
  });

  if (groups.length === 0 || (groups.length === 1 && groups[0].items.length === 0)) {
    return defaultRegulasiData;
  }

  return groups;
};
