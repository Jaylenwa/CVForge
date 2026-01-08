import { API_BASE } from '../config';

export const polishText = async (text: string, tone: 'professional' | 'creative' = 'professional'): Promise<string> => {
  try {
    const res = await fetch(`${API_BASE}/ai/polish`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ text, tone })
    });
    const data = await res.json();
    return (data.text as string) || text;
  } catch {
    return text;
  }
};

export const generateSummary = async (job: string, skills: string): Promise<string> => {
  try {
    const res = await fetch(`${API_BASE}/ai/summary`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ job, skills })
    });
    const data = await res.json();
    return (data.text as string) || "";
  } catch {
    return "";
  }
};
