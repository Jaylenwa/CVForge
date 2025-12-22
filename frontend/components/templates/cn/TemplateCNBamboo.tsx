import React from 'react';
import { ResumeData, ResumeSectionType } from '../../../types';
import { useSectionTitle } from '../../../hooks/useSectionTitle';
import { useLanguage } from '../../../contexts/LanguageContext';
import { sanitizeHtml } from '../../../utils/resume-helpers';

export const TemplateCNBamboo: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const { t } = useLanguage();
  const getSectionTitle = useSectionTitle();
  const accent = data.themeConfig?.color || '#166534';
  return (
    <div className={`w-full bg-white text-gray-900 h-full min-h-[1123px] ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none print:min-h-0 print:h-auto`} style={{ fontFamily: styles.fontFamily }}>
      <div className="p-10">
        <div className="mb-8">
          <div className="h-1 w-16 rounded-full" style={{ backgroundColor: accent }}></div>
          <div className="flex justify-between items-center mt-4">
            <div>
              <h1 className="text-4xl font-bold tracking-wide" style={{ color: accent }}>{data.personalInfo.fullName}</h1>
              <div className="text-base text-gray-600 mt-2">{data.personalInfo.jobTitle}</div>
              <div className="mt-3 flex flex-wrap gap-x-4 gap-y-1 text-sm text-gray-700">
                {data.personalInfo.email && <span>{data.personalInfo.email}</span>}
                {data.personalInfo.phone && <span>• {data.personalInfo.phone}</span>}
                {data.personalInfo.city && <span>• {data.personalInfo.city}</span>}
              </div>
            </div>
            {data.personalInfo.avatarUrl && (
              <img
                src={data.personalInfo.avatarUrl}
                alt={t('a11y.avatarAlt')}
                className="w-24 h-24 rounded-full object-cover border-2"
                style={{ borderColor: accent }}
              />
            )}
          </div>
        </div>

        <div className="space-y-8">
          {(data.sections || []).filter(s => s.isVisible).map(section => (
            <div key={section.id}>
              <h3 className="text-lg font-semibold mb-4 flex items-center" style={{ color: accent }}>
                <span className="inline-block w-2 h-2 rounded-full mr-2" style={{ backgroundColor: accent }}></span>
                {getSectionTitle(section)}
              </h3>
              {section.type === ResumeSectionType.Skills ? (
                <div className="flex flex-wrap gap-2">
                  {section.items.map(item => (
                    <span key={item.id} className="px-3 py-1 bg-green-50 text-green-900 rounded border" style={{ borderColor: accent }}>{item.description}</span>
                  ))}
                </div>
              ) : (
                <div className="space-y-6">
                  {section.items.map(item => (
                    <div key={item.id}>
                      <div className="flex justify-between items-baseline mb-1">
                        <h4 className="font-bold text-gray-800">{item.title}</h4>
                        {item.dateRange && <span className="text-xs text-gray-500">{item.dateRange}</span>}
                      </div>
                      {item.subtitle && <div className="text-sm text-gray-600 mb-2">{item.subtitle}</div>}
                      {item.description && (
                        <div className="text-sm leading-relaxed text-gray-700" dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description) }} />
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

