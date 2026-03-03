import { AuthConfig, SystemConfig } from '../types';
import { apiJson, apiVoid } from './apiClient';

export const getAuthConfig = async (): Promise<AuthConfig> => {
  try {
    return await apiJson<AuthConfig>('/auth/config');
  } catch {
    return {
      enableEmailVerification: false,
      enableWeChatMPLogin: false,
      enableGithubLogin: true,
      enablePricingPage: false,
      githubClientID: ''
    };
  }
};

export const getSystemConfigs = async (): Promise<SystemConfig[]> => {
  return apiJson<SystemConfig[]>('/admin/configs', { auth: true });
};

export const updateSystemConfigs = async (configs: SystemConfig[]): Promise<void> => {
  await apiVoid('/admin/configs', {
    method: 'PUT',
    auth: true,
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ configs })
  });
};
