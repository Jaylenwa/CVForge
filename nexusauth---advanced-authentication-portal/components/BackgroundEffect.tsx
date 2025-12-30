
import React from 'react';

export const BackgroundEffect: React.FC = () => {
  return (
    <div className="fixed inset-0 -z-0 bg-[#020617]">
      {/* Primary Orbs */}
      <div className="absolute top-[-10%] left-[-10%] w-[60%] h-[60%] bg-indigo-600/15 blur-[160px] rounded-full animate-pulse"></div>
      <div className="absolute bottom-[-10%] right-[-10%] w-[70%] h-[70%] bg-blue-600/15 blur-[160px] rounded-full animate-pulse" style={{ animationDelay: '3s' }}></div>
      <div className="absolute top-[20%] right-[10%] w-[30%] h-[30%] bg-violet-600/10 blur-[140px] rounded-full animate-pulse" style={{ animationDelay: '1.5s' }}></div>
      
      {/* Mesh/Grid layer */}
      <div className="absolute inset-0 opacity-[0.03] pointer-events-none" style={{ backgroundImage: `radial-gradient(circle at 2px 2px, white 1px, transparent 0)`, backgroundSize: '40px 40px' }}></div>
      
      {/* Additional Noise/Texture */}
      <div className="absolute inset-0 opacity-[0.02] pointer-events-none mix-blend-overlay" style={{ backgroundImage: `url('https://grainy-gradients.vercel.app/noise.svg')` }}></div>
    </div>
  );
};
