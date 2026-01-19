import { API_BASE } from '../config';

const authHeader = () => ({ Authorization: `Bearer ${localStorage.getItem('token') || ''}` });

export interface PageResp<T> {
  items: T[];
  page: number;
  pageSize: number;
  total: number;
  hasNext?: boolean;
  totalPages?: number;
}

// Stats
export interface AdminStats {
  totals: { users: number; resumes: number; templates: number; visitorsToday?: number };
  trend: { dates: string[]; users: number[]; resumes: number[]; templates: number[]; visitors?: number[] };
  generatedAt: number;
}

export const getAdminStats = async (days = 14): Promise<AdminStats> => {
  const res = await fetch(`${API_BASE}/admin/stats?days=${days}`, { headers: authHeader() });
  if (!res.ok) throw new Error('Failed to fetch stats');
  return res.json();
};

export const importCatalogSeed = async () => {
  const res = await fetch(`${API_BASE}/admin/catalog/import-seed`, { method: 'POST', headers: authHeader() });
  if (!res.ok) throw new Error('failed');
  return res.json() as Promise<{ success: boolean; counts?: any }>;
};

export interface AdminJobCategory {
  id?: number;
  ExternalID: string;
  Name: string;
  ParentExternalID?: string;
  OrderNum?: number;
  IsActive: boolean;
}

export interface AdminJobRole {
  id?: number;
  ExternalID: string;
  CategoryExternalID: string;
  Name: string;
  Tags?: string;
  OrderNum?: number;
  IsActive: boolean;
}

export interface AdminContentPreset {
  id?: number;
  ExternalID: string;
  Name: string;
  Language?: string;
  RoleExternalID?: string;
  Tags?: string;
  DataJSON?: string;
  IsActive: boolean;
}

export interface AdminTemplateVariant {
  id?: number;
  ExternalID: string;
  Name: string;
  LayoutTemplateExternalID: string;
  PresetExternalID: string;
  RoleExternalID: string;
  Tags?: string;
  UsageCount?: number;
  IsPremium?: boolean;
  IsActive: boolean;
}

export const adminListJobCategories = async (params: Record<string, string>) => {
  const q = new URLSearchParams(params).toString();
  const res = await fetch(`${API_BASE}/admin/catalog/job-categories?${q}`, { headers: authHeader() });
  if (!res.ok) throw new Error('failed');
  return res.json() as Promise<PageResp<AdminJobCategory>>;
};

export const adminCreateJobCategory = async (body: any) => {
  const res = await fetch(`${API_BASE}/admin/catalog/job-categories`, { method: 'POST', headers: { 'Content-Type': 'application/json', ...authHeader() }, body: JSON.stringify(body) });
  if (!res.ok) throw new Error('failed');
};

export const adminPatchJobCategory = async (externalId: string, body: any) => {
  const res = await fetch(`${API_BASE}/admin/catalog/job-categories/${externalId}`, { method: 'PATCH', headers: { 'Content-Type': 'application/json', ...authHeader() }, body: JSON.stringify(body) });
  if (!res.ok) throw new Error('failed');
};

export const adminDeleteJobCategory = async (externalId: string) => {
  const res = await fetch(`${API_BASE}/admin/catalog/job-categories/${externalId}`, { method: 'DELETE', headers: authHeader() });
  if (!res.ok) throw new Error('failed');
};

export const adminListJobRoles = async (params: Record<string, string>) => {
  const q = new URLSearchParams(params).toString();
  const res = await fetch(`${API_BASE}/admin/catalog/job-roles?${q}`, { headers: authHeader() });
  if (!res.ok) throw new Error('failed');
  return res.json() as Promise<PageResp<AdminJobRole>>;
};

export const adminCreateJobRole = async (body: any) => {
  const res = await fetch(`${API_BASE}/admin/catalog/job-roles`, { method: 'POST', headers: { 'Content-Type': 'application/json', ...authHeader() }, body: JSON.stringify(body) });
  if (!res.ok) throw new Error('failed');
};

export const adminPatchJobRole = async (externalId: string, body: any) => {
  const res = await fetch(`${API_BASE}/admin/catalog/job-roles/${externalId}`, { method: 'PATCH', headers: { 'Content-Type': 'application/json', ...authHeader() }, body: JSON.stringify(body) });
  if (!res.ok) throw new Error('failed');
};

export const adminDeleteJobRole = async (externalId: string) => {
  const res = await fetch(`${API_BASE}/admin/catalog/job-roles/${externalId}`, { method: 'DELETE', headers: authHeader() });
  if (!res.ok) throw new Error('failed');
};

export const adminListContentPresets = async (params: Record<string, string>) => {
  const q = new URLSearchParams(params).toString();
  const res = await fetch(`${API_BASE}/admin/catalog/content-presets?${q}`, { headers: authHeader() });
  if (!res.ok) throw new Error('failed');
  return res.json() as Promise<PageResp<AdminContentPreset>>;
};

export const adminCreateContentPreset = async (body: any) => {
  const res = await fetch(`${API_BASE}/admin/catalog/content-presets`, { method: 'POST', headers: { 'Content-Type': 'application/json', ...authHeader() }, body: JSON.stringify(body) });
  if (!res.ok) throw new Error(await res.text());
};

export const adminPatchContentPreset = async (externalId: string, body: any) => {
  const res = await fetch(`${API_BASE}/admin/catalog/content-presets/${externalId}`, { method: 'PATCH', headers: { 'Content-Type': 'application/json', ...authHeader() }, body: JSON.stringify(body) });
  if (!res.ok) throw new Error(await res.text());
};

export const adminDeleteContentPreset = async (externalId: string) => {
  const res = await fetch(`${API_BASE}/admin/catalog/content-presets/${externalId}`, { method: 'DELETE', headers: authHeader() });
  if (!res.ok) throw new Error('failed');
};

export const adminListTemplateVariants = async (params: Record<string, string>) => {
  const q = new URLSearchParams(params).toString();
  const res = await fetch(`${API_BASE}/admin/catalog/template-variants?${q}`, { headers: authHeader() });
  if (!res.ok) throw new Error('failed');
  return res.json() as Promise<PageResp<AdminTemplateVariant>>;
};

export const adminCreateTemplateVariant = async (body: any) => {
  const res = await fetch(`${API_BASE}/admin/catalog/template-variants`, { method: 'POST', headers: { 'Content-Type': 'application/json', ...authHeader() }, body: JSON.stringify(body) });
  if (!res.ok) throw new Error('failed');
};

export const adminPatchTemplateVariant = async (externalId: string, body: any) => {
  const res = await fetch(`${API_BASE}/admin/catalog/template-variants/${externalId}`, { method: 'PATCH', headers: { 'Content-Type': 'application/json', ...authHeader() }, body: JSON.stringify(body) });
  if (!res.ok) throw new Error('failed');
};

export const adminDeleteTemplateVariant = async (externalId: string) => {
  const res = await fetch(`${API_BASE}/admin/catalog/template-variants/${externalId}`, { method: 'DELETE', headers: authHeader() });
  if (!res.ok) throw new Error('failed');
};

export const adminGenerateTemplateVariants = async (body: any) => {
  const res = await fetch(`${API_BASE}/admin/catalog/template-variants/generate`, { method: 'POST', headers: { 'Content-Type': 'application/json', ...authHeader() }, body: JSON.stringify(body) });
  if (!res.ok) throw new Error(await res.text());
  return res.json();
};

// Users
export interface AdminUser {
  id: number;
  email: string;
  name: string;
  avatarUrl?: string;
  language?: string;
  role: string;
  isActive: boolean;
  lastLoginAt?: string;
}

export const listUsers = async (params: Record<string, string>) => {
  const q = new URLSearchParams(params).toString();
  const res = await fetch(`${API_BASE}/admin/users?${q}`, { headers: authHeader() });
  if (!res.ok) throw new Error('failed');
  return res.json() as Promise<PageResp<AdminUser>>;
};

export const getUser = async (id: number) => {
  const res = await fetch(`${API_BASE}/admin/users/${id}`, { headers: authHeader() });
  if (!res.ok) throw new Error('failed');
  return res.json() as Promise<AdminUser>;
};

export const updateUser = async (id: number, body: Partial<AdminUser>) => {
  const res = await fetch(`${API_BASE}/admin/users/${id}`, { method: 'PATCH', headers: { 'Content-Type': 'application/json', ...authHeader() }, body: JSON.stringify(body) });
  if (!res.ok) throw new Error('failed');
};

export const resetPassword = async (id: number, newPassword: string) => {
  const res = await fetch(`${API_BASE}/admin/users/${id}/reset-password`, { method: 'POST', headers: { 'Content-Type': 'application/json', ...authHeader() }, body: JSON.stringify({ newPassword }) });
  if (!res.ok) throw new Error('failed');
};

export const banUser = async (id: number) => {
  const res = await fetch(`${API_BASE}/admin/users/${id}/ban`, { method: 'POST', headers: authHeader() });
  if (!res.ok) throw new Error('failed');
};
export const unbanUser = async (id: number) => {
  const res = await fetch(`${API_BASE}/admin/users/${id}/unban`, { method: 'POST', headers: authHeader() });
  if (!res.ok) throw new Error('failed');
};

// Resumes
export interface AdminResume {
  id: string;
  userId: number;
  userName?: string;
  title: string;
  templateId: string;
  language?: 'en' | 'zh';
  Theme: { Color?: string; Font?: string; Spacing?: string; FontSize?: string };
  lastModified: number;
}

export const listResumes = async (params: Record<string, string>) => {
  const q = new URLSearchParams(params).toString();
  const res = await fetch(`${API_BASE}/admin/resumes?${q}`, { headers: authHeader() });
  if (!res.ok) throw new Error('failed');
  const data = await res.json();
  const items: AdminResume[] = (data.items || []).map((it: any) => {
    const r = it.resume || {};
    return {
      id: r.ExternalID,
      userId: it.userId,
      userName: it.userName,
      title: r.Title,
      templateId: r.TemplateID,
      language: (r.Language || '') === 'en' ? 'en' : 'zh',
      Theme: { Color: r.Theme?.Color, Font: r.Theme?.Font, Spacing: r.Theme?.Spacing, FontSize: r.Theme?.FontSize },
      lastModified: r.LastModified,
    } as AdminResume;
  });
  return { items, page: data.page, pageSize: data.pageSize, total: data.total, hasNext: data.hasNext, totalPages: data.totalPages } as PageResp<AdminResume>;
};

export const deleteResume = async (id: string) => {
  const res = await fetch(`${API_BASE}/admin/resumes/${id}`, { method: 'DELETE', headers: authHeader() });
  if (!res.ok) throw new Error('failed');
};

export const setResumeVisibility = async (id: string, isPublic: boolean) => {
  const res = await fetch(`${API_BASE}/admin/resumes/${id}/visibility`, { method: 'PATCH', headers: { 'Content-Type': 'application/json', ...authHeader() }, body: JSON.stringify({ isPublic }) });
  if (!res.ok) throw new Error('failed');
};

// Templates
export interface AdminTemplateReq {
  externalId: string;
  name: string;
  tags?: string;
}
export const createTemplate = async (body: AdminTemplateReq) => {
  const res = await fetch(`${API_BASE}/admin/templates`, { method: 'POST', headers: { 'Content-Type': 'application/json', ...authHeader() }, body: JSON.stringify(body) });
  if (!res.ok) throw new Error('failed');
};
export const updateTemplate = async (id: string, body: Partial<AdminTemplateReq>) => {
  const res = await fetch(`${API_BASE}/admin/templates/${id}`, { method: 'PATCH', headers: { 'Content-Type': 'application/json', ...authHeader() }, body: JSON.stringify(body) });
  if (!res.ok) throw new Error('failed');
};
export const deleteTemplate = async (id: string) => {
  const res = await fetch(`${API_BASE}/admin/templates/${id}`, { method: 'DELETE', headers: authHeader() });
  if (!res.ok) throw new Error('failed');
};

// Shares
export interface AdminShare {
  id: number;
  resumeId: number;
  userId?: number;
  userName?: string;
  slug: string;
  isPublic: boolean;
  url?: string;
}
export const listShares = async (params: Record<string, string>) => {
  const q = new URLSearchParams(params).toString();
  const res = await fetch(`${API_BASE}/admin/share-links?${q}`, { headers: authHeader() });
  if (!res.ok) throw new Error('failed');
  return res.json() as Promise<PageResp<AdminShare>>;
};
export const updateShare = async (slug: string, isPublic: boolean) => {
  const res = await fetch(`${API_BASE}/admin/share-links/${slug}`, { method: 'PATCH', headers: { 'Content-Type': 'application/json', ...authHeader() }, body: JSON.stringify({ isPublic }) });
  if (!res.ok) throw new Error('failed');
};
export const deleteShare = async (slug: string) => {
  const res = await fetch(`${API_BASE}/admin/share-links/${slug}`, { method: 'DELETE', headers: authHeader() });
  if (!res.ok) throw new Error('failed');
};
