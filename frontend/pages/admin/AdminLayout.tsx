import React from 'react';
import { NavLink, Outlet, useNavigate } from 'react-router-dom';
import { AppRoute } from '../../types';
import { useLanguage } from '../../contexts/LanguageContext';
import { BarChart3, Settings, Users, FileText, LayoutGrid, Share, ChevronRight, Globe } from 'lucide-react';

export const AdminLayout: React.FC = () => {
  const { t, language, setLanguage } = useLanguage();
  const navigate = useNavigate();

  const menu = [
    { to: AppRoute.Admin, label: t('admin.menu.dashboard'), icon: <BarChart3 size={18} /> },
    { to: AppRoute.AdminUsers, label: t('admin.menu.users'), icon: <Users size={18} /> },
    { to: AppRoute.AdminResumes, label: t('admin.menu.resumes'), icon: <FileText size={18} /> },
    { to: AppRoute.AdminTemplates, label: t('admin.menu.templates'), icon: <LayoutGrid size={18} /> },
    { to: AppRoute.AdminShares, label: t('admin.menu.shares'), icon: <Share size={18} /> },
    { to: AppRoute.AdminConfigs, label: t('admin.menu.settings'), icon: <Settings size={18} /> },
  ];

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="flex h-screen">
        {/* Sidebar */}
        <aside className="w-60 bg-white border-r border-gray-200 p-3">
          <div className="px-2 py-3 flex items-center justify-between">
            <div>
              <div className="text-xl font-bold text-gray-900">Admin</div>
              <div className="text-xs text-gray-500">{t('admin.title')}</div>
            </div>
            <button
              onClick={() => setLanguage(language === 'en' ? 'zh' : 'en')}
              className="p-2 text-gray-500 hover:text-gray-900 focus:outline-none rounded-md"
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
          <nav className="mt-2 space-y-1">
            {menu.map(item => (
              <NavLink
                key={item.to}
                to={item.to}
                className={({ isActive }) =>
                  `flex items-center justify-between px-3 py-2 rounded-md text-sm ${
                    isActive ? 'bg-blue-50 text-blue-700' : 'text-gray-700 hover:bg-gray-100'
                  }`
                }
              >
                <span className="flex items-center space-x-2">
                  <span className="text-gray-500">{item.icon}</span>
                  <span>{item.label}</span>
                </span>
                <ChevronRight size={16} className="text-gray-400" />
              </NavLink>
            ))}
          </nav>
          <div className="mt-6 px-2">
            <button
              onClick={() => navigate(AppRoute.Dashboard)}
              className="w-full px-3 py-2 text-sm text-blue-600 hover:bg-blue-50 rounded-md"
            >
              {t('admin.backDashboard')}
            </button>
          </div>
        </aside>

        {/* Content */}
        <main className="flex-1 overflow-y-auto">
          <div className="p-6">
            <Outlet />
          </div>
        </main>
      </div>
    </div>
  );
};
