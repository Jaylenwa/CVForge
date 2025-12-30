
import React from 'react';

export const SocialButtons: React.FC = () => {
  return (
    <div className="grid grid-cols-2 gap-4">
      <button 
        type="button"
        className="flex items-center justify-center space-x-3 bg-white/5 hover:bg-white/10 border border-white/10 rounded-2xl py-3.5 transition-all duration-300 group active:scale-95"
      >
        <GitHubIcon className="w-6 h-6 text-white transition-transform duration-300 group-hover:scale-110" />
        <span className="text-sm font-semibold text-slate-300">GitHub</span>
      </button>
      <button 
        type="button"
        className="flex items-center justify-center space-x-3 bg-white/5 hover:bg-white/10 border border-white/10 rounded-2xl py-3.5 transition-all duration-300 group active:scale-95"
      >
        <WeChatIcon className="w-6 h-6 text-emerald-500 transition-transform duration-300 group-hover:scale-110" />
        <span className="text-sm font-semibold text-slate-300">WeChat</span>
      </button>
    </div>
  );
};

const GitHubIcon = ({ className }: { className?: string }) => (
  <svg viewBox="0 0 24 24" fill="currentColor" className={className} xmlns="http://www.w3.org/2000/svg"><path d="M12 .297c-6.63 0-12 5.373-12 12 0 5.303 3.438 9.8 8.205 11.385.6.113.82-.258.82-.577 0-.285-.01-1.04-.015-2.04-3.338.724-4.042-1.61-4.042-1.61C4.422 18.07 3.633 17.7 3.633 17.7c-1.087-.744.084-.729.084-.729 1.205.084 1.838 1.236 1.838 1.236 1.07 1.835 2.809 1.305 3.495.998.108-.776.417-1.305.76-1.605-2.665-.3-5.466-1.332-5.466-5.93 0-1.31.465-2.38 1.235-3.22-.135-.303-.54-1.523.105-3.176 0 0 1.005-.322 3.3 1.23.96-.267 1.98-.399 3-.405 1.02.006 2.04.138 3 .405 2.28-1.552 3.285-1.23 3.285-1.23.645 1.653.24 2.873.12 3.176.765.84 1.23 1.91 1.23 3.22 0 4.61-2.805 5.625-5.475 5.92.43.372.823 1.102.823 2.222 0 1.606-.015 2.896-.015 3.286 0 .315.21.69.825.57C20.565 22.092 24 17.592 24 12.297c0-6.627-5.373-12-12-12"/></svg>
);

const WeChatIcon = ({ className }: { className?: string }) => (
  <svg viewBox="0 0 24 24" fill="currentColor" className={className} xmlns="http://www.w3.org/2000/svg"><path d="M8.5 13a.75.75 0 1 1 0-1.5.75.75 0 0 1 0 1.5zm4.5 0a.75.75 0 1 1 0-1.5.75.75 0 0 1 0 1.5zm-5-9.5c-4.418 0-8 3.134-8 7s3.582 7 8 7c.78 0 1.528-.1 2.241-.286L13 19.5l-.234-2.502c3.557-.59 6.234-3.415 6.234-6.748 0-3.866-3.582-7-8-7zm12 5c-3.314 0-6 2.462-6 5.5s2.686 5.5 6 5.5c.585 0 1.146-.079 1.68-.225l2.32 1.725-.175-1.966c2.668-.463 4.675-2.683 4.675-5.304 0-3.038-2.686-5.5-6-5.5zM17 13a.5.5 0 1 1 0-1 .5.5 0 0 1 0 1zm3 0a.5.5 0 1 1 0-1 .5.5 0 0 1 0 1z"/></svg>
);
