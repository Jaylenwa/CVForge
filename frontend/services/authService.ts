import { ApiError, apiJson, apiVoid } from './apiClient';

export const sendVerificationCode = async (email: string): Promise<void> => {
  await apiVoid('/auth/send-code', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email })
  });
};

export const loginUser = async (
  email: string,
  password: string
): Promise<{ success: boolean; token?: string; error?: 'invalid_credentials' | 'unknown' }> => {
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
  } catch (e) {
    if (e instanceof ApiError) {
      if (e.status === 401) return { success: false, error: 'invalid_credentials' };
    }
    return { success: false, error: 'unknown' };
  }
};

export const registerUser = async (
  email: string,
  code: string,
  password: string,
  name?: string
): Promise<{ success: boolean; token?: string; error?: 'invalid_code' | 'invalid_email' | 'email_exists' | 'unknown' }> => {
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
      if (e.status === 400) {
        const msg = String(e.message || '').toLowerCase();
        if (msg.includes('invalid email')) return { success: false, error: 'invalid_email' };
        if (msg.includes('invalid code')) return { success: false, error: 'invalid_code' };
      }
    }
    return { success: false, error: 'unknown' };
  }
};
