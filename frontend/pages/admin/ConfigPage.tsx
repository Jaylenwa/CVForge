import React, { useEffect, useMemo, useState } from 'react';
import { useLanguage } from '../../contexts/LanguageContext';
import { SystemConfig } from '../../types';
import { getSystemConfigs, updateSystemConfigs } from '../../services/configService';
import { Button } from '../../components/ui/Button';
import { Save, RefreshCw, ChevronDown } from 'lucide-react';
import { useToast } from '../../components/ui/Toast';

export const ConfigPage: React.FC = () => {
  const { t } = useLanguage();
  const [configs, setConfigs] = useState<SystemConfig[]>([]);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const { showToast } = useToast();
  const [activeTab, setActiveTab] = useState<string>('General');

  const getTabLabel = (name: string) => {
    const code = name.toLowerCase();
    return t(`admin.config.tab.${code}`);
  };

  const getConfigLabel = (key: string, description?: string) => {
    const map: Record<string, string> = {
      enable_email_verification: 'admin.config.key.enableEmailVerification',
      enabled_wechat_login: 'admin.config.key.enableWeChatLogin',
      enabled_github_login: 'admin.config.key.enableGithubLogin',
      wechat_app_id: 'admin.config.key.wechatAppId',
      wechat_appid: 'admin.config.key.wechatAppId',
      weChatAppID: 'admin.config.key.wechatAppId',
      wechat_app_secret: 'admin.config.key.wechatAppSecret',
      wechat_redirect_uri: 'admin.config.key.wechatRedirectUri',
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
      if (c.key === 'enable_email_verification' || c.key.startsWith('feature_') || c.key.startsWith('smtp_')) {
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

  return (
    <div className="flex-1 flex flex-col bg-white rounded-3xl m-2 overflow-hidden shadow-sm border border-gray-100">
      <div className="px-10 pt-10 pb-6">
        <div className="flex justify-between items-center">
          <h1 className="text-4xl font-bold text-gray-800">{t('admin.menu.settings')}</h1>
          <div className="flex space-x-2">
            <Button variant="outline" onClick={fetchConfigs} disabled={loading}>
              <RefreshCw size={16} className={`${loading ? 'animate-spin' : ''} mr-2`} /> {t('common.refresh') || 'Refresh'}
            </Button>
            <Button onClick={handleSave} isLoading={saving}>
              <Save size={16} className="mr-2" /> {t('common.save') || 'Save Changes'}
            </Button>
          </div>
        </div>
        <div className="flex border-b border-gray-100 overflow-x-auto mt-6">
          {Object.keys(groups).map(name => (
            <button
              key={name}
              onClick={() => setActiveTab(name)}
              className={`flex items-center space-x-2 px-4 pb-4 border-b-2 transition-all duration-200 whitespace-nowrap ${
                activeTab === name ? 'border-blue-600 text-blue-600' : 'border-transparent text-gray-500 hover:text-gray-700'
              }`}
            >
              <span className="text-[15px] font-medium">{getTabLabel(name)}</span>
            </button>
          ))}
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
                return items.filter(c => {
                  if (c.key === 'enable_email_verification') return true;
                  if (c.key.startsWith('smtp_')) return isOn;
                  return true;
                });
              }
              if (activeTab === 'OAuth') {
                const wechatVal = (configs.find(c => c.key === 'enabled_wechat_login')?.value || 'false');
                const githubVal = (configs.find(c => c.key === 'enabled_github_login')?.value || 'false');
                const wechatOn = wechatVal === 'true' || wechatVal === 'on';
                const githubOn = githubVal === 'true' || githubVal === 'on';
                return items.filter(c => {
                  if (c.key === 'enabled_wechat_login' || c.key === 'enabled_github_login') return true;
                  if (c.key === 'oauth_allowed_origins') return true;
                  if (c.key.startsWith('wechat_') || c.key.includes('wechat')) return wechatOn;
                  if (c.key.startsWith('github_') || c.key.includes('github')) return githubOn;
                  return true;
                });
              }
              return items;
            })()).map(config => {
              const isBool = config.type === 'bool' || config.value === 'true' || config.value === 'false' || config.value === 'on' || config.value === 'off';
              const enabled = config.value === 'true' || config.value === 'on';
              const indent =
                (activeTab === 'General' && config.key.startsWith('smtp_')) ||
                (activeTab === 'OAuth' &&
                  config.key !== 'enabled_wechat_login' &&
                  config.key !== 'enabled_github_login' &&
                  config.key !== 'oauth_allowed_origins' &&
                  (config.key.startsWith('wechat_') || config.key.includes('wechat') || config.key.startsWith('github_') || config.key.includes('github'))) ||
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
