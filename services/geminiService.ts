import { GoogleGenAI } from "@google/genai";

// Initialize the client safely
// Note: In a real production app, calls should go through a backend to protect the API Key.
// For this demo, we assume process.env.API_KEY is available.
const apiKey = process.env.API_KEY || ''; 
const ai = new GoogleGenAI({ apiKey });

export const polishText = async (text: string, tone: 'professional' | 'creative' = 'professional'): Promise<string> => {
  if (!apiKey) {
    console.warn("API Key is missing. Returning original text.");
    return text;
  }

  try {
    const prompt = `
      Act as a professional resume writer. Rewrite the following text to be more ${tone}, concise, and impactful. 
      Focus on action verbs and achievements. 
      Return ONLY the rewritten text, no explanations.
      
      Original Text: "${text}"
    `;

    const response = await ai.models.generateContent({
      model: 'gemini-2.5-flash',
      contents: prompt,
    });

    return response.text?.trim() || text;
  } catch (error) {
    console.error("Gemini Polish Error:", error);
    return text; // Fallback
  }
};

export const generateSummary = async (jobTitle: string, skills: string): Promise<string> => {
  if (!apiKey) return "API Key missing.";

  try {
    const prompt = `
      Write a professional resume summary (max 3 sentences) for a ${jobTitle} with skills in ${skills}.
      Return ONLY the summary text.
    `;

    const response = await ai.models.generateContent({
      model: 'gemini-2.5-flash',
      contents: prompt,
    });

    return response.text?.trim() || "";
  } catch (error) {
    console.error("Gemini Summary Error:", error);
    return "";
  }
};
