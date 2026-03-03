import { ApiError, apiJson } from './apiClient';

export interface WeChatMPScene {
  scene: string;
  qrUrl: string;
  expiresIn: number;
}

export interface WeChatMPSceneStatus {
  status: 'pending' | 'ok' | 'expired';
  ott?: string;
}

export const createWeChatMPScene = async (): Promise<WeChatMPScene> => {
  return apiJson<WeChatMPScene>('/auth/wechat-mp/scene/create', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' }
  });
};

export const getWeChatMPSceneStatus = async (scene: string): Promise<WeChatMPSceneStatus> => {
  try {
    return await apiJson<WeChatMPSceneStatus>(`/auth/wechat-mp/scene/${encodeURIComponent(scene)}/status`);
  } catch (e) {
    if (e instanceof ApiError && e.status === 404) return { status: 'expired' };
    throw e;
  }
};

export const consumeOtt = async (ott: string): Promise<{ accessToken: string; refreshToken: string }> => {
  const data = await apiJson<{ accessToken: string; refreshToken: string }>('/auth/wechat/consume-ott', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ ott })
  });
  return { accessToken: data.accessToken as string, refreshToken: data.refreshToken as string };
};
