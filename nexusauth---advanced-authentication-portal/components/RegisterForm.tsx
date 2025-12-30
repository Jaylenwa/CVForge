
import React, { useState, useEffect } from 'react';
import { Mail, ShieldCheck, Lock, Loader2 } from 'lucide-react';
import { SocialButtons } from './SocialButtons';

export const RegisterForm: React.FC = () => {
  const [isLoading, setIsLoading] = useState(false);
  const [isCounting, setIsCounting] = useState(false);
  const [timer, setTimer] = useState(60);
  const [formData, setFormData] = useState({ 
    email: '', 
    code: '', 
    password: '', 
    confirmPassword: '' 
  });

  useEffect(() => {
    let interval: any;
    if (isCounting && timer > 0) {
      interval = setInterval(() => setTimer(t => t - 1), 1000);
    } else if (timer === 0) {
      setIsCounting(false);
      setTimer(60);
    }
    return () => clearInterval(interval);
  }, [isCounting, timer]);

  const handleSendCode = () => {
    if (!formData.email) return;
    setIsCounting(true);
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setTimeout(() => setIsLoading(false), 2000);
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-5">
      <div className="space-y-2">
        <label className="text-xs font-medium text-slate-400 uppercase tracking-[0.1em] ml-1">Email</label>
        <div className="relative group">
          <Mail className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-slate-500 group-focus-within:text-indigo-400 transition-colors" />
          <input 
            type="email" 
            required
            className="w-full bg-white/5 border border-white/10 rounded-2xl py-3.5 pl-12 pr-4 text-white focus:outline-none focus:ring-2 focus:ring-indigo-500/50 transition-all"
            placeholder="your@email.com"
            value={formData.email}
            onChange={(e) => setFormData({ ...formData, email: e.target.value })}
          />
        </div>
      </div>

      <div className="space-y-2">
        <label className="text-xs font-medium text-slate-400 uppercase tracking-[0.1em] ml-1">Verification Code</label>
        <div className="flex gap-3">
          <div className="relative flex-1 group">
            <ShieldCheck className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-slate-500 group-focus-within:text-indigo-400 transition-colors" />
            <input 
              type="text" 
              required
              maxLength={6}
              className="w-full bg-white/5 border border-white/10 rounded-2xl py-3.5 pl-12 pr-4 text-white focus:outline-none focus:ring-2 focus:ring-indigo-500/50 transition-all tracking-[0.3em]"
              placeholder="000000"
              value={formData.code}
              onChange={(e) => setFormData({ ...formData, code: e.target.value })}
            />
          </div>
          <button 
            type="button"
            onClick={handleSendCode}
            disabled={isCounting || !formData.email}
            className="px-5 bg-white/5 hover:bg-white/10 disabled:opacity-30 text-indigo-400 text-sm font-semibold rounded-2xl border border-white/10 transition-all flex items-center justify-center min-w-[110px]"
          >
            {isCounting ? `${timer}s` : 'Get Code'}
          </button>
        </div>
      </div>

      <div className="grid grid-cols-2 gap-4">
        <div className="space-y-2">
          <label className="text-xs font-medium text-slate-400 uppercase tracking-[0.1em] ml-1">Password</label>
          <div className="relative group">
            <Lock className="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-500 group-focus-within:text-indigo-400 transition-colors" />
            <input 
              type="password" 
              required
              className="w-full bg-white/5 border border-white/10 rounded-2xl py-3.5 pl-11 pr-4 text-white text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500/50 transition-all"
              placeholder="8+ chars"
              value={formData.password}
              onChange={(e) => setFormData({ ...formData, password: e.target.value })}
            />
          </div>
        </div>
        <div className="space-y-2">
          <label className="text-xs font-medium text-slate-400 uppercase tracking-[0.1em] ml-1">Repeat</label>
          <div className="relative group">
            <Lock className="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-500 group-focus-within:text-indigo-400 transition-colors" />
            <input 
              type="password" 
              required
              className="w-full bg-white/5 border border-white/10 rounded-2xl py-3.5 pl-11 pr-4 text-white text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500/50 transition-all"
              placeholder="Confirm"
              value={formData.confirmPassword}
              onChange={(e) => setFormData({ ...formData, confirmPassword: e.target.value })}
            />
          </div>
        </div>
      </div>

      <button 
        type="submit" 
        disabled={isLoading}
        className="w-full bg-indigo-600 hover:bg-indigo-500 disabled:opacity-70 text-white font-bold py-4 rounded-2xl shadow-[0_10px_20px_rgba(79,70,229,0.3)] transition-all flex items-center justify-center space-x-2 active:scale-[0.97] mt-4"
      >
        {isLoading ? <Loader2 className="w-5 h-5 animate-spin" /> : <span>Start Journey</span>}
      </button>

      <div className="relative py-4">
        <div className="absolute inset-0 flex items-center"><div className="w-full border-t border-white/5"></div></div>
        <div className="relative flex justify-center text-[10px] uppercase font-bold tracking-widest text-slate-600"><span className="bg-transparent px-3">Join Via</span></div>
      </div>

      <SocialButtons />
    </form>
  );
};
