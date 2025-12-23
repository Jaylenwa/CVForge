import React, { useEffect, useLayoutEffect, useRef, useState } from 'react';
import { listResumes, AdminResume, deleteResume, setResumeVisibility } from '../../services/adminService';
import { Button } from '../../components/ui/Button';
import { Modal } from '../../components/ui/Modal';
import { useConfirm } from '../../components/ui/ConfirmDialog';
import { useToast } from '../../components/ui/Toast';
import { useLanguage } from '../../contexts/LanguageContext';
import { ResumeArtboard } from '../editor/ResumePreview';
import { INITIAL_RESUME } from '../../services/mockData';
import { RefreshCw, Search, ChevronLeft, ChevronRight, Eye, Trash2, Globe } from 'lucide-react';

export const ResumesPage: React.FC = () => {
  const [items, setItems] = useState<AdminResume[]>([]);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [total, setTotal] = useState(0);
  const [keyword, setKeyword] = useState('');
  const confirm = useConfirm();
  const { showToast } = useToast();
  const { t } = useLanguage();

  const [loading, setLoading] = useState(false);
  const mmToPx = 96 / 25.4;
  const a4w = 210 * mmToPx;
  const thumbnailWidth = 40;
  const thumbnailScale = thumbnailWidth / a4w;

  const [showPreview, setShowPreview] = useState(false);
  const [previewResume, setPreviewResume] = useState<AdminResume | null>(null);
  const previewContainerRef = useRef<HTMLDivElement | null>(null);
  const [previewScale, setPreviewScale] = useState<number | null>(null);
  const rafRef = useRef<number | null>(null);
  const roRef = useRef<ResizeObserver | null>(null);

  const load = async () => {
    setLoading(true);
    try {
      const resp = await listResumes({ page: String(page), pageSize: String(pageSize), title: keyword });
      setItems(resp.items);
      setTotal(resp.total);
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

  useLayoutEffect(() => {
    if (!showPreview) return;
    const schedule = () => {
      if (rafRef.current) cancelAnimationFrame(rafRef.current);
      rafRef.current = requestAnimationFrame(() => {
        const el = previewContainerRef.current;
        if (!el) return;
        const s = el.clientWidth / a4w;
        setPreviewScale(s);
      });
    };
    schedule();
    const onResize = () => schedule();
    window.addEventListener('resize', onResize);
    if (previewContainerRef.current) {
      roRef.current = new ResizeObserver(onResize);
      roRef.current.observe(previewContainerRef.current);
    }
    return () => {
      window.removeEventListener('resize', onResize);
      if (rafRef.current) cancelAnimationFrame(rafRef.current);
      if (roRef.current) roRef.current.disconnect();
    };
  }, [showPreview]);

  const Thumbnail: React.FC<{ r: AdminResume }> = ({ r }) => {
    return (
      <div className="w-10 h-12 bg-slate-100 rounded-md overflow-hidden flex-shrink-0 border border-slate-200 relative">
        <ResumeArtboard
          data={{ ...INITIAL_RESUME, templateId: r.templateId, themeConfig: r.themeConfig }}
          scale={thumbnailScale}
          disableShadow
          style={{ margin: 0 }}
          className="absolute top-0 left-0"
          showPageHint={false}
        />
      </div>
    );
  };

  return (
    <div className="min-h-screen p-4 md:p-8">
      <div className="max-w-7xl mx-auto mb-8 flex justify-between items-center">
        <h1 className="text-4xl font-bold text-gray-800 tracking-tight">{t('admin.menu.resumes')}</h1>
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
                          <div className="text-[10px] text-slate-400 uppercase tracking-tight">{r.themeConfig?.fontFamily || ''}</div>
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
                            onClick={() => { setPreviewResume(r); setShowPreview(true); }}
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

      <Modal isOpen={showPreview} onClose={() => setShowPreview(false)} title={t('admin.actions.preview')}>
        <div ref={previewContainerRef} className="aspect-[210/297] bg-gray-100 overflow-hidden relative">
          <div className="absolute inset-0 flex items-center justify-center">
            {previewResume && previewScale !== null ? (
              <div
                style={{ width: (96 / 25.4) * 210 * previewScale, height: (96 / 25.4) * 297 * previewScale }}
                className="relative select-none pointer-events-none shadow-sm bg-white"
              >
                <ResumeArtboard
                  data={{ ...INITIAL_RESUME, templateId: previewResume.templateId, themeConfig: previewResume.themeConfig }}
                  scale={previewScale}
                  disableShadow
                  style={{ margin: 0 }}
                />
              </div>
            ) : (
              <div className="w-full h-full bg-white" />
            )}
          </div>
        </div>
      </Modal>
    </div>
  );
};
