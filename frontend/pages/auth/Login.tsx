import React, { useState } from 'react';
import { Link, useNavigate, useLocation } from 'react-router-dom';
import { Mail, Lock } from 'lucide-react';
import { AuthLayout } from './AuthLayout';
import { Button } from '../../components/ui/Button';
import { useLanguage } from '../../contexts/LanguageContext';
import { useAuth } from '../../contexts/AuthContext';
import { loginUser } from '../../services/authService';
import { AppRoute } from '../../types';
import { QrCode } from 'lucide-react';

export const Login: React.FC = () => {
  const { t } = useLanguage();
  const { login, loginWithWeChat } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  
  const [formData, setFormData] = useState({
    email: '',
    password: ''
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    
    if (!formData.email || !formData.password) {
        setError(t('auth.error.fillAll'));
        return;
    }

    setLoading(true);
    try {
        const result = await loginUser(formData.email, formData.password);
        if (result.success) {
            await login(formData.email);
            const from = (location.state as any)?.from?.pathname || AppRoute.Home;
            navigate(from, { replace: true });
        } else {
            setError(t('auth.error.invalidCredentials'));
        }
    } catch (err) {
        setError(t('auth.error.general'));
    } finally {
        setLoading(false);
    }
  };

  return (
    <AuthLayout 
        image="https://images.unsplash.com/photo-1497215728101-856f4ea42174?ixlib=rb-1.2.1&auto=format&fit=crop&w=1950&q=80"
        quote={t('auth.quote')}
        author={t('auth.quoteAuthor')}
    >
      <div className="text-center">
        <h2 className="text-3xl font-extrabold text-gray-900">{t('auth.welcome')}</h2>
        <p className="mt-2 text-sm text-gray-600">{t('auth.welcomeDesc')}</p>
      </div>

      <form className="mt-8 space-y-6" onSubmit={handleSubmit}>
        <div className="space-y-4">
            <div>
                <label htmlFor="email" className="block text-sm font-medium text-gray-700">
                    {t('auth.email')}
                </label>
                <div className="mt-1 relative rounded-md shadow-sm">
                    <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                        <Mail className="h-5 w-5 text-gray-400" />
                    </div>
                    <input
                        id="email"
                        name="email"
                        type="email"
                        autoComplete="email"
                        required
                        className="block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                        placeholder={t('auth.placeholder.email')}
                        value={formData.email}
                        onChange={(e) => setFormData({...formData, email: e.target.value})}
                    />
                </div>
            </div>

            <div>
                <label htmlFor="password" className="block text-sm font-medium text-gray-700">
                    {t('auth.password')}
                </label>
                <div className="mt-1 relative rounded-md shadow-sm">
                    <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                        <Lock className="h-5 w-5 text-gray-400" />
                    </div>
                    <input
                        id="password"
                        name="password"
                        type="password"
                        autoComplete="current-password"
                        required
                        className="block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                        placeholder={t('auth.placeholder.password')}
                        value={formData.password}
                        onChange={(e) => setFormData({...formData, password: e.target.value})}
                    />
                </div>
            </div>
        </div>

        {error && (
            <div className="text-red-500 text-sm text-center bg-red-50 p-2 rounded">
                {error}
            </div>
        )}

        <div className="flex items-center justify-between">
          <div className="text-sm">
            <a href="#" className="font-medium text-blue-600 hover:text-blue-500">
              {t('auth.forgotPassword')}
            </a>
          </div>
        </div>

        <div>
          <Button type="submit" className="w-full" size="lg" isLoading={loading}>
            {t('auth.login')}
          </Button>
        </div>

        {/*
        <div className="mt-3">
          <Button type="button" className="w-full" variant="outline" onClick={async () => {
            setError('');
            setLoading(true);
            try {
              const ok = await loginWithWeChat();
              if (ok) {
                const from = (location.state as any)?.from?.pathname || AppRoute.Home;
                navigate(from, { replace: true });
              } else {
                setError(t('auth.error.general'));
              }
            } finally {
              setLoading(false);
            }
          }}>
            <QrCode className="mr-2 h-4 w-4" /> {t('auth.loginWithWeChat') || '微信登录'}
          </Button>
        </div>
        */}

        <div className="text-center mt-4">
            <span className="text-gray-600 text-sm">{t('auth.noAccount')} </span>
            <Link to={AppRoute.Register} className="text-blue-600 font-medium hover:underline text-sm">
                {t('auth.register')}
            </Link>
        </div>
      </form>
    </AuthLayout>
  );
};
