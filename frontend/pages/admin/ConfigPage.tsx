import React, { useCallback, useEffect, useLayoutEffect, useMemo, useRef, useState } from 'react';
import { useLanguage } from '../../contexts/LanguageContext';
import { SystemConfig } from '../../types';
import { getSystemConfigs, updateSystemConfigs } from '../../services/configService';
import { Button } from '../../components/ui/Button';
import { Save } from 'lucide-react';
import { useToast } from '../../components/ui/Toast';
import { motion } from 'framer-motion';

export const ConfigPage: React.FC = () => {
  const { t } = useLanguage();
  const [configs, setConfigs] = useState<SystemConfig[]>([]);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const { showToast } = useToast();
  const [activeTab, setActiveTab] = useState<string>('General');
  const tabNavRef = useRef<HTMLDivElement | null>(null);
  const tabButtonRefs = useRef<Record<string, HTMLButtonElement | null>>({});
  const [tabIndicator, setTabIndicator] = useState<{ left: number; width: number }>({ left: 0, width: 0 });

  const getTabLabel = (name: string) => {
    const code = name.toLowerCase();
    return t(`admin.config.tab.${code}`);
  };

  const getConfigLabel = (key: string, description?: string) => {
    const map: Record<string, string> = {
      enable_email_verification: 'admin.config.key.enableEmailVerification',
      enabled_wechat_login: 'admin.config.key.enableWeChatLogin',
      enabled_wechat_mp_login: 'admin.config.key.enableWeChatMPLogin',
      enabled_github_login: 'admin.config.key.enableGithubLogin',
      wechat_app_id: 'admin.config.key.wechatAppId',
      wechat_appid: 'admin.config.key.wechatAppId',
      weChatAppID: 'admin.config.key.wechatAppId',
      wechat_app_secret: 'admin.config.key.wechatAppSecret',
      wechat_redirect_uri: 'admin.config.key.wechatRedirectUri',
      wechat_mp_app_id: 'admin.config.key.wechatMPAppId',
      wechat_mp_app_secret: 'admin.config.key.wechatMPAppSecret',
      wechat_mp_token: 'admin.config.key.wechatMPToken',
      wechat_mp_aes_key: 'admin.config.key.wechatMPAESKey',
      github_client_id: 'admin.config.key.githubClientId',
      github_clientid: 'admin.config.key.githubClientId',
      githubClientID: 'admin.config.key.githubClientId',
      github_client_secret: 'admin.config.key.githubClientSecret',
      github_redirect_uri: 'admin.config.key.githubRedirectUri',
      smtp_host: 'admin.config.key.smtpHost',
      smtp_port: 'admin.config.key.smtpPort',
      smtp_user: 'admin.config.key.smtpUser',
      smtp_username: 'admin.config.key.smtpUser',
      smtp_pass: 'admin.config.key.smtpPass',
      smtp_password: 'admin.config.key.smtpPass',
      smtp_from_name: 'admin.config.key.smtpFromName',
      smtp_secure: 'admin.config.key.smtpSecure',
      oauth_allowed_origins: 'admin.config.key.oauthAllowedOrigins',
      cors_origins: 'admin.config.key.corsOrigins',
      frontend_base_url: 'admin.config.key.frontendBaseUrl',
      enabled_pricing_page: 'admin.config.key.enablePricingPage',
      enabled_storage_s3: 'admin.config.key.storageS3Enabled',
      storage_s3_bucket: 'admin.config.key.storageS3Bucket',
      storage_s3_region: 'admin.config.key.storageS3Region',
      storage_s3_endpoint: 'admin.config.key.storageS3Endpoint',
      storage_s3_access_key: 'admin.config.key.storageS3AccessKey',
      storage_s3_secret_key: 'admin.config.key.storageS3SecretKey',
      chrome_api_url: 'admin.config.key.chromeApiUrl'
    };
    const keyName = map[key];
    if (keyName) return t(keyName);
    return description || key;
  };

  const fetchConfigs = async () => {
    setLoading(true);
    try {
      const data = await getSystemConfigs();
      setConfigs(data);
    } catch (err) {
      showToast(t('admin.config.msg.loadFailed'), 'error');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchConfigs();
  }, []);

  const handleSave = async () => {
    setSaving(true);
    try {
      await updateSystemConfigs(configs);
      showToast(t('admin.config.msg.saveSuccess'), 'success');
    } catch (err) {
      showToast(t('admin.config.msg.saveFailed'), 'error');
    } finally {
      setSaving(false);
    }
  };

  const handleChange = (key: string, value: string) => {
    setConfigs(configs.map(c => c.key === key ? { ...c, value } : c));
  };

  const groups = useMemo(() => {
    const g: { [key: string]: SystemConfig[] } = {
      'General': [],
      'OAuth': [],
      'Storage': [],
      'Security': [],
      'Chrome': []
    };
    configs.forEach(c => {
      if (c.key === 'enable_email_verification' || c.key === 'enabled_pricing_page' || c.key.startsWith('smtp_')) {
        g['General'].push(c);
      } else if (
        c.key === 'oauth_allowed_origins' ||
        c.key.startsWith('wechat_') || c.key.includes('wechat') ||
        c.key.startsWith('github_') || c.key.includes('github') ||
        c.key === 'enabled_wechat_login' || c.key === 'enabled_github_login'
      ) {
        g['OAuth'].push(c);
      } else if (c.key.startsWith('storage_') || c.key === 'enabled_storage_s3') {
        g['Storage'].push(c);
      } else if (c.key.startsWith('cors_') || c.key.startsWith('frontend_')) {
        g['Security'].push(c);
      } else if (c.key.startsWith('chrome_')) {
        g['Chrome'].push(c);
      }
    });
    return g;
  }, [configs]);

  const tabKeys = useMemo(() => Object.keys(groups), [groups]);

  const updateTabIndicator = useCallback(() => {
    const nav = tabNavRef.current;
    const btn = tabButtonRefs.current[activeTab];
    if (!nav || !btn) return;

    const navRect = nav.getBoundingClientRect();
    const btnRect = btn.getBoundingClientRect();
    const left = btnRect.left - navRect.left + nav.scrollLeft;
    const width = btnRect.width;
    setTabIndicator({ left, width });
  }, [activeTab]);

  useLayoutEffect(() => {
    updateTabIndicator();
  }, [updateTabIndicator, tabKeys.length]);

  useEffect(() => {
    const handler = () => updateTabIndicator();
    window.addEventListener('resize', handler);
    return () => window.removeEventListener('resize', handler);
  }, [updateTabIndicator]);

  useEffect(() => {
    const nav = tabNavRef.current;
    const btn = tabButtonRefs.current[activeTab];
    if (!nav || !btn || typeof ResizeObserver === 'undefined') return;
    const ro = new ResizeObserver(() => updateTabIndicator());
    ro.observe(nav);
    ro.observe(btn);
    return () => ro.disconnect();
  }, [activeTab, updateTabIndicator]);

  return (
    <div className="flex-1 flex flex-col bg-white rounded-3xl m-2 overflow-hidden shadow-sm border border-gray-100">
      <div className="px-10 pt-10 pb-6">
        <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4 border-b border-gray-100 pb-4">
          <nav
            ref={tabNavRef}
            className="relative inline-flex items-center gap-1 p-1 bg-gray-100 rounded-xl self-start overflow-x-auto overflow-y-visible no-scrollbar"
            onScroll={updateTabIndicator}
          >
            <motion.div
              className="absolute top-1 bottom-1 left-0 bg-white rounded-lg shadow-sm pointer-events-none"
              initial={false}
              animate={{ x: tabIndicator.left, width: tabIndicator.width, opacity: tabIndicator.width ? 1 : 0 }}
              transition={{ type: 'spring', stiffness: 450, damping: 40, mass: 0.4 }}
            />
            {tabKeys.map((name) => (
              <button
                key={name}
                type="button"
                aria-pressed={activeTab === name}
                ref={(el) => { tabButtonRefs.current[name] = el; }}
                className={`relative z-10 flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-semibold transition-all duration-200 whitespace-nowrap ${
                  activeTab === name ? 'bg-white text-blue-600 shadow-sm' : 'text-gray-500 hover:text-gray-800 hover:bg-gray-200/60'
                }`}
                onClick={() => setActiveTab(name)}
              >
                {getTabLabel(name)}
              </button>
            ))}
          </nav>
          <div className="flex items-center gap-2 md:justify-end">
            <Button onClick={handleSave} isLoading={saving}>
              <Save size={16} className="mr-2" /> {t('common.save') || 'Save Changes'}
            </Button>
          </div>
        </div>
      </div>
      <div className="flex-1 overflow-y-auto px-10 pb-10">
        <div className="max-w-3xl space-y-8">
          {loading && configs.length === 0 ? (
            <div className="text-center py-12">{t('common.loading') || 'Loading...'}</div>
          ) : (
            ((() => {
              const items = groups[activeTab] || [];
              if (activeTab === 'Storage') {
                const s3Enabled = (configs.find(c => c.key === 'enabled_storage_s3')?.value || 'false');
                const isOn = s3Enabled === 'true' || s3Enabled === 'on';
                return items.filter(c => {
                  if (c.key === 'enabled_storage_s3') return true;
                  if (c.key.startsWith('storage_s3_')) return isOn;
                  return true;
                });
              }
              if (activeTab === 'General') {
                const ev = (configs.find(c => c.key === 'enable_email_verification')?.value || 'false');
                const isOn = ev === 'true' || ev === 'on';
                const visible = items.filter(c => {
                  if (c.key === 'enable_email_verification') return true;
                  if (c.key.startsWith('smtp_')) return isOn;
                  return true;
                });
                const smtpOrder: Record<string, number> = {
                  smtp_host: 0,
                  smtp_port: 1,
                  smtp_user: 2,
                  smtp_username: 2,
                  smtp_pass: 3,
                  smtp_password: 3,
                  smtp_from_name: 4,
                  smtp_secure: 5,
                };
                return visible
                  .map((c, idx) => ({ c, idx }))
                  .sort((a, b) => {
                    const rank = (k: string) => {
                      if (k === 'enable_email_verification') return 0;
                      if (k.startsWith('smtp_')) return 1;
                      if (k === 'enabled_pricing_page') return 2;
                      return 3;
                    };
                    const ra = rank(a.c.key);
                    const rb = rank(b.c.key);
                    if (ra !== rb) return ra - rb;
                    if (ra === 1) {
                      const oa = smtpOrder[a.c.key] ?? 999;
                      const ob = smtpOrder[b.c.key] ?? 999;
                      if (oa !== ob) return oa - ob;
                    }
                    return a.idx - b.idx;
                  })
                  .map(x => x.c);
              }
              if (activeTab === 'OAuth') {
                const wechatVal = (configs.find(c => c.key === 'enabled_wechat_login')?.value || 'false');
                const wechatMPVal = (configs.find(c => c.key === 'enabled_wechat_mp_login')?.value || 'false');
                const githubVal = (configs.find(c => c.key === 'enabled_github_login')?.value || 'false');
                const wechatOn = wechatVal === 'true' || wechatVal === 'on';
                const wechatMPOn = wechatMPVal === 'true' || wechatMPVal === 'on';
                const githubOn = githubVal === 'true' || githubVal === 'on';
                const isWeChatMPKey = (k: string) => k === 'enabled_wechat_mp_login' || k.startsWith('wechat_mp_') || k.includes('wechat_mp');
                const isWeChatOAuthKey = (k: string) => (k.startsWith('wechat_') || k.includes('wechat')) && !isWeChatMPKey(k) && k !== 'enabled_wechat_login';
                const isGithubKey = (k: string) => (k.startsWith('github_') || k.includes('github')) && k !== 'enabled_github_login';

                const filtered = items.filter(c => {
                  if (c.key === 'enabled_wechat_login' || c.key === 'enabled_wechat_mp_login' || c.key === 'enabled_github_login') return true;
                  if (c.key === 'oauth_allowed_origins') return true;
                  if (isWeChatMPKey(c.key)) return wechatMPOn;
                  if (isWeChatOAuthKey(c.key)) return wechatOn;
                  if (isGithubKey(c.key)) return githubOn;
                  return true;
                });

                const idx: Record<string, number> = {};
                items.forEach((c, i) => { idx[c.key] = i; });
                const groupRank = (k: string) => {
                  if (k === 'enabled_wechat_login') return 0;
                  if (isWeChatOAuthKey(k)) return 1;
                  if (k === 'enabled_wechat_mp_login') return 2;
                  if (isWeChatMPKey(k)) return 3;
                  if (k === 'enabled_github_login') return 4;
                  if (isGithubKey(k)) return 5;
                  if (k === 'oauth_allowed_origins') return 6;
                  return 7;
                };
                const fieldRank: Record<string, number> = {
                  wechat_app_id: 0,
                  wechat_app_secret: 1,
                  wechat_redirect_uri: 2,
                  wechat_mp_app_id: 0,
                  wechat_mp_app_secret: 1,
                  wechat_mp_token: 2,
                  wechat_mp_aes_key: 3,
                  github_client_id: 0,
                  github_client_secret: 1,
                  github_redirect_uri: 2,
                };
                return filtered
                  .map((c) => ({ c, i: idx[c.key] ?? 9999 }))
                  .sort((a, b) => {
                    const ra = groupRank(a.c.key);
                    const rb = groupRank(b.c.key);
                    if (ra !== rb) return ra - rb;
                    const fa = fieldRank[a.c.key];
                    const fb = fieldRank[b.c.key];
                    if (fa != null && fb != null && fa !== fb) return fa - fb;
                    return a.i - b.i;
                  })
                  .map(x => x.c);
              }
              return items;
            })()).map(config => {
              const isBool = config.type === 'bool' || config.value === 'true' || config.value === 'false' || config.value === 'on' || config.value === 'off';
              const enabled = config.value === 'true' || config.value === 'on';
              const indent =
                (activeTab === 'General' && config.key.startsWith('smtp_')) ||
                (activeTab === 'OAuth' &&
                  config.key !== 'enabled_wechat_login' &&
                  config.key !== 'enabled_wechat_mp_login' &&
                  config.key !== 'enabled_github_login' &&
                  config.key !== 'oauth_allowed_origins' &&
                  ((config.key.startsWith('wechat_mp_') || config.key.includes('wechat_mp')) ||
                    (config.key.startsWith('wechat_') || config.key.includes('wechat')) ||
                    (config.key.startsWith('github_') || config.key.includes('github')))) ||
                (activeTab === 'Storage' && config.key.startsWith('storage_s3_'));
              return (
                <div key={config.key} className={`space-y-2 ${indent ? 'ml-8' : ''}`}>
                  {isBool ? (
                    <div className="flex flex-col space-y-1 py-2">
                      <div className="flex items-center space-x-4">
                        <button
                          onClick={() => handleChange(config.key, enabled ? 'false' : 'true')}
                          className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors duration-200 focus:outline-none ${
                            enabled ? 'bg-blue-600' : 'bg-gray-300'
                          }`}
                        >
                          <span
                            className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform duration-200 ${
                              enabled ? 'translate-x-6' : 'translate-x-1'
                            }`}
                          />
                        </button>
                        <span className="text-gray-800 text-[15px] font-medium">{getConfigLabel(config.key, config.description)}</span>
                      </div>
                    </div>
                  ) : (
                    <div className="space-y-2">
                      <label className="text-gray-800 text-[14px] font-bold block">{getConfigLabel(config.key, config.description)}</label>
                      <div className="relative group">
                        <input
                          className="w-full h-11 px-4 bg-white border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 text-gray-800"
                          value={config.value}
                          onChange={(e) => handleChange(config.key, e.target.value)}
                        />
                      </div>
                    </div>
                  )}
                </div>
              );
            })
          )}
        </div>
      </div>
    </div>
  );
};
