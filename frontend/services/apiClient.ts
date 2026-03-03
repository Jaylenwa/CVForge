import { API_BASE } from '../config';

export class ApiError extends Error {
  status: number;
  data: unknown;

  constructor(message: string, status: number, data: unknown) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
    this.data = data;
  }
}

const buildUrl = (path: string) => {
  if (path.startsWith('http://') || path.startsWith('https://')) return path;
  if (path.startsWith('/')) return `${API_BASE}${path}`;
  return `${API_BASE}/${path}`;
};

const getAccessToken = () => localStorage.getItem('token') || '';

const withAuthHeader = (headers?: HeadersInit): HeadersInit => {
  const token = getAccessToken();
  const h: Record<string, string> = {};
  if (headers) {
    if (Array.isArray(headers)) {
      for (const [k, v] of headers) h[k] = v;
    } else if (headers instanceof Headers) {
      headers.forEach((v, k) => {
        h[k] = v;
      });
    } else {
      Object.assign(h, headers);
    }
  }
  if (token) h.Authorization = `Bearer ${token}`;
  return h;
};

const safeJson = async (res: Response): Promise<unknown> => {
  try {
    return await res.json();
  } catch {
    return null;
  }
};

const errorMessageFrom = (data: unknown): string => {
  if (!data || typeof data !== 'object') return '';
  const msg = (data as any).error;
  return typeof msg === 'string' ? msg.trim() : '';
};

export const apiRequest = async (
  path: string,
  init?: RequestInit & { auth?: boolean }
): Promise<Response> => {
  const url = buildUrl(path);
  const auth = !!init?.auth;
  const headers = auth ? withAuthHeader(init?.headers) : (init?.headers || undefined);
  return fetch(url, { ...init, headers });
};

export const apiJson = async <T>(
  path: string,
  init?: RequestInit & { auth?: boolean }
): Promise<T> => {
  const res = await apiRequest(path, init);
  if (res.ok) return (await safeJson(res)) as T;
  const data = await safeJson(res);
  const msg = errorMessageFrom(data);
  throw new ApiError(msg || `request failed: ${res.status}`, res.status, data);
};

export const apiVoid = async (
  path: string,
  init?: RequestInit & { auth?: boolean }
): Promise<void> => {
  const res = await apiRequest(path, init);
  if (res.ok) return;
  const data = await safeJson(res);
  const msg = errorMessageFrom(data);
  throw new ApiError(msg || `request failed: ${res.status}`, res.status, data);
};

export const apiJsonOrNull = async <T>(
  path: string,
  init?: RequestInit & { auth?: boolean }
): Promise<T | null> => {
  const res = await apiRequest(path, init);
  if (!res.ok) return null;
  return (await safeJson(res)) as T;
};
