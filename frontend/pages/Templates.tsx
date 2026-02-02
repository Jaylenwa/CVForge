import React, { useMemo, useRef, useState, useEffect } from 'react';
import { Button } from '../components/ui/Button';
import { AppRoute } from '../types';
import { useLanguage } from '../contexts/LanguageContext';
import { fetchContentPresetData, fetchTemplateCatalog, listTemplateLibraryItems } from '../services/catalogService';
import { JobSidebar } from '../components/templateLibrary/JobSidebar';
import { ResumeTemplateCard } from '../components/templateLibrary/ResumeTemplateCard';
  
export const Templates: React.FC = () => {
  React.useEffect(() => {
    document.body.classList.add('no-scrollbar');
    document.documentElement.classList.add('no-scrollbar');
    return () => {
      document.body.classList.remove('no-scrollbar');
      document.documentElement.classList.remove('no-scrollbar');
    };
  }, []);
  const { t, language } = useLanguage();
  
  const [sortMode, setSortMode] = useState<'hot' | 'new'>('hot');
  const [selectedJobCategory, setSelectedJobCategory] = useState<string>('');
  const [selectedJobRole, setSelectedJobRole] = useState<string>('');
  
  const [jobCategories, setJobCategories] = useState<Array<{ id: string; name: string; parentId?: string; orderNum?: number }>>([]);
  const [jobRoles, setJobRoles] = useState<Array<{ id: string; categoryId: string; name: string; orderNum?: number }>>([]);
  const [templates, setTemplates] = useState<Array<{ templateId: string; name: string; presetId?: string; roleId?: string; usageCount?: number; globalUsageCount?: number; isPremium?: boolean }>>([]);
  const [presetDataMap, setPresetDataMap] = useState<Record<string, any>>({});
  const inFlightPresetIdsRef = useRef<Set<string>>(new Set());

  useEffect(() => {
    (async () => {
      try {
        const { jobCategories, jobRoles } = await fetchTemplateCatalog();
        setJobCategories(jobCategories.map((c) => ({ id: String(c.id), name: c.name, parentId: c.parentId == null ? '' : String(c.parentId), orderNum: c.orderNum })));
        setJobRoles(jobRoles.map((r) => ({ id: String(r.id), categoryId: String(r.categoryId), name: r.name, orderNum: r.orderNum })));
      } catch {
        setJobCategories([]);
        setJobRoles([]);
        setPresetDataMap({});
      }
    })();
  }, []);

  useEffect(() => {
    (async () => {
      try {
        const roleId = selectedJobRole ? Number(selectedJobRole) : undefined;
        const apiSort = sortMode === 'hot' ? 'hot' : 'new';
        const items = await listTemplateLibraryItems(roleId ? { roleId, language, sort: apiSort } : { language, sort: apiSort });
        setTemplates(items.map((t) => ({ templateId: String(t.templateExternalId), name: t.name, presetId: t.presetId ? String(t.presetId) : undefined, roleId: t.roleId ? String(t.roleId) : undefined, usageCount: t.usageCount, globalUsageCount: t.globalUsageCount, isPremium: t.isPremium })));
      } catch {
        setTemplates([]);
      }
    })();
  }, [selectedJobRole, sortMode, language]);

  const childCategoryIdsForSelectedRoot = useMemo(() => {
    const rootId = selectedJobCategory;
    if (!rootId) return new Set<string>();
    return new Set(jobCategories.filter(c => c.parentId === rootId).map(c => c.id));
  }, [jobCategories, selectedJobCategory]);

  const roleMap = useMemo(() => {
    const m: Record<string, any> = {};
    for (const r of jobRoles) m[r.id] = r;
    return m;
  }, [jobRoles]);

  const sortedTemplates = useMemo(() => templates.slice(), [templates]);

  const handleUseTemplate = (templateId: string, presetId?: string, roleId?: string) => {
    const qs = new URLSearchParams();
    qs.set('template', templateId);
    if (presetId) qs.set('presetId', presetId);
    if (roleId) qs.set('roleId', roleId);
    qs.set('returnTo', AppRoute.Templates);
    window.open(`${window.location.origin}${window.location.pathname}#${AppRoute.Editor}?${qs.toString()}`, '_blank');
  };
  const handlePreviewTemplate = (templateId: string, presetId?: string) => {
    const qs = new URLSearchParams();
    qs.set('template', templateId);
    if (presetId) qs.set('presetId', presetId);
    window.open(`${window.location.origin}${window.location.pathname}#${AppRoute.Print}?${qs.toString()}`, '_blank');
  };

  useEffect(() => {
    const uniquePresetIds: string[] = [];
    const pid = sortedTemplates.find((t) => !!t.presetId)?.presetId;
    if (pid) uniquePresetIds.push(pid);

    const controller = new AbortController();
    const signal = controller.signal;

    const fetchOne = async (presetId: string) => {
      if (signal.aborted) return;
      if (presetDataMap[presetId]) return;
      if (inFlightPresetIdsRef.current.has(presetId)) return;
      inFlightPresetIdsRef.current.add(presetId);
      try {
        const parsed = await fetchContentPresetData(Number(presetId), signal);
        if (!parsed || typeof parsed !== 'object') return;
        setPresetDataMap((prev) => (prev[presetId] ? prev : { ...prev, [presetId]: parsed }));
      } catch {
      } finally {
        inFlightPresetIdsRef.current.delete(presetId);
      }
    };

    const run = async () => {
      const queue = uniquePresetIds.slice();
      const workers = Array.from({ length: Math.min(6, queue.length) }, async () => {
        while (queue.length && !signal.aborted) {
          const next = queue.shift();
          if (!next) return;
          await fetchOne(next);
        }
      });
      await Promise.all(workers);
    };

    run();

    return () => controller.abort();
  }, [sortedTemplates, presetDataMap]);

  const sidebarCategories = useMemo(() => {
    return jobCategories;
  }, [jobCategories]);

  const sidebarRoles = useMemo(() => {
    return jobRoles.slice().sort((a, b) => (a.orderNum ?? 0) - (b.orderNum ?? 0));
  }, [jobRoles]);

  const selectedCategoryIdForSidebar = useMemo(() => {
    if (!selectedJobCategory) return '';
    const cat = jobCategories.find(c => c.id === selectedJobCategory);
    if (cat?.parentId) return cat.parentId;
    return selectedJobCategory;
  }, [jobCategories, selectedJobCategory]);

  const clearAll = () => {
    setSelectedJobCategory('');
    setSelectedJobRole('');
    setSortMode('hot');
  };

  const rolesForSelectedCategory = useMemo(() => {
    if (!selectedJobCategory) return [];
    return jobRoles
      .filter(r => r.categoryId === selectedJobCategory || childCategoryIdsForSelectedRoot.has(r.categoryId))
      .slice()
      .sort((a, b) => (a.orderNum ?? 0) - (b.orderNum ?? 0));
  }, [jobRoles, selectedJobCategory, childCategoryIdsForSelectedRoot]);

  return (
    <div className="min-h-screen bg-[#f8fafc] text-slate-900 font-sans selection:bg-blue-100 selection:text-blue-700">
      <main className="max-w-[1440px] mx-auto px-6 py-8 flex gap-8">
        <JobSidebar
          title={t('templates.sidebar.navLabel')}
          categories={sidebarCategories}
          roles={sidebarRoles}
          selectedCategoryId={selectedCategoryIdForSidebar}
          onSelectCategory={(id) => {
            setSelectedJobCategory(id);
            setSelectedJobRole('');
          }}
          onSelectRole={(roleId) => {
            const role = roleMap[roleId];
            if (role?.categoryId) setSelectedJobCategory(role.categoryId);
            setSelectedJobRole(roleId);
          }}
        />

        <section className="flex-1 min-w-0">
          <div className="mb-10 space-y-6">
            <div className="flex flex-col md:flex-row md:items-end justify-between gap-6">
              <div className="max-w-xl flex-1">
                <h2 className="text-3xl font-black text-slate-900 mb-2">{t('templates.hero.title')}</h2>
                <p className="text-slate-500 font-medium">{t('templates.hero.desc')}</p>
              </div>
              <div className="flex items-center gap-3">
                <div className="flex bg-white p-1 rounded-xl border border-slate-200 shadow-sm">
                  <button
                    onClick={() => setSortMode('new')}
                    className={`px-4 py-1.5 rounded-lg text-xs font-bold ${sortMode === 'new' ? 'bg-slate-100 text-slate-900' : 'text-slate-400 hover:text-slate-600'}`}
                  >
                    {t('templates.sort.latest')}
                  </button>
                  <button
                    onClick={() => setSortMode('hot')}
                    className={`px-4 py-1.5 rounded-lg text-xs font-bold ${sortMode === 'hot' ? 'bg-slate-100 text-slate-900' : 'text-slate-400 hover:text-slate-600'}`}
                  >
                    {t('templates.sort.hot')}
                  </button>
                </div>
              </div>
            </div>

            <div className="md:hidden bg-white border border-slate-200 rounded-2xl p-4">
              <div className="grid grid-cols-1 gap-3">
                <div className="relative">
                  <select
                    value={selectedJobCategory}
                    onChange={(e) => {
                      setSelectedJobCategory(e.target.value);
                      setSelectedJobRole('');
                    }}
                    className="w-full bg-white border-2 border-slate-200 rounded-xl py-3 px-4 text-sm font-medium outline-none focus:border-blue-500 transition-all"
                  >
                    <option value="">{t('templates.category.all')}</option>
                    {jobCategories
                      .slice()
                      .filter(c => !c.parentId)
                      .sort((a, b) => (a.orderNum ?? 0) - (b.orderNum ?? 0))
                      .map(c => (
                        <option key={c.id} value={c.id}>
                          {c.name}
                        </option>
                      ))}
                  </select>
                </div>
                {selectedJobCategory ? (
                  <div className="relative">
                    <select
                      value={selectedJobRole}
                      onChange={(e) => setSelectedJobRole(e.target.value)}
                      className="w-full bg-white border-2 border-slate-200 rounded-xl py-3 px-4 text-sm font-medium outline-none focus:border-blue-500 transition-all"
                    >
                      <option value="">{t('templates.category.all')}</option>
                      {rolesForSelectedCategory.map(r => (
                        <option key={r.id} value={r.id}>
                          {r.name}
                        </option>
                      ))}
                    </select>
                  </div>
                ) : null}
              </div>
            </div>

          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-8">
            {sortedTemplates.map((it: any) => {
              const roleName = selectedJobRole ? (roleMap[selectedJobRole]?.name || '') : '';
              const presetId = it.presetId || '';
              return (
                <ResumeTemplateCard
                  key={it.templateId}
                  title={it.name}
                  templateId={it.templateId}
                  usageCount={it.usageCount ?? 0}
                  isPremium={it.isPremium ?? false}
                  tag={roleName || undefined}
                  presetData={presetId ? (presetDataMap[presetId] || null) : null}
                  onUse={() => handleUseTemplate(it.templateId, it.presetId, it.roleId || selectedJobRole || undefined)}
                  onPreview={() => handlePreviewTemplate(it.templateId, it.presetId)}
                />
              );
            })}
          </div>

          {sortedTemplates.length === 0 ? (
            <div className="mt-10 text-center py-20 bg-white rounded-2xl border border-slate-200">
              <p className="text-slate-500 text-lg font-medium">{t('templates.empty')}</p>
              <Button variant="ghost" onClick={clearAll} className="mt-4">
                {t('templates.actions.clearFilters')}
              </Button>
            </div>
          ) : null}
        </section>
      </main>
    </div>
  );
};
