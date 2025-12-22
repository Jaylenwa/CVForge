import React, { useEffect, useRef, useState } from 'react';
import { useLanguage } from '../../contexts/LanguageContext';
import { getAdminStats, AdminStats } from '../../services/adminService';
import { RefreshCw, Users, FileText } from 'lucide-react';
import { Button } from '../../components/ui/Button';

export const AdminHome: React.FC = () => {
  const { t } = useLanguage();
  const [stats, setStats] = useState<AdminStats | null>(null);
  const [loading, setLoading] = useState(true);

  const reload = async () => {
    setLoading(true);
    try {
      const s = await getAdminStats(14);
      setStats(s);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { reload(); }, []);

  const Chart: React.FC<{ dates: string[]; a: number[]; b: number[]; la: string; lb: string }> = ({ dates, a, b, la, lb }) => {
    const w = 720, h = 350, p = 28;
    const [hoverIndex, setHoverIndex] = useState<number | null>(null);
    const svgRef = useRef<SVGSVGElement | null>(null);
    const max = Math.max(1, ...a, ...b);
    const stepX = (w - p * 2) / Math.max(1, dates.length - 1);
    const toY = (v: number) => h - p - (v / max) * (h - p * 2);
    const toX = (i: number) => p + i * stepX;
    const path = (arr: number[]) => arr.map((v, i) => `${i === 0 ? 'M' : 'L'} ${toX(i)} ${toY(v)}`).join(' ');
    const handleMove = (e: React.MouseEvent<SVGSVGElement>) => {
      const rect = (e.target as Element).closest('svg')?.getBoundingClientRect();
      const x = rect ? e.clientX - rect.left : 0;
      const i = Math.max(0, Math.min(dates.length - 1, Math.round((x - p) / stepX)));
      setHoverIndex(i);
    };
    const handleLeave = () => setHoverIndex(null);
    return (
      <svg ref={svgRef} width={w} height={h} className="rounded-md border border-gray-200" onMouseMove={handleMove} onMouseLeave={handleLeave}>
        <rect x={0} y={0} width={w} height={h} fill="#ffffff" rx={12} />
        <line x1={p} y1={h - p} x2={w - p} y2={h - p} stroke="#e5e7eb" />
        <line x1={p} y1={p} x2={p} y2={h - p} stroke="#e5e7eb" />
        {Array.from({ length: 4 }).map((_, i) => {
          const y = p + ((h - p * 2) / 4) * i;
          return <line key={`hg-${i}`} x1={p} y1={y} x2={w - p} y2={y} stroke="#f3f4f6" strokeDasharray="3 3" />;
        })}
        {dates.map((_, i) => {
          const x = toX(i);
          return <line key={`vg-${i}`} x1={x} y1={p} x2={x} y2={h - p} stroke="#f3f4f6" strokeDasharray="3 3" />;
        })}
        <path d={path(a)} fill="none" stroke="#3b82f6" strokeWidth={2} />
        <path d={path(b)} fill="none" stroke="#22c55e" strokeWidth={2} />
        {a.map((v, i) => (
          <circle key={`ua-${i}`} cx={toX(i)} cy={toY(v)} r={hoverIndex === i ? 5 : 3} fill="#ffffff" stroke="#3b82f6" strokeWidth={hoverIndex === i ? 2 : 1.5} />
        ))}
        {b.map((v, i) => (
          <circle key={`ub-${i}`} cx={toX(i)} cy={toY(v)} r={hoverIndex === i ? 5 : 3} fill="#ffffff" stroke="#22c55e" strokeWidth={hoverIndex === i ? 2 : 1.5} />
        ))}
        <g>
          <circle cx={w - 200} cy={24} r={5} fill="#3b82f6" />
          <text x={w - 188} y={28} fontSize="13" fill="#374151">{la}</text>
          <circle cx={w - 120} cy={24} r={5} fill="#22c55e" />
          <text x={w - 108} y={28} fontSize="13" fill="#374151">{lb}</text>
        </g>
        {hoverIndex !== null && (
          <>
            <line x1={toX(hoverIndex)} y1={p} x2={toX(hoverIndex)} y2={h - p} stroke="#e5e7eb" />
            <rect x={Math.min(w - p - 160, Math.max(p, toX(hoverIndex) + 12))} y={p + 12} width={150} height={60} rx={8} fill="#ffffff" stroke="#e5e7eb" />
            <text x={Math.min(w - p - 150, Math.max(p + 10, toX(hoverIndex) + 22))} y={p + 32} fontSize="12" fill="#6b7280">{dates[hoverIndex]}</text>
            <circle cx={Math.min(w - p - 150, Math.max(p + 10, toX(hoverIndex) + 22))} cy={p + 48} r={4} fill="#3b82f6" />
            <text x={Math.min(w - p - 140, Math.max(p + 20, toX(hoverIndex) + 32))} y={p + 52} fontSize="12" fill="#374151">{la}: {a[hoverIndex]}</text>
            <circle cx={Math.min(w - p - 90, Math.max(p + 10, toX(hoverIndex) + 92))} cy={p + 48} r={4} fill="#22c55e" />
            <text x={Math.min(w - p - 80, Math.max(p + 20, toX(hoverIndex) + 102))} y={p + 52} fontSize="12" fill="#374151">{lb}: {b[hoverIndex]}</text>
          </>
        )}
        {dates.map((d, i) => (
          <text key={`tx-${i}`} x={toX(i)} y={h - p + 16} fontSize="12" fill="#9ca3af" textAnchor="middle">{d}</text>
        ))}
        {[0, Math.ceil(max / 4), Math.ceil(max / 2), Math.ceil((3 * max) / 4), max].map((tick, i) => (
          <text key={`ty-${i}`} x={p - 6} y={toY(tick)} fontSize="12" fill="#9ca3af" textAnchor="end">{tick}</text>
        ))}
      </svg>
    );
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-4xl font-bold text-gray-800 tracking-tight">{t('admin.menu.dashboard')}</h1>
        <Button variant="outline" onClick={reload} disabled={loading}>
          <RefreshCw size={16} className={`${loading ? 'animate-spin' : ''} mr-2`} /> {t('common.refresh') || 'Refresh'}
        </Button>
      </div>
      <div className="flex flex-col lg:flex-row gap-6 items-stretch">
        <div className="bg-white rounded-2xl border border-gray-100 p-6 flex-1 shadow-sm">
          <div className="flex justify-between items-center mb-6">
            <h3 className="text-gray-800 font-semibold text-lg">{t('admin.stats.trend')}</h3>
            {stats && <span className="text-gray-400 text-sm">{t('admin.stats.generated')} {new Date(stats.generatedAt * 1000).toLocaleString()}</span>}
          </div>
          <div className="h-[350px] w-full flex items-center justify-center">
            {stats ? (
              <Chart dates={stats.trend.dates} a={stats.trend.users} b={stats.trend.resumes} la={t('admin.stats.users')} lb={t('admin.stats.resumes')} />
            ) : (
              <div className="h-48 flex items-center justify-center text-gray-400">{t('common.loading') || 'Loading...'}</div>
            )}
          </div>
        </div>
        <div className="w-72 bg-white rounded-2xl border border-gray-100 p-6 shadow-sm flex flex-col">
          <h3 className="text-gray-800 font-bold text-lg mb-4">{t('admin.stats.totals')}</h3>
          <div className="h-px bg-gray-100 w-full mb-6"></div>
          <div className="space-y-8">
            <div className="flex items-center space-x-4">
              <div className="w-11 h-11 bg-blue-100 text-blue-500 rounded-full flex items-center justify-center">
                <Users size={20} />
              </div>
              <div>
                <p className="text-xl font-bold text-gray-800 leading-none mb-1">{stats?.totals.users ?? '-'}</p>
                <p className="text-gray-500 text-sm">{t('admin.stats.users')}</p>
              </div>
            </div>
            <div className="flex items-center space-x-4">
              <div className="w-11 h-11 bg-green-100 text-green-500 rounded-full flex items-center justify-center">
                <FileText size={20} />
              </div>
              <div>
                <p className="text-xl font-bold text-gray-800 leading-none mb-1">{stats?.totals.resumes ?? '-'}</p>
                <p className="text-gray-500 text-sm">{t('admin.stats.resumes')}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div className="mt-8 flex-1"></div>
    </div>
  );
};
