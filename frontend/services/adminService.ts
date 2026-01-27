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

export const importDefaultSeed = async () => {
  const res = await fetch(`${API_BASE}/admin/seed/import-default`, { method: 'POST', headers: authHeader() });
  if (!res.ok) throw new Error('failed');
  return res.json() as Promise<{ success: boolean; counts?: any }>;
};

export interface AdminJobCategory {
  ID: number;
  Name: string;
  ParentID?: number | null;
  OrderNum?: number;
  IsActive: boolean;
}

export interface AdminCreateJobCategoryReq {
  name: string;
  parentId?: number | null;
  orderNum?: number;
  isActive?: boolean;
}

export interface AdminJobRole {
  ID: number;
  CategoryID: number;
  Name: string;
  Tags?: string;
  OrderNum?: number;
  IsActive: boolean;
}

export interface AdminCreateJobRoleReq {
  categoryId: number;
  name: string;
  tags?: string;
  orderNum?: number;
  isActive?: boolean;
}

export interface AdminContentPreset {
  ID: number;
  Name: string;
  Language?: string;
  RoleID: number;
  Tags?: string;
  DataJSON?: string;
  IsActive: boolean;
}

export interface AdminCreateContentPresetReq {
  name: string;
  language?: string;
  roleId: number;
  tags?: string;
  dataJson?: string;
  isActive?: boolean;
}

export const adminListJobCategories = async (params: Record<string, string>) => {
  const q = new URLSearchParams(params).toString();
  const res = await fetch(`${API_BASE}/admin/taxonomy/categories?${q}`, { headers: authHeader() });
  if (!res.ok) throw new Error('failed');
  return res.json() as Promise<PageResp<AdminJobCategory>>;
};

export const adminCreateJobCategory = async (body: AdminCreateJobCategoryReq) => {
  const res = await fetch(`${API_BASE}/admin/taxonomy/categories`, { method: 'POST', headers: { 'Content-Type': 'application/json', ...authHeader() }, body: JSON.stringify(body) });
  if (!res.ok) throw new Error('failed');
};

export const adminPatchJobCategory = async (id: number, body: any) => {
  const res = await fetch(`${API_BASE}/admin/taxonomy/categories/${id}`, { method: 'PATCH', headers: { 'Content-Type': 'application/json', ...authHeader() }, body: JSON.stringify(body) });
  if (!res.ok) throw new Error('failed');
};

export const adminDeleteJobCategory = async (id: number) => {
  const res = await fetch(`${API_BASE}/admin/taxonomy/categories/${id}`, { method: 'DELETE', headers: authHeader() });
  if (!res.ok) throw new Error('failed');
};

export const adminListJobRoles = async (params: Record<string, string>) => {
  const q = new URLSearchParams(params).toString();
  const res = await fetch(`${API_BASE}/admin/taxonomy/roles?${q}`, { headers: authHeader() });
  if (!res.ok) throw new Error('failed');
  return res.json() as Promise<PageResp<AdminJobRole>>;
};

export const adminCreateJobRole = async (body: AdminCreateJobRoleReq) => {
  const res = await fetch(`${API_BASE}/admin/taxonomy/roles`, { method: 'POST', headers: { 'Content-Type': 'application/json', ...authHeader() }, body: JSON.stringify(body) });
  if (!res.ok) throw new Error('failed');
};

export const adminPatchJobRole = async (id: number, body: any) => {
  const res = await fetch(`${API_BASE}/admin/taxonomy/roles/${id}`, { method: 'PATCH', headers: { 'Content-Type': 'application/json', ...authHeader() }, body: JSON.stringify(body) });
  if (!res.ok) throw new Error('failed');
};

export const adminDeleteJobRole = async (id: number) => {
  const res = await fetch(`${API_BASE}/admin/taxonomy/roles/${id}`, { method: 'DELETE', headers: authHeader() });
  if (!res.ok) throw new Error('failed');
};

export const adminListContentPresets = async (params: Record<string, string>) => {
  const q = new URLSearchParams(params).toString();
  const res = await fetch(`${API_BASE}/admin/presets?${q}`, { headers: authHeader() });
  if (!res.ok) throw new Error('failed');
  return res.json() as Promise<PageResp<AdminContentPreset>>;
};

export const adminCreateContentPreset = async (body: AdminCreateContentPresetReq) => {
  const res = await fetch(`${API_BASE}/admin/presets`, { method: 'POST', headers: { 'Content-Type': 'application/json', ...authHeader() }, body: JSON.stringify(body) });
  if (!res.ok) throw new Error(await res.text());
};

export const adminPatchContentPreset = async (id: number, body: any) => {
  const res = await fetch(`${API_BASE}/admin/presets/${id}`, { method: 'PATCH', headers: { 'Content-Type': 'application/json', ...authHeader() }, body: JSON.stringify(body) });
  if (!res.ok) throw new Error(await res.text());
};

export const adminDeleteContentPreset = async (id: number) => {
  const res = await fetch(`${API_BASE}/admin/presets/${id}`, { method: 'DELETE', headers: authHeader() });
  if (!res.ok) throw new Error('failed');
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
  id: number;
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
      id: r.ID ?? r.id,
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

export const deleteResume = async (id: number) => {
  const res = await fetch(`${API_BASE}/admin/resumes/${id}`, { method: 'DELETE', headers: authHeader() });
  if (!res.ok) throw new Error('failed');
};

export const setResumeVisibility = async (id: number, isPublic: boolean) => {
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
