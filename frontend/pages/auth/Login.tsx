import React, { useEffect } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { AuthModalMode, useAuth } from '../../contexts/AuthContext';
import { AppRoute } from '../../types';

type Props = { initialMode?: AuthModalMode };

export const Login: React.FC<Props> = ({ initialMode = 'login' }) => {
  const { openAuthModal } = useAuth();
  const location = useLocation();
  const navigate = useNavigate();

  useEffect(() => {
    const from = (location.state as any)?.from;
    const returnTo = from ? (from.pathname + (from.search || '') + (from.hash || '')) : AppRoute.Home;
    openAuthModal({ mode: initialMode, returnTo, source: 'route' });
    navigate(returnTo, { replace: true });
  }, [initialMode, location.state, navigate, openAuthModal]);

  return null;
};

