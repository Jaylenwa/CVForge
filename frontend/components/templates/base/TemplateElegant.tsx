import React from 'react';
import { ResumeData, ResumeSectionType } from '../../../types';
import { useLanguage } from '../../../contexts/LanguageContext';
import { useSectionTitle } from '../../../hooks/useSectionTitle';
import { sanitizeHtml } from '../../../utils/resume-helpers';

export const TemplateElegant: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
    const getSectionTitle = useSectionTitle();
    const { t } = useLanguage();
    return (
    <div className={`w-full grid grid-cols-12 h-full min-h-[1123px] bg-white ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none print:min-h-0 print:h-auto print:bg-transparent`} style={{ fontFamily: styles.fontFamily, lineHeight: parseFloat(styles.spacingMultiplier) * 1.6 }}>
        <div className="hidden print:block fixed left-0 top-0 bottom-0 w-[70mm] -z-10" style={{ backgroundColor: data.themeConfig?.color || '#115e59' }}></div>
        <div className="col-span-4 text-white p-8 flex flex-col" style={{ backgroundColor: data.themeConfig?.color || '#115e59' }}>
            <div className="mb-10 text-center">
                 {data.personalInfo.avatarUrl && (
                     <img src={data.personalInfo.avatarUrl} alt="Profile" className="w-32 h-32 rounded-full mb-6 object-cover border-4 border-white/20 mx-auto shadow-md" style={{ borderRadius: '50%' }}/>
                 )}
                <h1 className="text-2xl font-serif font-bold leading-tight mb-2" style={{ fontFamily: styles.fontFamily }}>{data.personalInfo.fullName}</h1>
                 <p className="text-white/70 uppercase tracking-widest text-xs font-semibold">{data.personalInfo.jobTitle}</p>
            </div>

            <div className="space-y-6 text-sm text-white/90 flex-grow">
                 <div>
                    <span className="block text-white/60 text-xs font-bold uppercase mb-1">Contact</span>
                    <div className="space-y-1">
                        <div className="break-all">{data.personalInfo.email}</div>
                        <div>{data.personalInfo.phone}</div>
                        {/* website removed */}
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
                    </div>
                 </div>

                 {data.sections.filter(s => s.type === ResumeSectionType.Skills && s.isVisible).map(section => (
                    <div key={section.id}>
                        <span className="block text-white/60 text-xs font-bold uppercase mb-2">{getSectionTitle(section)}</span>
                        <div className="space-y-1 text-white/80">
                            {section.items.map(item => (
                                <div key={item.id}>{item.description}</div>
                            ))}
                        </div>
                    </div>
                 ))}
            </div>
        </div>
        
        <div className="col-span-8 p-10 bg-white">
             {data.sections.filter(s => s.type !== ResumeSectionType.Skills && s.isVisible).map(section => (
                 <div key={section.id} className="mb-10">
                      <h3 className="font-serif font-bold uppercase tracking-widest text-lg border-b pb-2 mb-6" style={{ borderColor: '#f3f4f6', color: data.themeConfig?.color || '#134e4a' }}>{getSectionTitle(section)}</h3>
                      <div className="space-y-6">
                          {section.items.map(item => (
                              <div key={item.id}>
                                  <div className="flex justify-between items-baseline mb-1">
                                      <h4 className="font-bold text-gray-900 text-lg">{item.title}</h4>
                                      <span className="text-sm font-medium" style={{ color: data.themeConfig?.color || '#115e59' }}>{item.dateRange}</span>
                                  </div>
                                  {item.subtitle && <p className="text-gray-600 italic mb-2">{item.subtitle}</p>}
                                  <div className="text-sm text-gray-600 leading-relaxed whitespace-pre-wrap font-light" dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description || '') }} />
                              </div>
                          ))}
                      </div>
                 </div>
             ))}
        </div>
    </div>
);
};
