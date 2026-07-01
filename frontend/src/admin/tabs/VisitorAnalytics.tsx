import React, { useState, useEffect } from 'react';
import { 
  Users, Eye, Compass, Globe, Laptop, Smartphone, Shield, RefreshCw, Activity
} from 'lucide-react';

interface VisitorAnalyticsProps {
  API_BASE_URL: string;
}

interface ActiveVisitor {
  ip_address: string;
  last_page: string;
  browser: string;
  os: string;
  last_active: string;
  country_code: string;
  full_name?: string;
  type: 'admin' | 'visitor';
}

interface MinuteBucket {
  time: string;
  count: number;
}

interface DailyTimeline {
  date: string;
  pageviews: number;
  visitors: number;
}

interface TopPage {
  page: string;
  count: number;
}

interface TopCountry {
  country: string;
  count: number;
}

interface TopBrowser {
  browser: string;
  count: number;
}

interface TopOS {
  os: string;
  count: number;
}

interface AnalyticsData {
  active_count: number;
  active_visitors: ActiveVisitor[];
  timeline_30m: MinuteBucket[];
  timeline_30d: DailyTimeline[];
  top_pages: TopPage[];
  top_countries: TopCountry[];
  top_browsers: TopBrowser[];
  top_os: TopOS[];
}

const countryNames: Record<string, string> = {
  ID: 'Indonesia',
  SG: 'Singapura',
  MY: 'Malaysia',
  US: 'Amerika Serikat',
  GB: 'Inggris Raya',
  AU: 'Australia',
  JP: 'Jepang',
  KR: 'Korea Selatan',
  DE: 'Jerman',
  TW: 'Taiwan',
  CN: 'Tiongkok',
  HK: 'Hong Kong',
  TH: 'Thailand',
  VN: 'Vietnam',
  PH: 'Filipina',
  IN: 'India',
  NL: 'Belanda',
  FR: 'Prancis',
};

export default function VisitorAnalytics({ API_BASE_URL }: VisitorAnalyticsProps) {
  const [data, setData] = useState<AnalyticsData | null>(null);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState<'realtime' | 'historical'>('realtime');
  const [hovered30dDay, setHovered30dDay] = useState<number | null>(null);

  const getEmojiFlag = (countryCode: string) => {
    if (!countryCode || countryCode.length !== 2) return '🌐';
    const codePoints = countryCode
      .toUpperCase()
      .split('')
      .map(char => 127397 + char.charCodeAt(0));
    try {
      return String.fromCodePoint(...codePoints);
    } catch (e) {
      return '🌐';
    }
  };

  const getCountryName = (code: string) => countryNames[code.toUpperCase()] || code;

  const fetchAnalytics = () => {
    const token = localStorage.getItem('auth_token');
    if (!token) {
      setError('Token tidak ditemukan. Harap login kembali.');
      setIsLoading(false);
      return;
    }

    fetch(`${API_BASE_URL}/admin/analytics/visitors`, {
      headers: { Authorization: `Bearer ${token}` }
    })
      .then(res => {
        if (!res.ok) throw new Error('Gagal mengambil data analisis pengunjung');
        return res.json();
      })
      .then(resData => {
        setData(resData);
        setError(null);
        setIsLoading(false);
      })
      .catch(err => {
        console.error(err);
        setError(err.message);
        setIsLoading(false);
      });
  };

  useEffect(() => {
    fetchAnalytics();
    
    // Poll data every 10 seconds for real-time responsiveness
    const interval = setInterval(fetchAnalytics, 10000);
    return () => clearInterval(interval);
  }, []);

  if (isLoading && !data) {
    return (
      <div className="py-20 flex flex-col items-center justify-center gap-3">
        <RefreshCw className="h-10 w-10 text-[#002147] animate-spin" />
        <p className="text-xs text-slate-400 font-extrabold uppercase tracking-widest">Memuat Analisis Pengunjung...</p>
      </div>
    );
  }

  if (error && !data) {
    return (
      <div className="py-16 text-center max-w-md mx-auto space-y-4">
        <div className="h-12 w-12 rounded-2xl bg-red-50 text-red-500 flex items-center justify-center mx-auto border border-red-100">
          <Shield className="h-6 w-6" />
        </div>
        <div className="space-y-1">
          <h3 className="text-sm font-black text-[#002147] uppercase tracking-wider">Terjadi Kesalahan</h3>
          <p className="text-xs text-slate-400 font-medium leading-relaxed">{error}</p>
        </div>
        <button 
          onClick={() => { setIsLoading(true); fetchAnalytics(); }}
          className="px-5 py-2.5 bg-[#002147] text-white hover:bg-amber-500 rounded-xl text-xs font-bold transition-all cursor-pointer shadow-md"
        >
          Coba Lagi
        </button>
      </div>
    );
  }

  const analytics = data!;
  
  // Calculate historical totals
  const totalPageviews30d = analytics.timeline_30d.reduce((sum, d) => sum + d.pageviews, 0);
  const totalUniqueVisitors30d = analytics.timeline_30d.reduce((sum, d) => sum + d.visitors, 0);
  const avgPageviewsPerDay = Math.round(totalPageviews30d / (analytics.timeline_30d.length || 1));

  // Determine maximum values for charting
  const max30mCount = Math.max(...analytics.timeline_30m.map(m => m.count), 5);
  const max30dPageviews = Math.max(...analytics.timeline_30d.map(d => d.pageviews), 10);
  const max30dVisitors = Math.max(...analytics.timeline_30d.map(d => d.visitors), 10);
  const max30dVal = Math.max(max30dPageviews, max30dVisitors, 10);

  // SVG Line Chart path generators
  const width30d = 600;
  const height30d = 220;
  const paddingLeft = 45;
  const paddingRight = 20;
  const paddingTop = 20;
  const paddingBottom = 30;

  const chartWidth = width30d - paddingLeft - paddingRight;
  const chartHeight = height30d - paddingTop - paddingBottom;
  const dataLength = analytics.timeline_30d.length || 1;

  const getCoordinates30d = (idx: number, val: number) => {
    const x = paddingLeft + (idx / (dataLength - 1)) * chartWidth;
    const y = paddingTop + chartHeight - (val / max30dVal) * chartHeight;
    return { x, y };
  };

  // Build line path
  let pvPath = '';
  let visPath = '';
  
  if (dataLength > 1) {
    pvPath = analytics.timeline_30d.map((d, i) => {
      const { x, y } = getCoordinates30d(i, d.pageviews);
      return `${i === 0 ? 'M' : 'L'} ${x.toFixed(1)} ${y.toFixed(1)}`;
    }).join(' ');

    visPath = analytics.timeline_30d.map((d, i) => {
      const { x, y } = getCoordinates30d(i, d.visitors);
      return `${i === 0 ? 'M' : 'L'} ${x.toFixed(1)} ${y.toFixed(1)}`;
    }).join(' ');
  }

  // Formatting date helper
  const formatTime = (timeStr: string) => {
    if (!timeStr) return '-';
    try {
      const d = new Date(timeStr);
      return d.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit', second: '2-digit' });
    } catch {
      return timeStr;
    }
  };

  return (
    <div className="space-y-6 text-left animate-fade-in">
      {/* Banner / Navigation */}
      <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 border-b border-slate-100 pb-5">
        <div>
          <h1 className="text-xl lg:text-2xl font-black text-[#002147]">Analisis Pengunjung</h1>
          <p className="text-xs text-slate-400 font-medium font-sans">
            Pelacakan pengunjung situs PPID secara realtime (mirip Google Analytics).
          </p>
        </div>
        
        {/* Toggle Mode */}
        <div className="bg-slate-150/80 p-1 rounded-2xl flex items-center gap-1 border border-slate-200">
          <button
            onClick={() => setActiveTab('realtime')}
            className={`px-4 py-2 rounded-xl text-[10px] font-black uppercase tracking-wider transition-all cursor-pointer border-0 ${
              activeTab === 'realtime'
                ? 'bg-[#002147] text-white shadow-sm'
                : 'text-slate-600 hover:text-slate-900 bg-transparent'
            }`}
          >
            Real-time
          </button>
          <button
            onClick={() => setActiveTab('historical')}
            className={`px-4 py-2 rounded-xl text-[10px] font-black uppercase tracking-wider transition-all cursor-pointer border-0 ${
              activeTab === 'historical'
                ? 'bg-[#002147] text-white shadow-sm'
                : 'text-slate-600 hover:text-slate-900 bg-transparent'
            }`}
          >
            30 Hari Terakhir
          </button>
        </div>
      </div>

      {activeTab === 'realtime' ? (
        <div className="space-y-6">
          {/* Real-time Widget Grid */}
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            
            {/* Live Count Card */}
            <div className="bg-gradient-to-br from-[#002147] to-[#0A3D73] text-white rounded-3xl p-6 shadow-sm border border-slate-800 flex flex-col justify-between relative overflow-hidden">
              <div className="absolute top-0 right-0 opacity-10 translate-x-3 -translate-y-3">
                <Activity className="h-44 w-44" />
              </div>
              <div className="space-y-2 relative z-10">
                <div className="flex items-center gap-2">
                  <span className="flex h-2.5 w-2.5 relative">
                    <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-red-400 opacity-75"></span>
                    <span className="relative inline-flex rounded-full h-2.5 w-2.5 bg-red-500"></span>
                  </span>
                  <span className="text-[10px] font-black text-amber-400 uppercase tracking-widest">Sedang Online</span>
                </div>
                <h2 className="text-6xl font-black">{analytics.active_count}</h2>
                <p className="text-[10px] text-slate-300 font-semibold uppercase tracking-wider leading-relaxed">
                  Pengguna aktif di situs dalam 5 menit terakhir
                </p>
              </div>
              <div className="pt-8 border-t border-white/10 flex items-center justify-between text-xs font-extrabold text-white relative z-10">
                <span>Rasio Admin vs Publik</span>
                <span className="bg-amber-400 text-[#002147] px-2 py-0.5 rounded-lg text-[10px] font-black">
                  {analytics.active_visitors.filter(v => v.type === 'admin').length} Admin / {analytics.active_visitors.filter(v => v.type === 'visitor').length} Publik
                </span>
              </div>
            </div>

            {/* Real-time SVG Minute Chart */}
            <div className="lg:col-span-2 bg-white border border-slate-200 rounded-3xl p-6 shadow-sm flex flex-col justify-between gap-4">
              <div className="space-y-1">
                <h3 className="text-xs font-black text-[#002147] uppercase tracking-wider flex items-center gap-1.5">
                  <Activity className="h-4 w-4 text-emerald-500" />
                  Aktivitas Pengunjung per Menit (30 Menit Terakhir)
                </h3>
                <p className="text-[10px] text-slate-400 font-medium">Beban hits request dari user secara langsung.</p>
              </div>

              {/* Bar Chart SVG */}
              <div className="h-32 w-full relative">
                <svg className="w-full h-full" viewBox="0 0 500 120" preserveAspectRatio="none">
                  {/* Grid Lines */}
                  {[0, 0.5, 1].map((ratio, idx) => {
                    const yVal = 10 + ratio * 85;
                    return (
                      <line 
                        key={idx} 
                        x1="0" 
                        y1={yVal} 
                        x2="500" 
                        y2={yVal} 
                        stroke="#E2E8F0" 
                        strokeDasharray="2 3" 
                        strokeWidth="0.8" 
                      />
                    );
                  })}
                  
                  {/* Bars */}
                  {analytics.timeline_30m.map((bucket, i) => {
                    const barWidth = 8;
                    const spacing = 15;
                    const x = 15 + i * spacing;
                    const barHeight = (bucket.count / max30mCount) * 85;
                    const y = 95 - barHeight;
                    
                    return (
                      <g key={i} className="group cursor-pointer">
                        <rect
                          x={x}
                          y={y}
                          width={barWidth}
                          height={barHeight}
                          rx="2"
                          fill={bucket.count > 0 ? '#10B981' : '#E2E8F0'}
                          className="transition-all group-hover:fill-amber-400"
                        />
                        {/* Custom tooltip on hover inside SVG */}
                        <title>{`${bucket.time}: ${bucket.count} hits`}</title>
                      </g>
                    );
                  })}
                </svg>
              </div>

              {/* Chart labels */}
              <div className="flex items-center justify-between text-[8px] font-bold text-slate-400 uppercase tracking-widest border-t border-slate-50 pt-2">
                <span>30 menit lalu</span>
                <span>15 menit lalu</span>
                <span>Sekarang ({analytics.timeline_30m[analytics.timeline_30m.length - 1]?.time || '00:00'})</span>
              </div>
            </div>
          </div>

          {/* Active Users Table Detail */}
          <div className="bg-white border border-slate-200 rounded-3xl p-6 shadow-sm space-y-4">
            <div className="space-y-1">
              <h3 className="text-sm font-black text-[#002147] uppercase tracking-wider">Rincian Pengguna Online</h3>
              <p className="text-[11px] text-slate-400 font-medium font-sans">
                Daftar IP address dan halaman yang sedang diakses saat ini secara real-time.
              </p>
            </div>

            <div className="overflow-x-auto">
              {analytics.active_visitors.length > 0 ? (
                <table className="w-full text-left border-collapse text-xs font-semibold text-slate-650">
                  <thead>
                    <tr className="border-b border-slate-100 text-[10px] uppercase text-slate-400">
                      <th className="py-2.5 font-extrabold">Negara</th>
                      <th className="py-2.5 font-extrabold">IP Address</th>
                      <th className="py-2.5 font-extrabold">Status/Pengguna</th>
                      <th className="py-2.5 font-extrabold">Halaman Diakses</th>
                      <th className="py-2.5 font-extrabold">Sistem / Browser</th>
                      <th className="py-2.5 font-extrabold text-right">Terakhir Aktif</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-slate-50">
                    {analytics.active_visitors.map((visitor, idx) => (
                      <tr key={idx} className="hover:bg-slate-50/50 transition-colors">
                        <td className="py-3">
                          <span className="inline-flex items-center gap-1.5 bg-slate-50 px-2 py-1 rounded-lg border border-slate-100">
                            <span className="text-base leading-none">{getEmojiFlag(visitor.country_code)}</span>
                            <span className="text-[9px] font-bold text-slate-500 uppercase">{visitor.country_code}</span>
                          </span>
                        </td>
                        <td className="py-3 font-mono font-bold text-[#002147]">{visitor.ip_address}</td>
                        <td className="py-3">
                          {visitor.type === 'admin' ? (
                            <span className="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full bg-amber-50 text-amber-700 text-[8px] font-black uppercase tracking-wide border border-amber-100">
                              🔐 {visitor.full_name || 'Admin'}
                            </span>
                          ) : (
                            <span className="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full bg-slate-100 text-slate-600 text-[8px] font-black uppercase tracking-wide border border-slate-200">
                              👥 Publik
                            </span>
                          )}
                        </td>
                        <td className="py-3">
                          <span className="text-slate-800 font-bold block truncate max-w-[200px]" title={visitor.last_page}>
                            {visitor.last_page}
                          </span>
                        </td>
                        <td className="py-3">
                          <span className="text-[10px] text-slate-450 block font-medium">
                            {visitor.os} · {visitor.browser}
                          </span>
                        </td>
                        <td className="py-3 text-right text-slate-400 font-bold font-sans">
                          {formatTime(visitor.last_active)}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              ) : (
                <div className="py-12 text-center text-slate-400 border border-dashed border-slate-200 rounded-2xl bg-slate-50/40 font-bold uppercase text-[10px] tracking-wider">
                  Tidak ada pengguna online yang terdeteksi
                </div>
              )}
            </div>
          </div>
        </div>
      ) : (
        <div className="space-y-6">
          {/* Key Metric cards */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            {[
              { label: 'Total Tayangan Halaman', val: totalPageviews30d.toLocaleString(), desc: 'Jumlah total halaman diklik', color: 'text-blue-500 bg-blue-50 border-blue-100', icon: Eye },
              { label: 'Pengunjung Unik', val: totalUniqueVisitors30d.toLocaleString(), desc: 'Berdasarkan IP unik user', color: 'text-purple-500 bg-purple-50 border-purple-100', icon: Users },
              { label: 'Rata-Rata Harian', val: avgPageviewsPerDay.toLocaleString(), desc: 'Trafik tayangan per hari', color: 'text-emerald-500 bg-emerald-50 border-emerald-100', icon: Compass }
            ].map((card, idx) => {
              const Icon = card.icon;
              return (
                <div key={idx} className="bg-white p-5 rounded-2xl border border-slate-200 shadow-sm flex items-start gap-4">
                  <div className={`flex h-11 w-11 shrink-0 items-center justify-center rounded-xl border ${card.color}`}>
                    <Icon className="h-5 w-5" />
                  </div>
                  <div className="space-y-1">
                    <span className="text-[10px] font-extrabold text-slate-400 uppercase tracking-wider block">{card.label}</span>
                    <span className="text-2xl font-black text-slate-800 block leading-none">{card.val}</span>
                    <span className="text-[9px] text-slate-400 font-semibold block">{card.desc}</span>
                  </div>
                </div>
              );
            })}
          </div>

          {/* Historical Trend Double-Line Chart */}
          <div className="bg-white border border-slate-200 rounded-3xl p-6 shadow-sm space-y-4">
            <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-2 border-b border-slate-50 pb-4">
              <div className="space-y-1">
                <h3 className="text-sm font-black text-[#002147] uppercase tracking-wider flex items-center gap-1.5">
                  Tren Kunjungan Harian (30 Hari Terakhir)
                </h3>
                <p className="text-[10px] text-slate-400 font-medium">Bandingan grafik Pageviews dan Pengunjung Unik harian.</p>
              </div>

              {/* Legend indicators */}
              <div className="flex items-center gap-4 text-[10px] font-bold uppercase tracking-wider">
                <div className="flex items-center gap-1.5">
                  <span className="h-2.5 w-2.5 rounded-full bg-[#002147]"></span>
                  <span className="text-slate-650">Pageviews</span>
                </div>
                <div className="flex items-center gap-1.5">
                  <span className="h-2.5 w-2.5 rounded-full bg-emerald-500"></span>
                  <span className="text-slate-650">Pengunjung Unik</span>
                </div>
              </div>
            </div>

            {/* Line Chart SVG */}
            <div className="w-full relative">
              <svg 
                className="w-full h-auto" 
                viewBox={`0 0 ${width30d} ${height30d}`} 
                fill="none" 
                xmlns="http://www.w3.org/2000/svg"
              >
                {/* Horizontal Grid lines */}
                {[0, 0.25, 0.5, 0.75, 1].map((ratio, idx) => {
                  const y = paddingTop + ratio * chartHeight;
                  const labelVal = Math.round(max30dVal * (1 - ratio));
                  return (
                    <g key={idx}>
                      <line 
                        x1={paddingLeft} 
                        y1={y} 
                        x2={width30d - paddingRight} 
                        y2={y} 
                        stroke="#F1F5F9" 
                        strokeWidth="1" 
                      />
                      <text 
                        x={paddingLeft - 8} 
                        y={y + 3} 
                        fill="#94A3B8" 
                        fontSize="9" 
                        fontWeight="bold" 
                        fontFamily="system-ui"
                        textAnchor="end"
                      >
                        {labelVal}
                      </text>
                    </g>
                  );
                })}

                {/* Draw Areas under lines */}
                {pvPath && (
                  <path
                    d={`${pvPath} L ${(paddingLeft + chartWidth).toFixed(1)} ${(paddingTop + chartHeight).toFixed(1)} L ${paddingLeft} ${(paddingTop + chartHeight).toFixed(1)} Z`}
                    fill="url(#pvGradient)"
                    opacity="0.04"
                  />
                )}
                {visPath && (
                  <path
                    d={`${visPath} L ${(paddingLeft + chartWidth).toFixed(1)} ${(paddingTop + chartHeight).toFixed(1)} L ${paddingLeft} ${(paddingTop + chartHeight).toFixed(1)} Z`}
                    fill="url(#visGradient)"
                    opacity="0.05"
                  />
                )}

                {/* Line definitions */}
                <defs>
                  <linearGradient id="pvGradient" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="0%" stopColor="#002147" />
                    <stop offset="100%" stopColor="#002147" stopOpacity="0" />
                  </linearGradient>
                  <linearGradient id="visGradient" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="0%" stopColor="#10B981" />
                    <stop offset="100%" stopColor="#10B981" stopOpacity="0" />
                  </linearGradient>
                </defs>

                {/* Draw Line paths */}
                {pvPath && (
                  <path 
                    d={pvPath} 
                    stroke="#002147" 
                    strokeWidth="2.5" 
                    strokeLinecap="round" 
                    strokeLinejoin="round" 
                  />
                )}
                {visPath && (
                  <path 
                    d={visPath} 
                    stroke="#10B981" 
                    strokeWidth="2" 
                    strokeLinecap="round" 
                    strokeLinejoin="round" 
                  />
                )}

                {/* Data points and hover catchers */}
                {analytics.timeline_30d.map((d, i) => {
                  const ptPV = getCoordinates30d(i, d.pageviews);
                  const ptVis = getCoordinates30d(i, d.visitors);
                  const isHovered = hovered30dDay === i;

                  // Render X Labels on alternate days to avoid clutter
                  const showLabel = i % 4 === 0 || i === dataLength - 1;

                  return (
                    <g key={i}>
                      {showLabel && (
                        <text
                          x={ptPV.x}
                          y={paddingTop + chartHeight + 16}
                          fill="#94A3B8"
                          fontSize="8"
                          fontWeight="bold"
                          fontFamily="system-ui"
                          textAnchor="middle"
                        >
                          {d.date}
                        </text>
                      )}

                      {/* Interactive dots on hover */}
                      {isHovered && (
                        <>
                          <line
                            x1={ptPV.x}
                            y1={paddingTop}
                            x2={ptPV.x}
                            y2={paddingTop + chartHeight}
                            stroke="#CBD5E1"
                            strokeWidth="1"
                            strokeDasharray="2 2"
                          />
                          <circle cx={ptPV.x} cy={ptPV.y} r="5" fill="#002147" stroke="#FFFFFF" strokeWidth="1.5" />
                          <circle cx={ptVis.x} cy={ptVis.y} r="4" fill="#10B981" stroke="#FFFFFF" strokeWidth="1.5" />
                        </>
                      )}

                      {/* Invisible wider slice to make hovering easy */}
                      <rect
                        x={ptPV.x - (chartWidth / (dataLength - 1)) / 2}
                        y={paddingTop}
                        width={chartWidth / (dataLength - 1)}
                        height={chartHeight}
                        fill="transparent"
                        className="cursor-pointer"
                        onMouseEnter={() => setHovered30dDay(i)}
                        onMouseLeave={() => setHovered30dDay(null)}
                      />
                    </g>
                  );
                })}
              </svg>
              
              {/* Tooltip Overlay */}
              {hovered30dDay !== null && analytics.timeline_30d[hovered30dDay] && (
                <div 
                  className="absolute bg-[#002147] text-white p-2.5 rounded-xl text-[9px] font-bold shadow-xl border border-white/10 pointer-events-none flex flex-col gap-1 z-25"
                  style={{
                    left: `${((getCoordinates30d(hovered30dDay, 0).x / width30d) * 100).toFixed(1)}%`,
                    top: '20px',
                    transform: hovered30dDay > dataLength / 2 ? 'translateX(-105%)' : 'translateX(5%)'
                  }}
                >
                  <span className="text-amber-400 font-extrabold uppercase border-b border-white/10 pb-1 block leading-none">
                    📅 {analytics.timeline_30d[hovered30dDay].date}
                  </span>
                  <span className="flex items-center justify-between gap-4 mt-1">
                    <span>Pageviews:</span>
                    <span className="font-mono text-white text-xs">{analytics.timeline_30d[hovered30dDay].pageviews}</span>
                  </span>
                  <span className="flex items-center justify-between gap-4">
                    <span>Pengunjung Unik:</span>
                    <span className="font-mono text-emerald-400 text-xs">{analytics.timeline_30d[hovered30dDay].visitors}</span>
                  </span>
                </div>
              )}
            </div>
          </div>

          {/* Breakdowns grid */}
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            
            {/* Top Countries Card */}
            <div className="bg-white border border-slate-200 rounded-3xl p-6 shadow-sm space-y-4">
              <div className="space-y-1">
                <h3 className="text-sm font-black text-[#002147] uppercase tracking-wider flex items-center gap-1.5">
                  <Globe className="h-4.5 w-4.5 text-blue-500" />
                  Negara Asal Pengunjung
                </h3>
                <p className="text-[10px] text-slate-400 font-medium">Berdasarkan 30 hari akumulasi IP geolocation.</p>
              </div>

              <div className="space-y-3.5 pt-2">
                {analytics.top_countries.length > 0 ? (
                  analytics.top_countries.map((item, idx) => {
                    // Calculate percentage of total hits
                    const totalCountryHits = analytics.top_countries.reduce((s, c) => s + c.count, 0);
                    const pct = Math.round((item.count / (totalCountryHits || 1)) * 100);
                    
                    return (
                      <div key={idx} className="space-y-1">
                        <div className="flex items-center justify-between text-xs font-bold text-slate-700">
                          <div className="flex items-center gap-2">
                            <span className="text-lg leading-none">{getEmojiFlag(item.country)}</span>
                            <span>{getCountryName(item.country)}</span>
                            <span className="text-[8px] font-black text-slate-400 bg-slate-100 px-1.5 py-0.5 rounded uppercase">
                              {item.country}
                            </span>
                          </div>
                          <span className="font-mono">{item.count.toLocaleString()} hits ({pct}%)</span>
                        </div>
                        {/* Progress Bar */}
                        <div className="h-2 w-full bg-slate-100 rounded-full overflow-hidden">
                          <div 
                            className="h-full bg-gradient-to-r from-[#002147] to-amber-500 rounded-full"
                            style={{ width: `${pct}%` }}
                          />
                        </div>
                      </div>
                    );
                  })
                ) : (
                  <div className="py-12 text-center text-slate-400 border border-dashed border-slate-200 rounded-2xl bg-slate-50/40 font-bold uppercase text-[9px] tracking-wider">
                    Belum ada data lokasi tercatat
                  </div>
                )}
              </div>
            </div>

            {/* Top Visited Pages Card */}
            <div className="bg-white border border-slate-200 rounded-3xl p-6 shadow-sm space-y-4">
              <div className="space-y-1">
                <h3 className="text-sm font-black text-[#002147] uppercase tracking-wider flex items-center gap-1.5">
                  <Compass className="h-4.5 w-4.5 text-rose-500" />
                  Halaman Paling Populer
                </h3>
                <p className="text-[10px] text-slate-400 font-medium">Halaman yang paling sering dikunjungi publik.</p>
              </div>

              <div className="space-y-1.5 pt-2">
                {analytics.top_pages.length > 0 ? (
                  <div className="border border-slate-100 rounded-2xl overflow-hidden divide-y divide-slate-100">
                    {analytics.top_pages.map((item, idx) => (
                      <div key={idx} className="p-3 hover:bg-slate-50 transition-colors flex items-center justify-between gap-4 text-xs font-bold">
                        <span className="text-[#002147] font-mono truncate max-w-[250px]" title={item.page}>
                          {item.page}
                        </span>
                        <span className="font-mono text-slate-500 bg-slate-100 px-2 py-0.5 rounded-lg shrink-0">
                          {item.count.toLocaleString()} tayangan
                        </span>
                      </div>
                    ))}
                  </div>
                ) : (
                  <div className="py-12 text-center text-slate-400 border border-dashed border-slate-200 rounded-2xl bg-slate-50/40 font-bold uppercase text-[9px] tracking-wider">
                    Belum ada data kunjungan halaman
                  </div>
                )}
              </div>
            </div>
          </div>

          {/* Browser & OS Card Grid */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            
            {/* Top Browsers */}
            <div className="bg-white border border-slate-200 rounded-3xl p-6 shadow-sm space-y-4">
              <div className="space-y-1">
                <h3 className="text-sm font-black text-[#002147] uppercase tracking-wider flex items-center gap-1.5">
                  <Laptop className="h-4.5 w-4.5 text-purple-500" />
                  Peramban (Browsers)
                </h3>
                <p className="text-[10px] text-slate-400 font-medium">Distribusi browser yang digunakan pengunjung.</p>
              </div>

              <div className="space-y-3 pt-2">
                {analytics.top_browsers.length > 0 ? (
                  analytics.top_browsers.map((item, idx) => {
                    const totalBrowsers = analytics.top_browsers.reduce((s, b) => s + b.count, 0);
                    const pct = Math.round((item.count / (totalBrowsers || 1)) * 100);
                    return (
                      <div key={idx} className="flex items-center justify-between text-xs font-bold text-slate-700">
                        <span className="capitalize">{item.browser}</span>
                        <div className="flex items-center gap-3 w-1/2">
                          <div className="h-1.5 flex-1 bg-slate-200 rounded-full overflow-hidden">
                            <div className="h-full bg-purple-500 rounded-full" style={{ width: `${pct}%` }} />
                          </div>
                          <span className="font-mono text-slate-500 shrink-0 w-12 text-right">{pct}%</span>
                        </div>
                      </div>
                    );
                  })
                ) : (
                  <div className="py-6 text-center text-slate-400 border border-dashed border-slate-250 rounded-2xl bg-slate-50/40 font-bold uppercase text-[8px] tracking-wider">
                    Belum ada data peramban
                  </div>
                )}
              </div>
            </div>

            {/* Top OS */}
            <div className="bg-white border border-slate-200 rounded-3xl p-6 shadow-sm space-y-4">
              <div className="space-y-1">
                <h3 className="text-sm font-black text-[#002147] uppercase tracking-wider flex items-center gap-1.5">
                  <Smartphone className="h-4.5 w-4.5 text-emerald-500" />
                  Sistem Operasi (OS)
                </h3>
                <p className="text-[10px] text-slate-400 font-medium">Distribusi OS yang digunakan pengunjung.</p>
              </div>

              <div className="space-y-3 pt-2">
                {analytics.top_os.length > 0 ? (
                  analytics.top_os.map((item, idx) => {
                    const totalOS = analytics.top_os.reduce((s, o) => s + o.count, 0);
                    const pct = Math.round((item.count / (totalOS || 1)) * 100);
                    return (
                      <div key={idx} className="flex items-center justify-between text-xs font-bold text-slate-700">
                        <span>{item.os}</span>
                        <div className="flex items-center gap-3 w-1/2">
                          <div className="h-1.5 flex-1 bg-slate-200 rounded-full overflow-hidden">
                            <div className="h-full bg-emerald-500 rounded-full" style={{ width: `${pct}%` }} />
                          </div>
                          <span className="font-mono text-slate-500 shrink-0 w-12 text-right">{pct}%</span>
                        </div>
                      </div>
                    );
                  })
                ) : (
                  <div className="py-6 text-center text-slate-400 border border-dashed border-slate-250 rounded-2xl bg-slate-50/40 font-bold uppercase text-[8px] tracking-wider">
                    Belum ada data OS
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
