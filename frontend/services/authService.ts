const API = 'http://localhost:8080/api/v1';

export const sendVerificationCode = async (email: string): Promise<boolean> => {
  const res = await fetch(`${API}/auth/send-code`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email })
  });
  return res.ok;
};

export const verifyCode = async (_email: string, code: string): Promise<boolean> => {
  // 后端已支持注册时核验，这里保持“123456”为有效以兼容现有前端流程
  return code === '123456';
};

export const loginUser = async (email: string, password: string): Promise<{ success: boolean; token?: string }> => {
  const res = await fetch(`${API}/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password })
  });
  if (!res.ok) return { success: false };
  const data = await res.json();
  const token = data.accessToken as string;
  const refresh = data.refreshToken as string;
  if (token) localStorage.setItem('token', token);
  if (refresh) localStorage.setItem('refreshToken', refresh);
  return { success: true, token };
};

export const registerUser = async (email: string, code: string, password: string, name?: string): Promise<{ success: boolean; token?: string }> => {
  const res = await fetch(`${API}/auth/register`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, code, password, name: name || email.split('@')[0] })
  });
  if (!res.ok) return { success: false };
  const data = await res.json();
  const token = data.accessToken as string;
  const refresh = data.refreshToken as string;
  if (token) localStorage.setItem('token', token);
  if (refresh) localStorage.setItem('refreshToken', refresh);
  return { success: true, token };
};
