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
  children, 
  image = "https://images.unsplash.com/photo-1586281380349-632531db7ed4?ixlib=rb-1.2.1&auto=format&fit=crop&w=1950&q=80",
  quote = "The future belongs to those who believe in the beauty of their dreams.",
  author = "Eleanor Roosevelt"
}) => {
  const { t, language, setLanguage } = useLanguage();
  return (
    <div className="min-h-screen flex bg-white">
      {/* Left Column - Image (Hidden on mobile) */}
      <div className="hidden lg:flex lg:w-1/2 relative bg-gray-900">
        <img 
          src={image} 
          alt={t('a11y.authBackgroundAlt')} 
          className="absolute inset-0 w-full h-full object-cover opacity-60"
        />
        <div className="relative z-10 w-full flex flex-col justify-between p-12 text-white">
            <div>
                 <Link to="/" className="flex items-center text-white/80 hover:text-white transition-colors">
                    <ArrowLeft size={20} className="mr-2" /> {t('common.backHome')}
                 </Link>
            </div>
            <div className="mb-10">
                <blockquote className="text-2xl font-light italic mb-4">
                    "{quote}"
                </blockquote>
                <p className="font-semibold text-lg">— {author}</p>
            </div>
        </div>
      </div>

      {/* Right Column - Form */}
      <div className="w-full lg:w-1/2 flex items-center justify-center p-8 sm:p-12 lg:p-16 overflow-y-auto relative">
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
            <div className="lg:hidden mb-8">
                 <Link to="/" className="flex items-center text-gray-500 hover:text-gray-900">
                    <ArrowLeft size={18} className="mr-2" /> {t('common.backHome')}
                 </Link>
            </div>
            {children}
        </div>
      </div>
    </div>
  );
};
