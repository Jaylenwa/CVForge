import React, { useEffect, useState } from 'react';
import { listUsers, AdminUser, updateUser, resetPassword, banUser, unbanUser } from '../../services/adminService';
import { Button } from '../../components/ui/Button';
import { useConfirm } from '../../components/ui/ConfirmDialog';
import { useToast } from '../../components/ui/Toast';
import { useLanguage } from '../../contexts/LanguageContext';
import { RefreshCw } from 'lucide-react';

export const UsersPage: React.FC = () => {
  const [items, setItems] = useState<AdminUser[]>([]);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [total, setTotal] = useState(0);
  const [keyword, setKeyword] = useState('');
  const confirm = useConfirm();
  const { showToast } = useToast();
  const { t } = useLanguage();

  const [loading, setLoading] = useState(false);
  const load = async () => {
    setLoading(true);
    try {
      const resp = await listUsers({ page: String(page), pageSize: String(pageSize), email: keyword, name: keyword });
      setItems(resp.items);
      setTotal(resp.total);
    } catch {
      showToast(t('admin.msg.loadUsersFailed'), 'error');
    } finally { setLoading(false); }
  };
  useEffect(() => { load(); }, [page, pageSize]);

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-4xl font-bold text-gray-800 tracking-tight">{t('admin.menu.users')}</h1>
        <Button variant="outline" onClick={() => load()} disabled={loading}>
          <RefreshCw size={16} className={`${loading ? 'animate-spin' : ''} mr-2`} /> {t('common.refresh') || 'Refresh'}
        </Button>
      </div>
      <div className="bg-white rounded-2xl border border-gray-100 p-6 shadow-sm">
        <div className="flex items-center mb-4 space-x-2">
          <input value={keyword} onChange={e => setKeyword(e.target.value)} placeholder={t('admin.keyword')} className="border rounded-lg px-3 py-2 text-sm shadow-sm border-gray-200" />
          <Button onClick={() => { setPage(1); load(); }}>{t('admin.search')}</Button>
        </div>
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-4 py-2 text-left text-xs font-medium text-gray-500">{t('admin.columns.id')}</th>
                <th className="px-4 py-2 text-left text-xs font-medium text-gray-500">{t('admin.columns.email')}</th>
                <th className="px-4 py-2 text-left text-xs font-medium text-gray-500">{t('admin.columns.name')}</th>
                <th className="px-4 py-2 text-left text-xs font-medium text-gray-500">{t('admin.columns.role')}</th>
                <th className="px-4 py-2 text-left text-xs font-medium text-gray-500">{t('admin.columns.status')}</th>
                <th className="px-4 py-2 text-right text-xs font-medium text-gray-500">{t('admin.columns.actions')}</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {items.map(u => (
                <tr key={u.id}>
                  <td className="px-4 py-2 text-sm text-gray-700">{u.id}</td>
                  <td className="px-4 py-2 text-sm text-gray-700">{u.email}</td>
                  <td className="px-4 py-2 text-sm text-gray-700">{u.name}</td>
                  <td className="px-4 py-2 text-sm">
                    <select className="border rounded-md px-2 py-1 text-sm" value={u.role} onChange={async e => { 
                      const role = e.target.value; 
                      try { await updateUser(u.id, { role }); showToast(t('admin.msg.roleUpdated'), 'success'); load(); } 
                      catch { showToast(t('admin.msg.roleUpdateFailed'), 'error'); } 
                    }}>
                      <option value="user">user</option>
                      <option value="admin">admin</option>
                    </select>
                  </td>
                  <td className="px-4 py-2 text-sm text-gray-700">{u.isActive ? t('admin.status.active') : t('admin.status.inactive')}</td>
                  <td className="px-4 py-2 text-sm text-right space-x-2">
                    {u.isActive ? (
                      <Button variant="outline" onClick={async () => { 
                        const ok = await confirm({ title: t('common.confirmAction'), message: t('admin.confirm.ban') }); 
                        if (!ok) return; 
                        try { await banUser(u.id); showToast(t('admin.msg.banned'), 'success'); load(); } 
                        catch { showToast(t('admin.msg.banFailed'), 'error'); } 
                      }}>{t('admin.actions.ban')}</Button>
                    ) : (
                      <Button variant="outline" onClick={async () => { 
                        try { await unbanUser(u.id); showToast(t('admin.msg.unbanned'), 'success'); load(); } 
                        catch { showToast(t('admin.msg.unbanFailed'), 'error'); } 
                      }}>{t('admin.actions.unban')}</Button>
                    )}
                    <Button variant="ghost" onClick={async () => {
                      const ok = await confirm({ title: t('common.confirmAction'), message: t('admin.confirm.resetPassword') });
                      if (!ok) return;
                      const newPwd = Math.random().toString(36).slice(2, 10);
                      try { await resetPassword(u.id, newPwd); showToast(`${t('admin.msg.newPassword')} ${newPwd}`, 'success'); } 
                      catch { showToast(t('admin.msg.resetFailed'), 'error'); }
                    }}>{t('admin.actions.resetPassword')}</Button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
        <div className="flex justify-between items-center mt-4">
          <div className="text-sm text-gray-500">{t('admin.total')} {total}</div>
          <div className="space-x-2">
            <Button variant="outline" disabled={page === 1} onClick={() => setPage(p => Math.max(p - 1, 1))}>{t('admin.prev')}</Button>
            <Button variant="outline" onClick={() => setPage(p => p + 1)} disabled={page * pageSize >= total}>{t('admin.next')}</Button>
            <select className="border rounded-md px-2 py-1 text-sm" value={pageSize} onChange={e => setPageSize(parseInt(e.target.value))}>
              <option value={10}>10</option>
              <option value={20}>20</option>
              <option value={50}>50</option>
            </select>
          </div>
        </div>
      </div>
    </div>
  );
};
