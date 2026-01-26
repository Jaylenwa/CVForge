import React, { Suspense } from 'react';
import { HashRouter as Router, Routes, Route, Outlet } from 'react-router-dom';
import { MainLayout } from './components/Layout';
import { AppRoute } from './types';
import { LanguageProvider } from './contexts/LanguageContext';
import { AuthProvider } from './contexts/AuthContext';
import { ToastProvider } from './components/ui/Toast';
import { ConfirmDialogProvider } from './components/ui/ConfirmDialog';
import { ErrorBoundary } from './components/ui/ErrorBoundary';
const Home = React.lazy(() => import('./pages/Home').then(m => ({ default: m.Home })));
const Templates = React.lazy(() => import('./pages/Templates').then(m => ({ default: m.Templates })));
const Dashboard = React.lazy(() => import('./pages/Dashboard').then(m => ({ default: m.Dashboard })));
const Editor = React.lazy(() => import('./pages/editor/Editor').then(m => ({ default: m.Editor })));
const Login = React.lazy(() => import('./pages/auth/Login').then(m => ({ default: m.Login })));
const Register = React.lazy(() => import('./pages/auth/Register').then(m => ({ default: m.Register })));
const OAuthCallback = React.lazy(() => import('./pages/auth/OAuthCallback').then(m => ({ default: m.OAuthCallback })));
const Settings = React.lazy(() => import('./pages/Settings').then(m => ({ default: m.Settings })));
const PrintResume = React.lazy(() => import('./pages/print/PrintResume').then(m => ({ default: m.PrintResume })));
const Pricing = React.lazy(() => import('./pages/Pricing').then(m => ({ default: m.Pricing })));
const PublicResume = React.lazy(() => import('./pages/public/PublicResume').then(m => ({ default: m.PublicResume })));

import { ProtectedRoute } from './components/ProtectedRoute';
import { ProtectedAdminRoute } from './components/ProtectedAdminRoute';
const AdminLayout = React.lazy(() => import('./pages/admin/AdminLayout').then(m => ({ default: m.AdminLayout })));
const UsersPage = React.lazy(() => import('./pages/admin/UsersPage').then(m => ({ default: m.UsersPage })));
const ResumesPage = React.lazy(() => import('./pages/admin/ResumesPage').then(m => ({ default: m.ResumesPage })));
const TemplatesPage = React.lazy(() => import('./pages/admin/TemplatesPage').then(m => ({ default: m.TemplatesPage })));
const SharesPage = React.lazy(() => import('./pages/admin/SharesPage').then(m => ({ default: m.SharesPage })));
const ConfigPage = React.lazy(() => import('./pages/admin/ConfigPage').then(m => ({ default: m.ConfigPage })));
const CatalogPage = React.lazy(() => import('./pages/admin/CatalogPage').then(m => ({ default: m.CatalogPage })));
const AdminHome = React.lazy(() => import('./pages/admin/AdminHome').then(m => ({ default: m.AdminHome })));

const LayoutWithFooter = () => (
  <MainLayout>
    <Outlet />
  </MainLayout>
);

const App: React.FC = () => {
  return (
    <LanguageProvider>
      <ToastProvider>
        <ConfirmDialogProvider>
          <AuthProvider>
            <Router>
              <ErrorBoundary>
                <Suspense fallback={<div className="min-h-screen flex items-center justify-center bg-gray-50"><div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div></div>}>
                  <Routes>
                    <Route element={<LayoutWithFooter />}>
                      <Route path={AppRoute.Home} element={<Home />} />
                    </Route>
                    <Route element={<MainLayout hideFooter><Outlet /></MainLayout>}>
                      <Route path={AppRoute.Templates} element={<Templates />} />
                      <Route path={AppRoute.Pricing} element={<Pricing />} />
                    </Route>

                    <Route element={<ProtectedRoute />}>
                      <Route element={<MainLayout hideFooter><Outlet /></MainLayout>}>
                        <Route path={AppRoute.Dashboard} element={<Dashboard />} />
                        <Route path={AppRoute.Settings} element={<Settings />} />
                      </Route>
                    </Route>
                    
                    <Route path={AppRoute.Login} element={<Login />} />
                    <Route path={AppRoute.Register} element={<Register />} />
                    <Route path={AppRoute.OAuthCallback} element={<OAuthCallback />} />
                    
                    <Route element={<ProtectedRoute />}>
                        <Route path={AppRoute.Editor} element={<Editor />} />
                    </Route>
                    
                    <Route element={<ProtectedAdminRoute />}>
                      <Route element={<AdminLayout />}>
                        <Route path={AppRoute.Admin} element={<AdminHome />} />
                        <Route path={AppRoute.AdminUsers} element={<UsersPage />} />
                        <Route path={AppRoute.AdminResumes} element={<ResumesPage />} />
                        <Route path={AppRoute.AdminTemplates} element={<TemplatesPage />} />
                        <Route path={AppRoute.AdminShares} element={<SharesPage />} />
                        <Route path={AppRoute.AdminCatalog} element={<CatalogPage />} />
                        <Route path={AppRoute.AdminConfigs} element={<ConfigPage />} />
                      </Route>
                    </Route>
                    
                    <Route path={AppRoute.Print} element={<PrintResume />} />
                    <Route path={AppRoute.Public} element={<PublicResume />} />
                  </Routes>
                </Suspense>
              </ErrorBoundary>
            </Router>
          </AuthProvider>
        </ConfirmDialogProvider>
      </ToastProvider>
    </LanguageProvider>
  );
};

export default App;
