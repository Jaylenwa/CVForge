import React, { useEffect, useState } from 'react';
import { useLanguage } from '../../contexts/LanguageContext';
import { SystemConfig } from '../../types';
import { getSystemConfigs, updateSystemConfigs } from '../../services/configService';
import { Button } from '../../components/ui/Button';
import { Save, RefreshCw } from 'lucide-react';
import { useToast } from '../../components/ui/Toast';

export const ConfigPage: React.FC = () => {
  const { t } = useLanguage();
  const [configs, setConfigs] = useState<SystemConfig[]>([]);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const { showToast } = useToast();

  const fetchConfigs = async () => {
    setLoading(true);
    try {
      const data = await getSystemConfigs();
      setConfigs(data);
    } catch (err) {
      showToast('Failed to load configs', 'error');
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
      showToast('Configs saved successfully', 'success');
    } catch (err) {
      showToast('Failed to save configs', 'error');
    } finally {
      setSaving(false);
    }
  };

  const handleChange = (key: string, value: string) => {
    setConfigs(configs.map(c => c.key === key ? { ...c, value } : c));
  };

  const groupConfigs = () => {
    const groups: { [key: string]: SystemConfig[] } = {
      'General': [],
      'SMTP': [],
      'WeChat': [],
      'GitHub': [],
      'Other': []
    };

    configs.forEach(c => {
      if (c.key.startsWith('smtp_')) groups['SMTP'].push(c);
      else if (c.key.startsWith('wechat_') || c.key.includes('wechat')) groups['WeChat'].push(c);
      else if (c.key.startsWith('github_') || c.key.includes('github')) groups['GitHub'].push(c);
      else if (c.key === 'enable_email_verification') groups['General'].push(c);
      else groups['Other'].push(c);
    });

    return groups;
  };

  const groups = groupConfigs();

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h2 className="text-xl font-semibold text-gray-800">System Settings</h2>
        <div className="flex space-x-2">
            <Button variant="outline" onClick={fetchConfigs} disabled={loading}>
                <RefreshCw size={16} className={`mr-2 ${loading ? 'animate-spin' : ''}`} /> Refresh
            </Button>
            <Button onClick={handleSave} isLoading={saving}>
                <Save size={16} className="mr-2" /> Save Changes
            </Button>
        </div>
      </div>

      {loading && configs.length === 0 ? (
        <div className="text-center py-12">Loading...</div>
      ) : (
        <div className="space-y-8">
          {Object.entries(groups).map(([groupName, groupConfigs]) => (
            groupConfigs.length > 0 && (
              <div key={groupName} className="bg-white shadow rounded-lg p-6">
                <h3 className="text-lg font-medium text-gray-900 mb-4 border-b pb-2">{groupName}</h3>
                <div className="grid grid-cols-1 gap-y-6 gap-x-4 sm:grid-cols-6">
                  {groupConfigs.map(config => (
                    <div key={config.key} className="sm:col-span-6">
                      <label htmlFor={config.key} className="block text-sm font-medium text-gray-700">
                        {config.description || config.key}
                      </label>
                      <div className="mt-1">
                        {config.type === 'bool' || config.value === 'true' || config.value === 'false' || config.value === 'on' || config.value === 'off' ? (
                           <select
                              id={config.key}
                              value={config.value}
                              onChange={(e) => handleChange(config.key, e.target.value)}
                              className="shadow-sm focus:ring-blue-500 focus:border-blue-500 block w-full sm:text-sm border-gray-300 rounded-md"
                           >
                               <option value="true">True / On</option>
                               <option value="false">False / Off</option>
                               <option value="on">On</option>
                               <option value="off">Off</option>
                           </select>
                        ) : (
                            <input
                            type="text"
                            name={config.key}
                            id={config.key}
                            value={config.value}
                            onChange={(e) => handleChange(config.key, e.target.value)}
                            className="shadow-sm focus:ring-blue-500 focus:border-blue-500 block w-full sm:text-sm border-gray-300 rounded-md p-2 border"
                            />
                        )}
                      </div>
                      <p className="mt-1 text-xs text-gray-500">{config.key}</p>
                    </div>
                  ))}
                </div>
              </div>
            )
          ))}
        </div>
      )}
    </div>
  );
};
