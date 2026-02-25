import React, { useEffect } from 'react';
import { Navigate, Outlet, useLocation } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import { AppRoute } from '../types';

export const ProtectedAdminRoute: React.FC = () => {
  const { isAuthenticated, loading, isAdmin, authModalOpen, openAuthModal } = useAuth();
  const location = useLocation();

  useEffect(() => {
    if (loading || isAuthenticated || authModalOpen) return;
    const returnTo = location.pathname + location.search + location.hash;
    openAuthModal({ mode: 'login', returnTo, source: 'protected' });
  }, [authModalOpen, isAuthenticated, loading, location.hash, location.pathname, location.search, openAuthModal]);

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
      </div>
    );
  }

  if (!isAuthenticated) {
    return null;
  }

  if (!isAdmin) {
    return <Navigate to={AppRoute.Dashboard} replace />;
  }

  return <Outlet />;
};
