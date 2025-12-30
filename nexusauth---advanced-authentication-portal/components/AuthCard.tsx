
import React from 'react';
import { LoginForm } from './LoginForm';
import { RegisterForm } from './RegisterForm';

interface AuthCardProps {
  isLogin: boolean;
  toggleAuthMode: () => void;
}

export const AuthCard: React.FC<AuthCardProps> = ({ isLogin, toggleAuthMode }) => {
  return (
    <div className="w-full max-w-[440px] mx-auto bg-slate-900/40 backdrop-blur-2xl border border-white/10 rounded-[2.5rem] shadow-[0_20px_50px_rgba(0,0,0,0.3)] overflow-hidden transition-all duration-700 hover:border-white/20">
      <div className="relative p-10 sm:p-12">
        {/* Header Section */}
        <div className="mb-10 text-center">
          <div className="inline-block p-3 rounded-2xl bg-indigo-500/10 border border-indigo-500/20 mb-6 animate-bounce-slow">
            <svg className="w-8 h-8 text-indigo-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 00-2 2zm10-10V7a4 4 0 00-8 0v4h8z" />
            </svg>
          </div>
          <h1 className="text-3xl font-bold text-white mb-3 tracking-tight">
            {isLogin ? 'Welcome Back' : 'Create Account'}
          </h1>
          <p className="text-slate-400 text-base font-light">
            {isLogin 
              ? 'Securely access your Nexus dashboard' 
              : 'Join thousands of creators worldwide'}
          </p>
        </div>

        {/* Form Container with smooth height and slide transition */}
        <div className="relative overflow-hidden transition-[height] duration-500 ease-in-out" style={{ minHeight: isLogin ? '380px' : '480px' }}>
           <div 
             className={`transition-all duration-700 ease-[cubic-bezier(0.23,1,0.32,1)] transform ${
               isLogin ? 'opacity-100 translate-x-0' : 'opacity-0 -translate-x-12 absolute w-full pointer-events-none'
             }`}
           >
             <LoginForm />
           </div>
           
           <div 
             className={`transition-all duration-700 ease-[cubic-bezier(0.23,1,0.32,1)] transform ${
               !isLogin ? 'opacity-100 translate-x-0' : 'opacity-0 translate-x-12 absolute w-full pointer-events-none'
             }`}
           >
             <RegisterForm />
           </div>
        </div>

        {/* Footer Toggle */}
        <div className="mt-10 text-center border-t border-white/5 pt-8">
          <button 
            onClick={toggleAuthMode}
            className="group text-slate-400 hover:text-white text-sm font-medium transition-all inline-flex items-center gap-2"
          >
            <span>{isLogin ? "New to our platform?" : "Already have an account?"}</span>
            <span className="text-indigo-400 group-hover:translate-x-1 transition-transform">
              {isLogin ? "Sign up now →" : "Log in here →"}
            </span>
          </button>
        </div>
      </div>
    </div>
  );
};

// Add custom animation to global styles via tailwind config in index.html if needed, 
// but here we use standard classes and some inline styles for simplicity.
