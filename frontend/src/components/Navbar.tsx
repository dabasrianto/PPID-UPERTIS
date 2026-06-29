import { useState } from 'react';
import {
  Menu, X, Home, ChevronDown, ChevronRight, Newspaper, Image, Download, LayoutDashboard, Lock, Globe, FileText, MessageSquare
} from 'lucide-react';
import {
  resolveImageUrl,
  getHeaderNavIcon,
  getPPIDMenuItemMeta
} from '../utils/helpers';

interface NavbarProps {
  siteConfig: any;
  currentPage: string;
  adminUser: any;
  menuGroups: any[];
  navigateToHome: () => void;
  navigateToNews: () => void;
  handleNavigation: (href: string) => void;
  setCurrentPage: (page: string) => void;
  setAdminActiveTab: (tab: string) => void;
}

const getMenuItemIcon = (slug: string, label: string) => {
  const s = slug.toLowerCase();
  const l = label.toLowerCase();
  if (s === '/' || s === 'home' || l.includes('home') || l.includes('beranda')) return Home;
  if (s.includes('berita') || l.includes('berita') || l.includes('warta')) return Newspaper;
  if (s.includes('galeri') || l.includes('galeri') || l.includes('foto')) return Image;
  if (s.includes('download') || l.includes('download') || l.includes('unduh')) return Download;
  if (s.includes('kontak') || l.includes('kontak') || l.includes('hubungi')) return Globe;
  return FileText;
};

export default function Navbar({
  siteConfig,
  currentPage,
  adminUser,
  menuGroups,
  navigateToHome,
  navigateToNews,
  handleNavigation,
  setCurrentPage,
  setAdminActiveTab
}: NavbarProps) {
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const [activeMenuDropdown, setActiveMenuDropdown] = useState<string | null>(null);
  
  // Mobile Bottom Sheet States
  const [activeBottomSheet, setActiveBottomSheet] = useState<'layanan' | 'menu' | null>(null);
  const [mobileAccordionActive, setMobileAccordionActive] = useState<string | null>(null);

  const onNavClick = (href: string) => {
    setMobileMenuOpen(false);
    setActiveBottomSheet(null);
    handleNavigation(href);
  };

  // Normalize dynamic menu structure
  const normalizedMenu = menuGroups.map(item => {
    if (item.group && !item.label) {
      return {
        label: item.group,
        type: 'dropdown',
        items: item.items || []
      };
    }
    return item;
  });

  // Extract service items for Layanan Bottom Sheet
  const layananMenu = normalizedMenu.find(item => 
    item.label && item.label.toLowerCase().includes('layanan')
  );
  
  const servicesList = layananMenu?.items || [
    { label: 'Permohonan Informasi', href: 'permohonan-informasi' },
    { label: 'Keberatan Informasi', href: 'keberatan-informasi' },
    { label: 'Penyelesaian Sengketa', href: 'Permohonan-penyelesaian-sengketa' },
    { label: 'Jadwal Layanan', href: 'jadwal-layanan-informasi' }
  ];

  // Rest of the menus (excluding Layanan and Home/Berita to show in Menu Bottom Sheet)
  const otherMenus = normalizedMenu.filter(item => {
    const label = (item.label || '').toLowerCase();
    const href = (item.href || '').toLowerCase();
    return !label.includes('layanan') && !label.includes('beranda') && !label.includes('home') && href !== '/' && href !== 'berita';
  });

  return (
    <header className="sticky top-0 z-50 bg-[#002147] text-white shadow-lg">
      <div className="container mx-auto px-4 py-4 flex items-center justify-between">
        
        {/* Logo and branding title */}
        <div className="flex items-center gap-3 cursor-pointer" onClick={navigateToHome}>
          {siteConfig?.logo_url ? (
            <img
              src={resolveImageUrl(siteConfig.logo_url)}
              alt="Logo"
              className="h-10 w-10 object-contain bg-white rounded-lg p-0.5"
              onError={(e) => {
                e.currentTarget.src = '/logo.png';
              }}
            />
          ) : (
            <div className="h-10 w-10 bg-amber-400 text-[#002147] rounded-lg flex items-center justify-center font-bold text-sm">
              PPID
            </div>
          )}
          <div className="text-left">
            <span className="text-xs text-amber-400 block font-bold tracking-widest uppercase leading-none">Portal Resmi</span>
            <span className="text-sm font-extrabold tracking-tight block mt-0.5">
              {siteConfig?.name || 'PPID Universitas Perintis Indonesia'}
            </span>
          </div>
        </div>

        {/* Desktop Navbar menu options */}
        <nav className="hidden lg:flex items-center gap-6">
          {normalizedMenu.map((item: any, idx: number) => {
            const isDropdown = item.items && Array.isArray(item.items) && item.items.length > 0;

            if (isDropdown) {
              const GroupIcon = getHeaderNavIcon(item.label);
              return (
                <div
                  key={item.label || idx}
                  className="relative group/nav"
                  onMouseEnter={() => setActiveMenuDropdown(item.label)}
                  onMouseLeave={() => setActiveMenuDropdown(null)}
                >
                  <button className="text-xs font-bold uppercase tracking-wider text-slate-200 hover:text-amber-400 flex items-center gap-1.5 cursor-pointer bg-transparent border-0 p-0">
                    <GroupIcon className="h-3.5 w-3.5 text-slate-400 group-hover/nav:text-amber-400 transition-colors" />
                    <span>{item.label}</span>
                    <ChevronDown className="h-3 w-3 transition-transform group-hover/nav:rotate-180" />
                  </button>

                  {/* Dropdown Menu */}
                  <div className="absolute top-full left-1/2 -translate-x-1/2 pt-3 w-80 opacity-0 invisible group-hover/nav:opacity-100 group-hover/nav:visible transition-all duration-250 z-50">
                    <div className="bg-white border border-slate-100 rounded-3xl shadow-2xl p-3.5 text-slate-800 space-y-1">
                      {item.items.map((subItem: any) => {
                        const meta = getPPIDMenuItemMeta(subItem.label);
                        const ItemIcon = meta.icon;
                        return (
                          <button
                            key={subItem.label}
                            onClick={() => onNavClick(subItem.href)}
                            className="w-full flex items-start gap-3 p-2.5 rounded-2xl text-left hover:bg-slate-50 transition-all cursor-pointer group/item border-0 bg-transparent"
                          >
                            <div className={`flex h-9 w-9 shrink-0 items-center justify-center rounded-xl border ${meta.color} transition-colors`}>
                              <ItemIcon className="h-4.5 w-4.5" />
                            </div>
                            <div className="space-y-0.5">
                              <div className="text-[11px] font-extrabold uppercase tracking-wide text-slate-700 group-hover/item:text-[#002147] transition-colors">
                                {subItem.label}
                              </div>
                              <div className="text-[9px] text-slate-400 font-medium leading-normal">
                                {meta.desc}
                              </div>
                            </div>
                          </button>
                        );
                      })}
                    </div>
                  </div>
                </div>
              );
            } else {
              // Direct dynamic link
              if (item.isHighlight) {
                return (
                  <button
                    key={item.label || idx}
                    onClick={() => onNavClick(item.href)}
                    className="px-4 py-2 bg-amber-400 hover:bg-amber-500 text-[#002147] rounded-xl text-xs font-bold uppercase tracking-wider transition-all shadow-md cursor-pointer shrink-0 border-0"
                  >
                    {item.label}
                  </button>
                );
              }

              const IconComponent = getMenuItemIcon(item.href, item.label);
              const isActive = (item.href === '/' && currentPage === 'home') ||
                               (item.href === 'berita' && currentPage.startsWith('berita')) ||
                               (currentPage === `page/${item.href}`);

              return (
                <button
                  key={item.label || idx}
                  onClick={() => {
                    if (item.href === '/' || item.href === 'home') navigateToHome();
                    else if (item.href === 'berita') navigateToNews();
                    else onNavClick(item.href);
                  }}
                  className={`text-xs font-bold uppercase tracking-wider hover:text-amber-400 transition-colors flex items-center gap-1.5 shrink-0 border-0 bg-transparent cursor-pointer ${isActive ? 'text-amber-400' : 'text-slate-200'}`}
                >
                  <IconComponent className={`h-3.5 w-3.5 ${isActive ? 'text-amber-400' : 'text-slate-400'} group-hover:text-amber-400 transition-colors`} />
                  <span>{item.label}</span>
                </button>
              );
            }
          })}

          {adminUser ? (
            <button
              onClick={() => { setCurrentPage('admin'); setAdminActiveTab('dashboard'); }}
              className="p-2 text-slate-200 hover:text-amber-400 border border-slate-700 hover:border-amber-400 rounded-xl transition-all cursor-pointer bg-transparent"
              title="Dashboard Admin"
            >
              <LayoutDashboard className="h-4 w-4" />
            </button>
          ) : (
            <button
              onClick={() => setCurrentPage('admin')}
              className="p-2 text-slate-200 hover:text-amber-400 border border-slate-700 hover:border-amber-400 rounded-xl transition-all cursor-pointer bg-transparent"
              title="Admin"
            >
              <Lock className="h-4 w-4" />
            </button>
          )}
        </nav>

        {/* Mobile Header Quick Actions (Lock/Dashboard icon on the top right) */}
        <div className="lg:hidden flex items-center gap-2">
          {adminUser ? (
            <button
              onClick={() => { setCurrentPage('admin'); setAdminActiveTab('dashboard'); }}
              className="p-2 text-slate-205 hover:text-amber-400 border border-white/10 rounded-xl transition-all cursor-pointer bg-transparent"
              title="Dashboard Admin"
            >
              <LayoutDashboard className="h-4 w-4" />
            </button>
          ) : (
            <button
              onClick={() => setCurrentPage('admin')}
              className="p-2 text-slate-205 hover:text-amber-400 border border-white/10 rounded-xl transition-all cursor-pointer bg-transparent"
              title="Admin"
            >
              <Lock className="h-4 w-4" />
            </button>
          )}
        </div>
      </div>

      {/* ======================================================== */}
      {/* MOBILE BOTTOM NAVIGATION BAR (Android style fixed bar)   */}
      {/* ======================================================== */}
      <div className="lg:hidden fixed bottom-0 left-0 right-0 h-16 bg-[#002147] border-t border-white/10 z-[999] shadow-2xl flex items-center justify-around px-2 pb-safe">
        
        {/* Tab 1: Beranda */}
        <button
          onClick={() => {
            setActiveBottomSheet(null);
            navigateToHome();
          }}
          className={`flex flex-col items-center justify-center gap-1 flex-1 py-1 cursor-pointer bg-transparent border-0 transition-colors ${
            currentPage === 'home' ? 'text-amber-400 font-extrabold' : 'text-slate-400'
          }`}
        >
          <Home className="h-5 w-5" />
          <span className="text-[8.5px] uppercase tracking-wider font-extrabold">Beranda</span>
        </button>

        {/* Tab 2: Layanan (Triggers Bottom Sheet) */}
        <button
          onClick={() => setActiveBottomSheet(activeBottomSheet === 'layanan' ? null : 'layanan')}
          className={`flex flex-col items-center justify-center gap-1 flex-1 py-1 cursor-pointer bg-transparent border-0 transition-colors ${
            activeBottomSheet === 'layanan' || currentPage.startsWith('page/permohonan') || currentPage.startsWith('page/keberatan') || currentPage.startsWith('page/Permohonan')
              ? 'text-amber-400 font-extrabold' : 'text-slate-400'
          }`}
        >
          <MessageSquare className="h-5 w-5" />
          <span className="text-[8.5px] uppercase tracking-wider font-extrabold">Layanan</span>
        </button>

        {/* Tab 3: Berita */}
        <button
          onClick={() => {
            setActiveBottomSheet(null);
            navigateToNews();
          }}
          className={`flex flex-col items-center justify-center gap-1 flex-1 py-1 cursor-pointer bg-transparent border-0 transition-colors ${
            currentPage.startsWith('berita') ? 'text-amber-400 font-extrabold' : 'text-slate-400'
          }`}
        >
          <Newspaper className="h-5 w-5" />
          <span className="text-[8.5px] uppercase tracking-wider font-extrabold">Berita</span>
        </button>

        {/* Tab 4: Regulasi */}
        <button
          onClick={() => {
            setActiveBottomSheet(null);
            onNavClick('regulasi');
          }}
          className={`flex flex-col items-center justify-center gap-1 flex-1 py-1 cursor-pointer bg-transparent border-0 transition-colors ${
            currentPage === 'page/regulasi' ? 'text-amber-400 font-extrabold' : 'text-slate-400'
          }`}
        >
          <FileText className="h-5 w-5" />
          <span className="text-[8.5px] uppercase tracking-wider font-extrabold">Regulasi</span>
        </button>

        {/* Tab 5: Lainnya / Menu (Triggers Menu Bottom Sheet) */}
        <button
          onClick={() => setActiveBottomSheet(activeBottomSheet === 'menu' ? null : 'menu')}
          className={`flex flex-col items-center justify-center gap-1 flex-1 py-1 cursor-pointer bg-transparent border-0 transition-colors ${
            activeBottomSheet === 'menu' ? 'text-amber-400 font-extrabold' : 'text-slate-400'
          }`}
        >
          <Menu className="h-5 w-5" />
          <span className="text-[8.5px] uppercase tracking-wider font-extrabold">Lainnya</span>
        </button>
      </div>

      {/* ======================================================== */}
      {/* MOBILE BOTTOM SHEETS (Slide-up native draw dialogs)     */}
      {/* ======================================================== */}
      
      {/* Backdrop overlay */}
      {activeBottomSheet && (
        <div 
          onClick={() => setActiveBottomSheet(null)} 
          className="fixed inset-0 bg-black/60 backdrop-blur-sm z-[998] lg:hidden animate-in fade-in duration-200" 
        />
      )}

      {/* Bottom Sheet Card */}
      {activeBottomSheet && (
        <div className="fixed bottom-16 left-0 right-0 max-h-[75vh] bg-white text-slate-800 rounded-t-[2rem] p-6 pb-8 z-[999] shadow-2xl overflow-y-auto text-left border-t border-slate-100 flex flex-col gap-4 animate-in slide-in-from-bottom duration-300">
          
          <div className="flex items-center justify-between border-b border-slate-100 pb-3.5">
            <div>
              <span className="text-[9px] text-amber-500 font-extrabold uppercase tracking-widest block">PORTAL UTAMA</span>
              <h3 className="text-sm font-black text-[#002147] uppercase tracking-wide mt-1">
                {activeBottomSheet === 'layanan' ? 'Kategori Layanan PPID' : 'Menu Navigasi Tambahan'}
              </h3>
            </div>
            <button 
              onClick={() => setActiveBottomSheet(null)} 
              className="p-1.5 hover:bg-slate-100 rounded-xl text-slate-400 hover:text-slate-700 cursor-pointer border-0 bg-transparent transition-all"
            >
              <X className="h-5 w-5" />
            </button>
          </div>

          {/* Sheet 1: Kategori Layanan */}
          {activeBottomSheet === 'layanan' && (
            <div className="grid grid-cols-1 gap-2.5">
              {servicesList.map((subItem: any) => {
                const meta = getPPIDMenuItemMeta(subItem.label);
                const ItemIcon = meta.icon;
                return (
                  <button
                    key={subItem.label}
                    onClick={() => onNavClick(subItem.href)}
                    className="w-full flex items-start gap-3.5 p-3 rounded-2xl text-left hover:bg-slate-50 transition-all cursor-pointer border border-slate-150 bg-white group"
                  >
                    <div className={`flex h-9 w-9 shrink-0 items-center justify-center rounded-xl border ${meta.color}`}>
                      <ItemIcon className="h-4.5 w-4.5" />
                    </div>
                    <div className="space-y-0.5">
                      <div className="text-[11px] font-extrabold uppercase tracking-wide text-slate-700 group-hover:text-[#002147]">
                        {subItem.label}
                      </div>
                      <div className="text-[9px] text-slate-400 font-medium leading-normal">
                        {meta.desc}
                      </div>
                    </div>
                  </button>
                );
              })}
            </div>
          )}

          {/* Sheet 2: Menu Navigasi Lainnya */}
          {activeBottomSheet === 'menu' && (
            <div className="space-y-3">
              {otherMenus.map((item: any, idx: number) => {
                const isDropdown = item.items && Array.isArray(item.items) && item.items.length > 0;

                if (isDropdown) {
                  const isAccordionOpen = mobileAccordionActive === item.label;
                  return (
                    <div key={item.label || idx} className="space-y-1">
                      <button
                        onClick={() => setMobileAccordionActive(isAccordionOpen ? null : item.label)}
                        className="w-full flex items-center justify-between font-black text-[11px] uppercase tracking-wide py-3 px-3.5 bg-slate-50 hover:bg-slate-100 rounded-xl border border-slate-150 cursor-pointer text-[#002147]"
                      >
                        <span>{item.label}</span>
                        {isAccordionOpen ? (
                          <ChevronDown className="h-4 w-4 text-amber-500" />
                        ) : (
                          <ChevronRight className="h-4 w-4 text-slate-400" />
                        )}
                      </button>

                      {isAccordionOpen && (
                        <div className="pl-3.5 space-y-2 py-2 bg-slate-50/50 rounded-2xl border border-slate-100 mt-1.5 animate-in fade-in duration-200">
                          {item.items.map((subItem: any) => {
                            const meta = getPPIDMenuItemMeta(subItem.label);
                            const SubIcon = meta.icon;
                            return (
                              <button
                                key={subItem.label}
                                onClick={() => onNavClick(subItem.href)}
                                className="w-full flex items-center gap-2.5 py-2 px-2 hover:bg-white rounded-xl text-left border-0 bg-transparent cursor-pointer text-slate-650"
                              >
                                <SubIcon className="h-4 w-4 text-slate-400" />
                                <span className="text-[10px] font-bold uppercase tracking-wider">{subItem.label}</span>
                              </button>
                            );
                          })}
                        </div>
                      )}
                    </div>
                  );
                } else {
                  const meta = getPPIDMenuItemMeta(item.label);
                  const MenuIcon = meta.icon;
                  return (
                    <button
                      key={item.label || idx}
                      onClick={() => onNavClick(item.href)}
                      className="w-full flex items-center gap-3.5 p-3 rounded-2xl text-left hover:bg-slate-50 transition-all cursor-pointer border border-slate-150 bg-white group"
                    >
                      <div className={`flex h-9 w-9 shrink-0 items-center justify-center rounded-xl border ${meta.color}`}>
                        <MenuIcon className="h-4.5 w-4.5" />
                      </div>
                      <div className="space-y-0.5">
                        <div className="text-[11px] font-extrabold uppercase tracking-wide text-slate-700 group-hover:text-[#002147]">
                          {item.label}
                        </div>
                        <div className="text-[9px] text-slate-400 font-medium leading-normal">
                          {meta.desc}
                        </div>
                      </div>
                    </button>
                  );
                }
              })}

              {/* Portal Admin button */}
              {!adminUser && (
                <button
                  onClick={() => { setActiveBottomSheet(null); setCurrentPage('admin'); }}
                  className="w-full inline-flex items-center justify-center gap-2 px-4 py-3 rounded-2xl bg-[#002147] hover:bg-[#003166] text-white text-xs font-black uppercase tracking-wider cursor-pointer border-0 shadow-sm mt-2"
                >
                  <Lock className="h-4 w-4 text-amber-400" /> Portal Admin
                </button>
              )}
            </div>
          )}
        </div>
      )}

    </header>
  );
}
