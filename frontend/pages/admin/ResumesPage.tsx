import React, { useEffect, useState } from 'react';
import { listResumes, AdminResume, deleteResume, transferResume, setResumeVisibility } from '../../services/adminService';
import { Button } from '../../components/ui/Button';
import { useConfirm } from '../../components/ui/ConfirmDialog';
import { useToast } from '../../components/ui/Toast';
import { useLanguage } from '../../contexts/LanguageContext';

export const ResumesPage: React.FC = () => {
  const [items, setItems] = useState<AdminResume[]>([]);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [total, setTotal] = useState(0);
  const [keyword, setKeyword] = useState('');
  const confirm = useConfirm();
  const { showToast } = useToast();
  const { t } = useLanguage();

  const load = async () => {
    try {
      const resp = await listResumes({ page: String(page), pageSize: String(pageSize), title: keyword });
      setItems(resp.items);
      setTotal(resp.total);
    } catch {
      showToast(t('admin.msg.loadResumesFailed'), 'error');
    }
  };
  useEffect(() => { load(); }, [page, pageSize]);

  return (
    <div>
      <div className="flex items-center mb-4 space-x-2">
        <input value={keyword} onChange={e => setKeyword(e.target.value)} placeholder={t('admin.titleKeyword')} className="border rounded-md px-3 py-2 text-sm" />
        <Button onClick={() => { setPage(1); load(); }}>{t('admin.search')}</Button>
      </div>
      <div className="bg-white shadow-sm rounded-md overflow-hidden">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-4 py-2 text-left text-xs font-medium text-gray-500">{t('admin.columns.id')}</th>
              <th className="px-4 py-2 text-left text-xs font-medium text-gray-500">{t('admin.columns.user')}</th>
              <th className="px-4 py-2 text-left text-xs font-medium text-gray-500">{t('admin.columns.title')}</th>
              <th className="px-4 py-2 text-left text-xs font-medium text-gray-500">{t('admin.columns.template')}</th>
              <th className="px-4 py-2 text-right text-xs font-medium text-gray-500">{t('admin.columns.actions')}</th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {items.map(r => (
              <tr key={r.id}>
                <td className="px-4 py-2 text-sm text-gray-700">{r.id}</td>
                <td className="px-4 py-2 text-sm text-gray-700">{r.userName || r.userId}</td>
                <td className="px-4 py-2 text-sm text-gray-700">{r.title}</td>
                <td className="px-4 py-2 text-sm text-gray-700">{r.templateId}</td>
                <td className="px-4 py-2 text-sm text-right space-x-2">
                  <Button variant="outline" onClick={async () => {
                    const ok = await confirm({ title: t('common.confirmAction'), message: t('admin.confirm.deleteResume') });
                    if (!ok) return;
                    try { await deleteResume(r.id); showToast(t('admin.msg.deleted'), 'success'); load(); }
                    catch { showToast(t('admin.msg.deleteFailed'), 'error'); }
                  }}>{t('admin.actions.delete')}</Button>
                  <Button variant="ghost" onClick={async () => {
                    const input = prompt(t('admin.input.newOwnerId'));
                    if (!input) return;
                    const uid = parseInt(input);
                    if (!uid) return;
                    try { await transferResume(r.id, uid); showToast(t('admin.msg.transferred'), 'success'); load(); }
                    catch { showToast(t('admin.msg.transferFailed'), 'error'); }
                  }}>{t('admin.actions.transferOwner')}</Button>
                  <Button variant="ghost" onClick={async () => {
                    const ok = await confirm({ title: t('common.confirmAction'), message: t('admin.confirm.setPublic') });
                    if (!ok) return;
                    try { await setResumeVisibility(r.id, true); showToast(t('admin.msg.setPublicSuccess'), 'success'); }
                    catch { showToast(t('admin.msg.setPublicFailed'), 'error'); }
                  }}>{t('admin.actions.setPublic')}</Button>
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
  );
};
