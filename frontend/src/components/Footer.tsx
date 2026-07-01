import { Heart } from 'lucide-react';

interface FooterProps {
  navigateToHome: () => void;
  navigateToPage: (slug: string) => void;
  setCurrentPage: (page: string) => void;
  playstoreLink?: string;
}

export default function Footer({
  navigateToHome,
  navigateToPage,
  setCurrentPage,
  playstoreLink
}: FooterProps) {

  return (
    <footer className="bg-[#001733] text-slate-400 py-8 px-4 border-t border-white/5 text-xs">
      <div className="container mx-auto max-w-6xl flex flex-col md:flex-row items-center justify-between gap-4">
        <div className="text-center md:text-left space-y-3">
          <div className="space-y-1">
            <span className="text-slate-200 block font-bold">PPID Universitas Perintis Indonesia (UPERTIS)</span>
            <p className="text-[11px] font-medium text-slate-500">&copy; {new Date().getFullYear()} Universitas Perintis Indonesia. All Rights Reserved.</p>
          </div>
          {playstoreLink && (
            <div className="pt-1 flex justify-center md:justify-start">
              <a
                href={playstoreLink}
                target="_blank"
                rel="noopener noreferrer"
                className="flex items-center gap-2 bg-black text-white px-3 py-1.5 rounded-xl border border-zinc-800 hover:border-zinc-700 hover:bg-zinc-900 transition-all text-left group shrink-0"
              >
                <svg
                  className="h-5 w-5 fill-current text-white group-hover:scale-110 transition-transform"
                  viewBox="0 0 512 512"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path d="M325.3 234.3L104.6 13l280.8 161.2-60.1 60.1zM47 0C34 0 25 9.3 25 22.3v467.4C25 502.7 34 512 47 512c5.8 0 11.6-2.2 16.5-6.6L293 276.5 63.5 6.6C58.6 2.2 52.8 0 47 0zm397.6 182.2L354.2 245l64.1 64.1L490 263.2c13.7-7.9 22-22.3 22-38.2 0-15.9-8.3-30.3-22-38.2l-45.4-24.6zM325.3 277.7l60.1 60.1L104.6 499l220.7-221.3z" />
                </svg>
                <div className="flex flex-col leading-none">
                  <span className="text-[7px] text-zinc-400 uppercase tracking-widest font-semibold">GET IT ON</span>
                  <span className="text-[12px] font-black tracking-tight text-white mt-0.5">Google Play</span>
                </div>
              </a>
            </div>
          )}
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
