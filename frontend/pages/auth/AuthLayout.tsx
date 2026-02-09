import React from 'react';
import { Link } from 'react-router-dom';
import { ArrowLeft, Globe } from 'lucide-react';
import { useLanguage } from '../../contexts/LanguageContext';

interface AuthLayoutProps {
  children: React.ReactNode;
  image?: string;
  quote?: string;
  author?: string;
}

export const AuthLayout: React.FC<AuthLayoutProps> = ({ 
  children
}) => {
  const { t, language, setLanguage } = useLanguage();
  return (
    <div className="min-h-screen bg-white flex items-center justify-center p-8 sm:p-12 lg:p-16 relative">
      <div className="absolute top-6 left-6">
        <Link to="/" className="flex items-center text-gray-500 hover:text-gray-900 transition-colors">
          <ArrowLeft size={20} className="mr-2" /> {t('common.backHome')}
        </Link>
      </div>
      <div className="absolute top-6 right-6">
        <button
          onClick={() => setLanguage(language === 'en' ? 'zh' : 'en')}
          className="p-2 text-gray-500 hover:text-gray-900 focus:outline-none"
          title={t('lang.switchTitle')}
        >
          <div className="flex items-center space-x-1">
            <Globe size={18} />
            <span className="text-sm font-medium">
              {language === 'en' ? t('lang.en_short') : t('lang.zh_short')}
            </span>
          </div>
        </button>
      </div>
  <div className="w-full max-w-md space-y-8">
        <div className="bg-white border border-gray-200 rounded-2xl shadow-sm p-6 sm:p-8">
          {children}
        </div>
      </div>
    </div>
  );
};
