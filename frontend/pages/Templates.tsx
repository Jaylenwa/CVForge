import React, { useState, useEffect, useRef, useLayoutEffect } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { Search, Star, Filter } from 'lucide-react';
import { Button } from '../components/ui/Button';
// 后端数据来源
import { API_BASE } from '../config';
import { AppRoute } from '../types';
import { useLanguage } from '../contexts/LanguageContext';
import { ResumeArtboard } from './editor/ResumePreview';
import { INITIAL_RESUME, MOCK_TEMPLATES } from '../services/mockData';
import { CONTENT_PRESETS_SEED, JOB_CATEGORIES_SEED, JOB_ROLES_SEED, TEMPLATE_VARIANTS_SEED } from '../services/catalogSeeds';
  
export const Templates: React.FC = () => {
  const navigate = useNavigate();
  React.useEffect(() => {
    document.body.classList.add('no-scrollbar');
    document.documentElement.classList.add('no-scrollbar');
    return () => {
      document.body.classList.remove('no-scrollbar');
      document.documentElement.classList.remove('no-scrollbar');
    };
  }, []);
  const [searchParams] = useSearchParams();
  const { t } = useLanguage();
  
  const [filter, setFilter] = useState('');
  const [selectedCategory, setSelectedCategory] = useState<string>('All');
  const [selectedJobCategory, setSelectedJobCategory] = useState<string>('');
  const [selectedJobRole, setSelectedJobRole] = useState<string>('');
  // preview moved to dedicated print page via router
  

  useEffect(() => {
    const tag = searchParams.get('tag');
    if (tag) {
        // Simple mapping for demo purposes
        if (['IT', 'Finance', 'Creative', 'General'].includes(tag)) {
            setSelectedCategory(tag);
        }
    }
  }, [searchParams]);

  const categories = ['All', 'IT', 'Finance', 'Creative', 'General'];

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

  const filteredTemplates = templates.filter((t: any) => {
    const matchesSearch = t.name.toLowerCase().includes(filter.toLowerCase());
    const matchesCategory = selectedCategory === 'All' || t.category === selectedCategory;
    return matchesSearch && matchesCategory;
  });

  const roleMap = React.useMemo(() => {
    const m: Record<string, any> = {};
    for (const r of jobRoles) m[r.id] = r;
    return m;
  }, [jobRoles]);

  const filteredRoles = jobRoles
    .filter(r => !selectedJobCategory || r.categoryId === selectedJobCategory)
    .sort((a, b) => (a.orderNum ?? 0) - (b.orderNum ?? 0));

  const filteredRolesBySearch = React.useMemo(() => {
    const q = filter.trim().toLowerCase();
    if (!q) return filteredRoles;
    return filteredRoles.filter(r => r.name.toLowerCase().includes(q) || (r.tags || []).some((x: string) => String(x).toLowerCase().includes(q)));
  }, [filteredRoles, filter]);

  const filteredVariants = variants.filter(v => {
    const role = roleMap[v.roleId];
    const matchesJobCategory = !selectedJobCategory || (role && role.categoryId === selectedJobCategory);
    const matchesJobRole = !selectedJobRole || v.roleId === selectedJobRole;
    const q = filter.trim().toLowerCase();
    const matchesSearch = !q || v.name.toLowerCase().includes(q) || (v.tags || []).some((x: string) => String(x).toLowerCase().includes(q));
    return matchesJobCategory && matchesJobRole && matchesSearch;
  });

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

  // preview is now handled by navigating to the print page

  const TemplateGridItem: React.FC<{ template: any; presetId?: string; variantId?: string }> = ({ template, presetId, variantId }) => {
    const containerRef = useRef<HTMLDivElement | null>(null);
    const rafRef = useRef<number | null>(null);
    const roRef = useRef<ResizeObserver | null>(null);
    const stableTimerRef = useRef<number | null>(null);
    const lastWidthRef = useRef<number>(0);
    const initializedRef = useRef(false);
    const [scale, setScale] = useState<number | null>(null);
    const [ready, setReady] = useState(false);
    useLayoutEffect(() => {
      const mmToPx = 96 / 25.4;
      const a4w = 210 * mmToPx;
      const scheduleUpdate = () => {
        if (rafRef.current) cancelAnimationFrame(rafRef.current);
        rafRef.current = requestAnimationFrame(() => {
          const el = containerRef.current;
          if (!el) return;
          lastWidthRef.current = el.clientWidth;
          if (stableTimerRef.current) {
            clearTimeout(stableTimerRef.current);
          }
          stableTimerRef.current = window.setTimeout(() => {
            const s = lastWidthRef.current / a4w;
            setScale(prev => (prev === null || Math.abs(prev - s) > 0.002) ? s : prev);
            setReady(true);
          }, 120);
        });
      };
      if (!initializedRef.current) {
        const el = containerRef.current;
        if (el) {
          const s = el.clientWidth / a4w;
          setScale(s);
          setReady(true);
          initializedRef.current = true;
        }
      } else {
        scheduleUpdate();
      }
      const onResize = () => scheduleUpdate();
      window.addEventListener('resize', onResize);
      if (containerRef.current) {
        roRef.current = new ResizeObserver(onResize);
        roRef.current.observe(containerRef.current);
      }
      return () => {
        window.removeEventListener('resize', onResize);
        if (rafRef.current) cancelAnimationFrame(rafRef.current);
        if (stableTimerRef.current) {
          clearTimeout(stableTimerRef.current);
        }
        if (roRef.current) {
          roRef.current.disconnect();
        }
      };
    }, []);
    const mmToPx = 96 / 25.4;
    const a4w = 210 * mmToPx;
    const a4h = 297 * mmToPx;
    const presetData = presetId ? presetDataMap[presetId] : null;
    const previewData = presetData ? { ...INITIAL_RESUME, ...presetData, templateId: template.id } : { ...INITIAL_RESUME, templateId: template.id };
    return (
      <div className="group relative bg-white border border-gray-200 rounded-xl overflow-hidden shadow-sm hover:shadow-lg">
        <div ref={containerRef} className="aspect-[210/297] w-full bg-gray-200 overflow-hidden relative">
          <div className="absolute inset-0 flex items-center justify-center">
            {ready && scale !== null ? (
              <div
                style={{ width: a4w * scale, height: a4h * scale }}
                className="relative select-none pointer-events-none shadow-sm bg-white"
              >
                <ResumeArtboard
                  data={previewData}
                  scale={scale}
                  disableShadow={true}
                  showPageHint={false}
                  style={{ margin: 0 }}
                />
              </div>
            ) : (
              <div className="w-full h-full bg-white" />
            )}
          </div>
          <div className="absolute inset-0 bg-black bg-opacity-0 group-hover:bg-opacity-40 flex items-center justify-center opacity-0 group-hover:opacity-100">
            <div className="flex flex-col items-center space-y-3">
              <Button className="w-40" onClick={() => handleUseTemplate(template.id, presetId, variantId)}>{t('templates.actions.useTemplate')}</Button>
              <Button className="w-40" variant="outline" onClick={() => handlePreviewTemplate(template.id)}>{t('common.preview')}</Button>
            </div>
          </div>
          {template.isPremium && (
            <div className="absolute top-2 right-2 bg-yellow-400 text-yellow-900 text-xs font-bold px-2 py-1 rounded flex items-center">
              <Star size={12} className="mr-1 fill-current" /> {t('templates.badge.premium')}
            </div>
          )}
        </div>
        <div className="p-4">
          <h3 className="text-lg font-medium text-gray-900">{template.name}</h3>
          <div className="mt-2 flex items-center justify-between text-sm text-gray-500">
            <span>{template.usageCount ?? 0} {t('templates.meta.usageCount')}</span>
          </div>
          <div className="mt-3 flex flex-wrap gap-1">
            {template.category ? (
              <span className="px-2 py-0.5 bg-blue-50 text-blue-700 text-xs rounded border border-blue-100">{template.category}</span>
            ) : null}
          </div>
        </div>
      </div>
    );
  };

  const useVariantMode = !!selectedJobCategory || !!selectedJobRole;
  const sortedJobCategories = React.useMemo(() => {
    return jobCategories.slice().sort((a, b) => (a.orderNum ?? 0) - (b.orderNum ?? 0));
  }, [jobCategories]);

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <div className="flex flex-col md:flex-row md:items-center md:justify-between mb-8">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">{t('templates.title')}</h1>
          <p className="mt-2 text-gray-500">{t('templates.desc')}</p>
        </div>
      </div>

      <div className="flex gap-6">
        <aside className="w-56 shrink-0">
          <div className="bg-white border border-gray-200 rounded-xl overflow-hidden">
            <div className="px-4 py-3 border-b border-gray-200 flex items-center justify-between">
              <div className="flex items-center space-x-2">
                <Filter size={16} className="text-gray-500" />
                <span className="font-medium text-gray-800 text-sm">{t('templates.filter.jobCategory')}</span>
              </div>
            </div>
            <div className="max-h-[520px] overflow-auto">
              <button
                className={`w-full text-left px-4 py-2.5 text-sm border-l-2 ${!selectedJobCategory ? 'bg-blue-50 border-blue-500 text-blue-800' : 'bg-white border-transparent text-gray-700 hover:bg-gray-50'}`}
                onClick={() => {
                  setSelectedJobCategory('');
                  setSelectedJobRole('');
                }}
              >
                {t('templates.category.all')}
              </button>
              {sortedJobCategories.map(c => (
                <button
                  key={c.id}
                  className={`w-full text-left px-4 py-2.5 text-sm border-l-2 ${selectedJobCategory === c.id ? 'bg-blue-50 border-blue-500 text-blue-800' : 'bg-white border-transparent text-gray-700 hover:bg-gray-50'}`}
                  onClick={() => {
                    setSelectedJobCategory(c.id);
                    setSelectedJobRole('');
                  }}
                >
                  {c.name}
                </button>
              ))}
            </div>
            <div className="p-4 border-t border-gray-200 bg-gray-50">
              <div className="flex items-center justify-between mb-2">
                <span className="text-xs font-medium text-gray-700">{t('templates.filter.industry')}</span>
                <button
                  onClick={() => { setSelectedCategory('All'); setSelectedJobCategory(''); setSelectedJobRole(''); setFilter(''); }}
                  className="text-xs text-blue-600 hover:text-blue-800"
                >
                  {t('templates.actions.clearAll')}
                </button>
              </div>
              <div className="relative">
                <select
                  value={selectedCategory}
                  onChange={(e) => setSelectedCategory(e.target.value)}
                  className="appearance-none block w-full pl-3 pr-8 py-2 text-sm border-gray-300 focus:outline-none focus:ring-blue-500 focus:border-blue-500 rounded-md border bg-white"
                >
                  {categories.map(c => {
                    const key = c === 'All' ? 'all' : c;
                    return <option key={c} value={c}>{t(`templates.category.${key}`)}</option>
                  })}
                </select>
                <span className="pointer-events-none absolute right-2 top-1/2 -translate-y-1/2 text-gray-500">
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M7 10l5 5 5-5" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                  </svg>
                </span>
              </div>
            </div>
          </div>
        </aside>

        <main className="min-w-0 flex-1">
          <div className="mb-4 flex flex-col md:flex-row md:items-center md:justify-between gap-3">
            <div className="relative rounded-md shadow-sm w-full md:w-80">
              <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <Search className="h-5 w-5 text-gray-400" />
              </div>
              <input
                type="text"
                className="focus:ring-blue-500 focus:border-blue-500 block w-full pl-10 sm:text-sm border-gray-300 rounded-md p-2 border"
                placeholder={t('templates.search')}
                value={filter}
                onChange={(e) => setFilter(e.target.value)}
              />
            </div>
          </div>

          {selectedJobCategory ? (
            <div className="mb-6 bg-white border border-gray-200 rounded-xl p-4">
              <div className="flex items-center justify-between mb-3">
                <span className="text-sm font-medium text-gray-800">{t('templates.filter.jobRole')}</span>
                {selectedJobRole ? (
                  <button className="text-sm text-blue-600 hover:text-blue-800" onClick={() => setSelectedJobRole('')}>
                    {t('templates.actions.clearFilters')}
                  </button>
                ) : null}
              </div>
              <div className="flex flex-wrap gap-2">
                {filteredRolesBySearch.map(r => {
                  const active = selectedJobRole === r.id;
                  return (
                    <button
                      key={r.id}
                      onClick={() => setSelectedJobRole(prev => (prev === r.id ? '' : r.id))}
                      className={`px-3 py-1.5 rounded-full text-sm border ${active ? 'bg-blue-600 border-blue-600 text-white' : 'bg-white border-gray-200 text-gray-700 hover:border-blue-300 hover:text-blue-700'}`}
                    >
                      {r.name}
                    </button>
                  );
                })}
                {filteredRolesBySearch.length === 0 ? (
                  <div className="text-sm text-gray-500">{t('templates.empty')}</div>
                ) : null}
              </div>
            </div>
          ) : null}

          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-8">
            {(useVariantMode ? filteredVariants : filteredTemplates).map((item: any) => {
              if (useVariantMode) {
                const variant = item;
                const tpl = {
                  id: variant.layoutTemplateId,
                  name: variant.name,
                  usageCount: variant.usageCount ?? 0,
                  isPremium: variant.isPremium ?? false,
                  category: ''
                };
                return <TemplateGridItem key={variant.id} template={tpl} presetId={variant.presetId} variantId={variant.id} />;
              }
              return <TemplateGridItem key={item.id} template={item} />;
            })}
          </div>
        </main>
      </div>
      
      {(useVariantMode ? filteredVariants.length === 0 : filteredTemplates.length === 0) && (
          <div className="text-center py-20 bg-gray-50 rounded-lg border-2 border-dashed border-gray-200">
              <p className="text-gray-500 text-lg">{t('templates.empty')}</p>
              <Button variant="ghost" onClick={() => {setFilter(''); setSelectedCategory('All'); setSelectedJobCategory(''); setSelectedJobRole('');}} className="mt-4">{t('templates.actions.clearFilters')}</Button>
          </div>
      )}
      
    </div>
  );
};
