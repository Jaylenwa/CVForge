import React from 'react';
import { ResumeData, ResumeSectionType } from '../../../types';
import { useLanguage } from '../../../contexts/LanguageContext';
import { useSectionTitle } from '../../../hooks/useSectionTitle';
import { sanitizeHtml, hasExtraPersonalInfo } from '../../../utils/resume-helpers';

export const TemplateMinimalist: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
    const getSectionTitle = useSectionTitle();
    const { t } = useLanguage();
    return (
    <div className={`w-full p-8 md:p-14 bg-white h-full min-h-[1123px] text-gray-800 ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none print:min-h-0 print:h-auto`} style={{ fontFamily: styles.fontFamily, lineHeight: parseFloat(styles.spacingMultiplier) * 1.4 }}>
        <header className="border-b-4 pb-6 mb-10 flex flex-col md:flex-row justify-between items-start md:items-end gap-6" style={{ borderColor: data.themeConfig?.color || 'black' }}>
            <div className="flex-1">
                <h1 className="text-5xl font-black tracking-tight uppercase mb-4 leading-none" style={{ color: data.themeConfig?.color || 'black' }}>{data.personalInfo.fullName}</h1>
                <div className="flex flex-wrap text-sm font-semibold gap-x-6 gap-y-2 text-gray-500 uppercase tracking-wide">
                    <span>{data.personalInfo.jobTitle}</span>
                    {data.personalInfo.email && <span>{data.personalInfo.email}</span>}
                    {data.personalInfo.phone && <span>{data.personalInfo.phone}</span>}
                </div>
                {hasExtraPersonalInfo(data) && (
                  <div className="mt-2 flex flex-wrap text-xs gap-x-6 gap-y-2 text-gray-600">
                      {data.personalInfo.gender && <span>{t('editor.fields.gender')}: {data.personalInfo.gender}</span>}
                      {data.personalInfo.age && <span>{t('editor.fields.age')}: {data.personalInfo.age}</span>}
                      {data.personalInfo.maritalStatus && <span>{t('editor.fields.maritalStatus')}: {data.personalInfo.maritalStatus}</span>}
                      {data.personalInfo.politicalStatus && <span>{t('editor.fields.politicalStatus')}: {data.personalInfo.politicalStatus}</span>}
                      {data.personalInfo.birthplace && <span>{t('editor.fields.birthplace')}: {data.personalInfo.birthplace}</span>}
                      {data.personalInfo.ethnicity && <span>{t('editor.fields.ethnicity')}: {data.personalInfo.ethnicity}</span>}
                      {data.personalInfo.height && <span>{t('editor.fields.height')}: {data.personalInfo.height}</span>}
                      {data.personalInfo.weight && <span>{t('editor.fields.weight')}: {data.personalInfo.weight}</span>}
                      {(data.personalInfo.customInfo || []).map((ci, idx) => (
                        <span key={idx}>{ci.label}: {ci.value}</span>
                      ))}
                  </div>
                )}
            </div>
             {data.personalInfo.avatarUrl && (
                <img 
                    src={data.personalInfo.avatarUrl} 
                    alt="Profile" 
                    className="w-24 h-24 object-cover border border-gray-200 grayscale shadow-sm flex-shrink-0 self-center md:self-end"
                />
            )}
        </header>

        <div className="grid grid-cols-1 gap-10">
             {data.sections.filter(s => s.isVisible).map(section => (
                 <div key={section.id}>
                     <h3 className="text-xs font-bold uppercase tracking-[0.2em] text-gray-400 mb-6">{getSectionTitle(section)}</h3>
                     <div className="space-y-8">
                         {section.items.map(item => (
                             <div key={item.id} className="grid grid-cols-1 md:grid-cols-12 gap-4">
                                 <div className="md:col-span-3 text-xs font-bold text-gray-400 pt-1 uppercase tracking-wide">
                                    {item.dateRange}
                                 </div>
                                 <div className="md:col-span-9">
                                     <h4 className="font-bold text-gray-900 text-lg leading-none mb-1">{item.title}</h4>
                                     {item.subtitle && <p className="text-sm font-semibold text-gray-600 mb-3">{item.subtitle}</p>}
                                     <div className="text-sm leading-relaxed text-gray-700 whitespace-pre-wrap" dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description || '') }} />
                                 </div>
                             </div>
                         ))}
                     </div>
                 </div>
             ))}
        </div>
    </div>
);
};
