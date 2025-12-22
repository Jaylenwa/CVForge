import React from 'react';
import { ResumeData, ResumeSectionType } from '../../../types';
import { useSectionTitle } from '../../../hooks/useSectionTitle';
import { useLanguage } from '../../../contexts/LanguageContext';
import { sanitizeHtml } from '../../../utils/resume-helpers';

export const TemplateCNPaper: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const { t } = useLanguage();
  const getSectionTitle = useSectionTitle();
  const accent = data.themeConfig?.color || '#8b5e3c';
  return (
    <div className={`w-full bg-[#f7f3e8] text-[#2c2c2c] h-full min-h-[1123px] ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none print:min-h-0 print:h-auto`} style={{ fontFamily: styles.fontFamily }}>
      <div className="p-12">
        <div className="flex justify-between items-center mb-10 border-b pb-6" style={{ borderColor: accent }}>
          <div>
            <h1 className="text-4xl font-semibold tracking-wide" style={{ color: accent }}>{data.personalInfo.fullName}</h1>
            <div className="text-base mt-2 text-[#5a5a5a]">{data.personalInfo.jobTitle}</div>
            <div className="mt-3 text-sm text-[#5a5a5a] space-x-3">
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

        <div className="space-y-10">
          {(data.sections || []).filter(s => s.isVisible).map(section => (
            <div key={section.id} className="pt-2">
              <h3 className="text-lg font-bold mb-4 inline-block border-b-2 pb-1" style={{ borderColor: accent, color: accent }}>{getSectionTitle(section)}</h3>
              {section.type === ResumeSectionType.Skills ? (
                <div className="flex flex-wrap gap-2">
                  {section.items.map(item => (
                    <span key={item.id} className="px-3 py-1 bg-[#efe7d8] rounded text-sm">{item.description}</span>
                  ))}
                </div>
              ) : (
                <div className="space-y-6">
                  {section.items.map(item => (
                    <div key={item.id}>
                      <div className="flex justify-between items-baseline mb-1">
                        <h4 className="text-lg font-medium">{item.title}</h4>
                        {item.dateRange && <span className="text-xs text-[#666]">{item.dateRange}</span>}
                      </div>
                      {item.subtitle && <div className="text-sm text-[#555] mb-2">{item.subtitle}</div>}
                      {item.description && (
                        <div className="text-sm leading-relaxed text-[#444]" dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description) }} />
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

