import React from 'react';
import { HashRouter as Router, Routes, Route, Outlet } from 'react-router-dom';
import { MainLayout } from './components/Layout';
import { Home } from './pages/Home';
import { Templates } from './pages/Templates';
import { Dashboard } from './pages/Dashboard';
import { Editor } from './pages/editor/Editor';
import { Login } from './pages/auth/Login';
import { Register } from './pages/auth/Register';
import { OAuthCallback } from './pages/auth/OAuthCallback';
import { Settings } from './pages/Settings';
import { AppRoute } from './types';
import { LanguageProvider } from './contexts/LanguageContext';
import { AuthProvider } from './contexts/AuthContext';
import { ToastProvider } from './components/ui/Toast';
import { ConfirmDialogProvider } from './components/ui/ConfirmDialog';
import { PrintResume } from './pages/print/PrintResume';
import { Pricing } from './pages/Pricing';

import { ProtectedRoute } from './components/ProtectedRoute';
import { ProtectedAdminRoute } from './components/ProtectedAdminRoute';
import { AdminLayout } from './pages/admin/AdminLayout';
import { UsersPage } from './pages/admin/UsersPage';
import { ResumesPage } from './pages/admin/ResumesPage';
import { TemplatesPage } from './pages/admin/TemplatesPage';
import { SharesPage } from './pages/admin/SharesPage';
import { ConfigPage } from './pages/admin/ConfigPage';
import { AdminHome } from './pages/admin/AdminHome';

const LayoutWrapper = () => (
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
          <Routes>
              {/* Public Routes with Main Navbar/Footer */}
              <Route element={<LayoutWrapper />}>
                  <Route path={AppRoute.Home} element={<Home />} />
                  <Route path={AppRoute.Templates} element={<Templates />} />
                  <Route path={AppRoute.Pricing} element={<Pricing />} />
              </Route>

              {/* Protected Routes with Main Navbar/Footer */}
              <Route element={<ProtectedRoute />}>
                  <Route element={<LayoutWrapper />}>
                      <Route path={AppRoute.Dashboard} element={<Dashboard />} />
                      <Route path={AppRoute.Settings} element={<Settings />} />
                  </Route>
              </Route>
              
              {/* Auth Routes */}
              <Route path={AppRoute.Login} element={<Login />} />
              <Route path={AppRoute.Register} element={<Register />} />
              <Route path={AppRoute.OAuthCallback} element={<OAuthCallback />} />
              
              {/* Protected Standalone Routes */}
              <Route element={<ProtectedRoute />}>
                  <Route path={AppRoute.Editor} element={<Editor />} />
              </Route>
              
              {/* Admin Routes */}
              <Route element={<ProtectedAdminRoute />}>
                <Route element={<AdminLayout />}>
                  <Route path={AppRoute.Admin} element={<AdminHome />} />
                  <Route path={AppRoute.AdminUsers} element={<UsersPage />} />
                  <Route path={AppRoute.AdminResumes} element={<ResumesPage />} />
                  <Route path={AppRoute.AdminTemplates} element={<TemplatesPage />} />
                  <Route path={AppRoute.AdminShares} element={<SharesPage />} />
                  <Route path={AppRoute.AdminConfigs} element={<ConfigPage />} />
                </Route>
              </Route>
              
              {/* Standalone Routes */}
              <Route path={AppRoute.Print} element={<PrintResume />} />
          </Routes>
            </Router>
          </AuthProvider>
        </ConfirmDialogProvider>
      </ToastProvider>
    </LanguageProvider>
  );
};

export default App;
