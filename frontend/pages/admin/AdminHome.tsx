import React, { useEffect, useRef, useState } from 'react';
import { useLanguage } from '../../contexts/LanguageContext';
import { getAdminStats, AdminStats, importDefaultSeed } from '../../services/adminService';
import { RefreshCw, Users, FileText, LayoutGrid, Eye } from 'lucide-react';
import { Button } from '../../components/ui/Button';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, Legend } from 'recharts';
import { useToast } from '../../components/ui/Toast';

export const AdminHome: React.FC = () => {
  const { t } = useLanguage();
  const { showToast } = useToast();
  const [stats, setStats] = useState<AdminStats | null>(null);
  const [loading, setLoading] = useState(true);
  const [importingSeed, setImportingSeed] = useState(false);

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

  const Chart: React.FC<{ dates: string[]; a: number[]; b: number[]; c?: number[]; d?: number[]; la: string; lb: string; lc?: string; ld?: string }> = ({ dates, a, b, c, d, la, lb, lc, ld }) => {
    const data = dates.map((date, i) => ({
      date,
      users: a[i] ?? 0,
      resumes: b[i] ?? 0,
      templates: c ? (c[i] ?? 0) : undefined,
      visitors: d ? (d[i] ?? 0) : undefined,
    }));
    return (
      <ResponsiveContainer width="100%" height="100%">
        <LineChart data={data}>
          <CartesianGrid vertical={true} horizontal={true} strokeDasharray="3 3" stroke="#f0f0f0" />
          <XAxis
            dataKey="date"
            axisLine={false}
            tickLine={false}
            tick={{ fill: '#9ca3af', fontSize: 13 }}
            dy={10}
          />
          <YAxis
            axisLine={false}
            tickLine={false}
            tick={{ fill: '#9ca3af', fontSize: 13 }}
            dx={-5}
            domain={[0, 'auto']}
            allowDecimals={false}
          />
          <Tooltip
            contentStyle={{ borderRadius: '8px', border: 'none', boxShadow: '0 4px 12px rgba(0,0,0,0.1)' }}
          />
          <Legend
            verticalAlign="bottom"
            height={36}
            iconType="circle"
            wrapperStyle={{ paddingTop: '20px' }}
          />
          <Line
            name={la}
            type="monotone"
            dataKey="users"
            stroke="#3b82f6"
            strokeWidth={1.5}
            dot={{ r: 3, fill: '#fff', strokeWidth: 1.5, stroke: '#3b82f6' }}
            activeDot={{ r: 5 }}
          />
          <Line
            name={lb}
            type="monotone"
            dataKey="resumes"
            stroke="#22c55e"
            strokeWidth={1.5}
            dot={{ r: 3, fill: '#fff', strokeWidth: 1.5, stroke: '#22c55e' }}
            activeDot={{ r: 5 }}
          />
          {c && lc && (
            <Line
              name={lc}
              type="monotone"
              dataKey="templates"
              stroke="#a855f7"
              strokeWidth={1.5}
              dot={{ r: 3, fill: '#fff', strokeWidth: 1.5, stroke: '#a855f7' }}
              activeDot={{ r: 5 }}
            />
          )}
          {d && ld && (
            <Line
              name={ld}
              type="monotone"
              dataKey="visitors"
              stroke="#f59e0b"
              strokeWidth={1.5}
              dot={{ r: 3, fill: '#fff', strokeWidth: 1.5, stroke: '#f59e0b' }}
              activeDot={{ r: 5 }}
            />
          )}
        </LineChart>
      </ResponsiveContainer>
    );
  };

  return (
    <div className="flex-1 flex flex-col bg-white rounded-3xl m-2 overflow-hidden shadow-sm border border-gray-100">
      <div className="px-10 pt-10 pb-6">
        <div className="flex items-center justify-end">
          <div className="flex items-center space-x-3">
            <Button
              variant="outline"
              onClick={async () => {
                if (importingSeed) return;
                setImportingSeed(true);
                try {
                  const r = await importDefaultSeed();
                  if (r?.success) {
                    showToast(t('admin.msg.seedImported'), 'success');
                  } else {
                    showToast(t('admin.msg.seedImportFailed'), 'error');
                  }
                } catch {
                  showToast(t('admin.msg.seedImportFailed'), 'error');
                } finally {
                  setImportingSeed(false);
                }
              }}
              disabled={importingSeed}
            >
              {t('admin.actions.importDefaultSeed')}
            </Button>
            <Button variant="outline" onClick={reload} disabled={loading}>
              <RefreshCw size={16} className={`${loading ? 'animate-spin' : ''} mr-2`} /> {t('common.refresh') || 'Refresh'}
            </Button>
          </div>
        </div>
      </div>
      <div className="flex-1 overflow-y-auto px-10 pb-10">
        <div className="flex flex-col lg:flex-row gap-6 items-stretch">
          <div className="bg-white rounded-2xl border border-gray-100 p-6 flex-1 shadow-sm">
            <div className="flex justify-between items-center mb-6">
              <h3 className="text-gray-800 font-semibold text-lg">{t('admin.stats.trend')}</h3>
              {stats && <span className="text-gray-400 text-sm">{t('admin.stats.generated')} {new Date(stats.generatedAt * 1000).toLocaleString()}</span>}
            </div>
            <div className="h-[350px] w-full flex items-center justify-center">
              {stats ? (
                <Chart dates={stats.trend.dates} a={stats.trend.users} b={stats.trend.resumes} c={stats.trend.templates} d={stats.trend.visitors} la={t('admin.stats.users')} lb={t('admin.stats.resumes')} lc={t('admin.stats.templates')} ld={t('admin.stats.visitors')} />
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
              <div className="flex items-center space-x-4">
                <div className="w-11 h-11 bg-purple-100 text-purple-500 rounded-full flex items-center justify-center">
                  <LayoutGrid size={20} />
                </div>
                <div>
                  <p className="text-xl font-bold text-gray-800 leading-none mb-1">{stats?.totals.templates ?? '-'}</p>
                  <p className="text-gray-500 text-sm">{t('admin.stats.templates')}</p>
                </div>
              </div>
              <div className="flex items-center space-x-4">
                <div className="w-11 h-11 bg-yellow-100 text-yellow-600 rounded-full flex items-center justify-center">
                  <Eye size={20} />
                </div>
                <div>
                  <p className="text-xl font-bold text-gray-800 leading-none mb-1">{stats?.totals.visitorsToday ?? '-'}</p>
                  <p className="text-gray-500 text-sm">{t('admin.stats.visitors')}</p>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div className="mt-8 flex-1"></div>
      </div>
    </div>
  );
};
