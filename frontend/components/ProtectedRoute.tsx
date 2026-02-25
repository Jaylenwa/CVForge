import React, { useEffect } from 'react';
import { Outlet, useLocation } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';

export const ProtectedRoute: React.FC = () => {
  const { isAuthenticated, loading, authModalOpen, openAuthModal } = useAuth();
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

  return <Outlet />;
};
