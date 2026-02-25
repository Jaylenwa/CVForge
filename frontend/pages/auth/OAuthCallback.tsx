import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { API_BASE } from '../../config';
import { AppRoute } from '../../types';
import { Button } from '../../components/ui/Button';
import { useAuth } from '../../contexts/AuthContext';

export const OAuthCallback: React.FC = () => {
  const navigate = useNavigate();
  const { openAuthModal } = useAuth();
  const [status, setStatus] = useState<'loading' | 'error' | 'ok'>('loading');
  const [message, setMessage] = useState('');
  useEffect(() => {
    const qs = new URLSearchParams(window.location.hash.split('?')[1] || '');
    const ott = qs.get('ott') || '';
    if (!ott) {
      setStatus('error');
      setMessage('缺少参数');
      return;
    }
    (async () => {
      try {
        const res = await fetch(`${API_BASE}/auth/wechat/consume-ott`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ ott })
        });
        if (!res.ok) {
          setStatus('error');
          setMessage('登录失败');
          return;
        }
        const data = await res.json();
        const at = data.accessToken as string;
        const rt = data.refreshToken as string;
        if (!at || !rt) {
          setStatus('error');
          setMessage('登录失败');
          return;
        }
        localStorage.setItem('token', at);
        localStorage.setItem('refreshToken', rt);
        setStatus('ok');
        navigate(AppRoute.Home, { replace: true });
      } catch {
        setStatus('error');
        setMessage('网络错误');
      }
    })();
  }, [navigate]);
  if (status === 'loading') return <div className="p-6 text-center">正在处理登录...</div>;
  if (status === 'error') return (
    <div className="p-6 text-center">
      <div className="text-red-600 mb-3">{message}</div>
      <Button
        onClick={() => {
          openAuthModal({ mode: 'login', returnTo: AppRoute.Home, source: 'route' });
          navigate(AppRoute.Home, { replace: true });
        }}
      >
        返回登录
      </Button>
    </div>
  );
  return <div className="p-6 text-center">登录成功，正在跳转...</div>;
};
