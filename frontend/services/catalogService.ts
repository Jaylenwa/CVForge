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
  tags?: string[];
  orderNum?: number;
};

export type TemplateVariant = {
  id: number;
  name: string;
  layoutTemplateId: string;
  presetId: number;
  roleId: number;
  tags?: string[];
  usageCount?: number;
  isPremium?: boolean;
};

export type ContentPreset = {
  id: number;
  name: string;
  language?: string;
  roleId?: number;
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
    id: c.id,
    name: c.name,
    parentId: c.parentId ?? null,
    orderNum: c.orderNum ?? 0
  }));

  const jobRoles: JobRole[] = (rolesJson.items || []).map((r: any) => ({
    id: r.id,
    categoryId: r.categoryId,
    name: r.name,
    tags: normalizeTags(r.tags),
    orderNum: r.orderNum ?? 0
  }));

  const variants: TemplateVariant[] = (varsJson.items || []).map((v: any) => ({
    id: v.id,
    name: v.name,
    layoutTemplateId: v.layoutTemplateExternalId,
    presetId: v.presetId,
    roleId: v.roleId,
    tags: normalizeTags(v.tags),
    usageCount: v.usageCount ?? 0,
    isPremium: !!v.isPremium
  }));

  return { jobCategories, jobRoles, variants };
};

export const listJobRoles = async (params?: { categoryId?: number; q?: string }) => {
  const qs = new URLSearchParams();
  if (params?.categoryId) qs.set('categoryId', String(params.categoryId));
  if (params?.q) qs.set('q', params.q);
  const json = await fetchJson<{ items: any[] }>(`${API_BASE}/taxonomy/roles?${qs.toString()}`);
  const jobRoles: JobRole[] = (json.items || []).map((r: any) => ({
    id: r.id,
    categoryId: r.categoryId,
    name: r.name,
    tags: normalizeTags(r.tags),
    orderNum: r.orderNum ?? 0
  }));
  return jobRoles;
};

export const listTemplateVariants = async (params?: { roleId?: number; categoryId?: number; q?: string }) => {
  const qs = new URLSearchParams();
  if (params?.roleId) qs.set('roleId', String(params.roleId));
  if (params?.categoryId) qs.set('categoryId', String(params.categoryId));
  if (params?.q) qs.set('q', params.q);
  const json = await fetchJson<{ items: any[] }>(`${API_BASE}/library/variants?${qs.toString()}`);
  const variants: TemplateVariant[] = (json.items || []).map((v: any) => ({
    id: v.id,
    name: v.name,
    layoutTemplateId: v.layoutTemplateExternalId,
    presetId: v.presetId,
    roleId: v.roleId,
    tags: normalizeTags(v.tags),
    usageCount: v.usageCount ?? 0,
    isPremium: !!v.isPremium
  }));
  return variants;
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
