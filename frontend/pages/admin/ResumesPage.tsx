import React, { useEffect, useState } from 'react';
import { listResumes, AdminResume, deleteResume, setResumeVisibility } from '../../services/adminService';
import { Button } from '../../components/ui/Button';
import { useConfirm } from '../../components/ui/ConfirmDialog';
import { useToast } from '../../components/ui/Toast';
import { useLanguage } from '../../contexts/LanguageContext';
import { ResumeArtboard } from '../editor/ResumePreview';
import { LanguageProvider } from '../../contexts/LanguageContext';
import { INITIAL_RESUME } from '../../services/mockData';
import { RefreshCw, Search, ChevronLeft, ChevronRight, Eye, Trash2, Globe } from 'lucide-react';
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
        <div className="flex justify-between items-center">
          <h1 className="text-4xl font-bold text-gray-800 tracking-tight">{t('admin.menu.resumes')}</h1>
          <Button variant="outline" onClick={() => load()} disabled={loading}>
            <RefreshCw size={16} className={`${loading ? 'animate-spin' : ''} mr-2`} /> {t('common.refresh') || 'Refresh'}
          </Button>
        </div>
        <div className="flex flex-col md:flex-row md:items-center gap-4 border-b border-gray-100 mt-6 pb-4">
          <div className="relative flex-1">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
            <input
              type="text"
              placeholder={t('admin.titleKeyword')}
              className="w-full pl-10 pr-4 py-2 bg-white border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all"
              value={keyword}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) => setKeyword(e.target.value)}
              onKeyDown={(e: React.KeyboardEvent<HTMLInputElement>) => {
                if (e.key === 'Enter') { setPage(1); load(); }
              }}
            />
          </div>
          <div className="flex items-center gap-2">
            <Button onClick={() => { setPage(1); load(); }}>{t('admin.search')}</Button>
          </div>
        </div>
      </div>

      <div className="flex-1 overflow-y-auto px-10 pb-10">
        <div className="overflow-x-auto">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="bg-slate-50/80 border-b border-slate-200">
                <th className="px-6 py-4 text-xs font-bold text-slate-500 uppercase tracking-wider">ID</th>
                <th className="px-6 py-4 text-xs font-bold text-slate-500 uppercase tracking-wider">{t('admin.columns.title')}</th>
                <th className="px-6 py-4 text-xs font-bold text-slate-500 uppercase tracking-wider">{t('admin.columns.user')}</th>
                <th className="px-6 py-4 text-xs font-bold text-slate-500 uppercase tracking-wider">{t('admin.columns.template')}</th>
                <th className="pl-6 pr-[88px] py-4 text-xs font-bold text-slate-500 uppercase tracking-wider text-right align-middle">{t('admin.columns.actions')}</th>
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

        <div className="mt-6 flex items-center justify-between text-sm text-gray-500">
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
              <span className="text-xs text-slate-500">/ {totalPages}</span>
            </div>
            <button
              onClick={() => setPage((p: number) => p + 1)}
              disabled={!hasNext}
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
  );
};
