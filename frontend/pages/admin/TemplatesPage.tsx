import React, { useEffect, useLayoutEffect, useRef, useState } from 'react';
import { API_BASE } from '../../config';
import { createTemplate, updateTemplate, deleteTemplate } from '../../services/adminService';
import { Button } from '../../components/ui/Button';
import { Modal } from '../../components/ui/Modal';
import { useToast } from '../../components/ui/Toast';
import { useConfirm } from '../../components/ui/ConfirmDialog';
import { useLanguage } from '../../contexts/LanguageContext';
import { ResumeArtboard } from '../editor/ResumePreview';
import { INITIAL_RESUME } from '../../services/mockData';

type Row = {
  id: string;
  name: string;
  thumbnail?: string;
  tags: string[];
  popularity: number;
  isPremium: boolean;
  category: string;
  level: string;
};

export const TemplatesPage: React.FC = () => {
  const { t } = useLanguage();
  const { showToast } = useToast();
  const confirm = useConfirm();

  const [items, setItems] = useState<Row[]>([]);
  const [keyword, setKeyword] = useState('');
  const [selectedCategory, setSelectedCategory] = useState('All');
  const [selectedLevel, setSelectedLevel] = useState('All');
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);

  const categories = ['All', 'IT', 'Finance', 'Creative', 'General'];
  const levels = ['All', 'Intern', 'Junior', 'Senior', 'Executive'];

  const [showForm, setShowForm] = useState(false);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [form, setForm] = useState({ externalId: '', name: '', tags: '', category: '', level: '', popularity: 0, isPremium: false });

  const [showPreview, setShowPreview] = useState(false);
  const [previewId, setPreviewId] = useState<string | null>(null);
  const previewContainerRef = useRef<HTMLDivElement | null>(null);
  const [previewScale, setPreviewScale] = useState<number | null>(null);
  const rafRef = useRef<number | null>(null);
  const roRef = useRef<ResizeObserver | null>(null);

  const load = async () => {
    try {
      const res = await fetch(`${API_BASE}/templates`);
      const data = await res.json();
      const mapped = (data.items || []).map((t: any) => ({
        id: t.ExternalID || t.id,
        name: t.Name || t.name,
        thumbnail: t.Thumbnail || t.thumbnail,
        tags: typeof t.Tags === 'string' ? (t.Tags as string).split(',').map((x: string) => x.trim()).filter(Boolean) : (t.tags || []),
        popularity: t.Popularity ?? t.popularity ?? 0,
        isPremium: t.IsPremium ?? t.isPremium ?? false,
        category: t.Category || t.category || '',
        level: t.Level || t.level || '',
      }));
      setItems(mapped);
    } catch {
      showToast(t('admin.msg.loadTemplatesFailed'), 'error');
    }
  };
  useEffect(() => { load(); }, []);

  const filtered = items.filter(i => {
    const s = keyword.trim().toLowerCase();
    const m1 = !s || i.name.toLowerCase().includes(s) || i.id.toLowerCase().includes(s);
    const m2 = selectedCategory === 'All' || i.category === selectedCategory;
    const m3 = selectedLevel === 'All' || i.level === selectedLevel;
    return m1 && m2 && m3;
  });
  const total = filtered.length;
  const pageItems = filtered.slice((page - 1) * pageSize, (page - 1) * pageSize + pageSize);

  const openCreate = () => {
    setEditingId(null);
    setForm({ externalId: '', name: '', thumbnail: '', tags: '', category: '', level: '', popularity: 0, isPremium: false });
    setShowForm(true);
  };
  const openEdit = (row: Row) => {
    setEditingId(row.id);
    setForm({
      externalId: row.id,
      name: row.name,
      tags: (row.tags || []).join(','),
      category: row.category || '',
      level: row.level || '',
      popularity: row.popularity || 0,
      isPremium: !!row.isPremium,
    });
    setShowForm(true);
  };
  const submitForm = async () => {
    try {
      if (editingId) {
        await updateTemplate(editingId, { name: form.name, tags: form.tags, category: form.category, level: form.level, popularity: form.popularity, isPremium: form.isPremium });
        showToast(t('admin.msg.templateUpdated'), 'success');
      } else {
        await createTemplate({ externalId: form.externalId, name: form.name, tags: form.tags, category: form.category, level: form.level, popularity: form.popularity, isPremium: form.isPremium });
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
    setPreviewId(id);
    setShowPreview(true);
  };

  useLayoutEffect(() => {
    if (!showPreview) return;
    const mmToPx = 96 / 25.4;
    const a4w = 210 * mmToPx;
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

  return (
    <div>
      <div className="flex items-center justify-between mb-4">
        <div className="flex items-center space-x-2">
          <input value={keyword} onChange={e => { setKeyword(e.target.value); setPage(1); }} placeholder={t('admin.keyword')} className="border rounded-md px-3 py-2 text-sm" />
          <div className="flex items-center space-x-2">
            <label className="text-sm text-gray-500">{t('templates.filter.industry')}:</label>
            <select className="border rounded-md px-2 py-1 text-sm" value={selectedCategory} onChange={e => { setSelectedCategory(e.target.value); setPage(1); }}>
              {categories.map(c => {
                const key = c === 'All' ? 'all' : c;
                return <option key={c} value={c}>{t(`templates.category.${key}`)}</option>;
              })}
            </select>
          </div>
          <div className="flex items-center space-x-2">
            <label className="text-sm text-gray-500">{t('templates.filter.level')}:</label>
            <select className="border rounded-md px-2 py-1 text-sm" value={selectedLevel} onChange={e => { setSelectedLevel(e.target.value); setPage(1); }}>
              {levels.map(l => {
                const key = l === 'All' ? 'all' : l;
                return <option key={l} value={l}>{t(`templates.level.${key}`)}</option>;
              })}
            </select>
          </div>
        </div>
        <Button onClick={openCreate}>{t('admin.actions.create')}</Button>
      </div>

      <div className="bg-white shadow-sm rounded-md overflow-hidden">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-4 py-2 text-left text-xs font-medium text-gray-500">ID</th>
              <th className="px-4 py-2 text-left text-xs font-medium text-gray-500">{t('admin.form.name')}</th>
              <th className="px-4 py-2 text-left text-xs font-medium text-gray-500">{t('templates.filter.industry')}</th>
              <th className="px-4 py-2 text-left text-xs font-medium text-gray-500">{t('templates.filter.level')}</th>
              <th className="px-4 py-2 text-left text-xs font-medium text-gray-500">{t('admin.form.popularity')}</th>
              <th className="px-4 py-2 text-left text-xs font-medium text-gray-500">{t('admin.form.isPremium')}</th>
              <th className="px-4 py-2 text-left text-xs font-medium text-gray-500">{t('admin.form.tags')}</th>
              <th className="px-4 py-2 text-right text-xs font-medium text-gray-500">{t('admin.columns.actions')}</th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {pageItems.map(r => (
              <tr key={r.id}>
                <td className="px-4 py-2 text-sm text-gray-700">{r.id}</td>
                <td className="px-4 py-2 text-sm text-gray-700">{r.name}</td>
                <td className="px-4 py-2 text-sm text-gray-700">{r.category}</td>
                <td className="px-4 py-2 text-sm text-gray-700">{r.level}</td>
                <td className="px-4 py-2 text-sm text-gray-700">{r.popularity}</td>
                <td className="px-4 py-2 text-sm text-gray-700">{r.isPremium ? t('admin.value.yes') : t('admin.value.no')}</td>
                <td className="px-4 py-2 text-sm text-gray-700">
                  <div className="flex flex-wrap gap-1">
                    {(r.tags || []).slice(0, 4).map(tag => (
                      <span key={tag} className="px-2 py-0.5 bg-gray-100 text-gray-600 text-xs rounded">{tag}</span>
                    ))}
                  </div>
                </td>
                <td className="px-4 py-2 text-sm text-right space-x-2">
                  <Button variant="outline" onClick={() => openEdit(r)}>{t('admin.actions.update')}</Button>
                  <Button variant="ghost" onClick={() => openPreview(r.id)}>{t('admin.actions.preview')}</Button>
                  <Button variant="danger" onClick={() => remove(r.id)}>{t('admin.actions.delete')}</Button>
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

      <Modal isOpen={showForm} onClose={() => setShowForm(false)} title={editingId ? t('admin.actions.update') : t('admin.actions.create')}>
        <div className="space-y-3">
          {!editingId && (
            <input className="border rounded-md px-3 py-2 text-sm w-full" placeholder={t('admin.form.externalId')} value={form.externalId} onChange={e => setForm({ ...form, externalId: e.target.value })} />
          )}
          <input className="border rounded-md px-3 py-2 text-sm w-full" placeholder={t('admin.form.name')} value={form.name} onChange={e => setForm({ ...form, name: e.target.value })} />
          <input className="border rounded-md px-3 py-2 text-sm w-full" placeholder={t('admin.form.tags')} value={form.tags} onChange={e => setForm({ ...form, tags: e.target.value })} />
          <div className="grid grid-cols-2 gap-3">
            <input className="border rounded-md px-3 py-2 text-sm w-full" placeholder={t('templates.filter.industry')} value={form.category} onChange={e => setForm({ ...form, category: e.target.value })} />
            <input className="border rounded-md px-3 py-2 text-sm w-full" placeholder={t('templates.filter.level')} value={form.level} onChange={e => setForm({ ...form, level: e.target.value })} />
          </div>
          <div className="grid grid-cols-2 gap-3 items-center">
            <input className="border rounded-md px-3 py-2 text-sm w-full" type="number" placeholder={t('admin.form.popularity')} value={form.popularity} onChange={e => setForm({ ...form, popularity: parseInt(e.target.value || '0') })} />
            <label className="inline-flex items-center space-x-2 text-sm">
              <input type="checkbox" checked={form.isPremium} onChange={e => setForm({ ...form, isPremium: e.target.checked })} />
              <span>{t('admin.form.isPremium')}</span>
            </label>
          </div>
          <div className="pt-2 flex justify-end space-x-2">
            <Button onClick={submitForm}>{editingId ? t('admin.actions.update') : t('admin.actions.create')}</Button>
          </div>
        </div>
      </Modal>

      <Modal isOpen={showPreview} onClose={() => setShowPreview(false)} title={t('admin.actions.preview')}>
        <div ref={previewContainerRef} className="aspect-[210/297] bg-gray-100 overflow-hidden relative">
          <div className="absolute inset-0 flex items-center justify-center">
            {previewId && previewScale !== null ? (
              <div
                style={{ width: (96 / 25.4) * 210 * previewScale, height: (96 / 25.4) * 297 * previewScale }}
                className="relative select-none pointer-events-none shadow-sm bg-white"
              >
                <ResumeArtboard
                  data={{ ...INITIAL_RESUME, templateId: previewId }}
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
