import React, { createContext, useContext, useState, ReactNode, useEffect, useRef } from 'react';
import { API_BASE } from '../config';

export interface User {
  email: string;
  name: string;
  avatarUrl?: string;
}

interface AuthContextType {
  user: User | null;
  isAuthenticated: boolean;
  loading: boolean;
  isAdmin: boolean;
  login: (email: string) => void;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [isAdmin, setIsAdmin] = useState(false);
  const refreshTimer = useRef<number | null>(null);

  const loadUser = async () => {
    const token = localStorage.getItem('token');
    if (!token) {
      setLoading(false);
      setIsAdmin(false);
      return;
    }
    try {
      const res = await fetch(`${API_BASE}/users/me`, { headers: { Authorization: `Bearer ${token}` } });
      if (res.ok) {
        const data = await res.json();
        setUser({ email: data.email, name: data.name || data.email.split('@')[0], avatarUrl: data.avatarUrl });
      }
      if (res.status === 401) {
        const ok = await refreshTokens();
        if (ok) {
          const res2 = await fetch(`${API_BASE}/users/me`, { headers: { Authorization: `Bearer ${localStorage.getItem('token') || ''}` } });
          if (res2.ok) {
            const data2 = await res2.json();
            setUser({ email: data2.email, name: data2.name || data2.email.split('@')[0], avatarUrl: data2.avatarUrl });
          }
        } else {
          localStorage.removeItem('token');
          localStorage.removeItem('refreshToken');
          setUser(null);
          setIsAdmin(false);
          window.location.hash = '#/login';
        }
      }
      // admin probing
      try {
        const probe = await fetch(`${API_BASE}/admin/users?page=1&pageSize=1`, { headers: { Authorization: `Bearer ${localStorage.getItem('token') || ''}` } });
        setIsAdmin(probe.ok);
      } catch {
        setIsAdmin(false);
      }
    } catch {}
    setLoading(false);
  };

  useEffect(() => { loadUser(); }, []);

  const login = (_email: string) => { loadUser(); };

  const logout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('refreshToken');
    if (refreshTimer.current) { window.clearTimeout(refreshTimer.current); refreshTimer.current = null; }
    setUser(null);
    setIsAdmin(false);
  };

  const decodeExp = (token: string) => {
    try {
      const parts = token.split('.');
      if (parts.length < 2) return 0;
      const payload = JSON.parse(atob(parts[1].replace(/-/g, '+').replace(/_/g, '/')));
      return typeof payload.exp === 'number' ? payload.exp : 0;
    } catch { return 0; }
  };

  const scheduleRefresh = () => {
    const token = localStorage.getItem('token');
    if (!token) return;
    const exp = decodeExp(token);
    if (!exp) return;
    const nowSec = Math.floor(Date.now() / 1000);
    const lead = 60;
    const ms = Math.max((exp - nowSec - lead) * 1000, 0);
    if (refreshTimer.current) window.clearTimeout(refreshTimer.current);
    refreshTimer.current = window.setTimeout(() => { refreshTokens(); }, ms);
  };

  const refreshTokens = async () => {
    const refresh = localStorage.getItem('refreshToken');
    if (!refresh) return false;
    try {
      const res = await fetch(`${API_BASE}/auth/refresh`, { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ refreshToken: refresh }) });
      if (!res.ok) return false;
      const data = await res.json();
      const at = data.accessToken as string;
      const rt = data.refreshToken as string;
      if (!at || !rt) return false;
      localStorage.setItem('token', at);
      localStorage.setItem('refreshToken', rt);
      scheduleRefresh();
      return true;
    } catch { return false; }
  };

  useEffect(() => { scheduleRefresh(); }, [user]);

  return (
    <AuthContext.Provider value={{ user, isAuthenticated: !!user, loading, isAdmin, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
