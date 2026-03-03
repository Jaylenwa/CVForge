import { ApiError, apiJson, apiVoid } from './apiClient';

export const sendVerificationCode = async (email: string): Promise<boolean> => {
  try {
    await apiVoid('/auth/send-code', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email })
    });
    return true;
  } catch {
    return false;
  }
};

export const loginUser = async (email: string, password: string): Promise<{ success: boolean; token?: string }> => {
  try {
    const data = await apiJson<{ accessToken?: string; refreshToken?: string }>('/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password })
    });
    const token = (data.accessToken || '') as string;
    const refresh = (data.refreshToken || '') as string;
    if (token) localStorage.setItem('token', token);
    if (refresh) localStorage.setItem('refreshToken', refresh);
    return { success: !!token, token };
  } catch {
    return { success: false };
  }
};

export const registerUser = async (
  email: string,
  code: string,
  password: string,
  name?: string
): Promise<{ success: boolean; token?: string; error?: 'invalid_code' | 'email_exists' | 'unknown' }> => {
  try {
    const data = await apiJson<{ accessToken?: string; refreshToken?: string }>('/auth/register', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, code, password, name: name || email.split('@')[0] })
    });
    const token = (data.accessToken || '') as string;
    const refresh = (data.refreshToken || '') as string;
    if (token) localStorage.setItem('token', token);
    if (refresh) localStorage.setItem('refreshToken', refresh);
    return { success: !!token, token };
  } catch (e) {
    if (e instanceof ApiError) {
      if (e.status === 409) return { success: false, error: 'email_exists' };
      if (e.status === 400) return { success: false, error: 'invalid_code' };
    }
    return { success: false, error: 'unknown' };
  }
};
