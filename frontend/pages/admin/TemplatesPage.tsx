import React, { useCallback, useEffect, useLayoutEffect, useMemo, useRef, useState } from 'react';
import { API_BASE } from '../../config';
import {
  AdminContentPreset,
  AdminJobRole,
  adminCreateContentPreset,
  adminDeleteContentPreset,
  adminListContentPresets,
  adminListJobRoles,
  adminPatchContentPreset,
  createTemplate,
  deleteTemplate,
  updateTemplate,
} from '../../services/adminService';
import { Button } from '../../components/ui/Button';
import { Modal } from '../../components/ui/Modal';
import { useToast } from '../../components/ui/Toast';
import { useConfirm } from '../../components/ui/ConfirmDialog';
import { useLanguage } from '../../contexts/LanguageContext';
import { Checkbox, Input, Label, Select, Textarea } from '../../components/ui/Form';
import { MultiSelect } from '../../components/ui/MultiSelect';
import { DataTable } from '../../components/ui/DataTable';
import { TableCard } from '../../components/ui/TableCard';
import { ResumeArtboard } from '../editor/ResumePreview';
import { AlertCircle, AlignLeft, Check, CheckCircle2, Copy, FileJson, FileText, Layers, LayoutGrid, Minimize2, Search, Trash2 } from 'lucide-react';
import { AppRoute } from '../../types';
import { motion } from 'framer-motion';

type Row = {
  id: string;
  name: string;
};

export const TemplatesPage: React.FC = () => {
  const { t, language } = useLanguage();
  const { showToast } = useToast();
  const confirm = useConfirm();

  type TabKey = 'templates' | 'presets';
  const [tab, setTab] = useState<TabKey>('templates');
  const tabNavRef = useRef<HTMLDivElement | null>(null);
  const tabButtonRefs = useRef<Record<TabKey, HTMLButtonElement | null>>({ templates: null, presets: null });
  const [tabIndicator, setTabIndicator] = useState<{ left: number; width: number }>({ left: 0, width: 0 });

  const parseJSON = (s: string) => {
    try {
      JSON.parse(s);
      return true;
    } catch {
      return false;
    }
  };

  const PresetJsonEditor: React.FC<{ value: string; onChange: (val: string) => void }> = ({ value, onChange }) => {
    const [isCopied, setIsCopied] = useState(false);

    const validation = useMemo(() => {
      const raw = String(value || '');
      const trimmed = raw.trim();
      if (!trimmed) return { kind: 'empty' as const, isValid: true, message: '' };
      try {
        JSON.parse(trimmed);
        return { kind: 'valid' as const, isValid: true, message: t('admin.catalog.msg.jsonOk') };
      } catch (err: any) {
        return {
          kind: 'invalid' as const,
          isValid: false,
          message: t('admin.catalog.msg.invalidJson'),
          detail: err?.message ? String(err.message) : '',
        };
      }
    }, [t, value]);

    const hasText = String(value || '').trim().length > 0;
    const canTransform = validation.kind === 'valid';

    const handleFormat = useCallback(() => {
      if (!canTransform) return;
      try {
        const parsed = JSON.parse(String(value || ''));
        const pretty = JSON.stringify(parsed, null, 2);
        onChange(`${pretty}\n`);
      } catch {
        showToast(t('admin.catalog.msg.invalidJson'), 'error');
      }
    }, [canTransform, onChange, showToast, t, value]);

    const handleMinify = useCallback(() => {
      if (!canTransform) return;
      try {
        const parsed = JSON.parse(String(value || ''));
        onChange(JSON.stringify(parsed));
      } catch {
        showToast(t('admin.catalog.msg.invalidJson'), 'error');
      }
    }, [canTransform, onChange, showToast, t, value]);

    const handleCopy = useCallback(async () => {
      if (!hasText) return;
      try {
        await navigator.clipboard.writeText(String(value || ''));
        setIsCopied(true);
        showToast(t('common.copied'), 'success');
        setTimeout(() => setIsCopied(false), 1600);
      } catch {
        showToast(t('common.copyFailed'), 'error');
      }
    }, [hasText, showToast, t, value]);

    const handleClear = useCallback(async () => {
      if (!hasText) return;
      const ok = await confirm({ title: t('admin.confirm.clearJsonTitle'), message: t('admin.confirm.clearJsonMsg') });
      if (!ok) return;
      onChange('');
    }, [confirm, hasText, onChange, t]);

    const lines = String(value || '').split('\n').length;

    return (
      <div className="border rounded-lg bg-white overflow-hidden flex flex-col shadow-sm">
        <div className="px-3 py-2 bg-slate-50 border-b border-slate-200 flex flex-wrap items-center justify-between gap-3">
          <div className="flex items-center gap-2">
            <div className="p-1.5 bg-indigo-100 text-indigo-600 rounded-md">
              <FileJson size={16} />
            </div>
            {validation.kind === 'empty' ? (
              <span className="text-xs font-medium text-slate-400">{t('admin.catalog.msg.jsonEmptyHint')}</span>
            ) : (
              <span className={`flex items-center text-xs font-medium ${validation.isValid ? 'text-emerald-600' : 'text-rose-600'}`}>
                {validation.isValid ? <CheckCircle2 size={14} className="mr-1.5" /> : <AlertCircle size={14} className="mr-1.5" />}
                {validation.isValid ? validation.message : `${validation.message}${(validation as any).detail ? `: ${(validation as any).detail}` : ''}`}
              </span>
            )}
          </div>

          <div className="flex items-center gap-2">
            <Button type="button" variant="outline" size="sm" icon={<AlignLeft size={14} />} onClick={handleFormat} disabled={!canTransform}>
              {t('admin.catalog.actions.formatJson')}
            </Button>
            <Button type="button" variant="outline" size="sm" icon={<Minimize2 size={14} />} onClick={handleMinify} disabled={!canTransform}>
              {t('admin.catalog.actions.minifyJson')}
            </Button>
            <Button
              type="button"
              variant="outline"
              size="sm"
              icon={isCopied ? <Check size={14} /> : <Copy size={14} />}
              onClick={handleCopy}
              disabled={!hasText}
            >
              {isCopied ? t('admin.catalog.actions.copied') : t('admin.catalog.actions.copyJson')}
            </Button>
            <Button type="button" variant="ghost" size="sm" icon={<Trash2 size={16} />} onClick={handleClear} disabled={!hasText} aria-label={t('admin.catalog.actions.clearJson')} />
          </div>
        </div>

        <div className="relative">
          <textarea
            value={value || ''}
            onChange={(e) => onChange(e.target.value)}
            spellCheck={false}
            className={`w-full min-h-[320px] p-3 font-mono text-xs leading-relaxed resize-y bg-white text-slate-900 focus:outline-none transition-all ${
              validation.kind !== 'invalid' ? 'focus:ring-2 focus:ring-blue-500/20' : 'focus:ring-2 focus:ring-rose-500/20'
            }`}
            placeholder={t('admin.catalog.form.dataJsonPlaceholder')}
          />
          {validation.kind === 'invalid' ? (
            <div className="absolute bottom-3 right-3 px-3 py-1.5 bg-rose-50 border border-rose-100 rounded-md text-rose-600 text-[11px] font-semibold shadow-sm flex items-center">
              <AlertCircle size={14} className="mr-2" />
              {t('admin.catalog.msg.invalidJson')}
            </div>
          ) : null}
        </div>

        <div className="px-3 py-1.5 bg-slate-50 border-t border-slate-200 flex items-center justify-between text-[10px] uppercase font-bold tracking-widest text-slate-400">
          <div className="flex gap-3">
            <span>{t('admin.catalog.stats.characters')}: {String(value || '').length}</span>
            <span>{t('admin.catalog.stats.lines')}: {lines}</span>
          </div>
        </div>
      </div>
    );
  };

  const [items, setItems] = useState<Row[]>([]);
  const [keyword, setKeyword] = useState('');
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);

  const [showForm, setShowForm] = useState(false);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [form, setForm] = useState<{ externalId: string; name: string }>({
    externalId: '',
    name: '',
  });
  const [saving, setSaving] = useState(false);

  const [loading, setLoading] = useState(false);

  const mmToPx = 96 / 25.4;
  const a4w = 210 * mmToPx;
  const thumbnailWidth = 40;
  const thumbnailScale = thumbnailWidth / a4w;

  const load = async () => {
    setLoading(true);
    try {
      const res = await fetch(`${API_BASE}/templates`);
      const data = await res.json();
      const mapped = (data.items || []).map((t: any) => ({
        id: t.ExternalID || t.id,
        name: t.Name || t.name,
      }));
      setItems(mapped);
    } catch {
      showToast(t('admin.msg.loadTemplatesFailed'), 'error');
    } finally {
      setLoading(false);
    }
  };
  useEffect(() => {
    load();
  }, []);

  const filtered = items.filter((i) => {
    const s = keyword.trim().toLowerCase();
    const m1 = !s || i.name.toLowerCase().includes(s) || i.id.toLowerCase().includes(s);
    return m1;
  });
  const total = filtered.length;
  const pageItems = filtered.slice((page - 1) * pageSize, (page - 1) * pageSize + pageSize);

  const openCreate = () => {
    setEditingId(null);
    setForm({ externalId: '', name: '' });
    setShowForm(true);
  };
  const openEdit = (row: Row) => {
    setEditingId(row.id);
    setForm({
      externalId: row.id,
      name: row.name,
    });
    setShowForm(true);
  };
  const submitForm = async () => {
    if (saving) return;
    setSaving(true);
    try {
      if (editingId) {
        await updateTemplate(editingId, { name: form.name });
        showToast(t('admin.msg.templateUpdated'), 'success');
      } else {
        await createTemplate({ externalId: form.externalId, name: form.name });
        showToast(t('admin.msg.templateCreated'), 'success');
      }
      setShowForm(false);
      await load();
    } catch {
      showToast(editingId ? t('admin.msg.templateUpdateFailed') : t('admin.msg.templateCreateFailed'), 'error');
    } finally {
      setSaving(false);
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

  const tabs = useMemo(
    () => [
      { key: 'templates' as const, label: t('admin.templates.tab.templates'), Icon: LayoutGrid },
      { key: 'presets' as const, label: t('admin.templates.tab.presets'), Icon: FileText },
    ],
    [t]
  );

  const updateTabIndicator = useCallback(() => {
    const nav = tabNavRef.current;
    const btn = tabButtonRefs.current[tab];
    if (!nav || !btn) return;

    const navRect = nav.getBoundingClientRect();
    const btnRect = btn.getBoundingClientRect();
    const left = btnRect.left - navRect.left + nav.scrollLeft;
    const width = btnRect.width;
    setTabIndicator({ left, width });
  }, [tab, tabs]);

  useLayoutEffect(() => {
    updateTabIndicator();
    requestAnimationFrame(() => updateTabIndicator());
  }, [updateTabIndicator]);

  useEffect(() => {
    const handler = () => updateTabIndicator();
    window.addEventListener('resize', handler);
    return () => window.removeEventListener('resize', handler);
  }, [updateTabIndicator]);

  useEffect(() => {
    const nav = tabNavRef.current;
    const btn = tabButtonRefs.current[tab];
    if (!nav || !btn || typeof ResizeObserver === 'undefined') return;
    const ro = new ResizeObserver(() => updateTabIndicator());
    ro.observe(nav);
    ro.observe(btn);
    return () => ro.disconnect();
  }, [tab, updateTabIndicator]);

  useEffect(() => {
    let cancelled = false;
    const fonts = (document as any)?.fonts;
    if (fonts?.ready?.then) {
      fonts.ready.then(() => {
        if (!cancelled) updateTabIndicator();
      });
    }
    return () => {
      cancelled = true;
    };
  }, [updateTabIndicator]);

  const Thumbnail: React.FC<{ templateId: string }> = ({ templateId }) => {
    return (
      <div className="w-10 h-12 bg-slate-100 rounded-md overflow-hidden flex-shrink-0 border border-slate-200 relative">
        <ResumeArtboard
          data={{ id: 'preview', title: '', templateId, lastModified: Date.now(), language: 'zh', Personal: {}, Theme: {}, sections: [] } as any}
          scale={thumbnailScale}
          disableShadow
          style={{ margin: 0 }}
          className="absolute top-0 left-0"
          showPageHint={false}
        />
      </div>
    );
  };

  const [presets, setPresets] = useState<AdminContentPreset[]>([]);
  const [catalogPage, setCatalogPage] = useState(1);
  const [catalogPageSize, setCatalogPageSize] = useState(20);
  const [catalogTotal, setCatalogTotal] = useState(0);
  const [catalogQ, setCatalogQ] = useState('');
  const [catalogLoading, setCatalogLoading] = useState(false);
  const [catalogSaving, setCatalogSaving] = useState(false);
  const [catalogShowForm, setCatalogShowForm] = useState(false);
  const [catalogEditingId, setCatalogEditingId] = useState<number | null>(null);
  const [catalogForm, setCatalogForm] = useState<any>({});

  const [allRoles, setAllRoles] = useState<AdminJobRole[]>([]);
  const [allPresets, setAllPresets] = useState<AdminContentPreset[]>([]);
  const [allTemplates, setAllTemplates] = useState<Array<{ id: string; name: string }>>([]);

  const roleNameById = useMemo(() => {
    const m = new Map<number, string>();
    for (const r of allRoles) m.set(r.id, String(r.name || '').trim());
    return m;
  }, [allRoles]);
  const presetNameById = useMemo(() => {
    const m = new Map<number, string>();
    for (const p of allPresets) m.set(p.ID, p.Name);
    return m;
  }, [allPresets]);
  const templateNameById = useMemo(() => {
    const m = new Map<string, string>();
    for (const tp of allTemplates) m.set(tp.id, tp.name);
    return m;
  }, [allTemplates]);

  const TextCell: React.FC<{ text?: string; className?: string }> = ({ text, className }) => {
    const v = (text || '').trim();
    const display = v || '-';
    return (
      <div className={`w-full max-w-full overflow-hidden truncate whitespace-nowrap text-sm text-gray-700 ${className || ''}`} title={v || ''}>
        {display}
      </div>
    );
  };

  const NameCell: React.FC<{ name?: string; id?: string }> = ({ name, id }) => {
    const display = (name || '').trim() || (id || '').trim() || '-';
    return (
      <div className="w-full max-w-full overflow-hidden truncate whitespace-nowrap text-sm text-gray-700" title={display}>
        {display}
      </div>
    );
  };

  const BoolBadge: React.FC<{ value: boolean; yes: string; no: string }> = ({ value, yes, no }) => {
    const cls = value ? 'bg-emerald-50 text-emerald-700 border-emerald-200' : 'bg-rose-50 text-rose-700 border-rose-200';
    return <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium border ${cls}`}>{value ? yes : no}</span>;
  };

  const SectionTitle: React.FC<{ children: React.ReactNode }> = ({ children }) => (
    <div className="flex items-center space-x-2">
      <div className="w-1 h-5 bg-blue-600 rounded-full"></div>
      <div className="text-sm font-semibold text-slate-900 uppercase tracking-wider">{children}</div>
    </div>
  );

  const IconButton: React.FC<{ title: string; onClick: () => void; kind: 'edit' | 'delete' }> = ({ title, onClick, kind }) => {
    const base = 'p-2 text-slate-400 rounded-lg transition-all';
    if (kind === 'edit') {
      return (
        <button onClick={onClick} className={`${base} hover:text-emerald-600 hover:bg-emerald-50`} title={title}>
          <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
            <path d="M17 3a2.85 2.83 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5Z" />
            <path d="m15 5 4 4" />
          </svg>
        </button>
      );
    }
    return (
      <button onClick={onClick} className={`${base} hover:text-rose-600 hover:bg-rose-50`} title={title}>
        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
          <path d="M3 6h18" />
          <path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6" />
          <path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2" />
          <line x1="10" x2="10" y1="11" y2="17" />
          <line x1="14" x2="14" y1="11" y2="17" />
        </svg>
      </button>
    );
  };

  const refreshCatalogRefs = async () => {
    try {
      const pageSize = 100;
      const maxPages = 50;
      let pageNum = 1;
      const collected: AdminJobRole[] = [];
      while (pageNum <= maxPages) {
        const rs = await adminListJobRoles({ page: String(pageNum), pageSize: String(pageSize), language: String(language || '') });
        const chunk = rs.items || [];
        collected.push(...chunk);
        if (!chunk.length) break;
        if (rs.total && collected.length >= rs.total) break;
        pageNum += 1;
      }
      setAllRoles(collected);
    } catch {}
    try {
      const pageSize = 100;
      const maxPages = 50;
      let pageNum = 1;
      const collected: AdminContentPreset[] = [];
      while (pageNum <= maxPages) {
        const ps = await adminListContentPresets({ page: String(pageNum), pageSize: String(pageSize), language: String(language || '') });
        const chunk = ps.items || [];
        collected.push(...chunk);
        if (!chunk.length) break;
        if (ps.total && collected.length >= ps.total) break;
        pageNum += 1;
      }
      setAllPresets(collected);
    } catch {}
    try {
      const res = await fetch(`${API_BASE}/templates`);
      const data = await res.json();
      const tpls = (data.items || []).map((x: any) => ({ id: x.ExternalID || x.id, name: x.Name || x.name || (x.ExternalID || x.id) }));
      setAllTemplates(tpls);
    } catch {}
  };

  useEffect(() => {
    refreshCatalogRefs();
  }, [language]);

  const loadCatalog = async () => {
    if (tab === 'templates') return;
    setCatalogLoading(true);
    try {
      const resp = await adminListContentPresets({
        page: String(catalogPage),
        pageSize: String(catalogPageSize),
        q: catalogQ,
        language: String(language || ''),
      });
      setPresets(resp.items || []);
      setCatalogTotal(resp.total || 0);
    } catch {
      showToast(t('admin.msg.loadFailed'), 'error');
    } finally {
      setCatalogLoading(false);
    }
  };

  useEffect(() => {
    if (tab === 'templates') return;
    loadCatalog();
  }, [tab, catalogPage, catalogPageSize, language]);

  useEffect(() => {
    if (tab === 'templates') return;
    const timer = setTimeout(() => {
      if (catalogPage !== 1) setCatalogPage(1);
      else loadCatalog();
    }, 300);
    return () => clearTimeout(timer);
  }, [catalogQ]);

  const openCreatePreset = () => {
    setCatalogEditingId(null);
    setCatalogForm({
      name: '',
      language: '',
      roleId: '',
      dataJson: '{\n  \"title\": \"\",\n  \"language\": \"zh\",\n  \"Personal\": {},\n  \"Theme\": {},\n  \"sections\": []\n}\n',
      isActive: true,
    });
    setCatalogShowForm(true);
  };

  const openEditPreset = (row: AdminContentPreset) => {
    setCatalogEditingId(row.ID);
    setCatalogForm({
      name: row.Name,
      language: row.Language || 'zh',
      roleId: String(row.RoleID),
      dataJson: row.DataJSON || '',
      isActive: !!row.IsActive,
    });
    setCatalogShowForm(true);
  };

  const openCatalogCreate = () => {
    openCreatePreset();
  };

  const saveCatalog = async () => {
    if (catalogSaving) return;
    setCatalogSaving(true);
    try {
      if (!String(catalogForm.name || '').trim() || !String(catalogForm.language || '').trim() || !Number(catalogForm.roleId || 0)) {
        showToast(t('auth.error.fillAll'), 'error');
        setCatalogSaving(false);
        return;
      }
      const dj = String(catalogForm.dataJson || '').trim();
      if (dj && !parseJSON(dj)) {
        showToast(t('admin.catalog.msg.invalidJson'), 'error');
        setCatalogSaving(false);
        return;
      }

      if (!catalogEditingId) {
        await adminCreateContentPreset({
          name: String(catalogForm.name || '').trim(),
          language: String(catalogForm.language || '').trim(),
          roleId: Number(catalogForm.roleId),
          dataJson: String(catalogForm.dataJson || ''),
          isActive: !!catalogForm.isActive,
        });
      } else {
        await adminPatchContentPreset(catalogEditingId, {
          name: String(catalogForm.name || '').trim(),
          language: String(catalogForm.language || '').trim(),
          roleId: Number(catalogForm.roleId),
          dataJson: String(catalogForm.dataJson || ''),
          isActive: !!catalogForm.isActive,
        });
      }
      setCatalogShowForm(false);
      showToast(t('admin.msg.saveSuccess'), 'success');
      await refreshCatalogRefs();
      await loadCatalog();
    } catch {
      showToast(t('admin.msg.saveFailed'), 'error');
    } finally {
      setCatalogSaving(false);
    }
  };

  const removeCatalogItem = async (id: number) => {
    const ok = await confirm({ title: t('admin.confirm.delete'), message: t('admin.confirm.deleteMsg') });
    if (!ok) return;
    try {
      await adminDeleteContentPreset(id);
      showToast(t('admin.msg.deleteSuccess'), 'success');
      await refreshCatalogRefs();
      await loadCatalog();
    } catch {
      showToast(t('admin.msg.deleteFailed') || 'Failed', 'error');
    }
  };

  const renderCatalogFormFields = () => {
    return (
      <div className="space-y-8">
        <section className="space-y-4">
          <SectionTitle>{t('admin.catalog.section.basic') || 'Basic'}</SectionTitle>
          <div>
            <Label required htmlFor="preset-name">{t('admin.form.name')}</Label>
            <Input id="preset-name" value={catalogForm.name || ''} onChange={(e) => setCatalogForm((p: any) => ({ ...p, name: e.target.value }))} />
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <Label required htmlFor="preset-language">{t('admin.catalog.form.language')}</Label>
              <Select
                value={catalogForm.language || ''}
                onChange={(e) => setCatalogForm((p: any) => ({ ...p, language: e.target.value }))}
                options={[
                  { label: t('admin.form.selectPlaceholder'), value: '', disabled: true, hidden: true },
                  { label: 'zh', value: 'zh' },
                  { label: 'en', value: 'en' },
                ]}
              />
            </div>
            <div>
              <Label required htmlFor="preset-role">{t('admin.catalog.form.role')}</Label>
              <Select
                value={catalogForm.roleId || ''}
                onChange={(e) => setCatalogForm((p: any) => ({ ...p, roleId: e.target.value }))}
                options={[
                  { label: t('admin.form.selectPlaceholder'), value: '', disabled: true, hidden: true },
                  ...allRoles.map((r) => ({ label: String(r.name || '').trim() || String(r.id), value: String(r.id) })),
                ]}
              />
            </div>
          </div>
          <Checkbox label={t('admin.catalog.form.active')} checked={!!catalogForm.isActive} onChange={(checked) => setCatalogForm((p: any) => ({ ...p, isActive: checked }))} />
        </section>
        <section className="space-y-4">
          <SectionTitle>{t('admin.catalog.section.data') || 'Data'}</SectionTitle>
          <div>
            <Label>{t('admin.catalog.form.dataJson')}</Label>
            <PresetJsonEditor value={String(catalogForm.dataJson || '')} onChange={(val) => setCatalogForm((p: any) => ({ ...p, dataJson: val }))} />
          </div>
        </section>
      </div>
    );
  };

  const CatalogPagination = () => {
    const totalPages = Math.max(1, Math.ceil(catalogTotal / catalogPageSize));
    return (
      <div className="flex items-center justify-between text-sm text-gray-500">
        <div>
          {t('admin.pagination.total') || 'Total'}: {catalogTotal}
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline" size="sm" disabled={catalogPage <= 1} onClick={() => setCatalogPage((p) => Math.max(1, p - 1))}>
            {t('admin.pagination.prev') || 'Prev'}
          </Button>
          <div className="min-w-[64px] text-center">
            {catalogPage} / {totalPages}
          </div>
          <Button variant="outline" size="sm" disabled={catalogPage >= totalPages} onClick={() => setCatalogPage((p) => Math.min(totalPages, p + 1))}>
            {t('admin.pagination.next') || 'Next'}
          </Button>
          <div className="w-[96px]">
            <Select
              value={String(catalogPageSize)}
              onChange={(e) => {
                setCatalogPage(1);
                setCatalogPageSize(Number(e.target.value));
              }}
              options={[10, 20, 50, 100].map((x) => ({ label: String(x), value: String(x) }))}
              className="py-1.5"
            />
          </div>
        </div>
      </div>
    );
  };

  const renderCatalogTable = () => {
    if (tab === 'presets') {
      return (
        <DataTable<AdminContentPreset>
          data={presets}
          getRowKey={(row) => String(row.ID)}
          emptyState={{
            title: t('admin.catalog.empty') || '暂无数据',
            description: t('admin.catalog.emptyDesc') || '点击右上角创建开始新增',
          }}
          columns={[
            { key: 'name', label: t('admin.form.name'), minWidth: 280, render: (p) => <TextCell text={p.Name} className="text-gray-900 font-medium" /> },
            { key: 'language', label: t('admin.catalog.form.language'), minWidth: 120, nowrap: true, render: (p) => <span className="text-sm text-gray-600">{p.Language || '-'}</span> },
            { key: 'role', label: t('admin.catalog.form.role'), minWidth: 240, render: (p) => <NameCell name={roleNameById.get(p.RoleID)} id={String(p.RoleID)} /> },
            { key: 'active', label: t('admin.catalog.form.active'), minWidth: 120, nowrap: true, render: (p) => <BoolBadge value={!!p.IsActive} yes={t('admin.catalog.enabled') || 'Enabled'} no={t('admin.catalog.disabled') || 'Disabled'} /> },
            {
              key: 'actions',
              label: t('admin.columns.actions'),
              minWidth: 140,
              fixed: 'right',
              headerClassName: 'text-right',
              cellClassName: 'text-right',
              render: (p) => (
                <div className="flex items-center justify-end gap-1">
                  <IconButton title={t('common.edit') || 'Edit'} onClick={() => openEditPreset(p)} kind="edit" />
                  <IconButton title={t('common.delete') || 'Delete'} onClick={() => removeCatalogItem(p.ID)} kind="delete" />
                </div>
              ),
            },
          ]}
        />
      );
    }
    return null;
  };

  return (
    <div className="flex-1 flex flex-col bg-white rounded-3xl m-2 overflow-hidden shadow-sm border border-gray-100">
      <div className="px-10 pt-10 pb-6">
        <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4 border-b border-gray-100 pb-4">
          <nav
            ref={tabNavRef}
            className="relative inline-flex items-center gap-1 p-1 bg-gray-100 rounded-xl self-start overflow-x-auto overflow-y-visible no-scrollbar"
            onScroll={updateTabIndicator}
          >
            <motion.div
              className="absolute top-1 bottom-1 left-0 bg-white rounded-lg shadow-sm pointer-events-none"
              initial={false}
              animate={{ x: tabIndicator.left, width: tabIndicator.width, opacity: tabIndicator.width ? 1 : 0 }}
              transition={{ type: 'spring', stiffness: 450, damping: 40, mass: 0.4 }}
            />
            {tabs.map((x) => (
              <button
                key={x.key}
                type="button"
                aria-pressed={tab === x.key}
                ref={(el) => {
                  tabButtonRefs.current[x.key] = el;
                }}
                className={`relative z-10 flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-semibold transition-all duration-200 whitespace-nowrap ${
                  tab === x.key ? 'text-blue-600' : 'text-gray-500 hover:text-gray-800 hover:bg-gray-200/60'
                }`}
                onClick={() => {
                  setTab(x.key);
                  if (x.key !== 'templates') {
                    setCatalogPage(1);
                    setCatalogTotal(0);
                  }
                }}
              >
                <x.Icon className="w-4 h-4" />
                {x.label}
              </button>
            ))}
          </nav>
        </div>
      </div>

      <div className="flex-1 overflow-hidden px-10 pb-10 flex flex-col">
        {tab === 'templates' ? (
          <div className="flex flex-col flex-1 overflow-hidden">
            <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4 pb-4 border-b border-gray-100">
              <div className="relative w-full md:w-[360px]">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
                <Input
                  value={keyword}
                  onChange={(e) => {
                    setKeyword(e.target.value);
                    setPage(1);
                  }}
                  placeholder={t('admin.keyword')}
                  className="pl-10"
                />
              </div>
              <div className="flex items-center gap-2 md:justify-end">
                <Button onClick={openCreate}>{t('admin.actions.create')}</Button>
              </div>
            </div>

            <div className="mt-6 flex-1 overflow-hidden">
              <TableCard className="flex flex-col h-full">
                <div className="flex-1 overflow-auto no-scrollbar">
                  <table className="w-full text-left border-collapse text-sm">
                    <thead className="sticky top-0 z-20 bg-slate-50 shadow-sm">
                      <tr className="border-b border-slate-200">
                        <th className="px-6 py-4 font-semibold text-gray-600">ID</th>
                        <th className="px-6 py-4 font-semibold text-gray-600">{t('admin.form.name')}</th>
                        <th className="px-6 py-4 font-semibold text-gray-600 text-right">{t('admin.columns.actions')}</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-slate-100">
                      {pageItems.map((r) => (
                        <tr key={r.id} className="hover:bg-indigo-50/30 transition-colors group">
                          <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-slate-400">#{r.id}</td>
                          <td className="px-6 py-4 whitespace-nowrap">
                            <div className="flex items-center gap-3">
                              <Thumbnail templateId={r.id} />
                              <div>
                                <div className="text-sm font-semibold text-slate-900">{r.name}</div>
                              </div>
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
                <div className="shrink-0 border-t border-gray-100 bg-white px-6 py-4">
                  <div className="flex justify-between items-center">
                    <div className="text-sm text-gray-500">{t('admin.total')} {total}</div>
                    <div className="flex items-center gap-2">
                      <Button variant="outline" size="sm" disabled={page === 1} onClick={() => setPage((p) => Math.max(p - 1, 1))}>
                        {t('admin.prev')}
                      </Button>
                      <div className="text-sm text-gray-500 min-w-[64px] text-center">
                        {page} / {Math.max(1, Math.ceil(total / pageSize))}
                      </div>
                      <Button variant="outline" size="sm" onClick={() => setPage((p) => p + 1)} disabled={page * pageSize >= total}>
                        {t('admin.next')}
                      </Button>
                      <div className="w-[96px]">
                        <Select
                          value={String(pageSize)}
                          onChange={(e) => {
                            setPageSize(parseInt(e.target.value, 10));
                            setPage(1);
                          }}
                          options={[10, 20, 50].map((x) => ({ label: String(x), value: String(x) }))}
                          className="py-1.5"
                        />
                      </div>
                    </div>
                  </div>
                </div>
              </TableCard>
            </div>

            <Modal isOpen={showForm} onClose={() => setShowForm(false)} title={editingId ? t('admin.actions.update') : t('admin.actions.create')} size="xl">
              <div className="flex flex-col max-h-[75vh]">
                <div className="flex-1 overflow-y-auto no-scrollbar pr-1 space-y-8">
                  <section className="space-y-4">
                    <div className="flex items-center space-x-2">
                      <div className="w-1 h-5 bg-blue-600 rounded-full"></div>
                      <div className="text-sm font-semibold text-slate-900 uppercase tracking-wider">{t('admin.catalog.section.basic')}</div>
                    </div>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                      {!editingId ? (
                        <div>
                          <Label required htmlFor="tpl-external">{t('admin.form.externalId')}</Label>
                          <Input id="tpl-external" placeholder={t('admin.form.externalId')} value={form.externalId} onChange={(e) => setForm({ ...form, externalId: e.target.value })} />
                        </div>
                      ) : (
                        <div>
                          <Label htmlFor="tpl-external">{t('admin.form.externalId')}</Label>
                          <Input id="tpl-external" value={form.externalId} disabled />
                        </div>
                      )}
                      <div>
                        <Label required htmlFor="tpl-name">{t('admin.form.name')}</Label>
                        <Input id="tpl-name" placeholder={t('admin.form.name')} value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} />
                      </div>
                    </div>
                  </section>
                </div>

                <div className="mt-6 pt-4 border-t border-slate-100 bg-slate-50/50 flex items-center justify-end space-x-3">
                  <Button variant="outline" onClick={() => setShowForm(false)} disabled={saving}>
                    {t('common.cancel')}
                  </Button>
                  <Button onClick={submitForm} disabled={saving}>
                    {editingId ? t('admin.actions.update') : t('admin.actions.create')}
                  </Button>
                </div>
              </div>
            </Modal>
          </div>
        ) : tab === 'presets' ? (
          <div className="flex flex-col flex-1 overflow-hidden">
            <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4 pb-4 border-b border-gray-100">
              <div className="relative w-full md:w-[360px]">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
                <Input value={catalogQ} onChange={(e) => setCatalogQ(e.target.value)} placeholder={t('admin.keyword')} className="pl-10" />
              </div>
              <div className="flex items-center gap-2 md:justify-end">
                <Button
                  onClick={() => {
                    setCatalogPage(1);
                    openCatalogCreate();
                  }}
                >
                  {t('admin.actions.create') || 'Create'}
                </Button>
              </div>
            </div>
            <div className="mt-6 flex-1 overflow-hidden">
              <TableCard className="flex flex-col h-full">
                {renderCatalogTable()}
                <div className="shrink-0 border-t border-gray-100 bg-white px-6 py-4">
                  <CatalogPagination />
                </div>
              </TableCard>
            </div>

            <Modal
              isOpen={catalogShowForm}
              onClose={() => setCatalogShowForm(false)}
              title={catalogEditingId ? (t('common.edit') || 'Edit') : (t('admin.actions.create') || 'Create')}
              size="lg"
              compact
            >
              <div className="flex flex-col max-h-[70vh]">
                <div className="flex-1 overflow-y-auto no-scrollbar pr-1 space-y-5">{renderCatalogFormFields()}</div>
                <div className="mt-4 pt-3 border-t border-slate-100 bg-slate-50/50 flex items-center justify-end gap-3">
                  <Button variant="outline" onClick={() => setCatalogShowForm(false)} disabled={catalogSaving}>
                    {t('common.cancel') || 'Cancel'}
                  </Button>
                  <Button onClick={saveCatalog} disabled={catalogSaving}>
                    {t('common.save') || 'Save'}
                  </Button>
                </div>
              </div>
            </Modal>
          </div>
        ) : null}
      </div>
    </div>
  );
};
