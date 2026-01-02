import React, { useState, useEffect, useRef } from 'react';
import { Trash2, Plus, Sparkles, ChevronDown, ChevronUp, Upload, X, Image as ImageIcon, Palette, Type, LayoutTemplate, Briefcase, GraduationCap, Wrench, User, Target, BookOpen, Layers, Award, Heart } from 'lucide-react';
import { DndContext, closestCenter, KeyboardSensor, PointerSensor, useSensor, useSensors, DragEndEvent } from '@dnd-kit/core';
import { SortableContext, arrayMove, verticalListSortingStrategy } from '@dnd-kit/sortable';
import { ResumeData, ResumeSection, ResumeItem, ResumeSectionType, ThemeConfig } from '../../types';
import { Button } from '../../components/ui/Button';
import { polishText, generateSummary } from '../../services/geminiService';
import { API_BASE } from '../../config';
import { useLanguage } from '../../contexts/LanguageContext';
import { useConfirm } from '../../components/ui/ConfirmDialog';
import { SortableSection } from './SortableSection';
import { RichTextEditor } from '../../components/ui/RichTextEditor';

interface EditorFormProps {
  data: ResumeData;
  onChange: (data: ResumeData) => void;
}

const generateUUID = () => {
    const c: any = (globalThis as any).crypto;
    if (c?.randomUUID) return c.randomUUID();
    const arr = new Uint32Array(4);
    if (c?.getRandomValues) {
      c.getRandomValues(arr);
      return Array.from(arr).map(n => n.toString(16)).join('');
    }
    return Math.random().toString(36).substr(2, 9);
};

export const EditorForm: React.FC<EditorFormProps> = ({ data, onChange }) => {
  const { t } = useLanguage();
  const confirm = useConfirm();
  const sensors = useSensors(
    useSensor(PointerSensor),
    useSensor(KeyboardSensor)
  );

  const fontOptions = [
    { group: t('editor.font.group.englishSans'), id: 'inter', label: t('font.inter') },
    { group: t('editor.font.group.englishSans'), id: 'roboto', label: t('font.roboto') },
    { group: t('editor.font.group.englishSerif'), id: 'merriweather', label: t('font.merriweather') },
    { group: t('editor.font.group.englishSerif'), id: 'playfair', label: t('font.playfair') },
    { group: t('editor.font.group.englishMono'), id: 'mono', label: t('font.mono') },
    { group: t('editor.font.group.chineseSans'), id: 'yahei', label: t('font.yahei') },
    { group: t('editor.font.group.chineseSans'), id: 'notosans', label: t('font.notosans') },
    { group: t('editor.font.group.chineseSerif'), id: 'simsun', label: t('font.simsun') },
    { group: t('editor.font.group.chineseSerif'), id: 'kaiti', label: t('font.kaiti') },
  ];
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

  const addCustomInfo = () => {
    const list = data.personalInfo.customInfo || [];
    onChange({
      ...data,
      personalInfo: { ...data.personalInfo, customInfo: [...list, { label: '', value: '' }] }
    });
  };

  const updateCustomInfo = (idx: number, key: 'label' | 'value', value: string) => {
    const list = (data.personalInfo.customInfo || []).map((item, i) => i === idx ? { ...item, [key]: value } : item);
    onChange({
      ...data,
      personalInfo: { ...data.personalInfo, customInfo: list }
    });
  };

  const removeCustomInfo = (idx: number) => {
    const list = (data.personalInfo.customInfo || []).filter((_, i) => i !== idx);
    onChange({
      ...data,
      personalInfo: { ...data.personalInfo, customInfo: list }
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
      const token = localStorage.getItem('token');
      const res = await fetch(`${API_BASE}/upload/avatar`, {
        method: 'POST',
        headers: token ? { Authorization: `Bearer ${token}` } : undefined,
        body: form
      });
      if (res.status === 401) {
        alert('登录已过期，请重新登录');
        return;
      }
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
      id: generateUUID(),
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

  useEffect(() => {
    const required: ResumeSectionType[] = [
      ResumeSectionType.JobApplication,
      ResumeSectionType.Exam,
      ResumeSectionType.Education,
      ResumeSectionType.Experience,
      ResumeSectionType.Projects,
      ResumeSectionType.Internships,
      ResumeSectionType.Portfolio,
      ResumeSectionType.Skills,
      ResumeSectionType.Awards,
      ResumeSectionType.SelfEvaluation,
      ResumeSectionType.Interests
    ];
    const existing = new Set(data.sections.map(s => s.type));
    const toAdd: ResumeSection[] = [];
    required.forEach(type => {
      if (!existing.has(type)) {
        toAdd.push({
          id: generateUUID(),
          type,
          title: '',
          isVisible: true,
          items: type === ResumeSectionType.SelfEvaluation ? [{ id: generateUUID(), description: '' }] : []
        });
      }
    });
    if (toAdd.length) {
      onChange({ ...data, sections: [...data.sections, ...toAdd] });
    }
  }, [data.sections]);

  useEffect(() => {
    const job = data.sections.find(s => s.type === ResumeSectionType.JobApplication);
    if (job && job.items.length === 0) {
      onChange({
        ...data,
        sections: data.sections.map(s =>
          s.id === job.id
            ? { ...s, items: [{ id: generateUUID(), title: '', subtitle: '', timeStart: '', description: '' }] }
            : s
        ),
      });
    }
  }, [data.sections]);

  useEffect(() => {
    const needTypesWithTime: ResumeSectionType[] = [
      ResumeSectionType.Education,
      ResumeSectionType.Experience,
      ResumeSectionType.Projects,
      ResumeSectionType.Internships,
    ];
    const needTypesDescOnly: ResumeSectionType[] = [
      ResumeSectionType.Portfolio,
      ResumeSectionType.Skills,
      ResumeSectionType.Interests,
      ResumeSectionType.Awards,
    ];
    let changed = false;
    const nextSections = data.sections.map(s => {
      if (s.type === ResumeSectionType.Summary || s.type === ResumeSectionType.SelfEvaluation || s.type === ResumeSectionType.JobApplication) {
        return s;
      }
      if (s.items.length === 0) {
        changed = true;
        if (needTypesWithTime.includes(s.type)) {
          return { ...s, items: [{ id: generateUUID(), title: '', subtitle: '', timeStart: '', timeEnd: '', today: false as any, description: '' }] };
        }
        if (needTypesDescOnly.includes(s.type)) {
          return { ...s, items: [{ id: generateUUID(), description: '' }] };
        }
        return { ...s, items: [{ id: generateUUID(), title: '', description: '' }] };
      }
      return s;
    });
    if (changed) {
      onChange({ ...data, sections: nextSections });
    }
  }, [data.sections]);

  const removeSection = async (sectionId: string) => {
      const ok = await confirm({ title: t('common.confirmAction'), message: t('editor.removeSection'), variant: 'danger' });
      if (!ok) return;
      onChange({
          ...data,
          sections: data.sections.filter(s => s.id !== sectionId)
      });
  };

  useEffect(() => {
    const summaries = data.sections.filter(s => s.type === ResumeSectionType.Summary);
    if (!summaries.length) return;
    const selfIndex = data.sections.findIndex(s => s.type === ResumeSectionType.SelfEvaluation);
    const otherSections = data.sections.filter(s => s.type !== ResumeSectionType.Summary);
    if (selfIndex >= 0) {
      const self = data.sections[selfIndex];
      const mergedItems = [...self.items, ...summaries.flatMap(s => s.items.length ? s.items : [{ id: generateUUID(), description: '' }])];
      const merged = otherSections.map(s => s.type === ResumeSectionType.SelfEvaluation ? { ...s, items: mergedItems, isVisible: true } : s);
      onChange({ ...data, sections: merged });
    } else {
      const newSelf: ResumeSection = {
        id: generateUUID(),
        type: ResumeSectionType.SelfEvaluation,
        title: '',
        isVisible: true,
        items: summaries.flatMap(s => s.items.length ? s.items : [{ id: generateUUID(), description: '' }])
      };
      onChange({ ...data, sections: [...otherSections, newSelf] });
    }
  }, [data.sections]);

  useEffect(() => {
    const order: ResumeSectionType[] = [
      ResumeSectionType.JobApplication,
      ResumeSectionType.Exam,
      ResumeSectionType.Education,
      ResumeSectionType.Experience,
      ResumeSectionType.Projects,
      ResumeSectionType.Internships,
      ResumeSectionType.Portfolio,
      ResumeSectionType.Skills,
      ResumeSectionType.Awards,
      ResumeSectionType.Interests,
      ResumeSectionType.SelfEvaluation
    ];
    const orderMap = new Map(order.map((t, i) => [t, i]));
    const origIndex = new Map(data.sections.map((s, i) => [s.id, i]));
    const sorted = [...data.sections].sort((a, b) => {
      const ai = orderMap.get(a.type);
      const bi = orderMap.get(b.type);
      const aKnown = ai !== undefined;
      const bKnown = bi !== undefined;
      if (aKnown && bKnown) return (ai as number) - (bi as number);
      if (aKnown && !bKnown) return -1;
      if (!aKnown && bKnown) return 1;
      return (origIndex.get(a.id) as number) - (origIndex.get(b.id) as number);
    });
    const changed = sorted.some((s, i) => s.id !== data.sections[i]?.id);
    if (changed) {
      onChange({ ...data, sections: sorted });
    }
  }, [data.sections]);

  const addCustomSection = () => {
    const newSection: ResumeSection = {
      id: generateUUID(),
      type: ResumeSectionType.Custom,
      title: '',
      isVisible: true,
      items: [{ id: generateUUID(), title: '', description: '' }]
    };
    onChange({ ...data, sections: [...data.sections, newSection] });
  };

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event;
    if (active.id !== over?.id) {
        const oldIndex = data.sections.findIndex((s) => s.id === active.id);
        const newIndex = data.sections.findIndex((s) => s.id === over?.id);
        onChange({
            ...data,
            sections: arrayMove(data.sections, oldIndex, newIndex),
        });
    }
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
      <div className="group rounded-2xl overflow-hidden border border-gray-200 bg-white shadow-sm mb-4">
        <div 
          className="w-full px-5 py-4 flex items-center justify-between text-gray-800"
          onClick={() => setActiveSection(activeSection === 'personal' ? null : 'personal')}
        >
            <div className="flex items-center">
              <div className="mr-3 w-6 h-6"></div>
              <div className="w-8 h-8 rounded-lg bg-indigo-50 text-indigo-600 flex items-center justify-center mr-3">
                <User size={18}/>
              </div>
              <span className="font-semibold">{t('editor.personal')}</span>
            </div>
            <button className="p-1 text-gray-500 hover:text-gray-700 transition-colors">
              {activeSection === 'personal' ? <ChevronUp size={16}/> : <ChevronDown size={16}/>}
            </button>
        </div>
        {activeSection === 'personal' && (
          <div className="px-6 pb-6 space-y-5">
             
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
             {/* removed website/linkedin fields */}

             <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700">{t('editor.fields.gender')}</label>
                  <input type="text" value={data.personalInfo.gender || ''} onChange={e => updatePersonalInfo('gender', e.target.value)} className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2 sm:text-sm"/>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700">{t('editor.fields.age')}</label>
                  <input type="text" value={data.personalInfo.age || ''} onChange={e => updatePersonalInfo('age', e.target.value)} className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2 sm:text-sm"/>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700">{t('editor.fields.maritalStatus')}</label>
                  <input type="text" value={data.personalInfo.maritalStatus || ''} onChange={e => updatePersonalInfo('maritalStatus', e.target.value)} className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2 sm:text-sm"/>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700">{t('editor.fields.politicalStatus')}</label>
                  <input type="text" value={data.personalInfo.politicalStatus || ''} onChange={e => updatePersonalInfo('politicalStatus', e.target.value)} className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2 sm:text-sm"/>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700">{t('editor.fields.birthplace')}</label>
                  <input type="text" value={data.personalInfo.birthplace || ''} onChange={e => updatePersonalInfo('birthplace', e.target.value)} className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2 sm:text-sm"/>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700">{t('editor.fields.ethnicity')}</label>
                  <input type="text" value={data.personalInfo.ethnicity || ''} onChange={e => updatePersonalInfo('ethnicity', e.target.value)} className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2 sm:text-sm"/>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700">{t('editor.fields.height')}</label>
                  <input type="text" value={data.personalInfo.height || ''} onChange={e => updatePersonalInfo('height', e.target.value)} className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2 sm:text-sm"/>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700">{t('editor.fields.weight')}</label>
                  <input type="text" value={data.personalInfo.weight || ''} onChange={e => updatePersonalInfo('weight', e.target.value)} className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2 sm:text-sm"/>
                </div>
             </div>

             <div className="bg-gray-50 p-4 rounded-lg border border-gray-200">
               <div className="flex items-center justify-between mb-3">
                 <span className="text-sm font-medium text-gray-700">{t('editor.customInfo.title')}</span>
                 <Button size="sm" variant="outline" icon={<Plus size={14}/>} onClick={addCustomInfo}>{t('editor.customInfo.add')}</Button>
               </div>
               <div className="space-y-3">
                 {(data.personalInfo.customInfo || []).map((ci, idx) => (
                   <div key={idx} className="grid grid-cols-12 gap-2 items-center">
                     <div className="col-span-5">
                       <input placeholder={t('editor.customInfo.label')} value={ci.label} onChange={e => updateCustomInfo(idx, 'label', e.target.value)} className="block w-full border border-gray-300 rounded-md shadow-sm p-2 sm:text-sm"/>
                     </div>
                     <div className="col-span-6">
                       <input placeholder={t('editor.customInfo.value')} value={ci.value} onChange={e => updateCustomInfo(idx, 'value', e.target.value)} className="block w-full border border-gray-300 rounded-md shadow-sm p-2 sm:text-sm"/>
                     </div>
                     <div className="col-span-1 flex justify-end">
                       <button onClick={() => removeCustomInfo(idx)} className="p-2 text-gray-400 hover:text-red-500 rounded hover:bg-red-50">
                         <Trash2 size={16}/>
                       </button>
                     </div>
                   </div>
                 ))}
               </div>
             </div>
          </div>
        )}
      </div>

      {/* Dynamic Sections */}
      <DndContext 
        sensors={sensors}
        collisionDetection={closestCenter}
        onDragEnd={handleDragEnd}
      >
          <SortableContext 
            items={data.sections.map(s => s.id)}
            strategy={verticalListSortingStrategy}
          >
            {data.sections.map((section) => (
                <SortableSection
                    key={section.id}
                    section={section}
                    isActive={activeSection === section.id}
                    onToggle={() => setActiveSection(activeSection === section.id ? null : section.id)}
                    onUpdate={(updates) => updateSection(section.id, updates)}
                    onRemove={() => removeSection(section.id)}
                    icon={
                      section.type === ResumeSectionType.Experience ? <Briefcase size={18}/> :
                      section.type === ResumeSectionType.Education ? <GraduationCap size={18}/> :
                      section.type === ResumeSectionType.Skills ? <Wrench size={18}/> :
                      section.type === ResumeSectionType.Projects ? <Layers size={18}/> :
                      section.type === ResumeSectionType.Internships ? <Briefcase size={18}/> :
                      section.type === ResumeSectionType.Portfolio ? <ImageIcon size={18}/> :
                      section.type === ResumeSectionType.Awards ? <Award size={18}/> :
                      section.type === ResumeSectionType.Interests ? <Heart size={18}/> :
                      section.type === ResumeSectionType.JobApplication ? <Target size={18}/> :
                      section.type === ResumeSectionType.Exam ? <BookOpen size={18}/> :
                      section.type === ResumeSectionType.SelfEvaluation ? <User size={18}/> :
                      section.type === ResumeSectionType.Summary ? <Sparkles size={18}/> :
                      section.type === ResumeSectionType.Custom ? <LayoutTemplate size={18}/> :
                      <Type size={18}/>
                    }
                    onAddItem={
                      section.type !== ResumeSectionType.Summary && section.type !== ResumeSectionType.JobApplication
                        ? () => addItem(section.id)
                        : undefined
                    }
                >
                    
                    {section.type === ResumeSectionType.JobApplication && (
                      <div className="bg-white p-4 rounded-lg shadow-sm border border-gray-200 relative group transition-all hover:shadow-md">
                        <div className="mt-2 grid grid-cols-1 gap-4 mb-3">
                          <input
                            placeholder="求职岗位"
                            className="w-full font-medium border-b border-gray-200 focus:border-blue-500 outline-none pb-1 bg-transparent transition-colors"
                            value={section.items[0]?.title || ''}
                            onChange={e =>
                              updateItem(section.id, section.items[0]?.id || '', 'title', e.target.value)
                            }
                          />
                          <div className="grid grid-cols-2 gap-4">
                            <input
                              placeholder="意向城市"
                              className="w-full text-sm border-b border-gray-200 focus:border-blue-500 outline-none pb-1 bg-transparent"
                              value={section.items[0]?.subtitle || ''}
                              onChange={e =>
                                updateItem(section.id, section.items[0]?.id || '', 'subtitle', e.target.value)
                              }
                            />
                            <input
                              placeholder="期望薪资"
                              className="w-full text-sm border-b border-gray-200 focus:border-blue-500 outline-none pb-1 bg-transparent"
                              value={section.items[0]?.description || ''}
                              onChange={e =>
                                updateItem(section.id, section.items[0]?.id || '', 'description', e.target.value)
                              }
                            />
                            <div className="flex items-center gap-2">
                              <input
                                type="month"
                                placeholder="入职时间"
                                className="w-full text-sm border-b border-gray-200 focus:border-blue-500 outline-none pb-1 bg-transparent"
                                value={section.items[0]?.timeStart || ''}
                                onChange={e =>
                                  updateItem(section.id, section.items[0]?.id || '', 'timeStart', e.target.value)
                                }
                              />
                              <span className="text-xs text-gray-500">-</span>
                              <input
                                type="month"
                                disabled
                                className="w-full text-sm border-b outline-none pb-1 bg-transparent border-gray-200 text-gray-400"
                                value=""
                                onChange={() => {}}
                              />
                              <label className="flex items-center ml-2 text-xs text-gray-600">
                                <input type="checkbox" className="mr-1" checked={false} readOnly />
                                {t('common.toPresent') || '至今'}
                              </label>
                            </div>
                          </div>
                        </div>
                      </div>
                    )}

                    {section.type !== ResumeSectionType.JobApplication && section.items.map((item) => (
                    <div key={item.id} className="bg-white p-4 rounded-lg shadow-sm border border-gray-200 relative group transition-all hover:shadow-md">
                        <button 
                            onClick={() => removeItem(section.id, item.id)}
                            className="absolute -top-3 -right-3 p-1.5 rounded-full bg-white border border-gray-200 text-gray-400 hover:text-red-600 hover:border-red-300 shadow-sm"
                            title={t('dashboard.confirm.delete')}
                        >
                            <Trash2 size={14} />
                        </button>
                        
                        {[
                          ResumeSectionType.Experience,
                          ResumeSectionType.Education,
                          ResumeSectionType.Projects,
                          ResumeSectionType.Internships,
                        ].includes(section.type) && (
                          <div className="mt-2 grid grid-cols-1 gap-4 mb-3">
                            <input
                              placeholder={
                                section.type === ResumeSectionType.Education
                                  ? '学校名称'
                                  : section.type === ResumeSectionType.Experience
                                  ? '公司名称'
                                  : section.type === ResumeSectionType.Projects
                                  ? '项目名称'
                                  : section.type === ResumeSectionType.Internships
                                  ? '公司名称'
                                  : t('editor.placeholder.titleExample')
                              }
                              className="w-full font-medium border-b border-gray-200 focus:border-blue-500 outline-none pb-1 bg-transparent transition-colors"
                              value={item.title}
                              onChange={e => updateItem(section.id, item.id, 'title', e.target.value)}
                            />
                            <div className="grid grid-cols-2 gap-4">
                              <input
                                placeholder={
                                  section.type === ResumeSectionType.Education
                                    ? '院系/学院（可选）'
                                    : section.type === ResumeSectionType.Experience
                                    ? '职位'
                                    : section.type === ResumeSectionType.Projects
                                    ? '参与角色'
                                    : section.type === ResumeSectionType.Internships
                                    ? '职位'
                                    : t('editor.placeholder.subtitleExample')
                                }
                                className="w-full text-sm border-b border-gray-200 focus:border-blue-500 outline-none pb-1 bg-transparent"
                                value={item.subtitle || ''}
                                onChange={e => updateItem(section.id, item.id, 'subtitle', e.target.value)}
                              />
                              {section.type === ResumeSectionType.Education && (
                                <input
                                  placeholder="所学专业"
                                  className="w-full text-sm border-b border-gray-200 focus:border-blue-500 outline-none pb-1 bg-transparent"
                                  value={item.major || ''}
                                  onChange={e => updateItem(section.id, item.id, 'major', e.target.value)}
                                />
                              )}
                              <div className="flex items-center gap-2">
                                <input
                                  type="month"
                                  className="w-full text-sm border-b border-gray-200 focus:border-blue-500 outline-none pb-1 bg-transparent"
                                  value={item.timeStart || ''}
                                  onChange={e => updateItem(section.id, item.id, 'timeStart', e.target.value)}
                                />
                                <span className="text-xs text-gray-500">-</span>
                                <input
                                  type="month"
                                  disabled={item.today}
                                  className={`w-full text-sm border-b outline-none pb-1 bg-transparent ${item.today ? 'border-gray-200 text-gray-400' : 'border-gray-200 focus:border-blue-500'}`}
                                  value={item.timeEnd || ''}
                                  onChange={e => updateItem(section.id, item.id, 'timeEnd', e.target.value)}
                                />
                                <label className="flex items-center ml-2 text-xs text-gray-600">
                                  <input
                                    type="checkbox"
                                    className="mr-1"
                                    checked={!!item.today}
                                    onChange={e => updateItem(section.id, item.id, 'today', e.target.checked ? true as any : false as any)}
                                  />
                                  {t('common.toPresent') || '至今'}
                                </label>
                              </div>
                              {section.type === ResumeSectionType.Education && (
                                <input
                                  placeholder="学历"
                                  className="w-full text-sm border-b border-gray-200 focus:border-blue-500 outline-none pb-1 bg-transparent"
                                  value={item.degree || ''}
                                  onChange={e => updateItem(section.id, item.id, 'degree', e.target.value)}
                                />
                              )}
                            </div>
                          </div>
                        )}

                        <div className="relative mt-2">
                            {section.type === ResumeSectionType.Skills ? (
                                <textarea
                                    rows={2}
                                    className="w-full text-sm border border-gray-300 rounded p-2 focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
                                    placeholder={t('editor.placeholder.skills')}
                                    value={item.description}
                                    onChange={e => updateItem(section.id, item.id, 'description', e.target.value)}
                                />
                            ) : (
                                <RichTextEditor
                                    value={item.description}
                                    valueFormat={/(<\/?[a-z][\s\S]*>)/i.test(item.description || '') ? 'html' : 'text'}
                                    outputFormat="html"
                                    onChange={(val) => updateItem(section.id, item.id, 'description', val)}
                                    aiContext={section.type === ResumeSectionType.Experience ? 'Work Experience' : (section.type === ResumeSectionType.Summary ? 'Resume Summary' : undefined)}
                                    minRows={4}
                                    maxHeight={300}
                                />
                            )}
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

                </SortableSection>
            ))}
          </SortableContext>
      </DndContext>

 
      <div className="mt-4">
        <Button
          variant="outline"
          size="sm"
          icon={<Plus size={14} />}
          onClick={addCustomSection}
          className="w-full"
        >
          {t('editor.addSection')}
        </Button>
      </div>

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
                <optgroup label={t('editor.font.group.englishSans')}>
                  {fontOptions.filter(f => f.group === t('editor.font.group.englishSans')).map(f => (
                    <option key={f.id} value={f.id}>{f.label}</option>
                  ))}
                </optgroup>
                <optgroup label={t('editor.font.group.englishSerif')}>
                  {fontOptions.filter(f => f.group === t('editor.font.group.englishSerif')).map(f => (
                    <option key={f.id} value={f.id}>{f.label}</option>
                  ))}
                </optgroup>
                <optgroup label={t('editor.font.group.englishMono')}>
                  {fontOptions.filter(f => f.group === t('editor.font.group.englishMono')).map(f => (
                    <option key={f.id} value={f.id}>{f.label}</option>
                  ))}
                </optgroup>
                <optgroup label={t('editor.font.group.chineseSans')}>
                  {fontOptions.filter(f => f.group === t('editor.font.group.chineseSans')).map(f => (
                    <option key={f.id} value={f.id}>{f.label}</option>
                  ))}
                </optgroup>
                <optgroup label={t('editor.font.group.chineseSerif')}>
                  {fontOptions.filter(f => f.group === t('editor.font.group.chineseSerif')).map(f => (
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
    <div className="h-full flex flex-col bg-white">
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
        <div className="flex-1 overflow-y-auto px-4 md:px-6 pt-4 pb-4">
            {activeTab === 'content' ? renderContentTab() : renderDesignTab()}
        </div>
    </div>
  );
};
