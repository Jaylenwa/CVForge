import React, { useState, useRef } from 'react';
import { User, Lock, Bell, Globe, Camera, Save } from 'lucide-react';
import { useAuth } from '../contexts/AuthContext';
import { useLanguage } from '../contexts/LanguageContext';
import { Button } from '../components/ui/Button';

export const Settings: React.FC = () => {
  const { user, login } = useAuth(); // In a real app, we'd have a specific updateProfile method
  const { t, language, setLanguage } = useLanguage();
  const fileInputRef = useRef<HTMLInputElement>(null);
  
  const [activeTab, setActiveTab] = useState<'profile' | 'security'>('profile');
  const [isLoading, setIsLoading] = useState(false);
  const [successMsg, setSuccessMsg] = useState('');

  // Local state for form
  const [formData, setFormData] = useState({
    name: user?.name || '',
    email: user?.email || '',
    avatarUrl: user?.avatarUrl || '',
    currentPassword: '',
    newPassword: '',
    confirmPassword: ''
  });

  const handleAvatarUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      const reader = new FileReader();
      reader.onloadend = () => {
        setFormData(prev => ({ ...prev, avatarUrl: reader.result as string }));
      };
      reader.readAsDataURL(file);
    }
  };

  const handleSaveProfile = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 800));
    
    // Update local context (mock)
    // We are re-using the login function to mock updating the user state
    // In a real app, useAuth would have an updateProfile function
    // For now, we simulate success
    setIsLoading(false);
    setSuccessMsg(t('settings.success.profileUpdated'));
    setTimeout(() => setSuccessMsg(''), 3000);
  };

  const handleSavePassword = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    await new Promise(resolve => setTimeout(resolve, 800));
    setIsLoading(false);
    setFormData(prev => ({ ...prev, currentPassword: '', newPassword: '', confirmPassword: '' }));
    setSuccessMsg(t('settings.success.passwordUpdated'));
    setTimeout(() => setSuccessMsg(''), 3000);
  };

  return (
    <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
      <h1 className="text-3xl font-bold text-gray-900 mb-8">{t('nav.settings')}</h1>

      <div className="bg-white rounded-lg shadow border border-gray-200 overflow-hidden flex flex-col md:flex-row min-h-[500px]">
        {/* Sidebar */}
        <div className="w-full md:w-64 bg-gray-50 border-r border-gray-200">
          <nav className="p-4 space-y-2">
            <button
              onClick={() => setActiveTab('profile')}
              className={`w-full flex items-center px-4 py-3 text-sm font-medium rounded-md transition-colors ${
                activeTab === 'profile' 
                  ? 'bg-blue-50 text-blue-700' 
                  : 'text-gray-700 hover:bg-gray-100'
              }`}
            >
              <User size={18} className="mr-3" />
              {t('nav.profile')}
            </button>
            <button
              onClick={() => setActiveTab('security')}
              className={`w-full flex items-center px-4 py-3 text-sm font-medium rounded-md transition-colors ${
                activeTab === 'security' 
                  ? 'bg-blue-50 text-blue-700' 
                  : 'text-gray-700 hover:bg-gray-100'
              }`}
            >
              <Lock size={18} className="mr-3" />
              {t('settings.security')}
            </button>
          </nav>
        </div>

        {/* Content */}
        <div className="flex-1 p-8">
          {successMsg && (
             <div className="mb-6 bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded relative">
                <span className="block sm:inline">{successMsg}</span>
            </div>
          )}

          {activeTab === 'profile' && (
            <form onSubmit={handleSaveProfile} className="space-y-6 max-w-lg">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">{t('settings.profile.photoLabel')}</label>
                <div className="flex items-center space-x-6">
                  <div className="relative">
                     <img 
                        src={formData.avatarUrl || user?.avatarUrl} 
                        alt="Profile" 
                        className="h-24 w-24 rounded-full object-cover border-4 border-white shadow-sm"
                     />
                     <button
                        type="button" 
                        onClick={() => fileInputRef.current?.click()}
                        className="absolute bottom-0 right-0 bg-blue-600 p-1.5 rounded-full text-white hover:bg-blue-700 shadow-sm border-2 border-white"
                     >
                        <Camera size={14} />
                     </button>
                  </div>
                  <div>
                    <Button type="button" variant="outline" size="sm" onClick={() => fileInputRef.current?.click()}>{t('settings.profile.changePhoto')}</Button>
                    <input 
                        ref={fileInputRef}
                        type="file" 
                        className="hidden" 
                        accept="image/*" 
                        onChange={handleAvatarUpload}
                    />
                    <p className="text-xs text-gray-500 mt-2">{t('settings.profile.photoTip')}</p>
                  </div>
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700">{t('settings.displayName')}</label>
                <input 
                  type="text" 
                  value={formData.name}
                  onChange={(e) => setFormData({...formData, name: e.target.value})}
                  className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2 focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700">{t('settings.emailAddress')}</label>
                <input 
                  type="email" 
                  value={formData.email}
                  disabled
                  className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2 bg-gray-50 text-gray-500 sm:text-sm cursor-not-allowed"
                />
                <p className="mt-1 text-xs text-gray-500">{t('settings.emailChangeTip')}</p>
              </div>
              
              <div className="pt-4 border-t border-gray-200">
                  <label className="block text-sm font-medium text-gray-700 mb-3">{t('settings.languagePreference')}</label>
                  <div className="flex items-center space-x-4">
                      <button 
                        type="button"
                        onClick={() => setLanguage('en')}
                        className={`px-4 py-2 border rounded-md text-sm flex items-center ${language === 'en' ? 'border-blue-500 bg-blue-50 text-blue-700' : 'border-gray-300 hover:bg-gray-50'}`}
                      >
                          <span className="mr-2">🇺🇸</span> {t('lang.en')}
                      </button>
                      <button 
                        type="button"
                        onClick={() => setLanguage('zh')}
                        className={`px-4 py-2 border rounded-md text-sm flex items-center ${language === 'zh' ? 'border-blue-500 bg-blue-50 text-blue-700' : 'border-gray-300 hover:bg-gray-50'}`}
                      >
                          <span className="mr-2">🇨🇳</span> {t('lang.zh')}
                      </button>
                  </div>
              </div>

              <div className="pt-4">
                <Button type="submit" isLoading={isLoading} icon={<Save size={16}/>}>{t('settings.saveChanges')}</Button>
              </div>
            </form>
          )}

          {activeTab === 'security' && (
            <form onSubmit={handleSavePassword} className="space-y-6 max-w-lg">
              <div>
                <label className="block text-sm font-medium text-gray-700">{t('settings.currentPassword')}</label>
                <input 
                  type="password" 
                  value={formData.currentPassword}
                  onChange={(e) => setFormData({...formData, currentPassword: e.target.value})}
                  className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2 focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700">{t('settings.newPassword')}</label>
                <input 
                  type="password" 
                  value={formData.newPassword}
                  onChange={(e) => setFormData({...formData, newPassword: e.target.value})}
                  className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2 focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700">{t('settings.confirmNewPassword')}</label>
                <input 
                  type="password" 
                  value={formData.confirmPassword}
                  onChange={(e) => setFormData({...formData, confirmPassword: e.target.value})}
                  className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2 focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                />
              </div>

              <div className="pt-4">
                <Button type="submit" isLoading={isLoading} icon={<Save size={16}/>}>{t('settings.updatePassword')}</Button>
              </div>
            </form>
          )}
        </div>
      </div>
    </div>
  );
};
