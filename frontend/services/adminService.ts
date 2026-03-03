import { apiJson, apiVoid } from './apiClient';

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
  return apiJson<AdminStats>(`/admin/stats?days=${days}`, { auth: true });
};

export interface AdminJobCategory {
  id: number;
  name?: string;
  names?: Record<string, string>;
  parentId?: number | null;
  orderNum?: number;
  isActive: boolean;
}

export interface AdminCreateJobCategoryReq {
  name?: string;
  names?: Record<string, string>;
  parentId?: number | null;
  orderNum?: number;
  isActive?: boolean;
}

export interface AdminJobRole {
  id: number;
  categoryId: number;
  name?: string;
  names?: Record<string, string>;
  orderNum?: number;
  isActive: boolean;
}

export interface AdminCreateJobRoleReq {
  categoryId: number;
  name?: string;
  names?: Record<string, string>;
  orderNum?: number;
  isActive?: boolean;
}

export interface AdminContentPreset {
  ID: number;
  Name: string;
  Language?: string;
  RoleID: number;
  DataJSON?: string;
  IsActive: boolean;
}

export interface AdminCreateContentPresetReq {
  name: string;
  language?: string;
  roleId: number;
  dataJson?: string;
  isActive?: boolean;
}

export const adminListJobCategories = async (params: Record<string, string>) => {
  const q = new URLSearchParams(params).toString();
  return apiJson<PageResp<AdminJobCategory>>(`/admin/taxonomy/categories?${q}`, { auth: true });
};

export const adminCreateJobCategory = async (body: AdminCreateJobCategoryReq) => {
  await apiVoid('/admin/taxonomy/categories', {
    method: 'POST',
    auth: true,
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body)
  });
};

export const adminPatchJobCategory = async (id: number, body: any) => {
  await apiVoid(`/admin/taxonomy/categories/${id}`, {
    method: 'PATCH',
    auth: true,
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body)
  });
};

export const adminDeleteJobCategory = async (id: number) => {
  await apiVoid(`/admin/taxonomy/categories/${id}`, { method: 'DELETE', auth: true });
};

export const adminListJobRoles = async (params: Record<string, string>) => {
  const q = new URLSearchParams(params).toString();
  return apiJson<PageResp<AdminJobRole>>(`/admin/taxonomy/roles?${q}`, { auth: true });
};

export const adminCreateJobRole = async (body: AdminCreateJobRoleReq) => {
  await apiVoid('/admin/taxonomy/roles', {
    method: 'POST',
    auth: true,
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body)
  });
};

export const adminPatchJobRole = async (id: number, body: any) => {
  await apiVoid(`/admin/taxonomy/roles/${id}`, {
    method: 'PATCH',
    auth: true,
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body)
  });
};

export const adminDeleteJobRole = async (id: number) => {
  await apiVoid(`/admin/taxonomy/roles/${id}`, { method: 'DELETE', auth: true });
};

export const adminListContentPresets = async (params: Record<string, string>) => {
  const q = new URLSearchParams(params).toString();
  return apiJson<PageResp<AdminContentPreset>>(`/admin/presets?${q}`, { auth: true });
};

export const adminCreateContentPreset = async (body: AdminCreateContentPresetReq) => {
  await apiVoid('/admin/presets', {
    method: 'POST',
    auth: true,
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body)
  });
};

export const adminPatchContentPreset = async (id: number, body: any) => {
  await apiVoid(`/admin/presets/${id}`, {
    method: 'PATCH',
    auth: true,
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body)
  });
};

export const adminDeleteContentPreset = async (id: number) => {
  await apiVoid(`/admin/presets/${id}`, { method: 'DELETE', auth: true });
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
  providers?: string[];
  loginProvider?: string;
}

export const listUsers = async (params: Record<string, string>) => {
  const q = new URLSearchParams(params).toString();
  return apiJson<PageResp<AdminUser>>(`/admin/users?${q}`, { auth: true });
};

export const getUser = async (id: number) => {
  return apiJson<AdminUser>(`/admin/users/${id}`, { auth: true });
};

export const updateUser = async (id: number, body: Partial<AdminUser>) => {
  await apiVoid(`/admin/users/${id}`, {
    method: 'PATCH',
    auth: true,
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body)
  });
};

export const resetPassword = async (id: number, newPassword: string) => {
  await apiVoid(`/admin/users/${id}/reset-password`, {
    method: 'POST',
    auth: true,
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ newPassword })
  });
};

export const banUser = async (id: number) => {
  await apiVoid(`/admin/users/${id}/ban`, { method: 'POST', auth: true });
};
export const unbanUser = async (id: number) => {
  await apiVoid(`/admin/users/${id}/unban`, { method: 'POST', auth: true });
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
  const data = await apiJson<any>(`/admin/resumes?${q}`, { auth: true });
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
  await apiVoid(`/admin/resumes/${id}`, { method: 'DELETE', auth: true });
};

export const setResumeVisibility = async (id: number, isPublic: boolean) => {
  await apiVoid(`/admin/resumes/${id}/visibility`, {
    method: 'PATCH',
    auth: true,
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ isPublic })
  });
};

// Templates
export interface AdminTemplateReq {
  externalId: string;
  name?: string;
  names?: Record<string, string>;
}
export const createTemplate = async (body: AdminTemplateReq) => {
  await apiVoid('/admin/templates', {
    method: 'POST',
    auth: true,
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body)
  });
};
export const updateTemplate = async (id: string, body: Partial<AdminTemplateReq>) => {
  await apiVoid(`/admin/templates/${id}`, {
    method: 'PATCH',
    auth: true,
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body)
  });
};
export const deleteTemplate = async (id: string) => {
  await apiVoid(`/admin/templates/${id}`, { method: 'DELETE', auth: true });
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
  return apiJson<PageResp<AdminShare>>(`/admin/share-links?${q}`, { auth: true });
};
export const updateShare = async (slug: string, isPublic: boolean) => {
  await apiVoid(`/admin/share-links/${slug}`, {
    method: 'PATCH',
    auth: true,
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ isPublic })
  });
};
export const deleteShare = async (slug: string) => {
  await apiVoid(`/admin/share-links/${slug}`, { method: 'DELETE', auth: true });
};
