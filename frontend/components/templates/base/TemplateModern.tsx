import React from 'react';
import { ResumeData, ResumeSectionType } from '../../../types';
import { useLanguage } from '../../../contexts/LanguageContext';
import { useSectionTitle } from '../../../hooks/useSectionTitle';
import { sanitizeHtml, hasExtraPersonalInfo } from '../../../utils/resume-helpers';

export const TemplateModern: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
    const getSectionTitle = useSectionTitle();
    const { t } = useLanguage();
    return (
    <div className={`w-full grid grid-cols-12 h-full min-h-[1123px] bg-white ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none print:min-h-0 print:h-auto print:bg-transparent`} style={{ fontFamily: styles.fontFamily, lineHeight: parseFloat(styles.spacingMultiplier) * 1.5 }}>
        <div className="hidden print:block fixed left-0 top-0 bottom-0 w-[70mm] -z-10" style={{ backgroundColor: '#0f172a' }}></div>
        <div className="col-span-4 text-white p-8" style={{ backgroundColor: '#0f172a' }}>
            <div className="mb-8 flex flex-col items-center md:items-start">
                 {data.personalInfo.avatarUrl && (
                     <img src={data.personalInfo.avatarUrl} alt={t('a11y.avatarAlt')} className="w-32 h-32 rounded-full mb-6 object-cover border-4 border-slate-700 shadow-lg" style={{ borderRadius: '50%' }}/>
                 )}
                 <h1 className="text-2xl font-bold leading-tight break-words" style={{ color: data.themeConfig?.color || 'white' }}>{data.personalInfo.fullName}</h1>
                 <p className="text-slate-300 mt-1 font-medium">{data.personalInfo.jobTitle}</p>
            </div>

                 <div className="space-y-4 text-sm text-slate-300 mb-8 break-all">
                     <div className="block">{data.personalInfo.email}</div>
                     <div className="block">{data.personalInfo.phone}</div>
                     {/* website removed */}
                 {data.personalInfo.gender && <div className="block">{t('editor.fields.gender')}: {data.personalInfo.gender}</div>}
                 {data.personalInfo.age && <div className="block">{t('editor.fields.age')}: {data.personalInfo.age}</div>}
                 {data.personalInfo.maritalStatus && <div className="block">{t('editor.fields.maritalStatus')}: {data.personalInfo.maritalStatus}</div>}
                 {data.personalInfo.politicalStatus && <div className="block">{t('editor.fields.politicalStatus')}: {data.personalInfo.politicalStatus}</div>}
                 {data.personalInfo.birthplace && <div className="block">{t('editor.fields.birthplace')}: {data.personalInfo.birthplace}</div>}
                 {data.personalInfo.ethnicity && <div className="block">{t('editor.fields.ethnicity')}: {data.personalInfo.ethnicity}</div>}
                 {data.personalInfo.height && <div className="block">{t('editor.fields.height')}: {data.personalInfo.height}</div>}
                 {data.personalInfo.weight && <div className="block">{t('editor.fields.weight')}: {data.personalInfo.weight}</div>}
                 {(data.personalInfo.customInfo || []).map((ci, idx) => (
                   <div key={idx} className="block">{ci.label}: {ci.value}</div>
                 ))}
            </div>

            {data.sections.filter(s => s.type === ResumeSectionType.Skills && s.isVisible).map(section => (
                <div key={section.id} className="mb-6">
                    <h3 className="text-white font-bold uppercase tracking-wider mb-4 text-sm border-b border-slate-700 pb-2">{getSectionTitle(section)}</h3>
                    <div className="space-y-1 text-slate-200">
                        {section.items.map(item => (
                            <div key={item.id} className="text-sm">{item.description}</div>
                        ))}
                    </div>
                </div>
            ))}
        </div>
        
        <div className="col-span-8 p-8">
             {data.sections.filter(s => s.type !== ResumeSectionType.Skills && s.isVisible).map(section => (
                 <div key={section.id} className="mb-8">
                      <h3 className="text-slate-900 font-bold uppercase tracking-wider mb-4 text-sm border-b-2 inline-block pb-1" style={{ borderColor: data.themeConfig?.color || '#3b82f6' }}>{getSectionTitle(section)}</h3>
                      <div className="space-y-5">
                          {section.items.map(item => (
                              <div key={item.id}>
                                  <div className="flex justify-between items-start">
                                      <div>
                                          {item.title && <h4 className="font-bold text-gray-800">{item.title}</h4>}
                                          {item.subtitle && <p className="font-medium text-sm" style={{ color: data.themeConfig?.color || '#2563eb' }}>{item.subtitle}</p>}
                                      </div>
                                      <div className="text-right shrink-0 ml-2">
                                           {item.dateRange && <p className="text-xs text-gray-500 font-medium bg-gray-100 px-2 py-1 rounded inline-block">{item.dateRange}</p>}
                                      </div>
                                  </div>
                                  <div className="mt-2 text-sm text-gray-600 leading-relaxed whitespace-pre-wrap" dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description || '') }} />
                              </div>
                          ))}
                      </div>
                 </div>
             ))}
        </div>
    </div>
)
};
