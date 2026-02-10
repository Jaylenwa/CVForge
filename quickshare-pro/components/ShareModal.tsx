
import React, { useState, useEffect } from 'react';
import { GoogleGenAI } from "@google/genai";

interface ShareModalProps {
  isOpen: boolean;
  onClose: () => void;
  fileTitle: string;
}

const ShareModal: React.FC<ShareModalProps> = ({ isOpen, onClose, fileTitle }) => {
  const [isCopied, setIsCopied] = useState(false);
  const [showSettings, setShowSettings] = useState(false);
  const [privacyMode, setPrivacyMode] = useState<'anyone' | 'restricted'>('anyone');
  const [hasPassword, setHasPassword] = useState(false);
  const [expiry, setExpiry] = useState('');
  const [aiSummary, setAiSummary] = useState<string | null>(null);
  const [isGeneratingSummary, setIsGeneratingSummary] = useState(false);

  const shareUrl = `https://share.pro/d/v/${Math.random().toString(36).substring(7)}`;

  useEffect(() => {
    if (isOpen) {
      document.body.style.overflow = 'hidden';
      generateSmartContext();
    } else {
      document.body.style.overflow = 'unset';
    }
  }, [isOpen]);

  const generateSmartContext = async () => {
    setIsGeneratingSummary(true);
    try {
      const ai = new GoogleGenAI({ apiKey: process.env.API_KEY });
      const response = await ai.models.generateContent({
        model: 'gemini-3-flash-preview',
        contents: `Generate a short, friendly, professional 1-sentence description for sharing a file named "${fileTitle}". Focus on making the recipient feel invited to view it. Keep it under 15 words.`,
      });
      setAiSummary(response.text || "Here's the document we discussed.");
    } catch (error) {
      console.error("AI summary failed", error);
    } finally {
      setIsGeneratingSummary(false);
    }
  };

  const handleCopy = () => {
    navigator.clipboard.writeText(shareUrl);
    setIsCopied(true);
    setTimeout(() => setIsCopied(false), 2000);
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 sm:p-6">
      <div 
        className="absolute inset-0 bg-slate-900/60 backdrop-blur-sm transition-opacity duration-300"
        onClick={onClose}
      />
      
      <div className="relative bg-white w-full max-w-lg rounded-[2rem] shadow-2xl overflow-hidden border border-slate-200 animate-in fade-in zoom-in duration-200">
        {/* Header */}
        <div className="px-8 pt-8 pb-4 flex items-center justify-between">
          <div>
            <h2 className="text-2xl font-bold text-slate-800">Share with others</h2>
            <p className="text-slate-500 text-sm mt-0.5">Control how others access this file</p>
          </div>
          <button 
            onClick={onClose}
            className="w-10 h-10 flex items-center justify-center rounded-full hover:bg-slate-100 text-slate-400 transition-colors"
          >
            <i className="fa-solid fa-xmark text-xl"></i>
          </button>
        </div>

        <div className="px-8 py-2">
            <div className="bg-indigo-50/50 rounded-2xl p-4 border border-indigo-100 mb-6 group transition-all hover:bg-indigo-50">
                <div className="flex items-center gap-3 mb-2">
                    <div className="bg-indigo-100 text-indigo-600 w-8 h-8 rounded-lg flex items-center justify-center">
                        <i className="fa-solid fa-wand-sparkles text-sm"></i>
                    </div>
                    <span className="text-xs font-bold text-indigo-600 uppercase tracking-tighter">AI Suggested Message</span>
                </div>
                {isGeneratingSummary ? (
                    <div className="h-5 w-3/4 bg-slate-200 rounded animate-pulse"></div>
                ) : (
                    <p className="text-slate-600 italic text-sm">"{aiSummary}"</p>
                )}
            </div>
        </div>

        {/* Main Link Section */}
        <div className="px-8 pb-6">
          <div className="space-y-4">
            <div className="relative group">
              <label className="text-xs font-bold text-slate-400 uppercase tracking-widest mb-1.5 block px-1">
                Link to share
              </label>
              <div className="flex items-center gap-3">
                <div className="flex-1 relative flex items-center">
                  <input 
                    readOnly
                    value={shareUrl}
                    className="w-full bg-slate-50 border-2 border-slate-100 rounded-2xl py-4 pl-5 pr-12 text-slate-700 font-medium focus:outline-none focus:border-indigo-500 transition-all cursor-default"
                  />
                  <div className="absolute right-4 text-slate-300 group-hover:text-indigo-400 transition-colors">
                    <i className="fa-solid fa-link"></i>
                  </div>
                </div>
                <button 
                  onClick={handleCopy}
                  className={`relative flex items-center justify-center gap-2 h-[58px] px-8 rounded-2xl font-bold text-white transition-all transform active:scale-95 shadow-lg ${
                    isCopied 
                    ? 'bg-emerald-500 shadow-emerald-200' 
                    : 'bg-indigo-600 hover:bg-indigo-700 shadow-indigo-200'
                  }`}
                >
                  {isCopied ? (
                    <>
                      <i className="fa-solid fa-check animate-bounce"></i>
                      <span>Copied!</span>
                    </>
                  ) : (
                    <>
                      <i className="fa-solid fa-copy"></i>
                      <span>Copy</span>
                    </>
                  )}
                </button>
              </div>
            </div>

            <div className="flex items-center justify-between px-1">
              <div className="flex items-center gap-2 text-sm text-slate-600 font-medium">
                <div className="w-2 h-2 rounded-full bg-emerald-500"></div>
                <span>Anyone with this link can view</span>
              </div>
              <button 
                onClick={() => setShowSettings(!showSettings)}
                className={`text-sm font-semibold flex items-center gap-1.5 transition-colors ${showSettings ? 'text-indigo-600' : 'text-slate-400 hover:text-slate-600'}`}
              >
                <i className={`fa-solid fa-sliders transition-transform ${showSettings ? 'rotate-90' : ''}`}></i>
                Settings
              </button>
            </div>
          </div>
        </div>

        {/* Expandable Advanced Settings */}
        <div className={`overflow-hidden transition-all duration-300 ease-in-out ${showSettings ? 'max-h-[500px] opacity-100 mb-6' : 'max-h-0 opacity-0'}`}>
          <div className="px-8 border-t border-slate-100 pt-6 space-y-6">
            
            {/* Setting: Privacy Mode */}
            <div className="grid grid-cols-2 gap-3 p-1 bg-slate-50 rounded-2xl border border-slate-200">
                <button 
                    onClick={() => setPrivacyMode('anyone')}
                    className={`flex items-center justify-center gap-2 py-2.5 rounded-xl text-sm font-semibold transition-all ${privacyMode === 'anyone' ? 'bg-white shadow-sm text-slate-800' : 'text-slate-400 hover:text-slate-600'}`}
                >
                    <i className="fa-solid fa-globe"></i> Public
                </button>
                <button 
                    onClick={() => setPrivacyMode('restricted')}
                    className={`flex items-center justify-center gap-2 py-2.5 rounded-xl text-sm font-semibold transition-all ${privacyMode === 'restricted' ? 'bg-white shadow-sm text-slate-800' : 'text-slate-400 hover:text-slate-600'}`}
                >
                    <i className="fa-solid fa-lock"></i> Restricted
                </button>
            </div>

            <div className="space-y-4">
                <div className="flex items-center justify-between">
                    <div className="flex flex-col">
                        <span className="text-sm font-bold text-slate-700">Password protection</span>
                        <span className="text-xs text-slate-500">Require a code to view the file</span>
                    </div>
                    <label className="relative inline-flex items-center cursor-pointer">
                        <input 
                            type="checkbox" 
                            className="sr-only peer" 
                            checked={hasPassword}
                            onChange={(e) => setHasPassword(e.target.checked)}
                        />
                        <div className="w-11 h-6 bg-slate-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-indigo-600"></div>
                    </label>
                </div>

                {hasPassword && (
                    <div className="animate-in slide-in-from-top-2 duration-200">
                        <input 
                            type="password"
                            placeholder="Enter sharing password..."
                            className="w-full bg-white border border-slate-200 rounded-xl py-3 px-4 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 transition-all"
                        />
                    </div>
                )}

                <div className="flex items-center justify-between">
                    <div className="flex flex-col">
                        <span className="text-sm font-bold text-slate-700">Self-destruct link</span>
                        <span className="text-xs text-slate-500">Deactivate link after a specific time</span>
                    </div>
                    <input 
                        type="date"
                        value={expiry}
                        onChange={(e) => setExpiry(e.target.value)}
                        className="text-xs bg-slate-50 border border-slate-200 rounded-lg px-2 py-1.5 font-medium text-slate-600 focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
                    />
                </div>
            </div>

            <div className="pt-2">
                <button 
                    onClick={() => setShowSettings(false)}
                    className="w-full py-3 rounded-xl bg-slate-800 hover:bg-slate-900 text-white font-bold text-sm transition-all shadow-lg shadow-slate-200"
                >
                    Update & Save Settings
                </button>
            </div>
          </div>
        </div>

        {/* Footer */}
        <div className="bg-slate-50/50 px-8 py-5 border-t border-slate-100 flex items-center justify-between">
          <div className="flex items-center gap-6">
            <div className="flex flex-col">
                <span className="text-[10px] font-bold text-slate-400 uppercase tracking-widest">Views</span>
                <span className="text-sm font-bold text-slate-700">124</span>
            </div>
            <div className="flex flex-col">
                <span className="text-[10px] font-bold text-slate-400 uppercase tracking-widest">Last Activity</span>
                <span className="text-sm font-bold text-slate-700">2 min ago</span>
            </div>
          </div>
          <div className="flex items-center gap-3">
            <button className="text-sm font-semibold text-slate-500 hover:text-slate-800 transition-colors">
                Delete Link
            </button>
            <div className="w-[1px] h-4 bg-slate-200"></div>
            <button className="flex items-center gap-1.5 text-sm font-bold text-indigo-600 hover:text-indigo-800 transition-colors">
                Preview <i className="fa-solid fa-arrow-up-right-from-square text-xs"></i>
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ShareModal;
