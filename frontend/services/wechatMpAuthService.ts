import { API_BASE } from '../config';

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
  const res = await fetch(`${API_BASE}/auth/wechat-mp/scene/create`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' }
  });
  if (!res.ok) {
    const msg = await safeReadError(res);
    throw new Error(msg || 'Failed to create scene');
  }
  return res.json();
};

export const getWeChatMPSceneStatus = async (scene: string): Promise<WeChatMPSceneStatus> => {
  const res = await fetch(`${API_BASE}/auth/wechat-mp/scene/${encodeURIComponent(scene)}/status`);
  if (res.status === 404) return { status: 'expired' };
  if (!res.ok) {
    const msg = await safeReadError(res);
    throw new Error(msg || 'Failed to fetch status');
  }
  return res.json();
};

export const consumeOtt = async (ott: string): Promise<{ accessToken: string; refreshToken: string }> => {
  const res = await fetch(`${API_BASE}/auth/wechat/consume-ott`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ ott })
  });
  if (!res.ok) {
    const msg = await safeReadError(res);
    throw new Error(msg || 'Failed to consume ott');
  }
  const data = await res.json();
  return { accessToken: data.accessToken as string, refreshToken: data.refreshToken as string };
};

const safeReadError = async (res: Response): Promise<string> => {
  try {
    const data = (await res.json()) as any;
    if (data && typeof data.error === 'string' && data.error.trim()) {
      return data.error.trim();
    }
  } catch {
  }
  return '';
};
