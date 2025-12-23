import React, { useEffect, useState } from 'react';
import { listUsers, AdminUser, updateUser, resetPassword, banUser, unbanUser } from '../../services/adminService';
import { Button } from '../../components/ui/Button';
import { useConfirm } from '../../components/ui/ConfirmDialog';
import { useToast } from '../../components/ui/Toast';
import { useLanguage } from '../../contexts/LanguageContext';
import { Search, RefreshCw, Mail, ChevronLeft, ChevronRight, Lock, Ban, ShieldAlert } from 'lucide-react';

export const UsersPage: React.FC = () => {
  const [items, setItems] = useState<AdminUser[]>([]);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [total, setTotal] = useState(0);
  const [keyword, setKeyword] = useState('');
  const [statusFilter, setStatusFilter] = useState<'all' | 'true' | 'false'>('all');
  const confirm = useConfirm();
  const { showToast } = useToast();
  const { t } = useLanguage();

  const [loading, setLoading] = useState(false);
  const load = async () => {
    setLoading(true);
    try {
      const params: Record<string, string> = { page: String(page), pageSize: String(pageSize), email: keyword, name: keyword };
      if (statusFilter !== 'all') params.isActive = statusFilter;
      const resp = await listUsers(params);
      setItems(resp.items);
      setTotal(resp.total);
    } catch {
      showToast(t('admin.msg.loadUsersFailed'), 'error');
    } finally { setLoading(false); }
  };
  useEffect(() => { load(); }, [page, pageSize, statusFilter]);

  

  useEffect(() => {
    const timer = setTimeout(() => {
      if (page !== 1) {
        setPage(1);
      } else {
        load();
      }
    }, 300);
    return () => clearTimeout(timer);
  }, [keyword]);

  const StatusBadge: React.FC<{ active: boolean }> = ({ active }) => {
    const cls = active ? 'bg-emerald-50 text-emerald-700 border-emerald-200' : 'bg-rose-50 text-rose-700 border-rose-200';
    const dot = active ? 'bg-emerald-500' : 'bg-rose-500';
    return (
      <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium border ${cls}`}>
        <span className={`w-1.5 h-1.5 rounded-full mr-1.5 ${dot}`}></span>
        {active ? t('admin.status.active') : t('admin.status.inactive')}
      </span>
    );
  };

  return (
    <div className="min-h-screen p-4 md:p-8">
      <div className="max-w-7xl mx-auto mb-8 flex justify-between items-center">
        <h1 className="text-4xl font-bold text-gray-800 tracking-tight">{t('admin.menu.users')}</h1>
        <Button variant="outline" onClick={() => load()} disabled={loading}>
          <RefreshCw size={16} className={`${loading ? 'animate-spin' : ''} mr-2`} /> {t('common.refresh') || 'Refresh'}
        </Button>
      </div>

      <div className="max-w-7xl mx-auto">
        <div className="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
            <div className="p-4 border-b border-gray-100 bg-gray-50/50 flex flex-col md:flex-row md:items-center gap-4">
              <div className="relative flex-1">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
                <input
                  type="text"
                  placeholder={t('admin.keyword')}
                  className="w-full pl-10 pr-4 py-2 bg-white border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all"
                  value={keyword}
                  onChange={(e: React.ChangeEvent<HTMLInputElement>) => setKeyword(e.target.value)}
                  onKeyDown={(e: React.KeyboardEvent<HTMLInputElement>) => {
                    if (e.key === 'Enter') { setPage(1); load(); }
                  }}
                />
              </div>
              <div className="flex items-center gap-2">
                <select
                  className="bg-white border border-gray-200 rounded-lg px-3 py-2 text-sm text-gray-600 focus:outline-none focus:ring-2 focus:ring-blue-500/20"
                  value={statusFilter}
                  onChange={(e: React.ChangeEvent<HTMLSelectElement>) => { setPage(1); setStatusFilter(e.target.value as any); }}
                >
                  <option value="all">{t('admin.filter.allStatuses')}</option>
                  <option value="true">{t('admin.status.active')}</option>
                  <option value="false">{t('admin.status.inactive')}</option>
                </select>
              </div>
            </div>

            <div className="overflow-x-auto">
              <table className="w-full text-left">
                <thead>
                  <tr className="bg-gray-50 text-gray-500 text-xs font-semibold uppercase tracking-wider">
                    <th className="px-6 py-4">{t('admin.columns.id')}</th>
                    <th className="px-6 py-4">{t('admin.columns.userDetails')}</th>
                    <th className="px-6 py-4">{t('admin.columns.role')}</th>
                    <th className="px-6 py-4">{t('admin.columns.status')}</th>
                    <th className="px-6 py-4">{t('admin.columns.lastActivity')}</th>
                    <th className="px-6 py-4 text-right">{t('admin.columns.actions')}</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-100">
                  {items.map((u: AdminUser) => (
                    <tr key={u.id} className="hover:bg-blue-50/30 transition-colors">
                      <td className="px-6 py-4 text-sm text-gray-500 font-mono">#{u.id}</td>
                      <td className="px-6 py-4">
                        <div className="flex items-center gap-3">
                          <div className="w-10 h-10 rounded-full bg-gradient-to-br from-blue-100 to-indigo-100 flex items-center justify-center text-blue-600 font-semibold border border-blue-200 shadow-sm">
                            {(u.name || u.email || '').charAt(0).toUpperCase()}
                          </div>
                          <div>
                            <div className="text-sm font-semibold text-gray-900">{u.name || '-'}</div>
                            <div className="text-xs text-gray-500 flex items-center gap-1">
                              <Mail className="w-3 h-3" />
                              {u.email}
                            </div>
                          </div>
                        </div>
                      </td>
                      <td className="px-6 py-4">
                        <select
                          className="bg-transparent border-0 text-sm font-medium text-gray-700 focus:ring-0 cursor-pointer hover:text-blue-600 transition-colors"
                          value={u.role}
                          onChange={async (e: React.ChangeEvent<HTMLSelectElement>) => {
                            const role = e.target.value;
                            try { await updateUser(u.id, { role }); showToast(t('admin.msg.roleUpdated'), 'success'); load(); }
                            catch { showToast(t('admin.msg.roleUpdateFailed'), 'error'); }
                          }}
                        >
                          <option value="user">user</option>
                          <option value="admin">admin</option>
                        </select>
                      </td>
                      <td className="px-6 py-4">
                        <StatusBadge active={u.isActive} />
                      </td>
                      <td className="px-6 py-4 text-sm text-gray-500">{u.lastLoginAt ? u.lastLoginAt : '-'}</td>
                      <td className="px-6 py-4 text-right">
                        <div className="flex items-center justify-end gap-2">
                          <div className="relative group">
                            <button
                            onClick={async () => {
                              const ok = await confirm({ title: t('common.confirmAction'), message: t('admin.confirm.resetPassword') });
                              if (!ok) return;
                              const newPwd = Math.random().toString(36).slice(2, 10);
                              try { await resetPassword(u.id, newPwd); showToast(`${t('admin.msg.newPassword')} ${newPwd}`, 'success'); }
                              catch { showToast(t('admin.msg.resetFailed'), 'error'); }
                            }}
                            className="p-2 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-colors"
                          >
                            <Lock className="w-4 h-4" />
                            </button>
                            <div className="pointer-events-none absolute bottom-full left-1/2 -translate-x-1/2 mb-1 hidden group-hover:block whitespace-nowrap rounded-md bg-gray-900 text-white text-xs px-2 py-1 shadow">
                              {t('admin.actions.resetPassword')}
                            </div>
                          </div>
                          <div className="relative group">
                            <button
                            onClick={async () => {
                              if (u.isActive) {
                                const ok = await confirm({ title: t('common.confirmAction'), message: t('admin.confirm.ban') });
                                if (!ok) return;
                                try { await banUser(u.id); showToast(t('admin.msg.banned'), 'success'); load(); }
                                catch { showToast(t('admin.msg.banFailed'), 'error'); }
                              } else {
                                try { await unbanUser(u.id); showToast(t('admin.msg.unbanned'), 'success'); load(); }
                                catch { showToast(t('admin.msg.unbanFailed'), 'error'); }
                              }
                            }}
                            className={`p-2 rounded-lg transition-colors ${!u.isActive ? 'text-emerald-600 hover:bg-emerald-50' : 'text-rose-600 hover:bg-rose-50'}`}
                          >
                            {u.isActive ? <Ban className="w-4 h-4" /> : <ShieldAlert className="w-4 h-4" />}
                            </button>
                            <div className="pointer-events-none absolute bottom-full left-1/2 -translate-x-1/2 mb-1 hidden group-hover:block whitespace-nowrap rounded-md bg-gray-900 text-white text-xs px-2 py-1 shadow">
                              {u.isActive ? t('admin.actions.ban') : t('admin.actions.unban')}
                            </div>
                          </div>
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            <div className="p-4 border-t border-gray-100 flex items-center justify-between text-sm text-gray-500">
              <div>
                {t('admin.total')} {total}
              </div>
              <div className="flex items-center gap-2">
                <button
                  onClick={() => setPage((p: number) => Math.max(p - 1, 1))}
                  disabled={page === 1}
                  className="px-3 py-1.5 border border-gray-200 rounded-lg hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-1 transition-colors"
                >
                  <ChevronLeft className="w-4 h-4" /> {t('admin.prev')}
                </button>
                <div className="flex items-center gap-1 px-4">
                  <span className="w-8 h-8 flex items-center justify-center bg-blue-600 text-white rounded-lg font-medium">{page}</span>
                </div>
                <button
                  onClick={() => setPage((p: number) => p + 1)}
                  disabled={page * pageSize >= total}
                  className="px-3 py-1.5 border border-gray-200 rounded-lg hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-1 transition-colors"
                >
                  {t('admin.next')} <ChevronRight className="w-4 h-4" />
                </button>
                <select className="border border-gray-200 rounded-lg px-3 py-1.5 text-sm text-gray-600" value={pageSize} onChange={(e: React.ChangeEvent<HTMLSelectElement>) => { setPage(1); setPageSize(parseInt(e.target.value)); }}>
                  <option value={10}>10</option>
                  <option value={20}>20</option>
                  <option value={50}>50</option>
                </select>
              </div>
            </div>
          </div>
      </div>
    </div>
  );
};
