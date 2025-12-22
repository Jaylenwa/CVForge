import React from 'react';
import { Link, NavLink, Outlet } from 'react-router-dom';
import { MainLayout } from '../../components/Layout';
import { AppRoute } from '../../types';
import { useLanguage } from '../../contexts/LanguageContext';

export const AdminLayout: React.FC = () => {
  const { t } = useLanguage();
  const tabs = [
    { to: AppRoute.AdminUsers, label: t('admin.users') },
    { to: AppRoute.AdminResumes, label: t('admin.resumes') },
    { to: AppRoute.AdminTemplates, label: t('admin.templates') },
    { to: AppRoute.AdminShares, label: t('admin.shares') },
    { to: AppRoute.AdminConfigs, label: 'Settings' },
  ];
  return (
    <MainLayout hideFooter>
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="flex items-center justify-between mb-6">
          <h1 className="text-2xl font-bold text-gray-900">{t('admin.title')}</h1>
          <Link to={AppRoute.Dashboard} className="text-sm text-blue-600 hover:text-blue-700">{t('admin.backDashboard')}</Link>
        </div>
        <div className="border-b border-gray-200">
          <nav className="-mb-px flex space-x-8">
            {tabs.map(tab => (
              <NavLink
                key={tab.to}
                to={tab.to}
                className={({ isActive }) =>
                  `whitespace-nowrap pb-2 px-1 border-b-2 text-sm font-medium ${
                    isActive ? 'border-blue-500 text-blue-600' : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                  }`
                }
              >
                {tab.label}
              </NavLink>
            ))}
          </nav>
        </div>
        <div className="mt-6">
          <Outlet />
        </div>
      </div>
    </MainLayout>
  );
};
