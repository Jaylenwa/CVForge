import React from 'react';
import { ResumeData, ResumeSectionType } from '../../types';
import { useSectionTitle } from '../../hooks/useSectionTitle';
import { useLanguage } from '../../contexts/LanguageContext';
import { hasExtraPersonalInfo, sanitizeHtml } from '../../utils/resume-helpers';

export const TemplateClassic: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const getSectionTitle = useSectionTitle();
  const { t } = useLanguage();
  return (
  <div className={`w-full bg-white text-gray-900 h-auto ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none p-10`} style={{ fontFamily: styles.fontFamily, lineHeight: parseFloat(styles.spacingMultiplier) * 1.5, fontSize: styles.fontSize }}>
  <div className="border-b-2 pb-6 mb-6 flex flex-col md:flex-row items-center md:items-start gap-6" style={{ borderColor: data.Theme?.Color || '#333' }}>
      <div className="flex-1 text-center md:text-left order-2 md:order-1">
          <h1 className="text-4xl font-bold uppercase tracking-wider" style={{ color: data.Theme?.Color }}>{data.Personal?.FullName}</h1>
          <p className="text-xl mt-2 text-gray-600">{data.Personal?.Job}</p>
          <div className="mt-4 flex flex-wrap justify-center md:justify-start gap-4 text-sm text-gray-600">
            {data.Personal?.Email && <span>{data.Personal.Email}</span>}
            {data.Personal?.Phone && <span>• {data.Personal.Phone}</span>}
          </div>
          {hasExtraPersonalInfo(data) && (
            <div className="mt-2 flex flex-wrap gap-x-4 gap-y-1 text-xs text-gray-600">
              {data.Personal?.Gender && <span>{t('editor.fields.gender')}: {data.Personal.Gender}</span>}
              {data.Personal?.Age && <span>{t('editor.fields.age')}: {data.Personal.Age}</span>}
              {data.Personal?.Degree && <span>{t('editor.fields.degree')}: {data.Personal.Degree}</span>}
              
              
              {(() => {
                try {
                  const raw = data.Personal?.CustomInfo;
                  if (raw) {
                    const arr = JSON.parse(raw);
                    if (Array.isArray(arr)) {
                      return arr.map((ci: any, idx: number) => <span key={idx}>{ci.label}: {ci.value}</span>);
                    }
                  }
                } catch {}
                return null;
              })()}
            </div>
          )}
      </div>
      {data.Personal?.AvatarURL && (
          <div className="order-1 md:order-2 flex-shrink-0">
          <img 
              src={data.Personal.AvatarURL} 
              alt={t('a11y.avatarAlt')} 
              className="w-32 h-32 rounded-lg object-cover border-2 shadow-sm"
              style={{ borderColor: data.Theme?.Color || '#e5e7eb' }}
          />
          </div>
      )}
  </div>

  <div className="space-y-6">
      <div className="grid grid-cols-1 md:grid-cols-12 gap-8">
          <div className="md:col-span-12">
              {data.sections.filter(s => s.isVisible).map(section => (
                  <div key={section.id} className="mb-8">
                      <h3 className="text-lg font-bold uppercase border-b mb-4 pb-2" style={{ borderColor: data.Theme?.Color || '#e5e7eb', color: data.Theme?.Color }}>{getSectionTitle(section)}</h3>
                      
                      {section.type === ResumeSectionType.Skills ? (
                          <div className="space-y-4">
                            {section.items.map(item => (
                              <div key={item.id}>
                                {item.description && (
                                  <div className="resume-rich-content text-gray-600 text-sm leading-relaxed" style={{ fontSize: styles.fontSize }} dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description) }} />
                                )}
                              </div>
                            ))}
                          </div>
                      ) : (
                          <div className="space-y-6">
                              {section.items.map(item => (
                                  <div key={item.id}>
                                      {section.type !== ResumeSectionType.Awards && (
                                        <div className="flex flex-col md:flex-row md:justify-between md:items-baseline mb-1">
                                          <h4 className="font-bold text-lg text-gray-800">{item.title}</h4>
                                          {(item.timeStart || item.timeEnd || item.today) && (
                                            <span className="text-sm text-gray-500 font-medium" style={{ fontSize: styles.fontSize }}>
                                              {item.timeStart || item.timeEnd}
                                              {' ~ '}
                                              {item.today ? t('common.toPresent') : (item.timeEnd || '')}
                                            </span>
                                          )}
                                        </div>
                                      )}
                                      {section.type !== ResumeSectionType.Awards && item.subtitle && <div className="text-gray-700 font-medium mb-2">{item.subtitle}</div>}
                                      {section.type === ResumeSectionType.Education && (item.major || item.degree) && (
                                        <div className="text-gray-700 text-sm mb-2" style={{ fontSize: styles.fontSize }}>
                                          {item.major}{item.major && item.degree ? ' • ' : ''}{item.degree}
                                        </div>
                                      )}
                                      {item.description && (
                                          <div className="resume-rich-content text-gray-600 text-sm leading-relaxed" style={{ fontSize: styles.fontSize }} dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description) }} />
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
