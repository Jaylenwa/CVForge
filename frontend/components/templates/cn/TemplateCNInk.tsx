import React from 'react';
import { ResumeData, ResumeSectionType } from '../../../types';
import { useSectionTitle } from '../../../hooks/useSectionTitle';
import { useLanguage } from '../../../contexts/LanguageContext';
import { sanitizeHtml } from '../../../utils/resume-helpers';

export const TemplateCNInk: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const { t } = useLanguage();
  const getSectionTitle = useSectionTitle();
  const accent = data.themeConfig?.color || '#111827';
  return (
    <div className={`w-full bg-white text-gray-900 h-full min-h-[1123px] ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none print:min-h-0 print:h-auto`} style={{ fontFamily: styles.fontFamily }}>
      <div className="p-10">
        <div className="mb-6">
          <div className="h-1 bg-gradient-to-r from-black to-transparent w-32"></div>
          <div className="flex justify-between items-center mt-4">
            <div>
              <h1 className="text-4xl font-extrabold tracking-tight" style={{ color: accent }}>{data.personalInfo.fullName}</h1>
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
                className="w-24 h-24 rounded-md object-cover border"
                style={{ borderColor: accent }}
              />
            )}
          </div>
        </div>

        <div className="space-y-8">
          {(data.sections || []).filter(s => s.isVisible).map(section => (
            <div key={section.id}>
              <div className="flex items-center mb-4">
                <div className="w-24 h-px bg-black mr-3"></div>
                <h3 className="text-lg font-bold tracking-wider" style={{ color: accent }}>{getSectionTitle(section)}</h3>
              </div>
              {section.type === ResumeSectionType.Skills ? (
                <div className="flex flex-wrap gap-2">
                  {section.items.map(item => (
                    <span key={item.id} className="px-3 py-1 bg-gray-100 rounded text-sm">{item.description}</span>
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

