import React, { useState, useRef, useEffect } from 'react';
import { Link, useLocation, useNavigate } from 'react-router-dom';
import { FileText, Grid, Menu, X, Star, Globe, LogOut, Settings } from 'lucide-react';
import { Button } from './ui/Button';
import { Avatar } from './ui/Avatar';
import { AppRoute } from '../types';
import { useLanguage } from '../contexts/LanguageContext';
import { useAuth } from '../contexts/AuthContext';
import { getAuthConfig } from '../services/configService';

export const Navbar: React.FC = () => {
  const [isOpen, setIsOpen] = useState(false);
  const [isProfileOpen, setIsProfileOpen] = useState(false);
  const location = useLocation();
  const navigate = useNavigate();
  const { t, language, setLanguage } = useLanguage();
  const { user, isAuthenticated, logout, isAdmin } = useAuth();
  const profileMenuRef = useRef<HTMLDivElement>(null);
  const hoverCloseTimerRef = useRef<number | null>(null);
  const [showPricing, setShowPricing] = useState(false);

  const toggleLanguage = () => {
    setLanguage(language === 'en' ? 'zh' : 'en');
  };

  const handleLogout = () => {
      logout();
      navigate(AppRoute.Home);
      setIsProfileOpen(false);
  };

  // Close profile menu when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (profileMenuRef.current && !profileMenuRef.current.contains(event.target as Node)) {
        setIsProfileOpen(false);
      }
    };
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  useEffect(() => {
    getAuthConfig().then(cfg => setShowPricing(!!cfg.enablePricingPage));
  }, []);

  // Standard nav items
  const navItems = [
    { label: t('nav.templates'), href: AppRoute.Templates, icon: <Grid size={18} /> },
    ...(showPricing ? [{ label: t('nav.pricing'), href: AppRoute.Pricing, icon: <Star size={18} /> }] : []),
  ];

  // Dashboard 入口仅保留在头像下拉菜单，不在顶部主导航显示

  return (
    <nav className="bg-white border-b border-gray-200 sticky top-0 z-50">
      <div className="max-w-[1440px] mx-auto px-6 lg:px-10">
        <div className="flex justify-between h-16">
          <div className="flex">
            <Link to="/" className="flex-shrink-0 flex items-center">
              <span className="text-2xl font-bold bg-gradient-to-r from-blue-600 to-indigo-600 text-transparent bg-clip-text">CVForge</span>
            </Link>
            <div className="hidden sm:ml-8 sm:flex sm:space-x-8">
              {navItems.map((item) => (
                <Link
                  key={item.href}
                  to={item.href}
                  className={`inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium ${
                    location.pathname === item.href
                      ? 'border-blue-500 text-gray-900'
                      : 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700'
                  }`}
                >
                  {item.label}
                </Link>
              ))}
            </div>
          </div>
          <div className="hidden sm:ml-6 sm:flex sm:items-center sm:space-x-4 sm:-mr-2 md:-mr-4">
            <button 
               onClick={toggleLanguage}
               className="p-2 text-gray-500 hover:text-gray-900 focus:outline-none"
               title={t('lang.switchTitle')}
           >
                <div className="flex items-center space-x-1">
                    <Globe size={18} />
                    <span className="text-sm font-medium">{language === 'en' ? t('lang.en_short') : t('lang.zh_short')}</span>
                </div>
            </button>

            {isAuthenticated && user ? (
                <div 
                    className="relative ml-3" 
                    ref={profileMenuRef}
                    onMouseEnter={() => {
                      if (hoverCloseTimerRef.current) { clearTimeout(hoverCloseTimerRef.current); hoverCloseTimerRef.current = null; }
                      setIsProfileOpen(true);
                    }}
                    onMouseLeave={() => {
                      if (hoverCloseTimerRef.current) { clearTimeout(hoverCloseTimerRef.current); }
                      hoverCloseTimerRef.current = window.setTimeout(() => setIsProfileOpen(false), 120);
                    }}
                    onTouchStart={() => setIsProfileOpen(true)}
                >
                    <div>
                        <button 
                            className="flex items-center max-w-xs text-sm rounded-full focus:outline-none" 
                            id="user-menu-button"
                        >
                            <span className="sr-only">{t('a11y.openUserMenu')}</span>
                            <Avatar className="h-8 w-8 rounded-full text-sm" src={user.avatarUrl} name={user.name} />
                        </button>
                    </div>
                    
                    {isProfileOpen && (
                        <div 
                            className="origin-top absolute left-1/2 -translate-x-1/2 mt-2 w-auto min-w-[122px] rounded-md shadow-lg py-1 bg-white ring-1 ring-black ring-opacity-5 focus:outline-none z-50 animate-fadeIn" 
                            role="menu"
                        >
                            
                            <Link to={AppRoute.Dashboard} target="_blank" rel="noopener noreferrer" onClick={() => setIsProfileOpen(false)} className="group flex items-center px-4 py-2 text-sm text-gray-700 hover:bg-gray-50" role="menuitem">
                                <FileText size={16} className="mr-3 text-gray-400 group-hover:text-gray-500"/>
                                {t('nav.dashboard')}
                            </Link>
                            
                            <Link to={AppRoute.Settings} onClick={() => setIsProfileOpen(false)} className="group flex items-center px-4 py-2 text-sm text-gray-700 hover:bg-gray-50" role="menuitem">
                                <Settings size={16} className="mr-3 text-gray-400 group-hover:text-gray-500"/>
                                {t('nav.settings')}
                            </Link>
                            
                            {isAdmin && (
                              <Link to={AppRoute.Admin} target="_blank" rel="noopener noreferrer" onClick={() => setIsProfileOpen(false)} className="group flex items-center px-4 py-2 text-sm text-gray-700 hover:bg-gray-50" role="menuitem">
                                  <Grid size={16} className="mr-3 text-gray-400 group-hover:text-gray-500"/>
                                  {t('nav.admin')}
                              </Link>
                            )}

                             <button 
                                onClick={handleLogout}
                                className="w-full group flex items-center px-4 py-2 text-sm text-red-600 hover:bg-red-50 border-t border-gray-100 mt-1" 
                                role="menuitem"
                            >
                                <LogOut size={16} className="mr-3 text-red-400 group-hover:text-red-500"/>
                                {t('nav.signout')}
                            </button>
                        </div>
                    )}
                </div>
            ) : (
                <>
                    <Link to={AppRoute.Login}>
                        <Button variant="ghost">{t('nav.signin')}</Button>
                    </Link>
                    <Link to={AppRoute.Templates}>
                        <Button>{t('nav.getStarted')}</Button>
                    </Link>
                </>
            )}
          </div>
          
          <div className="-mr-2 flex items-center sm:hidden">
            <button
              onClick={() => setIsOpen(!isOpen)}
              className="inline-flex items-center justify-center p-2 rounded-md text-gray-400 hover:text-gray-500 hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-blue-500"
            >
              {isOpen ? <X size={24} /> : <Menu size={24} />}
            </button>
          </div>
        </div>
      </div>

      {/* Mobile menu */}
      {isOpen && (
        <div className="sm:hidden bg-white border-b border-gray-200">
          <div className="pt-2 pb-3 space-y-1">
            {navItems.map((item) => (
              <Link
                key={item.href}
                to={item.href}
                onClick={() => setIsOpen(false)}
                className={`block pl-3 pr-4 py-2 border-l-4 text-base font-medium ${
                   location.pathname === item.href
                    ? 'bg-blue-50 border-blue-500 text-blue-700'
                    : 'border-transparent text-gray-500 hover:bg-gray-50 hover:border-gray-300 hover:text-gray-700'
                }`}
              >
                <div className="flex items-center">
                    <span className="mr-3">{item.icon}</span>
                    {item.label}
                </div>
              </Link>
            ))}
             <div className="pt-4 pb-3 border-t border-gray-200">
                <div className="mt-3 space-y-1 px-2">
                     <button onClick={() => { toggleLanguage(); setIsOpen(false); }} className="w-full text-left pl-3 pr-4 py-2 border-l-4 border-transparent text-base font-medium text-gray-600 hover:bg-gray-50">
                        {language === 'en' ? t('lang.switchToZh') : t('lang.switchToEn')}
                     </button>
                     
                     {isAuthenticated ? (
                         <>
                            <div className="px-4 py-2 flex items-center">
                                <Avatar className="h-8 w-8 rounded-full mr-3 text-sm" src={user?.avatarUrl} name={user?.name} />
                                <div>
                                    <div className="text-base font-medium text-gray-800">{user?.name}</div>
                                    <div className="text-sm font-medium text-gray-500">{user?.email}</div>
                                </div>
                            </div>
                             
                            <Button 
                                className="w-full justify-start mt-2" 
                                variant="ghost" 
                                icon={<FileText size={16}/>}
                                onClick={() => { window.open(`${window.location.origin}${window.location.pathname}#${AppRoute.Dashboard}`, '_blank'); setIsOpen(false); }}
                            >
                                {t('nav.dashboard')}
                            </Button>
                            
                            <Button 
                                className="w-full justify-start mt-2" 
                                variant="ghost" 
                                icon={<Settings size={16}/>}
                                onClick={() => { navigate(AppRoute.Settings); setIsOpen(false); }}
                            >
                                {t('nav.settings')}
                            </Button>
                            
                            {isAdmin && (
                              <Button 
                                  className="w-full justify-start mt-2" 
                                  variant="ghost" 
                                  icon={<Grid size={16}/>}
                                  onClick={() => { window.open(`${window.location.origin}${window.location.pathname}#${AppRoute.Admin}`, '_blank'); setIsOpen(false); }}
                              >
                                  {t('nav.admin')}
                              </Button>
                            )}

                             <Button 
                                className="w-full justify-start mt-2 text-red-600 hover:text-red-700 hover:bg-red-50" 
                                variant="ghost" 
                                icon={<LogOut size={16}/>} 
                                onClick={handleLogout}
                            >
                                {t('nav.signout')}
                            </Button>
                         </>
                     ) : (
                         <>
                             <Link to={AppRoute.Login} onClick={() => setIsOpen(false)}>
                                <Button className="w-full justify-center" variant="outline">{t('nav.signin')}</Button>
                             </Link>
                             <div className="h-2"></div>
                             <Link to={AppRoute.Templates} onClick={() => setIsOpen(false)}>
                                <Button className="w-full justify-center">{t('nav.getStarted')}</Button>
                             </Link>
                         </>
                     )}
                </div>
            </div>
          </div>
        </div>
      )}
    </nav>
  );
};

export const Footer: React.FC = () => {
    const { t, setLanguage } = useLanguage();
    const [showPricing, setShowPricing] = useState(false);
    useEffect(() => { getAuthConfig().then(cfg => setShowPricing(!!cfg.enablePricingPage)); }, []);

    return (
        <footer className="bg-slate-900 text-white py-6">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 grid grid-cols-1 md:grid-cols-4 gap-6">
                <div>
                    <h3 className="text-xl font-bold mb-4">CVForge</h3>
                    <p className="text-slate-400 text-sm">{t('hero.desc')}</p>
                </div>
                <div>
                    <h4 className="font-semibold mb-4">{t('footer.product')}</h4>
                    <ul className="space-y-2 text-slate-400 text-sm">
                        <li><Link to="/templates" className="hover:text-white">{t('nav.templates')}</Link></li>
                        {showPricing && <li><Link to="/pricing" className="hover:text-white">{t('nav.pricing')}</Link></li>}
                    </ul>
                </div>
                <div>
                    <h4 className="font-semibold mb-4">{t('footer.legal')}</h4>
                    <ul className="space-y-2 text-slate-400 text-sm">
                        <li>{t('footer.privacy')}</li>
                        <li>{t('footer.terms')}</li>
                    </ul>
                </div>
                <div>
                    <h4 className="font-semibold mb-4">{t('footer.language')}</h4>
                    <div className="flex space-x-2">
                        <button onClick={() => setLanguage('en')} className="text-sm text-slate-400 hover:text-white cursor-pointer">{t('lang.en')}</button>
                        <span className="text-sm text-slate-400">|</span>
                        <button onClick={() => setLanguage('zh')} className="text-sm text-slate-400 hover:text-white cursor-pointer">{t('lang.zh')}</button>
                    </div>
                </div>
            </div>
            <div className="max-w-7xl mx-auto px-4 mt-4 pt-4 border-t border-slate-800 text-center text-slate-500 text-sm">
                <div>
                  &copy; {new Date().getFullYear()} CVForge. {t('footer.rights')}
                </div>
                <a
                  href="https://beian.miit.gov.cn/"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="mt-1 inline-block hover:text-slate-300"
                >
                  赣ICP备2024029737号-2
                </a>
            </div>
        </footer>
    )
}

export const MainLayout: React.FC<{ children: React.ReactNode; hideFooter?: boolean }> = ({ children, hideFooter }) => {
  return (
    <div className="min-h-screen flex flex-col bg-gray-50">
      <Navbar />
      <main className="flex-grow">
        {children}
      </main>
      {!hideFooter && <Footer />}
    </div>
  );
};
