import React, { useState, useRef } from 'react';
import { Trash2, Plus, GripVertical, Sparkles, ChevronDown, ChevronUp, Upload, X, Image as ImageIcon, Palette, Type, LayoutTemplate } from 'lucide-react';
import { ResumeData, ResumeSection, ResumeItem, ResumeSectionType, ThemeConfig } from '../../types';
import { Button } from '../../components/ui/Button';
import { polishText, generateSummary } from '../../services/geminiService';
import { useLanguage } from '../../contexts/LanguageContext';

interface EditorFormProps {
  data: ResumeData;
  onChange: (data: ResumeData) => void;
}

const FONT_OPTIONS = [
  { group: 'English - Sans Serif', id: 'inter', label: 'Inter (Modern)' },
  { group: 'English - Sans Serif', id: 'roboto', label: 'Roboto (Technical)' },
  { group: 'English - Serif', id: 'merriweather', label: 'Merriweather (Elegant)' },
  { group: 'English - Serif', id: 'playfair', label: 'Playfair Display (Classy)' },
  { group: 'English - Mono', id: 'mono', label: 'Roboto Mono (Code)' },
  { group: 'Chinese - Sans (黑体)', id: 'yahei', label: '微软雅黑 (Microsoft YaHei)' },
  { group: 'Chinese - Sans (黑体)', id: 'notosans', label: '思源黑体 (Noto Sans SC)' },
  { group: 'Chinese - Serif (宋体/楷体)', id: 'simsun', label: '宋体 (SimSun)' },
  { group: 'Chinese - Serif (宋体/楷体)', id: 'kaiti', label: '楷体 (KaiTi)' },
];

export const EditorForm: React.FC<EditorFormProps> = ({ data, onChange }) => {
  const { t } = useLanguage();
  const [activeTab, setActiveTab] = useState<'content' | 'design'>('content');
  const [activeSection, setActiveSection] = useState<string | null>('personal');
  const [isAiLoading, setIsAiLoading] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const updatePersonalInfo = (field: string, value: string) => {
    onChange({
      ...data,
      personalInfo: { ...data.personalInfo, [field]: value }
    });
  };

  const updateTheme = (field: keyof ThemeConfig, value: string) => {
    onChange({
        ...data,
        themeConfig: { ...data.themeConfig, [field]: value }
    });
  };

  const handleAvatarUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;
    if (file.size > 2 * 1024 * 1024) {
        alert(t('editor.error.imageTooLarge'));
        return;
    }
    const form = new FormData();
    form.append('file', file);
    try {
      const res = await fetch('http://localhost:8080/api/v1/upload/avatar', {
        method: 'POST',
        body: form
      });
      if (res.ok) {
        const data = await res.json();
        updatePersonalInfo('avatarUrl', data.url);
      }
    } catch {}
  };

  const updateSection = (sectionId: string, updates: Partial<ResumeSection>) => {
    onChange({
      ...data,
      sections: data.sections.map(s => s.id === sectionId ? { ...s, ...updates } : s)
    });
  };

  const updateItem = (sectionId: string, itemId: string, field: keyof ResumeItem, value: string) => {
    onChange({
      ...data,
      sections: data.sections.map(s => {
        if (s.id !== sectionId) return s;
        return {
          ...s,
          items: s.items.map(item => item.id === itemId ? { ...item, [field]: value } : item)
        };
      })
    });
  };

  const addItem = (sectionId: string) => {
    const newItem: ResumeItem = {
      id: Math.random().toString(36).substr(2, 9),
      title: '',
      description: ''
    };
    onChange({
      ...data,
      sections: data.sections.map(s => 
        s.id === sectionId ? { ...s, items: [...s.items, newItem] } : s
      )
    });
  };

  const removeItem = (sectionId: string, itemId: string) => {
    onChange({
      ...data,
      sections: data.sections.map(s => 
        s.id === sectionId ? { ...s, items: s.items.filter(i => i.id !== itemId) } : s
      )
    });
  };

  const handleAiPolish = async (sectionId: string, itemId: string, text: string) => {
    if (!text) return;
    setIsAiLoading(true);
    const improved = await polishText(text);
    updateItem(sectionId, itemId, 'description', improved);
    setIsAiLoading(false);
  };

  const handleAiSummary = async () => {
     setIsAiLoading(true);
     // Find relevant skills for context
     const skillsSection = data.sections.find(s => s.type === ResumeSectionType.Skills);
     const skills = skillsSection?.items.map(i => i.description).join(', ') || 'General';
     
     const summary = await generateSummary(data.personalInfo.jobTitle, skills);
     
     // Update summary section
     const summarySection = data.sections.find(s => s.type === ResumeSectionType.Summary);
     if (summarySection && summarySection.items.length > 0) {
         updateItem(summarySection.id, summarySection.items[0].id, 'description', summary);
     }
     setIsAiLoading(false);
  }

  const renderContentTab = () => (
    <>
      {/* Personal Info */}
      <div className="border-b border-gray-200">
        <button 
          className="w-full px-6 py-4 flex items-center justify-between font-semibold hover:bg-gray-50 text-gray-700"
          onClick={() => setActiveSection(activeSection === 'personal' ? null : 'personal')}
        >
            <span>{t('editor.personal')}</span>
            {activeSection === 'personal' ? <ChevronUp size={16}/> : <ChevronDown size={16}/>}
        </button>
        
        {activeSection === 'personal' && (
          <div className="px-6 pb-6 space-y-5 animate-fadeIn">
             
             {/* Avatar Upload */}
             <div className="bg-gray-50 p-4 rounded-lg border border-gray-200">
                 <label className="block text-sm font-medium text-gray-700 mb-3 flex items-center">
                     <ImageIcon size={16} className="mr-2 text-gray-500" /> 
                    {t('editor.profile.label')}
                 </label>
                 <div className="flex items-center space-x-4">
                    {data.personalInfo.avatarUrl ? (
                        <div className="relative group shrink-0">
                            <img 
                                src={data.personalInfo.avatarUrl} 
                                alt={t('a11y.avatarAlt')} 
                                className="w-16 h-16 rounded-full object-cover border-2 border-white shadow-sm"
                            />
                            <button
                                onClick={() => updatePersonalInfo('avatarUrl', '')}
                                className="absolute -top-1 -right-1 bg-red-500 text-white rounded-full p-1 shadow-sm hover:bg-red-600 transition-colors"
                                title={t('a11y.removePhoto')}
                            >
                                <X size={10} />
                            </button>
                        </div>
                    ) : (
                        <div className="w-16 h-16 shrink-0 rounded-full bg-gray-100 flex items-center justify-center text-gray-400 border-2 border-gray-200 border-dashed">
                            <Upload size={20} />
                        </div>
                    )}
                    
                    <div className="flex-1">
                        <input 
                            ref={fileInputRef}
                            type="file" 
                            className="hidden" 
                            accept="image/png, image/jpeg, image/jpg" 
                            onChange={handleAvatarUpload} 
                        />
                        <Button 
                            variant="outline" 
                            size="sm"
                            onClick={() => fileInputRef.current?.click()}
                            className="w-full"
                        >
                            {data.personalInfo.avatarUrl ? t('editor.profile.change') : t('editor.profile.upload')}
                        </Button>
                        <p className="mt-2 text-xs text-gray-500">
                            {t('editor.profile.tipSquare')}<br/>
                            {t('editor.profile.tipMaxSize')}
                        </p>
                    </div>
                 </div>
             </div>

            <div>
              <label className="block text-sm font-medium text-gray-700">{t('editor.fields.fullName')}</label>
              <input 
                type="text" 
                value={data.personalInfo.fullName}
                onChange={e => updatePersonalInfo('fullName', e.target.value)}
                className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2 focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700">{t('editor.fields.jobTitle')}</label>
              <input 
                type="text" 
                value={data.personalInfo.jobTitle}
                onChange={e => updatePersonalInfo('jobTitle', e.target.value)}
                className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2 focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
              />
            </div>
             <div className="grid grid-cols-2 gap-4">
                <div>
                    <label className="block text-sm font-medium text-gray-700">{t('editor.fields.email')}</label>
                    <input type="email" value={data.personalInfo.email} onChange={e => updatePersonalInfo('email', e.target.value)} className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2 sm:text-sm"/>
                </div>
                <div>
                    <label className="block text-sm font-medium text-gray-700">{t('editor.fields.phone')}</label>
                    <input type="tel" value={data.personalInfo.phone} onChange={e => updatePersonalInfo('phone', e.target.value)} className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2 sm:text-sm"/>
                </div>
             </div>
             <div>
                <label className="block text-sm font-medium text-gray-700">{t('editor.fields.address')}</label>
                <input type="text" value={data.personalInfo.address} onChange={e => updatePersonalInfo('address', e.target.value)} className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2 sm:text-sm"/>
             </div>
             <div>
                <label className="block text-sm font-medium text-gray-700">{t('editor.fields.website')}</label>
                <input type="text" value={data.personalInfo.website} onChange={e => updatePersonalInfo('website', e.target.value)} className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2 sm:text-sm"/>
             </div>
          </div>
        )}
      </div>

      {/* Dynamic Sections */}
      {data.sections.map((section) => (
        <div key={section.id} className="border-b border-gray-200">
           <div className="flex items-center justify-between px-6 py-4 hover:bg-gray-50">
             <button 
                className="flex-grow text-left font-semibold flex items-center text-gray-700"
                onClick={() => setActiveSection(activeSection === section.id ? null : section.id)}
             >
                 <span className="mr-2 cursor-grab active:cursor-grabbing"><GripVertical size={16} className="text-gray-400"/></span>
                 {section.title}
             </button>
             <div className="flex items-center space-x-2">
                 {/* Toggle Visibility */}
                 <input 
                    type="checkbox" 
                    checked={section.isVisible} 
                    onChange={(e) => updateSection(section.id, { isVisible: e.target.checked })}
                    className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                 />
                  {activeSection === section.id ? <ChevronUp size={16}/> : <ChevronDown size={16}/>}
             </div>
           </div>

           {activeSection === section.id && (
             <div className="px-6 pb-6 space-y-6 bg-gray-50/50">
                {/* AI Summary Generator specific to Summary Section */}
                {section.type === ResumeSectionType.Summary && (
                    <div className="flex justify-end">
                         <Button 
                            size="sm" 
                            variant="secondary" 
                            icon={<Sparkles size={14}/>} 
                            onClick={handleAiSummary}
                            isLoading={isAiLoading}
                            className="mb-2"
                        >
                            {t('editor.ai_polish')}
                        </Button>
                    </div>
                )}

                {section.items.map((item, index) => (
                   <div key={item.id} className="bg-white p-4 rounded-lg shadow-sm border border-gray-200 relative group">
                      <button 
                        onClick={() => removeItem(section.id, item.id)}
                        className="absolute top-2 right-2 text-gray-400 hover:text-red-500 opacity-0 group-hover:opacity-100 transition-opacity"
                      >
                          <Trash2 size={16} />
                      </button>
                      
                      {section.type !== ResumeSectionType.Skills && section.type !== ResumeSectionType.Summary && (
                          <div className="grid grid-cols-1 gap-4 mb-3">
                              <input 
                                placeholder={t('editor.placeholder.titleExample')} 
                                className="w-full font-medium border-b border-gray-200 focus:border-blue-500 outline-none pb-1"
                                value={item.title}
                                onChange={e => updateItem(section.id, item.id, 'title', e.target.value)}
                              />
                              <div className="grid grid-cols-2 gap-4">
                                  <input 
                                    placeholder={t('editor.placeholder.subtitleExample')} 
                                    className="w-full text-sm border-b border-gray-200 focus:border-blue-500 outline-none pb-1"
                                    value={item.subtitle}
                                    onChange={e => updateItem(section.id, item.id, 'subtitle', e.target.value)}
                                  />
                                   <input 
                                    placeholder={t('editor.placeholder.dateRange')} 
                                    className="w-full text-sm border-b border-gray-200 focus:border-blue-500 outline-none pb-1"
                                    value={item.dateRange}
                                    onChange={e => updateItem(section.id, item.id, 'dateRange', e.target.value)}
                                  />
                              </div>
                          </div>
                      )}

                      <div className="relative">
                          <textarea 
                            rows={section.type === ResumeSectionType.Skills ? 2 : 4}
                            className="w-full text-sm border border-gray-300 rounded p-2 focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
                            placeholder={section.type === ResumeSectionType.Skills ? t('editor.placeholder.skills') : t('editor.placeholder.achievements')}
                            value={item.description}
                            onChange={e => updateItem(section.id, item.id, 'description', e.target.value)}
                          />
                          {section.type === ResumeSectionType.Experience && (
                              <button 
                                onClick={() => handleAiPolish(section.id, item.id, item.description)}
                                className="absolute bottom-2 right-2 text-xs bg-blue-100 text-blue-700 px-2 py-1 rounded flex items-center hover:bg-blue-200 transition-colors"
                                disabled={isAiLoading}
                              >
                                  <Sparkles size={12} className="mr-1"/> 
                                  {isAiLoading ? t('editor.ai.polishing') : t('editor.ai_polish')}
                              </button>
                          )}
                      </div>
                   </div>
                ))}

                <Button variant="outline" size="sm" className="w-full border-dashed" onClick={() => addItem(section.id)}>
                    <Plus size={16} className="mr-1"/> {t('editor.addItem')}
                </Button>
             </div>
           )}
        </div>
      ))}
    </>
  );

  const predefinedColors = ['#2563eb', '#0f172a', '#dc2626', '#16a34a', '#9333ea', '#ea580c', '#0891b2'];
  const isCustomColor = !predefinedColors.includes(data.themeConfig?.color || '');

  const renderDesignTab = () => (
    <div className="p-6 space-y-8">
        {/* Colors */}
        <div>
            <h3 className="text-sm font-semibold text-gray-900 mb-4 flex items-center">
                <Palette size={18} className="mr-2 text-gray-500"/> {t('editor.color')}
            </h3>
            <div className="flex flex-wrap gap-3 items-center">
                {predefinedColors.map(color => (
                    <button
                        key={color}
                        onClick={() => updateTheme('color', color)}
                        className={`w-8 h-8 rounded-full border-2 transition-transform hover:scale-110 ${data.themeConfig?.color === color ? 'border-gray-900 scale-110 ring-2 ring-gray-200' : 'border-transparent'}`}
                        style={{ backgroundColor: color }}
                    />
                ))}

                {/* Custom Color Picker */}
                <div className="relative group">
                    <div 
                        className={`w-8 h-8 rounded-full border-2 flex items-center justify-center overflow-hidden transition-transform group-hover:scale-110 cursor-pointer ${
                            isCustomColor 
                            ? 'border-gray-900 scale-110 ring-2 ring-gray-200' 
                            : 'border-gray-300'
                        }`}
                        style={{ 
                            background: isCustomColor 
                                ? data.themeConfig?.color 
                                : 'conic-gradient(from 180deg, red, yellow, lime, aqua, blue, magenta, red)'
                        }}
                    >
                         {!isCustomColor && (
                             <Plus size={14} className="text-white drop-shadow-md" />
                         )}
                    </div>
                    <input
                        type="color"
                        value={data.themeConfig?.color}
                        onChange={(e) => updateTheme('color', e.target.value)}
                        className="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
                        title={t('editor.color.customTitle')}
                    />
                </div>
            </div>
        </div>

        {/* Fonts */}
        <div>
            <h3 className="text-sm font-semibold text-gray-900 mb-4 flex items-center">
                <Type size={18} className="mr-2 text-gray-500"/> {t('editor.font')}
            </h3>
            <div className="relative">
              <select
                value={data.themeConfig.fontFamily}
                onChange={(e) => updateTheme('fontFamily', e.target.value)}
                className="block w-full pl-3 pr-10 py-2.5 text-base border-gray-300 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm rounded-md border"
              >
                {/* Grouping fonts manually since React doesn't support complex data structures in select directly easily with map inside optgroup without processing */}
                <optgroup label={t('editor.font.group.englishSans')}>
                  {FONT_OPTIONS.filter(f => f.group === 'English - Sans Serif').map(f => (
                    <option key={f.id} value={f.id}>{f.label}</option>
                  ))}
                </optgroup>
                <optgroup label={t('editor.font.group.englishSerif')}>
                  {FONT_OPTIONS.filter(f => f.group === 'English - Serif').map(f => (
                    <option key={f.id} value={f.id}>{f.label}</option>
                  ))}
                </optgroup>
                <optgroup label={t('editor.font.group.englishMono')}>
                  {FONT_OPTIONS.filter(f => f.group === 'English - Mono').map(f => (
                    <option key={f.id} value={f.id}>{f.label}</option>
                  ))}
                </optgroup>
                <optgroup label={t('editor.font.group.chineseSans')}>
                  {FONT_OPTIONS.filter(f => f.group === 'Chinese - Sans (黑体)').map(f => (
                    <option key={f.id} value={f.id}>{f.label}</option>
                  ))}
                </optgroup>
                <optgroup label={t('editor.font.group.chineseSerif')}>
                  {FONT_OPTIONS.filter(f => f.group === 'Chinese - Serif (宋体/楷体)').map(f => (
                    <option key={f.id} value={f.id}>{f.label}</option>
                  ))}
                </optgroup>
              </select>
            </div>
            <p className="mt-2 text-xs text-gray-500">
              {t('editor.font.note')}
            </p>
        </div>

        {/* Spacing */}
        <div>
             <h3 className="text-sm font-semibold text-gray-900 mb-4 flex items-center">
                <LayoutTemplate size={18} className="mr-2 text-gray-500"/> {t('editor.spacing')}
            </h3>
            <div className="space-y-3">
                <button 
                    onClick={() => updateTheme('spacing', 'compact')}
                    className={`w-full p-2 text-sm border rounded-md transition-colors ${data.themeConfig?.spacing === 'compact' ? 'bg-blue-50 border-blue-500 text-blue-700' : 'hover:bg-gray-50'}`}
                >
                    {t('editor.spacing.option.compact')}
                </button>
                 <button 
                    onClick={() => updateTheme('spacing', 'normal')}
                    className={`w-full p-2 text-sm border rounded-md transition-colors ${data.themeConfig?.spacing === 'normal' ? 'bg-blue-50 border-blue-500 text-blue-700' : 'hover:bg-gray-50'}`}
                >
                    {t('editor.spacing.option.normal')}
                </button>
                 <button 
                    onClick={() => updateTheme('spacing', 'spacious')}
                    className={`w-full p-2 text-sm border rounded-md transition-colors ${data.themeConfig?.spacing === 'spacious' ? 'bg-blue-50 border-blue-500 text-blue-700' : 'hover:bg-gray-50'}`}
                >
                    {t('editor.spacing.option.spacious')}
                </button>
            </div>
        </div>
    </div>
  );

  return (
    <div className="h-full flex flex-col bg-white border-r border-gray-200">
        {/* Tabs */}
        <div className="flex border-b border-gray-200">
            <button 
                onClick={() => setActiveTab('content')}
                className={`flex-1 py-4 text-sm font-medium border-b-2 transition-colors ${activeTab === 'content' ? 'border-blue-600 text-blue-600' : 'border-transparent text-gray-500 hover:text-gray-700'}`}
            >
                {t('editor.content')}
            </button>
             <button 
                onClick={() => setActiveTab('design')}
                className={`flex-1 py-4 text-sm font-medium border-b-2 transition-colors ${activeTab === 'design' ? 'border-blue-600 text-blue-600' : 'border-transparent text-gray-500 hover:text-gray-700'}`}
            >
                {t('editor.design')}
            </button>
        </div>

        {/* Content Area */}
        <div className="flex-1 overflow-y-auto">
            {activeTab === 'content' ? renderContentTab() : renderDesignTab()}
        </div>
    </div>
  );
};
