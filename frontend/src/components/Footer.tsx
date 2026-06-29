import { Heart } from 'lucide-react';

interface FooterProps {
  navigateToHome: () => void;
  navigateToPage: (slug: string) => void;
  setCurrentPage: (page: string) => void;
}

export default function Footer({
  navigateToHome,
  navigateToPage,
  setCurrentPage
}: FooterProps) {
  return (
    <footer className="bg-[#001733] text-slate-400 py-8 px-4 border-t border-white/5 text-xs">
      <div className="container mx-auto max-w-6xl flex flex-col md:flex-row items-center justify-between gap-4">
        <div className="text-center md:text-left space-y-1">
          <span className="text-slate-200 block font-bold">PPID Universitas Perintis Indonesia (UPERTIS)</span>
          <p className="text-[11px] font-medium text-slate-500">&copy; {new Date().getFullYear()} Universitas Perintis Indonesia. All Rights Reserved.</p>
        </div>
        <div className="flex items-center gap-4 font-bold text-[11px] uppercase tracking-wide">
          <button onClick={navigateToHome} className="hover:text-white transition-colors cursor-pointer">Beranda</button>
          <span>&bull;</span>
          <button onClick={() => navigateToPage('profil')} className="hover:text-white transition-colors cursor-pointer">Profil PPID</button>
          <span>&bull;</span>
          <button onClick={() => navigateToPage('regulasi')} className="hover:text-white transition-colors cursor-pointer">Regulasi</button>
          <span>&bull;</span>
          <button onClick={() => setCurrentPage('admin')} className="hover:text-white transition-colors cursor-pointer">Portal Admin</button>
          <span>&bull;</span>
          <a
            href="https://akhimedia.id"
            target="_blank"
            rel="noopener noreferrer"
            className="text-slate-650 hover:text-red-500 hover:scale-110 active:scale-95 transition-all duration-300 p-1 flex items-center justify-center rounded-lg hover:bg-white/5 shrink-0"
            title="Developer Website"
          >
            <Heart className="h-3.5 w-3.5 fill-current" />
          </a>
        </div>
      </div>
    </footer>
  );
}
