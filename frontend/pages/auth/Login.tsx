import React, { useState } from 'react';
import { Link, useNavigate, useLocation } from 'react-router-dom';
import { Mail, Lock, Github, ArrowLeft, Globe } from 'lucide-react';
import { Button } from '../../components/ui/Button';
import { useLanguage } from '../../contexts/LanguageContext';
import { useAuth } from '../../contexts/AuthContext';
import { loginUser } from '../../services/authService';
import { AppRoute } from '../../types';

export const Login: React.FC = () => {
  const { t, language, setLanguage } = useLanguage();
  const { login, loginWithGithub } = useAuth();
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

  const handleGithubLogin = async () => {
    setError('');
    setLoading(true);
    try {
        const success = await loginWithGithub();
        if (success) {
            const from = (location.state as any)?.from?.pathname || AppRoute.Home;
            navigate(from, { replace: true });
        } else {
            setError(t('auth.error.general'));
        }
    } catch (err) {
        setError(t('auth.error.general'));
    } finally {
        setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 flex flex-col justify-center py-12 sm:px-6 lg:px-8 relative">
       {/* Language Switcher & Back Home */}
       <div className="absolute top-6 right-6">
          <button
            onClick={() => setLanguage(language === 'en' ? 'zh' : 'en')}
            className="p-2 text-gray-500 hover:text-gray-900 focus:outline-none"
            title={t('lang.switchTitle')}
          >
            <div className="flex items-center space-x-1">
              <Globe size={18} />
              <span className="text-sm font-medium">
                {language === 'en' ? t('lang.en_short') : t('lang.zh_short')}
              </span>
            </div>
          </button>
        </div>
        <div className="absolute top-6 left-6">
             <Link to="/" className="flex items-center text-gray-500 hover:text-gray-900 transition-colors">
                <ArrowLeft size={20} className="mr-2" /> {t('common.backHome')}
             </Link>
        </div>

      <div className="sm:mx-auto sm:w-full sm:max-w-md">
        <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">
          {t('auth.welcome')}
        </h2>
        <p className="mt-2 text-center text-sm text-gray-600">
          {t('auth.welcomeDesc')}
        </p>
      </div>

      <div className="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
        <div className="bg-white py-8 px-4 shadow sm:rounded-lg sm:px-10 border border-gray-100">
          <form className="space-y-6" onSubmit={handleSubmit}>
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
                  className="appearance-none block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
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
                  className="appearance-none block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                  placeholder={t('auth.placeholder.password')}
                  value={formData.password}
                  onChange={(e) => setFormData({...formData, password: e.target.value})}
                />
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
          </form>

          <div className="mt-6">
            <div className="relative">
              <div className="absolute inset-0 flex items-center">
                <div className="w-full border-t border-gray-300" />
              </div>
              <div className="relative flex justify-center text-sm">
                <span className="px-2 bg-white text-gray-500">
                  {t('auth.orContinueWith') || 'Or continue with'}
                </span>
              </div>
            </div>

            <div className="mt-6 grid grid-cols-1 gap-3">
              <div>
                <Button
                  variant="outline"
                  className="w-full flex justify-center items-center"
                  onClick={handleGithubLogin}
                  isLoading={loading}
                >
                  <Github className="h-5 w-5 mr-2" />
                  GitHub
                </Button>
              </div>
            </div>
          </div>
          
           <div className="mt-6 text-center">
                <span className="text-gray-600 text-sm">{t('auth.noAccount')} </span>
                <Link to={AppRoute.Register} className="text-blue-600 font-medium hover:underline text-sm">
                    {t('auth.register')}
                </Link>
            </div>
        </div>
      </div>
    </div>
  );
};
