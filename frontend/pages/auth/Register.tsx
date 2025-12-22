import React, { useState, useEffect } from 'react';
import { Link, useNavigate, useLocation } from 'react-router-dom';
import { Mail, Lock, KeyRound, CheckCircle } from 'lucide-react';
import { AuthLayout } from './AuthLayout';
import { Button } from '../../components/ui/Button';
import { useLanguage } from '../../contexts/LanguageContext';
import { useAuth } from '../../contexts/AuthContext';
import { sendVerificationCode, verifyCode, registerUser } from '../../services/authService';
import { getAuthConfig } from '../../services/configService';
import { AppRoute, AuthConfig } from '../../types';

export const Register: React.FC = () => {
  const { t } = useLanguage();
  const { login } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  
  const [step, setStep] = useState<1 | 2>(1); // 1: Email, 2: Code & Password
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [authConfig, setAuthConfig] = useState<AuthConfig | null>(null);

  useEffect(() => {
    getAuthConfig().then(setAuthConfig);
  }, []);
  const emailVerificationEnabled = !!(authConfig && authConfig.enableEmailVerification);
  
  // Timer for code
  const [countdown, setCountdown] = useState(0);

  const [formData, setFormData] = useState({
    email: '',
    code: '',
    password: '',
    confirmPassword: ''
  });

  useEffect(() => {
    let timer: number;
    if (countdown > 0) {
      timer = window.setInterval(() => setCountdown(c => c - 1), 1000);
    }
    return () => clearInterval(timer);
  }, [countdown]);

  const handleSendCode = async () => {
    if (!formData.email || !formData.email.includes('@')) {
        setError(t('auth.error.invalidEmail'));
        return;
    }
    
    setError('');
    setLoading(true);
    
    try {
        await sendVerificationCode(formData.email);
        setStep(2);
        setCountdown(60); // 60s cooldown
    } catch (e) {
        setError(t('auth.error.sendCodeFailed'));
    } finally {
        setLoading(false);
    }
  };

  const handleRegister = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    if (formData.password !== formData.confirmPassword) {
        setError(t('auth.error.passwordMismatch'));
        return;
    }

    if (formData.password.length < 6) {
        setError(t('auth.error.passwordTooShort'));
        return;
    }

    setLoading(true);

  try {
        if (emailVerificationEnabled) {
            const isValid = await verifyCode(formData.email, formData.code);
            if (!isValid) {
                setError(t('auth.error.invalidCode'));
                setLoading(false);
                return;
            }
        }
        const res = await registerUser(formData.email, formData.code, formData.password);
        if (res.success) {
            await login(formData.email);
            const from = (location.state as any)?.from?.pathname || AppRoute.Home;
            navigate(from, { replace: true });
        } else {
            setError(t('auth.error.registrationFailed'));
        }
    } catch (e) {
        setError(t('auth.error.registrationFailed'));
    } finally {
        setLoading(false);
    }
  };

  return (
    <AuthLayout
        image="https://images.unsplash.com/photo-1517245386807-bb43f82c33c4?ixlib=rb-1.2.1&auto=format&fit=crop&w=1950&q=80"
        quote={t('auth.quote')}
        author={t('auth.quoteAuthor')}
    >
      <div className="text-center">
        <h2 className="text-3xl font-extrabold text-gray-900">{t('auth.createAccount')}</h2>
        <p className="mt-2 text-sm text-gray-600">{t('auth.createDesc')}</p>
      </div>

      <form className="mt-8 space-y-6" onSubmit={handleRegister}>
        {/* Step 1: Email Input */}
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
                    type="email"
                    required
                    disabled={step === 2 && countdown > 0}
                    className={`block w-full pl-10 pr-3 py-2 border rounded-md focus:outline-none sm:text-sm ${step === 2 ? 'bg-gray-100 text-gray-500 border-gray-200' : 'border-gray-300 focus:ring-blue-500 focus:border-blue-500'}`}
                    placeholder={t('auth.placeholder.email')}
                    value={formData.email}
                    onChange={(e) => setFormData({...formData, email: e.target.value})}
                />
                {step === 2 && (
                    <div className="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none">
                        <CheckCircle className="h-5 w-5 text-green-500" />
                    </div>
                )}
            </div>
        </div>

        {/* Verification Code (only when enabled and after sending) */}
        {emailVerificationEnabled && step === 2 && (
          <div className="space-y-4 animate-fadeIn">
            <div className="bg-green-50 text-green-700 p-3 rounded-md text-sm text-center">
              {t('auth.success.codeSent')}
            </div>
            <div>
              <label htmlFor="code" className="block text-sm font-medium text-gray-700">
                {t('auth.verificationCode')}
              </label>
              <div className="mt-1 flex space-x-2">
                <div className="relative rounded-md shadow-sm flex-grow">
                  <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                    <KeyRound className="h-5 w-5 text-gray-400" />
                  </div>
                  <input
                    id="code"
                    type="text"
                    required
                    className="block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                    placeholder={t('auth.placeholder.code')}
                    value={formData.code}
                    onChange={(e) => setFormData({ ...formData, code: e.target.value })}
                  />
                </div>
                <Button
                  type="button"
                  variant="outline"
                  disabled={countdown > 0}
                  onClick={handleSendCode}
                  className="whitespace-nowrap w-32"
                >
                  {countdown > 0 ? `${countdown}s` : t('auth.resend')}
                </Button>
              </div>
            </div>
          </div>
        )}

        {/* Password fields (always needed; show immediately if verification disabled) */}
        {(!emailVerificationEnabled || step === 2) && (
          <div className="space-y-4">
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
                  type="password"
                  required
                  className="block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                  placeholder={t('auth.placeholder.passwordMin')}
                  value={formData.password}
                  onChange={(e) => setFormData({ ...formData, password: e.target.value })}
                />
              </div>
            </div>

            <div>
              <label htmlFor="confirmPassword" className="block text-sm font-medium text-gray-700">
                {t('auth.confirmPassword')}
              </label>
              <div className="mt-1 relative rounded-md shadow-sm">
                <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                  <Lock className="h-5 w-5 text-gray-400" />
                </div>
                <input
                  id="confirmPassword"
                  type="password"
                  required
                  className="block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                  placeholder={t('auth.placeholder.reenterPassword')}
                  value={formData.confirmPassword}
                  onChange={(e) => setFormData({ ...formData, confirmPassword: e.target.value })}
                />
              </div>
            </div>
          </div>
        )}

        {error && (
            <div className="text-red-500 text-sm text-center bg-red-50 p-2 rounded">
                {error}
            </div>
        )}
        
        {emailVerificationEnabled && step === 1 ? (
             <Button 
                type="button" 
                onClick={handleSendCode} 
                className="w-full" 
                size="lg" 
                isLoading={loading}
            >
                {t('auth.sendCode')}
            </Button>
        ) : (
            <Button 
                type="submit" 
                className="w-full" 
                size="lg" 
                isLoading={loading}
            >
                {t('auth.register')}
            </Button>
        )}

        <div className="text-center mt-4">
             <span className="text-gray-600 text-sm">{t('auth.hasAccount')} </span>
             <Link to={AppRoute.Login} className="text-blue-600 font-medium hover:underline text-sm">
                {t('auth.login')}
            </Link>
        </div>
        
        <p className="text-xs text-center text-gray-400 mt-6">
            {t('auth.terms')}
        </p>
      </form>
    </AuthLayout>
  );
};
