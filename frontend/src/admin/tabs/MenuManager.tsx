import React, { useState } from 'react';
import {
  Plus, Trash2, ArrowUp, ArrowDown, Edit2, Check, X, FileText, Globe, RefreshCw, FolderPlus, Link
} from 'lucide-react';

interface MenuItem {
  label: string;
  href: string;
  isExternal?: boolean;
}

interface MenuConfigItem {
  // Direct links properties
  label?: string;
  href?: string;
  isExternal?: boolean;
  isHighlight?: boolean;
  type?: 'link';

  // Dropdown group properties
  group?: string;
  items?: MenuItem[];
}

interface MenuManagerProps {
  siteConfig: any;
  setSiteConfig: (config: any) => void;
  adminPages: any[];
  API_BASE_URL: string;
  setAdminGlobalMessage: (msg: string) => void;
  fetchAdminData: () => void;
}

const migrateMenuStructure = (menu: any[]): any[] => {
  if (!menu || !Array.isArray(menu) || menu.length === 0) {
    return [];
  }

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

  if (!hasHome) {
    migrated.push({ label: 'Home', href: '/', type: 'link' });
  }

  menu.forEach(item => {
    if (item.group && !item.label) {
      migrated.push({
        group: item.group,
        items: item.items || []
      });
    } else {
      migrated.push(item);
    }
  });

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

export default function MenuManager({
  siteConfig,
  setSiteConfig,
  adminPages,
  API_BASE_URL,
  setAdminGlobalMessage,
  fetchAdminData
}: MenuManagerProps) {
  const defaultMenuStructure: MenuConfigItem[] = [
    { label: 'Home', href: '/', type: 'link' },
    {
      group: 'Tentang',
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
      group: 'Info Publik',
      items: [
        { label: 'Info Publik Berkala', href: 'informasi-publik-berkala' },
        { label: 'Info Tersedia Setiap Saat', href: 'informasi-tersedia-setiap-saat' },
        { label: 'Info Serta Merta', href: 'info-serta-merta' },
        { label: 'Zona Integrasi', href: 'zona-integrasi' }
      ]
    },
    {
      group: 'Layanan',
      items: [
        { label: 'Jadwal Layanan', href: 'jadwal-layanan-informasi' },
        { label: 'Permohonan Informasi', href: 'permohonan-informasi' },
        { label: 'Pengajuan Keberatan', href: 'keberatan-informasi' },
        { label: 'Pengaduan Layanan', href: 'https://lapor.go.id/', isExternal: true },
        { label: 'Informasi Dikecualikan', href: 'informasi-dikecualikan' }
      ]
    },
    { label: 'Berita PPID', href: 'berita', type: 'link' },
    { label: 'Galeri', href: 'galeri', type: 'link' },
    { label: 'Download', href: 'download', type: 'link' },
    { label: 'Hubungi Kami', href: 'kontak', type: 'link', isHighlight: true }
  ];

  const getMigratedOrInitial = () => {
    if (siteConfig?.menu && Array.isArray(siteConfig.menu) && siteConfig.menu.length > 0) {
      return migrateMenuStructure(siteConfig.menu);
    }
    return defaultMenuStructure;
  };

  const [menuList, setMenuList] = useState<MenuConfigItem[]>(getMigratedOrInitial());
  const [isSaving, setIsSaving] = useState(false);

  // Sync menuList when siteConfig changes (e.g. from reload or initial load)
  React.useEffect(() => {
    if (siteConfig?.menu && Array.isArray(siteConfig.menu) && siteConfig.menu.length > 0) {
      setMenuList(migrateMenuStructure(siteConfig.menu));
    }
  }, [siteConfig]);

  // States for expanded/editing panels
  const [editingTopLevelIdx, setEditingTopLevelIdx] = useState<number | null>(null);
  const [editingSubIndex, setEditingSubIndex] = useState<{ groupIdx: number; itemIdx: number } | null>(null);

  // Adding new elements states
  const [isAddingGroup, setIsAddingGroup] = useState(false);
  const [newGroupName, setNewGroupName] = useState('');

  const [isAddingLink, setIsAddingLink] = useState(false);
  const [newLinkLabel, setNewLinkLabel] = useState('');
  const [newLinkHref, setNewLinkHref] = useState('');
  const [newLinkIsExternal, setNewLinkIsExternal] = useState(false);
  const [newLinkIsHighlight, setNewLinkIsHighlight] = useState(false);

  // Adding sub-item inside a dropdown group state
  const [newItemForGroup, setNewItemForGroup] = useState<number | null>(null);
  const [newSubLabel, setNewSubLabel] = useState('');
  const [newSubHref, setNewSubHref] = useState('');
  const [newSubIsExternal, setNewSubIsExternal] = useState(false);

  // Direct state handlers for inline editing
  const updateTopLevelLink = (idx: number, fields: Partial<MenuConfigItem>) => {
    const temp = [...menuList];
    temp[idx] = {
      ...temp[idx],
      ...fields
    };
    setMenuList(temp);
  };

  const updateTopLevelGroup = (idx: number, groupName: string) => {
    const temp = [...menuList];
    temp[idx] = {
      ...temp[idx],
      group: groupName
    };
    setMenuList(temp);
  };

  const updateSubItem = (groupIdx: number, itemIdx: number, fields: Partial<MenuItem>) => {
    const temp = [...menuList];
    const group = temp[groupIdx];
    if (group.items) {
      group.items[itemIdx] = {
        ...group.items[itemIdx],
        ...fields
      };
    }
    setMenuList(temp);
  };

  const handleSaveMenuStructure = async () => {
    const token = localStorage.getItem('auth_token');
    if (!token) return;
    setIsSaving(true);

    const siteId = siteConfig?.id || 'ppid';
    const payload = {
      ...siteConfig,
      menu: menuList
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
        setAdminGlobalMessage('Seluruh struktur menu navigasi dinamis berhasil diperbarui!');
        fetchAdminData();
        
        // Refresh site config in global App state
        const host = window.location.hostname;
        const refreshRes = await fetch(`${API_BASE_URL}/site-config?host=${host}`);
        if (refreshRes.ok) {
          const freshData = await refreshRes.json();
          if (freshData && !freshData.is_main) {
            setSiteConfig(freshData);
          } else {
            // Fallback for localhost
            const fallbackRes = await fetch(`${API_BASE_URL}/site-config?host=ppid.localhost`);
            if (fallbackRes.ok) {
              const fallbackData = await fallbackRes.json();
              if (fallbackData && !fallbackData.is_main) {
                setSiteConfig(fallbackData);
              }
            }
          }
        }
      } else {
        alert('Gagal menyimpan menu. Periksa log backend.');
      }
    } catch (err) {
      console.error(err);
      alert('Terjadi kesalahan koneksi.');
    } finally {
      setIsSaving(false);
    }
  };

  const addTopLevelGroup = () => {
    if (!newGroupName.trim()) return;
    setMenuList([...menuList, { group: newGroupName.trim(), items: [] }]);
    setNewGroupName('');
    setIsAddingGroup(false);
  };

  const addTopLevelLink = () => {
    if (!newLinkLabel.trim() || !newLinkHref.trim()) {
      alert('Label menu dan link/slug wajib diisi!');
      return;
    }
    setMenuList([...menuList, {
      label: newLinkLabel.trim(),
      href: newLinkHref.trim(),
      type: 'link',
      isExternal: newLinkIsExternal,
      isHighlight: newLinkIsHighlight
    }]);
    setNewLinkLabel('');
    setNewLinkHref('');
    setNewLinkIsExternal(false);
    setNewLinkIsHighlight(false);
    setIsAddingLink(false);
  };

  const deleteTopLevelElement = (idx: number) => {
    if (!confirm('Apakah Anda yakin ingin menghapus menu ini?')) return;
    const temp = [...menuList];
    temp.splice(idx, 1);
    setMenuList(temp);
    if (editingTopLevelIdx === idx) setEditingTopLevelIdx(null);
  };

  const moveTopLevelElement = (idx: number, direction: 'up' | 'down') => {
    if (direction === 'up' && idx === 0) return;
    if (direction === 'down' && idx === menuList.length - 1) return;

    const temp = [...menuList];
    const targetIdx = direction === 'up' ? idx - 1 : idx + 1;
    const hold = temp[idx];
    temp[idx] = temp[targetIdx];
    temp[targetIdx] = hold;
    setMenuList(temp);
    
    // Adjust edit index if active
    if (editingTopLevelIdx === idx) setEditingTopLevelIdx(targetIdx);
    else if (editingTopLevelIdx === targetIdx) setEditingTopLevelIdx(idx);
  };

  const addSubItemToGroup = (groupIdx: number) => {
    if (!newSubLabel.trim() || !newSubHref.trim()) {
      alert('Label submenu dan link/slug wajib diisi!');
      return;
    }
    const temp = [...menuList];
    const group = temp[groupIdx];
    if (group.items) {
      group.items.push({
        label: newSubLabel.trim(),
        href: newSubHref.trim(),
        isExternal: newSubIsExternal || newSubHref.startsWith('http')
      });
    }
    setMenuList(temp);
    setNewSubLabel('');
    setNewSubHref('');
    setNewSubIsExternal(false);
    setNewItemForGroup(null);
  };

  const deleteSubItem = (groupIdx: number, itemIdx: number) => {
    const temp = [...menuList];
    const group = temp[groupIdx];
    if (group.items) {
      group.items.splice(itemIdx, 1);
    }
    setMenuList(temp);
    if (editingSubIndex?.groupIdx === groupIdx && editingSubIndex?.itemIdx === itemIdx) {
      setEditingSubIndex(null);
    }
  };

  const moveSubItem = (groupIdx: number, itemIdx: number, direction: 'up' | 'down') => {
    const group = menuList[groupIdx];
    if (!group.items) return;
    if (direction === 'up' && itemIdx === 0) return;
    if (direction === 'down' && itemIdx === group.items.length - 1) return;

    const temp = [...menuList];
    const items = temp[groupIdx].items || [];
    const targetIdx = direction === 'up' ? itemIdx - 1 : itemIdx + 1;
    const hold = items[itemIdx];
    items[itemIdx] = items[targetIdx];
    items[targetIdx] = hold;
    setMenuList(temp);
  };

  const handleResetMenu = () => {
    if (confirm('Kembalikan menu ke pengaturan bawaan awal? (Perubahan belum disimpan ke database sebelum Anda mengklik Simpan Struktur Menu).')) {
      setMenuList(JSON.parse(JSON.stringify(defaultMenuStructure)));
      setEditingTopLevelIdx(null);
      setEditingSubIndex(null);
    }
  };

  return (
    <div className="space-y-6 text-left pb-16">
      {/* Header bar */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-slate-100 pb-5">
        <div>
          <h1 className="text-xl font-extrabold text-[#002147]">Pengaturan Menu Navigasi Dinamis</h1>
          <p className="text-xs text-slate-400 font-medium block mt-1">
            Menu statis kini sudah dinamis dan fleksibel. Anda dapat langsung mengedit isi menu dari panel di bawah ini tanpa menghapus apapun.
          </p>
        </div>
        <div className="flex gap-2 flex-wrap">
          <button
            onClick={handleResetMenu}
            className="px-4 py-2 border border-slate-200 text-slate-600 rounded-xl text-xs font-bold bg-white hover:bg-slate-50 transition-colors flex items-center gap-1.5 cursor-pointer shadow-sm"
          >
            <RefreshCw className="h-3.5 w-3.5" /> Reset Default
          </button>
          <button
            onClick={() => { setIsAddingLink(true); setIsAddingGroup(false); }}
            className="px-4 py-2 bg-slate-100 text-[#002147] border border-slate-200 rounded-xl text-xs font-bold hover:bg-slate-200 transition-colors flex items-center gap-1.5 cursor-pointer"
          >
            <Plus className="h-3.5 w-3.5" /> Tambah Link Menu
          </button>
          <button
            onClick={() => { setIsAddingGroup(true); setIsAddingLink(false); }}
            className="px-4 py-2 bg-slate-100 text-[#002147] border border-slate-200 rounded-xl text-xs font-bold hover:bg-slate-200 transition-colors flex items-center gap-1.5 cursor-pointer"
          >
            <FolderPlus className="h-3.5 w-3.5" /> Tambah Dropdown Grup
          </button>
          <button
            onClick={handleSaveMenuStructure}
            disabled={isSaving}
            className="px-5 py-2 bg-amber-400 text-[#002147] rounded-xl text-xs font-bold hover:bg-amber-500 transition-all flex items-center gap-1.5 shadow-sm disabled:opacity-50 cursor-pointer"
          >
            {isSaving ? 'Menyimpan...' : 'Simpan Struktur Menu'}
          </button>
        </div>
      </div>

      {/* Add Group Modal/Form */}
      {isAddingGroup && (
        <div className="bg-slate-50 border border-slate-200 rounded-3xl p-5 max-w-md animate-in slide-in-from-top duration-200">
          <div className="flex justify-between items-center pb-2 mb-3 border-b border-slate-200/50">
            <h4 className="text-xs font-bold text-[#002147] uppercase tracking-wider">Tambah Grup Dropdown Baru</h4>
            <button onClick={() => setIsAddingGroup(false)} className="text-slate-400 hover:text-slate-600">
              <X className="h-4 w-4" />
            </button>
          </div>
          <div className="flex gap-2">
            <input
              type="text"
              value={newGroupName}
              onChange={(e) => setNewGroupName(e.target.value)}
              placeholder="Nama Grup (e.g. Profil Kami)"
              className="flex-1 rounded-xl border border-slate-200 px-3.5 py-2 text-xs bg-white text-slate-800 focus:outline-none"
            />
            <button
              onClick={addTopLevelGroup}
              className="px-4 py-2 bg-[#002147] text-white rounded-xl text-xs font-bold hover:bg-[#00346c]"
            >
              Simpan
            </button>
          </div>
        </div>
      )}

      {/* Add Link Form */}
      {isAddingLink && (
        <div className="bg-slate-50 border border-slate-200 rounded-3xl p-5 max-w-md animate-in slide-in-from-top duration-200 space-y-3">
          <div className="flex justify-between items-center pb-2 border-b border-slate-200/50">
            <h4 className="text-xs font-bold text-[#002147] uppercase tracking-wider">Tambah Link Menu Baru</h4>
            <button onClick={() => setIsAddingLink(false)} className="text-slate-400 hover:text-slate-600">
              <X className="h-4 w-4" />
            </button>
          </div>
          <div className="space-y-3">
            <div className="space-y-1">
              <label className="text-[10px] font-bold text-slate-500 uppercase block">Label Link</label>
              <input
                type="text"
                placeholder="e.g. Berita KIP"
                value={newLinkLabel}
                onChange={(e) => setNewLinkLabel(e.target.value)}
                className="w-full rounded-xl border border-slate-200 px-3.5 py-2 text-xs bg-white text-slate-800"
              />
            </div>
            <div className="space-y-1">
              <label className="text-[10px] font-bold text-slate-500 uppercase block">Link / Slug</label>
              <div className="flex gap-1">
                <input
                  type="text"
                  placeholder="e.g. berita atau https://..."
                  value={newLinkHref}
                  onChange={(e) => setNewLinkHref(e.target.value)}
                  className="flex-1 rounded-xl border border-slate-200 px-3.5 py-2 text-xs bg-white text-slate-800"
                />
                <select
                  onChange={(e) => setNewLinkHref(e.target.value)}
                  className="rounded-xl border border-slate-200 text-[10px] bg-white text-slate-500 w-24"
                  defaultValue=""
                >
                  <option value="" disabled>Pilih Page</option>
                  <option value="/">Home</option>
                  <option value="berita">Berita</option>
                  {adminPages.map(p => (
                    <option key={p.slug} value={p.slug}>{p.title}</option>
                  ))}
                </select>
              </div>
            </div>
            <div className="flex items-center justify-between pt-1">
              <div className="flex items-center gap-1.5">
                <input
                  type="checkbox"
                  id="new-link-ext"
                  checked={newLinkIsExternal}
                  onChange={(e) => setNewLinkIsExternal(e.target.checked)}
                  className="rounded border-slate-300 text-blue-600 h-3.5 w-3.5"
                />
                <label htmlFor="new-link-ext" className="text-[10px] text-slate-500 font-semibold">
                  Buka di tab baru (External)
                </label>
              </div>
              <div className="flex items-center gap-1.5">
                <input
                  type="checkbox"
                  id="new-link-highlight"
                  checked={newLinkIsHighlight}
                  onChange={(e) => setNewLinkIsHighlight(e.target.checked)}
                  className="rounded border-slate-300 text-blue-600 h-3.5 w-3.5"
                />
                <label htmlFor="new-link-highlight" className="text-[10px] text-slate-500 font-semibold">
                  Jadikan Tombol Highlight
                </label>
              </div>
            </div>
            <button
              onClick={addTopLevelLink}
              className="w-full py-2 bg-[#002147] text-white rounded-xl text-xs font-bold hover:bg-[#00346c] mt-2"
            >
              Simpan Menu Link
            </button>
          </div>
        </div>
      )}

      {/* Main Drag-Sort/Layout Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 items-start">
        {menuList.map((item, idx) => {
          const isDropdown = !!item.group;
          const isEditing = editingTopLevelIdx === idx;

          return (
            <div
              key={idx}
              className="bg-white rounded-3xl border border-slate-200 shadow-sm overflow-hidden flex flex-col min-h-36"
            >
              {/* Group / Link Header */}
              <div className="bg-slate-50 border-b border-slate-100 p-4 flex items-center justify-between gap-3">
                <div className="flex-1 flex items-center gap-2 min-w-0">
                  {isDropdown ? (
                    <FolderPlus className="h-4 w-4 text-[#002147]/70 shrink-0" />
                  ) : (
                    <Link className="h-4 w-4 text-amber-500 shrink-0" />
                  )}

                  {isEditing && isDropdown ? (
                    <input
                      type="text"
                      value={item.group || ''}
                      onChange={(e) => updateTopLevelGroup(idx, e.target.value)}
                      className="rounded-lg border border-slate-300 px-2 py-0.5 text-xs text-slate-800 font-bold focus:outline-none w-32"
                    />
                  ) : (
                    <div className="min-w-0 text-left">
                      <span
                        onClick={() => setEditingTopLevelIdx(isEditing ? null : idx)}
                        className="text-xs font-extrabold text-[#002147] hover:underline cursor-pointer truncate block"
                      >
                        {isDropdown ? item.group : item.label}
                      </span>
                      {!isDropdown && (
                        <span className="text-[9px] text-slate-400 font-mono block truncate mt-0.5">
                          Link: {item.href} {item.isHighlight && '• Highlight'}
                        </span>
                      )}
                    </div>
                  )}
                </div>

                {/* Header Actions */}
                <div className="flex items-center gap-1 shrink-0">
                  <button
                    onClick={() => moveTopLevelElement(idx, 'up')}
                    disabled={idx === 0}
                    className="p-1 hover:bg-slate-200 rounded text-slate-400 hover:text-slate-700 disabled:opacity-40"
                    title="Geser Kiri"
                  >
                    <ArrowUp className="h-3.5 w-3.5 rotate-270" />
                  </button>
                  <button
                    onClick={() => moveTopLevelElement(idx, 'down')}
                    disabled={idx === menuList.length - 1}
                    className="p-1 hover:bg-slate-200 rounded text-slate-400 hover:text-slate-700 disabled:opacity-40"
                    title="Geser Kanan"
                  >
                    <ArrowDown className="h-3.5 w-3.5 rotate-270" />
                  </button>
                  <button
                    onClick={() => setEditingTopLevelIdx(isEditing ? null : idx)}
                    className={`p-1 rounded transition-colors ${isEditing ? 'bg-amber-100 text-amber-800' : 'hover:bg-slate-200 text-slate-400 hover:text-slate-700'}`}
                    title="Edit Menu"
                  >
                    {isEditing && isDropdown ? <Check className="h-3.5 w-3.5 text-green-600" /> : <Edit2 className="h-3.5 w-3.5" />}
                  </button>
                  <button
                    onClick={() => deleteTopLevelElement(idx)}
                    className="p-1 hover:bg-red-50 rounded text-slate-400 hover:text-red-600"
                    title="Hapus Menu"
                  >
                    <Trash2 className="h-3.5 w-3.5" />
                  </button>
                </div>
              </div>

              {/* Submenus items if group */}
              {isDropdown && (
                <div className="p-4 space-y-2 flex-1">
                  {(item.items || []).length === 0 ? (
                    <p className="text-[10px] text-slate-400 font-medium py-4 text-center">Belum ada submenu link.</p>
                  ) : (
                    (item.items || []).map((subItem, itemIdx) => {
                      const isSubEditing = editingSubIndex?.groupIdx === idx && editingSubIndex?.itemIdx === itemIdx;

                      if (isSubEditing) {
                        return (
                          <div key={itemIdx} className="bg-slate-50 border border-slate-200 rounded-xl p-3 space-y-2 text-left animate-in fade-in duration-200">
                            <div className="space-y-1">
                              <label className="text-[10px] font-bold text-slate-500 uppercase block">Label Menu</label>
                              <input
                                type="text"
                                value={subItem.label || ''}
                                onChange={(e) => updateSubItem(idx, itemIdx, { label: e.target.value })}
                                className="w-full rounded-lg border border-slate-200 px-2.5 py-1.5 text-xs text-slate-800 bg-white"
                              />
                            </div>
                            <div className="space-y-1">
                              <label className="text-[10px] font-bold text-slate-500 uppercase block">Link / Slug</label>
                              <div className="flex gap-1">
                                <input
                                  type="text"
                                  value={subItem.href || ''}
                                  onChange={(e) => updateSubItem(idx, itemIdx, { href: e.target.value })}
                                  className="flex-1 rounded-lg border border-slate-200 px-2.5 py-1.5 text-xs text-slate-800 bg-white"
                                />
                                <select
                                  onChange={(e) => updateSubItem(idx, itemIdx, { href: e.target.value })}
                                  className="rounded-lg border border-slate-200 text-[10px] bg-white text-slate-500 w-20"
                                  value={subItem.href || ''}
                                >
                                  <option value="" disabled>Pilih Page</option>
                                  {adminPages.map(p => (
                                    <option key={p.slug} value={p.slug}>{p.title}</option>
                                  ))}
                                </select>
                              </div>
                            </div>
                            <div className="flex items-center gap-1.5 pt-1">
                              <input
                                type="checkbox"
                                id={`edit-sub-ext-${idx}-${itemIdx}`}
                                checked={!!subItem.isExternal}
                                onChange={(e) => updateSubItem(idx, itemIdx, { isExternal: e.target.checked })}
                                className="rounded border-slate-300 text-blue-600 h-3 w-3"
                              />
                              <label htmlFor={`edit-sub-ext-${idx}-${itemIdx}`} className="text-[10px] text-slate-500 font-semibold">
                                Buka di tab baru (External)
                              </label>
                            </div>
                            <div className="pt-1">
                              <button
                                onClick={() => setEditingSubIndex(null)}
                                className="w-full py-1 bg-white border border-slate-200 rounded-lg text-xs text-slate-600 font-bold hover:bg-slate-100 flex items-center justify-center gap-1 cursor-pointer"
                              >
                                <Check className="h-3 w-3 text-green-600" /> Selesai
                              </button>
                            </div>
                          </div>
                        );
                      }

                      return (
                        <div
                          key={itemIdx}
                          className="group/item flex items-center justify-between border border-slate-100 rounded-xl px-3 py-2 bg-slate-50/50 hover:bg-slate-50 hover:border-slate-200 transition-all text-xs"
                        >
                          <div className="flex items-center gap-2 min-w-0 text-left">
                            {subItem.href.startsWith('http') || subItem.isExternal ? (
                              <Globe className="h-3.5 w-3.5 text-amber-500 shrink-0" />
                            ) : (
                              <FileText className="h-3.5 w-3.5 text-blue-500 shrink-0" />
                            )}
                            <div className="min-w-0">
                              <span className="font-bold text-slate-700 block truncate">{subItem.label}</span>
                              <span className="text-[9px] text-slate-400 block font-mono truncate">{subItem.href}</span>
                            </div>
                          </div>

                          {/* Sub Item Actions */}
                          <div className="flex items-center gap-1 shrink-0 opacity-0 group-hover/item:opacity-100 transition-opacity">
                            <button
                              onClick={() => moveSubItem(idx, itemIdx, 'up')}
                              disabled={itemIdx === 0}
                              className="p-1 hover:bg-slate-200 rounded text-slate-400 hover:text-slate-700 disabled:opacity-40"
                              title="Geser Atas"
                            >
                              <ArrowUp className="h-3 w-3" />
                            </button>
                            <button
                              onClick={() => moveSubItem(idx, itemIdx, 'down')}
                              disabled={itemIdx === (item.items || []).length - 1}
                              className="p-1 hover:bg-slate-200 rounded text-slate-400 hover:text-slate-700 disabled:opacity-40"
                              title="Geser Bawah"
                            >
                              <ArrowDown className="h-3 w-3" />
                            </button>
                            <button
                              onClick={() => setEditingSubIndex({ groupIdx: idx, itemIdx })}
                              className="p-1 hover:bg-slate-200 rounded text-slate-400 hover:text-slate-700"
                              title="Ubah Submenu"
                            >
                              <Edit2 className="h-3 w-3" />
                            </button>
                            <button
                              onClick={() => deleteSubItem(idx, itemIdx)}
                              className="p-1 hover:bg-red-50 rounded text-slate-400 hover:text-red-600"
                              title="Hapus Submenu"
                            >
                              <Trash2 className="h-3 w-3" />
                            </button>
                          </div>
                        </div>
                      );
                    })
                  )}
                </div>
              )}

              {/* Bottom Add Menu Form for this group */}
              {isDropdown && (
                <div className="border-t border-slate-100 p-4 bg-slate-50/50">
                  {newItemForGroup === idx ? (
                    <div className="space-y-2 text-left bg-white border border-slate-200 rounded-2xl p-3 animate-in fade-in duration-150">
                      <div className="space-y-1">
                        <label className="text-[9px] font-bold text-slate-400 uppercase tracking-wider block">Label Submenu</label>
                        <input
                          type="text"
                          placeholder="e.g. Buku Panduan"
                          value={newSubLabel}
                          onChange={(e) => setNewSubLabel(e.target.value)}
                          className="w-full rounded-lg border border-slate-200 px-2.5 py-1.5 text-xs text-slate-800"
                        />
                      </div>
                      <div className="space-y-1">
                        <label className="text-[9px] font-bold text-slate-400 uppercase tracking-wider block">Link / Slug</label>
                        <div className="flex gap-1">
                          <input
                            type="text"
                            placeholder="e.g. regulasi atau link luar"
                            value={newSubHref}
                            onChange={(e) => setNewSubHref(e.target.value)}
                            className="flex-1 rounded-lg border border-slate-200 px-2.5 py-1.5 text-xs text-slate-800"
                          />
                          <select
                            onChange={(e) => setNewSubHref(e.target.value)}
                            className="rounded-lg border border-slate-200 text-[10px] bg-slate-50 text-slate-500 w-20"
                            defaultValue=""
                          >
                            <option value="" disabled>Pilih Page</option>
                            {adminPages.map(p => (
                              <option key={p.slug} value={p.slug}>{p.title}</option>
                            ))}
                          </select>
                        </div>
                      </div>
                      <div className="flex items-center gap-1.5 pt-1">
                        <input
                          type="checkbox"
                          id={`new-sub-ext-${idx}`}
                          checked={newSubIsExternal}
                          onChange={(e) => setNewSubIsExternal(e.target.checked)}
                          className="rounded border-slate-300 text-blue-600 h-3 w-3"
                        />
                        <label htmlFor={`new-sub-ext-${idx}`} className="text-[9px] text-slate-400 font-semibold">
                          Buka di tab baru (External)
                        </label>
                      </div>
                      <div className="flex gap-1.5 pt-2">
                        <button
                          onClick={() => addSubItemToGroup(idx)}
                          className="flex-1 bg-[#002147] text-white py-1 rounded-lg text-xs font-bold hover:bg-[#00346c]"
                        >
                          Tambahkan
                        </button>
                        <button
                          onClick={() => setNewItemForGroup(null)}
                          className="px-3 border border-slate-200 rounded-lg text-xs text-slate-500 font-bold hover:bg-slate-100"
                        >
                          Batal
                        </button>
                      </div>
                    </div>
                  ) : (
                    <button
                      onClick={() => {
                        setNewItemForGroup(idx);
                        setNewSubLabel('');
                        setNewSubHref('');
                        setNewSubIsExternal(false);
                      }}
                      className="w-full py-2 bg-white border border-dashed border-slate-200 rounded-xl text-slate-500 hover:text-slate-700 hover:border-slate-300 transition-colors text-xs font-bold flex items-center justify-center gap-1 cursor-pointer"
                    >
                      <Plus className="h-3.5 w-3.5" /> Tambah Submenu
                    </button>
                  )}
                </div>
              )}

              {/* Inline Link Edit Form inside Card Body */}
              {!isDropdown && isEditing && (
                <div className="p-4 space-y-3 flex-1 text-left bg-slate-50/50 border-t border-slate-100">
                  <div className="space-y-1">
                    <label className="text-[10px] font-bold text-slate-500 uppercase block">Label Link</label>
                    <input
                      type="text"
                      value={item.label || ''}
                      onChange={(e) => updateTopLevelLink(idx, { label: e.target.value })}
                      className="w-full rounded-lg border border-slate-200 px-2.5 py-1.5 text-xs text-slate-800 bg-white"
                    />
                  </div>
                  <div className="space-y-1">
                    <label className="text-[10px] font-bold text-slate-500 uppercase block">Link / Slug</label>
                    <div className="flex gap-1">
                      <input
                        type="text"
                        value={item.href || ''}
                        onChange={(e) => updateTopLevelLink(idx, { href: e.target.value })}
                        className="flex-1 rounded-lg border border-slate-200 px-2.5 py-1.5 text-xs text-slate-800 bg-white"
                      />
                      <select
                        onChange={(e) => updateTopLevelLink(idx, { href: e.target.value })}
                        className="rounded-lg border border-slate-200 text-[10px] bg-white text-slate-500 w-20"
                        value={item.href || ''}
                      >
                        <option value="" disabled>Pilih Page</option>
                        <option value="/">Home</option>
                        <option value="berita">Berita</option>
                        {adminPages.map(p => (
                          <option key={p.slug} value={p.slug}>{p.title}</option>
                        ))}
                      </select>
                    </div>
                  </div>
                  <div className="flex items-center justify-between pt-1">
                    <div className="flex items-center gap-1.5">
                      <input
                        type="checkbox"
                        id={`edit-top-ext-${idx}`}
                        checked={!!item.isExternal}
                        onChange={(e) => updateTopLevelLink(idx, { isExternal: e.target.checked })}
                        className="rounded border-slate-300 text-blue-600 h-3.5 w-3.5"
                      />
                      <label htmlFor={`edit-top-ext-${idx}`} className="text-[10px] text-slate-500 font-semibold">
                        Buka di tab baru
                      </label>
                    </div>
                    <div className="flex items-center gap-1.5">
                      <input
                        type="checkbox"
                        id={`edit-top-highlight-${idx}`}
                        checked={!!item.isHighlight}
                        onChange={(e) => updateTopLevelLink(idx, { isHighlight: e.target.checked })}
                        className="rounded border-slate-300 text-blue-600 h-3.5 w-3.5"
                      />
                      <label htmlFor={`edit-top-highlight-${idx}`} className="text-[10px] text-slate-500 font-semibold">
                        Highlight Button
                      </label>
                    </div>
                  </div>
                  <div className="pt-2">
                    <button
                      onClick={() => setEditingTopLevelIdx(null)}
                      className="w-full py-1.5 border border-slate-200 rounded-lg text-xs text-slate-600 font-bold hover:bg-slate-100 flex items-center justify-center gap-1 bg-white cursor-pointer"
                    >
                      <Check className="h-3.5 w-3.5 text-green-600" /> Selesai Mengubah
                    </button>
                  </div>
                </div>
              )}
            </div>
          );
        })}
      </div>
    </div>
  );
}
