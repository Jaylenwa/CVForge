const API = 'http://localhost:8080/api/v1';

export const polishText = async (text: string, tone: 'professional' | 'creative' = 'professional'): Promise<string> => {
  try {
    const res = await fetch(`${API}/ai/polish`, {
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

export const generateSummary = async (jobTitle: string, skills: string): Promise<string> => {
  try {
    const res = await fetch(`${API}/ai/summary`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ jobTitle, skills })
    });
    const data = await res.json();
    return (data.text as string) || "";
  } catch {
    return "";
  }
};
