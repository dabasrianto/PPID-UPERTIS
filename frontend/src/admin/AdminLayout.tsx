import React, { useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import {
  LayoutDashboard, Layers, Download, Image, FileText, Newspaper, MessageSquare, Settings, LogOut, CheckCircle, X, Landmark, Menu
} from 'lucide-react';
import { resolveImageUrl } from '../utils/helpers';

interface AdminLayoutProps {
  adminUser: { full_name: string; role: string } | null;
  adminActiveTab: string;
  setAdminActiveTab: (tab: string) => void;
  fetchAdminData: () => void;
  setAdminUser: (val: any) => void;
  setCurrentPage: (page: string) => void;
  adminGlobalMessage: string;
  setAdminGlobalMessage: (msg: string) => void;
  siteConfig: any;
  children: React.ReactNode;
}

export default function AdminLayout({
  adminUser,
  adminActiveTab,
  setAdminActiveTab,
  fetchAdminData,
  setAdminUser,
  setCurrentPage,
  adminGlobalMessage,
  setAdminGlobalMessage,
  siteConfig,
  children
}: AdminLayoutProps) {
  const [isMobileOpen, setIsMobileOpen] = useState(false);

  if (!adminUser) return null;

  const menuList = [
    { key: 'dashboard', label: 'Dashboard Utama', icon: Landmark, action: () => { setAdminActiveTab('dashboard'); fetchAdminData(); }, desc: 'Statistik & ringkasan aktivitas', color: 'text-blue-500 bg-blue-50 border border-blue-100' },
    { key: 'slider', label: 'Slider Beranda', icon: Layers, action: () => setAdminActiveTab('slider'), desc: 'Banner promo & gambar utama', color: 'text-indigo-500 bg-indigo-50 border border-indigo-100' },
    { key: 'downloads', label: 'Kelola Unduhan', icon: Download, action: () => setAdminActiveTab('downloads'), desc: 'Unggah file & dokumen resmi', color: 'text-emerald-500 bg-emerald-50 border border-emerald-100' },
    { key: 'gallery', label: 'Kelola Galeri', icon: Image, action: () => setAdminActiveTab('gallery'), desc: 'Foto & dokumentasi kegiatan', color: 'text-amber-500 bg-amber-50 border border-amber-100' },
    { key: 'pages', label: 'Kelola Halaman', icon: FileText, action: () => setAdminActiveTab('pages'), desc: 'Profil, regulasi & sengketa', color: 'text-rose-500 bg-rose-50 border border-rose-100' },
    { key: 'posts', label: 'Kelola Berita', icon: Newspaper, action: () => setAdminActiveTab('posts'), desc: 'Tulis & rilis warta berita', color: 'text-purple-500 bg-purple-50 border border-purple-100' },
    { key: 'permohonan', label: 'Kelola Permohonan', icon: MessageSquare, action: () => setAdminActiveTab('permohonan'), desc: 'Tanggapan pemohon & sengketa', color: 'text-teal-500 bg-teal-50 border border-teal-100' },
    { key: 'menu-manager', label: 'Kelola Menu', icon: Layers, action: () => setAdminActiveTab('menu-manager'), desc: 'Navigasi portal & link menu', color: 'text-sky-500 bg-sky-50 border border-sky-100' },
    { key: 'settings', label: 'Pengaturan Portal', icon: Settings, action: () => setAdminActiveTab('settings'), desc: 'Nama instansi, kontak & medsos', color: 'text-slate-500 bg-slate-50 border border-slate-100' }
  ];

  return (
    <div className="flex-1 flex flex-col lg:flex-row gap-6 w-full px-6 lg:px-10 py-8 mx-auto text-left relative">
      
      {/* Mobile Top Header (Visible only on mobile/tablet) */}
      <div className="lg:hidden w-full bg-white border border-slate-200 rounded-3xl p-4 flex items-center justify-between shadow-sm mb-2">
        <div className="flex items-center gap-3">
          <button
            onClick={() => setIsMobileOpen(true)}
            className="p-2 hover:bg-slate-100 rounded-xl text-[#002147] transition-colors cursor-pointer border-0 bg-transparent"
            title="Buka Menu"
          >
            <Menu className="h-6 w-6" />
          </button>
          <div className="flex items-center gap-2">
            {siteConfig?.logo_url ? (
              <img
                src={resolveImageUrl(siteConfig.logo_url)}
                alt="Logo"
                className="h-8 w-8 object-contain bg-white rounded-lg p-0.5 border border-slate-100"
                onError={(e) => {
                  e.currentTarget.src = '/logo.png';
                }}
              />
            ) : (
              <div className="h-8 w-8 bg-[#002147] text-amber-400 rounded-lg flex items-center justify-center font-extrabold text-xs">
                PPID
              </div>
            )}
            <div>
              <span className="text-[8px] text-amber-500 font-extrabold uppercase tracking-wider block">ADMIN</span>
              <span className="text-xs font-black text-[#002147] leading-none">PPID UPERTIS</span>
            </div>
          </div>
        </div>
        <div className="text-right">
          <span className="text-[10px] text-slate-500 font-bold block leading-none">{adminUser.full_name.split(' ')[0]}</span>
          <span className="text-[8px] text-slate-400 capitalize mt-0.5 block">{adminUser.role}</span>
        </div>
      </div>

      {/* Mobile Drawer Navigation (Slide-over overlay panel) */}
      <AnimatePresence>
        {isMobileOpen && (
          <>
            {/* Backdrop Overlay */}
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              onClick={() => setIsMobileOpen(false)}
              className="fixed inset-0 bg-[#002147]/40 backdrop-blur-sm z-[999] lg:hidden"
            />
            
            {/* Drawer Panel Body */}
            <motion.aside
              initial={{ x: '-100%' }}
              animate={{ x: 0 }}
              exit={{ x: '-100%' }}
              transition={{ type: 'spring', damping: 25, stiffness: 220 }}
              className="fixed inset-y-0 left-0 w-80 bg-white p-5 flex flex-col justify-between gap-6 z-[1000] shadow-2xl lg:hidden overflow-y-auto border-r border-slate-100"
            >
              <div className="space-y-6">
                <div className="flex items-center justify-between border-b border-slate-100 pb-4">
                  <div className="flex items-center gap-3">
                    {siteConfig?.logo_url ? (
                      <img
                        src={resolveImageUrl(siteConfig.logo_url)}
                        alt="Logo"
                        className="h-9 w-9 object-contain bg-white border border-slate-150 rounded-xl p-1"
                        onError={(e) => {
                          e.currentTarget.src = '/logo.png';
                        }}
                      />
                    ) : (
                      <div className="h-9 w-9 bg-[#002147] text-amber-400 rounded-xl flex items-center justify-center font-extrabold text-sm">
                        PPID
                      </div>
                    )}
                    <div>
                      <span className="text-[8px] text-amber-500 font-extrabold uppercase tracking-widest block leading-none">ADMIN PORTAL</span>
                      <h2 className="text-xs font-black text-[#002147] tracking-tight leading-tight mt-1">PPID UPERTIS</h2>
                    </div>
                  </div>
                  <button
                    onClick={() => setIsMobileOpen(false)}
                    className="p-1.5 text-slate-400 hover:text-[#002147] hover:bg-slate-50 rounded-lg cursor-pointer border-0 bg-transparent"
                    title="Tutup"
                  >
                    <X className="h-5 w-5" />
                  </button>
                </div>

                <nav className="flex flex-col gap-1.5">
                  {menuList.map(menu => {
                    const Icon = menu.icon;
                    const isActive = adminActiveTab === menu.key;
                    return (
                      <button
                        key={menu.key}
                        onClick={() => {
                          menu.action();
                          setIsMobileOpen(false);
                        }}
                        className={`w-full flex items-start gap-3 rounded-2xl p-2.5 text-left transition-all cursor-pointer group border-0 ${
                          isActive
                            ? 'bg-[#002147] text-white shadow-md shadow-[#002147]/15'
                            : 'text-slate-700 hover:bg-slate-50 bg-transparent'
                        }`}
                      >
                        <div
                          className={`flex h-8.5 w-8.5 shrink-0 items-center justify-center rounded-xl transition-colors ${
                            isActive
                              ? 'bg-amber-400 text-[#002147] border border-amber-300 font-bold'
                              : `${menu.color}`
                          }`}
                        >
                          <Icon className="h-4 w-4" />
                        </div>
                        <div className="space-y-0.5">
                          <div
                            className={`text-[10.5px] font-extrabold uppercase tracking-wide transition-colors ${
                              isActive ? 'text-white' : 'text-slate-700 group-hover:text-[#002147]'
                            }`}
                          >
                            {menu.label}
                          </div>
                          <div
                            className={`text-[8.5px] font-medium leading-normal transition-colors ${
                              isActive ? 'text-slate-300' : 'text-slate-400'
                            }`}
                          >
                            {menu.desc}
                          </div>
                        </div>
                      </button>
                    );
                  })}
                </nav>
              </div>

              <div className="border-t border-slate-100 pt-4 space-y-3">
                <div className="bg-slate-50 border border-slate-200/60 rounded-2xl p-3 text-left">
                  <span className="text-[8px] text-slate-400 font-extrabold uppercase tracking-widest block">PENGGUNA AKTIF</span>
                  <span className="text-xs text-[#002147] font-black block truncate mt-1 leading-none">{adminUser.full_name}</span>
                  <span className="text-[9px] text-slate-500 font-bold capitalize block mt-1 bg-slate-200/50 px-2 py-0.5 rounded-md inline-block">{adminUser.role}</span>
                </div>

                <button
                  onClick={() => {
                    localStorage.removeItem('auth_token');
                    setAdminUser(null);
                    setCurrentPage('home');
                  }}
                  className="w-full inline-flex items-center justify-center gap-2 px-4 py-2.5 rounded-xl border border-red-200 hover:border-red-500 bg-red-50/30 hover:bg-red-500 hover:text-white text-xs font-bold text-red-600 transition-all cursor-pointer"
                >
                  <LogOut className="h-4 w-4" /> Keluar Sesi
                </button>
              </div>
            </motion.aside>
          </>
        )}
      </AnimatePresence>

      {/* Desktop Sidebar Navigation (Visible only on desktop screens) */}
      <aside className="hidden lg:flex w-full lg:w-72 shrink-0 bg-white border border-slate-200 shadow-sm rounded-[2rem] p-5 flex-col justify-between gap-6">
        <div className="space-y-6">
          <div className="flex items-center gap-3 border-b border-slate-100 pb-4">
            {siteConfig?.logo_url ? (
              <img
                src={resolveImageUrl(siteConfig.logo_url)}
                alt="Logo"
                className="h-10 w-10 object-contain bg-white border border-slate-150 rounded-xl p-1 shadow-sm"
                onError={(e) => {
                  e.currentTarget.src = '/logo.png';
                }}
              />
            ) : (
              <div className="h-10 w-10 bg-[#002147] text-amber-400 rounded-xl flex items-center justify-center font-extrabold text-sm">
                PPID
              </div>
            )}
            <div>
              <span className="text-[9px] text-amber-500 font-extrabold uppercase tracking-widest block leading-none">ADMIN PORTAL</span>
              <h2 className="text-sm font-black text-[#002147] tracking-tight leading-tight mt-1">PPID UPERTIS</h2>
            </div>
          </div>

          <nav className="flex flex-col gap-2">
            {menuList.map(menu => {
              const Icon = menu.icon;
              const isActive = adminActiveTab === menu.key;
              return (
                <button
                  key={menu.key}
                  onClick={menu.action}
                  className={`w-full flex items-start gap-3 rounded-2xl p-2.5 text-left transition-all cursor-pointer group border-0 ${
                    isActive
                      ? 'bg-[#002147] text-white shadow-md shadow-[#002147]/15'
                      : 'text-slate-700 hover:bg-slate-50 bg-transparent'
                  }`}
                >
                  <div
                    className={`flex h-9 w-9 shrink-0 items-center justify-center rounded-xl transition-colors ${
                      isActive
                        ? 'bg-amber-400 text-[#002147] border border-amber-300 font-bold'
                        : `${menu.color}`
                    }`}
                  >
                    <Icon className="h-4.5 w-4.5" />
                  </div>
                  <div className="space-y-0.5">
                    <div
                      className={`text-[11px] font-extrabold uppercase tracking-wide transition-colors ${
                        isActive ? 'text-white' : 'text-slate-700 group-hover:text-[#002147]'
                      }`}
                    >
                      {menu.label}
                    </div>
                    <div
                      className={`text-[9px] font-medium leading-normal transition-colors ${
                        isActive ? 'text-slate-300' : 'text-slate-450'
                      }`}
                    >
                      {menu.desc}
                    </div>
                  </div>
                </button>
              );
            })}
          </nav>
        </div>

        <div className="border-t border-slate-100 pt-4 space-y-3">
          <div className="bg-slate-50 border border-slate-200/60 rounded-2xl p-3.5 text-left">
            <span className="text-[8px] text-slate-400 font-extrabold uppercase tracking-widest block">PENGGUNA AKTIF</span>
            <span className="text-xs text-[#002147] font-black block truncate mt-1 leading-none">{adminUser.full_name}</span>
            <span className="text-[9px] text-slate-500 font-bold capitalize block mt-1 bg-slate-200/50 px-2 py-0.5 rounded-md inline-block">{adminUser.role}</span>
          </div>

          <button
            onClick={() => {
              localStorage.removeItem('auth_token');
              setAdminUser(null);
              setCurrentPage('home');
            }}
            className="w-full inline-flex items-center justify-center gap-2 px-4 py-3 rounded-xl border border-red-200 hover:border-red-500 bg-red-50/30 hover:bg-red-500 hover:text-white text-xs font-bold text-red-650 transition-all cursor-pointer"
          >
            <LogOut className="h-4 w-4" /> Keluar Sesi
          </button>
        </div>
      </aside>

      {/* Admin Workspace Content */}
      <div className="flex-1 min-w-0 space-y-6">
        <AnimatePresence>
          {adminGlobalMessage && (
            <motion.div
              initial={{ opacity: 0, y: -20, scale: 0.95 }}
              animate={{ opacity: 1, y: 0, scale: 1 }}
              exit={{ opacity: 0, y: -20, scale: 0.95 }}
              transition={{ duration: 0.25, ease: 'easeOut' }}
              className="fixed top-6 right-6 z-[9999] bg-[#002147] text-white border border-white/10 shadow-2xl rounded-2xl px-5 py-4 max-w-sm flex items-center justify-between gap-4 backdrop-blur-md"
            >
              <div className="flex items-center gap-3 text-left">
                <div className="h-8 w-8 bg-amber-400/20 text-amber-400 rounded-xl flex items-center justify-center shrink-0">
                  <CheckCircle className="h-5 w-5" />
                </div>
                <div>
                  <p className="text-xs font-extrabold tracking-wide text-white leading-tight">Sistem PPID</p>
                  <span className="text-[10px] text-slate-300 font-medium block mt-0.5 leading-normal">{adminGlobalMessage}</span>
                </div>
              </div>
              <button
                onClick={() => setAdminGlobalMessage('')}
                className="text-slate-400 hover:text-white transition-colors p-1 rounded-lg hover:bg-white/10 shrink-0 cursor-pointer"
                title="Tutup"
              >
                <X className="h-4 w-4" />
              </button>
            </motion.div>
          )}
        </AnimatePresence>

        {children}
      </div>
    </div>
  );
}
