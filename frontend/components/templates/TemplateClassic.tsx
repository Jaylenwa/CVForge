import React from 'react';
import { ResumeData, ResumeSectionType } from '../../types';
import { useSectionTitle } from '../../hooks/useSectionTitle';
import { useLanguage } from '../../contexts/LanguageContext';
import { ExamSection } from './shared/ExamSection';
import { RichText } from './shared/RichText';
import { formatDateRange, getAccentColor, getLineHeight, getOrderedItems, getOrderedVisibleSections, normalizeCustomPairs, parseCustomPairs } from './shared/templateTokens';

export const TemplateClassic: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const getSectionTitle = useSectionTitle();
  const { t } = useLanguage();
  const personal = (data.Personal || {}) as NonNullable<ResumeData['Personal']>;
  const customInfo = normalizeCustomPairs(parseCustomPairs(personal?.CustomInfo));
  const extraPairs: Array<{ label: string; value: string }> = [
    { label: t('editor.fields.gender'), value: personal?.Gender || '' },
    { label: t('editor.fields.age'), value: personal?.Age || '' },
    { label: t('editor.fields.degree'), value: personal?.Degree || '' },
    { label: t('editor.fields.city'), value: personal?.City || '' },
    { label: t('editor.fields.expectedSalary'), value: personal?.Money || '' },
    { label: t('editor.fields.joinTime'), value: personal?.JoinTime || '' },
  ].filter(p => p.value && String(p.value).trim());
  const hasCustomInfo = customInfo.length > 0;
  const lineHeight = getLineHeight(styles);
  const accent = getAccentColor(data, '#14b8a6');

  return (
  <div className={`w-full bg-white text-gray-900 h-auto ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none p-10`} style={{ fontFamily: styles.fontFamily, lineHeight, fontSize: styles.fontSize }}>
  <div className="pb-6 mb-6 flex flex-col md:flex-row items-center md:items-start gap-6">
      <div className="flex-1 text-center md:text-left order-2 md:order-1">
          <h1 className="text-4xl font-bold uppercase tracking-wider">{personal?.FullName}</h1>
          {personal?.Job && <p className="text-xl mt-2 text-gray-600">{personal.Job}</p>}
          {(personal?.Email || personal?.Phone) && (
            <div className="mt-4 flex flex-wrap justify-center md:justify-start gap-4 text-sm text-gray-600">
              {personal?.Email && <span>{personal.Email}</span>}
              {personal?.Phone && <span>{personal?.Email ? '• ' : ''}{personal.Phone}</span>}
            </div>
          )}
          {(extraPairs.length > 0 || hasCustomInfo) && (
            <div className="mt-2 flex flex-wrap gap-x-4 gap-y-1 text-xs text-gray-600">
              {extraPairs.map((p, idx) => (
                <span key={`${p.label}-${idx}`}>{p.label}: {p.value}</span>
              ))}
              {customInfo.map((ci, idx) => (
                <span key={`ci-${idx}`}>{ci.label ? `${ci.label}: ${ci.value}` : ci.value}</span>
              ))}
            </div>
          )}
      </div>
      {personal?.AvatarURL && (
          <div className="order-1 md:order-2 flex-shrink-0">
          <img 
              src={personal.AvatarURL} 
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
              {getOrderedVisibleSections([...data.sections]).map(section => (
                  <div key={section.id} className="mb-8">
                      <h3 className="text-lg font-bold uppercase border-b mb-4 pb-2" style={{ borderColor: data.Theme?.Color || '#e5e7eb', color: data.Theme?.Color }}>{getSectionTitle(section)}</h3>
                      
                      {section.type === ResumeSectionType.Exam ? (
                        <ExamSection section={section} color={accent} t={t} />
                      ) : section.type === ResumeSectionType.Skills ? (
                          <div className="space-y-4">
                            {getOrderedItems(section.items).map(item => (
                              <div key={item.id}>
                                {item.description && (
                                  <RichText html={item.description} className="text-gray-600" fontSize={styles.fontSize} />
                                )}
                              </div>
                            ))}
                          </div>
                      ) : (
                          <div className="space-y-6">
                              {getOrderedItems(section.items).map(item => (
                                  <div key={item.id}>
                                      {section.type !== ResumeSectionType.Awards && (
                                        <div className="flex flex-col md:flex-row md:justify-between md:items-baseline mb-1">
                                          <h4 className="font-bold text-lg text-gray-800">{item.title}</h4>
                                          {(item.timeStart || item.timeEnd || item.today) && (
                                            <span className="text-sm text-gray-500 font-medium" style={{ fontSize: styles.fontSize }}>
                                              {formatDateRange(item, t, { separatorVariant: 'tilde' })}
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
                                          <RichText html={item.description} className="text-gray-600" fontSize={styles.fontSize} />
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
