import { API_BASE } from '../config';

export type JobCategory = {
  id: number;
  name: string;
  parentId?: number | null;
  orderNum?: number;
};

export type JobRole = {
  id: number;
  categoryId: number;
  name: string;
  orderNum?: number;
};

export type TemplateLibraryItem = {
  templateExternalId: string;
  name: string;
  usageCount?: number;
  globalUsageCount?: number;
  presetId?: number;
  roleId?: number;
  isPremium?: boolean;
};

export type ContentPreset = {
  id: number;
  name: string;
  language?: string;
  roleId?: number;
  dataJson?: string;
};

const fetchJson = async <T>(url: string, init?: RequestInit): Promise<T> => {
  const res = await fetch(url, init);
  if (!res.ok) throw new Error(`request failed: ${res.status}`);
  return res.json();
};

export const fetchTemplateCatalog = async (params?: { language?: string; sort?: 'hot' | 'new' | 'name' }) => {
  const taxonomyQs = new URLSearchParams();
  if (params?.language) taxonomyQs.set('language', String(params.language));
  const taxonomySuffix = taxonomyQs.toString() ? `?${taxonomyQs.toString()}` : '';
  const [catsJson, rolesJson, templatesJson] = await Promise.all([
    fetchJson<{ items: any[] }>(`${API_BASE}/taxonomy/categories${taxonomySuffix}`),
    fetchJson<{ items: any[] }>(`${API_BASE}/taxonomy/roles${taxonomySuffix}`),
    (() => {
      const qs = new URLSearchParams();
      if (params?.language) qs.set('language', String(params.language));
      if (params?.sort) qs.set('sort', String(params.sort));
      return fetchJson<{ items: any[] }>(`${API_BASE}/library/templates?${qs.toString()}`);
    })()
  ]);

  const jobCategories: JobCategory[] = (catsJson.items || []).map((c: any) => ({
    id: c.id,
    name: c.name,
    parentId: c.parentId ?? null,
    orderNum: c.orderNum ?? 0
  }));

  const jobRoles: JobRole[] = (rolesJson.items || []).map((r: any) => ({
    id: r.id,
    categoryId: r.categoryId,
    name: r.name,
    orderNum: r.orderNum ?? 0
  }));

  const templates: TemplateLibraryItem[] = (templatesJson.items || []).map((t: any) => ({
    templateExternalId: t.templateExternalId,
    name: t.name,
    usageCount: t.usageCount ?? 0,
    globalUsageCount: t.globalUsageCount ?? 0,
    presetId: t.presetId || undefined,
    roleId: t.roleId || undefined,
    isPremium: !!t.isPremium
  }));

  return { jobCategories, jobRoles, templates };
};

export const listJobRoles = async (params?: { categoryId?: number; q?: string; language?: string }) => {
  const qs = new URLSearchParams();
  if (params?.categoryId) qs.set('categoryId', String(params.categoryId));
  if (params?.q) qs.set('q', params.q);
  if (params?.language) qs.set('language', params.language);
  const json = await fetchJson<{ items: any[] }>(`${API_BASE}/taxonomy/roles?${qs.toString()}`);
  const jobRoles: JobRole[] = (json.items || []).map((r: any) => ({
    id: r.id,
    categoryId: r.categoryId,
    name: r.name,
    orderNum: r.orderNum ?? 0
  }));
  return jobRoles;
};

export const listTemplateLibraryItems = async (params?: { roleId?: number; language?: string; sort?: 'hot' | 'new' | 'name' }) => {
  const qs = new URLSearchParams();
  if (params?.roleId) qs.set('roleId', String(params.roleId));
  if (params?.language) qs.set('language', String(params.language));
  if (params?.sort) qs.set('sort', String(params.sort));
  const json = await fetchJson<{ items: any[] }>(`${API_BASE}/library/templates?${qs.toString()}`);
  const templates: TemplateLibraryItem[] = (json.items || []).map((t: any) => ({
    templateExternalId: t.templateExternalId,
    name: t.name,
    usageCount: t.usageCount ?? 0,
    globalUsageCount: t.globalUsageCount ?? 0,
    presetId: t.presetId || undefined,
    roleId: t.roleId || undefined,
    isPremium: !!t.isPremium
  }));
  return templates;
};

export const fetchContentPresetData = async (presetId: number, signal?: AbortSignal): Promise<any | null> => {
  if (!presetId) return null;
  const res = await fetch(`${API_BASE}/presets/${encodeURIComponent(presetId)}`, { signal });
  if (!res.ok) return null;
  const p: any = await res.json();
  const dataJson = p?.dataJson;
  if (typeof dataJson !== 'string' || !dataJson.trim()) return null;
  try {
    const parsed = JSON.parse(dataJson);
    return parsed && typeof parsed === 'object' ? parsed : null;
  } catch {
    return null;
  }
};
