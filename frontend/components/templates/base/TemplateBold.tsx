import React from 'react';
import { ResumeData, ResumeSectionType } from '../../../types';
import { useLanguage } from '../../../contexts/LanguageContext';
import { useSectionTitle } from '../../../hooks/useSectionTitle';
import { sanitizeHtml, hasExtraPersonalInfo } from '../../../utils/resume-helpers';

export const TemplateBold: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
     const { t } = useLanguage();
     const getSectionTitle = useSectionTitle();
     return (
     <div className={`w-full bg-white h-full min-h-[1123px] ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none print:min-h-0 print:h-auto flex flex-col`} style={{ fontFamily: styles.fontFamily, lineHeight: parseFloat(styles.spacingMultiplier) * 1.5 }}>
        <div className="text-white p-10 print:text-white flex flex-col md:flex-row justify-between items-center gap-6" style={{ backgroundColor: data.themeConfig?.color || '#1d4ed8' }}>
             <div className="order-2 md:order-1 flex-1">
                <h1 className="text-5xl font-extrabold mb-2 tracking-tight">{data.personalInfo.fullName}</h1>
                <p className="text-white/80 text-2xl font-light">{data.personalInfo.jobTitle}</p>
                <div className="mt-6 flex flex-wrap gap-x-6 gap-y-2 text-sm text-white/90 font-medium">
                    {data.personalInfo.email && <div className="flex items-center gap-2">{data.personalInfo.email}</div>}
                    {data.personalInfo.phone && <div className="flex items-center gap-2">• {data.personalInfo.phone}</div>}
                </div>
                <div className="mt-2 flex flex-wrap gap-x-6 gap-y-1 text-xs text-white/80">
                    {hasExtraPersonalInfo(data) && (
                      <>
                        {data.personalInfo.gender && <div>{t('editor.fields.gender')}: {data.personalInfo.gender}</div>}
                        {data.personalInfo.age && <div>{t('editor.fields.age')}: {data.personalInfo.age}</div>}
                        {data.personalInfo.maritalStatus && <div>{t('editor.fields.maritalStatus')}: {data.personalInfo.maritalStatus}</div>}
                        {data.personalInfo.politicalStatus && <div>{t('editor.fields.politicalStatus')}: {data.personalInfo.politicalStatus}</div>}
                        {data.personalInfo.birthplace && <div>{t('editor.fields.birthplace')}: {data.personalInfo.birthplace}</div>}
                        {data.personalInfo.ethnicity && <div>{t('editor.fields.ethnicity')}: {data.personalInfo.ethnicity}</div>}
                        {data.personalInfo.height && <div>{t('editor.fields.height')}: {data.personalInfo.height}</div>}
                        {data.personalInfo.weight && <div>{t('editor.fields.weight')}: {data.personalInfo.weight}</div>}
                        {(data.personalInfo.customInfo || []).map((ci, idx) => (
                          <div key={idx}>{ci.label}: {ci.value}</div>
                        ))}
                      </>
                    )}
                </div>
             </div>
             {data.personalInfo.avatarUrl && (
                 <img src={data.personalInfo.avatarUrl} alt={t('a11y.avatarAlt')} className="order-1 md:order-2 w-32 h-32 rounded-full border-4 border-white object-cover shadow-xl flex-shrink-0" style={{ borderRadius: '50%' }} />
             )}
        </div>

        <div className="p-10 grid grid-cols-1 md:grid-cols-12 gap-8 flex-grow">
            <div className="md:col-span-8 pr-4">
                 {data.sections.filter(s => s.type !== ResumeSectionType.Skills && s.isVisible).map(section => (
                     <div key={section.id} className="mb-10">
                         <h3 className="font-bold text-lg uppercase mb-6 flex items-center tracking-wider" style={{ color: data.themeConfig?.color || '#1d4ed8' }}>
                             <span className="w-1.5 h-6 mr-3 rounded-sm" style={{ backgroundColor: data.themeConfig?.color || '#1d4ed8' }}></span>
                             {getSectionTitle(section)}
                         </h3>
                          <div className="space-y-8">
                            {section.items.map(item => (
                                <div key={item.id} className="relative">
                                    <div className="flex justify-between items-center mb-1">
                                        <h4 className="font-bold text-gray-900 text-lg">{item.title}</h4>
                                        {item.dateRange && <span className="text-xs font-bold bg-gray-100 text-gray-700 px-3 py-1 rounded-full whitespace-nowrap ml-4">{item.dateRange}</span>}
                                    </div>
                                    {item.subtitle && <p className="text-sm text-gray-600 font-medium mb-3">{item.subtitle}</p>}
                                    <div className="text-sm text-gray-600 whitespace-pre-wrap leading-relaxed" dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description || '') }} />
                                </div>
                            ))}
                        </div>
                     </div>
                 ))}
            </div>
            
            <div className="md:col-span-4 bg-gray-50 p-6 rounded-2xl h-fit border border-gray-100">
                 {data.sections.filter(s => s.type === ResumeSectionType.Skills && s.isVisible).map(section => (
                     <div key={section.id} className="mb-8">
                         <h3 className="text-gray-900 font-bold uppercase mb-4 text-sm tracking-wider">{getSectionTitle(section)}</h3>
                         <div className="space-y-1 text-gray-700">
                             {section.items.map(item => (
                                 <div key={item.id} className="text-sm">{item.description}</div>
                             ))}
                         </div>
                     </div>
                 ))}
            </div>
        </div>
     </div>
     );
};
