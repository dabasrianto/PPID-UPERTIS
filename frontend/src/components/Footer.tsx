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
                  className="h-5 w-5 group-hover:scale-110 transition-transform shrink-0"
                  viewBox="0 0 466 511.98"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <g fillRule="nonzero">
                    <path fill="#EA4335" d="M199.9 237.8 1.4 470.17c7.22 24.57 30.16 41.81 55.8 41.81 11.16 0 20.93-2.79 29.3-8.37l244.16-139.46L199.9 237.8z"/>
                    <path fill="#FBBC04" d="m433.91 205.1-104.65-60-111.61 110.22 113.01 108.83 104.64-58.6c18.14-9.77 30.7-29.3 30.7-50.23-1.4-20.93-13.95-40.46-32.09-50.22z"/>
                    <path fill="#34A853" d="M199.42 273.45 329.27 145.1 87.9 8.37C79.53 2.79 68.36 0 57.2 0 30.7 0 6.98 18.14 1.4 41.86l198.02 231.59z"/>
                    <path fill="#4285F4" d="M1.39 41.86C0 46.04 0 51.63 0 57.2v397.64c0 5.57 0 9.76 1.4 15.34l216.27-214.86L1.39 41.86z"/>
                  </g>
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
