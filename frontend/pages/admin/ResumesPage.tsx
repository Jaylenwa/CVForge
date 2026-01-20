import React, { useEffect, useState } from 'react';
import { listResumes, AdminResume, deleteResume, setResumeVisibility } from '../../services/adminService';
import { Button } from '../../components/ui/Button';
import { useConfirm } from '../../components/ui/ConfirmDialog';
import { useToast } from '../../components/ui/Toast';
import { useLanguage } from '../../contexts/LanguageContext';
import { TableCard } from '../../components/ui/TableCard';
import { Input, Select } from '../../components/ui/Form';
import { ResumeArtboard } from '../editor/ResumePreview';
import { LanguageProvider } from '../../contexts/LanguageContext';
import { INITIAL_RESUME } from '../../services/mockData';
import { Search, Eye, Trash2, Globe } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import { AppRoute } from '../../types';

export const ResumesPage: React.FC = () => {
  const navigate = useNavigate();
  const [items, setItems] = useState<AdminResume[]>([]);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [total, setTotal] = useState(0);
  const [totalPages, setTotalPages] = useState<number>(0);
  const [hasNext, setHasNext] = useState<boolean>(false);
  const [keyword, setKeyword] = useState('');
  const confirm = useConfirm();
  const { showToast } = useToast();
  const { t } = useLanguage();

  const [loading, setLoading] = useState(false);
  const thumbnailWidth = 40;
  const thumbnailScale = thumbnailWidth / (210 * (96 / 25.4));

  

  const load = async () => {
    setLoading(true);
    try {
      const resp = await listResumes({ page: String(page), pageSize: String(pageSize), title: keyword });
      setItems(resp.items);
      setTotal(resp.total);
      setTotalPages(resp.totalPages || Math.ceil(resp.total / resp.pageSize));
      setHasNext(!!resp.hasNext && page < (resp.totalPages || Math.ceil(resp.total / resp.pageSize)));
    } catch {
      showToast(t('admin.msg.loadResumesFailed'), 'error');
    } finally { setLoading(false); }
  };
  useEffect(() => { load(); }, [page, pageSize]);

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

  

  const Thumbnail: React.FC<{ r: AdminResume }> = ({ r }) => {
    return (
      <div className="w-10 h-12 bg-slate-100 rounded-md overflow-hidden flex-shrink-0 border border-slate-200 relative">
        <LanguageProvider languageOverride={r.language || 'zh'}>
          <ResumeArtboard
            data={{ ...INITIAL_RESUME, templateId: r.templateId, Theme: r.Theme, language: r.language || 'zh' }}
            scale={thumbnailScale}
            disableShadow
            style={{ margin: 0 }}
            className="absolute top-0 left-0"
            showPageHint={false}
          />
        </LanguageProvider>
      </div>
    );
  };

  return (
    <div className="flex-1 flex flex-col bg-white rounded-3xl m-2 overflow-hidden shadow-sm border border-gray-100">
      <div className="px-10 pt-10 pb-6">
        <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4 border-b border-gray-100 pb-4">
          <div className="relative flex-1">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
            <Input value={keyword} onChange={(e: React.ChangeEvent<HTMLInputElement>) => setKeyword(e.target.value)} placeholder={t('admin.titleKeyword')} className="pl-10" />
          </div>
          <div className="flex items-center gap-2 md:justify-end"></div>
        </div>
      </div>

      <div className="flex-1 overflow-y-auto px-10 pb-10">
        <TableCard>
          <div className="overflow-x-auto no-scrollbar">
            <table className="w-full text-left border-collapse text-sm">
            <thead>
              <tr className="bg-slate-50/80 border-b border-slate-200">
                <th className="px-6 py-4 font-semibold text-gray-600">ID</th>
                <th className="px-6 py-4 font-semibold text-gray-600">{t('admin.columns.title')}</th>
                <th className="px-6 py-4 font-semibold text-gray-600">{t('admin.columns.user')}</th>
                <th className="px-6 py-4 font-semibold text-gray-600">{t('admin.columns.template')}</th>
                <th className="pl-6 pr-[88px] py-4 font-semibold text-gray-600 text-right align-middle">{t('admin.columns.actions')}</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100">
              {items.map((r: AdminResume) => (
                <tr key={r.id} className="hover:bg-indigo-50/30 transition-colors">
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-slate-400">#{r.id}</td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="flex items-center gap-3">
                      <Thumbnail r={r} />
                      <div>
                        <div className="text-sm font-semibold text-slate-900">{r.title}</div>
                        <div className="text-[10px] text-slate-400 uppercase tracking-tight">{r.Theme?.Font || ''}</div>
                      </div>
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="text-sm font-semibold text-slate-900">{r.userName || r.userId}</div>
                    <div className="text-[10px] text-slate-400 uppercase tracking-tight">{new Date(r.lastModified).toLocaleString()}</div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-slate-700">{r.templateId}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-right align-middle">
                    <div className="flex items-center justify-end gap-0.5">
                      <div className="relative group">
                        <button
                          onClick={() => window.open(`${window.location.origin}${window.location.pathname}#${AppRoute.Print}?id=${r.id}`, '_blank')}
                          className="p-1.5 leading-none text-slate-400 hover:text-indigo-600 hover:bg-indigo-50 rounded-md transition-all"
                        >
                          <Eye className="w-4 h-4" />
                        </button>
                        <div className="pointer-events-none absolute bottom-full left-1/2 -translate-x-1/2 mb-1 hidden group-hover:block whitespace-nowrap rounded-md bg-gray-900 text-white text-xs px-2 py-1 shadow">
                          {t('admin.actions.preview')}
                        </div>
                      </div>
                      <div className="relative group">
                        <button
                          onClick={async () => {
                            const ok = await confirm({ title: t('common.confirmAction'), message: t('admin.confirm.setPublic') });
                            if (!ok) return;
                            try { await setResumeVisibility(r.id, true); showToast(t('admin.msg.setPublicSuccess'), 'success'); }
                            catch { showToast(t('admin.msg.setPublicFailed'), 'error'); }
                          }}
                          className="p-1.5 leading-none text-emerald-600 hover:text-emerald-700 hover:bg-emerald-50 rounded-md transition-all"
                        >
                          <Globe className="w-4 h-4" />
                        </button>
                        <div className="pointer-events-none absolute bottom-full left-1/2 -translate-x-1/2 mb-1 hidden group-hover:block whitespace-nowrap rounded-md bg-gray-900 text-white text-xs px-2 py-1 shadow">
                          {t('admin.actions.setPublic')}
                        </div>
                      </div>
                      <div className="relative group">
                        <button
                          onClick={async () => {
                            const ok = await confirm({ title: t('common.confirmAction'), message: t('admin.confirm.deleteResume') });
                            if (!ok) return;
                            try { await deleteResume(r.id); showToast(t('admin.msg.deleted'), 'success'); load(); }
                            catch { showToast(t('admin.msg.deleteFailed'), 'error'); }
                          }}
                          className="p-1.5 leading-none text-rose-600 hover:text-rose-700 hover:bg-rose-50 rounded-md transition-all"
                        >
                          <Trash2 className="w-4 h-4" />
                        </button>
                        <div className="pointer-events-none absolute bottom-full left-1/2 -translate-x-1/2 mb-1 hidden group-hover:block whitespace-nowrap rounded-md bg-gray-900 text-white text-xs px-2 py-1 shadow">
                          {t('admin.actions.delete')}
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
              {page} / {Math.max(1, totalPages || Math.ceil(total / pageSize))}
            </div>
            <Button variant="outline" size="sm" disabled={!hasNext} onClick={() => setPage((p: number) => p + 1)}>
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
