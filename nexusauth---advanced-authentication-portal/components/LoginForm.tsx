
import React, { useState } from 'react';
import { Mail, Lock, Loader2 } from 'lucide-react';
import { SocialButtons } from './SocialButtons';

export const LoginForm: React.FC = () => {
  const [isLoading, setIsLoading] = useState(false);
  const [formData, setFormData] = useState({ email: '', password: '' });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setTimeout(() => setIsLoading(false), 2000);
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div className="space-y-2">
        <label className="text-xs font-medium text-slate-400 uppercase tracking-[0.1em] ml-1">Email Address</label>
        <div className="relative group">
          <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none transition-colors group-focus-within:text-indigo-400 text-slate-500">
            <Mail className="w-5 h-5" />
          </div>
          <input 
            type="email" 
            required
            className="w-full bg-white/5 border border-white/10 rounded-2xl py-4 pl-12 pr-4 text-white placeholder:text-slate-600 focus:outline-none focus:ring-2 focus:ring-indigo-500/50 focus:border-indigo-500/50 transition-all duration-300"
            placeholder="example@email.com"
            value={formData.email}
            onChange={(e) => setFormData({ ...formData, email: e.target.value })}
          />
        </div>
      </div>

      <div className="space-y-2">
        <div className="flex justify-between items-center ml-1">
          <label className="text-xs font-medium text-slate-400 uppercase tracking-[0.1em]">Password</label>
          <button type="button" className="text-xs text-indigo-400 hover:text-indigo-300 transition-colors">Forgot Password?</button>
        </div>
        <div className="relative group">
          <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none transition-colors group-focus-within:text-indigo-400 text-slate-500">
            <Lock className="w-5 h-5" />
          </div>
          <input 
            type="password" 
            required
            className="w-full bg-white/5 border border-white/10 rounded-2xl py-4 pl-12 pr-4 text-white placeholder:text-slate-600 focus:outline-none focus:ring-2 focus:ring-indigo-500/50 focus:border-indigo-500/50 transition-all duration-300"
            placeholder="••••••••"
            value={formData.password}
            onChange={(e) => setFormData({ ...formData, password: e.target.value })}
          />
        </div>
      </div>

      <button 
        type="submit" 
        disabled={isLoading}
        className="w-full bg-indigo-600 hover:bg-indigo-500 disabled:opacity-70 text-white font-semibold py-4 rounded-2xl shadow-[0_10px_20px_rgba(79,70,229,0.3)] transition-all duration-300 flex items-center justify-center space-x-2 active:scale-[0.97] mt-2"
      >
        {isLoading ? <Loader2 className="w-5 h-5 animate-spin" /> : <span>Sign In</span>}
      </button>

      <div className="relative py-4">
        <div className="absolute inset-0 flex items-center">
          <div className="w-full border-t border-white/5"></div>
        </div>
        <div className="relative flex justify-center text-xs">
          <span className="bg-[#0f172a]/0 px-4 text-slate-500 font-medium">Alternative access</span>
        </div>
      </div>

      <SocialButtons />
    </form>
  );
};
