import React, { useEffect, useState } from 'react';
import { Mail, Lock, Github, KeyRound, CheckCircle, X } from 'lucide-react';
import { Button } from '../../components/ui/Button';
import { WeChatIcon } from '../../components/ui/WeChatIcon';
import { useLanguage } from '../../contexts/LanguageContext';
import { useAuth } from '../../contexts/AuthContext';
import { loginUser, sendVerificationCode, verifyCode, registerUser } from '../../services/authService';
import { getAuthConfig } from '../../services/configService';
import { AppRoute, AuthConfig } from '../../types';
import { useLocation, useNavigate } from 'react-router-dom';
import { AuthLayout } from './AuthLayout';
import { AnimatePresence, motion, LayoutGroup } from 'framer-motion';
import { useRef } from 'react';
import { Modal } from '../../components/ui/Modal';
import { consumeOtt, createWeChatMPScene, getWeChatMPSceneStatus } from '../../services/wechatMpAuthService';

type Mode = 'login' | 'register';
type Props = { initialMode?: Mode };

export const Login: React.FC<Props> = ({ initialMode = 'login' }) => {
  const { t } = useLanguage();
  const { login, loginWithGithub } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const [authConfig, setAuthConfig] = useState<AuthConfig | null>(null);
  const [mode, setMode] = useState<Mode>(initialMode);
  const hasMountedRef = useRef(false);
  useEffect(() => { hasMountedRef.current = true; }, []);
  const shouldInitialAnimate = hasMountedRef.current;
  const mpPollRef = useRef<number | null>(null);
  const [mpOpen, setMpOpen] = useState(false);
  const [mpQrUrl, setMpQrUrl] = useState('');
  const [mpError, setMpError] = useState('');
  const [mpExpired, setMpExpired] = useState(false);

  useEffect(() => {
    getAuthConfig().then(setAuthConfig);
  }, []);

  const emailVerificationEnabled = !!(authConfig && authConfig.enableEmailVerification);

  const [loginForm, setLoginForm] = useState({ email: '', password: '' });
  const [loginLoading, setLoginLoading] = useState(false);
  const [loginError, setLoginError] = useState('');

  const [registerForm, setRegisterForm] = useState({ email: '', code: '', password: '', confirmPassword: '' });
  const [registerLoading, setRegisterLoading] = useState(false);
  const [registerError, setRegisterError] = useState('');
  const [countdown, setCountdown] = useState(0);

  const clearLoginState = () => {
    setLoginForm({ email: '', password: '' });
    setLoginError('');
    setLoginLoading(false);
  };

  const clearRegisterState = () => {
    setRegisterForm({ email: '', code: '', password: '', confirmPassword: '' });
    setRegisterError('');
    setRegisterLoading(false);
    setCountdown(0);
  };

  const handleSwitchMode = (next: Mode) => {
    if (next === 'register') {
      clearLoginState();
    } else {
      clearRegisterState();
    }
    setMode(next);
  };

  useEffect(() => {
    let timer: number;
    if (countdown > 0) {
      timer = window.setInterval(() => setCountdown((c) => c - 1), 1000);
    }
    return () => clearInterval(timer);
  }, [countdown]);
  
  

  const handleLoginSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoginError('');
    if (!loginForm.email || !loginForm.password) {
      setLoginError(t('auth.error.fillAll'));
      return;
    }
    setLoginLoading(true);
    try {
      const result = await loginUser(loginForm.email, loginForm.password);
      if (result.success) {
        await login(loginForm.email);
        const from = (location.state as any)?.from?.pathname || AppRoute.Home;
        navigate(from, { replace: true });
      } else {
        setLoginError(t('auth.error.invalidCredentials'));
      }
    } catch {
      setLoginError(t('auth.error.general'));
    } finally {
      setLoginLoading(false);
    }
  };

  const handleGithubLogin = async () => {
    setLoginError('');
    setLoginLoading(true);
    try {
      const success = await loginWithGithub();
      if (success) {
        const from = (location.state as any)?.from?.pathname || AppRoute.Home;
        navigate(from, { replace: true });
      } else {
        setLoginError(t('auth.error.general'));
      }
    } catch {
      setLoginError(t('auth.error.general'));
    } finally {
      setLoginLoading(false);
    }
  };

  const stopMpPolling = () => {
    if (mpPollRef.current) {
      window.clearInterval(mpPollRef.current);
      mpPollRef.current = null;
    }
  };

  useEffect(() => {
    return () => stopMpPolling();
  }, []);

  const openWeChatMPModal = async () => {
    setLoginError('');
    setMpError('');
    setMpExpired(false);
    stopMpPolling();
    setMpQrUrl('');
    setMpOpen(true);
    try {
      const data = await createWeChatMPScene();
      const expiresAt = Date.now() + Math.max(1, data.expiresIn) * 1000;
      setMpQrUrl(data.qrUrl);
      const scene = data.scene;
      mpPollRef.current = window.setInterval(async () => {
        try {
          if (!scene) return;
          if (Date.now() > expiresAt) {
            setMpExpired(true);
            stopMpPolling();
            return;
          }
          const st = await getWeChatMPSceneStatus(scene);
          if (st.status === 'expired') {
            setMpExpired(true);
            stopMpPolling();
            return;
          }
          if (st.status === 'ok' && st.ott) {
            stopMpPolling();
            const tokens = await consumeOtt(st.ott);
            if (!tokens.accessToken || !tokens.refreshToken) {
              setMpError(t('auth.error.general'));
              return;
            }
            localStorage.setItem('token', tokens.accessToken);
            localStorage.setItem('refreshToken', tokens.refreshToken);
            await login('');
            setMpOpen(false);
            const from = (location.state as any)?.from?.pathname || AppRoute.Home;
            navigate(from, { replace: true });
          }
        } catch (e) {
          const msg = e instanceof Error ? e.message : '';
          setMpError(msg || t('auth.error.general'));
        }
      }, 1200);
    } catch (e) {
      const msg = e instanceof Error ? e.message : '';
      setMpError(msg || t('auth.error.general'));
      setMpExpired(true);
    }
  };

  const closeWeChatMPModal = () => {
    stopMpPolling();
    setMpOpen(false);
  };

  const handleSendCode = async () => {
    if (!registerForm.email || !registerForm.email.includes('@')) {
      setRegisterError(t('auth.error.invalidEmail'));
      return;
    }
    setRegisterError('');
    setRegisterLoading(true);
    try {
      await sendVerificationCode(registerForm.email);
      setCountdown(60);
    } catch {
      setRegisterError(t('auth.error.sendCodeFailed'));
    } finally {
      setRegisterLoading(false);
    }
  };

  const handleRegisterSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setRegisterError('');
    if (registerForm.password !== registerForm.confirmPassword) {
      setRegisterError(t('auth.error.passwordMismatch'));
      return;
    }
    if (registerForm.password.length < 6) {
      setRegisterError(t('auth.error.passwordTooShort'));
      return;
    }
    setRegisterLoading(true);
    try {
      if (emailVerificationEnabled) {
        const isValid = await verifyCode(registerForm.email, registerForm.code);
        if (!isValid) {
          setRegisterError(t('auth.error.invalidCode'));
          setRegisterLoading(false);
          return;
        }
      }
      const res = await registerUser(registerForm.email, registerForm.code, registerForm.password);
      if (res.success) {
        await login(registerForm.email);
        const from = (location.state as any)?.from?.pathname || AppRoute.Home;
        navigate(from, { replace: true });
      } else {
        setRegisterError(t('auth.error.registrationFailed'));
      }
    } catch {
      setRegisterError(t('auth.error.registrationFailed'));
    } finally {
      setRegisterLoading(false);
    }
  };

  return (
    <AuthLayout image="https://images.unsplash.com/photo-1517245386807-bb43f82c33c4?ixlib=rb-1.2.1&auto=format&fit=crop&w=1950&q=80" quote={t('auth.quote')} author={t('auth.quoteAuthor')}>
      <LayoutGroup>
        <div className="text-center">
          <motion.div className="min-h-[4.5rem]" layout>
            <AnimatePresence initial={false} mode="wait">
              {mode === 'login' ? (
                <motion.div
                  key="hdr-login"
                  initial={{ opacity: 0 }}
                  animate={{ opacity: 1 }}
                  exit={{ opacity: 0 }}
                  transition={{ duration: 0.16, ease: 'easeOut' }}
                  style={{ willChange: 'opacity', transform: 'translateZ(0)', WebkitFontSmoothing: 'antialiased', backfaceVisibility: 'hidden', textRendering: 'optimizeLegibility' }}
                  layout
                >
                  <h2 className="text-3xl font-extrabold leading-tight text-gray-900">{t('auth.welcome')}</h2>
                  <p className="mt-2 text-sm leading-relaxed text-gray-600">{t('auth.welcomeDesc')}</p>
                </motion.div>
              ) : (
                <motion.div
                  key="hdr-register"
                  initial={{ opacity: 0 }}
                  animate={{ opacity: 1 }}
                  exit={{ opacity: 0 }}
                  transition={{ duration: 0.16, ease: 'easeOut' }}
                  style={{ willChange: 'opacity', transform: 'translateZ(0)', WebkitFontSmoothing: 'antialiased', backfaceVisibility: 'hidden', textRendering: 'optimizeLegibility' }}
                  layout
                >
                  <h2 className="text-3xl font-extrabold leading-tight text-gray-900">{t('auth.createAccount')}</h2>
                  <p className="mt-2 text-sm leading-relaxed text-gray-600">{t('auth.createDesc')}</p>
                </motion.div>
              )}
            </AnimatePresence>
          </motion.div>
        </div>

        <motion.div className="mt-6 relative overflow-hidden" layout>
          <AnimatePresence initial={false} mode="wait">
            {mode === 'login' ? (
              <motion.form
                key="login"
                initial={shouldInitialAnimate ? { x: '-100%', opacity: 0.9 } : false}
                animate={{ x: 0, opacity: 1 }}
                exit={{ x: '-100%', opacity: 0.9 }}
                transition={{ duration: 0.18, ease: 'easeOut' }}
                className="space-y-6 w-full"
                onSubmit={handleLoginSubmit}
                layout
              >
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
                    className="appearance-none block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                    placeholder={t('auth.placeholder.email')}
                    value={loginForm.email}
                    onChange={(e) => setLoginForm({ ...loginForm, email: e.target.value })}
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
                    type="password"
                    required
                    className="appearance-none block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                    placeholder={t('auth.placeholder.password')}
                    value={loginForm.password}
                    onChange={(e) => setLoginForm({ ...loginForm, password: e.target.value })}
                  />
                </div>
              </div>

              {loginError && <div className="text-red-500 text-sm text-center bg-red-50 p-2 rounded">{loginError}</div>}

              <div>
                <Button type="submit" className="w-full" size="lg" isLoading={loginLoading}>
                  {t('auth.login')}
                </Button>
              </div>

              {(authConfig?.enableGithubLogin || authConfig?.enableWeChatMPLogin) && (
                <div className="mt-6">
                  <div className="relative">
                    <div className="absolute inset-0 flex items-center">
                      <div className="w-full border-t border-gray-300" />
                    </div>
                    <div className="relative flex justify-center text-sm">
                      <span className="px-2 bg-white text-gray-500">{t('auth.orContinueWith')}</span>
                    </div>
                  </div>
                  <div className="mt-6 grid grid-cols-2 gap-3">
                    {authConfig?.enableGithubLogin && (
                      <div>
                        <Button variant="outline" className="w-full flex justify-center items-center" onClick={handleGithubLogin} isLoading={loginLoading}>
                          <Github className="h-5 w-5 mr-2" />
                          {t('auth.provider.github')}
                        </Button>
                      </div>
                    )}
                    {authConfig?.enableWeChatMPLogin && (
                      <div>
                        <Button
                          variant="outline"
                          className="w-full flex justify-center items-center text-green-700 border-green-200 hover:bg-green-50"
                          onClick={openWeChatMPModal}
                          isLoading={loginLoading}
                        >
                          <WeChatIcon className="h-5 w-5 mr-2" />
                          {t('auth.provider.wechatMP')}
                        </Button>
                      </div>
                    )}
                  </div>
                </div>
              )}

              <div className="text-center mt-4">
                <span className="text-gray-600 text-sm">{t('auth.noAccount')} </span>
                <button type="button" className="text-blue-600 font-medium hover:underline text-sm" onClick={() => handleSwitchMode('register')}>
                  {t('auth.register')}
                </button>
              </div>
            </motion.form>
          ) : (
            <motion.form
              key="register"
              initial={shouldInitialAnimate ? { x: '100%', opacity: 0.9 } : false}
              animate={{ x: 0, opacity: 1 }}
              exit={{ x: '100%', opacity: 0.9 }}
              transition={{ duration: 0.18, ease: 'easeOut' }}
              className="space-y-6 w-full"
              onSubmit={handleRegisterSubmit}
              layout
            >
              <div>
                <label htmlFor="reg-email" className="block text-sm font-medium text-gray-700">
                  {t('auth.email')}
                </label>
                <div className="mt-1 relative rounded-md shadow-sm">
                  <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                    <Mail className="h-5 w-5 text-gray-400" />
                  </div>
                  <input
                    id="reg-email"
                    type="email"
                    required
                    disabled={countdown > 0}
                    className={`block w-full pl-10 pr-3 py-2 border rounded-md focus:outline-none sm:text-sm ${countdown > 0 ? 'bg-gray-100 text-gray-500 border-gray-200' : 'border-gray-300 focus:ring-blue-500 focus:border-blue-500'}`}
                    placeholder={t('auth.placeholder.email')}
                    value={registerForm.email}
                    onChange={(e) => setRegisterForm({ ...registerForm, email: e.target.value })}
                  />
                  {countdown > 0 && (
                    <div className="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none">
                      <CheckCircle className="h-5 w-5 text-green-500" />
                    </div>
                  )}
                </div>
              </div>

              {emailVerificationEnabled && (
                <div className="space-y-4">
                  {countdown > 0 && (
                    <div className="bg-green-50 text-green-700 p-3 rounded-md text-sm text-center">{t('auth.success.codeSent')}</div>
                  )}
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
                          value={registerForm.code}
                          onChange={(e) => setRegisterForm({ ...registerForm, code: e.target.value })}
                        />
                      </div>
                      <Button type="button" variant="outline" disabled={countdown > 0} onClick={handleSendCode} className="whitespace-nowrap w-32">
                        {countdown > 0 ? `${countdown}s` : t('auth.resend')}
                      </Button>
                    </div>
                  </div>
                </div>
              )}

              <div className="space-y-4">
                <div>
                  <label htmlFor="reg-password" className="block text-sm font-medium text-gray-700">
                    {t('auth.password')}
                  </label>
                  <div className="mt-1 relative rounded-md shadow-sm">
                    <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                      <Lock className="h-5 w-5 text-gray-400" />
                    </div>
                    <input
                      id="reg-password"
                      type="password"
                      required
                      className="block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                      placeholder={t('auth.placeholder.passwordMin')}
                      value={registerForm.password}
                      onChange={(e) => setRegisterForm({ ...registerForm, password: e.target.value })}
                    />
                  </div>
                </div>
                <div>
                  <label htmlFor="reg-confirm" className="block text-sm font-medium text-gray-700">
                    {t('auth.confirmPassword')}
                  </label>
                  <div className="mt-1 relative rounded-md shadow-sm">
                    <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                      <Lock className="h-5 w-5 text-gray-400" />
                    </div>
                    <input
                      id="reg-confirm"
                      type="password"
                      required
                      className="block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                      placeholder={t('auth.placeholder.reenterPassword')}
                      value={registerForm.confirmPassword}
                      onChange={(e) => setRegisterForm({ ...registerForm, confirmPassword: e.target.value })}
                    />
                  </div>
                </div>
              </div>

              {registerError && <div className="text-red-500 text-sm text-center bg-red-50 p-2 rounded">{registerError}</div>}

              <Button type="submit" className="w-full" size="lg" isLoading={registerLoading}>
                {t('auth.register')}
              </Button>

              <div className="text-center mt-4">
                <span className="text-gray-600 text-sm">{t('auth.hasAccount')} </span>
                <button type="button" className="text-blue-600 font-medium hover:underline text-sm" onClick={() => handleSwitchMode('login')}>
                  {t('auth.login')}
                </button>
              </div>
            </motion.form>
          )}
        </AnimatePresence>
        </motion.div>
      </LayoutGroup>
      <Modal
        isOpen={mpOpen}
        onClose={closeWeChatMPModal}
        hideHeader
        size="md"
        closeOnBackdrop
      >
        <div className="px-2 pb-2">
          <div className="relative pt-2">
            <div className="text-center text-2xl font-semibold text-gray-900">{t('auth.wechatMP.title')}</div>
            <button
              type="button"
              onClick={closeWeChatMPModal}
              aria-label={t('common.close')}
              className="absolute right-0 top-0 p-2 text-gray-400 hover:text-gray-600"
            >
              <X size={20} />
            </button>
          </div>

          <div className="mt-4 text-center text-sm text-gray-600">{t('auth.wechatMP.desc')}</div>

          <div className="mt-4 flex items-center justify-center">
            <div className="relative rounded-xl border border-gray-200 bg-white p-3">
              {mpQrUrl ? (
                <img src={mpQrUrl} alt="WeChat QR" className="h-64 w-64" />
              ) : (
                <div className="h-64 w-64 bg-gray-50" />
              )}
              {mpExpired && (
                <button
                  type="button"
                  onClick={openWeChatMPModal}
                  className="absolute inset-0 flex items-center justify-center rounded-xl bg-white/80 text-sm font-medium text-blue-600"
                >
                  {t('auth.wechatMP.refresh')}
                </button>
              )}
            </div>
          </div>

          {mpError && <div className="mt-3 text-center text-sm text-red-600">{mpError}</div>}

          <div className="mt-5 text-center text-sm text-gray-600">
            {t('auth.wechatMP.noAccount')}{' '}
            <button
              type="button"
              className="font-medium text-blue-600 hover:underline"
              onClick={() => {
                closeWeChatMPModal();
                navigate(AppRoute.Register);
              }}
            >
              {t('auth.wechatMP.toRegister')}
            </button>
          </div>

          <div className="mt-6">
            <div className="relative">
              <div className="absolute inset-0 flex items-center">
                <div className="w-full border-t border-gray-200" />
              </div>
              <div className="relative flex justify-center text-sm">
                <span className="bg-white px-3 text-gray-400">{t('auth.wechatMP.otherWays')}</span>
              </div>
            </div>

            <div className="mt-4 flex items-center justify-center gap-6">
              <button
                type="button"
                onClick={closeWeChatMPModal}
                className="flex h-12 w-12 items-center justify-center rounded-full border border-gray-200 text-gray-500 hover:bg-gray-50"
                aria-label={t('auth.login')}
              >
                <Mail className="h-5 w-5" />
              </button>
              {authConfig?.enableGithubLogin && (
                <button
                  type="button"
                  onClick={() => {
                    closeWeChatMPModal();
                    void handleGithubLogin();
                  }}
                  className="flex h-12 w-12 items-center justify-center rounded-full border border-gray-200 text-gray-700 hover:bg-gray-50"
                  aria-label={t('auth.provider.github')}
                >
                  <Github className="h-5 w-5" />
                </button>
              )}
            </div>
          </div>

          <div className="mt-6 text-center text-xs text-gray-400">{t('auth.wechatMP.terms')}</div>
        </div>
      </Modal>
    </AuthLayout>
  );
};
