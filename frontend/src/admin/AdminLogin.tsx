import React from 'react';
import { AlertTriangle, Eye, EyeOff } from 'lucide-react';

interface AdminLoginProps {
  adminLoginError: string | null;
  adminEmail: string;
  setAdminEmail: (val: string) => void;
  adminPassword: string;
  setAdminPassword: (val: string) => void;
  showPassword: boolean;
  setShowPassword: (val: boolean) => void;
  adminLoginLoading: boolean;
  handleAdminLogin: (e: React.FormEvent) => void;
}

export default function AdminLogin({
  adminLoginError,
  adminEmail,
  setAdminEmail,
  adminPassword,
  setAdminPassword,
  showPassword,
  setShowPassword,
  adminLoginLoading,
  handleAdminLogin
}: AdminLoginProps) {
  return (
    <div className="flex-1 flex items-center justify-center py-16 px-4">
      <div className="w-full max-w-md bg-white border border-slate-200 rounded-[2rem] p-8 shadow-xl text-left space-y-6 animate-in fade-in zoom-in-95 duration-200">
        <div className="flex items-center gap-3">
          <div className="h-10 w-10 bg-[#002147] text-white rounded-xl flex items-center justify-center font-bold text-sm">
            PPID
          </div>
          <div>
            <h1 className="font-extrabold text-sm text-[#002147] leading-none">Admin Portal PPID</h1>
            <span className="text-[10px] text-slate-400 font-semibold mt-1 block">Universitas Perintis Indonesia</span>
          </div>
        </div>

        <div className="space-y-1">
          <h2 className="text-xl font-extrabold text-slate-800">Masuk Akun</h2>
          <p className="text-xs text-slate-400 font-medium">Masukkan email dan password admin Anda.</p>
        </div>

        {adminLoginError && (
          <div className="bg-red-50 border border-red-200 text-red-600 rounded-xl p-3.5 text-xs font-bold flex items-center gap-2">
            <AlertTriangle className="h-4.5 w-4.5" />
            <span>{adminLoginError}</span>
          </div>
        )}

        <form onSubmit={handleAdminLogin} className="space-y-4">
          <div className="space-y-1">
            <label className="text-[10px] font-bold text-slate-600 uppercase tracking-wider block">Alamat Email</label>
            <input
              type="email"
              required
              placeholder="admin@upertis.ac.id"
              value={adminEmail}
              onChange={(e) => setAdminEmail(e.target.value)}
              className="w-full rounded-xl border border-slate-200 px-3.5 py-2.5 text-xs focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 transition-all bg-slate-50 font-medium text-slate-800"
            />
          </div>

          <div className="space-y-1">
            <label className="text-[10px] font-bold text-slate-600 uppercase tracking-wider block">Password Akun</label>
            <div className="relative">
              <input
                type={showPassword ? 'text' : 'password'}
                required
                placeholder="••••••••"
                value={adminPassword}
                onChange={(e) => setAdminPassword(e.target.value)}
                className="w-full rounded-xl border border-slate-200 pl-3.5 pr-10 py-2.5 text-xs focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 transition-all bg-slate-50 font-medium text-slate-800"
              />
              <button
                type="button"
                onClick={() => setShowPassword(!showPassword)}
                className="absolute right-3 top-3 text-slate-400 hover:text-slate-600 cursor-pointer"
              >
                {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
              </button>
            </div>
          </div>

          <button
            type="submit"
            disabled={adminLoginLoading}
            className="w-full py-3 bg-[#002147] hover:bg-amber-400 hover:text-[#002147] text-white font-bold text-xs uppercase tracking-wider rounded-xl transition-all shadow-md cursor-pointer"
          >
            {adminLoginLoading ? 'Menghubungkan...' : 'Masuk Dashboard'}
          </button>
        </form>
      </div>
    </div>
  );
}
