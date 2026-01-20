import React, { useEffect, useState } from 'react';
import { listUsers, AdminUser, updateUser, resetPassword, banUser, unbanUser } from '../../services/adminService';
import { Button } from '../../components/ui/Button';
import { useConfirm } from '../../components/ui/ConfirmDialog';
import { useToast } from '../../components/ui/Toast';
import { useLanguage } from '../../contexts/LanguageContext';
import { TableCard } from '../../components/ui/TableCard';
import { Input, Select } from '../../components/ui/Form';
import { Search, Mail, Lock, Ban, ShieldAlert } from 'lucide-react';

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
    <div className="flex-1 flex flex-col bg-white rounded-3xl m-2 overflow-hidden shadow-sm border border-gray-100">
      <div className="px-10 pt-10 pb-6">
        <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4 border-b border-gray-100 pb-4">
          <div className="relative flex-1">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
            <Input
              value={keyword}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) => setKeyword(e.target.value)}
              placeholder={t('admin.keyword')}
              className="pl-10"
            />
          </div>
          <div className="flex items-center gap-2 md:justify-end">
            <Select
              value={statusFilter}
              onChange={(e: React.ChangeEvent<HTMLSelectElement>) => {
                setPage(1);
                setStatusFilter(e.target.value as any);
              }}
              options={[
                { label: t('admin.filter.allStatuses'), value: 'all' },
                { label: t('admin.status.active'), value: 'true' },
                { label: t('admin.status.inactive'), value: 'false' },
              ]}
            />
            
          </div>
        </div>
      </div>
      <div className="flex-1 overflow-y-auto px-10 pb-10">
        <TableCard>
          <div className="overflow-x-auto no-scrollbar">
            <table className="w-full text-left text-sm">
            <thead>
              <tr className="bg-gray-50">
                <th className="px-6 py-4 font-semibold text-gray-600">{t('admin.columns.id')}</th>
                <th className="px-6 py-4 font-semibold text-gray-600">{t('admin.columns.userDetails')}</th>
                <th className="px-6 py-4 font-semibold text-gray-600">{t('admin.columns.role')}</th>
                <th className="px-6 py-4 font-semibold text-gray-600">{t('admin.columns.status')}</th>
                <th className="px-6 py-4 font-semibold text-gray-600">{t('admin.columns.lastActivity')}</th>
                <th className="px-6 py-4 font-semibold text-gray-600 text-right">{t('admin.columns.actions')}</th>
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
        </TableCard>
        <div className="mt-6 flex items-center justify-between text-sm text-gray-500">
          <div>
            {t('admin.total')} {total}
          </div>
          <div className="flex items-center gap-2">
            <Button variant="outline" size="sm" disabled={page === 1} onClick={() => setPage((p: number) => Math.max(p - 1, 1))}>
              {t('admin.prev')}
            </Button>
            <div className="min-w-[64px] text-center">
              {page} / {Math.max(1, Math.ceil(total / pageSize))}
            </div>
            <Button variant="outline" size="sm" disabled={page * pageSize >= total} onClick={() => setPage((p: number) => p + 1)}>
              {t('admin.next')}
            </Button>
            <div className="w-[96px]">
              <Select
                value={String(pageSize)}
                onChange={(e: React.ChangeEvent<HTMLSelectElement>) => {
                  setPage(1);
                  setPageSize(parseInt(e.target.value, 10));
                }}
                options={[10, 20, 50].map((x) => ({ label: String(x), value: String(x) }))}
                className="py-1.5"
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
