import React, { useEffect, useState } from 'react';
import { listShares, AdminShare, updateShare, deleteShare } from '../../services/adminService';
import { Button } from '../../components/ui/Button';
import { useToast } from '../../components/ui/Toast';
import { useConfirm } from '../../components/ui/ConfirmDialog';
import { useLanguage } from '../../contexts/LanguageContext';
import { TableCard } from '../../components/ui/TableCard';
import { RefreshCw } from 'lucide-react';

export const SharesPage: React.FC = () => {
  const [items, setItems] = useState<AdminShare[]>([]);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [total, setTotal] = useState(0);
  const [keyword, setKeyword] = useState('');
  const { showToast } = useToast();
  const confirm = useConfirm();
  const { t } = useLanguage();

  const [loading, setLoading] = useState(false);
  const load = async () => {
    setLoading(true);
    try {
      const resp = await listShares({ page: String(page), pageSize: String(pageSize), slug: keyword });
      setItems(resp.items);
      setTotal(resp.total);
    } catch {
      showToast(t('admin.msg.loadSharesFailed'), 'error');
    } finally { setLoading(false); }
  };
  useEffect(() => { load(); }, [page, pageSize]);

  return (
    <div className="flex-1 flex flex-col bg-white rounded-3xl m-2 overflow-hidden shadow-sm border border-gray-100">
      <div className="px-10 pt-10 pb-6">
        <div className="flex justify-between items-center">
          <h1 className="text-4xl font-bold text-gray-800 tracking-tight">{t('admin.menu.shares')}</h1>
          <Button variant="outline" onClick={() => load()} disabled={loading}>
            <RefreshCw size={16} className={`${loading ? 'animate-spin' : ''} mr-2`} /> {t('common.refresh') || 'Refresh'}
          </Button>
        </div>
        <div className="flex items-center mt-6 pb-4 border-b border-gray-100 space-x-2">
          <input value={keyword} onChange={e => setKeyword(e.target.value)} placeholder={t('admin.slugKeyword')} className="border rounded-lg px-3 py-2 text-sm shadow-sm border-gray-200" />
          <Button onClick={() => { setPage(1); load(); }}>{t('admin.search')}</Button>
        </div>
      </div>
      <div className="flex-1 overflow-y-auto px-10 pb-10">
        <TableCard>
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200 text-sm">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-4 py-2 text-left font-semibold text-gray-600">{t('admin.columns.slug')}</th>
                <th className="px-4 py-2 text-left font-semibold text-gray-600">{t('admin.columns.user')}</th>
                <th className="px-4 py-2 text-left font-semibold text-gray-600">{t('admin.columns.url')}</th>
                <th className="px-4 py-2 text-left font-semibold text-gray-600">{t('admin.columns.status')}</th>
                <th className="px-4 py-2 text-right font-semibold text-gray-600">{t('admin.columns.actions')}</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {items.map(s => (
                <tr key={s.slug}>
                  <td className="px-4 py-2 text-sm text-gray-700">{s.slug}</td>
                  <td className="px-4 py-2 text-sm text-gray-700">{s.userName || (s.userId ? `#${s.userId}` : '')}</td>
                  <td className="px-4 py-2 text-sm text-gray-700">
                    {(() => {
                      const rel = s.url ? (s.url.startsWith('/#/') ? s.url.slice(2) : s.url) : `#/public/${s.slug}`;
                      const uiUrl = `${window.location.origin}${window.location.pathname}${rel.startsWith('#') ? rel : `#${rel}`}`;
                      return (
                        <a href={uiUrl} target="_blank" rel="noreferrer" className="text-blue-600 hover:underline">
                          {uiUrl}
                        </a>
                      );
                    })()}
                  </td>
                  <td className="px-4 py-2 text-sm text-gray-700">{s.isPublic ? t('admin.actions.setPublic') : t('admin.actions.setPrivate')}</td>
                  <td className="px-4 py-2 text-sm text-right space-x-2">
                    <Button variant="outline" onClick={async () => {
                      try { await updateShare(s.slug, !s.isPublic); showToast(t('admin.msg.visibilityUpdated'), 'success'); load(); }
                      catch { showToast(t('admin.msg.visibilityUpdateFailed'), 'error'); }
                    }}>{s.isPublic ? t('admin.actions.setPrivate') : t('admin.actions.setPublic')}</Button>
                    <Button variant="danger" onClick={async () => {
                      const ok = await confirm({ title: t('common.confirmAction'), message: t('admin.confirm.deleteShare') });
                      if (!ok) return;
                      try { await deleteShare(s.slug); showToast(t('admin.msg.deleted'), 'success'); load(); }
                      catch { showToast(t('admin.msg.deleteFailed'), 'error'); }
                    }}>{t('admin.actions.delete')}</Button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
          </div>
        </TableCard>
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
