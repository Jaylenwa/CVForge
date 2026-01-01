import React from 'react';
import { ResumeData, ResumeSectionType } from '../../types';
import { useSectionTitle } from '../../hooks/useSectionTitle';
import { useLanguage } from '../../contexts/LanguageContext';
import { hasExtraPersonalInfo, sanitizeHtml } from '../../utils/resume-helpers';

export const TemplateClassic: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const getSectionTitle = useSectionTitle();
  const { t } = useLanguage();
  return (
  <div className={`w-full bg-white text-gray-900 h-auto ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none p-10`} style={{ fontFamily: styles.fontFamily, lineHeight: parseFloat(styles.spacingMultiplier) * 1.5 }}>
  <div className="border-b-2 pb-6 mb-6 flex flex-col md:flex-row items-center md:items-start gap-6" style={{ borderColor: data.themeConfig?.color || '#333' }}>
      <div className="flex-1 text-center md:text-left order-2 md:order-1">
          <h1 className="text-4xl font-bold uppercase tracking-wider" style={{ color: data.themeConfig?.color }}>{data.personalInfo.fullName}</h1>
          <p className="text-xl mt-2 text-gray-600">{data.personalInfo.jobTitle}</p>
          <div className="mt-4 flex flex-wrap justify-center md:justify-start gap-4 text-sm text-gray-600">
            {data.personalInfo.email && <span>{data.personalInfo.email}</span>}
            {data.personalInfo.phone && <span>• {data.personalInfo.phone}</span>}
          </div>
          {hasExtraPersonalInfo(data) && (
            <div className="mt-2 flex flex-wrap gap-x-4 gap-y-1 text-xs text-gray-600">
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
          <div className="order-1 md:order-2 flex-shrink-0">
          <img 
              src={data.personalInfo.avatarUrl} 
              alt={t('a11y.avatarAlt')} 
              className="w-32 h-32 rounded-lg object-cover border-2 shadow-sm"
              style={{ borderColor: data.themeConfig?.color || '#e5e7eb' }}
          />
          </div>
      )}
  </div>

  <div className="space-y-6">
      <div className="grid grid-cols-1 md:grid-cols-12 gap-8">
          <div className="md:col-span-12">
              {data.sections.filter(s => s.isVisible).map(section => (
                  <div key={section.id} className="mb-8">
                      <h3 className="text-lg font-bold uppercase border-b mb-4 pb-2" style={{ borderColor: data.themeConfig?.color || '#e5e7eb', color: data.themeConfig?.color }}>{getSectionTitle(section)}</h3>
                      
                      {section.type === ResumeSectionType.Skills ? (
                          <div className="flex flex-wrap gap-2">
                              {section.items.map(item => (
                                  <span key={item.id} className="px-3 py-1 bg-gray-100 rounded text-sm font-medium text-gray-700">{item.description}</span>
                              ))}
                          </div>
                      ) : (
                          <div className="space-y-6">
                              {section.items.map(item => (
                                  <div key={item.id}>
                                      <div className="flex flex-col md:flex-row md:justify-between md:items-baseline mb-1">
                                          <h4 className="font-bold text-lg text-gray-800">{item.title}</h4>
                                          {(item.timeStart || item.timeEnd || item.today) && (
                                            <span className="text-sm text-gray-500 font-medium">
                                              {item.timeStart || item.timeEnd}
                                              {' ~ '}
                                              {item.today ? t('common.toPresent') : (item.timeEnd || '')}
                                            </span>
                                          )}
                                      </div>
                                      {item.subtitle && <div className="text-gray-700 font-medium mb-2">{item.subtitle}</div>}
                                      {item.description && (
                                          <div className="text-gray-600 text-sm leading-relaxed whitespace-pre-wrap" dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description) }} />
                                      )}
                                  </div>
                              ))}
                          </div>
                      )}
                  </div>
              ))}
          </div>
      </div>
  </div>
  </div>
  );
};

export default TemplateClassic;
