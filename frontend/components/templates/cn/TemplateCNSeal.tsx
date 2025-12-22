import React from 'react';
import { ResumeData, ResumeSectionType } from '../../../types';
import { useSectionTitle } from '../../../hooks/useSectionTitle';
import { useLanguage } from '../../../contexts/LanguageContext';
import { hasExtraPersonalInfo, sanitizeHtml } from '../../../utils/resume-helpers';

export const TemplateCNSeal: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const { t } = useLanguage();
  const getSectionTitle = useSectionTitle();
  const accent = data.themeConfig?.color || '#b91c1c';
  const initial = (data.personalInfo.fullName || '').slice(0, 1);
  return (
    <div className={`w-full bg-white text-gray-900 h-full min-h-[1123px] ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none print:min-h-0 print:h-auto`} style={{ fontFamily: styles.fontFamily }}>
      <div className="p-10">
        <div className="flex justify-between items-center mb-8">
          <div>
            <h1 className="text-4xl font-bold tracking-wider" style={{ color: accent }}>{data.personalInfo.fullName}</h1>
            <div className="text-lg text-gray-600 mt-2">{data.personalInfo.jobTitle}</div>
            <div className="mt-3 flex flex-wrap gap-x-4 gap-y-1 text-sm text-gray-700">
              {data.personalInfo.email && <span>{data.personalInfo.email}</span>}
              {data.personalInfo.phone && <span>• {data.personalInfo.phone}</span>}
              {data.personalInfo.city && <span>• {data.personalInfo.city}</span>}
            </div>
            {hasExtraPersonalInfo(data) && (
              <div className="mt-1 flex flex-wrap gap-x-3 gap-y-1 text-xs text-gray-600">
                {data.personalInfo.gender && <span>{t('editor.fields.gender')}: {data.personalInfo.gender}</span>}
                {data.personalInfo.age && <span>{t('editor.fields.age')}: {data.personalInfo.age}</span>}
                {(data.personalInfo.customInfo || []).slice(0, 2).map((ci, idx) => (
                  <span key={idx}>{ci.label}: {ci.value}</span>
                ))}
              </div>
            )}
          </div>
          <div className="flex items-center gap-5">
            {data.personalInfo.avatarUrl && (
              <img
                src={data.personalInfo.avatarUrl}
                alt={t('a11y.avatarAlt')}
                className="w-24 h-24 rounded object-cover border-2"
                style={{ borderColor: accent }}
              />
            )}
            <div className="w-16 h-16 rounded-full border-4 flex items-center justify-center text-lg font-bold" style={{ borderColor: accent, color: accent }}>
              {initial || '简'}
            </div>
          </div>
        </div>

        <div className="space-y-8">
          {(data.sections || []).filter(s => s.isVisible).map(section => (
            <div key={section.id}>
              <h3 className="text-lg font-bold tracking-wider mb-4 pl-3 border-l-4" style={{ borderColor: accent, color: accent }}>{getSectionTitle(section)}</h3>
              {section.type === ResumeSectionType.Skills ? (
                <div className="flex flex-wrap gap-2">
                  {section.items.map(item => (
                    <span key={item.id} className="px-3 py-1 bg-red-50 text-red-800 rounded border" style={{ borderColor: accent }}>{item.description}</span>
                  ))}
                </div>
              ) : (
                <div className="space-y-6">
                  {section.items.map(item => (
                    <div key={item.id}>
                      <div className="flex justify-between items-baseline mb-1">
                        <h4 className="font-bold text-gray-800">{item.title}</h4>
                        {item.dateRange && <span className="text-sm text-gray-500 font-sans">{item.dateRange}</span>}
                      </div>
                      {item.subtitle && <div className="text-sm text-gray-600 mb-2">{item.subtitle}</div>}
                      {item.description && (
                        <div className="text-sm leading-relaxed text-gray-700 whitespace-pre-wrap" dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description) }} />
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
  );
};

