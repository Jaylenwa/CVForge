
import React, { useState } from 'react';
import { GoogleGenAI } from "@google/genai";
import { MessageCircle, Send, X, Bot } from 'lucide-react';

export const GeminiAssistant: React.FC = () => {
  const [isOpen, setIsOpen] = useState(false);
  const [query, setQuery] = useState('');
  const [response, setResponse] = useState('Hello! I am your AI security guide. How can I help you with your account?');
  const [isLoading, setIsLoading] = useState(false);

  const askGemini = async () => {
    if (!query.trim()) return;
    
    setIsLoading(true);
    setResponse('Analysing your request...');
    
    try {
      const ai = new GoogleGenAI({ apiKey: process.env.API_KEY });
      const res = await ai.models.generateContent({
        model: 'gemini-3-flash-preview',
        contents: query,
        config: {
          systemInstruction: "You are a friendly authentication and cybersecurity expert. Help the user understand login/signup issues, MFA benefits, or password security. Keep it brief.",
        }
      });
      setResponse(res.text || "I'm sorry, I couldn't process that.");
    } catch (error) {
      setResponse("My core processors are busy, please try again later.");
    } finally {
      setIsLoading(false);
      setQuery('');
    }
  };

  return (
    <div className="fixed bottom-6 right-6 z-50">
      {!isOpen ? (
        <button 
          onClick={() => setIsOpen(true)}
          className="w-14 h-14 bg-indigo-600 rounded-full flex items-center justify-center shadow-lg shadow-indigo-500/40 hover:scale-110 active:scale-95 transition-all text-white"
        >
          <MessageCircle className="w-7 h-7" />
        </button>
      ) : (
        <div className="bg-slate-900/95 backdrop-blur-xl border border-white/10 rounded-2xl w-80 shadow-2xl overflow-hidden animate-in slide-in-from-bottom-5 duration-300">
          <div className="bg-indigo-600 p-4 flex justify-between items-center">
            <div className="flex items-center space-x-2">
              <Bot className="w-5 h-5 text-white" />
              <span className="text-white font-medium text-sm">Security Guide</span>
            </div>
            <button onClick={() => setIsOpen(false)}><X className="w-5 h-5 text-white/80" /></button>
          </div>
          
          <div className="p-4 h-64 overflow-y-auto text-sm text-slate-300 whitespace-pre-wrap">
            {response}
          </div>

          <div className="p-3 border-t border-white/10 flex items-center space-x-2">
            <input 
              type="text" 
              className="flex-1 bg-white/5 rounded-lg py-1.5 px-3 text-xs text-white focus:outline-none focus:ring-1 focus:ring-indigo-500"
              placeholder="Ask me anything..."
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              onKeyDown={(e) => e.key === 'Enter' && askGemini()}
            />
            <button 
              onClick={askGemini}
              disabled={isLoading}
              className="p-1.5 bg-indigo-600 rounded-lg text-white hover:bg-indigo-500 disabled:opacity-50"
            >
              <Send className="w-4 h-4" />
            </button>
          </div>
        </div>
      )}
    </div>
  );
};
