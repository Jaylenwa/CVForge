import React, { useEffect, useState } from 'react';
import { useLanguage } from '../../contexts/LanguageContext';
import { getAdminStats, AdminStats } from '../../services/adminService';
import { RefreshCw } from 'lucide-react';
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

  // Simple SVG line chart
  const Chart: React.FC<{ dates: string[]; a: number[]; b: number[]; la: string; lb: string }> = ({ dates, a, b, la, lb }) => {
    const w = 640, h = 200, p = 24;
    const max = Math.max(1, ...a, ...b);
    const stepX = (w - p * 2) / Math.max(1, dates.length - 1);
    const toY = (v: number) => h - p - (v / max) * (h - p * 2);
    const toX = (i: number) => p + i * stepX;
    const path = (arr: number[]) => arr.map((v, i) => `${i === 0 ? 'M' : 'L'} ${toX(i)} ${toY(v)}`).join(' ');
    return (
      <svg width={w} height={h} className="bg-white rounded-md border border-gray-200">
        {/* axes */}
        <line x1={p} y1={h - p} x2={w - p} y2={h - p} stroke="#e5e7eb" />
        <line x1={p} y1={p} x2={p} y2={h - p} stroke="#e5e7eb" />
        {/* grid */}
        {Array.from({ length: 4 }).map((_, i) => {
          const y = p + ((h - p * 2) / 4) * i;
          return <line key={i} x1={p} y1={y} x2={w - p} y2={y} stroke="#f3f4f6" />;
        })}
        {/* lines */}
        <path d={path(a)} fill="none" stroke="#3b82f6" strokeWidth={2} />
        <path d={path(b)} fill="none" stroke="#22c55e" strokeWidth={2} />
        {/* legend */}
        <g>
          <circle cx={w - 180} cy={20} r={4} fill="#3b82f6" />
          <text x={w - 170} y={24} fontSize="12" fill="#374151">{la}</text>
          <circle cx={w - 100} cy={20} r={4} fill="#22c55e" />
          <text x={w - 90} y={24} fontSize="12" fill="#374151">{lb}</text>
        </g>
      </svg>
    );
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-gray-900">{t('admin.menu.dashboard')}</h1>
        <Button variant="outline" onClick={reload} disabled={loading}>
          <RefreshCw size={16} className={`${loading ? 'animate-spin' : ''} mr-2`} /> {t('common.refresh') || 'Refresh'}
        </Button>
      </div>

      {/* Cards */}
      <div className="grid grid-cols-12 gap-6">
        <div className="col-span-9">
          <div className="bg-white rounded-lg shadow p-4">
            <div className="flex items-center justify-between mb-2">
              <div className="text-lg font-medium">{t('admin.stats.trend')}</div>
              {stats && <div className="text-xs text-gray-500">{t('admin.stats.generated')} {new Date(stats.generatedAt * 1000).toLocaleString()}</div>}
            </div>
            {stats ? (
              <Chart dates={stats.trend.dates} a={stats.trend.users} b={stats.trend.resumes} la={t('admin.stats.users')} lb={t('admin.stats.resumes')} />
            ) : (
              <div className="h-48 flex items-center justify-center text-gray-400">{t('common.loading') || 'Loading...'}</div>
            )}
          </div>
        </div>
        <div className="col-span-3">
          <div className="bg-white rounded-lg shadow p-4">
            <div className="text-lg font-medium mb-2">{t('admin.stats.totals')}</div>
            <div className="space-y-3">
              <div className="flex items-center justify-between">
                <span className="flex items-center text-gray-700 text-sm">
                  <span className="w-2.5 h-2.5 bg-blue-500 rounded-full mr-2"></span>
                  {t('admin.stats.users')}
                </span>
                <span className="text-lg font-semibold">{stats?.totals.users ?? '-'}</span>
              </div>
              <div className="flex items-center justify-between">
                <span className="flex items-center text-gray-700 text-sm">
                  <span className="w-2.5 h-2.5 bg-green-500 rounded-full mr-2"></span>
                  {t('admin.stats.resumes')}
                </span>
                <span className="text-lg font-semibold">{stats?.totals.resumes ?? '-'}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
