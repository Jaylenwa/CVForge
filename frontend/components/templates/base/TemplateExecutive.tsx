import React from 'react';
import { ResumeData, ResumeSectionType } from '../../../types';
import { useLanguage } from '../../../contexts/LanguageContext';
import { useSectionTitle } from '../../../hooks/useSectionTitle';
import { sanitizeHtml, hasExtraPersonalInfo } from '../../../utils/resume-helpers';

export const TemplateExecutive: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
    const getSectionTitle = useSectionTitle();
    const { t } = useLanguage();
    return (
    <div className={`w-full p-10 md:p-12 bg-white h-full min-h-[1123px] text-gray-900 ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none print:min-h-0 print:h-auto`} style={{ fontFamily: styles.fontFamily, lineHeight: parseFloat(styles.spacingMultiplier) * 1.6 }}>
        <div className="text-center border-b pb-6 mb-8 flex flex-col items-center" style={{ borderColor: data.themeConfig?.color || '#111827' }}>
            {data.personalInfo.avatarUrl && (
                 <img 
                    src={data.personalInfo.avatarUrl} 
                    alt="Profile" 
                    className="w-32 h-32 rounded-full mb-6 object-cover border-4 border-double border-gray-200 shadow-sm"
                    style={{ borderRadius: '50%' }}
                />
            )}
            <h1 className="text-3xl font-bold uppercase mb-2 tracking-widest" style={{ color: data.themeConfig?.color }}>{data.personalInfo.fullName}</h1>
            <p className="italic text-lg text-gray-700 mb-3">{data.personalInfo.jobTitle}</p>
                <div className="text-sm text-gray-600 space-x-3 font-sans flex flex-wrap justify-center">
                 <span>{data.personalInfo.phone}</span> <span className="text-gray-300">•</span> <span>{data.personalInfo.email}</span>
            </div>
            {hasExtraPersonalInfo(data) && (
              <div className="mt-2 text-xs text-gray-600 space-x-3 font-sans flex flex-wrap justify-center">
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

        {data.sections.filter(s => s.isVisible).map(section => (
            <div key={section.id} className="mb-6">
                <h3 className="text-md font-bold uppercase border-b border-gray-400 mb-4 pb-1 flex justify-between items-end">
                    <span>{getSectionTitle(section)}</span>
                </h3>
                <div className="space-y-5">
                    {section.items.map(item => (
                        <div key={item.id}>
                            <div className="flex justify-between items-baseline mb-1">
                                <h4 className="font-bold text-gray-900 text-lg">{item.title}</h4>
                                <span className="text-sm font-bold font-sans text-gray-700">{item.dateRange}</span>
                            </div>
                            {item.subtitle && (
                                <div className="mb-2">
                                    <span className="italic text-gray-800">{item.subtitle}</span>
                                </div>
                            )}
                             <div className="text-sm leading-normal text-gray-800 whitespace-pre-wrap" dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description || '') }} />
                        </div>
                    ))}
                </div>
            </div>
        ))}
    </div>
);
};
