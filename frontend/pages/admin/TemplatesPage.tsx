import React, { useEffect, useState } from 'react';
import { API_BASE } from '../../config';
import { createTemplate, updateTemplate, deleteTemplate } from '../../services/adminService';
import { Button } from '../../components/ui/Button';
import { Modal } from '../../components/ui/Modal';
import { useToast } from '../../components/ui/Toast';
import { useConfirm } from '../../components/ui/ConfirmDialog';
import { useLanguage } from '../../contexts/LanguageContext';
import { ResumeArtboard } from '../editor/ResumePreview';
import { INITIAL_RESUME, MOCK_TEMPLATES } from '../../services/mockData';
import { RefreshCw } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import { AppRoute } from '../../types';

type Row = {
  id: string;
  name: string;
  tags: string[];
  usageCount: number;
  isPremium: boolean;
  category: string;
};

export const TemplatesPage: React.FC = () => {
  const { t } = useLanguage();
  const { showToast } = useToast();
  const confirm = useConfirm();
  const navigate = useNavigate();

  const [items, setItems] = useState<Row[]>([]);
  const [keyword, setKeyword] = useState('');
  const [selectedCategory, setSelectedCategory] = useState('All');
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);

  const categories = ['All', 'IT', 'Finance', 'Creative', 'General'];

  const [showForm, setShowForm] = useState(false);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [form, setForm] = useState({ externalId: '', name: '', tags: '', category: '', usageCount: 0, isPremium: false });

  const [syncing, setSyncing] = useState(false);
  const [syncDone, setSyncDone] = useState(0);
  const [syncTotal, setSyncTotal] = useState(0);
  const [loading, setLoading] = useState(false);

  const mmToPx = 96 / 25.4;
  const a4w = 210 * mmToPx;
  const thumbnailWidth = 40; // px
  const thumbnailScale = thumbnailWidth / a4w;

  const load = async () => {
    setLoading(true);
    try {
      const res = await fetch(`${API_BASE}/templates`);
      const data = await res.json();
      const mapped = (data.items || []).map((t: any) => ({
        id: t.ExternalID || t.id,
        name: t.Name || t.name,
        tags: typeof t.Tags === 'string' ? (t.Tags as string).split(',').map((x: string) => x.trim()).filter(Boolean) : (t.tags || []),
        usageCount: t.UsageCount ?? t.usageCount ?? t.Popularity ?? t.popularity ?? 0,
        isPremium: t.IsPremium ?? t.isPremium ?? false,
        category: t.Category || t.category || '',
      }));
      setItems(mapped);
    } catch {
      showToast(t('admin.msg.loadTemplatesFailed'), 'error');
    } finally { setLoading(false); }
  };
  useEffect(() => { load(); }, []);

  const handleSyncMockData = async () => {
    if (syncing) return;
    setSyncing(true);
    setSyncDone(0);
    const data = MOCK_TEMPLATES;
    setSyncTotal(data.length);
    const existing = new Set(items.map(i => i.id));
    let errors = 0;
    for (const t of data) {
      const body = { externalId: t.id, name: t.name, tags: (t.tags || []).join(','), category: t.category, usageCount: (t as any).usageCount ?? (t as any).popularity ?? 0, isPremium: t.isPremium };
      try {
        if (existing.has(t.id)) {
          await updateTemplate(t.id, { name: body.name, tags: body.tags, category: body.category, usageCount: body.usageCount, isPremium: body.isPremium });
        } else {
          await createTemplate(body);
          existing.add(t.id);
        }
      } catch {
        errors += 1;
      } finally {
        setSyncDone(x => x + 1);
      }
    }
    setSyncing(false);
    await load();
    if (errors === 0) {
      showToast(t('admin.sync.complete'), 'success');
    } else {
      showToast(t('admin.sync.partial').replace('{count}', String(errors)), 'warning');
    }
  };
  const filtered = items.filter(i => {
    const s = keyword.trim().toLowerCase();
    const m1 = !s || i.name.toLowerCase().includes(s) || i.id.toLowerCase().includes(s);
    const m2 = selectedCategory === 'All' || i.category === selectedCategory;
    return m1 && m2;
  });
  const total = filtered.length;
  const pageItems = filtered.slice((page - 1) * pageSize, (page - 1) * pageSize + pageSize);

  const openCreate = () => {
    setEditingId(null);
    setForm({ externalId: '', name: '', tags: '', category: '', usageCount: 0, isPremium: false });
    setShowForm(true);
  };
  const openEdit = (row: Row) => {
    setEditingId(row.id);
    setForm({
      externalId: row.id,
      name: row.name,
      tags: (row.tags || []).join(','),
      category: row.category || '',
      usageCount: row.usageCount || 0,
      isPremium: !!row.isPremium,
    });
    setShowForm(true);
  };
  const submitForm = async () => {
    try {
      if (editingId) {
        await updateTemplate(editingId, { name: form.name, tags: form.tags, category: form.category, usageCount: form.usageCount, isPremium: form.isPremium });
        showToast(t('admin.msg.templateUpdated'), 'success');
      } else {
        await createTemplate({ externalId: form.externalId, name: form.name, tags: form.tags, category: form.category, usageCount: form.usageCount, isPremium: form.isPremium });
        showToast(t('admin.msg.templateCreated'), 'success');
      }
      setShowForm(false);
      await load();
    } catch {
      showToast(editingId ? t('admin.msg.templateUpdateFailed') : t('admin.msg.templateCreateFailed'), 'error');
    }
  };

  const remove = async (id: string) => {
    const ok = await confirm({ title: t('common.confirmAction'), message: t('admin.actions.delete') });
    if (!ok) return;
    try {
      await deleteTemplate(id);
      showToast(t('admin.msg.templateDeleted'), 'success');
      await load();
    } catch {
      showToast(t('admin.msg.templateDeleteFailed'), 'error');
    }
  };

  const openPreview = (id: string) => {
    window.open(`${window.location.origin}${window.location.pathname}#${AppRoute.Print}?template=${id}`, '_blank');
  };

  

  const CategoryBadge: React.FC<{ value: string }> = ({ value }) => {
    const cls =
      value === 'IT'
        ? 'bg-blue-50 text-blue-700 border-blue-100'
        : value === 'Finance'
        ? 'bg-emerald-50 text-emerald-700 border-emerald-100'
        : value === 'Creative'
        ? 'bg-purple-50 text-purple-700 border-purple-100'
        : 'bg-slate-50 text-slate-700 border-slate-100';
    return <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium border ${cls}`}>{value || '-'}</span>;
  };

  const Thumbnail: React.FC<{ templateId: string }> = ({ templateId }) => {
    return (
      <div className="w-10 h-12 bg-slate-100 rounded-md overflow-hidden flex-shrink-0 border border-slate-200 relative">
        <ResumeArtboard
          data={{ ...INITIAL_RESUME, templateId }}
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
    <div className="flex-1 flex flex-col bg-white rounded-3xl m-2 overflow-hidden shadow-sm border border-gray-100">
      <div className="px-10 pt-10 pb-6">
        <div className="flex justify-between items-center">
          <h1 className="text-4xl font-bold text-gray-800 tracking-tight">{t('admin.menu.templates')}</h1>
          <Button variant="outline" onClick={() => load()} disabled={loading}>
            <RefreshCw size={16} className={`${loading ? 'animate-spin' : ''} mr-2`} /> {t('common.refresh') || 'Refresh'}
          </Button>
        </div>
        <div className="flex items-center justify-between mt-6 pb-4 border-b border-gray-100">
          <div className="flex items-center space-x-2">
            <input value={keyword} onChange={e => { setKeyword(e.target.value); setPage(1); }} placeholder={t('admin.keyword')} className="border rounded-lg px-3 py-2 text-sm shadow-sm border-gray-200" />
            <div className="flex items-center space-x-2">
              <label className="text-sm text-gray-500">{t('templates.filter.industry')}:</label>
              <select className="border rounded-md px-2 py-1 text-sm" value={selectedCategory} onChange={e => { setSelectedCategory(e.target.value); setPage(1); }}>
                {categories.map(c => {
                  const key = c === 'All' ? 'all' : c;
                  return <option key={c} value={c}>{t(`templates.category.${key}`)}</option>;
                })}
              </select>
            </div>
          </div>
          <div className="flex items-center space-x-2">
            <Button onClick={openCreate}>{t('admin.actions.create')}</Button>
            <Button variant="outline" onClick={handleSyncMockData} disabled={syncing}>
              {syncing ? `${syncDone}/${syncTotal} ${t('admin.sync.syncing')}` : t('admin.sync.syncMockData')}
            </Button>
          </div>
        </div>
      </div>
      <div className="flex-1 overflow-y-auto px-10 pb-10">
        <div className="overflow-x-auto">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="bg-slate-50/80 border-b border-slate-200">
                <th className="px-6 py-4 text-xs font-bold text-slate-500 uppercase tracking-wider">ID</th>
                <th className="px-6 py-4 text-xs font-bold text-slate-500 uppercase tracking-wider">{t('admin.form.name')}</th>
                <th className="px-6 py-4 text-xs font-bold text-slate-500 uppercase tracking-wider">{t('templates.filter.industry')}</th>
                <th className="px-6 py-4 text-xs font-bold text-slate-500 uppercase tracking-wider">{t('admin.form.usageCount')}</th>
                <th className="px-6 py-4 text-xs font-bold text-slate-500 uppercase tracking-wider">{t('admin.form.isPremium')}</th>
                <th className="px-6 py-4 text-xs font-bold text-slate-500 uppercase tracking-wider">{t('admin.form.tags')}</th>
                <th className="px-6 py-4 text-xs font-bold text-slate-500 uppercase tracking-wider text-right">{t('admin.columns.actions')}</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100">
              {pageItems.map(r => (
                <tr key={r.id} className="hover:bg-indigo-50/30 transition-colors group">
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-slate-400">#{r.id}</td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="flex items-center gap-3">
                      <Thumbnail templateId={r.id} />
                      <div>
                        <div className="text-sm font-semibold text-slate-900">{r.name}</div>
                        <div className="text-[10px] text-slate-400 uppercase tracking-tight">{r.category || '-'}</div>
                      </div>
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <CategoryBadge value={r.category} />
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <span className="text-sm font-semibold text-slate-700">{r.usageCount}</span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    {r.isPremium ? (
                      <span className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-bold bg-amber-50 text-amber-700 border border-amber-200 shadow-sm">
                        {t('home.badge.premium')}
                      </span>
                    ) : (
                      <span className="text-xs font-medium text-slate-400">Standard</span>
                    )}
                  </td>
                  <td className="px-6 py-4">
                    <div className="flex flex-wrap gap-1.5 max-w-[200px]">
                      {(r.tags || []).slice(0, 6).map(tag => (
                        <span key={tag} className="px-2 py-0.5 bg-white border border-slate-200 rounded text-[11px] text-slate-600 font-medium hover:border-indigo-300 hover:bg-indigo-50 transition-colors">
                          {tag}
                        </span>
                      ))}
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-right">
                    <div className="flex items-center justify-end gap-1">
                      <button
                        onClick={() => openPreview(r.id)}
                        className="p-2 text-slate-400 hover:text-indigo-600 hover:bg-indigo-50 rounded-lg transition-all"
                        title={t('admin.actions.preview')}
                      >
                        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7Z"/><circle cx="12" cy="12" r="3"/></svg>
                      </button>
                      <button
                        onClick={() => openEdit(r)}
                        className="p-2 text-slate-400 hover:text-emerald-600 hover:bg-emerald-50 rounded-lg transition-all"
                        title={t('admin.actions.update')}
                      >
                        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M17 3a2.85 2.83 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5Z"/><path d="m15 5 4 4"/></svg>
                      </button>
                      <button
                        onClick={() => remove(r.id)}
                        className="p-2 text-slate-400 hover:text-rose-600 hover:bg-rose-50 rounded-lg transition-all"
                        title={t('admin.actions.delete')}
                      >
                        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M3 6h18"/><path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/><path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/><line x1="10" x2="10" y1="11" y2="17"/><line x1="14" x2="14" y1="11" y2="17"/></svg>
                      </button>
                    </div>
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
            <select className="border rounded-md px-2 py-1 text-sm" value={pageSize} onChange={e => { setPageSize(parseInt(e.target.value)); setPage(1); }}>
              <option value={10}>10</option>
              <option value={20}>20</option>
              <option value={50}>50</option>
            </select>
          </div>
        </div>
      </div>

      <Modal isOpen={showForm} onClose={() => setShowForm(false)} title={editingId ? t('admin.actions.update') : t('admin.actions.create')}>
        <div className="space-y-3">
          {!editingId && (
            <div className="space-y-1">
              <label className="block text-sm font-medium text-gray-700">{t('admin.form.externalId')}</label>
              <input className="border rounded-md px-3 py-2 text-sm w-full" placeholder={t('admin.form.externalId')} value={form.externalId} onChange={e => setForm({ ...form, externalId: e.target.value })} />
            </div>
          )}
          <div className="space-y-1">
            <label className="block text-sm font-medium text-gray-700">{t('admin.form.name')}</label>
            <input className="border rounded-md px-3 py-2 text-sm w-full" placeholder={t('admin.form.name')} value={form.name} onChange={e => setForm({ ...form, name: e.target.value })} />
          </div>
          <div className="space-y-1">
            <label className="block text-sm font-medium text-gray-700">{t('admin.form.tags')}</label>
            <input className="border rounded-md px-3 py-2 text-sm w-full" placeholder={t('admin.form.tags')} value={form.tags} onChange={e => setForm({ ...form, tags: e.target.value })} />
          </div>
          <div className="grid grid-cols-2 gap-3">
            <div className="space-y-1">
              <label className="block text-sm font-medium text-gray-700">{t('admin.form.category')}</label>
              <input
                className="border rounded-md px-3 py-2 text-sm w-full"
                placeholder={t('admin.form.category')}
                value={form.category}
                onChange={e => setForm({ ...form, category: e.target.value })}
              />
            </div>
            <div className="space-y-1">
              <label className="block text-sm font-medium text-gray-700">{t('admin.form.usageCount')}</label>
              <input className="border rounded-md px-3 py-2 text-sm w-full" type="number" min={0} placeholder={t('admin.form.usageCount')} value={form.usageCount} onChange={e => setForm({ ...form, usageCount: parseInt(e.target.value || '0') })} />
            </div>
          </div>
          <div className="space-y-1">
            <label className="block text-sm font-medium text-gray-700">{t('admin.form.isPremium')}</label>
            <label className="inline-flex items-center space-x-2 text-sm">
              <input type="checkbox" checked={form.isPremium} onChange={e => setForm({ ...form, isPremium: e.target.checked })} />
              <span>{form.isPremium ? t('admin.value.yes') : t('admin.value.no')}</span>
            </label>
          </div>
          <div className="pt-2 flex justify-end space-x-2">
            <Button onClick={submitForm}>{editingId ? t('admin.actions.update') : t('admin.actions.create')}</Button>
          </div>
        </div>
      </Modal>
    </div>
  );
};
