import React, { useEffect, useMemo, useState } from 'react';
import { Button } from '../../components/ui/Button';
import { Checkbox, Input, Label, Select } from '../../components/ui/Form';
import { TableCard } from '../../components/ui/TableCard';
import { Modal } from '../../components/ui/Modal';
import { useConfirm } from '../../components/ui/ConfirmDialog';
import { useToast } from '../../components/ui/Toast';
import { useLanguage } from '../../contexts/LanguageContext';
import { Search } from 'lucide-react';
import { AnimatePresence, motion } from 'framer-motion';
import {
  adminCreateJobCategory,
  adminCreateJobRole,
  adminDeleteJobCategory,
  adminDeleteJobRole,
  adminListJobCategories,
  adminListJobRoles,
  adminPatchJobCategory,
  adminPatchJobRole,
  AdminJobCategory,
  AdminJobRole,
} from '../../services/adminService';

type FormKind = 'category' | 'role';

type CatalogPageProps = {
  embedded?: boolean;
};

export const CatalogPage: React.FC<CatalogPageProps> = ({ embedded }) => {
  const { t, language } = useLanguage();
  const { showToast } = useToast();
  const confirm = useConfirm();

  const [formKind, setFormKind] = useState<FormKind>('category');
  const [q, setQ] = useState('');
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);
  const [expandedCategories, setExpandedCategories] = useState<Record<string, boolean>>({});

  const [allCategories, setAllCategories] = useState<AdminJobCategory[]>([]);
  const [allRoles, setAllRoles] = useState<AdminJobRole[]>([]);

  const displayCategoryName = (c: AdminJobCategory) => (c.names?.[language] || c.names?.zh || c.names?.en || '').trim() || '-';
  const displayRoleName = (r: AdminJobRole) => (r.names?.[language] || r.names?.zh || r.names?.en || '').trim() || '-';

  const categoryNameById = useMemo(() => {
    const m = new Map<number, string>();
    for (const c of allCategories) m.set(c.id, displayCategoryName(c));
    return m;
  }, [allCategories, language]);

  const NameCell: React.FC<{ name?: string }> = ({ name }) => {
    const display = (name || '').trim() || '-';
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

  const load = async () => {
    setLoading(true);
    try {
      const cats = await adminListJobCategories({ page: '1', pageSize: '1000' });
      setAllCategories(cats.items || []);
      const rs = await adminListJobRoles({ page: '1', pageSize: '2000' });
      setAllRoles(rs.items || []);
    } catch {
      showToast(t('admin.msg.loadFailed'), 'error');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    load();
  }, []);

  const [showForm, setShowForm] = useState(false);
  const [editingId, setEditingId] = useState<number | null>(null);
  const [form, setForm] = useState<any>({});
  const openCreateCategory = () => {
    setEditingId(null);
    setFormKind('category');
    setForm({ names: { zh: '', en: '' }, parentId: null as number | null, orderNum: 0, isActive: true });
    setShowForm(true);
  };
  const openCreateSubCategory = (parentId: number) => {
    setEditingId(null);
    setFormKind('category');
    setForm({ names: { zh: '', en: '' }, parentId, orderNum: 0, isActive: true });
    setShowForm(true);
  };
  const openCreateRole = (categoryId: number) => {
    setEditingId(null);
    setFormKind('role');
    setForm({ categoryId, names: { zh: '', en: '' }, orderNum: 0, isActive: true });
    setShowForm(true);
  };

  const openEditCategory = (row: any) => {
    setEditingId(row.id);
    setFormKind('category');
    setForm({
      names: { zh: row.names?.zh ?? '', en: row.names?.en ?? '' },
      parentId: row.parentId ?? null,
      orderNum: row.orderNum ?? 0,
      isActive: !!row.isActive,
    });
    setShowForm(true);
  };
  const openEditRole = (row: any) => {
    setEditingId(row.id);
    setFormKind('role');
    setForm({
      categoryId: row.categoryId,
      names: { zh: row.names?.zh ?? '', en: row.names?.en ?? '' },
      orderNum: row.orderNum ?? 0,
      isActive: !!row.isActive,
    });
    setShowForm(true);
  };

  const openCreate = () => {
    openCreateCategory();
  };

  const save = async () => {
    if (saving) return;
    setSaving(true);
    try {
      const zhName = String(form?.names?.zh || '').trim();
      const enName = String(form?.names?.en || '').trim();
      const names: Record<string, string> = {};
      if (zhName) names.zh = zhName;
      if (enName) names.en = enName;

      if (formKind === 'category') {
        if (!Object.keys(names).length) {
          showToast(t('auth.error.fillAll'), 'error');
          setSaving(false);
          return;
        }
      }
      if (formKind === 'role') {
        if (!Number(form.categoryId || 0) || !Object.keys(names).length) {
          showToast(t('auth.error.fillAll'), 'error');
          setSaving(false);
          return;
        }
      }
      if (!editingId) {
        if (formKind === 'category') {
          await adminCreateJobCategory({ names, parentId: form.parentId ?? null, orderNum: form.orderNum ?? 0, isActive: !!form.isActive });
        } else if (formKind === 'role') {
          await adminCreateJobRole({ categoryId: Number(form.categoryId), names, orderNum: form.orderNum ?? 0, isActive: !!form.isActive });
        }
      } else {
        if (formKind === 'category') {
          await adminPatchJobCategory(editingId, { names, parentId: form.parentId ?? null, orderNum: form.orderNum ?? 0, isActive: !!form.isActive });
        } else if (formKind === 'role') {
          await adminPatchJobRole(editingId, { categoryId: Number(form.categoryId), names, orderNum: form.orderNum ?? 0, isActive: !!form.isActive });
        }
      }
      setShowForm(false);
      showToast(t('admin.msg.saveSuccess'), 'success');
      await load();
    } catch (e: any) {
      showToast(t('admin.msg.saveFailed'), 'error');
    } finally {
      setSaving(false);
    }
  };

  const removeItem = async (kind: FormKind, id: number) => {
    const ok = await confirm({ title: t('admin.confirm.delete'), message: t('admin.confirm.deleteMsg') });
    if (!ok) return;
    try {
      if (kind === 'category') await adminDeleteJobCategory(id);
      else if (kind === 'role') await adminDeleteJobRole(id);
      showToast(t('admin.msg.deleteSuccess'), 'success');
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
                <Label required htmlFor="cat-name-zh">{t('admin.form.name')} (ZH)</Label>
                <Input
                  id="cat-name-zh"
                  value={form?.names?.zh || ''}
                  onChange={(e) => setForm((p: any) => ({ ...p, names: { ...(p?.names || {}), zh: e.target.value } }))}
                />
              </div>
              <div>
                <Label htmlFor="cat-name-en">{t('admin.form.name')} (EN)</Label>
                <Input
                  id="cat-name-en"
                  value={form?.names?.en || ''}
                  onChange={(e) => setForm((p: any) => ({ ...p, names: { ...(p?.names || {}), en: e.target.value } }))}
                />
              </div>
            </div>
          </section>
          <section className="space-y-4">
            <SectionTitle>{t('admin.catalog.section.config') || 'Config'}</SectionTitle>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <Label htmlFor="cat-parent">{t('admin.catalog.form.parent')}</Label>
                {!editingId && Number(form.parentId || 0) ? (
                  <Input id="cat-parent" value={categoryNameById.get(Number(form.parentId)) || String(form.parentId)} disabled />
                ) : (
                  <Select
                    value={form.parentId == null ? '' : String(form.parentId)}
                    onChange={(e) => setForm((p: any) => ({ ...p, parentId: e.target.value ? Number(e.target.value) : null }))}
                    options={[
                      { label: t('admin.form.selectPlaceholder'), value: '', disabled: true, hidden: true },
                      ...allCategories.map(c => ({ label: displayCategoryName(c), value: String(c.id) })),
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
            <div>
              <Label required htmlFor="role-category">{t('admin.catalog.form.category')}</Label>
              {!editingId && Number(form.categoryId || 0) ? (
                <Input id="role-category" value={categoryNameById.get(Number(form.categoryId)) || String(form.categoryId)} disabled />
              ) : (
                <Select
                  value={form.categoryId ? String(form.categoryId) : ''}
                  onChange={(e) => setForm((p: any) => ({ ...p, categoryId: e.target.value ? Number(e.target.value) : 0 }))}
                  options={[
                    { label: t('admin.form.selectPlaceholder'), value: '', disabled: true, hidden: true },
                    ...allCategories.map(c => ({ label: displayCategoryName(c), value: String(c.id) })),
                  ]}
                />
              )}
            </div>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <Label required htmlFor="role-name-zh">{t('admin.form.name')} (ZH)</Label>
                <Input
                  id="role-name-zh"
                  value={form?.names?.zh || ''}
                  onChange={(e) => setForm((p: any) => ({ ...p, names: { ...(p?.names || {}), zh: e.target.value } }))}
                />
              </div>
              <div>
                <Label htmlFor="role-name-en">{t('admin.form.name')} (EN)</Label>
                <Input
                  id="role-name-en"
                  value={form?.names?.en || ''}
                  onChange={(e) => setForm((p: any) => ({ ...p, names: { ...(p?.names || {}), en: e.target.value } }))}
                />
              </div>
              <div>
                <Label htmlFor="role-order">{t('admin.catalog.form.order')}</Label>
                <Input id="role-order" type="number" value={String(form.orderNum ?? 0)} onChange={(e) => setForm((p: any) => ({ ...p, orderNum: Number(e.target.value) }))} />
              </div>
            </div>
            <Checkbox label={t('admin.catalog.form.active')} checked={!!form.isActive} onChange={(checked) => setForm((p: any) => ({ ...p, isActive: checked }))} />
          </section>
        </div>
      );
    }
    return null;
  };

  const renderTable = () => {
    const keyword = q.trim().toLowerCase();
    const rolesByCategory = new Map<number, AdminJobRole[]>();
    for (const r of allRoles) {
      const list = rolesByCategory.get(r.categoryId) || [];
      list.push(r);
      rolesByCategory.set(r.categoryId, list);
    }
    for (const [k, list] of rolesByCategory.entries()) {
      list.sort((a, b) => (a.orderNum ?? 0) - (b.orderNum ?? 0) || displayRoleName(a).localeCompare(displayRoleName(b)));
      rolesByCategory.set(k, list);
    }
    const categoriesByParent = new Map<number, AdminJobCategory[]>();
    for (const c of allCategories) {
      const parent = c.parentId ?? 0;
      const list = categoriesByParent.get(parent) || [];
      list.push(c);
      categoriesByParent.set(parent, list);
    }
    for (const [k, list] of categoriesByParent.entries()) {
      list.sort((a, b) => (a.orderNum ?? 0) - (b.orderNum ?? 0) || displayCategoryName(a).localeCompare(displayCategoryName(b)));
      categoriesByParent.set(k, list);
    }

    const rootCategories = categoriesByParent.get(0) || [];

    const matchCategory = (c: AdminJobCategory) => {
      if (!keyword) return true;
      return displayCategoryName(c).toLowerCase().includes(keyword);
    };
    const matchRole = (r: AdminJobRole) => {
      if (!keyword) return true;
      return displayRoleName(r).toLowerCase().includes(keyword);
    };
    const roleCountUnderCategory = (categoryId: number) => (rolesByCategory.get(categoryId) || []).length;
    const totalRoleCountUnderRoot = (rootId: number) => {
      const children = categoriesByParent.get(rootId) || [];
      let total = roleCountUnderCategory(rootId);
      for (const child of children) total += roleCountUnderCategory(child.id);
      return total;
    };

    const filteredRootCategories = keyword
      ? rootCategories.filter((root) => {
          if (matchCategory(root)) return true;
          const rootRoles = rolesByCategory.get(root.id) || [];
          if (rootRoles.some(matchRole)) return true;
          const children = categoriesByParent.get(root.id) || [];
          for (const child of children) {
            if (matchCategory(child)) return true;
            const childRoles = rolesByCategory.get(child.id) || [];
            if (childRoles.some(matchRole)) return true;
          }
          return false;
        })
      : rootCategories;

    return (
      <TableCard className="flex flex-col h-full">
        <div className="flex-1 overflow-y-scroll overflow-x-hidden no-scrollbar" style={{ scrollbarGutter: 'stable' }}>
          <table className="w-full text-left table-fixed text-sm">
            <colgroup>
              <col className="w-[300px]" />
              <col className="w-[240px]" />
              <col className="w-[80px]" />
              <col className="w-[90px]" />
              <col className="w-[140px]" />
            </colgroup>
            <thead className="sticky top-0 z-20 bg-gray-50 shadow-sm">
              <tr>
                <th className="px-4 py-4 font-semibold text-gray-600">{t('admin.catalog.jobs.name') || '岗位'}</th>
                <th className="px-4 py-4 font-semibold text-gray-600">{t('admin.catalog.jobs.stats')}</th>
                <th className="px-4 py-4 font-semibold text-gray-600">{t('admin.catalog.form.order')}</th>
                <th className="px-4 py-4 font-semibold text-gray-600">{t('admin.catalog.form.active')}</th>
                <th className="px-4 py-4 font-semibold text-gray-600 text-right sticky right-0 bg-gray-50">{t('admin.columns.actions')}</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100">
              {filteredRootCategories.map((c) => {
                const children = categoriesByParent.get(c.id) || [];
                const rootRoles = rolesByCategory.get(c.id) || [];
                const expandedKey = String(c.id);
                const expanded = keyword ? true : !!expandedCategories[expandedKey];
                const expandedRows = expanded
                  ? [
                      ...children.flatMap((child, idx) => {
                        const childRoles = rolesByCategory.get(child.id) || [];
                        const childExpandedKey = String(child.id);
                        const childExpanded = keyword ? true : !!expandedCategories[childExpandedKey];
                        const rows: React.ReactNode[] = [
                          <motion.tr
                            key={`category_${child.id}`}
                            initial={{ opacity: 0 }}
                            animate={{ opacity: 1 }}
                            exit={{ opacity: 0 }}
                            transition={{ duration: 0.16, delay: Math.min(0.06, idx * 0.01) }}
                            className="hover:bg-blue-50/20 transition-colors bg-white cursor-pointer"
                            onClick={() => setExpandedCategories((prev) => ({ ...prev, [childExpandedKey]: !prev[childExpandedKey] }))}
                            onKeyDown={(e) => {
                              if (e.key === 'Enter' || e.key === ' ') {
                                e.preventDefault();
                                setExpandedCategories((prev) => ({ ...prev, [childExpandedKey]: !prev[childExpandedKey] }));
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
                                    setExpandedCategories((prev) => ({ ...prev, [childExpandedKey]: !prev[childExpandedKey] }));
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
                                  <NameCell name={displayCategoryName(child)} />
                                </div>
                              </div>
                            </td>
                            <td className="px-4 py-4 text-sm text-slate-400">{childRoles.length}</td>
                            <td className="px-4 py-4 text-sm text-gray-500">{child.orderNum ?? 0}</td>
                            <td className="px-4 py-4">
                              <BoolBadge value={!!child.isActive} yes={t('admin.catalog.enabled') || 'Enabled'} no={t('admin.catalog.disabled') || 'Disabled'} />
                            </td>
                            <td className="px-4 py-4 text-right sticky right-0 bg-white" onClick={(e) => e.stopPropagation()}>
                              <div className="flex items-center justify-end gap-1">
                                <button
                                  type="button"
                                  onClick={() => openCreateRole(child.id)}
                                  className="p-2 text-slate-400 hover:text-indigo-600 hover:bg-indigo-50 rounded-lg transition-all"
                                  title={t('admin.catalog.actions.addRole') || 'Add role'}
                                >
                                  <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M5 12h14"/><path d="M12 5v14"/></svg>
                                </button>
                                <IconButton title={t('common.edit') || 'Edit'} onClick={() => openEditCategory(child)} kind="edit" />
                                <IconButton title={t('common.delete') || 'Delete'} onClick={() => removeItem('category', child.id)} kind="delete" />
                              </div>
                            </td>
                          </motion.tr>,
                        ];
                        if (childExpanded) {
                          rows.push(
                            ...childRoles.map((r, ridx) => (
                              <motion.tr
                                key={`role_${child.id}_${r.id}`}
                                initial={{ opacity: 0 }}
                                animate={{ opacity: 1 }}
                                exit={{ opacity: 0 }}
                                transition={{ duration: 0.16, delay: Math.min(0.06, ridx * 0.008) }}
                                className="hover:bg-blue-50/10 transition-colors bg-white"
                              >
                                <td className="px-4 py-4">
                                  <div className="flex items-start gap-2 pl-12">
                                    <div className="w-4 h-4 text-slate-200">
                                      <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M9 18h6"/><path d="M10 22h4"/><path d="M12 2v8"/><path d="M12 10l4 4"/><path d="M12 10l-4 4"/></svg>
                                    </div>
                                    <div className="min-w-0 flex-1">
                                      <NameCell name={displayRoleName(r)} />
                                    </div>
                                  </div>
                                </td>
                                <td className="px-4 py-4 text-sm text-gray-500">{r.orderNum ?? 0}</td>
                                <td className="px-4 py-4">
                                  <BoolBadge value={!!r.isActive} yes={t('admin.catalog.enabled') || 'Enabled'} no={t('admin.catalog.disabled') || 'Disabled'} />
                                </td>
                                <td className="px-4 py-4 text-right sticky right-0 bg-white" onClick={(e) => e.stopPropagation()}>
                                  <div className="flex items-center justify-end gap-1">
                                    <IconButton title={t('common.edit') || 'Edit'} onClick={() => openEditRole(r)} kind="edit" />
                                    <IconButton title={t('common.delete') || 'Delete'} onClick={() => removeItem('role', r.id)} kind="delete" />
                                  </div>
                                </td>
                              </motion.tr>
                            ))
                          );
                        }
                        return rows;
                      }),
                      ...(rootRoles.length
                        ? (() => {
                            const key = `${c.id}__ungrouped`;
                            const ungroupedExpanded = keyword ? true : !!expandedCategories[key];
                            const rows: React.ReactNode[] = [
                              <motion.tr
                                key={`ungrouped_${c.id}`}
                                initial={{ opacity: 0 }}
                                animate={{ opacity: 1 }}
                                exit={{ opacity: 0 }}
                                transition={{ duration: 0.16 }}
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
                              </motion.tr>,
                            ];
                            if (ungroupedExpanded) {
                              rows.push(
                                ...rootRoles.map((r, ridx) => (
                                  <motion.tr
                                    key={`root_role_${c.id}_${r.id}`}
                                    initial={{ opacity: 0 }}
                                    animate={{ opacity: 1 }}
                                    exit={{ opacity: 0 }}
                                    transition={{ duration: 0.16, delay: Math.min(0.06, ridx * 0.008) }}
                                    className="hover:bg-blue-50/10 transition-colors bg-white"
                                  >
                                    <td className="px-4 py-4">
                                      <div className="flex items-start gap-2 pl-12">
                                        <div className="w-4 h-4 text-slate-200">
                                          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M9 18h6"/><path d="M10 22h4"/><path d="M12 2v8"/><path d="M12 10l4 4"/><path d="M12 10l-4 4"/></svg>
                                        </div>
                                        <div className="min-w-0 flex-1">
                                          <NameCell name={displayRoleName(r)} />
                                        </div>
                                      </div>
                                    </td>
                                    <td className="px-4 py-4 text-sm text-gray-500">{r.orderNum ?? 0}</td>
                                    <td className="px-4 py-4">
                                      <BoolBadge value={!!r.isActive} yes={t('admin.catalog.enabled') || 'Enabled'} no={t('admin.catalog.disabled') || 'Disabled'} />
                                    </td>
                                    <td className="px-4 py-4 text-right sticky right-0 bg-white" onClick={(e) => e.stopPropagation()}>
                                      <div className="flex items-center justify-end gap-1">
                                        <IconButton title={t('common.edit') || 'Edit'} onClick={() => openEditRole(r)} kind="edit" />
                                        <IconButton title={t('common.delete') || 'Delete'} onClick={() => removeItem('role', r.id)} kind="delete" />
                                      </div>
                                    </td>
                                  </motion.tr>
                                ))
                              );
                            }
                            return rows;
                          })()
                        : []),
                    ]
                  : null;
                return (
                  <React.Fragment key={c.id}>
                      <tr
                        className="hover:bg-blue-50/30 transition-colors cursor-pointer"
                        onClick={() => setExpandedCategories((prev) => ({ ...prev, [expandedKey]: !prev[expandedKey] }))}
                        onKeyDown={(e) => {
                          if (e.key === 'Enter' || e.key === ' ') {
                            e.preventDefault();
                            setExpandedCategories((prev) => ({ ...prev, [expandedKey]: !prev[expandedKey] }));
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
                                setExpandedCategories((prev) => ({ ...prev, [expandedKey]: !prev[expandedKey] }));
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
                              <NameCell name={displayCategoryName(c)} />
                            </div>
                          </div>
                        </td>
                        <td className="px-4 py-4 text-sm text-slate-400 whitespace-nowrap">
                          {children.length ? `${children.length} ${t('admin.catalog.jobs.subcategories') || '小类'}` : '-'}
                          <span className="mx-2 text-slate-200">|</span>
                          {(totalRoleCountUnderRoot(c.id) || 0) ? `${totalRoleCountUnderRoot(c.id)} ${t('admin.catalog.jobs.roles') || '岗位'}` : '-'}
                        </td>
                        <td className="px-4 py-4 text-sm text-gray-500">{c.orderNum ?? 0}</td>
                        <td className="px-4 py-4">
                          <BoolBadge value={!!c.isActive} yes={t('admin.catalog.enabled') || 'Enabled'} no={t('admin.catalog.disabled') || 'Disabled'} />
                        </td>
                        <td className="px-4 py-4 text-right sticky right-0 bg-white" onClick={(e) => e.stopPropagation()}>
                          <div className="flex items-center justify-end gap-1">
                            <button
                              type="button"
                              onClick={() => openCreateSubCategory(c.id)}
                              className="p-2 text-slate-400 hover:text-indigo-600 hover:bg-indigo-50 rounded-lg transition-all"
                              title={t('admin.catalog.actions.addSubcategory') || 'Add subcategory'}
                            >
                              <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M5 12h14"/><path d="M12 5v14"/></svg>
                            </button>
                            <IconButton title={t('common.edit') || 'Edit'} onClick={() => openEditCategory(c)} kind="edit" />
                            <IconButton title={t('common.delete') || 'Delete'} onClick={() => removeItem('category', c.id)} kind="delete" />
                          </div>
                        </td>
                      </tr>
                      <AnimatePresence initial={false}>
                        {expandedRows}
                      </AnimatePresence>
                    </React.Fragment>
                  );
              })}
            </tbody>
          </table>
        </div>
      </TableCard>
    );
  };

  return (
    <div className={embedded ? 'flex flex-col' : 'flex-1 flex flex-col bg-white rounded-3xl m-2 overflow-hidden shadow-sm border border-gray-100'}>
      <div className={embedded ? 'pb-4 border-b border-gray-100' : 'px-10 pt-10 pb-6'}>
        {embedded ? (
          <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
            <div className="relative w-full md:w-[360px]">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
              <Input value={q} onChange={(e) => setQ(e.target.value)} placeholder={t('admin.keyword')} className="pl-10" />
            </div>
            <div className="flex items-center gap-2 md:justify-end">
              <Button onClick={openCreate}>{t('admin.actions.create') || 'Create'}</Button>
            </div>
          </div>
        ) : (
          <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4 border-b border-gray-100 pb-4">
            <div className="relative w-full md:w-[360px]">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
              <Input value={q} onChange={(e) => setQ(e.target.value)} placeholder={t('admin.keyword')} className="pl-10" />
            </div>
            <div className="flex items-center gap-2 md:justify-end">
              <Button onClick={openCreate}>{t('admin.actions.create') || 'Create'}</Button>
            </div>
          </div>
        )}
      </div>
      <div className={embedded ? 'flex-1 overflow-hidden pb-10 flex flex-col' : 'flex-1 overflow-hidden px-10 pb-10 flex flex-col'}>
        <div className="flex-1 overflow-hidden">
          {renderTable()}
        </div>
      </div>

      <Modal
        isOpen={showForm}
        onClose={() => setShowForm(false)}
        title={editingId ? (t('common.edit') || 'Edit') : (t('admin.actions.create') || 'Create')}
        size="lg"
        compact
      >
        <div className="flex flex-col max-h-[70vh]">
          <div className="flex-1 overflow-y-auto no-scrollbar pr-1 space-y-5">
            {renderFormFields()}
          </div>
          <div className="mt-4 pt-3 border-t border-slate-100 bg-slate-50/50 flex items-center justify-end gap-3">
            <Button variant="outline" onClick={() => setShowForm(false)} disabled={saving}>{t('common.cancel') || 'Cancel'}</Button>
            <Button onClick={save} disabled={saving}>{t('common.save') || 'Save'}</Button>
          </div>
        </div>
      </Modal>
    </div>
  );
};
