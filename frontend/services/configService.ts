import { API_BASE } from '../config';
import { AuthConfig, SystemConfig } from '../types';

export const getAuthConfig = async (): Promise<AuthConfig> => {
  try {
    const res = await fetch(`${API_BASE}/auth/config`);
    if (!res.ok) throw new Error('Network response was not ok');
    return await res.json();
  } catch (error) {
    console.error('Failed to fetch auth config:', error);
    return {
      enableEmailVerification: false,
      enableWeChatLogin: true,
      enableGithubLogin: true,
      weChatAppID: '',
      githubClientID: ''
    };
  }
};

export const getSystemConfigs = async (): Promise<SystemConfig[]> => {
  const token = localStorage.getItem('token');
  const res = await fetch(`${API_BASE}/admin/configs`, {
    headers: { 'Authorization': `Bearer ${token}` }
  });
  if (!res.ok) throw new Error('Failed to fetch configs');
  return res.json();
};

export const updateSystemConfigs = async (configs: SystemConfig[]): Promise<void> => {
  const token = localStorage.getItem('token');
  const res = await fetch(`${API_BASE}/admin/configs`, {
    method: 'PUT',
    headers: { 
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify({ configs })
  });
  if (!res.ok) throw new Error('Failed to update configs');
};
