import { apiJson } from './apiClient';

export const polishText = async (text: string, tone: 'professional' | 'creative' = 'professional'): Promise<string> => {
  try {
    const data = await apiJson<{ text?: string }>('/ai/polish', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ text, tone })
    });
    return (data.text as string) || text;
  } catch {
    return text;
  }
};

export const generateSummary = async (job: string, skills: string): Promise<string> => {
  try {
    const data = await apiJson<{ text?: string }>('/ai/summary', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ job, skills })
    });
    return (data.text as string) || '';
  } catch {
    return '';
  }
};
