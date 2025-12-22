import { API_BASE } from '../config';

const authHeader = () => ({ Authorization: `Bearer ${localStorage.getItem('token') || ''}` });

export interface PageResp<T> {
  items: T[];
  page: number;
  pageSize: number;
  total: number;
}

// Stats
export interface AdminStats {
  totals: { users: number; resumes: number };
  trend: { dates: string[]; users: number[]; resumes: number[] };
  generatedAt: number;
}

export const getAdminStats = async (days = 14): Promise<AdminStats> => {
  const res = await fetch(`${API_BASE}/admin/stats?days=${days}`, { headers: authHeader() });
  if (!res.ok) throw new Error('Failed to fetch stats');
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
  themeConfig: { color: string; fontFamily: string; spacing: string };
  lastModified: number;
}

export const listResumes = async (params: Record<string, string>) => {
  const q = new URLSearchParams(params).toString();
  const res = await fetch(`${API_BASE}/admin/resumes?${q}`, { headers: authHeader() });
  if (!res.ok) throw new Error('failed');
  return res.json() as Promise<PageResp<AdminResume>>;
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
  popularity?: number;
  isPremium?: boolean;
  category?: string;
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
  slug: string;
  isPublic: boolean;
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
