import React from 'react';
import { NavLink, Outlet, useNavigate } from 'react-router-dom';
import { AppRoute } from '../../types';
import { useLanguage } from '../../contexts/LanguageContext';
import { BarChart3, Settings, Users, FileText, LayoutGrid, Share, ChevronRight, Globe, Home } from 'lucide-react';

export const AdminLayout: React.FC = () => {
  const { t, language, setLanguage } = useLanguage();
  const navigate = useNavigate();

  const menu = [
    { to: AppRoute.Admin, label: t('admin.menu.dashboard'), icon: <BarChart3 size={18} />, end: true },
    { to: AppRoute.AdminUsers, label: t('admin.menu.users'), icon: <Users size={18} /> },
    { to: AppRoute.AdminResumes, label: t('admin.menu.resumes'), icon: <FileText size={18} /> },
    { to: AppRoute.AdminTemplates, label: t('admin.menu.templates'), icon: <LayoutGrid size={18} /> },
    { to: AppRoute.AdminShares, label: t('admin.menu.shares'), icon: <Share size={18} /> },
    { to: AppRoute.AdminCatalog, label: t('admin.menu.catalog'), icon: <LayoutGrid size={18} /> },
    { to: AppRoute.AdminConfigs, label: t('admin.menu.settings'), icon: <Settings size={18} /> },
  ];

  return (
    <div className="flex h-screen overflow-hidden bg-[#f4f7fa]">
      <aside className="w-64 bg-transparent h-screen flex flex-col py-6 px-4 overflow-y-auto">
        <div className="px-2 py-3 flex items-center justify-between">
          <div>
            <div className="text-xl font-bold text-gray-900">{t('nav.admin')}</div>
            <div className="text-xs text-gray-500">{t('admin.title')}</div>
          </div>
          <button
            onClick={() => setLanguage(language === 'en' ? 'zh' : 'en')}
            className="p-2 text-gray-500 hover:text-gray-900 rounded-md"
            title={t('lang.switchTitle')}
          >
            <div className="flex items-center space-x-1">
              <Globe size={18} />
              <span className="text-xs font-medium">
                {language === 'en' ? t('lang.en_short') : t('lang.zh_short')}
              </span>
            </div>
          </button>
        </div>
        <nav className="space-y-1">
          {menu.map(item => (
            <NavLink
              key={item.to}
              to={item.to}
              end={Boolean((item as any).end)}
              className={({ isActive }) =>
                `w-full flex items-center justify-between px-4 py-3 rounded-lg transition-all duration-200 group ${
                  isActive
                    ? 'bg-blue-100 text-blue-600 font-semibold'
                    : 'text-gray-500 hover:bg-gray-100 hover:text-gray-900'
                }`
              }
            >
              {({ isActive }) => (
                <>
                  <div className="flex items-center space-x-3">
                    <span className={`${isActive ? 'text-blue-600' : 'text-gray-400 group-hover:text-gray-600'}`}>
                      {item.icon}
                    </span>
                    <span className="text-[15px]">{item.label}</span>
                  </div>
                  <ChevronRight size={16} className="text-gray-400" />
                </>
              )}
            </NavLink>
          ))}
        </nav>
        <div className="mt-6">
          <NavLink
            to={AppRoute.Dashboard}
            className="w-full flex items-center justify-between px-4 py-3 rounded-lg transition-all duration-200 group text-gray-500 hover:bg-gray-100 hover:text-gray-900"
          >
            <div className="flex items-center space-x-3">
              <span className="text-gray-400 group-hover:text-gray-600">
                <Home size={18} />
              </span>
              <span className="text-[15px]">{t('admin.backDashboard')}</span>
            </div>
            <ChevronRight size={16} className="text-gray-400" />
          </NavLink>
        </div>
      </aside>
      <main className="flex-1 flex flex-col h-screen p-8 overflow-y-auto">
        <Outlet />
      </main>
    </div>
  );
};
