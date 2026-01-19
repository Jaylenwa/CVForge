import { API_BASE } from '../config';

export type JobCategory = {
  id: string;
  name: string;
  parentId?: string;
  orderNum?: number;
};

export type JobRole = {
  id: string;
  categoryId: string;
  name: string;
  tags?: string[];
  orderNum?: number;
};

export type TemplateVariant = {
  id: string;
  name: string;
  layoutTemplateId: string;
  presetId: string;
  roleId: string;
  tags?: string[];
  usageCount?: number;
  isPremium?: boolean;
};

export type ContentPreset = {
  id: string;
  name: string;
  language?: string;
  roleId?: string;
  tags?: string[];
  dataJson?: string;
};

const fetchJson = async <T>(url: string, init?: RequestInit): Promise<T> => {
  const res = await fetch(url, init);
  if (!res.ok) throw new Error(`request failed: ${res.status}`);
  return res.json();
};

const normalizeTags = (tags: unknown): string[] | undefined => {
  if (Array.isArray(tags)) {
    const out = tags.map((x) => String(x || '').trim()).filter(Boolean);
    return out.length ? out : undefined;
  }
  if (typeof tags === 'string') {
    const out = tags.split(',').map((x) => x.trim()).filter(Boolean);
    return out.length ? out : undefined;
  }
  return undefined;
};

export const fetchTemplateCatalog = async () => {
  const [catsJson, rolesJson, varsJson] = await Promise.all([
    fetchJson<{ items: any[] }>(`${API_BASE}/taxonomy/categories`),
    fetchJson<{ items: any[] }>(`${API_BASE}/taxonomy/roles`),
    fetchJson<{ items: any[] }>(`${API_BASE}/library/variants`)
  ]);

  const jobCategories: JobCategory[] = (catsJson.items || []).map((c: any) => ({
    id: c.externalId,
    name: c.name,
    parentId: c.parentExternalId || '',
    orderNum: c.orderNum ?? 0
  }));

  const jobRoles: JobRole[] = (rolesJson.items || []).map((r: any) => ({
    id: r.externalId,
    categoryId: r.categoryExternalId || '',
    name: r.name,
    tags: normalizeTags(r.tags),
    orderNum: r.orderNum ?? 0
  }));

  const variants: TemplateVariant[] = (varsJson.items || []).map((v: any) => ({
    id: v.externalId,
    name: v.name,
    layoutTemplateId: v.layoutTemplateExternalId,
    presetId: v.presetExternalId,
    roleId: v.roleExternalId,
    tags: normalizeTags(v.tags),
    usageCount: v.usageCount ?? 0,
    isPremium: !!v.isPremium
  }));

  return { jobCategories, jobRoles, variants };
};

export const listJobRoles = async (params?: { categoryId?: string; q?: string }) => {
  const qs = new URLSearchParams();
  if (params?.categoryId) qs.set('categoryId', params.categoryId);
  if (params?.q) qs.set('q', params.q);
  const json = await fetchJson<{ items: any[] }>(`${API_BASE}/taxonomy/roles?${qs.toString()}`);
  const jobRoles: JobRole[] = (json.items || []).map((r: any) => ({
    id: r.externalId,
    categoryId: r.categoryExternalId || '',
    name: r.name,
    tags: normalizeTags(r.tags),
    orderNum: r.orderNum ?? 0
  }));
  return jobRoles;
};

export const listTemplateVariants = async (params?: { roleId?: string; categoryId?: string; q?: string }) => {
  const qs = new URLSearchParams();
  if (params?.roleId) qs.set('roleId', params.roleId);
  if (params?.categoryId) qs.set('categoryId', params.categoryId);
  if (params?.q) qs.set('q', params.q);
  const json = await fetchJson<{ items: any[] }>(`${API_BASE}/library/variants?${qs.toString()}`);
  const variants: TemplateVariant[] = (json.items || []).map((v: any) => ({
    id: v.externalId,
    name: v.name,
    layoutTemplateId: v.layoutTemplateExternalId,
    presetId: v.presetExternalId,
    roleId: v.roleExternalId,
    tags: normalizeTags(v.tags),
    usageCount: v.usageCount ?? 0,
    isPremium: !!v.isPremium
  }));
  return variants;
};

export const fetchContentPresetData = async (presetId: string, signal?: AbortSignal): Promise<any | null> => {
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

