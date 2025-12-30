
import React, { useState } from 'react';
import { AuthCard } from './components/AuthCard';
import { BackgroundEffect } from './components/BackgroundEffect';
import { GeminiAssistant } from './components/GeminiAssistant';

const App: React.FC = () => {
  const [isLogin, setIsLogin] = useState(true);

  const toggleAuthMode = () => {
    setIsLogin(prev => !prev);
  };

  return (
    <div className="relative min-h-screen flex flex-col items-center justify-center overflow-hidden p-4">
      {/* Animated Abstract Background */}
      <BackgroundEffect />

      {/* Main Auth Container */}
      <div className="relative z-10 w-full max-w-md">
        <AuthCard isLogin={isLogin} toggleAuthMode={toggleAuthMode} />
      </div>

      {/* Intelligent Assistant Bubble */}
      <GeminiAssistant />

      {/* Footer Decoration */}
      <div className="mt-8 text-slate-400 text-sm z-10 font-light tracking-wide animate-pulse">
        © 2024 NexusAuth Secure Portal
      </div>
    </div>
  );
};

export default App;
