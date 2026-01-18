import React, { useEffect, useMemo, useState } from 'react';
import { Button } from '../../components/ui/Button';
import { Checkbox, Input, Label, Select, TagInput, Textarea } from '../../components/ui/Form';
import { MultiSelect } from '../../components/ui/MultiSelect';
import { DataTable } from '../../components/ui/DataTable';
import { TableCard } from '../../components/ui/TableCard';
import { Modal } from '../../components/ui/Modal';
import { useConfirm } from '../../components/ui/ConfirmDialog';
import { useToast } from '../../components/ui/Toast';
import { useLanguage } from '../../contexts/LanguageContext';
import { FileText, Layers, LayoutTemplate } from 'lucide-react';
import {
  adminCreateContentPreset,
  adminCreateJobCategory,
  adminCreateJobRole,
  adminCreateTemplateVariant,
  adminDeleteContentPreset,
  adminDeleteJobCategory,
  adminDeleteJobRole,
  adminDeleteTemplateVariant,
  adminGenerateTemplateVariants,
  adminListContentPresets,
  adminListJobCategories,
  adminListJobRoles,
  adminListTemplateVariants,
  adminPatchContentPreset,
  adminPatchJobCategory,
  adminPatchJobRole,
  adminPatchTemplateVariant,
  AdminContentPreset,
  AdminJobCategory,
  AdminJobRole,
  AdminTemplateVariant,
} from '../../services/adminService';
import { API_BASE } from '../../config';

type TabKey = 'jobs' | 'presets' | 'variants';

type FormKind = 'category' | 'role' | 'preset' | 'variant';

const parseJSON = (s: string) => {
  try {
    JSON.parse(s);
    return true;
  } catch {
    return false;
  }
};

const parseTags = (v: any): string[] => {
  if (Array.isArray(v)) return v.map(x => String(x).trim()).filter(Boolean);
  return String(v || '')
    .split(',')
    .map(x => x.trim())
    .filter(Boolean);
};

const joinTags = (tags: any): string => parseTags(tags).join(',');

export const CatalogPage: React.FC = () => {
  const { t } = useLanguage();
  const { showToast } = useToast();
  const confirm = useConfirm();

  const [tab, setTab] = useState<TabKey>('jobs');
  const [formKind, setFormKind] = useState<FormKind>('category');
  const [presets, setPresets] = useState<AdminContentPreset[]>([]);
  const [variants, setVariants] = useState<AdminTemplateVariant[]>([]);

  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const [total, setTotal] = useState(0);
  const [q, setQ] = useState('');
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);
  const [expandedCategories, setExpandedCategories] = useState<Record<string, boolean>>({});

  const [allCategories, setAllCategories] = useState<AdminJobCategory[]>([]);
  const [allRoles, setAllRoles] = useState<AdminJobRole[]>([]);
  const [allPresets, setAllPresets] = useState<AdminContentPreset[]>([]);
  const [allTemplates, setAllTemplates] = useState<Array<{ id: string; name: string }>>([]);

  const categoryNameById = useMemo(() => {
    const m = new Map<string, string>();
    for (const c of allCategories) m.set(c.ExternalID, c.Name);
    return m;
  }, [allCategories]);
  const roleNameById = useMemo(() => {
    const m = new Map<string, string>();
    for (const r of allRoles) m.set(r.ExternalID, r.Name);
    return m;
  }, [allRoles]);
  const presetNameById = useMemo(() => {
    const m = new Map<string, string>();
    for (const p of allPresets) m.set(p.ExternalID, p.Name);
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
    if (name && id && name !== id) {
      return (
        <div className="w-full max-w-full overflow-hidden" title={name}>
          <div className="text-sm text-gray-700 truncate whitespace-nowrap" title={name}>{name}</div>
          <div className="text-[11px] text-gray-400 font-mono truncate whitespace-nowrap" title={id}>{id}</div>
        </div>
      );
    }
    return <div className="w-full max-w-full overflow-hidden truncate whitespace-nowrap text-sm text-gray-700" title={display}>{display}</div>;
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
          <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M17 3a2.85 2.83 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5Z"/><path d="m15 5 4 4"/></svg>
        </button>
      );
    }
    return (
      <button onClick={onClick} className={`${base} hover:text-rose-600 hover:bg-rose-50`} title={title}>
        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M3 6h18"/><path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/><path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/><line x1="10" x2="10" y1="11" y2="17"/><line x1="14" x2="14" y1="11" y2="17"/></svg>
      </button>
    );
  };

  const [generateForm, setGenerateForm] = useState({
    roleId: '',
    presetId: '',
    layoutTemplateIds: [] as string[],
    namePrefix: '',
    tags: '',
    isPremium: false,
    isActive: true,
    mode: 'skip' as 'skip' | 'update',
  });
  const [generateResult, setGenerateResult] = useState<any>(null);
  const [generating, setGenerating] = useState(false);
  const [variantCreateMode, setVariantCreateMode] = useState<'single' | 'batch'>('single');

  const resetPaging = () => {
    setPage(1);
    setTotal(0);
  };

  useEffect(() => {
    resetPaging();
  }, [tab]);

  const refreshRefs = async () => {
    try {
      const cats = await adminListJobCategories({ page: '1', pageSize: '1000' });
      setAllCategories(cats.items || []);
    } catch {}
    try {
      const rs = await adminListJobRoles({ page: '1', pageSize: '2000' });
      setAllRoles(rs.items || []);
    } catch {}
    try {
      const ps = await adminListContentPresets({ page: '1', pageSize: '2000' });
      setAllPresets(ps.items || []);
    } catch {}
    try {
      const res = await fetch(`${API_BASE}/templates`);
      const data = await res.json();
      const items = (data.items || []).map((x: any) => ({ id: x.ExternalID || x.id, name: x.Name || x.name || (x.ExternalID || x.id) }));
      setAllTemplates(items);
    } catch {}
  };

  useEffect(() => {
    refreshRefs();
  }, []);

  const tabs = useMemo(
    () => [
      { key: 'jobs' as const, label: t('admin.catalog.tab.jobs'), Icon: Layers },
      { key: 'presets' as const, label: t('admin.catalog.tab.presets'), Icon: FileText },
      { key: 'variants' as const, label: t('admin.catalog.tab.variants'), Icon: LayoutTemplate },
    ],
    [t]
  );

  const load = async () => {
    setLoading(true);
    try {
      if (tab === 'jobs') {
        setTotal(allCategories.length);
      } else if (tab === 'presets') {
        const resp = await adminListContentPresets({ page: String(page), pageSize: String(pageSize), q });
        setPresets(resp.items || []);
        setTotal(resp.total || 0);
      } else {
        const resp = await adminListTemplateVariants({ page: String(page), pageSize: String(pageSize), q });
        setVariants(resp.items || []);
        setTotal(resp.total || 0);
      }
    } catch (e: any) {
      showToast(t('admin.msg.loadFailed'), 'error');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    load();
    if (tab === 'jobs') {
      setFormKind('category');
    } else if (tab === 'presets') {
      setFormKind('preset');
    } else {
      setFormKind('variant');
    }
  }, [tab, page, pageSize, allCategories.length]);

  useEffect(() => {
    const timer = setTimeout(() => {
      if (page !== 1) setPage(1);
      else load();
    }, 300);
    return () => clearTimeout(timer);
  }, [q]);

  const [showForm, setShowForm] = useState(false);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [form, setForm] = useState<any>({});
  const formatPresetDataJson = () => {
    const raw = String(form.dataJson || '');
    if (!raw.trim()) return;
    try {
      const parsed = JSON.parse(raw);
      const pretty = JSON.stringify(parsed, null, 2);
      setForm((p: any) => ({ ...p, dataJson: `${pretty}\n` }));
    } catch {
      showToast(t('admin.catalog.msg.invalidJson'), 'error');
    }
  };
  const openCreateCategory = () => {
    setEditingId(null);
    setFormKind('category');
    setForm({ externalId: '', name: '', parentExternalId: '', orderNum: 0, isActive: true });
    setShowForm(true);
  };
  const openCreateSubCategory = (parentExternalId: string) => {
    setEditingId(null);
    setFormKind('category');
    setForm({ externalId: '', name: '', parentExternalId, orderNum: 0, isActive: true });
    setShowForm(true);
  };
  const openCreateRole = (categoryExternalId: string) => {
    setEditingId(null);
    setFormKind('role');
    setForm({ externalId: '', categoryExternalId, name: '', tags: [] as string[], orderNum: 0, isActive: true });
    setShowForm(true);
  };
  const openCreatePreset = () => {
    setEditingId(null);
    setFormKind('preset');
    setForm({ externalId: '', name: '', language: '', roleExternalId: '', tags: [] as string[], dataJson: '{\n  \"title\": \"\",\n  \"language\": \"zh\",\n  \"Personal\": {},\n  \"Theme\": {},\n  \"sections\": []\n}\n', isActive: true });
    setShowForm(true);
  };
  const openCreateVariant = () => {
    setEditingId(null);
    setFormKind('variant');
    setForm({ externalId: '', name: '', layoutTemplateExternalId: '', presetExternalId: '', roleExternalId: '', tags: [] as string[], usageCount: 0, isPremium: false, isActive: true });
    setShowForm(true);
  };

  const openEditCategory = (row: any) => {
    setEditingId(row.ExternalID);
    setFormKind('category');
    setForm({ externalId: row.ExternalID, name: row.Name, parentExternalId: row.ParentExternalID || '', orderNum: row.OrderNum ?? 0, isActive: !!row.IsActive });
    setShowForm(true);
  };
  const openEditRole = (row: any) => {
    setEditingId(row.ExternalID);
    setFormKind('role');
    setForm({ externalId: row.ExternalID, categoryExternalId: row.CategoryExternalID, name: row.Name, tags: parseTags(row.Tags), orderNum: row.OrderNum ?? 0, isActive: !!row.IsActive });
    setShowForm(true);
  };
  const openEditPreset = (row: any) => {
    setEditingId(row.ExternalID);
    setFormKind('preset');
    setForm({ externalId: row.ExternalID, name: row.Name, language: row.Language || 'zh', roleExternalId: row.RoleExternalID || '', tags: parseTags(row.Tags), dataJson: row.DataJSON || '', isActive: !!row.IsActive });
    setShowForm(true);
  };
  const openEditVariant = (row: any) => {
    setEditingId(row.ExternalID);
    setFormKind('variant');
    setForm({ externalId: row.ExternalID, name: row.Name, layoutTemplateExternalId: row.LayoutTemplateExternalID, presetExternalId: row.PresetExternalID, roleExternalId: row.RoleExternalID, tags: parseTags(row.Tags), usageCount: row.UsageCount ?? 0, isPremium: !!row.IsPremium, isActive: !!row.IsActive });
    setShowForm(true);
  };

  const openCreate = () => {
    if (tab === 'jobs') {
      openCreateCategory();
      return;
    }
    if (tab === 'presets') {
      openCreatePreset();
      return;
    }
    setVariantCreateMode('single');
    setGenerateResult(null);
    openCreateVariant();
  };

  const runGenerateVariants = async () => {
    if (!String(generateForm.roleId || '').trim() || !String(generateForm.presetId || '').trim() || (generateForm.layoutTemplateIds || []).length === 0) {
      showToast(t('auth.error.fillAll'), 'error');
      return;
    }
    if (generating) return;
    setGenerating(true);
    setGenerateResult(null);
    try {
      const res = await adminGenerateTemplateVariants({
        roleId: generateForm.roleId,
        presetId: generateForm.presetId,
        layoutTemplateIds: generateForm.layoutTemplateIds,
        namePrefix: generateForm.namePrefix,
        tags: generateForm.tags,
        isPremium: generateForm.isPremium,
        isActive: generateForm.isActive,
        mode: generateForm.mode,
      });
      setGenerateResult(res?.result || null);
      showToast(t('admin.catalog.generate.done'), 'success');
      await load();
    } catch {
      showToast(t('admin.catalog.generate.failed'), 'error');
    } finally {
      setGenerating(false);
    }
  };

  const save = async () => {
    if (saving) return;
    setSaving(true);
    try {
      if (formKind === 'role') {
        if (!String(form.externalId || '').trim() || !String(form.categoryExternalId || '').trim() || !String(form.name || '').trim()) {
          showToast(t('auth.error.fillAll'), 'error');
          setSaving(false);
          return;
        }
      }
      if (formKind === 'preset') {
        if (!String(form.externalId || '').trim() || !String(form.name || '').trim() || !String(form.language || '').trim() || !String(form.roleExternalId || '').trim()) {
          showToast(t('auth.error.fillAll'), 'error');
          setSaving(false);
          return;
        }
      }
      if (formKind === 'variant') {
        if (!String(form.externalId || '').trim() || !String(form.name || '').trim() || !String(form.layoutTemplateExternalId || '').trim() || !String(form.presetExternalId || '').trim() || !String(form.roleExternalId || '').trim()) {
          showToast(t('auth.error.fillAll'), 'error');
          setSaving(false);
          return;
        }
      }
      if (formKind === 'preset') {
        const dj = String(form.dataJson || '').trim();
        if (dj && !parseJSON(dj)) {
          showToast(t('admin.catalog.msg.invalidJson'), 'error');
          setSaving(false);
          return;
        }
      }
      if (!editingId) {
        if (formKind === 'category') await adminCreateJobCategory(form);
        else if (formKind === 'role') await adminCreateJobRole({ ...form, tags: joinTags(form.tags) });
        else if (formKind === 'preset') await adminCreateContentPreset({ ...form, tags: joinTags(form.tags) });
        else await adminCreateTemplateVariant({ ...form, tags: joinTags(form.tags) });
      } else {
        const body = { ...form };
        delete body.externalId;
        if (formKind === 'category') await adminPatchJobCategory(editingId, body);
        else if (formKind === 'role') await adminPatchJobRole(editingId, { ...body, tags: joinTags(body.tags) });
        else if (formKind === 'preset') await adminPatchContentPreset(editingId, { ...body, tags: joinTags(body.tags) });
        else await adminPatchTemplateVariant(editingId, { ...body, tags: joinTags(body.tags) });
      }
      setShowForm(false);
      showToast(t('admin.msg.saveSuccess'), 'success');
      await refreshRefs();
      await load();
    } catch (e: any) {
      showToast(t('admin.msg.saveFailed'), 'error');
    } finally {
      setSaving(false);
    }
  };

  const removeItem = async (kind: FormKind, id: string) => {
    const ok = await confirm({ title: t('admin.confirm.delete'), message: t('admin.confirm.deleteMsg') });
    if (!ok) return;
    try {
      if (kind === 'category') await adminDeleteJobCategory(id);
      else if (kind === 'role') await adminDeleteJobRole(id);
      else if (kind === 'preset') await adminDeleteContentPreset(id);
      else await adminDeleteTemplateVariant(id);
      showToast(t('admin.msg.deleteSuccess'), 'success');
      await refreshRefs();
      await load();
    } catch {
      showToast(t('admin.msg.deleteFailed') || 'Failed', 'error');
    }
  };

  const renderFormFields = () => {
    if (formKind === 'category') {
      return (
        <div className="space-y-8">
          <section className="space-y-4">
            <SectionTitle>{t('admin.catalog.section.basic') || 'Basic'}</SectionTitle>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <Label required htmlFor="cat-external">{t('admin.form.externalId')}</Label>
                <Input id="cat-external" value={form.externalId || ''} disabled={!!editingId} onChange={(e) => setForm((p: any) => ({ ...p, externalId: e.target.value }))} />
              </div>
              <div>
                <Label required htmlFor="cat-name">{t('admin.form.name')}</Label>
                <Input id="cat-name" value={form.name || ''} onChange={(e) => setForm((p: any) => ({ ...p, name: e.target.value }))} />
              </div>
            </div>
          </section>
          <section className="space-y-4">
            <SectionTitle>{t('admin.catalog.section.config') || 'Config'}</SectionTitle>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <Label htmlFor="cat-parent">{t('admin.catalog.form.parent')}</Label>
                {!editingId && String(form.parentExternalId || '').trim() ? (
                  <div className="space-y-1">
                    <Input
                      id="cat-parent"
                      value={categoryNameById.get(String(form.parentExternalId)) || String(form.parentExternalId)}
                      disabled
                    />
                    <div className="text-[11px] text-slate-400 font-mono truncate whitespace-nowrap" title={String(form.parentExternalId)}>
                      {String(form.parentExternalId)}
                    </div>
                  </div>
                ) : (
                  <Select
                    value={form.parentExternalId || ''}
                    onChange={(e) => setForm((p: any) => ({ ...p, parentExternalId: e.target.value }))}
                    options={[
                      { label: t('admin.form.selectPlaceholder'), value: '', disabled: true, hidden: true },
                      ...allCategories.map(c => ({ label: c.Name, value: c.ExternalID })),
                    ]}
                  />
                )}
              </div>
              <div>
                <Label htmlFor="cat-order">{t('admin.catalog.form.order')}</Label>
                <Input id="cat-order" type="number" value={String(form.orderNum ?? 0)} onChange={(e) => setForm((p: any) => ({ ...p, orderNum: Number(e.target.value) }))} />
              </div>
            </div>
            <Checkbox label={t('admin.catalog.form.active')} checked={!!form.isActive} onChange={(checked) => setForm((p: any) => ({ ...p, isActive: checked }))} />
          </section>
        </div>
      );
    }
    if (formKind === 'role') {
      return (
        <div className="space-y-8">
          <section className="space-y-4">
            <SectionTitle>{t('admin.catalog.section.basic') || 'Basic'}</SectionTitle>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <Label required htmlFor="role-external">{t('admin.form.externalId')}</Label>
                <Input id="role-external" value={form.externalId || ''} disabled={!!editingId} onChange={(e) => setForm((p: any) => ({ ...p, externalId: e.target.value }))} />
              </div>
              <div>
                <Label required htmlFor="role-category">{t('admin.catalog.form.category')}</Label>
                {!editingId && String(form.categoryExternalId || '').trim() ? (
                  <div className="space-y-1">
                    <Input
                      id="role-category"
                      value={categoryNameById.get(String(form.categoryExternalId)) || String(form.categoryExternalId)}
                      disabled
                    />
                    <div className="text-[11px] text-slate-400 font-mono truncate whitespace-nowrap" title={String(form.categoryExternalId)}>
                      {String(form.categoryExternalId)}
                    </div>
                  </div>
                ) : (
                  <Select
                    value={form.categoryExternalId || ''}
                    onChange={(e) => setForm((p: any) => ({ ...p, categoryExternalId: e.target.value }))}
                    options={[
                      { label: t('admin.form.selectPlaceholder'), value: '', disabled: true, hidden: true },
                      ...allCategories.map(c => ({ label: c.Name, value: c.ExternalID })),
                    ]}
                  />
                )}
              </div>
            </div>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <Label required htmlFor="role-name">{t('admin.form.name')}</Label>
                <Input id="role-name" value={form.name || ''} onChange={(e) => setForm((p: any) => ({ ...p, name: e.target.value }))} />
              </div>
              <div>
                <Label htmlFor="role-order">{t('admin.catalog.form.order')}</Label>
                <Input id="role-order" type="number" value={String(form.orderNum ?? 0)} onChange={(e) => setForm((p: any) => ({ ...p, orderNum: Number(e.target.value) }))} />
              </div>
            </div>
            <Checkbox label={t('admin.catalog.form.active')} checked={!!form.isActive} onChange={(checked) => setForm((p: any) => ({ ...p, isActive: checked }))} />
          </section>
          <section className="space-y-4">
            <SectionTitle>{t('admin.catalog.section.tags') || 'Tags'}</SectionTitle>
            <div>
              <Label>{t('admin.form.tags')}</Label>
              <TagInput tags={parseTags(form.tags)} onChange={(tags) => setForm((p: any) => ({ ...p, tags }))} />
            </div>
          </section>
        </div>
      );
    }
    if (formKind === 'preset') {
      return (
        <div className="space-y-8">
          <section className="space-y-4">
            <SectionTitle>{t('admin.catalog.section.basic') || 'Basic'}</SectionTitle>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <Label required htmlFor="preset-external">{t('admin.form.externalId')}</Label>
                <Input id="preset-external" value={form.externalId || ''} disabled={!!editingId} onChange={(e) => setForm((p: any) => ({ ...p, externalId: e.target.value }))} />
              </div>
              <div>
                <Label required htmlFor="preset-name">{t('admin.form.name')}</Label>
                <Input id="preset-name" value={form.name || ''} onChange={(e) => setForm((p: any) => ({ ...p, name: e.target.value }))} />
              </div>
            </div>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <Label required htmlFor="preset-language">{t('admin.catalog.form.language')}</Label>
                <Select
                  value={form.language || ''}
                  onChange={(e) => setForm((p: any) => ({ ...p, language: e.target.value }))}
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
                  value={form.roleExternalId || ''}
                  onChange={(e) => setForm((p: any) => ({ ...p, roleExternalId: e.target.value }))}
                  options={[
                    { label: t('admin.form.selectPlaceholder'), value: '', disabled: true, hidden: true },
                    ...allRoles.map(r => ({ label: r.Name, value: r.ExternalID })),
                  ]}
                />
              </div>
            </div>
            <Checkbox label={t('admin.catalog.form.active')} checked={!!form.isActive} onChange={(checked) => setForm((p: any) => ({ ...p, isActive: checked }))} />
          </section>
          <section className="space-y-4">
            <SectionTitle>{t('admin.catalog.section.tags') || 'Tags'}</SectionTitle>
            <div>
              <Label>{t('admin.form.tags')}</Label>
              <TagInput tags={parseTags(form.tags)} onChange={(tags) => setForm((p: any) => ({ ...p, tags }))} />
            </div>
          </section>
          <section className="space-y-4">
            <SectionTitle>{t('admin.catalog.section.data') || 'Data'}</SectionTitle>
            <div>
              <Label
                right={(() => {
                  const raw = String(form.dataJson || '');
                  const trimmed = raw.trim();
                  const status = trimmed ? (parseJSON(raw) ? t('admin.catalog.msg.jsonOk') : t('admin.catalog.msg.invalidJson')) : null;
                  return (
                    <div className="flex items-center gap-2">
                      {status ? <span>{status}</span> : null}
                      <Button type="button" variant="outline" size="sm" onClick={formatPresetDataJson} disabled={!trimmed}>
                        {t('admin.catalog.actions.formatJson') || '格式化 JSON'}
                      </Button>
                    </div>
                  );
                })()}
              >
                {t('admin.catalog.form.dataJson')}
              </Label>
              <Textarea className="font-mono text-xs h-64" value={form.dataJson || ''} onChange={(e) => setForm((p: any) => ({ ...p, dataJson: e.target.value }))} />
            </div>
          </section>
        </div>
      );
    }
    return (
      <div className="space-y-8">
        {!editingId ? (
          <nav className="inline-flex items-center gap-1 p-1 bg-gray-100 rounded-xl">
            <button
              type="button"
              aria-pressed={variantCreateMode === 'single'}
              className={`px-4 py-2 rounded-lg text-sm font-semibold transition-all duration-200 ${variantCreateMode === 'single' ? 'bg-white text-blue-600 shadow-sm' : 'text-gray-500 hover:text-gray-800 hover:bg-gray-200/60'}`}
              onClick={() => setVariantCreateMode('single')}
            >
              {t('admin.catalog.generate.modeSingle')}
            </button>
            <button
              type="button"
              aria-pressed={variantCreateMode === 'batch'}
              className={`px-4 py-2 rounded-lg text-sm font-semibold transition-all duration-200 ${variantCreateMode === 'batch' ? 'bg-white text-blue-600 shadow-sm' : 'text-gray-500 hover:text-gray-800 hover:bg-gray-200/60'}`}
              onClick={() => setVariantCreateMode('batch')}
            >
              {t('admin.catalog.generate.title') || '批量生成'}
            </button>
          </nav>
        ) : null}

        {!editingId && variantCreateMode === 'batch' ? (
          <div className="space-y-8">
            <section className="space-y-4">
              <SectionTitle>{t('admin.catalog.section.config') || 'Config'}</SectionTitle>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                <div>
                  <Label required htmlFor="generate-role">{t('admin.catalog.form.role')}</Label>
                  <Select
                    value={generateForm.roleId}
                    onChange={(e) => setGenerateForm(p => ({ ...p, roleId: e.target.value }))}
                    options={[
                      { label: t('admin.form.selectPlaceholder'), value: '' },
                      ...allRoles.map(r => ({ label: r.Name, value: r.ExternalID })),
                    ]}
                  />
                </div>
                <div>
                  <Label required htmlFor="generate-preset">{t('admin.catalog.form.preset')}</Label>
                  <Select
                    value={generateForm.presetId}
                    onChange={(e) => setGenerateForm(p => ({ ...p, presetId: e.target.value }))}
                    options={[
                      { label: t('admin.form.selectPlaceholder'), value: '' },
                      ...allPresets.map(p => ({ label: p.Name, value: p.ExternalID })),
                    ]}
                  />
                </div>
                <div>
                  <Label htmlFor="generate-mode">{t('admin.catalog.generate.title')}</Label>
                  <Select
                    value={generateForm.mode}
                    onChange={(e) => setGenerateForm(p => ({ ...p, mode: e.target.value as any }))}
                    options={[
                      { label: t('admin.catalog.generate.modeSkip'), value: 'skip' },
                      { label: t('admin.catalog.generate.modeUpdate'), value: 'update' },
                    ]}
                  />
                </div>
              </div>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <Label htmlFor="generate-namePrefix">{t('admin.catalog.generate.namePrefix')}</Label>
                  <Input id="generate-namePrefix" value={generateForm.namePrefix} onChange={(e) => setGenerateForm(p => ({ ...p, namePrefix: e.target.value }))} />
                </div>
                <div>
                  <Label htmlFor="generate-tags">{t('admin.form.tags')}</Label>
                  <Input id="generate-tags" value={generateForm.tags} onChange={(e) => setGenerateForm(p => ({ ...p, tags: e.target.value }))} />
                </div>
              </div>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div className="md:col-span-2">
                  <Label required right={<span>{generateForm.layoutTemplateIds.length}</span>}>
                    {t('admin.catalog.form.template')}
                  </Label>
                  <MultiSelect
                    options={allTemplates.map((tp) => ({ value: tp.id, label: tp.name }))}
                    value={generateForm.layoutTemplateIds}
                    onChange={(vals) => setGenerateForm((p) => ({ ...p, layoutTemplateIds: vals }))}
                    placeholder={t('admin.form.selectPlaceholder')}
                    searchPlaceholder={t('admin.form.multiSelect.searchPlaceholder')}
                    selectAllLabel={t('admin.form.multiSelect.selectAll')}
                    unselectAllLabel={t('admin.form.multiSelect.unselectAll')}
                    availableSuffix={t('admin.form.multiSelect.availableSuffix')}
                    noResultsLabel={(q) => `${t('admin.form.multiSelect.noResults')}: “${q}”`}
                  />
                </div>
              </div>
              <div className="flex items-center gap-8">
                <Checkbox label={t('admin.form.isPremium')} checked={!!generateForm.isPremium} onChange={(checked) => setGenerateForm(p => ({ ...p, isPremium: checked }))} />
                <Checkbox label={t('admin.catalog.form.active')} checked={!!generateForm.isActive} onChange={(checked) => setGenerateForm(p => ({ ...p, isActive: checked }))} />
              </div>
            </section>

            {generateResult ? (
              <section className="space-y-4">
                <SectionTitle>{t('admin.catalog.generate.result')}</SectionTitle>
                <div className="text-sm text-gray-600 flex flex-wrap gap-x-4 gap-y-1">
                  <span>{t('admin.catalog.generate.created')}: {generateResult.created}</span>
                  <span>{t('admin.catalog.generate.updated')}: {generateResult.updated}</span>
                  <span>{t('admin.catalog.generate.skipped')}: {generateResult.skipped}</span>
                  <span>{t('admin.catalog.generate.failedCount')}: {generateResult.failed}</span>
                </div>
                <div className="overflow-x-auto border border-slate-200 rounded-2xl bg-white">
                  <table className="w-full text-left text-sm">
                    <thead>
                      <tr className="bg-slate-50/80 border-b border-slate-200">
                        <th className="px-4 py-3 font-semibold text-gray-600">{t('admin.catalog.form.template')}</th>
                        <th className="px-4 py-3 font-semibold text-gray-600">{t('admin.form.externalId')}</th>
                        <th className="px-4 py-3 font-semibold text-gray-600">{t('admin.catalog.generate.action')}</th>
                        <th className="px-4 py-3 font-semibold text-gray-600">{t('admin.catalog.generate.error')}</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-gray-100">
                      {(generateResult.items || []).map((it: any, idx: number) => (
                        <tr key={idx}>
                          <td className="px-4 py-3">
                            <NameCell name={templateNameById.get(it.layoutTemplateId)} id={it.layoutTemplateId} />
                          </td>
                          <td className="px-4 py-3 text-sm font-mono text-gray-600">{it.externalId}</td>
                          <td className="px-4 py-3 text-sm text-gray-600">{it.action}</td>
                          <td className="px-4 py-3 text-sm text-rose-600">{it.error || ''}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </section>
            ) : null}
          </div>
        ) : null}

        {!editingId && variantCreateMode === 'batch' ? null : (
          <>
        <section className="space-y-4">
          <SectionTitle>{t('admin.catalog.section.basic') || 'Basic'}</SectionTitle>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <Label required htmlFor="variant-external">{t('admin.form.externalId')}</Label>
              <Input id="variant-external" value={form.externalId || ''} disabled={!!editingId} onChange={(e) => setForm((p: any) => ({ ...p, externalId: e.target.value }))} />
            </div>
            <div>
              <Label required htmlFor="variant-name">{t('admin.form.name')}</Label>
              <Input id="variant-name" value={form.name || ''} onChange={(e) => setForm((p: any) => ({ ...p, name: e.target.value }))} />
            </div>
          </div>
        </section>
        <section className="space-y-4">
          <SectionTitle>{t('admin.catalog.section.config') || 'Config'}</SectionTitle>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div>
              <Label required htmlFor="variant-template">{t('admin.catalog.form.template')}</Label>
              <Select
                value={form.layoutTemplateExternalId || ''}
                onChange={(e) => setForm((p: any) => ({ ...p, layoutTemplateExternalId: e.target.value }))}
                options={[
                  { label: t('admin.form.selectPlaceholder'), value: '', disabled: true, hidden: true },
                  ...allTemplates.map(x => ({ label: x.name, value: x.id })),
                ]}
              />
            </div>
            <div>
              <Label required htmlFor="variant-preset">{t('admin.catalog.form.preset')}</Label>
              <Select
                value={form.presetExternalId || ''}
                onChange={(e) => setForm((p: any) => ({ ...p, presetExternalId: e.target.value }))}
                options={[
                  { label: t('admin.form.selectPlaceholder'), value: '', disabled: true, hidden: true },
                  ...allPresets.map(x => ({ label: x.Name, value: x.ExternalID })),
                ]}
              />
            </div>
            <div>
              <Label required htmlFor="variant-role">{t('admin.catalog.form.role')}</Label>
              <Select
                value={form.roleExternalId || ''}
                onChange={(e) => setForm((p: any) => ({ ...p, roleExternalId: e.target.value }))}
                options={[
                  { label: t('admin.form.selectPlaceholder'), value: '', disabled: true, hidden: true },
                  ...allRoles.map(x => ({ label: x.Name, value: x.ExternalID })),
                ]}
              />
            </div>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <Label htmlFor="variant-usage">{t('admin.form.usageCount')}</Label>
              <Input id="variant-usage" type="number" value={String(form.usageCount ?? 0)} onChange={(e) => setForm((p: any) => ({ ...p, usageCount: Number(e.target.value) }))} />
            </div>
            <div className="flex items-end gap-8">
              <Checkbox label={t('admin.form.isPremium')} checked={!!form.isPremium} onChange={(checked) => setForm((p: any) => ({ ...p, isPremium: checked }))} />
              <Checkbox label={t('admin.catalog.form.active')} checked={!!form.isActive} onChange={(checked) => setForm((p: any) => ({ ...p, isActive: checked }))} />
            </div>
          </div>
        </section>
        <section className="space-y-4">
          <SectionTitle>{t('admin.catalog.section.tags') || 'Tags'}</SectionTitle>
          <div>
            <Label>{t('admin.form.tags')}</Label>
            <TagInput tags={parseTags(form.tags)} onChange={(tags) => setForm((p: any) => ({ ...p, tags }))} />
          </div>
        </section>
          </>
        )}
      </div>
    );
  };

  const Pagination = () => {
    const totalPages = Math.max(1, Math.ceil(total / pageSize));
    return (
      <div className="flex items-center justify-between mt-4 text-sm text-gray-500">
        <div>
          {t('admin.pagination.total') || 'Total'}: {total}
        </div>
        <div className="flex items-center gap-2">
          <button className="px-3 py-1 border rounded" disabled={page <= 1} onClick={() => setPage(p => Math.max(1, p - 1))}>
            {t('admin.pagination.prev') || 'Prev'}
          </button>
          <div>
            {page} / {totalPages}
          </div>
          <button className="px-3 py-1 border rounded" disabled={page >= totalPages} onClick={() => setPage(p => Math.min(totalPages, p + 1))}>
            {t('admin.pagination.next') || 'Next'}
          </button>
          <select className="border rounded px-2 py-1" value={pageSize} onChange={(e) => { setPage(1); setPageSize(Number(e.target.value)); }}>
            {[10, 20, 50, 100].map(x => (
              <option key={x} value={x}>{x}</option>
            ))}
          </select>
        </div>
      </div>
    );
  };

  const renderTable = () => {
    if (tab === 'jobs') {
      const keyword = q.trim().toLowerCase();
      const rolesByCategory = new Map<string, AdminJobRole[]>();
      for (const r of allRoles) {
        const list = rolesByCategory.get(r.CategoryExternalID) || [];
        list.push(r);
        rolesByCategory.set(r.CategoryExternalID, list);
      }
      for (const [k, list] of rolesByCategory.entries()) {
        list.sort((a, b) => (a.OrderNum ?? 0) - (b.OrderNum ?? 0) || a.Name.localeCompare(b.Name));
        rolesByCategory.set(k, list);
      }
      const categoriesByParent = new Map<string, AdminJobCategory[]>();
      for (const c of allCategories) {
        const parent = c.ParentExternalID || '';
        const list = categoriesByParent.get(parent) || [];
        list.push(c);
        categoriesByParent.set(parent, list);
      }
      for (const [k, list] of categoriesByParent.entries()) {
        list.sort((a, b) => (a.OrderNum ?? 0) - (b.OrderNum ?? 0) || a.Name.localeCompare(b.Name));
        categoriesByParent.set(k, list);
      }

      const rootCategories = categoriesByParent.get('') || [];

      const matchCategory = (c: AdminJobCategory) => {
        if (!keyword) return true;
        return (c.Name || '').toLowerCase().includes(keyword) || (c.ExternalID || '').toLowerCase().includes(keyword);
      };
      const matchRole = (r: AdminJobRole) => {
        if (!keyword) return true;
        const name = (r.Name || '').toLowerCase();
        const id = (r.ExternalID || '').toLowerCase();
        const tags = (r.Tags || '').toLowerCase();
        return name.includes(keyword) || id.includes(keyword) || tags.includes(keyword);
      };
      const roleCountUnderCategory = (categoryId: string) => (rolesByCategory.get(categoryId) || []).length;
      const totalRoleCountUnderRoot = (rootId: string) => {
        const children = categoriesByParent.get(rootId) || [];
        let total = roleCountUnderCategory(rootId);
        for (const child of children) total += roleCountUnderCategory(child.ExternalID);
        return total;
      };

      const filteredRootCategories = keyword
        ? rootCategories.filter((root) => {
            if (matchCategory(root)) return true;
            const rootRoles = rolesByCategory.get(root.ExternalID) || [];
            if (rootRoles.some(matchRole)) return true;
            const children = categoriesByParent.get(root.ExternalID) || [];
            for (const child of children) {
              if (matchCategory(child)) return true;
              const childRoles = rolesByCategory.get(child.ExternalID) || [];
              if (childRoles.some(matchRole)) return true;
            }
            return false;
          })
        : rootCategories;

      return (
        <div className="space-y-3">
          <div className="flex items-center gap-2">
            <Button
              variant="outline"
              onClick={() => {
                const next: Record<string, boolean> = {};
                for (const c of filteredRootCategories) {
                  next[c.ExternalID] = true;
                  for (const child of categoriesByParent.get(c.ExternalID) || []) next[child.ExternalID] = true;
                  next[`${c.ExternalID}__ungrouped`] = true;
                }
                setExpandedCategories(next);
              }}
            >
              {t('admin.catalog.actions.expandAll') || 'Expand all'}
            </Button>
            <Button variant="outline" onClick={() => setExpandedCategories({})}>
              {t('admin.catalog.actions.collapseAll') || 'Collapse all'}
            </Button>
          </div>

          <TableCard>
            <div className="overflow-x-auto">
              <table className="w-full text-left table-fixed text-sm">
                <colgroup>
                  <col className="w-[300px]" />
                  <col className="w-[320px]" />
                  <col className="w-[80px]" />
                  <col className="w-[90px]" />
                  <col className="w-[140px]" />
                </colgroup>
                <thead className="sticky top-0 z-20 bg-gray-50 shadow-sm">
                  <tr>
                    <th className="px-4 py-4 font-semibold text-gray-600">{t('admin.catalog.jobs.name') || '岗位'}</th>
                    <th className="px-4 py-4 font-semibold text-gray-600">{t('admin.form.tags')}</th>
                    <th className="px-4 py-4 font-semibold text-gray-600">{t('admin.catalog.form.order')}</th>
                    <th className="px-4 py-4 font-semibold text-gray-600">{t('admin.catalog.form.active')}</th>
                    <th className="px-4 py-4 font-semibold text-gray-600 text-right sticky right-0 bg-gray-50">{t('admin.columns.actions')}</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-100">
                  {filteredRootCategories.map((c) => {
                    const children = categoriesByParent.get(c.ExternalID) || [];
                    const rootRoles = rolesByCategory.get(c.ExternalID) || [];
                    const expanded = keyword ? true : !!expandedCategories[c.ExternalID];
                    return (
                      <React.Fragment key={c.ExternalID}>
                        <tr
                          className="hover:bg-blue-50/30 transition-colors cursor-pointer"
                          onClick={() => setExpandedCategories((prev) => ({ ...prev, [c.ExternalID]: !prev[c.ExternalID] }))}
                          onKeyDown={(e) => {
                            if (e.key === 'Enter' || e.key === ' ') {
                              e.preventDefault();
                              setExpandedCategories((prev) => ({ ...prev, [c.ExternalID]: !prev[c.ExternalID] }));
                            }
                          }}
                          tabIndex={0}
                        >
                          <td className="px-4 py-4">
                            <div className="flex items-center gap-2">
                              <button
                                type="button"
                                className="p-1 rounded hover:bg-slate-100 text-slate-500"
                                onClick={(e) => {
                                  e.stopPropagation();
                                  setExpandedCategories((prev) => ({ ...prev, [c.ExternalID]: !prev[c.ExternalID] }));
                                }}
                                title={expanded ? (t('admin.catalog.actions.collapse') || 'Collapse') : (t('admin.catalog.actions.expand') || 'Expand')}
                              >
                                {expanded ? (
                                  <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="m6 9 6 6 6-6"/></svg>
                                ) : (
                                  <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="m9 18 6-6-6-6"/></svg>
                                )}
                              </button>
                              <div className="min-w-0 flex-1">
                                <NameCell name={c.Name} id={c.ExternalID} />
                              </div>
                            </div>
                          </td>
                          <td className="px-4 py-4 text-sm text-slate-400 whitespace-nowrap">
                            {children.length ? `${children.length} ${t('admin.catalog.jobs.subcategories') || '小类'}` : '-'}
                            <span className="mx-2 text-slate-200">|</span>
                            {(totalRoleCountUnderRoot(c.ExternalID) || 0) ? `${totalRoleCountUnderRoot(c.ExternalID)} ${t('admin.catalog.jobs.roles') || '岗位'}` : '-'}
                          </td>
                          <td className="px-4 py-4 text-sm text-gray-500">{c.OrderNum ?? 0}</td>
                          <td className="px-4 py-4">
                            <BoolBadge value={!!c.IsActive} yes={t('admin.catalog.enabled') || 'Enabled'} no={t('admin.catalog.disabled') || 'Disabled'} />
                          </td>
                          <td className="px-4 py-4 text-right sticky right-0 bg-white" onClick={(e) => e.stopPropagation()}>
                            <div className="flex items-center justify-end gap-1">
                              <button
                                type="button"
                                onClick={() => openCreateSubCategory(c.ExternalID)}
                                className="p-2 text-slate-400 hover:text-indigo-600 hover:bg-indigo-50 rounded-lg transition-all"
                                title={t('admin.catalog.actions.addSubcategory') || 'Add subcategory'}
                              >
                                <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M5 12h14"/><path d="M12 5v14"/></svg>
                              </button>
                              <IconButton title={t('common.edit') || 'Edit'} onClick={() => openEditCategory(c)} kind="edit" />
                              <IconButton title={t('common.delete') || 'Delete'} onClick={() => removeItem('category', c.ExternalID)} kind="delete" />
                            </div>
                          </td>
                        </tr>
                        {expanded ? (
                          <>
                            {children.map((child) => {
                              const childRoles = rolesByCategory.get(child.ExternalID) || [];
                              const childExpanded = keyword ? true : !!expandedCategories[child.ExternalID];
                              return (
                                <React.Fragment key={child.ExternalID}>
                                  <tr
                                    className="hover:bg-blue-50/20 transition-colors bg-white cursor-pointer"
                                    onClick={() => setExpandedCategories((prev) => ({ ...prev, [child.ExternalID]: !prev[child.ExternalID] }))}
                                    onKeyDown={(e) => {
                                      if (e.key === 'Enter' || e.key === ' ') {
                                        e.preventDefault();
                                        setExpandedCategories((prev) => ({ ...prev, [child.ExternalID]: !prev[child.ExternalID] }));
                                      }
                                    }}
                                    tabIndex={0}
                                  >
                                    <td className="px-4 py-4">
                                      <div className="flex items-center gap-2 pl-5">
                                        <button
                                          type="button"
                                          className="p-1 rounded hover:bg-slate-100 text-slate-400"
                                          onClick={(e) => {
                                            e.stopPropagation();
                                            setExpandedCategories((prev) => ({ ...prev, [child.ExternalID]: !prev[child.ExternalID] }));
                                          }}
                                          title={childExpanded ? (t('admin.catalog.actions.collapse') || 'Collapse') : (t('admin.catalog.actions.expand') || 'Expand')}
                                        >
                                          {childExpanded ? (
                                            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="m6 9 6 6 6-6"/></svg>
                                          ) : (
                                            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="m9 18 6-6-6-6"/></svg>
                                          )}
                                        </button>
                                        <div className="min-w-0 flex-1">
                                          <NameCell name={child.Name} id={child.ExternalID} />
                                        </div>
                                      </div>
                                    </td>
                                    <td className="px-4 py-4 text-sm text-slate-400">{childRoles.length}</td>
                                    <td className="px-4 py-4 text-sm text-gray-500">{child.OrderNum ?? 0}</td>
                                    <td className="px-4 py-4">
                                      <BoolBadge value={!!child.IsActive} yes={t('admin.catalog.enabled') || 'Enabled'} no={t('admin.catalog.disabled') || 'Disabled'} />
                                    </td>
                                    <td className="px-4 py-4 text-right sticky right-0 bg-white" onClick={(e) => e.stopPropagation()}>
                                      <div className="flex items-center justify-end gap-1">
                                        <button
                                          type="button"
                                          onClick={() => openCreateRole(child.ExternalID)}
                                          className="p-2 text-slate-400 hover:text-indigo-600 hover:bg-indigo-50 rounded-lg transition-all"
                                          title={t('admin.catalog.actions.addRole') || 'Add role'}
                                        >
                                          <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M5 12h14"/><path d="M12 5v14"/></svg>
                                        </button>
                                        <IconButton title={t('common.edit') || 'Edit'} onClick={() => openEditCategory(child)} kind="edit" />
                                        <IconButton title={t('common.delete') || 'Delete'} onClick={() => removeItem('category', child.ExternalID)} kind="delete" />
                                      </div>
                                    </td>
                                  </tr>
                                  {childExpanded
                                    ? childRoles.map((r) => (
                                        <tr key={r.ExternalID} className="hover:bg-blue-50/10 transition-colors bg-white">
                                          <td className="px-4 py-4">
                                            <div className="flex items-start gap-2 pl-12">
                                              <div className="w-4 h-4 text-slate-200">
                                                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M9 18h6"/><path d="M10 22h4"/><path d="M12 2v8"/><path d="M12 10l4 4"/><path d="M12 10l-4 4"/></svg>
                                              </div>
                                              <div className="min-w-0 flex-1">
                                                <NameCell name={r.Name} id={r.ExternalID} />
                                              </div>
                                            </div>
                                          </td>
                                          <td className="px-4 py-4">
                                            <div className="flex flex-wrap gap-1.5">
                                              {parseTags(r.Tags).slice(0, 6).map((tag) => (
                                                <span
                                                  key={tag}
                                                  title={tag}
                                                  className="px-2 py-0.5 bg-white border border-slate-200 rounded text-[11px] text-slate-600 font-medium hover:border-indigo-300 hover:bg-indigo-50 transition-colors max-w-[200px] truncate"
                                                >
                                                  {tag}
                                                </span>
                                              ))}
                                              {parseTags(r.Tags).length > 6 ? (
                                                <span className="px-2 py-0.5 bg-slate-50 border border-slate-200 rounded text-[11px] text-slate-500 font-medium">
                                                  +{parseTags(r.Tags).length - 6}
                                                </span>
                                              ) : null}
                                            </div>
                                          </td>
                                          <td className="px-4 py-4 text-sm text-gray-500">{r.OrderNum ?? 0}</td>
                                          <td className="px-4 py-4">
                                            <BoolBadge value={!!r.IsActive} yes={t('admin.catalog.enabled') || 'Enabled'} no={t('admin.catalog.disabled') || 'Disabled'} />
                                          </td>
                                          <td className="px-4 py-4 text-right sticky right-0 bg-white" onClick={(e) => e.stopPropagation()}>
                                            <div className="flex items-center justify-end gap-1">
                                              <IconButton title={t('common.edit') || 'Edit'} onClick={() => openEditRole(r)} kind="edit" />
                                              <IconButton title={t('common.delete') || 'Delete'} onClick={() => removeItem('role', r.ExternalID)} kind="delete" />
                                            </div>
                                          </td>
                                        </tr>
                                      ))
                                    : null}
                                </React.Fragment>
                              );
                            })}
                            {rootRoles.length ? (
                              (() => {
                                const key = `${c.ExternalID}__ungrouped`;
                                const ungroupedExpanded = keyword ? true : !!expandedCategories[key];
                                return (
                                  <React.Fragment key={key}>
                                    <tr
                                      className="hover:bg-blue-50/20 transition-colors bg-white cursor-pointer"
                                      onClick={() => setExpandedCategories((prev) => ({ ...prev, [key]: !prev[key] }))}
                                      onKeyDown={(e) => {
                                        if (e.key === 'Enter' || e.key === ' ') {
                                          e.preventDefault();
                                          setExpandedCategories((prev) => ({ ...prev, [key]: !prev[key] }));
                                        }
                                      }}
                                      tabIndex={0}
                                    >
                                      <td className="px-4 py-4">
                                        <div className="flex items-center gap-2 pl-5">
                                          <button
                                            type="button"
                                            className="p-1 rounded hover:bg-slate-100 text-slate-400"
                                            onClick={(e) => {
                                              e.stopPropagation();
                                              setExpandedCategories((prev) => ({ ...prev, [key]: !prev[key] }));
                                            }}
                                          >
                                            {ungroupedExpanded ? (
                                              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="m6 9 6 6 6-6"/></svg>
                                            ) : (
                                              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="m9 18 6-6-6-6"/></svg>
                                            )}
                                          </button>
                                          <div className="text-sm font-medium text-slate-500">{t('admin.catalog.jobs.ungrouped') || '未分组'}</div>
                                        </div>
                                      </td>
                                      <td className="px-4 py-4 text-sm text-slate-400">{rootRoles.length}</td>
                                      <td className="px-4 py-4 text-sm text-gray-300">-</td>
                                      <td className="px-4 py-4 text-sm text-gray-300">-</td>
                                      <td className="px-4 py-4 sticky right-0 bg-white" onClick={(e) => e.stopPropagation()} />
                                    </tr>
                                    {ungroupedExpanded
                                      ? rootRoles.map((r) => (
                                          <tr key={r.ExternalID} className="hover:bg-blue-50/10 transition-colors bg-white">
                                            <td className="px-4 py-4">
                                              <div className="flex items-start gap-2 pl-12">
                                                <div className="w-4 h-4 text-slate-200">
                                                  <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M9 18h6"/><path d="M10 22h4"/><path d="M12 2v8"/><path d="M12 10l4 4"/><path d="M12 10l-4 4"/></svg>
                                                </div>
                                                <div className="min-w-0 flex-1">
                                                  <NameCell name={r.Name} id={r.ExternalID} />
                                                </div>
                                              </div>
                                            </td>
                                            <td className="px-4 py-4">
                                              <div className="flex flex-wrap gap-1.5">
                                                {parseTags(r.Tags).slice(0, 6).map((tag) => (
                                                  <span
                                                    key={tag}
                                                    title={tag}
                                                    className="px-2 py-0.5 bg-white border border-slate-200 rounded text-[11px] text-slate-600 font-medium hover:border-indigo-300 hover:bg-indigo-50 transition-colors max-w-[200px] truncate"
                                                  >
                                                    {tag}
                                                  </span>
                                                ))}
                                                {parseTags(r.Tags).length > 6 ? (
                                                  <span className="px-2 py-0.5 bg-slate-50 border border-slate-200 rounded text-[11px] text-slate-500 font-medium">
                                                    +{parseTags(r.Tags).length - 6}
                                                  </span>
                                                ) : null}
                                              </div>
                                            </td>
                                            <td className="px-4 py-4 text-sm text-gray-500">{r.OrderNum ?? 0}</td>
                                            <td className="px-4 py-4">
                                              <BoolBadge value={!!r.IsActive} yes={t('admin.catalog.enabled') || 'Enabled'} no={t('admin.catalog.disabled') || 'Disabled'} />
                                            </td>
                                            <td className="px-4 py-4 text-right sticky right-0 bg-white" onClick={(e) => e.stopPropagation()}>
                                              <div className="flex items-center justify-end gap-1">
                                                <IconButton title={t('common.edit') || 'Edit'} onClick={() => openEditRole(r)} kind="edit" />
                                                <IconButton title={t('common.delete') || 'Delete'} onClick={() => removeItem('role', r.ExternalID)} kind="delete" />
                                              </div>
                                            </td>
                                          </tr>
                                        ))
                                      : null}
                                  </React.Fragment>
                                );
                              })()
                            ) : null}
                          </>
                        ) : null}
                      </React.Fragment>
                    );
                  })}
                </tbody>
              </table>
            </div>
          </TableCard>
        </div>
      );
    }
    if (tab === 'presets') {
      return (
        <TableCard>
          <DataTable<AdminContentPreset>
            data={presets}
            getRowKey={(row) => row.ExternalID}
            emptyState={{
              title: t('admin.catalog.empty') || '暂无数据',
              description: t('admin.catalog.emptyDesc') || '点击右上角创建开始新增',
            }}
            columns={[
              {
                key: 'externalId',
                label: 'ExternalID',
                minWidth: 200,
                nowrap: true,
                render: (p) => <TextCell text={p.ExternalID} className="font-mono text-gray-600" />,
              },
              { key: 'name', label: t('admin.form.name'), minWidth: 280, render: (p) => <TextCell text={p.Name} className="text-gray-900 font-medium" /> },
              { key: 'language', label: t('admin.catalog.form.language'), minWidth: 120, nowrap: true, render: (p) => <span className="text-sm text-gray-600">{p.Language || '-'}</span> },
              { key: 'role', label: t('admin.catalog.form.role'), minWidth: 240, render: (p) => <NameCell name={p.RoleExternalID ? roleNameById.get(p.RoleExternalID) : ''} id={p.RoleExternalID || ''} /> },
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
                    <IconButton title={t('common.delete') || 'Delete'} onClick={() => removeItem('preset', p.ExternalID)} kind="delete" />
                  </div>
                ),
              },
            ]}
          />
        </TableCard>
      );
    }
    return (
      <TableCard>
          <DataTable<AdminTemplateVariant>
            data={variants}
            getRowKey={(row) => row.ExternalID}
            emptyState={{
              title: t('admin.catalog.empty') || '暂无数据',
              description: t('admin.catalog.emptyDesc') || '点击右上角创建开始新增',
            }}
            columns={[
              {
                key: 'externalId',
                label: 'ExternalID',
                minWidth: 200,
                nowrap: true,
                render: (v) => (
                  <div className="text-sm font-mono text-gray-600 truncate whitespace-nowrap" title={v.ExternalID}>
                    {v.ExternalID}
                  </div>
                ),
              },
              {
                key: 'name',
                label: t('admin.form.name'),
                minWidth: 240,
                render: (v) => (
                  <div className="min-w-0">
                    <div className="text-sm font-medium text-gray-900 truncate whitespace-nowrap" title={v.Name}>
                      {v.Name}
                    </div>
                  </div>
                ),
              },
              {
                key: 'template',
                label: t('admin.catalog.form.template'),
                minWidth: 240,
                render: (v) => <NameCell name={templateNameById.get(v.LayoutTemplateExternalID)} id={v.LayoutTemplateExternalID} />,
              },
              {
                key: 'preset',
                label: t('admin.catalog.form.preset'),
                minWidth: 220,
                render: (v) => <NameCell name={presetNameById.get(v.PresetExternalID)} id={v.PresetExternalID} />,
              },
              {
                key: 'role',
                label: t('admin.catalog.form.role'),
                minWidth: 220,
                render: (v) => <NameCell name={roleNameById.get(v.RoleExternalID)} id={v.RoleExternalID} />,
              },
              {
                key: 'tags',
                label: t('admin.form.tags'),
                minWidth: 240,
                render: (v) => {
                  const tags = parseTags(v.Tags);
                  return (
                    <div className="flex flex-wrap gap-1.5 py-1">
                      {tags.slice(0, 6).map((tag) => (
                        <span
                          key={tag}
                          title={tag}
                          className="px-2 py-0.5 bg-white border border-slate-200 rounded text-[11px] text-slate-600 font-medium hover:border-indigo-300 hover:bg-indigo-50 transition-colors max-w-[200px] truncate"
                        >
                          {tag}
                        </span>
                      ))}
                      {tags.length > 6 ? (
                        <span className="px-2 py-0.5 bg-slate-50 border border-slate-200 rounded text-[11px] text-slate-500 font-medium">
                          +{tags.length - 6}
                        </span>
                      ) : null}
                    </div>
                  );
                },
              },
              {
                key: 'usage',
                label: t('admin.form.usageCount'),
                minWidth: 110,
                nowrap: true,
                render: (v) => <span className="text-sm text-gray-600">{v.UsageCount ?? 0}</span>,
              },
              {
                key: 'isPremium',
                label: t('admin.form.isPremium'),
                minWidth: 110,
                nowrap: true,
                render: (v) => <BoolBadge value={!!v.IsPremium} yes={t('admin.catalog.yes') || 'Yes'} no={t('admin.catalog.no') || 'No'} />,
              },
              {
                key: 'isActive',
                label: t('admin.catalog.form.active'),
                minWidth: 120,
                nowrap: true,
                render: (v) => <BoolBadge value={!!v.IsActive} yes={t('admin.catalog.enabled') || 'Enabled'} no={t('admin.catalog.disabled') || 'Disabled'} />,
              },
              {
                key: 'actions',
                label: t('admin.columns.actions'),
                minWidth: 140,
                fixed: 'right',
                headerClassName: 'text-right',
                cellClassName: 'text-right',
                render: (v) => (
                  <div className="flex items-center justify-end gap-1">
                    <IconButton title={t('common.edit') || 'Edit'} onClick={() => openEditVariant(v)} kind="edit" />
                    <IconButton title={t('common.delete') || 'Delete'} onClick={() => removeItem('variant', v.ExternalID)} kind="delete" />
                  </div>
                ),
              },
            ]}
          />
      </TableCard>
    );
  };

  return (
    <div className="flex-1 flex flex-col bg-white rounded-3xl m-2 overflow-hidden shadow-sm border border-gray-100">
      <div className="px-10 pt-10 pb-6">
        <div className="flex justify-between items-center">
          <h1 className="text-4xl font-bold text-gray-800 tracking-tight">{t('admin.menu.catalog')}</h1>
          <div className="flex items-center gap-2">
            <Button variant="outline" onClick={() => load()} disabled={loading}>{t('common.refresh') || 'Refresh'}</Button>
            <Button onClick={openCreate}>{t('admin.actions.create') || 'Create'}</Button>
          </div>
        </div>
        <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4 border-b border-gray-100 mt-6 pb-4">
          <nav className="inline-flex items-center gap-1 p-1 bg-gray-100 rounded-xl self-start">
            {tabs.map(x => (
              <button
                key={x.key}
                type="button"
                aria-pressed={tab === x.key}
                className={`flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-semibold transition-all duration-200 ${tab === x.key ? 'bg-white text-blue-600 shadow-sm' : 'text-gray-500 hover:text-gray-800 hover:bg-gray-200/60'}`}
                onClick={() => setTab(x.key)}
              >
                <x.Icon className="w-4 h-4" />
                {x.label}
              </button>
            ))}
          </nav>
          <div className="flex w-full md:w-auto md:justify-end">
            <input
              type="text"
              placeholder={t('admin.keyword')}
              className="w-full md:w-[360px] px-4 py-2 bg-white border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all"
              value={q}
              onChange={(e) => setQ(e.target.value)}
            />
          </div>
        </div>
      </div>
      <div className="flex-1 overflow-y-auto px-10 pb-10">
        <div className="overflow-x-auto">
          {renderTable()}
        </div>
        {tab === 'jobs' ? null : <Pagination />}
      </div>

      <Modal
        isOpen={showForm}
        onClose={() => setShowForm(false)}
        title={editingId ? (t('common.edit') || 'Edit') : (t('admin.actions.create') || 'Create')}
        size="lg"
        compact
      >
        <div className="flex flex-col max-h-[70vh]">
          <div className="flex-1 overflow-y-auto pr-1 space-y-5">
            {renderFormFields()}
          </div>
          <div className="mt-4 pt-3 border-t border-slate-100 bg-slate-50/50 flex items-center justify-end gap-3">
            <Button variant="outline" onClick={() => setShowForm(false)} disabled={saving || generating}>{t('common.cancel') || 'Cancel'}</Button>
            {formKind === 'variant' && !editingId && variantCreateMode === 'batch' ? (
              <Button onClick={runGenerateVariants} disabled={generating}>
                {t('admin.catalog.generate.submit')}
              </Button>
            ) : (
              <Button onClick={save} disabled={saving}>{t('common.save') || 'Save'}</Button>
            )}
          </div>
        </div>
      </Modal>
    </div>
  );
};
