import React, { useMemo, useState, useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';
import { Button } from '../components/ui/Button';
import { API_BASE } from '../config';
import { AppRoute } from '../types';
import { useLanguage } from '../contexts/LanguageContext';
import { MOCK_TEMPLATES } from '../services/mockData';
import { CONTENT_PRESETS_SEED, JOB_CATEGORIES_SEED, JOB_ROLES_SEED, TEMPLATE_VARIANTS_SEED } from '../services/catalogSeeds';
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
  const [searchParams] = useSearchParams();
  const { t, language } = useLanguage();
  
  const [sortMode, setSortMode] = useState<'hot' | 'new'>('hot');
  const [selectedJobCategory, setSelectedJobCategory] = useState<string>('');
  const [selectedJobRole, setSelectedJobRole] = useState<string>('');
  

  useEffect(() => {
    const tag = searchParams.get('tag');
    if (tag) {
        // Simple mapping for demo purposes
        if (['IT', 'Finance', 'Creative', 'General'].includes(tag)) {
            // legacy param retained for backward compatibility; no-op for now
        }
    }
  }, [searchParams]);

  const [templates, setTemplates] = useState(Array<{id?:string; ExternalID?:string; name?:string; Name?:string; tags?:string[]; Tags?:string; usageCount?:number; UsageCount?:number; Popularity?:number; isPremium?:boolean; IsPremium?:boolean; category?:string; Category?:string}>());
  const [jobCategories, setJobCategories] = useState<Array<{ id: string; name: string; parentId?: string; orderNum?: number }>>([]);
  const [jobRoles, setJobRoles] = useState<Array<{ id: string; categoryId: string; name: string; tags?: string[]; orderNum?: number }>>([]);
  const [variants, setVariants] = useState<Array<{ id: string; name: string; layoutTemplateId: string; presetId: string; roleId: string; tags?: string[]; usageCount?: number; isPremium?: boolean }>>([]);
  const [presetDataMap, setPresetDataMap] = useState<Record<string, any>>({});

  useEffect(() => {
    (async () => {
      try {
        const res = await fetch(`${API_BASE}/templates`);
        if (res.ok) {
          const data = await res.json();
          const items = (data.items || []).map((t: any) => ({
            id: t.ExternalID || t.id,
            name: t.Name || t.name,
            tags: typeof t.Tags === 'string' ? (t.Tags as string).split(',') : (t.tags || []),
            usageCount: t.UsageCount ?? t.usageCount ?? t.Popularity ?? t.popularity,
            isPremium: t.IsPremium ?? t.isPremium,
            category: t.Category || t.category,
          }));
          setTemplates(items);
        } else {
          throw new Error('Network response was not ok');
        }
      } catch (error) {
        console.warn('Failed to fetch templates from API, falling back to mock data:', error);
        setTemplates(MOCK_TEMPLATES);
      }
    })();
  }, []);

  useEffect(() => {
    (async () => {
      try {
        const [catsRes, rolesRes, varsRes] = await Promise.all([
          fetch(`${API_BASE}/job-categories`),
          fetch(`${API_BASE}/job-roles`),
          fetch(`${API_BASE}/template-variants`)
        ]);
        if (!catsRes.ok || !rolesRes.ok || !varsRes.ok) {
          throw new Error('catalog api not ok');
        }
        const catsJson = await catsRes.json();
        const rolesJson = await rolesRes.json();
        const varsJson = await varsRes.json();
        const cats = (catsJson.items || []).map((c: any) => ({
          id: c.ExternalID || c.externalId || c.id,
          name: c.Name || c.name,
          parentId: c.ParentExternalID || c.parentId || '',
          orderNum: c.OrderNum ?? c.orderNum ?? 0
        }));
        const roles = (rolesJson.items || []).map((r: any) => ({
          id: r.ExternalID || r.externalId || r.id,
          categoryId: r.CategoryExternalID || r.categoryId || '',
          name: r.Name || r.name,
          tags: typeof r.Tags === 'string' ? (r.Tags as string).split(',') : (r.tags || []),
          orderNum: r.OrderNum ?? r.orderNum ?? 0
        }));
        const vars = (varsJson.items || []).map((v: any) => ({
          id: v.ExternalID || v.externalId || v.id,
          name: v.Name || v.name,
          layoutTemplateId: v.LayoutTemplateExternalID || v.layoutTemplateId,
          presetId: v.PresetExternalID || v.presetId,
          roleId: v.RoleExternalID || v.roleId,
          tags: typeof v.Tags === 'string' ? (v.Tags as string).split(',') : (v.tags || []),
          usageCount: v.UsageCount ?? v.usageCount ?? 0,
          isPremium: v.IsPremium ?? v.isPremium ?? false
        }));
        setJobCategories(cats);
        setJobRoles(roles);
        setVariants(vars);
        const presetMap: Record<string, any> = {};
        for (const p of CONTENT_PRESETS_SEED) {
          presetMap[p.id] = p.data;
        }
        setPresetDataMap(presetMap);
      } catch {
        setJobCategories(JOB_CATEGORIES_SEED);
        setJobRoles(JOB_ROLES_SEED);
        setVariants(TEMPLATE_VARIANTS_SEED);
        const presetMap: Record<string, any> = {};
        for (const p of CONTENT_PRESETS_SEED) {
          presetMap[p.id] = p.data;
        }
        setPresetDataMap(presetMap);
      }
    })();
  }, []);

  const rootJobCategories = useMemo(() => {
    return jobCategories.filter(c => !c.parentId);
  }, [jobCategories]);

  const childCategoryIdsForSelectedRoot = useMemo(() => {
    const rootId = selectedJobCategory;
    if (!rootId) return new Set<string>();
    return new Set(jobCategories.filter(c => c.parentId === rootId).map(c => c.id));
  }, [jobCategories, selectedJobCategory]);

  useEffect(() => {
    if (!selectedJobCategory && rootJobCategories.length > 0) {
      setSelectedJobCategory(rootJobCategories.slice().sort((a, b) => (a.orderNum ?? 0) - (b.orderNum ?? 0))[0].id);
    }
  }, [rootJobCategories, selectedJobCategory]);

  const roleMap = useMemo(() => {
    const m: Record<string, any> = {};
    for (const r of jobRoles) m[r.id] = r;
    return m;
  }, [jobRoles]);

  const filteredVariants = useMemo(() => {
    return variants.filter(v => {
    const role = roleMap[v.roleId];
    const matchesJobCategory = !selectedJobCategory || (role && (role.categoryId === selectedJobCategory || childCategoryIdsForSelectedRoot.has(role.categoryId)));
    const matchesJobRole = !selectedJobRole || v.roleId === selectedJobRole;
    return matchesJobCategory && matchesJobRole;
    });
  }, [variants, roleMap, selectedJobCategory, selectedJobRole, childCategoryIdsForSelectedRoot]);

  const filteredTemplates = useMemo(() => {
    return templates.filter((tpl: any) => {
      return true;
    });
  }, [templates]);

  const sortedVariants = useMemo(() => {
    const list = filteredVariants.slice();
    if (sortMode === 'hot') {
      return list.sort((a, b) => (b.usageCount ?? 0) - (a.usageCount ?? 0));
    }
    return list.sort((a, b) => String(a.name || '').localeCompare(String(b.name || ''), language === 'zh' ? 'zh' : 'en'));
  }, [filteredVariants, sortMode, language]);

  const sortedTemplates = useMemo(() => {
    const list = filteredTemplates.slice();
    if (sortMode === 'hot') {
      return list.sort((a: any, b: any) => (b.usageCount ?? 0) - (a.usageCount ?? 0));
    }
    return list.sort((a: any, b: any) => String(a.name || '').localeCompare(String(b.name || ''), language === 'zh' ? 'zh' : 'en'));
  }, [filteredTemplates, sortMode, language]);

  const handleUseTemplate = (templateId: string, presetId?: string, variantId?: string) => {
    const qs = new URLSearchParams();
    qs.set('template', templateId);
    if (presetId) qs.set('preset', presetId);
    if (variantId) qs.set('variant', variantId);
    qs.set('returnTo', AppRoute.Templates);
    window.open(`${window.location.origin}${window.location.pathname}#${AppRoute.Editor}?${qs.toString()}`, '_blank');
  };
  const handlePreviewTemplate = (templateId: string) => {
    window.open(`${window.location.origin}${window.location.pathname}#${AppRoute.Print}?template=${templateId}`, '_blank');
  };

  const useVariantMode = !!selectedJobCategory || !!selectedJobRole;
  const sidebarCategories = useMemo(() => {
    return jobCategories;
  }, [jobCategories]);

  const sidebarRoles = useMemo(() => {
    return jobRoles.slice().sort((a, b) => (a.orderNum ?? 0) - (b.orderNum ?? 0));
  }, [jobRoles]);

  const selectedCategoryIdForSidebar = selectedJobCategory || (rootJobCategories[0]?.id ?? '');

  const clearAll = () => {
    if (rootJobCategories.length > 0) {
      setSelectedJobCategory(rootJobCategories.slice().sort((a, b) => (a.orderNum ?? 0) - (b.orderNum ?? 0))[0].id);
    }
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
            {(useVariantMode ? sortedVariants : sortedTemplates).map((item: any) => {
              if (useVariantMode) {
                const variant = item;
                const roleName = roleMap[variant.roleId]?.name || '';
                return (
                  <ResumeTemplateCard
                    key={variant.id}
                    title={variant.name}
                    templateId={variant.layoutTemplateId}
                    usageCount={variant.usageCount ?? 0}
                    isPremium={variant.isPremium ?? false}
                    tag={roleName}
                    presetData={presetDataMap[variant.presetId] || null}
                    onUse={() => handleUseTemplate(variant.layoutTemplateId, variant.presetId, variant.id)}
                    onPreview={() => handlePreviewTemplate(variant.layoutTemplateId)}
                  />
                );
              }
              const tpl = item;
              return (
                <ResumeTemplateCard
                  key={tpl.id}
                  title={tpl.name}
                  templateId={tpl.id}
                  usageCount={tpl.usageCount ?? 0}
                  isPremium={tpl.isPremium ?? false}
                  tag={tpl.category}
                  presetData={null}
                  onUse={() => handleUseTemplate(tpl.id)}
                  onPreview={() => handlePreviewTemplate(tpl.id)}
                />
              );
            })}
          </div>

          {(useVariantMode ? sortedVariants.length === 0 : sortedTemplates.length === 0) ? (
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
