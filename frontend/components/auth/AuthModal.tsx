import React, { useEffect, useRef, useState } from 'react';
import { Mail, Lock, KeyRound, CheckCircle, X } from 'lucide-react';
import { useLocation, useNavigate } from 'react-router-dom';
import { AnimatePresence, LayoutGroup, motion } from 'framer-motion';
import { Button } from '../ui/Button';
import { Modal } from '../ui/Modal';
import { useLanguage } from '../../contexts/LanguageContext';
import { useAuth } from '../../contexts/AuthContext';
import { loginUser, registerUser, sendVerificationCode, verifyCode } from '../../services/authService';
import { getAuthConfig } from '../../services/configService';
import { consumeOtt, createWeChatMPScene, getWeChatMPSceneStatus } from '../../services/wechatMpAuthService';
import { AppRoute, AuthConfig } from '../../types';

type Mode = 'login' | 'register';
type View = 'main' | 'wechat_mp';

export const AuthModal: React.FC = () => {
  const { t } = useLanguage();
  const navigate = useNavigate();
  const location = useLocation();
  const { isAuthenticated, login, loginWithGithub, authModalOpen, authModalMode, authModalReturnTo, authModalSource, closeAuthModal } = useAuth();
  const [authConfig, setAuthConfig] = useState<AuthConfig | null>(null);
  const [booting, setBooting] = useState(false);
  const [mode, setMode] = useState<Mode>('login');
  const [view, setView] = useState<View>('main');

  const mpPollRef = useRef<number | null>(null);
  const [mpQrUrl, setMpQrUrl] = useState('');
  const [mpError, setMpError] = useState('');
  const [mpExpired, setMpExpired] = useState(false);

  const emailVerificationEnabled = !!(authConfig && authConfig.enableEmailVerification);
  const extraProvidersCount = (authConfig?.enableGithubLogin ? 1 : 0) + (authConfig?.enableWeChatMPLogin ? 1 : 0);

  const [loginForm, setLoginForm] = useState({ email: '', password: '' });
  const [loginLoading, setLoginLoading] = useState(false);
  const [loginError, setLoginError] = useState('');

  const [registerForm, setRegisterForm] = useState({ email: '', code: '', password: '', confirmPassword: '' });
  const [registerLoading, setRegisterLoading] = useState(false);
  const [registerError, setRegisterError] = useState('');
  const [countdown, setCountdown] = useState(0);

  const stopMpPolling = () => {
    if (mpPollRef.current) {
      window.clearInterval(mpPollRef.current);
      mpPollRef.current = null;
    }
  };

  useEffect(() => {
    return () => stopMpPolling();
  }, []);

  useEffect(() => {
    let timer: number;
    if (countdown > 0) {
      timer = window.setInterval(() => setCountdown((c) => c - 1), 1000);
    }
    return () => clearInterval(timer);
  }, [countdown]);

  const resetStates = (nextMode: Mode) => {
    setMode(nextMode);
    setView('main');
    stopMpPolling();
    setMpQrUrl('');
    setMpError('');
    setMpExpired(false);
    setLoginForm({ email: '', password: '' });
    setLoginLoading(false);
    setLoginError('');
    setRegisterForm({ email: '', code: '', password: '', confirmPassword: '' });
    setRegisterLoading(false);
    setRegisterError('');
    setCountdown(0);
  };

  const currentPath = location.pathname + location.search + location.hash;

  const finishAuth = async () => {
    const target = authModalReturnTo || currentPath || AppRoute.Home;
    closeAuthModal();
    if (target && target !== currentPath) {
      navigate(target, { replace: true });
    }
  };

  const handleClose = () => {
    stopMpPolling();
    closeAuthModal();
    if (!isAuthenticated && authModalSource === 'protected') {
      navigate(AppRoute.Home, { replace: true });
    }
  };

  const handleSwitchMode = (next: Mode) => {
    if (next === 'register') {
      setLoginForm({ email: '', password: '' });
      setLoginError('');
      setLoginLoading(false);
    } else {
      setRegisterForm({ email: '', code: '', password: '', confirmPassword: '' });
      setRegisterError('');
      setRegisterLoading(false);
      setCountdown(0);
    }
    setMode(next);
    setView('main');
  };

  const openWeChatMP = async () => {
    setLoginError('');
    setMpError('');
    setMpExpired(false);
    stopMpPolling();
    setMpQrUrl('');
    setView('wechat_mp');
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
            await finishAuth();
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

  useEffect(() => {
    if (!authModalOpen) return;
    const nextMode: Mode = authModalMode === 'register' ? 'register' : 'login';
    setBooting(nextMode === 'login');
    resetStates(nextMode);
    getAuthConfig().then((cfg) => {
      setAuthConfig(cfg);
      if (cfg.enableWeChatMPLogin && nextMode === 'login') {
        void openWeChatMP();
        setBooting(false);
        return;
      }
      setBooting(false);
    });
  }, [authModalOpen, authModalMode]);

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
        await finishAuth();
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
        await finishAuth();
      } else {
        setLoginError(t('auth.error.general'));
      }
    } catch {
      setLoginError(t('auth.error.general'));
    } finally {
      setLoginLoading(false);
    }
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
        await finishAuth();
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
    <Modal
      isOpen={authModalOpen}
      onClose={handleClose}
      hideHeader
      size="md"
      bodyClassName="p-0"
    >
      <div className="relative p-6">
        <button
          type="button"
          onClick={handleClose}
          aria-label={t('common.close')}
          className="absolute right-4 top-4 p-2 text-gray-400 hover:text-gray-600"
        >
          <X size={20} />
        </button>

        <LayoutGroup>
          <AnimatePresence initial={false} mode="wait">
            {booting && mode === 'login' ? (
              <motion.div
                key="view-loading"
                initial={{ opacity: 0, y: 8 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: 8 }}
                transition={{ duration: 0.16, ease: 'easeOut' }}
              >
                <div className="text-center text-2xl font-semibold text-gray-900">{t('auth.login')}</div>
                <div className="mt-3 flex items-center justify-center">
                  <div className="h-64 w-64 rounded-xl border border-gray-200 bg-gray-50" />
                </div>
              </motion.div>
            ) : view === 'wechat_mp' ? (
              <motion.div
                key="view-wechatmp"
                initial={{ opacity: 0, y: 8 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: 8 }}
                transition={{ duration: 0.16, ease: 'easeOut' }}
              >
                <div className="text-center text-2xl font-semibold text-gray-900">{t('auth.wechatMP.title')}</div>
                <div className="mt-3 text-center text-sm text-gray-600">{t('auth.wechatMP.desc')}</div>

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
                        onClick={() => void openWeChatMP()}
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
                      stopMpPolling();
                      setView('main');
                      handleSwitchMode('register');
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
                      onClick={() => {
                        stopMpPolling();
                        setView('main');
                        handleSwitchMode('login');
                      }}
                      className="flex h-12 w-12 items-center justify-center rounded-full border border-gray-200 text-gray-500 hover:bg-gray-50"
                      aria-label={t('auth.login')}
                    >
                      <Mail className="h-6 w-6" />
                    </button>
                    {authConfig?.enableGithubLogin && (
                      <button
                        type="button"
                        onClick={() => {
                          stopMpPolling();
                          void handleGithubLogin();
                        }}
                        className="flex h-12 w-12 items-center justify-center rounded-full border border-gray-200 text-gray-700 hover:bg-gray-50"
                        aria-label={t('auth.provider.github')}
                      >
                        <img src="/github.svg" alt={t('auth.provider.github')} className="h-7 w-7" />
                      </button>
                    )}
                  </div>
                </div>

                <div className="mt-6 text-center text-xs text-gray-400">{t('auth.wechatMP.terms')}</div>
              </motion.div>
            ) : (
              <motion.div
                key="view-main"
                initial={{ opacity: 0, y: 8 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: 8 }}
                transition={{ duration: 0.16, ease: 'easeOut' }}
              >
                <div className="text-center">
                  <AnimatePresence initial={false} mode="wait">
                    {mode === 'login' ? (
                      <motion.div
                        key="hdr-login"
                        initial={{ opacity: 0 }}
                        animate={{ opacity: 1 }}
                        exit={{ opacity: 0 }}
                        transition={{ duration: 0.16, ease: 'easeOut' }}
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
                      >
                        <h2 className="text-3xl font-extrabold leading-tight text-gray-900">{t('auth.createAccount')}</h2>
                        <p className="mt-2 text-sm leading-relaxed text-gray-600">{t('auth.createDesc')}</p>
                      </motion.div>
                    )}
                  </AnimatePresence>
                </div>

                <div className="mt-6 relative overflow-hidden">
                  <AnimatePresence initial={false} mode="wait">
                    {mode === 'login' ? (
                      <motion.form
                        key="login"
                        initial={{ x: '-100%', opacity: 0.9 }}
                        animate={{ x: 0, opacity: 1 }}
                        exit={{ x: '-100%', opacity: 0.9 }}
                        transition={{ duration: 0.18, ease: 'easeOut' }}
                        className="space-y-6 w-full"
                        onSubmit={handleLoginSubmit}
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
                            <div className="mt-6 flex items-center justify-center gap-6">
                              {authConfig?.enableGithubLogin && (
                                <button
                                  type="button"
                                  onClick={handleGithubLogin}
                                  disabled={loginLoading}
                                  aria-label={t('auth.provider.github')}
                                  className="flex h-12 w-12 items-center justify-center rounded-full border border-gray-200 text-gray-700 hover:bg-gray-50 disabled:opacity-60 disabled:cursor-not-allowed"
                                >
                                  <img src="/github.svg" alt={t('auth.provider.github')} className="h-7 w-7" />
                                </button>
                              )}
                              {authConfig?.enableWeChatMPLogin && (
                                <button
                                  type="button"
                                  onClick={() => void openWeChatMP()}
                                  disabled={loginLoading}
                                  aria-label={t('auth.provider.wechatMP')}
                                  className="flex h-12 w-12 items-center justify-center rounded-full border border-gray-200 text-gray-700 hover:bg-gray-50 disabled:opacity-60 disabled:cursor-not-allowed"
                                >
                                  <img src="/wechat.svg" alt={t('auth.provider.wechatMP')} className="h-7 w-7" />
                                </button>
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
                        initial={{ x: '100%', opacity: 0.9 }}
                        animate={{ x: 0, opacity: 1 }}
                        exit={{ x: '100%', opacity: 0.9 }}
                        transition={{ duration: 0.18, ease: 'easeOut' }}
                        className="space-y-6 w-full"
                        onSubmit={handleRegisterSubmit}
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
                                placeholder={t('auth.placeholder.password')}
                                value={registerForm.password}
                                onChange={(e) => setRegisterForm({ ...registerForm, password: e.target.value })}
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
                </div>
              </motion.div>
            )}
          </AnimatePresence>
        </LayoutGroup>
      </div>
    </Modal>
  );
};
