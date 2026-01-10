import React, { useState, useEffect, useRef } from 'react';
import { Trash2, Plus, Sparkles, ChevronDown, ChevronUp, Upload, X, Image as ImageIcon, Palette, Type, LayoutTemplate, Briefcase, GraduationCap, Wrench, User, BookOpen, Layers, Award, Heart } from 'lucide-react';
import { DndContext, closestCenter, KeyboardSensor, PointerSensor, useSensor, useSensors, DragEndEvent } from '@dnd-kit/core';
import { SortableContext, arrayMove, verticalListSortingStrategy } from '@dnd-kit/sortable';
import { ResumeData, ResumeSection, ResumeItem, ResumeSectionType } from '../../types';
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
    { group: t('editor.font.group.chineseSans'), id: 'notosans', label: t('font.notosans') },
  ];
  const [activeTab, setActiveTab] = useState<'content' | 'design'>('content');
  const [activeSection, setActiveSection] = useState<string | null>('personal');
  const [isAiLoading, setIsAiLoading] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const personal = (data.Personal || {}) as NonNullable<ResumeData['Personal']>;
  const theme = (data.Theme || {}) as NonNullable<ResumeData['Theme']>;
  const parseCustomList = (): Array<{ label: string; value: string }> => {
    try {
      const raw = personal?.CustomInfo;
      if (raw) {
        const parsed = JSON.parse(raw);
        if (Array.isArray(parsed)) return parsed;
      }
    } catch {}
    return [];
  };
  const updatePersonal = (key: keyof NonNullable<ResumeData['Personal']>, value: string) => {
    onChange({ ...data, Personal: { ...personal, [key]: value } });
  };

  const addCustomInfo = () => {
    const list = parseCustomList();
    const next = [...list, { label: '', value: '' }];
    onChange({ ...data, Personal: { ...personal, CustomInfo: JSON.stringify(next) } });
  };

  const updateCustomInfo = (idx: number, key: 'label' | 'value', value: string) => {
    const list = parseCustomList().map((item, i) => i === idx ? { ...item, [key]: value } : item);
    onChange({ ...data, Personal: { ...personal, CustomInfo: JSON.stringify(list) } });
  };

  const removeCustomInfo = (idx: number) => {
    const list = parseCustomList().filter((_, i) => i !== idx);
    onChange({ ...data, Personal: { ...personal, CustomInfo: JSON.stringify(list) } });
  };

  const updateTheme = (key: keyof NonNullable<ResumeData['Theme']>, value: string) => {
    onChange({ ...data, Theme: { ...theme, [key]: value } });
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
        updatePersonal('AvatarURL', data.url);
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
    const nextSections = data.sections.map(s => {
      if (s.id !== sectionId) return s;
      const items = [...s.items, newItem].map((it, idx) => ({ ...it, orderNum: idx }));
      return { ...s, items };
    });
    onChange({ ...data, sections: nextSections.map((s, idx) => ({ ...s, orderNum: idx })) });
  };

  const removeItem = (sectionId: string, itemId: string) => {
    const nextSections = data.sections.map(s => {
      if (s.id !== sectionId) return s;
      const items = s.items.filter(i => i.id !== itemId).map((it, idx) => ({ ...it, orderNum: idx }));
      return { ...s, items };
    });
    onChange({ ...data, sections: nextSections.map((s, idx) => ({ ...s, orderNum: idx })) });
  };

  useEffect(() => {
    const required: ResumeSectionType[] = [
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
      if (s.type === ResumeSectionType.Summary || s.type === ResumeSectionType.SelfEvaluation) {
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
      const next = data.sections.filter(s => s.id !== sectionId).map((s, idx) => ({ ...s, orderNum: idx, items: s.items.map((it, ii) => ({ ...it, orderNum: ii })) }));
      onChange({ ...data, sections: next });
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
  }, [data.sections]);

  const addCustomSection = () => {
    const newSection: ResumeSection = {
      id: generateUUID(),
      type: ResumeSectionType.Custom,
      title: '',
      isVisible: true,
      items: [{ id: generateUUID(), title: '', description: '', orderNum: 0 }],
      orderNum: (data.sections.length || 0)
    };
    const next = [...data.sections, newSection].map((s, idx) => ({ ...s, orderNum: idx, items: s.items.map((it, ii) => ({ ...it, orderNum: ii })) }));
    onChange({ ...data, sections: next });
  };

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event;
    if (active.id !== over?.id) {
        const oldIndex = data.sections.findIndex((s) => s.id === active.id);
        const newIndex = data.sections.findIndex((s) => s.id === over?.id);
        const movedArr = arrayMove(data.sections, oldIndex, newIndex) as ResumeSection[];
        const moved = movedArr.map((s, idx) => ({ ...s, orderNum: idx, items: s.items.map((it, ii) => ({ ...it, orderNum: ii })) }));
        onChange({ ...data, sections: moved });
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
     
    const summary = await generateSummary(personal?.Job || '', skills);
     
     // Update summary section
     const summarySection = data.sections.find(s => s.type === ResumeSectionType.Summary);
     if (summarySection && summarySection.items.length > 0) {
         updateItem(summarySection.id, summarySection.items[0].id, 'description', summary);
     }
     setIsAiLoading(false);
  }

  const getLabels = (type: ResumeSectionType) => {
    return {
      title:
        type === ResumeSectionType.Education
          ? t('editor.labels.education.school')
          : type === ResumeSectionType.Experience
          ? t('editor.labels.experience.title')
          : type === ResumeSectionType.Projects
          ? t('editor.labels.projects.title')
          : type === ResumeSectionType.Internships
          ? t('editor.labels.internships.company')
          : t('editor.placeholder.titleExample'),
      subtitle:
        type === ResumeSectionType.Education
          ? t('editor.labels.education.department')
          : type === ResumeSectionType.Experience
          ? t('editor.labels.experience.role')
          : type === ResumeSectionType.Projects
          ? t('editor.labels.projects.role')
          : type === ResumeSectionType.Internships
          ? t('editor.labels.experience.role')
          : t('editor.placeholder.subtitleExample'),
      major: t('editor.labels.education.major'),
      degree: t('editor.labels.education.degree'),
      start: t('editor.labels.common.startDate'),
      end: t('editor.labels.common.endDate'),
      desc: t('editor.labels.common.description'),
    };
  };

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
            {personal?.AvatarURL ? (
                        <div className="relative group shrink-0">
                            <img 
                                src={personal.AvatarURL || ''} 
                                alt={t('a11y.avatarAlt')} 
                                className="w-16 h-16 rounded-full object-cover border-2 border-white shadow-sm"
                            />
                            <button
                                onClick={() => updatePersonal('AvatarURL', '')}
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
                            {personal?.AvatarURL ? t('editor.profile.change') : t('editor.profile.upload')}
                        </Button>
                        <p className="mt-2 text-xs text-gray-500">
                            {t('editor.profile.tipSquare')}<br/>
                            {t('editor.profile.tipMaxSize')}
                        </p>
                     </div>
                 </div>
             </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700">{t('editor.fields.fullName')}</label>
                <input 
                  type="text" 
                  value={personal?.FullName || ''}
                  onChange={e => updatePersonal('FullName', e.target.value)}
                  className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2 focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700">{t('editor.fields.degree')}</label>
                <input 
                  type="text" 
                  value={personal?.Degree || ''}
                  onChange={e => updatePersonal('Degree', e.target.value)}
                  className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2 focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                />
              </div>
            </div>
             <div className="grid grid-cols-2 gap-4">
                <div>
                    <label className="block text-sm font-medium text-gray-700">{t('editor.fields.email')}</label>
                    <input type="email" value={personal?.Email || ''} onChange={e => updatePersonal('Email', e.target.value)} className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2 sm:text-sm"/>
                </div>
                <div>
                    <label className="block text-sm font-medium text-gray-700">{t('editor.fields.phone')}</label>
                    <input type="tel" value={personal?.Phone || ''} onChange={e => updatePersonal('Phone', e.target.value)} className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2 sm:text-sm"/>
                </div>
             </div>
             {/* removed website/linkedin fields */}

             <div className="grid grid-cols-2 gap-4">
               <div>
                 <label className="block text-sm font-medium text-gray-700">{t('editor.fields.gender')}</label>
                 <input type="text" value={personal?.Gender || ''} onChange={e => updatePersonal('Gender', e.target.value)} className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2 sm:text-sm"/>
               </div>
               <div>
                 <label className="block text-sm font-medium text-gray-700">{t('editor.fields.age')}</label>
                 <input type="text" value={personal?.Age || ''} onChange={e => updatePersonal('Age', e.target.value)} className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2 sm:text-sm"/>
               </div>
             </div>

             <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
               <div>
                  <label className="block text-sm font-medium text-gray-700">{t('editor.fields.jobApplication')}</label>
                  <input 
                    type="text" 
                    value={personal?.Job || ''} 
                    onChange={e => updatePersonal('Job', e.target.value)} 
                    className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2 sm:text-sm"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700">{t('editor.fields.city')}</label>
                  <input 
                    type="text" 
                    value={personal?.City || ''} 
                    onChange={e => updatePersonal('City', e.target.value)} 
                    className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2 sm:text-sm"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700">{t('editor.fields.expectedSalary')}</label>
                  <input 
                    type="text" 
                    value={personal?.Money || ''} 
                    onChange={e => updatePersonal('Money', e.target.value)} 
                    className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2 sm:text-sm"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700">{t('editor.fields.joinTime')}</label>
                  <input
                    type="text"
                    value={personal?.JoinTime || ''}
                    onChange={e => updatePersonal('JoinTime', e.target.value)}
                    className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2 sm:text-sm"
                  />
                </div>
             </div>
             
            
 
             <div className="bg-gray-50 p-4 rounded-lg border border-gray-200">
               <div className="flex items-center justify-between mb-3">
                 <span className="text-sm font-medium text-gray-700">{t('editor.customInfo.title')}</span>
                 <Button size="sm" variant="outline" icon={<Plus size={14}/>} onClick={addCustomInfo}>{t('editor.customInfo.add')}</Button>
               </div>
               <div className="space-y-3">
                 {parseCustomList().map((ci, idx) => (
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
                      section.type === ResumeSectionType.Exam ? <BookOpen size={18}/> :
                      section.type === ResumeSectionType.SelfEvaluation ? <User size={18}/> :
                      section.type === ResumeSectionType.Summary ? <Sparkles size={18}/> :
                      section.type === ResumeSectionType.Custom ? <LayoutTemplate size={18}/> :
                      <Type size={18}/>
                    }
                    onAddItem={
                      section.type !== ResumeSectionType.Summary
                        ? () => addItem(section.id)
                        : undefined
                    }
                >

                    {section.items.map((item) => {
                    const isComplex = [
                      ResumeSectionType.Experience,
                      ResumeSectionType.Education,
                      ResumeSectionType.Projects,
                      ResumeSectionType.Internships,
                    ].includes(section.type);
                    const wrapperClass = isComplex 
                      ? "bg-white p-4 rounded-lg shadow-sm border border-gray-200 relative group/item transition-all hover:shadow-md"
                      : "relative group/item";
                    return (
                    <div key={item.id} className={wrapperClass}>
                        <button 
                            onClick={() => removeItem(section.id, item.id)}
                            className="absolute -top-3 -right-3 p-1.5 rounded-full bg-white border border-gray-200 text-gray-400 hover:text-red-600 hover:border-red-300 shadow-sm opacity-0 group-hover/item:opacity-100 pointer-events-none group-hover/item:pointer-events-auto transition-opacity z-50"
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
                          <div className="mt-2 space-y-4 mb-3">
                            {(() => {
                              const labels = getLabels(section.type);
                              return (
                                <>
                                  <div>
                                    <span className="text-xs text-gray-600 mb-1 block">{labels.title}</span>
                                    <input
                                      className="w-full font-medium border-b border-gray-200 focus:border-blue-500 outline-none pb-1 bg-transparent transition-colors"
                                      value={item.title || ''}
                                      onChange={e => updateItem(section.id, item.id, 'title', e.target.value)}
                                    />
                                  </div>
                                  <div className="grid grid-cols-2 gap-4">
                                    <div>
                                      <span className="text-xs text-gray-600 mb-1 block">{labels.subtitle}</span>
                                      <input
                                        className="w-full text-sm border-b border-gray-200 focus:border-blue-500 outline-none pb-1 bg-transparent"
                                        value={item.subtitle || ''}
                                        onChange={e => updateItem(section.id, item.id, 'subtitle', e.target.value)}
                                      />
                                    </div>
                                    {section.type === ResumeSectionType.Education && (
                                      <div>
                                        <span className="text-xs text-gray-600 mb-1 block">{labels.major}</span>
                                        <input
                                          className="w-full text-sm border-b border-gray-200 focus:border-blue-500 outline-none pb-1 bg-transparent"
                                          value={item.major || ''}
                                          onChange={e => updateItem(section.id, item.id, 'major', e.target.value)}
                                        />
                                      </div>
                                    )}
                                  </div>
                                  <div className="grid grid-cols-2 gap-4">
                                    <div>
                                      <span className="text-xs text-gray-600 mb-1 block">{labels.start}</span>
                                      <input
                                        type="month"
                                        className="w-full text-sm border-b border-gray-200 focus:border-blue-500 outline-none pb-1 bg-transparent"
                                        value={item.timeStart || ''}
                                        onChange={e => updateItem(section.id, item.id, 'timeStart', e.target.value)}
                                      />
                                    </div>
                                    <div>
                                      <span className="text-xs text-gray-600 mb-1 block">{labels.end}</span>
                                      <div className="flex items-center gap-2">
                                        <input
                                          type="month"
                                          disabled={item.today}
                                          className={`w-full text-sm border-b outline-none pb-1 bg-transparent ${item.today ? 'border-gray-200 text-gray-400' : 'border-gray-200 focus:border-blue-500'}`}
                                          value={item.timeEnd || ''}
                                          onChange={e => updateItem(section.id, item.id, 'timeEnd', e.target.value)}
                                        />
                                        <label className="flex items-center text-xs text-gray-600">
                                          <input
                                            type="checkbox"
                                            className="mr-1"
                                            checked={!!item.today}
                                            onChange={e => updateItem(section.id, item.id, 'today', e.target.checked ? true as any : false as any)}
                                          />
                                          {t('common.toPresent') || '至今'}
                                        </label>
                                      </div>
                                    </div>
                                  </div>
                                  {section.type === ResumeSectionType.Education && (
                                    <div>
                                      <span className="text-xs text-gray-600 mb-1 block">{labels.degree}</span>
                                      <input
                                        className="w-full text-sm border-b border-gray-200 focus:border-blue-500 outline-none pb-1 bg-transparent"
                                        value={item.degree || ''}
                                        onChange={e => updateItem(section.id, item.id, 'degree', e.target.value)}
                                      />
                                    </div>
                                  )}
                                </>
                              );
                            })()}
                          </div>
                        )}

                        <div className="relative mt-2">
                            <RichTextEditor
                                value={item.description}
                                onChange={(val) => updateItem(section.id, item.id, 'description', val)}
                                aiContext={section.type === ResumeSectionType.Experience ? 'Work Experience' : (section.type === ResumeSectionType.Summary ? 'Resume Summary' : undefined)}
                                minRows={section.type === ResumeSectionType.Skills ? 3 : 4}
                                maxHeight={300}
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
                    )})}

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
  const isCustomColor = !predefinedColors.includes(theme?.Color || '');

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
                        onClick={() => updateTheme('Color', color)}
                        className={`w-8 h-8 rounded-full border-2 transition-transform hover:scale-110 ${theme?.Color === color ? 'border-gray-900 scale-110 ring-2 ring-gray-200' : 'border-transparent'}`}
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
                                ? (theme?.Color || '')
                                : 'conic-gradient(from 180deg, red, yellow, lime, aqua, blue, magenta, red)'
                        }}
                    >
                         {!isCustomColor && (
                             <Plus size={14} className="text-white drop-shadow-md" />
                         )}
                    </div>
                    <input
                        type="color"
                        value={theme?.Color || ''}
                        onChange={(e) => updateTheme('Color', e.target.value)}
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
                value={fontOptions.some(f => f.id === (theme?.Font || '')) ? (theme?.Font as string) : 'notosans'}
                onChange={(e) => updateTheme('Font', e.target.value)}
                className="block w-full pl-3 pr-10 py-2.5 text-base border-gray-300 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm rounded-md border"
              >
                <optgroup label={t('editor.font.group.chineseSans')}>
                  {fontOptions.filter(f => f.group === t('editor.font.group.chineseSans')).map(f => (
                    <option key={f.id} value={f.id}>{f.label}</option>
                  ))}
                </optgroup>
              </select>
            </div>
            <p className="mt-2 text-xs text-gray-500">
              {t('editor.font.note')}
            </p>
        </div>

        {/* Font Size */}
        <div>
            <h3 className="text-sm font-semibold text-gray-900 mb-4 flex items-center">
                <Type size={18} className="mr-2 text-gray-500"/> {t('editor.fontSize')}
            </h3>
            <div className="grid grid-cols-5 gap-2">
                {[12,13,14,15,16].map(fs => (
                    <button
                        key={fs}
                        onClick={() => updateTheme('FontSize', String(fs))}
                        className={`w-full p-2 text-sm border rounded-md transition-colors ${String(fs) === (theme?.FontSize || '13') ? 'bg-blue-50 border-blue-500 text-blue-700' : 'hover:bg-gray-50'}`}
                    >
                        {fs}
                    </button>
                ))}
            </div>
        </div>

        {/* Spacing */}
        <div>
             <h3 className="text-sm font-semibold text-gray-900 mb-4 flex items-center">
                <LayoutTemplate size={18} className="mr-2 text-gray-500"/> {t('editor.spacing')}
            </h3>
            <div className="space-y-3">
                <button 
                    onClick={() => updateTheme('Spacing', 'compact')}
                    className={`w-full p-2 text-sm border rounded-md transition-colors ${theme?.Spacing === 'compact' ? 'bg-blue-50 border-blue-500 text-blue-700' : 'hover:bg-gray-50'}`}
                >
                    {t('editor.spacing.option.compact')}
                </button>
                 <button 
                    onClick={() => updateTheme('Spacing', 'normal')}
                    className={`w-full p-2 text-sm border rounded-md transition-colors ${theme?.Spacing === 'normal' ? 'bg-blue-50 border-blue-500 text-blue-700' : 'hover:bg-gray-50'}`}
                >
                    {t('editor.spacing.option.normal')}
                </button>
                 <button 
                    onClick={() => updateTheme('Spacing', 'spacious')}
                    className={`w-full p-2 text-sm border rounded-md transition-colors ${theme?.Spacing === 'spacious' ? 'bg-blue-50 border-blue-500 text-blue-700' : 'hover:bg-gray-50'}`}
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
